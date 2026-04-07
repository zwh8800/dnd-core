package engine

import (
	"context"
	"testing"

	"github.com/zwh8800/dnd-core/internal/model"
)

// TestAddItem 测试添加物品
func TestAddItem(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	// 创建PC
	pc := &model.PlayerCharacter{
		Actor: model.Actor{
			Name:  "Hero",
			Size:  model.SizeMedium,
			Speed: 30,
			AbilityScores: model.AbilityScores{
				Strength: 16, Dexterity: 12, Constitution: 14,
				Intelligence: 10, Wisdom: 8, Charisma: 13,
			},
		},
		Race: model.RaceReference{Name: "Human"},
		Classes: []model.ClassLevel{
			{ClassName: "Fighter", Level: 1},
		},
		TotalLevel: 1,
	}
	pcResult, err := engine.CreatePC(ctx, gameID, pc)
	if err != nil {
		t.Fatalf("Failed to create PC: %v", err)
	}

	// 添加物品
	item := &model.Item{
		Name:     "Health Potion",
		Type:     model.ItemTypePotion,
		Weight:   0.5,
		Quantity: 3,
		Value:    50,
	}

	result, err := engine.AddItem(ctx, gameID, pcResult.ID, item)
	if err != nil {
		t.Fatalf("Failed to add item: %v", err)
	}

	if !result.Success {
		t.Fatal("Expected add item to succeed")
	}
	if result.Quantity != 3 {
		t.Errorf("Expected quantity 3, got %d", result.Quantity)
	}
}

// TestRemoveItem 测试移除物品
func TestRemoveItem(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	// 创建PC
	pc := &model.PlayerCharacter{
		Actor: model.Actor{
			Name:  "Hero",
			Size:  model.SizeMedium,
			Speed: 30,
			AbilityScores: model.AbilityScores{
				Strength: 16, Dexterity: 12, Constitution: 14,
				Intelligence: 10, Wisdom: 8, Charisma: 13,
			},
		},
		Race: model.RaceReference{Name: "Human"},
		Classes: []model.ClassLevel{
			{ClassName: "Fighter", Level: 1},
		},
		TotalLevel: 1,
	}
	pcResult, err := engine.CreatePC(ctx, gameID, pc)
	if err != nil {
		t.Fatalf("Failed to create PC: %v", err)
	}

	// 添加物品
	item := &model.Item{
		Name:     "Health Potion",
		Type:     model.ItemTypePotion,
		Weight:   0.5,
		Quantity: 3,
		Value:    50,
	}
	addResult, err := engine.AddItem(ctx, gameID, pcResult.ID, item)
	if err != nil {
		t.Fatalf("Failed to add item: %v", err)
	}

	// 移除物品
	result, err := engine.RemoveItem(ctx, gameID, pcResult.ID, addResult.ItemID, 2)
	if err != nil {
		t.Fatalf("Failed to remove item: %v", err)
	}

	if !result.Success {
		t.Fatal("Expected remove item to succeed")
	}
	if result.Quantity != 2 {
		t.Errorf("Expected removed quantity 2, got %d", result.Quantity)
	}
}

// TestGetInventory 测试获取库存
func TestGetInventory(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	// 创建PC
	pc := &model.PlayerCharacter{
		Actor: model.Actor{
			Name:  "Hero",
			Size:  model.SizeMedium,
			Speed: 30,
			AbilityScores: model.AbilityScores{
				Strength: 16, Dexterity: 12, Constitution: 14,
				Intelligence: 10, Wisdom: 8, Charisma: 13,
			},
		},
		Race: model.RaceReference{Name: "Human"},
		Classes: []model.ClassLevel{
			{ClassName: "Fighter", Level: 1},
		},
		TotalLevel: 1,
	}
	pcResult, err := engine.CreatePC(ctx, gameID, pc)
	if err != nil {
		t.Fatalf("Failed to create PC: %v", err)
	}

	// 获取库存（应该是空的）
	inv, err := engine.GetInventory(ctx, gameID, pcResult.ID)
	if err != nil {
		t.Fatalf("Failed to get inventory: %v", err)
	}

	if inv.OwnerID != pcResult.ID {
		t.Errorf("Expected owner ID %s, got %s", pcResult.ID, inv.OwnerID)
	}
	if len(inv.Items) != 0 {
		t.Errorf("Expected 0 items, got %d", len(inv.Items))
	}
}

