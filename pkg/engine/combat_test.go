package engine

import (
	"context"
	"testing"

	"github.com/zwh8800/dnd-core/internal/model"
)

// TestStartCombat 测试开始战斗
func TestStartCombat(t *testing.T) {
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
	combat, err := engine.StartCombat(ctx, gameID, scene.Scene.ID, []model.ID{pcResult.ID, enemyResult.ID})
	if err != nil {
		t.Fatalf("Failed to start combat: %v", err)
	}

	if combat == nil {
		t.Fatal("Combat is nil")
	}
	if combat.Status != model.CombatStatusActive {
		t.Errorf("Expected combat status %s, got %s", model.CombatStatusActive, combat.Status)
	}
	if len(combat.Initiative) != 2 {
		t.Errorf("Expected 2 combatants, got %d", len(combat.Initiative))
	}
}

// TestEndCombat 测试结束战斗
func TestEndCombat(t *testing.T) {
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

	// 结束战斗
	err = engine.EndCombat(ctx, gameID)
	if err != nil {
		t.Fatalf("Failed to end combat: %v", err)
	}

	// 验证战斗已结束 - 尝试获取战斗应该返回ErrCombatNotActive
	_, err = engine.GetCombatSummary(ctx, gameID)
	if err == nil {
		t.Fatal("Expected error when getting combat after ending, got nil")
	}
	if err != ErrCombatNotActive {
		t.Errorf("Expected ErrCombatNotActive, got %v", err)
	}
}

// TestCombatActions 测试战斗动作
func TestCombatActions(t *testing.T) {
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

	// 获取当前回合的角色
	turn, err := engine.GetCurrentTurn(ctx, gameID)
	if err != nil {
		t.Fatalf("Failed to get current turn: %v", err)
	}

	if turn == nil {
		t.Fatal("Turn is nil")
	}
	if turn.ActorID == "" {
		t.Error("Expected current actor ID to be set")
	}
}

// TestNextTurn 测试下一个回合
func TestNextTurn(t *testing.T) {
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

	// 获取第一个回合
	turn1, err := engine.GetCurrentTurn(ctx, gameID)
	if err != nil {
		t.Fatalf("Failed to get current turn: %v", err)
	}

	// 下一个回合
	_, err = engine.NextTurn(ctx, gameID)
	if err != nil {
		t.Fatalf("Failed to next turn: %v", err)
	}

	// 获取新的回合
	turn2, err := engine.GetCurrentTurn(ctx, gameID)
	if err != nil {
		t.Fatalf("Failed to get current turn after next: %v", err)
	}

	// 应该是不同的角色
	if turn1.ActorID == turn2.ActorID {
		t.Error("Expected different actor after next turn")
	}
}
