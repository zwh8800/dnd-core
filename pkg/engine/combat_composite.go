package engine

import (
	"context"
	"fmt"
	"sort"

	"github.com/zwh8800/dnd-core/pkg/model"
	"github.com/zwh8800/dnd-core/pkg/rules"
)

// ============================================================================
// 组合型战斗API — 类型定义
// ============================================================================

// SetupCombatRequest 组合型战斗初始化请求
type SetupCombatRequest struct {
	GameID         model.ID   `json:"game_id"`
	SceneID        model.ID   `json:"scene_id,omitempty"`
	ParticipantIDs []model.ID `json:"participant_ids"`
	IsSurprise     bool       `json:"is_surprise"`
	StealthySide   []model.ID `json:"stealthy_side,omitempty"`
	Observers      []model.ID `json:"observers,omitempty"`
}

// SetupCombatResult 组合型战斗初始化结果
type SetupCombatResult struct {
	Combat    *CombatInfo       `json:"combat"`
	FirstTurn *EnhancedTurnInfo `json:"first_turn"`
}

// EnhancedTurnInfo 增强版回合信息（包含可用动作和参与者状态）
type EnhancedTurnInfo struct {
	ActorID          model.ID                `json:"actor_id"`
	ActorName        string                  `json:"actor_name"`
	ActorType        string                  `json:"actor_type"`
	ActorHP          int                     `json:"actor_hp"`
	ActorMaxHP       int                     `json:"actor_max_hp"`
	ActorAC          int                     `json:"actor_ac"`
	ActorConditions  []string                `json:"actor_conditions"`
	Round            int                     `json:"round"`
	AvailableActions *AvailableActionsResult `json:"available_actions"`
	Participants     []CombatantStatus       `json:"participants"`
	CombatEnd        *CombatEndState         `json:"combat_end,omitempty"`
}

// CombatantStatus 战斗参与者状态快照
type CombatantStatus struct {
	ActorID    model.ID `json:"actor_id"`
	ActorName  string   `json:"actor_name"`
	ActorType  string   `json:"actor_type"`
	HP         int      `json:"hp"`
	MaxHP      int      `json:"max_hp"`
	AC         int      `json:"ac"`
	Conditions []string `json:"conditions"`
	IsDefeated bool     `json:"is_defeated"`
	IsAlly     bool     `json:"is_ally"`
}

// CombatEndState 战斗结束状态
type CombatEndState struct {
	Reason  string `json:"reason"`  // "victory", "defeat", "flee", "manual"
	Winners string `json:"winners"` // "players", "enemies", ""
}

// ExecuteTurnActionRequest 统一动作执行请求
type ExecuteTurnActionRequest struct {
	GameID     model.ID       `json:"game_id"`
	ActorID    model.ID       `json:"actor_id"`
	ActionID   string         `json:"action_id"`
	TargetID   model.ID       `json:"target_id,omitempty"`
	TargetIDs  []model.ID     `json:"target_ids,omitempty"`
	Parameters map[string]any `json:"parameters,omitempty"`
}

// ExecuteTurnActionResult 统一动作执行结果
type ExecuteTurnActionResult struct {
	Success          bool                    `json:"success"`
	ActionName       string                  `json:"action_name"`
	Narrative        string                  `json:"narrative"`
	AttackResult     *AttackResult           `json:"attack_result,omitempty"`
	DamageResult     *DamageResult           `json:"damage_result,omitempty"`
	HealResult       *HealResult             `json:"heal_result,omitempty"`
	SpellResult      *SpellResult            `json:"spell_result,omitempty"`
	ActionResult     *ActionResult           `json:"action_result,omitempty"`
	MoveResult       *MoveResult             `json:"move_result,omitempty"`
	RemainingActions *AvailableActionsResult `json:"remaining_actions"`
	TurnComplete     bool                    `json:"turn_complete"`
	CombatEnd        *CombatEndState         `json:"combat_end,omitempty"`
	Combat           *CombatInfo             `json:"combat"`
}

// ============================================================================
// 操作常量
// ============================================================================

const (
	OpSetupCombat       Operation = "setup_combat"
	OpNextTurnComposite Operation = "next_turn_composite"
	OpExecuteTurnAction Operation = "execute_turn_action"
)

// ============================================================================
// API 1: SetupCombat — 组合型战斗初始化
// ============================================================================

