package engine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zwh8800/dnd-core/pkg/model"
)

func TestCreateScene(t *testing.T) {
	t.Run("creates scene successfully", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for scenes",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		result, err := e.CreateScene(ctx, CreateSceneRequest{GameID: gameID,
			Name:        "Tavern",
			Description: "The Rusty Nail - a cozy tavern",
			SceneType:   model.SceneTypeIndoor,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.Scene)
		assert.Equal(t, "Tavern", result.Scene.Name)
		assert.Equal(t, model.SceneTypeIndoor, result.Scene.Type)
	})
}

func TestGetScene(t *testing.T) {
	t.Run("gets scene info", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for scenes",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create a scene
		createResult, err := e.CreateScene(ctx, CreateSceneRequest{GameID: gameID,
			Name:        "Blacksmith",
			Description: "The village blacksmith shop",
			SceneType:   model.SceneTypeIndoor,
		})
		require.NoError(t, err)

		// Get scene
		result, err := e.GetScene(ctx, GetSceneRequest{GameID: gameID,
			SceneID: createResult.Scene.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "Blacksmith", result.Name)
		assert.Equal(t, "The village blacksmith shop", result.Description)
	})
}

func TestUpdateScene(t *testing.T) {
	t.Run("updates scene successfully", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for scenes",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create a scene
		createResult, err := e.CreateScene(ctx, CreateSceneRequest{GameID: gameID,
			Name:        "Cave",
			Description: "A dark cave",
			SceneType:   model.SceneTypeDungeon,
		})
		require.NoError(t, err)

		// Update scene
		err = e.UpdateScene(ctx, UpdateSceneRequest{GameID: gameID,
			SceneID: createResult.Scene.ID,
			Updates: SceneUpdate{
				Description: "A well-lit cave with glowing crystals",
				IsDark:      false,
				IsDarkSet:   true,
				LightLevel:  "bright",
			},
		})

		require.NoError(t, err)

		// Verify update
		result, err := e.GetScene(ctx, GetSceneRequest{GameID: gameID,
			SceneID: createResult.Scene.ID,
		})
		require.NoError(t, err)
		assert.Equal(t, "A well-lit cave with glowing crystals", result.Description)
		assert.False(t, result.IsDark)
		assert.Equal(t, "bright", result.LightLevel)
	})
}

func TestListScenes(t *testing.T) {
	t.Run("lists all scenes", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for scenes",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create scenes
		_, err = e.CreateScene(ctx, CreateSceneRequest{GameID: gameID,
			Name:        "Scene 1",
			Description: "Description 1",
			SceneType:   model.SceneTypeIndoor,
		})
		require.NoError(t, err)

		_, err = e.CreateScene(ctx, CreateSceneRequest{GameID: gameID,
			Name:        "Scene 2",
			Description: "Description 2",
			SceneType:   model.SceneTypeOutdoor,
		})
		require.NoError(t, err)

		// List scenes
		result, err := e.ListScenes(ctx, ListScenesRequest{GameID: gameID})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Scenes, 2)
	})
}

func TestSetCurrentScene(t *testing.T) {
	t.Run("sets current scene", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for scenes",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create a scene
		createResult, err := e.CreateScene(ctx, CreateSceneRequest{GameID: gameID,
			Name:        "Tavern",
			Description: "The Rusty Nail",
			SceneType:   model.SceneTypeIndoor,
		})
		require.NoError(t, err)

		// Set current scene
		err = e.SetCurrentScene(ctx, SetCurrentSceneRequest{GameID: gameID,
			SceneID: createResult.Scene.ID,
		})

		require.NoError(t, err)

		// Get current scene
		result, err := e.GetCurrentScene(ctx, GetCurrentSceneRequest{GameID: gameID})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "Tavern", result.Name)
	})
}

func TestMoveActorToScene(t *testing.T) {
	t.Run("moves actor to scene", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for scenes",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create scene
		scene, err := e.CreateScene(ctx, CreateSceneRequest{GameID: gameID,
			Name:        "Tavern",
			Description: "The Rusty Nail",
			SceneType:   model.SceneTypeIndoor,
		})
		require.NoError(t, err)

		// Create PC
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

		// Move to scene
		result, err := e.MoveActorToScene(ctx, MoveActorToSceneRequest{GameID: gameID,
			ActorID: pc.Actor.ID,
			SceneID: scene.Scene.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.SceneMoveResult)
		assert.True(t, result.SceneMoveResult.Success)
		assert.Equal(t, scene.Scene.ID, result.SceneMoveResult.ToScene)
	})
}

func TestGetSceneActors(t *testing.T) {
	t.Run("gets actors in scene", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for scenes",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create scene
		scene, err := e.CreateScene(ctx, CreateSceneRequest{GameID: gameID,
			Name:        "Tavern",
			Description: "The Rusty Nail",
			SceneType:   model.SceneTypeIndoor,
		})
		require.NoError(t, err)

		// Create PC and move to scene
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

		_, err = e.MoveActorToScene(ctx, MoveActorToSceneRequest{GameID: gameID,
			ActorID: pc.Actor.ID,
			SceneID: scene.Scene.ID,
		})
		require.NoError(t, err)

		// Get scene actors
		result, err := e.GetSceneActors(ctx, GetSceneActorsRequest{GameID: gameID,
			SceneID: scene.Scene.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Actors, 1)
		assert.Equal(t, pc.Actor.ID, result.Actors[0].ActorID)
	})
}

func TestAddSceneConnection(t *testing.T) {
	t.Run("adds scene connection", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for scenes",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create scenes
		scene1, err := e.CreateScene(ctx, CreateSceneRequest{GameID: gameID,
			Name:        "Room A",
			Description: "First room",
			SceneType:   model.SceneTypeIndoor,
		})
		require.NoError(t, err)

		scene2, err := e.CreateScene(ctx, CreateSceneRequest{GameID: gameID,
			Name:        "Room B",
			Description: "Second room",
			SceneType:   model.SceneTypeIndoor,
		})
		require.NoError(t, err)

		// Add connection
		err = e.AddSceneConnection(ctx, AddSceneConnectionRequest{
			GameID:        gameID,
			SceneID:       scene1.Scene.ID,
			TargetSceneID: scene2.Scene.ID,
			Description:   "A wooden door",
			Locked:        false,
			DC:            10,
		})

		require.NoError(t, err)

		// Verify connection
		result, err := e.GetScene(ctx, GetSceneRequest{GameID: gameID,
			SceneID: scene1.Scene.ID,
		})
		require.NoError(t, err)
		assert.Len(t, result.Connections, 1)
		assert.Equal(t, scene2.Scene.ID, result.Connections[0].TargetSceneID)
	})
}
