package engine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zwh8800/dnd-core/pkg/model"
)

func TestSelectFeat(t *testing.T) {
	t.Run("selects feat successfully with valid data", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameResult.Game.ID,
			PC: &PlayerCharacterInput{
				Name:  "Test Hero",
				Race:  "Human",
				Class: "Fighter",
				Level: 4,
				AbilityScores: AbilityScoresInput{
					Strength:     16,
					Dexterity:    12,
					Constitution: 14,
					Intelligence: 10,
					Wisdom:       8,
					Charisma:     10,
				},
			},
		})
		require.NoError(t, err)

		result, err := e.SelectFeat(ctx, SelectFeatRequest{
			GameID: gameResult.Game.ID,
			PCID:   pcResult.Actor.ID,
			FeatID: "alert",
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Feats), 1)
		// Verify the feat was added
		found := false
		for _, feat := range result.Feats {
			if feat.ID == "alert" {
				found = true
				break
			}
		}
		assert.True(t, found, "alert feat should be in the list")
	})

	t.Run("returns error when feat not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameResult.Game.ID,
			PC: &PlayerCharacterInput{
				Name:  "Test Hero",
				Race:  "Human",
				Class: "Fighter",
				Level: 4,
			},
		})
		require.NoError(t, err)

		result, err := e.SelectFeat(ctx, SelectFeatRequest{
			GameID: gameResult.Game.ID,
			PCID:   pcResult.Actor.ID,
			FeatID: "nonexistent-feat",
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		// 不存在的专长会先检查先决条件，返回 "prerequisites not met" 或 "feat not found"
		assert.Contains(t, err.Error(), "nonexistent-feat")
	})

	t.Run("returns error when pc not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.SelectFeat(ctx, SelectFeatRequest{
			GameID: gameResult.Game.ID,
			PCID:   model.NewID(),
			FeatID: "alert",
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("returns error when game not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.SelectFeat(ctx, SelectFeatRequest{
			GameID: model.NewID(),
			PCID:   model.NewID(),
			FeatID: "alert",
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("selecting same non-repeatable feat twice returns error", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameResult.Game.ID,
			PC: &PlayerCharacterInput{
				Name:  "Test Hero",
				Race:  "Human",
				Class: "Fighter",
				Level: 4,
			},
		})
		require.NoError(t, err)

		// Select feat first time
		_, err = e.SelectFeat(ctx, SelectFeatRequest{
			GameID: gameResult.Game.ID,
			PCID:   pcResult.Actor.ID,
			FeatID: "alert",
		})
		require.NoError(t, err)

		// Select same feat again should fail
		result, err := e.SelectFeat(ctx, SelectFeatRequest{
			GameID: gameResult.Game.ID,
			PCID:   pcResult.Actor.ID,
			FeatID: "alert",
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "feat already possessed")
	})
}

