package engine

import (
	"context"
	"fmt"
	"strings"

	"github.com/zwh8800/dnd-core/pkg/data"
	"github.com/zwh8800/dnd-core/pkg/model"
)

// ============================================================================
// 可用动作计算 API 类型定义
// ============================================================================

// GetAvailableActionsRequest 获取可用动作请求
type GetAvailableActionsRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID
	ActorID model.ID `json:"actor_id"` // 角色ID
}

// AvailableActionsResult 当前角色可用动作的完整列表
type AvailableActionsResult struct {
	ActorID      model.ID          `json:"actor_id"`
	ActorName    string            `json:"actor_name"`
	ActorType    string            `json:"actor_type"` // "pc", "enemy", "npc", "companion"
	Movement     MovementOption    `json:"movement"`
	Actions      []AvailableAction `json:"actions"`
	BonusActions []AvailableAction `json:"bonus_actions"`
	Reactions    []AvailableAction `json:"reactions"`
	FreeActions  []AvailableAction `json:"free_actions"`
	Conditions   []string          `json:"conditions"`
}

// IsEmpty 检查是否没有任何可用动作（包括移动）
func (r *AvailableActionsResult) IsEmpty() bool {
	return len(r.Actions) == 0 &&
		len(r.BonusActions) == 0 &&
		len(r.FreeActions) == 0 &&
		(!r.Movement.Available || r.Movement.RemainingFeet <= 0)
}

// AvailableAction 单个可用动作
type AvailableAction struct {
	ID             string         `json:"id"`                         // 唯一标识，如 "attack_longsword", "cast_fireball"
	Category       string         `json:"category"`                   // "attack", "spell", "standard_action", "class_feature", "item_use"
	Name           string         `json:"name"`                       // 显示名称
	Description    string         `json:"description"`                // 规则描述
	CostType       string         `json:"cost_type"`                  // "action", "bonus_action", "reaction", "free_action"
	RequiresTarget bool           `json:"requires_target"`            // 是否需要目标
	TargetType     string         `json:"target_type"`                // "single", "area", "self", "none"
	ValidTargetIDs []model.ID     `json:"valid_target_ids,omitempty"` // 合法目标列表
	Range          string         `json:"range"`                      // "5尺触及", "120尺"
	ResourceCost   string         `json:"resource_cost,omitempty"`    // "消耗3环法术位"
	DamagePreview  string         `json:"damage_preview,omitempty"`   // "1d8+3 挥砍"
	Metadata       map[string]any `json:"metadata"`                   // 路由信息：_route, weapon_id, spell_id 等
}

// MovementOption 移动选项
type MovementOption struct {
	Available     bool `json:"available"`
	RemainingFeet int  `json:"remaining_feet"`
}

// ============================================================================
// 可用动作计算 — 公开 API
// ============================================================================

// OpGetAvailableActions 获取可用动作操作标识
const OpGetAvailableActions Operation = "get_available_actions"

// GetAvailableActions 获取指定角色当前可用的所有动作
// 根据角色类型、职业、装备、法术、状态效果等动态计算可用动作列表。
// 必须在战斗阶段且有活跃战斗时调用。
func (e *Engine) GetAvailableActions(ctx context.Context, req GetAvailableActionsRequest) (*AvailableActionsResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if game.Combat == nil || game.Combat.Status != model.CombatStatusActive {
		return nil, ErrCombatNotActive
	}

	return e.computeAvailableActions(game, req.ActorID)
}

// ============================================================================
// 可用动作计算 — 内部实现（不获取锁，由调用方负责）
// ============================================================================

