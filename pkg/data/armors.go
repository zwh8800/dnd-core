package data

import (
	"github.com/zwh8800/dnd-core/pkg/model"
)

// init 注册所有护甲
func init() {
	for _, armor := range Armors {
		GlobalRegistry.RegisterArmor(armor)
	}
}

// Armors 护甲数据库（16 种护甲 + 盾牌）
var Armors = []*model.Item{
	// ========== 轻甲 ==========
	{
		ID:          "padded-armor",
		Name:        "棉甲",
		Description: "由多层填充物和绗缝制成的护甲",
		Type:        model.ItemTypeArmor,
		Weight:      8.0,
		Value:       500, // 5 gp
		ArmorProps: &model.ArmorProperties{
			BaseAC:              11,
			StealthDisadvantage: false,
		},
	},
	{
		ID:          "leather-armor",
		Name:        "皮甲",
		Description: "经过硬化处理的皮革制成的护甲",
		Type:        model.ItemTypeArmor,
		Weight:      10.0,
		Value:       1000, // 10 gp
		ArmorProps: &model.ArmorProperties{
			BaseAC:              11,
			StealthDisadvantage: false,
		},
	},
	{
		ID:          "studded-leather",
		Name:        "镶嵌皮甲",
		Description: "带有铆钉加固的皮革护甲",
		Type:        model.ItemTypeArmor,
		Weight:      13.0,
		Value:       4500, // 45 gp
		ArmorProps: &model.ArmorProperties{
			BaseAC:              12,
			StealthDisadvantage: false,
		},
	},

	// ========== 中甲 ==========
	{
		ID:          "hide-armor",
		Name:        "毛皮甲",
		Description: "由野兽毛皮制成的粗糙护甲",
		Type:        model.ItemTypeArmor,
		Weight:      12.0,
		Value:       1000, // 10 gp
		ArmorProps: &model.ArmorProperties{
			BaseAC:              12,
			MaxDexModifier:      ptrInt(2),
			StealthDisadvantage: false,
		},
	},
	{
		ID:          "chain-shirt",
		Name:        "链甲衫",
		Description: "由金属环互相连接的轻型链甲",
		Type:        model.ItemTypeArmor,
		Weight:      20.0,
		Value:       5000, // 50 gp
		ArmorProps: &model.ArmorProperties{
			BaseAC:              13,
			MaxDexModifier:      ptrInt(2),
			StealthDisadvantage: false,
		},
	},
	{
		ID:          "scale-mail",
		Name:        "鳞甲",
		Description: "由许多小金属片缝制在衣服上的护甲",
		Type:        model.ItemTypeArmor,
		Weight:      45.0,
		Value:       5000, // 50 gp
		ArmorProps: &model.ArmorProperties{
			BaseAC:              14,
			MaxDexModifier:      ptrInt(2),
			StealthDisadvantage: true,
		},
	},
	{
		ID:          "breastplate",
		Name:        "胸甲",
		Description: "贴合身体的金属胸甲",
		Type:        model.ItemTypeArmor,
		Weight:      20.0,
		Value:       40000, // 400 gp
		ArmorProps: &model.ArmorProperties{
			BaseAC:              14,
			MaxDexModifier:      ptrInt(2),
			StealthDisadvantage: false,
		},
	},
	{
		ID:          "half-plate-armor",
		Name:        "半身甲",
		Description: "覆盖大部分身体的塑形金属板",
		Type:        model.ItemTypeArmor,
		Weight:      40.0,
		Value:       75000, // 750 gp
		ArmorProps: &model.ArmorProperties{
			BaseAC:              15,
			MaxDexModifier:      ptrInt(2),
			StealthDisadvantage: true,
		},
	},

	// ========== 重甲 ==========
	{
		ID:          "ring-mail",
		Name:        "环甲",
		Description: "由金属环缝制在皮革上的护甲",
		Type:        model.ItemTypeArmor,
		Weight:      40.0,
		Value:       3000, // 30 gp
		ArmorProps: &model.ArmorProperties{
			BaseAC:              14,
			StealthDisadvantage: true,
			StrengthRequirement: 0,
		},
	},
	{
		ID:          "chain-mail",
		Name:        "锁甲",
		Description: "由金属环互相链接而成的护甲",
		Type:        model.ItemTypeArmor,
		Weight:      55.0,
		Value:       7500, // 75 gp
		ArmorProps: &model.ArmorProperties{
			BaseAC:              16,
			StealthDisadvantage: true,
			StrengthRequirement: 13,
		},
	},
	{
		ID:          "splint-armor",
		Name:        "夹板甲",
		Description: "由垂直金属条固定在皮革衬里上的护甲",
		Type:        model.ItemTypeArmor,
		Weight:      60.0,
		Value:       20000, // 200 gp
		ArmorProps: &model.ArmorProperties{
			BaseAC:              17,
			StealthDisadvantage: true,
			StrengthRequirement: 15,
		},
	},
	{
		ID:          "plate-armor",
		Name:        "全身甲",
		Description: "由定制锻造的互锁钢板组成的护甲",
		Type:        model.ItemTypeArmor,
		Weight:      65.0,
		Value:       150000, // 1500 gp
		ArmorProps: &model.ArmorProperties{
			BaseAC:              18,
			StealthDisadvantage: true,
			StrengthRequirement: 15,
		},
	},

	// ========== 盾牌 ==========
	{
		ID:          "shield",
		Name:        "盾牌",
		Description: "一面木制或金属制的盾牌",
		Type:        model.ItemTypeArmor,
		Weight:      6.0,
		Value:       1000, // 10 gp
		ArmorProps: &model.ArmorProperties{
			BaseAC:              2, // +2 AC
			StealthDisadvantage: false,
		},
	},
}

func ptrInt(i int) *int {
	return &i
}
