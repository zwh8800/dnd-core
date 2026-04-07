package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/zwh8800/dnd-core/internal/model"
)

// CreateQuestRequest 创建任务请求（整合 QuestInput）
type CreateQuestRequest struct {
	GameID      model.ID           `json:"game_id"`     // 游戏会话ID
	Name        string             `json:"name"`        // 任务名称
	Description string             `json:"description"` // 任务描述
	GiverID     model.ID           `json:"giver_id"`    // 任务发布者ID
	GiverName   string             `json:"giver_name"`  // 任务发布者名称
	Objectives  []ObjectiveInput   `json:"objectives"`  // 任务目标列表
	Rewards     *QuestRewardsInput `json:"rewards"`     // 任务奖励（可选）
}

// ObjectiveInput 目标输入
type ObjectiveInput struct {
	ID          string `json:"id"`          // 目标ID
	Description string `json:"description"` // 目标描述
	Required    int    `json:"required"`    // 完成所需数量
}

// QuestRewardsInput 任务奖励输入
type QuestRewardsInput struct {
	Experience int               `json:"experience"` // 经验奖励
	Gold       int               `json:"gold"`       // 金币奖励
	Items      []ItemRewardInput `json:"items"`      // 物品奖励列表
}

// ItemRewardInput 物品奖励输入
type ItemRewardInput struct {
	Name        string `json:"name"`        // 物品名称
	Description string `json:"description"` // 物品描述
	Quantity    int    `json:"quantity"`    // 数量
}

// GetQuestRequest 获取任务请求
type GetQuestRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID
	QuestID model.ID `json:"quest_id"` // 任务ID
}

// ListQuestsRequest 列出任务请求
type ListQuestsRequest struct {
	GameID model.ID           `json:"game_id"`          // 游戏会话ID
	Status *model.QuestStatus `json:"status,omitempty"` // 按状态过滤（可选）
}

// ListQuestsResult 列出任务结果
type ListQuestsResult struct {
	Quests []QuestInfo `json:"quests"` // 任务列表
}

// AcceptQuestRequest 接受任务请求
type AcceptQuestRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID
	QuestID model.ID `json:"quest_id"` // 任务ID
	ActorID model.ID `json:"actor_id"` // 接受任务的角色ID
}

// UpdateQuestObjectiveRequest 更新任务目标请求
type UpdateQuestObjectiveRequest struct {
	GameID      model.ID `json:"game_id"`      // 游戏会话ID
	QuestID     model.ID `json:"quest_id"`     // 任务ID
	ObjectiveID string   `json:"objective_id"` // 目标ID
	Progress    int      `json:"progress"`     // 进度增量
}

// CompleteQuestRequest 完成任务请求
type CompleteQuestRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID
	QuestID model.ID `json:"quest_id"` // 任务ID
}

// FailQuestRequest 任务失败请求
type FailQuestRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID
	QuestID model.ID `json:"quest_id"` // 任务ID
}

// DeleteQuestRequest 删除任务请求
type DeleteQuestRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID
	QuestID model.ID `json:"quest_id"` // 任务ID
}

// GetActorQuestsRequest 获取角色任务请求
type GetActorQuestsRequest struct {
	GameID  model.ID           `json:"game_id"`          // 游戏会话ID
	ActorID model.ID           `json:"actor_id"`         // 角色ID
	Status  *model.QuestStatus `json:"status,omitempty"` // 按状态过滤（可选）
}

// GetActorQuestsResult 获取角色任务结果
type GetActorQuestsResult struct {
	Quests []QuestInfo `json:"quests"` // 任务列表
}

// GetQuestGiverQuestsRequest 获取任务发布者任务请求
type GetQuestGiverQuestsRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID
	GiverID model.ID `json:"giver_id"` // 任务发布者ID
}

// GetQuestGiverQuestsResult 获取任务发布者任务结果
type GetQuestGiverQuestsResult struct {
	Quests []QuestInfo `json:"quests"` // 任务列表
}

// QuestInput 任务输入
type QuestInput struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	GiverID     model.ID            `json:"giver_id"`
	GiverName   string              `json:"giver_name"`
	Objectives  []ObjectiveInput    `json:"objectives"`
	Rewards     *model.QuestRewards `json:"rewards,omitempty"`
}