// computeAvailableActions 计算指定角色的可用动作列表
// 此方法不获取锁，必须在已持有锁的上下文中调用
func (e *Engine) computeAvailableActions(game *model.GameState, actorID model.ID) (*AvailableActionsResult, error) {
	actorAny, ok := game.GetActor(actorID)
	if !ok {
		return nil, fmt.Errorf("actor %s not found", actorID)
	}

	// 提取基础 Actor 和类型信息
	var baseActor *model.Actor
	var actorName string
	var actorType string

	switch a := actorAny.(type) {
	case *model.PlayerCharacter:
		baseActor = &a.Actor
		actorName = a.Name
		actorType = string(model.ActorTypePC)
	case *model.Enemy:
		baseActor = &a.Actor
		actorName = a.Name
		actorType = string(model.ActorTypeEnemy)
	case *model.NPC:
		baseActor = &a.Actor
		actorName = a.Name
		actorType = string(model.ActorTypeNPC)
	case *model.Companion:
		baseActor = &a.Actor
		actorName = a.Name
		actorType = string(model.ActorTypeCompanion)
	default:
		return nil, fmt.Errorf("unknown actor type for %s", actorID)
	}

	result := &AvailableActionsResult{
		ActorID:      actorID,
		ActorName:    actorName,
		ActorType:    actorType,
		Actions:      make([]AvailableAction, 0),
		BonusActions: make([]AvailableAction, 0),
		Reactions:    make([]AvailableAction, 0),
		FreeActions:  make([]AvailableAction, 0),
		Conditions:   make([]string, 0),
	}

	// 收集当前状态效果
	conditionEffects := collectConditionEffects(baseActor)
	for _, c := range baseActor.Conditions {
		result.Conditions = append(result.Conditions, string(c.Type))
	}

	// Layer 1: 状态无力化检查 — 无法行动则返回空
	if conditionEffects.cantTakeActions {
		result.Movement = MovementOption{Available: false, RemainingFeet: 0}
		return result, nil
	}

	// 获取回合状态
	turnState := game.Combat.CurrentTurn
	if turnState == nil || turnState.ActorID != actorID {
		// 不是当前回合的角色，返回空动作
		result.Movement = MovementOption{Available: false, RemainingFeet: 0}
		return result, nil
	}

	// Layer 2: 计算移动
	movementLeft := baseActor.Speed - turnState.MovementUsed
	if conditionEffects.speedZero {
		movementLeft = 0
	}
	result.Movement = MovementOption{
		Available:     movementLeft > 0 && !conditionEffects.speedZero,
		RemainingFeet: movementLeft,
	}

	// 计算战场上的敌友目标
	enemies, allies := classifyTargets(game, actorID, actorType)

	// Layer 3-6: 根据角色类型计算具体可用动作
	switch a := actorAny.(type) {
	case *model.PlayerCharacter:
		e.computePCActions(game, a, turnState, conditionEffects, enemies, allies, result)
	case *model.Enemy:
		e.computeEnemyActions(game, a, turnState, conditionEffects, enemies, allies, result)
	case *model.NPC:
		e.computeNPCActions(a, turnState, conditionEffects, enemies, allies, result)
	case *model.Companion:
		e.computeCompanionActions(a, turnState, conditionEffects, enemies, allies, result)
	}

	return result, nil
}

// ============================================================================
// PC 可用动作计算
// ============================================================================

