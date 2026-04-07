package model

// ActionType 代表动作类型
type ActionType string

const (
	// 标准动作
	ActionAttack    ActionType = "attack"     // 攻击
	ActionCastSpell ActionType = "cast_spell" // 施法
	ActionDash      ActionType = "dash"       // 冲刺
	ActionDisengage ActionType = "disengage"  // 脱离
	ActionDodge     ActionType = "dodge"      // 闪避
	ActionHelp      ActionType = "help"       // 协助
	ActionHide      ActionType = "hide"       // 躲藏
	ActionReady     ActionType = "ready"      // 准备
	ActionSearch    ActionType = "search"     // 搜索
	ActionUseObject ActionType = "use_object" // 使用物品

	// 特殊动作
	ActionSecondWind      ActionType = "second_wind"      // 复苏之风（战士）
	ActionActionSurge     ActionType = "action_surge"     // 动作如潮（战士）
	ActionRage            ActionType = "rage"             // 狂暴（野蛮人）
	ActionWildShape       ActionType = "wild_shape"       // 野兽形态（德鲁伊）
	ActionChannelDivinity ActionType = "channel_divinity" // 引导神力（牧师）

	// 自定义动作
	ActionCustom ActionType = "custom"
)

// Action 代表一个可以执行的动作
type Action struct {
	Type        ActionType     `json:"type"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Details     map[string]any `json:"details,omitempty"`

	// 目标
	TargetIDs   []ID   `json:"target_ids,omitempty"`
	TargetPoint *Point `json:"target_point,omitempty"`

	// 武器/法术相关
	WeaponID  *ID    `json:"weapon_id,omitempty"`
	SpellID   string `json:"spell_id,omitempty"`
	SlotLevel int    `json:"slot_level,omitempty"` // 使用的法术位等级

	// 移动相关
	Distance int `json:"distance,omitempty"` // 移动距离（英尺）

	// 其他
	IsBonusAction bool `json:"is_bonus_action"` // 是否是附赠动作
	IsReaction    bool `json:"is_reaction"`     // 是否是反应
	IsFreeAction  bool `json:"is_free_action"`  // 是否是自由动作
}

// Reaction 代表一个反应动作
type Reaction struct {
	ID          ID     `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Trigger     string `json:"trigger"`  // 触发条件
	Action      Action `json:"action"`   // 触发的动作
	Used        bool   `json:"used"`     // 是否已使用
	ActorID     ID     `json:"actor_id"` // 拥有此反应的角色
}

// CommonReactions 常见的反应
var CommonReactions = map[string]string{
	"opportunity_attack": "借机攻击",
	"shield_spell":       "庇护术（Shield）",
	"feather_fall":       "羽落术",
	"deflection":         "偏转（蒙克）",
}

// BonusActionType 代表附赠动作类型
type BonusActionType string

const (
	BonusActionOffHandAttack BonusActionType = "off_hand_attack" // 副手攻击
	BonusActionDash          BonusActionType = "dash"            // 附赠冲刺（如盗贼）
	BonusActionHide          BonusActionType = "hide"            // 附赠躲藏
	BonusActionDisengage     BonusActionType = "disengage"       // 附赠脱离
	BonusActionHeal          BonusActionType = "heal"            // 附赠治疗
	BonusActionSpell         BonusActionType = "spell"           // 附赠施法
	BonusActionCustom        BonusActionType = "custom"          // 自定义
)

// MovementAction 代表移动动作
type MovementAction struct {
	From           Point `json:"from"`
	To             Point `json:"to"`
	Distance       int   `json:"distance"`
	SpeedUsed      int   `json:"speed_used"`
	SpeedRemaining int   `json:"speed_remaining"`
}
