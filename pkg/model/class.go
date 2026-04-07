package model

// ClassID 代表 D&D 5e 官方职业类型
type ClassID string

// D&D 5e 官方 12 种职业
const (
	ClassBarbarian ClassID = "野蛮人"
	ClassBard      ClassID = "吟游诗人"
	ClassCleric    ClassID = "牧师"
	ClassDruid     ClassID = "德鲁伊"
	ClassFighter   ClassID = "战士"
	ClassMonk      ClassID = "武僧"
	ClassPaladin   ClassID = "圣武士"
	ClassRanger    ClassID = "游侠"
	ClassRogue     ClassID = "游荡者"
	ClassSorcerer  ClassID = "术士"
	ClassWarlock   ClassID = "邪术师"
	ClassWizard    ClassID = "法师"
)

// String 返回职业的字符串表示
func (id ClassID) String() string {
	return string(id)
}

// IsValid 检查是否为有效的 D&D 5e 官方职业
func (id ClassID) IsValid() bool {
	switch id {
	case ClassBarbarian, ClassBard, ClassCleric, ClassDruid,
		ClassFighter, ClassMonk, ClassPaladin, ClassRanger,
		ClassRogue, ClassSorcerer, ClassWarlock, ClassWizard:
		return true
	default:
		return false
	}
}

// AllClasses 返回所有 D&D 5e 官方职业
func AllClasses() []ClassID {
	return []ClassID{
		ClassBarbarian,
		ClassBard,
		ClassCleric,
		ClassDruid,
		ClassFighter,
		ClassMonk,
		ClassPaladin,
		ClassRanger,
		ClassRogue,
		ClassSorcerer,
		ClassWarlock,
		ClassWizard,
	}
}

// ClassLevel 代表一个职业的等级信息
type ClassLevel struct {
	Class    ClassID  `json:"class"`
	Level    int      `json:"level"`
	Features []string `json:"features,omitempty"` // 该等级获得的职业特性
}

// CasterType 代表施法者类型
type CasterType string

const (
	CasterTypeNone  CasterType = ""      // 非施法者
	CasterTypeFull  CasterType = "full"  // 全施法者（法师、牧师、德鲁伊、术士、吟游诗人）
	CasterTypeHalf  CasterType = "half"  // 半施法者（圣武士、游侠）
	CasterTypeThird CasterType = "third" // 1/3施法者（奥法骑士、诡术师）
)

// FightingStyle 战士战斗风格
type FightingStyle string

const (
	FightingStyleArchery     FightingStyle = "箭术"
	FightingStyleDefense     FightingStyle = "防御"
	FightingStyleDueling     FightingStyle = "对决"
	FightingStyleGreatWeapon FightingStyle = "巨武器战斗"
	FightingStyleProtection  FightingStyle = "守护"
	FightingStyleTwoWeapon   FightingStyle = "双武器战斗"
)

// MartialArchetype 战士武术范型
type MartialArchetype string

const (
	MartialArchetypeChampion       MartialArchetype = "勇士"
	MartialArchetypeBattleMaster   MartialArchetype = "战斗大师"
	MartialArchetypeEldritchKnight MartialArchetype = "奥法骑士"
)

// FighterFeatures 战士职业特性跟踪
type FighterFeatures struct {
	// 战斗风格和武术范型
	SelectedFightingStyle FightingStyle    `json:"selected_fighting_style,omitempty"`
	SelectedArchetype     MartialArchetype `json:"selected_archetype,omitempty"`

	// 回气（Second Wind）
	SecondWindMax  int `json:"second_wind_max"`  // 最大使用次数（通常为1）
	SecondWindUsed int `json:"second_wind_used"` // 已使用次数

	// 动作如潮（Action Surge）
	ActionSurgeMax  int `json:"action_surge_max"`  // 最大使用次数（2级1次，17级2次）
	ActionSurgeUsed int `json:"action_surge_used"` // 已使用次数

	// 额外攻击（Extra Attack）
	ExtraAttacks int `json:"extra_attacks"` // 额外攻击次数（5级+1，11级+2，20级+3）

	// 不屈（Indomitable）
	IndomitableMax  int `json:"indomitable_max"`  // 最大使用次数（9级1次，13级2次，17级3次）
	IndomitableUsed int `json:"indomitable_used"` // 已使用次数
}

// ClassState 职业特性状态接口，用于扩展其他职业的特性跟踪
type ClassState interface {
	ClassID() ClassID
}

// ClassID 实现 ClassState 接口
func (id ClassID) ClassID() ClassID {
	return id
}
