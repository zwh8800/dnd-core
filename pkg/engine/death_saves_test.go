package engine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zwh8800/dnd-core/pkg/model"
)

func TestPerformDeathSave(t *testing.T) {
	t.Run("perform death save successfully", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		// Create a game
		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test death saves",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create a PC
		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:      "Test Character",
				Race:      "Human",
				Class:     "战士",
				Level:     1,
				HitPoints: 10,
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

		// Reduce HP to 0
		_, err = e.ExecuteDamage(ctx, ExecuteDamageRequest{
			GameID:   gameID,
			TargetID: pcResult.Actor.ID,
			Damage: DamageInput{
				Amount: 10,
				Type:   model.DamageTypeSlashing,
				Source: pcResult.Actor.ID,
			},
		})
		require.NoError(t, err)

		// Perform death save
		result, err := e.PerformDeathSave(ctx, PerformDeathSaveRequest{
			GameID:  gameID,
			ActorID: pcResult.Actor.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, result.Roll, 1)
		assert.LessOrEqual(t, result.Roll, 20)
		assert.Contains(t, result.Message, "死亡豁免掷骰")
	})

	t.Run("perform death save with HP > 0 fails", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test death saves",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:      "Test Character",
				Race:      "Human",
				Class:     "战士",
				Level:     1,
				HitPoints: 10,
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

		// Try to perform death save with full HP - should fail
		_, err = e.PerformDeathSave(ctx, PerformDeathSaveRequest{
			GameID:  gameID,
			ActorID: pcResult.Actor.ID,
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "HP大于0")
	})

	t.Run("perform death save for non-PC fails", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test death saves",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create an enemy
		enemyResult, err := e.CreateEnemy(ctx, CreateEnemyRequest{
			GameID: gameID,
			Enemy: &EnemyInput{
				Name:        "Goblin",
				Description: "A goblin",
				Size:        model.SizeSmall,
				Speed:       30,
				HitPoints:   7,
				ArmorClass:  15,
				AbilityScores: AbilityScoresInput{
					Strength:     8,
					Dexterity:    14,
					Constitution: 10,
					Intelligence: 10,
					Wisdom:       8,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)

		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		// Reduce enemy HP to 0
		_, err = e.ExecuteDamage(ctx, ExecuteDamageRequest{
			GameID:   gameID,
			TargetID: enemyResult.Actor.ID,
			Damage: DamageInput{
				Amount: 7,
				Type:   model.DamageTypeSlashing,
				Source: enemyResult.Actor.ID,
			},
		})
		require.NoError(t, err)

		// Try to perform death save for enemy - should fail
		_, err = e.PerformDeathSave(ctx, PerformDeathSaveRequest{
			GameID:  gameID,
			ActorID: enemyResult.Actor.ID,
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "only player characters")
	})

	t.Run("perform death save for stabilized creature fails", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test death saves",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:      "Test Character",
				Race:      "Human",
				Class:     "战士",
				Level:     1,
				HitPoints: 10,
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

		// Reduce HP to 0
		_, err = e.ExecuteDamage(ctx, ExecuteDamageRequest{
			GameID:   gameID,
			TargetID: pcResult.Actor.ID,
			Damage: DamageInput{
				Amount: 10,
				Type:   model.DamageTypeSlashing,
				Source: pcResult.Actor.ID,
			},
		})
		require.NoError(t, err)

		// Stabilize the creature
		_, err = e.StabilizeCreature(ctx, StabilizeCreatureRequest{
			GameID:  gameID,
			ActorID: pcResult.Actor.ID,
		})
		require.NoError(t, err)

		// Try to perform death save for stabilized creature - should fail
		_, err = e.PerformDeathSave(ctx, PerformDeathSaveRequest{
			GameID:  gameID,
			ActorID: pcResult.Actor.ID,
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "已稳定")
	})
}

func TestStabilizeCreature(t *testing.T) {
	t.Run("stabilize creature successfully", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test death saves",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:      "Test Character",
				Race:      "Human",
				Class:     "战士",
				Level:     1,
				HitPoints: 10,
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

		// Reduce HP to 0
		_, err = e.ExecuteDamage(ctx, ExecuteDamageRequest{
			GameID:   gameID,
			TargetID: pcResult.Actor.ID,
			Damage: DamageInput{
				Amount: 10,
				Type:   model.DamageTypeSlashing,
				Source: pcResult.Actor.ID,
			},
		})
		require.NoError(t, err)

		// Stabilize the creature
		result, err := e.StabilizeCreature(ctx, StabilizeCreatureRequest{
			GameID:  gameID,
			ActorID: pcResult.Actor.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, pcResult.Actor.ID, result.ActorID)
		assert.Contains(t, result.Message, "稳定")
	})

	t.Run("stabilize creature with HP > 0 fails", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test death saves",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:      "Test Character",
				Race:      "Human",
				Class:     "战士",
				Level:     1,
				HitPoints: 10,
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

		// Try to stabilize with full HP - should fail
		_, err = e.StabilizeCreature(ctx, StabilizeCreatureRequest{
			GameID:  gameID,
			ActorID: pcResult.Actor.ID,
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "HP大于0")
	})

	t.Run("stabilize already stabilized creature fails", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test death saves",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:      "Test Character",
				Race:      "Human",
				Class:     "战士",
				Level:     1,
				HitPoints: 10,
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

		// Reduce HP to 0
		_, err = e.ExecuteDamage(ctx, ExecuteDamageRequest{
			GameID:   gameID,
			TargetID: pcResult.Actor.ID,
			Damage: DamageInput{
				Amount: 10,
				Type:   model.DamageTypeSlashing,
				Source: pcResult.Actor.ID,
			},
		})
		require.NoError(t, err)

		// Stabilize first time
		_, err = e.StabilizeCreature(ctx, StabilizeCreatureRequest{
			GameID:  gameID,
			ActorID: pcResult.Actor.ID,
		})
		require.NoError(t, err)

		// Try to stabilize again - should fail
		_, err = e.StabilizeCreature(ctx, StabilizeCreatureRequest{
			GameID:  gameID,
			ActorID: pcResult.Actor.ID,
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "已稳定")
	})
}

func TestGetDeathSaveStatus(t *testing.T) {
	t.Run("get death save status for healthy character", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test death saves",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:      "Test Character",
				Race:      "Human",
				Class:     "战士",
				Level:     1,
				HitPoints: 10,
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

		result, err := e.GetDeathSaveStatus(ctx, GetDeathSaveStatusRequest{
			GameID:  gameID,
			ActorID: pcResult.Actor.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, pcResult.Actor.ID, result.ActorID)
		assert.False(t, result.IsUnconscious)
		assert.False(t, result.IsStable)
		assert.False(t, result.IsDead)
		assert.Equal(t, 10, result.CurrentHP)
	})

	t.Run("get death save status for unconscious character", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test death saves",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:      "Test Character",
				Race:      "Human",
				Class:     "战士",
				Level:     1,
				HitPoints: 10,
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

		// Reduce HP to 0
		_, err = e.ExecuteDamage(ctx, ExecuteDamageRequest{
			GameID:   gameID,
			TargetID: pcResult.Actor.ID,
			Damage: DamageInput{
				Amount: 10,
				Type:   model.DamageTypeSlashing,
				Source: pcResult.Actor.ID,
			},
		})
		require.NoError(t, err)

		result, err := e.GetDeathSaveStatus(ctx, GetDeathSaveStatusRequest{
			GameID:  gameID,
			ActorID: pcResult.Actor.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.IsUnconscious)
		assert.False(t, result.IsStable)
		// 注意: Actor.IsDead() 返回 HP<=0 && !IsStabilized，所以濒死状态也被视为"死亡"
		assert.True(t, result.IsDead)
		assert.Equal(t, 0, result.CurrentHP)
	})

	t.Run("get death save status for stabilized character", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test death saves",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:      "Test Character",
				Race:      "Human",
				Class:     "战士",
				Level:     1,
				HitPoints: 10,
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

		// Reduce HP to 0
		_, err = e.ExecuteDamage(ctx, ExecuteDamageRequest{
			GameID:   gameID,
			TargetID: pcResult.Actor.ID,
			Damage: DamageInput{
				Amount: 10,
				Type:   model.DamageTypeSlashing,
				Source: pcResult.Actor.ID,
			},
		})
		require.NoError(t, err)

		// Stabilize
		_, err = e.StabilizeCreature(ctx, StabilizeCreatureRequest{
			GameID:  gameID,
			ActorID: pcResult.Actor.ID,
		})
		require.NoError(t, err)

		result, err := e.GetDeathSaveStatus(ctx, GetDeathSaveStatusRequest{
			GameID:  gameID,
			ActorID: pcResult.Actor.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.IsUnconscious)
		assert.True(t, result.IsStable)
		assert.False(t, result.IsDead)
		assert.Equal(t, 0, result.CurrentHP)
	})

	t.Run("get death save status for non-existent actor fails", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test death saves",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		_, err = e.GetDeathSaveStatus(ctx, GetDeathSaveStatusRequest{
			GameID:  gameID,
			ActorID: "non-existent",
		})

		assert.Error(t, err)
	})

	t.Run("get death save status for non-PC fails", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test death saves",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		enemyResult, err := e.CreateEnemy(ctx, CreateEnemyRequest{
			GameID: gameID,
			Enemy: &EnemyInput{
				Name:        "Goblin",
				Description: "A goblin",
				Size:        model.SizeSmall,
				Speed:       30,
				HitPoints:   7,
				ArmorClass:  15,
				AbilityScores: AbilityScoresInput{
					Strength:     8,
					Dexterity:    14,
					Constitution: 10,
					Intelligence: 10,
					Wisdom:       8,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)

		_, err = e.GetDeathSaveStatus(ctx, GetDeathSaveStatusRequest{
			GameID:  gameID,
			ActorID: enemyResult.Actor.ID,
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "only player characters")
	})
}
