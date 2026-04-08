package engine

import (
	"context"
	"fmt"

	"github.com/zwh8800/dnd-core/pkg/model"
)

// StartCraftingRequest 开始制作请求
type StartCraftingRequest struct {
	GameID   model.ID `json:"game_id"`
	ActorID  model.ID `json:"actor_id"`
	RecipeID string   `json:"recipe_id"`
}

// StartCraftingResult 开始制作结果
type StartCraftingResult struct {
	Progress *model.CraftingProgress `json:"progress"`
	Message  string                  `json:"message"`
}

// AdvanceCraftingRequest 推进制作请求
type AdvanceCraftingRequest struct {
	GameID  model.ID `json:"game_id"`
	ActorID model.ID `json:"actor_id"`
	Days    int      `json:"days"`
}

// AdvanceCraftingResult 推进制作结果
type AdvanceCraftingResult struct {
	Progress   *model.CraftingProgress `json:"progress"`
	IsComplete bool                    `json:"is_complete"`
	Message    string                  `json:"message"`
}

// CompleteCraftingRequest 完成制作请求
type CompleteCraftingRequest struct {
	GameID   model.ID `json:"game_id"`
	ActorID  model.ID `json:"actor_id"`
	RecipeID string   `json:"recipe_id"`
}

// CompleteCraftingResult 完成制作结果
type CompleteCraftingResult struct {
	Success bool   `json:"success"`
	ItemID  string `json:"item_id,omitempty"`
	Message string `json:"message"`
}

// StartCrafting 开始制作
func (e *Engine) StartCrafting(ctx context.Context, req StartCraftingRequest) (*StartCraftingResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	_, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	// 简化实现
	progress := &model.CraftingProgress{
		RecipeID:   req.RecipeID,
		DaysWorked: 0,
		TotalDays:  7,
	}

	result := &StartCraftingResult{
		Progress: progress,
		Message:  fmt.Sprintf("开始制作：%s（需要%d天）", req.RecipeID, progress.TotalDays),
	}

	return result, nil
}

// AdvanceCrafting 推进制作进度
func (e *Engine) AdvanceCrafting(ctx context.Context, req AdvanceCraftingRequest) (*AdvanceCraftingResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	_, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	// 简化实现
	progress := &model.CraftingProgress{
		DaysWorked: req.Days,
		TotalDays:  7,
	}

	isComplete := progress.DaysWorked >= progress.TotalDays

	result := &AdvanceCraftingResult{
		Progress:   progress,
		IsComplete: isComplete,
		Message:    fmt.Sprintf("制作进度：%d/%d天", progress.DaysWorked, progress.TotalDays),
	}

	return result, nil
}

// CompleteCrafting 完成制作
func (e *Engine) CompleteCrafting(ctx context.Context, req CompleteCraftingRequest) (*CompleteCraftingResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	_, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	result := &CompleteCraftingResult{
		Success: true,
		ItemID:  req.RecipeID,
		Message: fmt.Sprintf("制作完成：获得 %s", req.RecipeID),
	}

	return result, nil
}
