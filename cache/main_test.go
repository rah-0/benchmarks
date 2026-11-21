package cache

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	redisImage     = "redis:8.10.3@sha256:1f88291b4c57bdf249f02eb085a8669889b2d2954e91d4de5118a5da6b78dd53"
	memcachedImage = "memcached:1.6.45@sha256:a395a1bcf5e0e2b0fa4a0142e12b4e24a2320ba4b7c6ffa8655a1a86e8f29dee"
)

func TestMain(m *testing.M) {
	os.Exit(runCacheTests(m))
}

func runCacheTests(m *testing.M) (code int) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	code = 1
	cleanup := func(container testcontainers.Container) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := testcontainers.TerminateContainer(container, testcontainers.StopContext(ctx)); err != nil {
			fmt.Fprintln(os.Stderr, err)
			code = 1
		}
	}

	redisContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        redisImage,
			ExposedPorts: []string{"6379/tcp"},
			Cmd:          []string{"redis-server", "--save", "", "--appendonly", "no"},
			WaitingFor:   wait.ForListeningPort("6379/tcp"),
		},
		Started: true,
	})
	defer cleanup(redisContainer)
	if err != nil {
		fmt.Fprintln(os.Stderr, "start Redis:", err)
		return
	}
	redisAddress, err := redisContainer.Endpoint(ctx, "")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Redis endpoint:", err)
		return
	}
	redisClient = redis.NewClient(&redis.Options{Addr: redisAddress})
	defer redisClient.Close()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Redis readiness:", err)
		return
	}

	memcachedContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        memcachedImage,
			ExposedPorts: []string{"11211/tcp"},
			// Size tests include 10 MiB values, above the default 1 MiB limit.
			Cmd:        []string{"memcached", "-I", "16m"},
			WaitingFor: wait.ForListeningPort("11211/tcp"),
		},
		Started: true,
	})
	defer cleanup(memcachedContainer)
	if err != nil {
		fmt.Fprintln(os.Stderr, "start Memcached:", err)
		return
	}
	memcachedAddress, err := memcachedContainer.Endpoint(ctx, "")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Memcached endpoint:", err)
		return
	}
	memcacheClient = memcache.New(memcachedAddress)
	defer memcacheClient.Close()
	if err := memcacheClient.Ping(); err != nil {
		fmt.Fprintln(os.Stderr, "Memcached readiness:", err)
		return
	}

	return m.Run()
}
