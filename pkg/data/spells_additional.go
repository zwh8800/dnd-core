package data

import "github.com/zwh8800/dnd-core/pkg/model"

// AdditionalSpells 补充法术数据
var AdditionalSpells = []*model.SpellDefinition{
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
