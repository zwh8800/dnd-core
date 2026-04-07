package model

import "time"

// QuestStatus 代表任务状态
type QuestStatus string

const (
	QuestStatusAvailable QuestStatus = "available" // 可用（未接受）
	QuestStatusActive    QuestStatus = "active"    // 进行中
	QuestStatusCompleted QuestStatus = "completed" // 已完成
	QuestStatusFailed    QuestStatus = "failed"    // 已失败
)

// ObjectiveStatus 代表目标状态
type ObjectiveStatus string

const (
	ObjectiveStatusIncomplete ObjectiveStatus = "incomplete" // 未完成
	ObjectiveStatusComplete   ObjectiveStatus = "complete"   // 已完成
)

// QuestRewards 代表任务奖励
type QuestRewards struct {
	Experience int      `json:"experience"`  // 经验值
	Gold       int      `json:"gold"`        // 金币
	Currency   Currency `json:"currency"`    // 货币
	Items      []ID     `json:"items"`       // 物品ID列表
	Features   []string `json:"features"`    // 特性/能力
	FactionRep int      `json:"faction_rep"` // 声望
}

// QuestObjective 代表任务目标
type QuestObjective struct {
	ID          string          `json:"id"`
	Description string          `json:"description"`
	Status      ObjectiveStatus `json:"status"`
	Progress    int             `json:"progress"`     // 当前进度
	Required    int             `json:"required"`     // 需要的进度（如杀10个怪）
	Optional    bool            `json:"optional"`     // 是否是可选目标
	CompletedAt *time.Time      `json:"completed_at"` // 完成时间
}

// Quest 代表一个任务
type Quest struct {
	ID          ID          `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Status      QuestStatus `json:"status"`

	// 任务 giver
	GiverID   ID     `json:"giver_id"`   // 给任务的NPC ID
	GiverName string `json:"giver_name"` // 给任务的NPC名字

	// 目标
	Objectives []*QuestObjective `json:"objectives"`

	// 奖励
	Rewards QuestRewards `json:"rewards"`

	// 接受任务的角色
	AcceptedBy []ID `json:"accepted_by"`

	// 时间
	CreatedAt   time.Time  `json:"created_at"`
	AcceptedAt  *time.Time `json:"accepted_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`

	// 自定义数据
	CustomData map[string]any `json:"custom_data,omitempty"`
}

// NewQuest 创建新任务
func NewQuest(name, description string) *Quest {
	return &Quest{
		ID:          NewID(),
		Name:        name,
		Description: description,
		Status:      QuestStatusAvailable,
		Objectives:  make([]*QuestObjective, 0),
		AcceptedBy:  make([]ID, 0),
		Rewards: QuestRewards{
			Currency: Currency{},
			Items:    make([]ID, 0),
			Features: make([]string, 0),
		},
		CreatedAt:  time.Now(),
		CustomData: make(map[string]any),
	}
}

// AddObjective 添加任务目标
func (q *Quest) AddObjective(id, description string, required int, optional bool) {
	q.Objectives = append(q.Objectives, &QuestObjective{
		ID:          id,
		Description: description,
		Status:      ObjectiveStatusIncomplete,
		Progress:    0,
		Required:    required,
		Optional:    optional,
	})
}

// UpdateProgress 更新目标进度
func (q *Quest) UpdateProgress(objectiveID string, progress int) {
	for _, obj := range q.Objectives {
		if obj.ID == objectiveID {
			obj.Progress = progress
			if obj.Progress >= obj.Required {
				obj.Status = ObjectiveStatusComplete
				now := time.Now()
				obj.CompletedAt = &now
			}
			return
		}
	}
}

// IsComplete 检查任务是否完成（所有必须目标都完成）
func (q *Quest) IsComplete() bool {
	for _, obj := range q.Objectives {
		if !obj.Optional && obj.Status != ObjectiveStatusComplete {
			return false
		}
	}
	return len(q.Objectives) > 0
}

// Accept 接受任务
func (q *Quest) Accept(actorID ID) {
	q.Status = QuestStatusActive
	now := time.Now()
	q.AcceptedAt = &now
	q.AcceptedBy = append(q.AcceptedBy, actorID)
}

// Complete 完成任务
func (q *Quest) Complete() {
	q.Status = QuestStatusCompleted
	now := time.Now()
	q.CompletedAt = &now
}

// Fail 标记任务失败
func (q *Quest) Fail() {
	q.Status = QuestStatusFailed
}
