package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/zwh8800/dnd-core/pkg/data"
	"github.com/zwh8800/dnd-core/pkg/model"
	"github.com/zwh8800/dnd-core/pkg/rules"
)

// StartCraftingRequest 开始制作请求
type StartCraftingRequest struct {
	GameID   model.ID `json:"game_id"`
	ActorID  model.ID `json:"actor_id"`
	RecipeID string   `json:"recipe_id"`
}

// StartCraftingResult 开始制作结果
type StartCraftingResult struct {
	Progress *model.CraftingProgress `json:"progress"`
	Message  string                  `json:"message"`
}

// AdvanceCraftingRequest 推进制作请求
type AdvanceCraftingRequest struct {
	GameID  model.ID `json:"game_id"`
	ActorID model.ID `json:"actor_id"`
	Days    int      `json:"days"`
}

// AdvanceCraftingResult 推进制作结果
type AdvanceCraftingResult struct {
	Progress   *model.CraftingProgress `json:"progress"`
	IsComplete bool                    `json:"is_complete"`
	Message    string                  `json:"message"`
}

// CompleteCraftingRequest 完成制作请求
type CompleteCraftingRequest struct {
	GameID   model.ID `json:"game_id"`
	ActorID  model.ID `json:"actor_id"`
	RecipeID string   `json:"recipe_id"`
}

// CompleteCraftingResult 完成制作结果
type CompleteCraftingResult struct {
	Success bool   `json:"success"`
	ItemID  string `json:"item_id,omitempty"`
	Message string `json:"message"`
}