func (e *Engine) computePCActions(
	game *model.GameState,
	pc *model.PlayerCharacter,
	turn *model.TurnState,
	conditions conditionSummary,
	enemies, allies []model.ID,
	result *AvailableActionsResult,
) {
	actionAvailable := !turn.ActionUsed
	bonusActionAvailable := !turn.BonusActionUsed

	// === 标准动作（Action）===
	if actionAvailable {
		// 基础 D&D 战斗动作
		result.Actions = append(result.Actions, makeStandardActions(enemies)...)

		// 武器攻击
		weaponActions := e.computeWeaponAttacks(game, pc, enemies)
		result.Actions = append(result.Actions, weaponActions...)

		// 法术（action 施法时间）
		if pc.Spellcasting != nil {
			spellActions := computeSpellActions(pc.Spellcasting, conditions, enemies, allies, "action")
			result.Actions = append(result.Actions, spellActions...)
		}
	}

	// === 附赠动作（Bonus Action）===
	if bonusActionAvailable {
		// 副手攻击（双持轻武器）
		offHandAction := e.computeOffHandAttack(game, pc, enemies)
		if offHandAction != nil {
			result.BonusActions = append(result.BonusActions, *offHandAction)
		}

		// 法术（bonus action 施法时间）
		if pc.Spellcasting != nil {
			spellBonusActions := computeSpellActions(pc.Spellcasting, conditions, enemies, allies, "bonus_action")
			result.BonusActions = append(result.BonusActions, spellBonusActions...)
		}
	}

	// === 职业特性动作 ===
	if pc.FeatureHooks != nil {
		for _, hook := range pc.FeatureHooks {
			templates := hook.GetAvailableActions()
			for _, t := range templates {
				if t.CurrentUses <= 0 {
					continue
				}
				action := AvailableAction{
					ID:          fmt.Sprintf("class_%s", strings.ToLower(string(t.Type))),
					Category:    "class_feature",
					Name:        t.Name,
					Description: fmt.Sprintf("剩余 %d/%d 次", t.CurrentUses, t.UsesPerRest),
					Metadata:    map[string]any{"_route": "class_feature", "action_type": string(t.Type)},
				}

				switch {
				case t.IsFreeAction:
					action.CostType = "free_action"
					action.TargetType = "self"
					result.FreeActions = append(result.FreeActions, action)
				case t.IsBonusAction && bonusActionAvailable:
					action.CostType = "bonus_action"
					action.TargetType = "self"
					result.BonusActions = append(result.BonusActions, action)
				case !t.IsBonusAction && !t.IsFreeAction && actionAvailable:
					action.CostType = "action"
					action.TargetType = "self"
					result.Actions = append(result.Actions, action)
				}
			}
		}
	}

	// === 反应（Reaction）===
	if !turn.ReactionUsed {
		result.Reactions = append(result.Reactions, AvailableAction{
			ID:             "reaction_opportunity_attack",
			Category:       "attack",
			Name:           "借机攻击",
			Description:    "当敌人离开你的触及范围时可执行",
			CostType:       "reaction",
			RequiresTarget: true,
			TargetType:     "single",
			Range:          "5尺触及",
			Metadata:       map[string]any{"_route": "reaction", "reaction_type": "opportunity_attack"},
		})
	}
}

// ============================================================================
// Enemy 可用动作计算
// ============================================================================

func (e *Engine) computeEnemyActions(
	game *model.GameState,
	enemy *model.Enemy,
	turn *model.TurnState,
	conditions conditionSummary,
	enemies, allies []model.ID,
	result *AvailableActionsResult,
) {
	actionAvailable := !turn.ActionUsed
	bonusActionAvailable := !turn.BonusActionUsed

	if actionAvailable {
		// 基础 D&D 战斗动作
		result.Actions = append(result.Actions, makeStandardActions(enemies)...)

		// 敌人的武器攻击（基于 AttackBonus 和基础属性）
		if enemy.AttackBonus > 0 || enemy.DamagePerRound > 0 {
			result.Actions = append(result.Actions, AvailableAction{
				ID:             "attack_melee",
				Category:       "attack",
				Name:           "近战攻击",
				Description:    fmt.Sprintf("攻击加值 +%d", enemy.AttackBonus),
				CostType:       "action",
				RequiresTarget: true,
				TargetType:     "single",
				ValidTargetIDs: enemies,
				Range:          "5尺触及",
				DamagePreview:  fmt.Sprintf("约 %d 伤害", enemy.DamagePerRound),
				Metadata:       map[string]any{"_route": "attack", "is_unarmed": true},
			})
		}

		// 敌人装备的武器攻击
		weaponActions := e.computeEnemyWeaponAttacks(game, enemy, enemies)
		result.Actions = append(result.Actions, weaponActions...)
	}

	// 附赠动作
	if bonusActionAvailable {
		// 敌人一般没有复杂的附赠动作，但可扩展
	}

	// 反应
	if !turn.ReactionUsed {
		result.Reactions = append(result.Reactions, AvailableAction{
			ID:             "reaction_opportunity_attack",
			Category:       "attack",
			Name:           "借机攻击",
			Description:    "当敌人离开触及范围时可执行",
			CostType:       "reaction",
			RequiresTarget: true,
			TargetType:     "single",
			Range:          "5尺触及",
			Metadata:       map[string]any{"_route": "reaction", "reaction_type": "opportunity_attack"},
		})
	}
}

// ============================================================================
// NPC / Companion 可用动作计算
// ============================================================================

