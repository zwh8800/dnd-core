package testsuite

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zwh8800/dnd-core/pkg/engine"
	"github.com/zwh8800/dnd-core/pkg/model"
)

func TestRestingSystem(t *testing.T) {
	t.Run("short rest", func(t *testing.T) {
		e := engine.NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, engine.NewGameRequest{
			Name:        "Campfire Tales",
			Description: "Resting at camp",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		fighter, err := e.CreatePC(ctx, engine.CreatePCRequest{
			GameID: gameID,
			PC: &engine.PlayerCharacterInput{
				Name:  "Conan",
				Race:  "Human",
				Class: "Fighter",
				Level: 3,
				AbilityScores: engine.AbilityScoresInput{
					Strength:     16,
					Dexterity:    14,
					Constitution: 16,
					Intelligence: 10,
					Wisdom:       12,
					Charisma:     10,
				},
			},
		})
		require.NoError(t, err)

		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "Resting")
		require.NoError(t, err)

		restResult, err := e.ShortRest(ctx, engine.ShortRestRequest{
			GameID:   gameID,
			ActorIDs: []model.ID{fighter.Actor.ID},
		})
		require.NoError(t, err)
		require.NotNil(t, restResult)
		require.GreaterOrEqual(t, len(restResult.ActorResults), 1)
		t.Logf("Short rest: %s", restResult.Message)
	})

	t.Run("long rest", func(t *testing.T) {
		e := engine.NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, engine.NewGameRequest{
			Name:        "Inn Stay",
			Description: "Resting at an inn",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		cleric, err := e.CreatePC(ctx, engine.CreatePCRequest{
			GameID: gameID,
			PC: &engine.PlayerCharacterInput{
				Name:  "Healer",
				Race:  "Dwarf",
				Class: "Cleric",
				Level: 2,
				AbilityScores: engine.AbilityScoresInput{
					Strength:     14,
					Dexterity:    10,
					Constitution: 14,
					Intelligence: 10,
					Wisdom:       16,
					Charisma:     12,
				},
			},
		})
		require.NoError(t, err)

		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "Resting")
		require.NoError(t, err)

		restResult, err := e.StartLongRest(ctx, engine.StartLongRestRequest{
			GameID:   gameID,
			ActorIDs: []model.ID{cleric.Actor.ID},
		})
		require.NoError(t, err)
		require.NotNil(t, restResult)
		t.Logf("Long rest started: %s", restResult.Message)

		endResult, err := e.EndLongRest(ctx, engine.EndLongRestRequest{
			GameID: gameID,
		})
		require.NoError(t, err)
		require.NotNil(t, endResult)
		require.GreaterOrEqual(t, len(endResult.ActorResults), 1)

		actorResult := endResult.ActorResults[0]
		assert.GreaterOrEqual(t, actorResult.HPRecovered, 0)
		t.Logf("Long rest: HP recovered=%d, SpellSlots=%v",
			actorResult.HPRecovered, actorResult.SpellSlotsRestored)
	})

	t.Run("party rest", func(t *testing.T) {
		e := engine.NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, engine.NewGameRequest{
			Name:        "Group Rest",
			Description: "Whole party resting together",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		pc1, err := e.CreatePC(ctx, engine.CreatePCRequest{
			GameID: gameID,
			PC: &engine.PlayerCharacterInput{
				Name:  "Warrior",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: engine.AbilityScoresInput{
					Strength: 16, Dexterity: 14, Constitution: 14,
					Intelligence: 10, Wisdom: 12, Charisma: 10,
				},
			},
		})
		require.NoError(t, err)

		pc2, err := e.CreatePC(ctx, engine.CreatePCRequest{
			GameID: gameID,
			PC: &engine.PlayerCharacterInput{
				Name:  "Mage",
				Race:  "Elf",
				Class: "Wizard",
				Level: 1,
				AbilityScores: engine.AbilityScoresInput{
					Strength: 8, Dexterity: 14, Constitution: 10,
					Intelligence: 16, Wisdom: 12, Charisma: 12,
				},
			},
		})
		require.NoError(t, err)

		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "Party rests")
		require.NoError(t, err)

		restResult, err := e.StartLongRest(ctx, engine.StartLongRestRequest{
			GameID:   gameID,
			ActorIDs: []model.ID{pc1.Actor.ID, pc2.Actor.ID},
		})
		require.NoError(t, err)
		require.NotNil(t, restResult)

		endResult, err := e.EndLongRest(ctx, engine.EndLongRestRequest{
			GameID: gameID,
		})
		require.NoError(t, err)
		require.NotNil(t, endResult)
		require.Equal(t, 2, len(endResult.ActorResults))

		t.Logf("Party rested: %s recovered %d HP, %s recovered %d HP",
			pc1.Actor.Name, endResult.ActorResults[0].HPRecovered,
			pc2.Actor.Name, endResult.ActorResults[1].HPRecovered)
	})
}
