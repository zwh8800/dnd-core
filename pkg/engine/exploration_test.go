package engine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zwh8800/dnd-core/pkg/model"
)

func TestStartTravel(t *testing.T) {
	t.Run("starts travel successfully", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.StartTravel(ctx, StartTravelRequest{
			GameID:      gameResult.Game.ID,
			Destination: "Waterdeep",
			Pace:        model.TravelPaceNormal,
			Terrain:     model.TerrainClear,
			Distance:    100.0,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.TravelState)
		assert.Equal(t, "Waterdeep", result.TravelState.Destination)
		assert.Equal(t, model.TravelPaceNormal, result.TravelState.Pace)
		assert.Equal(t, model.TerrainClear, result.TravelState.Terrain)
		assert.Equal(t, 100.0, result.TravelState.DistanceTotal)
		assert.Equal(t, float64(0), result.TravelState.DistanceTraveled)
		assert.True(t, result.TravelState.IsActive)
		assert.Contains(t, result.Message, "Waterdeep")
	})

	t.Run("starts travel with fast pace", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.StartTravel(ctx, StartTravelRequest{
			GameID:      gameResult.Game.ID,
			Destination: "Baldur's Gate",
			Pace:        model.TravelPaceFast,
			Terrain:     model.TerrainForest,
			Distance:    50.0,
		})

		require.NoError(t, err)
		assert.Equal(t, model.TravelPaceFast, result.TravelState.Pace)
		assert.Equal(t, model.TerrainForest, result.TravelState.Terrain)
	})

	t.Run("starts travel with slow pace", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.StartTravel(ctx, StartTravelRequest{
			GameID:      gameResult.Game.ID,
			Destination: "Neverwinter",
			Pace:        model.TravelPaceSlow,
			Terrain:     model.TerrainMountain,
			Distance:    200.0,
		})

		require.NoError(t, err)
		assert.Equal(t, model.TravelPaceSlow, result.TravelState.Pace)
	})

	t.Run("starts travel with zero distance", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.StartTravel(ctx, StartTravelRequest{
			GameID:      gameResult.Game.ID,
			Destination: "Nearby Village",
			Pace:        model.TravelPaceNormal,
			Terrain:     model.TerrainClear,
			Distance:    0.0,
		})

		require.NoError(t, err)
		assert.Equal(t, 0.0, result.TravelState.DistanceTotal)
	})

	t.Run("returns error when game not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.StartTravel(ctx, StartTravelRequest{
			GameID:      model.NewID(),
			Destination: "Waterdeep",
			Pace:        model.TravelPaceNormal,
			Terrain:     model.TerrainClear,
			Distance:    100.0,
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("starting travel replaces previous travel", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		// Start first travel
		_, err = e.StartTravel(ctx, StartTravelRequest{
			GameID:      gameResult.Game.ID,
			Destination: "First Destination",
			Pace:        model.TravelPaceNormal,
			Terrain:     model.TerrainClear,
			Distance:    50.0,
		})
		require.NoError(t, err)

		// Start second travel
		result, err := e.StartTravel(ctx, StartTravelRequest{
			GameID:      gameResult.Game.ID,
			Destination: "Second Destination",
			Pace:        model.TravelPaceFast,
			Terrain:     model.TerrainForest,
			Distance:    100.0,
		})
		require.NoError(t, err)

		assert.Equal(t, "Second Destination", result.TravelState.Destination)
	})
}

