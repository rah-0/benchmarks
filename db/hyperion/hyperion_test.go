package hyperion

import (
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	SampleV1 "github.com/rah-0/hyperion/entities/Sample/v1"
	"github.com/rah-0/hyperion/hconn"
	"github.com/rah-0/hyperion/node"
	"github.com/rah-0/hyperion/query"
	"github.com/rah-0/testmark/testutil"
)

func TestMain(m *testing.M) {
	testutil.TestMainWrapper(testutil.TestConfig{
		M: m,
		LoadResources: func() (err error) {
			hconn.Timeout = 2 * time.Minute
			connection, err = node.ConnectToNodeWithHostAndPort("127.0.0.1", "5000")
			return
		},
		UnloadResources: func() error {
			return connection.Close()
		},
	})
}

func BenchmarkHyperionSingleInsertFixedData(b *testing.B) {
	defer testutil.RecoverBenchHandler(b)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := dbSampleInsert("9xKf3QpLm2Ry7UbHt6NwEjVg8As5OcIy4B")
		if err != nil {
			b.Fatalf("Insert failed on iteration %d: %v", i, err)
		}
	}
}

func BenchmarkHyperionSingleInsertRandomData(b *testing.B) {
	defer testutil.RecoverBenchHandler(b)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := dbSampleInsert(uuid.NewString())
		if err != nil {
			b.Fatalf("Insert failed on iteration %d: %v", i, err)
		}
	}
}

func BenchmarkHyperionInsert1MilAndFindMiddle(b *testing.B) {
	defer testutil.RecoverBenchHandler(b)

	const totalRows = 1_000_000
	middleIndex := totalRows / 2
	uuids := make([]string, totalRows)

	for i := range uuids {
		uuids[i] = uuid.NewString()
	}

	for _, id := range uuids {
		if err := dbSampleInsert(id); err != nil {
			b.Fatalf("Insert failed: %v", err)
		}
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q := query.NewQuery().
			SetFilters(query.FilterTypeAnd, []query.Filter{
				{Field: SampleV1.FieldName, Op: query.OperatorTypeEqual, Value: uuids[middleIndex]},
			}).
			SetLimit(1)

		results, err := SampleV1.DbQuery(connection, q)
		if err != nil {
			b.Fatal(err)
		}

		if results[0].Name != uuids[middleIndex] {
			b.Fatalf("Expected %s, got %s", uuids[middleIndex], results[0].Name)
		}
	}
}

var inserted bool

func BenchmarkHyperionInsert100kAndSort(b *testing.B) {
	defer testutil.RecoverBenchHandler(b)

	const totalRows = 100_000
	if !inserted {
		uuids := make([]string, totalRows)
		for i := range uuids {
			uuids[i] = uuid.NewString()
			if err := dbSampleInsert(uuids[i]); err != nil {
				b.Fatalf("Insert failed: %v", err)
			}
		}
		inserted = true
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		q := query.NewQuery().
			AddOrder(query.OrderTypeAsc, SampleV1.FieldName)

		results, err := SampleV1.DbQuery(connection, q)
		if err != nil {
			b.Fatal(err)
		}

		if len(results) != totalRows {
			b.Fatalf("Expected %d results, got %d", totalRows, len(results))
		}

		for j := 1; j < len(results); j++ {
			if results[j].Name < results[j-1].Name {
				b.Fatalf("Sort order incorrect at index %d: %s < %s", j, results[j].Name, results[j-1].Name)
			}
		}
	}
}

func BenchmarkHyperionInsert100kAndQueryOlderThan15Min(b *testing.B) {
	defer testutil.RecoverBenchHandler(b)

	const totalRows = 100_000
	if !inserted {
		now := time.Now()

		for i := 0; i < totalRows; i++ {
			var birth time.Time
			if i < int(float64(totalRows)*0.8) {
				birth = now.Add(-time.Duration(16+rand.Intn(10)) * time.Minute) // Older than 15 min
			} else {
				birth = now.Add(-time.Duration(rand.Intn(10)) * time.Minute) // Within 15 min
			}

			err := dbSampleInsertWithDate(uuid.NewString(), birth)
			if err != nil {
				b.Fatalf("Insert failed: %v", err)
			}
		}
		inserted = true
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cutoff := time.Now().Add(-15 * time.Minute)
		q := query.NewQuery().
			SetFilters(query.FilterTypeAnd, []query.Filter{
				{Field: SampleV1.FieldBirth, Op: query.OperatorTypeLessThan, Value: cutoff},
			}).
			SetOrders([]query.Order{
				{Type: query.OrderTypeAsc, Field: SampleV1.FieldBirth},
			})

		results, err := SampleV1.DbQuery(connection, q)
		if err != nil {
			b.Fatal(err)
		}

		var last time.Time
		for j, r := range results {
			if r.Birth.After(cutoff) {
				b.Fatalf("Found entry newer than 15min cutoff: %v", r.Birth)
			}
			if j > 0 && r.Birth.Before(last) {
				b.Fatalf("Result not ordered correctly at index %d: %v < %v", j, r.Birth, last)
			}
			last = r.Birth
		}
	}
}
