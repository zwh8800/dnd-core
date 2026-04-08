package engine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zwh8800/dnd-core/pkg/model"
)

func TestApplyPoison(t *testing.T) {
	t.Run("applies poison to weapon successfully", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.ApplyPoison(ctx, ApplyPoisonRequest{
			GameID:   gameResult.Game.ID,
			ActorID:  model.NewID(),
			PoisonID: "basic-poison",
			WeaponID: "sword-1",
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.PoisonInstance)
		assert.Equal(t, "basic-poison", result.PoisonInstance.PoisonID)
		assert.Equal(t, "sword-1", result.PoisonInstance.AppliedTo)
		assert.Equal(t, 1, result.PoisonInstance.RemainingUses)
		assert.Equal(t, "1 minute", result.PoisonInstance.ExpiresAfter)
		assert.Contains(t, result.Message, "基础毒药")
	})

	t.Run("applies different poison type", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.ApplyPoison(ctx, ApplyPoisonRequest{
			GameID:   gameResult.Game.ID,
			PoisonID: "drow-poison",
			WeaponID: "crossbow-1",
		})

		require.NoError(t, err)
		assert.Equal(t, "drow-poison", result.PoisonInstance.PoisonID)
		assert.Contains(t, result.Message, "卓尔")
	})

	t.Run("returns error for nonexistent poison", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.ApplyPoison(ctx, ApplyPoisonRequest{
			GameID:   gameResult.Game.ID,
			PoisonID: "nonexistent-poison",
			WeaponID: "sword-1",
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "poison data")
	})

	t.Run("returns error when game not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.ApplyPoison(ctx, ApplyPoisonRequest{
			GameID:   model.NewID(),
			PoisonID: "basic-poison",
			WeaponID: "sword-1",
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("applies powerful poison", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.ApplyPoison(ctx, ApplyPoisonRequest{
			GameID:   gameResult.Game.ID,
			PoisonID: "purple-worm-poison",
			WeaponID: "spear-1",
		})

		require.NoError(t, err)
		assert.Equal(t, "purple-worm-poison", result.PoisonInstance.PoisonID)
		assert.Contains(t, result.Message, "紫虫")
	})

	t.Run("applies poison with empty weapon ID", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.ApplyPoison(ctx, ApplyPoisonRequest{
			GameID:   gameResult.Game.ID,
			PoisonID: "basic-poison",
			WeaponID: "",
		})

		require.NoError(t, err)
		assert.Equal(t, "", result.PoisonInstance.AppliedTo)
	})
}

func TestResolvePoisonEffect(t *testing.T) {
	t.Run("resolves poison effect with save DC", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.ResolvePoisonEffect(ctx, ResolvePoisonEffectRequest{
			GameID:  gameResult.Game.ID,
			ActorID: model.NewID(),
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Greater(t, result.SaveDC, 0)
		assert.False(t, result.SaveSuccess)
		assert.Contains(t, result.Message, "毒药发作")
	})

	t.Run("returns error when game not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.ResolvePoisonEffect(ctx, ResolvePoisonEffectRequest{
			GameID:  model.NewID(),
			ActorID: model.NewID(),
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("multiple poison resolutions return consistent results", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		for i := 0; i < 3; i++ {
			result, err := e.ResolvePoisonEffect(ctx, ResolvePoisonEffectRequest{
				GameID:  gameResult.Game.ID,
				ActorID: model.NewID(),
			})
			require.NoError(t, err)
			assert.Greater(t, result.SaveDC, 0)
			assert.False(t, result.SaveSuccess)
		}
	})

	t.Run("poison effect message contains description", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.ResolvePoisonEffect(ctx, ResolvePoisonEffectRequest{
			GameID:  gameResult.Game.ID,
			ActorID: model.NewID(),
		})

		require.NoError(t, err)
		assert.NotEmpty(t, result.Message)
	})

	t.Run("resolve poison uses write lock", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		done := make(chan bool, 2)
		go func() {
			_, err := e.ResolvePoisonEffect(ctx, ResolvePoisonEffectRequest{
				GameID:  gameResult.Game.ID,
				ActorID: model.NewID(),
			})
			assert.NoError(t, err)
			done <- true
		}()
		go func() {
			_, err := e.ResolvePoisonEffect(ctx, ResolvePoisonEffectRequest{
				GameID:  gameResult.Game.ID,
				ActorID: model.NewID(),
			})
			assert.NoError(t, err)
			done <- true
		}()

		<-done
		<-done
	})
}

func TestRemovePoison(t *testing.T) {
	t.Run("removes poison from weapon successfully", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.RemovePoison(ctx, RemovePoisonRequest{
			GameID:   gameResult.Game.ID,
			ActorID:  model.NewID(),
			WeaponID: "sword-1",
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.Success)
		assert.Contains(t, result.Message, "已移除")
	})

	t.Run("removes poison from different weapon", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.RemovePoison(ctx, RemovePoisonRequest{
			GameID:   gameResult.Game.ID,
			WeaponID: "crossbow-1",
		})

		require.NoError(t, err)
		assert.True(t, result.Success)
	})

	t.Run("removes poison with empty weapon ID", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.RemovePoison(ctx, RemovePoisonRequest{
			GameID:   gameResult.Game.ID,
			WeaponID: "",
		})

		require.NoError(t, err)
		assert.True(t, result.Success)
	})

	t.Run("returns error when game not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.RemovePoison(ctx, RemovePoisonRequest{
			GameID:   model.NewID(),
			WeaponID: "sword-1",
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("multiple removes return consistent results", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		for i := 0; i < 3; i++ {
			result, err := e.RemovePoison(ctx, RemovePoisonRequest{
				GameID:   gameResult.Game.ID,
				WeaponID: "weapon",
			})
			require.NoError(t, err)
			assert.True(t, result.Success)
		}
	})
}
