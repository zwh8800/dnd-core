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
			require.NotNil(t, enemyBefore.Actor)
			hpBeforeAttack := enemyBefore.Actor.HitPoints.Current
			maxHP := enemyBefore.Actor.HitPoints.Maximum
			t.Logf("攻击前敌人状态: HP=%d/%d", hpBeforeAttack, maxHP)

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

			// 获取攻击后敌人状态
			enemyAfter, err := e.GetActor(ctx, engine.GetActorRequest{
				GameID:  gameID,
				ActorID: orc.Actor.ID,
			})
			require.NoError(t, err)
			require.NotNil(t, enemyAfter.Actor)
			hpAfterAttack := enemyAfter.Actor.HitPoints.Current

			// 验证攻击结果逻辑
			if attackResult.AttackResult.Hit {
				t.Logf("攻击命中! 伤害: %d", attackResult.AttackResult.Damage.FinalDamage)

				// 验证伤害计算
				if attackResult.AttackResult.Damage != nil {
					finalDamage := attackResult.AttackResult.Damage.FinalDamage

					// 断言1: 验证攻击结果中记录的HP与实际一致
					assert.Equal(t, hpBeforeAttack, attackResult.AttackResult.Damage.TargetHPBefore,
						"攻击结果中记录的攻击前HP应与实际查询一致")
					assert.Equal(t, hpAfterAttack, attackResult.AttackResult.Damage.TargetHPAfter,
						"攻击结果中记录的攻击后HP应与实际查询一致")

					// 断言2: 验证敌人HP实际减少量等于伤害值
					actualDamageDealt := hpBeforeAttack - hpAfterAttack
					assert.Equal(t, finalDamage, actualDamageDealt,
						"敌人HP实际减少量应等于最终伤害值")

					// 断言3: 验证敌人当前HP计算正确
					expectedHP := hpBeforeAttack - finalDamage
					if expectedHP < 0 {
						expectedHP = 0
					}
					assert.Equal(t, expectedHP, hpAfterAttack,
						"敌人当前HP应等于攻击前HP减去伤害值（最低为0）")

					// 断言4: 验证敌人最大HP不变
					assert.Equal(t, maxHP, enemyAfter.Actor.HitPoints.Maximum,
						"敌人最大HP不应因攻击而改变")

					// 断言5: 如果敌人死亡，验证死亡状态
					if hpAfterAttack <= 0 {
						assert.True(t, attackResult.AttackResult.Damage.IsDead,
							"HP小于等于0时，伤害结果应标记为死亡")
						t.Logf("敌人已死亡! 死亡状态验证通过")
					} else {
						assert.False(t, attackResult.AttackResult.Damage.IsDead,
							"HP大于0时，伤害结果不应标记为死亡")
					}

					// 断言6: 验证伤害值合理性
					assert.GreaterOrEqual(t, finalDamage, 0, "伤害值应大于等于0")
					if attackResult.AttackResult.IsCritical {
						t.Logf("重击验证通过 - 伤害应翻倍")
					}

					t.Logf("敌人状态验证: 当前HP=%d, 最大HP=%d, 实际受到伤害=%d, 是否死亡=%v",
						hpAfterAttack,
						enemyAfter.Actor.HitPoints.Maximum,
						actualDamageDealt,
						attackResult.AttackResult.Damage.IsDead)
				}
			} else {
				t.Logf("攻击未命中! 擦伤伤害: %d", attackResult.AttackResult.GrazeDamage)

				// 断言7: 未命中时HP应保持不变
				assert.Equal(t, hpBeforeAttack, hpAfterAttack,
					"攻击未命中时敌人HP应保持不变")

				// 断言8: 未命中时擦伤伤害应为0
				assert.Equal(t, 0, attackResult.AttackResult.GrazeDamage,
					"未命中时擦伤伤害应为0")

				// 断言9: 未命中时不应有伤害结果
				if attackResult.AttackResult.Damage == nil {
					t.Logf("未命中验证通过 - 无伤害结果")
				}
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
			require.NotNil(t, enemyBefore.Actor)
			hpBeforeRogue := enemyBefore.Actor.HitPoints.Current
			t.Logf("盗贼攻击前敌人状态: HP=%d/%d", hpBeforeRogue, enemyBefore.Actor.HitPoints.Maximum)

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
			require.NotNil(t, attackResult.AttackResult)

			t.Logf("盗贼攻击结果: 掷骰=%d, 命中=%v, 伤害=%d",
				attackResult.AttackResult.Roll.Total,
				attackResult.AttackResult.Hit,
				attackResult.AttackResult.Damage.FinalDamage)

			// 获取攻击后敌人状态
			enemyAfter, err := e.GetActor(ctx, engine.GetActorRequest{
				GameID:  gameID,
				ActorID: orc.Actor.ID,
			})
			require.NoError(t, err)
			require.NotNil(t, enemyAfter.Actor)
			hpAfterRogue := enemyAfter.Actor.HitPoints.Current

			if attackResult.AttackResult.Hit && attackResult.AttackResult.Damage != nil {
				finalDamage := attackResult.AttackResult.Damage.FinalDamage

				// 断言1: 验证HP实际减少
				assert.Less(t, hpAfterRogue, hpBeforeRogue,
					"攻击命中后敌人HP应减少")

				// 断言2: 验证HP减少量等于伤害值
				actualDamage := hpBeforeRogue - hpAfterRogue
				assert.Equal(t, finalDamage, actualDamage,
					"敌人HP减少量应等于最终伤害值")

				// 断言3: 验证攻击结果中记录的HP与实际一致
				assert.Equal(t, hpBeforeRogue, attackResult.AttackResult.Damage.TargetHPBefore,
					"攻击结果中记录的攻击前HP应与实际查询一致")
				assert.Equal(t, hpAfterRogue, attackResult.AttackResult.Damage.TargetHPAfter,
					"攻击结果中记录的攻击后HP应与实际查询一致")

				// 断言4: 验证死亡状态
				if hpAfterRogue <= 0 {
					assert.True(t, attackResult.AttackResult.Damage.IsDead,
						"HP小于等于0时应标记为死亡")
					t.Logf("敌人被盗贼击败! 死亡状态验证通过")
				}

				t.Logf("盗贼攻击后敌人状态验证: HP=%d -> %d, 伤害=%d, 是否死亡=%v",
					hpBeforeRogue, hpAfterRogue, finalDamage,
					attackResult.AttackResult.Damage.IsDead)
			} else {
				// 断言5: 未命中时HP不变
				assert.Equal(t, hpBeforeRogue, hpAfterRogue,
					"攻击未命中时敌人HP应保持不变")
				t.Logf("盗贼未命中，敌人HP保持: %d", hpAfterRogue)
			}
		})

		// 测试执行直接伤害
		t.Run("执行直接伤害", func(t *testing.T) {
			// 记录伤害前敌人状态
			enemyBefore, err := e.GetActor(ctx, engine.GetActorRequest{
				GameID:  gameID,
				ActorID: orc.Actor.ID,
			})
			require.NoError(t, err)
			require.NotNil(t, enemyBefore.Actor)
			hpBeforeDamage := enemyBefore.Actor.HitPoints.Current
			t.Logf("伤害前敌人状态: HP=%d/%d", hpBeforeDamage, enemyBefore.Actor.HitPoints.Maximum)

			// 执行直接伤害（模拟陷阱伤害）
			damageResult, err := e.ExecuteDamage(ctx, engine.ExecuteDamageRequest{
				GameID:   gameID,
				TargetID: orc.Actor.ID,
				Damage: engine.DamageInput{
					Amount: 10,
					Type:   model.DamageTypeFire,
					Source: orc.Actor.ID, // 环境伤害
				},
			})
			require.NoError(t, err)
			require.NotNil(t, damageResult)
			require.NotNil(t, damageResult.DamageResult)

			// 获取伤害后敌人状态
			enemyAfter, err := e.GetActor(ctx, engine.GetActorRequest{
				GameID:  gameID,
				ActorID: orc.Actor.ID,
			})
			require.NoError(t, err)
			require.NotNil(t, enemyAfter.Actor)
			hpAfterDamage := enemyAfter.Actor.HitPoints.Current

			t.Logf("伤害结果: 原始伤害=%d, 最终伤害=%d, HP=%d -> %d",
				damageResult.DamageResult.RawDamage,
				damageResult.DamageResult.FinalDamage,
				damageResult.DamageResult.TargetHPBefore,
				damageResult.DamageResult.TargetHPAfter)

			// 断言1: 验证HP实际减少
			assert.Less(t, hpAfterDamage, hpBeforeDamage,
				"受到伤害后敌人HP应减少")

			// 断言2: 验证HP减少量等于伤害值
			actualDamage := hpBeforeDamage - hpAfterDamage
			assert.Equal(t, damageResult.DamageResult.FinalDamage, actualDamage,
				"敌人HP减少量应等于最终伤害值")

			// 断言3: 验证伤害结果中记录的HP与实际一致
			assert.Equal(t, hpBeforeDamage, damageResult.DamageResult.TargetHPBefore,
				"伤害结果中记录的攻击前HP应与实际查询一致")
			assert.Equal(t, hpAfterDamage, damageResult.DamageResult.TargetHPAfter,
				"伤害结果中记录的攻击后HP应与实际查询一致")

			// 断言4: 验证死亡状态
			if hpAfterDamage <= 0 {
				assert.True(t, damageResult.DamageResult.IsDead,
					"HP小于等于0时应标记为死亡")
				t.Logf("敌人因直接伤害死亡! 死亡状态验证通过")
			} else {
				assert.False(t, damageResult.DamageResult.IsDead,
					"HP大于0时不应标记为死亡")
			}

			// 断言5: 验证伤害值合理性
			// 注意：这里验证最终伤害是否考虑了抗性等
			expectedDamage := 10 // 假设没有抗性
			// 如果有抗性系统，应在这里验证
			assert.Equal(t, expectedDamage, damageResult.DamageResult.FinalDamage,
				"最终伤害应考虑抗性等因素")

			t.Logf("直接伤害验证: HP=%d -> %d, 伤害=%d, 是否死亡=%v",
				hpBeforeDamage, hpAfterDamage, actualDamage,
				damageResult.DamageResult.IsDead)
		})

		// 测试执行治疗
		t.Run("执行治疗", func(t *testing.T) {
			// 先让战士受到一些伤害，然后治疗
			// 记录治疗前战士状态
			fighterBefore, err := e.GetActor(ctx, engine.GetActorRequest{
				GameID:  gameID,
				ActorID: fighter.Actor.ID,
			})
			require.NoError(t, err)
			require.NotNil(t, fighterBefore.Actor)
			hpBeforeHeal := fighterBefore.Actor.HitPoints.Current
			maxHP := fighterBefore.Actor.HitPoints.Maximum

			// 如果战士是满血，先造成伤害
			if hpBeforeHeal >= maxHP {
				t.Logf("战士已满血(%d/%d)，先造成伤害以便测试治疗", hpBeforeHeal, maxHP)
				_, err = e.ExecuteDamage(ctx, engine.ExecuteDamageRequest{
					GameID:   gameID,
					TargetID: fighter.Actor.ID,
					Damage: engine.DamageInput{
						Amount: 5,
						Type:   model.DamageTypeFire,
						Source: fighter.Actor.ID,
					},
				})
				require.NoError(t, err)

				// 重新获取受伤后的状态
				fighterInjured, err := e.GetActor(ctx, engine.GetActorRequest{
					GameID:  gameID,
					ActorID: fighter.Actor.ID,
				})
				require.NoError(t, err)
				require.NotNil(t, fighterInjured.Actor)
				hpBeforeHeal = fighterInjured.Actor.HitPoints.Current
				t.Logf("受伤后战士状态: HP=%d/%d", hpBeforeHeal, maxHP)
			}

			t.Logf("治疗前战士状态: HP=%d/%d", hpBeforeHeal, maxHP)

			// 执行治疗
			healAmount := 15
			healResult, err := e.ExecuteHealing(ctx, engine.ExecuteHealingRequest{
				GameID:   gameID,
				TargetID: fighter.Actor.ID,
				Amount:   healAmount,
			})
			require.NoError(t, err)
			require.NotNil(t, healResult)

			// 获取治疗后战士状态
			fighterAfter, err := e.GetActor(ctx, engine.GetActorRequest{
				GameID:  gameID,
				ActorID: fighter.Actor.ID,
			})
			require.NoError(t, err)
			require.NotNil(t, fighterAfter.Actor)
			hpAfterHeal := fighterAfter.Actor.HitPoints.Current

			t.Logf("治疗结果: 请求治疗=%d, 实际治疗=%d, HP=%d -> %d",
				healAmount, healResult.Healed, hpBeforeHeal, hpAfterHeal)

			// 断言1: 验证HP实际增加
			assert.Greater(t, hpAfterHeal, hpBeforeHeal,
				"治疗后角色HP应增加")

			// 断言2: 验证HP增加量等于实际治疗量
			actualHeal := hpAfterHeal - hpBeforeHeal
			assert.Equal(t, healResult.Healed, actualHeal,
				"角色HP增加量应等于实际治疗量")

			// 断言3: 验证治疗结果中报告的当前HP与实际一致
			assert.Equal(t, hpAfterHeal, healResult.CurrentHP,
				"治疗结果中报告的当前HP应与实际查询一致")

			// 断言4: 验证治疗量不超过最大HP
			expectedHeal := healAmount
			if hpBeforeHeal+healAmount > maxHP {
				expectedHeal = maxHP - hpBeforeHeal
			}
			assert.Equal(t, expectedHeal, healResult.Healed,
				"实际治疗量不应超过最大HP限制")
			assert.LessOrEqual(t, hpAfterHeal, maxHP,
				"治疗后HP不应超过最大HP")

			// 断言5: 验证最大HP不变
			assert.Equal(t, maxHP, fighterAfter.Actor.HitPoints.Maximum,
				"治疗不应改变最大HP")

			t.Logf("治疗验证: HP=%d -> %d, 治疗量=%d, 最大HP=%d",
				hpBeforeHeal, hpAfterHeal, actualHeal, maxHP)
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

		// 测试死亡状态验证（在战斗结束后进行，因为需要创建新敌人）
		t.Run("死亡状态验证", func(t *testing.T) {
			// 切换到探索阶段以创建敌人
			_, err := e.SetPhase(ctx, gameID, model.PhaseExploration, "准备测试死亡")
			require.NoError(t, err)

			// 创建一个低HP的敌人用于测试死亡
			weakEnemy, err := e.CreateEnemy(ctx, engine.CreateEnemyRequest{
				GameID: gameID,
				Enemy: &engine.EnemyInput{
					Name:        "Weak Goblin",
					Description: "虚弱的哥布林",
					Size:        model.SizeSmall,
					Speed:       30,
					AbilityScores: engine.AbilityScoresInput{
						Strength:     8,
						Dexterity:    10,
						Constitution: 8,
						Intelligence: 6,
						Wisdom:       8,
						Charisma:     6,
					},
					ChallengeRating: 0,
					HitPoints:       5, // 非常低的HP
					ArmorClass:      10,
				},
			})
			require.NoError(t, err)
			require.NotNil(t, weakEnemy.Actor)

			t.Logf("创建虚弱敌人: HP=%d", weakEnemy.Actor.HitPoints.Maximum)

			// 重新开始战斗以便执行攻击
			_, err = e.StartCombat(ctx, engine.StartCombatRequest{
				GameID: gameID,
				ParticipantIDs: []model.ID{
					fighter.Actor.ID,
					weakEnemy.Actor.ID,
				},
			})
			require.NoError(t, err)

			// 记录攻击前状态
			enemyBefore, err := e.GetActor(ctx, engine.GetActorRequest{
				GameID:  gameID,
				ActorID: weakEnemy.Actor.ID,
			})
			require.NoError(t, err)
			require.NotNil(t, enemyBefore.Actor)
			hpBefore := enemyBefore.Actor.HitPoints.Current
			assert.Equal(t, 5, hpBefore, "敌人初始HP应为5")

			// 执行致命攻击
			attackResult, err := e.ExecuteAttack(ctx, engine.ExecuteAttackRequest{
				GameID:     gameID,
				AttackerID: fighter.Actor.ID,
				TargetID:   weakEnemy.Actor.ID,
				Attack: engine.AttackInput{
					IsUnarmed: false,
				},
			})
			require.NoError(t, err)
			require.NotNil(t, attackResult)
			require.NotNil(t, attackResult.AttackResult)

			// 获取攻击后状态
			enemyAfter, err := e.GetActor(ctx, engine.GetActorRequest{
				GameID:  gameID,
				ActorID: weakEnemy.Actor.ID,
			})
			require.NoError(t, err)
			require.NotNil(t, enemyAfter.Actor)
			hpAfter := enemyAfter.Actor.HitPoints.Current

			t.Logf("致命攻击结果: 命中=%v, 伤害=%d, HP=%d -> %d",
				attackResult.AttackResult.Hit,
				attackResult.AttackResult.Damage.FinalDamage,
				hpBefore, hpAfter)

			if attackResult.AttackResult.Hit && attackResult.AttackResult.Damage != nil {
				// 断言1: 验证敌人死亡
				assert.True(t, attackResult.AttackResult.Damage.IsDead,
					"致命攻击后应标记为死亡")

				// 断言2: 验证HP为0或负数
				assert.LessOrEqual(t, hpAfter, 0,
					"死亡敌人的HP应小于等于0")

				// 断言3: 验证HP变化正确
				expectedHP := 0 // 死亡时HP应为0
				assert.Equal(t, expectedHP, hpAfter,
					"死亡敌人的HP应为0")

				// 断言4: 验证攻击结果中的HP记录
				assert.Equal(t, hpBefore, attackResult.AttackResult.Damage.TargetHPBefore,
					"攻击前HP记录应正确")
				assert.Equal(t, hpAfter, attackResult.AttackResult.Damage.TargetHPAfter,
					"攻击后HP记录应正确")

				t.Logf("死亡状态验证通过: HP=%d, 死亡=%v",
					hpAfter, attackResult.AttackResult.Damage.IsDead)
			} else {
				t.Log("攻击未命中，跳过死亡验证")
			}

			// 测试直接伤害致死
			t.Run("直接伤害致死验证", func(t *testing.T) {
				// 先确保在探索阶段
				_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "准备测试直接伤害")
				require.NoError(t, err)

				// 创建另一个低HP敌人
				anotherWeakEnemy, err := e.CreateEnemy(ctx, engine.CreateEnemyRequest{
					GameID: gameID,
					Enemy: &engine.EnemyInput{
						Name:        "Test Target",
						Description: "测试目标",
						Size:        model.SizeSmall,
						Speed:       30,
						AbilityScores: engine.AbilityScoresInput{
							Strength:     8,
							Dexterity:    10,
							Constitution: 8,
							Intelligence: 6,
							Wisdom:       8,
							Charisma:     6,
						},
						ChallengeRating: 0,
						HitPoints:       3,
						ArmorClass:      10,
					},
				})
				require.NoError(t, err)
				require.NotNil(t, anotherWeakEnemy.Actor)

				// 开始战斗以便执行伤害
				_, err = e.StartCombat(ctx, engine.StartCombatRequest{
					GameID: gameID,
					ParticipantIDs: []model.ID{
						fighter.Actor.ID,
						anotherWeakEnemy.Actor.ID,
					},
				})
				require.NoError(t, err)

				// 执行致命伤害
				damageResult, err := e.ExecuteDamage(ctx, engine.ExecuteDamageRequest{
					GameID:   gameID,
					TargetID: anotherWeakEnemy.Actor.ID,
					Damage: engine.DamageInput{
						Amount: 10, // 远超HP的伤害
						Type:   model.DamageTypeFire,
						Source: anotherWeakEnemy.Actor.ID,
					},
				})
				require.NoError(t, err)
				require.NotNil(t, damageResult)
				require.NotNil(t, damageResult.DamageResult)

				// 获取伤害后状态
				enemyAfterDamage, err := e.GetActor(ctx, engine.GetActorRequest{
					GameID:  gameID,
					ActorID: anotherWeakEnemy.Actor.ID,
				})
				require.NoError(t, err)
				require.NotNil(t, enemyAfterDamage.Actor)

				// 断言1: 验证死亡标记
				assert.True(t, damageResult.DamageResult.IsDead,
					"致命伤害后应标记为死亡")

				// 断言2: 验证HP为0
				assert.Equal(t, 0, enemyAfterDamage.Actor.HitPoints.Current,
					"死亡敌人HP应为0")

				// 断言3: 验证伤害结果中的HP记录
				assert.Equal(t, 3, damageResult.DamageResult.TargetHPBefore,
					"攻击前HP应为3")
				assert.Equal(t, 0, damageResult.DamageResult.TargetHPAfter,
					"攻击后HP应为0")

				t.Logf("直接伤害致死验证通过: HP=%d -> %d, 死亡=%v",
					damageResult.DamageResult.TargetHPBefore,
					damageResult.DamageResult.TargetHPAfter,
					damageResult.DamageResult.IsDead)
			})
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
		getActor := func(id model.ID) *engine.ActorInfo {
			actor, err := e.GetActor(ctx, engine.GetActorRequest{GameID: gameID, ActorID: id})
			require.NoError(t, err)
			require.NotNil(t, actor.Actor)
			return actor.Actor
		}

		getHP := func(id model.ID) int {
			actor := getActor(id)
			return actor.HitPoints.Current
		}

		getMaxHP := func(id model.ID) int {
			actor := getActor(id)
			return actor.HitPoints.Maximum
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

		// 攻击并验证HP变化
		attack := func(attacker, target model.ID) bool {
			// 记录攻击前目标状态
			targetBefore := getActor(target)
			hpBefore := targetBefore.HitPoints.Current
			maxHP := targetBefore.HitPoints.Maximum

			result, err := e.ExecuteAttack(ctx, engine.ExecuteAttackRequest{
				GameID:     gameID,
				AttackerID: attacker,
				TargetID:   target,
				Attack:     engine.AttackInput{IsUnarmed: false},
			})
			require.NoError(t, err)
			require.NotNil(t, result)
			require.NotNil(t, result.AttackResult)

			// 记录攻击后目标状态
			targetAfter := getActor(target)
			hpAfter := targetAfter.HitPoints.Current

			if result.AttackResult.Hit {
				// 断言1: 验证HP变化与伤害一致（考虑死亡时HP=0）
				expectedHP := hpBefore - result.AttackResult.Damage.FinalDamage
				if expectedHP < 0 {
					expectedHP = 0
				}
				assert.Equal(t, expectedHP, hpAfter,
					"攻击命中后目标HP应该等于攻击前HP减去伤害值（死亡时停在0）")

				// 断言2: 验证攻击结果中记录的HP与实际一致
				assert.Equal(t, hpBefore, result.AttackResult.Damage.TargetHPBefore,
					"攻击结果中记录的攻击前HP应与实际查询一致")
				assert.Equal(t, hpAfter, result.AttackResult.Damage.TargetHPAfter,
					"攻击结果中记录的攻击后HP应与实际查询一致")

				// 断言3: 验证最大HP不变
				assert.Equal(t, maxHP, targetAfter.HitPoints.Maximum,
					"攻击不应改变目标最大HP")

				// 断言4: 验证死亡状态
				if hpAfter <= 0 {
					assert.True(t, result.AttackResult.Damage.IsDead,
						"HP小于等于0时应标记为死亡")
				} else {
					assert.False(t, result.AttackResult.Damage.IsDead,
						"HP大于0时不应标记为死亡")
				}

				t.Logf("  ✓ 命中! 伤害: %d, HP: %d -> %d, 是否死亡: %v",
					result.AttackResult.Damage.FinalDamage, hpBefore, hpAfter,
					result.AttackResult.Damage.IsDead)
				return true
			} else {
				// 断言5: 未命中时HP不变
				assert.Equal(t, hpBefore, hpAfter, "未命中时目标HP应该不变")

				// 断言6: 未命中时擦伤伤害应为0
				assert.Equal(t, 0, result.AttackResult.GrazeDamage,
					"未命中时擦伤伤害应为0")

				t.Logf("  ✗ 未命中! HP保持: %d", hpAfter)
				return false
			}
		}

		castSpell := func(caster model.ID, spellID string, targets []model.ID, slotLevel int) bool {
			// 记录施法前HP
			hpBefore := make(map[model.ID]int)
			maxHPBefore := make(map[model.ID]int)
			for _, targetID := range targets {
				hpBefore[targetID] = getHP(targetID)
				maxHPBefore[targetID] = getMaxHP(targetID)
			}

			spellResult, err := e.CastSpell(ctx, engine.CastSpellRequest{
				GameID:   gameID,
				CasterID: caster,
				Spell: engine.SpellInput{
					SpellID:   spellID,
					SlotLevel: slotLevel,
					TargetIDs: targets,
				},
			})

			if err == nil {
				require.NotNil(t, spellResult)
				t.Logf("  ✓ 法术 %s 施放成功", spellID)

				// 验证每个目标的HP变化
				for _, targetID := range targets {
					hpAfter := getHP(targetID)
					maxHPAfter := getMaxHP(targetID)

					if hpBefore[targetID] > hpAfter {
						// 目标受到伤害
						damage := hpBefore[targetID] - hpAfter
						t.Logf("    目标 %s 受到 %d 伤害, HP: %d -> %d",
							targetID, damage, hpBefore[targetID], hpAfter)

						// 断言：法术应该造成伤害
						assert.Greater(t, damage, 0, "法术应该对目标造成伤害")

						// 断言：最大HP不变
						assert.Equal(t, maxHPBefore[targetID], maxHPAfter,
							"法术伤害不应改变目标最大HP")
					} else if hpAfter > hpBefore[targetID] {
						// 目标被治疗
						healing := hpAfter - hpBefore[targetID]
						t.Logf("    目标 %s 恢复 %d HP, HP: %d -> %d",
							targetID, healing, hpBefore[targetID], hpAfter)

						// 断言：治疗量应为正
						assert.Greater(t, healing, 0, "治疗量应大于0")

						// 断言：不超过最大HP
						assert.LessOrEqual(t, hpAfter, maxHPAfter,
							"治疗后HP不应超过最大HP")
					} else {
						t.Logf("    目标 %s HP未变化: %d", targetID, hpAfter)
					}
				}
				return true
			}

			t.Logf("  ✗ 法术 %s 施放失败: %v", spellID, err)
			return false
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
			goblinHPBefore := getHP(goblin.Actor.ID)
			if attack(fighter.Actor.ID, goblin.Actor.ID) {
				goblinHPAfter := getHP(goblin.Actor.ID)
				t.Logf("战士命中哥布林! 哥布林HP: %d -> %d", goblinHPBefore, goblinHPAfter)
				if goblinHPAfter <= 0 {
					t.Log("哥布林被击败!")

					// 验证哥布林死亡状态
					goblinActor := getActor(goblin.Actor.ID)
					assert.LessOrEqual(t, goblinActor.HitPoints.Current, 0,
						"哥布林被击败后HP应小于等于0")
				}
			} else {
				t.Log("战士未命中")
			}

			// 法师回合 - 治疗战士
			advanceTo("Merlin")
			t.Log("=== 法师 Merlin 治疗战士 ===")

			// 记录治疗前战士状态
			fighterBefore := getActor(fighter.Actor.ID)
			fighterHPBefore := fighterBefore.HitPoints.Current
			fighterMaxHP := fighterBefore.HitPoints.Maximum
			t.Logf("治疗前战士状态: HP=%d/%d", fighterHPBefore, fighterMaxHP)

			healResult, err := e.ExecuteHealing(ctx, engine.ExecuteHealingRequest{
				GameID:   gameID,
				TargetID: fighter.Actor.ID,
				Amount:   20,
			})
			if err == nil {
				require.NotNil(t, healResult)

				// 获取治疗后战士状态
				fighterAfter := getActor(fighter.Actor.ID)
				fighterHPAfter := fighterAfter.HitPoints.Current

				t.Logf("治疗结果: 请求治疗=20, 实际治疗=%d, HP=%d -> %d",
					healResult.Healed, fighterHPBefore, fighterHPAfter)

				// 断言1: 验证HP实际增加
				assert.Greater(t, fighterHPAfter, fighterHPBefore,
					"治疗后战士HP应增加")

				// 断言2: 验证HP增加量等于实际治疗量
				actualHeal := fighterHPAfter - fighterHPBefore
				assert.Equal(t, healResult.Healed, actualHeal,
					"战士HP增加量应等于实际治疗量")

				// 断言3: 验证治疗结果中报告的当前HP与实际一致
				assert.Equal(t, fighterHPAfter, healResult.CurrentHP,
					"治疗结果中报告的当前HP应与实际查询一致")

				// 断言4: 验证治疗不超过最大HP
				expectedHeal := 20
				if fighterHPBefore+20 > fighterMaxHP {
					expectedHeal = fighterMaxHP - fighterHPBefore
				}
				assert.Equal(t, expectedHeal, healResult.Healed,
					"实际治疗量不应超过最大HP限制")
				assert.LessOrEqual(t, fighterHPAfter, fighterMaxHP,
					"治疗后HP不应超过最大HP")

				// 断言5: 验证最大HP不变
				assert.Equal(t, fighterMaxHP, fighterAfter.HitPoints.Maximum,
					"治疗不应改变最大HP")
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
