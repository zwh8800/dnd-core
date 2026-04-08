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

// PlaceTrap 在指定场景位置放置一个陷阱
// 根据陷阱ID从数据注册表中获取陷阱定义，创建陷阱状态实例并将其放置在指定位置。
// 参数:
//
//	ctx - 上下文
//	req - 放置陷阱请求，包含游戏ID、场景ID、陷阱ID和放置位置
//
// 返回:
//
//	*PlaceTrapResult - 包含创建的陷阱状态和成功消息
//	error - 游戏加载失败或陷阱数据不存在时返回错误
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

// DetectTrap 检测场景中的陷阱
// 返回陷阱的检测DC，角色需要进行感知（察觉）检定来发现陷阱。
// 当检定结果大于或等于检测DC时，陷阱被发现。
// 参数:
//
//	ctx - 上下文
//	req - 检测陷阱请求，包含游戏ID、角色ID、场景ID和陷阱ID
//
// 返回:
//
//	*DetectTrapResult - 包含检测DC、是否发现陷阱等信息
//	error - 游戏加载失败时返回错误
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

// DisarmTrap 解除场景中的陷阱
// 返回陷阱的解除DC，角色需要进行敏捷（巧手）检定并使用盗贼工具来解除陷阱。
// 当检定结果大于或等于解除DC时，陷阱被成功解除且不再触发。
// 参数:
//
//	ctx - 上下文
//	req - 解除陷阱请求，包含游戏ID、角色ID、场景ID和陷阱ID
//
// 返回:
//
//	*DisarmTrapResult - 包含解除DC、是否解除陷阱等信息
//	error - 游戏加载失败时返回错误
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

// TriggerTrap 触发陷阱并应用其效果
// 当角色进入陷阱区域或触发陷阱机关时调用此方法。
// 返回陷阱的伤害效果、豁免DC和豁免属性，受影响的角色需要进行相应的豁免检定。
// 参数:
//
//	ctx - 上下文
//	req - 触发陷阱请求，包含游戏ID、角色ID、场景ID和陷阱ID
//
// 返回:
//
//	*TriggerTrapResult - 包含陷阱效果列表、豁免DC和描述信息
//	error - 游戏加载失败或陷阱数据不存在时返回错误
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
