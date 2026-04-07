package rules

import "github.com/zwh8800/dnd-core/pkg/model"

// AttackResult 代表攻击检定的结果
type AttackResult struct {
	Roll        int  // 攻击掷骰结果
	AttackBonus int  // 攻击加值
	Total       int  // 攻击总计
	TargetAC    int  // 目标AC
	Hit         bool // 是否命中
	IsCritical  bool // 是否暴击（天然20）
	IsFumble    bool // 是否失误（天然1）
}

// PerformAttackRoll 执行攻击检定
func PerformAttackRoll(roll int, attackBonus int, targetAC int) *AttackResult {
	isNat20 := roll == 20
	isNat1 := roll == 1
	total := roll + attackBonus

	result := &AttackResult{
		Roll:        roll,
		AttackBonus: attackBonus,
		Total:       total,
		TargetAC:    targetAC,
		Hit:         total >= targetAC || isNat20,
		IsCritical:  isNat20,
		IsFumble:    isNat1,
	}

	// 天然1总是失误（即使AC很低）
	if isNat1 {
		result.Hit = false
	}

	return result
}

// DamageCalculation 伤害计算结果
type DamageCalculation struct {
	BaseDamage      int                // 基础伤害
	CriticalBonus   int                // 暴击额外伤害
	Modifier        int                // 修正值
	RawTotal        int                // 原始总伤害
	Resistances     []model.DamageType // 应用的抗性
	Vulnerabilities []model.DamageType // 应用的弱点
	FinalDamage     int                // 最终伤害
}

// CalculateDamage 计算伤害
// 顺序：掷骰 → 修正 → 弱点(×2) → 抗性(÷2) → 免疫(→0)
func CalculateDamage(
	baseDiceTotal int,
	modifier int,
	damageType model.DamageType,
	resistances *model.DamageResistances,
	isCritical bool,
) *DamageCalculation {
	calc := &DamageCalculation{
		BaseDamage: baseDiceTotal,
		Modifier:   modifier,
	}

	// 暴击时额外掷一次伤害骰
	if isCritical {
		calc.CriticalBonus = baseDiceTotal // 简单实现：暴击伤害翻倍
	}

	// 原始总伤害
	calc.RawTotal = baseDiceTotal + calc.CriticalBonus + modifier

	// 先检查免疫
	if resistances != nil && resistances.HasImmunity(damageType) {
		calc.FinalDamage = 0
		calc.Resistances = append(calc.Resistances, damageType)
		return calc
	}

	// 应用弱点
	if resistances != nil && resistances.HasVulnerability(damageType) {
		calc.RawTotal *= 2
		calc.Vulnerabilities = append(calc.Vulnerabilities, damageType)
	}

	// 应用抗性
	if resistances != nil && resistances.HasResistance(damageType) {
		calc.RawTotal /= 2
		calc.Resistances = append(calc.Resistances, damageType)
	}

	calc.FinalDamage = calc.RawTotal
	if calc.FinalDamage < 0 {
		calc.FinalDamage = 0
	}

	return calc
}

// ApplyDamage 应用伤害到目标
func ApplyDamage(currentHP int, tempHP int, damage int) (newHP int, newTempHP int, damageToHP int) {
	// 先扣除临时HP
	if tempHP > 0 {
		if damage <= tempHP {
			newTempHP = tempHP - damage
			newHP = currentHP
			damageToHP = 0
			return
		}
		damage -= tempHP
		newTempHP = 0
	}

	// 扣除实际HP
	damageToHP = damage
	newHP = currentHP - damage
	if newHP < 0 {
		newHP = 0 // 不会降到0以下（死亡豁免另外处理）
	}

	return
}

// ApplyHealing 应用治疗
func ApplyHealing(currentHP int, maxHP int, healing int) (newHP int, overheal int) {
	newHP = currentHP + healing
	if newHP > maxHP {
		overheal = newHP - maxHP
		newHP = maxHP
	}
	return
}
