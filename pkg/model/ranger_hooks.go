package model

// RangerFeatureHooks 游侠特性钩子实现
type RangerFeatureHooks struct {
	Features *RangerFeatures
	Level    int // 游侠等级
}

// ClassID 实现 ClassState 接口
func (h *RangerFeatureHooks) ClassID() ClassID {
	return ClassRanger
}

// OnAttackRoll 处理游侠攻击加值
func (h *RangerFeatureHooks) OnAttackRoll(ctx *AttackContext) {
	// 宿敌：对特定类型生物攻击有优势（由战斗系统处理）
}

// OnDamageCalc 处理游侠伤害计算
func (h *RangerFeatureHooks) OnDamageCalc(ctx *DamageContext) {
	// 游侠通常不直接增加伤害
}

// OnACCalc 处理游侠AC计算
func (h *RangerFeatureHooks) OnACCalc(ctx *ACContext) {
	// 游侠有中等护甲和盾牌熟练
}

// OnSpellCalc 处理游侠法术计算
func (h *RangerFeatureHooks) OnSpellCalc(ctx *SpellContext) {
	// 游侠使用感知作为施法属性
	// 法术DC和攻击加值由基础系统计算
}

// GetAvailableActions 返回游侠可用的特殊动作
func (h *RangerFeatureHooks) GetAvailableActions() []ActionTemplate {
	actions := []ActionTemplate{}

	// 猎杀：3级获得（猎人游侠）
	if h.Level >= 3 && h.Features != nil && h.Features.HuntersPreyAvailable {
		actions = append(actions, ActionTemplate{
			Type:          ActionAttack,
			Name:          "猎杀攻击",
			IsBonusAction: false,
			UsesPerRest:   1,
			CurrentUses:   1,
		})
	}

	return actions
}

// OnShortRest 短休时恢复游侠资源
func (h *RangerFeatureHooks) OnShortRest() {
	// 游侠资源通过长休恢复
}

// OnLongRest 长休时恢复游侠资源
func (h *RangerFeatureHooks) OnLongRest() {
	if h.Features != nil {
		h.Features.HuntersPreyAvailable = true
	}
}

// UpdateRangerFeatures 根据游侠等级更新特性状态
func UpdateRangerFeatures(features *RangerFeatures, level int) {
	// 宿敌：1级获得
	if level >= 1 {
		features.HasFavoredEnemy = true
	}

	// 自然探索：1级获得
	if level >= 1 {
		features.HasNaturalExplorer = true
	}

	// 战斗风格：2级获得
	if level >= 2 {
		features.HasFightingStyle = true
	}

	// 法术施放：2级获得
	if level >= 2 {
		features.HasSpellcasting = true
	}

	// 猎杀：3级获得
	if level >= 3 {
		features.HuntersPreyAvailable = true
	}

	// 额外攻击：5级获得
	if level >= 5 {
		features.ExtraAttacks = 1
	}

	// 陆地漫步：8级获得
	if level >= 8 {
		features.HasLandStride = true
	}

	// 隐藏：10级获得
	if level >= 10 {
		features.HasHideInPlainSight = true
	}
}
