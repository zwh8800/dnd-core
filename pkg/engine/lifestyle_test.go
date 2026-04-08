package engine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zwh8800/dnd-core/pkg/model"
)

func TestSetLifestyle(t *testing.T) {
	t.Run("sets lifestyle to comfortable successfully", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.SetLifestyle(ctx, SetLifestyleRequest{
			GameID: gameResult.Game.ID,
			Tier:   model.LifestyleComfortable,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, model.LifestyleComfortable, result.Tier)
		assert.Greater(t, result.DailyCost, 0)
		assert.NotEmpty(t, result.Description)
		assert.Contains(t, result.Message, "comfortable")
	})

	t.Run("sets lifestyle to wealthy", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.SetLifestyle(ctx, SetLifestyleRequest{
			GameID: gameResult.Game.ID,
			Tier:   model.LifestyleWealthy,
		})

		require.NoError(t, err)
		assert.Equal(t, model.LifestyleWealthy, result.Tier)
		assert.Greater(t, result.DailyCost, 0)
	})

	t.Run("sets lifestyle to aristocratic", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.SetLifestyle(ctx, SetLifestyleRequest{
			GameID: gameResult.Game.ID,
			Tier:   model.LifestyleAristocratic,
		})

		require.NoError(t, err)
		assert.Equal(t, model.LifestyleAristocratic, result.Tier)
	})

	t.Run("returns error for invalid lifestyle tier", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.SetLifestyle(ctx, SetLifestyleRequest{
			GameID: gameResult.Game.ID,
			Tier:   model.LifestyleTier("invalid"),
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid lifestyle tier")
	})

	t.Run("returns error when game not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.SetLifestyle(ctx, SetLifestyleRequest{
			GameID: model.NewID(),
			Tier:   model.LifestyleComfortable,
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("can update lifestyle multiple times", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		tiers := []model.LifestyleTier{
			model.LifestylePoor,
			model.LifestyleModest,
			model.LifestyleComfortable,
		}

		for _, tier := range tiers {
			result, err := e.SetLifestyle(ctx, SetLifestyleRequest{
				GameID: gameResult.Game.ID,
				Tier:   tier,
			})
			require.NoError(t, err)
			assert.Equal(t, tier, result.Tier)
		}
	})

	t.Run("sets lifestyle to wretched", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.SetLifestyle(ctx, SetLifestyleRequest{
			GameID: gameResult.Game.ID,
			Tier:   model.LifestyleWretched,
		})

		require.NoError(t, err)
		assert.Equal(t, model.LifestyleWretched, result.Tier)
	})
}

func TestAdvanceGameTime(t *testing.T) {
	t.Run("advances time by one day", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		_, err = e.SetLifestyle(ctx, SetLifestyleRequest{
			GameID: gameResult.Game.ID,
			Tier:   model.LifestyleComfortable,
		})
		require.NoError(t, err)

		result, err := e.AdvanceGameTime(ctx, AdvanceGameTimeRequest{
			GameID: gameResult.Game.ID,
			Days:   1,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 1, result.DaysAdvanced)
		assert.Greater(t, result.TotalCost, 0)
		assert.True(t, result.PaymentSuccess)
		assert.Contains(t, result.Message, "1 天")
	})

	t.Run("advances time by multiple days", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		_, err = e.SetLifestyle(ctx, SetLifestyleRequest{
			GameID: gameResult.Game.ID,
			Tier:   model.LifestyleComfortable,
		})
		require.NoError(t, err)

		result, err := e.AdvanceGameTime(ctx, AdvanceGameTimeRequest{
			GameID: gameResult.Game.ID,
			Days:   7,
		})

		require.NoError(t, err)
		assert.Equal(t, 7, result.DaysAdvanced)
		assert.Contains(t, result.Message, "7 天")
	})

	t.Run("returns error when no lifestyle set", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.AdvanceGameTime(ctx, AdvanceGameTimeRequest{
			GameID: gameResult.Game.ID,
			Days:   1,
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "no lifestyle set")
	})

	t.Run("returns error when game not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.AdvanceGameTime(ctx, AdvanceGameTimeRequest{
			GameID: model.NewID(),
			Days:   1,
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("advances time with wealthy lifestyle costs more", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		_, err = e.SetLifestyle(ctx, SetLifestyleRequest{
			GameID: gameResult.Game.ID,
			Tier:   model.LifestyleWealthy,
		})
		require.NoError(t, err)

		result, err := e.AdvanceGameTime(ctx, AdvanceGameTimeRequest{
			GameID: gameResult.Game.ID,
			Days:   1,
		})

		require.NoError(t, err)
		wealthyCost := result.TotalCost

		// Now set to poor and compare
		_, err = e.SetLifestyle(ctx, SetLifestyleRequest{
			GameID: gameResult.Game.ID,
			Tier:   model.LifestylePoor,
		})
		require.NoError(t, err)

		result, err = e.AdvanceGameTime(ctx, AdvanceGameTimeRequest{
			GameID: gameResult.Game.ID,
			Days:   1,
		})
		require.NoError(t, err)

		assert.Greater(t, wealthyCost, result.TotalCost)
	})

	t.Run("advances time by zero days", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		_, err = e.SetLifestyle(ctx, SetLifestyleRequest{
			GameID: gameResult.Game.ID,
			Tier:   model.LifestyleComfortable,
		})
		require.NoError(t, err)

		result, err := e.AdvanceGameTime(ctx, AdvanceGameTimeRequest{
			GameID: gameResult.Game.ID,
			Days:   0,
		})

		require.NoError(t, err)
		assert.Equal(t, 0, result.DaysAdvanced)
		assert.Equal(t, 0, result.TotalCost)
	})
}
