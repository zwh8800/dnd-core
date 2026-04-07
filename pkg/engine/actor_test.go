package engine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zwh8800/dnd-core/pkg/model"
)

func TestCreatePC(t *testing.T) {
	t.Run("creates player character successfully", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		// Create a game first
		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for CreatePC",
		})
		require.NoError(t, err)
		require.NotNil(t, gameResult)
		require.NotNil(t, gameResult.Game)
		gameID := gameResult.Game.ID

		// Create a player character
		result, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:      "Test Character",
				Race:      "Human",
				Class:     "Fighter",
				Level:     1,
				Alignment: "Neutral",
				AbilityScores: AbilityScoresInput{
					Strength:     15,
					Dexterity:    14,
					Constitution: 13,
					Intelligence: 12,
					Wisdom:       10,
					Charisma:     8,
				},
			},
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.Actor)
		assert.Equal(t, "Test Character", result.Actor.Name)
		assert.Equal(t, model.ActorTypePC, result.Actor.Type)
	})

	t.Run("creates PC with custom HP", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for CreatePC",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		result, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:      "High HP Character",
				Race:      "Dwarf",
				Class:     "Barbarian",
				Level:     1,
				HitPoints: 20,
				AbilityScores: AbilityScoresInput{
					Strength:     16,
					Dexterity:    10,
					Constitution: 16,
					Intelligence: 8,
					Wisdom:       12,
					Charisma:     10,
				},
			},
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 20, result.Actor.HitPoints.Maximum)
	})
}

func TestCreateNPC(t *testing.T) {
	t.Run("creates NPC successfully", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for CreateNPC",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		result, err := e.CreateNPC(ctx, CreateNPCRequest{
			GameID: gameID,
			NPC: &NPCInput{
				Name:        "Village Elder",
				Description: "An old and wise elder",
				Size:        model.SizeMedium,
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
		require.NotNil(t, result)
		require.NotNil(t, result.Actor)
		assert.Equal(t, "Village Elder", result.Actor.Name)
		assert.Equal(t, model.ActorTypeNPC, result.Actor.Type)
	})
}

func TestCreateEnemy(t *testing.T) {
	t.Run("creates enemy successfully", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for CreateEnemy",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		result, err := e.CreateEnemy(ctx, CreateEnemyRequest{
			GameID: gameID,
			Enemy: &EnemyInput{
				Name:        "Goblin",
				Description: "A small green creature",
				Size:        model.SizeSmall,
				Speed:       30,
				AbilityScores: AbilityScoresInput{
					Strength:     8,
					Dexterity:    14,
					Constitution: 10,
					Intelligence: 10,
					Wisdom:       8,
					Charisma:     8,
				},
				ChallengeRating: 0.25,
				HitPoints:       7,
				ArmorClass:      15,
			},
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.Actor)
		assert.Equal(t, "Goblin", result.Actor.Name)
		assert.Equal(t, model.ActorTypeEnemy, result.Actor.Type)
	})
}

func TestGetActor(t *testing.T) {
	t.Run("gets actor successfully", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for GetActor",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create an actor first
		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "Rogue",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     12,
					Dexterity:    16,
					Constitution: 14,
					Intelligence: 10,
					Wisdom:       12,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		// Get the actor
		result, err := e.GetActor(ctx, GetActorRequest{
			GameID:  gameID,
			ActorID: actorID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.Actor)
		assert.Equal(t, "Test Character", result.Actor.Name)
		assert.Equal(t, actorID, result.Actor.ID)
	})

	t.Run("returns error for non-existent actor", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for GetActor",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		_, err = e.GetActor(ctx, GetActorRequest{
			GameID:  gameID,
			ActorID: model.NewID(),
		})

		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})
}

func TestGetPC(t *testing.T) {
	t.Run("gets player character info", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for GetPC",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create a PC
		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:       "Detailed Character",
				Race:       "Elf",
				Background: "Soldier",
				Class:      "Wizard",
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
		pcID := createResult.Actor.ID

		// Get PC info
		result, err := e.GetPC(ctx, GetPCRequest{
			GameID: gameID,
			PCID:   pcID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.PC)
		assert.Equal(t, "Detailed Character", result.PC.Name)
		assert.Equal(t, "Elf", result.PC.Race)
		assert.Equal(t, "Soldier", result.PC.Background)
		assert.Equal(t, 3, result.PC.TotalLevel)
	})
}

