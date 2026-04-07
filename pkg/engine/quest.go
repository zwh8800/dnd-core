package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/zwh8800/dnd-core/internal/model"
)

// QuestInput 任务输入
type QuestInput struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	GiverID     model.ID            `json:"giver_id"`
	GiverName   string              `json:"giver_name"`
	Objectives  []ObjectiveInput    `json:"objectives"`
	Rewards     *model.QuestRewards `json:"rewards,omitempty"`
}

// ObjectiveInput 目标输入
type ObjectiveInput struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Required    int    `json:"required"`
	Optional    bool   `json:"optional"`
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
func (e *Engine) CreateQuest(ctx context.Context, gameID model.ID, input QuestInput) (*QuestResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	quest := model.NewQuest(input.Name, input.Description)
	quest.GiverID = input.GiverID
	quest.GiverName = input.GiverName

	// 添加目标
	for _, objInput := range input.Objectives {
		quest.AddObjective(objInput.ID, objInput.Description, objInput.Required, objInput.Optional)
	}

	// 设置奖励
	if input.Rewards != nil {
		quest.Rewards = *input.Rewards
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
func (e *Engine) GetQuest(ctx context.Context, gameID model.ID, questID model.ID) (*QuestInfo, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	quest, ok := game.Quests[questID]
	if !ok {
		return nil, ErrNotFound
	}

	return questToInfo(quest), nil
}

// ListQuests 列出所有任务
func (e *Engine) ListQuests(ctx context.Context, gameID model.ID, status *model.QuestStatus) ([]QuestInfo, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	result := make([]QuestInfo, 0, len(game.Quests))
	for _, quest := range game.Quests {
		if status != nil && quest.Status != *status {
			continue
		}
		result = append(result, *questToInfo(quest))
	}

	return result, nil
}

// AcceptQuest 接受任务
func (e *Engine) AcceptQuest(ctx context.Context, gameID model.ID, questID model.ID, actorID model.ID) (*QuestResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	quest, ok := game.Quests[questID]
	if !ok {
		return nil, ErrNotFound
	}

	// 检查角色是否存在
	if _, ok := game.GetActor(actorID); !ok {
		return nil, ErrNotFound
	}

	quest.Accept(actorID)

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &QuestResult{
		Quest:   quest,
		Message: fmt.Sprintf("%s 接受了任务: %s", actorID, quest.Name),
	}, nil
}

// UpdateQuestObjective 更新任务目标进度
func (e *Engine) UpdateQuestObjective(ctx context.Context, gameID model.ID, questID model.ID, objectiveID string, progress int) (*QuestResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	quest, ok := game.Quests[questID]
	if !ok {
		return nil, ErrNotFound
	}

	quest.UpdateProgress(objectiveID, progress)

	// 检查任务是否完成
	if quest.IsComplete() && quest.Status == model.QuestStatusActive {
		quest.Complete()
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	message := fmt.Sprintf("更新了任务目标进度: %s", objectiveID)
	if quest.Status == model.QuestStatusCompleted {
		message = fmt.Sprintf("任务完成: %s", quest.Name)
	}

	return &QuestResult{
		Quest:   quest,
		Message: message,
	}, nil
}

// CompleteQuest 完成任务并发放奖励
func (e *Engine) CompleteQuest(ctx context.Context, gameID model.ID, questID model.ID) (*QuestResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	quest, ok := game.Quests[questID]
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
func (e *Engine) FailQuest(ctx context.Context, gameID model.ID, questID model.ID) (*QuestResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	quest, ok := game.Quests[questID]
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
func (e *Engine) DeleteQuest(ctx context.Context, gameID model.ID, questID model.ID) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return err
	}

	if _, ok := game.Quests[questID]; !ok {
		return ErrNotFound
	}

	delete(game.Quests, questID)

	if err := e.saveGame(ctx, game); err != nil {
		return err
	}

	return nil
}

// GetActorQuests 获取角色的任务列表
func (e *Engine) GetActorQuests(ctx context.Context, gameID model.ID, actorID model.ID, status *model.QuestStatus) ([]QuestInfo, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	result := make([]QuestInfo, 0)
	for _, quest := range game.Quests {
		// 检查角色是否接受了该任务
		accepted := false
		for _, id := range quest.AcceptedBy {
			if id == actorID {
				accepted = true
				break
			}
		}

		if !accepted {
			continue
		}

		if status != nil && quest.Status != *status {
			continue
		}

		result = append(result, *questToInfo(quest))
	}

	return result, nil
}

// GetQuestGiverQuests 获取NPC发布的任务列表
func (e *Engine) GetQuestGiverQuests(ctx context.Context, gameID model.ID, giverID model.ID) ([]QuestInfo, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	result := make([]QuestInfo, 0)
	for _, quest := range game.Quests {
		if quest.GiverID == giverID {
			result = append(result, *questToInfo(quest))
		}
	}

	return result, nil
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
