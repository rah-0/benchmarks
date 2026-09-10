package log

import (
	"context"
	"log/slog"
	"testing"

	"github.com/rah-0/slogx"
)

// Context preparation is excluded: these cases measure enrichment and encoding
// when logging with a context that already carries request metadata.
func BenchmarkSlogxFeatures_Context(b *testing.B) {
	attrs := []slog.Attr{
		slog.String("request_id", "request-123"),
		slog.Int("user_id", 12345),
		slog.String("route", "/items"),
		slog.String("method", "GET"),
		slog.String("region", "eu-west"),
		slog.String("service", "api"),
		slog.String("version", "1.0.0"),
		slog.String("tenant", "example"),
		slog.Int("attempt", 1),
		slog.Bool("cached", false),
	}
	base := context.Background()
	for _, scenario := range []struct {
		name string
		ctx  context.Context
	}{
		{name: "NoAttrs", ctx: base},
		{name: "Attrs3", ctx: slogx.WithAttrs(base, attrs[:3]...)},
		{name: "Attrs10", ctx: slogx.WithAttrs(base, attrs...)},
		{name: "TraceOnly", ctx: slogx.ContextWithSpanContext(base, slogx.SpanContext{
			TraceID:    "0123456789abcdef0123456789abcdef",
			SpanID:     "0123456789abcdef",
			TraceFlags: 1,
		})},
	} {
		b.Run(scenario.name, func(b *testing.B) {
			logger := setupSlogx()
			b.ReportAllocs()
			for b.Loop() {
				logger.LogAttrs(scenario.ctx, slog.LevelInfo, "operation completed", slog.String("action", "read"))
			}
		})
	}
}

// One operation creates and ends one span. Recording includes JSON emission;
// non-recording still includes IDs, context allocation, and lifecycle work.
func BenchmarkSlogxFeatures_Span(b *testing.B) {
	base := context.Background()
	parent := slogx.SpanContext{
		TraceID:    "0123456789abcdef0123456789abcdef",
		SpanID:     "0123456789abcdef",
		TraceFlags: 1,
		Remote:     true,
	}
	child := slogx.ContextWithSpanContext(base, parent)
	parent.TraceFlags = 0
	unsampledChild := slogx.ContextWithSpanContext(base, parent)
	for _, scenario := range []struct {
		name          string
		ctx           context.Context
		levelDisabled bool
	}{
		{name: "RootRecording", ctx: base},
		{name: "ChildRecording", ctx: child},
		{name: "RootLevelDisabled", ctx: base, levelDisabled: true},
		{name: "ChildLevelDisabled", ctx: child, levelDisabled: true},
		{name: "ChildUnsampled", ctx: unsampledChild},
	} {
		b.Run(scenario.name, func(b *testing.B) {
			logger := setupSlogx()
			if scenario.levelDisabled {
				logger = slogx.New(slogx.Options{
					Format: slogx.JSON,
					Level:  slog.LevelWarn,
					Writer: discard,
				})
			}
			previous := slog.Default()
			slog.SetDefault(logger)
			b.Cleanup(func() { slog.SetDefault(previous) })
			b.ReportAllocs()
			for b.Loop() {
				_, span := slogx.StartSpan(scenario.ctx, "read items",
					slog.String("table", "items"),
					slog.Int("limit", 20),
				)
				span.End()
			}
		})
	}
}

func BenchmarkSlogxFeatures_StructuredError(b *testing.B) {
	b.Run("WrapThreeLayers", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			err := slogx.Wrap(testError, "database error", "table", "items")
			err = slogx.Wrap(err, "service error", "operation", "read")
			_ = slogx.Wrap(err, "handler error", "route", "/items")
		}
	})
	b.Run("LogThreeLayers", func(b *testing.B) {
		logger := setupSlogx()
		err := slogx.Wrap(testError, "database error", "table", "items")
		err = slogx.Wrap(err, "service error", "operation", "read")
		err = slogx.Wrap(err, "handler error", "route", "/items")
		b.ReportAllocs()
		for b.Loop() {
			logger.Error("operation failed", "err", err)
		}
	})
}
