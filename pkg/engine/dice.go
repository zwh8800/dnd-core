package engine

import (
	"context"
	"fmt"

	"github.com/zwh8800/dnd-core/internal/model"
)

// RollRequest 掷骰请求
type RollRequest struct {
	Expression string `json:"expression"` // 骰子表达式，如 "2d6+3", "1d20", "advantage"
	Modifier   int    `json:"modifier"`   // 额外修正值
	Reason     string `json:"reason"`     // 掷骰原因
}

// RollResult 掷骰结果
type RollResult struct {
	Expression string           `json:"expression"`
	Total      int              `json:"total"`
	Rolls      []model.DiceRoll `json:"rolls"`
	Modifier   int              `json:"modifier"`
	Reason     string           `json:"reason,omitempty"`
	Message    string           `json:"message"`
}

// RollAdvantageRequest 优势掷骰请求
type RollAdvantageRequest struct {
	Modifier int    `json:"modifier"`
	Reason   string `json:"reason"`
}

// RollDisadvantageRequest 劣势掷骰请求
type RollDisadvantageRequest struct {
	Modifier int    `json:"modifier"`
	Reason   string `json:"reason"`
}

// Roll 执行骰子投掷
func (e *Engine) Roll(ctx context.Context, req RollRequest) (*RollResult, error) {
	result, err := e.roller.Roll(req.Expression)
	if err != nil {
		return nil, fmt.Errorf("invalid dice expression: %w", err)
	}

	// 应用额外修正值
	if req.Modifier != 0 {
		result.Total += req.Modifier
		result.Modifier += req.Modifier
	}

	return &RollResult{
		Expression: result.Expression,
		Total:      result.Total,
		Rolls:      result.Rolls,
		Modifier:   result.Modifier,
		Reason:     req.Reason,
		Message:    fmt.Sprintf("掷骰 %s = %d", result.Expression, result.Total),
	}, nil
}

// RollAdvantage 执行优势掷骰（2d20取高）
func (e *Engine) RollAdvantage(ctx context.Context, req RollAdvantageRequest) (*RollResult, error) {
	result, err := e.roller.RollAdvantage(req.Modifier)
	if err != nil {
		return nil, err
	}

	return &RollResult{
		Expression: result.Expression,
		Total:      result.Total,
		Rolls:      result.Rolls,
		Modifier:   result.Modifier,
		Reason:     req.Reason,
		Message:    fmt.Sprintf("优势掷骰: %d (取高)", result.Total),
	}, nil
}

// RollDisadvantage 执行劣势掷骰（2d20取低）
func (e *Engine) RollDisadvantage(ctx context.Context, req RollDisadvantageRequest) (*RollResult, error) {
	result, err := e.roller.RollDisadvantage(req.Modifier)
	if err != nil {
		return nil, err
	}

	return &RollResult{
		Expression: result.Expression,
		Total:      result.Total,
		Rolls:      result.Rolls,
		Modifier:   result.Modifier,
		Reason:     req.Reason,
		Message:    fmt.Sprintf("劣势掷骰: %d (取低)", result.Total),
	}, nil
}

// RollAbility 属性掷骰（4d6去掉最低值）
func (e *Engine) RollAbility(ctx context.Context) (*RollResult, error) {
	// 投掷4d6
	result, err := e.roller.Roll("4d6")
	if err != nil {
		return nil, err
	}

	// 找到最低值并标记为丢弃
	minIdx := 0
	minVal := result.Rolls[0].Value
	for i, roll := range result.Rolls {
		if roll.Value < minVal {
			minVal = roll.Value
			minIdx = i
		}
	}
	result.Rolls[minIdx].Kept = false
	result.Rolls[minIdx].Dropped = true

	// 重新计算总和（去掉最低值）
	total := 0
	for _, roll := range result.Rolls {
		if roll.Kept {
			total += roll.Value
		}
	}
	result.Total = total

	return &RollResult{
		Expression: "4d6 (去掉最低)",
		Total:      total,
		Rolls:      result.Rolls,
		Message:    fmt.Sprintf("属性掷骰: %d", total),
	}, nil
}

// RollHitDice 生命骰掷骰
func (e *Engine) RollHitDice(ctx context.Context, diceType int, modifier int) (*RollResult, error) {
	expr := fmt.Sprintf("1d%d", diceType)
	result, err := e.roller.Roll(expr)
	if err != nil {
		return nil, err
	}

	if modifier != 0 {
		result.Total += modifier
		result.Modifier = modifier
	}

	return &RollResult{
		Expression: result.Expression,
		Total:      result.Total,
		Rolls:      result.Rolls,
		Modifier:   modifier,
		Message:    fmt.Sprintf("生命骰 %s = %d", expr, result.Total),
	}, nil
}