// SetupCombat 初始化一场战斗，返回战斗状态和第一个行动者的完整回合信息。
// 内部组合了战斗创建、先攻掷骰、突袭判定和可用动作计算。
func (e *Engine) SetupCombat(ctx context.Context, req SetupCombatRequest) (*SetupCombatResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpStartCombat); err != nil {
		return nil, err
	}

	if game.Combat != nil && game.Combat.Status == model.CombatStatusActive {
		return nil, ErrCombatAlreadyActive
	}

	// 自动创建场景（如果未指定）
	if req.SceneID == "" {
		sceneID := model.NewID()
		if game.Scenes == nil {
			game.Scenes = make(map[model.ID]*model.Scene)
		}
		game.Scenes[sceneID] = &model.Scene{
			ID:          sceneID,
			Name:        "Combat Arena",
			Description: "Auto-created combat scene",
			Type:        model.SceneTypeOutdoor,
		}
		req.SceneID = sceneID
	}

	// 确定参与者列表
	allParticipants := req.ParticipantIDs
	if req.IsSurprise {
		allParticipants = append(req.StealthySide, req.Observers...)
	}

	// 验证参与者
	for _, pid := range allParticipants {
		if _, ok := game.GetActor(pid); !ok {
			return nil, fmt.Errorf("actor %s not found", pid)
		}
	}

	// 创建战斗状态
	combat := &model.CombatState{
		ID:      model.NewID(),
		SceneID: req.SceneID,
		Status:  model.CombatStatusActive,
		Round:   1,
		Log:     make([]model.CombatLogEntry, 0),
	}

	// 处理突袭判定
	surprisedMap := make(map[model.ID]bool)
	if req.IsSurprise {
		surprisedMap = e.computeSurprise(game, req.StealthySide, req.Observers)
	}

	// 掷先攻并排序
	combat.Initiative = make([]model.CombatantEntry, 0, len(allParticipants))
	for _, actorID := range allParticipants {
		entry, err := e.rollInitiative(game, actorID)
		if err != nil {
			return nil, err
		}
		if surprisedMap[actorID] {
			entry.IsSurprised = true
		}
		combat.Initiative = append(combat.Initiative, entry)
	}

	sort.Slice(combat.Initiative, func(i, j int) bool {
		return combat.Initiative[i].InitiativeTotal > combat.Initiative[j].InitiativeTotal
	})

	// 设置第一个回合
	combat.CurrentIndex = 0
	combat.CurrentTurn = &model.TurnState{
		ActorID:         combat.Initiative[0].ActorID,
		Round:           1,
		ActionUsed:      false,
		BonusActionUsed: false,
		ReactionUsed:    false,
	}

	// 处理第一个角色是否被突袭
	if combat.Initiative[0].IsSurprised && combat.Round == 1 {
		combat.CurrentTurn.ActionUsed = true
		combat.CurrentTurn.BonusActionUsed = true
	}

	game.Combat = combat
	game.Phase = model.PhaseCombat

	// 计算第一个行动者的可用动作
	firstActorID := combat.Initiative[0].ActorID
	actions, err := e.computeAvailableActions(game, firstActorID)
	if err != nil {
		actions = &AvailableActionsResult{ActorID: firstActorID}
	}

	// 构建增强回合信息
	firstTurn := e.buildEnhancedTurnInfo(game, firstActorID, combat.Round, actions)

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &SetupCombatResult{
		Combat:    combatStateToInfo(combat),
		FirstTurn: firstTurn,
	}, nil
}

// ============================================================================
// API 2: NextTurnWithActions — 回合推进 + 完整态势
// ============================================================================

// NextTurnWithActions 推进到下一个角色的回合，并返回增强版回合信息（含可用动作和参与者状态）。
// 如果检测到战斗结束条件（所有敌人击败或所有玩家击败），会在返回结果中标记。
func (e *Engine) NextTurnWithActions(ctx context.Context, req NextTurnRequest) (*NextTurnResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpNextTurn); err != nil {
		return nil, err
	}

	if game.Combat == nil || game.Combat.Status != model.CombatStatusActive {
		return nil, ErrCombatNotActive
	}

	// 检测战斗结束
	if endState := checkCombatEnd(game); endState != nil {
		combatCopy := *game.Combat
		return &NextTurnResult{
			Combat: combatStateToInfo(&combatCopy),
			Turn: &EnhancedTurnInfo{
				CombatEnd: endState,
			},
		}, nil
	}

	// 推进回合
	game.Combat.AdvanceTurn()

	// 处理突袭
	currentCombatant := game.Combat.GetCurrentCombatant()
	if currentCombatant != nil && currentCombatant.IsSurprised && game.Combat.Round == 1 {
		game.Combat.CurrentTurn.ActionUsed = true
		game.Combat.CurrentTurn.BonusActionUsed = true
	}

	// 计算新行动者的可用动作
	actorID := game.Combat.CurrentTurn.ActorID
	actions, err := e.computeAvailableActions(game, actorID)
	if err != nil {
		actions = &AvailableActionsResult{ActorID: actorID}
	}

	// 构建增强回合信息
	turnInfo := e.buildEnhancedTurnInfo(game, actorID, game.Combat.Round, actions)

	// 再次检测战斗结束
	turnInfo.CombatEnd = checkCombatEnd(game)

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	combatCopy := *game.Combat
	return &NextTurnResult{
		Combat: combatStateToInfo(&combatCopy),
		Turn:   turnInfo,
	}, nil
}

