package engine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zwh8800/dnd-core/internal/model"
)

func TestStartCombat(t *testing.T) {
	t.Run("starts combat successfully", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for combat",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create combatants
		pc1, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Fighter",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     16,
					Dexterity:    14,
					Constitution: 15,
					Intelligence: 10,
					Wisdom:       12,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)

		pc2, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Wizard",
				Race:  "Elf",
				Class: "Wizard",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     8,
					Dexterity:    14,
					Constitution: 12,
					Intelligence: 16,
					Wisdom:       10,
					Charisma:     12,
				},
			},
		})
		require.NoError(t, err)

		enemy, err := e.CreateEnemy(ctx, CreateEnemyRequest{
			GameID: gameID,
			Enemy: &EnemyInput{
				Name:        "Goblin",
				Description: "A sneaky goblin",
				Size:        model.SizeSmall,
				Speed:       30,
				AbilityScores: AbilityScoresInput{
					Strength:     8,
					Dexterity:    14,
					Constitution: 10,
					Intelligence: 10,
					Wisdom:       8,
					Charisma:     8,
				},
				ChallengeRating: 0.25,
				HitPoints:       7,
				ArmorClass:      15,
			},
		})
		require.NoError(t, err)

		// Switch to exploration phase
		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "test")
		require.NoError(t, err)

		result, err := e.StartCombat(ctx, StartCombatRequest{
			GameID:         gameID,
			ParticipantIDs: []model.ID{pc1.Actor.ID, pc2.Actor.ID, enemy.Actor.ID},
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Combat)
	})

	t.Run("starts combat with surprise", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for combat with surprise",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create combatants
		pc, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Fighter",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     16,
					Dexterity:    14,
					Constitution: 15,
					Intelligence: 10,
					Wisdom:       12,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)

		enemy, err := e.CreateEnemy(ctx, CreateEnemyRequest{
			GameID: gameID,
			Enemy: &EnemyInput{
				Name:        "Assassin",
				Description: "A deadly assassin",
				Size:        model.SizeMedium,
				Speed:       30,
				AbilityScores: AbilityScoresInput{
					Strength:     12,
					Dexterity:    16,
					Constitution: 12,
					Intelligence: 10,
					Wisdom:       10,
					Charisma:     8,
				},
				ChallengeRating: 0.5,
				HitPoints:       12,
				ArmorClass:      14,
			},
		})
		require.NoError(t, err)

		// Switch to exploration phase
		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "test")
		require.NoError(t, err)

		result, err := e.StartCombatWithSurprise(ctx, StartCombatWithSurpriseRequest{
			GameID:       gameID,
			StealthySide: []model.ID{pc.Actor.ID},
			Observers:    []model.ID{enemy.Actor.ID},
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Combat)
	})
}

func TestEndCombat(t *testing.T) {
	t.Run("ends combat successfully", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for combat",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create combatants and start combat
		pc, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Fighter",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     16,
					Dexterity:    14,
					Constitution: 15,
					Intelligence: 10,
					Wisdom:       12,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)

		enemy, err := e.CreateEnemy(ctx, CreateEnemyRequest{
			GameID: gameID,
			Enemy: &EnemyInput{
				Name:        "Goblin",
				Description: "A sneaky goblin",
				Size:        model.SizeSmall,
				Speed:       30,
				AbilityScores: AbilityScoresInput{
					Strength:     8,
					Dexterity:    14,
					Constitution: 10,
					Intelligence: 10,
					Wisdom:       8,
					Charisma:     8,
				},
				ChallengeRating: 0.25,
				HitPoints:       7,
				ArmorClass:      15,
			},
		})
		require.NoError(t, err)

		// Switch to exploration phase
		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "test")
		require.NoError(t, err)

		_, err = e.StartCombat(ctx, StartCombatRequest{
			GameID:         gameID,
			ParticipantIDs: []model.ID{pc.Actor.ID, enemy.Actor.ID},
		})
		require.NoError(t, err)

		// End combat
		err = e.EndCombat(ctx, EndCombatRequest{GameID: gameID})

		require.NoError(t, err)
	})
}

