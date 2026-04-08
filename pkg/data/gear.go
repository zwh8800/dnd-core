package data

import (
	"github.com/zwh8800/dnd-core/pkg/model"
)

// init 注册所有冒险装备和工具
func init() {
	for _, item := range AdventuringGear {
		GlobalRegistry.RegisterGear(item)
	}
	for _, tool := range Tools {
		GlobalRegistry.RegisterTool(tool)
	}
}

// AdventuringGear 冒险装备数据库（40+ 物品）
var AdventuringGear = []*model.Item{
	// 背包和容器
	{
		ID:          "backpack",
		Name:        "背包",
		Description: "一个可以背在肩上的皮制背包",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      5.0,
		Value:       200, // 2 gp
	},
	{
		ID:          "pouch",
		Name:        "腰包",
		Description: "一个系在腰带上的小袋子",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      1.0,
		Value:       50, // 5 sp
	},
	{
		ID:          "belt-pouch",
		Name:        "钱袋",
		Description: "用于装钱币的小袋子",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      0.0,
		Value:       50, // 5 sp
	},
	{
		ID:          "chest",
		Name:        "木箱",
		Description: "一个用于储存的大型木箱",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      25.0,
		Value:       500, // 5 gp
	},
	{
		ID:          "flask",
		Name:        "烧瓶",
		Description: "一个玻璃容器，容量1品脱",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      0.0,
		Value:       2, // 2 cp
	},
	{
		ID:          "jug",
		Name:        "水壶",
		Description: "一个陶制容器，容量1加仑",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      4.0,
		Value:       20, // 2 sp
	},
	{
		ID:          "pot",
		Name:        "铁锅",
		Description: "一个用于烹饪的金属锅",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      10.0,
		Value:       200, // 2 gp
	},
	{
		ID:          "pouch-leather",
		Name:        "皮袋",
		Description: "一个用于储存小物品的皮袋",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      1.0,
		Value:       50, // 5 sp
	},
	{
		ID:          "sack",
		Name:        "麻袋",
		Description: "一个用于装物品的大麻袋",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      0.5,
		Value:       1, // 1 cp
	},
	{
		ID:          "vial",
		Name:        "玻璃瓶",
		Description: "一个小型玻璃容器",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      0.0,
		Value:       100, // 1 gp
	},
	{
		ID:          "waterskin",
		Name:        "水袋",
		Description: "一个用于装水的皮制水袋（装满时重4磅）",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      1.0, // 空的时候
		Value:       20,  // 2 sp
	},

	// 照明
	{
		ID:          "bullseye-lantern",
		Name:        "提灯",
		Description: "发出60尺锥形明亮光照",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      2.0,
		Value:       1000, // 10 gp
	},
	{
		ID:          "candle",
		Name:        "蜡烛",
		Description: "发出5尺明亮光照，燃烧1小时",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      0.0,
		Value:       1, // 1 cp
	},
	{
		ID:          "hooded-lantern",
		Name:        "罩灯",
		Description: "发出30尺明亮光照，可以遮蔽",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      2.0,
		Value:       500, // 5 gp
	},
	{
		ID:          "lamp",
		Name:        "油灯",
		Description: "发出15尺明亮光照，燃烧6小时",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      1.0,
		Value:       50, // 5 sp
	},
	{
		ID:          "torch",
		Name:        "火把",
		Description: "发出20尺明亮光照，燃烧1小时",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      1.0,
		Value:       1, // 1 cp
	},
	{
		ID:          "oil-flask",
		Name:        "油瓶",
		Description: "一瓶灯油，可燃烧6小时",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      1.0,
		Value:       10, // 1 sp
	},

	// 生存装备
	{
		ID:          "bedroll",
		Name:        "卷铺盖",
		Description: "用于户外睡眠的铺盖卷",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      7.0,
		Value:       100, // 1 gp
	},
	{
		ID:          "blanket",
		Name:        "毯子",
		Description: "一条普通的毯子",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      3.0,
		Value:       50, // 5 sp
	},
	{
		ID:          "rations",
		Name:        "口粮",
		Description: "一天的干粮（包括肉干、饼干、坚果等）",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      2.0,
		Value:       50, // 5 sp
	},
	{
		ID:          "tent",
		Name:        "帐篷",
		Description: "一个可容纳两人的简易帐篷",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      20.0,
		Value:       2000, // 20 gp
	},
	{
		ID:          "fishing-tackle",
		Name:        "钓鱼用具",
		Description: "包含鱼线、鱼钩和浮标的套件",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      4.0,
		Value:       100, // 1 gp
	},
	{
		ID:          "trap-jaw",
		Name:        "捕兽夹",
		Description: "一个铁制捕兽夹",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      25.0,
		Value:       500, // 5 gp
	},

	// 绳索和攀爬
	{
		ID:          "chain",
		Name:        "铁链",
		Description: "10尺长的铁链（DC 20 力量检定可挣脱）",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      10.0,
		Value:       500, // 5 gp
	},
	{
		ID:          "grappling-hook",
		Name:        "抓钩",
		Description: "一个用于攀爬的金属抓钩",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      4.0,
		Value:       200, // 2 gp
	},
	{
		ID:          "rope-hempen",
		Name:        "麻绳",
		Description: "2磅承重的50尺麻绳",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      10.0,
		Value:       100, // 1 gp
	},
	{
		ID:          "rope-silk",
		Name:        "丝绳",
		Description: "可承重4磅的50尺丝绳",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      5.0,
		Value:       1000, // 10 gp
	},

	// 工具和杂物
	{
		ID:          "block-tackle",
		Name:        "滑轮组",
		Description: "一组用于提升重物的滑轮",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      5.0,
		Value:       100, // 1 gp
	},
	{
		ID:          "crowbar",
		Name:        "撬棍",
		Description: "用于撬开物品的铁棍",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      5.0,
		Value:       200, // 2 gp
	},
	{
		ID:          "hammer",
		Name:        "锤子",
		Description: "一把普通的锤子",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      3.0,
		Value:       100, // 1 gp
	},
	{
		ID:          "sledgehammer",
		Name:        "大锤",
		Description: "一把重型锤子",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      10.0,
		Value:       200, // 2 gp
	},
	{
		ID:          "shovel",
		Name:        "铁锹",
		Description: "一把普通的铁锹",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      5.0,
		Value:       200, // 2 gp
	},
	{
		ID:          "pickaxe",
		Name:        "十字镐",
		Description: "用于挖矿和攀岩的工具",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      10.0,
		Value:       200, // 2 gp
	},
	{
		ID:          "whetstone",
		Name:        "磨刀石",
		Description: "用于磨利武器的石头",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      1.0,
		Value:       1, // 1 cp
	},
	{
		ID:          "pliers",
		Name:        "钳子",
		Description: "一把普通的钳子",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      1.0,
		Value:       50, // 5 sp
	},
	{
		ID:          "mirror-steel",
		Name:        "钢镜",
		Description: "一面小型钢镜",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      0.5,
		Value:       500, // 5 gp
	},
	{
		ID:          "magnifying-glass",
		Name:        "放大镜",
		Description: "一个用于检查细节的放大镜",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      0.0,
		Value:       10000, // 100 gp
	},
	{
		ID:          "signal-whistle",
		Name:        "信号哨",
		Description: "一个用于发出信号的哨子",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      0.0,
		Value:       50, // 5 sp
	},
	{
		ID:          "spyglass",
		Name:        "望远镜",
		Description: "一个使远处物体看起来更近的望远镜",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      1.0,
		Value:       100000, // 1000 gp
	},
	{
		ID:          "two-person-saw",
		Name:        "双人锯",
		Description: "一把需要两人操作的大锯",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      15.0,
		Value:       300, // 3 gp
	},

	// 书写和记录
	{
		ID:          "ink",
		Name:        "墨水",
		Description: "一瓶1盎司的墨水",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      0.0,
		Value:       1000, // 10 gp
	},
	{
		ID:          "ink-pen",
		Name:        "羽毛笔",
		Description: "一支用于书写的羽毛笔",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      0.0,
		Value:       2, // 2 cp
	},
	{
		ID:          "paper",
		Name:        "纸",
		Description: "一张空白的纸",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      0.0,
		Value:       200, // 2 gp
	},
	{
		ID:          "parchment",
		Name:        "羊皮纸",
		Description: "一张空白的羊皮纸",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      0.0,
		Value:       10, // 1 sp
	},
	{
		ID:          "sealing-wax",
		Name:        "封蜡",
		Description: "用于密封信件的蜡",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      0.0,
		Value:       50, // 5 sp
	},
	{
		ID:          "signet-ring",
		Name:        "印章戒指",
		Description: "一枚带有家族纹章或个人标记的戒指",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      0.0,
		Value:       500, // 5 gp
	},
	{
		ID:          "soap",
		Name:        "肥皂",
		Description: "一块普通的肥皂",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      0.0,
		Value:       2, // 2 cp
	},

	// 弹药
	{
		ID:          "arrows-20",
		Name:        "箭（20支）",
		Description: "一匣20支箭",
		Type:        model.ItemTypeAmmunition,
		Weight:      1.0,
		Value:       100, // 1 gp
	},
	{
		ID:          "bolts-20",
		Name:        "弩矢（20支）",
		Description: "一匣20支弩矢",
		Type:        model.ItemTypeAmmunition,
		Weight:      1.5,
		Value:       100, // 1 gp
	},
	{
		ID:          "bullets-10",
		Name:        "弹丸（10颗）",
		Description: "10颗投石索用的弹丸",
		Type:        model.ItemTypeAmmunition,
		Weight:      1.5,
		Value:       4, // 4 cp
	},
	{
		ID:          "blowgun-needles-50",
		Name:        "吹箭针（50根）",
		Description: "50根吹箭用的针",
		Type:        model.ItemTypeAmmunition,
		Weight:      1.0,
		Value:       100, // 1 gp
	},
	{
		ID:          "net",
		Name:        "网",
		Description: "一张用于捕捉的网",
		Type:        model.ItemTypeWeapon,
		Weight:      3.0,
		Value:       100, // 1 gp
	},

	// 特殊装备
	{
		ID:          "acid-vial",
		Name:        "强酸瓶",
		Description: "作为动作掷出，造成2d6强酸伤害",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      0.0,
		Value:       2500, // 25 gp
		Consumable:  true,
	},
	{
		ID:          "alchemists-fire-flask",
		Name:        "炼金火焰瓶",
		Description: "作为动作掷出，每回合开始造成1d4火焰伤害",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      1.0,
		Value:       5000, // 50 gp
		Consumable:  true,
	},
	{
		ID:          "antitoxin-vial",
		Name:        "解毒剂",
		Description: "饮用后1小时内对中毒状态免疫",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      0.0,
		Value:       5000, // 50 gp
		Consumable:  true,
	},
	{
		ID:          "ball-bearings",
		Name:        "铁蒺藜",
		Description: "一袋1000个铁球，覆盖10尺见方区域",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      2.0,
		Value:       100, // 1 gp
	},
	{
		ID:          "caltrops",
		Name:        "铁蒺藜",
		Description: "一袋20个铁蒺藜，覆盖5尺见方区域",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      2.0,
		Value:       100, // 1 gp
	},
	{
		ID:          "component-pouch",
		Name:        "法器包",
		Description: "一个小袋子，包含施法所需的所有材料成分",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      2.0,
		Value:       2500, // 25 gp
	},
	{
		ID:          "healers-kit",
		Name:        "医疗包",
		Description: "一个包含绷带和药膏的套件，可进行10次稳定伤势",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      3.0,
		Value:       500, // 5 gp
	},
	{
		ID:          "holy-water-flask",
		Name:        "圣水",
		Description: "对邪魔和不死生物造成2d6光耀伤害",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      1.0,
		Value:       2500, // 25 gp
		Consumable:  true,
	},
	{
		ID:          "hunting-trap",
		Name:        "狩猎陷阱",
		Description: "一个夹住生物的陷阱",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      25.0,
		Value:       500, // 5 gp
	},
	{
		ID:          "poison-basic-vial",
		Name:        "基础毒药",
		Description: "涂抹在武器上，造成1d4毒素伤害",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      0.0,
		Value:       10000, // 100 gp
		Consumable:  true,
	},
	{
		ID:          "tinderbox",
		Name:        "火绒盒",
		Description: "包含打火石、铁片和引火物，用于生火",
		Type:        model.ItemTypeAdventuringGear,
		Weight:      1.0,
		Value:       50, // 5 sp
	},
}

