package model

// WizardFeatureHooks 法师特性钩子实现
type WizardFeatureHooks struct {
	Features *WizardFeatures
	Level    int // 法师等级
}

// ClassID 实现 ClassState 接口
func (h *WizardFeatureHooks) ClassID() ClassID {
	return ClassWizard
}

// OnAttackRoll 处理法师攻击加值
func (h *WizardFeatureHooks) OnAttackRoll(ctx *AttackContext) {
	// 法师通常不使用武器攻击
}

// OnDamageCalc 处理法师伤害计算
func (h *WizardFeatureHooks) OnDamageCalc(ctx *DamageContext) {
	// 法师通过法术造成伤害
}

// OnACCalc 处理法师AC计算
func (h *WizardFeatureHooks) OnACCalc(ctx *ACContext) {
	// 法师无特殊AC加成
}

// OnSpellCalc 处理法师法术计算
func (h *WizardFeatureHooks) OnSpellCalc(ctx *SpellContext) {
	// 法师使用智力作为施法属性
	// 法术DC和攻击加值由基础系统计算
}

// GetAvailableActions 返回法师可用的特殊动作
func (h *WizardFeatureHooks) GetAvailableActions() []ActionTemplate {
	actions := []ActionTemplate{}

	// 奥术恢复：1级获得
	if h.Level >= 1 && h.Features != nil && h.Features.ArcaneRecoveryUsed < h.Features.ArcaneRecoveryMax {
		actions = append(actions, ActionTemplate{
			Type:          ActionCustom,
			Name:          "奥术恢复",
			IsBonusAction: false,
			UsesPerRest:   h.Features.ArcaneRecoveryMax,
			CurrentUses:   h.Features.ArcaneRecoveryMax - h.Features.ArcaneRecoveryUsed,
		})
	}

	return actions
}

// OnShortRest 短休时恢复法师资源
func (h *WizardFeatureHooks) OnShortRest() {
	// 奥术恢复通过长休恢复
}

// OnLongRest 长休时恢复法师资源
func (h *WizardFeatureHooks) OnLongRest() {
	if h.Features != nil {
		h.Features.ArcaneRecoveryUsed = 0
	}
}

// UpdateWizardFeatures 根据法师等级更新特性状态
func UpdateWizardFeatures(features *WizardFeatures, level int) {
	// 奥术恢复：可恢复的法术位等级总和
	if level >= 1 {
		features.ArcaneRecoveryMax = 1
	}

	// 奥术传统：2级获得
	if level >= 2 {
		features.HasArcaneTradition = true
	}

	// 法术书：1级获得
	if level >= 1 {
		features.HasSpellbook = true
		features.SpellbookSpells = 6 // 初始法术
	}

	// 额外法术：每升一级添加2个法术
	if level >= 1 {
		features.SpellbookSpells = 6 + (level-1)*2
	}

	// 法术精通：20级获得
	if level >= 20 {
		features.HasSpellMastery = true
	}

	// 标志法术：18级获得
	if level >= 18 {
		features.HasSignatureSpells = true
	}
}