func (e *Engine) computeNPCActions(
	npc *model.NPC,
	turn *model.TurnState,
	conditions conditionSummary,
	enemies, allies []model.ID,
	result *AvailableActionsResult,
) {
	if !turn.ActionUsed {
		result.Actions = append(result.Actions, makeStandardActions(enemies)...)
		// NPC 基础攻击
		result.Actions = append(result.Actions, AvailableAction{
			ID:             "attack_melee",
			Category:       "attack",
			Name:           "近战攻击",
			Description:    "基础近战攻击",
			CostType:       "action",
			RequiresTarget: true,
			TargetType:     "single",
			ValidTargetIDs: enemies,
			Range:          "5尺触及",
			Metadata:       map[string]any{"_route": "attack", "is_unarmed": true},
		})
	}
}

func (e *Engine) computeCompanionActions(
	companion *model.Companion,
	turn *model.TurnState,
	conditions conditionSummary,
	enemies, allies []model.ID,
	result *AvailableActionsResult,
) {
	if !turn.ActionUsed {
		result.Actions = append(result.Actions, makeStandardActions(enemies)...)
		result.Actions = append(result.Actions, AvailableAction{
			ID:             "attack_melee",
			Category:       "attack",
			Name:           "近战攻击",
			Description:    "伙伴近战攻击",
			CostType:       "action",
			RequiresTarget: true,
			TargetType:     "single",
			ValidTargetIDs: enemies,
			Range:          "5尺触及",
			Metadata:       map[string]any{"_route": "attack", "is_unarmed": true},
		})
	}
}

// ============================================================================
// 武器攻击计算
// ============================================================================

func (e *Engine) computeWeaponAttacks(game *model.GameState, pc *model.PlayerCharacter, enemies []model.ID) []AvailableAction {
	actions := make([]AvailableAction, 0)

	inv := game.Inventories[pc.InventoryID]
	if inv == nil || inv.Equipment == nil {
		// 无装备，提供徒手攻击
		actions = append(actions, AvailableAction{
			ID:             "attack_unarmed",
			Category:       "attack",
			Name:           "徒手攻击",
			Description:    "1+力量修正 钝击伤害",
			CostType:       "action",
			RequiresTarget: true,
			TargetType:     "single",
			ValidTargetIDs: enemies,
			Range:          "5尺触及",
			DamagePreview:  "1+STR 钝击",
			Metadata:       map[string]any{"_route": "attack", "is_unarmed": true},
		})
		return actions
	}

	// 主手武器
	if mainHand, exists := inv.Equipment.Slots[model.SlotMainHand]; exists && mainHand != nil && mainHand.WeaponProps != nil {
		wp := mainHand.WeaponProps
		rangeStr := "5尺触及"
		if wp.Reach {
			rangeStr = "10尺触及"
		}
		if wp.WeaponType == "ranged" && wp.Range > 0 {
			rangeStr = fmt.Sprintf("%d/%d尺", wp.Range, wp.LongRange)
		}

		actions = append(actions, AvailableAction{
			ID:             fmt.Sprintf("attack_%s", strings.ReplaceAll(strings.ToLower(mainHand.Name), " ", "_")),
			Category:       "attack",
			Name:           mainHand.Name + "攻击",
			Description:    formatWeaponDescription(mainHand),
			CostType:       "action",
			RequiresTarget: true,
			TargetType:     "single",
			ValidTargetIDs: enemies,
			Range:          rangeStr,
			DamagePreview:  fmt.Sprintf("%s %s", wp.DamageDice, wp.DamageType),
			Metadata: map[string]any{
				"_route":    "attack",
				"weapon_id": string(mainHand.ID),
			},
		})
	}

	// 徒手攻击始终可用
	actions = append(actions, AvailableAction{
		ID:             "attack_unarmed",
		Category:       "attack",
		Name:           "徒手攻击",
		Description:    "1+力量修正 钝击伤害",
		CostType:       "action",
		RequiresTarget: true,
		TargetType:     "single",
		ValidTargetIDs: enemies,
		Range:          "5尺触及",
		DamagePreview:  "1+STR 钝击",
		Metadata:       map[string]any{"_route": "attack", "is_unarmed": true},
	})

	return actions
}

