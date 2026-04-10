package engine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zwh8800/dnd-core/pkg/model"
)

// setupTestPC 创建测试用的游戏和角色
func setupTestPC(t *testing.T, e *Engine, name string, strength, constitution int) (model.ID, model.ID) {
	t.Helper()
	ctx := context.Background()

	gameResult, err := e.NewGame(ctx, NewGameRequest{
		Name:        "Test Game",
		Description: "Test movement and exploration",
	})
	require.NoError(t, err)
	gameID := gameResult.Game.ID

	pcResult, err := e.CreatePC(ctx, CreatePCRequest{
		GameID: gameID,
		PC: &PlayerCharacterInput{
			Name:  name,
			Race:  "Human",
			Class: "Fighter",
			Level: 1,
			AbilityScores: AbilityScoresInput{
				Strength:     strength,
				Dexterity:    14,
				Constitution: constitution,
				Intelligence: 12,
				Wisdom:       10,
				Charisma:     8,
			},
		},
	})
	require.NoError(t, err)

	return gameID, pcResult.Actor.ID
}

// ============================================================================
// PerformJump Tests
// ============================================================================

func TestPerformJump(t *testing.T) {
	t.Run("long jump with running start", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameID, actorID := setupTestPC(t, e, "Jumper", 15, 13)

		// Switch to exploration phase
		_, err := e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		result, err := e.PerformJump(ctx, PerformJumpRequest{
			GameID:          gameID,
			ActorID:         actorID,
			JumpType:        model.JumpTypeLong,
			HasRunningStart: true,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, model.JumpTypeLong, result.JumpType)
		assert.Equal(t, 15, result.Distance) // 力量值 15 = 15 尺
		assert.True(t, result.HasRunningStart)
		assert.Equal(t, 15, result.Strength)
		assert.Equal(t, 2, result.StrengthMod) // (15-10)/2 = 2
		assert.Contains(t, result.Message, "跳远")
		assert.Contains(t, result.Message, "有助跑")
	})

	t.Run("long jump without running start", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameID, actorID := setupTestPC(t, e, "Jumper", 16, 13)

		_, err := e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		result, err := e.PerformJump(ctx, PerformJumpRequest{
			GameID:          gameID,
			ActorID:         actorID,
			JumpType:        model.JumpTypeLong,
			HasRunningStart: false,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 8, result.Distance) // 16/2 = 8 尺
		assert.False(t, result.HasRunningStart)
	})

	t.Run("high jump with running start", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameID, actorID := setupTestPC(t, e, "Jumper", 18, 13)

		_, err := e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		result, err := e.PerformJump(ctx, PerformJumpRequest{
			GameID:          gameID,
			ActorID:         actorID,
			JumpType:        model.JumpTypeHigh,
			HasRunningStart: true,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, model.JumpTypeHigh, result.JumpType)
		// 跳高: 3 + 力量修正(4) = 7 尺
		assert.Equal(t, 7, result.Distance)
		assert.Equal(t, 4, result.StrengthMod) // (18-10)/2 = 4
	})

	t.Run("high jump without running start", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameID, actorID := setupTestPC(t, e, "Jumper", 14, 13)

		_, err := e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		result, err := e.PerformJump(ctx, PerformJumpRequest{
			GameID:          gameID,
			ActorID:         actorID,
			JumpType:        model.JumpTypeHigh,
			HasRunningStart: false,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		// 跳高: (3 + 力量修正(2)) / 2 = 2 尺
		assert.Equal(t, 2, result.Distance)
	})

	t.Run("jump with low strength", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameID, actorID := setupTestPC(t, e, "Weak Jumper", 6, 13)

		_, err := e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		result, err := e.PerformJump(ctx, PerformJumpRequest{
			GameID:          gameID,
			ActorID:         actorID,
			JumpType:        model.JumpTypeHigh,
			HasRunningStart: true,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		// 3 + (-2) = 1 尺
		assert.Equal(t, 1, result.Distance)
		assert.Equal(t, -2, result.StrengthMod)
	})

	t.Run("jump with non-existent actor fails", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test",
		})
		require.NoError(t, err)

		_, err = e.SetPhase(ctx, gameResult.Game.ID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		_, err = e.PerformJump(ctx, PerformJumpRequest{
			GameID:  gameResult.Game.ID,
			ActorID: "non-existent",
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("jump in wrong phase fails", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameID, actorID := setupTestPC(t, e, "Jumper", 15, 13)
		// Stay in character creation phase

		_, err := e.PerformJump(ctx, PerformJumpRequest{
			GameID:  gameID,
			ActorID: actorID,
		})

		assert.Error(t, err) // Phase not allowed
	})
}

// ============================================================================
// ApplyFallDamage Tests
// ============================================================================

func TestApplyFallDamage(t *testing.T) {
	t.Run("fall damage 30 feet", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameID, actorID := setupTestPC(t, e, "Faller", 15, 13)

		_, err := e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		// Set initial HP
		hp := 50
		err = e.UpdateActor(ctx, UpdateActorRequest{
			GameID:  gameID,
			ActorID: actorID,
			Update: ActorUpdate{
				HitPoints: &HitPointUpdate{
					Current: &hp,
				},
			},
		})
		require.NoError(t, err)

		result, err := e.ApplyFallDamage(ctx, ApplyFallDamageRequest{
			GameID:       gameID,
			ActorID:      actorID,
			FallDistance: 30,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 30, result.FallDistance)
		assert.Equal(t, 3, result.DamageDice)           // 30/10 = 3d6
		assert.GreaterOrEqual(t, result.DamageTaken, 3) // 最小 3
		assert.LessOrEqual(t, result.DamageTaken, 18)   // 最大 3*6 = 18
		assert.Equal(t, 120, result.MaxPossible)        // 20d6 最大 = 120
		assert.Contains(t, result.Message, "跌落")
		assert.Contains(t, result.Message, "3d6")
	})

	t.Run("fall damage no damage under 10 feet", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameID, actorID := setupTestPC(t, e, "Faller", 15, 13)

		_, err := e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		result, err := e.ApplyFallDamage(ctx, ApplyFallDamageRequest{
			GameID:       gameID,
			ActorID:      actorID,
			FallDistance: 5,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 0, result.DamageDice)
		assert.Equal(t, 0, result.DamageTaken)
		assert.Contains(t, result.Message, "不足以造成伤害")
	})

	t.Run("fall damage maximum 20d6", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameID, actorID := setupTestPC(t, e, "Faller", 15, 13)

		_, err := e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		// Set high HP to survive
		hp := 200
		err = e.UpdateActor(ctx, UpdateActorRequest{
			GameID:  gameID,
			ActorID: actorID,
			Update: ActorUpdate{
				HitPoints: &HitPointUpdate{
					Current: &hp,
				},
			},
		})
		require.NoError(t, err)

		result, err := e.ApplyFallDamage(ctx, ApplyFallDamageRequest{
			GameID:       gameID,
			ActorID:      actorID,
			FallDistance: 300, // 应该限制为 20d6
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 20, result.DamageDice)           //  capped at 20
		assert.GreaterOrEqual(t, result.DamageTaken, 20) // 最小 20
		assert.LessOrEqual(t, result.DamageTaken, 120)   // 最大 20*6 = 120
	})

	t.Run("fall damage reduces HP", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameID, actorID := setupTestPC(t, e, "Faller", 15, 13)

		_, err := e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		// Set known HP
		hp := 30
		err = e.UpdateActor(ctx, UpdateActorRequest{
			GameID:  gameID,
			ActorID: actorID,
			Update: ActorUpdate{
				HitPoints: &HitPointUpdate{
					Current: &hp,
				},
			},
		})
		require.NoError(t, err)

		result, err := e.ApplyFallDamage(ctx, ApplyFallDamageRequest{
			GameID:       gameID,
			ActorID:      actorID,
			FallDistance: 20, // 2d6
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 2, result.DamageDice)
		assert.Greater(t, result.DamageTaken, 0)
		// HP should be reduced
		assert.Less(t, result.CurrentHP, 30)
	})

	t.Run("fall damage does not go below 0 HP", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameID, actorID := setupTestPC(t, e, "Faller", 15, 13)

		_, err := e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		// Set low HP
		hp := 5
		err = e.UpdateActor(ctx, UpdateActorRequest{
			GameID:  gameID,
			ActorID: actorID,
			Update: ActorUpdate{
				HitPoints: &HitPointUpdate{
					Current: &hp,
				},
			},
		})
		require.NoError(t, err)

		result, err := e.ApplyFallDamage(ctx, ApplyFallDamageRequest{
			GameID:       gameID,
			ActorID:      actorID,
			FallDistance: 50, // 5d6, likely to exceed 5 HP
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, result.CurrentHP, 0)
	})

	t.Run("fall damage with non-existent actor fails", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test",
		})
		require.NoError(t, err)

		_, err = e.SetPhase(ctx, gameResult.Game.ID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		_, err = e.ApplyFallDamage(ctx, ApplyFallDamageRequest{
			GameID:  gameResult.Game.ID,
			ActorID: "non-existent",
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

// ============================================================================
// CalculateBreathHolding Tests
// ============================================================================

func TestCalculateBreathHolding(t *testing.T) {
	t.Run("calculate breath holding with normal constitution", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameID, actorID := setupTestPC(t, e, "Breath Holder", 15, 14)

		_, err := e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		result, err := e.CalculateBreathHolding(ctx, CalculateBreathHoldingRequest{
			GameID:  gameID,
			ActorID: actorID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 14, result.Constitution)
		assert.Equal(t, 2, result.ConstitutionMod) // (14-10)/2 = 2
		// 闭气时间 = (1 + 2) * 60 = 180 秒
		assert.Equal(t, 180, result.CanHoldBreathSecs)
		// 窒息轮数 = 2
		assert.Equal(t, 2, result.RoundsUntilUnconscious)
		assert.Contains(t, result.Message, "闭气")
		assert.Contains(t, result.Message, "180 秒")
	})

	t.Run("calculate breath holding with high constitution", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameID, actorID := setupTestPC(t, e, "Tough Guy", 15, 20)

		_, err := e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		result, err := e.CalculateBreathHolding(ctx, CalculateBreathHoldingRequest{
			GameID:  gameID,
			ActorID: actorID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 20, result.Constitution)
		assert.Equal(t, 5, result.ConstitutionMod)
		// (1 + 5) * 60 = 360 秒
		assert.Equal(t, 360, result.CanHoldBreathSecs)
		assert.Equal(t, 5, result.RoundsUntilUnconscious)
	})

	t.Run("calculate breath holding with low constitution", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameID, actorID := setupTestPC(t, e, "Weak Lungs", 15, 8)

		_, err := e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		result, err := e.CalculateBreathHolding(ctx, CalculateBreathHoldingRequest{
			GameID:  gameID,
			ActorID: actorID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 8, result.Constitution)
		assert.Equal(t, -1, result.ConstitutionMod)
		// 最少 1 分钟 = 60 秒
		assert.Equal(t, 60, result.CanHoldBreathSecs)
		// 最少 1 轮
		assert.Equal(t, 1, result.RoundsUntilUnconscious)
	})

	t.Run("calculate breath holding with very low constitution", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameID, actorID := setupTestPC(t, e, "Very Weak", 15, 3)

		_, err := e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		result, err := e.CalculateBreathHolding(ctx, CalculateBreathHoldingRequest{
			GameID:  gameID,
			ActorID: actorID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		// (3-10)/2 = -3.5, 向下取整 = -4
		assert.Equal(t, -4, result.ConstitutionMod)
		// 仍然最少 1 分钟
		assert.Equal(t, 60, result.CanHoldBreathSecs)
		assert.Equal(t, 1, result.RoundsUntilUnconscious)
	})

	t.Run("calculate breath holding with non-existent actor fails", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test",
		})
		require.NoError(t, err)

		_, err = e.SetPhase(ctx, gameResult.Game.ID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		_, err = e.CalculateBreathHolding(ctx, CalculateBreathHoldingRequest{
			GameID:  gameResult.Game.ID,
			ActorID: "non-existent",
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

// ============================================================================
// ApplySuffocation Tests
// ============================================================================

func TestApplySuffocation(t *testing.T) {
	t.Run("apply suffocation starts suffocation", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameID, actorID := setupTestPC(t, e, "Drowning", 15, 14)

		_, err := e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		result, err := e.ApplySuffocation(ctx, ApplySuffocationRequest{
			GameID:  gameID,
			ActorID: actorID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 14, result.Constitution)
		assert.Equal(t, 2, result.SuffocationRounds)
		assert.False(t, result.IsUnconscious)
		assert.Contains(t, result.Message, "开始窒息")
		assert.Contains(t, result.Message, "2 轮")
	})

	t.Run("apply suffocation with high constitution", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameID, actorID := setupTestPC(t, e, "Tough", 15, 18)

		_, err := e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		result, err := e.ApplySuffocation(ctx, ApplySuffocationRequest{
			GameID:  gameID,
			ActorID: actorID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		// (18-10)/2 = 4
		assert.Equal(t, 4, result.SuffocationRounds)
	})

	t.Run("apply suffocation with non-existent actor fails", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test",
		})
		require.NoError(t, err)

		_, err = e.SetPhase(ctx, gameResult.Game.ID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		_, err = e.ApplySuffocation(ctx, ApplySuffocationRequest{
			GameID:  gameResult.Game.ID,
			ActorID: "non-existent",
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

// ============================================================================
// PerformEncounterCheck Tests
// ============================================================================

func TestPerformEncounterCheck(t *testing.T) {
	t.Run("perform encounter check returns result", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test encounters",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		result, err := e.PerformEncounterCheck(ctx, PerformEncounterCheckRequest{
			GameID: gameID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		// Roll should be 1-6
		assert.GreaterOrEqual(t, result.Roll, 1)
		assert.LessOrEqual(t, result.Roll, 6)
		// If encountered, type should be one of the valid types
		if result.Encountered {
			assert.Contains(t, []string{"monster", "npc", "treasure", "trap"}, result.EncounterType)
		}
		assert.NotEmpty(t, result.Message)
	})

	t.Run("perform encounter check multiple times", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test encounters",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		// Run multiple checks to verify it works repeatedly
		for i := 0; i < 10; i++ {
			result, err := e.PerformEncounterCheck(ctx, PerformEncounterCheckRequest{
				GameID: gameID,
			})

			require.NoError(t, err)
			require.NotNil(t, result)
			assert.GreaterOrEqual(t, result.Roll, 1)
			assert.LessOrEqual(t, result.Roll, 6)
		}
	})

	t.Run("perform encounter check with invalid game fails", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		_, err := e.PerformEncounterCheck(ctx, PerformEncounterCheckRequest{
			GameID: "non-existent",
		})

		assert.Error(t, err)
	})
}
