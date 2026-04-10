package engine

import (
	"context"
	"fmt"

	"github.com/zwh8800/dnd-core/pkg/model"
	"github.com/zwh8800/dnd-core/pkg/rules"
)

// ============================================================================
// 信息查询API（只读，使用RLock）
// ============================================================================

// GetLifestyleRequest 获取生活方式信息请求
type GetLifestyleRequest struct {
	GameID    model.ID            `json:"game_id"`   // 游戏会话ID（必填）
	Lifestyle model.LifestyleTier `json:"lifestyle"` // 生活方式类型（必填）
}

// GetLifestyleResult 生活方式信息结果
type GetLifestyleResult struct {
	Lifestyle   model.LifestyleTier `json:"lifestyle"`    // 生活方式类型
	DailyCost   int                 `json:"daily_cost"`   // 每日花费（GP）
	MonthlyCost int                 `json:"monthly_cost"` // 每月花费（GP）
	Description string              `json:"description"`  // 描述
}

// GetLifestyleInfo 获取生活方式信息
// 返回指定生活方式的详细信息，包括每日花费、每月花费和描述
// 参数:
//
//	ctx - 上下文
//	req - 生活方式信息请求，包含游戏ID和生活方式类型
//
// 返回:
//
//	*GetLifestyleResult - 生活方式信息，包含每日花费、描述等
//	error - 游戏不存在时返回错误
func (e *Engine) GetLifestyleInfo(ctx context.Context, req GetLifestyleRequest) (*GetLifestyleResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	_, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	// 如果未指定，使用默认
	lifestyle := req.Lifestyle
	if lifestyle == "" {
		lifestyle = model.LifestyleModest
	}

	// 获取生活方式信息
	cost := rules.CalculateLifestyleCost(lifestyle, 1)
	description := rules.GetLifestyleDescription(lifestyle)

	return &GetLifestyleResult{
		Lifestyle:   lifestyle,
		DailyCost:   cost.DailyCost,
		MonthlyCost: cost.MonthlyCost,
		Description: description,
	}, nil
}

// GetCraftingInfoRequest 获取工艺制作信息请求
type GetCraftingInfoRequest struct {
	ItemName string `json:"item_name"` // 物品名称（必填）
}

// GetCraftingInfoResult 工艺制作信息结果
type GetCraftingInfoResult struct {
	ItemName       string `json:"item_name"`       // 物品名称
	CraftingDC     int    `json:"crafting_dc"`     // 制作DC
	CraftingTime   int    `json:"crafting_time"`   // 制作时间（工作日）
	CraftingCost   int    `json:"crafting_cost"`   // 制作成本（GP）
	ProficiencyReq string `json:"proficiency_req"` // 所需熟练项
	Description    string `json:"description"`     // 描述
}

// GetCraftingInfo 获取工艺制作信息
// 返回指定物品的制作DC、时间、成本和所需熟练项
// 参数:
//
//	ctx - 上下文
//	req - 工艺制作信息请求，包含物品名称
//
// 返回:
//
//	*GetCraftingInfoResult - 工艺制作信息，包含DC、时间、成本等
//	error - 物品不存在时返回错误
func (e *Engine) GetCraftingInfo(ctx context.Context, req GetCraftingInfoRequest) (*GetCraftingInfoResult, error) {
	// 注意：这是一个只读查询，不需要加载游戏状态
	// 但为了保持一致性，我们仍然使用引擎方法

	// 从数据注册中心获取物品定义
	// TODO: 实现物品查询逻辑
	// 目前返回一个通用的制作信息

	info := &GetCraftingInfoResult{
		ItemName:    req.ItemName,
		CraftingDC:  rules.DCMedium, // 默认DC 15
		Description: fmt.Sprintf("制作 %s 所需的信息", req.ItemName),
	}

	// TODO: 根据物品类型获取具体制作信息
	// 这里需要扩展物品数据模型以包含制作信息

	return info, nil
}

// GetCarryingCapacityRequest 获取负重能力请求
type GetCarryingCapacityRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID（必填）
	ActorID model.ID `json:"actor_id"` // 角色ID（必填）
}

