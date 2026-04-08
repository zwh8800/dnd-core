package rules

import (
	"fmt"

	"github.com/zwh8800/dnd-core/pkg/dice"
	"github.com/zwh8800/dnd-core/pkg/model"
)

// TerrainDifficulty 返回地形难度系数
func TerrainDifficulty(terrain model.TerrainType) float64 {
	switch terrain {
	case model.TerrainClear, model.TerrainGrassland:
		return 1.0
	case model.TerrainForest, model.TerrainDesert:
		return 0.75
	case model.TerrainMountain, model.TerrainSwamp, model.TerrainArctic:
		return 0.5
	default:
		return 1.0
	}
}

// CalculateTravelDistance 根据速度、地形、步伐计算日行距离（英里）
func CalculateTravelDistance(speed int, pace model.TravelPace, terrain model.TerrainType) float64 {
	// 基础日行距离 = 速度 / 5 * 8 小时
	baseDistance := float64(speed) / 5.0 * 8.0

	// 步伐修正
	paceMultiplier := 1.0
	switch pace {
	case model.TravelPaceFast:
		paceMultiplier = 1.33
	case model.TravelPaceNormal:
		paceMultiplier = 1.0
	case model.TravelPaceSlow:
		paceMultiplier = 0.67
	}

	// 地形修正
	terrainMultiplier := TerrainDifficulty(terrain)

	return baseDistance * paceMultiplier * terrainMultiplier
}

// ForagingCheck 觅食检定（WIS 生存检定）
func ForagingCheck(abilityScore int, proficiencyBonus int, hasProficiency bool) (*model.ForageResult, error) {
	roller := dice.New(0)

	// 生存检定 DC 15
	dc := 15

	// 掷骰
	rollResult, err := roller.Roll("1d20")
	if err != nil {
		return nil, err
	}

	// 计算加值
	abilityMod := AbilityModifier(abilityScore)
	profBonus := 0
	if hasProficiency {
		profBonus = proficiencyBonus
	}

	total := rollResult.Total + abilityMod + profBonus
	success := total >= dc

	result := &model.ForageResult{
		Success:   success,
		RollTotal: total,
		DC:        dc,
	}

	if success {
		// 成功：获得 1d6 份食物
		foodRoll, _ := roller.Roll("1d6")
		result.FoodObtained = foodRoll.Total
		result.WaterObtained = true
		result.Message = fmt.Sprintf("觅食成功！获得 %d 份食物和充足水源", result.FoodObtained)
	} else {
		result.FoodObtained = 0
		result.WaterObtained = false
		result.Message = "觅食失败，未找到足够的食物和水源"
	}

	return result, nil
}

// NavigationCheck 导航检定（WIS 生存检定）
func NavigationCheck(abilityScore int, proficiencyBonus int, hasProficiency bool, terrain model.TerrainType) (*model.NavigationCheck, error) {
	roller := dice.New(0)

	// 导航 DC 根据地 形变化
	dc := 10
	switch terrain {
	case model.TerrainClear, model.TerrainGrassland:
		dc = 10
	case model.TerrainForest, model.TerrainDesert:
		dc = 15
	case model.TerrainMountain, model.TerrainSwamp, model.TerrainArctic:
		dc = 20
	}

	// 掷骰
	rollResult, err := roller.Roll("1d20")
	if err != nil {
		return nil, err
	}

	// 计算加值
	abilityMod := AbilityModifier(abilityScore)
	profBonus := 0
	if hasProficiency {
		profBonus = proficiencyBonus
	}

	total := rollResult.Total + abilityMod + profBonus
	success := total >= dc

	result := &model.NavigationCheck{
		Success:   success,
		RollTotal: total,
		DC:        dc,
		Lost:      !success,
	}

	if success {
		result.Message = "导航成功，队伍沿着正确的方向前进"
	} else {
		result.Message = "导航失败，队伍迷路了！"
	}

	return result, nil
}

// EncounterCheck 随机遭遇检定
func EncounterCheck() (*model.EncounterCheck, error) {
	roller := dice.New(0)

	// 遭遇检定：1d6，1-2 表示遭遇
	roll, err := roller.Roll("1d6")
	if err != nil {
		return nil, err
	}

	encountered := roll.Total <= 2

	result := &model.EncounterCheck{
		Encountered: encountered,
		DC:          2,
		Roll:        roll.Total,
	}

	if encountered {
		// 随机决定遭遇类型
		typeRoll, _ := roller.Roll("1d4")
		switch typeRoll.Total {
		case 1:
			result.EncounterType = "monster"
			result.Message = "遭遇怪物！"
		case 2:
			result.EncounterType = "npc"
			result.Message = "遇到 NPC！"
		case 3:
			result.EncounterType = "treasure"
			result.Message = "发现宝藏！"
		case 4:
			result.EncounterType = "trap"
			result.Message = "触发陷阱！"
		}
	} else {
		result.Message = "旅途平静，没有遭遇"
	}

	return result, nil
}
