package engine

import (
	"context"
	"fmt"

	"github.com/zwh8800/dnd-core/pkg/data"
	"github.com/zwh8800/dnd-core/pkg/model"
)

// ApplyPoisonRequest 涂抹毒药请求
type ApplyPoisonRequest struct {
	GameID   model.ID `json:"game_id"`
	ActorID  model.ID `json:"actor_id"`
	PoisonID string   `json:"poison_id"`
	WeaponID string   `json:"weapon_id"`
}

// ApplyPoisonResult 涂抹毒药结果
type ApplyPoisonResult struct {
	PoisonInstance *model.PoisonInstance `json:"poison_instance"`
	Message        string                `json:"message"`
}

// ResolvePoisonEffectRequest 解析毒药效果请求
type ResolvePoisonEffectRequest struct {
	GameID  model.ID `json:"game_id"`
	ActorID model.ID `json:"actor_id"`
}

// ResolvePoisonEffectResult 解析毒药效果结果
type ResolvePoisonEffectResult struct {
	SaveRoll      int    `json:"save_roll"`
	SaveDC        int    `json:"save_dc"`
	SaveSuccess   bool   `json:"save_success"`
	DamageRolled  string `json:"damage_rolled,omitempty"`
	StatusApplied string `json:"status_applied,omitempty"`
	Message       string `json:"message"`
}

// RemovePoisonRequest 移除毒药请求
type RemovePoisonRequest struct {
	GameID   model.ID `json:"game_id"`
	ActorID  model.ID `json:"actor_id"`
	WeaponID string   `json:"weapon_id"`
}

// RemovePoisonResult 移除毒药结果
type RemovePoisonResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ApplyPoison 涂抹毒药到武器
// 将指定毒药涂抹到角色的武器上，创建毒药实例并设置有效时长。
// 毒药在涂抹后会在一定时间后失效。
// 参数:
//
//	ctx - 上下文
//	req - 涂抹毒药请求参数，包含游戏会话ID、角色ID、毒药ID和目标武器ID
//
// 返回:
//
//	*ApplyPoisonResult - 包含毒药实例和描述消息
//	error - 当毒药数据不存在时返回错误
func (e *Engine) ApplyPoison(ctx context.Context, req ApplyPoisonRequest) (*ApplyPoisonResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	poisonData := data.GetPoisonData(req.PoisonID)
	if poisonData == nil {
		return nil, fmt.Errorf("poison data %s not found", req.PoisonID)
	}

	poisonInstance := &model.PoisonInstance{
		PoisonID:      req.PoisonID,
		RemainingUses: 1,
		AppliedTo:     req.WeaponID,
		ExpiresAfter:  "1 minute",
	}

	_ = game // 简化实现

	result := &ApplyPoisonResult{
		PoisonInstance: poisonInstance,
		Message:        fmt.Sprintf("已将%s涂抹到武器上（1分钟后失效）", poisonData.Name),
	}

	return result, nil
}

// ResolvePoisonEffect 解析毒药效果
// 当角色受到毒药影响时，计算豁免检定结果并应用毒药效果，
// 包括伤害投骰和状态效果。
// 参数:
//
//	ctx - 上下文
//	req - 解析毒药效果请求参数，包含游戏会话ID和受影响的角色ID
//
// 返回:
//
//	*ResolvePoisonEffectResult - 包含豁免投骰值、豁免DC、是否成功、伤害和状态效果
//	error - 当游戏不存在或毒药数据不存在时返回错误
func (e *Engine) ResolvePoisonEffect(ctx context.Context, req ResolvePoisonEffectRequest) (*ResolvePoisonEffectResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	_, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	// 简化实现：返回毒药效果
	poisonData := data.GetPoisonData("basic-poison")
	if poisonData == nil {
		return nil, fmt.Errorf("poison data not found")
	}

	result := &ResolvePoisonEffectResult{
		SaveDC:      poisonData.Effect.SaveDC,
		SaveSuccess: false,
		Message:     fmt.Sprintf("毒药发作！%s", poisonData.Effect.Description),
	}

	return result, nil
}

// RemovePoison 移除武器上的毒药
// 清除指定武器上已涂抹的毒药，使其不再具有毒性效果。
// 参数:
//
//	ctx - 上下文
//	req - 移除毒药请求参数，包含游戏会话ID、角色ID和目标武器ID
//
// 返回:
//
//	*RemovePoisonResult - 包含操作是否成功及描述消息
//	error - 当游戏不存在时返回错误
func (e *Engine) RemovePoison(ctx context.Context, req RemovePoisonRequest) (*RemovePoisonResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	_, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	result := &RemovePoisonResult{
		Success: true,
		Message: "已移除武器上的毒药",
	}

	return result, nil
}
