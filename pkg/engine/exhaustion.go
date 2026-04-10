package engine

import (
	"context"
	"fmt"
	"strings"

	"github.com/zwh8800/dnd-core/pkg/model"
	"github.com/zwh8800/dnd-core/pkg/rules"
)

// ============================================================================
// 力竭管理API结构体定义
// ============================================================================

// ApplyExhaustionRequest 应用力竭请求
// 用于对角色施加力竭等级（如长途跋涉、恶劣环境等）
type ApplyExhaustionRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID（必填）
	ActorID model.ID `json:"actor_id"` // 角色ID（必填）
	Levels  int      `json:"levels"`   // 增加的力竭等级（必填，通常为1）
}

// ApplyExhaustionResult 应用力竭结果
// 描述力竭应用后的状态和效果
type ApplyExhaustionResult struct {
	ActorID  model.ID `json:"actor_id"`  // 角色ID
	NewLevel int      `json:"new_level"` // 新的力竭等级
	Effects  []string `json:"effects"`   // 当前所有力竭效果
	IsDead   bool     `json:"is_dead"`   // 是否因力竭死亡（6级）
	Message  string   `json:"message"`   // 人类可读消息
}

// RemoveExhaustionRequest 移除力竭请求
// 用于减少角色的力竭等级（如长休后）
type RemoveExhaustionRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID（必填）
	ActorID model.ID `json:"actor_id"` // 角色ID（必填）
	Levels  int      `json:"levels"`   // 移除的力竭等级（必填，通常为1）
}

// RemoveExhaustionResult 移除力竭结果
type RemoveExhaustionResult struct {
	ActorID  model.ID `json:"actor_id"`  // 角色ID
	NewLevel int      `json:"new_level"` // 新的力竭等级
	Message  string   `json:"message"`   // 人类可读消息
}

// GetExhaustionStatusRequest 获取力竭状态请求
type GetExhaustionStatusRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID（必填）
	ActorID model.ID `json:"actor_id"` // 角色ID（必填）
}

// GetExhaustionStatusResult 力竭状态结果
type GetExhaustionStatusResult struct {
	ActorID      model.ID `json:"actor_id"`      // 角色ID
	CurrentLevel int      `json:"current_level"` // 当前力竭等级
	Effects      []string `json:"effects"`       // 当前所有力竭效果
	IsDead       bool     `json:"is_dead"`       // 是否因力竭死亡
	Message      string   `json:"message"`       // 人类可读消息
}

// ============================================================================
// 力竭管理API
// ============================================================================

// ApplyExhaustion 对角色施加力竭等级
// 力竭是D&D 5e中的累积负面状态，共6级：
//   - 1级：检定劣势
//   - 2级：速度减半
//   - 3级：攻击劣势（叠加1级的效果）
//   - 4级：HP最大值减半
//   - 5级：速度降为0
//   - 6级：死亡
//
// 参数:
//
//	ctx - 上下文
//	req - 应用力竭请求，包含游戏ID、角色ID和增加的力竭等级
//
// 返回:
//
//	*ApplyExhaustionResult - 应用力竭结果，包含新等级、效果和消息
//	error - 角色不存在、参数无效或保存失败时返回错误
func (e *Engine) ApplyExhaustion(ctx context.Context, req ApplyExhaustionRequest) (*ApplyExhaustionResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpApplyExhaustion); err != nil {
		return nil, err
	}

	if req.Levels < 1 {
		return nil, fmt.Errorf("力竭等级必须大于0")
	}

	actor, ok := game.GetActor(req.ActorID)
	if !ok {
		return nil, ErrNotFound
	}

	var baseActor *model.Actor
	switch a := actor.(type) {
	case *model.PlayerCharacter:
		baseActor = &a.Actor
	case *model.NPC:
		baseActor = &a.Actor
	case *model.Enemy:
		baseActor = &a.Actor
	case *model.Companion:
		baseActor = &a.Actor
	}

	// 增加力竭等级
	baseActor.Exhaustion += req.Levels

	// 检查是否超过6级
	if baseActor.Exhaustion > 6 {
		baseActor.Exhaustion = 6
	}

	// 应用力竭效果
	effects := rules.ApplyExhaustionEffects(baseActor.Exhaustion, baseActor)
	isDead := baseActor.Exhaustion >= 6

	// 如果6级力竭，HP归零
	if isDead {
		baseActor.HitPoints.Current = 0
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	message := fmt.Sprintf("%s 力竭等级提升至 %d", baseActor.Name, baseActor.Exhaustion)
	if len(effects) > 0 {
		message += fmt.Sprintf(" - 效果: %s", strings.Join(effects, ", "))
	}
	if isDead {
		message += " - 6级力竭，角色死亡！"
	}

	return &ApplyExhaustionResult{
		ActorID:  req.ActorID,
		NewLevel: baseActor.Exhaustion,
		Effects:  effects,
		IsDead:   isDead,
		Message:  message,
	}, nil
}

