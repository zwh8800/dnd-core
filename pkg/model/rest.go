package model

// RestType 代表休息类型
type RestType string

const (
	RestTypeShort RestType = "short" // 短休（至少1小时）
	RestTypeLong  RestType = "long"  // 长休（至少8小时）
)

// RestPhase 代表长休的阶段
type RestPhase string

const (
	RestPhaseNotStarted  RestPhase = "not_started" // 未开始
	RestPhaseInProgress  RestPhase = "in_progress" // 进行中
	RestPhaseCompleted   RestPhase = "completed"   // 已完成
	RestPhaseInterrupted RestPhase = "interrupted" // 被打断
)

// RestState 代表休息状态
type RestState struct {
	Type     RestType  `json:"type"`
	Phase    RestPhase `json:"phase"`
	Started  bool      `json:"started"`
	Progress int       `json:"progress"` // 进度（小时）
	Required int       `json:"required"` // 需要的小时数

	// 参与的角色
	ParticipantIDs []ID `json:"participant_ids"`

	// 恢复信息
	HitDiceRecovered  map[ID]int  `json:"hit_dice_recovered"` // 每个角色恢复的生命骰
	HPRecovered       map[ID]int  `json:"hp_recovered"`       // 每个角色恢复的HP
	SlotsRestored     map[ID]bool `json:"slots_restored"`     // 每个角色是否恢复法术位
	ExhaustionReduced map[ID]bool `json:"exhaustion_reduced"` // 每个角色是否减少力竭
}

// NewShortRest 创建新的短休
func NewShortRest(participants []ID) *RestState {
	return &RestState{
		Type:              RestTypeShort,
		Phase:             RestPhaseNotStarted,
		Progress:          0,
		Required:          1, // 短休需要1小时
		ParticipantIDs:    participants,
		HitDiceRecovered:  make(map[ID]int),
		HPRecovered:       make(map[ID]int),
		SlotsRestored:     make(map[ID]bool),
		ExhaustionReduced: make(map[ID]bool),
	}
}

// NewLongRest 创建新的长休
func NewLongRest(participants []ID) *RestState {
	return &RestState{
		Type:              RestTypeLong,
		Phase:             RestPhaseNotStarted,
		Progress:          0,
		Required:          8, // 长休需要8小时
		ParticipantIDs:    participants,
		HitDiceRecovered:  make(map[ID]int),
		HPRecovered:       make(map[ID]int),
		SlotsRestored:     make(map[ID]bool),
		ExhaustionReduced: make(map[ID]bool),
	}
}

// Start 开始休息
func (r *RestState) Start() {
	r.Started = true
	r.Phase = RestPhaseInProgress
}

// Complete 完成休息
func (r *RestState) Complete() {
	r.Phase = RestPhaseCompleted
	r.Progress = r.Required
}

// Interrupt 打断休息
func (r *RestState) Interrupt() {
	r.Phase = RestPhaseInterrupted
}

// IsComplete 检查休息是否完成
func (r *RestState) IsComplete() bool {
	return r.Phase == RestPhaseCompleted
}

// IsInterrupted 检查休息是否被打断
func (r *RestState) IsInterrupted() bool {
	return r.Phase == RestPhaseInterrupted
}

// CanRecoverHitDice 检查是否可以恢复生命骰（长休）
func (r *RestState) CanRecoverHitDice() bool {
	return r.Type == RestTypeLong && r.IsComplete()
}

// GetHitDiceRecoveryAmount 获取长休可以恢复的生命骰数量
// 规则：长休结束时，你可以恢复一些已耗尽的生命骰。
// 你能恢复的总数最多等于你总等级的一半（至少一个）。
func GetHitDiceRecoveryAmount(totalLevel int) int {
	recovery := totalLevel / 2
	if recovery < 1 {
		return 1
	}
	return recovery
}
