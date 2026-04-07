package model

// ConditionType 代表状态效果的类型
type ConditionType string

const (
	ConditionBlinded       ConditionType = "blinded"       // 目盲
	ConditionCharmed       ConditionType = "charmed"       // 魅惑
	ConditionDeafened      ConditionType = "deafened"      // 耳聋
	ConditionExhaustion    ConditionType = "exhaustion"    // 力竭
	ConditionFrightened    ConditionType = "frightened"    // 恐慌
	ConditionGrappled      ConditionType = "grappled"      // 擒抱
	ConditionIncapacitated ConditionType = "incapacitated" // 失能
	ConditionInvisible     ConditionType = "invisible"     // 隐形
	ConditionParalyzed     ConditionType = "paralyzed"     // 麻痹
	ConditionPetrified     ConditionType = "petrified"     // 石化
	ConditionPoisoned      ConditionType = "poisoned"      // 中毒
	ConditionProne         ConditionType = "prone"         // 倒地
	ConditionRestrained    ConditionType = "restrained"    // 束缚
	ConditionStunned       ConditionType = "stunned"       // 眩晕
	ConditionUnconscious   ConditionType = "unconscious"   // 昏迷
	ConditionStabilized    ConditionType = "stabilized"    // 稳定（死亡豁免成功后）
)

// ConditionInstance 代表一个状态效果的实例
type ConditionInstance struct {
	Type            ConditionType `json:"type"`
	Source          string        `json:"source"`                 // 来源（法术名、能力名等）
	Duration        string        `json:"duration"`               // 持续时间描述
	RemainingRounds int           `json:"remaining_rounds"`       // 剩余回合数（-1表示永久或直到触发条件）
	SaveDC          int           `json:"save_dc,omitempty"`      // 豁免DC（如果可以豁免结束）
	SaveAbility     Ability       `json:"save_ability,omitempty"` // 豁免属性
}

// ConditionEffect 描述状态效果的具体影响
type ConditionEffect struct {
	Type ConditionType `json:"type"`

	// 攻击相关
	AttackDisadvantage     bool `json:"attack_disadvantage"`      // 攻击检定劣势
	AttackAgainstAdvantage bool `json:"attack_against_advantage"` // 被攻击时对方有优势

	// 检定相关
	AbilityCheckDisadvantage bool `json:"ability_check_disadvantage"` // 属性检定劣势
	SkillCheckDisadvantage   bool `json:"skill_check_disadvantage"`   // 技能检定劣势
	SavingThrowDisadvantage  bool `json:"saving_throw_disadvantage"`  // 豁免检定劣势

	// 移动相关
	SpeedZero bool `json:"speed_zero"` // 速度降为0
	CantMove  bool `json:"cant_move"`  // 无法移动

	// 行动相关
	CantTakeActions   bool `json:"cant_take_actions"`   // 无法执行动作
	CantTakeReactions bool `json:"cant_take_reactions"` // 无法执行反应
	CantSpeak         bool `json:"cant_speak"`          // 无法说话
	CastSpellFail     bool `json:"cast_spell_fail"`     // 施法失败

	// 感官相关
	CantSee  bool `json:"cant_see"`  // 无法看见
	CantHear bool `json:"cant_hear"` // 无法听见

	// 自动失败
	AutoFailStrength bool `json:"auto_fail_strength"` // 力量检定自动失败
	AutoFailDex      bool `json:"auto_fail_dex"`      // 敏捷检定自动失败

	// 其他
	Description string `json:"description"` // 效果描述
}

// GetConditionEffect 获取特定状态的效果描述
func GetConditionEffect(conditionType ConditionType) ConditionEffect {
	switch conditionType {
	case ConditionBlinded:
		return ConditionEffect{
			Type:                   conditionType,
			AttackDisadvantage:     true,
			AttackAgainstAdvantage: true,
			CantSee:                true,
			Description:            "无法看见，攻击检定有劣势，被攻击时对方有优势",
		}
	case ConditionCharmed:
		return ConditionEffect{
			Type:                     conditionType,
			AbilityCheckDisadvantage: false,
			Description:              "无法攻击魅惑来源，魅惑源对你有优势",
		}
	case ConditionDeafened:
		return ConditionEffect{
			Type:        conditionType,
			CantHear:    true,
			Description: "无法听见，依赖听觉的检定自动失败",
		}
	case ConditionFrightened:
		return ConditionEffect{
			Type:                     conditionType,
			AbilityCheckDisadvantage: true,
			AttackDisadvantage:       true,
			Description:              "对恐惧源的攻击检定和属性检定有劣势",
		}
	case ConditionGrappled:
		return ConditionEffect{
			Type:        conditionType,
			SpeedZero:   true,
			Description: "速度变为0",
		}
	case ConditionIncapacitated:
		return ConditionEffect{
			Type:              conditionType,
			CantTakeActions:   true,
			CantTakeReactions: true,
			Description:       "无法执行动作或反应",
		}
	case ConditionInvisible:
		return ConditionEffect{
			Type:                   conditionType,
			AttackDisadvantage:     false,
			AttackAgainstAdvantage: true,
			Description:            "攻击有优势，被攻击时有劣势",
		}
	case ConditionParalyzed:
		return ConditionEffect{
			Type:                   conditionType,
			CantTakeActions:        true,
			CantTakeReactions:      true,
			AttackAgainstAdvantage: true,
			AutoFailStrength:       true,
			AutoFailDex:            true,
			Description:            "无法行动，攻击自动失败，被暴击",
		}
	case ConditionPetrified:
		return ConditionEffect{
			Type:            conditionType,
			CantTakeActions: true,
			CantMove:        true,
			Description:     "变成石头，无法行动",
		}
	case ConditionPoisoned:
		return ConditionEffect{
			Type:                     conditionType,
			AbilityCheckDisadvantage: true,
			AttackDisadvantage:       true,
			Description:              "攻击检定和属性检定有劣势",
		}
	case ConditionProne:
		return ConditionEffect{
			Type:               conditionType,
			AttackDisadvantage: true,
			Description:        "攻击有劣势，近战攻击对方有优势",
		}
	case ConditionRestrained:
		return ConditionEffect{
			Type:                   conditionType,
			AttackDisadvantage:     true,
			AttackAgainstAdvantage: true,
			SpeedZero:              true,
			Description:            "速度为0，攻击有劣势，被攻击有劣势",
		}
	case ConditionStunned:
		return ConditionEffect{
			Type:                   conditionType,
			CantTakeActions:        true,
			CantTakeReactions:      true,
			AttackAgainstAdvantage: true,
			AutoFailStrength:       true,
			AutoFailDex:            true,
			Description:            "无法行动，攻击自动失败",
		}
	case ConditionUnconscious:
		return ConditionEffect{
			Type:                   conditionType,
			CantTakeActions:        true,
			CantTakeReactions:      true,
			AttackAgainstAdvantage: true,
			AutoFailStrength:       true,
			AutoFailDex:            true,
			CantSee:                true,
			Description:            "昏迷，无法行动，攻击自动失败",
		}
	case ConditionStabilized:
		return ConditionEffect{
			Type:        conditionType,
			Description: "不再进行死亡豁免检定",
		}
	default:
		return ConditionEffect{
			Type:        conditionType,
			Description: "未知状态",
		}
	}
}
