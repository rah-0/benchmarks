package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/moby/moby/api/types/network"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const databaseImage = "postgres:19.0-trixie@sha256:bb6e70add742cd868fb62cb3b2f67d1dd2bdd2b65a521e778420748753687c60"

func databaseURL(host string, port network.Port) string {
	return "postgres://bench:bench@" + net.JoinHostPort(host, port.Port()) + "/bench?sslmode=disable"
}

func TestMain(m *testing.M) { os.Exit(runTests(m)) }

func runTests(m *testing.M) (code int) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        databaseImage,
			ExposedPorts: []string{"5432/tcp"},
			Env:          map[string]string{"POSTGRES_DB": "bench", "POSTGRES_USER": "bench", "POSTGRES_PASSWORD": "bench", "TZ": "UTC"},
			Cmd:          []string{"postgres", "-c", "fsync=off", "-c", "synchronous_commit=off", "-c", "shared_buffers=1GB"},
			WaitingFor:   wait.ForSQL("5432/tcp", "postgres", databaseURL).WithStartupTimeout(2 * time.Minute),
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
	port, err := container.MappedPort(ctx, "5432/tcp")
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
	stmt, err := tx.Prepare("INSERT INTO sample_a (FieldA, FirstInsert) VALUES ($1, $2)")
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
		rows, err := db.Query("SELECT FieldA, FirstInsert FROM sample_a WHERE FirstInsert < $1 ORDER BY FirstInsert", fixture.cutoff)
		if err != nil {
			t.Fatal(err)
		}
		checkRows(t, rows, fixture.matchingRows, true)
	}
}

func BenchmarkPostgresSingleInsertFixedData(b *testing.B) {
	resetTable(b)
	b.ReportAllocs()
	for b.Loop() {
		if err := dbTableSampleAInsert("9xKf3QpLm2Ry7UbHt6NwEjVg8As5OcIy4B"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPostgresSingleInsertRandomData(b *testing.B) {
	resetTable(b)
	b.ReportAllocs()
	for b.Loop() {
		if err := dbTableSampleAInsert(uuid.NewString()); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPostgresInsert1MilAndFindMiddle(b *testing.B) {
	fixture := prepareFixture(b, 1_000_000)
	b.ReportAllocs()
	for b.Loop() {
		var found string
		row := db.QueryRow("SELECT FieldA FROM sample_a WHERE FieldA = $1 LIMIT 1", fixture.middleID)
		if err := row.Scan(&found); err != nil {
			b.Fatal(err)
		}
		if found != fixture.middleID {
			b.Fatalf("expected %s, got %s", fixture.middleID, found)
		}
	}
}

func BenchmarkPostgresInsert100kAndSort(b *testing.B) {
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

func BenchmarkPostgresInsert100kAndQueryOlderThan15Min(b *testing.B) {
	fixture := prepareFixture(b, 100_000)
	b.ReportAllocs()
	for b.Loop() {
		rows, err := db.Query("SELECT FieldA, FirstInsert FROM sample_a WHERE FirstInsert < $1 ORDER BY FirstInsert", fixture.cutoff)
		if err != nil {
			b.Fatal(err)
		}
		checkRows(b, rows, fixture.matchingRows, true)
	}
}
