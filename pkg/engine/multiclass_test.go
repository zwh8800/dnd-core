package engine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zwh8800/dnd-core/pkg/data"
	"github.com/zwh8800/dnd-core/pkg/model"
)

func init() {
	// 初始化数据
	data.InitDefaultData()
}

// TestValidateMulticlassChoice 测试多职业验证功能
func TestValidateMulticlassChoice(t *testing.T) {
	t.Run("validates multiclass successfully with sufficient ability score", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test multiclass validation",
		})
		require.NoError(t, err)

		// Create a character with high Strength for Fighter multiclass
		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameResult.Game.ID,
			PC: &PlayerCharacterInput{
				Name:  "Test Fighter",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     15, // >= 13 for Fighter multiclass
					Dexterity:    12,
					Constitution: 14,
					Intelligence: 10,
					Wisdom:       8,
					Charisma:     10,
				},
			},
		})
		require.NoError(t, err)

		result, err := e.ValidateMulticlassChoice(ctx, ValidateMulticlassRequest{
			GameID:   gameResult.Game.ID,
			PCID:     pcResult.Actor.ID,
			NewClass: model.ClassFighter,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.Valid)
		assert.True(t, result.MeetsRequirements)
		assert.Contains(t, result.Message, "满足")
	})

	t.Run("fails validation with insufficient ability score", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test multiclass validation failure",
		})
		require.NoError(t, err)

		// Create a character with low Intelligence for Wizard multiclass
		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameResult.Game.ID,
			PC: &PlayerCharacterInput{
				Name:  "Test Fighter",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     16,
					Dexterity:    12,
					Constitution: 14,
					Intelligence: 8, // < 13, fails Wizard requirement
					Wisdom:       10,
					Charisma:     10,
				},
			},
		})
		require.NoError(t, err)

		result, err := e.ValidateMulticlassChoice(ctx, ValidateMulticlassRequest{
			GameID:   gameResult.Game.ID,
			PCID:     pcResult.Actor.ID,
			NewClass: model.ClassWizard,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.False(t, result.Valid)
		assert.False(t, result.MeetsRequirements)
		assert.Equal(t, 13, result.RequiredScore)
		assert.Equal(t, 8, result.CurrentScore)
	})

	t.Run("returns error when class not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test invalid class",
		})
		require.NoError(t, err)

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameResult.Game.ID,
			PC: &PlayerCharacterInput{
				Name:  "Test Hero",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
			},
		})
		require.NoError(t, err)

		result, err := e.ValidateMulticlassChoice(ctx, ValidateMulticlassRequest{
			GameID:   gameResult.Game.ID,
			PCID:     pcResult.Actor.ID,
			NewClass: "nonexistent-class",
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "class definition not found")
	})

	t.Run("returns error when pc not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test PC not found",
		})
		require.NoError(t, err)

		result, err := e.ValidateMulticlassChoice(ctx, ValidateMulticlassRequest{
			GameID:   gameResult.Game.ID,
			PCID:     model.NewID(),
			NewClass: model.ClassFighter,
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("returns error when game not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.ValidateMulticlassChoice(ctx, ValidateMulticlassRequest{
			GameID:   model.NewID(),
			PCID:     model.NewID(),
			NewClass: model.ClassFighter,
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})
}

// TestGetMulticlassSpellSlots 测试多职业法术位计算
func TestGetMulticlassSpellSlots(t *testing.T) {
	t.Run("gets spell slots for single class caster", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test single class spell slots",
		})
		require.NoError(t, err)

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameResult.Game.ID,
			PC: &PlayerCharacterInput{
				Name:  "Test Wizard",
				Race:  "Human",
				Class: "法师",
				Level: 3,
				AbilityScores: AbilityScoresInput{
					Intelligence: 16,
				},
			},
		})
		require.NoError(t, err)

		result, err := e.GetMulticlassSpellSlots(ctx, GetMulticlassSpellSlotsRequest{
			GameID: gameResult.Game.ID,
			PCID:   pcResult.Actor.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 3, result.EffectiveCasterLevel)
		assert.False(t, result.IsMulticlassCaster)
		assert.Greater(t, len(result.SpellSlots), 0)
		// Level 3 caster should have 1st and 2nd level slots
		assert.GreaterOrEqual(t, len(result.SpellSlots), 2)
	})

	t.Run("returns no slots for non-caster", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test non-caster",
		})
		require.NoError(t, err)

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameResult.Game.ID,
			PC: &PlayerCharacterInput{
				Name:  "Test Fighter",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength: 16,
				},
			},
		})
		require.NoError(t, err)

		result, err := e.GetMulticlassSpellSlots(ctx, GetMulticlassSpellSlotsRequest{
			GameID: gameResult.Game.ID,
			PCID:   pcResult.Actor.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 0, result.EffectiveCasterLevel)
		assert.Nil(t, result.SpellSlots)
		assert.False(t, result.IsMulticlassCaster)
	})

	t.Run("returns error when pc not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test PC not found",
		})
		require.NoError(t, err)

		result, err := e.GetMulticlassSpellSlots(ctx, GetMulticlassSpellSlotsRequest{
			GameID: gameResult.Game.ID,
			PCID:   model.NewID(),
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("returns error when game not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.GetMulticlassSpellSlots(ctx, GetMulticlassSpellSlotsRequest{
			GameID: model.NewID(),
			PCID:   model.NewID(),
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("spell slot info has correct structure", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test spell slot structure",
		})
		require.NoError(t, err)

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameResult.Game.ID,
			PC: &PlayerCharacterInput{
				Name:  "Test Cleric",
				Race:  "Human",
				Class: "牧师",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Wisdom: 16,
				},
			},
		})
		require.NoError(t, err)

		result, err := e.GetMulticlassSpellSlots(ctx, GetMulticlassSpellSlotsRequest{
			GameID: gameResult.Game.ID,
			PCID:   pcResult.Actor.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		require.Greater(t, len(result.SpellSlots), 0)

		// Check first slot structure
		firstSlot := result.SpellSlots[0]
		assert.Equal(t, 1, firstSlot.Level) // 1st level
		assert.Greater(t, firstSlot.Total, 0)
		assert.GreaterOrEqual(t, firstSlot.Used, 0)
	})
}
