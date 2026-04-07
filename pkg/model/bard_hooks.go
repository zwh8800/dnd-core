package model

// BardFeatureHooks 吟游诗人特性钩子实现
type BardFeatureHooks struct {
	Features *BardFeatures
	Level    int // 吟游诗人等级
}

// ClassID 实现 ClassState 接口
func (h *BardFeatureHooks) ClassID() ClassID {
	return ClassBard
}

// OnAttackRoll 处理吟游诗人攻击加值
func (h *BardFeatureHooks) OnAttackRoll(ctx *AttackContext) {
	// 剑诗人学院：近战武器攻击加值（需要额外实现）
}

// OnDamageCalc 处理吟游诗人伤害计算
func (h *BardFeatureHooks) OnDamageCalc(ctx *DamageContext) {
	// 诗人通常不直接增加伤害
}

// OnACCalc 处理吟游诗人AC计算
func (h *BardFeatureHooks) OnACCalc(ctx *ACContext) {
	// 诗人无特殊AC加成
}

// OnSpellCalc 处理吟游诗人法术计算
func (h *BardFeatureHooks) OnSpellCalc(ctx *SpellContext) {
	// 诗人使用魅力作为施法属性
	// 法术DC和攻击加值由基础系统计算
}

// GetAvailableActions 返回吟游诗人可用的特殊动作
func (h *BardFeatureHooks) GetAvailableActions() []ActionTemplate {
	actions := []ActionTemplate{}

	// 吟游 Inspiration：1级获得
	if h.Level >= 1 && h.Features != nil && h.Features.InspirationUsesRemaining > 0 {
		actions = append(actions, ActionTemplate{
			Type:          ActionCustom,
			Name:          "吟游激励",
			IsBonusAction: true,
			UsesPerRest:   h.Features.InspirationMaxUses,
			CurrentUses:   h.Features.InspirationUsesRemaining,
		})
	}

	// 割讥：3级获得（剑诗人学院）
	if h.Level >= 3 && h.Features != nil && h.Features.SlashingUseRemaining > 0 {
		actions = append(actions, ActionTemplate{
			Type:          ActionAttack,
			Name:          " slashing 攻击",
			IsBonusAction: false,
			UsesPerRest:   h.Features.SlashingMaxUses,
			CurrentUses:   h.Features.SlashingUseRemaining,
		})
	}

	return actions
}

// OnShortRest 短休时恢复吟游诗人资源
func (h *BardFeatureHooks) OnShortRest() {
	// 吟游诗人资源通过长休恢复
}

// OnLongRest 长休时恢复吟游诗人资源
func (h *BardFeatureHooks) OnLongRest() {
	if h.Features != nil {
		h.Features.InspirationUsesRemaining = h.Features.InspirationMaxUses
		h.Features.SlashingUseRemaining = h.Features.SlashingMaxUses
	}
}

// UpdateBardFeatures 根据吟游诗人等级更新特性状态
func UpdateBardFeatures(features *BardFeatures, level int) {
	// 吟游激励次数
	if level >= 15 {
		features.InspirationMaxUses = 5
	} else if level >= 10 {
		features.InspirationMaxUses = 4
	} else if level >= 5 {
		features.InspirationMaxUses = 3
	} else if level >= 1 {
		features.InspirationMaxUses = 2
	}

	features.InspirationUsesRemaining = features.InspirationMaxUses

	// 割讥次数（剑诗人学院）
	if level >= 3 {
		features.SlashingMaxUses = 1
		features.SlashingUseRemaining = 1
	}

	// 专家：3级获得
	if level >= 3 {
		features.HasExpertise = true
	}

	// 反制法术：6级获得
	if level >= 6 {
		features.HasCountercharm = true
	}
}
