package data

import "github.com/zwh8800/dnd-core/internal/model"

// GetSkillDefinition 获取技能定义
func GetSkillDefinition(skill model.Skill) (name string, ability model.Ability, description string) {
	switch skill {
	case model.SkillAcrobatics:
		return "Acrobatics", model.AbilityDexterity, "保持平衡、翻滚、灵巧动作"
	case model.SkillAnimalHandling:
		return "Animal Handling", model.AbilityWisdom, "驯服、照顾动物"
	case model.SkillArcana:
		return "Arcana", model.AbilityIntelligence, "魔法、法术、魔法物品知识"
	case model.SkillAthletics:
		return "Athletics", model.AbilityStrength, "攀爬、跳跃、游泳"
	case model.SkillDeception:
		return "Deception", model.AbilityCharisma, "欺骗、误导他人"
	case model.SkillHistory:
		return "History", model.AbilityIntelligence, "历史事件、文明知识"
	case model.SkillInsight:
		return "Insight", model.AbilityWisdom, "洞察他人意图、情感"
	case model.SkillIntimidation:
		return "Intimidation", model.AbilityCharisma, "威胁、恐吓他人"
	case model.SkillInvestigation:
		return "Investigation", model.AbilityIntelligence, "寻找线索、推理分析"
	case model.SkillMedicine:
		return "Medicine", model.AbilityWisdom, "医疗、稳定伤势"
	case model.SkillNature:
		return "Nature", model.AbilityIntelligence, "自然、地形、天气知识"
	case model.SkillPerception:
		return "Perception", model.AbilityWisdom, "察觉周围环境"
	case model.SkillPerformance:
		return "Performance", model.AbilityCharisma, "表演、娱乐他人"
	case model.SkillPersuasion:
		return "Persuasion", model.AbilityCharisma, "说服、影响他人"
	case model.SkillReligion:
		return "Religion", model.AbilityIntelligence, "宗教、神祇知识"
	case model.SkillSleightOfHand:
		return "Sleight of Hand", model.AbilityDexterity, "巧手、偷窃、魔术"
	case model.SkillStealth:
		return "Stealth", model.AbilityDexterity, "潜行、躲藏"
	case model.SkillSurvival:
		return "Survival", model.AbilityWisdom, "野外求生、追踪"
	default:
		return string(skill), "", "未知技能"
	}
}

// GetAllSkills 获取所有技能列表
func GetAllSkills() []model.Skill {
	return model.AllSkills()
}
