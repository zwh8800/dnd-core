package rules

import (
	"github.com/zwh8800/dnd-core/pkg/dice"
	"github.com/zwh8800/dnd-core/pkg/model"
)

// CalculateNPCReaction 基于 NPC 倾向和角色表现计算反应 DC
func CalculateNPCReaction(disposition model.NPCDisposition, checkType model.SocialCheckType) int {
	// 基础 DC
	dc := 10

	// 根据倾向调整
	switch disposition {
	case model.DispositionHelpful:
		dc = 5
	case model.DispositionFriendly:
		dc = 10
	case model.DispositionIndifferent:
		dc = 15
	case model.DispositionSuspicious:
		dc = 18
	case model.DispositionHostile:
		dc = 20
	}

	// 根据检定类型微调
	if checkType == model.SocialCheckIntimidation && disposition == model.DispositionHostile {
		dc += 2 // 威吓敌对 NPC 更难
	}
	if checkType == model.SocialCheckDeception && disposition == model.DispositionSuspicious {
		dc += 3 // 欺骗多疑 NPC 更难
	}

	return dc
}

// DetermineAttitudeChange 判定态度变化
func DetermineAttitudeChange(current model.NPCAttitude, success bool, checkType model.SocialCheckType) model.NPCAttitude {
	if success {
		// 检定成功：态度改善
		switch current {
		case model.AttitudeHostile:
			return model.AttitudeIndifferent
		case model.AttitudeIndifferent:
			return model.AttitudeFriendly
		case model.AttitudeFriendly:
			return model.AttitudeFriendly // 已经是最高
		}
	} else {
		// 检定失败：态度恶化
		switch current {
		case model.AttitudeFriendly:
			return model.AttitudeIndifferent
		case model.AttitudeIndifferent:
			return model.AttitudeHostile
		case model.AttitudeHostile:
			return model.AttitudeHostile // 已经是最低
		}
	}

	return current
}

// PerformSocialCheck 执行社交检定
func PerformSocialCheck(abilityScore int, proficiencyBonus int, hasProficiency bool, disposition model.NPCDisposition, checkType model.SocialCheckType) (*model.SocialInteractionResult, error) {
	roller := dice.New(0)

	// 计算 DC
	dc := CalculateNPCReaction(disposition, checkType)

	// 掷骰
	rollResult, err := roller.Roll("1d20")
	if err != nil {
		return nil, err
	}

	// 计算加值
	abilityMod := AbilityModifier(abilityScore)
	profBonus := 0
	if hasProficiency {
		profBonus = proficiencyBonus
	}

	total := rollResult.Total + abilityMod + profBonus
	success := total >= dc

	// 判定态度变化
	attitudeChange := DetermineAttitudeChange(model.AttitudeIndifferent, success, checkType)

	result := &model.SocialInteractionResult{
		Success:        success,
		AttitudeChange: attitudeChange,
		RollTotal:      total,
		DC:             dc,
		CheckType:      checkType,
	}

	if success {
		result.Message = "社交检定成功！NPC 态度改善"
	} else {
		result.Message = "社交检定失败，NPC 态度恶化"
	}

	return result, nil
}
