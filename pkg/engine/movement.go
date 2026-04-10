package engine

import (
	"context"
	"fmt"

	"github.com/zwh8800/dnd-core/pkg/dice"
	"github.com/zwh8800/dnd-core/pkg/model"
	"github.com/zwh8800/dnd-core/pkg/rules"
)

// getBaseActor 从任意 actor 类型中提取基础 Actor
func getBaseActor(actor any) (*model.Actor, error) {
	switch a := actor.(type) {
	case *model.PlayerCharacter:
		return &a.Actor, nil
	case *model.NPC:
		return &a.Actor, nil
	case *model.Enemy:
		return &a.Actor, nil
	case *model.Companion:
		return &a.Actor, nil
	default:
		return nil, fmt.Errorf("unknown actor type")
	}
}

// PerformJumpRequest 执行跳跃请求
type PerformJumpRequest struct {
	GameID          model.ID       `json:"game_id"`           // 游戏会话ID（必填）
	ActorID         model.ID       `json:"actor_id"`          // 跳跃者ID（必填）
	JumpType        model.JumpType `json:"jump_type"`         // 跳跃类型：long（跳远）或 high（跳高）（必填）
	HasRunningStart bool           `json:"has_running_start"` // 是否有助跑（可选，默认 false）
}

// PerformJumpResult 执行跳跃结果
type PerformJumpResult struct {
	JumpType        model.JumpType `json:"jump_type"`         // 跳跃类型
	Distance        int            `json:"distance"`          // 跳跃距离（尺）
	HasRunningStart bool           `json:"has_running_start"` // 是否有助跑
	Strength        int            `json:"strength"`          // 力量值
	StrengthMod     int            `json:"strength_mod"`      // 力量修正
	Message         string         `json:"message"`           // 描述消息
}

// PerformJump 执行跳跃动作
// 根据 D&D 5e 规则计算跳跃距离：
//   - 跳远：有助跑时 = 力量值尺数，立定跳远 = 力量值/2
//   - 跳高：有助跑时 = 3+力量修正尺数，立定跳高 = 一半
//
// 参数:
//
//	ctx - 上下文
//	req - 执行跳跃请求，包含游戏ID、角色ID、跳跃类型和是否助跑
//
// 返回:
//
//	*PerformJumpResult - 跳跃结果，包含跳跃距离和详细信息
//	error - 错误信息，如游戏加载失败、角色不存在等
func (e *Engine) PerformJump(ctx context.Context, req PerformJumpRequest) (*PerformJumpResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpPerformJump); err != nil {
		return nil, err
	}

	actorAny, ok := game.GetActor(req.ActorID)
	if !ok {
		return nil, fmt.Errorf("actor %s not found", req.ActorID)
	}

	baseActor, err := getBaseActor(actorAny)
	if err != nil {
		return nil, err
	}

	strength := baseActor.AbilityScores.Strength
	distance := rules.CalculateJumpDistance(strength, req.JumpType, req.HasRunningStart)

	strMod := rules.AbilityModifier(strength)

	jumpTypeDesc := "跳远"
	if req.JumpType == model.JumpTypeHigh {
		jumpTypeDesc = "跳高"
	}
	runningDesc := ""
	if req.HasRunningStart {
		runningDesc = "（有助跑）"
	}

	result := &PerformJumpResult{
		JumpType:        req.JumpType,
		Distance:        distance,
		HasRunningStart: req.HasRunningStart,
		Strength:        strength,
		StrengthMod:     strMod,
		Message: fmt.Sprintf("%s %s%s 跳跃了 %d 尺（力量：%d，修正：%d）",
			baseActor.Name, jumpTypeDesc, runningDesc, distance, strength, strMod),
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return result, nil
}

// ApplyFallDamageRequest 应用跌落伤害请求
type ApplyFallDamageRequest struct {
	GameID       model.ID `json:"game_id"`       // 游戏会话ID（必填）
	ActorID      model.ID `json:"actor_id"`      // 受击者ID（必填）
	FallDistance int      `json:"fall_distance"` // 跌落距离（尺）（必填）
}

// ApplyFallDamageResult 应用跌落伤害结果
type ApplyFallDamageResult struct {
	FallDistance int    `json:"fall_distance"` // 跌落距离（尺）
	DamageDice   int    `json:"damage_dice"`   // 伤害骰数（如 3d6 中的 3）
	DamageTaken  int    `json:"damage_taken"`  // 实际受到的伤害
	MaxPossible  int    `json:"max_possible"`  // 最大可能伤害（20d6 = 120）
	CurrentHP    int    `json:"current_hp"`    // 当前生命值
	Message      string `json:"message"`       // 描述消息
}

