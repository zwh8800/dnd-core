package engine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zwh8800/dnd-core/pkg/model"
)

func TestCreateQuest(t *testing.T) {
	t.Run("creates quest successfully", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for quests",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		npc, err := e.CreateNPC(ctx, CreateNPCRequest{
			GameID: gameID,
			NPC: &NPCInput{
				Name:        "Quest Giver",
				Description: "A merchant with a job",
				Size:        string(model.SizeMedium),
				Speed:       30,
				AbilityScores: AbilityScoresInput{
					Strength: 12, Dexterity: 12, Constitution: 12,
					Intelligence: 14, Wisdom: 12, Charisma: 16,
				},
			},
		})
		require.NoError(t, err)

		result, err := e.CreateQuest(ctx, CreateQuestRequest{GameID: gameID,
			Name:        "Retrieve the Lost Artifact",
			Description: "The merchant's precious artifact was stolen",
			GiverID:     npc.Actor.ID,
			GiverName:   "Merchant John",
			Objectives: []ObjectiveInput{
				{ID: "find_artifact", Description: "Find the artifact", Required: 1},
			},
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.Quest)
		assert.Equal(t, "Retrieve the Lost Artifact", result.Quest.Name)
	})
}

func TestGetQuest(t *testing.T) {
	t.Run("gets quest info", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for quests",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		npc, err := e.CreateNPC(ctx, CreateNPCRequest{
			GameID: gameID,
			NPC: &NPCInput{
				Name:        "Quest Giver",
				Description: "A merchant",
				Size:        string(model.SizeMedium),
				Speed:       30,
				AbilityScores: AbilityScoresInput{
					Strength: 12, Dexterity: 12, Constitution: 12,
					Intelligence: 14, Wisdom: 12, Charisma: 16,
				},
			},
		})
		require.NoError(t, err)

		createResult, err := e.CreateQuest(ctx, CreateQuestRequest{GameID: gameID,
			Name:        "Simple Task",
			Description: "Help the merchant",
			GiverID:     npc.Actor.ID,
			GiverName:   "Merchant John",
			Objectives: []ObjectiveInput{
				{ID: "help", Description: "Help", Required: 1},
			},
		})
		require.NoError(t, err)

		result, err := e.GetQuest(ctx, GetQuestRequest{GameID: gameID,
			QuestID: createResult.Quest.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "Simple Task", result.Name)
	})
}

func TestListQuests(t *testing.T) {
	t.Run("lists all quests", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for quests",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		npc, err := e.CreateNPC(ctx, CreateNPCRequest{
			GameID: gameID,
			NPC: &NPCInput{
				Name:        "Quest Giver",
				Description: "A merchant",
				Size:        string(model.SizeMedium),
				Speed:       30,
				AbilityScores: AbilityScoresInput{
					Strength: 12, Dexterity: 12, Constitution: 12,
					Intelligence: 14, Wisdom: 12, Charisma: 16,
				},
			},
		})
		require.NoError(t, err)

		_, err = e.CreateQuest(ctx, CreateQuestRequest{GameID: gameID,
			Name:        "Quest 1",
			Description: "First quest",
			GiverID:     npc.Actor.ID,
			GiverName:   "Merchant John",
			Objectives: []ObjectiveInput{
				{ID: "obj1", Description: "Objective 1", Required: 1},
			},
		})
		require.NoError(t, err)

		_, err = e.CreateQuest(ctx, CreateQuestRequest{GameID: gameID,
			Name:        "Quest 2",
			Description: "Second quest",
			GiverID:     npc.Actor.ID,
			GiverName:   "Merchant John",
			Objectives: []ObjectiveInput{
				{ID: "obj2", Description: "Objective 2", Required: 1},
			},
		})
		require.NoError(t, err)

		result, err := e.ListQuests(ctx, ListQuestsRequest{GameID: gameID})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Quests, 2)
	})
}

func TestAcceptQuest(t *testing.T) {
	t.Run("accepts quest", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for quests",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		npc, err := e.CreateNPC(ctx, CreateNPCRequest{
			GameID: gameID,
			NPC: &NPCInput{
				Name:        "Quest Giver",
				Description: "A merchant",
				Size:        string(model.SizeMedium),
				Speed:       30,
				AbilityScores: AbilityScoresInput{
					Strength: 12, Dexterity: 12, Constitution: 12,
					Intelligence: 14, Wisdom: 12, Charisma: 16,
				},
			},
		})
		require.NoError(t, err)

		createResult, err := e.CreateQuest(ctx, CreateQuestRequest{GameID: gameID,
			Name:        "Simple Task",
			Description: "Help the merchant",
			GiverID:     npc.Actor.ID,
			GiverName:   "Merchant John",
			Objectives: []ObjectiveInput{
				{ID: "help", Description: "Help", Required: 1},
			},
		})
		require.NoError(t, err)

		pc, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Adventurer",
				Race:  "Human",
				Class: "战士",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength: 16, Dexterity: 14, Constitution: 15,
					Intelligence: 10, Wisdom: 12, Charisma: 8,
				},
			},
		})
		require.NoError(t, err)

		result, err := e.AcceptQuest(ctx, AcceptQuestRequest{GameID: gameID,
			QuestID: createResult.Quest.ID,
			ActorID: pc.Actor.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Contains(t, result.Quest.AcceptedBy, pc.Actor.ID)
	})
}

