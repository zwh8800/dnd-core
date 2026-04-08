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

// SetEnvironment 设置环境
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

// ResolveEnvironmentalDamage 解析环境伤害
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