// computeOffHandAttack 计算副手攻击（双持轻武器时的附赠动作攻击）
func (e *Engine) computeOffHandAttack(game *model.GameState, pc *model.PlayerCharacter, enemies []model.ID) *AvailableAction {
	inv := game.Inventories[pc.InventoryID]
	if inv == nil || inv.Equipment == nil {
		return nil
	}

	mainHand, mainExists := inv.Equipment.Slots[model.SlotMainHand]
	offHand, offExists := inv.Equipment.Slots[model.SlotOffHand]

	// 双持条件：主手和副手都是轻武器
	if !mainExists || mainHand == nil || mainHand.WeaponProps == nil || !mainHand.WeaponProps.Light {
		return nil
	}
	if !offExists || offHand == nil || offHand.WeaponProps == nil || !offHand.WeaponProps.Light {
		return nil
	}

	return &AvailableAction{
		ID:             fmt.Sprintf("offhand_%s", strings.ReplaceAll(strings.ToLower(offHand.Name), " ", "_")),
		Category:       "attack",
		Name:           offHand.Name + "副手攻击",
		Description:    "双持武器附赠动作攻击（不加属性修正到伤害）",
		CostType:       "bonus_action",
		RequiresTarget: true,
		TargetType:     "single",
		ValidTargetIDs: enemies,
		Range:          "5尺触及",
		DamagePreview:  fmt.Sprintf("%s %s", offHand.WeaponProps.DamageDice, offHand.WeaponProps.DamageType),
		Metadata: map[string]any{
			"_route":      "attack",
			"weapon_id":   string(offHand.ID),
			"is_off_hand": true,
		},
	}
}

// computeEnemyWeaponAttacks 计算敌人装备的武器攻击
func (e *Engine) computeEnemyWeaponAttacks(game *model.GameState, enemy *model.Enemy, targets []model.ID) []AvailableAction {
	actions := make([]AvailableAction, 0)

	// 遍历游戏中的库存找到敌人的装备
	for _, inv := range game.Inventories {
		if inv.OwnerID != enemy.ID || inv.Equipment == nil {
			continue
		}
		if mainHand, exists := inv.Equipment.Slots[model.SlotMainHand]; exists && mainHand != nil && mainHand.WeaponProps != nil {
			wp := mainHand.WeaponProps
			rangeStr := "5尺触及"
			if wp.Reach {
				rangeStr = "10尺触及"
			}
			if wp.WeaponType == "ranged" && wp.Range > 0 {
				rangeStr = fmt.Sprintf("%d/%d尺", wp.Range, wp.LongRange)
			}
			actions = append(actions, AvailableAction{
				ID:             fmt.Sprintf("attack_%s", strings.ReplaceAll(strings.ToLower(mainHand.Name), " ", "_")),
				Category:       "attack",
				Name:           mainHand.Name + "攻击",
				Description:    formatWeaponDescription(mainHand),
				CostType:       "action",
				RequiresTarget: true,
				TargetType:     "single",
				ValidTargetIDs: targets,
				Range:          rangeStr,
				DamagePreview:  fmt.Sprintf("%s %s", wp.DamageDice, wp.DamageType),
				Metadata: map[string]any{
					"_route":    "attack",
					"weapon_id": string(mainHand.ID),
				},
			})
		}
		break // 一个敌人只有一个库存
	}

	return actions
}

// ============================================================================
// 法术动作计算
// ============================================================================

