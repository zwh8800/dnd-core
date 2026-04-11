package data

import "github.com/zwh8800/dnd-core/pkg/model"

func init() {
	for _, trap := range TrapDataList {
		GlobalRegistry.RegisterTrap(&trap)
	}
}

// TrapData 陷阱数据
type TrapData = model.TrapDefinition

// TrapDataList 陷阱数据列表
var TrapDataList = []TrapData{
	{
		ID:          "poison-needle",
		Name:        "毒针陷阱",
		Type:        model.TrapTypeMechanical,
		Description: "一个隐藏的毒针从锁孔或把手中弹出",
		Trigger:     model.TrapTriggerManual,
		DetectDC:    15,
		DisarmDC:    15,
		DisarmSkill: "sleight-of-hand",
		Effects: []model.TrapEffect{
			{
				Type:        model.TrapEffectDamage,
				DamageDice:  "1d4",
				DamageType:  "poison",
				SaveDC:      15,
				SaveAbility: "con",
				Description: "受到1d4穿刺伤害，必须进行DC 15的体质豁免，失败则中毒1小时",
			},
		},
		Resettable: true,
		CR:         0.25,
		Value:      25000,
	},
	{
		ID:          "falling-rock",
		Name:        "落石陷阱",
		Type:        model.TrapTypeMechanical,
		Description: "踩到压力板后，天花板上的岩石坠落",
		Trigger:     model.TrapTriggerPressure,
		DetectDC:    15,
		DisarmDC:    20,
		DisarmSkill: "thieves-tools",
		Effects: []model.TrapEffect{
			{
				Type:        model.TrapEffectDamage,
				DamageDice:  "3d6",
				DamageType:  "bludgeoning",
				SaveDC:      15,
				SaveAbility: "dex",
				Description: "受到3d6钝击伤害，DC 15敏捷豁免减半",
			},
		},
		Resettable: false,
		CR:         2,
		Value:      50000,
	},
	{
		ID:          "pit-trap",
		Name:        "陷坑陷阱",
		Type:        model.TrapTypeMechanical,
		Description: "伪装的坑洞，底部可能有尖刺",
		Trigger:     model.TrapTriggerPressure,
		DetectDC:    15,
		DisarmDC:    0,
		DisarmSkill: "",
		Effects: []model.TrapEffect{
			{
				Type:        model.TrapEffectDamage,
				DamageDice:  "1d6",
				DamageType:  "falling",
				SaveDC:      15,
				SaveAbility: "dex",
				Description: "坠落10尺受到1d6钝击伤害，DC 15敏捷豁免避免坠落",
			},
		},
		Resettable: false,
		CR:         1,
		Value:      10000,
	},
	{
		ID:          "fire-dart",
		Name:        "火焰镖陷阱",
		Type:        model.TrapTypeMagical,
		Description: "魔法火焰镖射向触发者",
		Trigger:     model.TrapTriggerProximity,
		DetectDC:    18,
		DisarmDC:    18,
		DisarmSkill: "arcana",
		Effects: []model.TrapEffect{
			{
				Type:        model.TrapEffectDamage,
				DamageDice:  "2d6",
				DamageType:  "fire",
				SaveDC:      15,
				SaveAbility: "dex",
				Description: "受到2d6火焰伤害，DC 15敏捷豁免减半",
			},
		},
		Resettable: true,
		CR:         3,
		Value:      75000,
	},
	{
		ID:          "gas-vent",
		Name:        "毒气喷射陷阱",
		Type:        model.TrapTypeMechanical,
		Description: "从墙壁或地板喷出有毒气体",
		Trigger:     model.TrapTriggerPressure,
		DetectDC:    15,
		DisarmDC:    15,
		DisarmSkill: "thieves-tools",
		Effects: []model.TrapEffect{
			{
				Type:         model.TrapEffectStatus,
				SaveDC:       15,
				SaveAbility:  "con",
				StatusEffect: "poisoned",
				Description:  "10尺锥状区域内所有生物必须进行DC 15体质豁免，失败则中毒1分钟",
			},
		},
		Resettable: true,
		CR:         2,
		Value:      40000,
	},
	{
		ID:          "lightning-wire",
		Name:        "闪电丝线陷阱",
		Type:        model.TrapTypeMagical,
		Description: "细丝上带有魔法闪电能量",
		Trigger:     model.TrapTriggerTripwire,
		DetectDC:    18,
		DisarmDC:    18,
		DisarmSkill: "thieves-tools",
		Effects: []model.TrapEffect{
			{
				Type:        model.TrapEffectDamage,
				DamageDice:  "4d6",
				DamageType:  "lightning",
				SaveDC:      18,
				SaveAbility: "dex",
				Description: "受到4d6闪电伤害，DC 18敏捷豁免减半",
			},
		},
		Resettable: true,
		CR:         5,
		Value:      120000,
	},
	{
		ID:          "symbol-pain",
		Name:        "痛苦徽记陷阱",
		Type:        model.TrapTypeMagical,
		Description: "刻有痛苦徽记的魔法符号",
		Trigger:     model.TrapTriggerVisual,
		DetectDC:    18,
		DisarmDC:    0,
		DisarmSkill: "",
		Effects: []model.TrapEffect{
			{
				Type:         model.TrapEffectStatus,
				SaveDC:       18,
				SaveAbility:  "con",
				StatusEffect: "incapacitated",
				Description:  "60尺内所有能看到徽记的生物必须进行DC 18体质豁免，失败则因剧痛失去能力1小时",
			},
		},
		Resettable: false,
		CR:         6,
		Value:      150000,
	},
	{
		ID:          "scythe-blade",
		Name:        "镰刀摆荡陷阱",
		Type:        model.TrapTypeMechanical,
		Description: "巨大的镰刀从墙壁摆荡而出",
		Trigger:     model.TrapTriggerPressure,
		DetectDC:    15,
		DisarmDC:    20,
		DisarmSkill: "thieves-tools",
		Effects: []model.TrapEffect{
			{
				Type:        model.TrapEffectDamage,
				DamageDice:  "4d6",
				DamageType:  "slashing",
				SaveDC:      15,
				SaveAbility: "dex",
				Description: "受到4d6挥砍伤害，DC 15敏捷豁免减半",
			},
		},
		Resettable: true,
		CR:         4,
		Value:      90000,
	},
	{
		ID:          "sleep-gas",
		Name:        "催眠气体陷阱",
		Type:        model.TrapTypeMagical,
		Description: "释放魔法催眠气体",
		Trigger:     model.TrapTriggerProximity,
		DetectDC:    15,
		DisarmDC:    15,
		DisarmSkill: "thieves-tools",
		Effects: []model.TrapEffect{
			{
				Type:         model.TrapEffectStatus,
				SaveDC:       15,
				SaveAbility:  "wis",
				StatusEffect: "unconscious",
				Description:  "10尺半径球体内所有生物必须进行DC 15感知豁免，失败则陷入昏迷10分钟",
			},
		},
		Resettable: true,
		CR:         3,
		Value:      60000,
	},
	{
		ID:          "collapsing-ceiling",
		Name:        "天花板坍塌陷阱",
		Type:        model.TrapTypeMechanical,
		Description: "整个天花板坍塌下来",
		Trigger:     model.TrapTriggerPressure,
		DetectDC:    20,
		DisarmDC:    0,
		DisarmSkill: "",
		Effects: []model.TrapEffect{
			{
				Type:        model.TrapEffectDamage,
				DamageDice:  "8d6",
				DamageType:  "bludgeoning",
				SaveDC:      20,
				SaveAbility: "dex",
				Description: "受到8d6钝击伤害，DC 20敏捷豁免减半",
			},
		},
		Resettable: false,
		CR:         7,
		Value:      200000,
	},
}

// GetTrapData 获取陷阱数据
func GetTrapData(trapID string) *TrapData {
	for _, trap := range TrapDataList {
		if trap.ID == trapID {
			return &trap
		}
	}
	return nil
}

// ListTrapData 列出所有陷阱数据
func ListTrapData() []TrapData {
	return TrapDataList
}
