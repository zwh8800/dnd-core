package engine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoll(t *testing.T) {
	t.Run("rolls dice successfully", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.Roll(ctx, RollRequest{
			Expression: "1d20",
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, result.Total, 1)
		assert.LessOrEqual(t, result.Total, 20)
	})

	t.Run("rolls multiple dice", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.Roll(ctx, RollRequest{
			Expression: "2d6",
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, result.Total, 2)
		assert.LessOrEqual(t, result.Total, 12)
	})

	t.Run("rolls with modifier", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.Roll(ctx, RollRequest{
			Expression: "1d20",
			Modifier:   5,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, result.Total, 6)
		assert.LessOrEqual(t, result.Total, 25)
	})
}

func TestRollAdvantage(t *testing.T) {
	t.Run("rolls with advantage", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.RollAdvantage(ctx, RollAdvantageRequest{
			Modifier: 0,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, result.Total, 1)
		assert.LessOrEqual(t, result.Total, 20)
	})
}

func TestRollDisadvantage(t *testing.T) {
	t.Run("rolls with disadvantage", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.RollDisadvantage(ctx, RollDisadvantageRequest{
			Modifier: 0,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, result.Total, 1)
		assert.LessOrEqual(t, result.Total, 20)
	})
}

func TestRollAbility(t *testing.T) {
	t.Run("rolls ability check", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.RollAbility(ctx, RollAbilityRequest{})

		require.NoError(t, err)
		require.NotNil(t, result)
		// 4d6 drop lowest, min is 3 (1+1+1+0dropped)
		assert.GreaterOrEqual(t, result.Total, 3)
		// max is 18 (6+6+6+0dropped)
		assert.LessOrEqual(t, result.Total, 18)
	})
}

func TestRollHitDice(t *testing.T) {
	t.Run("rolls hit dice", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.RollHitDice(ctx, RollHitDiceRequest{
			DiceType: 8,
			Modifier: 2,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		// d8 + modifier
		assert.GreaterOrEqual(t, result.Total, 3)
		assert.LessOrEqual(t, result.Total, 10)
	})
}
