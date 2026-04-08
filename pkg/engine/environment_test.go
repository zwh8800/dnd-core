package engine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zwh8800/dnd-core/pkg/model"
)

func TestSetEnvironment2(t *testing.T) {
	t.Run("sets environment to bright light", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.SetEnvironment(ctx, SetEnvironmentRequest{
			GameID:  gameResult.Game.ID,
			SceneID: model.NewID(),
			EnvType: model.EnvBrightLight,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, model.EnvBrightLight, result.Environment.Type)
		assert.NotEmpty(t, result.Environment.Description)
		assert.Contains(t, result.Message, "bright_light")
	})

	t.Run("sets environment to darkness", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.SetEnvironment(ctx, SetEnvironmentRequest{
			GameID:  gameResult.Game.ID,
			SceneID: model.NewID(),
			EnvType: model.EnvDarkness,
		})

		require.NoError(t, err)
		assert.Equal(t, model.EnvDarkness, result.Environment.Type)
	})

	t.Run("sets environment to extreme cold", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.SetEnvironment(ctx, SetEnvironmentRequest{
			GameID:  gameResult.Game.ID,
			SceneID: model.NewID(),
			EnvType: model.EnvExtremeCold,
		})

		require.NoError(t, err)
		assert.Equal(t, model.EnvExtremeCold, result.Environment.Type)
		assert.GreaterOrEqual(t, result.Environment.SaveDC, 0)
	})

	t.Run("sets environment to extreme heat", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.SetEnvironment(ctx, SetEnvironmentRequest{
			GameID:  gameResult.Game.ID,
			SceneID: model.NewID(),
			EnvType: model.EnvExtremeHeat,
		})

		require.NoError(t, err)
		assert.Equal(t, model.EnvExtremeHeat, result.Environment.Type)
	})

	t.Run("returns error when game not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.SetEnvironment(ctx, SetEnvironmentRequest{
			GameID:  model.NewID(),
			EnvType: model.EnvBrightLight,
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("can set environment multiple times", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		envs := []model.EnvironmentType{
			model.EnvBrightLight,
			model.EnvDarkness,
			model.EnvSmoke,
		}

		for _, env := range envs {
			result, err := e.SetEnvironment(ctx, SetEnvironmentRequest{
				GameID:  gameResult.Game.ID,
				EnvType: env,
			})
			require.NoError(t, err)
			assert.Equal(t, env, result.Environment.Type)
		}
	})
}

func TestResolveEnvironmentalDamage2(t *testing.T) {
	t.Run("resolves damage for extreme cold", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.ResolveEnvironmentalDamage(ctx, ResolveEnvironmentalDamageRequest{
			GameID:          gameResult.Game.ID,
			ActorID:         model.NewID(),
			EnvType:         model.EnvExtremeCold,
			ExposureMinutes: 60,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, result.SaveDC, 0)
		assert.False(t, result.SaveSuccess)
		assert.Contains(t, result.Message, "环境伤害")
	})

	t.Run("resolves damage for extreme heat", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.ResolveEnvironmentalDamage(ctx, ResolveEnvironmentalDamageRequest{
			GameID:          gameResult.Game.ID,
			ActorID:         model.NewID(),
			EnvType:         model.EnvExtremeHeat,
			ExposureMinutes: 30,
		})

		require.NoError(t, err)
		assert.GreaterOrEqual(t, result.SaveDC, 0)
	})

	t.Run("resolves damage with zero exposure", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.ResolveEnvironmentalDamage(ctx, ResolveEnvironmentalDamageRequest{
			GameID:          gameResult.Game.ID,
			ActorID:         model.NewID(),
			EnvType:         model.EnvExtremeCold,
			ExposureMinutes: 0,
		})

		require.NoError(t, err)
		assert.GreaterOrEqual(t, result.SaveDC, 0)
	})

	t.Run("returns error when game not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.ResolveEnvironmentalDamage(ctx, ResolveEnvironmentalDamageRequest{
			GameID:  model.NewID(),
			EnvType: model.EnvExtremeCold,
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("resolves damage for high altitude", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.ResolveEnvironmentalDamage(ctx, ResolveEnvironmentalDamageRequest{
			GameID:          gameResult.Game.ID,
			ActorID:         model.NewID(),
			EnvType:         model.EnvHighAltitude,
			ExposureMinutes: 120,
		})

		require.NoError(t, err)
		assert.GreaterOrEqual(t, result.SaveDC, 0)
		assert.NotEmpty(t, result.Message)
	})
}
