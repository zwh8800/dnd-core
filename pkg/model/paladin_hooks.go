package model

// PaladinFeatureHooks 圣武士特性钩子实现
type PaladinFeatureHooks struct {
	Features *PaladinFeatures
	Level    int // 圣武士等级
}

// ClassID 实现 ClassState 接口
func (h *PaladinFeatureHooks) ClassID() ClassID {
	return ClassPaladin
}

// OnAttackRoll 处理圣武士攻击加值
func (h *PaladinFeatureHooks) OnAttackRoll(ctx *AttackContext) {
	// 防御风格：AC+1（在OnACCalc处理）
}

// OnDamageCalc 处理圣武士伤害计算
func (h *PaladinFeatureHooks) OnDamageCalc(ctx *DamageContext) {
	// 至圣斩：消耗法术位增加伤害
	// 由战斗系统调用
}

// OnACCalc 处理圣武士AC计算
func (h *PaladinFeatureHooks) OnACCalc(ctx *ACContext) {
	// 防御风格：着装护甲时AC+1
	if h.Features != nil && h.Features.SelectedFightingStyle == FightingStyleDefense && ctx.HasArmor {
		ctx.Bonus += 1
	}
}

// OnSpellCalc 处理圣武士法术计算
func (h *PaladinFeatureHooks) OnSpellCalc(ctx *SpellContext) {
	// 圣武士使用魅力作为施法属性
	// 法术DC和攻击加值由基础系统计算
}

// GetAvailableActions 返回圣武士可用的特殊动作
func (h *PaladinFeatureHooks) GetAvailableActions() []ActionTemplate {
	actions := []ActionTemplate{}

	// 至圣斩：2级获得
	if h.Level >= 2 && h.Features != nil && h.Features.DivineSmiteAvailable {
		actions = append(actions, ActionTemplate{
			Type:          ActionCustom,
			Name:          "至圣斩",
			IsBonusAction: false,
			UsesPerRest:   99, // 消耗法术位，无次数限制
			CurrentUses:   99,
		})
	}

	// 圣疗：1级获得
	if h.Level >= 1 && h.Features != nil && h.Features.LayOnHandsRemaining > 0 {
		actions = append(actions, ActionTemplate{
			Type:          ActionCustom,
			Name:          "圣疗",
			IsBonusAction: false,
			UsesPerRest:   h.Features.LayOnHandsMax,
			CurrentUses:   h.Features.LayOnHandsRemaining,
		})
	}

	// 神恩：3级获得
	if h.Level >= 3 && h.Features != nil && h.Features.DivineSenseRemaining > 0 {
		actions = append(actions, ActionTemplate{
			Type:          ActionCustom,
			Name:          "神恩",
			IsBonusAction: false,
			UsesPerRest:   h.Features.DivineSenseMax,
			CurrentUses:   h.Features.DivineSenseRemaining,
		})
	}

	return actions
}

// OnShortRest 短休时恢复圣武士资源
func (h *PaladinFeatureHooks) OnShortRest() {
	// 圣武士大部分资源通过长休恢复
}

// OnLongRest 长休时恢复圣武士资源
func (h *PaladinFeatureHooks) OnLongRest() {
	if h.Features != nil {
		h.Features.LayOnHandsRemaining = h.Features.LayOnHandsMax
		h.Features.DivineSenseRemaining = h.Features.DivineSenseMax
	}
}

// UpdatePaladinFeatures 根据圣武士等级更新特性状态
func UpdatePaladinFeatures(features *PaladinFeatures, level int) {
	// 圣疗池
	features.LayOnHandsMax = level * 5
	features.LayOnHandsRemaining = features.LayOnHandsMax

	// 神恩次数
	if level >= 1 {
		features.DivineSenseMax = 1 + 3 // 1 + 魅力调整值
		features.DivineSenseRemaining = features.DivineSenseMax
	}

	// 至圣斩：2级获得
	if level >= 2 {
		features.DivineSmiteAvailable = true
	}

	// 灵光：6级获得（勇气灵光）
	if level >= 6 {
		features.HasAuraOfProtection = true
	}

	// 高等至圣斩：11级获得
	if level >= 11 {
		features.HasImprovedDivineSmite = true
	}

	// 清洁：14级获得
	if level >= 14 {
		features.HasCleansingTouch = true
	}
}
