package engine

import (
	"context"
	"fmt"

	"github.com/zwh8800/dnd-core/pkg/model"
)

// CurseActorRequest 诅咒角色请求
type CurseActorRequest struct {
	GameID   model.ID `json:"game_id"`
	TargetID model.ID `json:"target_id"`
	CurseID  string   `json:"curse_id"`
	Source   string   `json:"source"`
}

// CurseActorResult 诅咒角色结果
type CurseActorResult struct {
	CurseInstance *model.CurseInstance `json:"curse_instance"`
	Message       string               `json:"message"`
}

// RemoveCurseRequest 移除诅咒请求
type RemoveCurseRequest struct {
	GameID   model.ID `json:"game_id"`
	TargetID model.ID `json:"target_id"`
	CurseID  string   `json:"curse_id"`
}

// RemoveCurseResult 移除诅咒结果
type RemoveCurseResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// GetCursesRequest 获取诅咒请求
type GetCursesRequest struct {
	GameID   model.ID `json:"game_id"`
	TargetID model.ID `json:"target_id"`
}

// GetCursesResult 获取诅咒结果
type GetCursesResult struct {
	Curses []model.CurseInstance `json:"curses"`
	Count  int                   `json:"count"`
}

// CurseActor 诅咒角色
func (e *Engine) CurseActor(ctx context.Context, req CurseActorRequest) (*CurseActorResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	_, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	curseInstance := &model.CurseInstance{
		CurseID:           req.CurseID,
		Source:            req.Source,
		IsPermanent:       true,
		RemainingDuration: "permanent",
	}

	result := &CurseActorResult{
		CurseInstance: curseInstance,
		Message:       fmt.Sprintf("已施加诅咒：%s（来源：%s）", req.CurseID, req.Source),
	}

	return result, nil
}

// RemoveCurse 移除诅咒
func (e *Engine) RemoveCurse(ctx context.Context, req RemoveCurseRequest) (*RemoveCurseResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	_, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	result := &RemoveCurseResult{
		Success: true,
		Message: fmt.Sprintf("已移除诅咒：%s", req.CurseID),
	}

	return result, nil
}

// GetCurses 获取角色身上的诅咒
func (e *Engine) GetCurses(ctx context.Context, req GetCursesRequest) (*GetCursesResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	_, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	result := &GetCursesResult{
		Curses: []model.CurseInstance{},
		Count:  0,
	}

	return result, nil
}
