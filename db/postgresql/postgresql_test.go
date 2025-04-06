package postgresql

import (
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/rah-0/benchmarks/util/testutil"
)

func TestMain(m *testing.M) {
	testutil.TestMainWrapper(testutil.TestConfig{
		M: m,
		LoadResources: func() error {
			if err := dbConnect(); err != nil {
				return err
			}
			if err := dbTableSampleACreate(); err != nil {
				return err
			}
			return nil
		},
		UnloadResources: func() error {
			if err := dbTableSampleADrop(); err != nil {
				return err
			}
			if err := dbDisconnect(); err != nil {
				return err
			}
			return nil
		},
	})
}

func BenchmarkPostgresSingleInsertFixedData(b *testing.B) {
	defer testutil.RecoverBenchHandler(b)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := dbTableSampleAInsert("9xKf3QpLm2Ry7UbHt6NwEjVg8As5OcIy4B")
		if err != nil {
			b.Fatalf("Insert failed on iteration %d: %v", i, err)
		}
	}
}

func BenchmarkPostgresSingleInsertRandomData(b *testing.B) {
	defer testutil.RecoverBenchHandler(b)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := dbTableSampleAInsert(uuid.NewString())
		if err != nil {
			b.Fatalf("Insert failed on iteration %d: %v", i, err)
		}
	}
}

func BenchmarkPostgresInsert1MilAndFindMiddle(b *testing.B) {
	defer testutil.RecoverBenchHandler(b)

	const totalRows = 1_000_000
	middleIndex := totalRows / 2
	uuids := make([]string, totalRows)

	for i := range uuids {
		uuids[i] = uuid.NewString()
	}

	for _, id := range uuids {
		if err := dbTableSampleAInsert(id); err != nil {
			b.Fatalf("Insert failed: %v", err)
		}
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var found string
		row := db.QueryRow(`SELECT FieldA FROM sample_a WHERE FieldA = $1 LIMIT 1`, uuids[middleIndex])
		if err := row.Scan(&found); err != nil {
			b.Fatalf("Select failed: %v", err)
		}
		if found != uuids[middleIndex] {
			b.Fatalf("Expected %s, got %s", uuids[middleIndex], found)
		}
	}
}

var inserted bool

func BenchmarkPostgresInsert100kAndSort(b *testing.B) {
	defer testutil.RecoverBenchHandler(b)

	const totalRows = 100_000
	uuids := make([]string, totalRows)

	if !inserted {
		for i := range uuids {
			uuids[i] = uuid.NewString()
			if err := dbTableSampleAInsert(uuids[i]); err != nil {
				b.Fatalf("Insert failed: %v", err)
			}
		}

		inserted = true
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rows, err := db.Query(`SELECT FieldA FROM sample_a ORDER BY FieldA ASC`)
		if err != nil {
			b.Fatalf("Select failed: %v", err)
		}

		var last string
		first := true
		count := 0

		for rows.Next() {
			var field string
			if err := rows.Scan(&field); err != nil {
				b.Fatalf("Row scan failed: %v", err)
			}
			if !first && field < last {
				b.Fatalf("Sort order incorrect at row %d: %s < %s", count, field, last)
			}
			last = field
			first = false
			count++
		}
		rows.Close()

		if count != totalRows {
			b.Fatalf("Expected %d rows, got %d", totalRows, count)
		}
		if err := rows.Err(); err != nil {
			b.Fatalf("Row iteration error: %v", err)
		}
	}
}

func BenchmarkPostgresInsert100kAndQueryOlderThan15Min(b *testing.B) {
	defer testutil.RecoverBenchHandler(b)

	const totalRows = 100_000

	if !inserted {
		now := time.Now()

		for i := 0; i < totalRows; i++ {
			// ~80% of records will be older than 15 minutes
			var insertTime time.Time
			if i < int(float64(totalRows)*0.8) {
				insertTime = now.Add(-time.Duration(16+rand.Intn(10)) * time.Minute)
			} else {
				insertTime = now.Add(-time.Duration(rand.Intn(10)) * time.Minute)
			}

			_, err := db.Exec(`INSERT INTO sample_a (FieldA, FirstInsert) VALUES ($1, $2)`, uuid.NewString(), insertTime)
			if err != nil {
				b.Fatalf("Insert failed: %v", err)
			}
		}

		inserted = true
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		rows, err := db.Query(`
			SELECT FieldA, FirstInsert
			FROM sample_a
			WHERE FirstInsert < (NOW() - INTERVAL '15 minutes')
			ORDER BY FirstInsert
		`)
		if err != nil {
			b.Fatalf("Query failed: %v", err)
		}

		var count int
		var last time.Time
		for rows.Next() {
			var fieldA string
			var ts time.Time
			if err := rows.Scan(&fieldA, &ts); err != nil {
				b.Fatalf("Scan failed: %v", err)
			}
			if count > 0 && ts.Before(last) {
				b.Fatalf("Sort order invalid at row %d: %v before %v", count, ts, last)
			}
			last = ts
			count++
		}
		rows.Close()

		if err := rows.Err(); err != nil {
			b.Fatalf("Row error: %v", err)
		}

		if count == 0 {
			b.Fatalf("Expected some rows older than 15 minutes, got 0")
		}
	}
}
