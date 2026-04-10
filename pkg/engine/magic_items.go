package engine

import (
	"context"
	"fmt"

	"github.com/zwh8800/dnd-core/pkg/model"
	"github.com/zwh8800/dnd-core/pkg/rules"
)

// UseMagicItemRequest 使用魔法物品请求
type UseMagicItemRequest struct {
	GameID    model.ID   `json:"game_id"`              // 游戏会话ID（必填）
	ActorID   model.ID   `json:"actor_id"`             // 角色ID（必填）
	ItemID    model.ID   `json:"item_id"`              // 物品ID（必填）
	TargetIDs []model.ID `json:"target_ids,omitempty"` // 目标ID列表（可选）
}

// UseMagicItemResult 使用魔法物品结果
type UseMagicItemResult struct {
	ItemName         string   `json:"item_name"`                   // 物品名称
	Messages         []string `json:"messages"`                    // 使用效果消息
	Consumed         bool     `json:"consumed"`                    // 物品是否被消耗
	ChargesRemaining int      `json:"charges_remaining,omitempty"` // 剩余充能（如果有）
}

// UnattuneItemRequest 取消调音物品请求
type UnattuneItemRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID（必填）
	ActorID model.ID `json:"actor_id"` // 角色ID（必填）
	ItemID  model.ID `json:"item_id"`  // 物品ID（必填）
}

// UnattuneItemResult 取消调音结果
type UnattuneItemResult struct {
	ItemName     string `json:"item_name"`     // 物品名称
	Success      bool   `json:"success"`       // 是否成功
	AttunedCount int    `json:"attuned_count"` // 当前调音物品数量
	Message      string `json:"message"`       // 人类可读消息
}

// RechargeMagicItemsRequest 充能魔法物品请求
type RechargeMagicItemsRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID（必填）
	ActorID model.ID `json:"actor_id"` // 角色ID（必填）
}

// RechargeMagicItemsResult 充能魔法物品结果
type RechargeMagicItemsResult struct {
	RechargedItems []string `json:"recharged_items"` // 已充能的物品名称列表
	Message        string   `json:"message"`         // 人类可读消息
}

// GetMagicItemBonusRequest 获取魔法物品加值请求
type GetMagicItemBonusRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID（必填）
	ActorID model.ID `json:"actor_id"` // 角色ID（必填）
	ItemID  model.ID `json:"item_id"`  // 物品ID（必填）
}

// GetMagicItemBonusResult 获取魔法物品加值结果
type GetMagicItemBonusResult struct {
	ItemName     string         `json:"item_name"`         // 物品名称
	MagicBonus   int            `json:"magic_bonus"`       // 魔法加值
	MagicEffects []string       `json:"magic_effects"`     // 魔法效果列表
	Attuned      bool           `json:"attuned"`           // 是否已调音
	Bonuses      map[string]int `json:"bonuses,omitempty"` // 各项加值明细
}

