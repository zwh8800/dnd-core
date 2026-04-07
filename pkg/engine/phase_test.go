package engine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zwh8800/dnd-core/internal/model"
)

func TestSetPhase(t *testing.T) {
	t.Run("sets phase", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for phases",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		result, err := e.SetPhase(ctx, gameID, model.PhaseExploration, "Starting exploration")

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, model.PhaseExploration, result.NewPhase)
	})

	t.Run("transitions to combat phase", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for phases",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		_, err = e.SetPhase(ctx, gameID, model.PhaseCombat, "Combat starting")
		require.NoError(t, err)

		result, err := e.SetPhase(ctx, gameID, model.PhaseExploration, "Combat ended")
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, model.PhaseExploration, result.NewPhase)
	})
}

func TestGetPhase(t *testing.T) {
	t.Run("gets current phase", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for phases",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		phase, err := e.GetPhase(ctx, gameID)

		require.NoError(t, err)
		assert.NotEmpty(t, phase)
	})
}

func TestGetAllowedOperations(t *testing.T) {
	t.Run("gets allowed operations", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for phases",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		ops, err := e.GetAllowedOperations(ctx, gameID)

		require.NoError(t, err)
		assert.NotEmpty(t, ops)
	})
}
