package model

// SpellTargetType 代表法术目标类型
type SpellTargetType string

const (
	SpellTargetSingleTarget SpellTargetType = "single_target" // 单一目标
	SpellTargetCone         SpellTargetType = "cone"          // 锥形区域
	SpellTargetSphere       SpellTargetType = "sphere"        // 球形区域
	SpellTargetLine         SpellTargetType = "line"          // 线形区域
	SpellTargetCube         SpellTargetType = "cube"          // 立方体区域
	SpellTargetEmanation    SpellTargetType = "emanation"     // 扩散区域
	SpellTargetSelf         SpellTargetType = "self"          // 自身
	SpellTargetTouch        SpellTargetType = "touch"         // 接触
)

// SpellDamageEntry 代表法术伤害条目
type SpellDamageEntry struct {
	// BaseDice 基础伤害骰（如 "8d6"）
	BaseDice string `json:"base_dice"`
	// DamageType 伤害类型
	DamageType DamageType `json:"damage_type"`
	// UpcastDicePerLevel 每升一级的伤害骰增量
	UpcastDicePerLevel string `json:"upcast_dice_per_level,omitempty"`
	// UpcastStartLevel 开始升环增加伤害的等级
	UpcastStartLevel int `json:"upcast_start_level,omitempty"`
}

// SpellEffect 代表法术效果
type SpellEffect struct {
	// Type 效果类型
	Type SpellEffectType `json:"type"`
	// TargetType 目标类型
	TargetType SpellTargetType `json:"target_type"`
	// Range 射程
	Range string `json:"range"`
	// AreaSize 区域大小（尺）
	AreaSize int `json:"area_size,omitempty"`
	// Damage 伤害条目
	Damage *SpellDamageEntry `json:"damage,omitempty"`
	// HealingDice 治疗骰（如 "1d8+3"）
	HealingDice string `json:"healing_dice,omitempty"`
	// SaveDC 豁免DC（0表示不需要豁免）
	SaveDC int `json:"save_dc,omitempty"`
	// SaveAbility 豁免属性
	SaveAbility Ability `json:"save_ability,omitempty"`
	// SaveSuccessEffect 豁免成功效果（"half"=半伤, "none"=无效, "partial"=部分效果）
	SaveSuccessEffect string `json:"save_success_effect,omitempty"`
	// ConditionApplied 施加的状态
	ConditionApplied ConditionType `json:"condition_applied,omitempty"`
	// ConditionDuration 状态持续时间
	ConditionDuration string `json:"condition_duration,omitempty"`
	// SummonCreature 召唤的生物
	SummonCreature string `json:"summon_creature,omitempty"`
	// SummonCount 召唤数量
	SummonCount int `json:"summon_count,omitempty"`
	// SummonDuration 召唤持续时间
	SummonDuration string `json:"summon_duration,omitempty"`
	// Description 效果描述
	Description string `json:"description"`
}

// SpellEffectType 代表法术效果类型
type SpellEffectType string

const (
	SpellEffectDamage    SpellEffectType = "damage"        // 伤害
	SpellEffectHealing   SpellEffectType = "healing"       // 治疗
	SpellEffectCondition SpellEffectType = "condition"     // 施加状态
	SpellEffectBuff      SpellEffectType = "buff"          // 增益效果
	SpellEffectDebuff    SpellEffectType = "debuff"        // 减益效果
	SpellEffectSummon    SpellEffectType = "summon"        // 召唤
	SpellEffectTeleport  SpellEffectType = "teleport"      // 传送
	SpellEffectUtility   SpellEffectType = "utility"       // 实用效果
	SpellEffectCreateObj SpellEffectType = "create_object" // 创造物体
)

// SpellDefinition 完整的法术定义（扩展 Spell）
type SpellDefinition struct {
	Spell
	// Effects 法术效果列表
	Effects []SpellEffect `json:"effects,omitempty"`
	// RequiresAttackRoll 是否需要攻击掷骰
	RequiresAttackRoll bool `json:"requires_attack_roll,omitempty"`
	// UpcastEffects 升环效果
	UpcastEffects map[int][]SpellEffect `json:"upcast_effects,omitempty"`
}