// ============================================================================
// API 3: ExecuteTurnAction — 统一动作执行器
// ============================================================================

// ExecuteTurnAction 执行当前回合的一个动作，根据动作类型自动路由到对应的处理器。
// 执行后重新计算可用动作列表，检测回合是否完成和战斗是否结束。
func (e *Engine) ExecuteTurnAction(ctx context.Context, req ExecuteTurnActionRequest) (*ExecuteTurnActionResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpExecuteAction); err != nil {
		return nil, err
	}

	if game.Combat == nil || game.Combat.Status != model.CombatStatusActive {
		return nil, ErrCombatNotActive
	}

	if !game.Combat.IsActorTurn(req.ActorID) {
		return nil, ErrNotYourTurn
	}

	// 查找动作定义：先计算可用动作，再匹配 ActionID
	availableActions, err := e.computeAvailableActions(game, req.ActorID)
	if err != nil {
		return nil, fmt.Errorf("failed to compute available actions: %w", err)
	}

	actionDef := findActionByID(availableActions, req.ActionID)
	if actionDef == nil {
		return nil, fmt.Errorf("action %q is not available for actor %s", req.ActionID, req.ActorID)
	}

	// 根据 _route 分派执行
	route, _ := actionDef.Metadata["_route"].(string)
	result := &ExecuteTurnActionResult{
		ActionName: actionDef.Name,
	}

	switch route {
	case "attack":
		err = e.executeAttackRoute(game, req, actionDef, result)
	case "spell":
		err = e.executeSpellRoute(ctx, game, req, actionDef, result)
	case "action":
		err = e.executeStandardActionRoute(game, req, actionDef, result)
	case "class_feature":
		err = e.executeClassFeatureRoute(game, req, actionDef, result)
	case "reaction":
		err = e.executeReactionRoute(game, req, actionDef, result)
	default:
		// 未知路由，尝试作为通用动作执行
		err = e.executeGenericRoute(game, req, actionDef, result)
	}

	if err != nil {
		return nil, err
	}

	result.Success = true

	// 重新计算可用动作
	remaining, err := e.computeAvailableActions(game, req.ActorID)
	if err != nil {
		remaining = &AvailableActionsResult{ActorID: req.ActorID}
	}
	result.RemainingActions = remaining

	// 判断回合是否完成
	result.TurnComplete = remaining.IsEmpty()

	// 检测战斗结束
	result.CombatEnd = checkCombatEnd(game)

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	combatCopy := *game.Combat
	result.Combat = combatStateToInfo(&combatCopy)

	return result, nil
}

// ============================================================================
// ExecuteTurnAction 路由处理器
// ============================================================================

