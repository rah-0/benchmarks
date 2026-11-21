package mariadb

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/moby/moby/api/types/network"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const databaseImage = "mariadb:13.0.2@sha256:54f493ea0d60b621a5bff84e5b20100f8b86b6efa6b79158e04756c5194a5e4d"

func databaseURL(host string, port network.Port) string {
	return "bench:bench@tcp(" + net.JoinHostPort(host, port.Port()) + ")/bench?parseTime=true&loc=UTC"
}

func TestMain(m *testing.M) { os.Exit(runTests(m)) }

func runTests(m *testing.M) (code int) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        databaseImage,
			ExposedPorts: []string{"3306/tcp"},
			Env:          map[string]string{"MARIADB_ROOT_PASSWORD": "bench-root", "MARIADB_DATABASE": "bench", "MARIADB_USER": "bench", "MARIADB_PASSWORD": "bench", "TZ": "UTC"},
			Cmd:          []string{"--innodb-flush-log-at-trx-commit=0", "--sync-binlog=0", "--innodb-buffer-pool-size=1G"},
			WaitingFor:   wait.ForSQL("3306/tcp", "mysql", databaseURL).WithStartupTimeout(2 * time.Minute),
		},
		Started: true,
	})
	if container != nil {
		defer func() {
			cleanupCtx, stop := context.WithTimeout(context.Background(), 30*time.Second)
			defer stop()
			if err := container.Terminate(cleanupCtx); err != nil {
				fmt.Fprintln(os.Stderr, err)
				code = 1
			}
		}()
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	host, err := container.Host(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	port, err := container.MappedPort(ctx, "3306/tcp")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if err := dbConnect(databaseURL(host, port)); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer func() {
		if err := dbDisconnect(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			code = 1
		}
	}()
	return m.Run()
}

type fixture struct {
	middleID     string
	cutoff       time.Time
	matchingRows int
}

func resetTable(tb testing.TB) {
	tb.Helper()
	if err := dbTableSampleADrop(); err != nil {
		tb.Fatal(err)
	}
	if err := dbTableSampleACreate(); err != nil {
		tb.Fatal(err)
	}
}

// Read benchmarks populate an independent fixture before B.Loop starts timing.
// A transaction and prepared insert keep population outside the measured work.
func prepareFixture(tb testing.TB, totalRows int) fixture {
	tb.Helper()
	resetTable(tb)
	reference := time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC)
	result := fixture{cutoff: reference.Add(-15 * time.Minute), matchingRows: totalRows * 4 / 5}
	tx, err := db.Begin()
	if err != nil {
		tb.Fatal(err)
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare("INSERT INTO sample_a (FieldA, FirstInsert) VALUES (?, ?)")
	if err != nil {
		tb.Fatal(err)
	}
	defer stmt.Close()
	for i := range totalRows {
		id := uuid.NewString()
		if i == totalRows/2 {
			result.middleID = id
		}
		minutes := i % 10
		if i < result.matchingRows {
			minutes += 16
		}
		insertedAt := reference.Add(-time.Duration(minutes) * time.Minute)
		if _, err := stmt.Exec(id, insertedAt); err != nil {
			tb.Fatal(err)
		}
	}
	if err := stmt.Close(); err != nil {
		tb.Fatal(err)
	}
	if err := tx.Commit(); err != nil {
		tb.Fatal(err)
	}
	return result
}

func checkRows(tb testing.TB, rows *sql.Rows, expected int, byTime bool) {
	tb.Helper()
	defer rows.Close()
	var previous string
	var previousTime time.Time
	count := 0
	for rows.Next() {
		var id string
		if byTime {
			var insertedAt time.Time
			if err := rows.Scan(&id, &insertedAt); err != nil {
				tb.Fatal(err)
			}
			if count > 0 && insertedAt.Before(previousTime) {
				tb.Fatal("timestamps are not sorted")
			}
			previousTime = insertedAt
		} else {
			if err := rows.Scan(&id); err != nil {
				tb.Fatal(err)
			}
			if count > 0 && id < previous {
				tb.Fatal("IDs are not sorted")
			}
			previous = id
		}
		count++
	}
	if err := rows.Err(); err != nil {
		tb.Fatal(err)
	}
	if err := rows.Close(); err != nil {
		tb.Fatal(err)
	}
	if count != expected {
		tb.Fatalf("expected %d rows, got %d", expected, count)
	}
}

func TestFixturesAreIndependent(t *testing.T) {
	for _, size := range []int{10, 5} {
		fixture := prepareFixture(t, size)
		var count int
		if err := db.QueryRow("SELECT COUNT(*) FROM sample_a").Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count != size {
			t.Fatalf("expected %d fixture rows, got %d", size, count)
		}
		rows, err := db.Query("SELECT FieldA, FirstInsert FROM sample_a WHERE FirstInsert < ? ORDER BY FirstInsert", fixture.cutoff)
		if err != nil {
			t.Fatal(err)
		}
		checkRows(t, rows, fixture.matchingRows, true)
	}
}

func BenchmarkMariaDBSingleInsertFixedData(b *testing.B) {
	resetTable(b)
	b.ReportAllocs()
	for b.Loop() {
		if err := dbTableSampleAInsert("9xKf3QpLm2Ry7UbHt6NwEjVg8As5OcIy4B"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMariaDBSingleInsertRandomData(b *testing.B) {
	resetTable(b)
	b.ReportAllocs()
	for b.Loop() {
		if err := dbTableSampleAInsert(uuid.NewString()); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMariaDBInsert1MilAndFindMiddle(b *testing.B) {
	fixture := prepareFixture(b, 1_000_000)
	b.ReportAllocs()
	for b.Loop() {
		var found string
		row := db.QueryRow("SELECT FieldA FROM sample_a WHERE FieldA = ? LIMIT 1", fixture.middleID)
		if err := row.Scan(&found); err != nil {
			b.Fatal(err)
		}
		if found != fixture.middleID {
			b.Fatalf("expected %s, got %s", fixture.middleID, found)
		}
	}
}

func BenchmarkMariaDBInsert100kAndSort(b *testing.B) {
	const totalRows = 100_000
	prepareFixture(b, totalRows)
	b.ReportAllocs()
	for b.Loop() {
		rows, err := db.Query("SELECT FieldA FROM sample_a ORDER BY FieldA ASC")
		if err != nil {
			b.Fatal(err)
		}
		checkRows(b, rows, totalRows, false)
	}
}

func BenchmarkMariaDBInsert100kAndQueryOlderThan15Min(b *testing.B) {
	fixture := prepareFixture(b, 100_000)
	b.ReportAllocs()
	for b.Loop() {
		rows, err := db.Query("SELECT FieldA, FirstInsert FROM sample_a WHERE FirstInsert < ? ORDER BY FirstInsert", fixture.cutoff)
		if err != nil {
			b.Fatal(err)
		}
		checkRows(b, rows, fixture.matchingRows, true)
	}
}
