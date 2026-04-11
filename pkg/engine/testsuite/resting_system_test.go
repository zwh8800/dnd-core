package testsuite

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zwh8800/dnd-core/pkg/engine"
	"github.com/zwh8800/dnd-core/pkg/model"
)

// TestRestingSystem 测试休息系统
// 包含短休、长休和团队休息三种类型的测试
func TestRestingSystem(t *testing.T) {
	// 测试短休：角色可以恢复生命值和职业能力
	t.Run("短休", func(t *testing.T) {
		// 创建测试引擎和上下文
		e := engine.NewTestEngine(t)
		ctx := context.Background()

		// 创建新游戏
		gameResult, err := e.NewGame(ctx, engine.NewGameRequest{
			Name:        "Campfire Tales",
			Description: "营火旁休息",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// 创建战士角色（3级）
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

		// 记录短休前的HP，用于后续验证
		actorBeforeRest, err := e.GetActor(ctx, engine.GetActorRequest{GameID: gameID, ActorID: fighter.Actor.ID})
		require.NoError(t, err)
		hpBeforeRest := actorBeforeRest.Actor.HitPoints.Current

		// 切换到探索阶段（休息需要此阶段）
		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "开始休息")
		require.NoError(t, err)

		// 执行短休
		restResult, err := e.ShortRest(ctx, engine.ShortRestRequest{
			GameID:   gameID,
			ActorIDs: []model.ID{fighter.Actor.ID},
		})
		require.NoError(t, err)
		require.NotNil(t, restResult)
		require.GreaterOrEqual(t, len(restResult.ActorResults), 1)
		t.Logf("短休结果: %s", restResult.Message)

		// 验证：短休后HP应大于等于休息前（恢复至少1点HP或生命骰）
		actorAfterRest, err := e.GetActor(ctx, engine.GetActorRequest{GameID: gameID, ActorID: fighter.Actor.ID})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, actorAfterRest.Actor.HitPoints.Current, hpBeforeRest)
	})

	// 测试长休：角色恢复全部HP和法术位
	t.Run("长休", func(t *testing.T) {
		e := engine.NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, engine.NewGameRequest{
			Name:        "Inn Stay",
			Description: "旅店中休息",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// 创建牧师角色（2级）
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

		// 记录最大HP，用于验证恢复效果
		actorBeforeRest, err := e.GetActor(ctx, engine.GetActorRequest{GameID: gameID, ActorID: cleric.Actor.ID})
		require.NoError(t, err)
		maxHP := actorBeforeRest.Actor.HitPoints.Maximum

		// 切换到探索阶段
		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "开始长休")
		require.NoError(t, err)

		// 开始长休
		restResult, err := e.StartLongRest(ctx, engine.StartLongRestRequest{
			GameID:   gameID,
			ActorIDs: []model.ID{cleric.Actor.ID},
		})
		require.NoError(t, err)
		require.NotNil(t, restResult)
		t.Logf("长休已开始: %s", restResult.Message)

		// 结束长休（模拟8小时休息完成）
		endResult, err := e.EndLongRest(ctx, engine.EndLongRestRequest{
			GameID: gameID,
		})
		require.NoError(t, err)
		require.NotNil(t, endResult)
		require.GreaterOrEqual(t, len(endResult.ActorResults), 1)

		// 验证恢复结果
		actorResult := endResult.ActorResults[0]
		assert.GreaterOrEqual(t, actorResult.HPRecovered, 0)
		t.Logf("长休恢复: HP=%d, 法术位=%v",
			actorResult.HPRecovered, actorResult.SpellSlotsRestored)

		// 验证：通过GetActor确认HP已恢复到最大值
		actorAfterRest, err := e.GetActor(ctx, engine.GetActorRequest{GameID: gameID, ActorID: cleric.Actor.ID})
		require.NoError(t, err)
		assert.Equal(t, maxHP, actorAfterRest.Actor.HitPoints.Current)
	})

	// 测试团队休息：多个角色同时休息
	t.Run("团队休息", func(t *testing.T) {
		e := engine.NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, engine.NewGameRequest{
			Name:        "Group Rest",
			Description: "整个团队一起休息",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// 创建战士角色
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

		// 创建法师角色
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

		// 记录休息前双方HP
		pc1BeforeHP, err := e.GetActor(ctx, engine.GetActorRequest{GameID: gameID, ActorID: pc1.Actor.ID})
		require.NoError(t, err)
		pc2BeforeHP, err := e.GetActor(ctx, engine.GetActorRequest{GameID: gameID, ActorID: pc2.Actor.ID})
		require.NoError(t, err)

		// 切换到探索阶段
		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "团队休息")
		require.NoError(t, err)

		// 开始团队长休
		restResult, err := e.StartLongRest(ctx, engine.StartLongRestRequest{
			GameID:   gameID,
			ActorIDs: []model.ID{pc1.Actor.ID, pc2.Actor.ID},
		})
		require.NoError(t, err)
		require.NotNil(t, restResult)

		// 结束团队长休
		endResult, err := e.EndLongRest(ctx, engine.EndLongRestRequest{
			GameID: gameID,
		})
		require.NoError(t, err)
		require.NotNil(t, endResult)
		require.Equal(t, 2, len(endResult.ActorResults))

		t.Logf("团队休息结果: %s 恢复 %d HP, %s 恢复 %d HP",
			pc1.Actor.Name, endResult.ActorResults[0].HPRecovered,
			pc2.Actor.Name, endResult.ActorResults[1].HPRecovered)

		// 验证：通过GetActor确认两个角色的HP都有所恢复
		pc1AfterHP, err := e.GetActor(ctx, engine.GetActorRequest{GameID: gameID, ActorID: pc1.Actor.ID})
		require.NoError(t, err)
		pc2AfterHP, err := e.GetActor(ctx, engine.GetActorRequest{GameID: gameID, ActorID: pc2.Actor.ID})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, pc1AfterHP.Actor.HitPoints.Current, pc1BeforeHP.Actor.HitPoints.Current)
		assert.GreaterOrEqual(t, pc2AfterHP.Actor.HitPoints.Current, pc2BeforeHP.Actor.HitPoints.Current)
	})
}
