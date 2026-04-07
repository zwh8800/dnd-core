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

// BarbarianFeatures 野蛮人特性状态
type BarbarianFeatures struct {
	IsRaging              bool `json:"is_raging"`
	RageMaxUses           int  `json:"rage_max_uses"`
	RageUsesRemaining     int  `json:"rage_uses_remaining"`
	HasUnarmoredDefense   bool `json:"has_unarmored_defense"`
	UnarmoredDefenseBonus int  `json:"unarmored_defense_bonus"`
	HasDangerSense        bool `json:"has_danger_sense"`
	HasRecklessAttack     bool `json:"has_reckless_attack"`
	ExtraAttacks          int  `json:"extra_attacks"`
}

// BardFeatures 吟游诗人特性状态
type BardFeatures struct {
	InspirationMaxUses       int  `json:"inspiration_max_uses"`
	InspirationUsesRemaining int  `json:"inspiration_uses_remaining"`
	HasExpertise             bool `json:"has_expertise"`
	HasCountercharm          bool `json:"has_countercharm"`
	SlashingMaxUses          int  `json:"slashing_max_uses"`
	SlashingUseRemaining     int  `json:"slashing_use_remaining"`
}

// ClericFeatures 牧师特性状态
type ClericFeatures struct {
	ChannelDivinityMaxUses          int  `json:"channel_divinity_max_uses"`
	ChannelDivinityUsesRemaining    int  `json:"channel_divinity_uses_remaining"`
	DivineInterventionMaxUses       int  `json:"divine_intervention_max_uses"`
	DivineInterventionUsesRemaining int  `json:"divine_intervention_uses_remaining"`
	HasDestroyUndead                bool `json:"has_destroy_undead"`
}

// DruidFeatures 德鲁伊特性状态
type DruidFeatures struct {
	WildShapeMaxUses        int  `json:"wild_shape_max_uses"`
	WildShapeUsesRemaining  int  `json:"wild_shape_uses_remaining"`
	HasDruidic              bool `json:"has_druidic"`
	HasWildShapeImprovement bool `json:"has_wild_shape_improvement"`
	HasBeastSpells          bool `json:"has_beast_spells"`
}

// MonkFeatures 武僧特性状态
type MonkFeatures struct {
	KiPointsMax           int  `json:"ki_points_max"`
	KiPointsRemaining     int  `json:"ki_points_remaining"`
	MartialArtsBonus      int  `json:"martial_arts_bonus"`
	HasUnarmoredDefense   bool `json:"has_unarmored_defense"`
	UnarmoredDefenseBonus int  `json:"unarmored_defense_bonus"`
	HasUnarmoredMovement  bool `json:"has_unarmored_movement"`
	HasStepOfTheWind      bool `json:"has_step_of_the_wind"`
	HasStunningStrike     bool `json:"has_stunning_strike"`
	HasEvasion            bool `json:"has_evasion"`
	HasUncannyDodge       bool `json:"has_uncanny_dodge"`
}

// PaladinFeatures 圣武士特性状态
type PaladinFeatures struct {
	LayOnHandsMax          int           `json:"lay_on_hands_max"`
	LayOnHandsRemaining    int           `json:"lay_on_hands_remaining"`
	DivineSenseMax         int           `json:"divine_sense_max"`
	DivineSenseRemaining   int           `json:"divine_sense_remaining"`
	DivineSmiteAvailable   bool          `json:"divine_smite_available"`
	SelectedFightingStyle  FightingStyle `json:"selected_fighting_style"`
	HasAuraOfProtection    bool          `json:"has_aura_of_protection"`
	HasImprovedDivineSmite bool          `json:"has_improved_divine_smite"`
	HasCleansingTouch      bool          `json:"has_cleansing_touch"`
}

// RangerFeatures 游侠特性状态
type RangerFeatures struct {
	HasFavoredEnemy      bool `json:"has_favored_enemy"`
	HasNaturalExplorer   bool `json:"has_natural_explorer"`
	HasFightingStyle     bool `json:"has_fighting_style"`
	HasSpellcasting      bool `json:"has_spellcasting"`
	HuntersPreyAvailable bool `json:"hunters_prey_available"`
	ExtraAttacks         int  `json:"extra_attacks"`
	HasLandStride        bool `json:"has_land_stride"`
	HasHideInPlainSight  bool `json:"has_hide_in_plain_sight"`
}

// RogueFeatures 游荡者特性状态
type RogueFeatures struct {
	SneakAttackDamage      int  `json:"sneak_attack_damage"`
	SneakAttackAvailable   bool `json:"sneak_attack_available"`
	CunningActionAvailable bool `json:"cunning_action_available"`
	HasUncannyDodge        bool `json:"has_uncanny_dodge"`
	HasEvasion             bool `json:"has_evasion"`
	HasReliableTalent      bool `json:"has_reliable_talent"`
	HasBlindsense          bool `json:"has_blindsense"`
	HasSlipperyMind        bool `json:"has_slippery_mind"`
	HasElusive             bool `json:"has_elusive"`
	HasStrokeOfLuck        bool `json:"has_stroke_of_luck"`
}

// SorcererFeatures 术士特性状态
type SorcererFeatures struct {
	SorceryPointsMax        int  `json:"sorcery_points_max"`
	SorceryPointsRemaining  int  `json:"sorcery_points_remaining"`
	HasMetamagic            bool `json:"has_metamagic"`
	MetamagicOptionsKnown   int  `json:"metamagic_options_known"`
	HasSorcerousRestoration bool `json:"has_sorcerous_restoration"`
}

// WarlockFeatures 邪术师特性状态
type WarlockFeatures struct {
	SpellSlotLevel         int  `json:"spell_slot_level"`
	SpellSlotsMax          int  `json:"spell_slots_max"`
	SpellSlotsUsed         int  `json:"spell_slots_used"`
	HasEldritchInvocations bool `json:"has_eldritch_invocations"`
	InvocationsKnown       int  `json:"invocations_known"`
	InvocationsUsed        int  `json:"invocations_used"`
	HasPactBoon            bool `json:"has_pact_boon"`
	HasMysticArcanum       bool `json:"has_mystic_arcanum"`
}

// WizardFeatures 法师特性状态
type WizardFeatures struct {
	ArcaneRecoveryMax  int  `json:"arcane_recovery_max"`
	ArcaneRecoveryUsed int  `json:"arcane_recovery_used"`
	HasArcaneTradition bool `json:"has_arcane_tradition"`
	HasSpellbook       bool `json:"has_spellbook"`
	SpellbookSpells    int  `json:"spellbook_spells"`
	HasSpellMastery    bool `json:"has_spell_mastery"`
	HasSignatureSpells bool `json:"has_signature_spells"`
}
