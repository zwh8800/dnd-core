package model

// DruidFeatureHooks 德鲁伊特性钩子实现
type DruidFeatureHooks struct {
	Features *DruidFeatures
	Level    int // 德鲁伊等级
}

// ClassID 实现 ClassState 接口
func (h *DruidFeatureHooks) ClassID() ClassID {
	return ClassDruid
}

// OnAttackRoll 处理德鲁伊攻击加值
func (h *DruidFeatureHooks) OnAttackRoll(ctx *AttackContext) {
	// 月亮德鲁伊：野兽形态攻击（由野兽形态处理）
}

// OnDamageCalc 处理德鲁伊伤害计算
func (h *DruidFeatureHooks) OnDamageCalc(ctx *DamageContext) {
	// 德鲁伊通常不直接增加伤害
}

// OnACCalc 处理德鲁伊AC计算
func (h *DruidFeatureHooks) OnACCalc(ctx *ACContext) {
	// 德鲁伊有中等护甲和盾牌熟练
}

// OnSpellCalc 处理德鲁伊法术计算
func (h *DruidFeatureHooks) OnSpellCalc(ctx *SpellContext) {
	// 德鲁伊使用感知作为施法属性
	// 法术DC和攻击加值由基础系统计算
}

// GetAvailableActions 返回德鲁伊可用的特殊动作
func (h *DruidFeatureHooks) GetAvailableActions() []ActionTemplate {
	actions := []ActionTemplate{}

	// 荒野形态：2级获得
	if h.Level >= 2 && h.Features != nil && h.Features.WildShapeUsesRemaining > 0 {
		actions = append(actions, ActionTemplate{
			Type:          ActionCustom,
			Name:          "荒野形态",
			IsBonusAction: true,
			UsesPerRest:   h.Features.WildShapeMaxUses,
			CurrentUses:   h.Features.WildShapeUsesRemaining,
		})
	}

	return actions
}

// OnShortRest 短休时恢复德鲁伊资源
func (h *DruidFeatureHooks) OnShortRest() {
	// 荒野形态通过短休恢复
	if h.Features != nil {
		h.Features.WildShapeUsesRemaining = h.Features.WildShapeMaxUses
	}
}

// OnLongRest 长休时恢复德鲁伊资源
func (h *DruidFeatureHooks) OnLongRest() {
	if h.Features != nil {
		h.Features.WildShapeUsesRemaining = h.Features.WildShapeMaxUses
	}
}

// UpdateDruidFeatures 根据德鲁伊等级更新特性状态
func UpdateDruidFeatures(features *DruidFeatures, level int) {
	// 荒野形态次数
	if level >= 18 {
		features.WildShapeMaxUses = 3
	} else if level >= 8 {
		features.WildShapeMaxUses = 2
	} else if level >= 2 {
		features.WildShapeMaxUses = 2
	}

	features.WildShapeUsesRemaining = features.WildShapeMaxUses

	// 德鲁伊语言：1级获得
	if level >= 1 {
		features.HasDruidic = true
	}

	// 无甲变形：4级获得（月亮德鲁伊可变形为元素生物）
	if level >= 4 {
		features.HasWildShapeImprovement = true
	}

	// 时间躯体：18级获得
	if level >= 18 {
		features.HasBeastSpells = true
	}
}
