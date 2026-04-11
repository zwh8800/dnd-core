package data

import "github.com/zwh8800/dnd-core/pkg/model"

func init() {
	for _, bg := range Backgrounds {
		GlobalRegistry.RegisterBackground(bg)
	}
}

// Backgrounds 包含 4 个 SRD 背景数据
var Backgrounds = []*model.BackgroundDefinition{
	{
		ID:          model.BackgroundAcolyte,
		Name:        "侍僧",
		Description: "你一生都在侍奉一座特定神祇或神系的神庙。",
		SkillProficiencies: []model.Skill{
			model.SkillInsight,
			model.SkillReligion,
		},
		LanguageProficiencies: []string{"任选两种语言"},
		AssociatedFeat:        "魔法入门（牧师）",
		StartingEquipmentChoices: []model.BackgroundChoice{
			{OptionA: "圣徽", OptionB: "祈祷轮"},
		},
		StartingGold:       1500, // 15 GP，单位为银币
		FeatureName:        "忠信者的庇护所",
		FeatureDescription: "作为一名侍僧，你能赢得与你共享信仰者的尊敬，并且你可以执行你神祇的宗教仪式。",
		SuggestedCharacteristics: &model.BackgroundCharacteristics{
			PersonalityTraits: []string{
				"我崇拜我信仰中的一位特定英雄，并不断提及此人的事迹和榜样。",
				"我能在最激烈的敌人之间找到共同点，同理他们并始终为和平而努力。",
			},
			Ideals: []string{
				"传统。古老的崇拜和牺牲传统必须被保存和维护。",
				"慈善。我总是尽力帮助那些需要帮助的人，无论个人代价是什么。",
			},
			Bonds: []string{
				"我愿付出生命去找回我信仰中失落已久的古代圣物。",
				"总有一天，我要向那些给我打上异端标签的腐败神庙 hierarchy 复仇。",
			},
			Flaws: []string{
				"我对他人评判严苛，对自己更加严厉。",
				"我对神庙的 hierarchy 太过信任。",
			},
		},
	},
	{
		ID:          model.BackgroundCriminal,
		Name:        "罪犯",
		Description: "你是一名有着违法前科的老练罪犯。",
		SkillProficiencies: []model.Skill{
			model.SkillDeception,
			model.SkillStealth,
		},
		ToolProficiencies: []string{"一种赌具", "盗贼工具"},
		AssociatedFeat:    "警觉",
		StartingEquipmentChoices: []model.BackgroundChoice{
			{OptionA: "撬棍", OptionB: "锤子"},
		},
		StartingEquipment:  []string{"带兜帽的深色衣服"},
		StartingGold:       1500,
		FeatureName:        "罪犯线人",
		FeatureDescription: "你有一个可靠且值得信赖的线人，他充当着你与其他罪犯网络之间的联络人。",
		SuggestedCharacteristics: &model.BackgroundCharacteristics{
			PersonalityTraits: []string{
				"当事情出错时，我总是有应对计划。",
				"无论在什么情况下，我总是保持冷静。",
			},
			Ideals: []string{
				"荣誉。我不偷同行的人。",
				"自由。锁链注定要被打破，那些锻造锁链的人也一样。",
			},
			Bonds: []string{
				"我正努力偿还欠一位慷慨恩人的旧债。",
				"我的不义之财用于养活我的家人。",
			},
			Flaws: []string{
				"我忍不住要欺骗那些比我强大的人。",
				"我太贪婪了，这对我不利。",
			},
		},
	},
	{
		ID:          model.BackgroundSage,
		Name:        "学者",
		Description: "你花了数年时间通过学习和研究来了解多元宇宙的知识。",
		SkillProficiencies: []model.Skill{
			model.SkillArcana,
			model.SkillHistory,
		},
		LanguageProficiencies: []string{"任选两种语言"},
		AssociatedFeat:        "魔法入门（法师）",
		StartingEquipmentChoices: []model.BackgroundChoice{
			{OptionA: "一瓶墨水", OptionB: "一支羽毛笔"},
		},
		StartingEquipment:  []string{"小刀", "普通衣服", "已故同事的来信"},
		StartingGold:       1000,
		FeatureName:        "研究者",
		FeatureDescription: "当你试图了解或回忆某条知识时，如果你不知道那条信息，你通常知道从哪里以及从谁那里可以获得它。",
		SuggestedCharacteristics: &model.BackgroundCharacteristics{
			PersonalityTraits: []string{
				"我使用多音节词来给人一种博学多才的印象。",
				"我读过世界上最伟大图书馆里的每一本书。",
			},
			Ideals: []string{
				"知识。通往权力和自我提升的道路是知识。",
				"美丽。美丽的事物指引我们超越自身，走向真理。",
			},
			Bonds: []string{
				"保护学生是我的职责。",
				"我有一本包含可怕秘密的古籍，绝不能落入坏人之手。",
			},
			Flaws: []string{
				"我很容易被知识的承诺分散注意力。",
				"我说话不经大脑，总是冒犯那些跟不上我思路的人。",
			},
		},
	},
	{
		ID:          model.BackgroundSoldier,
		Name:        "士兵",
		Description: "战争是你大部分人生的主要职业。",
		SkillProficiencies: []model.Skill{
			model.SkillAthletics,
			model.SkillIntimidation,
		},
		ToolProficiencies: []string{"一种赌具", "载具（陆地）"},
		AssociatedFeat:    "野蛮攻击者",
		StartingEquipmentChoices: []model.BackgroundChoice{
			{OptionA: "军衔徽章", OptionB: "阵亡战友的纪念品"},
		},
		StartingEquipment:  []string{"骨制骰子", "一副牌", "普通衣服"},
		StartingGold:       1000,
		FeatureName:        "军衔",
		FeatureDescription: "你拥有作为士兵生涯的军衔。忠于你前军事组织的士兵仍然认可你的权威和影响力。",
		SuggestedCharacteristics: &model.BackgroundCharacteristics{
			PersonalityTraits: []string{
				"我总是礼貌而恭敬。",
				"我失去了太多朋友，结交新朋友很慢。",
			},
			Ideals: []string{
				"大义。我们的使命是为保卫他人而献身。",
				"责任。我做我必须做的事，并服从公正的权威。",
			},
			Bonds: []string{
				"我仍然愿意与我一起服役的人同生共死。",
				"我的荣誉就是我的生命。",
			},
			Flaws: []string{
				"即使法律造成痛苦，我也会遵守法律。",
				"我宁可吞下盔甲也不愿承认自己错了。",
			},
		},
	},
}
