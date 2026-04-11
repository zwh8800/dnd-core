package data

import (
	"github.com/zwh8800/dnd-core/pkg/model"
)

// init 注册所有法术
func init() {
	for _, spell := range Spells {
		GlobalRegistry.RegisterSpell(&spell.Spell)
	}
}

// Spells 核心法术数据库（50个法术）
var Spells = []*model.SpellDefinition{
	// ========== 戏法 (Cantrips) ==========
	{
		Spell: model.Spell{
			ID:          "fire-bolt",
			Name:        "火焰箭",
			Level:       0,
			School:      model.SpellSchoolEvocation,
			CastTime:    model.SpellCastTime{Value: 1, Unit: "action"},
			Range:       "120尺",
			Components:  []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic},
			Duration:    "立即",
			Description: "你向射程内的一个生物或物体投掷一团火焰。对该目标进行一次远程法术攻击。若命中，目标受到1d10火焰伤害。该法术的伤害在你5级时变为2d10，11级时3d10，17级时4d10。",
			DamageDice:  "1d10",
			DamageType:  model.DamageTypeFire,
			Classes:     []string{"Sorcerer", "Wizard"},
		},
		Effects: []model.SpellEffect{
			{
				Type:        model.SpellEffectDamage,
				TargetType:  model.SpellTargetSingleTarget,
				Range:       "120尺",
				Damage:      &model.SpellDamageEntry{BaseDice: "1d10", DamageType: model.DamageTypeFire, UpcastDicePerLevel: "1d10", UpcastStartLevel: 5},
				Description: "1d10 火焰伤害",
			},
		},
		RequiresAttackRoll: true,
	},
	{
		Spell: model.Spell{
			ID:          "light",
			Name:        "光亮术",
			Level:       0,
			School:      model.SpellSchoolEvocation,
			CastTime:    model.SpellCastTime{Value: 1, Unit: "action"},
			Range:       "接触",
			Components:  []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentMaterial},
			Materials:   "萤火虫或发光苔藓",
			Duration:    "1小时",
			Description: "你接触一个不大于10立方尺的物体。该物体发出20尺半径的明亮光照，并在其后10尺内发出微光光照。",
			Classes:     []string{"Bard", "Cleric", "Sorcerer", "Wizard"},
		},
		Effects: []model.SpellEffect{
			{
				Type:        model.SpellEffectUtility,
				TargetType:  model.SpellTargetTouch,
				Description: "物体发出20尺明亮光照",
			},
		},
	},
	{
		Spell: model.Spell{
			ID:          "mage-hand",
			Name:        "法师之手",
			Level:       0,
			School:      model.SpellSchoolConjuration,
			CastTime:    model.SpellCastTime{Value: 1, Unit: "action"},
			Range:       "30尺",
			Components:  []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic},
			Duration:    "1分钟",
			Description: "一只幽灵般的漂浮手出现在你选择的射程内的一个点上。你可以用这只手操纵物体、进行简单的操作。",
			Classes:     []string{"Bard", "Sorcerer", "Warlock", "Wizard"},
		},
		Effects: []model.SpellEffect{
			{
				Type:        model.SpellEffectUtility,
				TargetType:  model.SpellTargetSingleTarget,
				Range:       "30尺",
				Description: "创造一只可操纵物体的魔法手",
			},
		},
	},
	{
		Spell: model.Spell{
			ID:          "sacred-flame",
			Name:        "圣火术",
			Level:       0,
			School:      model.SpellSchoolEvocation,
			CastTime:    model.SpellCastTime{Value: 1, Unit: "action"},
			Range:       "60尺",
			Components:  []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic},
			Duration:    "立即",
			Description: "火焰般的辐射从目标上方的无形之火中降下。目标必须进行一次敏捷豁免，失败则受到1d8光耀伤害。",
			DamageDice:  "1d8",
			DamageType:  model.DamageTypeRadiant,
			SaveDC:      model.AbilityDexterity,
			Classes:     []string{"Cleric"},
		},
		Effects: []model.SpellEffect{
			{
				Type:        model.SpellEffectDamage,
				TargetType:  model.SpellTargetSingleTarget,
				Range:       "60尺",
				Damage:      &model.SpellDamageEntry{BaseDice: "1d8", DamageType: model.DamageTypeRadiant, UpcastDicePerLevel: "1d8", UpcastStartLevel: 5},
				SaveAbility: model.AbilityDexterity,
				Description: "1d8 光耀伤害，敏捷豁免",
			},
		},
	},
	{
		Spell: model.Spell{
			ID:          "shocking-grasp",
			Name:        "电击掌",
			Level:       0,
			School:      model.SpellSchoolEvocation,
			CastTime:    model.SpellCastTime{Value: 1, Unit: "action"},
			Range:       "接触",
			Components:  []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic},
			Duration:    "立即",
			Description: "闪电从你手中跃出。对目标进行一次近战法术攻击，命中则受到1d8闪电伤害，且该目标本回合无法进行反应动作。",
			DamageDice:  "1d8",
			DamageType:  model.DamageTypeLightning,
			Classes:     []string{"Sorcerer", "Wizard"},
		},
		Effects: []model.SpellEffect{
			{
				Type:        model.SpellEffectDamage,
				TargetType:  model.SpellTargetTouch,
				Damage:      &model.SpellDamageEntry{BaseDice: "1d8", DamageType: model.DamageTypeLightning, UpcastDicePerLevel: "1d8", UpcastStartLevel: 5},
				Description: "1d8 闪电伤害，目标无法进行反应动作",
			},
		},
		RequiresAttackRoll: true,
	},
	{
		Spell: model.Spell{
			ID:          "ray-of-frost",
			Name:        "冰霜射线",
			Level:       0,
			School:      model.SpellSchoolEvocation,
			CastTime:    model.SpellCastTime{Value: 1, Unit: "action"},
			Range:       "60尺",
			Components:  []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic},
			Duration:    "立即",
			Description: "一道冰冷的射线射向目标。进行一次远程法术攻击，命中则受到1d8寒冷伤害，且目标速度减少10尺直到你的下一回合开始。",
			DamageDice:  "1d8",
			DamageType:  model.DamageTypeCold,
			Classes:     []string{"Sorcerer", "Wizard"},
		},
		Effects: []model.SpellEffect{
			{
				Type:        model.SpellEffectDamage,
				TargetType:  model.SpellTargetSingleTarget,
				Range:       "60尺",
				Damage:      &model.SpellDamageEntry{BaseDice: "1d8", DamageType: model.DamageTypeCold, UpcastDicePerLevel: "1d8", UpcastStartLevel: 5},
				Description: "1d8 寒冷伤害，速度-10尺",
			},
		},
		RequiresAttackRoll: true,
	},

	// ========== 1环法术 ==========
	{
		Spell: model.Spell{
			ID:          "magic-missile",
			Name:        "魔法飞弹",
			Level:       1,
			School:      model.SpellSchoolEvocation,
			CastTime:    model.SpellCastTime{Value: 1, Unit: "action"},
			Range:       "120尺",
			Components:  []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic},
			Duration:    "立即",
			Description: "你创造出三根旋转的魔法飞弹。每根飞弹造成1d4+1力场伤害，自动命中。",
			DamageDice:  "3d4+3",
			DamageType:  model.DamageTypeForce,
			Classes:     []string{"Sorcerer", "Wizard"},
		},
		Effects: []model.SpellEffect{
			{
				Type:        model.SpellEffectDamage,
				TargetType:  model.SpellTargetSingleTarget,
				Range:       "120尺",
				Damage:      &model.SpellDamageEntry{BaseDice: "3d4+3", DamageType: model.DamageTypeForce},
				Description: "3根飞弹，每根1d4+1力场伤害，自动命中",
			},
		},
	},
	{
		Spell: model.Spell{
			ID:          "shield",
			Name:        "护盾术",
			Level:       1,
			School:      model.SpellSchoolAbjuration,
			CastTime:    model.SpellCastTime{Value: 1, Unit: "reaction"},
			Range:       "自身",
			Components:  []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic},
			Duration:    "1轮",
			Description: "一道无形的屏障出现，直到你的下一回合开始前AC获得+5加值，且不受魔法飞弹影响。",
			Classes:     []string{"Sorcerer", "Wizard"},
		},
		Effects: []model.SpellEffect{
			{
				Type:        model.SpellEffectBuff,
				TargetType:  model.SpellTargetSelf,
				Description: "AC +5，免疫魔法飞弹，持续1轮",
			},
		},
	},
	{
		Spell: model.Spell{
			ID:          "cure-wounds",
			Name:        "疗伤术",
			Level:       1,
			School:      model.SpellSchoolEvocation,
			CastTime:    model.SpellCastTime{Value: 1, Unit: "action"},
			Range:       "接触",
			Components:  []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic},
			Duration:    "立即",
			Description: "你接触一个生物，恢复1d8+你的施法属性修正值的生命值。",
			HealingDice: "1d8",
			Classes:     []string{"Bard", "Cleric", "Druid", "Paladin", "Ranger"},
		},
		Effects: []model.SpellEffect{
			{
				Type:        model.SpellEffectHealing,
				TargetType:  model.SpellTargetTouch,
				HealingDice: "1d8",
				Description: "恢复1d8+施法属性修正值HP",
			},
		},
	},
	{
		Spell: model.Spell{
			ID:            "bless",
			Name:          "祝福术",
			Level:         1,
			School:        model.SpellSchoolEnchantment,
			CastTime:      model.SpellCastTime{Value: 1, Unit: "action"},
			Range:         "30尺",
			Components:    []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic, model.SpellComponentMaterial},
			Materials:     "一滴圣水",
			Duration:      "专注，最多1分钟",
			Description:   "你祝福射程内至多三个生物。在持续时间内，目标每次进行攻击掷骰或豁免检定时可以额外掷一枚d4并加到结果上。",
			Concentration: true,
			Classes:       []string{"Cleric", "Paladin"},
		},
		Effects: []model.SpellEffect{
			{
				Type:        model.SpellEffectBuff,
				TargetType:  model.SpellTargetSingleTarget,
				Range:       "30尺",
				Description: "攻击和豁免+1d4，专注最多1分钟",
			},
		},
	},
	{
		Spell: model.Spell{
			ID:          "burning-hands",
			Name:        "燃烧之手",
			Level:       1,
			School:      model.SpellSchoolEvocation,
			CastTime:    model.SpellCastTime{Value: 1, Unit: "action"},
			Range:       "自身（15尺锥形）",
			Components:  []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic},
			Duration:    "立即",
			Description: "你的拇指互相摩擦，一道15尺锥形的火焰喷射而出。每个区域内的生物必须进行敏捷豁免，失败受到3d6火焰伤害，成功则减半。",
			DamageDice:  "3d6",
			DamageType:  model.DamageTypeFire,
			SaveDC:      model.AbilityDexterity,
			Classes:     []string{"Sorcerer", "Wizard"},
		},
		Effects: []model.SpellEffect{
			{
				Type:              model.SpellEffectDamage,
				TargetType:        model.SpellTargetCone,
				AreaSize:          15,
				Damage:            &model.SpellDamageEntry{BaseDice: "3d6", DamageType: model.DamageTypeFire},
				SaveAbility:       model.AbilityDexterity,
				SaveSuccessEffect: "half",
				Description:       "3d6 火焰伤害，敏捷豁免成功减半",
			},
		},
	},
	{
		Spell: model.Spell{
			ID:          "sleep",
			Name:        "睡眠术",
			Level:       1,
			School:      model.SpellSchoolEnchantment,
			CastTime:    model.SpellCastTime{Value: 1, Unit: "action"},
			Range:       "90尺",
			Components:  []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic, model.SpellComponentMaterial},
			Materials:   "细沙或玫瑰花瓣",
			Duration:    "1分钟",
			Description: "roll 5d8，结果决定法术影响的HP总和。从HP最低的生物开始，范围内的生物陷入魔法睡眠状态。",
			Classes:     []string{"Bard", "Sorcerer", "Wizard"},
		},
		Effects: []model.SpellEffect{
			{
				Type:             model.SpellEffectCondition,
				TargetType:       model.SpellTargetSphere,
				Range:            "90尺",
				ConditionApplied: model.ConditionUnconscious,
				Description:      "5d8 HP总和的生物陷入魔法睡眠",
			},
		},
	},
	{
		Spell: model.Spell{
			ID:          "thunderwave",
			Name:        "雷鸣波",
			Level:       1,
			School:      model.SpellSchoolEvocation,
			CastTime:    model.SpellCastTime{Value: 1, Unit: "action"},
			Range:       "自身（15尺立方体）",
			Components:  []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic},
			Duration:    "立即",
			Description: "一道15尺立方体的冲击波从你身上爆发。每个生物必须进行体质豁免，失败受到2d8雷鸣伤害并被推离10尺。",
			DamageDice:  "2d8",
			DamageType:  model.DamageTypeThunder,
			SaveDC:      model.AbilityConstitution,
			Classes:     []string{"Bard", "Druid", "Sorcerer", "Wizard"},
		},
		Effects: []model.SpellEffect{
			{
				Type:              model.SpellEffectDamage,
				TargetType:        model.SpellTargetCube,
				AreaSize:          15,
				Damage:            &model.SpellDamageEntry{BaseDice: "2d8", DamageType: model.DamageTypeThunder},
				SaveAbility:       model.AbilityConstitution,
				SaveSuccessEffect: "half",
				Description:       "2d8 雷鸣伤害，体质豁免，失败推离10尺",
			},
		},
	},

	// ========== 2环法术 ==========
	{
		Spell: model.Spell{
			ID:          "scorching-ray",
			Name:        "灼热射线",
			Level:       2,
			School:      model.SpellSchoolEvocation,
			CastTime:    model.SpellCastTime{Value: 1, Unit: "action"},
			Range:       "120尺",
			Components:  []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic},
			Duration:    "立即",
			Description: "你射出三条火焰射线。每条进行一次远程法术攻击，命中造成2d6火焰伤害。",
			DamageDice:  "2d6",
			DamageType:  model.DamageTypeFire,
			Classes:     []string{"Sorcerer", "Wizard"},
		},
		Effects: []model.SpellEffect{
			{
				Type:        model.SpellEffectDamage,
				TargetType:  model.SpellTargetSingleTarget,
				Range:       "120尺",
				Damage:      &model.SpellDamageEntry{BaseDice: "2d6", DamageType: model.DamageTypeFire},
				Description: "3条射线，每条2d6火焰伤害",
			},
		},
		RequiresAttackRoll: true,
	},
	{
		Spell: model.Spell{
			ID:          "misty-step",
			Name:        "迷踪步",
			Level:       2,
			School:      model.SpellSchoolConjuration,
			CastTime:    model.SpellCastTime{Value: 1, Unit: "bonus_action"},
			Range:       "自身",
			Components:  []model.SpellComponent{model.SpellComponentVerbal},
			Duration:    "立即",
			Description: "你短暂地被银雾包围，然后传送到30尺内一个你能看见的位置。",
			Classes:     []string{"Sorcerer", "Warlock", "Wizard"},
		},
		Effects: []model.SpellEffect{
			{
				Type:        model.SpellEffectTeleport,
				TargetType:  model.SpellTargetSelf,
				Range:       "30尺",
				Description: "传送30尺到可见位置",
			},
		},
	},
	{
		Spell: model.Spell{
			ID:            "hold-person",
			Name:          "人类定身术",
			Level:         2,
			School:        model.SpellSchoolEnchantment,
			CastTime:      model.SpellCastTime{Value: 1, Unit: "action"},
			Range:         "60尺",
			Components:    []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic, model.SpellComponentMaterial},
			Materials:     "一小块直铁片",
			Duration:      "专注，最多1分钟",
			Description:   "选择一个射程内的人形生物。目标必须进行感知豁免，失败则被麻痹。",
			Concentration: true,
			SaveDC:        model.AbilityWisdom,
			Classes:       []string{"Bard", "Cleric", "Druid", "Sorcerer", "Warlock", "Wizard"},
		},
		Effects: []model.SpellEffect{
			{
				Type:             model.SpellEffectCondition,
				TargetType:       model.SpellTargetSingleTarget,
				Range:            "60尺",
				SaveAbility:      model.AbilityWisdom,
				ConditionApplied: model.ConditionParalyzed,
				Description:      "目标感知豁免失败则被麻痹，专注最多1分钟",
			},
		},
	},
	{
		Spell: model.Spell{
			ID:            "invisibility",
			Name:          "隐形术",
			Level:         2,
			School:        model.SpellSchoolIllusion,
			CastTime:      model.SpellCastTime{Value: 1, Unit: "action"},
			Range:         "接触",
			Components:    []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic, model.SpellComponentMaterial},
			Materials:     "眼睑粉和涂胶的阿拉伯胶",
			Duration:      "专注，最多1小时",
			Description:   "你接触一个自愿生物，使其隐形直到法术结束。目标攻击或施法后隐形结束。",
			Concentration: true,
			Classes:       []string{"Bard", "Sorcerer", "Warlock", "Wizard"},
		},
		Effects: []model.SpellEffect{
			{
				Type:        model.SpellEffectBuff,
				TargetType:  model.SpellTargetTouch,
				Description: "目标隐形，攻击或施法后结束",
			},
		},
	},
	{
		Spell: model.Spell{
			ID:            "web",
			Name:          "蛛网术",
			Level:         2,
			School:        model.SpellSchoolConjuration,
			CastTime:      model.SpellCastTime{Value: 1, Unit: "action"},
			Range:         "60尺",
			Components:    []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic, model.SpellComponentMaterial},
			Materials:     "一点蜘蛛丝",
			Duration:      "专注，最多1小时",
			Description:   "你在射程内创造一个20尺立方体的蛛网区域。区域内的生物必须进行敏捷豁免，失败则被束缚。",
			Concentration: true,
			SaveDC:        model.AbilityDexterity,
			Classes:       []string{"Sorcerer", "Wizard"},
		},
		Effects: []model.SpellEffect{
			{
				Type:             model.SpellEffectCondition,
				TargetType:       model.SpellTargetCube,
				AreaSize:         20,
				SaveAbility:      model.AbilityDexterity,
				ConditionApplied: model.ConditionRestrained,
				Description:      "20尺立方体蛛网，敏捷豁免失败被束缚",
			},
		},
	},

	// ========== 3环法术 ==========
	{
		Spell: model.Spell{
			ID:             "fireball",
			Name:           "火球术",
			Level:          3,
			School:         model.SpellSchoolEvocation,
			CastTime:       model.SpellCastTime{Value: 1, Unit: "action"},
			Range:          "150尺",
			Components:     []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic, model.SpellComponentMaterial},
			Materials:      "一小块蝙蝠粪便和硫磺",
			Duration:       "立即",
			Description:    "一道明亮的射线射向目标点并爆炸。60尺半径球形区域内的每个生物必须进行敏捷豁免，失败受到8d6火焰伤害，成功则减半。",
			DamageDice:     "8d6",
			DamageType:     model.DamageTypeFire,
			SaveDC:         model.AbilityDexterity,
			AtHigherLevels: "使用4环或更高法术位施放时，法术位每比3环高一环，伤害增加1d6。",
			Classes:        []string{"Sorcerer", "Wizard"},
		},
		Effects: []model.SpellEffect{
			{
				Type:              model.SpellEffectDamage,
				TargetType:        model.SpellTargetSphere,
				Range:             "150尺",
				AreaSize:          60,
				Damage:            &model.SpellDamageEntry{BaseDice: "8d6", DamageType: model.DamageTypeFire, UpcastDicePerLevel: "1d6"},
				SaveAbility:       model.AbilityDexterity,
				SaveSuccessEffect: "half",
				Description:       "60尺半径，8d6火焰伤害，敏捷豁免成功减半，升环+1d6/环",
			},
		},
	},
	{
		Spell: model.Spell{
			ID:          "counterspell",
			Name:        "反制法术",
			Level:       3,
			School:      model.SpellSchoolAbjuration,
			CastTime:    model.SpellCastTime{Value: 1, Unit: "reaction"},
			Range:       "60尺",
			Components:  []model.SpellComponent{model.SpellComponentSomatic},
			Duration:    "立即",
			Description: "你试图打断一个正在施法的生物。如果该法术为3环或更低，自动失败。如果为4环或更高，进行施法属性检定。",
			Classes:     []string{"Sorcerer", "Warlock", "Wizard"},
		},
		Effects: []model.SpellEffect{
			{
				Type:        model.SpellEffectDebuff,
				TargetType:  model.SpellTargetSingleTarget,
				Range:       "60尺",
				Description: "打断施法，3环及以下是自动，4环以上需检定",
			},
		},
	},
	{
		Spell: model.Spell{
			ID:          "revivify",
			Name:        "复生术",
			Level:       3,
			School:      model.SpellSchoolNecromancy,
			CastTime:    model.SpellCastTime{Value: 1, Unit: "action"},
			Range:       "接触",
			Components:  []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic, model.SpellComponentMaterial},
			Materials:   "价值300GP的钻石，施法时消耗",
			Duration:    "立即",
			Description: "你接触一个在3轮内死亡的生物。该生物以1HP复活。",
			Classes:     []string{"Cleric", "Paladin"},
		},
		Effects: []model.SpellEffect{
			{
				Type:        model.SpellEffectHealing,
				TargetType:  model.SpellTargetTouch,
				Description: "复活3轮内死亡的生物，恢复1HP",
			},
		},
	},
	{
		Spell: model.Spell{
			ID:             "lightning-bolt",
			Name:           "闪电束",
			Level:          3,
			School:         model.SpellSchoolEvocation,
			CastTime:       model.SpellCastTime{Value: 1, Unit: "action"},
			Range:          "自身（100尺线形）",
			Components:     []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic, model.SpellComponentMaterial},
			Materials:      "一小块毛皮和一根玻璃棒",
			Duration:       "立即",
			Description:    "一道100尺长、5尺宽的闪电束从你手中射出。每个区域内的生物必须进行敏捷豁免，失败受到8d6闪电伤害。",
			DamageDice:     "8d6",
			DamageType:     model.DamageTypeLightning,
			SaveDC:         model.AbilityDexterity,
			AtHigherLevels: "使用4环或更高法术位施放时，法术位每比3环高一环，伤害增加1d6。",
			Classes:        []string{"Sorcerer", "Wizard"},
		},
		Effects: []model.SpellEffect{
			{
				Type:              model.SpellEffectDamage,
				TargetType:        model.SpellTargetLine,
				AreaSize:          100,
				Damage:            &model.SpellDamageEntry{BaseDice: "8d6", DamageType: model.DamageTypeLightning, UpcastDicePerLevel: "1d6"},
				SaveAbility:       model.AbilityDexterity,
				SaveSuccessEffect: "half",
				Description:       "100尺×5尺，8d6闪电伤害，敏捷豁免成功减半",
			},
		},
	},
	{
		Spell: model.Spell{
			ID:            "fly",
			Name:          "飞行术",
			Level:         3,
			School:        model.SpellSchoolTransmutation,
			CastTime:      model.SpellCastTime{Value: 1, Unit: "action"},
			Range:         "接触",
			Components:    []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic, model.SpellComponentMaterial},
			Materials:     "鸟类的羽毛",
			Duration:      "专注，最多10分钟",
			Description:   "你接触一个自愿生物，获得60尺飞行速度。",
			Concentration: true,
			Classes:       []string{"Sorcerer", "Warlock", "Wizard"},
		},
		Effects: []model.SpellEffect{
			{
				Type:        model.SpellEffectBuff,
				TargetType:  model.SpellTargetTouch,
				Description: "目标获得60尺飞行速度",
			},
		},
	},
	{
		Spell: model.Spell{
			ID:            "haste",
			Name:          "加速术",
			Level:         3,
			School:        model.SpellSchoolTransmutation,
			CastTime:      model.SpellCastTime{Value: 1, Unit: "action"},
			Range:         "30尺",
			Components:    []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic, model.SpellComponentMaterial},
			Materials:     "削成片的甘草根",
			Duration:      "专注，最多1分钟",
			Description:   "选择一个自愿生物。目标速度翻倍，AC+2，敏捷豁免具有优势，每回合额外获得一个动作。",
			Concentration: true,
			Classes:       []string{"Sorcerer", "Wizard"},
		},
		Effects: []model.SpellEffect{
			{
				Type:        model.SpellEffectBuff,
				TargetType:  model.SpellTargetSingleTarget,
				Range:       "30尺",
				Description: "速度翻倍，AC+2，敏捷豁免优势，额外动作",
			},
		},
	},

	// ========== 4-5环法术 ==========
	{
		Spell: model.Spell{
			ID:            "polymorph",
			Name:          "变身术",
			Level:         4,
			School:        model.SpellSchoolTransmutation,
			CastTime:      model.SpellCastTime{Value: 1, Unit: "action"},
			Range:         "60尺",
			Components:    []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic, model.SpellComponentMaterial},
			Materials:     "毛毛虫茧",
			Duration:      "专注，最多1小时",
			Description:   "将一个生物变形为另一种生物。新形态的挑战等级不能超过目标的等级或CR。",
			Concentration: true,
			SaveDC:        model.AbilityWisdom,
			Classes:       []string{"Bard", "Druid", "Sorcerer", "Wizard"},
		},
		Effects: []model.SpellEffect{
			{
				Type:        model.SpellEffectDebuff,
				TargetType:  model.SpellTargetSingleTarget,
				Range:       "60尺",
				SaveAbility: model.AbilityWisdom,
				Description: "目标变形为新形态，专注最多1小时",
			},
		},
	},
	{
		Spell: model.Spell{
			ID:             "cone-of-cold",
			Name:           "寒冰锥",
			Level:          5,
			School:         model.SpellSchoolEvocation,
			CastTime:       model.SpellCastTime{Value: 1, Unit: "action"},
			Range:          "自身（60尺锥形）",
			Components:     []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic, model.SpellComponentMaterial},
			Materials:      "一小块水晶或玻璃锥",
			Duration:       "立即",
			Description:    "从你手中释放60尺锥形的极寒。每个生物必须进行体质豁免，失败受到8d8寒冷伤害，成功则减半。",
			DamageDice:     "8d8",
			DamageType:     model.DamageTypeCold,
			SaveDC:         model.AbilityConstitution,
			AtHigherLevels: "使用6环或更高法术位施放时，法术位每比5环高一环，伤害增加1d8。",
			Classes:        []string{"Sorcerer", "Wizard"},
		},
		Effects: []model.SpellEffect{
			{
				Type:              model.SpellEffectDamage,
				TargetType:        model.SpellTargetCone,
				AreaSize:          60,
				Damage:            &model.SpellDamageEntry{BaseDice: "8d8", DamageType: model.DamageTypeCold, UpcastDicePerLevel: "1d8"},
				SaveAbility:       model.AbilityConstitution,
				SaveSuccessEffect: "half",
				Description:       "60尺锥形，8d8寒冷伤害，体质豁免成功减半",
			},
		},
	},
	{
		Spell: model.Spell{
			ID:            "wall-of-force",
			Name:          "力场墙",
			Level:         5,
			School:        model.SpellSchoolEvocation,
			CastTime:      model.SpellCastTime{Value: 1, Unit: "action"},
			Range:         "120尺",
			Components:    []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic, model.SpellComponentMaterial},
			Materials:     "一小块 powdered glass",
			Duration:      "专注，最多10分钟",
			Description:   "一道无形的力场墙在射程内出现。墙免疫所有伤害，可以创造平面或穹顶形状。",
			Concentration: true,
			Classes:       []string{"Wizard"},
		},
		Effects: []model.SpellEffect{
			{
				Type:        model.SpellEffectUtility,
				TargetType:  model.SpellTargetSingleTarget,
				Range:       "120尺",
				Description: "创造不可见的力场墙，免疫所有伤害",
			},
		},
	},
	{
		Spell: model.Spell{
			ID:            "dominate-person",
			Name:          "支配人类",
			Level:         5,
			School:        model.SpellSchoolEnchantment,
			CastTime:      model.SpellCastTime{Value: 1, Unit: "action"},
			Range:         "60尺",
			Components:    []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic},
			Duration:      "专注，最多1分钟",
			Description:   "你通过魔法控制射程内一个人形生物。目标进行感知豁免，失败则被你魅惑。",
			Concentration: true,
			SaveDC:        model.AbilityWisdom,
			Classes:       []string{"Bard", "Sorcerer", "Wizard"},
		},
		Effects: []model.SpellEffect{
			{
				Type:             model.SpellEffectCondition,
				TargetType:       model.SpellTargetSingleTarget,
				Range:            "60尺",
				SaveAbility:      model.AbilityWisdom,
				ConditionApplied: model.ConditionCharmed,
				Description:      "目标感知豁免失败被魅惑，你可以控制其行动",
			},
		},
	},
	{
		Spell: model.Spell{
			ID:          "greater-restoration",
			Name:        "高等复原术",
			Level:       5,
			School:      model.SpellSchoolAbjuration,
			CastTime:    model.SpellCastTime{Value: 1, Unit: "action"},
			Range:       "接触",
			Components:  []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic},
			Materials:   "价值至少100GP的钻石尘，施法时消耗",
			Duration:    "立即",
			Description: "你解除一个目标生物身上的所有诅咒、石化、魅惑、麻痹或中毒状态。",
			Classes:     []string{"Bard", "Cleric", "Druid"},
		},
		Effects: []model.SpellEffect{
			{
				Type:        model.SpellEffectUtility,
				TargetType:  model.SpellTargetTouch,
				Description: "移除所有诅咒、石化、魅惑、麻痹、中毒",
			},
		},
	},
	// 6环
	{
		Spell:              model.Spell{ID: "disintegrate", Name: "灰飞烟灭", Level: 6, School: model.SpellSchoolTransmutation, CastTime: model.SpellCastTime{Value: 1, Unit: "action"}, Range: "60尺", Components: []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic, model.SpellComponentMaterial}, Duration: "立即", Description: "绿色光线。远程法术攻击命中受到10d6+40力场伤害，0HP则化为灰烬。", Classes: []string{"Sorcerer", "Wizard"}},
		Effects:            []model.SpellEffect{{Type: model.SpellEffectDamage, TargetType: model.SpellTargetSingleTarget, Range: "60尺", Damage: &model.SpellDamageEntry{BaseDice: "10d6+40", DamageType: model.DamageTypeForce}, Description: "10d6+40力场伤害"}},
		RequiresAttackRoll: true,
	},
	{
		Spell:   model.Spell{ID: "circle-of-death", Name: "死亡法阵", Level: 6, School: model.SpellSchoolNecromancy, CastTime: model.SpellCastTime{Value: 1, Unit: "action"}, Range: "150尺", Components: []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic, model.SpellComponentMaterial}, Duration: "立即", Description: "60尺半径黑色能量球。体质豁免失败8d6黯蚀伤害，成功减半。", Classes: []string{"Sorcerer", "Warlock", "Wizard"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectDamage, TargetType: model.SpellTargetSphere, Range: "150尺", AreaSize: 60, Damage: &model.SpellDamageEntry{BaseDice: "8d6", DamageType: model.DamageTypeNecrotic}, SaveDC: 15, SaveAbility: model.AbilityConstitution, SaveSuccessEffect: "half", Description: "8d6黯蚀伤害，体质豁免减半"}},
	},
	{
		Spell:   model.Spell{ID: "heal", Name: "医疗术", Level: 6, School: model.SpellSchoolEvocation, CastTime: model.SpellCastTime{Value: 1, Unit: "action"}, Range: "60尺", Components: []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic}, Duration: "立即", Description: "恢复70HP，移除致盲、魅惑、耳聋、恐慌、中毒。", Classes: []string{"Cleric", "Druid"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectHealing, TargetType: model.SpellTargetSingleTarget, Range: "60尺", HealingDice: "70", Description: "恢复70HP"}},
	},
	{
		Spell:   model.Spell{ID: "harm", Name: "伤害术", Level: 6, School: model.SpellSchoolNecromancy, CastTime: model.SpellCastTime{Value: 1, Unit: "action"}, Range: "60尺", Components: []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic}, Duration: "立即", Description: "体质豁免失败14d6黯蚀伤害，成功减半，HP上限减少。", Classes: []string{"Cleric"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectDamage, TargetType: model.SpellTargetSingleTarget, Range: "60尺", Damage: &model.SpellDamageEntry{BaseDice: "14d6", DamageType: model.DamageTypeNecrotic}, SaveDC: 15, SaveAbility: model.AbilityConstitution, SaveSuccessEffect: "half", Description: "14d6黯蚀伤害，HP上限减少"}},
	},
	{
		Spell:   model.Spell{ID: "globe-of-invulnerability", Name: "防魔法力场", Level: 6, School: model.SpellSchoolAbjuration, CastTime: model.SpellCastTime{Value: 1, Unit: "action"}, Range: "自我", Components: []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic, model.SpellComponentMaterial}, Duration: "专注，最多1分钟", Description: "10尺半径力场，5环或更低法术无法影响内部。", Concentration: true, Classes: []string{"Sorcerer", "Wizard"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectBuff, TargetType: model.SpellTargetEmanation, AreaSize: 10, Description: "5环或更低法术无效"}},
	},
	{
		Spell:   model.Spell{ID: "contingency", Name: "连锁意外术", Level: 6, School: model.SpellSchoolEvocation, CastTime: model.SpellCastTime{Value: 10, Unit: "minutes"}, Range: "自我", Components: []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic, model.SpellComponentMaterial}, Duration: "10天", Description: "预设条件触发时自动施放另一个法术。", Classes: []string{"Wizard"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectUtility, TargetType: model.SpellTargetSelf, Description: "条件触发自动施法"}},
	},
	{
		Spell:   model.Spell{ID: "mass-suggestion", Name: "群体暗示术", Level: 6, School: model.SpellSchoolEnchantment, CastTime: model.SpellCastTime{Value: 1, Unit: "action"}, Range: "60尺", Components: []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentMaterial}, Duration: "24小时", Description: "最多12个生物，感知豁免失败被暗示。", Classes: []string{"Bard", "Sorcerer", "Warlock", "Wizard"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectCondition, TargetType: model.SpellTargetSingleTarget, Range: "60尺", ConditionApplied: model.ConditionCharmed, ConditionDuration: "24小时", SaveDC: 15, SaveAbility: model.AbilityWisdom, Description: "最多12生物感知豁免失败被魅惑"}},
	},
	{
		Spell:   model.Spell{ID: "true-seeing", Name: "真知术", Level: 6, School: model.SpellSchoolDivination, CastTime: model.SpellCastTime{Value: 1, Unit: "action"}, Range: "接触", Components: []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic, model.SpellComponentMaterial}, Duration: "1小时", Description: "120尺真实视觉，看穿黑暗、隐形、幻象。", Classes: []string{"Bard", "Cleric", "Sorcerer", "Wizard"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectUtility, TargetType: model.SpellTargetTouch, Description: "120尺真实视觉"}},
	},
	{
		Spell:   model.Spell{ID: "otilukes-freezing-sphere", Name: "欧提路克的冰冻球", Level: 6, School: model.SpellSchoolEvocation, CastTime: model.SpellCastTime{Value: 1, Unit: "action"}, Range: "300尺", Components: []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic, model.SpellComponentMaterial}, Duration: "立即", Description: "60尺锥状或半径10d6寒冷伤害。", Classes: []string{"Wizard"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectDamage, TargetType: model.SpellTargetCone, Range: "60尺", Damage: &model.SpellDamageEntry{BaseDice: "10d6", DamageType: model.DamageTypeCold}, Description: "10d6寒冷伤害"}},
	},
	{
		Spell:   model.Spell{ID: "sunbeam", Name: "日光束", Level: 6, School: model.SpellSchoolEvocation, CastTime: model.SpellCastTime{Value: 1, Unit: "action"}, Range: "自我", Components: []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic, model.SpellComponentMaterial}, Duration: "专注，最多1分钟", Description: "60尺长5尺宽光束。体质豁免失败6d8光耀伤害，成功减半。", Concentration: true, Classes: []string{"Druid", "Sorcerer", "Wizard"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectDamage, TargetType: model.SpellTargetLine, Range: "60尺", Damage: &model.SpellDamageEntry{BaseDice: "6d8", DamageType: model.DamageTypeRadiant}, SaveDC: 15, SaveAbility: model.AbilityConstitution, SaveSuccessEffect: "half", Description: "6d8光耀伤害，体质豁免减半"}},
	},
	// 7环
	{
		Spell:   model.Spell{ID: "fire-storm", Name: "火焰风暴", Level: 7, School: model.SpellSchoolEvocation, CastTime: model.SpellCastTime{Value: 1, Unit: "action"}, Range: "150尺", Components: []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic}, Duration: "立即", Description: "10个10尺立方体。敏捷豁免失败7d10火焰伤害，成功减半。", Classes: []string{"Cleric", "Druid", "Sorcerer"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectDamage, TargetType: model.SpellTargetCube, Range: "150尺", Damage: &model.SpellDamageEntry{BaseDice: "7d10", DamageType: model.DamageTypeFire}, SaveDC: 15, SaveAbility: model.AbilityDexterity, SaveSuccessEffect: "half", Description: "7d10火焰伤害，敏捷豁免减半"}},
	},
	{
		Spell:   model.Spell{ID: "finger-of-death", Name: "死亡一指", Level: 7, School: model.SpellSchoolNecromancy, CastTime: model.SpellCastTime{Value: 1, Unit: "action"}, Range: "60尺", Components: []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic}, Duration: "立即", Description: "体质豁免失败7d8+30黯蚀伤害，成功减半，死亡变僵尸。", Classes: []string{"Sorcerer", "Warlock", "Wizard"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectDamage, TargetType: model.SpellTargetSingleTarget, Range: "60尺", Damage: &model.SpellDamageEntry{BaseDice: "7d8+30", DamageType: model.DamageTypeNecrotic}, SaveDC: 15, SaveAbility: model.AbilityConstitution, SaveSuccessEffect: "half", Description: "7d8+30黯蚀伤害"}},
	},
	{
		Spell:   model.Spell{ID: "forcecage", Name: "力场牢笼", Level: 7, School: model.SpellSchoolEvocation, CastTime: model.SpellCastTime{Value: 1, Unit: "action"}, Range: "100尺", Components: []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic, model.SpellComponentMaterial}, Duration: "1小时", Description: "15尺牢笼，魅力豁免失败无法传送离开。", Classes: []string{"Bard", "Warlock", "Wizard"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectCondition, TargetType: model.SpellTargetCube, Range: "100尺", Description: "15尺牢笼困住生物"}},
	},
	{
		Spell:   model.Spell{ID: "resurrection", Name: "复生术", Level: 7, School: model.SpellSchoolNecromancy, CastTime: model.SpellCastTime{Value: 1, Unit: "hour"}, Range: "接触", Components: []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic, model.SpellComponentMaterial}, Duration: "立即", Description: "复活死亡不超过100年的生物，恢复所有HP。", Classes: []string{"Bard", "Cleric"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectHealing, TargetType: model.SpellTargetTouch, Description: "复活并恢复所有HP"}},
	},
	{
		Spell:   model.Spell{ID: "teleport", Name: "传送术", Level: 7, School: model.SpellSchoolConjuration, CastTime: model.SpellCastTime{Value: 1, Unit: "action"}, Range: "10尺", Components: []model.SpellComponent{model.SpellComponentVerbal}, Duration: "立即", Description: "最多9个生物立即传送至目的地。", Classes: []string{"Bard", "Sorcerer", "Wizard"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectTeleport, TargetType: model.SpellTargetSingleTarget, Range: "10尺", Description: "传送到目的地"}},
	},
	{
		Spell:   model.Spell{ID: "plane-shift", Name: "位面转移", Level: 7, School: model.SpellSchoolConjuration, CastTime: model.SpellCastTime{Value: 1, Unit: "action"}, Range: "接触", Components: []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic, model.SpellComponentMaterial}, Duration: "立即", Description: "最多9个生物传送到另一个位面。", Classes: []string{"Cleric", "Druid", "Sorcerer", "Warlock", "Wizard"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectTeleport, TargetType: model.SpellTargetTouch, Description: "传送到另一位面"}},
	},
	{
		Spell:   model.Spell{ID: "prismatic-spray", Name: "虹光喷射", Level: 7, School: model.SpellSchoolEvocation, CastTime: model.SpellCastTime{Value: 1, Unit: "action"}, Range: "自我", Components: []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic}, Duration: "立即", Description: "60尺锥状，1d8决定颜色效果。", Classes: []string{"Sorcerer", "Wizard"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectDamage, TargetType: model.SpellTargetCone, Range: "60尺", Description: "60尺锥状虹光"}},
	},
	{
		Spell:   model.Spell{ID: "reverse-gravity", Name: "反转重力", Level: 7, School: model.SpellSchoolTransmutation, CastTime: model.SpellCastTime{Value: 1, Unit: "action"}, Range: "100尺", Components: []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic, model.SpellComponentMaterial}, Duration: "专注，最多1分钟", Description: "50尺半径圆柱内重力反转。", Concentration: true, Classes: []string{"Druid", "Wizard"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectUtility, TargetType: model.SpellTargetSphere, Range: "100尺", AreaSize: 50, Description: "圆柱内重力反转"}},
	},
	// 8环
	{
		Spell:   model.Spell{ID: "meteor-swarm", Name: "流星爆", Level: 8, School: model.SpellSchoolEvocation, CastTime: model.SpellCastTime{Value: 1, Unit: "action"}, Range: "1英里", Components: []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic}, Duration: "立即", Description: "四个40尺球。敏捷豁免失败40d6火焰+20d6钝击，成功减半。", Classes: []string{"Sorcerer", "Wizard"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectDamage, TargetType: model.SpellTargetSphere, Range: "1英里", AreaSize: 40, Damage: &model.SpellDamageEntry{BaseDice: "40d6+20d6", DamageType: model.DamageTypeFire}, SaveDC: 15, SaveAbility: model.AbilityDexterity, SaveSuccessEffect: "half", Description: "40d6火焰+20d6钝击"}},
	},
	{
		Spell:   model.Spell{ID: "power-word-stun", Name: "律令：震晕", Level: 8, School: model.SpellSchoolEnchantment, CastTime: model.SpellCastTime{Value: 1, Unit: "action"}, Range: "60尺", Components: []model.SpellComponent{model.SpellComponentVerbal}, Duration: "可变", Description: "HP<=150自动震晕。", Classes: []string{"Bard", "Sorcerer", "Warlock", "Wizard"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectCondition, TargetType: model.SpellTargetSingleTarget, Range: "60尺", ConditionApplied: model.ConditionStunned, Description: "HP<=150自动震晕"}},
	},
	{
		Spell:   model.Spell{ID: "earthquake", Name: "地震术", Level: 8, School: model.SpellSchoolEvocation, CastTime: model.SpellCastTime{Value: 1, Unit: "action"}, Range: "500尺", Components: []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic, model.SpellComponentMaterial}, Duration: "专注，最多1分钟", Description: "100尺半径。敏捷豁免失败50钝击伤害并击倒。", Concentration: true, Classes: []string{"Cleric", "Druid", "Sorcerer"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectDamage, TargetType: model.SpellTargetSphere, Range: "500尺", AreaSize: 100, Damage: &model.SpellDamageEntry{BaseDice: "50", DamageType: model.DamageTypeBludgeoning}, SaveDC: 15, SaveAbility: model.AbilityDexterity, SaveSuccessEffect: "half", Description: "50钝击伤害，敏捷豁免减半"}},
	},
	{
		Spell:   model.Spell{ID: "incendiary-cloud", Name: "燃烧云", Level: 8, School: model.SpellSchoolConjuration, CastTime: model.SpellCastTime{Value: 1, Unit: "action"}, Range: "150尺", Components: []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic}, Duration: "专注，最多1分钟", Description: "20尺半径云雾。每回合开始10d8火焰伤害。", Concentration: true, Classes: []string{"Sorcerer", "Wizard"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectDamage, TargetType: model.SpellTargetSphere, Range: "150尺", AreaSize: 20, Damage: &model.SpellDamageEntry{BaseDice: "10d8", DamageType: model.DamageTypeFire}, Description: "云雾内每回合10d8火焰伤害"}},
	},
	{
		Spell:   model.Spell{ID: "sunburst", Name: "阳炎爆", Level: 8, School: model.SpellSchoolEvocation, CastTime: model.SpellCastTime{Value: 1, Unit: "action"}, Range: "150尺", Components: []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic}, Duration: "立即", Description: "60尺半径。体质豁免失败12d6光耀伤害，成功减半。", Classes: []string{"Druid", "Sorcerer", "Wizard"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectDamage, TargetType: model.SpellTargetSphere, Range: "150尺", AreaSize: 60, Damage: &model.SpellDamageEntry{BaseDice: "12d6", DamageType: model.DamageTypeRadiant}, SaveDC: 15, SaveAbility: model.AbilityConstitution, SaveSuccessEffect: "half", Description: "12d6光耀伤害，体质豁免减半"}},
	},
	{
		Spell:   model.Spell{ID: "clone", Name: "克隆术", Level: 8, School: model.SpellSchoolNecromancy, CastTime: model.SpellCastTime{Value: 1, Unit: "hour"}, Range: "接触", Components: []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic, model.SpellComponentMaterial}, Duration: "立即", Description: "创造克隆体，死亡时复活。", Classes: []string{"Wizard"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectUtility, TargetType: model.SpellTargetTouch, Description: "死亡时复活到克隆体"}},
	},
	{
		Spell:   model.Spell{ID: "mind-blank", Name: "心灵屏障", Level: 8, School: model.SpellSchoolAbjuration, CastTime: model.SpellCastTime{Value: 1, Unit: "action"}, Range: "接触", Components: []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic}, Duration: "24小时", Description: "免疫心灵伤害、魅惑、恐慌、探知。", Classes: []string{"Bard", "Wizard"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectBuff, TargetType: model.SpellTargetTouch, Description: "免疫心灵伤害和状态"}},
	},
	{
		Spell:   model.Spell{ID: "holy-aura", Name: "圣洁灵光", Level: 8, School: model.SpellSchoolAbjuration, CastTime: model.SpellCastTime{Value: 1, Unit: "action"}, Range: "自我", Components: []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic}, Duration: "专注，最多1分钟", Description: "30尺灵光。友方豁免优势，攻击者可能目盲。", Concentration: true, Classes: []string{"Cleric"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectBuff, TargetType: model.SpellTargetEmanation, AreaSize: 30, Description: "友方豁免优势"}},
	},
	// 9环
	{
		Spell:   model.Spell{ID: "power-word-kill", Name: "律令：死", Level: 9, School: model.SpellSchoolEnchantment, CastTime: model.SpellCastTime{Value: 1, Unit: "action"}, Range: "60尺", Components: []model.SpellComponent{model.SpellComponentVerbal}, Duration: "立即", Description: "HP<=100立即死亡，无豁免。", Classes: []string{"Bard", "Sorcerer", "Warlock", "Wizard"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectDamage, TargetType: model.SpellTargetSingleTarget, Range: "60尺", Description: "HP<=100立即死亡"}},
	},
	{
		Spell:   model.Spell{ID: "wish", Name: "祈愿术", Level: 9, School: model.SpellSchoolConjuration, CastTime: model.SpellCastTime{Value: 1, Unit: "action"}, Range: "自我", Components: []model.SpellComponent{model.SpellComponentVerbal}, Duration: "立即", Description: "复制任何8环或更低法术，或创造几乎任何效果。", Classes: []string{"Sorcerer", "Wizard"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectUtility, TargetType: model.SpellTargetSelf, Description: "万能法术"}},
	},
	{
		Spell:   model.Spell{ID: "time-stop", Name: "时间停止", Level: 9, School: model.SpellSchoolTransmutation, CastTime: model.SpellCastTime{Value: 1, Unit: "action"}, Range: "自我", Components: []model.SpellComponent{model.SpellComponentVerbal}, Duration: "立即", Description: "获得1d4+1个额外回合。", Classes: []string{"Sorcerer", "Wizard"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectBuff, TargetType: model.SpellTargetSelf, Description: "1d4+1额外回合"}},
	},
	{
		Spell:   model.Spell{ID: "mass-heal", Name: "群体医疗术", Level: 9, School: model.SpellSchoolConjuration, CastTime: model.SpellCastTime{Value: 1, Unit: "action"}, Range: "60尺", Components: []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic}, Duration: "立即", Description: "最多6个生物分配300点治疗。", Classes: []string{"Cleric"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectHealing, TargetType: model.SpellTargetSingleTarget, Range: "60尺", HealingDice: "300", Description: "分配300点治疗"}},
	},
	{
		Spell:   model.Spell{ID: "foresight", Name: "预警术", Level: 9, School: model.SpellSchoolDivination, CastTime: model.SpellCastTime{Value: 1, Unit: "minute"}, Range: "接触", Components: []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic, model.SpellComponentMaterial}, Duration: "8小时", Description: "攻击和豁免优势，被攻击劣势。", Classes: []string{"Bard", "Druid", "Warlock", "Wizard"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectBuff, TargetType: model.SpellTargetTouch, Description: "攻击和豁免优势"}},
	},
	{
		Spell:   model.Spell{ID: "shapechange", Name: "变形术（高级）", Level: 9, School: model.SpellSchoolTransmutation, CastTime: model.SpellCastTime{Value: 1, Unit: "action"}, Range: "自我", Components: []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic}, Duration: "专注，最多1小时", Description: "变成CR<=你等级的任何生物。", Concentration: true, Classes: []string{"Druid", "Wizard"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectBuff, TargetType: model.SpellTargetSelf, Description: "变形成任何生物"}},
	},
	{
		Spell:   model.Spell{ID: "true-polymorph", Name: "完全变形术", Level: 9, School: model.SpellSchoolTransmutation, CastTime: model.SpellCastTime{Value: 1, Unit: "action"}, Range: "30尺", Components: []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic, model.SpellComponentMaterial}, Duration: "专注，最多1小时", Description: "将生物或物体变成另一个，专注1小时则永久。", Concentration: true, Classes: []string{"Bard", "Warlock", "Wizard"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectBuff, TargetType: model.SpellTargetSingleTarget, Range: "30尺", Description: "目标变成另一物体"}},
	},
	{
		Spell:   model.Spell{ID: "imprisonment", Name: "监禁术", Level: 9, School: model.SpellSchoolAbjuration, CastTime: model.SpellCastTime{Value: 1, Unit: "minute"}, Range: "30尺", Components: []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic}, Duration: "直到被解除", Description: "永久监禁一个生物。", Classes: []string{"Warlock", "Wizard"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectCondition, TargetType: model.SpellTargetSingleTarget, Range: "30尺", Description: "永久监禁"}},
	},
	{
		Spell:   model.Spell{ID: "astral-projection", Name: "星界投射", Level: 9, School: model.SpellSchoolNecromancy, CastTime: model.SpellCastTime{Value: 1, Unit: "hour"}, Range: "10尺", Components: []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic, model.SpellComponentMaterial}, Duration: "特殊", Description: "最多9个生物投射到星界位面。", Classes: []string{"Cleric", "Warlock", "Wizard"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectTeleport, TargetType: model.SpellTargetSingleTarget, Range: "10尺", Description: "投射到星界"}},
	},
	{
		Spell:   model.Spell{ID: "gate", Name: "异界之门", Level: 9, School: model.SpellSchoolConjuration, CastTime: model.SpellCastTime{Value: 1, Unit: "action"}, Range: "60尺", Components: []model.SpellComponent{model.SpellComponentVerbal, model.SpellComponentSomatic, model.SpellComponentMaterial}, Duration: "专注，最多1分钟", Description: "打开传送到其他位面或召唤强大生物。", Concentration: true, Classes: []string{"Cleric", "Sorcerer", "Wizard"}},
		Effects: []model.SpellEffect{{Type: model.SpellEffectSummon, TargetType: model.SpellTargetSingleTarget, Range: "60尺", Description: "异界之门传送或召唤"}},
	},
}
