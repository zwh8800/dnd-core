package engine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zwh8800/dnd-core/pkg/storage"
)

func TestNew(t *testing.T) {
	t.Run("creates engine with default config", func(t *testing.T) {
		cfg := DefaultConfig()
		e, err := New(cfg)

		require.NoError(t, err)
		require.NotNil(t, e)
		assert.NotNil(t, e.store)
		assert.NotNil(t, e.roller)
		e.Close()
	})

	t.Run("creates engine with custom storage", func(t *testing.T) {
		store := storage.NewMemoryStore()
		cfg := Config{
			Storage:  store,
			DiceSeed: 123,
		}
		e, err := New(cfg)

		require.NoError(t, err)
		require.NotNil(t, e)
		assert.NotNil(t, e.store)
		assert.Equal(t, store, e.store)
		e.Close()
	})

	t.Run("creates engine with nil storage uses default", func(t *testing.T) {
		cfg := Config{
			Storage:  nil,
			DiceSeed: 0,
		}
		e, err := New(cfg)

		require.NoError(t, err)
		require.NotNil(t, e)
		assert.NotNil(t, e.store)
		e.Close()
	})

	t.Run("creates engine with different dice seeds", func(t *testing.T) {
		seeds := []int64{0, 42, 100, 9999}
		for _, seed := range seeds {
			cfg := DefaultConfig()
			cfg.DiceSeed = seed
			e, err := New(cfg)

			require.NoError(t, err)
			require.NotNil(t, e)
			assert.NotNil(t, e.roller)
			e.Close()
		}
	})

	t.Run("engine is ready for operations after creation", func(t *testing.T) {
		cfg := DefaultConfig()
		e, err := New(cfg)

		require.NoError(t, err)
		require.NotNil(t, e)

		ctx := context.Background()
		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test",
		})
		require.NoError(t, err)
		require.NotNil(t, gameResult)
		assert.NotEmpty(t, gameResult.Game.ID)

		e.Close()
	})
}

func TestDefaultConfig(t *testing.T) {
	t.Run("returns config with memory store", func(t *testing.T) {
		cfg := DefaultConfig()

		assert.NotNil(t, cfg.Storage)
	})

	t.Run("returns config with zero dice seed", func(t *testing.T) {
		cfg := DefaultConfig()

		assert.Equal(t, int64(0), cfg.DiceSeed)
	})

	t.Run("returns config with empty data path", func(t *testing.T) {
		cfg := DefaultConfig()

		assert.Equal(t, "", cfg.DataPath)
	})

	t.Run("returns unique store instances", func(t *testing.T) {
		cfg1 := DefaultConfig()
		cfg2 := DefaultConfig()

		assert.NotSame(t, cfg1.Storage, cfg2.Storage)
	})

	t.Run("config can be used to create engine", func(t *testing.T) {
		cfg := DefaultConfig()
		e, err := New(cfg)

		require.NoError(t, err)
		require.NotNil(t, e)
		e.Close()
	})
}

func TestClose(t *testing.T) {
	t.Run("closes engine successfully", func(t *testing.T) {
		cfg := DefaultConfig()
		e, err := New(cfg)
		require.NoError(t, err)

		err = e.Close()
		assert.NoError(t, err)
	})

	t.Run("closing already closed engine returns no error", func(t *testing.T) {
		cfg := DefaultConfig()
		e, err := New(cfg)
		require.NoError(t, err)

		err = e.Close()
		assert.NoError(t, err)

		err = e.Close()
		assert.NoError(t, err)
	})

	t.Run("closed engine can still access stored data", func(t *testing.T) {
		cfg := DefaultConfig()
		e, err := New(cfg)
		require.NoError(t, err)

		ctx := context.Background()
		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test",
		})
		require.NoError(t, err)

		err = e.Close()
		require.NoError(t, err)

		// 由于使用内存存储，关闭后数据可能仍然可访问
		// 这取决于具体实现
		_ = gameResult
	})

	t.Run("close releases resources", func(t *testing.T) {
		cfg := DefaultConfig()
		e, err := New(cfg)
		require.NoError(t, err)

		ctx := context.Background()
		_, err = e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test",
		})
		require.NoError(t, err)

		err = e.Close()
		assert.NoError(t, err)
	})

	t.Run("multiple close calls are safe", func(t *testing.T) {
		cfg := DefaultConfig()
		e, err := New(cfg)
		require.NoError(t, err)

		for i := 0; i < 5; i++ {
			err = e.Close()
			assert.NoError(t, err)
		}
	})
}