// Tools 工具数据库
var Tools = []*model.Item{
	{
		ID:     "alchemists-supplies",
		Name:   "炼金工具",
		Type:   model.ItemTypeTool,
		Weight: 8.0,
		Value:  5000, // 50 gp
	},
	{
		ID:     "brewers-supplies",
		Name:   "酿酒工具",
		Type:   model.ItemTypeTool,
		Weight: 9.0,
		Value:  2000, // 20 gp
	},
	{
		ID:     "calligraphers-supplies",
		Name:   "书法工具",
		Type:   model.ItemTypeTool,
		Weight: 5.0,
		Value:  1000, // 10 gp
	},
	{
		ID:     "carpenters-tools",
		Name:   "木工工具",
		Type:   model.ItemTypeTool,
		Weight: 6.0,
		Value:  800, // 8 gp
	},
	{
		ID:     "cartographers-tools",
		Name:   "制图工具",
		Type:   model.ItemTypeTool,
		Weight: 6.0,
		Value:  1500, // 15 gp
	},
	{
		ID:     "cobblers-tools",
		Name:   "制鞋工具",
		Type:   model.ItemTypeTool,
		Weight: 5.0,
		Value:  500, // 5 gp
	},
	{
		ID:     "cooks-utensils",
		Name:   "厨师用具",
		Type:   model.ItemTypeTool,
		Weight: 8.0,
		Value:  100, // 1 gp
	},
	{
		ID:     "glassblowers-tools",
		Name:   "玻璃吹制工具",
		Type:   model.ItemTypeTool,
		Weight: 5.0,
		Value:  3000, // 30 gp
	},
	{
		ID:     "jewelers-tools",
		Name:   "珠宝匠工具",
		Type:   model.ItemTypeTool,
		Weight: 2.0,
		Value:  2500, // 25 gp
	},
	{
		ID:     "leatherworkers-tools",
		Name:   "皮匠工具",
		Type:   model.ItemTypeTool,
		Weight: 5.0,
		Value:  500, // 5 gp
	},
	{
		ID:     "masons-tools",
		Name:   "石匠工具",
		Type:   model.ItemTypeTool,
		Weight: 8.0,
		Value:  1000, // 10 gp
	},
	{
		ID:     "painters-supplies",
		Name:   "画家用品",
		Type:   model.ItemTypeTool,
		Weight: 5.0,
		Value:  1000, // 10 gp
	},
	{
		ID:     "potters-tools",
		Name:   "陶匠工具",
		Type:   model.ItemTypeTool,
		Weight: 3.0,
		Value:  1000, // 10 gp
	},
	{
		ID:     "smiths-tools",
		Name:   "铁匠工具",
		Type:   model.ItemTypeTool,
		Weight: 8.0,
		Value:  2000, // 20 gp
	},
	{
		ID:     "tinkers-tools",
		Name:   "修补匠工具",
		Type:   model.ItemTypeTool,
		Weight: 10.0,
		Value:  5000, // 50 gp
	},
	{
		ID:     "weavers-tools",
		Name:   "织工工具",
		Type:   model.ItemTypeTool,
		Weight: 5.0,
		Value:  100, // 1 gp
	},
	{
		ID:     "woodcarvers-tools",
		Name:   "木雕工具",
		Type:   model.ItemTypeTool,
		Weight: 5.0,
		Value:  100, // 1 gp
	},
	{
		ID:     "dice-set",
		Name:   "骰子套装",
		Type:   model.ItemTypeTool,
		Weight: 0.0,
		Value:  10, // 1 sp
	},
	{
		ID:     "playing-card-set",
		Name:   "纸牌套装",
		Type:   model.ItemTypeTool,
		Weight: 0.0,
		Value:  50, // 5 sp
	},
	{
		ID:     "bagpipes",
		Name:   "风笛",
		Type:   model.ItemTypeTool,
		Weight: 6.0,
		Value:  3000, // 30 gp
	},
	{
		ID:     "drum",
		Name:   "鼓",
		Type:   model.ItemTypeTool,
		Weight: 3.0,
		Value:  600, // 6 gp
	},
	{
		ID:     "dulcimer",
		Name:   "扬琴",
		Type:   model.ItemTypeTool,
		Weight: 10.0,
		Value:  2500, // 25 gp
	},
	{
		ID:     "flute",
		Name:   "笛子",
		Type:   model.ItemTypeTool,
		Weight: 1.0,
		Value:  200, // 2 gp
	},
	{
		ID:     "lute",
		Name:   "鲁特琴",
		Type:   model.ItemTypeTool,
		Weight: 2.0,
		Value:  3500, // 35 gp
	},
	{
		ID:     "lyre",
		Name:   "里拉琴",
		Type:   model.ItemTypeTool,
		Weight: 2.0,
		Value:  3000, // 30 gp
	},
	{
		ID:     "horn",
		Name:   "号角",
		Type:   model.ItemTypeTool,
		Weight: 2.0,
		Value:  300, // 3 gp
	},
	{
		ID:     "pan-flute",
		Name:   "排箫",
		Type:   model.ItemTypeTool,
		Weight: 2.0,
		Value:  1200, // 12 gp
	},
	{
		ID:     "shawm",
		Name:   "肖姆管",
		Type:   model.ItemTypeTool,
		Weight: 1.0,
		Value:  200, // 2 gp
	},
	{
		ID:     "viol",
		Name:   "维奥尔琴",
		Type:   model.ItemTypeTool,
		Weight: 1.0,
		Value:  3000, // 30 gp
	},
	{
		ID:     "navigator-tools",
		Name:   "导航工具",
		Type:   model.ItemTypeTool,
		Weight: 2.0,
		Value:  2500, // 25 gp
	},
	{
		ID:     "thieves-tools",
		Name:   "盗贼工具",
		Type:   model.ItemTypeTool,
		Weight: 1.0,
		Value:  2500, // 25 gp
	},
	{
		ID:     "disguise-kit",
		Name:   "易容工具",
		Type:   model.ItemTypeTool,
		Weight: 3.0,
		Value:  2500, // 25 gp
	},
	{
		ID:     "forgery-kit",
		Name:   "伪造工具",
		Type:   model.ItemTypeTool,
		Weight: 5.0,
		Value:  1500, // 15 gp
	},
	{
		ID:     "gaming-set-dice",
		Name:   "赌具（骰子）",
		Type:   model.ItemTypeTool,
		Weight: 0.0,
		Value:  10, // 1 sp
	},
	{
		ID:     "gaming-set-cards",
		Name:   "赌具（纸牌）",
		Type:   model.ItemTypeTool,
		Weight: 0.0,
		Value:  50, // 5 sp
	},
	{
		ID:     "herbalism-kit",
		Name:   "草药包",
		Type:   model.ItemTypeTool,
		Weight: 3.0,
		Value:  500, // 5 gp
	},
	{
		ID:     "poisoners-kit",
		Name:   "制毒工具",
		Type:   model.ItemTypeTool,
		Weight: 5.0,
		Value:  5000, // 50 gp
	},
}