// executeAttackRoute 攻击路由
func (e *Engine) executeAttackRoute(
	game *model.GameState,
	req ExecuteTurnActionRequest,
	action *AvailableAction,
	result *ExecuteTurnActionResult,
) error {
	actor, ok := game.GetActor(req.ActorID)
	if !ok {
		return ErrNotFound
	}

	target, ok := game.GetActor(req.TargetID)
	if !ok {
		return ErrNotFound
	}

	var attackerActor, targetActor *model.Actor
	switch a := actor.(type) {
	case *model.PlayerCharacter:
		attackerActor = &a.Actor
	case *model.Enemy:
		attackerActor = &a.Actor
	case *model.NPC:
		attackerActor = &a.Actor
	case *model.Companion:
		attackerActor = &a.Actor
	}
	switch a := target.(type) {
	case *model.PlayerCharacter:
		targetActor = &a.Actor
	case *model.Enemy:
		targetActor = &a.Actor
	case *model.NPC:
		targetActor = &a.Actor
	case *model.Companion:
		targetActor = &a.Actor
	}

	// 构建 AttackInput
	attackInput := AttackInput{}
	if weaponID, ok := action.Metadata["weapon_id"].(string); ok {
		wid := model.ID(weaponID)
		attackInput.WeaponID = &wid
	}
	if isUnarmed, ok := action.Metadata["is_unarmed"].(bool); ok && isUnarmed {
		attackInput.IsUnarmed = true
	}
	if isOffHand, ok := action.Metadata["is_off_hand"].(bool); ok && isOffHand {
		attackInput.IsOffHand = true
	}

	// 计算攻击加值
	attackBonus := rules.CalcAttachBonus(actor, attackerActor.AbilityScores.Strength)

	// 掷攻击骰
	rollResult, _ := e.roller.Roll("1d20")
	rollValue := rollResult.Rolls[0].Value

	// 攻击检定
	attackCheck := rules.PerformAttackRoll(rollValue, attackBonus, targetActor.ArmorClass)

	attackResult := &AttackResult{
		Roll: &model.DiceResult{
			Rolls: []model.DiceRoll{{Value: rollValue}},
			Total: rollValue,
		},
		AttackTotal: attackCheck.Total,
		TargetAC:    attackCheck.TargetAC,
		Hit:         attackCheck.Hit,
		IsCritical:  attackCheck.IsCritical,
		IsFumble:    attackCheck.IsFumble,
		Message:     fmt.Sprintf("%s 攻击 %s: 掷骰 %d (总计 %d) vs AC %d", attackerActor.Name, targetActor.Name, rollValue, attackCheck.Total, targetActor.ArmorClass),
	}

	if attackCheck.Hit {
		damageResult, err := e.calculateAndApplyDamage(game, req.ActorID, req.TargetID, attackInput, attackCheck.IsCritical)
		if err != nil {
			return err
		}
		attackResult.Damage = damageResult
		attackResult.Message += fmt.Sprintf(" - 命中！造成 %d 点伤害", damageResult.FinalDamage)
		result.DamageResult = damageResult

		// 检查是否击杀
		if damageResult.IsDead {
			// 标记战败
			if game.Combat != nil {
				combatant := game.Combat.GetCombatantByActorID(req.TargetID)
				if combatant != nil {
					combatant.IsDefeated = true
				}
			}
		}
	} else {
		attackResult.Message += " - 未命中"
	}

	result.AttackResult = attackResult
	result.Narrative = attackResult.Message

	// 消耗动作/附赠动作
	if action.CostType == "bonus_action" {
		game.Combat.CurrentTurn.BonusActionUsed = true
	} else {
		game.Combat.CurrentTurn.ActionUsed = true
	}

	return nil
}

// executeSpellRoute 施法路由
func (e *Engine) executeSpellRoute(
	ctx context.Context,
	game *model.GameState,
	req ExecuteTurnActionRequest,
	action *AvailableAction,
	result *ExecuteTurnActionResult,
) error {
	spellID, _ := action.Metadata["spell_id"].(string)
	if spellID == "" {
		return fmt.Errorf("spell_id not found in action metadata")
	}

	slotLevel := 0
	if sl, ok := action.Metadata["slot_level"]; ok {
		switch v := sl.(type) {
		case int:
			slotLevel = v
		case float64:
			slotLevel = int(v)
		}
	}

	// 构建目标列表
	targetIDs := req.TargetIDs
	if req.TargetID != "" && len(targetIDs) == 0 {
		targetIDs = []model.ID{req.TargetID}
	}

	// 使用现有的 CastSpell API（释放锁后调用会死锁，所以需要内联）
	// 由于 CastSpell 也会获取锁，这里直接调用内部逻辑
	spellReq := CastSpellRequest{
		GameID:   req.GameID,
		CasterID: req.ActorID,
		Spell: SpellInput{
			SpellID:   spellID,
			SlotLevel: slotLevel,
			TargetIDs: targetIDs,
		},
	}

	// 释放锁，调用公开 API（避免重复实现施法逻辑）
	// 注意：这里不能直接调用 CastSpell 因为我们已持有锁
	// 因此我们用内部方法处理施法
	spellResult, err := e.castSpellInternal(ctx, game, spellReq)
	if err != nil {
		return err
	}

	result.SpellResult = spellResult
	result.Narrative = spellResult.Message

	// 消耗动作/附赠动作
	if action.CostType == "bonus_action" {
		game.Combat.CurrentTurn.BonusActionUsed = true
	} else {
		game.Combat.CurrentTurn.ActionUsed = true
	}

	return nil
}

