package rules

import (
	"github.com/zwh8800/dnd-core/pkg/data"
	"github.com/zwh8800/dnd-core/pkg/model"
)

// FeatBonuses 汇总所有专长的加值
type FeatBonuses struct {
	InitiativeBonus int
	AttackBonus     int
	ACBonus         int
}

// CheckFeatPrerequisites 验证角色是否满足专长的先决条件
func CheckFeatPrerequisites(pc *model.PlayerCharacter, featID string) bool {
	feat, exists := data.GlobalRegistry.GetFeat(featID)
	if !exists {
		return false
	}

	if feat.Prerequisite == nil {
		return true // 无先决条件
	}

	prereq := feat.Prerequisite

	// 检查属性分数要求
	if prereq.MinimumAbilityScores != nil {
		for ability, minScore := range prereq.MinimumAbilityScores {
			if pc.AbilityScores.Get(ability) < minScore {
				return false
			}
		}
	}

	// 检查职业要求
	if prereq.RequiredClass != "" {
		hasClass := false
		for _, classLevel := range pc.Classes {
			if classLevel.Class == prereq.RequiredClass {
				hasClass = true
				break
			}
		}
		if !hasClass {
			return false
		}
	}

	// 检查等级要求
	if prereq.MinimumLevel > 0 && pc.TotalLevel < prereq.MinimumLevel {
		return false
	}

	// 检查前置专长要求
	if prereq.RequiredFeat != "" {
		hasFeat := false
		for _, featInstance := range pc.Feats {
			if featInstance.FeatID == prereq.RequiredFeat {
				hasFeat = true
				break
			}
		}
		if !hasFeat {
			return false
		}
	}

	return true
}

// ApplyFeatEffects 应用专长效果到角色
func ApplyFeatEffects(pc *model.PlayerCharacter, featID string) {
	feat, exists := data.GlobalRegistry.GetFeat(featID)
	if !exists {
		return
	}

	effects := feat.Effects

	// 应用属性加成
	if effects.AbilityScoreIncrease != nil {
		for ability, increase := range effects.AbilityScoreIncrease {
			current := pc.AbilityScores.Get(ability)
			pc.AbilityScores.Set(ability, current+increase)
		}
	}

	// 应用属性最大值提升
	if effects.AbilityScoreMax != nil {
		// SRD 5.2.1 中属性最大值通常为 20，专长可突破此限制
		// 这里暂不实现，留待后续扩展
	}

	// 应用技能熟练
	if len(effects.SkillProficiencies) > 0 {
		if pc.Proficiencies.ProficientSkills == nil {
			pc.Proficiencies.ProficientSkills = make(map[model.Skill]bool)
		}
		for _, skill := range effects.SkillProficiencies {
			pc.Proficiencies.ProficientSkills[skill] = true
		}
	}

	// 应用工具熟练
	if len(effects.ToolProficiencies) > 0 {
		if pc.Proficiencies.ToolProficiencies == nil {
			pc.Proficiencies.ToolProficiencies = make(map[string]bool)
		}
		for _, tool := range effects.ToolProficiencies {
			pc.Proficiencies.ToolProficiencies[tool] = true
		}
	}
}

// GetFeatBonuses 汇总角色所有专长的加值
func GetFeatBonuses(pc *model.PlayerCharacter) FeatBonuses {
	bonuses := FeatBonuses{}

	for _, featInstance := range pc.Feats {
		feat, exists := data.GlobalRegistry.GetFeat(featInstance.FeatID)
		if !exists {
			continue
		}

		bonuses.InitiativeBonus += feat.Effects.InitiativeBonus
		bonuses.AttackBonus += feat.Effects.AttackBonus
		bonuses.ACBonus += feat.Effects.ACBonus
	}

	return bonuses
}

// HasFeatAbility 检查角色是否拥有专长赋予的特殊能力
func HasFeatAbility(pc *model.PlayerCharacter, ability string) bool {
	for _, featInstance := range pc.Feats {
		feat, exists := data.GlobalRegistry.GetFeat(featInstance.FeatID)
		if !exists {
			continue
		}

		for _, specialAbility := range feat.Effects.SpecialAbilities {
			if specialAbility == ability {
				return true
			}
		}
	}

	return false
}
