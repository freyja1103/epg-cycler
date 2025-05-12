package logging

import (
	"context"
	"log/slog"
)

func SrcLog(title, basename, number string) {
	slog.LogAttrs(context.Background(), slog.LevelInfo, "TS info from EPG", slog.String("Title", title), slog.String("Basename", basename), slog.String("Episode", number))
}

func InfoAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	slog.LogAttrs(ctx, level, msg, attrs...)
}

func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

func Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}
