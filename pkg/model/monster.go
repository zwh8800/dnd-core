package model

// RechargeInfo 代表怪物动作的充能机制
type RechargeInfo struct {
	// RollRange 充能掷骰范围，例如 5-6 表示 d6 掷出 5 或 6 时充能
	RollRange [2]int `json:"roll_range"` // [最小值, 最大值]
	// UsesPerDay 每日使用次数（如果不为 0 则使用此机制而非掷骰）
	UsesPerDay int `json:"uses_per_day,omitempty"`
	// CurrentUses 当前剩余使用次数
	CurrentUses int `json:"current_uses,omitempty"`
}

// IsRecharged 检查动作是否已充能
func (r *RechargeInfo) IsRecharged() bool {
	if r.UsesPerDay > 0 {
		return r.CurrentUses > 0
	}
	// 如果 RollRange 未设置，默认已充能
	if r.RollRange[0] == 0 && r.RollRange[1] == 0 {
		return true
	}
	return true // 充能状态由外部管理
}

// MonsterSaveBonus 代表怪物的豁免熟练
type MonsterSaveBonus struct {
	Ability Ability `json:"ability"`
	Bonus   int     `json:"bonus"`
}

// MonsterSkillBonus 代表怪物的技能熟练
type MonsterSkillBonus struct {
	Skill Skill `json:"skill"`
	Bonus int   `json:"bonus"`
}

// DamageImmunity 代表伤害免疫/抗性/易伤
type DamageImmunity struct {
	DamageTypes []DamageType `json:"damage_types"`
	// NonMagical 是否仅对非魔法攻击免疫
	NonMagical bool `json:"non_magical,omitempty"`
}

// ConditionImmunity 代表状态免疫
type ConditionImmunity struct {
	Conditions []ConditionType `json:"conditions"`
}

// Senses 代表怪物的感官
type Senses struct {
	// Darkvision 黑暗视觉距离（英尺），0 表示无
	Darkvision int `json:"darkvision,omitempty"`
	// Blindsight 盲视距离（英尺），0 表示无
	Blindsight int `json:"blindsight,omitempty"`
	// Tremorsense 震颤感知距离（英尺），0 表示无
	Tremorsense int `json:"tremorsense,omitempty"`
	// Truesight 真实视觉距离（英尺），0 表示无
	Truesight int `json:"truesight,omitempty"`
	// PassivePerception 被动察觉
	PassivePerception int `json:"passive_perception"`
}

// MonsterTrait 代表怪物的被动特性（如传奇抗性、魔法抗力等）
type MonsterTrait struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	// Source 特性来源（用于分类）
	Source string `json:"source,omitempty"`
}

// SaveEffect 代表豁免失败/成功/减半的效果
type SaveEffect struct {
	// DC 豁免难度等级
	DC int `json:"dc"`
	// Ability 豁免属性
	Ability Ability `json:"ability"`
	// OnSuccess 成功时的效果描述
	OnSuccess string `json:"on_success,omitempty"`
	// OnFailure 失败时的效果描述
	OnFailure string `json:"on_failure"`
	// HalfDamage 失败时伤害减半
	HalfDamage bool `json:"half_damage,omitempty"`
}

// MonsterAttackEffect 代表怪物攻击的额外效果
type MonsterAttackEffect struct {
	// DamageDice 额外伤害骰（如 "2d6 poison"）
	DamageDice string `json:"damage_dice,omitempty"`
	// DamageType 额外伤害类型
	DamageType DamageType `json:"damage_type,omitempty"`
	// SaveEffect 豁免效果（如中毒、倒地等）
	SaveEffect *SaveEffect `json:"save_effect,omitempty"`
	// ConditionApplied 施加的状态
	ConditionApplied ConditionType `json:"condition_applied,omitempty"`
}

// MonsterAction 代表怪物的动作（攻击、法术、特殊能力）
type MonsterAction struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	// ActionType 动作类型
	ActionType ActionType `json:"action_type"`
	// AttackBonus 攻击加值（如果是攻击动作）
	AttackBonus int `json:"attack_bonus,omitempty"`
	// Reach 触及范围（英尺，近战攻击）
	Reach int `json:"reach,omitempty"`
	// Range 射程（英尺，远程攻击）
	Range int `json:"range,omitempty"`
	// LongRange 远程攻击的长射程
	LongRange int `json:"long_range,omitempty"`
	// TargetCount 目标数量
	TargetCount int `json:"target_count,omitempty"`
	// DamageDice 伤害骰表达式（如 "1d8+2"）
	DamageDice string `json:"damage_dice,omitempty"`
	// DamageType 伤害类型
	DamageType DamageType `json:"damage_type,omitempty"`
	// Effects 额外效果
	Effects []MonsterAttackEffect `json:"effects,omitempty"`
	// Recharge 充能信息
	Recharge *RechargeInfo `json:"recharge,omitempty"`
}

