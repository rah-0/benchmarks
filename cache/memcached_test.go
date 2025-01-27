package cache

import (
	"strconv"
	"testing"

	"github.com/google/uuid"

	"github.com/rah-0/benchmarks/util/testutil"
)

func TestCacheMemcacheSetGet(t *testing.T) {
	key := "test_key"
	value := "test_value"
	err := CacheMemcacheSet(key, value)
	if err != nil {
		t.Fatalf("Memcache Set failed: %v", err)
	}

	retrievedValue, err := CacheMemcacheGet(key)
	if err != nil {
		t.Fatalf("Memcache Get failed: %v", err)
	}

	if retrievedValue != value {
		t.Fatalf("Memcache Get returned unexpected value: got %v, want %v", retrievedValue, value)
	}
}

func TestCacheMemcacheSetGetSizes(t *testing.T) {
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

			if err = CacheMemcacheSet(key, value); err != nil {
				t.Fatalf("Memcache Set failed for size %d: %v", size, err)
			}

			retrievedValue, err := CacheMemcacheGet(key)
			if err != nil {
				t.Fatalf("Memcache Get failed for size %d: %v", size, err)
			}

			if retrievedValue != value {
				t.Fatalf("Memcache Get returned unexpected value for size %d: got %v, want %v", size, retrievedValue, value)
			}
		})
	}
}

func BenchmarkMemcache2KB(b *testing.B)   { benchmarkMemcache(b, 1) }
func BenchmarkMemcache4KB(b *testing.B)   { benchmarkMemcache(b, 2) }
func BenchmarkMemcache8KB(b *testing.B)   { benchmarkMemcache(b, 3) }
func BenchmarkMemcache16KB(b *testing.B)  { benchmarkMemcache(b, 4) }
func BenchmarkMemcache32KB(b *testing.B)  { benchmarkMemcache(b, 5) }
func BenchmarkMemcache64KB(b *testing.B)  { benchmarkMemcache(b, 6) }
func BenchmarkMemcache128KB(b *testing.B) { benchmarkMemcache(b, 7) }
func BenchmarkMemcache256KB(b *testing.B) { benchmarkMemcache(b, 8) }
func BenchmarkMemcache512KB(b *testing.B) { benchmarkMemcache(b, 9) }

func benchmarkMemcache(b *testing.B, sizeType int) {
	key := uuid.NewString()

	message, err := testutil.GenerateMessage(sizeType)
	if err != nil {
		b.Fatalf("Failed to generate message: %v", err)
	}
	value := string(message)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err = CacheMemcacheSet(key, value); err != nil {
			b.Fatalf("Memcache Set failed: %v", err)
		}

		v, err := CacheMemcacheGet(key)
		if err != nil {
			b.Fatalf("Memcache Get failed: %v", err)
		}

		if v != value {
			b.Fatalf("Memcache Get returned unexpected value: got %v, want %v", v, value)
		}
	}
}
