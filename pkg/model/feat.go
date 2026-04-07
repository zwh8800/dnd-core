package model

// FeatType 代表专长的类型
type FeatType string

const (
	FeatTypeOrigin  FeatType = "Origin"  // 起源专长（角色创建时获得）
	FeatTypeGeneral FeatType = "General" // 通用专长
	FeatTypeCombat  FeatType = "Combat"  // 战斗专长
	FeatTypeEpic    FeatType = "Epic"    // 史诗恩赐
)

// FeatPrerequisite 代表专长的先决条件
type FeatPrerequisite struct {
	// MinimumAbilityScores 最低属性分数要求
	MinimumAbilityScores map[Ability]int `json:"minimum_ability_scores,omitempty"`
	// RequiredClass 需要的职业
	RequiredClass ClassID `json:"required_class,omitempty"`
	// MinimumLevel 最低等级
	MinimumLevel int `json:"minimum_level,omitempty"`
	// RequiredFeat 需要的专长（前置专长）
	RequiredFeat string `json:"required_feat,omitempty"`
	// Description 先决条件描述（用于复杂条件）
	Description string `json:"description,omitempty"`
}

// FeatEffect 代表专长应用的效果
type FeatEffect struct {
	// AbilityScoreIncrease 属性值增加
	AbilityScoreIncrease map[Ability]int `json:"ability_score_increase,omitempty"`
	// AbilityScoreMax 属性最大值提升
	AbilityScoreMax map[Ability]int `json:"ability_score_max,omitempty"`

	// AttackBonus 攻击加值
	AttackBonus int `json:"attack_bonus,omitempty"`
	// ACBonus AC 加值
	ACBonus int `json:"ac_bonus,omitempty"`
	// InitiativeBonus 先攻加值
	InitiativeBonus int `json:"initiative_bonus,omitempty"`

	// SkillProficiencies 获得的技能熟练
	SkillProficiencies []Skill `json:"skill_proficiencies,omitempty"`
	// ToolProficiencies 获得的工具熟练
	ToolProficiencies []string `json:"tool_proficiencies,omitempty"`

	// DamageResistances 获得的伤害抗性
	DamageResistances []DamageType `json:"damage_resistances,omitempty"`

	// SpecialAbilities 特殊能力标记（由 FeatureHook 处理）
	SpecialAbilities []string `json:"special_abilities,omitempty"`

	// Description 效果描述
	Description string `json:"description,omitempty"`
}

// FeatDefinition 代表 SRD 5.2.1 中的专长定义
type FeatDefinition struct {
	// ID 专长唯一标识
	ID string `json:"id"`
	// Name 专长名称
	Name string `json:"name"`
	// Type 专长类型
	Type FeatType `json:"type"`
	// Prerequisite 先决条件
	Prerequisite *FeatPrerequisite `json:"prerequisite,omitempty"`
	// Effects 专长效果
	Effects FeatEffect `json:"effects"`
	// Repeatable 是否可重复选择
	Repeatable bool `json:"repeatable"`
	// Description 专长描述
	Description string `json:"description"`
}

// FeatSource 代表专长的获得来源
type FeatSource string

const (
	FeatSourceBackground FeatSource = "background" // 来自背景
	FeatSourceLevelUp    FeatSource = "level_up"   // 来自升级（ASI）
	FeatSourceVariant    FeatSource = "variant"    // 来自变体人类
)

// FeatInstance 代表角色已获得的一个专长实例
type FeatInstance struct {
	// FeatID 专长 ID
	FeatID string `json:"feat_id"`
	// Source 获得来源
	Source FeatSource `json:"source"`
	// AcquiredLevel 获得时的等级
	AcquiredLevel int `json:"acquired_level"`
}
