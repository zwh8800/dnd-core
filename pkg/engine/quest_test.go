package engine

import (
	"context"
	"testing"

	"github.com/zwh8800/dnd-core/internal/model"
)

// TestCreateQuest 测试创建任务
func TestCreateQuest(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	// 创建NPC作为任务发布者
	npc := &model.NPC{
		Actor: model.Actor{
			Name: "Quest Giver",
			AbilityScores: model.AbilityScores{
				Strength: 10, Dexterity: 10, Constitution: 10,
				Intelligence: 12, Wisdom: 14, Charisma: 16,
			},
		},
	}
	npcResult, err := engine.CreateNPC(ctx, gameID, npc)
	if err != nil {
		t.Fatalf("Failed to create NPC: %v", err)
	}

	// 创建任务
	input := QuestInput{
		Name:        "Find the Lost Artifact",
		Description: "Search for the ancient artifact in the ruins",
		GiverID:     npcResult.ID,
		GiverName:   "Quest Giver",
		Objectives: []ObjectiveInput{
			{
				ID:          "find_artifact",
				Description: "Locate the lost artifact",
				Required:    1,
				Optional:    false,
			},
			{
				ID:          "return_safely",
				Description: "Return to the quest giver safely",
				Required:    1,
				Optional:    false,
			},
		},
		Rewards: &model.QuestRewards{
			Experience: 500,
			Gold:       100,
		},
	}

	result, err := engine.CreateQuest(ctx, gameID, input)
	if err != nil {
		t.Fatalf("Failed to create quest: %v", err)
	}

	if result.Quest == nil {
		t.Fatal("Result quest is nil")
	}
	if result.Quest.Name != "Find the Lost Artifact" {
		t.Errorf("Expected quest name 'Find the Lost Artifact', got %s", result.Quest.Name)
	}
	if len(result.Quest.Objectives) != 2 {
		t.Errorf("Expected 2 objectives, got %d", len(result.Quest.Objectives))
	}
}