// ApplyFallDamage 应用跌落伤害
// 根据 D&D 5e 规则：每 10 尺跌落造成 1d6 伤害，最多 20d6
//
// 参数:
//
//	ctx - 上下文
//	req - 应用跌落伤害请求，包含游戏ID、角色ID和跌落距离
//
// 返回:
//
//	*ApplyFallDamageResult - 跌落伤害结果，包含伤害详细信息
//	error - 错误信息，如游戏加载失败、角色不存在等
func (e *Engine) ApplyFallDamage(ctx context.Context, req ApplyFallDamageRequest) (*ApplyFallDamageResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpApplyFallDamage); err != nil {
		return nil, err
	}

	actorAny, ok := game.GetActor(req.ActorID)
	if !ok {
		return nil, fmt.Errorf("actor %s not found", req.ActorID)
	}

	baseActor, err := getBaseActor(actorAny)
	if err != nil {
		return nil, err
	}

	_, diceCount, maxDamage := rules.CalculateFallDamage(req.FallDistance)

	if diceCount == 0 {
		return &ApplyFallDamageResult{
			FallDistance: req.FallDistance,
			DamageDice:   0,
			DamageTaken:  0,
			MaxPossible:  maxDamage,
			CurrentHP:    baseActor.HitPoints.Current,
			Message:      fmt.Sprintf("%s 从 %d 尺高处跌落，但高度不足以造成伤害", baseActor.Name, req.FallDistance),
		}, nil
	}

	// 使用骰子掷骰计算实际伤害
	roller := dice.New(0)
	diceNotation := fmt.Sprintf("%dd6", diceCount)
	rollResult, err := roller.Roll(diceNotation)
	if err != nil {
		return nil, fmt.Errorf("failed to roll damage dice: %w", err)
	}

	actualDamage := rollResult.Total
	if actualDamage > maxDamage {
		actualDamage = maxDamage
	}

	// 应用伤害
	baseActor.HitPoints.Current -= actualDamage
	if baseActor.HitPoints.Current < 0 {
		baseActor.HitPoints.Current = 0
	}

	result := &ApplyFallDamageResult{
		FallDistance: req.FallDistance,
		DamageDice:   diceCount,
		DamageTaken:  actualDamage,
		MaxPossible:  maxDamage,
		CurrentHP:    baseActor.HitPoints.Current,
		Message: fmt.Sprintf("%s 从 %d 尺高处跌落，受到 %d 点伤害（%s = %d），当前 HP：%d",
			baseActor.Name, req.FallDistance, actualDamage, diceNotation, actualDamage, baseActor.HitPoints.Current),
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return result, nil
}

// CalculateBreathHoldingRequest 计算闭气时间请求
type CalculateBreathHoldingRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID（必填）
	ActorID model.ID `json:"actor_id"` // 角色ID（必填）
}

// CalculateBreathHoldingResult 计算闭气时间结果
type CalculateBreathHoldingResult struct {
	Constitution           int    `json:"constitution"`             // 体质值
	ConstitutionMod        int    `json:"constitution_mod"`         // 体质修正
	CanHoldBreathSecs      int    `json:"can_hold_breath_secs"`     // 可闭气时间（秒）
	RoundsUntilUnconscious int    `json:"rounds_until_unconscious"` // 失去意识前的轮数
	Message                string `json:"message"`                  // 描述消息
}

// CalculateBreathHolding 计算角色的闭气能力
// 根据 D&D 5e 规则：
//   - 可闭气时间 = 1 + 体质修正（分钟），最少 30 秒
//   - 窒息后可存活轮数 = 体质修正（最少 1 轮）
//
// 参数:
//
//	ctx - 上下文
//	req - 计算闭气时间请求，包含游戏ID和角色ID
//
// 返回:
//
//	*CalculateBreathHoldingResult - 闭气能力结果，包含详细信息
//	error - 错误信息，如游戏加载失败、角色不存在等
func (e *Engine) CalculateBreathHolding(ctx context.Context, req CalculateBreathHoldingRequest) (*CalculateBreathHoldingResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpCalculateBreathHolding); err != nil {
		return nil, err
	}

	actorAny, ok := game.GetActor(req.ActorID)
	if !ok {
		return nil, fmt.Errorf("actor %s not found", req.ActorID)
	}

	baseActor, err := getBaseActor(actorAny)
	if err != nil {
		return nil, err
	}

	constitution := baseActor.AbilityScores.Constitution
	holdBreathTime := rules.CalculateHoldBreathTime(constitution)
	suffocationRounds := rules.CalculateSuffocationRounds(constitution)
	conMod := rules.AbilityModifier(constitution)

	result := &CalculateBreathHoldingResult{
		Constitution:           constitution,
		ConstitutionMod:        conMod,
		CanHoldBreathSecs:      holdBreathTime,
		RoundsUntilUnconscious: suffocationRounds,
		Message: fmt.Sprintf("%s 可以闭气 %d 秒（%d 分钟），窒息后还能存活 %d 轮（体质：%d，修正：%d）",
			baseActor.Name, holdBreathTime, holdBreathTime/60, suffocationRounds, constitution, conMod),
	}

	return result, nil
}

