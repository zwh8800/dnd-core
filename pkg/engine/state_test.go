package engine

import (
	"context"
	"testing"

	"github.com/zwh8800/dnd-core/internal/model"
)

// TestGetStateSummary 测试获取状态摘要
func TestGetStateSummary(t *testing.T) {
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
	_, err := engine.CreatePC(ctx, gameID, pc)
	if err != nil {
		t.Fatalf("Failed to create PC: %v", err)
	}

	// 获取状态摘要
	summary, err := engine.GetStateSummary(ctx, gameID)
	if err != nil {
		t.Fatalf("Failed to get state summary: %v", err)
	}

	if summary.GameName != "Test Game" {
		t.Errorf("Expected game name 'Test Game', got %s", summary.GameName)
	}
	if len(summary.PartyMembers) != 1 {
		t.Errorf("Expected 1 party member, got %d", len(summary.PartyMembers))
	}
}

// TestGetActorSheet 测试获取角色卡
func TestGetActorSheet(t *testing.T) {
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

	// 获取角色卡
	sheet, err := engine.GetActorSheet(ctx, gameID, pcResult.ID)
	if err != nil {
		t.Fatalf("Failed to get actor sheet: %v", err)
	}

	if sheet.AbilityScores["STR"] != 8 {
		t.Errorf("Expected STR 8, got %d", sheet.AbilityScores["STR"])
	}
	if sheet.AbilityScores["INT"] != 16 {
		t.Errorf("Expected INT 16, got %d", sheet.AbilityScores["INT"])
	}
	if len(sheet.Skills) == 0 {
		t.Error("Expected skills to be populated")
	}
	if len(sheet.SavingThrows) == 0 {
		t.Error("Expected saving throws to be populated")
	}
}

// TestGetCombatSummary 测试获取战斗摘要
func TestGetCombatSummary(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	// 创建场景
	scene, err := engine.CreateScene(ctx, gameID, "Battlefield", "A battlefield", model.SceneTypeOutdoor)
	if err != nil {
		t.Fatalf("Failed to create scene: %v", err)
	}

	// 创建PC
	pc := &model.PlayerCharacter{
		Actor: model.Actor{
			Name:  "Fighter",
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

	// 创建敌人
	enemy := &model.Enemy{
		Actor: model.Actor{
			Name:  "Goblin",
			Size:  model.SizeSmall,
			Speed: 30,
			AbilityScores: model.AbilityScores{
				Strength: 8, Dexterity: 14, Constitution: 10,
				Intelligence: 10, Wisdom: 8, Charisma: 8,
			},
			HitPoints:  model.HitPoints{Current: 7, Maximum: 7},
			ArmorClass: 15,
		},
		ChallengeRating: 0.25,
	}
	enemyResult, err := engine.CreateEnemy(ctx, gameID, enemy)
	if err != nil {
		t.Fatalf("Failed to create enemy: %v", err)
	}

	// 开始战斗
	_, err = engine.StartCombat(ctx, gameID, scene.Scene.ID, []model.ID{pcResult.ID, enemyResult.ID})
	if err != nil {
		t.Fatalf("Failed to start combat: %v", err)
	}

	// 获取战斗摘要
	combat, err := engine.GetCombatSummary(ctx, gameID)
	if err != nil {
		t.Fatalf("Failed to get combat summary: %v", err)
	}

	if combat.Round != 1 {
		t.Errorf("Expected round 1, got %d", combat.Round)
	}
	if len(combat.TurnOrder) != 2 {
		t.Errorf("Expected 2 in turn order, got %d", len(combat.TurnOrder))
	}
	if len(combat.Combatants) != 2 {
		t.Errorf("Expected 2 combatants, got %d", len(combat.Combatants))
	}
	if combat.CurrentActor == "" {
		t.Error("Expected current actor to be set")
	}
}

// TestGetCombatSummaryNotActive 测试无战斗时获取战斗摘要
func TestGetCombatSummaryNotActive(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	_, err := engine.GetCombatSummary(ctx, gameID)
	if err == nil {
		t.Fatal("Expected error when no combat is active, got nil")
	}
	if err != ErrCombatNotActive {
		t.Errorf("Expected ErrCombatNotActive, got %v", err)
	}
}
