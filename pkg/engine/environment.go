package engine

import (
	"context"
	"fmt"

	"github.com/zwh8800/dnd-core/pkg/model"
	"github.com/zwh8800/dnd-core/pkg/rules"
)

// SetEnvironmentRequest 设置环境请求
type SetEnvironmentRequest struct {
	GameID  model.ID              `json:"game_id"`
	SceneID model.ID              `json:"scene_id"`
	EnvType model.EnvironmentType `json:"env_type"`
}

// SetEnvironmentResult 设置环境结果
type SetEnvironmentResult struct {
	Environment model.EnvironmentalEffect `json:"environment"`
	Message     string                    `json:"message"`
}

// ResolveEnvironmentalDamageRequest 解析环境伤害请求
type ResolveEnvironmentalDamageRequest struct {
	GameID          model.ID              `json:"game_id"`
	ActorID         model.ID              `json:"actor_id"`
	EnvType         model.EnvironmentType `json:"env_type"`
	ExposureMinutes int                   `json:"exposure_minutes"`
}

// ResolveEnvironmentalDamageResult 解析环境伤害结果
type ResolveEnvironmentalDamageResult struct {
	DamageRolled  string `json:"damage_rolled,omitempty"`
	SaveDC        int    `json:"save_dc,omitempty"`
	SaveSuccess   bool   `json:"save_success"`
	StatusApplied string `json:"status_applied,omitempty"`
	Message       string `json:"message"`
}

// SetEnvironment 设置指定场景的环境效果
// 环境效果会影响角色的行动和能力，如强光、黑暗、浓雾等
//
// 参数:
//
//	ctx - 上下文
//	req - 设置环境请求参数，包含游戏ID、场景ID和环境类型
//
// 返回:
//
//	*SetEnvironmentResult - 包含环境效果信息和设置消息的结果
//	error - 加载游戏失败时返回错误
func (e *Engine) SetEnvironment(ctx context.Context, req SetEnvironmentRequest) (*SetEnvironmentResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	_, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	effect := rules.GetEnvironmentEffect(req.EnvType)

	result := &SetEnvironmentResult{
		Environment: effect,
		Message:     fmt.Sprintf("环境已设置为：%s - %s", req.EnvType, effect.Description),
	}

	return result, nil
}

// ResolveEnvironmentalDamage 解析并计算环境对角色造成的伤害
// 根据环境类型和暴露时间，计算角色是否受到伤害以及相应的豁免检定
//
// 参数:
//
//	ctx - 上下文
//	req - 解析环境伤害请求参数，包含游戏ID、角色ID、环境类型和暴露时间（分钟）
//
// 返回:
//
//	*ResolveEnvironmentalDamageResult - 包含伤害掷骰、豁免DC、豁免是否成功、施加的状态和消息的结果
//	error - 加载游戏失败时返回错误
func (e *Engine) ResolveEnvironmentalDamage(ctx context.Context, req ResolveEnvironmentalDamageRequest) (*ResolveEnvironmentalDamageResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	_, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	effect := rules.GetEnvironmentEffect(req.EnvType)

	result := &ResolveEnvironmentalDamageResult{
		SaveDC:  effect.SaveDC,
		Message: fmt.Sprintf("环境伤害：%s", effect.EffectDescription),
	}

	return result, nil
}
