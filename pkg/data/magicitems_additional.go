package data

import "github.com/zwh8800/dnd-core/pkg/model"

// AdditionalMagicItems 补充魔法物品数据
var AdditionalMagicItems = []*model.Item{
	// 常见
	{ID: "common-1", Name: "魔法火花", Type: model.ItemTypeWondrousItem, Rarity: model.RarityCommon, Value: 1000, Description: "指尖产生小火花"},
	{ID: "common-2", Name: "永不熄灭的蜡烛", Type: model.ItemTypeWondrousItem, Rarity: model.RarityCommon, Value: 500, Description: "永远燃烧的蜡烛"},
	{ID: "common-3", Name: "魔法粉笔", Type: model.ItemTypeWondrousItem, Rarity: model.RarityCommon, Value: 800, Description: "画出暂时发光的线条"},
	// 非普通
	{ID: "uncommon-1", Name: "敏捷手套", Type: model.ItemTypeWondrousItem, Rarity: model.RarityUncommon, Value: 5000, Attunement: "需要调音", MagicEffects: []string{"dex+1"}, Description: "敏捷+1"},
	{ID: "uncommon-2", Name: "体质护符", Type: model.ItemTypeWondrousItem, Rarity: model.RarityUncommon, Value: 5000, Attunement: "需要调音", MagicEffects: []string{"con+1"}, Description: "体质+1"},
	{ID: "uncommon-3", Name: "智力头环", Type: model.ItemTypeWondrousItem, Rarity: model.RarityUncommon, Value: 5000, Attunement: "需要调音", MagicEffects: []string{"int+1"}, Description: "智力+1"},
	{ID: "uncommon-4", Name: "感知头带", Type: model.ItemTypeWondrousItem, Rarity: model.RarityUncommon, Value: 5000, Attunement: "需要调音", MagicEffects: []string{"wis+1"}, Description: "感知+1"},
	{ID: "uncommon-5", Name: "魅力徽章", Type: model.ItemTypeWondrousItem, Rarity: model.RarityUncommon, Value: 5000, Attunement: "需要调音", MagicEffects: []string{"cha+1"}, Description: "魅力+1"},
	{ID: "uncommon-6", Name: "治疗戒指", Type: model.ItemTypeRing, Rarity: model.RarityUncommon, Value: 10000, Attunement: "需要调音", Charges: 3, MaxCharges: 3, Recharge: "dawn", Description: "消耗1充能施展治疗术"},
	{ID: "uncommon-7", Name: "防护卷轴", Type: model.ItemTypeScroll, Rarity: model.RarityUncommon, Value: 8000, Consumable: true, Description: "展开后获得临时AC+2"},
	{ID: "uncommon-8", Name: "火焰匕首", Type: model.ItemTypeWeapon, Rarity: model.RarityUncommon, Value: 12000, Attunement: "需要调音", MagicBonus: 1, Description: "+1匕首，命中额外1d4火焰伤害"},
	{ID: "uncommon-9", Name: "冰冻长剑", Type: model.ItemTypeWeapon, Rarity: model.RarityUncommon, Value: 12000, Attunement: "需要调音", MagicBonus: 1, Description: "+1长剑，命中额外1d4寒冷伤害"},
	{ID: "uncommon-10", Name: "闪电战锤", Type: model.ItemTypeWeapon, Rarity: model.RarityUncommon, Value: 12000, Attunement: "需要调音", MagicBonus: 1, Description: "+1战锤，命中额外1d4闪电伤害"},
	// 稀有
	{ID: "rare-1", Name: "力量护腕", Type: model.ItemTypeWondrousItem, Rarity: model.RarityRare, Value: 25000, Attunement: "需要调音", MagicEffects: []string{"str-19"}, Description: "力量设为19"},
	{ID: "rare-2", Name: "敏捷斗篷", Type: model.ItemTypeWondrousItem, Rarity: model.RarityRare, Value: 25000, Attunement: "需要调音", MagicEffects: []string{"dex-19"}, Description: "敏捷设为19"},
	{ID: "rare-3", Name: "体质护符", Type: model.ItemTypeWondrousItem, Rarity: model.RarityRare, Value: 25000, Attunement: "需要调音", MagicEffects: []string{"con-19"}, Description: "体质设为19"},
	{ID: "rare-4", Name: "智力护符", Type: model.ItemTypeWondrousItem, Rarity: model.RarityRare, Value: 25000, Attunement: "需要调音", MagicEffects: []string{"int-19"}, Description: "智力设为19"},
	{ID: "rare-5", Name: "感知护符", Type: model.ItemTypeWondrousItem, Rarity: model.RarityRare, Value: 25000, Attunement: "需要调音", MagicEffects: []string{"wis-19"}, Description: "感知设为19"},
	{ID: "rare-6", Name: "魅力护符", Type: model.ItemTypeWondrousItem, Rarity: model.RarityRare, Value: 25000, Attunement: "需要调音", MagicEffects: []string{"cha-19"}, Description: "魅力设为19"},
	{ID: "rare-7", Name: "飞行斗篷", Type: model.ItemTypeWondrousItem, Rarity: model.RarityRare, Value: 30000, Attunement: "需要调音", MagicEffects: []string{"fly-speed-60"}, Description: "获得60尺飞行速度"},
	{ID: "rare-8", Name: "隐形斗篷", Type: model.ItemTypeWondrousItem, Rarity: model.RarityRare, Value: 35000, Attunement: "需要调音", MagicEffects: []string{"invisibility"}, Description: "可用动作变为隐形"},
	{ID: "rare-9", Name: "防护戒指", Type: model.ItemTypeRing, Rarity: model.RarityRare, Value: 30000, Attunement: "需要调音", MagicEffects: []string{"ac+1"}, Description: "AC+1"},
	{ID: "rare-10", Name: "抵抗戒指", Type: model.ItemTypeRing, Rarity: model.RarityRare, Value: 30000, Attunement: "需要调音", MagicEffects: []string{"all-saves+1"}, Description: "所有豁免+1"},
	{ID: "rare-11", Name: "火球术魔杖", Type: model.ItemTypeWand, Rarity: model.RarityRare, Value: 25000, Attunement: "需要调音", Charges: 7, MaxCharges: 7, Recharge: "dawn", Description: "消耗1-3充能施放火球术"},
	{ID: "rare-12", Name: "闪电魔杖", Type: model.ItemTypeWand, Rarity: model.RarityRare, Value: 25000, Attunement: "需要调音", Charges: 7, MaxCharges: 7, Recharge: "dawn", Description: "消耗1-2充能施放闪电束"},
	{ID: "rare-13", Name: "治疗药水（上级）", Type: model.ItemTypePotion, Rarity: model.RarityRare, Value: 15000, Consumable: true, Description: "恢复8d4+8 HP"},
	{ID: "rare-14", Name: "隐形药水", Type: model.ItemTypePotion, Rarity: model.RarityRare, Value: 12000, Consumable: true, Description: "饮用后隐形1小时"},
	{ID: "rare-15", Name: "飞行药水", Type: model.ItemTypePotion, Rarity: model.RarityRare, Value: 15000, Consumable: true, Description: "饮用后获得飞行速度60尺1小时"},
	// 极稀有
	{ID: "very-rare-1", Name: "防护披风", Type: model.ItemTypeWondrousItem, Rarity: model.RarityVeryRare, Value: 60000, Attunement: "需要调音", MagicEffects: []string{"ac+2"}, Description: "AC+2"},
	{ID: "very-rare-2", Name: "抵抗披风", Type: model.ItemTypeWondrousItem, Rarity: model.RarityVeryRare, Value: 60000, Attunement: "需要调音", MagicEffects: []string{"all-saves+2"}, Description: "所有豁免+2"},
	{ID: "very-rare-3", Name: "速度之靴", Type: model.ItemTypeWondrousItem, Rarity: model.RarityVeryRare, Value: 70000, Attunement: "需要调音", MagicEffects: []string{"double-speed"}, Description: "速度翻倍"},
	{ID: "very-rare-4", Name: "巨力手套", Type: model.ItemTypeWondrousItem, Rarity: model.RarityVeryRare, Value: 60000, Attunement: "需要调音", MagicEffects: []string{"str-21"}, Description: "力量设为21"},
	{ID: "very-rare-5", Name: "敏捷手套", Type: model.ItemTypeWondrousItem, Rarity: model.RarityVeryRare, Value: 60000, Attunement: "需要调音", MagicEffects: []string{"dex-21"}, Description: "敏捷设为21"},
	{ID: "very-rare-6", Name: "体质手套", Type: model.ItemTypeWondrousItem, Rarity: model.RarityVeryRare, Value: 60000, Attunement: "需要调音", MagicEffects: []string{"con-21"}, Description: "体质设为21"},
	{ID: "very-rare-7", Name: "智力手套", Type: model.ItemTypeWondrousItem, Rarity: model.RarityVeryRare, Value: 60000, Attunement: "需要调音", MagicEffects: []string{"int-21"}, Description: "智力设为21"},
	{ID: "very-rare-8", Name: "感知手套", Type: model.ItemTypeWondrousItem, Rarity: model.RarityVeryRare, Value: 60000, Attunement: "需要调音", MagicEffects: []string{"wis-21"}, Description: "感知设为21"},
	{ID: "very-rare-9", Name: "魅力手套", Type: model.ItemTypeWondrousItem, Rarity: model.RarityVeryRare, Value: 60000, Attunement: "需要调音", MagicEffects: []string{"cha-21"}, Description: "魅力设为21"},
	{ID: "very-rare-10", Name: "+2武器", Type: model.ItemTypeWeapon, Rarity: model.RarityVeryRare, Value: 80000, Attunement: "需要调音", MagicBonus: 2, Description: "+2武器"},
	{ID: "very-rare-11", Name: "治疗药水（最高级）", Type: model.ItemTypePotion, Rarity: model.RarityVeryRare, Value: 50000, Consumable: true, Description: "恢复10d4+20 HP"},
	{ID: "very-rare-12", Name: "复原药水", Type: model.ItemTypePotion, Rarity: model.RarityVeryRare, Value: 40000, Consumable: true, Description: "移除所有负面状态"},
	// 传说
	{ID: "legendary-1", Name: "力量腰带", Type: model.ItemTypeWondrousItem, Rarity: model.RarityLegendary, Value: 200000, Attunement: "需要调音", MagicEffects: []string{"str-24"}, Description: "力量设为24"},
	{ID: "legendary-2", Name: "敏捷腰带", Type: model.ItemTypeWondrousItem, Rarity: model.RarityLegendary, Value: 200000, Attunement: "需要调音", MagicEffects: []string{"dex-24"}, Description: "敏捷设为24"},
	{ID: "legendary-3", Name: "体质腰带", Type: model.ItemTypeWondrousItem, Rarity: model.RarityLegendary, Value: 200000, Attunement: "需要调音", MagicEffects: []string{"con-24"}, Description: "体质设为24"},
	{ID: "legendary-4", Name: "智力腰带", Type: model.ItemTypeWondrousItem, Rarity: model.RarityLegendary, Value: 200000, Attunement: "需要调音", MagicEffects: []string{"int-24"}, Description: "智力设为24"},
	{ID: "legendary-5", Name: "感知腰带", Type: model.ItemTypeWondrousItem, Rarity: model.RarityLegendary, Value: 200000, Attunement: "需要调音", MagicEffects: []string{"wis-24"}, Description: "感知设为24"},
	{ID: "legendary-6", Name: "魅力腰带", Type: model.ItemTypeWondrousItem, Rarity: model.RarityLegendary, Value: 200000, Attunement: "需要调音", MagicEffects: []string{"cha-24"}, Description: "魅力设为24"},
	{ID: "legendary-7", Name: "+3武器", Type: model.ItemTypeWeapon, Rarity: model.RarityLegendary, Value: 250000, Attunement: "需要调音", MagicBonus: 3, Description: "+3武器"},
	{ID: "legendary-8", Name: "神圣复仇者", Type: model.ItemTypeWeapon, Rarity: model.RarityLegendary, Value: 300000, Attunement: "需要调音", MagicBonus: 3, Description: "+3长剑，对邪恶生物额外2d8光耀伤害"},
	{ID: "legendary-9", Name: "世界之柱", Type: model.ItemTypeWeapon, Rarity: model.RarityLegendary, Value: 300000, Attunement: "需要调音", MagicBonus: 3, Description: "+3巨剑，命中附加目标最大HP的1d10黯蚀伤害"},
	{ID: "legendary-10", Name: "全能法杖", Type: model.ItemTypeStaff, Rarity: model.RarityLegendary, Value: 400000, Attunement: "需要调音", Charges: 50, MaxCharges: 50, Recharge: "dawn", Description: "可以施放几乎所有法术"},
}
