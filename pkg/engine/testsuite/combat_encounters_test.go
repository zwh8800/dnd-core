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

		// 验证先攻顺序
		t.Run("验证先攻顺序", func(t *testing.T) {
			summary, err := e.GetCombatSummary(ctx, gameID)
			require.NoError(t, err)
			require.NotNil(t, summary)

			t.Logf("先攻顺序:")
			for i, entry := range summary.TurnOrder {
				t.Logf("  %d. %s (先攻: %d)", i+1, entry.ActorName, entry.Initiative)
			}

			// 验证所有参与者都在先攻列表中
			participantNames := make(map[string]bool)
			for _, entry := range summary.TurnOrder {
				participantNames[entry.ActorName] = true
			}
			assert.True(t, participantNames["Marcus"], "战士 Marcus 应在先攻列表中")
			assert.True(t, participantNames["Shadow"], "盗贼 Shadow 应在先攻列表中")
			assert.True(t, participantNames["Orc Brute"], "兽人 Orc Brute 应在先攻列表中")
		})

		// 验证：通过GetCurrentTurn确认当前回合信息正确
		currentTurn, err := e.GetCurrentTurn(ctx, engine.GetCurrentTurnRequest{
			GameID: gameID,
		})
		require.NoError(t, err)
		require.NotNil(t, currentTurn)
		t.Logf("当前回合: %s", currentTurn.ActorName)

		// 测试执行攻击动作
		t.Run("执行攻击动作", func(t *testing.T) {
			// 记录攻击前敌人状态
			enemyBefore, err := e.GetActor(ctx, engine.GetActorRequest{
				GameID:  gameID,
				ActorID: orc.Actor.ID,
			})
			require.NoError(t, err)
			t.Logf("攻击前敌人状态: HP=%d/%d", enemyBefore.Actor.HitPoints.Current, enemyBefore.Actor.HitPoints.Maximum)

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

			t.Logf("攻击结果: 掷骰=%d, 攻击总值=%d, 目标AC=%d, 命中=%v, 重击=%v, 大失败=%v",
				attackResult.AttackResult.Roll.Total,
				attackResult.AttackResult.AttackTotal,
				attackResult.AttackResult.TargetAC,
				attackResult.AttackResult.Hit,
				attackResult.AttackResult.IsCritical,
				attackResult.AttackResult.IsFumble)

			// 验证攻击结果逻辑
			if attackResult.AttackResult.Hit {
				t.Logf("攻击命中! 伤害: %d", attackResult.AttackResult.Damage.FinalDamage)

				// 验证伤害计算
				if attackResult.AttackResult.Damage != nil {
					assert.GreaterOrEqual(t, attackResult.AttackResult.Damage.FinalDamage, 0, "伤害值应大于等于0")
					assert.Equal(t, attackResult.AttackResult.Damage.TargetHPBefore, enemyBefore.Actor.HitPoints.Current,
						"攻击前HP应与记录一致")

					// 验证敌人状态与攻击结果一致
					enemyState, err := e.GetActor(ctx, engine.GetActorRequest{
						GameID:  gameID,
						ActorID: orc.Actor.ID,
					})
					require.NoError(t, err)
					require.NotNil(t, enemyState.Actor)

					assert.Equal(t, attackResult.AttackResult.Damage.TargetHPAfter, enemyState.Actor.HitPoints.Current,
						"敌人当前生命值应与攻击结果中的目标HP一致")
					assert.Equal(t, attackResult.AttackResult.Damage.TargetHPBefore, enemyState.Actor.HitPoints.Maximum,
						"敌人最大生命值应与攻击结果中的目标HP前一致")

					t.Logf("敌人状态验证: 当前HP=%d, 最大HP=%d, 是否死亡=%v",
						enemyState.Actor.HitPoints.Current,
						enemyState.Actor.HitPoints.Maximum,
						attackResult.AttackResult.Damage.IsDead)

					// 验证HP变化
					expectedHP := enemyBefore.Actor.HitPoints.Current - attackResult.AttackResult.Damage.FinalDamage
					assert.Equal(t, expectedHP, enemyState.Actor.HitPoints.Current,
						"敌人HP应等于攻击前HP减去伤害值")
				}
			} else {
				t.Logf("攻击未命中! 擦伤伤害: %d", attackResult.AttackResult.GrazeDamage)
				assert.Equal(t, 0, attackResult.AttackResult.GrazeDamage, "未命中时擦伤伤害应为0")
			}
		})

		// 测试盗贼攻击
		t.Run("盗贼执行攻击动作", func(t *testing.T) {
			// 推进到盗贼回合
			for {
				currentTurn, err := e.GetCurrentTurn(ctx, engine.GetCurrentTurnRequest{
					GameID: gameID,
				})
				require.NoError(t, err)

				if currentTurn.ActorName == "Shadow" {
					break
				}

				_, err = e.NextTurn(ctx, engine.NextTurnRequest{
					GameID: gameID,
				})
				require.NoError(t, err)
			}

			// 记录攻击前敌人状态
			enemyBefore, err := e.GetActor(ctx, engine.GetActorRequest{
				GameID:  gameID,
				ActorID: orc.Actor.ID,
			})
			require.NoError(t, err)
			t.Logf("盗贼攻击前敌人状态: HP=%d/%d", enemyBefore.Actor.HitPoints.Current, enemyBefore.Actor.HitPoints.Maximum)

			attackResult, err := e.ExecuteAttack(ctx, engine.ExecuteAttackRequest{
				GameID:     gameID,
				AttackerID: rogue.Actor.ID,
				TargetID:   orc.Actor.ID,
				Attack: engine.AttackInput{
					IsUnarmed: false,
				},
			})
			require.NoError(t, err)
			require.NotNil(t, attackResult)

			t.Logf("盗贼攻击结果: 掷骰=%d, 命中=%v, 伤害=%d",
				attackResult.AttackResult.Roll.Total,
				attackResult.AttackResult.Hit,
				attackResult.AttackResult.Damage.FinalDamage)

			if attackResult.AttackResult.Hit && attackResult.AttackResult.Damage != nil {
				// 验证敌人状态
				enemyState, err := e.GetActor(ctx, engine.GetActorRequest{
					GameID:  gameID,
					ActorID: orc.Actor.ID,
				})
				require.NoError(t, err)

				assert.LessOrEqual(t, enemyState.Actor.HitPoints.Current, enemyBefore.Actor.HitPoints.Current,
					"敌人HP不应增加")
				t.Logf("盗贼攻击后敌人状态: HP=%d/%d, 是否死亡=%v",
					enemyState.Actor.HitPoints.Current,
					enemyState.Actor.HitPoints.Maximum,
					attackResult.AttackResult.Damage.IsDead)
			}
		})

		// 测试推进回合
		t.Run("推进下一回合", func(t *testing.T) {
			// 记录当前回合
			currentTurn, err := e.GetCurrentTurn(ctx, engine.GetCurrentTurnRequest{
				GameID: gameID,
			})
			require.NoError(t, err)
			t.Logf("推进前当前回合: %s", currentTurn.ActorName)

			nextResult, err := e.NextTurn(ctx, engine.NextTurnRequest{
				GameID: gameID,
			})
			require.NoError(t, err)
			require.NotNil(t, nextResult)

			// 验证：通过GetCombatSummary确认战斗状态更新
			summary, err := e.GetCombatSummary(ctx, gameID)
			require.NoError(t, err)
			require.NotNil(t, summary)
			t.Logf("战斗状态: 第%d回合, 当前回合: %s", summary.Round, summary.CurrentActor)

			// 验证回合推进
			assert.NotEqual(t, currentTurn.ActorName, summary.CurrentActor,
				"回合推进后当前回合角色应改变")
		})

		// 测试结束战斗
		t.Run("结束战斗", func(t *testing.T) {
			// 记录战斗结束前状态
			summary, err := e.GetCombatSummary(ctx, gameID)
			require.NoError(t, err)
			t.Logf("战斗结束前状态: 第%d回合, 当前回合: %s", summary.Round, summary.CurrentActor)

			err = e.EndCombat(ctx, engine.EndCombatRequest{
				GameID: gameID,
			})
			require.NoError(t, err)

			// 验证：战斗结束后Phase应恢复为Exploration
			phase, err := e.GetPhase(ctx, gameID)
			require.NoError(t, err)
			assert.Equal(t, model.PhaseExploration, phase, "战斗结束后阶段应恢复为探索阶段")

			t.Logf("战斗已结束, 当前阶段: %s", phase)
		})
	})
}
