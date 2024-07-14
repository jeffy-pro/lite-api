package log

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantLevel slog.Level
		wantErr   bool
	}{
		{
			name:      "Valid DEBUG level",
			input:     "DEBUG",
			wantLevel: slog.LevelDebug,
			wantErr:   false,
		},
		{
			name:      "Valid INFO level",
			input:     "INFO",
			wantLevel: slog.LevelInfo,
			wantErr:   false,
		},
		{
			name:      "Valid WARN level",
			input:     "WARN",
			wantLevel: slog.LevelWarn,
			wantErr:   false,
		},
		{
			name:      "Valid ERROR level",
			input:     "ERROR",
			wantLevel: slog.LevelError,
			wantErr:   false,
		},
		{
			name:      "Case insensitive input",
			input:     "debug",
			wantLevel: slog.LevelDebug,
			wantErr:   false,
		},
		{
			name:      "Invalid level",
			input:     "INVALID",
			wantLevel: slog.Level(0),
			wantErr:   true,
		},
		{
			name:      "Empty string",
			input:     "",
			wantLevel: slog.Level(0),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLevel, err := ParseLevel(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.Equal(t, tt.wantLevel, gotLevel)

		})
	}
}
