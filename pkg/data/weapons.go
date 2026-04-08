package data

import (
	"github.com/zwh8800/dnd-core/pkg/model"
)

// init 注册所有武器
func init() {
	for _, weapon := range Weapons {
		GlobalRegistry.RegisterWeapon(weapon)
	}
}

// Weapons 武器数据库（30+ 武器）
var Weapons = []*model.Item{
	// ========== 简易近战武器 ==========
	{
		ID:          "club",
		Name:        "木棒",
		Description: "一根沉重的木棍",
		Type:        model.ItemTypeWeapon,
		Weight:      2.0,
		Value:       10, // 1 sp = 10 cp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1d4",
			DamageType: model.DamageTypeBludgeoning,
			WeaponType: "melee",
		},
	},
	{
		ID:          "dagger",
		Name:        "匕首",
		Description: "一把锋利的短刀",
		Type:        model.ItemTypeWeapon,
		Weight:      1.0,
		Value:       200, // 2 gp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1d4",
			DamageType: model.DamageTypePiercing,
			WeaponType: "melee",
			Range:      20,
			LongRange:  60,
			Finesse:    true,
			Light:      true,
			Thrown:     true,
		},
	},
	{
		ID:          "greatclub",
		Name:        "巨木棒",
		Description: "一根粗大的木棒",
		Type:        model.ItemTypeWeapon,
		Weight:      10.0,
		Value:       20, // 2 sp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1d8",
			DamageType: model.DamageTypeBludgeoning,
			WeaponType: "melee",
			TwoHanded:  true,
		},
	},
	{
		ID:          "handaxe",
		Name:        "手斧",
		Description: "一把适合投掷的小斧头",
		Type:        model.ItemTypeWeapon,
		Weight:      2.0,
		Value:       500, // 5 gp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1d6",
			DamageType: model.DamageTypeSlashing,
			WeaponType: "melee",
			Range:      20,
			LongRange:  60,
			Light:      true,
			Thrown:     true,
		},
	},
	{
		ID:          "javelin",
		Name:        "标枪",
		Description: "一根平衡的投掷矛",
		Type:        model.ItemTypeWeapon,
		Weight:      2.0,
		Value:       50, // 5 sp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1d6",
			DamageType: model.DamageTypePiercing,
			WeaponType: "melee",
			Range:      30,
			LongRange:  120,
			Thrown:     true,
		},
	},
	{
		ID:          "light-hammer",
		Name:        "轻锤",
		Description: "一把小巧的锤子",
		Type:        model.ItemTypeWeapon,
		Weight:      2.0,
		Value:       200, // 2 gp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1d4",
			DamageType: model.DamageTypeBludgeoning,
			WeaponType: "melee",
			Range:      20,
			LongRange:  60,
			Light:      true,
			Thrown:     true,
		},
	},
	{
		ID:          "mace",
		Name:        "锤矛",
		Description: "一把带有金属头的权杖",
		Type:        model.ItemTypeWeapon,
		Weight:      4.0,
		Value:       500, // 5 gp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1d6",
			DamageType: model.DamageTypeBludgeoning,
			WeaponType: "melee",
		},
	},
	{
		ID:          "quarterstaff",
		Name:        "木杖",
		Description: "一根坚固的长木棍",
		Type:        model.ItemTypeWeapon,
		Weight:      4.0,
		Value:       20, // 2 sp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1d6",
			DamageType: model.DamageTypeBludgeoning,
			WeaponType: "melee",
			Versatile:  "1d8",
			TwoHanded:  true,
		},
	},
	{
		ID:          "sickle",
		Name:        "镰刀",
		Description: "一把小型农用镰刀",
		Type:        model.ItemTypeWeapon,
		Weight:      2.0,
		Value:       100, // 1 gp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1d4",
			DamageType: model.DamageTypeSlashing,
			WeaponType: "melee",
			Light:      true,
		},
	},
	{
		ID:          "spear",
		Name:        "长矛",
		Description: "一根带尖的长杆武器",
		Type:        model.ItemTypeWeapon,
		Weight:      3.0,
		Value:       100, // 1 gp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1d6",
			DamageType: model.DamageTypePiercing,
			WeaponType: "melee",
			Range:      20,
			LongRange:  60,
			Versatile:  "1d8",
			Thrown:     true,
		},
	},

	// ========== 简易远程武器 ==========
	{
		ID:          "light-crossbow",
		Name:        "轻弩",
		Description: "一种小型弩",
		Type:        model.ItemTypeWeapon,
		Weight:      5.0,
		Value:       2500, // 25 gp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1d8",
			DamageType: model.DamageTypePiercing,
			WeaponType: "ranged",
			Range:      80,
			LongRange:  320,
			TwoHanded:  true,
			Loading:    true,
			Ammunition: true,
		},
	},
	{
		ID:          "dart",
		Name:        "飞镖",
		Description: "一根小型投掷武器",
		Type:        model.ItemTypeWeapon,
		Weight:      0.25,
		Value:       5, // 5 cp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1d4",
			DamageType: model.DamageTypePiercing,
			WeaponType: "ranged",
			Range:      20,
			LongRange:  60,
			Finesse:    true,
			Thrown:     true,
		},
	},
	{
		ID:          "shortbow",
		Name:        "短弓",
		Description: "一张小巧的弓",
		Type:        model.ItemTypeWeapon,
		Weight:      2.0,
		Value:       2500, // 25 gp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1d6",
			DamageType: model.DamageTypePiercing,
			WeaponType: "ranged",
			Range:      80,
			LongRange:  320,
			TwoHanded:  true,
			Ammunition: true,
		},
	},
	{
		ID:          "sling",
		Name:        "投石索",
		Description: "一根用于投掷石块的皮带",
		Type:        model.ItemTypeWeapon,
		Weight:      0.0,
		Value:       10, // 1 sp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1d4",
			DamageType: model.DamageTypeBludgeoning,
			WeaponType: "ranged",
			Range:      30,
			LongRange:  120,
			Ammunition: true,
		},
	},

	// ========== 军用近战武器 ==========
	{
		ID:          "battleaxe",
		Name:        "战斧",
		Description: "一把重型战斗斧",
		Type:        model.ItemTypeWeapon,
		Weight:      4.0,
		Value:       1000, // 10 gp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1d8",
			DamageType: model.DamageTypeSlashing,
			WeaponType: "melee",
			Versatile:  "1d10",
		},
	},
	{
		ID:          "flail",
		Name:        "连枷",
		Description: "一个带链条的钉头锤",
		Type:        model.ItemTypeWeapon,
		Weight:      2.0,
		Value:       1000, // 10 gp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1d8",
			DamageType: model.DamageTypeBludgeoning,
			WeaponType: "melee",
		},
	},
	{
		ID:          "glaive",
		Name:        "长刀",
		Description: "一根带刃的长杆武器",
		Type:        model.ItemTypeWeapon,
		Weight:      6.0,
		Value:       2000, // 20 gp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1d10",
			DamageType: model.DamageTypeSlashing,
			WeaponType: "melee",
			Heavy:      true,
			TwoHanded:  true,
			Reach:      true,
		},
	},
	{
		ID:          "greataxe",
		Name:        "巨斧",
		Description: "一把巨大的双手斧",
		Type:        model.ItemTypeWeapon,
		Weight:      7.0,
		Value:       3000, // 30 gp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1d12",
			DamageType: model.DamageTypeSlashing,
			WeaponType: "melee",
			Heavy:      true,
			TwoHanded:  true,
		},
	},
	{
		ID:          "greatsword",
		Name:        "巨剑",
		Description: "一把大型双手剑",
		Type:        model.ItemTypeWeapon,
		Weight:      6.0,
		Value:       5000, // 50 gp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "2d6",
			DamageType: model.DamageTypeSlashing,
			WeaponType: "melee",
			Heavy:      true,
			TwoHanded:  true,
		},
	},
	{
		ID:          "halberd",
		Name:        "戟",
		Description: "一种多功能的长柄武器",
		Type:        model.ItemTypeWeapon,
		Weight:      6.0,
		Value:       2000, // 20 gp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1d10",
			DamageType: model.DamageTypeSlashing,
			WeaponType: "melee",
			Heavy:      true,
			TwoHanded:  true,
			Reach:      true,
		},
	},
	{
		ID:          "lance",
		Name:        "骑枪",
		Description: "一根骑乘用的长矛",
		Type:        model.ItemTypeWeapon,
		Weight:      6.0,
		Value:       1000, // 10 gp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1d12",
			DamageType: model.DamageTypePiercing,
			WeaponType: "melee",
			Reach:      true,
			Special:    []string{"骑乘时单手使用，否则需双手"},
		},
	},
	{
		ID:          "longsword",
		Name:        "长剑",
		Description: "一把标准的军用剑",
		Type:        model.ItemTypeWeapon,
		Weight:      3.0,
		Value:       1500, // 15 gp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1d8",
			DamageType: model.DamageTypeSlashing,
			WeaponType: "melee",
			Versatile:  "1d10",
		},
	},
	{
		ID:          "maul",
		Name:        "巨锤",
		Description: "一个大型的战锤",
		Type:        model.ItemTypeWeapon,
		Weight:      10.0,
		Value:       1000, // 10 gp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "2d6",
			DamageType: model.DamageTypeBludgeoning,
			WeaponType: "melee",
			Heavy:      true,
			TwoHanded:  true,
		},
	},
	{
		ID:          "morningstar",
		Name:        "流星锤",
		Description: "一个带刺的锤头武器",
		Type:        model.ItemTypeWeapon,
		Weight:      4.0,
		Value:       1500, // 15 gp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1d8",
			DamageType: model.DamageTypePiercing,
			WeaponType: "melee",
		},
	},
	{
		ID:          "pike",
		Name:        "长枪",
		Description: "一根超长的刺矛",
		Type:        model.ItemTypeWeapon,
		Weight:      18.0,
		Value:       500, // 5 gp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1d10",
			DamageType: model.DamageTypePiercing,
			WeaponType: "melee",
			Heavy:      true,
			TwoHanded:  true,
			Reach:      true,
		},
	},
	{
		ID:          "rapier",
		Name:        "细剑",
		Description: "一把纤细的刺剑",
		Type:        model.ItemTypeWeapon,
		Weight:      2.0,
		Value:       2500, // 25 gp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1d8",
			DamageType: model.DamageTypePiercing,
			WeaponType: "melee",
			Finesse:    true,
		},
	},
	{
		ID:          "scimitar",
		Name:        "弯刀",
		Description: "一把弯曲的单刃剑",
		Type:        model.ItemTypeWeapon,
		Weight:      3.0,
		Value:       2500, // 25 gp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1d6",
			DamageType: model.DamageTypeSlashing,
			WeaponType: "melee",
			Light:      true,
			Finesse:    true,
		},
	},
	{
		ID:          "shortsword",
		Name:        "短剑",
		Description: "一把小型的剑",
		Type:        model.ItemTypeWeapon,
		Weight:      2.0,
		Value:       1000, // 10 gp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1d6",
			DamageType: model.DamageTypePiercing,
			WeaponType: "melee",
			Light:      true,
			Finesse:    true,
		},
	},
	{
		ID:          "trident",
		Name:        "三叉戟",
		Description: "一把三叉的矛",
		Type:        model.ItemTypeWeapon,
		Weight:      4.0,
		Value:       500, // 5 gp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1d6",
			DamageType: model.DamageTypePiercing,
			WeaponType: "melee",
			Range:      20,
			LongRange:  60,
			Versatile:  "1d8",
			Thrown:     true,
		},
	},
	{
		ID:          "war-pick",
		Name:        "战镐",
		Description: "一把用于战斗的镐",
		Type:        model.ItemTypeWeapon,
		Weight:      2.0,
		Value:       500, // 5 gp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1d8",
			DamageType: model.DamageTypePiercing,
			WeaponType: "melee",
		},
	},
	{
		ID:          "warhammer",
		Name:        "战锤",
		Description: "一把重型战斗锤",
		Type:        model.ItemTypeWeapon,
		Weight:      2.0,
		Value:       1500, // 15 gp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1d8",
			DamageType: model.DamageTypeBludgeoning,
			WeaponType: "melee",
			Versatile:  "1d10",
		},
	},
	{
		ID:          "whip",
		Name:        "长鞭",
		Description: "一条带刺的长鞭",
		Type:        model.ItemTypeWeapon,
		Weight:      3.0,
		Value:       200, // 2 gp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1d4",
			DamageType: model.DamageTypeSlashing,
			WeaponType: "melee",
			Finesse:    true,
			Reach:      true,
		},
	},

	// ========== 军用远程武器 ==========
	{
		ID:          "blowgun",
		Name:        "吹箭筒",
		Description: "一根用于吹射毒针的管子",
		Type:        model.ItemTypeWeapon,
		Weight:      1.0,
		Value:       1000, // 10 gp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1",
			DamageType: model.DamageTypePiercing,
			WeaponType: "ranged",
			Range:      25,
			LongRange:  100,
			Loading:    true,
			Ammunition: true,
		},
	},
	{
		ID:          "hand-crossbow",
		Name:        "手弩",
		Description: "一种单手弩",
		Type:        model.ItemTypeWeapon,
		Weight:      3.0,
		Value:       7500, // 75 gp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1d6",
			DamageType: model.DamageTypePiercing,
			WeaponType: "ranged",
			Range:      30,
			LongRange:  120,
			Light:      true,
			Loading:    true,
			Ammunition: true,
		},
	},
	{
		ID:          "heavy-crossbow",
		Name:        "重弩",
		Description: "一种大型强弩",
		Type:        model.ItemTypeWeapon,
		Weight:      18.0,
		Value:       5000, // 50 gp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1d10",
			DamageType: model.DamageTypePiercing,
			WeaponType: "ranged",
			Range:      100,
			LongRange:  400,
			Heavy:      true,
			TwoHanded:  true,
			Loading:    true,
			Ammunition: true,
		},
	},
	{
		ID:          "longbow",
		Name:        "长弓",
		Description: "一张大型强力弓",
		Type:        model.ItemTypeWeapon,
		Weight:      2.0,
		Value:       5000, // 50 gp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "1d8",
			DamageType: model.DamageTypePiercing,
			WeaponType: "ranged",
			Range:      150,
			LongRange:  600,
			Heavy:      true,
			TwoHanded:  true,
			Ammunition: true,
		},
	},
	{
		ID:          "net",
		Name:        "网",
		Description: "一张用于捕捉的网",
		Type:        model.ItemTypeWeapon,
		Weight:      3.0,
		Value:       100, // 1 gp
		WeaponProps: &model.WeaponProperties{
			DamageDice: "",
			WeaponType: "ranged",
			Range:      5,
			LongRange:  15,
			Thrown:     true,
			Special:    []string{"命中时目标受束"},
		},
	},
}
