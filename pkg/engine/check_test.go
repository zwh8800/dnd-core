package engine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zwh8800/dnd-core/pkg/model"
)

func TestPerformAbilityCheck(t *testing.T) {
	t.Run("performs ability check successfully", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for ability checks",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "游荡者",
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
		assert.Equal(t, "Test Character", result.ActorName)
	})

	t.Run("ability check with DC succeeds", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for ability checks",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "战士",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     20,
					Dexterity:    12,
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
			DC:      10,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.Success)
		assert.GreaterOrEqual(t, result.Margin, 0)
	})

	t.Run("ability check with advantage", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for ability checks",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "游荡者",
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
			GameID:    gameID,
			ActorID:   actorID,
			Ability:   model.AbilityDexterity,
			Advantage: model.RollModifier{Advantage: true},
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Roll)
	})

	t.Run("ability check with disadvantage", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for ability checks",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "游荡者",
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
			GameID:    gameID,
			ActorID:   actorID,
			Ability:   model.AbilityDexterity,
			Advantage: model.RollModifier{Disadvantage: true},
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Roll)
	})

	t.Run("ability check fails with invalid actor", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for ability checks",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		_, err = e.PerformAbilityCheck(ctx, AbilityCheckRequest{
			GameID:  gameID,
			ActorID: model.ID("invalid-actor-id"),
			Ability: model.AbilityStrength,
		})

		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("ability check fails with invalid game", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		_, err := e.PerformAbilityCheck(ctx, AbilityCheckRequest{
			GameID:  model.ID("invalid-game-id"),
			ActorID: model.ID("actor-id"),
			Ability: model.AbilityStrength,
		})

		assert.Error(t, err)
	})

	t.Run("ability check with different actor types", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for ability checks",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreateEnemy(ctx, CreateEnemyRequest{
			GameID: gameID,
			Enemy: &EnemyInput{
				Name:            "Goblin",
				Description:     "A small goblin",
				AbilityScores:   AbilityScoresInput{Strength: 10, Dexterity: 12, Constitution: 10, Intelligence: 8, Wisdom: 10, Charisma: 8},
				ChallengeRating: "1/4",
				HitPoints:       7,
				ArmorClass:      15,
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		result, err := e.PerformAbilityCheck(ctx, AbilityCheckRequest{
			GameID:  gameID,
			ActorID: actorID,
			Ability: model.AbilityDexterity,
			DC:      12,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, model.AbilityDexterity, result.Ability)
		assert.True(t, result.RollTotal > 0)
	})
}

func TestPerformSkillCheck(t *testing.T) {
	t.Run("performs skill check successfully", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for skill checks",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "游荡者",
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
		assert.Equal(t, "Test Character", result.ActorName)
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

		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "游荡者",
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
		assert.NotNil(t, result.Roll)
	})

	t.Run("skill check with disadvantage", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for skill checks",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "游荡者",
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
			Advantage: model.RollModifier{Disadvantage: true},
		})

		require.NoError(t, err)
		require.NotNil(t, result)
	})

	t.Run("skill check with DC", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for skill checks",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "战士",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     18,
					Dexterity:    12,
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
			Skill:   model.SkillAthletics,
			DC:      12,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Roll)
		assert.Equal(t, model.SkillAthletics, result.Skill)
	})

	t.Run("skill check fails with invalid actor", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for skill checks",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		_, err = e.PerformSkillCheck(ctx, SkillCheckRequest{
			GameID:  gameID,
			ActorID: model.ID("invalid-actor-id"),
			Skill:   model.SkillStealth,
		})

		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("skill check with different skills", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for skill checks",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "法师",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     8,
					Dexterity:    12,
					Constitution: 10,
					Intelligence: 18,
					Wisdom:       12,
					Charisma:     10,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		skills := []model.Skill{model.SkillArcana, model.SkillHistory, model.SkillInvestigation}
		for _, skill := range skills {
			result, err := e.PerformSkillCheck(ctx, SkillCheckRequest{
				GameID:  gameID,
				ActorID: actorID,
				Skill:   skill,
			})

			require.NoError(t, err)
			require.NotNil(t, result)
			assert.Equal(t, skill, result.Skill)
		}
	})
}

