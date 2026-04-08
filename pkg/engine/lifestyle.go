package engine

import (
	"context"
	"fmt"

	"github.com/zwh8800/dnd-core/pkg/data"
	"github.com/zwh8800/dnd-core/pkg/model"
	"github.com/zwh8800/dnd-core/pkg/rules"
)

// SetLifestyleRequest 设置生活方式请求
type SetLifestyleRequest struct {
	GameID model.ID            `json:"game_id"`
	Tier   model.LifestyleTier `json:"tier"`
}

// SetLifestyleResult 设置生活方式结果
type SetLifestyleResult struct {
	Tier        model.LifestyleTier `json:"tier"`
	DailyCost   int                 `json:"daily_cost"`
	Description string              `json:"description"`
	Message     string              `json:"message"`
}

// SetLifestyle 设置生活方式
func (e *Engine) SetLifestyle(ctx context.Context, req SetLifestyleRequest) (*SetLifestyleResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if game.Lifestyle == nil {
		game.Lifestyle = &model.LifestyleState{}
	}

	game.Lifestyle.CurrentTier = req.Tier

	lifestyleData := data.GetLifestyleData(req.Tier)
	if lifestyleData == nil {
		return nil, fmt.Errorf("invalid lifestyle tier: %s", req.Tier)
	}

	game.Lifestyle.DailyCost = lifestyleData.DailyCost

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &SetLifestyleResult{
		Tier:        req.Tier,
		DailyCost:   lifestyleData.DailyCost,
		Description: lifestyleData.Description,
		Message:     fmt.Sprintf("生活方式设置为 %s，每日花费 %d 铜币", req.Tier, lifestyleData.DailyCost),
	}, nil
}

// AdvanceGameTimeRequest 推进游戏时间请求
type AdvanceGameTimeRequest struct {
	GameID model.ID `json:"game_id"`
	Days   int      `json:"days"`
}

// AdvanceGameTimeResult 推进游戏时间结果
type AdvanceGameTimeResult struct {
	DaysAdvanced   int    `json:"days_advanced"`
	TotalCost      int    `json:"total_cost"`
	PaymentSuccess bool   `json:"payment_success"`
	Message        string `json:"message"`
}

// AdvanceGameTime 推进游戏时间并扣除开销
func (e *Engine) AdvanceGameTime(ctx context.Context, req AdvanceGameTimeRequest) (*AdvanceGameTimeResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if game.Lifestyle == nil {
		return nil, fmt.Errorf("no lifestyle set")
	}

	cost := rules.CalculateLifestyleCost(game.Lifestyle.CurrentTier, req.Days)
	totalCost := cost.DailyCost

	game.Lifestyle.TotalSpent += totalCost
	game.Lifestyle.DaysElapsed += req.Days
	game.Lifestyle.LastPayment = totalCost

	result := &AdvanceGameTimeResult{
		DaysAdvanced:   req.Days,
		TotalCost:      totalCost,
		PaymentSuccess: true,
		Message:        fmt.Sprintf("时间推进 %d 天，生活方式费用 %d 铜币", req.Days, totalCost),
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return result, nil
}
