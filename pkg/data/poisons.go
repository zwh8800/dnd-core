package data

import "github.com/zwh8800/dnd-core/pkg/model"

func init() {
	for _, poison := range PoisonDataList {
		GlobalRegistry.RegisterPoison(poison)
	}
}

// PoisonData 毒药数据
type PoisonData = model.PoisonDefinition

// PoisonDataList 毒药数据列表
var PoisonDataList = []PoisonData{
	{
		ID:          "assassins-blood",
		Name:        "刺客之血",
		Type:        model.PoisonInjury,
		Description: "这种毒药用于暗杀，涂抹在武器上",
		Effect: model.PoisonEffect{
			SaveDC:       15,
			Duration:     "24 hours",
			StatusEffect: "poisoned",
			Description:  "中毒生物在24小时内处于中毒状态，除非进行DC 15体质豁免",
		},
		Price:  15000,
		Rarity: "uncommon",
	},
	{
		ID:          "basic-poison",
		Name:        "基础毒药",
		Type:        model.PoisonInjury,
		Description: "一种简单的毒药，可以涂抹在武器上",
		Effect: model.PoisonEffect{
			DamageDice:  "1d4",
			DamageType:  "poison",
			SaveDC:      10,
			Duration:    "1 minute",
			Description: "命中后目标受到1d4毒素伤害，必须进行DC 10体质豁免，失败则在1分钟内持续受到1d4毒素伤害",
		},
		Price:  10000,
		Rarity: "common",
	},
	{
		ID:          "blue-whinnis",
		Name:        "蓝色温尼斯",
		Type:        model.PoisonInjury,
		Description: "这种蓝色液体涂抹在武器上时非常有效",
		Effect: model.PoisonEffect{
			DamageDice:   "1d4",
			DamageType:   "poison",
			SaveDC:       14,
			Duration:     "1 hour",
			StatusEffect: "poisoned",
			Description:  "目标受到1d4毒素伤害，必须进行DC 14体质豁免，失败则中毒1小时",
		},
		Price:  12000,
		Rarity: "common",
	},
	{
		ID:          "blood-root-paste",
		Name:        "血根糊",
		Type:        model.PoisonInjury,
		Description: "从血根植物提取的毒药",
		Effect: model.PoisonEffect{
			DamageDice:  "2d6",
			DamageType:  "poison",
			SaveDC:      15,
			Duration:    "1 minute",
			Description: "目标受到2d6毒素伤害，必须进行DC 15体质豁免，失败则在1分钟内每回合开始时受到1d6毒素伤害",
		},
		Price:  25000,
		Rarity: "uncommon",
	},
	{
		ID:          "carrion-crawler-mucus",
		Name:        "食腐虫黏液",
		Type:        model.PoisonInjury,
		Description: "从食腐虫腺体中提取的麻痹性黏液",
		Effect: model.PoisonEffect{
			SaveDC:       13,
			Duration:     "1 minute",
			StatusEffect: "paralyzed",
			Description:  "目标必须进行DC 13体质豁免，失败则麻痹1分钟，可以每回合结束时重新豁免",
		},
		Price:  15000,
		Rarity: "uncommon",
	},
	{
		ID:          "drow-poison",
		Name:        "卓尔毒药",
		Type:        model.PoisonInjury,
		Description: "卓尔精灵使用的特殊毒药，能使人失去意识",
		Effect: model.PoisonEffect{
			SaveDC:       13,
			Duration:     "1 hour",
			StatusEffect: "unconscious",
			Description:  "目标必须进行DC 13体质豁免，失败则失去意识1小时，受到任何伤害或被人用动作摇晃会唤醒",
		},
		Price:  20000,
		Rarity: "uncommon",
	},
	{
		ID:          "essence-of-ether",
		Name:        "乙醚精华",
		Type:        model.PoisonInhaled,
		Description: "吸入后会导致昏迷的挥发性液体",
		Effect: model.PoisonEffect{
			SaveDC:       15,
			Duration:     "8 hours",
			StatusEffect: "unconscious",
			Description:  "目标必须进行DC 15体质豁免，失败则失去意识8小时",
		},
		Price:  30000,
		Rarity: "rare",
	},
	{
		ID:          "malice",
		Name:        "恶意毒药",
		Type:        model.PoisonInhaled,
		Description: "卓尔精灵使用的吸入性毒药",
		Effect: model.PoisonEffect{
			SaveDC:       15,
			Duration:     "1 hour",
			StatusEffect: "poisoned",
			Description:  "目标必须进行DC 15体质豁免，失败则中毒1小时，期间具有中毒状态",
		},
		Price:  25000,
		Rarity: "uncommon",
	},
	{
		ID:          "midnight-tears",
		Name:        "午夜之泪",
		Type:        model.PoisonIngested,
		Description: "只有在午夜才会发作的强力毒药",
		Effect: model.PoisonEffect{
			DamageDice:  "10d6",
			DamageType:  "poison",
			SaveDC:      17,
			Duration:    "instant",
			Description: "目标在下一个午夜时受到10d6毒素伤害，必须进行DC 17体质豁免，失败则死亡",
		},
		Price:  100000,
		Rarity: "very rare",
	},
	{
		ID:          "purple-worm-poison",
		Name:        "紫虫毒液",
		Type:        model.PoisonInjury,
		Description: "从紫虫提取的致命毒液",
		Effect: model.PoisonEffect{
			DamageDice:  "12d6",
			DamageType:  "poison",
			SaveDC:      19,
			Duration:    "instant",
			Description: "目标受到12d6毒素伤害，必须进行DC 19体质豁免，失败则中毒并在1小时内每回合开始时受到6d6毒素伤害",
		},
		Price:  200000,
		Rarity: "very rare",
	},
}

// GetPoisonData 获取毒药数据
func GetPoisonData(poisonID string) *PoisonData {
	for _, poison := range PoisonDataList {
		if poison.ID == poisonID {
			return &poison
		}
	}
	return nil
}

// ListPoisonData 列出所有毒药数据
func ListPoisonData() []PoisonData {
	return PoisonDataList
}
