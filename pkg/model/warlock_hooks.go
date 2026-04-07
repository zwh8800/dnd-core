package model

// WarlockFeatureHooks 邪术师特性钩子实现
type WarlockFeatureHooks struct {
	Features *WarlockFeatures
	Level    int // 邪术师等级
}

// ClassID 实现 ClassState 接口
func (h *WarlockFeatureHooks) ClassID() ClassID {
	return ClassWarlock
}

// OnAttackRoll 处理邪术师攻击加值
func (h *WarlockFeatureHooks) OnAttackRoll(ctx *AttackContext) {
	// 魔能祈唤：地狱火召唤（增加攻击）
}

// OnDamageCalc 处理邪术师伤害计算
func (h *WarlockFeatureHooks) OnDamageCalc(ctx *DamageContext) {
	// 邪术师通常不直接增加武器伤害
	// 魔能爆伤害由法术系统处理
}

// OnACCalc 处理邪术师AC计算
func (h *WarlockFeatureHooks) OnACCalc(ctx *ACContext) {
	// 邪术师无特殊AC加成
}

// OnSpellCalc 处理邪术师法术计算
func (h *WarlockFeatureHooks) OnSpellCalc(ctx *SpellContext) {
	// 邪术师使用魅力作为施法属性
	// 法术DC和攻击加值由基础系统计算
}

// GetAvailableActions 返回邪术师可用的特殊动作
func (h *WarlockFeatureHooks) GetAvailableActions() []ActionTemplate {
	actions := []ActionTemplate{}

	// 魔能爆：戏法
	if h.Level >= 1 {
		actions = append(actions, ActionTemplate{
			Type:          ActionCastSpell,
			Name:          "魔能爆",
			IsBonusAction: false,
			UsesPerRest:   99, // 戏法无限
			CurrentUses:   99,
		})
	}

	return actions
}

// OnShortRest 短休时恢复邪术师资源
func (h *WarlockFeatureHooks) OnShortRest() {
	// 邪术师法术位通过短休恢复
	if h.Features != nil {
		h.Features.SpellSlotsUsed = 0
	}
}

// OnLongRest 长休时恢复邪术师资源
func (h *WarlockFeatureHooks) OnLongRest() {
	if h.Features != nil {
		h.Features.SpellSlotsUsed = 0
		h.Features.InvocationsUsed = 0
	}
}

// UpdateWarlockFeatures 根据邪术师等级更新特性状态
func UpdateWarlockFeatures(features *WarlockFeatures, level int) {
	// 邪术师法术位（特殊：所有法术位总是最高环阶）
	if level >= 17 {
		features.SpellSlotLevel = 5
		features.SpellSlotsMax = 4
	} else if level >= 11 {
		features.SpellSlotLevel = 5
		features.SpellSlotsMax = 3
	} else if level >= 9 {
		features.SpellSlotLevel = 5
		features.SpellSlotsMax = 2
	} else if level >= 7 {
		features.SpellSlotLevel = 4
		features.SpellSlotsMax = 2
	} else if level >= 5 {
		features.SpellSlotLevel = 3
		features.SpellSlotsMax = 2
	} else if level >= 1 {
		features.SpellSlotLevel = 1
		features.SpellSlotsMax = 1
	}

	features.SpellSlotsUsed = 0

	// 魔能祈唤：2级获得
	if level >= 2 {
		features.HasEldritchInvocations = true
		features.InvocationsKnown = 2
	}

	// 额外魔能祈唤
	if level >= 15 {
		features.InvocationsKnown = 7
	} else if level >= 12 {
		features.InvocationsKnown = 6
	} else if level >= 9 {
		features.InvocationsKnown = 5
	} else if level >= 7 {
		features.InvocationsKnown = 4
	} else if level >= 5 {
		features.InvocationsKnown = 3
	}

	// 契约恩赐：3级获得
	if level >= 3 {
		features.HasPactBoon = true
	}

	// 神秘恢复：5级获得
	if level >= 5 {
		features.HasMysticArcanum = true
	}
}