func TestNextTurn(t *testing.T) {
	t.Run("advances to next turn", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for combat",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create combatants and start combat
		pc1, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Fighter",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     16,
					Dexterity:    14,
					Constitution: 15,
					Intelligence: 10,
					Wisdom:       12,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)

		pc2, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Wizard",
				Race:  "Elf",
				Class: "Wizard",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     8,
					Dexterity:    14,
					Constitution: 12,
					Intelligence: 16,
					Wisdom:       10,
					Charisma:     12,
				},
			},
		})
		require.NoError(t, err)

		// Switch to exploration phase
		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "test")
		require.NoError(t, err)

		_, err = e.StartCombat(ctx, StartCombatRequest{
			GameID:         gameID,
			ParticipantIDs: []model.ID{pc1.Actor.ID, pc2.Actor.ID},
		})
		require.NoError(t, err)

		// Next turn
		result, err := e.NextTurn(ctx, NextTurnRequest{GameID: gameID})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Combat)
	})
}

func TestGetCurrentCombat(t *testing.T) {
	t.Run("gets current combat info", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for combat",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create combatants and start combat
		pc, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Fighter",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     16,
					Dexterity:    14,
					Constitution: 15,
					Intelligence: 10,
					Wisdom:       12,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)

		// Switch to exploration phase
		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "test")
		require.NoError(t, err)

		_, err = e.StartCombat(ctx, StartCombatRequest{
			GameID:         gameID,
			ParticipantIDs: []model.ID{pc.Actor.ID},
		})
		require.NoError(t, err)

		// Get current combat
		result, err := e.GetCurrentCombat(ctx, GetCurrentCombatRequest{GameID: gameID})

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.Combat)
		assert.Equal(t, 1, result.Combat.Round)
	})
}

func TestExecuteAttack(t *testing.T) {
	t.Run("executes attack", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for combat",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create combatants and start combat
		attacker, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Fighter",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     16,
					Dexterity:    14,
					Constitution: 15,
					Intelligence: 10,
					Wisdom:       12,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)

		target, err := e.CreateEnemy(ctx, CreateEnemyRequest{
			GameID: gameID,
			Enemy: &EnemyInput{
				Name:        "Goblin",
				Description: "A sneaky goblin",
				Size:        model.SizeSmall,
				Speed:       30,
				AbilityScores: AbilityScoresInput{
					Strength:     8,
					Dexterity:    14,
					Constitution: 10,
					Intelligence: 10,
					Wisdom:       8,
					Charisma:     8,
				},
				ChallengeRating: 0.25,
				HitPoints:       7,
				ArmorClass:      15,
			},
		})
		require.NoError(t, err)

		// Switch to exploration phase
		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "test")
		require.NoError(t, err)

		_, err = e.StartCombat(ctx, StartCombatRequest{
			GameID:         gameID,
			ParticipantIDs: []model.ID{attacker.Actor.ID, target.Actor.ID},
		})
		require.NoError(t, err)

		// Execute attack
		result, err := e.ExecuteAttack(ctx, ExecuteAttackRequest{
			GameID:     gameID,
			AttackerID: attacker.Actor.ID,
			TargetID:   target.Actor.ID,
			Attack:     AttackInput{IsUnarmed: true},
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.AttackResult)
	})
}

func TestExecuteDamage(t *testing.T) {
	t.Run("executes damage", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for combat",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create target
		target, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Target",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     16,
					Dexterity:    14,
					Constitution: 15,
					Intelligence: 10,
					Wisdom:       12,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)

		// Execute damage
		result, err := e.ExecuteDamage(ctx, ExecuteDamageRequest{
			GameID:   gameID,
			TargetID: target.Actor.ID,
			Damage: DamageInput{
				Amount: 5,
				Type:   model.DamageTypeSlashing,
			},
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 5, result.DamageResult.FinalDamage)
	})
}

func TestExecuteHealing(t *testing.T) {
	t.Run("executes healing", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for healing",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create wounded target
		target, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Wounded Target",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     16,
					Dexterity:    14,
					Constitution: 15,
					Intelligence: 10,
					Wisdom:       12,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)

		// Damage the target
		_, err = e.ExecuteDamage(ctx, ExecuteDamageRequest{
			GameID:   gameID,
			TargetID: target.Actor.ID,
			Damage: DamageInput{
				Amount: 10,
				Type:   model.DamageTypeSlashing,
			},
		})
		require.NoError(t, err)

		// Heal the target
		result, err := e.ExecuteHealing(ctx, ExecuteHealingRequest{
			GameID:   gameID,
			TargetID: target.Actor.ID,
			Amount:   5,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 5, result.Healed)
	})
}
