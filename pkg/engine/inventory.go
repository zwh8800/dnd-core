package engine

import (
	"context"
	"fmt"

	"github.com/zwh8800/dnd-core/pkg/model"
	"github.com/zwh8800/dnd-core/pkg/rules"
)

// ItemInput 物品创建输入（替代 *model.Item 入参）
type ItemInput struct {
	Name               string         `json:"name"`                // 物品名称
	Description        string         `json:"description"`         // 物品描述
	Type               model.ItemType `json:"type"`                // 物品类型
	Rarity             model.Rarity   `json:"rarity"`              // 稀有度
	Weight             float64        `json:"weight"`              // 重量（磅）
	Quantity           int            `json:"quantity"`            // 数量
	Value              int            `json:"value"`               // 价值（铜币）
	Attunement         string         `json:"attunement"`          // 调谐要求描述
	RequiresAttunement bool           `json:"requires_attunement"` // 是否需要调谐
	ArmorClass         int            `json:"armor_class"`         // 护甲等级（防具）
	AttackBonus        int            `json:"attack_bonus"`        // 攻击加值
	Damage             string         `json:"damage"`              // 伤害描述
	Range              string         `json:"range"`               // 射程
}

// ItemSummary 物品摘要信息（替代 *model.Item 返回值）
type ItemSummary struct {
	ID          model.ID       `json:"id"`          // 物品唯一标识
	Name        string         `json:"name"`        // 物品名称
	Description string         `json:"description"` // 物品描述
	Type        model.ItemType `json:"type"`        // 物品类型
	Rarity      model.Rarity   `json:"rarity"`      // 稀有度
	Weight      float64        `json:"weight"`      // 重量（磅）
	Quantity    int            `json:"quantity"`    // 数量
	Value       int            `json:"value"`       // 价值（铜币）
	Attuned     bool           `json:"attuned"`     // 是否已调谐
	MagicBonus  int            `json:"magic_bonus"` // 魔法加值
}

// AddItemRequest 添加物品请求
type AddItemRequest struct {
	GameID  model.ID   `json:"game_id"`  // 游戏会话ID
	ActorID model.ID   `json:"actor_id"` // 角色ID
	Item    *ItemInput `json:"item"`     // 物品信息
}

// RemoveItemRequest 移除物品请求
type RemoveItemRequest struct {
	GameID   model.ID `json:"game_id"`  // 游戏会话ID
	ActorID  model.ID `json:"actor_id"` // 角色ID
	ItemID   model.ID `json:"item_id"`  // 物品ID
	Quantity int      `json:"quantity"` // 移除数量
}

// GetInventoryRequest 获取库存请求
type GetInventoryRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID
	ActorID model.ID `json:"actor_id"` // 角色ID
}

// EquipItemRequest 装备物品请求
type EquipItemRequest struct {
	GameID  model.ID            `json:"game_id"`  // 游戏会话ID
	ActorID model.ID            `json:"actor_id"` // 角色ID
	ItemID  model.ID            `json:"item_id"`  // 物品ID
	Slot    model.EquipmentSlot `json:"slot"`     // 装备槽位
}

// UnequipItemRequest 卸下装备请求
type UnequipItemRequest struct {
	GameID  model.ID            `json:"game_id"`  // 游戏会话ID
	ActorID model.ID            `json:"actor_id"` // 角色ID
	Slot    model.EquipmentSlot `json:"slot"`     // 装备槽位
}

// GetEquipmentRequest 获取装备请求
type GetEquipmentRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID
	ActorID model.ID `json:"actor_id"` // 角色ID
}

// AttuneItemRequest 调谐物品请求
type AttuneItemRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID
	ActorID model.ID `json:"actor_id"` // 角色ID
	ItemID  model.ID `json:"item_id"`  // 物品ID
}

// TransferItemRequest 转移物品请求
type TransferItemRequest struct {
	GameID      model.ID `json:"game_id"`       // 游戏会话ID
	FromActorID model.ID `json:"from_actor_id"` // 源角色ID
	ToActorID   model.ID `json:"to_actor_id"`   // 目标角色ID
	ItemID      model.ID `json:"item_id"`       // 物品ID
	Quantity    int      `json:"quantity"`      // 转移数量
}

