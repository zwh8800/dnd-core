package data

import "github.com/zwh8800/dnd-core/pkg/model"

// Feats 包含 6 个 SRD 起源专长数据
var Feats = []*model.FeatDefinition{
	{
		ID:   "alert",
		Name: "警觉",
		Type: model.FeatTypeOrigin,
		Effects: model.FeatEffect{
			InitiativeBonus: 5,
			SpecialAbilities: []string{
				"清醒时不会被突袭",
				"其他生物因不可见而获得的攻击优势对你无效",
			},
			Description: "你的先攻获得+5加值。你在清醒时不会被突袭。其他生物因不可见而对你进行攻击检定时不会获得优势。",
		},
		Description: "你时刻保持警惕，在感知和先攻方面获得加成。",
	},
	{
		ID:   "magic-initiate-cleric",
		Name: "魔法入门（牧师）",
		Type: model.FeatTypeOrigin,
		Effects: model.FeatEffect{
			SpecialAbilities: []string{
				"学习 2 个牧师戏法",
				"学习 1 个 1 环牧师法术（每次长休可施放一次）",
			},
			Description: "选择牧师法术列表。从该列表中学习 2 个戏法和 1 个 1 环法术。该 1 环法术每次长休可施放一次。",
		},
		Description: "你研习了特定的魔法传统，并从牧师法术列表中学习了法术。",
	},
	{
		ID:   "magic-initiate-wizard",
		Name: "魔法入门（法师）",
		Type: model.FeatTypeOrigin,
		Effects: model.FeatEffect{
			SpecialAbilities: []string{
				"学习 2 个法师戏法",
				"学习 1 个 1 环法师法术（每次长休可施放一次）",
			},
			Description: "选择法师法术列表。从该列表中学习 2 个戏法和 1 个 1 环法术。该 1 环法术每次长休可施放一次。",
		},
		Description: "你研习了特定的魔法传统，并从法师法术列表中学习了法术。",
	},
	{
		ID:   "savage-attacker",
		Name: "野蛮攻击者",
		Type: model.FeatTypeOrigin,
		Effects: model.FeatEffect{
			SpecialAbilities: []string{
				"每回合一次，可重掷武器伤害骰并选择任意结果使用",
			},
			Description: "每回合一次，当你为近战武器攻击掷伤害骰时，你可以重掷武器的伤害骰，并使用两次结果中的任意一个。",
		},
		Description: "你能以特别残暴的力量进行打击，最大化你的伤害输出。",
	},
	{
		ID:   "skilled",
		Name: "多才多艺",
		Type: model.FeatTypeOrigin,
		Effects: model.FeatEffect{
			SkillProficiencies: []model.Skill{
				model.SkillAthletics,
				model.SkillStealth,
				model.SkillPerception,
			},
			Description: "选择任意三种技能或工具的组合。你获得所选技能或工具的熟练。",
		},
		Description: "你研习了各种学科，并在你选择的三项技能或工具上获得了熟练。",
	},
	{
		ID:   "tough",
		Name: "健壮",
		Type: model.FeatTypeOrigin,
		Effects: model.FeatEffect{
			SpecialAbilities: []string{
				"生命值上限增加 2 点，之后每升一级再增加 2 点",
			},
			Description: "你的生命值上限增加 2 点，并且你每次获得等级时，生命值上限再增加 2 点。",
		},
		Description: "你的生命值上限增加，使你在战斗中更具韧性。",
	},
}