// executeStandardActionRoute 标准动作路由（冲刺、闪避等）
func (e *Engine) executeStandardActionRoute(
	game *model.GameState,
	req ExecuteTurnActionRequest,
	action *AvailableAction,
	result *ExecuteTurnActionResult,
) error {
	actor, ok := game.GetActor(req.ActorID)
	if !ok {
		return ErrNotFound
	}

	var baseActor *model.Actor
	switch a := actor.(type) {
	case *model.PlayerCharacter:
		baseActor = &a.Actor
	case *model.Enemy:
		baseActor = &a.Actor
	case *model.NPC:
		baseActor = &a.Actor
	case *model.Companion:
		baseActor = &a.Actor
	}

	actionType, _ := action.Metadata["action_type"].(string)

	actionResult := &ActionResult{
		Success: true,
		Effects: make([]EffectDetail, 0),
	}

	switch model.ActionType(actionType) {
	case model.ActionDash:
		actionResult.Message = fmt.Sprintf("%s 执行冲刺动作，获得额外 %d 尺移动速度", baseActor.Name, baseActor.Speed)
		actionResult.Effects = append(actionResult.Effects, EffectDetail{
			Type: "movement", Description: "额外移动速度", Value: baseActor.Speed,
		})
	case model.ActionDisengage:
		actionResult.Message = fmt.Sprintf("%s 执行撤离动作，本回合移动不引发借机攻击", baseActor.Name)
	case model.ActionDodge:
		actionResult.Message = fmt.Sprintf("%s 执行闪避动作，对其攻击有劣势", baseActor.Name)
	case model.ActionHelp:
		actionResult.Message = fmt.Sprintf("%s 执行协助动作", baseActor.Name)
	case model.ActionHide:
		actionResult.Message = fmt.Sprintf("%s 尝试躲藏", baseActor.Name)
	case model.ActionReady:
		actionResult.Message = fmt.Sprintf("%s 执行预备动作", baseActor.Name)
	case model.ActionSearch:
		actionResult.Message = fmt.Sprintf("%s 进行搜索", baseActor.Name)
	default:
		actionResult.Message = fmt.Sprintf("%s 执行动作: %s", baseActor.Name, actionType)
	}

	game.Combat.CurrentTurn.ActionUsed = true

	result.ActionResult = actionResult
	result.Narrative = actionResult.Message

	return nil
}

// executeClassFeatureRoute 职业特性路由
func (e *Engine) executeClassFeatureRoute(
	game *model.GameState,
	req ExecuteTurnActionRequest,
	action *AvailableAction,
	result *ExecuteTurnActionResult,
) error {
	actor, ok := game.GetActor(req.ActorID)
	if !ok {
		return ErrNotFound
	}

	var baseActor *model.Actor
	switch a := actor.(type) {
	case *model.PlayerCharacter:
		baseActor = &a.Actor
	case *model.Enemy:
		baseActor = &a.Actor
	case *model.NPC:
		baseActor = &a.Actor
	case *model.Companion:
		baseActor = &a.Actor
	}

	actionResult := &ActionResult{
		Success: true,
		Message: fmt.Sprintf("%s 使用 %s", baseActor.Name, action.Name),
		Effects: make([]EffectDetail, 0),
	}

	// 消耗动作资源
	switch action.CostType {
	case "bonus_action":
		game.Combat.CurrentTurn.BonusActionUsed = true
	case "free_action":
		// 自由动作不消耗
	default:
		game.Combat.CurrentTurn.ActionUsed = true
	}

	result.ActionResult = actionResult
	result.Narrative = actionResult.Message

	return nil
}

// executeReactionRoute 反应路由
func (e *Engine) executeReactionRoute(
	game *model.GameState,
	req ExecuteTurnActionRequest,
	action *AvailableAction,
	result *ExecuteTurnActionResult,
) error {
	game.Combat.CurrentTurn.ReactionUsed = true

	result.Narrative = fmt.Sprintf("准备反应: %s", action.Name)
	result.ActionResult = &ActionResult{
		Success: true,
		Message: result.Narrative,
		Effects: make([]EffectDetail, 0),
	}

	return nil
}

// executeGenericRoute 通用路由
func (e *Engine) executeGenericRoute(
	game *model.GameState,
	req ExecuteTurnActionRequest,
	action *AvailableAction,
	result *ExecuteTurnActionResult,
) error {
	// 对未知类型的动作，按照 CostType 消耗资源
	switch action.CostType {
	case "action":
		game.Combat.CurrentTurn.ActionUsed = true
	case "bonus_action":
		game.Combat.CurrentTurn.BonusActionUsed = true
	case "reaction":
		game.Combat.CurrentTurn.ReactionUsed = true
	}

	result.Narrative = fmt.Sprintf("执行: %s", action.Name)
	result.ActionResult = &ActionResult{
		Success: true,
		Message: result.Narrative,
		Effects: make([]EffectDetail, 0),
	}

	return nil
}

// ============================================================================
// 内部辅助函数
// ============================================================================

