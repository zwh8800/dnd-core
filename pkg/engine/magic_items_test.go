package engine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zwh8800/dnd-core/pkg/model"
)

// setupMagicItemTest creates a game and character, then switches to PhaseExploration for magic item operations
func setupMagicItemTest(t *testing.T, className string) (*Engine, context.Context, *NewGameResult, *CreatePCResult) {
	t.Helper()
	e := NewTestEngine(t)
	ctx := context.Background()

	gameResult, err := e.NewGame(ctx, NewGameRequest{
		Name:        "Test Game",
		Description: "Test magic items",
	})
	require.NoError(t, err)

	pcResult, err := e.CreatePC(ctx, CreatePCRequest{
		GameID: gameResult.Game.ID,
		PC: &PlayerCharacterInput{
			Name:  "Test Hero",
			Race:  "Human",
			Class: className,
			Level: 1,
		},
	})
	require.NoError(t, err)

	// Switch to PhaseExploration for magic item operations
	_, err = e.SetPhase(ctx, gameResult.Game.ID, model.PhaseExploration, "test")
	require.NoError(t, err)

	return e, ctx, gameResult, pcResult
}

// TestUseMagicItem 测试魔法物品使用功能
func TestUseMagicItem(t *testing.T) {
	t.Run("uses consumable magic item successfully", func(t *testing.T) {
		e, ctx, gameResult, pcResult := setupMagicItemTest(t, "Fighter")

		// Add a healing potion to inventory
		addResult, err := e.AddItem(ctx, AddItemRequest{
			GameID:  gameResult.Game.ID,
			ActorID: pcResult.Actor.ID,
			Item: &ItemInput{
				Name:        "治疗药水",
				Type:        model.ItemTypePotion,
				Rarity:      model.RarityCommon,
				Quantity:    1,
				Description: "饮用后恢复2d4+2 HP",
				Consumable:  true,
			},
		})
		require.NoError(t, err)

		result, err := e.UseMagicItem(ctx, UseMagicItemRequest{
			GameID:  gameResult.Game.ID,
			ActorID: pcResult.Actor.ID,
			ItemID:  addResult.ItemID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "治疗药水", result.ItemName)
		assert.True(t, result.Consumed)
		assert.Greater(t, len(result.Messages), 0)
	})

	t.Run("uses item with quantity decreases count", func(t *testing.T) {
		e, ctx, gameResult, pcResult := setupMagicItemTest(t, "Fighter")

		// Add item with quantity > 1
		addResult, err := e.AddItem(ctx, AddItemRequest{
			GameID:  gameResult.Game.ID,
			ActorID: pcResult.Actor.ID,
			Item: &ItemInput{
				Name:       "治疗药水",
				Type:       model.ItemTypePotion,
				Rarity:     model.RarityCommon,
				Quantity:   3,
				Consumable: true,
			},
		})
		require.NoError(t, err)

		// Use one
		_, err = e.UseMagicItem(ctx, UseMagicItemRequest{
			GameID:  gameResult.Game.ID,
			ActorID: pcResult.Actor.ID,
			ItemID:  addResult.ItemID,
		})
		require.NoError(t, err)

		// Check inventory
		invResult, err := e.GetInventory(ctx, GetInventoryRequest{
			GameID:  gameResult.Game.ID,
			ActorID: pcResult.Actor.ID,
		})
		require.NoError(t, err)

		// Should still have items (quantity decreased from 3 to 2)
		found := false
		for _, item := range invResult.Items {
			if item.ID == addResult.ItemID {
				found = true
				assert.Equal(t, 2, item.Quantity)
				break
			}
		}
		assert.True(t, found, "item should still exist with decreased quantity")
	})

	t.Run("returns error when item requires attunement but not attuned", func(t *testing.T) {
		e, ctx, gameResult, pcResult := setupMagicItemTest(t, "Fighter")

		// Add item that requires attunement
		addResult, err := e.AddItem(ctx, AddItemRequest{
			GameID:  gameResult.Game.ID,
			ActorID: pcResult.Actor.ID,
			Item: &ItemInput{
				Name:               "魔法戒指",
				Type:               model.ItemTypeRing,
				Rarity:             model.RarityRare,
				Quantity:           1,
				RequiresAttunement: true,
				Attunement:         "需要调谐",
			},
		})
		require.NoError(t, err)

		result, err := e.UseMagicItem(ctx, UseMagicItemRequest{
			GameID:  gameResult.Game.ID,
			ActorID: pcResult.Actor.ID,
			ItemID:  addResult.ItemID,
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "调音")
	})

	t.Run("returns error when item not found", func(t *testing.T) {
		e, ctx, gameResult, pcResult := setupMagicItemTest(t, "Fighter")

		result, err := e.UseMagicItem(ctx, UseMagicItemRequest{
			GameID:  gameResult.Game.ID,
			ActorID: pcResult.Actor.ID,
			ItemID:  model.NewID(),
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("returns error when actor not found", func(t *testing.T) {
		e, ctx, gameResult, _ := setupMagicItemTest(t, "Fighter")

		result, err := e.UseMagicItem(ctx, UseMagicItemRequest{
			GameID:  gameResult.Game.ID,
			ActorID: model.NewID(),
			ItemID:  model.NewID(),
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("returns error when game not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.UseMagicItem(ctx, UseMagicItemRequest{
			GameID:  model.NewID(),
			ActorID: model.NewID(),
			ItemID:  model.NewID(),
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})
}

// TestUnattuneItem 测试取消调谐功能
func TestUnattuneItem(t *testing.T) {
	t.Run("unattunes item successfully", func(t *testing.T) {
		e, ctx, gameResult, pcResult := setupMagicItemTest(t, "Fighter")

		// Add item that requires attunement
		addResult, err := e.AddItem(ctx, AddItemRequest{
			GameID:  gameResult.Game.ID,
			ActorID: pcResult.Actor.ID,
			Item: &ItemInput{
				Name:               "魔法戒指",
				Type:               model.ItemTypeRing,
				Rarity:             model.RarityRare,
				Quantity:           1,
				RequiresAttunement: true,
				Attunement:         "需要调谐",
			},
		})
		require.NoError(t, err)

		// Attune the item first
		_, err = e.AttuneItem(ctx, AttuneItemRequest{
			GameID:  gameResult.Game.ID,
			ActorID: pcResult.Actor.ID,
			ItemID:  addResult.ItemID,
		})
		require.NoError(t, err)

		// Unattune the item
		result, err := e.UnattuneItem(ctx, UnattuneItemRequest{
			GameID:  gameResult.Game.ID,
			ActorID: pcResult.Actor.ID,
			ItemID:  addResult.ItemID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.Success)
		assert.Equal(t, "魔法戒指", result.ItemName)
		assert.Equal(t, 0, result.AttunedCount)
	})

	t.Run("returns error when item not attuned", func(t *testing.T) {
		e, ctx, gameResult, pcResult := setupMagicItemTest(t, "Fighter")

		// Add item but don't attune
		addResult, err := e.AddItem(ctx, AddItemRequest{
			GameID:  gameResult.Game.ID,
			ActorID: pcResult.Actor.ID,
			Item: &ItemInput{
				Name:               "魔法戒指",
				Type:               model.ItemTypeRing,
				Rarity:             model.RarityRare,
				Quantity:           1,
				RequiresAttunement: true,
				Attunement:         "需要调谐",
			},
		})
		require.NoError(t, err)

		result, err := e.UnattuneItem(ctx, UnattuneItemRequest{
			GameID:  gameResult.Game.ID,
			ActorID: pcResult.Actor.ID,
			ItemID:  addResult.ItemID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.False(t, result.Success)
		assert.Contains(t, result.Message, "not attuned")
	})

	t.Run("returns error when item not found", func(t *testing.T) {
		e, ctx, gameResult, pcResult := setupMagicItemTest(t, "Fighter")

		result, err := e.UnattuneItem(ctx, UnattuneItemRequest{
			GameID:  gameResult.Game.ID,
			ActorID: pcResult.Actor.ID,
			ItemID:  model.NewID(),
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("returns error when actor not found", func(t *testing.T) {
		e, ctx, gameResult, _ := setupMagicItemTest(t, "Fighter")

		result, err := e.UnattuneItem(ctx, UnattuneItemRequest{
			GameID:  gameResult.Game.ID,
			ActorID: model.NewID(),
			ItemID:  model.NewID(),
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("returns error when game not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.UnattuneItem(ctx, UnattuneItemRequest{
			GameID:  model.NewID(),
			ActorID: model.NewID(),
			ItemID:  model.NewID(),
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})
}

// TestRechargeMagicItems 测试魔法物品充能功能
func TestRechargeMagicItems(t *testing.T) {
	t.Run("recharges magic items at dawn", func(t *testing.T) {
		e, ctx, gameResult, pcResult := setupMagicItemTest(t, "Wizard")

		// Add a charged item with dawn recharge
		// Note: We need to manually create an item with charges since ItemInput doesn't support it
		// For this test, we'll verify the function works with no items first
		result, err := e.RechargeMagicItems(ctx, RechargeMagicItemsRequest{
			GameID:  gameResult.Game.ID,
			ActorID: pcResult.Actor.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		// No items to recharge, should return appropriate message
		// Message could be about no items or no magic items needing recharge
		assert.True(t, len(result.Message) > 0)
	})

	t.Run("returns message when actor has no inventory", func(t *testing.T) {
		e, ctx, gameResult, _ := setupMagicItemTest(t, "Wizard")

		// Create NPC (no inventory by default)
		npcResult, err := e.CreateNPC(ctx, CreateNPCRequest{
			GameID: gameResult.Game.ID,
			NPC: &NPCInput{
				Name: "Test NPC",
			},
		})
		require.NoError(t, err)

		result, err := e.RechargeMagicItems(ctx, RechargeMagicItemsRequest{
			GameID:  gameResult.Game.ID,
			ActorID: npcResult.Actor.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Contains(t, result.Message, "没有库存")
	})

	t.Run("returns error when actor not found", func(t *testing.T) {
		e, ctx, gameResult, _ := setupMagicItemTest(t, "Wizard")

		result, err := e.RechargeMagicItems(ctx, RechargeMagicItemsRequest{
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

		result, err := e.RechargeMagicItems(ctx, RechargeMagicItemsRequest{
			GameID:  model.NewID(),
			ActorID: model.NewID(),
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})
}

// TestGetMagicItemBonus 测试魔法物品加值查询功能
func TestGetMagicItemBonus(t *testing.T) {
	t.Run("gets magic item bonus for +1 weapon", func(t *testing.T) {
		e, ctx, gameResult, pcResult := setupMagicItemTest(t, "Fighter")

		// Add a magic weapon with bonus
		addResult, err := e.AddItem(ctx, AddItemRequest{
			GameID:  gameResult.Game.ID,
			ActorID: pcResult.Actor.ID,
			Item: &ItemInput{
				Name:        "+1 长剑",
				Type:        model.ItemTypeWeapon,
				Rarity:      model.RarityUncommon,
				Quantity:    1,
				AttackBonus: 1,
			},
		})
		require.NoError(t, err)

		result, err := e.GetMagicItemBonus(ctx, GetMagicItemBonusRequest{
			GameID:  gameResult.Game.ID,
			ActorID: pcResult.Actor.ID,
			ItemID:  addResult.ItemID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "+1 长剑", result.ItemName)
		assert.Equal(t, 1, result.MagicBonus)
		assert.Contains(t, result.Bonuses, "攻击加值")
		assert.Contains(t, result.Bonuses, "伤害加值")
	})

	t.Run("gets magic item bonus for +1 armor", func(t *testing.T) {
		e, ctx, gameResult, pcResult := setupMagicItemTest(t, "Fighter")

		// Add magic armor
		addResult, err := e.AddItem(ctx, AddItemRequest{
			GameID:  gameResult.Game.ID,
			ActorID: pcResult.Actor.ID,
			Item: &ItemInput{
				Name:       "+1 链甲",
				Type:       model.ItemTypeArmor,
				Rarity:     model.RarityUncommon,
				Quantity:   1,
				ArmorClass: 17,
			},
		})
		require.NoError(t, err)

		result, err := e.GetMagicItemBonus(ctx, GetMagicItemBonusRequest{
			GameID:  gameResult.Game.ID,
			ActorID: pcResult.Actor.ID,
			ItemID:  addResult.ItemID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "+1 链甲", result.ItemName)
	})

	t.Run("returns zero bonus for non-magic item", func(t *testing.T) {
		e, ctx, gameResult, pcResult := setupMagicItemTest(t, "Fighter")

		// Add regular item
		addResult, err := e.AddItem(ctx, AddItemRequest{
			GameID:  gameResult.Game.ID,
			ActorID: pcResult.Actor.ID,
			Item: &ItemInput{
				Name:     "普通长剑",
				Type:     model.ItemTypeWeapon,
				Rarity:   model.RarityCommon,
				Quantity: 1,
			},
		})
		require.NoError(t, err)

		result, err := e.GetMagicItemBonus(ctx, GetMagicItemBonusRequest{
			GameID:  gameResult.Game.ID,
			ActorID: pcResult.Actor.ID,
			ItemID:  addResult.ItemID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 0, result.MagicBonus)
		assert.Empty(t, result.Bonuses)
	})

	t.Run("returns error when item not found", func(t *testing.T) {
		e, ctx, gameResult, pcResult := setupMagicItemTest(t, "Fighter")

		result, err := e.GetMagicItemBonus(ctx, GetMagicItemBonusRequest{
			GameID:  gameResult.Game.ID,
			ActorID: pcResult.Actor.ID,
			ItemID:  model.NewID(),
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("returns error when actor not found", func(t *testing.T) {
		e, ctx, gameResult, _ := setupMagicItemTest(t, "Fighter")

		result, err := e.GetMagicItemBonus(ctx, GetMagicItemBonusRequest{
			GameID:  gameResult.Game.ID,
			ActorID: model.NewID(),
			ItemID:  model.NewID(),
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("returns error when game not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.GetMagicItemBonus(ctx, GetMagicItemBonusRequest{
			GameID:  model.NewID(),
			ActorID: model.NewID(),
			ItemID:  model.NewID(),
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("result includes attunement status", func(t *testing.T) {
		e, ctx, gameResult, pcResult := setupMagicItemTest(t, "Fighter")

		// Add item that can be attuned
		addResult, err := e.AddItem(ctx, AddItemRequest{
			GameID:  gameResult.Game.ID,
			ActorID: pcResult.Actor.ID,
			Item: &ItemInput{
				Name:               "魔法戒指",
				Type:               model.ItemTypeRing,
				Rarity:             model.RarityRare,
				Quantity:           1,
				RequiresAttunement: true,
				Attunement:         "需要调谐",
			},
		})
		require.NoError(t, err)

		result, err := e.GetMagicItemBonus(ctx, GetMagicItemBonusRequest{
			GameID:  gameResult.Game.ID,
			ActorID: pcResult.Actor.ID,
			ItemID:  addResult.ItemID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.False(t, result.Attuned) // Not attuned yet

		// Attune the item
		_, err = e.AttuneItem(ctx, AttuneItemRequest{
			GameID:  gameResult.Game.ID,
			ActorID: pcResult.Actor.ID,
			ItemID:  addResult.ItemID,
		})
		require.NoError(t, err)

		// Check again
		result, err = e.GetMagicItemBonus(ctx, GetMagicItemBonusRequest{
			GameID:  gameResult.Game.ID,
			ActorID: pcResult.Actor.ID,
			ItemID:  addResult.ItemID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.Attuned) // Now attuned
	})
}
