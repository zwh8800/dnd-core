package rules

import (
	"github.com/zwh8800/dnd-core/pkg/data"
	"github.com/zwh8800/dnd-core/pkg/model"
)

// CalculateLifestyleCost 计算生活方式费用
func CalculateLifestyleCost(tier model.LifestyleTier, days int) *model.LifestyleCost {
	lifestyleData := data.GetLifestyleData(tier)
	if lifestyleData == nil {
		return &model.LifestyleCost{
			DailyCost:   0,
			MonthlyCost: 0,
		}
	}

	return &model.LifestyleCost{
		DailyCost:   lifestyleData.DailyCost * days,
		MonthlyCost: lifestyleData.MonthlyCost,
	}
}

// DeductLifestyle 扣除生活方式费用
func DeductLifestyle(gold int, tier model.LifestyleTier, days int) (int, int, bool) {
	lifestyleData := data.GetLifestyleData(tier)
	if lifestyleData == nil {
		return gold, 0, false
	}

	cost := lifestyleData.DailyCost * days

	if gold >= cost {
		gold -= cost
		return gold, cost, true
	}

	// 金钱不足
	remaining := gold
	gold = 0
	return gold, remaining, false
}

// GetLifestyleDescription 获取生活方式描述
func GetLifestyleDescription(tier model.LifestyleTier) string {
	lifestyleData := data.GetLifestyleData(tier)
	if lifestyleData == nil {
		return "未知的生活方式"
	}
	return lifestyleData.Description
}
