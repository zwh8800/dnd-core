package rules

import (
	"github.com/zwh8800/dnd-core/pkg/model"
)

// ApplyWeaponMastery 在攻击时应用武器掌控效果
func ApplyWeaponMastery(attackResult *AttackResult, masteryType model.WeaponMasteryType, target *model.Actor) {
	if masteryType == "" {
		return
	}

	masteryEffect := model.GetMasteryEffect(masteryType)

	// 根据掌控类型应用效果
	switch masteryType {
	case model.MasterySlow:
		// 命中时降低目标速度（由战斗系统处理）
		if attackResult.Hit {
			attackResult.Effects = append(attackResult.Effects, AttackEffect{
				Type:        "slow",
				Description: masteryEffect.Description,
			})
		}

	case model.MasteryTopple:
		// 命中时尝试击倒目标
		if attackResult.Hit {
			attackResult.Effects = append(attackResult.Effects, AttackEffect{
				Type:        "topple",
				Description: masteryEffect.Description,
			})
		}

	case model.MasteryPush:
		// 命中时推离目标
		if attackResult.Hit {
			attackResult.Effects = append(attackResult.Effects, AttackEffect{
				Type:        "push",
				Description: masteryEffect.Description,
			})
		}

	case model.MasteryNick:
		// 轻武器额外攻击可作为附赠动作（由战斗系统处理）
		// 这里不添加效果，由战斗引擎特殊处理

	case model.MasteryVex:
		// 命中后下次攻击有优势
		if attackResult.Hit {
			attackResult.Effects = append(attackResult.Effects, AttackEffect{
				Type:        "vex",
				Description: masteryEffect.Description,
			})
		}

	case model.MasteryCleave:
		// 击杀目标后可攻击邻近生物
		if attackResult.Hit && attackResult.Killed {
			attackResult.Effects = append(attackResult.Effects, AttackEffect{
				Type:        "cleave",
				Description: masteryEffect.Description,
			})
		}

	case model.MasterySap:
		// 命中时目标下次攻击有劣势
		if attackResult.Hit {
			attackResult.Effects = append(attackResult.Effects, AttackEffect{
				Type:        "sap",
				Description: masteryEffect.Description,
			})
		}

	case model.MasteryGraze:
		// 未命中时仍造成属性修正值伤害
		if !attackResult.Hit {
			attackResult.GrazeDamage = attackResult.AbilityModifier
		}
	}
}
