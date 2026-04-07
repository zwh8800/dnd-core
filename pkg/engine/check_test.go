package engine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zwh8800/dnd-core/internal/model"
)

func TestAbilityCheck(t *testing.T) {
	t.Run("performs ability check", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for ability checks",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create an actor
		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "Rogue",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     12,
					Dexterity:    16,
					Constitution: 14,
					Intelligence: 10,
					Wisdom:       12,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		result, err := e.PerformAbilityCheck(ctx, AbilityCheckRequest{
			GameID:  gameID,
			ActorID: actorID,
			Ability: model.AbilityStrength,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Roll)
		assert.Equal(t, model.AbilityStrength, result.Ability)
	})

	t.Run("ability check with DC", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for ability checks",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create an actor
		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "Rogue",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     12,
					Dexterity:    16,
					Constitution: 14,
					Intelligence: 10,
					Wisdom:       12,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		result, err := e.PerformAbilityCheck(ctx, AbilityCheckRequest{
			GameID:  gameID,
			ActorID: actorID,
			Ability: model.AbilityDexterity,
			DC:      15,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Success)
		assert.NotNil(t, result.Roll)
	})
}

func TestSkillCheck(t *testing.T) {
	t.Run("performs skill check", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for skill checks",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create an actor
		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "Rogue",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     12,
					Dexterity:    16,
					Constitution: 14,
					Intelligence: 10,
					Wisdom:       12,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		result, err := e.PerformSkillCheck(ctx, SkillCheckRequest{
			GameID:  gameID,
			ActorID: actorID,
			Skill:   model.SkillStealth,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Roll)
		assert.Equal(t, model.SkillStealth, result.Skill)
	})

	t.Run("skill check with advantage", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for skill checks",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create an actor
		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "Rogue",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     12,
					Dexterity:    16,
					Constitution: 14,
					Intelligence: 10,
					Wisdom:       12,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		result, err := e.PerformSkillCheck(ctx, SkillCheckRequest{
			GameID:    gameID,
			ActorID:   actorID,
			Skill:     model.SkillStealth,
			Advantage: model.RollModifier{Advantage: true},
		})

		require.NoError(t, err)
		require.NotNil(t, result)
	})
}

func TestSavingThrow(t *testing.T) {
	t.Run("performs saving throw", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for saving throws",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create an actor
		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     16,
					Dexterity:    12,
					Constitution: 15,
					Intelligence: 10,
					Wisdom:       12,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		result, err := e.PerformSavingThrow(ctx, SavingThrowRequest{
			GameID:  gameID,
			ActorID: actorID,
			Ability: model.AbilityConstitution,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Roll)
		assert.Equal(t, model.AbilityConstitution, result.Ability)
	})

	t.Run("saving throw with DC", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for saving throws",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create an actor
		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     16,
					Dexterity:    12,
					Constitution: 15,
					Intelligence: 10,
					Wisdom:       12,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		result, err := e.PerformSavingThrow(ctx, SavingThrowRequest{
			GameID:  gameID,
			ActorID: actorID,
			Ability: model.AbilityStrength,
			DC:      14,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Success)
	})
}

func TestGetPassivePerception(t *testing.T) {
	t.Run("gets passive perception", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for passive perception",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create an actor
		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "Rogue",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     12,
					Dexterity:    16,
					Constitution: 14,
					Intelligence: 10,
					Wisdom:       14,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		result, err := e.GetPassivePerception(ctx, GetPassivePerceptionRequest{
			GameID:  gameID,
			ActorID: actorID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		// Passive perception = 10 + WIS modifier + proficiency (if proficient in Perception)
		// WIS 14 = +2 modifier, so base is 12
		assert.GreaterOrEqual(t, result.PassivePerception, 10)
	})
}