// TestEquipItem 测试装备物品
func TestEquipItem(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	// 创建PC
	pc := &model.PlayerCharacter{
		Actor: model.Actor{
			Name:  "Hero",
			Size:  model.SizeMedium,
			Speed: 30,
			AbilityScores: model.AbilityScores{
				Strength: 16, Dexterity: 12, Constitution: 14,
				Intelligence: 10, Wisdom: 8, Charisma: 13,
			},
		},
		Race: model.RaceReference{Name: "Human"},
		Classes: []model.ClassLevel{
			{ClassName: "Fighter", Level: 1},
		},
		TotalLevel: 1,
	}
	pcResult, err := engine.CreatePC(ctx, gameID, pc)
	if err != nil {
		t.Fatalf("Failed to create PC: %v", err)
	}

	// 添加护甲
	armor := &model.Item{
		Name:   "Chain Mail",
		Type:   model.ItemTypeArmor,
		Weight: 55,
		Value:  75,
		ArmorProps: &model.ArmorProperties{
			BaseAC: 16,
		},
	}
	addResult, err := engine.AddItem(ctx, gameID, pcResult.ID, armor)
	if err != nil {
		t.Fatalf("Failed to add armor: %v", err)
	}

	// 装备护甲
	result, err := engine.EquipItem(ctx, gameID, pcResult.ID, addResult.ItemID, model.SlotChest)
	if err != nil {
		t.Fatalf("Failed to equip item: %v", err)
	}

	if !result.Success {
		t.Fatal("Expected equip item to succeed")
	}
	if result.ItemName != "Chain Mail" {
		t.Errorf("Expected item name 'Chain Mail', got %s", result.ItemName)
	}
	if result.Slot != model.SlotChest {
		t.Errorf("Expected slot %s, got %s", model.SlotChest, result.Slot)
	}
}

// TestUnequipItem 测试卸下装备
func TestUnequipItem(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	// 创建PC
	pc := &model.PlayerCharacter{
		Actor: model.Actor{
			Name:  "Hero",
			Size:  model.SizeMedium,
			Speed: 30,
			AbilityScores: model.AbilityScores{
				Strength: 16, Dexterity: 12, Constitution: 14,
				Intelligence: 10, Wisdom: 8, Charisma: 13,
			},
		},
		Race: model.RaceReference{Name: "Human"},
		Classes: []model.ClassLevel{
			{ClassName: "Fighter", Level: 1},
		},
		TotalLevel: 1,
	}
	pcResult, err := engine.CreatePC(ctx, gameID, pc)
	if err != nil {
		t.Fatalf("Failed to create PC: %v", err)
	}

	// 添加并装备护甲
	armor := &model.Item{
		Name:   "Chain Mail",
		Type:   model.ItemTypeArmor,
		Weight: 55,
		Value:  75,
		ArmorProps: &model.ArmorProperties{
			BaseAC: 16,
		},
	}
	addResult, err := engine.AddItem(ctx, gameID, pcResult.ID, armor)
	if err != nil {
		t.Fatalf("Failed to add armor: %v", err)
	}

	_, err = engine.EquipItem(ctx, gameID, pcResult.ID, addResult.ItemID, model.SlotChest)
	if err != nil {
		t.Fatalf("Failed to equip armor: %v", err)
	}

	// 卸下护甲
	result, err := engine.UnequipItem(ctx, gameID, pcResult.ID, model.SlotChest)
	if err != nil {
		t.Fatalf("Failed to unequip item: %v", err)
	}

	if !result.Success {
		t.Fatal("Expected unequip item to succeed")
	}
	if result.ItemName != "Chain Mail" {
		t.Errorf("Expected item name 'Chain Mail', got %s", result.ItemName)
	}
}

