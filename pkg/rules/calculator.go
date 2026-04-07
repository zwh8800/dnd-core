package rules

import "math"

// D&D 5e 核心规则计算函数（纯函数）

// AbilityModifier 计算属性修正值: floor((score - 10) / 2)
func AbilityModifier(score int) int {
	return int(math.Floor(float64(score-10) / 2))
}

// ProficiencyBonus 计算熟练加值
// 1-4级: +2
// 5-8级: +3
// 9-12级: +4
// 13-16级: +5
// 17-20级: +6
func ProficiencyBonus(totalLevel int) int {
	switch {
	case totalLevel >= 17:
		return 6
	case totalLevel >= 13:
		return 5
	case totalLevel >= 9:
		return 4
	case totalLevel >= 5:
		return 3
	case totalLevel >= 1:
		return 2
	default:
		return 0
	}
}

// SkillModifier 计算技能修正值
// 公式: 属性修正 + (熟练加值 if 熟练) + (2×熟练加值 if 专家) + 其他加值
func SkillModifier(abilityScore int, isProficient bool, isExpert bool, otherBonus int) int {
	mod := AbilityModifier(abilityScore)
	profBonus := ProficiencyBonus(0) // 调用方应传入总等级

	if isExpert {
		mod += profBonus * 2
	} else if isProficient {
		mod += profBonus
	}

	mod += otherBonus
	return mod
}

// SkillModifierWithLevel 带等级参数的技能修正值计算
func SkillModifierWithLevel(abilityScore int, level int, isProficient bool, isExpert bool, otherBonus int) int {
	mod := AbilityModifier(abilityScore)
	profBonus := ProficiencyBonus(level)

	if isExpert {
		mod += profBonus * 2
	} else if isProficient {
		mod += profBonus
	}

	mod += otherBonus
	return mod
}

// SpellSaveDC 计算法术豁免DC
// 公式: 8 + 熟练加值 + 施法属性修正
func SpellSaveDC(spellcastingAbilityScore int, proficiencyBonus int) int {
	return 8 + proficiencyBonus + AbilityModifier(spellcastingAbilityScore)
}

// SpellAttackBonus 计算法术攻击加值
// 公式: 熟练加值 + 施法属性修正
func SpellAttackBonus(spellcastingAbilityScore int, proficiencyBonus int) int {
	return proficiencyBonus + AbilityModifier(spellcastingAbilityScore)
}

// ArmorClass 计算护甲等级
// 公式取决于护甲类型
func ArmorClass(
	armorType string,
	armorBaseAC int,
	maxDexMod *int,
	dexModifier int,
	shieldBonus int,
	otherBonus int,
) int {
	ac := armorBaseAC

	// 应用敏捷修正（受护甲类型限制）
	dexMod := dexModifier
	if maxDexMod != nil {
		if dexMod > *maxDexMod {
			dexMod = *maxDexMod
		}
	}
	// 重型护甲不允许敏捷修正
	if armorType == "heavy" {
		dexMod = 0
	}

	ac += dexMod + shieldBonus + otherBonus
	return ac
}

// PassiveScore 计算被动检定分数
// 公式: 10 + 总修正值
func PassiveScore(totalModifier int) int {
	return 10 + totalModifier
}

// InitiativeModifier 计算先攻修正值（通常等于敏捷修正）
func InitiativeModifier(dexScore int) int {
	return AbilityModifier(dexScore)
}

// DeathSaveDC 死亡豁免的固定DC
func DeathSaveDC() int {
	return 10
}

// StabilizeDeathSaves 检查是否稳定（3次成功）
func StabilizeDeathSaves(successes int) bool {
	return successes >= 3
}

// IsDeadFromDeathSaves 检查是否死亡（3次失败）
func IsDeadFromDeathSaves(failures int) bool {
	return failures >= 3
}

// CalculateRestHPRecovery 计算休息时的HP恢复
// 短休: 使用生命骰恢复
// 长休: 恢复所有HP
func CalculateRestHPRecovery(isLongRest bool, currentHP int, maxHP int) int {
	if isLongRest {
		return maxHP - currentHP // 长休恢复所有HP
	}
	// 短休需要通过生命骰计算，这里返回0表示需要额外计算
	return 0
}

// CalculateHitDiceRecovery 计算长休后恢复的生命骰数量
// 规则：最多恢复总等级一半的生命骰（至少1个）
func CalculateHitDiceRecovery(totalLevel int) int {
	recovery := totalLevel / 2
	if recovery < 1 {
		return 1
	}
	return recovery
}

// CalculateExhaustionReduction 检查长休是否减少力竭
// 长休结束时，如果角色有1级或更高级力竭，则减少1级
func CalculateExhaustionReduction(currentExhaustion int) int {
	if currentExhaustion > 0 {
		return currentExhaustion - 1
	}
	return 0
}