// AddCurrencyRequest 添加货币请求
type AddCurrencyRequest struct {
	GameID   model.ID       `json:"game_id"`  // 游戏会话ID
	ActorID  model.ID       `json:"actor_id"` // 角色ID
	Currency model.Currency `json:"currency"` // 货币信息
}

// InventoryResult 库存操作结果
type InventoryResult struct {
	Success     bool     `json:"success"`
	ItemID      model.ID `json:"item_id,omitempty"`
	Quantity    int      `json:"quantity"`
	TotalWeight float64  `json:"total_weight"`
	Message     string   `json:"message"`
}

// EquipResult 装备操作结果
type EquipResult struct {
	Success      bool                `json:"success"`
	ItemName     string              `json:"item_name"`
	Slot         model.EquipmentSlot `json:"slot"`
	PreviousItem *model.Item         `json:"previous_item,omitempty"`
	ACChanged    bool                `json:"ac_changed"`
	NewAC        int                 `json:"new_ac,omitempty"`
	Message      string              `json:"message"`
}

// InventoryInfo 库存信息
type InventoryInfo struct {
	OwnerID     model.ID       `json:"owner_id"`
	Items       []*model.Item  `json:"items"`
	TotalWeight float64        `json:"total_weight"`
	MaxWeight   float64        `json:"max_weight"`
	Encumbered  bool           `json:"encumbered"`
	Currency    model.Currency `json:"currency"`
}

// EquipmentInfo 装备信息
type EquipmentInfo struct {
	OwnerID       model.ID                            `json:"owner_id"`
	EquippedSlots map[model.EquipmentSlot]*model.Item `json:"equipped_slots"`
	TotalACBonus  int                                 `json:"total_ac_bonus"`
	MagicBonuses  map[string]int                      `json:"magic_bonuses"`
}

// AttuneResult 调谐结果
type AttuneResult struct {
	Success       bool   `json:"success"`
	ItemName      string `json:"item_name"`
	Attuned       bool   `json:"attuned"`
	AttunedCount  int    `json:"attuned_count"`
	MaxAttunement int    `json:"max_attunement"`
	Message       string `json:"message"`
}

// TransferResult 物品转移结果
type TransferResult struct {
	Success   bool     `json:"success"`
	ItemID    model.ID `json:"item_id"`
	FromActor model.ID `json:"from_actor"`
	ToActor   model.ID `json:"to_actor"`
	Quantity  int      `json:"quantity"`
	Message   string   `json:"message"`
}

// getInventoryHelper 辅助函数：获取角色库存
func getInventoryHelper(game *model.GameState, actorID model.ID) *model.Inventory {
	if inv, ok := game.Inventories[actorID]; ok {
		return inv
	}
	return nil
}

// itemToSummary 将 model.Item 转换为 ItemSummary
func itemToSummary(item *model.Item) *ItemSummary {
	if item == nil {
		return nil
	}
	return &ItemSummary{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		Type:        item.Type,
		Rarity:      item.Rarity,
		Weight:      item.Weight,
		Quantity:    item.Quantity,
		Value:       item.Value,
		Attuned:     item.Attuned,
		MagicBonus:  item.MagicBonus,
	}
}

// itemInputToItem 将 ItemInput 转换为 model.Item
func itemInputToItem(input *ItemInput) *model.Item {
	if input == nil {
		return nil
	}

	item := &model.Item{
		ID:          model.NewID(),
		Name:        input.Name,
		Description: input.Description,
		Type:        input.Type,
		Rarity:      input.Rarity,
		Weight:      input.Weight,
		Quantity:    input.Quantity,
		Value:       input.Value,
		Attunement:  input.Attunement,
		Attuned:     false, // 新添加的物品默认未调谐
		MagicBonus:  0,     // 默认无魔法加值
	}

	// 如果需要调谐但未提供描述，使用默认值
	if input.RequiresAttunement && item.Attunement == "" {
		item.Attunement = "需要调谐"
	}

	// 设置武器属性
	if input.Type == model.ItemTypeWeapon {
		item.WeaponProps = &model.WeaponProperties{}
		// 魔法加值设置到物品级别
		if input.AttackBonus != 0 {
			item.MagicBonus = input.AttackBonus
		}
	}

	// 设置护甲属性
	if input.Type == model.ItemTypeArmor && input.ArmorClass != 0 {
		item.ArmorProps = &model.ArmorProperties{
			BaseAC: input.ArmorClass,
		}
	}

	return item
}