func computeSpellActions(
	sc *model.SpellcasterState,
	conditions conditionSummary,
	enemies, allies []model.ID,
	costTypeFilter string, // "action" or "bonus_action"
) []AvailableAction {
	if sc == nil || sc.Slots == nil {
		return nil
	}

	actions := make([]AvailableAction, 0)

	// 合并可用法术列表（PreparedSpells 优先，否则用 KnownSpells）
	spellIDs := sc.PreparedSpells
	if len(spellIDs) == 0 {
		spellIDs = sc.KnownSpells
	}

	for _, spellID := range spellIDs {
		spell, exists := data.GlobalRegistry.GetSpell(spellID)
		if !exists {
			continue
		}

		// 检查施法时间是否匹配
		castTimeUnit := strings.ToLower(spell.CastTime.Unit)
		if costTypeFilter == "action" && castTimeUnit != "action" {
			continue
		}
		if costTypeFilter == "bonus_action" && castTimeUnit != "bonus_action" && castTimeUnit != "bonus action" {
			continue
		}

		// 检查法术位可用性
		if spell.Level > 0 { // 非戏法需要法术位
			if sc.Slots.GetAvailableSlots(spell.Level) <= 0 {
				// 尝试更高级法术位
				if sc.Slots.GetLowestAvailableSlot(spell.Level) == 0 {
					continue
				}
			}
		}

		// 检查言语成分 vs 沉默状态
		if conditions.cantSpeak {
			hasVerbal := false
			for _, c := range spell.Components {
				if c == model.SpellComponentVerbal {
					hasVerbal = true
					break
				}
			}
			if hasVerbal {
				continue // 沉默状态下无法施放含言语成分的法术
			}
		}

		// 构建法术动作
		action := AvailableAction{
			ID:       fmt.Sprintf("cast_%s", strings.ReplaceAll(strings.ToLower(spell.Name), " ", "_")),
			Category: "spell",
			Name:     spell.Name,
			CostType: costTypeFilter,
			Range:    spell.Range,
			Metadata: map[string]any{
				"_route":     "spell",
				"spell_id":   spell.ID,
				"spell_name": spell.Name,
			},
		}

		// 设置描述
		desc := fmt.Sprintf("%d环", spell.Level)
		if spell.Level == 0 {
			desc = "戏法"
		}
		if spell.Concentration {
			desc += " 专注"
			if sc.ConcentrationSpell != "" {
				desc += "（将终止当前专注）"
			}
		}
		action.Description = desc

		// 法术位消耗说明
		if spell.Level > 0 {
			lowestSlot := sc.Slots.GetLowestAvailableSlot(spell.Level)
			remaining := sc.Slots.GetAvailableSlots(lowestSlot)
			action.ResourceCost = fmt.Sprintf("消耗%d环法术位（剩余%d个）", lowestSlot, remaining)
			action.Metadata["slot_level"] = lowestSlot
		}

		// 伤害/治疗预览
		if spell.DamageDice != "" {
			action.DamagePreview = fmt.Sprintf("%s %s", spell.DamageDice, spell.DamageType)
		}
		if spell.HealingDice != "" {
			action.DamagePreview = fmt.Sprintf("治疗 %s", spell.HealingDice)
		}

		// 目标类型
		rangeStr := strings.ToLower(spell.Range)
		switch {
		case rangeStr == "self" || rangeStr == "自身":
			action.TargetType = "self"
			action.RequiresTarget = false
		case rangeStr == "touch" || rangeStr == "触及":
			action.TargetType = "single"
			action.RequiresTarget = true
			// 触及法术目标包含敌友
			allTargets := append([]model.ID{}, enemies...)
			allTargets = append(allTargets, allies...)
			action.ValidTargetIDs = allTargets
		case strings.Contains(rangeStr, "半径") || strings.Contains(rangeStr, "radius"):
			action.TargetType = "area"
			action.RequiresTarget = false
		default:
			// 定向法术
			if spell.SaveDC != "" || spell.DamageDice != "" {
				action.TargetType = "single"
				action.RequiresTarget = true
				action.ValidTargetIDs = enemies
			} else if spell.HealingDice != "" {
				action.TargetType = "single"
				action.RequiresTarget = true
				action.ValidTargetIDs = allies
			} else {
				action.TargetType = "single"
				action.RequiresTarget = true
				allTargets := append([]model.ID{}, enemies...)
				allTargets = append(allTargets, allies...)
				action.ValidTargetIDs = allTargets
			}
		}

		actions = append(actions, action)
	}

	return actions
}

// ============================================================================
// 辅助函数
// ============================================================================

// conditionSummary 状态效果汇总
type conditionSummary struct {
	cantTakeActions bool
	cantSpeak       bool
	speedZero       bool
	attackDisadv    bool
}

// collectConditionEffects 汇总角色的所有状态效果
func collectConditionEffects(actor *model.Actor) conditionSummary {
	summary := conditionSummary{}
	for _, c := range actor.Conditions {
		effect := model.GetConditionEffect(c.Type)
		if effect.CantTakeActions {
			summary.cantTakeActions = true
		}
		if effect.CantSpeak {
			summary.cantSpeak = true
		}
		if effect.SpeedZero || effect.CantMove {
			summary.speedZero = true
		}
		if effect.AttackDisadvantage {
			summary.attackDisadv = true
		}
	}
	return summary
}

