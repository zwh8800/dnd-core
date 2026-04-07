package model

// SorcererFeatureHooks 术士特性钩子实现
type SorcererFeatureHooks struct {
	Features *SorcererFeatures
	Level    int // 术士等级
}

// ClassID 实现 ClassState 接口
func (h *SorcererFeatureHooks) ClassID() ClassID {
	return ClassSorcerer
}

// OnAttackRoll 处理术士攻击加值
func (h *SorcererFeatureHooks) OnAttackRoll(ctx *AttackContext) {
	// 术士通常不使用武器攻击
}

// OnDamageCalc 处理术士伤害计算
func (h *SorcererFeatureHooks) OnDamageCalc(ctx *DamageContext) {
	// 超魔：增效法术（增加伤害）
	// 由超魔系统处理
}

// OnACCalc 处理术士AC计算
func (h *SorcererFeatureHooks) OnACCalc(ctx *ACContext) {
	// 术士无特殊AC加成
}

// OnSpellCalc 处理术士法术计算
func (h *SorcererFeatureHooks) OnSpellCalc(ctx *SpellContext) {
	// 术士使用魅力作为施法属性
	// 法术DC和攻击加值由基础系统计算
}

// GetAvailableActions 返回术士可用的特殊动作
func (h *SorcererFeatureHooks) GetAvailableActions() []ActionTemplate {
	actions := []ActionTemplate{}

	// 超魔：3级获得
	if h.Level >= 3 && h.Features != nil && h.Features.SorceryPointsRemaining > 0 {
		actions = append(actions, ActionTemplate{
			Type:          ActionCustom,
			Name:          "超魔 - 增效",
			IsBonusAction: false,
			UsesPerRest:   h.Features.SorceryPointsMax,
			CurrentUses:   h.Features.SorceryPointsRemaining,
		})

		actions = append(actions, ActionTemplate{
			Type:          ActionCustom,
			Name:          "超魔 - 延时",
			IsBonusAction: false,
			UsesPerRest:   h.Features.SorceryPointsMax,
			CurrentUses:   h.Features.SorceryPointsRemaining,
		})
	}

	return actions
}

// OnShortRest 短休时恢复术士资源
func (h *SorcererFeatureHooks) OnShortRest() {
	// 术力点通过长休恢复
}

// OnLongRest 长休时恢复术士资源
func (h *SorcererFeatureHooks) OnLongRest() {
	if h.Features != nil {
		h.Features.SorceryPointsRemaining = h.Features.SorceryPointsMax
	}
}

// UpdateSorcererFeatures 根据术士等级更新特性状态
func UpdateSorcererFeatures(features *SorcererFeatures, level int) {
	// 术力点
	features.SorceryPointsMax = level
	features.SorceryPointsRemaining = level

	// 超魔：3级获得
	if level >= 3 {
		features.HasMetamagic = true
		features.MetamagicOptionsKnown = 2
	}

	// 额外超魔选项
	if level >= 10 {
		features.MetamagicOptionsKnown = 4
	} else if level >= 7 {
		features.MetamagicOptionsKnown = 3
	}

	// 魔法之源：20级获得
	if level >= 20 {
		features.HasSorcerousRestoration = true
	}
}