// AddItem 添加物品到角色库存
func (e *Engine) AddItem(ctx context.Context, req AddItemRequest) (*InventoryResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

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

	// 获取或创建库存
	inventory := getInventoryHelper(game, baseActor.ID)
	if inventory == nil {
		inventory = model.NewInventory(baseActor.ID)
		game.Inventories[baseActor.ID] = inventory
	}

	// 将 ItemInput 转换为 model.Item
	item := itemInputToItem(req.Item)

	// 检查负重
	newWeight := calculateTotalWeight(inventory) + item.Weight*float64(item.Quantity)
	maxWeight := inventory.MaxWeight
	if maxWeight == 0 {
		maxWeight = calculateMaxWeight(baseActor)
		inventory.MaxWeight = maxWeight
	}

	if newWeight > maxWeight {
		return &InventoryResult{
			Success: false,
			Message: "物品太重，无法携带",
		}, nil
	}

	// 添加物品
	item.ID = model.NewID()
	inventory.Items = append(inventory.Items, item)

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &InventoryResult{
		Success:     true,
		ItemID:      item.ID,
		Quantity:    item.Quantity,
		TotalWeight: newWeight,
		Message:     fmt.Sprintf("添加了 %s x%d", item.Name, item.Quantity),
	}, nil
}

// RemoveItem 从角色库存移除物品
func (e *Engine) RemoveItem(ctx context.Context, req RemoveItemRequest) (*InventoryResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

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

	inventory := getInventoryHelper(game, baseActor.ID)
	if inventory == nil {
		return nil, ErrNotFound
	}

	// 查找物品
	for i, item := range inventory.Items {
		if item.ID == req.ItemID {
			if item.Quantity <= req.Quantity {
				// 移除整个物品
				inventory.Items = append(inventory.Items[:i], inventory.Items[i+1:]...)
				req.Quantity = item.Quantity
			} else {
				// 减少数量
				item.Quantity -= req.Quantity
			}

			if err := e.saveGame(ctx, game); err != nil {
				return nil, err
			}

			return &InventoryResult{
				Success:     true,
				ItemID:      req.ItemID,
				Quantity:    req.Quantity,
				TotalWeight: calculateTotalWeight(inventory),
				Message:     fmt.Sprintf("移除了 %s x%d", item.Name, req.Quantity),
			}, nil
		}
	}

	return nil, ErrNotFound
}

// GetInventory 获取角色库存信息
func (e *Engine) GetInventory(ctx context.Context, req GetInventoryRequest) (*InventoryInfo, error) {
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

	inventory := getInventoryHelper(game, baseActor.ID)
	if inventory == nil {
		return &InventoryInfo{
			OwnerID:     req.ActorID,
			Items:       make([]*model.Item, 0),
			TotalWeight: 0,
			MaxWeight:   calculateMaxWeight(baseActor),
			Encumbered:  false,
			Currency:    model.Currency{},
		}, nil
	}

	totalWeight := calculateTotalWeight(inventory)
	maxWeight := inventory.MaxWeight
	if maxWeight == 0 {
		maxWeight = calculateMaxWeight(baseActor)
	}

	return &InventoryInfo{
		OwnerID:     req.ActorID,
		Items:       inventory.Items,
		TotalWeight: totalWeight,
		MaxWeight:   maxWeight,
		Encumbered:  totalWeight > maxWeight*0.5,
		Currency:    inventory.Currency,
	}, nil
}

