package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/zwh8800/dnd-core/internal/model"
)

// Operation 定义引擎支持的所有操作类型
type Operation string

const (
	// 角色操作
	OpCreatePC        Operation = "create_pc"
	OpCreateNPC       Operation = "create_npc"
	OpCreateEnemy     Operation = "create_enemy"
	OpCreateCompanion Operation = "create_companion"
	OpGetActor        Operation = "get_actor"
	OpUpdateActor     Operation = "update_actor"
	OpRemoveActor     Operation = "remove_actor"

	// 骰子与检定
	OpRoll         Operation = "roll"
	OpAbilityCheck Operation = "ability_check"
	OpSkillCheck   Operation = "skill_check"
	OpSavingThrow  Operation = "saving_throw"

	// 战斗操作
	OpStartCombat    Operation = "start_combat"
	OpEndCombat      Operation = "end_combat"
	OpNextTurn       Operation = "next_turn"
	OpExecuteAction  Operation = "execute_action"
	OpExecuteAttack  Operation = "execute_attack"
	OpExecuteDamage  Operation = "execute_damage"
	OpExecuteHealing Operation = "execute_healing"

	// 法术操作
	OpCastSpell     Operation = "cast_spell"
	OpGetSpellSlots Operation = "get_spell_slots"
	OpPrepareSpells Operation = "prepare_spells"

	// 物品操作
	OpAddItem      Operation = "add_item"
	OpEquipItem    Operation = "equip_item"
	OpGetInventory Operation = "get_inventory"
	OpTransferItem Operation = "transfer_item"

	// 场景操作
	OpCreateScene      Operation = "create_scene"
	OpMoveActorToScene Operation = "move_actor_to_scene"
	OpSetCurrentScene  Operation = "set_current_scene"

	// 任务操作
	OpCreateQuest     Operation = "create_quest"
	OpAcceptQuest     Operation = "accept_quest"
	OpUpdateObjective Operation = "update_objective"
	OpCompleteQuest   Operation = "complete_quest"

	// 休息操作
	OpShortRest     Operation = "short_rest"
	OpStartLongRest Operation = "start_long_rest"
	OpEndLongRest   Operation = "end_long_rest"

	// 经验与升级
	OpAddExperience Operation = "add_experience"
	OpLevelUp       Operation = "level_up"

	// 状态查询（所有阶段都允许）
	OpGetStateSummary Operation = "get_state_summary"
	OpGetActorSheet   Operation = "get_actor_sheet"
	OpGetPhase        Operation = "get_phase"
)

// PhaseTransitionResult 阶段转换的结果
type PhaseTransitionResult struct {
	OldPhase    model.Phase `json:"old_phase"`
	NewPhase    model.Phase `json:"new_phase"`
	Reason      string      `json:"reason"`
	Timestamp   time.Time   `json:"timestamp"`
	AutoActions []string    `json:"auto_actions"`
	Message     string      `json:"message"`
}

// phasePermissions 定义每个阶段允许的操作
var phasePermissions = map[model.Phase]map[Operation]bool{
	model.PhaseCharacterCreation: {
		OpCreatePC: true, OpCreateNPC: true, OpCreateEnemy: true, OpCreateCompanion: true,
		OpGetActor: true, OpUpdateActor: true, OpRemoveActor: true,
		OpRoll: true, OpAbilityCheck: true, OpSkillCheck: true,
		OpGetSpellSlots: true, OpPrepareSpells: true,
		OpAddItem: true, OpEquipItem: true, OpGetInventory: true, OpTransferItem: true,
		OpCreateScene:     true,
		OpExecuteHealing:  true,
		OpLevelUp:         true,
		OpGetStateSummary: true, OpGetActorSheet: true, OpGetPhase: true,
	},
	model.PhaseExploration: {
		OpCreatePC: true, OpCreateNPC: true, OpCreateEnemy: true, OpCreateCompanion: true,
		OpGetActor: true, OpUpdateActor: true, OpRemoveActor: true,
		OpRoll: true, OpAbilityCheck: true, OpSkillCheck: true, OpSavingThrow: true,
		OpStartCombat:    true,
		OpExecuteHealing: true,
		OpCastSpell:      true, OpGetSpellSlots: true, OpPrepareSpells: true,
		OpAddItem: true, OpEquipItem: true, OpGetInventory: true, OpTransferItem: true,
		OpCreateScene: true, OpMoveActorToScene: true, OpSetCurrentScene: true,
		OpCreateQuest: true, OpAcceptQuest: true, OpUpdateObjective: true, OpCompleteQuest: true,
		OpShortRest: true, OpStartLongRest: true,
		OpAddExperience:   true,
		OpGetStateSummary: true, OpGetActorSheet: true, OpGetPhase: true,
	},
	model.PhaseCombat: {
		OpGetActor: true, OpUpdateActor: true,
		OpRoll: true, OpAbilityCheck: true, OpSkillCheck: true, OpSavingThrow: true,
		OpEndCombat: true, OpNextTurn: true, OpExecuteAction: true,
		OpExecuteAttack: true, OpExecuteDamage: true, OpExecuteHealing: true,
		OpCastSpell: true, OpGetSpellSlots: true,
		OpGetInventory:    true,
		OpGetStateSummary: true, OpGetActorSheet: true, OpGetPhase: true,
	},
	model.PhaseRest: {
		OpGetActor: true, OpUpdateActor: true,
		OpRoll:           true,
		OpExecuteHealing: true,
		OpGetSpellSlots:  true, OpGetInventory: true,
		OpShortRest: true, OpStartLongRest: true, OpEndLongRest: true,
		OpGetStateSummary: true, OpGetActorSheet: true, OpGetPhase: true,
	},
}

