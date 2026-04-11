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
	// 初始化背景数据
	data.InitDefaultData()
}

// TestApplyBackground 测试背景应用功能
func TestApplyBackground(t *testing.T) {
	t.Run("applies background successfully", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test background application",
		})
		require.NoError(t, err)

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameResult.Game.ID,
			PC: &PlayerCharacterInput{
				Name:  "Test Hero",
				Race:  "Human",
				Class: "战士",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:  16,
					Dexterity: 14,
				},
			},
		})
		require.NoError(t, err)

		result, err := e.ApplyBackground(ctx, ApplyBackgroundRequest{
			GameID:       gameResult.Game.ID,
			PCID:         pcResult.Actor.ID,
			BackgroundID: model.BackgroundSoldier,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, string(model.BackgroundSoldier), result.BackgroundID)
		assert.Equal(t, "士兵", result.BackgroundName)
		assert.NotEmpty(t, result.Message)
		assert.Contains(t, result.Message, "士兵")
	})

	t.Run("applies acolyte background with skill proficiencies", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test acolyte background",
		})
		require.NoError(t, err)

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameResult.Game.ID,
			PC: &PlayerCharacterInput{
				Name:  "Test Acolyte",
				Race:  "Human",
				Class: "牧师",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Wisdom: 16,
				},
			},
		})
		require.NoError(t, err)

		result, err := e.ApplyBackground(ctx, ApplyBackgroundRequest{
			GameID:       gameResult.Game.ID,
			PCID:         pcResult.Actor.ID,
			BackgroundID: model.BackgroundAcolyte,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, string(model.BackgroundAcolyte), result.BackgroundID)
		assert.Equal(t, "侍僧", result.BackgroundName)
		// Acolyte should grant skill proficiencies
		assert.NotEmpty(t, result.SkillProficiencies)
	})

	t.Run("returns error when background not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test invalid background",
		})
		require.NoError(t, err)

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameResult.Game.ID,
			PC: &PlayerCharacterInput{
				Name:  "Test Hero",
				Race:  "Human",
				Class: "战士",
				Level: 1,
			},
		})
		require.NoError(t, err)

		result, err := e.ApplyBackground(ctx, ApplyBackgroundRequest{
			GameID:       gameResult.Game.ID,
			PCID:         pcResult.Actor.ID,
			BackgroundID: "nonexistent-background",
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "background not found")
	})

	t.Run("returns error when pc not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test PC not found",
		})
		require.NoError(t, err)

		result, err := e.ApplyBackground(ctx, ApplyBackgroundRequest{
			GameID:       gameResult.Game.ID,
			PCID:         model.NewID(),
			BackgroundID: model.BackgroundSoldier,
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("returns error when game not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.ApplyBackground(ctx, ApplyBackgroundRequest{
			GameID:       model.NewID(),
			PCID:         model.NewID(),
			BackgroundID: model.BackgroundSoldier,
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("applying background grants associated feat", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test background feat",
		})
		require.NoError(t, err)

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameResult.Game.ID,
			PC: &PlayerCharacterInput{
				Name:  "Test Criminal",
				Race:  "Human",
				Class: "游荡者",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Dexterity: 16,
				},
			},
		})
		require.NoError(t, err)

		result, err := e.ApplyBackground(ctx, ApplyBackgroundRequest{
			GameID:       gameResult.Game.ID,
			PCID:         pcResult.Actor.ID,
			BackgroundID: model.BackgroundCriminal,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		// Criminal background has an associated feat
		if result.AssociatedFeat != "" {
			assert.NotEmpty(t, result.AssociatedFeat)
		}
	})

	t.Run("can apply different backgrounds to different characters", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test multiple backgrounds",
		})
		require.NoError(t, err)

		// Create first character
		pc1Result, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameResult.Game.ID,
			PC: &PlayerCharacterInput{
				Name:  "Character 1",
				Race:  "Human",
				Class: "战士",
				Level: 1,
			},
		})
		require.NoError(t, err)

		// Create second character
		pc2Result, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameResult.Game.ID,
			PC: &PlayerCharacterInput{
				Name:  "Character 2",
				Race:  "Human",
				Class: "法师",
				Level: 1,
			},
		})
		require.NoError(t, err)

		// Apply different backgrounds
		_, err = e.ApplyBackground(ctx, ApplyBackgroundRequest{
			GameID:       gameResult.Game.ID,
			PCID:         pc1Result.Actor.ID,
			BackgroundID: model.BackgroundSoldier,
		})
		require.NoError(t, err)

		_, err = e.ApplyBackground(ctx, ApplyBackgroundRequest{
			GameID:       gameResult.Game.ID,
			PCID:         pc2Result.Actor.ID,
			BackgroundID: model.BackgroundSage,
		})
		require.NoError(t, err)

		// Verify backgrounds were applied correctly
		bg1Result, err := e.GetBackgroundFeatures(ctx, GetBackgroundFeaturesRequest{
			GameID: gameResult.Game.ID,
			PCID:   pc1Result.Actor.ID,
		})
		require.NoError(t, err)
		assert.Equal(t, string(model.BackgroundSoldier), bg1Result.BackgroundID)

		bg2Result, err := e.GetBackgroundFeatures(ctx, GetBackgroundFeaturesRequest{
			GameID: gameResult.Game.ID,
			PCID:   pc2Result.Actor.ID,
		})
		require.NoError(t, err)
		assert.Equal(t, string(model.BackgroundSage), bg2Result.BackgroundID)
	})
}

