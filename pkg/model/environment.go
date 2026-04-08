package model

// EnvironmentType 环境类型
type EnvironmentType string

const (
	EnvExtremeCold  EnvironmentType = "extreme_cold"  // 极寒
	EnvExtremeHeat  EnvironmentType = "extreme_heat"  // 极热
	EnvHighAltitude EnvironmentType = "high_altitude" // 高海拔
	EnvDeepWater    EnvironmentType = "deep_water"    // 深水
	EnvUnderwater   EnvironmentType = "underwater"    // 水下
	EnvSmoke        EnvironmentType = "smoke"         // 烟雾
	EnvBrightLight  EnvironmentType = "bright_light"  // 强光
	EnvDarkness     EnvironmentType = "darkness"      // 黑暗
)

// EnvironmentalEffect 环境效果
type EnvironmentalEffect struct {
	Type              EnvironmentType `json:"type"`
	Description       string          `json:"description"`
	ExposureTime      string          `json:"exposure_time"`           // 暴露时间间隔
	SaveDC            int             `json:"save_dc,omitempty"`       // 豁免DC
	SaveAbility       string          `json:"save_ability,omitempty"`  // 豁免属性
	DamageDice        string          `json:"damage_dice,omitempty"`   // 伤害骰
	DamageType        string          `json:"damage_type,omitempty"`   // 伤害类型
	StatusEffect      string          `json:"status_effect,omitempty"` // 状态效果
	EffectDescription string          `json:"effect_description"`      // 效果描述
}

// EnvironmentState 环境状态
type EnvironmentState struct {
	CurrentEnvironment EnvironmentType       `json:"current_environment"`
	Effects            []EnvironmentalEffect `json:"effects"`
	IsActive           bool                  `json:"is_active"`
}
