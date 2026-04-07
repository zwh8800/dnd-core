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
	"Human": {
		Name: "Human",
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
		Languages:   []string{"Common"},
		Traits:      []string{"Extra Language"},
		Description: "人类是最常见、最多样化的种族",
	},
	"Elf": {
		Name:     "Elf",
		Subraces: []string{"High Elf", "Wood Elf", "Drow"},
		AbilityBonuses: map[model.Ability]int{
			model.AbilityDexterity: 2,
		},
		Speed:       30,
		Size:        model.SizeMedium,
		Languages:   []string{"Common", "Elvish"},
		Traits:      []string{"Darkvision", "Keen Senses", "Fey Ancestry", "Trance"},
		Description: "优雅长寿的森林种族",
	},
	"Dwarf": {
		Name:     "Dwarf",
		Subraces: []string{"Hill Dwarf", "Mountain Dwarf"},
		AbilityBonuses: map[model.Ability]int{
			model.AbilityConstitution: 2,
		},
		Speed:       25,
		Size:        model.SizeMedium,
		Languages:   []string{"Common", "Dwarvish"},
		Traits:      []string{"Darkvision", "Dwarven Resilience", "Combat Training", "Tool Proficiency", "Stonecunning"},
		Description: "坚韧不拔的山地民族",
	},
	"Halfling": {
		Name:     "Halfling",
		Subraces: []string{"Lightfoot", "Stout"},
		AbilityBonuses: map[model.Ability]int{
			model.AbilityDexterity: 2,
		},
		Speed:       25,
		Size:        model.SizeSmall,
		Languages:   []string{"Common", "Halfling"},
		Traits:      []string{"Lucky", "Brave", "Halfling Nimbleness"},
		Description: "小巧而勇敢的民族",
	},
	"Dragonborn": {
		Name: "Dragonborn",
		AbilityBonuses: map[model.Ability]int{
			model.AbilityStrength: 2,
			model.AbilityCharisma: 1,
		},
		Speed:       30,
		Size:        model.SizeMedium,
		Languages:   []string{"Common", "Draconic"},
		Traits:      []string{"Draconic Ancestry", "Breath Weapon", "Damage Resistance"},
		Description: "拥有龙族血统的战士",
	},
	"Gnome": {
		Name:     "Gnome",
		Subraces: []string{"Forest Gnome", "Rock Gnome"},
		AbilityBonuses: map[model.Ability]int{
			model.AbilityIntelligence: 2,
		},
		Speed:       25,
		Size:        model.SizeSmall,
		Languages:   []string{"Common", "Gnomish"},
		Traits:      []string{"Darkvision", "Gnome Cunning"},
		Description: "聪明好奇的小型种族",
	},
	"Half-Elf": {
		Name: "Half-Elf",
		AbilityBonuses: map[model.Ability]int{
			model.AbilityCharisma: 2,
		},
		Speed:       30,
		Size:        model.SizeMedium,
		Languages:   []string{"Common", "Elvish"},
		Traits:      []string{"Darkvision", "Fey Ancestry", "Skill Versatility"},
		Description: "人类与精灵的混血",
	},
	"Half-Orc": {
		Name: "Half-Orc",
		AbilityBonuses: map[model.Ability]int{
			model.AbilityStrength:     2,
			model.AbilityConstitution: 1,
		},
		Speed:       30,
		Size:        model.SizeMedium,
		Languages:   []string{"Common", "Orc"},
		Traits:      []string{"Darkvision", "Menacing", "Relentless Endurance", "Savage Attacks"},
		Description: "人类与兽人的混血",
	},
	"Tiefling": {
		Name: "Tiefling",
		AbilityBonuses: map[model.Ability]int{
			model.AbilityCharisma:     2,
			model.AbilityIntelligence: 1,
		},
		Speed:       30,
		Size:        model.SizeMedium,
		Languages:   []string{"Common", "Infernal"},
		Traits:      []string{"Darkvision", "Hellish Resistance", "Infernal Legacy"},
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
