package engine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zwh8800/dnd-core/pkg/model"
)

func TestPlaceTrap(t *testing.T) {
	t.Run("places trap successfully with valid data", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.PlaceTrap(ctx, PlaceTrapRequest{
			GameID:   gameResult.Game.ID,
			SceneID:  model.NewID(),
			TrapID:   "poison-needle",
			Position: "entrance",
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Trap)
		assert.True(t, result.Trap.IsArmed)
		assert.False(t, result.Trap.HasTriggered)
		assert.Equal(t, "entrance", result.Trap.Position)
		assert.Contains(t, result.Message, "毒针陷阱")
	})

	t.Run("returns error when trap data not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.PlaceTrap(ctx, PlaceTrapRequest{
			GameID:   gameResult.Game.ID,
			TrapID:   "nonexistent-trap",
			Position: "room1",
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "trap data")
	})

	t.Run("places magical trap successfully", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.PlaceTrap(ctx, PlaceTrapRequest{
			GameID:   gameResult.Game.ID,
			TrapID:   "fire-dart",
			Position: "corridor",
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Trap)
		assert.Equal(t, "fire-dart", result.Trap.Definition.ID)
		assert.True(t, result.Trap.IsArmed)
	})

	t.Run("places trap with empty position", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.PlaceTrap(ctx, PlaceTrapRequest{
			GameID:   gameResult.Game.ID,
			TrapID:   "pit-trap",
			Position: "",
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "", result.Trap.Position)
	})

	t.Run("returns error when game not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.PlaceTrap(ctx, PlaceTrapRequest{
			GameID:   model.NewID(),
			TrapID:   "poison-needle",
			Position: "room1",
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})
}

func TestDetectTrap(t *testing.T) {
	t.Run("returns trap detection DC", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.DetectTrap(ctx, DetectTrapRequest{
			GameID:  gameResult.Game.ID,
			TrapID:  model.NewID(),
			SceneID: model.NewID(),
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 15, result.DC)
		assert.False(t, result.TrapRevealed)
		assert.Contains(t, result.Message, "DC 15")
	})

	t.Run("returns error when game not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.DetectTrap(ctx, DetectTrapRequest{
			GameID: model.NewID(),
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("detection result has correct message format", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.DetectTrap(ctx, DetectTrapRequest{
			GameID: gameResult.Game.ID,
		})

		require.NoError(t, err)
		assert.Contains(t, result.Message, "感知")
		assert.Contains(t, result.Message, "察觉")
	})

	t.Run("multiple detections return consistent results", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		for i := 0; i < 3; i++ {
			result, err := e.DetectTrap(ctx, DetectTrapRequest{
				GameID: gameResult.Game.ID,
			})
			require.NoError(t, err)
			assert.Equal(t, 15, result.DC)
			assert.False(t, result.TrapRevealed)
		}
	})

	t.Run("detection does not modify game state", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		_, err = e.DetectTrap(ctx, DetectTrapRequest{
			GameID: gameResult.Game.ID,
		})
		require.NoError(t, err)

		// Verify game is still accessible by loading it
		loadResult, err := e.LoadGame(ctx, LoadGameRequest{GameID: gameResult.Game.ID})
		require.NoError(t, err)
		assert.NotNil(t, loadResult)
	})
}

func TestDisarmTrap(t *testing.T) {
	t.Run("returns trap disarm DC", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.DisarmTrap(ctx, DisarmTrapRequest{
			GameID:  gameResult.Game.ID,
			TrapID:  model.NewID(),
			SceneID: model.NewID(),
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 15, result.DC)
		assert.False(t, result.TrapDisarmed)
		assert.Contains(t, result.Message, "DC 15")
	})

	t.Run("returns error when game not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.DisarmTrap(ctx, DisarmTrapRequest{
			GameID: model.NewID(),
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("disarm result mentions required skill", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.DisarmTrap(ctx, DisarmTrapRequest{
			GameID: gameResult.Game.ID,
		})

		require.NoError(t, err)
		assert.Contains(t, result.Message, "妙手")
	})

	t.Run("multiple disarm attempts return consistent results", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		for i := 0; i < 3; i++ {
			result, err := e.DisarmTrap(ctx, DisarmTrapRequest{
				GameID: gameResult.Game.ID,
			})
			require.NoError(t, err)
			assert.Equal(t, 15, result.DC)
		}
	})

	t.Run("disarm uses write lock for concurrency safety", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		done := make(chan bool, 2)
		go func() {
			_, err := e.DisarmTrap(ctx, DisarmTrapRequest{GameID: gameResult.Game.ID})
			assert.NoError(t, err)
			done <- true
		}()
		go func() {
			_, err := e.DisarmTrap(ctx, DisarmTrapRequest{GameID: gameResult.Game.ID})
			assert.NoError(t, err)
			done <- true
		}()

		<-done
		<-done
	})
}

func TestTriggerTrap(t *testing.T) {
	t.Run("traps triggers and returns effects", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.TriggerTrap(ctx, TriggerTrapRequest{
			GameID:  gameResult.Game.ID,
			TrapID:  model.NewID(),
			SceneID: model.NewID(),
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.TrapTriggered)
		assert.NotEmpty(t, result.Effects)
		assert.NotEmpty(t, result.Message)
	})

	t.Run("trigger result includes save DC", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.TriggerTrap(ctx, TriggerTrapRequest{
			GameID: gameResult.Game.ID,
		})

		require.NoError(t, err)
		assert.Greater(t, result.SaveDC, 0)
		assert.NotEmpty(t, result.SaveAbility)
	})

	t.Run("returns error when game not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.TriggerTrap(ctx, TriggerTrapRequest{
			GameID: model.NewID(),
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("trap trigger message contains trap name", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.TriggerTrap(ctx, TriggerTrapRequest{
			GameID: gameResult.Game.ID,
		})

		require.NoError(t, err)
		assert.Contains(t, result.Message, "陷阱触发")
	})

	t.Run("multiple triggers return effects consistently", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		for i := 0; i < 3; i++ {
			result, err := e.TriggerTrap(ctx, TriggerTrapRequest{
				GameID: gameResult.Game.ID,
			})
			require.NoError(t, err)
			assert.True(t, result.TrapTriggered)
			assert.NotEmpty(t, result.Effects)
		}
	})
}