// EquipItem 装备物品
func (e *Engine) EquipItem(ctx context.Context, req EquipItemRequest) (*EquipResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

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

	inventory := getInventoryHelper(game, baseActor.ID)
	if inventory == nil {
		return nil, ErrNotFound
	}

	// 查找物品
	var itemToEquip *model.Item
	for _, item := range inventory.Items {
		if item.ID == req.ItemID {
			itemToEquip = item
			break
		}
	}

	if itemToEquip == nil {
		return nil, ErrNotFound
	}

	// 验证物品可以装备到该槽位
	if !canEquipToSlot(itemToEquip, req.Slot) {
		return &EquipResult{
			Success: false,
			Message: fmt.Sprintf("%s 不能装备到 %s 槽位", itemToEquip.Name, req.Slot),
		}, nil
	}

	// 获取当前装备
	previousItem := inventory.Equipment.Slots[req.Slot]

	// 装备新物品
	inventory.Equipment.Slots[req.Slot] = itemToEquip

	// 计算AC变化
	oldAC := baseActor.ArmorClass
	newAC := calculateAC(baseActor, inventory.Equipment)
	baseActor.ArmorClass = newAC

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	result := &EquipResult{
		Success:      true,
		ItemName:     itemToEquip.Name,
		Slot:         req.Slot,
		PreviousItem: previousItem,
		ACChanged:    oldAC != newAC,
		NewAC:        newAC,
		Message:      fmt.Sprintf("装备了 %s 到 %s", itemToEquip.Name, req.Slot),
	}

	if previousItem != nil {
		result.Message += fmt.Sprintf(" (替换了 %s)", previousItem.Name)
	}

	return result, nil
}

// UnequipItem 卸下装备
func (e *Engine) UnequipItem(ctx context.Context, req UnequipItemRequest) (*EquipResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

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

	inventory := getInventoryHelper(game, baseActor.ID)
	if inventory == nil {
		return nil, ErrNotFound
	}

	// 获取当前装备
	item := inventory.Equipment.Slots[req.Slot]
	if item == nil {
		return &EquipResult{
			Success: false,
			Message: fmt.Sprintf("%s 槽位没有装备", req.Slot),
		}, nil
	}

	// 卸下装备
	delete(inventory.Equipment.Slots, req.Slot)

	// 重新计算AC
	newAC := calculateAC(baseActor, inventory.Equipment)
	baseActor.ArmorClass = newAC

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &EquipResult{
		Success:   true,
		ItemName:  item.Name,
		Slot:      req.Slot,
		ACChanged: true,
		NewAC:     newAC,
		Message:   fmt.Sprintf("卸下了 %s", item.Name),
	}, nil
}

// GetEquipment 获取角色装备信息
func (e *Engine) GetEquipment(ctx context.Context, req GetEquipmentRequest) (*EquipmentInfo, error) {
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

	inventory := getInventoryHelper(game, baseActor.ID)
	if inventory == nil {
		return &EquipmentInfo{
			OwnerID:       req.ActorID,
			EquippedSlots: make(map[model.EquipmentSlot]*model.Item),
			TotalACBonus:  0,
			MagicBonuses:  make(map[string]int),
		}, nil
	}

	// 计算魔法加值
	magicBonuses := make(map[string]int)
	for _, item := range inventory.Equipment.Slots {
		if item != nil && item.Attuned {
			magicBonuses[item.Name] = item.MagicBonus
		}
	}

	return &EquipmentInfo{
		OwnerID:       req.ActorID,
		EquippedSlots: inventory.Equipment.Slots,
		TotalACBonus:  calculateACBonus(inventory.Equipment),
		MagicBonuses:  magicBonuses,
	}, nil
}