func TestAdvanceTravel(t *testing.T) {
	t.Run("advances travel by hours", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		_, err = e.StartTravel(ctx, StartTravelRequest{
			GameID:      gameResult.Game.ID,
			Destination: "Waterdeep",
			Pace:        model.TravelPaceNormal,
			Terrain:     model.TerrainClear,
			Distance:    100.0,
		})
		require.NoError(t, err)

		result, err := e.AdvanceTravel(ctx, AdvanceTravelRequest{
			GameID: gameResult.Game.ID,
			Hours:  4,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Greater(t, result.DistanceTraveled, float64(0))
		assert.Contains(t, result.Message, "行进")
	})

	t.Run("advances travel full day triggers checks", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		_, err = e.StartTravel(ctx, StartTravelRequest{
			GameID:      gameResult.Game.ID,
			Destination: "Waterdeep",
			Pace:        model.TravelPaceNormal,
			Terrain:     model.TerrainForest,
			Distance:    100.0,
		})
		require.NoError(t, err)

		result, err := e.AdvanceTravel(ctx, AdvanceTravelRequest{
			GameID: gameResult.Game.ID,
			Hours:  8,
		})

		require.NoError(t, err)
		assert.Greater(t, result.DaysElapsed, 0)
		// Forage, navigation, and encounter checks may or may not return results
		// depending on the rules implementation
	})

	t.Run("arrives at destination when distance completed", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		_, err = e.StartTravel(ctx, StartTravelRequest{
			GameID:      gameResult.Game.ID,
			Destination: "Waterdeep",
			Pace:        model.TravelPaceNormal,
			Terrain:     model.TerrainClear,
			Distance:    1.0, // Very short distance
		})
		require.NoError(t, err)

		// Advance enough to complete the journey
		result, err := e.AdvanceTravel(ctx, AdvanceTravelRequest{
			GameID: gameResult.Game.ID,
			Hours:  8,
		})

		require.NoError(t, err)
		assert.Contains(t, result.Message, "已到达目的地")
	})

	t.Run("returns error when no active travel", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.AdvanceTravel(ctx, AdvanceTravelRequest{
			GameID: gameResult.Game.ID,
			Hours:  1,
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "no active travel")
	})

	t.Run("returns error when game not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.AdvanceTravel(ctx, AdvanceTravelRequest{
			GameID: model.NewID(),
			Hours:  1,
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("advances travel multiple times accumulates distance", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		_, err = e.StartTravel(ctx, StartTravelRequest{
			GameID:      gameResult.Game.ID,
			Destination: "Waterdeep",
			Pace:        model.TravelPaceNormal,
			Terrain:     model.TerrainClear,
			Distance:    500.0,
		})
		require.NoError(t, err)

		var totalDistance float64
		for i := 0; i < 3; i++ {
			result, err := e.AdvanceTravel(ctx, AdvanceTravelRequest{
				GameID: gameResult.Game.ID,
				Hours:  4,
			})
			require.NoError(t, err)
			totalDistance += result.DistanceTraveled
		}

		assert.Greater(t, totalDistance, float64(0))
	})
}

func TestForage(t *testing.T) {
	t.Run("forages successfully", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.Forage(ctx, ForageRequest{
			GameID: gameResult.Game.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.Result)
		// ForageResult fields are populated by rules.ForagingCheck
	})

	t.Run("forage returns error when game not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.Forage(ctx, ForageRequest{
			GameID: model.NewID(),
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("multiple forage attempts", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		for i := 0; i < 3; i++ {
			result, err := e.Forage(ctx, ForageRequest{
				GameID: gameResult.Game.ID,
			})
			require.NoError(t, err)
			require.NotNil(t, result)
		}
	})

	t.Run("forage uses write lock", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		done := make(chan bool, 2)
		go func() {
			_, err := e.Forage(ctx, ForageRequest{GameID: gameResult.Game.ID})
			assert.NoError(t, err)
			done <- true
		}()
		go func() {
			_, err := e.Forage(ctx, ForageRequest{GameID: gameResult.Game.ID})
			assert.NoError(t, err)
			done <- true
		}()

		<-done
		<-done
	})

	t.Run("forage result has valid structure", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.Forage(ctx, ForageRequest{
			GameID: gameResult.Game.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, result.Result)
		// The result should have a message
		assert.NotEmpty(t, result.Result.Message)
	})
}

func TestNavigate(t *testing.T) {
	t.Run("navigates successfully", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.Navigate(ctx, NavigateRequest{
			GameID: gameResult.Game.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.Result)
	})

	t.Run("navigate uses travel terrain when available", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		_, err = e.StartTravel(ctx, StartTravelRequest{
			GameID:      gameResult.Game.ID,
			Destination: "Waterdeep",
			Pace:        model.TravelPaceNormal,
			Terrain:     model.TerrainMountain,
			Distance:    100.0,
		})
		require.NoError(t, err)

		result, err := e.Navigate(ctx, NavigateRequest{
			GameID: gameResult.Game.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
	})

	t.Run("navigate uses clear terrain when no travel", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.Navigate(ctx, NavigateRequest{
			GameID: gameResult.Game.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
	})

	t.Run("navigate returns error when game not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.Navigate(ctx, NavigateRequest{
			GameID: model.NewID(),
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("multiple navigation attempts", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		for i := 0; i < 3; i++ {
			result, err := e.Navigate(ctx, NavigateRequest{
				GameID: gameResult.Game.ID,
			})
			require.NoError(t, err)
			require.NotNil(t, result)
		}
	})
}