// buildEnhancedTurnInfo 构建增强版回合信息
func (e *Engine) buildEnhancedTurnInfo(
	game *model.GameState,
	actorID model.ID,
	round int,
	actions *AvailableActionsResult,
) *EnhancedTurnInfo {
	info := &EnhancedTurnInfo{
		ActorID:          actorID,
		Round:            round,
		AvailableActions: actions,
		ActorConditions:  make([]string, 0),
		Participants:     make([]CombatantStatus, 0),
	}

	// 填充行动者信息
	actorAny, ok := game.GetActor(actorID)
	if ok {
		var baseActor *model.Actor
		switch a := actorAny.(type) {
		case *model.PlayerCharacter:
			baseActor = &a.Actor
			info.ActorName = a.Name
			info.ActorType = string(model.ActorTypePC)
		case *model.Enemy:
			baseActor = &a.Actor
			info.ActorName = a.Name
			info.ActorType = string(model.ActorTypeEnemy)
		case *model.NPC:
			baseActor = &a.Actor
			info.ActorName = a.Name
			info.ActorType = string(model.ActorTypeNPC)
		case *model.Companion:
			baseActor = &a.Actor
			info.ActorName = a.Name
			info.ActorType = string(model.ActorTypeCompanion)
		}
		if baseActor != nil {
			info.ActorHP = baseActor.HitPoints.Current
			info.ActorMaxHP = baseActor.HitPoints.Maximum
			info.ActorAC = baseActor.ArmorClass
			for _, c := range baseActor.Conditions {
				info.ActorConditions = append(info.ActorConditions, string(c.Type))
			}
		}
	}

	// 构建参与者状态列表
	if game.Combat != nil {
		isPlayerSide := info.ActorType == string(model.ActorTypePC) ||
			info.ActorType == string(model.ActorTypeCompanion)

		for _, entry := range game.Combat.Initiative {
			participantAny, ok := game.GetActor(entry.ActorID)
			if !ok {
				continue
			}

			status := CombatantStatus{
				ActorID:    entry.ActorID,
				ActorName:  entry.ActorName,
				IsDefeated: entry.IsDefeated,
				Conditions: make([]string, 0),
			}

			var pActor *model.Actor
			switch a := participantAny.(type) {
			case *model.PlayerCharacter:
				pActor = &a.Actor
				status.ActorType = string(model.ActorTypePC)
			case *model.Enemy:
				pActor = &a.Actor
				status.ActorType = string(model.ActorTypeEnemy)
			case *model.NPC:
				pActor = &a.Actor
				status.ActorType = string(model.ActorTypeNPC)
			case *model.Companion:
				pActor = &a.Actor
				status.ActorType = string(model.ActorTypeCompanion)
			}

			if pActor != nil {
				status.HP = pActor.HitPoints.Current
				status.MaxHP = pActor.HitPoints.Maximum
				status.AC = pActor.ArmorClass
				for _, c := range pActor.Conditions {
					status.Conditions = append(status.Conditions, string(c.Type))
				}
			}

			// 判断是否为盟友
			pIsPlayerSide := status.ActorType == string(model.ActorTypePC) ||
				status.ActorType == string(model.ActorTypeCompanion)
			status.IsAlly = isPlayerSide == pIsPlayerSide

			info.Participants = append(info.Participants, status)
		}
	}

	return info
}

// computeSurprise 计算突袭判定
func (e *Engine) computeSurprise(game *model.GameState, stealthySide, observers []model.ID) map[model.ID]bool {
	surprisedMap := make(map[model.ID]bool)

	// 潜行方进行隐匿检定
	highestStealth := 0
	for _, actorID := range stealthySide {
		actor, ok := game.GetActor(actorID)
		if !ok {
			continue
		}

		var baseActor *model.Actor
		var pc *model.PlayerCharacter
		switch a := actor.(type) {
		case *model.PlayerCharacter:
			baseActor = &a.Actor
			pc = a
		case *model.Enemy:
			baseActor = &a.Actor
		case *model.NPC:
			baseActor = &a.Actor
		case *model.Companion:
			baseActor = &a.Actor
		}

		dexMod := rules.AbilityModifier(baseActor.AbilityScores.Dexterity)
		stealthBonus := dexMod
		if pc != nil && pc.Proficiencies.ProficientSkills != nil {
			if pc.Proficiencies.ProficientSkills[model.SkillStealth] {
				stealthBonus += rules.ProficiencyBonus(pc.TotalLevel)
			}
		}

		stealthRoll, err := e.roller.Roll("1d20")
		if err != nil {
			continue
		}
		stealthTotal := stealthRoll.Total + stealthBonus
		if stealthTotal > highestStealth {
			highestStealth = stealthTotal
		}
	}

	// 观察方被动察觉 vs 隐匿
	for _, observerID := range observers {
		observer, ok := game.GetActor(observerID)
		if !ok {
			continue
		}

		var baseActor *model.Actor
		var pc *model.PlayerCharacter
		switch a := observer.(type) {
		case *model.PlayerCharacter:
			baseActor = &a.Actor
			pc = a
		case *model.Enemy:
			baseActor = &a.Actor
		case *model.NPC:
			baseActor = &a.Actor
		case *model.Companion:
			baseActor = &a.Actor
		}

		wisMod := rules.AbilityModifier(baseActor.AbilityScores.Wisdom)
		perceptionBonus := wisMod
		if pc != nil && pc.Proficiencies.ProficientSkills != nil {
			if pc.Proficiencies.ProficientSkills[model.SkillPerception] {
				perceptionBonus += rules.ProficiencyBonus(pc.TotalLevel)
			}
		}

		passivePerception := rules.PassiveScore(perceptionBonus)
		if passivePerception < highestStealth {
			surprisedMap[observerID] = true
		}
	}

	return surprisedMap
}

