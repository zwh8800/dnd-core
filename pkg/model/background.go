package model

// BackgroundID 代表 SRD 5.2.1 中的背景标识
type BackgroundID string

const (
	BackgroundAcolyte  BackgroundID = "acolyte"  // 侍僧
	BackgroundCriminal BackgroundID = "criminal" // 罪犯
	BackgroundSage     BackgroundID = "sage"     // 学者
	BackgroundSoldier  BackgroundID = "soldier"  // 士兵
)

// BackgroundChoice 代表背景起始装备的选择项
type BackgroundChoice struct {
	// OptionA 选项 A 描述
	OptionA string `json:"option_a"`
	// OptionB 选项 B 描述
	OptionB string `json:"option_b"`
}

// BackgroundDefinition 代表 SRD 5.2.1 中的背景定义
type BackgroundDefinition struct {
	// ID 背景唯一标识
	ID BackgroundID `json:"id"`
	// Name 背景名称
	Name string `json:"name"`
	// Description 背景描述
	Description string `json:"description"`

	// AbilityScoreIncreases 属性值增加选项（SRD 5.2.1 背景可能提供属性加成）
	AbilityScoreIncreases map[Ability]int `json:"ability_score_increases,omitempty"`

	// SkillProficiencies 技能熟练
	SkillProficiencies []Skill `json:"skill_proficiencies"`
	// ToolProficiencies 工具熟练
	ToolProficiencies []string `json:"tool_proficiencies,omitempty"`
	// LanguageProficiencies 语言熟练
	LanguageProficiencies []string `json:"language_proficiencies,omitempty"`

	// AssociatedFeat 关联的起源专长
	AssociatedFeat string `json:"associated_feat,omitempty"`

	// StartingEquipmentChoices 起始装备选择
	StartingEquipmentChoices []BackgroundChoice `json:"starting_equipment_choices,omitempty"`
	// StartingEquipment 固定起始装备
	StartingEquipment []string `json:"starting_equipment,omitempty"`
	// StartingGold 起始金币（银币数量）
	StartingGold int `json:"starting_gold,omitempty"`

	// FeatureName 背景特性名称
	FeatureName string `json:"feature_name,omitempty"`
	// FeatureDescription 背景特性描述
	FeatureDescription string `json:"feature_description,omitempty"`

	// SuggestedCharacteristics 建议的角色特征
	SuggestedCharacteristics *BackgroundCharacteristics `json:"suggested_characteristics,omitempty"`
}

// BackgroundCharacteristics 代表背景建议的角色特征
type BackgroundCharacteristics struct {
	// PersonalityTraits 性格特征表
	PersonalityTraits []string `json:"personality_traits,omitempty"`
	// Ideals 理想表
	Ideals []string `json:"ideals,omitempty"`
	// Bonds 羁绊表
	Bonds []string `json:"bonds,omitempty"`
	// Flaws 缺陷表
	Flaws []string `json:"flaws,omitempty"`
}