// StartCrafting 开始制作
func (e *Engine) StartCrafting(ctx context.Context, req StartCraftingRequest) (*StartCraftingResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	// 获取角色
	actor, ok := game.GetActor(req.ActorID)
	if !ok {
		return nil, ErrNotFound
	}

	pc, ok := actor.(*model.PlayerCharacter)
	if !ok {
		return nil, fmt.Errorf("只有玩家角色可以制作物品")
	}

	// 获取配方
	recipe, exists := data.GetCraftingRecipe(req.RecipeID)
	if !exists {
		return nil, fmt.Errorf("配方不存在: %s", req.RecipeID)
	}

	// 验证角色等级
	if pc.TotalLevel < recipe.MinLevel {
		return nil, fmt.Errorf("等级不足: 需要%d级,当前%d级", recipe.MinLevel, pc.TotalLevel)
	}

	// 验证工具熟练
	if len(recipe.ToolsRequired) > 0 {
		hasProficiency := false
		for _, tool := range recipe.ToolsRequired {
			if pc.Proficiencies.ToolProficiencies != nil && pc.Proficiencies.ToolProficiencies[tool] {
				hasProficiency = true
				break
			}
		}
		if !hasProficiency {
			return nil, fmt.Errorf("缺少所需的工具熟练: %v", recipe.ToolsRequired)
		}
	}

	// 验证法术能力
	if recipe.SpellRequired != "" {
		if pc.Spellcasting == nil {
			return nil, fmt.Errorf("需要能施放法术: %s", recipe.SpellRequired)
		}
		// 检查是否已知或准备了该法术
		hasSpell := false
		for _, spell := range pc.Spellcasting.KnownSpells {
			if spell == recipe.SpellRequired {
				hasSpell = true
				break
			}
		}
		if !hasSpell && pc.Spellcasting.PreparedSpells != nil {
			for _, spell := range pc.Spellcasting.PreparedSpells {
				if spell == recipe.SpellRequired {
					hasSpell = true
					break
				}
			}
		}
		if !hasSpell {
			return nil, fmt.Errorf("需要能施放法术: %s", recipe.SpellRequired)
		}
	}

	// 验证金币是否足够
	totalCost := recipe.Cost
	if pc.Gold < totalCost {
		return nil, fmt.Errorf("金币不足: 需要%dgp, 当前%dgp", totalCost, pc.Gold)
	}

	// 计算制作时间
	hasProficiency := false
	if len(recipe.ToolsRequired) > 0 && pc.Proficiencies.ToolProficiencies != nil {
		for _, tool := range recipe.ToolsRequired {
			if pc.Proficiencies.ToolProficiencies[tool] {
				hasProficiency = true
				break
			}
		}
	}
	totalDays := rules.CalculateCraftingTime(recipe, hasProficiency)

	// 扣除金币
	pc.Gold -= totalCost

	// 创建制作进度
	progress := &model.CraftingProgress{
		RecipeID:    req.RecipeID,
		DaysWorked:  0,
		TotalDays:   totalDays,
		MoneySpent:  totalCost,
		LastWorkDay: time.Now().Format("2006-01-02"),
		IsComplete:  false,
	}

	// 存储进度到角色的制作列表中
	if pc.CraftingProgress == nil {
		pc.CraftingProgress = make(map[string]*model.CraftingProgress)
	}
	pc.CraftingProgress[req.RecipeID] = progress

	result := &StartCraftingResult{
		Progress: progress,
		Message:  fmt.Sprintf("开始制作: %s (需要%d天, 已扣除%dgp)", recipe.Name, totalDays, totalCost),
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return result, nil
}

// AdvanceCrafting 推进制作进度
func (e *Engine) AdvanceCrafting(ctx context.Context, req AdvanceCraftingRequest) (*AdvanceCraftingResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	// 获取角色
	actor, ok := game.GetActor(req.ActorID)
	if !ok {
		return nil, ErrNotFound
	}

	pc, ok := actor.(*model.PlayerCharacter)
	if !ok {
		return nil, fmt.Errorf("只有玩家角色可以制作物品")
	}

	// 查找正在进行的制作
	if pc.CraftingProgress == nil {
		return nil, fmt.Errorf("没有正在进行的制作")
	}

	// 找到第一个未完成的制作
	var progress *model.CraftingProgress
	var recipeID string
	for id, p := range pc.CraftingProgress {
		if !p.IsComplete {
			progress = p
			recipeID = id
			break
		}
	}

	if progress == nil {
		return nil, fmt.Errorf("没有正在进行的制作")
	}

	// 更新进度
	progress.DaysWorked += req.Days
	progress.LastWorkDay = time.Now().Format("2006-01-02")

	// 检查是否完成
	isComplete := progress.DaysWorked >= progress.TotalDays
	if isComplete {
		progress.IsComplete = true
	}

	recipe, _ := data.GetCraftingRecipe(recipeID)

	result := &AdvanceCraftingResult{
		Progress:   progress,
		IsComplete: isComplete,
		Message:    fmt.Sprintf("制作进度: %s %d/%d天", recipe.Name, progress.DaysWorked, progress.TotalDays),
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return result, nil
}

// CompleteCrafting 完成制作
func (e *Engine) CompleteCrafting(ctx context.Context, req CompleteCraftingRequest) (*CompleteCraftingResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	// 获取角色
	actor, ok := game.GetActor(req.ActorID)
	if !ok {
		return nil, ErrNotFound
	}

	pc, ok := actor.(*model.PlayerCharacter)
	if !ok {
		return nil, fmt.Errorf("只有玩家角色可以制作物品")
	}

	// 查找制作进度
	if pc.CraftingProgress == nil {
		return nil, fmt.Errorf("没有正在进行的制作")
	}

	progress, exists := pc.CraftingProgress[req.RecipeID]
	if !exists {
		return nil, fmt.Errorf("未找到该制作进度: %s", req.RecipeID)
	}

	if !progress.IsComplete {
		return nil, fmt.Errorf("制作尚未完成: %d/%d天", progress.DaysWorked, progress.TotalDays)
	}

	// 获取配方
	recipe, exists := data.GetCraftingRecipe(req.RecipeID)
	if !exists {
		return nil, fmt.Errorf("配方不存在: %s", req.RecipeID)
	}

	// 创建物品
	item := model.Item{
		ID:          model.ID(recipe.ID),
		Name:        recipe.Name,
		Description: recipe.Description,
		Type:        model.ItemTypeWondrousItem,
		Weight:      0,
		Value:       recipe.Cost,
		Quantity:    1,
	}

	// 根据类型设置物品属性
	switch recipe.Type {
	case model.CraftingTypePotion:
		item.Consumable = true
		item.Effect = "恢复生命值"
	case model.CraftingTypeScroll:
		item.Consumable = true
		item.Effect = "施放法术"
	case model.CraftingTypeMagicItem:
		item.Rarity = model.RarityUncommon
		item.MagicEffects = []string{"魔法物品效果"}
	}

	// 添加到角色库存
	if pc.InventoryID == "" {
		invID := model.NewID()
		pc.InventoryID = invID
		game.Inventories[invID] = &model.Inventory{
			ID:    invID,
			Items: []*model.Item{&item},
		}
	} else {
		if inv, exists := game.Inventories[pc.InventoryID]; exists {
			inv.Items = append(inv.Items, &item)
		}
	}

	// 移除制作进度
	delete(pc.CraftingProgress, req.RecipeID)

	result := &CompleteCraftingResult{
		Success: true,
		ItemID:  recipe.ID,
		Message: fmt.Sprintf("制作完成: 获得 %s", recipe.Name),
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return result, nil
}

// GetCraftingRecipes 获取所有可用配方
func (e *Engine) GetCraftingRecipes(ctx context.Context) ([]CraftingRecipeInfo, error) {
	recipes := data.GetAllCraftingRecipes()
	result := make([]CraftingRecipeInfo, 0, len(recipes))

	for _, recipe := range recipes {
		result = append(result, CraftingRecipeInfo{
			ID:          recipe.ID,
			Name:        recipe.Name,
			Type:        string(recipe.Type),
			Description: recipe.Description,
			TimeDays:    recipe.TimeDays,
			DC:          recipe.DC,
			MinLevel:    recipe.MinLevel,
			Cost:        recipe.Cost,
		})
	}

	return result, nil
}

// CraftingRecipeInfo 配方信息
type CraftingRecipeInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	TimeDays    int    `json:"time_days"`
	DC          int    `json:"dc"`
	MinLevel    int    `json:"min_level"`
	Cost        int    `json:"cost"`
}