// QuestResult 任务操作结果
type QuestResult struct {
	Quest   *model.Quest `json:"quest"`
	Message string       `json:"message"`
}

// QuestInfo 任务信息
type QuestInfo struct {
	ID          model.ID           `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Status      model.QuestStatus  `json:"status"`
	GiverID     model.ID           `json:"giver_id"`
	GiverName   string             `json:"giver_name"`
	Objectives  []ObjectiveInfo    `json:"objectives"`
	Rewards     model.QuestRewards `json:"rewards"`
	AcceptedBy  []model.ID         `json:"accepted_by"`
	CreatedAt   time.Time          `json:"created_at"`
	AcceptedAt  *time.Time         `json:"accepted_at,omitempty"`
	CompletedAt *time.Time         `json:"completed_at,omitempty"`
}

// ObjectiveInfo 目标信息
type ObjectiveInfo struct {
	ID          string                `json:"id"`
	Description string                `json:"description"`
	Status      model.ObjectiveStatus `json:"status"`
	Progress    int                   `json:"progress"`
	Required    int                   `json:"required"`
	Optional    bool                  `json:"optional"`
	CompletedAt *time.Time            `json:"completed_at"`
}

// CreateQuest 创建新任务
func (e *Engine) CreateQuest(ctx context.Context, req CreateQuestRequest) (*QuestResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	quest := model.NewQuest(req.Name, req.Description)
	quest.GiverID = req.GiverID
	quest.GiverName = req.GiverName

	// 添加目标
	for _, objInput := range req.Objectives {
		quest.AddObjective(objInput.ID, objInput.Description, objInput.Required, false)
	}

	// 设置奖励
	if req.Rewards != nil {
		quest.Rewards.Experience = req.Rewards.Experience
		quest.Rewards.Gold = req.Rewards.Gold
		// 物品奖励将在后续添加
		for range req.Rewards.Items {
		}
	}

	game.Quests[quest.ID] = quest

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &QuestResult{
		Quest:   quest,
		Message: fmt.Sprintf("创建了任务: %s", quest.Name),
	}, nil
}

// GetQuest 获取任务信息
func (e *Engine) GetQuest(ctx context.Context, req GetQuestRequest) (*QuestInfo, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	quest, ok := game.Quests[req.QuestID]
	if !ok {
		return nil, ErrNotFound
	}

	return questToInfo(quest), nil
}

// ListQuests 列出所有任务
func (e *Engine) ListQuests(ctx context.Context, req ListQuestsRequest) (*ListQuestsResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	result := make([]QuestInfo, 0, len(game.Quests))
	for _, quest := range game.Quests {
		if req.Status != nil && quest.Status != *req.Status {
			continue
		}
		result = append(result, *questToInfo(quest))
	}

	return &ListQuestsResult{Quests: result}, nil
}

// AcceptQuest 接受任务
func (e *Engine) AcceptQuest(ctx context.Context, req AcceptQuestRequest) (*QuestResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	quest, ok := game.Quests[req.QuestID]
	if !ok {
		return nil, ErrNotFound
	}

	// 检查角色是否存在
	if _, ok := game.GetActor(req.ActorID); !ok {
		return nil, ErrNotFound
	}

	quest.Accept(req.ActorID)

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &QuestResult{
		Quest:   quest,
		Message: fmt.Sprintf("%s 接受了任务: %s", req.ActorID, quest.Name),
	}, nil
}

// UpdateQuestObjective 更新任务目标进度
func (e *Engine) UpdateQuestObjective(ctx context.Context, req UpdateQuestObjectiveRequest) (*QuestResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	quest, ok := game.Quests[req.QuestID]
	if !ok {
		return nil, ErrNotFound
	}

	quest.UpdateProgress(req.ObjectiveID, req.Progress)

	// 检查任务是否完成
	if quest.IsComplete() && quest.Status == model.QuestStatusActive {
		quest.Complete()
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	message := fmt.Sprintf("更新了任务目标进度: %s", req.ObjectiveID)
	if quest.Status == model.QuestStatusCompleted {
		message = fmt.Sprintf("任务完成: %s", quest.Name)
	}

	return &QuestResult{
		Quest:   quest,
		Message: message,
	}, nil
}

// CompleteQuest 完成任务并发放奖励
func (e *Engine) CompleteQuest(ctx context.Context, req CompleteQuestRequest) (*QuestResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	quest, ok := game.Quests[req.QuestID]
	if !ok {
		return nil, ErrNotFound
	}

	if !quest.IsComplete() {
		return nil, fmt.Errorf("quest is not yet complete")
	}

	quest.Complete()

	// 发放奖励给接受任务的角色
	for _, actorID := range quest.AcceptedBy {
		actor, ok := game.GetActor(actorID)
		if !ok {
			continue
		}

		switch a := actor.(type) {
		case *model.PlayerCharacter:
			// 发放经验
			if quest.Rewards.Experience > 0 {
				a.Experience += quest.Rewards.Experience
			}
		}
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &QuestResult{
		Quest:   quest,
		Message: fmt.Sprintf("完成任务: %s，发放奖励", quest.Name),
	}, nil
}

// FailQuest 标记任务失败
func (e *Engine) FailQuest(ctx context.Context, req FailQuestRequest) (*QuestResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	quest, ok := game.Quests[req.QuestID]
	if !ok {
		return nil, ErrNotFound
	}

	quest.Fail()

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &QuestResult{
		Quest:   quest,
		Message: fmt.Sprintf("任务失败: %s", quest.Name),
	}, nil
}

// DeleteQuest 删除任务
func (e *Engine) DeleteQuest(ctx context.Context, req DeleteQuestRequest) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return err
	}

	if _, ok := game.Quests[req.QuestID]; !ok {
		return ErrNotFound
	}

	delete(game.Quests, req.QuestID)

	if err := e.saveGame(ctx, game); err != nil {
		return err
	}

	return nil
}

// GetActorQuests 获取角色的任务列表
func (e *Engine) GetActorQuests(ctx context.Context, req GetActorQuestsRequest) (*GetActorQuestsResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	result := make([]QuestInfo, 0)
	for _, quest := range game.Quests {
		// 检查角色是否接受了该任务
		accepted := false
		for _, id := range quest.AcceptedBy {
			if id == req.ActorID {
				accepted = true
				break
			}
		}

		if !accepted {
			continue
		}

		if req.Status != nil && quest.Status != *req.Status {
			continue
		}

		result = append(result, *questToInfo(quest))
	}

	return &GetActorQuestsResult{Quests: result}, nil
}

// GetQuestGiverQuests 获取NPC发布的任务列表
func (e *Engine) GetQuestGiverQuests(ctx context.Context, req GetQuestGiverQuestsRequest) (*GetQuestGiverQuestsResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	result := make([]QuestInfo, 0)
	for _, quest := range game.Quests {
		if quest.GiverID == req.GiverID {
			result = append(result, *questToInfo(quest))
		}
	}

	return &GetQuestGiverQuestsResult{Quests: result}, nil
}

// questToInfo 将任务模型转换为信息结构
func questToInfo(quest *model.Quest) *QuestInfo {
	objectives := make([]ObjectiveInfo, 0, len(quest.Objectives))
	for _, obj := range quest.Objectives {
		objectives = append(objectives, ObjectiveInfo{
			ID:          obj.ID,
			Description: obj.Description,
			Status:      obj.Status,
			Progress:    obj.Progress,
			Required:    obj.Required,
			Optional:    obj.Optional,
			CompletedAt: obj.CompletedAt,
		})
	}

	return &QuestInfo{
		ID:          quest.ID,
		Name:        quest.Name,
		Description: quest.Description,
		Status:      quest.Status,
		GiverID:     quest.GiverID,
		GiverName:   quest.GiverName,
		Objectives:  objectives,
		Rewards:     quest.Rewards,
		AcceptedBy:  quest.AcceptedBy,
		CreatedAt:   quest.CreatedAt,
		AcceptedAt:  quest.AcceptedAt,
		CompletedAt: quest.CompletedAt,
	}
}