// RemoveExhaustion 移除角色的力竭等级
// 通常通过长休移除1级力竭。
// 参数:
//
//	ctx - 上下文
//	req - 移除力竭请求，包含游戏ID、角色ID和移除的力竭等级
//
// 返回:
//
//	*RemoveExhaustionResult - 移除力竭结果
//	error - 角色不存在、参数无效或保存失败时返回错误
func (e *Engine) RemoveExhaustion(ctx context.Context, req RemoveExhaustionRequest) (*RemoveExhaustionResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpRemoveExhaustion); err != nil {
		return nil, err
	}

	if req.Levels < 1 {
		return nil, fmt.Errorf("移除的力竭等级必须大于0")
	}

	actor, ok := game.GetActor(req.ActorID)
	if !ok {
		return nil, ErrNotFound
	}

	var baseActor *model.Actor
	switch a := actor.(type) {
	case *model.PlayerCharacter:
		baseActor = &a.Actor
	case *model.NPC:
		baseActor = &a.Actor
	case *model.Enemy:
		baseActor = &a.Actor
	case *model.Companion:
		baseActor = &a.Actor
	}

	if baseActor.Exhaustion <= 0 {
		return nil, fmt.Errorf("角色当前没有力竭等级")
	}

	// 使用 rules 包函数移除力竭
	newLevel := rules.RemoveExhaustion(baseActor.Exhaustion, req.Levels)
	baseActor.Exhaustion = newLevel

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	message := fmt.Sprintf("%s 力竭等级降低至 %d", baseActor.Name, baseActor.Exhaustion)
	if baseActor.Exhaustion == 0 {
		message += " - 已完全恢复"
	}

	return &RemoveExhaustionResult{
		ActorID:  req.ActorID,
		NewLevel: baseActor.Exhaustion,
		Message:  message,
	}, nil
}

// GetExhaustionStatus 获取角色的力竭状态
// 查询角色当前的力竭等级和所有累积效果。
// 参数:
//
//	ctx - 上下文
//	req - 获取力竭状态请求，包含游戏ID和角色ID
//
// 返回:
//
//	*GetExhaustionStatusResult - 力竭状态结果
//	error - 角色不存在或加载游戏失败时返回错误
func (e *Engine) GetExhaustionStatus(ctx context.Context, req GetExhaustionStatusRequest) (*GetExhaustionStatusResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	actor, ok := game.GetActor(req.ActorID)
	if !ok {
		return nil, ErrNotFound
	}

	var baseActor *model.Actor
	switch a := actor.(type) {
	case *model.PlayerCharacter:
		baseActor = &a.Actor
	case *model.NPC:
		baseActor = &a.Actor
	case *model.Enemy:
		baseActor = &a.Actor
	case *model.Companion:
		baseActor = &a.Actor
	}

	currentLevel := baseActor.Exhaustion
	isDead := currentLevel >= 6

	// 获取力竭效果描述
	effects := []string{}
	if currentLevel > 0 {
		for level := 1; level <= currentLevel; level++ {
			effect := rules.GetExhaustionEffect(level)
			desc := rules.GetExhaustionDescription(level)
			effects = append(effects, fmt.Sprintf("力竭%d级: %s (%s)", level, effect, desc))
		}
	}

	message := fmt.Sprintf("%s - 力竭等级: %d/6", baseActor.Name, currentLevel)
	if currentLevel == 0 {
		message += " - 无"
	} else if isDead {
		message += " - 已死亡"
	}

	return &GetExhaustionStatusResult{
		ActorID:      req.ActorID,
		CurrentLevel: currentLevel,
		Effects:      effects,
		IsDead:       isDead,
		Message:      message,
	}, nil
}