// checkPermission 检查当前阶段是否允许执行指定操作
func (e *Engine) checkPermission(phase model.Phase, op Operation) error {
	allowed, ok := phasePermissions[phase]
	if !ok {
		return fmt.Errorf("unknown phase: %s", phase)
	}
	if !allowed[op] {
		return &EngineError{
			Op:    string(op),
			Err:   ErrPhaseNotAllowed,
			Phase: phase,
			Details: map[string]any{
				"operation": op,
				"phase":     phase,
			},
		}
	}
	return nil
}

// SetPhase 切换游戏当前阶段
func (e *Engine) SetPhase(ctx context.Context, gameID model.ID, phase model.Phase, reason string) (*PhaseTransitionResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	oldPhase := game.Phase
	result := &PhaseTransitionResult{
		OldPhase:  oldPhase,
		NewPhase:  phase,
		Reason:    reason,
		Timestamp: time.Now(),
	}

	autoActions, message := handlePhaseTransition(game, oldPhase, phase)
	result.AutoActions = autoActions
	result.Message = message

	game.Phase = phase
	game.UpdatedAt = time.Now()

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return result, nil
}

// GetPhase 获取游戏当前阶段
func (e *Engine) GetPhase(ctx context.Context, gameID model.ID) (model.Phase, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return "", err
	}
	return game.Phase, nil
}

// GetAllowedOperations 获取当前阶段允许的所有操作
func (e *Engine) GetAllowedOperations(ctx context.Context, gameID model.ID) ([]Operation, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	allowed := phasePermissions[game.Phase]
	result := make([]Operation, 0)
	for op, ok := range allowed {
		if ok {
			result = append(result, op)
		}
	}
	return result, nil
}

// handlePhaseTransition 处理阶段转换时的自动操作
func handlePhaseTransition(game *model.GameState, oldPhase, newPhase model.Phase) ([]string, string) {
	autoActions := make([]string, 0)
	message := ""

	switch {
	case oldPhase == model.PhaseExploration && newPhase == model.PhaseCombat:
		autoActions = append(autoActions, "初始化战斗状态")
		message = "战斗开始！所有角色需要进行先攻检定"

	case oldPhase == model.PhaseCombat && newPhase == model.PhaseExploration:
		if game.Combat != nil {
			game.Combat.Status = model.CombatStatusFinished
			autoActions = append(autoActions, "结束战斗状态")
		}
		game.Combat = nil
		message = "战斗结束，队伍恢复探索状态"

	case oldPhase == model.PhaseExploration && newPhase == model.PhaseRest:
		message = "队伍开始长休，需要8小时才能完成"

	case oldPhase == model.PhaseRest && newPhase == model.PhaseExploration:
		if game.ActiveRest != nil && game.ActiveRest.Type == model.RestTypeLong {
			autoActions = append(autoActions, "应用长休恢复效果")
			message = "长休完成，队伍恢复活力"
		} else {
			message = "休息结束，队伍继续探索"
		}
		game.ActiveRest = nil

	case oldPhase == model.PhaseCharacterCreation && newPhase == model.PhaseExploration:
		message = "角色创建完成，冒险开始！"
	}

	return autoActions, message
}