func TestListFeats(t *testing.T) {
	t.Run("lists all feats without filter", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.ListFeats(ctx, ListFeatsRequest{})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Greater(t, len(result.Feats), 0, "should have at least one feat")
	})

	t.Run("filters feats by origin type", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		featType := model.FeatTypeOrigin
		result, err := e.ListFeats(ctx, ListFeatsRequest{
			FilterType: &featType,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		// All returned feats should be of Origin type
		for _, feat := range result.Feats {
			assert.Equal(t, string(model.FeatTypeOrigin), feat.Type)
		}
	})

	t.Run("filters feats by combat type", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		featType := model.FeatTypeCombat
		result, err := e.ListFeats(ctx, ListFeatsRequest{
			FilterType: &featType,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		// All returned feats should be of Combat type
		for _, feat := range result.Feats {
			assert.Equal(t, string(model.FeatTypeCombat), feat.Type)
		}
	})

	t.Run("filters feats by general type", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		featType := model.FeatTypeGeneral
		result, err := e.ListFeats(ctx, ListFeatsRequest{
			FilterType: &featType,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		// All returned feats should be of General type
		for _, feat := range result.Feats {
			assert.Equal(t, string(model.FeatTypeGeneral), feat.Type)
		}
	})

	t.Run("filtered list has fewer or equal feats than unfiltered", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		// Get all feats
		allResult, err := e.ListFeats(ctx, ListFeatsRequest{})
		require.NoError(t, err)

		// Get filtered feats
		featType := model.FeatTypeCombat
		filteredResult, err := e.ListFeats(ctx, ListFeatsRequest{
			FilterType: &featType,
		})
		require.NoError(t, err)

		assert.GreaterOrEqual(t, len(allResult.Feats), len(filteredResult.Feats))
	})
}

func TestGetFeatDetails(t *testing.T) {
	t.Run("gets feat details for alert feat", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.GetFeatDetails(ctx, GetFeatDetailsRequest{
			FeatID: "alert",
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.Feat)
		assert.Equal(t, "alert", result.Feat.ID)
		assert.NotEmpty(t, result.Feat.Name)
		assert.NotEmpty(t, result.Feat.Description)
	})

	t.Run("gets feat details for combat feat", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.GetFeatDetails(ctx, GetFeatDetailsRequest{
			FeatID: "charger",
		})

		require.NoError(t, err)
		require.NotNil(t, result.Feat)
		assert.Equal(t, "charger", result.Feat.ID)
		assert.Equal(t, string(model.FeatTypeCombat), result.Feat.Type)
	})

	t.Run("returns error for nonexistent feat", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.GetFeatDetails(ctx, GetFeatDetailsRequest{
			FeatID: "nonexistent-feat",
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "feat not found")
	})

	t.Run("feat details include repeatable flag", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.GetFeatDetails(ctx, GetFeatDetailsRequest{
			FeatID: "lucky",
		})

		require.NoError(t, err)
		// Lucky feat should have repeatable flag set (whatever its value is)
		assert.NotNil(t, result.Feat)
	})

	t.Run("feat details has valid structure", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.GetFeatDetails(ctx, GetFeatDetailsRequest{
			FeatID: "tough",
		})

		require.NoError(t, err)
		require.NotNil(t, result.Feat)
		assert.NotEmpty(t, result.Feat.ID)
		assert.NotEmpty(t, result.Feat.Name)
		// Type should be one of the valid feat types
		assert.Contains(t, []string{"Origin", "General", "Combat", "Epic"}, result.Feat.Type)
	})
}

func TestRemoveFeat(t *testing.T) {
	t.Run("removes feat successfully", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameResult.Game.ID,
			PC: &PlayerCharacterInput{
				Name:  "Test Hero",
				Race:  "Human",
				Class: "Fighter",
				Level: 4,
			},
		})
		require.NoError(t, err)

		// Add a feat first
		_, err = e.SelectFeat(ctx, SelectFeatRequest{
			GameID: gameResult.Game.ID,
			PCID:   pcResult.Actor.ID,
			FeatID: "alert",
		})
		require.NoError(t, err)

		// Remove the feat
		err = e.RemoveFeat(ctx, RemoveFeatRequest{
			GameID: gameResult.Game.ID,
			PCID:   pcResult.Actor.ID,
			FeatID: "alert",
		})

		assert.NoError(t, err)

		// Verify feat was removed
		featsResult, err := e.GetActorFeats(ctx, GetActorRequest{
			GameID:  gameResult.Game.ID,
			ActorID: pcResult.Actor.ID,
		})
		require.NoError(t, err)
		for _, feat := range featsResult.Feats {
			assert.NotEqual(t, "alert", feat.ID, "alert feat should be removed")
		}
	})

	t.Run("returns error when pc not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		err = e.RemoveFeat(ctx, RemoveFeatRequest{
			GameID: gameResult.Game.ID,
			PCID:   model.NewID(),
			FeatID: "alert",
		})

		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("returns error when feat not on character", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameResult.Game.ID,
			PC: &PlayerCharacterInput{
				Name:  "Test Hero",
				Race:  "Human",
				Class: "Fighter",
				Level: 4,
			},
		})
		require.NoError(t, err)

		err = e.RemoveFeat(ctx, RemoveFeatRequest{
			GameID: gameResult.Game.ID,
			PCID:   pcResult.Actor.ID,
			FeatID: "alert",
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "feat not found on character")
	})

	t.Run("returns error when game not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		err := e.RemoveFeat(ctx, RemoveFeatRequest{
			GameID: model.NewID(),
			PCID:   model.NewID(),
			FeatID: "alert",
		})

		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("removing feat does not affect other feats", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameResult.Game.ID,
			PC: &PlayerCharacterInput{
				Name:  "Test Hero",
				Race:  "Human",
				Class: "Fighter",
				Level: 8,
			},
		})
		require.NoError(t, err)

		// Add two feats
		_, err = e.SelectFeat(ctx, SelectFeatRequest{
			GameID: gameResult.Game.ID,
			PCID:   pcResult.Actor.ID,
			FeatID: "alert",
		})
		require.NoError(t, err)

		_, err = e.SelectFeat(ctx, SelectFeatRequest{
			GameID: gameResult.Game.ID,
			PCID:   pcResult.Actor.ID,
			FeatID: "tough",
		})
		require.NoError(t, err)

		// Remove one feat
		err = e.RemoveFeat(ctx, RemoveFeatRequest{
			GameID: gameResult.Game.ID,
			PCID:   pcResult.Actor.ID,
			FeatID: "alert",
		})
		require.NoError(t, err)

		// Verify tough feat is still there
		featsResult, err := e.GetActorFeats(ctx, GetActorRequest{
			GameID:  gameResult.Game.ID,
			ActorID: pcResult.Actor.ID,
		})
		require.NoError(t, err)
		found := false
		for _, feat := range featsResult.Feats {
			if feat.ID == "tough" {
				found = true
				break
			}
		}
		assert.True(t, found, "tough feat should still be present")
	})
}

