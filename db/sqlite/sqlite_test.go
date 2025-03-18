package sqlite

import (
	"os"
	"testing"

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
			os.Remove("./benchmark.db")
			return nil
		},
	})
}

func BenchmarkSQLiteSingleInsertFixedData(b *testing.B) {
	defer testutil.RecoverBenchHandler(b)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := dbTableSampleAInsert("9xKf3QpLm2Ry7UbHt6NwEjVg8As5OcIy4B"); err != nil {
			b.Fatalf("Insert failed on iteration %d: %v", i, err)
		}
	}
}

func BenchmarkSQLiteSingleInsertRandomData(b *testing.B) {
	defer testutil.RecoverBenchHandler(b)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := dbTableSampleAInsert(uuid.NewString()); err != nil {
			b.Fatalf("Insert failed on iteration %d: %v", i, err)
		}
	}
}