// UseMagicItem 使用魔法物品
// 根据物品类型执行不同效果：消耗品直接使用，充能物品消耗充能，被动物品提供持续效果。
// 使用后物品可能被消耗或减少充能次数。
// 参数:
//
//	ctx - 上下文
//	req - 使用请求，包含游戏会话ID、角色ID、物品ID和目标列表
//
// 返回:
//
//	*UseMagicItemResult - 使用结果，包含效果消息和物品状态
//	error - 可能返回 ErrNotFound（角色或物品不存在）、游戏不存在、权限错误或物品无法使用
func (e *Engine) UseMagicItem(ctx context.Context, req UseMagicItemRequest) (*UseMagicItemResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpUseMagicItem); err != nil {
		return nil, err
	}

	// 获取角色
	actor, ok := game.GetActor(req.ActorID)
	if !ok {
		return nil, ErrNotFound
	}

	var baseActor *model.Actor
	var pc *model.PlayerCharacter
	switch a := actor.(type) {
	case *model.PlayerCharacter:
		baseActor = &a.Actor
		pc = a
	case *model.NPC:
		baseActor = &a.Actor
	case *model.Enemy:
		baseActor = &a.Actor
	case *model.Companion:
		baseActor = &a.Actor
	}

	// 获取库存
	inventory := getInventoryHelper(game, baseActor.ID)
	if inventory == nil {
		return nil, ErrNotFound
	}

	// 查找物品
	var itemToUse *model.Item
	var itemIndex int
	found := false
	for i, item := range inventory.Items {
		if item.ID == req.ItemID {
			itemToUse = item
			itemIndex = i
			found = true
			break
		}
	}
	if !found || itemToUse == nil {
		return nil, ErrNotFound
	}

	// 检查是否需要调音
	if itemToUse.Attunement != "" && !itemToUse.Attuned {
		return nil, fmt.Errorf("物品 %s 需要调音后才能使用", itemToUse.Name)
	}

	// 创建物品副本用于使用（避免直接修改库存中的物品）
	itemCopy := *itemToUse

	// 调用规则层使用物品
	var userPC *model.PlayerCharacter
	if pc != nil {
		userPC = pc
	} else {
		// 非PC角色创建临时PC用于规则计算
		userPC = &model.PlayerCharacter{
			Actor: *baseActor,
		}
	}

	result, err := rules.UseMagicItem(&itemCopy, userPC, req.TargetIDs)
	if err != nil {
		return nil, fmt.Errorf("使用魔法物品失败: %w", err)
	}

	// 应用物品使用效果到库存中的物品
	if itemToUse.Consumable {
		// 消耗品：从库存中移除
		if itemToUse.Quantity > 1 {
			itemToUse.Quantity--
		} else {
			inventory.Items = append(inventory.Items[:itemIndex], inventory.Items[itemIndex+1:]...)
		}
	} else if itemToUse.Charges > 0 {
		// 充能物品：已在UseMagicItem中减少充能，同步回库存
		itemToUse.Charges = itemCopy.Charges
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &UseMagicItemResult{
		ItemName:         result.ItemName,
		Messages:         result.Messages,
		Consumed:         itemToUse.Consumable,
		ChargesRemaining: itemToUse.Charges,
	}, nil
}

// UnattuneItem 取消调音魔法物品
// 解除角色与指定魔法物品的调音关系，释放调音槽位。调音解除后物品的魔法效果不再生效。
// 参数:
//
//	ctx - 上下文
//	req - 取消调音请求，包含游戏会话ID、角色ID和物品ID
//
// 返回:
//
//	*UnattuneItemResult - 取消调音结果
//	error - 可能返回 ErrNotFound（角色或物品不存在）、游戏不存在、权限错误或物品未调音
func (e *Engine) UnattuneItem(ctx context.Context, req UnattuneItemRequest) (*UnattuneItemResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpUnattuneItem); err != nil {
		return nil, err
	}

	// 获取角色
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

	// 获取库存
	inventory := getInventoryHelper(game, baseActor.ID)
	if inventory == nil {
		return nil, ErrNotFound
	}

	// 查找物品
	var itemToUnattune *model.Item
	found := false
	for _, item := range inventory.Items {
		if item.ID == req.ItemID {
			itemToUnattune = item
			found = true
			break
		}
	}
	if !found || itemToUnattune == nil {
		return nil, ErrNotFound
	}

	// 调用规则层取消调音
	err = rules.UnattuneItem(itemToUnattune)
	if err != nil {
		return &UnattuneItemResult{
			ItemName: itemToUnattune.Name,
			Success:  false,
			Message:  err.Error(),
		}, nil
	}

	// 计算当前调音数量
	attunedCount := rules.GetAttunedItemCount(inventory)

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &UnattuneItemResult{
		ItemName:     itemToUnattune.Name,
		Success:      true,
		AttunedCount: attunedCount,
		Message:      fmt.Sprintf("已解除与 %s 的调音", itemToUnattune.Name),
	}, nil
}

