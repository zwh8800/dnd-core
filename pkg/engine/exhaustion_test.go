package engine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zwh8800/dnd-core/pkg/model"
)

func TestApplyExhaustion(t *testing.T) {
	t.Run("apply exhaustion successfully", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		// Create a game
		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test exhaustion",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create a PC
		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     15,
					Dexterity:    14,
					Constitution: 13,
					Intelligence: 12,
					Wisdom:       10,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)

		// Switch to exploration phase
		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		// Apply exhaustion level 1
		result, err := e.ApplyExhaustion(ctx, ApplyExhaustionRequest{
			GameID:  gameID,
			ActorID: pcResult.Actor.ID,
			Levels:  1,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 1, result.NewLevel)
		assert.False(t, result.IsDead)
		assert.Contains(t, result.Message, "力竭等级提升至 1")
		assert.Len(t, result.Effects, 1)
	})

	t.Run("apply multiple exhaustion levels", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test exhaustion",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     15,
					Dexterity:    14,
					Constitution: 13,
					Intelligence: 12,
					Wisdom:       10,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)

		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		// Apply 3 levels at once
		result, err := e.ApplyExhaustion(ctx, ApplyExhaustionRequest{
			GameID:  gameID,
			ActorID: pcResult.Actor.ID,
			Levels:  3,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 3, result.NewLevel)
		assert.False(t, result.IsDead)
		assert.Len(t, result.Effects, 3) // Should have 3 cumulative effects
	})

	t.Run("apply exhaustion reaches level 6 causes death", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test exhaustion",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     15,
					Dexterity:    14,
					Constitution: 13,
					Intelligence: 12,
					Wisdom:       10,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)

		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		// Apply 6 levels - should cause death
		result, err := e.ApplyExhaustion(ctx, ApplyExhaustionRequest{
			GameID:  gameID,
			ActorID: pcResult.Actor.ID,
			Levels:  6,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 6, result.NewLevel)
		assert.True(t, result.IsDead)
		assert.Contains(t, result.Message, "角色死亡")
	})

	t.Run("apply exhaustion with invalid levels fails", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test exhaustion",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     15,
					Dexterity:    14,
					Constitution: 13,
					Intelligence: 12,
					Wisdom:       10,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)

		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		// Try to apply 0 levels - should fail
		_, err = e.ApplyExhaustion(ctx, ApplyExhaustionRequest{
			GameID:  gameID,
			ActorID: pcResult.Actor.ID,
			Levels:  0,
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "力竭等级必须大于0")
	})

	t.Run("apply exhaustion with non-existent actor fails", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test exhaustion",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		_, err = e.ApplyExhaustion(ctx, ApplyExhaustionRequest{
			GameID:  gameID,
			ActorID: "non-existent",
			Levels:  1,
		})

		assert.Error(t, err)
	})

	t.Run("apply exhaustion in wrong phase fails", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test exhaustion",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     15,
					Dexterity:    14,
					Constitution: 13,
					Intelligence: 12,
					Wisdom:       10,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)

		// Stay in character creation phase
		_, err = e.ApplyExhaustion(ctx, ApplyExhaustionRequest{
			GameID:  gameID,
			ActorID: pcResult.Actor.ID,
			Levels:  1,
		})

		assert.Error(t, err) // Should fail due to phase restriction
	})
}

func TestRemoveExhaustion(t *testing.T) {
	t.Run("remove exhaustion successfully", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test exhaustion",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     15,
					Dexterity:    14,
					Constitution: 13,
					Intelligence: 12,
					Wisdom:       10,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)

		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		// First apply 2 levels
		_, err = e.ApplyExhaustion(ctx, ApplyExhaustionRequest{
			GameID:  gameID,
			ActorID: pcResult.Actor.ID,
			Levels:  2,
		})
		require.NoError(t, err)

		// Remove 1 level
		result, err := e.RemoveExhaustion(ctx, RemoveExhaustionRequest{
			GameID:  gameID,
			ActorID: pcResult.Actor.ID,
			Levels:  1,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 1, result.NewLevel)
		assert.Contains(t, result.Message, "力竭等级降低至 1")
	})

	t.Run("remove exhaustion completely recovers", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test exhaustion",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     15,
					Dexterity:    14,
					Constitution: 13,
					Intelligence: 12,
					Wisdom:       10,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)

		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		// Apply 1 level
		_, err = e.ApplyExhaustion(ctx, ApplyExhaustionRequest{
			GameID:  gameID,
			ActorID: pcResult.Actor.ID,
			Levels:  1,
		})
		require.NoError(t, err)

		// Remove 1 level - should completely recover
		result, err := e.RemoveExhaustion(ctx, RemoveExhaustionRequest{
			GameID:  gameID,
			ActorID: pcResult.Actor.ID,
			Levels:  1,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 0, result.NewLevel)
		assert.Contains(t, result.Message, "已完全恢复")
	})

	t.Run("remove exhaustion with no exhaustion fails", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test exhaustion",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     15,
					Dexterity:    14,
					Constitution: 13,
					Intelligence: 12,
					Wisdom:       10,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)

		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		// Try to remove when no exhaustion - should fail
		_, err = e.RemoveExhaustion(ctx, RemoveExhaustionRequest{
			GameID:  gameID,
			ActorID: pcResult.Actor.ID,
			Levels:  1,
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "没有力竭等级")
	})

	t.Run("remove exhaustion with invalid levels fails", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test exhaustion",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     15,
					Dexterity:    14,
					Constitution: 13,
					Intelligence: 12,
					Wisdom:       10,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)

		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		// Try to remove 0 levels - should fail
		_, err = e.RemoveExhaustion(ctx, RemoveExhaustionRequest{
			GameID:  gameID,
			ActorID: pcResult.Actor.ID,
			Levels:  0,
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "移除的力竭等级必须大于0")
	})
}

func TestGetExhaustionStatus(t *testing.T) {
	t.Run("get exhaustion status with no exhaustion", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test exhaustion",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     15,
					Dexterity:    14,
					Constitution: 13,
					Intelligence: 12,
					Wisdom:       10,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)

		result, err := e.GetExhaustionStatus(ctx, GetExhaustionStatusRequest{
			GameID:  gameID,
			ActorID: pcResult.Actor.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 0, result.CurrentLevel)
		assert.False(t, result.IsDead)
		assert.Empty(t, result.Effects)
		assert.Contains(t, result.Message, "力竭等级: 0")
	})

	t.Run("get exhaustion status with exhaustion", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test exhaustion",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     15,
					Dexterity:    14,
					Constitution: 13,
					Intelligence: 12,
					Wisdom:       10,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)

		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		// Apply 3 levels
		_, err = e.ApplyExhaustion(ctx, ApplyExhaustionRequest{
			GameID:  gameID,
			ActorID: pcResult.Actor.ID,
			Levels:  3,
		})
		require.NoError(t, err)

		result, err := e.GetExhaustionStatus(ctx, GetExhaustionStatusRequest{
			GameID:  gameID,
			ActorID: pcResult.Actor.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 3, result.CurrentLevel)
		assert.False(t, result.IsDead)
		assert.Len(t, result.Effects, 3)
		assert.Contains(t, result.Effects[0], "力竭1级")
		assert.Contains(t, result.Effects[1], "力竭2级")
		assert.Contains(t, result.Effects[2], "力竭3级")
	})

	t.Run("get exhaustion status with non-existent actor fails", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test exhaustion",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		_, err = e.GetExhaustionStatus(ctx, GetExhaustionStatusRequest{
			GameID:  gameID,
			ActorID: "non-existent",
		})

		assert.Error(t, err)
	})
}
