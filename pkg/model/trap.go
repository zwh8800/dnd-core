package model

// TrapType 陷阱类型
type TrapType string

const (
	TrapTypeMechanical TrapType = "mechanical" // 机械陷阱
	TrapTypeMagical    TrapType = "magical"    // 魔法陷阱
	TrapTypeNatural    TrapType = "natural"    // 自然陷阱
	TrapTypeHybrid     TrapType = "hybrid"     // 混合陷阱
)

// TrapTriggerType 触发条件类型
type TrapTriggerType string

const (
	TrapTriggerProximity TrapTriggerType = "proximity" // 接近触发
	TrapTriggerPressure  TrapTriggerType = "pressure"  // 压力板触发
	TrapTriggerTripwire  TrapTriggerType = "tripwire"  // 绊线触发
	TrapTriggerVisual    TrapTriggerType = "visual"    // 视觉触发
	TrapTriggerAuditory  TrapTriggerType = "auditory"  // 听觉触发
	TrapTriggerManual    TrapTriggerType = "manual"    // 手动触发
)

// TrapEffectType 陷阱效果类型
type TrapEffectType string

const (
	TrapEffectDamage     TrapEffectType = "damage"     // 伤害
	TrapEffectStatus     TrapEffectType = "status"     // 状态效果
	TrapEffectTransport  TrapEffectType = "transport"  // 传送/位移
	TrapEffectSummon     TrapEffectType = "summon"     // 召唤
	TrapEffectAlteration TrapEffectType = "alteration" // 环境改变
)

// TrapEffect 陷阱效果
type TrapEffect struct {
	Type         TrapEffectType `json:"type"`
	DamageDice   string         `json:"damage_dice,omitempty"`   // 伤害骰，如 "2d6"
	DamageType   string         `json:"damage_type,omitempty"`   // 伤害类型
	SaveDC       int            `json:"save_dc,omitempty"`       // 豁免DC
	SaveAbility  string         `json:"save_ability,omitempty"`  // 豁免属性
	StatusEffect string         `json:"status_effect,omitempty"` // 状态效果ID
	Description  string         `json:"description"`             // 效果描述
}

// TrapDefinition 陷阱定义
type TrapDefinition struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Type        TrapType        `json:"type"`
	Description string          `json:"description"`
	Trigger     TrapTriggerType `json:"trigger"`
	DetectDC    int             `json:"detect_dc"`    // 察觉检测DC
	DisarmDC    int             `json:"disarm_dc"`    // 解除陷阱DC
	DisarmSkill string          `json:"disarm_skill"` // 解除技能（通常是Sleight of Hand）
	Effects     []TrapEffect    `json:"effects"`      // 陷阱效果
	Resettable  bool            `json:"resettable"`   // 是否可重置
	CR          float64         `json:"cr"`           // 挑战等级
	Value       int             `json:"value"`        // 制作价值（铜币）
}

// TrapState 陷阱状态
type TrapState struct {
	Definition   *TrapDefinition `json:"definition"`
	IsArmed      bool            `json:"is_armed"`      // 是否已激活
	HasTriggered bool            `json:"has_triggered"` // 是否已触发
	Remaining    int             `json:"remaining"`     // 剩余触发次数（0=无限）
	Position     string          `json:"position"`      // 位置描述
}