func TestUpdateQuestObjective(t *testing.T) {
	t.Run("updates objective progress", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for quests",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		npc, err := e.CreateNPC(ctx, CreateNPCRequest{
			GameID: gameID,
			NPC: &NPCInput{
				Name:        "Quest Giver",
				Description: "A merchant",
				Size:        string(model.SizeMedium),
				Speed:       30,
				AbilityScores: AbilityScoresInput{
					Strength: 12, Dexterity: 12, Constitution: 12,
					Intelligence: 14, Wisdom: 12, Charisma: 16,
				},
			},
		})
		require.NoError(t, err)

		createResult, err := e.CreateQuest(ctx, CreateQuestRequest{GameID: gameID,
			Name:        "Hunt the Goblins",
			Description: "Clear out the goblin camp",
			GiverID:     npc.Actor.ID,
			GiverName:   "Merchant John",
			Objectives: []ObjectiveInput{
				{ID: "kill_goblins", Description: "Kill goblins", Required: 5},
			},
		})
		require.NoError(t, err)

		result, err := e.UpdateQuestObjective(ctx, UpdateQuestObjectiveRequest{GameID: gameID,
			QuestID:     createResult.Quest.ID,
			ObjectiveID: "kill_goblins",
			Progress:    3,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 3, result.Quest.Objectives[0].Progress)
	})
}

func TestCompleteQuest(t *testing.T) {
	t.Run("completes quest", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for quests",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		npc, err := e.CreateNPC(ctx, CreateNPCRequest{
			GameID: gameID,
			NPC: &NPCInput{
				Name:        "Quest Giver",
				Description: "A merchant",
				Size:        string(model.SizeMedium),
				Speed:       30,
				AbilityScores: AbilityScoresInput{
					Strength: 12, Dexterity: 12, Constitution: 12,
					Intelligence: 14, Wisdom: 12, Charisma: 16,
				},
			},
		})
		require.NoError(t, err)

		createResult, err := e.CreateQuest(ctx, CreateQuestRequest{GameID: gameID,
			Name:        "Simple Delivery",
			Description: "Deliver a package",
			GiverID:     npc.Actor.ID,
			GiverName:   "Merchant John",
			Objectives: []ObjectiveInput{
				{ID: "deliver", Description: "Deliver the package", Required: 1},
			},
		})
		require.NoError(t, err)

		_, err = e.UpdateQuestObjective(ctx, UpdateQuestObjectiveRequest{GameID: gameID,
			QuestID:     createResult.Quest.ID,
			ObjectiveID: "deliver",
			Progress:    1,
		})
		require.NoError(t, err)

		result, err := e.CompleteQuest(ctx, CompleteQuestRequest{GameID: gameID,
			QuestID: createResult.Quest.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, model.QuestStatusCompleted, result.Quest.Status)
	})
}

func TestFailQuest(t *testing.T) {
	t.Run("marks quest as failed", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for quests",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		npc, err := e.CreateNPC(ctx, CreateNPCRequest{
			GameID: gameID,
			NPC: &NPCInput{
				Name:        "Quest Giver",
				Description: "A merchant",
				Size:        string(model.SizeMedium),
				Speed:       30,
				AbilityScores: AbilityScoresInput{
					Strength: 12, Dexterity: 12, Constitution: 12,
					Intelligence: 14, Wisdom: 12, Charisma: 16,
				},
			},
		})
		require.NoError(t, err)

		createResult, err := e.CreateQuest(ctx, CreateQuestRequest{GameID: gameID,
			Name:        "Time-sensitive Quest",
			Description: "Complete before time runs out",
			GiverID:     npc.Actor.ID,
			GiverName:   "Merchant John",
			Objectives: []ObjectiveInput{
				{ID: "urgent_task", Description: "Complete the task", Required: 1},
			},
		})
		require.NoError(t, err)

		result, err := e.FailQuest(ctx, FailQuestRequest{GameID: gameID,
			QuestID: createResult.Quest.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, model.QuestStatusFailed, result.Quest.Status)
	})
}

func TestDeleteQuest(t *testing.T) {
	t.Run("deletes quest", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for quests",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		npc, err := e.CreateNPC(ctx, CreateNPCRequest{
			GameID: gameID,
			NPC: &NPCInput{
				Name:        "Quest Giver",
				Description: "A merchant",
				Size:        string(model.SizeMedium),
				Speed:       30,
				AbilityScores: AbilityScoresInput{
					Strength: 12, Dexterity: 12, Constitution: 12,
					Intelligence: 14, Wisdom: 12, Charisma: 16,
				},
			},
		})
		require.NoError(t, err)

		createResult, err := e.CreateQuest(ctx, CreateQuestRequest{GameID: gameID,
			Name:        "Obsolete Quest",
			Description: "This quest is no longer needed",
			GiverID:     npc.Actor.ID,
			GiverName:   "Merchant John",
			Objectives: []ObjectiveInput{
				{ID: "old_task", Description: "Old task", Required: 1},
			},
		})
		require.NoError(t, err)

		err = e.DeleteQuest(ctx, DeleteQuestRequest{GameID: gameID,
			QuestID: createResult.Quest.ID,
		})

		require.NoError(t, err)

		_, err = e.GetQuest(ctx, GetQuestRequest{GameID: gameID,
			QuestID: createResult.Quest.ID,
		})
		assert.Error(t, err)
	})
}