// checkCombatEnd 检测战斗是否应该结束
func checkCombatEnd(game *model.GameState) *CombatEndState {
	if game.Combat == nil {
		return nil
	}

	allEnemiesDefeated := true
	allPlayersDefeated := true
	hasEnemies := false
	hasPlayers := false

	for _, entry := range game.Combat.Initiative {
		actorAny, ok := game.GetActor(entry.ActorID)
		if !ok {
			continue
		}

		switch actorAny.(type) {
		case *model.Enemy:
			hasEnemies = true
			if !entry.IsDefeated {
				allEnemiesDefeated = false
			}
		case *model.PlayerCharacter, *model.Companion:
			hasPlayers = true
			if !entry.IsDefeated {
				allPlayersDefeated = false
			}
		}
	}

	if hasEnemies && allEnemiesDefeated {
		return &CombatEndState{
			Reason:  "victory",
			Winners: "players",
		}
	}

	if hasPlayers && allPlayersDefeated {
		return &CombatEndState{
			Reason:  "defeat",
			Winners: "enemies",
		}
	}

	return nil
}

// findActionByID 在可用动作列表中查找指定ID的动作
func findActionByID(actions *AvailableActionsResult, actionID string) *AvailableAction {
	if actions == nil {
		return nil
	}

	// 搜索所有动作类别
	allActions := make([]AvailableAction, 0, len(actions.Actions)+len(actions.BonusActions)+len(actions.Reactions)+len(actions.FreeActions))
	allActions = append(allActions, actions.Actions...)
	allActions = append(allActions, actions.BonusActions...)
	allActions = append(allActions, actions.Reactions...)
	allActions = append(allActions, actions.FreeActions...)

	for i := range allActions {
		if allActions[i].ID == actionID {
			return &allActions[i]
		}
	}

	// 特殊动作ID "move" 用于移动
	if actionID == "move" {
		return &AvailableAction{
			ID:       "move",
			Category: "movement",
			Name:     "移动",
			CostType: "free_action",
			Metadata: map[string]any{"_route": "move"},
		}
	}

	return nil
}

