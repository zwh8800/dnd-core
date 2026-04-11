package data

import "github.com/zwh8800/dnd-core/pkg/model"

// RaceDefinition 种族定义
type RaceDefinition struct {
	Name           string
	Subraces       []string
	AbilityBonuses map[model.Ability]int
	Speed          int
	Size           model.Size
	Languages      []string
	Traits         []string
	Description    string
}

// Races 所有种族定义
var Races = map[string]*RaceDefinition{
	"人类": {
		Name: "人类",
		AbilityBonuses: map[model.Ability]int{
			model.AbilityStrength:     1,
			model.AbilityDexterity:    1,
			model.AbilityConstitution: 1,
			model.AbilityIntelligence: 1,
			model.AbilityWisdom:       1,
			model.AbilityCharisma:     1,
		},
		Speed:       30,
		Size:        model.SizeMedium,
		Languages:   []string{"通用语"},
		Traits:      []string{"额外语言"},
		Description: "人类是最常见、最多样化的种族",
	},
	"精灵": {
		Name:     "精灵",
		Subraces: []string{"高等精灵", "木精灵", "卓尔"},
		AbilityBonuses: map[model.Ability]int{
			model.AbilityDexterity: 2,
		},
		Speed:       30,
		Size:        model.SizeMedium,
		Languages:   []string{"通用语", "精灵语"},
		Traits:      []string{"黑暗视觉", "敏锐感官", "精类血统", "冥想"},
		Description: "优雅长寿的森林种族",
	},
	"矮人": {
		Name:     "矮人",
		Subraces: []string{"丘陵矮人", "山地矮人"},
		AbilityBonuses: map[model.Ability]int{
			model.AbilityConstitution: 2,
		},
		Speed:       25,
		Size:        model.SizeMedium,
		Languages:   []string{"通用语", "矮人语"},
		Traits:      []string{"黑暗视觉", "矮人韧性", "战斗训练", "工具熟练", "石中精妙"},
		Description: "坚韧不拔的山地民族",
	},
	"半身人": {
		Name:     "半身人",
		Subraces: []string{"轻足", "强魄"},
		AbilityBonuses: map[model.Ability]int{
			model.AbilityDexterity: 2,
		},
		Speed:       25,
		Size:        model.SizeSmall,
		Languages:   []string{"通用语", "半身人语"},
		Traits:      []string{"幸运", "勇敢", "半身人敏捷"},
		Description: "小巧而勇敢的民族",
	},
	"龙裔": {
		Name: "龙裔",
		AbilityBonuses: map[model.Ability]int{
			model.AbilityStrength: 2,
			model.AbilityCharisma: 1,
		},
		Speed:       30,
		Size:        model.SizeMedium,
		Languages:   []string{"通用语", "龙语"},
		Traits:      []string{"龙族先祖", "吐息武器", "伤害抗性"},
		Description: "拥有龙族血统的战士",
	},
	"侏儒": {
		Name:     "侏儒",
		Subraces: []string{"森林侏儒", "岩侏儒"},
		AbilityBonuses: map[model.Ability]int{
			model.AbilityIntelligence: 2,
		},
		Speed:       25,
		Size:        model.SizeSmall,
		Languages:   []string{"通用语", "侏儒语"},
		Traits:      []string{"黑暗视觉", "侏儒聪明"},
		Description: "聪明好奇的小型种族",
	},
	"半精灵": {
		Name: "半精灵",
		AbilityBonuses: map[model.Ability]int{
			model.AbilityCharisma: 2,
		},
		Speed:       30,
		Size:        model.SizeMedium,
		Languages:   []string{"通用语", "精灵语"},
		Traits:      []string{"黑暗视觉", "精类血统", "技能多才多艺"},
		Description: "人类与精灵的混血",
	},
	"半兽人": {
		Name: "半兽人",
		AbilityBonuses: map[model.Ability]int{
			model.AbilityStrength:     2,
			model.AbilityConstitution: 1,
		},
		Speed:       30,
		Size:        model.SizeMedium,
		Languages:   []string{"通用语", "兽人语"},
		Traits:      []string{"黑暗视觉", "恐吓", "坚韧不屈", "野蛮攻击"},
		Description: "人类与兽人的混血",
	},
	"提夫林": {
		Name: "提夫林",
		AbilityBonuses: map[model.Ability]int{
			model.AbilityCharisma:     2,
			model.AbilityIntelligence: 1,
		},
		Speed:       30,
		Size:        model.SizeMedium,
		Languages:   []string{"通用语", "炼狱语"},
		Traits:      []string{"黑暗视觉", "地狱抗性", "炼狱传承"},
		Description: "拥有魔族血统的种族",
	},
}

// GetRace 获取种族定义
func GetRace(name string) *RaceDefinition {
	return Races[name]
}

// GetRaceNames 获取所有种族名称
func GetRaceNames() []string {
	names := make([]string, 0, len(Races))
	for name := range Races {
		names = append(names, name)
	}
	return names
}
