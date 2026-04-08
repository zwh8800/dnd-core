package rules

import "github.com/zwh8800/dnd-core/pkg/model"

// CalculateCraftingTime 计算制作时间
func CalculateCraftingTime(recipe model.CraftingRecipe, hasProficiency bool) int {
	days := recipe.TimeDays

	// 如果有工具熟练，时间减半
	if hasProficiency {
		days = days / 2
		if days < 1 {
			days = 1
		}
	}

	return days
}

// CalculateCraftingCost 计算制作成本
func CalculateCraftingCost(recipe model.CraftingRecipe) int {
	return recipe.Cost
}

// CanCraftRecipe 检查是否可以制作配方
func CanCraftRecipe(actorLevel int, hasToolsProficiency bool, hasSpell bool, recipe model.CraftingRecipe) bool {
	// 检查等级
	if actorLevel < recipe.MinLevel {
		return false
	}

	// 检查工具熟练
	if len(recipe.ToolsRequired) > 0 && !hasToolsProficiency {
		return false
	}

	// 检查法术
	if recipe.SpellRequired != "" && !hasSpell {
		return false
	}

	return true
}

// GetCraftingDC 获取制作DC
func GetCraftingDC(recipe model.CraftingRecipe) int {
	return recipe.DC
}