// MonsterActionType 怪物动作类型
type MonsterActionType string

const (
	MonsterActionMeleeAttack  MonsterActionType = "melee_attack"  // 近战攻击
	MonsterActionRangedAttack MonsterActionType = "ranged_attack" // 远程攻击
	MonsterActionSpell        MonsterActionType = "spell"         // 法术
	MonsterActionSpecial      MonsterActionType = "special"       // 特殊能力
)

// MonsterStatBlock 代表 SRD 5.2.1 中的完整怪物数据块（模板，不可变）
type MonsterStatBlock struct {
	// ID 怪物模板唯一标识
	ID string `json:"id"`
	// Name 怪物名称
	Name string `json:"name"`
	// Size 体型
	Size Size `json:"size"`
	// CreatureType 生物类型
	CreatureType CreatureType `json:"creature_type"`
	// CreatureTags 生物标签（如 "Human", "Goblinoid"）
	CreatureTags []CreatureTag `json:"creature_tags,omitempty"`
	// Alignment 阵营
	Alignment string `json:"alignment"`

	// ArmorClass 盔甲等级
	ArmorClass int `json:"armor_class"`
	// ArmorType 护甲类型描述（如 "natural armor", "leather armor"）
	ArmorType string `json:"armor_type,omitempty"`

	// InitiativeBonus 先攻调整值
	InitiativeBonus int `json:"initiative_bonus"`

	// HitDice 生命骰（如 "5d8+10"）
	HitDice string `json:"hit_dice"`
	// HitPointsAverage 平均生命值
	HitPointsAverage int `json:"hit_points_average"`

	// Speed 速度（包含多种移动方式）
	Speed SpeedTypes `json:"speed"`

	// AbilityScores 属性分数
	AbilityScores AbilityScores `json:"ability_scores"`

	// SaveBonuses 豁免熟练加值
	SaveBonuses []MonsterSaveBonus `json:"save_bonuses,omitempty"`
	// SkillBonuses 技能熟练加值
	SkillBonuses []MonsterSkillBonus `json:"skill_bonuses,omitempty"`

	// DamageVulnerabilities 伤害易伤
	DamageVulnerabilities []DamageImmunity `json:"damage_vulnerabilities,omitempty"`
	// DamageResistances 伤害抗性
	DamageResistances []DamageImmunity `json:"damage_resistances,omitempty"`
	// DamageImmunities 伤害免疫
	DamageImmunities []DamageImmunity `json:"damage_immunities,omitempty"`
	// ConditionImmunities 状态免疫
	ConditionImmunities []ConditionType `json:"condition_immunities,omitempty"`

	// Senses 感官
	Senses Senses `json:"senses"`
	// Languages 语言
	Languages string `json:"languages"`

	// ChallengeRating 挑战等级
	ChallengeRating string `json:"challenge_rating"`
	// ExperiencePoints 经验值
	ExperiencePoints int `json:"experience_points"`
	// ProficiencyBonus 熟练加值（已根据 CR 计算）
	ProficiencyBonus int `json:"proficiency_bonus"`

	// Traits 被动特性
	Traits []MonsterTrait `json:"traits,omitempty"`
	// Actions 动作
	Actions []MonsterAction `json:"actions,omitempty"`
	// BonusActions 附赠动作
	BonusActions []MonsterAction `json:"bonus_actions,omitempty"`
	// Reactions 反应
	Reactions []MonsterAction `json:"reactions,omitempty"`
	// LegendaryActions 传说动作
	LegendaryActions []MonsterAction `json:"legendary_actions,omitempty"`
	// LegendaryActionsPerRound 每轮可执行的传说动作数量
	LegendaryActionsPerRound int `json:"legendary_actions_per_round,omitempty"`

	// Description 描述文本
	Description string `json:"description,omitempty"`
}

// IsSpellcaster 检查怪物是否是施法者
func (msb *MonsterStatBlock) IsSpellcaster() bool {
	for _, action := range msb.Actions {
		if action.ActionType == ActionCastSpell {
			return true
		}
	}
	return false
}

// HasLegendaryActions 检查怪物是否有传说动作
func (msb *MonsterStatBlock) HasLegendaryActions() bool {
	return len(msb.LegendaryActions) > 0
}
