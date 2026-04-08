package rules

import (
	"fmt"

	"github.com/zwh8800/dnd-core/pkg/model"
)

// ExhaustionLevel 力竭等级效果
type ExhaustionLevel struct {
	Level       int    `json:"level"`
	Effect      string `json:"effect"`
	Description string `json:"description"`
}

// ExhaustionTable 力竭效果表（SRD 5.2.1）
var ExhaustionTable = []ExhaustionLevel{
	{
		Level:       1,
		Effect:      "检定劣势",
		Description: "属性检定、攻击掷骰和豁免检定具有劣势",
	},
	{
		Level:       2,
		Effect:      "速度减半",
		Description: "速度减半",
	},
	{
		Level:       3,
		Effect:      "攻击劣势",
		Description: "攻击掷骰具有劣势（叠加1级的效果）",
	},
	{
		Level:       4,
		Effect:      "HP上限减半",
		Description: "HP最大值减半",
	},
	{
		Level:       5,
		Effect:      "速度降为0",
		Description: "速度降为0，无法移动",
	},
	{
		Level:       6,
		Effect:      "死亡",
		Description: "死亡",
	},
}

// GetExhaustionEffect 获取力竭等级效果
func GetExhaustionEffect(level int) string {
	if level < 1 || level > 6 {
		return "无效力竭等级"
	}
	return ExhaustionTable[level-1].Effect
}

// GetExhaustionDescription 获取力竭等级详细描述
func GetExhaustionDescription(level int) string {
	if level < 1 || level > 6 {
		return "无效力竭等级"
	}
	return ExhaustionTable[level-1].Description
}

// ApplyExhaustionEffects 应用力竭效果
func ApplyExhaustionEffects(currentExhaustion int, actor *model.Actor) []string {
	effects := []string{}

	if currentExhaustion < 1 || currentExhaustion > 6 {
		return effects
	}

	// 累积所有等级的效果
	for level := 1; level <= currentExhaustion; level++ {
		effects = append(effects, fmt.Sprintf("力竭%d级: %s", level, GetExhaustionEffect(level)))
	}

	// 6级力竭直接死亡
	if currentExhaustion >= 6 {
		actor.HitPoints.Current = 0
	}

	return effects
}

// RemoveExhaustion 移除力竭等级
func RemoveExhaustion(currentLevel int, removeAmount int) int {
	newLevel := currentLevel - removeAmount
	if newLevel < 0 {
		newLevel = 0
	}
	return newLevel
}

// HasLongRestRemovedExhaustion 长休是否移除了力竭
func HasLongRestRemovedExhaustion(currentLevel int) int {
	// 长休移除1级力竭
	return RemoveExhaustion(currentLevel, 1)
}
