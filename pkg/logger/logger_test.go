package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid console config",
			config: Config{
				Level:    "info",
				Encoding: "console",
			},
			wantErr: false,
		},
		{
			name: "valid json config",
			config: Config{
				Level:    "debug",
				Encoding: "json",
			},
			wantErr: false,
		},
		{
			name: "invalid level",
			config: Config{
				Level: "unknown",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := New(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, logger)

				// Verify level
				if tt.config.Level != "" {
					expectedLevel, _ := zapcore.ParseLevel(tt.config.Level)
					assert.True(t, logger.Core().Enabled(expectedLevel))
				}
			}
		})
	}
}
