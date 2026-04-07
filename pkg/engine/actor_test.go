package engine

import (
	"context"
	"testing"

	"github.com/zwh8800/dnd-core/internal/model"
)

// TestCreatePC 测试创建玩家角色
func TestCreatePC(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	pc := &model.PlayerCharacter{
		Actor: model.Actor{
			Name:  "Test Hero",
			Size:  model.SizeMedium,
			Speed: 30,
			AbilityScores: model.AbilityScores{
				Strength: 15, Dexterity: 12, Constitution: 14,
				Intelligence: 10, Wisdom: 8, Charisma: 13,
			},
		},
		Race: model.RaceReference{Name: "Human"},
		Classes: []model.ClassLevel{
			{ClassName: "Fighter", Level: 1},
		},
		TotalLevel: 1,
	}

	result, err := engine.CreatePC(ctx, gameID, pc)
	if err != nil {
		t.Fatalf("Failed to create PC: %v", err)
	}

	if result == nil {
		t.Fatal("Result is nil")
	}
	if result.Name != "Test Hero" {
		t.Errorf("Expected name 'Test Hero', got %s", result.Name)
	}
	if result.Race.Name != "Human" {
		t.Errorf("Expected race 'Human', got %s", result.Race.Name)
	}
	if result.TotalLevel != 1 {
		t.Errorf("Expected level 1, got %d", result.TotalLevel)
	}
}

// TestCreateNPC 测试创建NPC
func TestCreateNPC(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	npc := &model.NPC{
		Actor: model.Actor{
			Name:  "Village Elder",
			Size:  model.SizeMedium,
			Speed: 30,
			AbilityScores: model.AbilityScores{
				Strength: 10, Dexterity: 10, Constitution: 10,
				Intelligence: 12, Wisdom: 14, Charisma: 16,
			},
		},
	}

	result, err := engine.CreateNPC(ctx, gameID, npc)
	if err != nil {
		t.Fatalf("Failed to create NPC: %v", err)
	}

	if result == nil {
		t.Fatal("Result is nil")
	}
	if result.Name != "Village Elder" {
		t.Errorf("Expected name 'Village Elder', got %s", result.Name)
	}
}

// TestCreateEnemy 测试创建敌人
func TestCreateEnemy(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

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

	result, err := engine.CreateEnemy(ctx, gameID, enemy)
	if err != nil {
		t.Fatalf("Failed to create enemy: %v", err)
	}

	if result == nil {
		t.Fatal("Result is nil")
	}
	if result.Name != "Goblin" {
		t.Errorf("Expected name 'Goblin', got %s", result.Name)
	}
	if result.ChallengeRating != 0.25 {
		t.Errorf("Expected CR 0.25, got %f", result.ChallengeRating)
	}
}

// TestGetActor 测试获取角色
func TestGetActor(t *testing.T) {
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

	// 获取角色
	actor, err := engine.GetActor(ctx, gameID, pcResult.ID)
	if err != nil {
		t.Fatalf("Failed to get actor: %v", err)
	}

	if actor.ID != pcResult.ID {
		t.Errorf("Expected actor ID %s, got %s", pcResult.ID, actor.ID)
	}
	if actor.Name != "Hero" {
		t.Errorf("Expected name 'Hero', got %s", actor.Name)
	}
}

// TestListActors 测试列出角色
func TestListActors(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	// 创建多个PC
	for i := 0; i < 3; i++ {
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
			t.Fatalf("Failed to create PC %d: %v", i, err)
		}
	}

	// 列出所有角色
	actors, err := engine.ListActors(ctx, gameID, nil)
	if err != nil {
		t.Fatalf("Failed to list actors: %v", err)
	}

	if len(actors) != 3 {
		t.Errorf("Expected 3 actors, got %d", len(actors))
	}
}

// TestRemoveActor 测试删除角色
func TestRemoveActor(t *testing.T) {
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

	// 删除角色
	err = engine.RemoveActor(ctx, gameID, pcResult.ID)
	if err != nil {
		t.Fatalf("Failed to remove actor: %v", err)
	}

	// 验证角色已删除
	_, err = engine.GetActor(ctx, gameID, pcResult.ID)
	if err == nil {
		t.Fatal("Expected error when getting deleted actor, got nil")
	}
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}