// GetCarryingCapacityResult 负重能力结果
type GetCarryingCapacityResult struct {
	ActorID             model.ID `json:"actor_id"`              // 角色ID
	Strength            int      `json:"strength"`              // 力量值
	CarryingCapacity    int      `json:"carrying_capacity"`     // 负重能力（磅）
	PushDragLift        int      `json:"push_drag_lift"`        // 推/拖/举能力（磅）
	CurrentWeight       int      `json:"current_weight"`        // 当前负重（磅）
	IsEncumbered        bool     `json:"is_encumbered"`         // 是否超重
	IsHeavilyEncumbered bool     `json:"is_heavily_encumbered"` // 是否严重超重
}

// GetCarryingCapacity 获取角色负重能力
// 返回角色的负重能力、当前负重状态和是否超重
// 参数:
//
//	ctx - 上下文
//	req - 负重能力请求，包含游戏ID和角色ID
//
// 返回:
//
//	*GetCarryingCapacityResult - 负重能力信息，包含负重能力、当前状态等
//	error - 角色不存在时返回错误
func (e *Engine) GetCarryingCapacity(ctx context.Context, req GetCarryingCapacityRequest) (*GetCarryingCapacityResult, error) {
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

	// 计算负重能力
	carryingCapacity := rules.CalculateCarryingCapacity(baseActor.AbilityScores.Strength)
	pushDragLift := rules.CalculatePushDragLift(baseActor.AbilityScores.Strength)

	// 计算当前负重（需要遍历背包物品）
	currentWeight := 0
	// TODO: 实现物品重量累加逻辑

	// 判断是否超重（可选变体规则）
	isEncumbered := currentWeight > carryingCapacity/2
	isHeavilyEncumbered := currentWeight > carryingCapacity

	return &GetCarryingCapacityResult{
		ActorID:             req.ActorID,
		Strength:            baseActor.AbilityScores.Strength,
		CarryingCapacity:    carryingCapacity,
		PushDragLift:        pushDragLift,
		CurrentWeight:       currentWeight,
		IsEncumbered:        isEncumbered,
		IsHeavilyEncumbered: isHeavilyEncumbered,
	}, nil
}

// GetExhaustionEffectsRequest 获取力竭效果请求
type GetExhaustionEffectsRequest struct {
	ExhaustionLevel int `json:"exhaustion_level"` // 力竭等级（必填，0-6）
}

// GetExhaustionEffectsResult 力竭效果结果
type GetExhaustionEffectsResult struct {
	ExhaustionLevel int      `json:"exhaustion_level"` // 力竭等级
	Effects         []string `json:"effects"`          // 效果列表
	IsDead          bool     `json:"is_dead"`          // 是否死亡（6级力竭）
	Description     string   `json:"description"`      // 详细描述
}

// GetExhaustionEffects 获取力竭效果描述
// 返回指定力竭等级的所有累积效果
// 参数:
//
//	ctx - 上下文
//	req - 力竭效果请求，包含力竭等级
//
// 返回:
//
//	*GetExhaustionEffectsResult - 力竭效果信息，包含效果列表和描述
//	error - 力竭等级无效时返回错误
func (e *Engine) GetExhaustionEffects(ctx context.Context, req GetExhaustionEffectsRequest) (*GetExhaustionEffectsResult, error) {
	if req.ExhaustionLevel < 0 || req.ExhaustionLevel > 6 {
		return nil, fmt.Errorf("invalid exhaustion level: %d (must be 0-6)", req.ExhaustionLevel)
	}

	// 获取累积效果
	effects := make([]string, 0, req.ExhaustionLevel)
	for i := 1; i <= req.ExhaustionLevel; i++ {
		if effect, ok := rules.ExhaustionEffects[i]; ok {
			effects = append(effects, effect)
		}
	}

	isDead := req.ExhaustionLevel >= 6

	description := fmt.Sprintf("力竭等级 %d", req.ExhaustionLevel)
	if isDead {
		description += " - 死亡"
	} else if len(effects) > 0 {
		description += fmt.Sprintf(" - 效果: %s", rules.JoinStrings(effects, ", "))
	}

	return &GetExhaustionEffectsResult{
		ExhaustionLevel: req.ExhaustionLevel,
		Effects:         effects,
		IsDead:          isDead,
		Description:     description,
	}, nil
}
