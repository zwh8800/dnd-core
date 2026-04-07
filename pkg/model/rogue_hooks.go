package model

// RogueFeatureHooks 游荡者特性钩子实现
type RogueFeatureHooks struct {
	Features *RogueFeatures
	Level    int // 游荡者等级
}

// ClassID 实现 ClassState 接口
func (h *RogueFeatureHooks) ClassID() ClassID {
	return ClassRogue
}

// OnAttackRoll 处理游荡者攻击加值
func (h *RogueFeatureHooks) OnAttackRoll(ctx *AttackContext) {
	// 偷袭：优势时攻击（由战斗系统处理）
}

// OnDamageCalc 处理游荡者伤害计算
func (h *RogueFeatureHooks) OnDamageCalc(ctx *DamageContext) {
	// 偷袭：额外伤害
	if h.Features != nil && h.Features.SneakAttackAvailable && ctx.IsMelee {
		ctx.Bonus += h.Features.SneakAttackDamage
	}
}

// OnACCalc 处理游荡者AC计算
func (h *RogueFeatureHooks) OnACCalc(ctx *ACContext) {
	// 游荡者无特殊AC加成
}

// OnSpellCalc 处理游荡者法术计算（诡术师子职业）
func (h *RogueFeatureHooks) OnSpellCalc(ctx *SpellContext) {
	// 游荡者通常不是施法者（诡术师除外）
}

// GetAvailableActions 返回游荡者可用的特殊动作
func (h *RogueFeatureHooks) GetAvailableActions() []ActionTemplate {
	actions := []ActionTemplate{}

	// 巧计动作：2级获得
	if h.Level >= 2 && h.Features != nil && h.Features.CunningActionAvailable {
		actions = append(actions, ActionTemplate{
			Type:          ActionCustom,
			Name:          "巧计动作 - 撤离",
			IsBonusAction: true,
			UsesPerRest:   99,
			CurrentUses:   99,
		})

		actions = append(actions, ActionTemplate{
			Type:          ActionCustom,
			Name:          "巧计动作 - 躲藏",
			IsBonusAction: true,
			UsesPerRest:   99,
			CurrentUses:   99,
		})

		actions = append(actions, ActionTemplate{
			Type:          ActionCustom,
			Name:          "巧计动作 - 疾走",
			IsBonusAction: true,
			UsesPerRest:   99,
			CurrentUses:   99,
		})
	}

	return actions
}

// OnShortRest 短休时恢复游荡者资源
func (h *RogueFeatureHooks) OnShortRest() {
	// 游荡者资源通常不通过短休恢复
}

// OnLongRest 长休时恢复游荡者资源
func (h *RogueFeatureHooks) OnLongRest() {
	if h.Features != nil {
		h.Features.SneakAttackAvailable = true
		h.Features.CunningActionAvailable = true
	}
}

// UpdateRogueFeatures 根据游荡者等级更新特性状态
func UpdateRogueFeatures(features *RogueFeatures, level int) {
	// 偷袭伤害
	if level >= 17 {
		features.SneakAttackDamage = 10 // 10d6
	} else if level >= 15 {
		features.SneakAttackDamage = 9 // 9d6
	} else if level >= 13 {
		features.SneakAttackDamage = 8 // 8d6
	} else if level >= 11 {
		features.SneakAttackDamage = 7 // 7d6
	} else if level >= 9 {
		features.SneakAttackDamage = 6 // 6d6
	} else if level >= 7 {
		features.SneakAttackDamage = 5 // 5d6
	} else if level >= 5 {
		features.SneakAttackDamage = 4 // 4d6
	} else if level >= 3 {
		features.SneakAttackDamage = 3 // 3d6
	} else if level >= 1 {
		features.SneakAttackDamage = 2 // 2d6
	}

	features.SneakAttackAvailable = true

	// 巧计动作：2级获得
	if level >= 2 {
		features.CunningActionAvailable = true
	}

	// 直觉闪避：5级获得
	if level >= 5 {
		features.HasUncannyDodge = true
	}

	// 闪避攻击：7级获得
	if level >= 7 {
		features.HasEvasion = true
	}

	// 可靠才能：11级获得
	if level >= 11 {
		features.HasReliableTalent = true
	}

	// 模糊思维：14级获得
	if level >= 14 {
		features.HasBlindsense = true
	}

	// 滑溜头脑：15级获得
	if level >= 15 {
		features.HasSlipperyMind = true
	}

	// 躲避大师：18级获得
	if level >= 18 {
		features.HasElusive = true
	}

	// 神偷：20级获得
	if level >= 20 {
		features.HasStrokeOfLuck = true
	}
}
