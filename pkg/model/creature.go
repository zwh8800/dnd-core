package model

// CreatureType 代表 SRD 5.2.1 中的生物类型
type CreatureType string

const (
	CreatureTypeAberration  CreatureType = "Aberration"  // 异怪
	CreatureTypeBeast       CreatureType = "Beast"       // 野兽
	CreatureTypeCelestial   CreatureType = "Celestial"   // 天界生物
	CreatureTypeConstruct   CreatureType = "Construct"   // 构装体
	CreatureTypeDragon      CreatureType = "Dragon"      // 龙类
	CreatureTypeElemental   CreatureType = "Elemental"   // 元素生物
	CreatureTypeFey         CreatureType = "Fey"         // 精类
	CreatureTypeFiend       CreatureType = "Fiend"       // 邪魔
	CreatureTypeGiant       CreatureType = "Giant"       // 巨人
	CreatureTypeHumanoid    CreatureType = "Humanoid"    // 类人生物
	CreatureTypeMonstrosity CreatureType = "Monstrosity" // 怪兽
	CreatureTypeOoze        CreatureType = "Ooze"        // 泥怪
	CreatureTypePlant       CreatureType = "Plant"       // 植物
	CreatureTypeUndead      CreatureType = "Undead"      // 不死生物
)

// CreatureTag 代表生物的标签（用于更细粒度的分类）
type CreatureTag string

const (
	CreatureTagHuman  CreatureTag = "Human"
	CreatureTagElf    CreatureTag = "Elf"
	CreatureTagDwarf  CreatureTag = "Dwarf"
	CreatureTagGoblin CreatureTag = "Goblinoid"
	CreatureTagOrc    CreatureTag = "Orc"
	CreatureTagKobold CreatureTag = "Kobold"
	CreatureTagDemon  CreatureTag = "Demon"
	CreatureTagDevil  CreatureTag = "Devil"
	CreatureTagAngel  CreatureTag = "Angel"
	CreatureTagGnome  CreatureTag = "Gnome"
)
