package data

import (
	"fmt"

	"github.com/zwh8800/dnd-core/pkg/model"
)

// init 注册所有职业
func init() {
	for id, class := range Classes {
		class.ID = id
		GlobalRegistry.RegisterClass(class)
	}
}

// ClassDefinition 职业定义
type ClassDefinition struct {
	ID                  model.ClassID     // 职业ID
	Name                string            // 中文名
	HitDie              int               // 生命骰: 6/8/10/12
	PrimaryAbilities    []model.Ability   // 主要属性(按重要性排序)
	SavingThrows        []model.Ability   // 豁免熟练(2项)
	SkillChoices        []model.Skill     // 可选技能列表
	NumberOfSkills      int               // 1级可选技能数量
	ArmorProficiencies  []model.ArmorType // 护甲熟练
	WeaponProficiencies []string          // 武器熟练
	ToolProficiencies   []string          // 工具熟练
	SpellcastingAbility model.Ability     // 施法属性(非施法职业为空)
	CasterType          model.CasterType  // 施法者类型
	Description         string            // 中文描述
}

// Classes 所有职业定义
var Classes = map[model.ClassID]*ClassDefinition{
	model.ClassBarbarian: {
		ID:     model.ClassBarbarian,
		Name:   "野蛮人",
		HitDie: 12,
		PrimaryAbilities: []model.Ability{
			model.AbilityStrength,
			model.AbilityConstitution,
			model.AbilityDexterity,
		},
		SavingThrows: []model.Ability{
			model.AbilityStrength,
			model.AbilityConstitution,
		},
		SkillChoices: []model.Skill{
			model.SkillAnimalHandling,
			model.SkillAthletics,
			model.SkillIntimidation,
			model.SkillNature,
			model.SkillPerception,
			model.SkillSurvival,
		},
		NumberOfSkills:      2,
		ArmorProficiencies:  []model.ArmorType{model.ArmorTypeLight, model.ArmorTypeMedium, model.ArmorTypeShield},
		WeaponProficiencies: []string{"简易武器", "军用武器"},
		ToolProficiencies:   []string{},
		SpellcastingAbility: "",
		CasterType:          model.CasterTypeNone,
		Description:         "野蛮人凭借原始的愤怒和直觉在战斗中狂暴冲锋,拥有强大的生命力和破坏力",
	},

	model.ClassBard: {
		ID:     model.ClassBard,
		Name:   "吟游诗人",
		HitDie: 8,
		PrimaryAbilities: []model.Ability{
			model.AbilityCharisma,
			model.AbilityDexterity,
			model.AbilityConstitution,
		},
		SavingThrows: []model.Ability{
			model.AbilityDexterity,
			model.AbilityCharisma,
		},
		SkillChoices: []model.Skill{
			model.SkillAcrobatics,
			model.SkillAnimalHandling,
			model.SkillArcana,
			model.SkillAthletics,
			model.SkillDeception,
			model.SkillHistory,
			model.SkillInsight,
			model.SkillIntimidation,
			model.SkillInvestigation,
			model.SkillMedicine,
			model.SkillNature,
			model.SkillPerception,
			model.SkillPerformance,
			model.SkillPersuasion,
			model.SkillReligion,
			model.SkillSleightOfHand,
			model.SkillStealth,
			model.SkillSurvival,
		},
		NumberOfSkills:      3,
		ArmorProficiencies:  []model.ArmorType{model.ArmorTypeLight},
		WeaponProficiencies: []string{"简易武器", "手弩", "长剑", "细剑", "短剑", "短弓", "长弓"},
		ToolProficiencies:   []string{"三种乐器"},
		SpellcastingAbility: model.AbilityCharisma,
		CasterType:          model.CasterTypeFull,
		Description:         "吟游诗人用音乐和魔法激励盟友、迷惑敌人,是全能的支持者和施法者",
	},

	model.ClassCleric: {
		ID:     model.ClassCleric,
		Name:   "牧师",
		HitDie: 8,
		PrimaryAbilities: []model.Ability{
			model.AbilityWisdom,
			model.AbilityStrength,
			model.AbilityConstitution,
		},
		SavingThrows: []model.Ability{
			model.AbilityWisdom,
			model.AbilityCharisma,
		},
		SkillChoices: []model.Skill{
			model.SkillHistory,
			model.SkillInsight,
			model.SkillMedicine,
			model.SkillPersuasion,
			model.SkillReligion,
		},
		NumberOfSkills:      2,
		ArmorProficiencies:  []model.ArmorType{model.ArmorTypeLight, model.ArmorTypeMedium, model.ArmorTypeShield},
		WeaponProficiencies: []string{"简易武器"},
		ToolProficiencies:   []string{},
		SpellcastingAbility: model.AbilityWisdom,
		CasterType:          model.CasterTypeFull,
		Description:         "牧师是神祗的凡人代表,拥有强大的治疗和战斗法术",
	},

	model.ClassDruid: {
		ID:     model.ClassDruid,
		Name:   "德鲁伊",
		HitDie: 8,
		PrimaryAbilities: []model.Ability{
			model.AbilityWisdom,
			model.AbilityIntelligence,
			model.AbilityConstitution,
		},
		SavingThrows: []model.Ability{
			model.AbilityIntelligence,
			model.AbilityWisdom,
		},
		SkillChoices: []model.Skill{
			model.SkillArcana,
			model.SkillAnimalHandling,
			model.SkillInsight,
			model.SkillMedicine,
			model.SkillNature,
			model.SkillPerception,
			model.SkillReligion,
			model.SkillSurvival,
		},
		NumberOfSkills:      2,
		ArmorProficiencies:  []model.ArmorType{model.ArmorTypeLight, model.ArmorTypeMedium, model.ArmorTypeShield},
		WeaponProficiencies: []string{"木棍", "匕首", "飞镖", "木棒", "弯刀", "手斧", "轻锤", "长矛"},
		ToolProficiencies:   []string{"草药师工具"},
		SpellcastingAbility: model.AbilityWisdom,
		CasterType:          model.CasterTypeFull,
		Description:         "德鲁伊是自然的守护者,能化身为野兽并操控自然之力",
	},

	model.ClassFighter: {
		ID:     model.ClassFighter,
		Name:   "战士",
		HitDie: 10,
		PrimaryAbilities: []model.Ability{
			model.AbilityStrength,
			model.AbilityDexterity,
			model.AbilityConstitution,
		},
		SavingThrows: []model.Ability{
			model.AbilityStrength,
			model.AbilityConstitution,
		},
		SkillChoices: []model.Skill{
			model.SkillAcrobatics,
			model.SkillAnimalHandling,
			model.SkillAthletics,
			model.SkillHistory,
			model.SkillInsight,
			model.SkillIntimidation,
			model.SkillPerception,
			model.SkillSurvival,
		},
		NumberOfSkills:      2,
		ArmorProficiencies:  []model.ArmorType{model.ArmorTypeLight, model.ArmorTypeMedium, model.ArmorTypeHeavy, model.ArmorTypeShield},
		WeaponProficiencies: []string{"简易武器", "军用武器"},
		ToolProficiencies:   []string{},
		SpellcastingAbility: "",
		CasterType:          model.CasterTypeNone,
		Description:         "战士是精通各种战斗方式的大师,从剑术到箭术,从重装到双持",
	},

	model.ClassMonk: {
		ID:     model.ClassMonk,
		Name:   "武僧",
		HitDie: 8,
		PrimaryAbilities: []model.Ability{
			model.AbilityDexterity,
			model.AbilityWisdom,
			model.AbilityStrength,
		},
		SavingThrows: []model.Ability{
			model.AbilityStrength,
			model.AbilityDexterity,
		},
		SkillChoices: []model.Skill{
			model.SkillAcrobatics,
			model.SkillAthletics,
			model.SkillHistory,
			model.SkillInsight,
			model.SkillReligion,
			model.SkillStealth,
		},
		NumberOfSkills:      2,
		ArmorProficiencies:  []model.ArmorType{},
		WeaponProficiencies: []string{"简易武器", "短剑"},
		ToolProficiencies:   []string{"一种乐器或一种工匠工具"},
		SpellcastingAbility: "",
		CasterType:          model.CasterTypeNone,
		Description:         "武僧通过修行掌握体内的气,能以惊人的速度和力量战斗",
	},

	model.ClassPaladin: {
		ID:     model.ClassPaladin,
		Name:   "圣武士",
		HitDie: 10,
		PrimaryAbilities: []model.Ability{
			model.AbilityStrength,
			model.AbilityCharisma,
			model.AbilityConstitution,
		},
		SavingThrows: []model.Ability{
			model.AbilityWisdom,
			model.AbilityCharisma,
		},
		SkillChoices: []model.Skill{
			model.SkillAthletics,
			model.SkillInsight,
			model.SkillIntimidation,
			model.SkillMedicine,
			model.SkillPersuasion,
			model.SkillReligion,
		},
		NumberOfSkills:      2,
		ArmorProficiencies:  []model.ArmorType{model.ArmorTypeLight, model.ArmorTypeMedium, model.ArmorTypeHeavy, model.ArmorTypeShield},
		WeaponProficiencies: []string{"简易武器", "军用武器"},
		ToolProficiencies:   []string{},
		SpellcastingAbility: model.AbilityCharisma,
		CasterType:          model.CasterTypeHalf,
		Description:         "圣武士是神圣誓言的战士,拥有强大的近战能力和神术",
	},

	model.ClassRanger: {
		ID:     model.ClassRanger,
		Name:   "游侠",
		HitDie: 10,
		PrimaryAbilities: []model.Ability{
			model.AbilityDexterity,
			model.AbilityWisdom,
			model.AbilityStrength,
		},
		SavingThrows: []model.Ability{
			model.AbilityStrength,
			model.AbilityDexterity,
		},
		SkillChoices: []model.Skill{
			model.SkillAnimalHandling,
			model.SkillAthletics,
			model.SkillInsight,
			model.SkillInvestigation,
			model.SkillNature,
			model.SkillPerception,
			model.SkillStealth,
			model.SkillSurvival,
		},
		NumberOfSkills:      3,
		ArmorProficiencies:  []model.ArmorType{model.ArmorTypeLight, model.ArmorTypeMedium, model.ArmorTypeShield},
		WeaponProficiencies: []string{"简易武器", "军用武器"},
		ToolProficiencies:   []string{},
		SpellcastingAbility: model.AbilityWisdom,
		CasterType:          model.CasterTypeHalf,
		Description:         "游侠是荒野的守护者,精通追踪、潜行和自然法术",
	},

	model.ClassRogue: {
		ID:     model.ClassRogue,
		Name:   "游荡者",
		HitDie: 8,
		PrimaryAbilities: []model.Ability{
			model.AbilityDexterity,
			model.AbilityIntelligence,
			model.AbilityCharisma,
		},
		SavingThrows: []model.Ability{
			model.AbilityDexterity,
			model.AbilityIntelligence,
		},
		SkillChoices: []model.Skill{
			model.SkillAcrobatics,
			model.SkillAthletics,
			model.SkillDeception,
			model.SkillInsight,
			model.SkillIntimidation,
			model.SkillInvestigation,
			model.SkillPerception,
			model.SkillPerformance,
			model.SkillPersuasion,
			model.SkillSleightOfHand,
			model.SkillStealth,
		},
		NumberOfSkills:      4,
		ArmorProficiencies:  []model.ArmorType{model.ArmorTypeLight},
		WeaponProficiencies: []string{"简易武器", "手弩", "长剑", "刺剑", "短剑", "短弓", "长弓"},
		ToolProficiencies:   []string{"盗贼工具"},
		SpellcastingAbility: "",
		CasterType:          model.CasterTypeNone,
		Description:         "游荡者依靠诡计和精准打击,擅长偷袭、开锁和解除陷阱",
	},

	model.ClassSorcerer: {
		ID:     model.ClassSorcerer,
		Name:   "术士",
		HitDie: 6,
		PrimaryAbilities: []model.Ability{
			model.AbilityCharisma,
			model.AbilityConstitution,
			model.AbilityDexterity,
		},
		SavingThrows: []model.Ability{
			model.AbilityConstitution,
			model.AbilityCharisma,
		},
		SkillChoices: []model.Skill{
			model.SkillArcana,
			model.SkillDeception,
			model.SkillInsight,
			model.SkillIntimidation,
			model.SkillPersuasion,
			model.SkillReligion,
		},
		NumberOfSkills:      2,
		ArmorProficiencies:  []model.ArmorType{},
		WeaponProficiencies: []string{"匕首", "飞镖", "木棒", "轻弩", "长弓"},
		ToolProficiencies:   []string{},
		SpellcastingAbility: model.AbilityCharisma,
		CasterType:          model.CasterTypeFull,
		Description:         "术士天生具有魔法天赋,血脉中流淌着原始的魔法力量",
	},

	model.ClassWarlock: {
		ID:     model.ClassWarlock,
		Name:   "邪术师",
		HitDie: 8,
		PrimaryAbilities: []model.Ability{
			model.AbilityCharisma,
			model.AbilityWisdom,
			model.AbilityConstitution,
		},
		SavingThrows: []model.Ability{
			model.AbilityWisdom,
			model.AbilityCharisma,
		},
		SkillChoices: []model.Skill{
			model.SkillArcana,
			model.SkillDeception,
			model.SkillHistory,
			model.SkillIntimidation,
			model.SkillInvestigation,
			model.SkillNature,
			model.SkillReligion,
		},
		NumberOfSkills:      2,
		ArmorProficiencies:  []model.ArmorType{model.ArmorTypeLight},
		WeaponProficiencies: []string{"简易武器"},
		ToolProficiencies:   []string{},
		SpellcastingAbility: model.AbilityCharisma,
		CasterType:          model.CasterTypeFull,
		Description:         "邪术师与强大的存在缔结契约,借此获得独特的魔法能力",
	},

	model.ClassWizard: {
		ID:     model.ClassWizard,
		Name:   "法师",
		HitDie: 6,
		PrimaryAbilities: []model.Ability{
			model.AbilityIntelligence,
			model.AbilityConstitution,
			model.AbilityDexterity,
		},
		SavingThrows: []model.Ability{
			model.AbilityIntelligence,
			model.AbilityWisdom,
		},
		SkillChoices: []model.Skill{
			model.SkillArcana,
			model.SkillHistory,
			model.SkillInsight,
			model.SkillInvestigation,
			model.SkillMedicine,
			model.SkillReligion,
		},
		NumberOfSkills:      2,
		ArmorProficiencies:  []model.ArmorType{},
		WeaponProficiencies: []string{"匕首", "飞镖", "木棒", "轻弩", "长弓"},
		ToolProficiencies:   []string{},
		SpellcastingAbility: model.AbilityIntelligence,
		CasterType:          model.CasterTypeFull,
		Description:         "法师通过钻研和学习掌握深奥的魔法知识,是最全能的施法者",
	},
}

