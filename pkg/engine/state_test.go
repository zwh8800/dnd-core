package engine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zwh8800/dnd-core/pkg/model"
)

func TestGetStateSummary(t *testing.T) {
	t.Run("gets state summary", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Epic Campaign",
			Description: "A grand adventure",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create a PC
		pc, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Hero",
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

		// Create and set current scene
		scene, err := e.CreateScene(ctx, CreateSceneRequest{
			GameID:      gameID,
			Name:        "Tavern",
			Description: "The Rusty Nail",
			SceneType:   model.SceneTypeIndoor,
		})
		require.NoError(t, err)

		err = e.SetCurrentScene(ctx, SetCurrentSceneRequest{
			GameID:  gameID,
			SceneID: scene.Scene.ID,
		})
		require.NoError(t, err)

		// Move PC to scene
		_, err = e.MoveActorToScene(ctx, MoveActorToSceneRequest{
			GameID:  gameID,
			ActorID: pc.Actor.ID,
			SceneID: scene.Scene.ID,
		})
		require.NoError(t, err)

		// Get state summary
		summary, err := e.GetStateSummary(ctx, gameID)

		require.NoError(t, err)
		require.NotNil(t, summary)
		assert.Equal(t, "Epic Campaign", summary.GameName)
		assert.Equal(t, model.PhaseCharacterCreation, summary.Phase)
		assert.NotNil(t, summary.CurrentScene)
		assert.Equal(t, "Tavern", summary.CurrentScene.Name)
		assert.Len(t, summary.PartyMembers, 1)
	})
}

func TestGetActorSheet(t *testing.T) {
	t.Run("gets actor sheet for PC", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create PC
		pc, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:       "Detailed Hero",
				Race:       "Elf",
				Background: "Soldier",
				Class:      "法师",
				Level:      3,
				AbilityScores: AbilityScoresInput{
					Strength:     8,
					Dexterity:    14,
					Constitution: 12,
					Intelligence: 16,
					Wisdom:       10,
					Charisma:     12,
				},
			},
		})
		require.NoError(t, err)

		// Get actor sheet
		sheet, err := e.GetActorSheet(ctx, gameID, pc.Actor.ID)

		require.NoError(t, err)
		require.NotNil(t, sheet)
		assert.Contains(t, sheet.BasicInfo, "Detailed Hero")
		assert.Equal(t, 16, sheet.AbilityScores["INT"])
		assert.Equal(t, 8, sheet.AbilityScores["STR"])
		assert.Equal(t, 14, sheet.AbilityScores["DEX"])
		assert.Equal(t, 12, sheet.AbilityScores["CON"])
		assert.NotNil(t, sheet.Combat)
		assert.NotNil(t, sheet.Skills)
		assert.NotNil(t, sheet.SavingThrows)
	})

	t.Run("gets actor sheet for NPC", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create NPC
		npc, err := e.CreateNPC(ctx, CreateNPCRequest{
			GameID: gameID,
			NPC: &NPCInput{
				Name:        "Village Elder",
				Description: "A wise elder",
				Size:        string(model.SizeMedium),
				Speed:       30,
				AbilityScores: AbilityScoresInput{
					Strength:     10,
					Dexterity:    12,
					Constitution: 11,
					Intelligence: 14,
					Wisdom:       16,
					Charisma:     13,
				},
			},
		})
		require.NoError(t, err)

		// Get actor sheet
		sheet, err := e.GetActorSheet(ctx, gameID, npc.Actor.ID)

		require.NoError(t, err)
		require.NotNil(t, sheet)
		assert.Contains(t, sheet.BasicInfo, "Village Elder")
		assert.Contains(t, sheet.BasicInfo, "npc")
		assert.Equal(t, 14, sheet.AbilityScores["INT"])
		assert.Equal(t, 16, sheet.AbilityScores["WIS"])
	})

	t.Run("returns error for non-existent actor", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		_, err = e.GetActorSheet(ctx, gameID, model.NewID())

		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})
}

func TestGetCombatSummary(t *testing.T) {
	t.Run("gets combat summary", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for combat",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create combatants
		pc1, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Fighter",
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

		pc2, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Wizard",
				Race:  "Elf",
				Class: "法师",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     8,
					Dexterity:    14,
					Constitution: 12,
					Intelligence: 16,
					Wisdom:       10,
					Charisma:     12,
				},
			},
		})
		require.NoError(t, err)

		// Create scene
		scene, err := e.CreateScene(ctx, CreateSceneRequest{
			GameID:      gameID,
			Name:        "Battlefield",
			Description: "A test battlefield",
			SceneType:   model.SceneTypeOutdoor,
		})
		require.NoError(t, err)

		// Switch to exploration phase
		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "test")
		require.NoError(t, err)

		// Start combat
		_, err = e.StartCombat(ctx, StartCombatRequest{
			GameID:         gameID,
			SceneID:        scene.Scene.ID,
			ParticipantIDs: []model.ID{pc1.Actor.ID, pc2.Actor.ID},
		})
		require.NoError(t, err)

		// Get combat summary
		summary, err := e.GetCombatSummary(ctx, gameID)

		require.NoError(t, err)
		require.NotNil(t, summary)
		assert.Equal(t, 1, summary.Round)
		assert.Len(t, summary.TurnOrder, 2)
		assert.Len(t, summary.Combatants, 2)
	})

	t.Run("returns error when no combat active", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		_, err = e.GetCombatSummary(ctx, gameID)

		assert.Error(t, err)
		assert.Equal(t, ErrCombatNotActive, err)
	})
}
