package model

// FeatureHook 职业特性钩子接口
// 用于实现职业特性与战斗、法术等系统的交互
type FeatureHook interface {
	// OnAttackRoll 攻击掷骰时调用,可修改攻击加值和暴击范围
	OnAttackRoll(ctx *AttackContext)

	// OnDamageCalc 伤害计算时调用,可修改伤害值
	OnDamageCalc(ctx *DamageContext)

	// OnACCalc 护甲等级计算时调用,可修改AC
	OnACCalc(ctx *ACContext)

	// OnSpellCalc 法术计算时调用,可修改法术DC/攻击加值
	OnSpellCalc(ctx *SpellContext)

	// GetAvailableActions 返回可用的特殊动作
	GetAvailableActions() []ActionTemplate

	// OnShortRest 短休时调用,恢复资源
	OnShortRest()

	// OnLongRest 长休时调用,恢复资源
	OnLongRest()
}

// AttackContext 攻击上下文
type AttackContext struct {
	BaseBonus     int    // 基础攻击加值
	Bonus         int    // 额外加值(可修改)
	WeaponType    string // 武器类型
	IsRanged      bool   // 是否远程
	CriticalRange int    // 暴击范围(默认20,勇士范型可改为19或18)
}

// DamageContext 伤害上下文
type DamageContext struct {
	BaseDamage int        // 基础伤害
	Bonus      int        // 额外伤害(可修改)
	DamageType DamageType // 伤害类型
	IsMelee    bool       // 是否近战
	IsRanged   bool       // 是否远程
}

// ACContext 护甲等级上下文
type ACContext struct {
	BaseAC    int  // 基础AC
	Bonus     int  // 额外AC(可修改)
	HasShield bool // 是否持盾
	HasArmor  bool // 是否着装护甲
}

// SpellContext 法术上下文
type SpellContext struct {
	SpellSaveDC      int // 法术豁免DC
	SpellAttackBonus int // 法术攻击加值
}

// ActionTemplate 特殊动作模板
type ActionTemplate struct {
	Type          ActionType // 动作类型
	Name          string     // 动作名称
	IsBonusAction bool       // 是否附赠动作
	IsFreeAction  bool       // 是否自由动作
	UsesPerRest   int        // 每rest可用次数
	CurrentUses   int        // 当前剩余次数
}

// FighterFeatureHooks 战士特性钩子实现
type FighterFeatureHooks struct {
	Features *FighterFeatures
	Level    int // 战士等级
}

// ClassID 实现 ClassState 接口
func (h *FighterFeatureHooks) ClassID() ClassID {
	return ClassFighter
}

// OnAttackRoll 处理战士攻击加值
func (h *FighterFeatureHooks) OnAttackRoll(ctx *AttackContext) {
	// 箭术：远程武器攻击检定+2
	if h.Features.SelectedFightingStyle == FightingStyleArchery && ctx.IsRanged {
		ctx.Bonus += 2
	}

	// 勇士范型：改进暴击范围
	if h.Features.SelectedArchetype == MartialArchetypeChampion {
		if h.Level >= 15 {
			ctx.CriticalRange = 18 // 18-20暴击
		} else if h.Level >= 7 {
			ctx.CriticalRange = 19 // 19-20暴击
		}
	}
}

// OnDamageCalc 处理战士伤害计算
func (h *FighterFeatureHooks) OnDamageCalc(ctx *DamageContext) {
	// 对决：单手持握近战武器且未持用其他武器时+2伤害
	if h.Features.SelectedFightingStyle == FightingStyleDueling && ctx.IsMelee {
		// 注意：这里需要检查武器配置，简化处理
		ctx.Bonus += 2
	}
}

// OnACCalc 处理战士AC计算
func (h *FighterFeatureHooks) OnACCalc(ctx *ACContext) {
	// 防御：着装护甲时AC+1
	if h.Features.SelectedFightingStyle == FightingStyleDefense && ctx.HasArmor {
		ctx.Bonus += 1
	}
}

// OnSpellCalc 战士无施法能力，空实现
func (h *FighterFeatureHooks) OnSpellCalc(ctx *SpellContext) {
	// 战士不是施法者（奥法骑士除外，需要额外实现）
}

// GetAvailableActions 返回战士可用的特殊动作
func (h *FighterFeatureHooks) GetAvailableActions() []ActionTemplate {
	actions := []ActionTemplate{}

	// 回气：1级获得，附赠动作
	if h.Level >= 1 && h.Features.SecondWindUsed < h.Features.SecondWindMax {
		actions = append(actions, ActionTemplate{
			Type:          ActionSecondWind,
			Name:          "复苏之风",
			IsBonusAction: true,
			UsesPerRest:   h.Features.SecondWindMax,
			CurrentUses:   h.Features.SecondWindMax - h.Features.SecondWindUsed,
		})
	}

	// 动作如潮：2级获得，自由动作
	if h.Level >= 2 && h.Features.ActionSurgeUsed < h.Features.ActionSurgeMax {
		actions = append(actions, ActionTemplate{
			Type:         ActionActionSurge,
			Name:         "动作如潮",
			IsFreeAction: true,
			UsesPerRest:  h.Features.ActionSurgeMax,
			CurrentUses:  h.Features.ActionSurgeMax - h.Features.ActionSurgeUsed,
		})
	}

	return actions
}

// OnShortRest 短休时恢复战士资源
func (h *FighterFeatureHooks) OnShortRest() {
	h.Features.SecondWindUsed = 0
	h.Features.ActionSurgeUsed = 0
	h.Features.IndomitableUsed = 0
}

// OnLongRest 长休时恢复战士资源
func (h *FighterFeatureHooks) OnLongRest() {
	h.Features.SecondWindUsed = 0
	h.Features.ActionSurgeUsed = 0
	h.Features.IndomitableUsed = 0
}

// UpdateFighterFeatures 根据战士等级更新特性状态
func UpdateFighterFeatures(features *FighterFeatures, level int) {
	features.SecondWindMax = 1

	// 动作如潮：2级1次，17级2次
	if level >= 17 {
		features.ActionSurgeMax = 2
	} else if level >= 2 {
		features.ActionSurgeMax = 1
	}

	// 额外攻击：5级+1，11级+2，20级+3
	if level >= 20 {
		features.ExtraAttacks = 3
	} else if level >= 11 {
		features.ExtraAttacks = 2
	} else if level >= 5 {
		features.ExtraAttacks = 1
	}

	// 不屈：9级1次，13级2次，17级3次
	if level >= 17 {
		features.IndomitableMax = 3
	} else if level >= 13 {
		features.IndomitableMax = 2
	} else if level >= 9 {
		features.IndomitableMax = 1
	}
}
