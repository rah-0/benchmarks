package log

import (
	"errors"
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/rah-0/nabu"
	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	testError      = errors.New("test error message")
	discard        = io.Discard
	originalStdout *os.File
	originalStderr *os.File
	devNull        *os.File
)

func init() {
	// Save original stdout/stderr
	originalStdout = os.Stdout
	originalStderr = os.Stderr

	// Open /dev/null for discarding output
	var err error
	devNull, err = os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		panic(err)
	}
}

// Setup functions for each logger
func setupNabu() {
	os.Stdout = devNull
	nabu.SetLogLevel(nabu.LevelInfo)
	nabu.SetLogOutput(nabu.OutputStdout)
}

func restoreNabu() {
	os.Stdout = originalStdout
}

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

func BenchmarkNabu_SimpleMessage(b *testing.B) {
	setupNabu()
	defer restoreNabu()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nabu.FromMessage("simple log message").Log()
	}
}

func BenchmarkZerolog_SimpleMessage(b *testing.B) {
	logger := setupZerolog()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info().Msg("simple log message")
	}
}

func BenchmarkZap_SimpleMessage(b *testing.B) {
	logger := setupZap()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("simple log message")
	}
}

func BenchmarkLogrus_SimpleMessage(b *testing.B) {
	logger := setupLogrus()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("simple log message")
	}
}

func BenchmarkSlog_SimpleMessage(b *testing.B) {
	logger := setupSlog()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("simple log message")
	}
}

func BenchmarkNabu_WithFields(b *testing.B) {
	setupNabu()
	defer restoreNabu()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nabu.FromMessage("log with fields").
			WithArgs("userID", 12345, "action", "login", "ip", "192.168.1.1").
			Log()
	}
}

func BenchmarkZerolog_WithFields(b *testing.B) {
	logger := setupZerolog()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info().
			Int("userID", 12345).
			Str("action", "login").
			Str("ip", "192.168.1.1").
			Msg("log with fields")
	}
}

func BenchmarkZap_WithFields(b *testing.B) {
	logger := setupZap()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("log with fields",
			zap.Int("userID", 12345),
			zap.String("action", "login"),
			zap.String("ip", "192.168.1.1"),
		)
	}
}

func BenchmarkLogrus_WithFields(b *testing.B) {
	logger := setupLogrus()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.WithFields(logrus.Fields{
			"userID": 12345,
			"action": "login",
			"ip":     "192.168.1.1",
		}).Info("log with fields")
	}
}

func BenchmarkSlog_WithFields(b *testing.B) {
	logger := setupSlog()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("log with fields",
			"userID", 12345,
			"action", "login",
			"ip", "192.168.1.1",
		)
	}
}

func BenchmarkNabu_Error(b *testing.B) {
	setupNabu()
	defer restoreNabu()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nabu.FromError(testError).WithMessage("operation failed").Log()
	}
}

func BenchmarkZerolog_Error(b *testing.B) {
	logger := setupZerolog()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Error().Err(testError).Msg("operation failed")
	}
}

func BenchmarkZap_Error(b *testing.B) {
	logger := setupZap()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Error("operation failed", zap.Error(testError))
	}
}

func BenchmarkLogrus_Error(b *testing.B) {
	logger := setupLogrus()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.WithError(testError).Error("operation failed")
	}
}

func BenchmarkSlog_Error(b *testing.B) {
	logger := setupSlog()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Error("operation failed", "error", testError)
	}
}

func BenchmarkNabu_ErrorWithFields(b *testing.B) {
	setupNabu()
	defer restoreNabu()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		nabu.FromError(testError).
			WithMessage("database operation failed").
			WithArgs("table", "users", "operation", "insert", "retries", 3).
			Log()
	}
}

func BenchmarkZerolog_ErrorWithFields(b *testing.B) {
	logger := setupZerolog()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
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
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.WithError(testError).WithFields(logrus.Fields{
			"table":     "users",
			"operation": "insert",
			"retries":   3,
		}).Error("database operation failed")
	}
}

func BenchmarkSlog_ErrorWithFields(b *testing.B) {
	logger := setupSlog()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Error("database operation failed",
			"error", testError,
			"table", "users",
			"operation", "insert",
			"retries", 3,
		)
	}
}

func BenchmarkNabu_ErrorChain(b *testing.B) {
	setupNabu()
	defer restoreNabu()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err1 := nabu.FromError(testError).WithMessage("database error").Log()
		err2 := nabu.FromError(err1).WithMessage("service error").Log()
		nabu.FromError(err2).WithMessage("handler error").Log()
	}
}

func BenchmarkZerolog_ErrorChain(b *testing.B) {
	logger := setupZerolog()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Error().Err(testError).Msg("database error")
		logger.Error().Err(testError).Msg("service error")
		logger.Error().Err(testError).Msg("handler error")
	}
}

func BenchmarkZap_ErrorChain(b *testing.B) {
	logger := setupZap()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Error("database error", zap.Error(testError))
		logger.Error("service error", zap.Error(testError))
		logger.Error("handler error", zap.Error(testError))
	}
}

func BenchmarkLogrus_ErrorChain(b *testing.B) {
	logger := setupLogrus()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.WithError(testError).Error("database error")
		logger.WithError(testError).Error("service error")
		logger.WithError(testError).Error("handler error")
	}
}

func BenchmarkSlog_ErrorChain(b *testing.B) {
	logger := setupSlog()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Error("database error", "error", testError)
		logger.Error("service error", "error", testError)
		logger.Error("handler error", "error", testError)
	}
}
