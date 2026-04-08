package data

import (
	"github.com/zwh8800/dnd-core/pkg/model"
)

// CraftingRecipes 制作配方数据库
var CraftingRecipes = map[string]model.CraftingRecipe{
	// === 非魔法物品制作 ===
	"antitoxin": {
		ID:          "antitoxin",
		Name:        "解毒剂",
		Type:        model.CraftingTypeNonMagical,
		Description: "一剂解毒剂，服用后1小时内对中毒状态免疫",
		Materials: []model.CraftingMaterial{
			{Name: "草药", Quantity: 2, Cost: 50},
			{Name: "纯净水", Quantity: 1, Cost: 10},
		},
		ToolsRequired: []string{"炼金工具"},
		TimeDays:      1,
		DC:            10,
		MinLevel:      1,
		Cost:          500, // 5gp = 500cp
	},
	"tinker-tools": {
		ID:          "tinker-tools",
		Name:        "修补工具",
		Type:        model.CraftingTypeNonMagical,
		Description: "一套修补工具，用于制作小型机械装置",
		Materials: []model.CraftingMaterial{
			{Name: "铁锭", Quantity: 5, Cost: 20},
			{Name: "皮革", Quantity: 2, Cost: 40},
		},
		ToolsRequired: []string{"锻造工具"},
		TimeDays:      3,
		DC:            12,
		MinLevel:      1,
		Cost:          5000, // 50gp
	},
	"longsword": {
		ID:          "longsword",
		Name:        "长剑",
		Type:        model.CraftingTypeNonMagical,
		Description: "军用近战武器，1d8挥砍伤害",
		Materials: []model.CraftingMaterial{
			{Name: "钢锭", Quantity: 10, Cost: 10},
			{Name: "木材", Quantity: 2, Cost: 5},
			{Name: "皮革", Quantity: 1, Cost: 10},
		},
		ToolsRequired: []string{"锻造工具"},
		TimeDays:      5,
		DC:            15,
		MinLevel:      1,
		Cost:          1500, // 15gp
	},
	"chain-mail": {
		ID:          "chain-mail",
		Name:        "链甲",
		Type:        model.CraftingTypeNonMagical,
		Description: "重甲，AC 16，隐匿劣势",
		Materials: []model.CraftingMaterial{
			{Name: "铁环", Quantity: 50, Cost: 5},
			{Name: "皮革", Quantity: 5, Cost: 10},
		},
		ToolsRequired: []string{"锻造工具", "锁甲匠工具"},
		TimeDays:      10,
		DC:            15,
		MinLevel:      3,
		Cost:          7500, // 75gp
	},

	// === 药水制作 ===
	"potion-healing": {
		ID:          "potion-healing",
		Name:        "治疗药水",
		Type:        model.CraftingTypePotion,
		Description: "饮用后恢复2d4+2生命值",
		Materials: []model.CraftingMaterial{
			{Name: "治疗草药", Quantity: 3, Cost: 25},
			{Name: "纯净水", Quantity: 1, Cost: 10},
			{Name: "玻璃瓶", Quantity: 1, Cost: 5},
		},
		ToolsRequired: []string{"炼金工具"},
		TimeDays:      1,
		DC:            10,
		MinLevel:      1,
		Cost:          5000, // 50gp
	},
	"potion-greater-healing": {
		ID:          "potion-greater-healing",
		Name:        "强效治疗药水",
		Type:        model.CraftingTypePotion,
		Description: "饮用后恢复4d4+4生命值",
		Materials: []model.CraftingMaterial{
			{Name: "稀有草药", Quantity: 5, Cost: 100},
			{Name: "魔法精华", Quantity: 2, Cost: 200},
			{Name: "水晶瓶", Quantity: 1, Cost: 50},
		},
		ToolsRequired: []string{"炼金工具"},
		TimeDays:      3,
		DC:            15,
		MinLevel:      5,
		Cost:          15000, // 150gp
	},
	"potion-fire-breath": {
		ID:          "potion-fire-breath",
		Name:        "火焰吐息药水",
		Type:        model.CraftingTypePotion,
		Description: "饮用后可喷出15尺锥状火焰，造成4d6火焰伤害",
		Materials: []model.CraftingMaterial{
			{Name: "火蜥蜴血", Quantity: 1, Cost: 200},
			{Name: "硫磺", Quantity: 2, Cost: 50},
			{Name: "水晶瓶", Quantity: 1, Cost: 50},
		},
		ToolsRequired: []string{"炼金工具"},
		TimeDays:      2,
		DC:            13,
		MinLevel:      3,
		Cost:          10000, // 100gp
	},
	"potion-invisibility": {
		ID:          "potion-invisibility",
		Name:        "隐形药水",
		Type:        model.CraftingTypePotion,
		Description: "饮用后获得隐形状态1小时",
		Materials: []model.CraftingMaterial{
			{Name: "隐形菇", Quantity: 2, Cost: 300},
			{Name: "魔法精华", Quantity: 3, Cost: 200},
			{Name: "水晶瓶", Quantity: 1, Cost: 50},
		},
		ToolsRequired: []string{"炼金工具"},
		TimeDays:      5,
		DC:            17,
		MinLevel:      7,
		Cost:          25000, // 250gp
	},

	// === 法术卷轴制作 ===
	"scroll-fireball": {
		ID:          "scroll-fireball",
		Name:        "火球术卷轴",
		Type:        model.CraftingTypeScroll,
		Description: "施放后可释放3环火球术",
		Materials: []model.CraftingMaterial{
			{Name: "优质羊皮纸", Quantity: 1, Cost: 100},
			{Name: "魔法墨水", Quantity: 1, Cost: 200},
			{Name: "火元素精华", Quantity: 1, Cost: 300},
		},
		ToolsRequired: []string{"书法工具"},
		TimeDays:      3,
		DC:            15,
		MinLevel:      5,
		SpellRequired: "火球术",
		Cost:          30000, // 300gp
	},
	"scroll-shield": {
		ID:          "scroll-shield",
		Name:        "护盾术卷轴",
		Type:        model.CraftingTypeScroll,
		Description: "施放后可释放1环护盾术",
		Materials: []model.CraftingMaterial{
			{Name: "羊皮纸", Quantity: 1, Cost: 50},
			{Name: "魔法墨水", Quantity: 1, Cost: 100},
		},
		ToolsRequired: []string{"书法工具"},
		TimeDays:      1,
		DC:            11,
		MinLevel:      1,
		SpellRequired: "护盾术",
		Cost:          5000, // 50gp
	},

	// === 魔法物品制作 ===
	"wand-magic-missiles": {
		ID:          "wand-magic-missiles",
		Name:        "魔法飞弹魔杖",
		Type:        model.CraftingTypeMagicItem,
		Description: "7充能魔杖，每日黎明恢复1d6+1充能。消耗1充能释放魔法飞弹",
		Materials: []model.CraftingMaterial{
			{Name: "魔法木材", Quantity: 1, Cost: 500},
			{Name: "奥术水晶", Quantity: 3, Cost: 300},
			{Name: "魔法精华", Quantity: 5, Cost: 200},
		},
		ToolsRequired: []string{"木工工具"},
		TimeDays:      14,
		DC:            15,
		MinLevel:      5,
		SpellRequired: "魔法飞弹",
		Cost:          25000, // 250gp
	},
	"weapon-plus1": {
		ID:          "weapon-plus1",
		Name:        "+1武器",
		Type:        model.CraftingTypeMagicItem,
		Description: "攻击和伤害掷骰获得+1加值的魔法武器",
		Materials: []model.CraftingMaterial{
			{Name: "精金", Quantity: 5, Cost: 500},
			{Name: "魔法精华", Quantity: 10, Cost: 200},
			{Name: "附魔粉末", Quantity: 3, Cost: 400},
		},
		ToolsRequired: []string{"锻造工具"},
		TimeDays:      30,
		DC:            17,
		MinLevel:      5,
		SpellRequired: "魔法武器",
		Cost:          50000, // 500gp
	},
	"armor-plus1": {
		ID:          "armor-plus1",
		Name:        "+1护甲",
		Type:        model.CraftingTypeMagicItem,
		Description: "AC获得+1加值的魔法护甲",
		Materials: []model.CraftingMaterial{
			{Name: "精金", Quantity: 10, Cost: 500},
			{Name: "魔法精华", Quantity: 10, Cost: 200},
			{Name: "防护符文", Quantity: 2, Cost: 600},
		},
		ToolsRequired: []string{"锻造工具", "锁甲匠工具"},
		TimeDays:      30,
		DC:            17,
		MinLevel:      5,
		SpellRequired: "魔法武器",
		Cost:          50000, // 500gp
	},
}

// GetCraftingRecipe 获取指定配方
func GetCraftingRecipe(recipeID string) (model.CraftingRecipe, bool) {
	recipe, exists := CraftingRecipes[recipeID]
	return recipe, exists
}

// GetAllCraftingRecipes 获取所有配方
func GetAllCraftingRecipes() []model.CraftingRecipe {
	recipes := make([]model.CraftingRecipe, 0, len(CraftingRecipes))
	for _, recipe := range CraftingRecipes {
		recipes = append(recipes, recipe)
	}
	return recipes
}

// GetRecipesByType 按类型获取配方
func GetRecipesByType(craftingType model.CraftingType) []model.CraftingRecipe {
	recipes := make([]model.CraftingRecipe, 0)
	for _, recipe := range CraftingRecipes {
		if recipe.Type == craftingType {
			recipes = append(recipes, recipe)
		}
	}
	return recipes
}
