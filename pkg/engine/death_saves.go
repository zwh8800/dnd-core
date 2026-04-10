package engine

import (
	"context"
	"fmt"

	"github.com/zwh8800/dnd-core/pkg/model"
	"github.com/zwh8800/dnd-core/pkg/rules"
)

// ============================================================================
// 死亡豁免API结构体定义
// ============================================================================

// PerformDeathSaveRequest 执行死亡豁免请求
// 当角色HP降至0且未立即死亡时，需要进行死亡豁免检定
type PerformDeathSaveRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID（必填）
	ActorID model.ID `json:"actor_id"` // 角色ID（必填）
}

// PerformDeathSaveResult 死亡豁免结果
// 描述死亡豁免检定的完整结果和当前状态
type PerformDeathSaveResult struct {
	Roll            int    `json:"roll"`              // d20掷骰值
	Success         bool   `json:"success"`           // 本次检定是否成功
	IsCritical      bool   `json:"is_critical"`       // 是否天然20（立即恢复1HP）
	IsCriticalFail  bool   `json:"is_critical_fail"`  // 是否天然1（2次失败）
	TotalSuccesses  int    `json:"total_successes"`   // 累计成功次数
	TotalFailures   int    `json:"total_failures"`    // 累计失败次数
	IsStable        bool   `json:"is_stable"`         // 是否已稳定
	IsDead          bool   `json:"is_dead"`           // 是否死亡
	AutoStabilizeHP int    `json:"auto_stabilize_hp"` // 自动恢复的HP（仅天然20）
	Message         string `json:"message"`           // 人类可读消息
}

// StabilizeCreatureRequest 稳定生物请求
// 用于对0HP的濒死生物进行急救，使其不再进行死亡豁免
type StabilizeCreatureRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID（必填）
	ActorID model.ID `json:"actor_id"` // 目标角色ID（必填）
}

// StabilizeCreatureResult 稳定生物结果
type StabilizeCreatureResult struct {
	ActorID model.ID `json:"actor_id"` // 角色ID
	Message string   `json:"message"`  // 人类可读消息
}

// GetDeathSaveStatusRequest 获取死亡豁免状态请求
type GetDeathSaveStatusRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID（必填）
	ActorID model.ID `json:"actor_id"` // 角色ID（必填）
}

// GetDeathSaveStatusResult 死亡豁免状态结果
type GetDeathSaveStatusResult struct {
	ActorID       model.ID `json:"actor_id"`       // 角色ID
	IsUnconscious bool     `json:"is_unconscious"` // 是否 unconscious
	IsStable      bool     `json:"is_stable"`      // 是否已稳定
	IsDead        bool     `json:"is_dead"`        // 是否死亡
	Successes     int      `json:"successes"`      // 成功次数
	Failures      int      `json:"failures"`       // 失败次数
	CurrentHP     int      `json:"current_hp"`     // 当前HP
	Message       string   `json:"message"`        // 人类可读消息
}

// ============================================================================
// 死亡豁免API
// ============================================================================

// PerformDeathSave 执行死亡豁免检定
// 当角色HP为0时，每回合结束时必须进行死亡豁免检定（DC 10）。
// 规则：
// - 掷骰 >= 10：成功
// - 掷骰 < 10：失败
// - 天然20：立即恢复1HP，脱离濒死状态
// - 天然1：计为2次失败
// - 累计3次成功：生物稳定（不再检定，但仍 unconscious）
// - 累计3次失败：生物死亡
// 参数:
//
//	ctx - 上下文
//	req - 死亡豁免请求，包含游戏ID和角色ID
//
// 返回:
//
//	*PerformDeathSaveResult - 死亡豁免结果，包含掷骰、状态和消息
//	error - 角色不存在、HP不为0或保存失败时返回错误
func (e *Engine) PerformDeathSave(ctx context.Context, req PerformDeathSaveRequest) (*PerformDeathSaveResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpPerformDeathSave); err != nil {
		return nil, err
	}

	actor, ok := game.GetActor(req.ActorID)
	if !ok {
		return nil, ErrNotFound
	}

	var pc *model.PlayerCharacter
	var baseActor *model.Actor
	switch a := actor.(type) {
	case *model.PlayerCharacter:
		pc = a
		baseActor = &a.Actor
	default:
		// NPC/敌人/同伴在0HP时直接死亡，不需要死亡豁免
		return nil, fmt.Errorf("only player characters make death saves")
	}

	// 检查是否在0HP
	if baseActor.HitPoints.Current > 0 {
		return nil, fmt.Errorf("角色HP大于0，不需要进行死亡豁免")
	}

	// 检查是否已稳定
	if baseActor.HasCondition(model.ConditionStabilized) {
		return nil, fmt.Errorf("角色已稳定，不需要进行死亡豁免")
	}

	// 执行死亡豁免检定
	deathSaveResult, err := rules.MakeDeathSave()
	if err != nil {
		return nil, fmt.Errorf("死亡豁免掷骰失败: %w", err)
	}

	// 更新计数器
	if deathSaveResult.IsCritical {
		// 天然20：立即恢复1HP
		pc.DeathSaveSuccesses += 2 // 算2次成功
		baseActor.HitPoints.Current = 1
		// 移除稳定状态（如果有的话）
		newConditions := make([]model.ConditionInstance, 0)
		for _, c := range baseActor.Conditions {
			if c.Type != model.ConditionStabilized {
				newConditions = append(newConditions, c)
			}
		}
		baseActor.Conditions = newConditions
	} else if deathSaveResult.IsCriticalFail {
		// 天然1：2次失败
		pc.DeathSaveFailures += 2
	} else if deathSaveResult.Success {
		pc.DeathSaveSuccesses++
	} else {
		pc.DeathSaveFailures++
	}

	// 检查死亡/稳定状态
	isStable, isDead, statusMessage := rules.CheckDeathStatus(pc.DeathSaveSuccesses, pc.DeathSaveFailures)

	if isStable {
		// 添加稳定状态
		baseActor.Conditions = append(baseActor.Conditions, model.ConditionInstance{
			Type: model.ConditionStabilized,
		})
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	message := fmt.Sprintf("死亡豁免掷骰: %d", deathSaveResult.Roll)
	if deathSaveResult.IsCritical {
		message += " (天然20！恢复1HP)"
	} else if deathSaveResult.IsCriticalFail {
		message += " (天然1！2次失败)"
	} else if deathSaveResult.Success {
		message += " (成功)"
	} else {
		message += " (失败)"
	}
	message += fmt.Sprintf(" [成功: %d/3, 失败: %d/3]", pc.DeathSaveSuccesses, pc.DeathSaveFailures)
	message += fmt.Sprintf(" - %s", statusMessage)

	return &PerformDeathSaveResult{
		Roll:            deathSaveResult.Roll,
		Success:         deathSaveResult.Success,
		IsCritical:      deathSaveResult.IsCritical,
		IsCriticalFail:  deathSaveResult.IsCriticalFail,
		TotalSuccesses:  pc.DeathSaveSuccesses,
		TotalFailures:   pc.DeathSaveFailures,
		IsStable:        isStable,
		IsDead:          isDead,
		AutoStabilizeHP: baseActor.HitPoints.Current,
		Message:         message,
	}, nil
}

