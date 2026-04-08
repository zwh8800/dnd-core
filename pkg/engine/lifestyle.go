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

// SetLifestyle 设置角色的生活方式等级
// 根据指定的生活方式等级更新角色的日常开销，并保存游戏状态。
// 参数:
//
//	ctx - 上下文
//	req - 设置生活方式请求，包含游戏ID和生活方式等级
//
// 返回:
//
//	*SetLifestyleResult - 设置结果，包含生活方式等级、每日花费和描述信息
//	error - 错误信息，如游戏加载失败或生活方式等级无效
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

// AdvanceGameTime 推进游戏时间并扣除生活方式开销
// 根据当前生活方式等级计算指定天数的总费用，更新已过去的天数和累计花费。
// 参数:
//
//	ctx - 上下文
//	req - 推进游戏时间请求，包含游戏ID和要推进的天数
//
// 返回:
//
//	*AdvanceGameTimeResult - 推进结果，包含推进天数、总费用和支付状态
//	error - 错误信息，如游戏加载失败或未设置生活方式
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
