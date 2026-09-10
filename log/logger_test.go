package log

import (
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/rah-0/slogx"
	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// The comparison measures correlation fields, excluding identifier generation.
const traceID = "0123456789abcdef0123456789abcdef"

var (
	testError = errors.New("test error message")
	discard   = io.Discard
)

// Setup functions for each logger
func setupZerolog() zerolog.Logger {
	return zerolog.New(discard).With().Timestamp().Logger()
}

func setupZap() *zap.Logger {
	config := zap.NewProductionEncoderConfig()
	config.TimeKey = "timestamp"
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		zapcore.AddSync(discard),
		zapcore.InfoLevel,
	)
	return zap.New(core)
}

func setupLogrus() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(discard)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)
	return logger
}

func setupSlog() *slog.Logger {
	return slog.New(slog.NewJSONHandler(discard, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

func setupSlogx() *slog.Logger {
	return slogx.New(slogx.Options{
		Format: slogx.JSON,
		Level:  slog.LevelInfo,
		Writer: discard,
	})
}

func BenchmarkZerolog_SimpleMessage(b *testing.B) {
	logger := setupZerolog()
	for b.Loop() {
		logger.Info().Msg("simple log message")
	}
}

func BenchmarkZap_SimpleMessage(b *testing.B) {
	logger := setupZap()
	for b.Loop() {
		logger.Info("simple log message")
	}
}

func BenchmarkLogrus_SimpleMessage(b *testing.B) {
	logger := setupLogrus()
	for b.Loop() {
		logger.Info("simple log message")
	}
}

func BenchmarkSlog_SimpleMessage(b *testing.B) {
	logger := setupSlog()
	for b.Loop() {
		logger.Info("simple log message")
	}
}

func BenchmarkSlogx_SimpleMessage(b *testing.B) {
	logger := setupSlogx()
	for b.Loop() {
		logger.Info("simple log message")
	}
}

func BenchmarkZerolog_WithFields(b *testing.B) {
	logger := setupZerolog()
	for b.Loop() {
		logger.Info().
			Int("userID", 12345).
			Str("action", "login").
			Str("ip", "192.168.1.1").
			Msg("log with fields")
	}
}

func BenchmarkZap_WithFields(b *testing.B) {
	logger := setupZap()
	for b.Loop() {
		logger.Info("log with fields",
			zap.Int("userID", 12345),
			zap.String("action", "login"),
			zap.String("ip", "192.168.1.1"),
		)
	}
}

func BenchmarkLogrus_WithFields(b *testing.B) {
	logger := setupLogrus()
	for b.Loop() {
		logger.WithFields(logrus.Fields{
			"userID": 12345,
			"action": "login",
			"ip":     "192.168.1.1",
		}).Info("log with fields")
	}
}

func BenchmarkSlog_WithFields(b *testing.B) {
	logger := setupSlog()
	for b.Loop() {
		logger.Info("log with fields",
			"userID", 12345,
			"action", "login",
			"ip", "192.168.1.1",
		)
	}
}

func BenchmarkSlogx_WithFields(b *testing.B) {
	logger := setupSlogx()
	for b.Loop() {
		logger.Info("log with fields",
			"userID", 12345,
			"action", "login",
			"ip", "192.168.1.1",
		)
	}
}

func BenchmarkZerolog_Error(b *testing.B) {
	logger := setupZerolog()
	for b.Loop() {
		logger.Error().Err(testError).Msg("operation failed")
	}
}

func BenchmarkZap_Error(b *testing.B) {
	logger := setupZap()
	for b.Loop() {
		logger.Error("operation failed", zap.Error(testError))
	}
}

func BenchmarkLogrus_Error(b *testing.B) {
	logger := setupLogrus()
	for b.Loop() {
		logger.WithError(testError).Error("operation failed")
	}
}

func BenchmarkSlog_Error(b *testing.B) {
	logger := setupSlog()
	for b.Loop() {
		logger.Error("operation failed", "error", testError)
	}
}

func BenchmarkSlogx_Error(b *testing.B) {
	logger := setupSlogx()
	for b.Loop() {
		logger.Error("operation failed", "error", testError)
	}
}

func BenchmarkZerolog_ErrorWithFields(b *testing.B) {
	logger := setupZerolog()
	for b.Loop() {
		logger.Error().
			Err(testError).
			Str("table", "users").
			Str("operation", "insert").
			Int("retries", 3).
			Msg("database operation failed")
	}
}

func BenchmarkZap_ErrorWithFields(b *testing.B) {
	logger := setupZap()
	for b.Loop() {
		logger.Error("database operation failed",
			zap.Error(testError),
			zap.String("table", "users"),
			zap.String("operation", "insert"),
			zap.Int("retries", 3),
		)
	}
}

func BenchmarkLogrus_ErrorWithFields(b *testing.B) {
	logger := setupLogrus()
	for b.Loop() {
		logger.WithError(testError).WithFields(logrus.Fields{
			"table":     "users",
			"operation": "insert",
			"retries":   3,
		}).Error("database operation failed")
	}
}

func BenchmarkSlog_ErrorWithFields(b *testing.B) {
	logger := setupSlog()
	for b.Loop() {
		logger.Error("database operation failed",
			"error", testError,
			"table", "users",
			"operation", "insert",
			"retries", 3,
		)
	}
}

func BenchmarkSlogx_ErrorWithFields(b *testing.B) {
	logger := setupSlogx()
	for b.Loop() {
		logger.Error("database operation failed",
			"error", testError,
			"table", "users",
			"operation", "insert",
			"retries", 3,
		)
	}
}

func BenchmarkZerolog_ErrorChain(b *testing.B) {
	logger := setupZerolog()
	for b.Loop() {
		logger.Error().Str("trace_id", traceID).Err(testError).Msg("database error")
		logger.Error().Str("trace_id", traceID).Msg("service error")
		logger.Error().Str("trace_id", traceID).Msg("handler error")
	}
}

func BenchmarkZap_ErrorChain(b *testing.B) {
	logger := setupZap()
	for b.Loop() {
		logger.Error("database error", zap.String("trace_id", traceID), zap.Error(testError))
		logger.Error("service error", zap.String("trace_id", traceID))
		logger.Error("handler error", zap.String("trace_id", traceID))
	}
}

func BenchmarkLogrus_ErrorChain(b *testing.B) {
	logger := setupLogrus()
	for b.Loop() {
		logger.WithField("trace_id", traceID).WithError(testError).Error("database error")
		logger.WithField("trace_id", traceID).Error("service error")
		logger.WithField("trace_id", traceID).Error("handler error")
	}
}

func BenchmarkSlog_ErrorChain(b *testing.B) {
	logger := setupSlog()
	for b.Loop() {
		logger.Error("database error", "trace_id", traceID, "error", testError)
		logger.Error("service error", "trace_id", traceID)
		logger.Error("handler error", "trace_id", traceID)
	}
}

func BenchmarkSlogx_ErrorChain(b *testing.B) {
	logger := setupSlogx()
	for b.Loop() {
		logger.Error("database error", "trace_id", traceID, "error", testError)
		logger.Error("service error", "trace_id", traceID)
		logger.Error("handler error", "trace_id", traceID)
	}
}
