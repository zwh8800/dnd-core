package testsuite

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zwh8800/dnd-core/pkg/engine"
	"github.com/zwh8800/dnd-core/pkg/model"
)

func TestCombatEncounters(t *testing.T) {
	t.Run("full combat round with multiple participants", func(t *testing.T) {
		e := engine.NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, engine.NewGameRequest{
			Name:        "Arena Battle",
			Description: "A gladiatorial combat",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		fighter, err := e.CreatePC(ctx, engine.CreatePCRequest{
			GameID: gameID,
			PC: &engine.PlayerCharacterInput{
				Name:  "Marcus",
				Race:  "Human",
				Class: "Fighter",
				Level: 5,
				AbilityScores: engine.AbilityScoresInput{
					Strength:     18,
					Dexterity:    14,
					Constitution: 16,
					Intelligence: 10,
					Wisdom:       12,
					Charisma:     10,
				},
			},
		})
		require.NoError(t, err)

		rogue, err := e.CreatePC(ctx, engine.CreatePCRequest{
			GameID: gameID,
			PC: &engine.PlayerCharacterInput{
				Name:  "Shadow",
				Race:  "Halfling",
				Class: "Rogue",
				Level: 5,
				AbilityScores: engine.AbilityScoresInput{
					Strength:     10,
					Dexterity:    18,
					Constitution: 12,
					Intelligence: 14,
					Wisdom:       12,
					Charisma:     14,
				},
			},
		})
		require.NoError(t, err)

		orc, err := e.CreateEnemy(ctx, engine.CreateEnemyRequest{
			GameID: gameID,
			Enemy: &engine.EnemyInput{
				Name:        "Orc Brute",
				Description: "A fearsome orc warrior",
				Size:        model.SizeMedium,
				Speed:       30,
				AbilityScores: engine.AbilityScoresInput{
					Strength:     18,
					Dexterity:    12,
					Constitution: 16,
					Intelligence: 6,
					Wisdom:       10,
					Charisma:     12,
				},
				ChallengeRating: 1,
				HitPoints:       43,
				ArmorClass:      13,
			},
		})
		require.NoError(t, err)

		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "Begin combat")
		require.NoError(t, err)

		combatResult, err := e.StartCombat(ctx, engine.StartCombatRequest{
			GameID: gameID,
			ParticipantIDs: []model.ID{
				fighter.Actor.ID,
				rogue.Actor.ID,
				orc.Actor.ID,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, combatResult.Combat)
		assert.Equal(t, model.CombatStatusActive, combatResult.Combat.Status)
		assert.Equal(t, 3, len(combatResult.Combat.Initiative))

		t.Logf("Combat started - Round %d, %d combatants", combatResult.Combat.Round, len(combatResult.Combat.Initiative))

		currentTurn, err := e.GetCurrentTurn(ctx, engine.GetCurrentTurnRequest{
			GameID: gameID,
		})
		require.NoError(t, err)
		require.NotNil(t, currentTurn)
		t.Logf("Current turn: %s", currentTurn.ActorName)

		t.Run("execute attack action", func(t *testing.T) {
			attackResult, err := e.ExecuteAttack(ctx, engine.ExecuteAttackRequest{
				GameID:     gameID,
				AttackerID: fighter.Actor.ID,
				TargetID:   orc.Actor.ID,
				Attack: engine.AttackInput{
					IsUnarmed: false,
				},
			})
			require.NoError(t, err)
			require.NotNil(t, attackResult)
			require.NotNil(t, attackResult.AttackResult)
			t.Logf("Attack: Roll=%d, Hit=%v, Critical=%v",
				attackResult.AttackResult.Roll.Total,
				attackResult.AttackResult.Hit,
				attackResult.AttackResult.IsCritical)
		})

		t.Run("next turn advances", func(t *testing.T) {
			nextResult, err := e.NextTurn(ctx, engine.NextTurnRequest{
				GameID: gameID,
			})
			require.NoError(t, err)
			require.NotNil(t, nextResult)

			summary, err := e.GetCombatSummary(ctx, gameID)
			require.NoError(t, err)
			require.NotNil(t, summary)
			t.Logf("Combat: Round %d", summary.Round)
		})

		t.Run("end combat", func(t *testing.T) {
			err := e.EndCombat(ctx, engine.EndCombatRequest{
				GameID: gameID,
			})
			require.NoError(t, err)
		})
	})
}
