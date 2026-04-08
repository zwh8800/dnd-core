package engine

import (
	"context"
	"fmt"

	"github.com/zwh8800/dnd-core/pkg/model"
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
	Modifier int    `json:"modifier"` // 额外修正值
	Reason   string `json:"reason"`   // 掷骰原因
}

// RollAbilityRequest 属性掷骰请求（4d6去掉最低值）
type RollAbilityRequest struct {
	// 属性掷骰无需参数，使用空结构体以统一API签名
}

// RollHitDiceRequest 生命骰掷骰请求
type RollHitDiceRequest struct {
	DiceType int `json:"dice_type"` // 生命骰类型，如 6(d6), 8(d8), 10(d10), 12(d12)
	Modifier int `json:"modifier"`  // 额外修正值（通常为体质修正）
}

// Roll 执行骰子投掷
// 解析并执行指定的骰子表达式（如 "2d6+3"、"1d20"），支持额外的修正值
// 参数:
//
//	ctx - 上下文
//	req - 掷骰请求，包含骰子表达式、额外修正值和掷骰原因
//
// 返回:
//
//	*RollResult - 掷骰结果，包含表达式、总值、各次投掷明细和消息
//	error - 骰子表达式无效时返回错误
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
// 投掷两枚d20骰子并取较高值，加上修正值后返回结果
// 参数:
//
//	ctx - 上下文
//	req - 优势掷骰请求，包含额外修正值和掷骰原因
//
// 返回:
//
//	*RollResult - 掷骰结果，包含表达式、总值、各次投掷明细和消息
//	error - 掷骰失败时返回错误
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
// 投掷两枚d20骰子并取较低值，加上修正值后返回结果
// 参数:
//
//	ctx - 上下文
//	req - 劣势掷骰请求，包含额外修正值和掷骰原因
//
// 返回:
//
//	*RollResult - 掷骰结果，包含表达式、总值、各次投掷明细和消息
//	error - 掷骰失败时返回错误
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
// 按照D&D标准属性生成规则，投掷4枚d6骰子并去掉最低值，用于生成角色属性值
// 参数:
//
//	ctx - 上下文
//	req - 属性掷骰请求（空结构体，无需参数）
//
// 返回:
//
//	*RollResult - 掷骰结果，包含表达式、总值（去掉最低值后的总和）、各次投掷明细（标记被丢弃的最低值）
//	error - 掷骰失败时返回错误
func (e *Engine) RollAbility(ctx context.Context, req RollAbilityRequest) (*RollResult, error) {
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
// 根据指定的生命骰类型（如d6、d8、d10、d12）投掷1枚骰子，通常用于角色短休恢复生命值
// 参数:
//
//	ctx - 上下文
//	req - 生命骰掷骰请求，包含生命骰类型和额外修正值（通常为体质修正）
//
// 返回:
//
//	*RollResult - 掷骰结果，包含表达式、总值（骰子结果+修正值）、各次投掷明细和消息
//	error - 骰子表达式无效时返回错误
func (e *Engine) RollHitDice(ctx context.Context, req RollHitDiceRequest) (*RollResult, error) {
	expr := fmt.Sprintf("1d%d", req.DiceType)
	result, err := e.roller.Roll(expr)
	if err != nil {
		return nil, err
	}

	if req.Modifier != 0 {
		result.Total += req.Modifier
		result.Modifier = req.Modifier
	}

	return &RollResult{
		Expression: result.Expression,
		Total:      result.Total,
		Rolls:      result.Rolls,
		Modifier:   req.Modifier,
		Message:    fmt.Sprintf("生命骰 %s = %d", expr, result.Total),
	}, nil
}
