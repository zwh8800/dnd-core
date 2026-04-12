package testsuite

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zwh8800/dnd-core/pkg/engine"
	"github.com/zwh8800/dnd-core/pkg/model"
)

// TestCombatEncounters 测试战斗遭遇系统
// 包含多角色参与完整战斗回合的测试
func TestCombatEncounters(t *testing.T) {
	// 测试完整战斗回合：多个参与者按先攻顺序进行战斗
	t.Run("多角色参与完整战斗回合", func(t *testing.T) {
		// 创建测试引擎和上下文
		e := engine.NewTestEngine(t)
		ctx := context.Background()

		// 创建新游戏
		gameResult, err := e.NewGame(ctx, engine.NewGameRequest{
			Name:        "Arena Battle",
			Description: "角斗场战斗",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// 创建战士角色（5级）
		fighter, err := e.CreatePC(ctx, engine.CreatePCRequest{
			GameID: gameID,
			PC: &engine.PlayerCharacterInput{
				Name:  "Marcus",
				Race:  "Human",
				Class: "战士",
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

		// 创建盗贼角色（5级）
		rogue, err := e.CreatePC(ctx, engine.CreatePCRequest{
			GameID: gameID,
			PC: &engine.PlayerCharacterInput{
				Name:  "Shadow",
				Race:  "Halfling",
				Class: "游荡者",
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

		// 创建敌人：兽人战士
		orc, err := e.CreateEnemy(ctx, engine.CreateEnemyRequest{
			GameID: gameID,
			Enemy: &engine.EnemyInput{
				Name:        "Orc Brute",
				Description: "凶猛的兽人战士",
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

		// 切换到探索阶段
		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "开始战斗")
		require.NoError(t, err)

		// 开始战斗
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

		// 验证战斗状态
		assert.Equal(t, model.CombatStatusActive, combatResult.Combat.Status)
		assert.Equal(t, 3, len(combatResult.Combat.Initiative))
		t.Logf("战斗已开始 - 第%d回合, %d名参战者", combatResult.Combat.Round, len(combatResult.Combat.Initiative))

		// 验证：通过GetCurrentTurn确认当前回合信息正确
		currentTurn, err := e.GetCurrentTurn(ctx, engine.GetCurrentTurnRequest{
			GameID: gameID,
		})
		require.NoError(t, err)
		require.NotNil(t, currentTurn)
		t.Logf("当前回合: %s", currentTurn.ActorName)

		// 测试执行攻击动作
		t.Run("执行攻击动作", func(t *testing.T) {
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
			t.Logf("攻击结果: 掷骰=%d, 命中=%v, 重击=%v",
				attackResult.AttackResult.Roll.Total,
				attackResult.AttackResult.Hit,
				attackResult.AttackResult.IsCritical)

			// 验证敌人状态与攻击结果一致
			if attackResult.AttackResult.Hit {
				// 获取敌人当前状态
				enemyState, err := e.GetActor(ctx, engine.GetActorRequest{
					GameID:  gameID,
					ActorID: orc.Actor.ID,
				})
				require.NoError(t, err)
				require.NotNil(t, enemyState.Actor)

				// 验证生命值变化与攻击结果一致
				if attackResult.AttackResult.Damage != nil {
					assert.Equal(t, attackResult.AttackResult.Damage.TargetHPAfter, enemyState.Actor.HitPoints.Current,
						"敌人当前生命值应与攻击结果中的目标HP一致")
					assert.Equal(t, attackResult.AttackResult.Damage.TargetHPBefore, enemyState.Actor.HitPoints.Maximum,
						"敌人最大生命值应与攻击结果中的目标HP前一致")

					t.Logf("敌人状态验证: 当前HP=%d, 最大HP=%d, 是否死亡=%v",
						enemyState.Actor.HitPoints.Current,
						enemyState.Actor.HitPoints.Maximum,
						attackResult.AttackResult.Damage.IsDead)
				}
			}
		})

		// 测试下一回合推进
		t.Run("推进下一回合", func(t *testing.T) {
			nextResult, err := e.NextTurn(ctx, engine.NextTurnRequest{
				GameID: gameID,
			})
			require.NoError(t, err)
			require.NotNil(t, nextResult)

			// 验证：通过GetCombatSummary确认战斗状态更新
			summary, err := e.GetCombatSummary(ctx, gameID)
			require.NoError(t, err)
			require.NotNil(t, summary)
			t.Logf("战斗状态: 第%d回合", summary.Round)
		})

		// 测试结束战斗
		t.Run("结束战斗", func(t *testing.T) {
			err := e.EndCombat(ctx, engine.EndCombatRequest{
				GameID: gameID,
			})
			require.NoError(t, err)

			// 验证：战斗结束后Phase应恢复为Exploration
			phase, err := e.GetPhase(ctx, gameID)
			require.NoError(t, err)
			assert.Equal(t, model.PhaseExploration, phase)
		})
	})
}
