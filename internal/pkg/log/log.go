package log

import (
	"log"
	"log/slog"
)

var (
	Fatal = log.Fatal
)

func ParseLevel(s string) (slog.Level, error) {
	var level slog.Level
	var err = level.UnmarshalText([]byte(s))
	return level, err
}