// GetClass 获取职业定义
func GetClass(id model.ClassID) *ClassDefinition {
	return Classes[id]
}

// GetClassID 根据名称获取职业ID,支持中英文
func GetClassID(name string) (model.ClassID, error) {
	// 先尝试直接匹配
	for id := range Classes {
		if string(id) == name {
			return id, nil
		}
	}

	// 尝试英文名映射
	englishToClassID := map[string]model.ClassID{
		"Barbarian": model.ClassBarbarian,
		"Bard":      model.ClassBard,
		"Cleric":    model.ClassCleric,
		"Druid":     model.ClassDruid,
		"Fighter":   model.ClassFighter,
		"Monk":      model.ClassMonk,
		"Paladin":   model.ClassPaladin,
		"Ranger":    model.ClassRanger,
		"Rogue":     model.ClassRogue,
		"Sorcerer":  model.ClassSorcerer,
		"Warlock":   model.ClassWarlock,
		"Wizard":    model.ClassWizard,
	}

	if id, ok := englishToClassID[name]; ok {
		return id, nil
	}

	return "", fmt.Errorf("未知的职业: %s", name)
}

// GetClassNames 获取所有职业名称
func GetClassNames() []string {
	names := make([]string, 0, len(Classes))
	for _, def := range Classes {
		names = append(names, def.Name)
	}
	return names
}