// StabilizeCreature 稳定濒死生物
// 对HP为0的生物进行急救（如使用医疗技能），使其不再进行死亡豁免检定。
// 生物仍然 unconscious，但不会继续恶化。
// 参数:
//
//	ctx - 上下文
//	req - 稳定生物请求，包含游戏ID和目标角色ID
//
// 返回:
//
//	*StabilizeCreatureResult - 稳定结果
//	error - 角色不存在、HP不为0或保存失败时返回错误
func (e *Engine) StabilizeCreature(ctx context.Context, req StabilizeCreatureRequest) (*StabilizeCreatureResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpStabilizeCreature); err != nil {
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

	// 检查是否在0HP
	if baseActor.HitPoints.Current > 0 {
		return nil, fmt.Errorf("角色HP大于0，不需要稳定")
	}

	// 检查是否已稳定
	if baseActor.HasCondition(model.ConditionStabilized) {
		return nil, fmt.Errorf("角色已稳定")
	}

	// 应用稳定效果
	message := rules.StabilizeCreature()

	// 添加稳定状态
	baseActor.Conditions = append(baseActor.Conditions, model.ConditionInstance{
		Type: model.ConditionStabilized,
	})

	// 重置死亡豁免计数（已稳定后不再需要）
	if pc, ok := actor.(*model.PlayerCharacter); ok {
		pc.DeathSaveSuccesses = 3 // 设为3表示已稳定
		pc.DeathSaveFailures = 0
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &StabilizeCreatureResult{
		ActorID: req.ActorID,
		Message: message,
	}, nil
}

// GetDeathSaveStatus 获取死亡豁免状态
// 查询角色当前的死亡豁免状态，包括成功/失败计数和稳定状态。
// 参数:
//
//	ctx - 上下文
//	req - 获取死亡豁免状态请求，包含游戏ID和角色ID
//
// 返回:
//
//	*GetDeathSaveStatusResult - 死亡豁免状态结果
//	error - 角色不存在或加载游戏失败时返回错误
func (e *Engine) GetDeathSaveStatus(ctx context.Context, req GetDeathSaveStatusRequest) (*GetDeathSaveStatusResult, error) {
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

	var pc *model.PlayerCharacter
	var baseActor *model.Actor
	switch a := actor.(type) {
	case *model.PlayerCharacter:
		pc = a
		baseActor = &a.Actor
	default:
		return nil, fmt.Errorf("only player characters have death save status")
	}

	isUnconscious := baseActor.HitPoints.Current <= 0
	isStable := baseActor.HasCondition(model.ConditionStabilized)
	isDead := baseActor.IsDead()

	message := fmt.Sprintf("%s - HP: %d/%d", baseActor.Name, baseActor.HitPoints.Current, baseActor.HitPoints.Maximum)
	if isStable {
		message += " [已稳定]"
	} else if isUnconscious {
		message += fmt.Sprintf(" [死亡豁免: 成功 %d/3, 失败 %d/3]", pc.DeathSaveSuccesses, pc.DeathSaveFailures)
	}

	return &GetDeathSaveStatusResult{
		ActorID:       req.ActorID,
		IsUnconscious: isUnconscious,
		IsStable:      isStable,
		IsDead:        isDead,
		Successes:     pc.DeathSaveSuccesses,
		Failures:      pc.DeathSaveFailures,
		CurrentHP:     baseActor.HitPoints.Current,
		Message:       message,
	}, nil
}
