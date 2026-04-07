package rules

import (
	"github.com/zwh8800/dnd-core/pkg/data"
	"github.com/zwh8800/dnd-core/pkg/model"
)

// ApplyBackground 应用背景效果到角色
func ApplyBackground(pc *model.PlayerCharacter, backgroundID model.BackgroundID) error {
	bg, exists := data.GlobalRegistry.GetBackground(string(backgroundID))
	if !exists {
		return nil // 背景不存在，静默处理
	}

	// 设置背景 ID
	pc.BackgroundID = string(backgroundID)

	// 应用技能熟练
	if len(bg.SkillProficiencies) > 0 {
		if pc.Proficiencies.ProficientSkills == nil {
			pc.Proficiencies.ProficientSkills = make(map[model.Skill]bool)
		}
		for _, skill := range bg.SkillProficiencies {
			pc.Proficiencies.ProficientSkills[skill] = true
		}
	}

	// 应用工具熟练
	if len(bg.ToolProficiencies) > 0 {
		if pc.Proficiencies.ToolProficiencies == nil {
			pc.Proficiencies.ToolProficiencies = make(map[string]bool)
		}
		for _, tool := range bg.ToolProficiencies {
			pc.Proficiencies.ToolProficiencies[tool] = true
		}
	}

	// 应用语言熟练
	if len(bg.LanguageProficiencies) > 0 {
		if pc.Proficiencies.LanguageProficiencies == nil {
			pc.Proficiencies.LanguageProficiencies = make(map[string]bool)
		}
		for _, lang := range bg.LanguageProficiencies {
			pc.Proficiencies.LanguageProficiencies[lang] = true
		}
	}

	// 应用关联专长
	if bg.AssociatedFeat != "" {
		featInstance := model.FeatInstance{
			FeatID:        bg.AssociatedFeat,
			Source:        model.FeatSourceBackground,
			AcquiredLevel: 1,
		}
		pc.Feats = append(pc.Feats, featInstance)

		// 应用专长效果
		ApplyFeatEffects(pc, bg.AssociatedFeat)
	}

	// 应用背景特性
	if bg.FeatureName != "" {
		pc.Features = append(pc.Features, bg.FeatureName)
	}

	return nil
}

// GetBackgroundFeatures 获取角色背景的特性列表
func GetBackgroundFeatures(pc *model.PlayerCharacter) []string {
	if pc.BackgroundID == "" {
		return nil
	}

	bg, exists := data.GlobalRegistry.GetBackground(pc.BackgroundID)
	if !exists {
		return nil
	}

	features := []string{}
	if bg.FeatureName != "" {
		features = append(features, bg.FeatureName)
	}

	return features
}
