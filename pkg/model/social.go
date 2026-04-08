package model

// NPCAttitude 代表 NPC 态度
type NPCAttitude string

const (
	AttitudeFriendly    NPCAttitude = "friendly"    // 友好
	AttitudeIndifferent NPCAttitude = "indifferent" // 冷漠
	AttitudeHostile     NPCAttitude = "hostile"     // 敌对
)

// NPCDisposition 代表 NPC 倾向
type NPCDisposition string

const (
	DispositionHelpful     NPCDisposition = "helpful"     // 乐于助人
	DispositionFriendly    NPCDisposition = "friendly"    // 友好
	DispositionIndifferent NPCDisposition = "indifferent" // 冷漠
	DispositionSuspicious  NPCDisposition = "suspicious"  // 多疑
	DispositionHostile     NPCDisposition = "hostile"     // 敌对
)

// SocialInteractionState 代表社交互动状态
type SocialInteractionState struct {
	CurrentAttitude  NPCAttitude    `json:"current_attitude"`  // 当前态度
	Disposition      NPCDisposition `json:"disposition"`       // 倾向
	Impressions      []string       `json:"impressions"`       // 已建立的印象
	InteractionCount int            `json:"interaction_count"` // 互动次数
	LastInteraction  string         `json:"last_interaction"`  // 上次互动类型
}

// SocialCheckType 代表社交检定类型
type SocialCheckType string

const (
	SocialCheckPersuasion   SocialCheckType = "persuasion"   // 游说
	SocialCheckDeception    SocialCheckType = "deception"    // 欺骗
	SocialCheckIntimidation SocialCheckType = "intimidation" // 威吓
	SocialCheckPerformance  SocialCheckType = "performance"  // 表演
)

// SocialInteractionResult 代表社交互动结果
type SocialInteractionResult struct {
	Success        bool            `json:"success"`
	AttitudeChange NPCAttitude     `json:"attitude_change,omitempty"`
	RollTotal      int             `json:"roll_total"`
	DC             int             `json:"dc"`
	CheckType      SocialCheckType `json:"check_type"`
	Message        string          `json:"message"`
}