func TestUpdateActor(t *testing.T) {
	t.Run("updates actor HP", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for UpdateActor",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create an actor
		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Human",
				Class: "Cleric",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     12,
					Dexterity:    10,
					Constitution: 14,
					Intelligence: 8,
					Wisdom:       16,
					Charisma:     10,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		// Update HP
		currentHP := 5
		err = e.UpdateActor(ctx, UpdateActorRequest{
			GameID:  gameID,
			ActorID: actorID,
			Update: ActorUpdate{
				HitPoints: &HitPointUpdate{
					Current: &currentHP,
				},
			},
		})

		require.NoError(t, err)

		// Verify the update
		getResult, err := e.GetActor(ctx, GetActorRequest{
			GameID:  gameID,
			ActorID: actorID,
		})
		require.NoError(t, err)
		assert.Equal(t, 5, getResult.Actor.HitPoints.Current)
	})

	t.Run("adds conditions to actor", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for UpdateActor",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create an actor
		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Test Character",
				Race:  "Halfling",
				Class: "Rogue",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     10,
					Dexterity:    16,
					Constitution: 12,
					Intelligence: 10,
					Wisdom:       12,
					Charisma:     14,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		// Add condition
		err = e.UpdateActor(ctx, UpdateActorRequest{
			GameID:  gameID,
			ActorID: actorID,
			Update: ActorUpdate{
				Conditions: &ConditionUpdate{
					Add: []model.ConditionInstance{
						{Type: model.ConditionPoisoned},
					},
				},
			},
		})

		require.NoError(t, err)

		// Verify the condition
		getResult, err := e.GetActor(ctx, GetActorRequest{
			GameID:  gameID,
			ActorID: actorID,
		})
		require.NoError(t, err)
		assert.Contains(t, getResult.Actor.Conditions, string(model.ConditionPoisoned))
	})
}

func TestListActors(t *testing.T) {
	t.Run("lists all actors", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for ListActors",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create multiple actors
		_, err = e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Character 1",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength: 15, Dexterity: 13, Constitution: 14,
					Intelligence: 10, Wisdom: 12, Charisma: 8,
				},
			},
		})
		require.NoError(t, err)

		_, err = e.CreateNPC(ctx, CreateNPCRequest{
			GameID: gameID,
			NPC: &NPCInput{
				Name:        "Merchant",
				Description: "A friendly merchant",
				Size:        model.SizeMedium,
				Speed:       30,
				AbilityScores: AbilityScoresInput{
					Strength: 10, Dexterity: 12, Constitution: 11,
					Intelligence: 14, Wisdom: 10, Charisma: 16,
				},
			},
		})
		require.NoError(t, err)

		// List all actors
		result, err := e.ListActors(ctx, ListActorsRequest{
			GameID: gameID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Actors, 2)
	})

	t.Run("filters actors by type", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for ListActors",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create actors
		_, err = e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Player",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength: 15, Dexterity: 13, Constitution: 14,
					Intelligence: 10, Wisdom: 12, Charisma: 8,
				},
			},
		})
		require.NoError(t, err)

		_, err = e.CreateNPC(ctx, CreateNPCRequest{
			GameID: gameID,
			NPC: &NPCInput{
				Name:        "NPC",
				Description: "An NPC",
				Size:        model.SizeMedium,
				Speed:       30,
				AbilityScores: AbilityScoresInput{
					Strength: 10, Dexterity: 12, Constitution: 11,
					Intelligence: 14, Wisdom: 10, Charisma: 16,
				},
			},
		})
		require.NoError(t, err)

		// Filter by PC type
		alive := true
		result, err := e.ListActors(ctx, ListActorsRequest{
			GameID: gameID,
			Filter: &ActorFilter{
				Types: []model.ActorType{model.ActorTypePC},
				Alive: &alive,
			},
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Actors, 1)
		assert.Equal(t, model.ActorTypePC, result.Actors[0].Type)
	})
}

func TestRemoveActor(t *testing.T) {
	t.Run("removes actor successfully", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for RemoveActor",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create an actor
		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "To Be Removed",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength: 15, Dexterity: 13, Constitution: 14,
					Intelligence: 10, Wisdom: 12, Charisma: 8,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		// Remove the actor
		err = e.RemoveActor(ctx, RemoveActorRequest{
			GameID:  gameID,
			ActorID: actorID,
		})

		require.NoError(t, err)

		// Verify actor is removed
		_, err = e.GetActor(ctx, GetActorRequest{
			GameID:  gameID,
			ActorID: actorID,
		})
		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})
}

func TestAddExperience(t *testing.T) {
	t.Run("adds experience to PC", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for AddExperience",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create a PC
		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Adventurer",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength: 16, Dexterity: 12, Constitution: 15,
					Intelligence: 10, Wisdom: 11, Charisma: 8,
				},
			},
		})
		require.NoError(t, err)
		pcID := createResult.Actor.ID

		// Switch to exploration phase
		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "test")
		require.NoError(t, err)

		// Add experience
		result, err := e.AddExperience(ctx, AddExperienceRequest{
			GameID: gameID,
			PCID:   pcID,
			XP:     300,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.LeveledUp) // Level 1 -> 2 requires 300 XP
		assert.Equal(t, 1, result.OldLevel)
		assert.Equal(t, 2, result.NewLevel)
	})

	t.Run("levels up PC", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for AddExperience",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create a PC
		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Veteran",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength: 16, Dexterity: 12, Constitution: 15,
					Intelligence: 10, Wisdom: 11, Charisma: 8,
				},
			},
		})
		require.NoError(t, err)
		pcID := createResult.Actor.ID

		// Switch to exploration phase
		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "test")
		require.NoError(t, err)

		// Add enough experience to level up
		result, err := e.AddExperience(ctx, AddExperienceRequest{
			GameID: gameID,
			PCID:   pcID,
			XP:     400,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.LeveledUp)
		assert.Equal(t, 1, result.OldLevel)
		assert.Equal(t, 2, result.NewLevel)
	})
}
