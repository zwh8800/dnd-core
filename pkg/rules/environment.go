package rules

import (
	"github.com/zwh8800/dnd-core/pkg/model"
)

// GetEnvironmentEffect 获取环境效果
func GetEnvironmentEffect(envType model.EnvironmentType) model.EnvironmentalEffect {
	switch envType {
	case model.EnvExtremeCold:
		return model.EnvironmentalEffect{
			Type:              model.EnvExtremeCold,
			Description:       "极寒环境",
			ExposureTime:      "10 minutes",
			SaveDC:            10,
			SaveAbility:       "con",
			DamageDice:        "1d6",
			DamageType:        "cold",
			EffectDescription: "暴露在极寒环境中每10分钟必须进行DC 10体质豁免，失败则受到1d6寒冷伤害，DC每次增加1",
		}
	case model.EnvExtremeHeat:
		return model.EnvironmentalEffect{
			Type:              model.EnvExtremeHeat,
			Description:       "极热环境",
			ExposureTime:      "1 hour",
			SaveDC:            10,
			SaveAbility:       "con",
			DamageDice:        "1d6",
			DamageType:        "fire",
			StatusEffect:      "exhaustion",
			EffectDescription: "暴露在极热环境中每小时必须进行DC 10体质豁免，失败则受到1d6火焰伤害并获得1级力竭",
		}
	case model.EnvHighAltitude:
		return model.EnvironmentalEffect{
			Type:              model.EnvHighAltitude,
			Description:       "高海拔环境（10000尺以上）",
			ExposureTime:      "1 minute",
			StatusEffect:      "exhaustion",
			EffectDescription: "在10000尺以上高度旅行，每24小时结束时获得1级力竭，力竭等级等于旅行天数",
		}
	case model.EnvDeepWater:
		return model.EnvironmentalEffect{
			Type:              model.EnvDeepWater,
			Description:       "深水环境（60尺以下）",
			EffectDescription: "在深水下需要憋气，能憋气的分钟数等于10 + 体质修正（最小30秒）",
		}
	case model.EnvUnderwater:
		return model.EnvironmentalEffect{
			Type:              model.EnvUnderwater,
			Description:       "水下环境",
			EffectDescription: "水下攻击有劣势（除非使用穿刺武器），远程攻击自动失手，火焰抗性免疫，寒冷抗性易伤",
		}
	case model.EnvSmoke:
		return model.EnvironmentalEffect{
			Type:              model.EnvSmoke,
			Description:       "烟雾环境",
			EffectDescription: "烟雾提供重度遮蔽，所有依赖视觉的检定有劣势",
		}
	case model.EnvBrightLight:
		return model.EnvironmentalEffect{
			Type:              model.EnvBrightLight,
			Description:       "强光环境",
			EffectDescription: "强光环境下，拥有黑暗视觉的生物攻击有劣势",
		}
	case model.EnvDarkness:
		return model.EnvironmentalEffect{
			Type:              model.EnvDarkness,
			Description:       "黑暗环境",
			EffectDescription: "黑暗环境中，没有黑暗视觉或类似能力的生物视同目盲",
		}
	default:
		return model.EnvironmentalEffect{}
	}
}

// ApplyEnvironmentalEffects 应用环境效果
func ApplyEnvironmentalEffects(envType model.EnvironmentType, exposureMinutes int) string {
	effect := GetEnvironmentEffect(envType)
	if effect.Type == "" {
		return "无环境效果"
	}

	return effect.EffectDescription
}
