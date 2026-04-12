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

	// 测试完整D&D战斗流程：包含魔法、移动、战斗动作等，6回合以上
	t.Run("完整DnD战斗模拟 - 多回合战术战斗", func(t *testing.T) {
		e := engine.NewTestEngine(t)
		ctx := context.Background()

		// 创建游戏
		gameResult, err := e.NewGame(ctx, engine.NewGameRequest{
			Name:        "Tactical Combat",
			Description: "完整战术战斗测试",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// 创建战士（5级）
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
		t.Logf("创建战士: %s", fighter.Actor.Name)

		// 创建法师（5级）
		wizard, err := e.CreatePC(ctx, engine.CreatePCRequest{
			GameID: gameID,
			PC: &engine.PlayerCharacterInput{
				Name:  "Merlin",
				Race:  "HighElf",
				Class: "法师",
				Level: 5,
				AbilityScores: engine.AbilityScoresInput{
					Strength:     8,
					Dexterity:    14,
					Constitution: 12,
					Intelligence: 18,
					Wisdom:       13,
					Charisma:     10,
				},
			},
		})
		require.NoError(t, err)
		t.Logf("创建法师: %s", wizard.Actor.Name)

		// 创建盗贼（5级）
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
		t.Logf("创建盗贼: %s", rogue.Actor.Name)

		// 创建兽人战士
		orc, err := e.CreateEnemy(ctx, engine.CreateEnemyRequest{
			GameID: gameID,
			Enemy: &engine.EnemyInput{
				Name:        "Orc Warrior",
				Description: "兽人战士",
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
		t.Logf("创建兽人: %s, HP=%d", orc.Actor.Name, orc.Actor.HitPoints.Maximum)

		// 创建哥布林弓箭手
		goblin, err := e.CreateEnemy(ctx, engine.CreateEnemyRequest{
			GameID: gameID,
			Enemy: &engine.EnemyInput{
				Name:        "Goblin Archer",
				Description: "哥布林弓箭手",
				Size:        model.SizeSmall,
				Speed:       30,
				AbilityScores: engine.AbilityScoresInput{
					Strength:     8,
					Dexterity:    14,
					Constitution: 10,
					Intelligence: 10,
					Wisdom:       12,
					Charisma:     12,
				},
				ChallengeRating: 1,
				HitPoints:       22,
				ArmorClass:      14,
			},
		})
		require.NoError(t, err)
		t.Logf("创建哥布林: %s, HP=%d", goblin.Actor.Name, goblin.Actor.HitPoints.Maximum)

		// 切换到探索阶段
		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "准备战斗")
		require.NoError(t, err)

		// 开始战斗
		combatResult, err := e.StartCombat(ctx, engine.StartCombatRequest{
			GameID: gameID,
			ParticipantIDs: []model.ID{
				fighter.Actor.ID,
				wizard.Actor.ID,
				rogue.Actor.ID,
				orc.Actor.ID,
				goblin.Actor.ID,
			},
		})
		require.NoError(t, err)
		t.Logf("战斗开始 - 回合: %d, 参战者: %d", combatResult.Combat.Round, len(combatResult.Combat.Initiative))

		// 辅助函数
		getHP := func(id model.ID) int {
			actor, _ := e.GetActor(ctx, engine.GetActorRequest{GameID: gameID, ActorID: id})
			return actor.Actor.HitPoints.Current
		}

		advanceTo := func(name string) {
			for i := 0; i < 20; i++ {
				turn, _ := e.GetCurrentTurn(ctx, engine.GetCurrentTurnRequest{GameID: gameID})
				if turn.ActorName == name {
					return
				}
				e.NextTurn(ctx, engine.NextTurnRequest{GameID: gameID})
			}
		}

		attack := func(attacker, target model.ID) bool {
			result, err := e.ExecuteAttack(ctx, engine.ExecuteAttackRequest{
				GameID:     gameID,
				AttackerID: attacker,
				TargetID:   target,
				Attack:     engine.AttackInput{IsUnarmed: false},
			})
			if err == nil && result.AttackResult.Hit {
				return true
			}
			return false
		}

		castSpell := func(caster model.ID, spellID string, targets []model.ID, slotLevel int) bool {
			_, err := e.CastSpell(ctx, engine.CastSpellRequest{
				GameID:   gameID,
				CasterID: caster,
				Spell: engine.SpellInput{
					SpellID:   spellID,
					SlotLevel: slotLevel,
					TargetIDs: targets,
				},
			})
			return err == nil
		}

		doAction := func(actor model.ID, actionType string, params map[string]any) bool {
			action := engine.ActionInput{
				Type: model.ActionType(actionType),
			}
			if params != nil {
				action.Details = params
			}
			_, err := e.ExecuteAction(ctx, engine.ExecuteActionRequest{
				GameID:  gameID,
				ActorID: actor,
				Action:  action,
			})
			return err == nil
		}

		// ========== 第1回合：初始交锋 ==========
		t.Run("回合1 - 初始交锋", func(t *testing.T) {
			// 战士回合 - 攻击兽人
			advanceTo("Marcus")
			t.Log("=== 战士 Marcus 攻击兽人 ===")
			if attack(fighter.Actor.ID, orc.Actor.ID) {
				t.Logf("战士命中兽人! 兽人HP: %d", getHP(orc.Actor.ID))
			} else {
				t.Log("战士未命中")
			}

			// 法师回合 - 施放火球术
			advanceTo("Merlin")
			t.Log("=== 法师 Merlin 施放火球术 ===")
			if castSpell(wizard.Actor.ID, "火球术", []model.ID{orc.Actor.ID, goblin.Actor.ID}, 3) {
				t.Logf("火球术命中! 兽人HP: %d, 哥布林HP: %d", getHP(orc.Actor.ID), getHP(goblin.Actor.ID))
			} else {
				t.Log("火球术施放失败，改用魔法飞弹")
				castSpell(wizard.Actor.ID, "魔法飞弹", []model.ID{orc.Actor.ID}, 1)
			}

			// 盗贼回合 - 隐藏后攻击
			advanceTo("Shadow")
			t.Log("=== 盗贼 Shadow 隐藏并攻击 ===")
			doAction(rogue.Actor.ID, "Hide", nil)
			if attack(rogue.Actor.ID, orc.Actor.ID) {
				t.Logf("盗贼命中兽人! 兽人HP: %d", getHP(orc.Actor.ID))
			} else {
				t.Log("盗贼未命中")
			}

			// 兽人回合 - 攻击战士
			advanceTo("Orc Warrior")
			t.Log("=== 兽人 Orc Warrior 攻击战士 ===")
			if attack(orc.Actor.ID, fighter.Actor.ID) {
				t.Logf("兽人命中战士! 战士HP: %d", getHP(fighter.Actor.ID))
			} else {
				t.Log("兽人未命中")
			}

			// 哥布林回合 - 攻击法师
			advanceTo("Goblin Archer")
			t.Log("=== 哥布林 Goblin Archer 攻击法师 ===")
			if attack(goblin.Actor.ID, wizard.Actor.ID) {
				t.Logf("哥布林命中法师! 法师HP: %d", getHP(wizard.Actor.ID))
			} else {
				t.Log("哥布林未命中")
			}

			t.Logf("回合1结束 - 战士HP: %d, 法师HP: %d, 盗贼HP: %d, 兽人HP: %d, 哥布林HP: %d",
				getHP(fighter.Actor.ID), getHP(wizard.Actor.ID), getHP(rogue.Actor.ID),
				getHP(orc.Actor.ID), getHP(goblin.Actor.ID))
		})

		// ========== 第2回合：战术机动 ==========
		t.Run("回合2 - 战术机动", func(t *testing.T) {
			// 战士回合 - 使用冲刺并攻击哥布林
			advanceTo("Marcus")
			t.Log("=== 战士 Marcus 冲刺并攻击哥布林 ===")
			doAction(fighter.Actor.ID, "Dash", nil)
			if attack(fighter.Actor.ID, goblin.Actor.ID) {
				t.Logf("战士命中哥布林! 哥布林HP: %d", getHP(goblin.Actor.ID))
			} else {
				t.Log("战士未命中")
			}

			// 法师回合 - 闪避
			advanceTo("Merlin")
			t.Log("=== 法师 Merlin 执行闪避 ===")
			if !doAction(wizard.Actor.ID, "Dodge", nil) {
				t.Log("闪避失败，改用魔法飞弹")
				castSpell(wizard.Actor.ID, "魔法飞弹", []model.ID{orc.Actor.ID}, 1)
			}

			// 盗贼回合 - 脱离并移动
			advanceTo("Shadow")
			t.Log("=== 盗贼 Shadow 脱离战斗 ===")
			doAction(rogue.Actor.ID, "Disengage", nil)

			// 兽人回合 - 攻击战士
			advanceTo("Orc Warrior")
			t.Log("=== 兽人 Orc Warrior 攻击战士 ===")
			if attack(orc.Actor.ID, fighter.Actor.ID) {
				t.Logf("兽人命中战士! 战士HP: %d", getHP(fighter.Actor.ID))
			} else {
				t.Log("兽人未命中")
			}

			// 哥布林回合 - 近战攻击
			advanceTo("Goblin Archer")
			t.Log("=== 哥布林 Goblin Archer 近战攻击战士 ===")
			if attack(goblin.Actor.ID, fighter.Actor.ID) {
				t.Logf("哥布林命中战士! 战士HP: %d", getHP(fighter.Actor.ID))
			} else {
				t.Log("哥布林未命中")
			}

			t.Logf("回合2结束 - 战士HP: %d, 兽人HP: %d, 哥布林HP: %d",
				getHP(fighter.Actor.ID), getHP(orc.Actor.ID), getHP(goblin.Actor.ID))
		})

		// ========== 第3回合：激烈战斗 ==========
		t.Run("回合3 - 激烈战斗", func(t *testing.T) {
			// 战士回合
			advanceTo("Marcus")
			t.Log("=== 战士 Marcus 攻击兽人 ===")
			result, _ := e.ExecuteAttack(ctx, engine.ExecuteAttackRequest{
				GameID:     gameID,
				AttackerID: fighter.Actor.ID,
				TargetID:   orc.Actor.ID,
				Attack:     engine.AttackInput{IsUnarmed: false},
			})
			if result.AttackResult.Hit {
				t.Logf("战士命中兽人! 伤害: %d, 暴击: %v, 兽人HP: %d",
					result.AttackResult.Damage.FinalDamage,
					result.AttackResult.IsCritical,
					getHP(orc.Actor.ID))
			} else {
				t.Logf("战士未命中, 大失败: %v", result.AttackResult.IsFumble)
			}

			// 法师回合 - 施放闪电束
			advanceTo("Merlin")
			t.Log("=== 法师 Merlin 施放闪电束 ===")
			if !castSpell(wizard.Actor.ID, "闪电束", []model.ID{orc.Actor.ID, goblin.Actor.ID}, 3) {
				t.Log("闪电束失败，改用灼热射线")
				castSpell(wizard.Actor.ID, "灼热射线", []model.ID{orc.Actor.ID}, 2)
			}

			// 盗贼回合 - 协助战士
			advanceTo("Shadow")
			t.Log("=== 盗贼 Shadow 协助战士 ===")
			doAction(rogue.Actor.ID, "Help", map[string]any{"target_id": fighter.Actor.ID})
			if attack(rogue.Actor.ID, orc.Actor.ID) {
				t.Logf("盗贼命中兽人! 兽人HP: %d", getHP(orc.Actor.ID))
			}

			// 兽人回合 - 推撞战士
			advanceTo("Orc Warrior")
			t.Log("=== 兽人 Orc Warrior 推撞战士 ===")
			if !doAction(orc.Actor.ID, "Shove", map[string]any{"target_id": fighter.Actor.ID, "knock_prone": true}) {
				t.Log("推撞失败，改为攻击")
				attack(orc.Actor.ID, fighter.Actor.ID)
			}

			// 哥布林回合 - 脱离并撤退
			advanceTo("Goblin Archer")
			t.Log("=== 哥布林 Goblin Archer 脱离战斗 ===")
			doAction(goblin.Actor.ID, "Disengage", nil)

			t.Logf("回合3结束 - 战士HP: %d, 兽人HP: %d, 哥布林HP: %d",
				getHP(fighter.Actor.ID), getHP(orc.Actor.ID), getHP(goblin.Actor.ID))
		})

		// ========== 第4回合：战斗高潮 ==========
		t.Run("回合4 - 战斗高潮", func(t *testing.T) {
			// 战士回合
			advanceTo("Marcus")
			t.Log("=== 战士 Marcus 攻击哥布林 ===")
			if attack(fighter.Actor.ID, goblin.Actor.ID) {
				t.Logf("战士命中哥布林! 哥布林HP: %d", getHP(goblin.Actor.ID))
				if getHP(goblin.Actor.ID) <= 0 {
					t.Log("哥布林被击败!")
				}
			}

			// 法师回合 - 治疗战士
			advanceTo("Merlin")
			t.Log("=== 法师 Merlin 治疗战士 ===")
			healResult, err := e.ExecuteHealing(ctx, engine.ExecuteHealingRequest{
				GameID:   gameID,
				TargetID: fighter.Actor.ID,
				Amount:   20,
			})
			if err == nil {
				t.Logf("治疗成功! 恢复 %d HP, 战士HP: %d",
					healResult.Healed, getHP(fighter.Actor.ID))
			} else {
				t.Log("治疗失败，改用魔法飞弹")
				castSpell(wizard.Actor.ID, "魔法飞弹", []model.ID{orc.Actor.ID}, 1)
			}

			// 盗贼回合
			advanceTo("Shadow")
			t.Log("=== 盗贼 Shadow 攻击兽人 ===")
			if attack(rogue.Actor.ID, orc.Actor.ID) {
				t.Logf("盗贼命中兽人! 兽人HP: %d", getHP(orc.Actor.ID))
				if getHP(orc.Actor.ID) <= 0 {
					t.Log("兽人被击败!")
				}
			}

			// 兽人回合（如果还活着）
			if getHP(orc.Actor.ID) > 0 {
				advanceTo("Orc Warrior")
				t.Log("=== 兽人 Orc Warrior 攻击盗贼 ===")
				if attack(orc.Actor.ID, rogue.Actor.ID) {
					t.Logf("兽人命中盗贼! 盗贼HP: %d", getHP(rogue.Actor.ID))
				}
			}

			// 哥布林回合（如果还活着）
			if getHP(goblin.Actor.ID) > 0 {
				advanceTo("Goblin Archer")
				t.Log("=== 哥布林 Goblin Archer 擒抱战士 ===")
				if !doAction(goblin.Actor.ID, "Grapple", map[string]any{"target_id": fighter.Actor.ID}) {
					t.Log("擒抱失败，改为攻击")
					attack(goblin.Actor.ID, fighter.Actor.ID)
				}
			}

			t.Logf("回合4结束 - 战士HP: %d, 法师HP: %d, 盗贼HP: %d, 兽人HP: %d, 哥布林HP: %d",
				getHP(fighter.Actor.ID), getHP(wizard.Actor.ID), getHP(rogue.Actor.ID),
				getHP(orc.Actor.ID), getHP(goblin.Actor.ID))
		})

		// ========== 第5回合：决胜时刻 ==========
		t.Run("回合5 - 决胜时刻", func(t *testing.T) {
			// 战士回合
			advanceTo("Marcus")
			t.Log("=== 战士 Marcus 攻击 ===")
			if getHP(orc.Actor.ID) > 0 {
				result, _ := e.ExecuteAttack(ctx, engine.ExecuteAttackRequest{
					GameID:     gameID,
					AttackerID: fighter.Actor.ID,
					TargetID:   orc.Actor.ID,
					Attack:     engine.AttackInput{IsUnarmed: false},
				})
				if result.AttackResult.Hit {
					t.Logf("战士命中兽人! 伤害: %d, 兽人HP: %d",
						result.AttackResult.Damage.FinalDamage, getHP(orc.Actor.ID))
				}
			}
			if getHP(goblin.Actor.ID) > 0 {
				attack(fighter.Actor.ID, goblin.Actor.ID)
			}

			// 法师回合
			advanceTo("Merlin")
			t.Log("=== 法师 Merlin 施放火球术 ===")
			targets := []model.ID{}
			if getHP(orc.Actor.ID) > 0 {
				targets = append(targets, orc.Actor.ID)
			}
			if getHP(goblin.Actor.ID) > 0 {
				targets = append(targets, goblin.Actor.ID)
			}
			if len(targets) > 0 {
				castSpell(wizard.Actor.ID, "火球术", targets, 3)
			}

			// 盗贼回合
			advanceTo("Shadow")
			t.Log("=== 盗贼 Shadow 准备动作 ===")
			if !doAction(rogue.Actor.ID, "Ready", map[string]any{"trigger": "敌人移动时攻击"}) {
				t.Log("准备失败，直接攻击")
				if getHP(orc.Actor.ID) > 0 {
					attack(rogue.Actor.ID, orc.Actor.ID)
				}
			}

			// 兽人回合
			if getHP(orc.Actor.ID) > 0 {
				advanceTo("Orc Warrior")
				t.Log("=== 兽人 Orc Warrior 攻击 ===")
				attack(orc.Actor.ID, fighter.Actor.ID)
			}

			// 哥布林回合
			if getHP(goblin.Actor.ID) > 0 {
				advanceTo("Goblin Archer")
				t.Log("=== 哥布林 Goblin Archer 攻击 ===")
				attack(goblin.Actor.ID, wizard.Actor.ID)
			}

			t.Logf("回合5结束 - 战士HP: %d, 兽人HP: %d, 哥布林HP: %d",
				getHP(fighter.Actor.ID), getHP(orc.Actor.ID), getHP(goblin.Actor.ID))
		})

		// ========== 第6回合：清理战场 ==========
		t.Run("回合6 - 清理战场", func(t *testing.T) {
			orcAlive := getHP(orc.Actor.ID) > 0
			goblinAlive := getHP(goblin.Actor.ID) > 0

			if !orcAlive && !goblinAlive {
				t.Log("所有敌人都被击败，跳过本回合")
				return
			}

			// 战士回合
			advanceTo("Marcus")
			t.Log("=== 战士 Marcus 攻击 ===")
			if orcAlive {
				if attack(fighter.Actor.ID, orc.Actor.ID) {
					t.Logf("战士命中兽人! 兽人HP: %d", getHP(orc.Actor.ID))
				}
			}

			// 法师回合
			advanceTo("Merlin")
			t.Log("=== 法师 Merlin 攻击 ===")
			if goblinAlive {
				castSpell(wizard.Actor.ID, "魔法飞弹", []model.ID{goblin.Actor.ID}, 1)
			} else if orcAlive {
				castSpell(wizard.Actor.ID, "魔法飞弹", []model.ID{orc.Actor.ID}, 1)
			}

			// 盗贼回合
			advanceTo("Shadow")
			t.Log("=== 盗贼 Shadow 攻击 ===")
			if orcAlive && getHP(orc.Actor.ID) > 0 {
				if attack(rogue.Actor.ID, orc.Actor.ID) {
					t.Logf("盗贼命中兽人! 兽人HP: %d", getHP(orc.Actor.ID))
				}
			}

			t.Logf("回合6结束")
		})

		// ========== 结束战斗 ==========
		t.Run("结束战斗", func(t *testing.T) {
			summary, _ := e.GetCombatSummary(ctx, gameID)

			t.Logf("=== 战斗结束 ===")
			t.Logf("战斗持续了 %d 回合", summary.Round)
			t.Logf("最终状态:")
			t.Logf("  战士 Marcus: %d/%d HP", getHP(fighter.Actor.ID), fighter.Actor.HitPoints.Maximum)
			t.Logf("  法师 Merlin: %d/%d HP", getHP(wizard.Actor.ID), wizard.Actor.HitPoints.Maximum)
			t.Logf("  盗贼 Shadow: %d/%d HP", getHP(rogue.Actor.ID), rogue.Actor.HitPoints.Maximum)
			t.Logf("  兽人 Orc Warrior: %d/%d HP", getHP(orc.Actor.ID), orc.Actor.HitPoints.Maximum)
			t.Logf("  哥布林 Goblin Archer: %d/%d HP", getHP(goblin.Actor.ID), goblin.Actor.HitPoints.Maximum)

			err := e.EndCombat(ctx, engine.EndCombatRequest{GameID: gameID})
			require.NoError(t, err)

			phase, _ := e.GetPhase(ctx, gameID)
			assert.Equal(t, model.PhaseExploration, phase)
			t.Logf("当前阶段: %s", phase)
		})
	})
}
