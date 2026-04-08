package data

import "github.com/zwh8800/dnd-core/pkg/model"

// AdditionalFeats 补充专长数据（战斗专长和通用专长）
var AdditionalFeats = []*model.FeatDefinition{
	// 战斗专长
	{
		ID:   "charger",
		Name: "冲锋者",
		Type: model.FeatTypeCombat,
		Prerequisite: &model.FeatPrerequisite{
			MinimumAbilityScores: map[model.Ability]int{model.AbilityStrength: 13},
		},
		Effects: model.FeatEffect{
			SpecialAbilities: []string{
				"使用动作Dash后，可用附赠动作进行一次近战武器攻击",
				"命中时目标需进行力量豁免，失败则被击倒",
			},
			Description: "当使用动作进行Dash后，你可以用附赠动作进行一次近战武器攻击。如果命中，目标必须进行力量豁免（DC = 8 + 力量调整 + 熟练加值），失败则被击倒。",
		},
		Description: "你能在冲锋后发动强力的攻击。",
	},
	{
		ID:   "dual-wielder",
		Name: "双武器战斗者",
		Type: model.FeatTypeCombat,
		Effects: model.FeatEffect{
			ACBonus: 1,
			SpecialAbilities: []string{
				"双持武器时AC获得+1加值",
				"可以使用非轻型武器进行双持",
				"可以用一个动作同时拔出或收起两把武器",
			},
			Description: "你掌握双持战斗的技巧。双持武器时AC获得+1加值。你可以使用非轻型单手武器进行双持。",
		},
		Description: "你能同时使用两把武器进行战斗。",
	},
	{
		ID:   "great-weapon-master",
		Name: "巨武器大师",
		Type: model.FeatTypeCombat,
		Prerequisite: &model.FeatPrerequisite{
			MinimumAbilityScores: map[model.Ability]int{model.AbilityStrength: 13},
		},
		Effects: model.FeatEffect{
			SpecialAbilities: []string{
				"杀死生物或暴击后，可用附赠动作进行一次近战武器攻击",
				"可以在攻击前选择-5命中，使伤害+10",
			},
			Description: "当你使用具有Heavy属性的武器杀死一个生物或进行暴击时，你可以用附赠动作进行一次近战武器攻击。此外，你可以在攻击前选择使攻击掷骰有-5减值，如果命中则伤害获得+10加值。",
		},
		Description: "你能使用重型武器造成毁灭性打击。",
	},
	{
		ID:   "sharpshooter",
		Name: "神射手",
		Type: model.FeatTypeCombat,
		Prerequisite: &model.FeatPrerequisite{
			MinimumAbilityScores: map[model.Ability]int{model.AbilityDexterity: 13},
		},
		Effects: model.FeatEffect{
			SpecialAbilities: []string{
				"远程攻击忽略半身掩蔽和四分之三掩蔽",
				"攻击拥有掩护的目标时不会处于劣势",
				"可以在攻击前选择-5命中，使伤害+10",
			},
			Description: "你的远程攻击忽略半身掩蔽和四分之三掩蔽。你可以在攻击前选择使攻击掷骰有-5减值，如果命中则伤害获得+10加值。",
		},
		Description: "你能精准地射击远处的目标。",
	},
	{
		ID:   "sentinel",
		Name: "哨兵",
		Type: model.FeatTypeCombat,
		Effects: model.FeatEffect{
			SpecialAbilities: []string{
				"使用机会攻击命中后，目标速度变为0直到回合结束",
				"即使目标使用Disengage动作，仍可进行机会攻击",
				"当敌人攻击你5尺内的其他生物时，可用反应对该敌人进行近战攻击",
			},
			Description: "当你使用机会攻击命中一个生物时，该生物的速度变为0直到当前回合结束。即使目标生物使用Disengage动作，你仍可对其进行机会攻击。",
		},
		Description: "你能有效地控制战场并保护盟友。",
	},
	{
		ID:   "polearm-master",
		Name: "长柄武器大师",
		Type: model.FeatTypeCombat,
		Prerequisite: &model.FeatPrerequisite{
			MinimumAbilityScores: map[model.Ability]int{model.AbilityStrength: 13},
		},
		Effects: model.FeatEffect{
			SpecialAbilities: []string{
				"使用长柄武器进行攻击动作时，可用附赠动作用武器另一端攻击（1d4钝击）",
				"生物进入你的触及范围时，可用反应对其进行近战攻击",
			},
			Description: "当你进行攻击动作并使用长柄武器时，你可以用附赠动作进行额外的近战攻击，伤害为1d4 + 你的属性调整值（钝击伤害）。",
		},
		Description: "你能用长柄武器的两端进行攻击。",
	},
	{
		ID:   "crossbow-expert",
		Name: "弩专家",
		Type: model.FeatTypeCombat,
		Effects: model.FeatEffect{
			SpecialAbilities: []string{
				"忽略弩武器的装填属性",
				"在敌对生物5尺内进行远程攻击掷骰时不会处于劣势",
				"使用单手武器进行攻击动作时，可用附赠动作用手弩攻击",
			},
			Description: "你忽略弩武器的装填属性。你在敌对生物5尺内进行远程攻击掷骰时不会处于劣势。当你使用攻击动作攻击并持用单手武器时，你可以用附赠动作使用手弩进行攻击。",
		},
		Description: "你能快速装填和射击弩。",
	},
	{
		ID:   "grappler",
		Name: "擒抱者",
		Type: model.FeatTypeCombat,
		Prerequisite: &model.FeatPrerequisite{
			MinimumAbilityScores: map[model.Ability]int{model.AbilityStrength: 13},
		},
		Effects: model.FeatEffect{
			SpecialAbilities: []string{
				"对你擒抱的生物进行攻击时具有优势",
				"你可以使用动作尝试压制被擒抱的生物",
			},
			Description: "你对被你擒抱的生物进行攻击掷骰时具有优势。你可以使用动作尝试压制被擒抱的生物。",
		},
		Description: "你在近战擒抱方面技巧娴熟。",
	},
	// 通用专长
	{
		ID:   "actor",
		Name: "演员",
		Type: model.FeatTypeGeneral,
		Effects: model.FeatEffect{
			AbilityScoreIncrease: map[model.Ability]int{model.AbilityCharisma: 1},
			SpecialAbilities: []string{
				"魅力增加1（最高20）",
				"擅长欺骗和表演检定",
				"模仿他人说话和声音时有优势",
			},
			Description: "魅力增加1（最高20）。当你试图模仿他人说话或声音时，你的欺瞒和表演检定具有优势。",
		},
		Description: "你善于模仿和表演。",
	},
	{
		ID:   "healer",
		Name: "治疗者",
		Type: model.FeatTypeGeneral,
		Effects: model.FeatEffect{
			SpecialAbilities: []string{
				"使用医疗包稳定濒死生物时，该生物恢复1HP",
				"使用动作和治疗包，可使一个生物恢复1d6+4 HP",
			},
			Description: "你使用医疗包稳定一个濒死生物时，该生物恢复1HP。你可以使用一个动作和治疗包中的一个耗材来为一个生物恢复1d6 + 4 HP。",
		},
		Description: "你是有效的治疗者。",
	},
	{
		ID:   "lucky",
		Name: "幸运",
		Type: model.FeatTypeGeneral,
		Effects: model.FeatEffect{
			SpecialAbilities: []string{
				"拥有3点幸运点",
				"可以在掷骰后但在DM宣布结果前使用幸运点",
				"使用幸运点可以重掷一次d20",
				"每次长休后恢复所有未使用的幸运点",
			},
			Description: "你拥有3点幸运点。你可以在掷骰后但在地城主宣布结果前使用。每次长休后恢复所有未使用的幸运点。",
		},
		Description: "你有着非凡的运气。",
	},
	{
		ID:   "observant",
		Name: "观察力敏锐",
		Type: model.FeatTypeGeneral,
		Effects: model.FeatEffect{
			AbilityScoreIncrease: map[model.Ability]int{model.AbilityIntelligence: 1},
			SpecialAbilities: []string{
				"智力增加1（最高20）",
				"可以阅读唇语",
				"被动察觉和被动洞悉+5",
			},
			Description: "智力增加1（最高20）。你可以阅读唇语。你的被动察觉和被动洞悉获得+5加值。",
		},
		Description: "你能迅速注意到细节。",
	},
	{
		ID:   "resilient",
		Name: "坚韧",
		Type: model.FeatTypeGeneral,
		Effects: model.FeatEffect{
			AbilityScoreIncrease: map[model.Ability]int{model.AbilityConstitution: 1},
			SpecialAbilities: []string{
				"选择一项属性，该属性增加1（最高20）",
				"获得该属性的豁免熟练",
			},
			Description: "选择一项属性，该属性增加1（最高20），并获得该属性的豁免检定熟练。",
		},
		Description: "你在某一属性方面特别坚韧。",
	},
	{
		ID:   "war-caster",
		Name: "战斗施法者",
		Type: model.FeatTypeGeneral,
		Prerequisite: &model.FeatPrerequisite{
			Description: "必须能够施放至少一个法术",
		},
		Effects: model.FeatEffect{
			SpecialAbilities: []string{
				"进行专注检定时具有优势",
				"即使双手被占用也可以执行法术的成分",
				"可以用反应施放法术替代机会攻击",
			},
			Description: "你在维持专注的体质豁免中具有优势。即使你的双手被武器或盾牌占用，你也可以执行法术所需的成分。",
		},
		Description: "你能在战斗中有效施法。",
	},
	{
		ID:   "mobile",
		Name: "灵活移动",
		Type: model.FeatTypeGeneral,
		Effects: model.FeatEffect{
			SpecialAbilities: []string{
				"速度增加10尺",
				"使用Dash动作时，本回合困难地形不会消耗额外移动力",
				"进行近战攻击后，不会引发目标的机会攻击",
			},
			Description: "速度增加10尺。当你使用Dash动作时，困难地形在该回合内不会消耗你的额外移动力。当你进行近战攻击后，不会引发该目标的机会攻击。",
		},
		Description: "你在移动和闪避方面异常灵活。",
	},
	{
		ID:   "skill-expert",
		Name: "技能专家",
		Type: model.FeatTypeGeneral,
		Effects: model.FeatEffect{
			AbilityScoreIncrease: map[model.Ability]int{model.AbilityStrength: 1},
			SpecialAbilities: []string{
				"一项属性增加1（最高20）",
				"选择两个熟练项，获得专家熟练（加倍熟练加值）",
			},
			Description: "一项属性增加1（最高20）。选择你熟练的两个技能或工具，你获得专家熟练。",
		},
		Description: "你在某些技能上达到了专家水平。",
	},
	{
		ID:   "heavy-armor-master",
		Name: "重甲大师",
		Type: model.FeatTypeGeneral,
		Prerequisite: &model.FeatPrerequisite{
			Description: "必须穿戴重甲",
		},
		Effects: model.FeatEffect{
			AbilityScoreIncrease: map[model.Ability]int{model.AbilityStrength: 1},
			SpecialAbilities: []string{
				"力量增加1（最高20）",
				"来自非魔法武器的钝击、穿刺、挥砍伤害减少3点",
			},
			Description: "力量增加1（最高20）。当你穿着重甲时，来自非魔法武器的钝击、穿刺、挥砍伤害减少3点。",
		},
		Description: "你能利用重甲来减少伤害。",
	},
	{
		ID:   "lightly-armored",
		Name: "轻甲熟练",
		Type: model.FeatTypeGeneral,
		Effects: model.FeatEffect{
			SpecialAbilities: []string{
				"获得轻甲和盾牌熟练",
			},
			Description: "你获得轻甲和盾牌熟练。",
		},
		Description: "你学会了如何有效使用轻甲。",
	},
	{
		ID:   "moderately-armored",
		Name: "中甲熟练",
		Type: model.FeatTypeGeneral,
		Prerequisite: &model.FeatPrerequisite{
			Description: "必须具有轻甲熟练",
		},
		Effects: model.FeatEffect{
			SpecialAbilities: []string{
				"获得中甲熟练",
			},
			Description: "你获得中甲熟练。",
		},
		Description: "你学会了如何有效使用中甲。",
	},
	{
		ID:   "heavily-armored",
		Name: "重甲熟练",
		Type: model.FeatTypeGeneral,
		Prerequisite: &model.FeatPrerequisite{
			Description: "必须具有中甲熟练",
		},
		Effects: model.FeatEffect{
			SpecialAbilities: []string{
				"获得重甲熟练",
			},
			Description: "你获得重甲熟练。",
		},
		Description: "你学会了如何有效使用重甲。",
	},
	{
		ID:   "weapon-master",
		Name: "武器大师",
		Type: model.FeatTypeGeneral,
		Effects: model.FeatEffect{
			AbilityScoreIncrease: map[model.Ability]int{model.AbilityStrength: 1},
			SpecialAbilities: []string{
				"获得两种武器类型的熟练",
			},
			Description: "选择两种武器类型，你获得这些武器的熟练项。",
		},
		Description: "你精通多种武器。",
	},
	{
		ID:   "ritual-caster",
		Name: "仪式施法者",
		Type: model.FeatTypeGeneral,
		Prerequisite: &model.FeatPrerequisite{
			MinimumAbilityScores: map[model.Ability]int{model.AbilityIntelligence: 13},
		},
		Effects: model.FeatEffect{
			SpecialAbilities: []string{
				"获得一本包含两个1环法术的仪式法术书",
				"可以选择牧师、德鲁伊、法师法术列表",
				"可以施放拥有仪式标签的法术而无需法术位",
			},
			Description: "你获得一本包含两个1环法术的仪式法术书。这些法术必须拥有仪式标签。你可以施放这些法术作为仪式施法。",
		},
		Description: "你学会了仪式施法的艺术。",
	},
	{
		ID:   "mage-slayer",
		Name: "法术杀手",
		Type: model.FeatTypeGeneral,
		Effects: model.FeatEffect{
			SpecialAbilities: []string{
				"对5尺内的生物施法时，可用反应进行近战攻击",
				"对专注法术进行伤害时，目标的专注豁免具有劣势",
				"5尺内生物施法时，你对该生物的下一次攻击具有优势",
			},
			Description: "当一个生物在5尺内施放法术时，你可以用反应对该生物进行近战武器攻击。当你对一个专注法术的生物造成伤害时，该生物的专注豁免具有劣势。",
		},
		Description: "你专门猎杀施法者。",
	},
	{
		ID:   "spell-sniper",
		Name: "法术狙击手",
		Type: model.FeatTypeGeneral,
		Prerequisite: &model.FeatPrerequisite{
			Description: "必须能够施放至少一个法术",
		},
		Effects: model.FeatEffect{
			SpecialAbilities: []string{
				"需要攻击掷骰的法术射程翻倍",
				"忽视掩蔽（除了全掩蔽）",
				"学习一个需要攻击掷骰的戏法",
			},
			Description: "你学习一个需要攻击掷骰的戏法。当你施放需要攻击掷骰的法术时，射程翻倍。你的法术忽视掩蔽（除了全掩蔽）。",
		},
		Description: "你的法术精准而远程。",
	},
	{
		ID:   "elemental-adept",
		Name: "元素专家",
		Type: model.FeatTypeGeneral,
		Prerequisite: &model.FeatPrerequisite{
			Description: "必须能够施放至少一个法术",
		},
		Effects: model.FeatEffect{
			SpecialAbilities: []string{
				"选择一种伤害类型：强酸、寒冷、火焰、闪电、雷鸣",
				"忽略该类型伤害的抗性",
				"该类型伤害骰的1视为2",
			},
			Description: "选择强酸、寒冷、火焰、闪电或雷鸣。你施放的法术忽视目标对该伤害类型的抗性。你掷该类型伤害骰时，1视为2。",
		},
		Description: "你精通一种元素伤害。",
	},
	{
		ID:   " Inspiring-leader",
		Name: "激励领袖",
		Type: model.FeatTypeGeneral,
		Prerequisite: &model.FeatPrerequisite{
			MinimumAbilityScores: map[model.Ability]int{model.AbilityCharisma: 13},
		},
		Effects: model.FeatEffect{
			SpecialAbilities: []string{
				"通过10分钟的激励演讲，使最多6个盟友获得临时HP",
				"临时HP = 你的等级 + 魅力调整值",
			},
			Description: "你可以通过10分钟的激励演讲或布道来鼓舞盟友。选择最多6个友善生物，这些生物获得等于你的等级 + 你的魅力调整值的临时HP。",
		},
		Description: "你能通过言语激励他人。",
	},
	{
		ID:   "chef",
		Name: "厨师",
		Type: model.FeatTypeGeneral,
		Effects: model.FeatEffect{
			AbilityScoreIncrease: map[model.Ability]int{model.AbilityConstitution: 1},
			SpecialAbilities: []string{
				"体质增加1（最高20）",
				"在短休结束时，如果拥有厨师工具熟练，可使最多6个生物各恢复1d8 HP",
				"花费额外时间制作特殊口粮，食用后获得临时HP",
			},
			Description: "体质增加1（最高20）。在短休结束时，如果你拥有厨师工具熟练且拥有足够的口粮，你可以为最多6个生物准备特殊口粮，这些生物各恢复1d8 HP。",
		},
		Description: "你是熟练的厨师。",
	},
	{
		ID:   "crusher",
		Name: "粉碎者",
		Type: model.FeatTypeGeneral,
		Effects: model.FeatEffect{
			AbilityScoreIncrease: map[model.Ability]int{model.AbilityStrength: 1},
			SpecialAbilities: []string{
				"力量增加1（最高20）",
				"钝击伤害的攻击可以将目标移动5尺到相邻位置",
				"对目标进行暴击时，下次攻击具有优势",
			},
			Description: "力量增加1（最高20）。当你使用钝击伤害命中一个生物时，你可以将该生物移动5尺到你触及范围内的一个空间。此外，当你对一个生物进行暴击时，下次对该生物的攻击具有优势。",
		},
		Description: "你能用钝击力量粉碎敌人。",
	},
	{
		ID:   "piercer",
		Name: "穿刺者",
		Type: model.FeatTypeGeneral,
		Effects: model.FeatEffect{
			AbilityScoreIncrease: map[model.Ability]int{model.AbilityStrength: 1},
			SpecialAbilities: []string{
				"力量或敏捷增加1（最高20）",
				"穿刺伤害骰可以重掷一次",
				"暴击时额外一个伤害骰",
			},
			Description: "力量或敏捷增加1（最高20）。当你掷穿刺伤害骰时，可以重掷一个骰子。当你对一个生物进行暴击时，额外掷一个伤害骰。",
		},
		Description: "你能用穿刺攻击造成致命伤害。",
	},
	{
		ID:   "slasher",
		Name: "挥砍者",
		Type: model.FeatTypeGeneral,
		Effects: model.FeatEffect{
			AbilityScoreIncrease: map[model.Ability]int{model.AbilityStrength: 1},
			SpecialAbilities: []string{
				"力量或敏捷增加1（最高20）",
				"挥砍伤害命中后，目标速度减少10尺直到你的下回合开始",
				"暴击时，目标速度减为0",
			},
			Description: "力量或敏捷增加1（最高20）。当你使用挥砍伤害命中一个生物时，该生物的速度减少10尺直到你的下回合开始。当你进行暴击时，该生物速度减为0。",
		},
		Description: "你能用挥砍攻击减缓敌人。",
	},
}