// TestGetEquipment 测试获取装备
func TestGetEquipment(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	// 创建PC
	pc := &model.PlayerCharacter{
		Actor: model.Actor{
			Name:  "Hero",
			Size:  model.SizeMedium,
			Speed: 30,
			AbilityScores: model.AbilityScores{
				Strength: 16, Dexterity: 12, Constitution: 14,
				Intelligence: 10, Wisdom: 8, Charisma: 13,
			},
		},
		Race: model.RaceReference{Name: "Human"},
		Classes: []model.ClassLevel{
			{ClassName: "Fighter", Level: 1},
		},
		TotalLevel: 1,
	}
	pcResult, err := engine.CreatePC(ctx, gameID, pc)
	if err != nil {
		t.Fatalf("Failed to create PC: %v", err)
	}

	// 获取装备（应该是空的）
	equip, err := engine.GetEquipment(ctx, gameID, pcResult.ID)
	if err != nil {
		t.Fatalf("Failed to get equipment: %v", err)
	}

	if equip.OwnerID != pcResult.ID {
		t.Errorf("Expected owner ID %s, got %s", pcResult.ID, equip.OwnerID)
	}
	if len(equip.EquippedSlots) != 0 {
		t.Errorf("Expected 0 equipped items, got %d", len(equip.EquippedSlots))
	}
}

// TestAttuneItem 测试调谐魔法物品
func TestAttuneItem(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	// 创建PC
	pc := &model.PlayerCharacter{
		Actor: model.Actor{
			Name:  "Wizard",
			Size:  model.SizeMedium,
			Speed: 30,
			AbilityScores: model.AbilityScores{
				Strength: 8, Dexterity: 14, Constitution: 10,
				Intelligence: 16, Wisdom: 12, Charisma: 9,
			},
		},
		Race: model.RaceReference{Name: "Elf"},
		Classes: []model.ClassLevel{
			{ClassName: "Wizard", Level: 1},
		},
		TotalLevel: 1,
	}
	pcResult, err := engine.CreatePC(ctx, gameID, pc)
	if err != nil {
		t.Fatalf("Failed to create PC: %v", err)
	}

	// 添加需要调谐的魔法物品
	item := &model.Item{
		Name:       "Wand of Magic Missiles",
		Type:       model.ItemTypeWand,
		Weight:     1,
		Attunement: "requires attunement by a spellcaster",
		MagicBonus: 0,
	}
	addResult, err := engine.AddItem(ctx, gameID, pcResult.ID, item)
	if err != nil {
		t.Fatalf("Failed to add item: %v", err)
	}

	// 调谐物品
	result, err := engine.AttuneItem(ctx, gameID, pcResult.ID, addResult.ItemID)
	if err != nil {
		t.Fatalf("Failed to attune item: %v", err)
	}

	if !result.Success {
		t.Fatal("Expected attune item to succeed")
	}
	if !result.Attuned {
		t.Error("Expected item to be attuned")
	}
	if result.AttunedCount != 1 {
		t.Errorf("Expected attuned count 1, got %d", result.AttunedCount)
	}
}

// TestTransferItem 测试转移物品
func TestTransferItem(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	// 创建两个PC
	pc1 := &model.PlayerCharacter{
		Actor: model.Actor{
			Name:  "Hero 1",
			Size:  model.SizeMedium,
			Speed: 30,
			AbilityScores: model.AbilityScores{
				Strength: 16, Dexterity: 12, Constitution: 14,
				Intelligence: 10, Wisdom: 8, Charisma: 13,
			},
		},
		Race: model.RaceReference{Name: "Human"},
		Classes: []model.ClassLevel{
			{ClassName: "Fighter", Level: 1},
		},
		TotalLevel: 1,
	}
	pc1Result, err := engine.CreatePC(ctx, gameID, pc1)
	if err != nil {
		t.Fatalf("Failed to create PC 1: %v", err)
	}

	pc2 := &model.PlayerCharacter{
		Actor: model.Actor{
			Name:  "Hero 2",
			Size:  model.SizeMedium,
			Speed: 30,
			AbilityScores: model.AbilityScores{
				Strength: 14, Dexterity: 14, Constitution: 12,
				Intelligence: 10, Wisdom: 10, Charisma: 12,
			},
		},
		Race: model.RaceReference{Name: "Human"},
		Classes: []model.ClassLevel{
			{ClassName: "Rogue", Level: 1},
		},
		TotalLevel: 1,
	}
	pc2Result, err := engine.CreatePC(ctx, gameID, pc2)
	if err != nil {
		t.Fatalf("Failed to create PC 2: %v", err)
	}

	// 给PC1添加物品
	item := &model.Item{
		Name:     "Gold Coins",
		Type:     model.ItemTypeTreasure,
		Weight:   0,
		Quantity: 100,
		Value:    100,
	}
	addResult, err := engine.AddItem(ctx, gameID, pc1Result.ID, item)
	if err != nil {
		t.Fatalf("Failed to add item: %v", err)
	}

	// 转移物品给PC2
	result, err := engine.TransferItem(ctx, gameID, pc1Result.ID, pc2Result.ID, addResult.ItemID, 50)
	if err != nil {
		t.Fatalf("Failed to transfer item: %v", err)
	}

	if !result.Success {
		t.Fatal("Expected transfer item to succeed")
	}
	if result.FromActor != pc1Result.ID {
		t.Errorf("Expected from actor %s, got %s", pc1Result.ID, result.FromActor)
	}
	if result.ToActor != pc2Result.ID {
		t.Errorf("Expected to actor %s, got %s", pc2Result.ID, result.ToActor)
	}
	if result.Quantity != 50 {
		t.Errorf("Expected quantity 50, got %d", result.Quantity)
	}
}