// AttuneItem 调谐魔法物品
func (e *Engine) AttuneItem(ctx context.Context, req AttuneItemRequest) (*AttuneResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

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

	inventory := getInventoryHelper(game, baseActor.ID)
	if inventory == nil {
		return nil, ErrNotFound
	}

	// 查找物品
	var itemToAttune *model.Item
	for _, item := range inventory.Items {
		if item.ID == req.ItemID {
			itemToAttune = item
			break
		}
	}

	if itemToAttune == nil {
		return nil, ErrNotFound
	}

	// 检查是否需要调谐
	if itemToAttune.Attunement == "" {
		return &AttuneResult{
			Success: false,
			Message: fmt.Sprintf("%s 不需要调谐", itemToAttune.Name),
		}, nil
	}

	// 计算当前调谐数量
	attunedCount := 0
	for _, item := range inventory.Items {
		if item.Attuned {
			attunedCount++
		}
	}

	// 最大调谐数量为3
	maxAttunement := 3

	if !itemToAttune.Attuned && attunedCount >= maxAttunement {
		return &AttuneResult{
			Success:       false,
			ItemName:      itemToAttune.Name,
			AttunedCount:  attunedCount,
			MaxAttunement: maxAttunement,
			Message:       "已达到最大调谐数量（3个）",
		}, nil
	}

	// 切换调谐状态
	itemToAttune.Attuned = !itemToAttune.Attuned
	if itemToAttune.Attuned {
		attunedCount++
	} else {
		attunedCount--
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	message := fmt.Sprintf("与 %s 的调谐已", itemToAttune.Name)
	if itemToAttune.Attuned {
		message += "建立"
	} else {
		message += "解除"
	}

	return &AttuneResult{
		Success:       true,
		ItemName:      itemToAttune.Name,
		Attuned:       itemToAttune.Attuned,
		AttunedCount:  attunedCount,
		MaxAttunement: maxAttunement,
		Message:       message,
	}, nil
}

// TransferItem 转移物品给另一个角色
func (e *Engine) TransferItem(ctx context.Context, req TransferItemRequest) (*TransferResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	// 获取源角色
	fromActor, ok := game.GetActor(req.FromActorID)
	if !ok {
		return nil, ErrNotFound
	}

	// 获取目标角色
	toActor, ok := game.GetActor(req.ToActorID)
	if !ok {
		return nil, ErrNotFound
	}

	var fromBaseActor *model.Actor
	switch a := fromActor.(type) {
	case *model.PlayerCharacter:
		fromBaseActor = &a.Actor
	case *model.NPC:
		fromBaseActor = &a.Actor
	case *model.Enemy:
		fromBaseActor = &a.Actor
	case *model.Companion:
		fromBaseActor = &a.Actor
	}

	var toBaseActor *model.Actor
	switch a := toActor.(type) {
	case *model.PlayerCharacter:
		toBaseActor = &a.Actor
	case *model.NPC:
		toBaseActor = &a.Actor
	case *model.Enemy:
		toBaseActor = &a.Actor
	case *model.Companion:
		toBaseActor = &a.Actor
	}

	fromInventory := getInventoryHelper(game, fromBaseActor.ID)
	if fromInventory == nil {
		return nil, ErrNotFound
	}

	toInventory := getInventoryHelper(game, toBaseActor.ID)
	if toInventory == nil {
		toInventory = model.NewInventory(toBaseActor.ID)
		game.Inventories[toBaseActor.ID] = toInventory
	}

	// 查找物品
	var itemToTransfer *model.Item
	var itemIndex int
	for i, item := range fromInventory.Items {
		if item.ID == req.ItemID {
			itemToTransfer = item
			itemIndex = i
			break
		}
	}

	if itemToTransfer == nil {
		return nil, ErrNotFound
	}

	// 检查数量
	if itemToTransfer.Quantity < req.Quantity {
		return &TransferResult{
			Success: false,
			Message: "物品数量不足",
		}, nil
	}

	// 创建转移的物品副本
	transferredItem := &model.Item{
		ID:          model.NewID(),
		Name:        itemToTransfer.Name,
		Description: itemToTransfer.Description,
		Type:        itemToTransfer.Type,
		Rarity:      itemToTransfer.Rarity,
		Weight:      itemToTransfer.Weight,
		Quantity:    req.Quantity,
		Value:       itemToTransfer.Value,
		Attuned:     false, // 转移后调谐解除
		Attunement:  itemToTransfer.Attunement,
		WeaponProps: itemToTransfer.WeaponProps,
		ArmorProps:  itemToTransfer.ArmorProps,
		MagicBonus:  itemToTransfer.MagicBonus,
	}

	// 添加到目标库存
	toInventory.Items = append(toInventory.Items, transferredItem)

	// 从源库存移除
	if itemToTransfer.Quantity == req.Quantity {
		// 移除整个物品
		fromInventory.Items = append(fromInventory.Items[:itemIndex], fromInventory.Items[itemIndex+1:]...)
	} else {
		// 减少数量
		itemToTransfer.Quantity -= req.Quantity
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &TransferResult{
		Success:   true,
		ItemID:    transferredItem.ID,
		FromActor: req.FromActorID,
		ToActor:   req.ToActorID,
		Quantity:  req.Quantity,
		Message:   fmt.Sprintf("转移了 %s x%d", transferredItem.Name, req.Quantity),
	}, nil
}

// AddCurrency 添加货币
func (e *Engine) AddCurrency(ctx context.Context, req AddCurrencyRequest) (*InventoryResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

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

	inventory := getInventoryHelper(game, baseActor.ID)
	if inventory == nil {
		inventory = model.NewInventory(baseActor.ID)
		game.Inventories[baseActor.ID] = inventory
	}

	// 添加货币
	inventory.Currency.Platinum += req.Currency.Platinum
	inventory.Currency.Gold += req.Currency.Gold
	inventory.Currency.Electrum += req.Currency.Electrum
	inventory.Currency.Silver += req.Currency.Silver
	inventory.Currency.Copper += req.Currency.Copper

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	totalGold := inventory.Currency.TotalInGold()

	return &InventoryResult{
		Success: true,
		Message: fmt.Sprintf("添加了货币。当前总价值: %.2f gp", totalGold),
	}, nil
}

// 辅助函数

// calculateTotalWeight 计算库存总重量
func calculateTotalWeight(inventory *model.Inventory) float64 {
	total := 0.0
	for _, item := range inventory.Items {
		total += item.Weight * float64(item.Quantity)
	}
	// 货币重量（简化：每50枚硬币=1磅）
	coinCount := inventory.Currency.Platinum + inventory.Currency.Gold +
		inventory.Currency.Electrum + inventory.Currency.Silver + inventory.Currency.Copper
	total += float64(coinCount) / 50.0
	return total
}

// calculateMaxWeight 计算最大负重
func calculateMaxWeight(actor *model.Actor) float64 {
	// 基础负重 = 力量值 × 15（单位：磅）
	strength := actor.AbilityScores.Strength
	return float64(strength) * 15.0
}

// canEquipToSlot 检查物品是否可以装备到指定槽位
func canEquipToSlot(item *model.Item, slot model.EquipmentSlot) bool {
	switch item.Type {
	case model.ItemTypeWeapon:
		return slot == model.SlotMainHand || slot == model.SlotOffHand
	case model.ItemTypeArmor:
		return slot == model.SlotChest
	case model.ItemTypeRing:
		return slot == model.SlotFinger1 || slot == model.SlotFinger2
	default:
		return false
	}
}

// calculateAC 计算角色的AC
func calculateAC(actor *model.Actor, equipment *model.Equipment) int {
	// 基础AC = 10 + 敏捷修正
	dexMod := rules.AbilityModifier(actor.AbilityScores.Dexterity)
	baseAC := 10 + dexMod

	// 如果装备了护甲，使用护甲的AC
	if armor := equipment.Slots[model.SlotChest]; armor != nil && armor.ArmorProps != nil {
		baseAC = armor.ArmorProps.BaseAC
		// 应用敏捷修正（受最大敏捷修正限制）
		if armor.ArmorProps.MaxDexModifier != nil {
			if dexMod > *armor.ArmorProps.MaxDexModifier {
				dexMod = *armor.ArmorProps.MaxDexModifier
			}
		}
		baseAC += dexMod
	}

	// 盾牌加值
	if shield := equipment.Slots[model.SlotOffHand]; shield != nil && shield.ArmorProps != nil {
		baseAC += 2 // 标准盾牌加值
	}

	// 魔法加值
	for _, item := range equipment.Slots {
		if item != nil && item.Attuned {
			baseAC += item.MagicBonus
		}
	}

	return baseAC
}

// calculateACBonus 计算装备提供的AC加值
func calculateACBonus(equipment *model.Equipment) int {
	bonus := 0
	for _, item := range equipment.Slots {
		if item != nil && item.ArmorProps != nil {
			bonus += item.ArmorProps.BaseAC - 10 // 相对于基础AC的加值
		}
		if item != nil && item.Attuned {
			bonus += item.MagicBonus
		}
	}
	return bonus
}
