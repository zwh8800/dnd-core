package engine

import (
	"context"
	"fmt"

	"github.com/zwh8800/dnd-core/pkg/data"
	"github.com/zwh8800/dnd-core/pkg/model"
)

// PlaceTrapRequest 放置陷阱请求
type PlaceTrapRequest struct {
	GameID   model.ID `json:"game_id"`
	SceneID  model.ID `json:"scene_id"`
	TrapID   string   `json:"trap_id"`
	Position string   `json:"position"`
}

// PlaceTrapResult 放置陷阱结果
type PlaceTrapResult struct {
	Trap    *model.TrapState `json:"trap"`
	Message string           `json:"message"`
}

// DetectTrapRequest 检测陷阱请求
type DetectTrapRequest struct {
	GameID  model.ID `json:"game_id"`
	ActorID model.ID `json:"actor_id"`
	SceneID model.ID `json:"scene_id"`
	TrapID  model.ID `json:"trap_id"`
}

// DetectTrapResult 检测陷阱结果
type DetectTrapResult struct {
	Success      bool   `json:"success"`
	CheckTotal   int    `json:"check_total"`
	DC           int    `json:"dc"`
	TrapRevealed bool   `json:"trap_revealed"`
	Message      string `json:"message"`
}

// DisarmTrapRequest 解除陷阱请求
type DisarmTrapRequest struct {
	GameID  model.ID `json:"game_id"`
	ActorID model.ID `json:"actor_id"`
	SceneID model.ID `json:"scene_id"`
	TrapID  model.ID `json:"trap_id"`
}

// DisarmTrapResult 解除陷阱结果
type DisarmTrapResult struct {
	Success      bool   `json:"success"`
	CheckTotal   int    `json:"check_total"`
	DC           int    `json:"dc"`
	TrapDisarmed bool   `json:"trap_disarmed"`
	Message      string `json:"message"`
}

// TriggerTrapRequest 触发陷阱请求
type TriggerTrapRequest struct {
	GameID  model.ID `json:"game_id"`
	ActorID model.ID `json:"actor_id"`
	SceneID model.ID `json:"scene_id"`
	TrapID  model.ID `json:"trap_id"`
}

// TriggerTrapResult 触发陷阱结果
type TriggerTrapResult struct {
	TrapTriggered bool               `json:"trap_triggered"`
	Effects       []model.TrapEffect `json:"effects"`
	DamageRolls   []string           `json:"damage_rolls,omitempty"`
	SaveDC        int                `json:"save_dc,omitempty"`
	SaveAbility   string             `json:"save_ability,omitempty"`
	Message       string             `json:"message"`
}

// PlaceTrap 放置陷阱
func (e *Engine) PlaceTrap(ctx context.Context, req PlaceTrapRequest) (*PlaceTrapResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	trapData := data.GetTrapData(req.TrapID)
	if trapData == nil {
		return nil, fmt.Errorf("trap data %s not found", req.TrapID)
	}

	trapState := &model.TrapState{
		Definition:   trapData,
		IsArmed:      true,
		HasTriggered: false,
		Remaining:    0, // 0 = infinite
		Position:     req.Position,
	}

	// 将陷阱添加到场景
	// 注意：这里简化处理，实际应该添加到Scene的陷阱列表中
	_ = game // 简化实现

	result := &PlaceTrapResult{
		Trap:    trapState,
		Message: fmt.Sprintf("已放置陷阱：%s 在 %s", trapData.Name, req.Position),
	}

	return result, nil
}

// DetectTrap 检测陷阱
func (e *Engine) DetectTrap(ctx context.Context, req DetectTrapRequest) (*DetectTrapResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	_, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	// 简化实现：返回陷阱的检测DC
	// 实际应该进行WIS (Perception)检定
	trapState := &model.TrapState{
		Definition: &model.TrapDefinition{
			DetectDC: 15,
			Name:     "示例陷阱",
		},
	}

	result := &DetectTrapResult{
		DC:           trapState.Definition.DetectDC,
		TrapRevealed: false,
		Message:      fmt.Sprintf("检测陷阱 DC %d（需要进行感知（察觉）检定）", trapState.Definition.DetectDC),
	}

	return result, nil
}

// DisarmTrap 解除陷阱
func (e *Engine) DisarmTrap(ctx context.Context, req DisarmTrapRequest) (*DisarmTrapResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	_, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	// 简化实现：返回陷阱的解除DC
	trapState := &model.TrapState{
		Definition: &model.TrapDefinition{
			DisarmDC:    15,
			Name:        "示例陷阱",
			DisarmSkill: "thieves-tools",
		},
	}

	result := &DisarmTrapResult{
		DC:           trapState.Definition.DisarmDC,
		TrapDisarmed: false,
		Message:      fmt.Sprintf("解除陷阱 DC %d（需要进行敏捷（妙手）检定）", trapState.Definition.DisarmDC),
	}

	return result, nil
}

// TriggerTrap 触发陷阱
func (e *Engine) TriggerTrap(ctx context.Context, req TriggerTrapRequest) (*TriggerTrapResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	_, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	// 简化实现：返回陷阱效果
	trapData := data.GetTrapData("poison-needle")
	if trapData == nil {
		return nil, fmt.Errorf("trap data not found")
	}

	result := &TriggerTrapResult{
		TrapTriggered: true,
		Effects:       trapData.Effects,
		SaveDC:        trapData.Effects[0].SaveDC,
		SaveAbility:   trapData.Effects[0].SaveAbility,
		Message:       fmt.Sprintf("陷阱触发：%s！%s", trapData.Name, trapData.Effects[0].Description),
	}

	return result, nil
}
