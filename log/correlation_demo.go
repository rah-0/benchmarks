// +build ignore

package main

import (
	"errors"
	"log/slog"
	"os"

	"github.com/rah-0/nabu"
	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var testError = errors.New("database connection failed")

func main() {
	println("=== NABU: Automatic UUID Correlation ===")
	nabu.SetLogLevel(nabu.LevelInfo)
	nabu.SetLogOutput(nabu.OutputStdout)
	
	// nabu automatically correlates via UUID when wrapping errors
	err1 := nabu.FromError(testError).WithMessage("database error").Log()
	err2 := nabu.FromError(err1).WithMessage("service error").Log()
	nabu.FromError(err2).WithMessage("handler error").Log()
	
	println("\n=== ZEROLOG: Manual trace_id Correlation ===")
	zerologLogger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	traceID := "trace-12345"
	zerologLogger.Error().Str("trace_id", traceID).Err(testError).Msg("database error")
	zerologLogger.Error().Str("trace_id", traceID).Msg("service error")
	zerologLogger.Error().Str("trace_id", traceID).Msg("handler error")
	
	println("\n=== ZAP: Manual trace_id Correlation ===")
	config := zap.NewProductionEncoderConfig()
	config.TimeKey = "timestamp"
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		zapcore.AddSync(os.Stdout),
		zapcore.InfoLevel,
	)
	zapLogger := zap.New(core)
	zapLogger.Error("database error", zap.String("trace_id", traceID), zap.Error(testError))
	zapLogger.Error("service error", zap.String("trace_id", traceID))
	zapLogger.Error("handler error", zap.String("trace_id", traceID))
	
	println("\n=== LOGRUS: Manual trace_id Correlation ===")
	logrusLogger := logrus.New()
	logrusLogger.SetOutput(os.Stdout)
	logrusLogger.SetFormatter(&logrus.JSONFormatter{})
	logrusLogger.SetLevel(logrus.InfoLevel)
	logrusLogger.WithField("trace_id", traceID).WithError(testError).Error("database error")
	logrusLogger.WithField("trace_id", traceID).Error("service error")
	logrusLogger.WithField("trace_id", traceID).Error("handler error")
	
	println("\n=== SLOG: Manual trace_id Correlation ===")
	slogLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slogLogger.Error("database error", "trace_id", traceID, "error", testError)
	slogLogger.Error("service error", "trace_id", traceID)
	slogLogger.Error("handler error", "trace_id", traceID)
}

