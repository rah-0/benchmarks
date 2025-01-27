package cache

import (
	"strconv"
	"testing"

	"github.com/google/uuid"

	"github.com/rah-0/benchmarks/util/testutil"
)

func TestCacheRedisSetGet(t *testing.T) {
	key := "test_key"
	value := "test_value"
	err := CacheRedisSet(key, value)
	if err != nil {
		t.Fatalf("Redis Set failed: %v", err)
	}

	retrievedValue, err := CacheRedisGet(key)
	if err != nil {
		t.Fatalf("Redis Get failed: %v", err)
	}

	if retrievedValue != value {
		t.Fatalf("Redis Get returned unexpected value: got %v, want %v", retrievedValue, value)
	}
}

func TestCacheRedisSetGetSizes(t *testing.T) {
	for sizeType, size := range testutil.Sizes {
		if size > 11 {
			continue
		}

		t.Run("Size: "+strconv.Itoa(sizeType), func(t *testing.T) {
			key := "test_key_size"
			message, err := testutil.GenerateMessage(sizeType)
			if err != nil {
				t.Fatalf("Failed to generate message of size %d: %v", size, err)
			}
			value := string(message)

			if err = CacheRedisSet(key, value); err != nil {
				t.Fatalf("Redis Set failed for size %d: %v", size, err)
			}

			retrievedValue, err := CacheRedisGet(key)
			if err != nil {
				t.Fatalf("Redis Get failed for size %d: %v", size, err)
			}

			if retrievedValue != value {
				t.Fatalf("Redis Get returned unexpected value for size %d: got %v, want %v", size, retrievedValue, value)
			}
		})
	}
}

func BenchmarkRedis2KB(b *testing.B)   { benchmarkRedis(b, 1) }
func BenchmarkRedis4KB(b *testing.B)   { benchmarkRedis(b, 2) }
func BenchmarkRedis8KB(b *testing.B)   { benchmarkRedis(b, 3) }
func BenchmarkRedis16KB(b *testing.B)  { benchmarkRedis(b, 4) }
func BenchmarkRedis32KB(b *testing.B)  { benchmarkRedis(b, 5) }
func BenchmarkRedis64KB(b *testing.B)  { benchmarkRedis(b, 6) }
func BenchmarkRedis128KB(b *testing.B) { benchmarkRedis(b, 7) }
func BenchmarkRedis256KB(b *testing.B) { benchmarkRedis(b, 8) }
func BenchmarkRedis512KB(b *testing.B) { benchmarkRedis(b, 9) }

func benchmarkRedis(b *testing.B, sizeType int) {
	key := uuid.NewString()

	message, err := testutil.GenerateMessage(sizeType)
	if err != nil {
		b.Fatalf("Failed to generate message: %v", err)
	}
	value := string(message)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err = CacheRedisSet(key, value); err != nil {
			b.Fatalf("Redis Set failed: %v", err)
		}

		v, err := CacheRedisGet(key)
		if err != nil {
			b.Fatalf("Redis Get failed: %v", err)
		}

		if v != value {
			b.Fatalf("Redis Get returned unexpected value: got %v, want %v", v, value)
		}
	}
}
