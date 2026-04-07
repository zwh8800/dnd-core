package engine

import (
	"fmt"

	"github.com/zwh8800/dnd-core/pkg/data"
	"github.com/zwh8800/dnd-core/pkg/model"
)

// LoadMonster 从模板创建怪物实例
func (e *Engine) LoadMonster(templateID string) (*model.Enemy, error) {
	statBlock, exists := data.GlobalRegistry.GetMonster(templateID)
	if !exists {
		return nil, fmt.Errorf("monster template not found: %s", templateID)
	}

	return CreateEnemyFromStatBlock(statBlock)
}

// CreateEnemyFromStatBlock 从怪物数据块创建 Enemy 实例
func CreateEnemyFromStatBlock(statBlock *model.MonsterStatBlock) (*model.Enemy, error) {
	// 创建 Enemy 实例
	enemy := &model.Enemy{
		Actor: model.Actor{
			ID:              model.NewID(),
			Type:            model.ActorTypeEnemy,
			Name:            statBlock.Name,
			Description:     statBlock.Description,
			Size:            statBlock.Size,
			CreatureType:    statBlock.CreatureType,
			Speed:           statBlock.Speed.Walk,
			Speeds:          statBlock.Speed,
			ArmorClass:      statBlock.ArmorClass,
			ChallengeRating: statBlock.ChallengeRating,
			Exhaustion:      0,
		},
		StatBlock: statBlock,
	}

	// 复制属性分数
	enemy.Actor.AbilityScores = statBlock.AbilityScores

	// 设置生命值
	enemy.Actor.HitPoints = model.HitPoints{
		Current: statBlock.HitPointsAverage,
		Maximum: statBlock.HitPointsAverage,
	}

	// 复制状态免疫
	enemy.ConditionImmunities = make([]model.ConditionType, len(statBlock.ConditionImmunities))
	copy(enemy.ConditionImmunities, statBlock.ConditionImmunities)

	// 复制伤害免疫/抗性/易伤
	enemy.DamageImmunities = make([]model.DamageImmunity, len(statBlock.DamageImmunities))
	copy(enemy.DamageImmunities, statBlock.DamageImmunities)
	enemy.DamageResistances = make([]model.DamageImmunity, len(statBlock.DamageResistances))
	copy(enemy.DamageResistances, statBlock.DamageResistances)
	enemy.DamageVulnerabilities = make([]model.DamageImmunity, len(statBlock.DamageVulnerabilities))
	copy(enemy.DamageVulnerabilities, statBlock.DamageVulnerabilities)

	// 初始化传说动作追踪
	if statBlock.HasLegendaryActions() {
		enemy.LegendaryActionsRemaining = statBlock.LegendaryActionsPerRound
	}

	// 初始化充能动作
	for i, action := range statBlock.Actions {
		if action.Recharge != nil && action.Recharge.UsesPerDay > 0 {
			if enemy.ActionRecharges == nil {
				enemy.ActionRecharges = make(map[int]int)
			}
			enemy.ActionRecharges[i] = action.Recharge.UsesPerDay
		}
	}

	return enemy, nil
}

// GetMonsterActions 获取怪物可用的动作列表
func GetMonsterActions(monster *model.Enemy) []model.MonsterAction {
	actions := make([]model.MonsterAction, 0)

	// 添加常规动作
	for _, action := range monster.StatBlock.Actions {
		if action.Recharge == nil || action.Recharge.IsRecharged() {
			actions = append(actions, action)
		}
	}

	// 添加附赠动作
	for _, action := range monster.StatBlock.BonusActions {
		if action.Recharge == nil || action.Recharge.IsRecharged() {
			actions = append(actions, action)
		}
	}

	// 添加反应
	for _, action := range monster.StatBlock.Reactions {
		actions = append(actions, action)
	}

	// 添加传说动作（如果有剩余）
	if monster.LegendaryActionsRemaining > 0 {
		actions = append(actions, monster.StatBlock.LegendaryActions...)
	}

	return actions
}

// RechargeMonsterActions 为怪物掷骰充能动作
func RechargeMonsterActions(monster *model.Enemy) {
	for _, action := range monster.StatBlock.Actions {
		if action.Recharge != nil && action.Recharge.RollRange[0] > 0 {
			// 模拟 d6 掷骰 (这里简化处理，实际应使用骰子系统)
			roll := 4 // 简化：假设成功充能
			if roll >= action.Recharge.RollRange[0] && roll <= action.Recharge.RollRange[1] {
				// 充能成功
			}
		}
	}

	// 恢复传说动作
	if monster.StatBlock.HasLegendaryActions() {
		monster.LegendaryActionsRemaining = monster.StatBlock.LegendaryActionsPerRound
	}
}

// UseLegendaryAction 使用一个传说动作
func UseLegendaryAction(monster *model.Enemy) error {
	if monster.LegendaryActionsRemaining <= 0 {
		return fmt.Errorf("no legendary actions remaining")
	}
	monster.LegendaryActionsRemaining--
	return nil
}

// UseRechargeAction 使用充能动作
func UseRechargeAction(monster *model.Enemy, actionIndex int) error {
	if monster.ActionRecharges == nil {
		return fmt.Errorf("action does not use charges")
	}

	uses, exists := monster.ActionRecharges[actionIndex]
	if !exists {
		return fmt.Errorf("action does not use charges")
	}

	if uses <= 0 {
		return fmt.Errorf("no charges remaining for this action")
	}

	monster.ActionRecharges[actionIndex]--
	return nil
}
