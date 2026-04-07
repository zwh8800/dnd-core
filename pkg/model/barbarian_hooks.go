package model

// BarbarianFeatureHooks 野蛮人特性钩子实现
type BarbarianFeatureHooks struct {
	Features *BarbarianFeatures
	Level    int // 野蛮人等级
}

// ClassID 实现 ClassState 接口
func (h *BarbarianFeatureHooks) ClassID() ClassID {
	return ClassBarbarian
}

// OnAttackRoll 处理野蛮人攻击加值
func (h *BarbarianFeatureHooks) OnAttackRoll(ctx *AttackContext) {
	// 狂暴时力量攻击有优势（由战斗系统处理）
	// 这里可以添加其他攻击加值
}

// OnDamageCalc 处理野蛮人伤害计算
func (h *BarbarianFeatureHooks) OnDamageCalc(ctx *DamageContext) {
	// 狂暴时近战攻击获得额外伤害
	if h.Features != nil && h.Features.IsRaging && ctx.IsMelee {
		ctx.Bonus += h.GetRageDamageBonus()
	}
}

// OnACCalc 处理野蛮人AC计算
func (h *BarbarianFeatureHooks) OnACCalc(ctx *ACContext) {
	// 无甲防御：未着装护甲时，AC = 10 + 敏捷调整值 + 体质调整值
	if h.Features != nil && h.Features.HasUnarmoredDefense && !ctx.HasArmor {
		// 这个计算需要在外部获取属性值，这里只标记
		ctx.Bonus += h.Features.UnarmoredDefenseBonus
	}
}

// OnSpellCalc 野蛮人无施法能力（部分子职业除外）
func (h *BarbarianFeatureHooks) OnSpellCalc(ctx *SpellContext) {
	// 野蛮人不是施法者
}

// GetAvailableActions 返回野蛮人可用的特殊动作
func (h *BarbarianFeatureHooks) GetAvailableActions() []ActionTemplate {
	actions := []ActionTemplate{}

	// 狂暴：1级获得，附赠动作
	if h.Level >= 1 && h.Features != nil && h.Features.RageUsesRemaining > 0 {
		actions = append(actions, ActionTemplate{
			Type:          ActionCustom,
			Name:          "狂暴",
			IsBonusAction: true,
			UsesPerRest:   h.Features.RageMaxUses,
			CurrentUses:   h.Features.RageUsesRemaining,
		})
	}

	return actions
}

// OnShortRest 短休时恢复野蛮人资源
func (h *BarbarianFeatureHooks) OnShortRest() {
	// 野蛮人资源通过长休恢复
}

// OnLongRest 长休时恢复野蛮人资源
func (h *BarbarianFeatureHooks) OnLongRest() {
	if h.Features != nil {
		h.Features.RageUsesRemaining = h.Features.RageMaxUses
		h.Features.IsRaging = false
	}
}

// GetRageDamageBonus 根据等级获取狂暴伤害加值
func (h *BarbarianFeatureHooks) GetRageDamageBonus() int {
	if h.Level >= 16 {
		return 4
	} else if h.Level >= 9 {
		return 3
	} else if h.Level >= 1 {
		return 2
	}
	return 0
}

// UpdateBarbarianFeatures 根据野蛮人等级更新特性状态
func UpdateBarbarianFeatures(features *BarbarianFeatures, level int) {
	// 狂暴次数
	if level >= 18 {
		features.RageMaxUses = 6
	} else if level >= 16 {
		features.RageMaxUses = 5
	} else if level >= 12 {
		features.RageMaxUses = 4
	} else if level >= 9 {
		features.RageMaxUses = 3
	} else if level >= 1 {
		features.RageMaxUses = 2
	}

	features.RageUsesRemaining = features.RageMaxUses

	// 无甲防御：1级获得
	if level >= 1 {
		features.HasUnarmoredDefense = true
	}

	// 危险感知：7级获得
	if level >= 7 {
		features.HasDangerSense = true
	}

	// 鲁莽攻击：2级获得
	if level >= 2 {
		features.HasRecklessAttack = true
	}

	// 额外攻击：5级获得
	if level >= 5 {
		features.ExtraAttacks = 1
	}
}