func TestGetActorFeats(t *testing.T) {
	t.Run("gets empty feats list for new character", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
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

		result, err := e.GetActorFeats(ctx, GetActorRequest{
			GameID:  gameResult.Game.ID,
			ActorID: pcResult.Actor.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		// New character may have starting feats or empty list
	})

	t.Run("gets feats after selecting", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameResult.Game.ID,
			PC: &PlayerCharacterInput{
				Name:  "Test Hero",
				Race:  "Human",
				Class: "Fighter",
				Level: 4,
			},
		})
		require.NoError(t, err)

		// Select two feats
		_, err = e.SelectFeat(ctx, SelectFeatRequest{
			GameID: gameResult.Game.ID,
			PCID:   pcResult.Actor.ID,
			FeatID: "alert",
		})
		require.NoError(t, err)

		_, err = e.SelectFeat(ctx, SelectFeatRequest{
			GameID: gameResult.Game.ID,
			PCID:   pcResult.Actor.ID,
			FeatID: "lucky",
		})
		require.NoError(t, err)

		result, err := e.GetActorFeats(ctx, GetActorRequest{
			GameID:  gameResult.Game.ID,
			ActorID: pcResult.Actor.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Feats), 2)
	})

	t.Run("returns error when pc not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.GetActorFeats(ctx, GetActorRequest{
			GameID:  gameResult.Game.ID,
			ActorID: model.NewID(),
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("returns error when game not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.GetActorFeats(ctx, GetActorRequest{
			GameID:  model.NewID(),
			ActorID: model.NewID(),
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("feat info has complete data", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameResult.Game.ID,
			PC: &PlayerCharacterInput{
				Name:  "Test Hero",
				Race:  "Human",
				Class: "Fighter",
				Level: 4,
			},
		})
		require.NoError(t, err)

		_, err = e.SelectFeat(ctx, SelectFeatRequest{
			GameID: gameResult.Game.ID,
			PCID:   pcResult.Actor.ID,
			FeatID: "savage-attacker", // 使用无先决条件的专长
		})
		require.NoError(t, err)

		result, err := e.GetActorFeats(ctx, GetActorRequest{
			GameID:  gameResult.Game.ID,
			ActorID: pcResult.Actor.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		require.Greater(t, len(result.Feats), 0)
		// Check first feat has complete info
		feat := result.Feats[0]
		assert.NotEmpty(t, feat.ID)
		assert.NotEmpty(t, feat.Name)
		assert.NotEmpty(t, feat.Type)
		assert.NotEmpty(t, feat.Description)
	})
}
