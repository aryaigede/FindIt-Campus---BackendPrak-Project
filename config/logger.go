package config

import (
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(env string) *slog.Logger {
	_ = os.MkdirAll("logs", 0755)

	fileRotator := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     14,
		Compress:   true,
	}

	multiWriter := io.MultiWriter(os.Stdout, fileRotator)

	level := slog.LevelInfo
	if env == "development" {
		level = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level: level,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}