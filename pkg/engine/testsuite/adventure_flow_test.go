package testsuite

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zwh8800/dnd-core/pkg/engine"
	"github.com/zwh8800/dnd-core/pkg/model"
)

// TestFullAdventureFlow 测试完整冒险流程
// 模拟从角色创建到战斗的完整游戏流程
func TestFullAdventureFlow(t *testing.T) {
	// 测试完整冒险流程：创建角色 -> 创建任务 -> 创建场景 -> 遭遇敌人 -> 战斗
	t.Run("从创建到战斗的完整冒险", func(t *testing.T) {
		// 创建测试引擎和上下文
		e := engine.NewTestEngine(t)
		ctx := context.Background()

		// 创建新游戏
		gameResult, err := e.NewGame(ctx, engine.NewGameRequest{
			Name:        "The Lost Temple",
			Description: "探索古代神庙的冒险",
		})
		require.NoError(t, err)
		require.NotNil(t, gameResult)
		gameID := gameResult.Game.ID

		// 阶段1：创建玩家角色
		t.Run("阶段1: 创建玩家角色", func(t *testing.T) {
			// 创建精灵游侠
			pc1, err := e.CreatePC(ctx, engine.CreatePCRequest{
				GameID: gameID,
				PC: &engine.PlayerCharacterInput{
					Name:       "Aelindra",
					Race:       "Elf",
					Class:      "游侠",
					Level:      3,
					Background: "Outlander",
					AbilityScores: engine.AbilityScoresInput{
						Strength:     12,
						Dexterity:    16,
						Constitution: 14,
						Intelligence: 10,
						Wisdom:       14,
						Charisma:     10,
					},
				},
			})
			require.NoError(t, err)
			require.NotNil(t, pc1)
			assert.NotEmpty(t, pc1.Actor.ID)

			// 创建矮人战士
			pc2, err := e.CreatePC(ctx, engine.CreatePCRequest{
				GameID: gameID,
				PC: &engine.PlayerCharacterInput{
					Name:       "Thorin",
					Race:       "Dwarf",
					Class:      "战士",
					Level:      3,
					Background: "Soldier",
					AbilityScores: engine.AbilityScoresInput{
						Strength:     16,
						Dexterity:    12,
						Constitution: 15,
						Intelligence: 10,
						Wisdom:       12,
						Charisma:     8,
					},
				},
			})
			require.NoError(t, err)
			require.NotNil(t, pc2)
			assert.NotEmpty(t, pc2.Actor.ID)

			t.Logf("已创建角色: %s (HP: %d), %s (HP: %d)",
				pc1.Actor.Name, pc1.Actor.HitPoints.Current,
				pc2.Actor.Name, pc2.Actor.HitPoints.Current)
		})

		// 阶段2：创建任务给予者NPC和任务
		t.Run("阶段2: 创建任务给予者和任务", func(t *testing.T) {
			// 创建NPC：智慧的长者
			npc, err := e.CreateNPC(ctx, engine.CreateNPCRequest{
				GameID: gameID,
				NPC: &engine.NPCInput{
					Name:        "Elder Theron",
					Description: "村庄里的智慧老人",
					Size:        string(model.SizeMedium),
					Speed:       25,
					AbilityScores: engine.AbilityScoresInput{
						Strength:     10,
						Dexterity:    10,
						Constitution: 12,
						Intelligence: 16,
						Wisdom:       14,
						Charisma:     14,
					},
				},
			})
			require.NoError(t, err)
			require.NotNil(t, npc)

			// 创建任务
			questResult, err := e.CreateQuest(ctx, engine.CreateQuestRequest{
				GameID:      gameID,
				Name:        "The Lost Temple",
				Description: "在森林中找到古代神庙并恢复神圣文物",
				GiverID:     npc.Actor.ID,
				GiverName:   "Elder Theron",
				Objectives: []engine.ObjectiveInput{
					{ID: "find_temple", Description: "找到古代神庙", Required: 1},
					{ID: "defeat_guardians", Description: "击败神庙守护者", Required: 3},
					{ID: "retrieve_artifact", Description: "取回神圣文物", Required: 1},
				},
			})
			require.NoError(t, err)
			require.NotNil(t, questResult)
			assert.NotEmpty(t, questResult.Quest.ID)
			assert.Equal(t, "The Lost Temple", questResult.Quest.Name)

			t.Logf("已创建任务: %s", questResult.Quest.Name)
		})

		// 阶段3：创建探索场景
		t.Run("阶段3: 创建探索场景", func(t *testing.T) {
			// 创建村庄场景
			villageScene, err := e.CreateScene(ctx, engine.CreateSceneRequest{
				GameID:      gameID,
				Name:        "Willowbrook Village",
				Description: "Darkwood边缘的宁静农业村落",
				SceneType:   model.SceneTypeOutdoor,
			})
			require.NoError(t, err)
			require.NotNil(t, villageScene)

			// 创建森林场景
			forestScene, err := e.CreateScene(ctx, engine.CreateSceneRequest{
				GameID:      gameID,
				Name:        "Darkwood Forest",
				Description: "黑暗神秘的森林",
				SceneType:   model.SceneTypeWilderness,
			})
			require.NoError(t, err)
			require.NotNil(t, forestScene)

			// 创建神庙场景
			templeScene, err := e.CreateScene(ctx, engine.CreateSceneRequest{
				GameID:      gameID,
				Name:        "Temple of the Ancients",
				Description: "被藤蔓覆盖的古代神庙",
				SceneType:   model.SceneTypeDungeon,
			})
			require.NoError(t, err)
			require.NotNil(t, templeScene)

			// 创建场景连接：村庄 -> 森林
			err = e.AddSceneConnection(ctx, engine.AddSceneConnectionRequest{
				GameID:        gameID,
				SceneID:       villageScene.Scene.ID,
				TargetSceneID: forestScene.Scene.ID,
				Description:   "通往森林的土路",
				Locked:        false,
				DC:            0,
				Hidden:        false,
			})
			require.NoError(t, err)

			// 创建场景连接：森林 -> 神庙（隐藏且需要DC15）
			err = e.AddSceneConnection(ctx, engine.AddSceneConnectionRequest{
				GameID:        gameID,
				SceneID:       forestScene.Scene.ID,
				TargetSceneID: templeScene.Scene.ID,
				Description:   "通往地下室的古老石阶",
				Locked:        true,
				DC:            15,
				Hidden:        true,
			})
			require.NoError(t, err)

			t.Logf("已创建场景及连接")
		})

		// 阶段4：添加敌人怪物
		t.Run("阶段4: 添加敌人怪物", func(t *testing.T) {
			// 创建骷髅战士
			skeleton1, err := e.CreateEnemy(ctx, engine.CreateEnemyRequest{
				GameID: gameID,
				Enemy: &engine.EnemyInput{
					Name:        "Skeleton Warrior",
					Description: "手持生锈剑的活化骷髅",
					Size:        string(model.SizeMedium),
					Speed:       30,
					AbilityScores: engine.AbilityScoresInput{
						Strength:     14,
						Dexterity:    12,
						Constitution: 10,
						Intelligence: 6,
						Wisdom:       8,
						Charisma:     5,
					},
					ChallengeRating: "1/4",
					HitPoints:       13,
					ArmorClass:      13,
				},
			})
			require.NoError(t, err)
			require.NotNil(t, skeleton1)
			assert.NotEmpty(t, skeleton1.Actor.ID)

			// 创建骷髅弓箭手
			skeleton2, err := e.CreateEnemy(ctx, engine.CreateEnemyRequest{
				GameID: gameID,
				Enemy: &engine.EnemyInput{
					Name:        "Skeleton Archer",
					Description: "持弓的活化骷髅",
					Size:        string(model.SizeMedium),
					Speed:       30,
					AbilityScores: engine.AbilityScoresInput{
						Strength:     10,
						Dexterity:    14,
						Constitution: 10,
						Intelligence: 6,
						Wisdom:       8,
						Charisma:     5,
					},
					ChallengeRating: "1/4",
					HitPoints:       10,
					ArmorClass:      13,
				},
			})
			require.NoError(t, err)
			require.NotNil(t, skeleton2)
			assert.NotEmpty(t, skeleton2.Actor.ID)

			t.Logf("已创建敌人: %s (AC: %d, HP: %d), %s (AC: %d, HP: %d)",
				skeleton1.Actor.Name, skeleton1.Actor.ArmorClass, skeleton1.Actor.HitPoints.Current,
				skeleton2.Actor.Name, skeleton2.Actor.ArmorClass, skeleton2.Actor.HitPoints.Current)
		})

		// 阶段5：列出角色用于战斗
		t.Run("阶段5: 列出角色用于战斗", func(t *testing.T) {
			// 列出所有玩家角色
			actorsResult, err := e.ListActors(ctx, engine.ListActorsRequest{
				GameID: gameID,
				Filter: &engine.ActorFilter{
					Types: []model.ActorType{model.ActorTypePC},
				},
			})
			require.NoError(t, err)
			require.GreaterOrEqual(t, len(actorsResult.Actors), 2)
			t.Logf("发现 %d 名玩家角色", len(actorsResult.Actors))

			// 列出所有敌人
			enemiesResult, err := e.ListActors(ctx, engine.ListActorsRequest{
				GameID: gameID,
				Filter: &engine.ActorFilter{
					Types: []model.ActorType{model.ActorTypeEnemy},
				},
			})
			require.NoError(t, err)
			require.GreaterOrEqual(t, len(enemiesResult.Actors), 2)
			t.Logf("发现 %d 名敌人", len(enemiesResult.Actors))
		})

		// 阶段6：切换到探索阶段
		t.Run("阶段6: 切换到探索阶段", func(t *testing.T) {
			result, err := e.SetPhase(ctx, gameID, model.PhaseExploration, "开始探索森林")
			require.NoError(t, err)
			require.NotNil(t, result)

			// 验证：通过GetPhase确认阶段已切换
			phase, err := e.GetPhase(ctx, gameID)
			require.NoError(t, err)
			assert.Equal(t, model.PhaseExploration, phase)
			t.Logf("阶段已设置为: %s", phase)
		})

		// 阶段7：开始带突袭的战斗
		t.Run("阶段7: 开始带突袭的战斗", func(t *testing.T) {
			// 获取所有参与者
			actorsResult, err := e.ListActors(ctx, engine.ListActorsRequest{
				GameID: gameID,
				Filter: &engine.ActorFilter{
					Types: []model.ActorType{model.ActorTypePC, model.ActorTypeEnemy},
				},
			})
			require.NoError(t, err)
			require.GreaterOrEqual(t, len(actorsResult.Actors), 4)

			// 构建参与者列表，前2个角色为潜行方
			var participantIDs []model.ID
			var stealthyIDs []model.ID
			for i, actor := range actorsResult.Actors {
				participantIDs = append(participantIDs, actor.ID)
				if i < 2 {
					stealthyIDs = append(stealthyIDs, actor.ID)
				}
			}

			// 开始带突袭判定的战斗
			surpriseResult, err := e.StartCombatWithSurprise(ctx, engine.StartCombatWithSurpriseRequest{
				GameID:       gameID,
				SceneID:      actorsResult.Actors[0].SceneID,
				StealthySide: stealthyIDs,
				Observers:    participantIDs[len(stealthyIDs):],
			})
			require.NoError(t, err)
			require.NotNil(t, surpriseResult)
			assert.Equal(t, model.CombatStatusActive, surpriseResult.Combat.Status)

			t.Logf("战斗已开始，包含 %d 名参与者（可能有突袭）",
				len(participantIDs))
		})
	})

	// 测试带技能检定的冒险
	t.Run("带技能检定的冒险", func(t *testing.T) {
		e := engine.NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, engine.NewGameRequest{
			Name:        "Diplomatic Mission",
			Description: "社交和探索为核心的冒险",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// 创建吟游诗人角色
		pc, err := e.CreatePC(ctx, engine.CreatePCRequest{
			GameID: gameID,
			PC: &engine.PlayerCharacterInput{
				Name:       "Elara",
				Race:       "Half-Elf",
				Class:      "吟游诗人",
				Level:      2,
				Background: "Entertainer",
				AbilityScores: engine.AbilityScoresInput{
					Strength:     10,
					Dexterity:    14,
					Constitution: 10,
					Intelligence: 12,
					Wisdom:       10,
					Charisma:     16,
				},
			},
		})
		require.NoError(t, err)
		require.NotNil(t, pc)

		// 切换到探索阶段
		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "开始冒险")
		require.NoError(t, err)

		// 测试探索中进行属性检定
		t.Run("探索中进行属性检定", func(t *testing.T) {
			checkResult, err := e.PerformAbilityCheck(ctx, engine.AbilityCheckRequest{
				GameID:  gameID,
				ActorID: pc.Actor.ID,
				Ability: model.AbilityCharisma,
				DC:      12,
				Reason:  "说服守卫",
			})
			require.NoError(t, err)
			require.NotNil(t, checkResult)
			assert.Equal(t, model.AbilityCharisma, checkResult.Ability)
			t.Logf("魅力检定: 掷骰=%d, 总值=%d, DC=%d, 成功=%v",
				checkResult.Roll.Total, checkResult.RollTotal, checkResult.DC, checkResult.Success)
		})

		// 测试进行技能检定
		t.Run("进行技能检定", func(t *testing.T) {
			checkResult, err := e.PerformSkillCheck(ctx, engine.SkillCheckRequest{
				GameID:  gameID,
				ActorID: pc.Actor.ID,
				Skill:   model.SkillPersuasion,
				DC:      14,
				Reason:  "向商人推销",
			})
			require.NoError(t, err)
			require.NotNil(t, checkResult)
			assert.Equal(t, model.SkillPersuasion, checkResult.Skill)
			t.Logf("说服检定: 掷骰=%d, 总值=%d, 成功=%v",
				checkResult.Roll.Total, checkResult.RollTotal, checkResult.Success)
		})

		// 测试进行豁免检定
		t.Run("进行豁免检定", func(t *testing.T) {
			saveResult, err := e.PerformSavingThrow(ctx, engine.SavingThrowRequest{
				GameID:  gameID,
				ActorID: pc.Actor.ID,
				Ability: model.AbilityCharisma,
				DC:      10,
				Reason:  "抵抗魅惑效果",
			})
			require.NoError(t, err)
			require.NotNil(t, saveResult)
			assert.Equal(t, model.AbilityCharisma, saveResult.Ability)
			t.Logf("魅力豁免: 掷骰=%d, 总值=%d, 成功=%v",
				saveResult.Roll.Total, saveResult.RollTotal, saveResult.Success)
		})
	})
}