// TestGetBackgroundFeatures 测试背景特性查询功能
func TestGetBackgroundFeatures(t *testing.T) {
	t.Run("gets background features after applying background", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test get background features",
		})
		require.NoError(t, err)

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameResult.Game.ID,
			PC: &PlayerCharacterInput{
				Name:  "Test Sage",
				Race:  "Human",
				Class: "法师",
				Level: 1,
			},
		})
		require.NoError(t, err)

		// Apply background first
		_, err = e.ApplyBackground(ctx, ApplyBackgroundRequest{
			GameID:       gameResult.Game.ID,
			PCID:         pcResult.Actor.ID,
			BackgroundID: model.BackgroundSage,
		})
		require.NoError(t, err)

		// Get background features
		result, err := e.GetBackgroundFeatures(ctx, GetBackgroundFeaturesRequest{
			GameID: gameResult.Game.ID,
			PCID:   pcResult.Actor.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, string(model.BackgroundSage), result.BackgroundID)
		assert.Equal(t, "学者", result.BackgroundName)
	})

	t.Run("returns empty features for character without background", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test no background",
		})
		require.NoError(t, err)

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameResult.Game.ID,
			PC: &PlayerCharacterInput{
				Name:  "Test Hero",
				Race:  "Human",
				Class: "战士",
				Level: 1,
			},
		})
		require.NoError(t, err)

		result, err := e.GetBackgroundFeatures(ctx, GetBackgroundFeaturesRequest{
			GameID: gameResult.Game.ID,
			PCID:   pcResult.Actor.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Empty(t, result.BackgroundID)
	})

	t.Run("returns error when pc not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test PC not found",
		})
		require.NoError(t, err)

		result, err := e.GetBackgroundFeatures(ctx, GetBackgroundFeaturesRequest{
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

		result, err := e.GetBackgroundFeatures(ctx, GetBackgroundFeaturesRequest{
			GameID: model.NewID(),
			PCID:   model.NewID(),
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("background features has complete data", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test background data completeness",
		})
		require.NoError(t, err)

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameResult.Game.ID,
			PC: &PlayerCharacterInput{
				Name:  "Test Acolyte",
				Race:  "Human",
				Class: "牧师",
				Level: 1,
			},
		})
		require.NoError(t, err)

		// Apply background
		_, err = e.ApplyBackground(ctx, ApplyBackgroundRequest{
			GameID:       gameResult.Game.ID,
			PCID:         pcResult.Actor.ID,
			BackgroundID: model.BackgroundAcolyte,
		})
		require.NoError(t, err)

		result, err := e.GetBackgroundFeatures(ctx, GetBackgroundFeaturesRequest{
			GameID: gameResult.Game.ID,
			PCID:   pcResult.Actor.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotEmpty(t, result.BackgroundID)
		assert.NotEmpty(t, result.BackgroundName)
	})
}
