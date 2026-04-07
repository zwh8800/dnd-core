package model

// ClericFeatureHooks 牧师特性钩子实现
type ClericFeatureHooks struct {
	Features *ClericFeatures
	Level    int // 牧师等级
}

// ClassID 实现 ClassState 接口
func (h *ClericFeatureHooks) ClassID() ClassID {
	return ClassCleric
}

// OnAttackRoll 处理牧师攻击加值
func (h *ClericFeatureHooks) OnAttackRoll(ctx *AttackContext) {
	// 战争领域：神性打击（需要额外实现）
}

// OnDamageCalc 处理牧师伤害计算
func (h *ClericFeatureHooks) OnDamageCalc(ctx *DamageContext) {
	// 引导神力：毁灭打击（需要额外实现）
}

// OnACCalc 处理牧师AC计算
func (h *ClericFeatureHooks) OnACCalc(ctx *ACContext) {
	// 战争领域：重甲熟练（由熟练系统处理）
}

// OnSpellCalc 处理牧师法术计算
func (h *ClericFeatureHooks) OnSpellCalc(ctx *SpellContext) {
	// 牧师使用感知作为施法属性
	// 法术DC和攻击加值由基础系统计算
}

// GetAvailableActions 返回牧师可用的特殊动作
func (h *ClericFeatureHooks) GetAvailableActions() []ActionTemplate {
	actions := []ActionTemplate{}

	// 引导神力：2级获得
	if h.Level >= 2 && h.Features != nil && h.Features.ChannelDivinityUsesRemaining > 0 {
		actions = append(actions, ActionTemplate{
			Type:          ActionCustom,
			Name:          "引导神力",
			IsBonusAction: false,
			UsesPerRest:   h.Features.ChannelDivinityMaxUses,
			CurrentUses:   h.Features.ChannelDivinityUsesRemaining,
		})
	}

	return actions
}

// OnShortRest 短休时恢复牧师资源
func (h *ClericFeatureHooks) OnShortRest() {
	// 引导神力通过短休恢复
	if h.Features != nil {
		h.Features.ChannelDivinityUsesRemaining = h.Features.ChannelDivinityMaxUses
	}
}

// OnLongRest 长休时恢复牧师资源
func (h *ClericFeatureHooks) OnLongRest() {
	if h.Features != nil {
		h.Features.ChannelDivinityUsesRemaining = h.Features.ChannelDivinityMaxUses
		h.Features.DivineInterventionUsesRemaining = h.Features.DivineInterventionMaxUses
	}
}

// UpdateClericFeatures 根据牧师等级更新特性状态
func UpdateClericFeatures(features *ClericFeatures, level int) {
	// 引导神力次数
	if level >= 18 {
		features.ChannelDivinityMaxUses = 3
	} else if level >= 6 {
		features.ChannelDivinityMaxUses = 2
	} else if level >= 2 {
		features.ChannelDivinityMaxUses = 1
	}

	features.ChannelDivinityUsesRemaining = features.ChannelDivinityMaxUses

	// 神圣干预：10级获得
	if level >= 10 {
		features.DivineInterventionMaxUses = 1
		features.DivineInterventionUsesRemaining = 1
	}

	// 破坏不死生物：1级获得
	if level >= 1 {
		features.HasDestroyUndead = true
	}
}
