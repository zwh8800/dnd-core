package model

// Ability 代表D&D 5e的六大属性
type Ability string

const (
	AbilityStrength     Ability = "STR" // 力量
	AbilityDexterity    Ability = "DEX" // 敏捷
	AbilityConstitution Ability = "CON" // 体质
	AbilityIntelligence Ability = "INT" // 智力
	AbilityWisdom       Ability = "WIS" // 感知
	AbilityCharisma     Ability = "CHA" // 魅力
)

// AllAbilities 返回所有属性
func AllAbilities() []Ability {
	return []Ability{
		AbilityStrength,
		AbilityDexterity,
		AbilityConstitution,
		AbilityIntelligence,
		AbilityWisdom,
		AbilityCharisma,
	}
}

// AbilityScores 存储六个属性值
type AbilityScores struct {
	Strength     int `json:"strength"`
	Dexterity    int `json:"dexterity"`
	Constitution int `json:"constitution"`
	Intelligence int `json:"intelligence"`
	Wisdom       int `json:"wisdom"`
	Charisma     int `json:"charisma"`
}

// Get 获取指定属性的值
func (s *AbilityScores) Get(ability Ability) int {
	switch ability {
	case AbilityStrength:
		return s.Strength
	case AbilityDexterity:
		return s.Dexterity
	case AbilityConstitution:
		return s.Constitution
	case AbilityIntelligence:
		return s.Intelligence
	case AbilityWisdom:
		return s.Wisdom
	case AbilityCharisma:
		return s.Charisma
	default:
		return 0
	}
}

// Set 设置指定属性的值
func (s *AbilityScores) Set(ability Ability, value int) {
	switch ability {
	case AbilityStrength:
		s.Strength = value
	case AbilityDexterity:
		s.Dexterity = value
	case AbilityConstitution:
		s.Constitution = value
	case AbilityIntelligence:
		s.Intelligence = value
	case AbilityWisdom:
		s.Wisdom = value
	case AbilityCharisma:
		s.Charisma = value
	}
}

// Skill 代表D&D 5e的技能
type Skill string

const (
	SkillAcrobatics     Skill = "Acrobatics"      // 特技（DEX）
	SkillAnimalHandling Skill = "Animal Handling" // 驯兽（WIS）
	SkillArcana         Skill = "Arcana"          // 奥秘（INT）
	SkillAthletics      Skill = "Athletics"       // 运动（STR）
	SkillDeception      Skill = "Deception"       // 欺瞒（CHA）
	SkillHistory        Skill = "History"         // 历史（INT）
	SkillInsight        Skill = "Insight"         // 洞察（WIS）
	SkillIntimidation   Skill = "Intimidation"    // 威吓（CHA）
	SkillInvestigation  Skill = "Investigation"   // 调查（INT）
	SkillMedicine       Skill = "Medicine"        // 医药（WIS）
	SkillNature         Skill = "Nature"          // 自然（INT）
	SkillPerception     Skill = "Perception"      // 察觉（WIS）
	SkillPerformance    Skill = "Performance"     // 表演（CHA）
	SkillPersuasion     Skill = "Persuasion"      // 游说（CHA）
	SkillReligion       Skill = "Religion"        // 宗教（INT）
	SkillSleightOfHand  Skill = "Sleight of Hand" // 巧手（DEX）
	SkillStealth        Skill = "Stealth"         // 隐匿（DEX）
	SkillSurvival       Skill = "Survival"        // 求生（WIS）
)

// SkillAbilityMap 定义技能与属性的对应关系
var SkillAbilityMap = map[Skill]Ability{
	SkillAcrobatics:     AbilityDexterity,
	SkillAnimalHandling: AbilityWisdom,
	SkillArcana:         AbilityIntelligence,
	SkillAthletics:      AbilityStrength,
	SkillDeception:      AbilityCharisma,
	SkillHistory:        AbilityIntelligence,
	SkillInsight:        AbilityWisdom,
	SkillIntimidation:   AbilityCharisma,
	SkillInvestigation:  AbilityIntelligence,
	SkillMedicine:       AbilityWisdom,
	SkillNature:         AbilityIntelligence,
	SkillPerception:     AbilityWisdom,
	SkillPerformance:    AbilityCharisma,
	SkillPersuasion:     AbilityCharisma,
	SkillReligion:       AbilityIntelligence,
	SkillSleightOfHand:  AbilityDexterity,
	SkillStealth:        AbilityDexterity,
	SkillSurvival:       AbilityWisdom,
}

// AllSkills 返回所有技能
func AllSkills() []Skill {
	return []Skill{
		SkillAcrobatics,
		SkillAnimalHandling,
		SkillArcana,
		SkillAthletics,
		SkillDeception,
		SkillHistory,
		SkillInsight,
		SkillIntimidation,
		SkillInvestigation,
		SkillMedicine,
		SkillNature,
		SkillPerception,
		SkillPerformance,
		SkillPersuasion,
		SkillReligion,
		SkillSleightOfHand,
		SkillStealth,
		SkillSurvival,
	}
}

// Proficiencies 存储熟练度信息
type Proficiencies struct {
	// ProficientSkills 熟练的技能集合
	ProficientSkills map[Skill]bool `json:"proficient_skills"`
	// ExpertiseSkills 拥有专家技能加成的技能集合（两倍熟练加值）
	ExpertiseSkills map[Skill]bool `json:"expertise_skills"`
	// SavingThrowProficiencies 熟练的豁免检定属性
	SavingThrowProficiencies map[Ability]bool `json:"saving_throw_proficiencies"`
	// ArmorProficiencies 熟练的护甲类型
	ArmorProficiencies map[ArmorType]bool `json:"armor_proficiencies"`
	// WeaponProficiencies 熟练的武器类型
	WeaponProficiencies map[string]bool `json:"weapon_proficiencies"`
	// ToolProficiencies 熟练的工具
	ToolProficiencies map[string]bool `json:"tool_proficiencies"`
	// LanguageProficiencies 掌握的语言
	LanguageProficiencies map[string]bool `json:"language_proficiencies"`
}

// NewProficiencies 创建新的熟练度配置
func NewProficiencies() *Proficiencies {
	return &Proficiencies{
		ProficientSkills:         make(map[Skill]bool),
		ExpertiseSkills:          make(map[Skill]bool),
		SavingThrowProficiencies: make(map[Ability]bool),
		ArmorProficiencies:       make(map[ArmorType]bool),
		WeaponProficiencies:      make(map[string]bool),
		ToolProficiencies:        make(map[string]bool),
		LanguageProficiencies:    make(map[string]bool),
	}
}

// IsProficient 检查是否对某项技能熟练
func (p *Proficiencies) IsProficient(skill Skill) bool {
	return p.ProficientSkills[skill]
}

// HasExpertise 检查是否对某项技能有专家加成
func (p *Proficiencies) HasExpertise(skill Skill) bool {
	return p.ExpertiseSkills[skill]
}

// IsSavingThrowProficient 检查是否对某属性豁免熟练
func (p *Proficiencies) IsSavingThrowProficient(ability Ability) bool {
	return p.SavingThrowProficiencies[ability]
}