func TestPerformSavingThrow(t *testing.T) {
	t.Run("performs saving throw successfully", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for saving throws",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "战士",
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
			DC:      10,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Roll)
		assert.Equal(t, model.AbilityConstitution, result.Ability)
		assert.Equal(t, "Test Character", result.ActorName)
	})

	t.Run("saving throw succeeds", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for saving throws",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "战士",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     20,
					Dexterity:    12,
					Constitution: 18,
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
			DC:      10,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Roll)
		assert.Equal(t, model.AbilityStrength, result.Ability)
	})

	t.Run("saving throw with advantage", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for saving throws",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "战士",
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
			GameID:    gameID,
			ActorID:   actorID,
			Ability:   model.AbilityConstitution,
			DC:        12,
			Advantage: model.RollModifier{Advantage: true},
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Roll)
	})

	t.Run("saving throw with disadvantage", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for saving throws",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "战士",
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
			GameID:    gameID,
			ActorID:   actorID,
			Ability:   model.AbilityConstitution,
			DC:        12,
			Advantage: model.RollModifier{Disadvantage: true},
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Roll)
	})

	t.Run("saving throw fails with invalid actor", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for saving throws",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		_, err = e.PerformSavingThrow(ctx, SavingThrowRequest{
			GameID:  gameID,
			ActorID: model.ID("invalid-actor-id"),
			Ability: model.AbilityConstitution,
			DC:      10,
		})

		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("saving throw fails with invalid game", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		_, err := e.PerformSavingThrow(ctx, SavingThrowRequest{
			GameID:  model.ID("invalid-game-id"),
			ActorID: model.ID("actor-id"),
			Ability: model.AbilityConstitution,
			DC:      10,
		})

		assert.Error(t, err)
	})

	t.Run("saving throw with enemy actor", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for saving throws",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreateEnemy(ctx, CreateEnemyRequest{
			GameID: gameID,
			Enemy: &EnemyInput{
				Name:            "Orc",
				Description:     "A fierce orc warrior",
				AbilityScores:   AbilityScoresInput{Strength: 16, Dexterity: 12, Constitution: 16, Intelligence: 8, Wisdom: 10, Charisma: 8},
				ChallengeRating: "1/2",
				HitPoints:       15,
				ArmorClass:      13,
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		result, err := e.PerformSavingThrow(ctx, SavingThrowRequest{
			GameID:  gameID,
			ActorID: actorID,
			Ability: model.AbilityConstitution,
			DC:      13,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, model.AbilityConstitution, result.Ability)
		assert.NotNil(t, result.Roll)
	})
}

func TestGetSkillAbility(t *testing.T) {
	e := NewTestEngine(t)

	t.Run("stealth maps to dexterity", func(t *testing.T) {
		ability := e.GetSkillAbility(model.SkillStealth)
		assert.Equal(t, model.AbilityDexterity, ability)
	})

	t.Run("athletics maps to strength", func(t *testing.T) {
		ability := e.GetSkillAbility(model.SkillAthletics)
		assert.Equal(t, model.AbilityStrength, ability)
	})

	t.Run("perception maps to wisdom", func(t *testing.T) {
		ability := e.GetSkillAbility(model.SkillPerception)
		assert.Equal(t, model.AbilityWisdom, ability)
	})

	t.Run("arcana maps to intelligence", func(t *testing.T) {
		ability := e.GetSkillAbility(model.SkillArcana)
		assert.Equal(t, model.AbilityIntelligence, ability)
	})

	t.Run("persuasion maps to charisma", func(t *testing.T) {
		ability := e.GetSkillAbility(model.SkillPersuasion)
		assert.Equal(t, model.AbilityCharisma, ability)
	})

	t.Run("investigation maps to intelligence", func(t *testing.T) {
		ability := e.GetSkillAbility(model.SkillInvestigation)
		assert.Equal(t, model.AbilityIntelligence, ability)
	})
}

func TestGetPassivePerception(t *testing.T) {
	t.Run("gets passive perception for PC", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for passive perception",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "游荡者",
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
		assert.GreaterOrEqual(t, result.PassivePerception, 10)
	})

	t.Run("gets passive perception for NPC", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for passive perception",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreateNPC(ctx, CreateNPCRequest{
			GameID: gameID,
			NPC: &NPCInput{
				Name: "Villager",
				AbilityScores: AbilityScoresInput{
					Strength:     10,
					Dexterity:    10,
					Constitution: 10,
					Intelligence: 10,
					Wisdom:       12,
					Charisma:     10,
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
		assert.GreaterOrEqual(t, result.PassivePerception, 10)
	})

	t.Run("gets passive perception for enemy", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for passive perception",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreateEnemy(ctx, CreateEnemyRequest{
			GameID: gameID,
			Enemy: &EnemyInput{
				Name:            "Goblin",
				AbilityScores:   AbilityScoresInput{Strength: 10, Dexterity: 12, Constitution: 10, Intelligence: 8, Wisdom: 14, Charisma: 8},
				ChallengeRating: "1/4",
				HitPoints:       7,
				ArmorClass:      15,
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
		assert.GreaterOrEqual(t, result.PassivePerception, 10)
	})

	t.Run("gets passive perception for companion", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for passive perception",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Leader",
				Race:  "Human",
				Class: "游侠",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     12,
					Dexterity:    14,
					Constitution: 12,
					Intelligence: 10,
					Wisdom:       14,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)

		createResult, err := e.CreateCompanion(ctx, CreateCompanionRequest{
			GameID: gameID,
			Companion: &CompanionInput{
				Name:          "Wolf",
				LeaderID:      string(pcResult.Actor.ID),
				AbilityScores: AbilityScoresInput{Strength: 12, Dexterity: 14, Constitution: 12, Intelligence: 4, Wisdom: 12, Charisma: 6},
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
		assert.GreaterOrEqual(t, result.PassivePerception, 10)
	})

	t.Run("fails with invalid actor", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for passive perception",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		_, err = e.GetPassivePerception(ctx, GetPassivePerceptionRequest{
			GameID:  gameID,
			ActorID: model.ID("invalid-actor-id"),
		})

		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("fails with invalid game", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		_, err := e.GetPassivePerception(ctx, GetPassivePerceptionRequest{
			GameID:  model.ID("invalid-game-id"),
			ActorID: model.ID("actor-id"),
		})

		assert.Error(t, err)
	})
}
