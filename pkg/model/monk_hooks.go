package model

// MonkFeatureHooks 武僧特性钩子实现
type MonkFeatureHooks struct {
	Features *MonkFeatures
	Level    int // 武僧等级
}

// ClassID 实现 ClassState 接口
func (h *MonkFeatureHooks) ClassID() ClassID {
	return ClassMonk
}

// OnAttackRoll 处理武僧攻击加值
func (h *MonkFeatureHooks) OnAttackRoll(ctx *AttackContext) {
	// 武僧使用敏捷或力量进行攻击
	// 武术：可以使用敏捷替代力量
}

// OnDamageCalc 处理武僧伤害计算
func (h *MonkFeatureHooks) OnDamageCalc(ctx *DamageContext) {
	// 武术：徒手打击伤害
	if h.Features != nil && ctx.IsMelee {
		ctx.Bonus += h.Features.MartialArtsBonus
	}

	// 疾风连击：附赠攻击
}

// OnACCalc 处理武僧AC计算
func (h *MonkFeatureHooks) OnACCalc(ctx *ACContext) {
	// 无甲防御：未着装护甲时，AC = 10 + 敏捷调整值 + 感知调整值
	if h.Features != nil && h.Features.HasUnarmoredDefense && !ctx.HasArmor {
		ctx.Bonus += h.Features.UnarmoredDefenseBonus
	}
}

// OnSpellCalc 处理武僧法术计算（部分子职业有法术）
func (h *MonkFeatureHooks) OnSpellCalc(ctx *SpellContext) {
	// 武僧通常不是施法者（四象宗除外）
}

// GetAvailableActions 返回武僧可用的特殊动作
func (h *MonkFeatureHooks) GetAvailableActions() []ActionTemplate {
	actions := []ActionTemplate{}

	// 气点：2级获得
	if h.Level >= 2 && h.Features != nil && h.Features.KiPointsRemaining > 0 {
		actions = append(actions, ActionTemplate{
			Type:          ActionCustom,
			Name:          "气点 - 疾风连击",
			IsBonusAction: true,
			UsesPerRest:   h.Features.KiPointsMax,
			CurrentUses:   h.Features.KiPointsRemaining,
		})

		actions = append(actions, ActionTemplate{
			Type:          ActionCustom,
			Name:          "气点 - 患者防御",
			IsBonusAction: true,
			UsesPerRest:   h.Features.KiPointsMax,
			CurrentUses:   h.Features.KiPointsRemaining,
		})
	}

	return actions
}

// OnShortRest 短休时恢复武僧资源
func (h *MonkFeatureHooks) OnShortRest() {
	// 气点通过短休或长休恢复
	if h.Features != nil {
		h.Features.KiPointsRemaining = h.Features.KiPointsMax
	}
}

// OnLongRest 长休时恢复武僧资源
func (h *MonkFeatureHooks) OnLongRest() {
	if h.Features != nil {
		h.Features.KiPointsRemaining = h.Features.KiPointsMax
	}
}

// UpdateMonkFeatures 根据武僧等级更新特性状态
func UpdateMonkFeatures(features *MonkFeatures, level int) {
	// 气点数量
	features.KiPointsMax = level
	features.KiPointsRemaining = level

	// 武术伤害
	if level >= 17 {
		features.MartialArtsBonus = 4
	} else if level >= 11 {
		features.MartialArtsBonus = 3
	} else if level >= 5 {
		features.MartialArtsBonus = 2
	} else if level >= 1 {
		features.MartialArtsBonus = 1
	}

	// 无甲防御：1级获得
	if level >= 1 {
		features.HasUnarmoredDefense = true
	}

	// 无甲移动：2级获得
	if level >= 2 {
		features.HasUnarmoredMovement = true
	}

	// 轻功：2级获得
	if level >= 2 {
		features.HasStepOfTheWind = true
	}

	// 震慑拳：5级获得
	if level >= 5 {
		features.HasStunningStrike = true
	}

	// 闪避攻击：7级获得
	if level >= 7 {
		features.HasEvasion = true
	}

	// 直觉闪避：7级获得
	if level >= 7 {
		features.HasUncannyDodge = false // 这是游侠的特性
	}
}