// TestAddCurrency 测试添加货币
func TestAddCurrency(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	// 创建PC
	pc := &model.PlayerCharacter{
		Actor: model.Actor{
			Name:  "Hero",
			Size:  model.SizeMedium,
			Speed: 30,
			AbilityScores: model.AbilityScores{
				Strength: 16, Dexterity: 12, Constitution: 14,
				Intelligence: 10, Wisdom: 8, Charisma: 13,
			},
		},
		Race: model.RaceReference{Name: "Human"},
		Classes: []model.ClassLevel{
			{ClassName: "Fighter", Level: 1},
		},
		TotalLevel: 1,
	}
	pcResult, err := engine.CreatePC(ctx, gameID, pc)
	if err != nil {
		t.Fatalf("Failed to create PC: %v", err)
	}

	// 添加货币
	currency := model.Currency{
		Gold:   100,
		Silver: 50,
		Copper: 25,
	}
	result, err := engine.AddCurrency(ctx, gameID, pcResult.ID, currency)
	if err != nil {
		t.Fatalf("Failed to add currency: %v", err)
	}

	if !result.Success {
		t.Fatal("Expected add currency to succeed")
	}

	// 验证货币已添加
	inv, err := engine.GetInventory(ctx, gameID, pcResult.ID)
	if err != nil {
		t.Fatalf("Failed to get inventory: %v", err)
	}

	if inv.Currency.Gold != 100 {
		t.Errorf("Expected 100 gold, got %d", inv.Currency.Gold)
	}
	if inv.Currency.Silver != 50 {
		t.Errorf("Expected 50 silver, got %d", inv.Currency.Silver)
	}
}

// TestCalculateTotalWeight 测试计算总重量
func TestCalculateTotalWeight(t *testing.T) {
	inventory := &model.Inventory{
		Items: []*model.Item{
			{Name: "Item 1", Weight: 2.0, Quantity: 3},
			{Name: "Item 2", Weight: 1.5, Quantity: 2},
		},
	}

	total := calculateTotalWeight(inventory)
	expected := 9.0 // (2.0 * 3) + (1.5 * 2)
	if total != expected {
		t.Errorf("Expected total weight %f, got %f", expected, total)
	}
}

// TestCalculateMaxWeight 测试计算最大负重
func TestCalculateMaxWeight(t *testing.T) {
	actor := &model.Actor{
		AbilityScores: model.AbilityScores{
			Strength: 16,
		},
	}

	maxWeight := calculateMaxWeight(actor)
	expected := 240.0 // 16 * 15
	if maxWeight != expected {
		t.Errorf("Expected max weight %f, got %f", expected, maxWeight)
	}
}

// TestCanEquipToSlot 测试装备槽位验证
func TestCanEquipToSlot(t *testing.T) {
	weapon := &model.Item{Type: model.ItemTypeWeapon}
	if !canEquipToSlot(weapon, model.SlotMainHand) {
		t.Error("Expected weapon to be equippable to main hand")
	}
	if !canEquipToSlot(weapon, model.SlotOffHand) {
		t.Error("Expected weapon to be equippable to off hand")
	}

	armor := &model.Item{Type: model.ItemTypeArmor}
	if !canEquipToSlot(armor, model.SlotChest) {
		t.Error("Expected armor to be equippable to chest")
	}

	ring := &model.Item{Type: model.ItemTypeRing}
	if !canEquipToSlot(ring, model.SlotFinger1) {
		t.Error("Expected ring to be equippable to finger1")
	}
	if !canEquipToSlot(ring, model.SlotFinger2) {
		t.Error("Expected ring to be equippable to finger2")
	}
}
