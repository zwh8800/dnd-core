package engine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zwh8800/dnd-core/pkg/model"
)

func TestCurseActor(t *testing.T) {
	t.Run("curses actor successfully", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.CurseActor(ctx, CurseActorRequest{
			GameID:   gameResult.Game.ID,
			TargetID: model.NewID(),
			CurseID:  "bestow-curse",
			Source:   "spell",
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.CurseInstance)
		assert.Equal(t, "bestow-curse", result.CurseInstance.CurseID)
		assert.Equal(t, "spell", result.CurseInstance.Source)
		assert.True(t, result.CurseInstance.IsPermanent)
		assert.Equal(t, "permanent", result.CurseInstance.RemainingDuration)
		assert.Contains(t, result.Message, "bestow-curse")
	})

	t.Run("curse with different source", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.CurseActor(ctx, CurseActorRequest{
			GameID:   gameResult.Game.ID,
			TargetID: model.NewID(),
			CurseID:  "lycanthropy",
			Source:   "trap",
		})

		require.NoError(t, err)
		assert.Equal(t, "lycanthropy", result.CurseInstance.CurseID)
		assert.Contains(t, result.Message, "trap")
	})

	t.Run("curse with empty source", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.CurseActor(ctx, CurseActorRequest{
			GameID:   gameResult.Game.ID,
			TargetID: model.NewID(),
			CurseID:  "mummy-rot",
			Source:   "",
		})

		require.NoError(t, err)
		assert.Equal(t, "", result.CurseInstance.Source)
	})

	t.Run("returns error when game not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.CurseActor(ctx, CurseActorRequest{
			GameID:   model.NewID(),
			TargetID: model.NewID(),
			CurseID:  "bestow-curse",
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("multiple curses can be applied", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		curses := []string{"curse1", "curse2", "curse3"}
		for _, curseID := range curses {
			result, err := e.CurseActor(ctx, CurseActorRequest{
				GameID:   gameResult.Game.ID,
				TargetID: model.NewID(),
				CurseID:  curseID,
				Source:   "test",
			})
			require.NoError(t, err)
			assert.Equal(t, curseID, result.CurseInstance.CurseID)
		}
	})
}

func TestRemoveCurse(t *testing.T) {
	t.Run("removes curse successfully", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.RemoveCurse(ctx, RemoveCurseRequest{
			GameID:   gameResult.Game.ID,
			TargetID: model.NewID(),
			CurseID:  "bestow-curse",
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.Success)
		assert.Contains(t, result.Message, "bestow-curse")
	})

	t.Run("removes curse with different ID", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.RemoveCurse(ctx, RemoveCurseRequest{
			GameID:   gameResult.Game.ID,
			TargetID: model.NewID(),
			CurseID:  "lycanthropy",
		})

		require.NoError(t, err)
		assert.True(t, result.Success)
		assert.Contains(t, result.Message, "lycanthropy")
	})

	t.Run("removing non-existent curse succeeds (simplified)", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.RemoveCurse(ctx, RemoveCurseRequest{
			GameID:   gameResult.Game.ID,
			TargetID: model.NewID(),
			CurseID:  "nonexistent-curse",
		})

		require.NoError(t, err)
		assert.True(t, result.Success)
	})

	t.Run("returns error when game not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.RemoveCurse(ctx, RemoveCurseRequest{
			GameID:   model.NewID(),
			TargetID: model.NewID(),
			CurseID:  "bestow-curse",
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("multiple removes return consistent results", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		for i := 0; i < 3; i++ {
			result, err := e.RemoveCurse(ctx, RemoveCurseRequest{
				GameID:   gameResult.Game.ID,
				TargetID: model.NewID(),
				CurseID:  "curse",
			})
			require.NoError(t, err)
			assert.True(t, result.Success)
		}
	})
}

func TestGetCurses(t *testing.T) {
	t.Run("returns empty curses list", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.GetCurses(ctx, GetCursesRequest{
			GameID:   gameResult.Game.ID,
			TargetID: model.NewID(),
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Empty(t, result.Curses)
		assert.Equal(t, 0, result.Count)
	})

	t.Run("returns curses for different targets", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		for i := 0; i < 3; i++ {
			result, err := e.GetCurses(ctx, GetCursesRequest{
				GameID:   gameResult.Game.ID,
				TargetID: model.NewID(),
			})
			require.NoError(t, err)
			assert.Empty(t, result.Curses)
			assert.Equal(t, 0, result.Count)
		}
	})

	t.Run("returns error when game not found", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.GetCurses(ctx, GetCursesRequest{
			GameID:   model.NewID(),
			TargetID: model.NewID(),
		})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("get curses uses read lock", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		done := make(chan bool, 2)
		go func() {
			_, err := e.GetCurses(ctx, GetCursesRequest{
				GameID:   gameResult.Game.ID,
				TargetID: model.NewID(),
			})
			assert.NoError(t, err)
			done <- true
		}()
		go func() {
			_, err := e.GetCurses(ctx, GetCursesRequest{
				GameID:   gameResult.Game.ID,
				TargetID: model.NewID(),
			})
			assert.NoError(t, err)
			done <- true
		}()

		<-done
		<-done
	})

	t.Run("curses list is initialized not nil", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{Name: "Test", Description: "Test"})
		require.NoError(t, err)

		result, err := e.GetCurses(ctx, GetCursesRequest{
			GameID:   gameResult.Game.ID,
			TargetID: model.NewID(),
		})

		require.NoError(t, err)
		assert.NotNil(t, result.Curses)
	})
}