// castSpellInternal 内部施法逻辑（不获取锁，由调用方负责）
// 这是 CastSpell 的无锁版本，供组合API在已持有锁的上下文中调用
func (e *Engine) castSpellInternal(_ context.Context, game *model.GameState, req CastSpellRequest) (*SpellResult, error) {
	caster, ok := game.GetActor(req.CasterID)
	if !ok {
		return nil, ErrNotFound
	}

	var casterActor *model.Actor
	var spellcaster *model.SpellcasterState
	switch c := caster.(type) {
	case *model.PlayerCharacter:
		casterActor = &c.Actor
		spellcaster = c.Spellcasting
	default:
		return nil, fmt.Errorf("only player characters can cast spells")
	}

	if spellcaster == nil {
		return nil, fmt.Errorf("actor %s is not a spellcaster", casterActor.Name)
	}

	spellDef := findSpellDefinition(req.Spell.SpellID)
	if spellDef == nil {
		return nil, fmt.Errorf("spell %s not found", req.Spell.SpellID)
	}

	if !canCastSpell(spellcaster, req.Spell.SpellID) {
		return nil, fmt.Errorf("caster does not know or have prepared spell: %s", spellDef.Name)
	}

	// 消耗法术位
	if spellDef.Level > 0 {
		slotLevel := req.Spell.SlotLevel
		if slotLevel == 0 {
			slotLevel = spellDef.Level
		}
		if slotLevel < spellDef.Level {
			return nil, fmt.Errorf("slot level %d is lower than spell level %d", slotLevel, spellDef.Level)
		}
		if spellcaster.Slots.GetAvailableSlots(slotLevel) <= 0 {
			return nil, ErrInsufficientSlots
		}
		spellcaster.Slots.UseSlot(slotLevel)
	}

	// 处理专注
	if spellcaster.IsConcentrating() && spellDef.Concentration {
		spellcaster.ConcentrationSpell = ""
	}

	result := &SpellResult{
		SpellName:     spellDef.Name,
		SlotLevel:     req.Spell.SlotLevel,
		CasterSaveDC:  spellcaster.SpellSaveDC,
		Targets:       make([]SpellTargetResult, 0),
		Concentration: spellDef.Concentration,
	}

	// 攻击类法术
	if spellDef.DamageDice != "" && spellDef.SaveDC == "" {
		attackBonus := spellcaster.SpellAttackBonus
		rollResult, _ := e.roller.Roll("1d20")
		attackTotal := rollResult.Total + attackBonus
		result.AttackRoll = rollResult
		result.AttackTotal = attackTotal

		for _, targetID := range req.Spell.TargetIDs {
			target, ok := game.GetActor(targetID)
			if !ok {
				continue
			}
			var targetActor *model.Actor
			switch t := target.(type) {
			case *model.PlayerCharacter:
				targetActor = &t.Actor
			case *model.NPC:
				targetActor = &t.Actor
			case *model.Enemy:
				targetActor = &t.Actor
			case *model.Companion:
				targetActor = &t.Actor
			}

			isNat20 := rules.IsCriticalHit(rollResult.Rolls[0].Value)
			isNat1 := rules.IsCriticalFumble(rollResult.Rolls[0].Value)
			hit := attackTotal >= targetActor.ArmorClass || isNat20
			if isNat1 {
				hit = false
			}

			targetResult := SpellTargetResult{ActorID: targetID}
			if hit {
				damageResult, err := e.applySpellDamage(game, req.CasterID, targetID, spellDef, req.Spell.UpcastLevel, isNat20)
				if err != nil {
					return nil, err
				}
				targetResult.Damage = damageResult
			}
			result.Targets = append(result.Targets, targetResult)
		}
	} else if spellDef.DamageDice != "" && spellDef.SaveDC != "" {
		// 豁免类伤害法术
		for _, targetID := range req.Spell.TargetIDs {
			target, ok := game.GetActor(targetID)
			if !ok {
				continue
			}
			var targetActor *model.Actor
			switch t := target.(type) {
			case *model.PlayerCharacter:
				targetActor = &t.Actor
			case *model.NPC:
				targetActor = &t.Actor
			case *model.Enemy:
				targetActor = &t.Actor
			case *model.Companion:
				targetActor = &t.Actor
			}

			saveRoll, _ := e.roller.Roll("1d20")
			saveAbility := targetActor.AbilityScores.Get(spellDef.SaveDC)
			saveBonus := rules.AbilityModifier(saveAbility)
			saveTotal := saveRoll.Total + saveBonus
			saveSuccess := saveTotal >= spellcaster.SpellSaveDC

			targetResult := SpellTargetResult{
				ActorID:     targetID,
				SaveRoll:    saveRoll,
				SaveTotal:   saveTotal,
				SaveSuccess: saveSuccess,
			}

			if spellDef.DamageDice != "" {
				roll, err := e.roller.Roll(spellDef.DamageDice)
				baseDamage := 0
				if err == nil {
					baseDamage = roll.Total
				}
				if saveSuccess {
					baseDamage /= 2
				}
				damageResult, err := e.applySpellDamageDirect(game, req.CasterID, targetID, baseDamage, spellDef.DamageType, false)
				if err != nil {
					return nil, err
				}
				targetResult.Damage = damageResult
			}
			result.Targets = append(result.Targets, targetResult)
		}
	} else if spellDef.HealingDice != "" {
		// 治疗法术
		roll, err := e.roller.Roll(spellDef.HealingDice)
		healingAmount := 0
		if err == nil {
			healingAmount = roll.Total
		}
		for _, targetID := range req.Spell.TargetIDs {
			healResult, err := e.applyHealingInSpell(game, targetID, healingAmount)
			if err != nil {
				return nil, err
			}
			result.Targets = append(result.Targets, SpellTargetResult{
				ActorID: targetID,
				Healing: healResult,
			})
		}
	}

	// 设置专注
	if spellDef.Concentration {
		spellcaster.ConcentrationSpell = req.Spell.SpellID
	}

	result.Message = fmt.Sprintf("%s 施展了 %s", casterActor.Name, spellDef.Name)
	if spellDef.Concentration {
		result.Message += " (专注)"
	}

	return result, nil
}
