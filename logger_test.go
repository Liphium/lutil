package util

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExistingLoggerHonorsMagicEnvironment(t *testing.T) {
	t.Setenv("MINIMUM_LOG_LEVEL", "")
	logger := NewLogger("test")
	ctx := context.Background()
	assert.False(t, logger.Enabled(ctx, slog.LevelDebug))
	assert.True(t, logger.Enabled(ctx, slog.LevelInfo))

	// Magic sets this after packages have initialized their loggers.
	t.Setenv("MINIMUM_LOG_LEVEL", "debug")
	assert.True(t, logger.Enabled(ctx, slog.LevelDebug))
	assert.True(t, logger.With("field", "value").Enabled(ctx, slog.LevelDebug))

	t.Setenv("MINIMUM_LOG_LEVEL", "info")
	assert.False(t, logger.Enabled(ctx, slog.LevelDebug))
	t.Setenv("MINIMUM_LOG_LEVEL", "invalid")
	assert.False(t, logger.Enabled(ctx, slog.LevelDebug))
	assert.True(t, logger.Enabled(ctx, slog.LevelInfo))
}