// FighterFeaturesByLevel 战士每级获得的特性
var FighterFeaturesByLevel = map[int][]string{
	1:  {"战斗风格", "复苏之风"},
	2:  {"动作如潮"},
	3:  {"武术范型"},
	4:  {"属性值提升"},
	5:  {"额外攻击"},
	6:  {"属性值提升"},
	7:  {"范型特性"},
	8:  {"属性值提升"},
	9:  {"不屈"},
	10: {"范型特性"},
	11: {"额外攻击(2)"},
	12: {"属性值提升"},
	13: {"不屈(2)"},
	14: {"属性值提升"},
	15: {"范型特性"},
	16: {"属性值提升"},
	17: {"动作如潮(2)", "不屈(3)"},
	18: {"范型特性"},
	19: {"属性值提升"},
	20: {"额外攻击(3)"},
}

// FightingStyleDescriptions 战斗风格描述
var FightingStyleDescriptions = map[model.FightingStyle]string{
	model.FightingStyleArchery:     "使用远程武器进行攻击检定时获得+2加值",
	model.FightingStyleDefense:     "着装护甲时,AC获得+1加值",
	model.FightingStyleDueling:     "单手持握一把近战武器且未持用其他武器时,该武器伤害掷骰获得+2加值",
	model.FightingStyleGreatWeapon: "双手持握近战武器攻击时,伤害掷骰掷出1或2可重掷",
	model.FightingStyleProtection:  "持盾时,可用反应使5尺内盟友受到的攻击检定具有劣势",
	model.FightingStyleTwoWeapon:   "两武器战斗时,可将属性调整值加到副手攻击的伤害上",
}

// MartialArchetypeDescriptions 武术范型描述
var MartialArchetypeDescriptions = map[model.MartialArchetype]string{
	model.MartialArchetypeChampion:       "勇士:专注于提升物理战斗能力,扩大暴击范围",
	model.MartialArchetypeBattleMaster:   "战斗大师:掌握各种战斗机动动作,使用 superiority dice",
	model.MartialArchetypeEldritchKnight: "奥法骑士:将魔法与武艺结合,学习法师法术",
}