// TestAcceptQuest 测试接受任务
func TestAcceptQuest(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	// 创建PC
	pc := &model.PlayerCharacter{
		Actor: model.Actor{
			Name: "Hero",
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

	// 创建NPC
	npc := &model.NPC{
		Actor: model.Actor{
			Name: "Quest Giver",
			AbilityScores: model.AbilityScores{
				Strength: 10, Dexterity: 10, Constitution: 10,
				Intelligence: 12, Wisdom: 14, Charisma: 16,
			},
		},
	}
	npcResult, err := engine.CreateNPC(ctx, gameID, npc)
	if err != nil {
		t.Fatalf("Failed to create NPC: %v", err)
	}

	// 创建任务
	questInput := QuestInput{
		Name:        "Test Quest",
		Description: "A test quest",
		GiverID:     npcResult.ID,
		GiverName:   "Quest Giver",
		Objectives: []ObjectiveInput{
			{
				ID:          "objective1",
				Description: "Complete objective 1",
				Required:    1,
				Optional:    false,
			},
		},
	}
	questResult, err := engine.CreateQuest(ctx, gameID, questInput)
	if err != nil {
		t.Fatalf("Failed to create quest: %v", err)
	}

	// 接受任务
	result, err := engine.AcceptQuest(ctx, gameID, questResult.Quest.ID, pcResult.ID)
	if err != nil {
		t.Fatalf("Failed to accept quest: %v", err)
	}

	if result.Quest.Status != model.QuestStatusActive {
		t.Errorf("Expected quest status %s, got %s", model.QuestStatusActive, result.Quest.Status)
	}
	if len(result.Quest.AcceptedBy) != 1 {
		t.Errorf("Expected 1 accepted by, got %d", len(result.Quest.AcceptedBy))
	}
}

// TestUpdateQuestObjective 测试更新任务目标
func TestUpdateQuestObjective(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	// 创建PC
	pc := &model.PlayerCharacter{
		Actor: model.Actor{
			Name: "Hero",
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

	// 创建NPC
	npc := &model.NPC{
		Actor: model.Actor{
			Name: "Quest Giver",
			AbilityScores: model.AbilityScores{
				Strength: 10, Dexterity: 10, Constitution: 10,
				Intelligence: 12, Wisdom: 14, Charisma: 16,
			},
		},
	}
	npcResult, err := engine.CreateNPC(ctx, gameID, npc)
	if err != nil {
		t.Fatalf("Failed to create NPC: %v", err)
	}

	// 创建任务
	questInput := QuestInput{
		Name:        "Test Quest",
		Description: "A test quest",
		GiverID:     npcResult.ID,
		GiverName:   "Quest Giver",
		Objectives: []ObjectiveInput{
			{
				ID:          "objective1",
				Description: "Complete objective 1",
				Required:    1,
				Optional:    false,
			},
		},
	}
	questResult, err := engine.CreateQuest(ctx, gameID, questInput)
	if err != nil {
		t.Fatalf("Failed to create quest: %v", err)
	}

	// 接受任务
	_, err = engine.AcceptQuest(ctx, gameID, questResult.Quest.ID, pcResult.ID)
	if err != nil {
		t.Fatalf("Failed to accept quest: %v", err)
	}

	// 更新目标进度
	result, err := engine.UpdateQuestObjective(ctx, gameID, questResult.Quest.ID, "objective1", 1)
	if err != nil {
		t.Fatalf("Failed to update quest objective: %v", err)
	}

	// 任务应该自动完成
	if result.Quest.Status != model.QuestStatusCompleted {
		t.Errorf("Expected quest status %s after completing all objectives, got %s",
			model.QuestStatusCompleted, result.Quest.Status)
	}
}

// TestCompleteQuest 测试完成任务
func TestCompleteQuest(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	// 创建PC
	pc := &model.PlayerCharacter{
		Actor: model.Actor{
			Name: "Hero",
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

	// 创建NPC
	npc := &model.NPC{
		Actor: model.Actor{
			Name: "Quest Giver",
			AbilityScores: model.AbilityScores{
				Strength: 10, Dexterity: 10, Constitution: 10,
				Intelligence: 12, Wisdom: 14, Charisma: 16,
			},
		},
	}
	npcResult, err := engine.CreateNPC(ctx, gameID, npc)
	if err != nil {
		t.Fatalf("Failed to create NPC: %v", err)
	}

	// 创建任务
	questInput := QuestInput{
		Name:        "Test Quest",
		Description: "A test quest",
		GiverID:     npcResult.ID,
		GiverName:   "Quest Giver",
		Objectives: []ObjectiveInput{
			{
				ID:          "objective1",
				Description: "Complete objective 1",
				Required:    1,
				Optional:    false,
			},
		},
		Rewards: &model.QuestRewards{
			Experience: 500,
		},
	}
	questResult, err := engine.CreateQuest(ctx, gameID, questInput)
	if err != nil {
		t.Fatalf("Failed to create quest: %v", err)
	}

	// 接受任务
	_, err = engine.AcceptQuest(ctx, gameID, questResult.Quest.ID, pcResult.ID)
	if err != nil {
		t.Fatalf("Failed to accept quest: %v", err)
	}

	// 更新目标进度
	_, err = engine.UpdateQuestObjective(ctx, gameID, questResult.Quest.ID, "objective1", 1)
	if err != nil {
		t.Fatalf("Failed to update quest objective: %v", err)
	}

	// 完成任务
	result, err := engine.CompleteQuest(ctx, gameID, questResult.Quest.ID)
	if err != nil {
		t.Fatalf("Failed to complete quest: %v", err)
	}

	if result.Quest.Status != model.QuestStatusCompleted {
		t.Errorf("Expected quest status %s, got %s", model.QuestStatusCompleted, result.Quest.Status)
	}
}

// TestListQuests 测试列出任务
func TestListQuests(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	// 创建NPC
	npc := &model.NPC{
		Actor: model.Actor{
			Name: "Quest Giver",
			AbilityScores: model.AbilityScores{
				Strength: 10, Dexterity: 10, Constitution: 10,
				Intelligence: 12, Wisdom: 14, Charisma: 16,
			},
		},
	}
	npcResult, err := engine.CreateNPC(ctx, gameID, npc)
	if err != nil {
		t.Fatalf("Failed to create NPC: %v", err)
	}

	// 创建多个任务
	for i := 0; i < 3; i++ {
		questInput := QuestInput{
			Name:        "Quest",
			Description: "A quest",
			GiverID:     npcResult.ID,
			GiverName:   "Quest Giver",
			Objectives: []ObjectiveInput{
				{
					ID:          "obj1",
					Description: "Objective 1",
					Required:    1,
					Optional:    false,
				},
			},
		}
		_, err := engine.CreateQuest(ctx, gameID, questInput)
		if err != nil {
			t.Fatalf("Failed to create quest %d: %v", i, err)
		}
	}

	// 列出所有任务
	quests, err := engine.ListQuests(ctx, gameID, nil)
	if err != nil {
		t.Fatalf("Failed to list quests: %v", err)
	}

	if len(quests) != 3 {
		t.Errorf("Expected 3 quests, got %d", len(quests))
	}
}
