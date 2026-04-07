package engine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zwh8800/dnd-core/internal/model"
)

func TestNewGame(t *testing.T) {
	t.Run("creates new game successfully", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Epic Campaign",
			Description: "A grand adventure awaits",
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.Game)
		assert.Equal(t, "Epic Campaign", result.Game.Name)
		assert.Equal(t, "A grand adventure awaits", result.Game.Description)
		assert.NotEmpty(t, result.Game.ID)
		assert.Equal(t, model.PhaseCharacterCreation, result.Game.Phase)
	})

	t.Run("creates game with default values", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.NewGame(ctx, NewGameRequest{
			Name: "Simple Game",
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.Game)
		assert.Equal(t, "Simple Game", result.Game.Name)
		assert.Empty(t, result.Game.Description)
	})
}

func TestLoadGame(t *testing.T) {
	t.Run("loads existing game", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		// Create a game first
		createResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Game to Load",
			Description: "Will be loaded later",
		})
		require.NoError(t, err)
		gameID := createResult.Game.ID

		// Load the game
		result, err := e.LoadGame(ctx, LoadGameRequest{
			GameID: gameID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.Game)
		assert.Equal(t, gameID, result.Game.ID)
		assert.Equal(t, "Game to Load", result.Game.Name)
	})

	t.Run("returns error for non-existent game", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		_, err := e.LoadGame(ctx, LoadGameRequest{
			GameID: model.NewID(),
		})

		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})
}

func TestSaveGame(t *testing.T) {
	t.Run("saves game successfully", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		// Create a game
		createResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Game to Save",
			Description: "Will be saved",
		})
		require.NoError(t, err)
		gameID := createResult.Game.ID

		// Save the game
		err = e.SaveGame(ctx, SaveGameRequest{
			GameID: gameID,
		})

		require.NoError(t, err)

		// Verify by loading
		loadResult, err := e.LoadGame(ctx, LoadGameRequest{
			GameID: gameID,
		})
		require.NoError(t, err)
		assert.Equal(t, "Game to Save", loadResult.Game.Name)
	})

	t.Run("returns error for non-existent game", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		err := e.SaveGame(ctx, SaveGameRequest{
			GameID: model.NewID(),
		})

		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})
}

func TestDeleteGame(t *testing.T) {
	t.Run("deletes game successfully", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		// Create a game
		createResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Game to Delete",
			Description: "Will be deleted",
		})
		require.NoError(t, err)
		gameID := createResult.Game.ID

		// Delete the game
		err = e.DeleteGame(ctx, DeleteGameRequest{
			GameID: gameID,
		})

		require.NoError(t, err)

		// Verify game is deleted
		_, err = e.LoadGame(ctx, LoadGameRequest{
			GameID: gameID,
		})
		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("returns error for non-existent game", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		err := e.DeleteGame(ctx, DeleteGameRequest{
			GameID: model.NewID(),
		})

		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})
}

func TestListGames(t *testing.T) {
	t.Run("lists all games", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		// Create multiple games
		_, err := e.NewGame(ctx, NewGameRequest{
			Name: "Game One",
		})
		require.NoError(t, err)

		_, err = e.NewGame(ctx, NewGameRequest{
			Name: "Game Two",
		})
		require.NoError(t, err)

		_, err = e.NewGame(ctx, NewGameRequest{
			Name: "Game Three",
		})
		require.NoError(t, err)

		// List games
		result, err := e.ListGames(ctx, ListGamesRequest{})

		require.NoError(t, err)
		assert.Len(t, result, 3)
	})

	t.Run("returns empty list when no games exist", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.ListGames(ctx, ListGamesRequest{})

		require.NoError(t, err)
		assert.Len(t, result, 0)
	})
}

func TestGameInfo(t *testing.T) {
	t.Run("game info contains correct counts", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		// Create a game
		createResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Campaign Game",
			Description: "A game with characters",
		})
		require.NoError(t, err)
		gameID := createResult.Game.ID

		// Add some actors
		_, err = e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Hero 1",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength: 15, Dexterity: 14, Constitution: 13,
					Intelligence: 10, Wisdom: 12, Charisma: 8,
				},
			},
		})
		require.NoError(t, err)

		_, err = e.CreateNPC(ctx, CreateNPCRequest{
			GameID: gameID,
			NPC: &NPCInput{
				Name:        "Guard",
				Description: "A city guard",
				Size:        model.SizeMedium,
				Speed:       30,
				AbilityScores: AbilityScoresInput{
					Strength: 14, Dexterity: 12, Constitution: 14,
					Intelligence: 10, Wisdom: 12, Charisma: 10,
				},
			},
		})
		require.NoError(t, err)

		// Load and check counts
		result, err := e.LoadGame(ctx, LoadGameRequest{
			GameID: gameID,
		})
		require.NoError(t, err)
		assert.Equal(t, 1, result.Game.PCCount)
		assert.Equal(t, 1, result.Game.NPCCount)
		assert.Equal(t, 0, result.Game.EnemyCount)
		assert.Equal(t, 0, result.Game.CompanionCount)
		assert.Equal(t, 0, result.Game.SceneCount)
		assert.False(t, result.Game.InCombat)
	})
}