// ApplySuffocationRequest 应用窒息伤害请求
type ApplySuffocationRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID（必填）
	ActorID model.ID `json:"actor_id"` // 受窒息影响的角色ID（必填）
}

// ApplySuffocationResult 应用窒息伤害结果
type ApplySuffocationResult struct {
	Constitution      int    `json:"constitution"`       // 体质值
	SuffocationRounds int    `json:"suffocation_rounds"` // 窒息后可存活轮数
	IsUnconscious     bool   `json:"is_unconscious"`     // 是否已失去意识
	CurrentHP         int    `json:"current_hp"`         // 当前生命值
	Message           string `json:"message"`            // 描述消息
}

// ApplySuffocation 应用窒息效果
// 当生物无法呼吸时，在体质修正数量的轮数后失去意识（最少1轮）
// 失去意识后，如果仍然无法呼吸，会在后续回合中死亡
//
// 参数:
//
//	ctx - 上下文
//	req - 应用窒息伤害请求，包含游戏ID和角色ID
//
// 返回:
//
//	*ApplySuffocationResult - 窒息结果，包含详细信息
//	error - 错误信息，如游戏加载失败、角色不存在等
func (e *Engine) ApplySuffocation(ctx context.Context, req ApplySuffocationRequest) (*ApplySuffocationResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpApplySuffocation); err != nil {
		return nil, err
	}

	actorAny, ok := game.GetActor(req.ActorID)
	if !ok {
		return nil, fmt.Errorf("actor %s not found", req.ActorID)
	}

	baseActor, err := getBaseActor(actorAny)
	if err != nil {
		return nil, err
	}

	constitution := baseActor.AbilityScores.Constitution
	suffocationRounds := rules.CalculateSuffocationRounds(constitution)

	// 检查是否已经窒息
	isUnconscious := baseActor.HasCondition(model.ConditionUnconscious)

	// 如果已经失去意识，窒息会导致进一步恶化
	// 这里简化处理：每次调用表示又过了一轮窒息
	if isUnconscious {
		// 窒息中的生物每轮受到 1d6 伤害（简化处理）
		roller := dice.New(0)
		rollResult, err := roller.Roll("1d6")
		if err != nil {
			return nil, fmt.Errorf("failed to roll suffocation damage: %w", err)
		}

		baseActor.HitPoints.Current -= rollResult.Total
		if baseActor.HitPoints.Current < 0 {
			baseActor.HitPoints.Current = 0
		}

		result := &ApplySuffocationResult{
			Constitution:      constitution,
			SuffocationRounds: suffocationRounds,
			IsUnconscious:     true,
			CurrentHP:         baseActor.HitPoints.Current,
			Message: fmt.Sprintf("%s 窒息中，受到 %d 点伤害，当前 HP：%d",
				baseActor.Name, rollResult.Total, baseActor.HitPoints.Current),
		}

		if err := e.saveGame(ctx, game); err != nil {
			return nil, err
		}

		return result, nil
	}

	result := &ApplySuffocationResult{
		Constitution:      constitution,
		SuffocationRounds: suffocationRounds,
		IsUnconscious:     false,
		CurrentHP:         baseActor.HitPoints.Current,
		Message: fmt.Sprintf("%s 开始窒息，还能坚持 %d 轮后失去意识（体质：%d）",
			baseActor.Name, suffocationRounds, constitution),
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return result, nil
}

// PerformEncounterCheckRequest 执行遭遇检定请求
type PerformEncounterCheckRequest struct {
	GameID model.ID `json:"game_id"` // 游戏会话ID（必填）
}

// PerformEncounterCheckResult 执行遭遇检定结果
type PerformEncounterCheckResult struct {
	Encountered   bool   `json:"encountered"`    // 是否遭遇
	EncounterType string `json:"encounter_type"` // 遭遇类型：monster/npc/treasure/trap
	Roll          int    `json:"roll"`           // 掷骰结果
	Message       string `json:"message"`        // 描述消息
}

// PerformEncounterCheck 执行随机遭遇检定
// 根据 D&D 5e 规则：掷 1d6，结果为 1-2 时发生遭遇
// 遭遇类型随机决定：怪物、NPC、宝藏或陷阱
//
// 参数:
//
//	ctx - 上下文
//	req - 执行遭遇检定请求，包含游戏ID
//
// 返回:
//
//	*PerformEncounterCheckResult - 遭遇检定结果
//	error - 错误信息，如游戏加载失败等
func (e *Engine) PerformEncounterCheck(ctx context.Context, req PerformEncounterCheckRequest) (*PerformEncounterCheckResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpPerformEncounterCheck); err != nil {
		return nil, err
	}

	encounterResult, err := rules.EncounterCheck()
	if err != nil {
		return nil, fmt.Errorf("failed to perform encounter check: %w", err)
	}

	result := &PerformEncounterCheckResult{
		Encountered:   encounterResult.Encountered,
		EncounterType: encounterResult.EncounterType,
		Roll:          encounterResult.Roll,
		Message:       encounterResult.Message,
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return result, nil
}
