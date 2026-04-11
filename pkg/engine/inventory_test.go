package engine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zwh8800/dnd-core/pkg/model"
)

func TestAddItem(t *testing.T) {
	t.Run("adds item to inventory", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for inventory",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		pc, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "战士",
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

		result, err := e.AddItem(ctx, AddItemRequest{
			GameID:  gameID,
			ActorID: pc.Actor.ID,
			Item: &ItemInput{
				Name:        "Longsword",
				Description: "A fine steel longsword",
				Type:        model.ItemTypeWeapon,
				Rarity:      model.RarityCommon,
				Quantity:    1,
			},
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.Success)
		assert.NotEqual(t, model.ID(""), result.ItemID)
	})
}

func TestGetInventory(t *testing.T) {
	t.Run("gets inventory", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for inventory",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		pc, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "战士",
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

		_, err = e.AddItem(ctx, AddItemRequest{
			GameID:  gameID,
			ActorID: pc.Actor.ID,
			Item: &ItemInput{
				Name:     "Longsword",
				Type:     model.ItemTypeWeapon,
				Rarity:   model.RarityCommon,
				Quantity: 1,
			},
		})
		require.NoError(t, err)

		result, err := e.GetInventory(ctx, GetInventoryRequest{
			GameID:  gameID,
			ActorID: pc.Actor.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 1)
	})
}

func TestEquipItem(t *testing.T) {
	t.Run("equips item", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for equipping",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		pc, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "战士",
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

		// Add item first
		addResult, err := e.AddItem(ctx, AddItemRequest{
			GameID:  gameID,
			ActorID: pc.Actor.ID,
			Item: &ItemInput{
				Name:     "Longsword",
				Type:     model.ItemTypeWeapon,
				Rarity:   model.RarityCommon,
				Quantity: 1,
			},
		})
		require.NoError(t, err)

		// Equip the item
		result, err := e.EquipItem(ctx, EquipItemRequest{
			GameID:  gameID,
			ActorID: pc.Actor.ID,
			ItemID:  addResult.ItemID,
			Slot:    model.SlotMainHand,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.Success)
	})
}