// classifyTargets 将战场上的角色分为敌友两方
func classifyTargets(game *model.GameState, myID model.ID, myType string) (enemies []model.ID, allies []model.ID) {
	if game.Combat == nil {
		return
	}

	isPlayerSide := myType == string(model.ActorTypePC) ||
		myType == string(model.ActorTypeCompanion)

	for _, entry := range game.Combat.Initiative {
		if entry.ActorID == myID || entry.IsDefeated {
			continue
		}

		// 确定此角色的阵营
		actorAny, ok := game.GetActor(entry.ActorID)
		if !ok {
			continue
		}
		var targetIsPlayerSide bool
		switch actorAny.(type) {
		case *model.PlayerCharacter, *model.Companion:
			targetIsPlayerSide = true
		case *model.Enemy:
			targetIsPlayerSide = false
		case *model.NPC:
			// NPC 默认视为中立，敌人方可攻击
			targetIsPlayerSide = false
		}

		if isPlayerSide == targetIsPlayerSide {
			allies = append(allies, entry.ActorID)
		} else {
			enemies = append(enemies, entry.ActorID)
		}
	}
	return
}

// makeStandardActions 生成基础 D&D 战斗动作
func makeStandardActions(enemies []model.ID) []AvailableAction {
	return []AvailableAction{
		{
			ID: "action_dash", Category: "standard_action",
			Name: "冲刺", Description: "本回合获得额外移动速度",
			CostType: "action", TargetType: "self",
			Metadata: map[string]any{"_route": "action", "action_type": "dash"},
		},
		{
			ID: "action_disengage", Category: "standard_action",
			Name: "撤离", Description: "本回合的移动不会引发借机攻击",
			CostType: "action", TargetType: "self",
			Metadata: map[string]any{"_route": "action", "action_type": "disengage"},
		},
		{
			ID: "action_dodge", Category: "standard_action",
			Name: "闪避", Description: "直到下回合开始，对你的攻击检定有劣势",
			CostType: "action", TargetType: "self",
			Metadata: map[string]any{"_route": "action", "action_type": "dodge"},
		},
		{
			ID: "action_help", Category: "standard_action",
			Name: "协助", Description: "给予盟友在下次攻击或检定上的优势",
			CostType: "action", RequiresTarget: true, TargetType: "single",
			Metadata: map[string]any{"_route": "action", "action_type": "help"},
		},
		{
			ID: "action_hide", Category: "standard_action",
			Name: "躲藏", Description: "进行敏捷（隐匿）检定尝试躲藏",
			CostType: "action", TargetType: "self",
			Metadata: map[string]any{"_route": "action", "action_type": "hide"},
		},
		{
			ID: "action_ready", Category: "standard_action",
			Name: "预备", Description: "准备一个触发条件和反应动作",
			CostType: "action", TargetType: "self",
			Metadata: map[string]any{"_route": "action", "action_type": "ready"},
		},
		{
			ID: "action_search", Category: "standard_action",
			Name: "搜索", Description: "进行感知（察觉）或智力（调查）检定",
			CostType: "action", TargetType: "self",
			Metadata: map[string]any{"_route": "action", "action_type": "search"},
		},
	}
}

// formatWeaponDescription 格式化武器描述
func formatWeaponDescription(item *model.Item) string {
	if item.WeaponProps == nil {
		return item.Description
	}
	wp := item.WeaponProps
	parts := []string{}
	if wp.Finesse {
		parts = append(parts, "灵巧")
	}
	if wp.Light {
		parts = append(parts, "轻型")
	}
	if wp.Heavy {
		parts = append(parts, "重型")
	}
	if wp.TwoHanded {
		parts = append(parts, "双手")
	}
	if wp.Versatile != "" {
		parts = append(parts, fmt.Sprintf("多用(%s)", wp.Versatile))
	}
	if wp.Reach {
		parts = append(parts, "长柄")
	}
	if wp.Thrown {
		parts = append(parts, "投掷")
	}
	if len(parts) > 0 {
		return strings.Join(parts, "、")
	}
	return fmt.Sprintf("%s %s", wp.DamageDice, wp.DamageType)
}