// RechargeMagicItems 在黎明时恢复角色所有魔法物品的充能
// 恢复所有标记为"dawn"充能的物品到最大充能次数。通常在长休结束后自动调用。
// 参数:
//
//	ctx - 上下文
//	req - 充能请求，包含游戏会话ID和角色ID
//
// 返回:
//
//	*RechargeMagicItemsResult - 充能结果，包含已充能的物品列表
//	error - 可能返回 ErrNotFound（角色不存在）、游戏不存在或权限错误
func (e *Engine) RechargeMagicItems(ctx context.Context, req RechargeMagicItemsRequest) (*RechargeMagicItemsResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpRechargeMagicItems); err != nil {
		return nil, err
	}

	// 获取角色
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

	// 获取库存
	inventory := getInventoryHelper(game, baseActor.ID)
	if inventory == nil {
		return &RechargeMagicItemsResult{
			RechargedItems: nil,
			Message:        "角色没有库存",
		}, nil
	}

	// 收集需要充能的物品
	itemsToRecharge := make([]*model.Item, 0)
	for _, item := range inventory.Items {
		if item.MaxCharges > 0 && item.Recharge == "dawn" {
			itemsToRecharge = append(itemsToRecharge, item)
		}
	}

	if len(itemsToRecharge) == 0 {
		return &RechargeMagicItemsResult{
			RechargedItems: nil,
			Message:        "没有需要充能的魔法物品",
		}, nil
	}

	// 调用规则层充能
	messages := rules.RechargeMagicItems(itemsToRecharge)

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	message := "黎明时分，魔法物品已充能"
	if len(messages) > 0 {
		message += fmt.Sprintf("：%s", messages[0])
	}

	return &RechargeMagicItemsResult{
		RechargedItems: messages,
		Message:        message,
	}, nil
}

// GetMagicItemBonus 获取魔法物品的加值信息
// 返回指定魔法物品的魔法加值、效果列表以及调音状态。
// 参数:
//
//	ctx - 上下文
//	req - 获取请求，包含游戏会话ID、角色ID和物品ID
//
// 返回:
//
//	*GetMagicItemBonusResult - 魔法物品加值信息
//	error - 可能返回 ErrNotFound（角色或物品不存在）、游戏不存在或权限错误
func (e *Engine) GetMagicItemBonus(ctx context.Context, req GetMagicItemBonusRequest) (*GetMagicItemBonusResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpGetMagicItemBonus); err != nil {
		return nil, err
	}

	// 获取角色
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

	// 获取库存
	inventory := getInventoryHelper(game, baseActor.ID)
	if inventory == nil {
		return nil, ErrNotFound
	}

	// 查找物品
	var itemToCheck *model.Item
	found := false
	for _, item := range inventory.Items {
		if item.ID == req.ItemID {
			itemToCheck = item
			found = true
			break
		}
	}
	if !found || itemToCheck == nil {
		return nil, ErrNotFound
	}

	// 获取魔法加值
	magicBonus := rules.GetMagicItemBonus(itemToCheck)

	// 构建加值明细
	bonuses := make(map[string]int)
	if magicBonus != 0 {
		// 根据物品类型确定加值类型
		switch itemToCheck.Type {
		case model.ItemTypeWeapon:
			bonuses["攻击加值"] = magicBonus
			bonuses["伤害加值"] = magicBonus
		case model.ItemTypeArmor:
			bonuses["AC加值"] = magicBonus
		default:
			bonuses["魔法加值"] = magicBonus
		}
	}

	return &GetMagicItemBonusResult{
		ItemName:     itemToCheck.Name,
		MagicBonus:   magicBonus,
		MagicEffects: itemToCheck.MagicEffects,
		Attuned:      itemToCheck.Attuned,
		Bonuses:      bonuses,
	}, nil
}
