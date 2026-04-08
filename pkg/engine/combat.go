package engine

import (
	"context"
	"fmt"
	"sort"

	"github.com/zwh8800/dnd-core/pkg/model"
	"github.com/zwh8800/dnd-core/pkg/rules"
)

// ============================================================================
// 战斗API结构体定义
// ============================================================================

// CombatInfo 战斗状态信息（替代 *model.CombatState）
// 用于API返回时封装战斗的核心状态数据
type CombatInfo struct {
	ID           model.ID             `json:"id"`            // 战斗唯一标识
	SceneID      model.ID             `json:"scene_id"`      // 战斗所在场景
	Status       model.CombatStatus   `json:"status"`        // 战斗状态
	Round        int                  `json:"round"`         // 当前回合数
	CurrentIndex int                  `json:"current_index"` // 当前先攻索引
	Initiative   []CombatantEntryInfo `json:"initiative"`    // 先攻顺序
	CurrentTurn  *TurnInfo            `json:"current_turn"`  // 当前回合信息
}

// CombatantEntryInfo 战斗者条目信息
// 描述参与战斗的角色在先攻序列中的信息
type CombatantEntryInfo struct {
	ActorID         model.ID `json:"actor_id"`         // 角色ID
	ActorName       string   `json:"actor_name"`       // 角色名称
	InitiativeRoll  int      `json:"initiative_roll"`  // 先攻掷骰值（原始d20结果）
	InitiativeTotal int      `json:"initiative_total"` // 先攻总值（d20+敏捷修正+加值）
	IsSurprised     bool     `json:"is_surprised"`     // 是否处于被突袭状态（第一回合无法行动）
	IsDefeated      bool     `json:"is_defeated"`      // 是否已被击败（HP降至0）
}

// StartCombatRequest 开始战斗请求
// 用于启动一场标准战斗遭遇
type StartCombatRequest struct {
	GameID         model.ID   `json:"game_id"`         // 游戏会话ID
	SceneID        model.ID   `json:"scene_id"`        // 战斗所在场景ID
	ParticipantIDs []model.ID `json:"participant_ids"` // 参战角色ID列表
}

// StartCombatResult 开始战斗结果
// 返回新创建的战斗状态信息
type StartCombatResult struct {
	Combat *CombatInfo `json:"combat"` // 战斗状态信息
}

// StartCombatWithSurpriseRequest 突袭战斗请求
// 用于启动带突袭判定的战斗，区分潜行方和观察方
type StartCombatWithSurpriseRequest struct {
	GameID       model.ID   `json:"game_id"`       // 游戏会话ID
	SceneID      model.ID   `json:"scene_id"`      // 战斗所在场景ID
	StealthySide []model.ID `json:"stealthy_side"` // 潜行方角色ID列表（不会第一回合被突袭）
	Observers    []model.ID `json:"observers"`     // 被观察方角色ID列表（第一回合被突袭）
}

// StartCombatWithSurpriseResult 突袭战斗结果
// 返回新创建的带突袭判定的战斗状态信息
type StartCombatWithSurpriseResult struct {
	Combat *CombatInfo `json:"combat"` // 战斗状态信息
}

// EndCombatRequest 结束战斗请求
// 用于主动结束当前活跃的战斗
type EndCombatRequest struct {
	GameID model.ID `json:"game_id"` // 游戏会话ID
}

// GetCurrentCombatRequest 获取当前战斗请求
// 用于查询当前活跃战斗的状态
type GetCurrentCombatRequest struct {
	GameID model.ID `json:"game_id"` // 游戏会话ID
}

// GetCurrentCombatResult 获取当前战斗结果
// 返回当前活跃战斗的状态信息
type GetCurrentCombatResult struct {
	Combat *CombatInfo `json:"combat"` // 战斗状态信息
}

// NextTurnRequest 下一回合请求
// 用于推进战斗到下一个角色的回合
type NextTurnRequest struct {
	GameID model.ID `json:"game_id"` // 游戏会话ID
}

// NextTurnResult 下一回合结果
// 返回推进后的战斗状态信息
type NextTurnResult struct {
	Combat *CombatInfo `json:"combat"` // 战斗状态信息
}

// GetCurrentTurnRequest 获取当前回合请求
// 用于查询当前回合角色的详细信息
type GetCurrentTurnRequest struct {
	GameID model.ID `json:"game_id"` // 游戏会话ID
}

// ExecuteActionRequest 执行动作请求
// 用于在当前回合执行一个动作（冲刺、脱离、闪避等）
type ExecuteActionRequest struct {
	GameID  model.ID    `json:"game_id"`  // 游戏会话ID
	ActorID model.ID    `json:"actor_id"` // 角色ID
	Action  ActionInput `json:"action"`   // 动作输入
}

// ExecuteActionResult 执行动作结果
// 返回动作执行后的结果
type ExecuteActionResult struct {
	ActionResult *ActionResult `json:"action_result"` // 动作执行结果
	Combat       *CombatInfo   `json:"combat"`        // 当前战斗状态
}

// ExecuteAttackRequest 执行攻击请求
// 用于执行一次攻击动作
type ExecuteAttackRequest struct {
	GameID     model.ID    `json:"game_id"`     // 游戏会话ID
	AttackerID model.ID    `json:"attacker_id"` // 攻击者ID
	TargetID   model.ID    `json:"target_id"`   // 目标ID
	Attack     AttackInput `json:"attack"`      // 攻击输入
}

// ExecuteAttackResult 执行攻击结果
// 返回攻击的完整结果，包括命中判定和伤害
type ExecuteAttackResult struct {
	AttackResult *AttackResult `json:"attack_result"` // 攻击结果
	Combat       *CombatInfo   `json:"combat"`        // 当前战斗状态
}

// ExecuteDamageRequest 执行伤害请求
// 用于直接对角色施加伤害（如环境伤害、陷阱伤害等）
type ExecuteDamageRequest struct {
	GameID   model.ID    `json:"game_id"`   // 游戏会话ID
	TargetID model.ID    `json:"target_id"` // 目标ID
	Damage   DamageInput `json:"damage"`    // 伤害输入
}

// ExecuteDamageResult 执行伤害结果
// 返回应用伤害后的结果
type ExecuteDamageResult struct {
	DamageResult *DamageResult `json:"damage_result"` // 伤害结果
	Combat       *CombatInfo   `json:"combat"`        // 当前战斗状态
}

// ExecuteHealingRequest 执行治疗请求
// 用于对角色进行治疗
type ExecuteHealingRequest struct {
	GameID   model.ID `json:"game_id"`   // 游戏会话ID
	TargetID model.ID `json:"target_id"` // 目标ID
	Amount   int      `json:"amount"`    // 治疗量
}

// ExecuteHealingResult 执行治疗结果
// 返回治疗后的角色状态变化
type ExecuteHealingResult struct {
	TargetID  model.ID `json:"target_id"`  // 目标ID
	Healed    int      `json:"healed"`     // 实际治疗量
	CurrentHP int      `json:"current_hp"` // 治疗后当前HP
	Message   string   `json:"message"`    // 人类可读消息
}

// MoveActorRequest 移动角色请求
// 用于在场景中移动角色位置
type MoveActorRequest struct {
	GameID  model.ID    `json:"game_id"`  // 游戏会话ID
	ActorID model.ID    `json:"actor_id"` // 角色ID
	To      model.Point `json:"to"`       // 目标位置
}

// MoveActorResult 移动角色结果
// 返回移动操作的结果
type MoveActorResult struct {
	MoveResult      *MoveResult      `json:"move_result"`       // 移动结果
	SceneMoveResult *SceneMoveResult `json:"scene_move_result"` // 场景移动结果（如果是场景间移动）
	Combat          *CombatInfo      `json:"combat"`            // 当前战斗状态（如果战斗中有移动）
}

// ActionInput 动作输入
type ActionInput struct {
	Type    model.ActionType `json:"type"`              // 动作类型
	Details map[string]any   `json:"details,omitempty"` // 动作详情（可选）
}

// ActionResult 动作执行结果
type ActionResult struct {
	Success bool              `json:"success"`        // 是否成功
	Message string            `json:"message"`        // 人类可读消息
	Roll    *model.DiceResult `json:"roll,omitempty"` // 相关掷骰结果
	Effects []EffectDetail    `json:"effects"`        // 产生的效果列表
}

// EffectDetail 效果详情
// 描述动作产生的具体效果
type EffectDetail struct {
	Type        string `json:"type"`            // 效果类型
	Description string `json:"description"`     // 效果描述
	Value       int    `json:"value,omitempty"` // 效果数值（如果有）
}

// AttackInput 攻击输入
// 描述一次攻击的所有参数
type AttackInput struct {
	WeaponID      *model.ID          `json:"weapon_id,omitempty"`    // 武器ID（如果使用武器攻击）
	SpellID       *string            `json:"spell_id,omitempty"`     // 法术ID（如果施法攻击）
	IsUnarmed     bool               `json:"is_unarmed"`             // 是否徒手攻击
	IsOffHand     bool               `json:"is_off_hand"`            // 是否为副手攻击
	WeaponMastery string             `json:"weapon_mastery"`         // 武器掌控类型
	Advantage     model.RollModifier `json:"advantage"`              // 攻击掷骰的优劣势修正
	ExtraDamage   []DamageInput      `json:"extra_damage,omitempty"` // 额外伤害（如偷袭、爆发等）
}

// AttackResult 攻击结果
// 描述一次攻击的完整结果
type AttackResult struct {
	Roll        *model.DiceResult `json:"roll"`              // 攻击掷骰结果
	AttackTotal int               `json:"attack_total"`      // 攻击总值
	TargetAC    int               `json:"target_ac"`         // 目标护甲等级
	Hit         bool              `json:"hit"`               // 是否命中
	IsCritical  bool              `json:"is_critical"`       // 是否为重击（自然20）
	IsFumble    bool              `json:"is_fumble"`         // 是否为大失败（自然1）
	Damage      *DamageResult     `json:"damage,omitempty"`  // 伤害结果（如果命中）
	Effects     []AttackEffect    `json:"effects,omitempty"` // 额外效果（如武器掌控）
	GrazeDamage int               `json:"graze_damage"`      // 擦伤伤害（未命中时）
	Message     string            `json:"message"`           // 人类可读消息
}

// AttackEffect 攻击额外效果
type AttackEffect struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

// DamageInput 伤害输入
// 描述伤害的来源和类型
type DamageInput struct {
	Amount int              `json:"amount"`         // 伤害数量
	Type   model.DamageType `json:"type"`           // 伤害类型
	Dice   string           `json:"dice,omitempty"` // 伤害骰子表达式（如"2d6"）
	Source model.ID         `json:"source"`         // 伤害来源ID
}

// DamageResult 伤害结果
// 描述应用伤害后的完整结果
type DamageResult struct {
	RawDamage       int                `json:"raw_damage"`              // 原始伤害值
	Resistances     []model.DamageType `json:"resistances_applied"`     // 应用的伤害抗性
	Vulnerabilities []model.DamageType `json:"vulnerabilities_applied"` // 应用的伤害易伤
	FinalDamage     int                `json:"final_damage"`            // 最终伤害值
	TargetHPBefore  int                `json:"target_hp_before"`        // 目标攻击前HP
	TargetHPAfter   int                `json:"target_hp_after"`         // 目标攻击后HP
	IsDead          bool               `json:"is_dead"`                 // 是否死亡
	IsStabilized    bool               `json:"is_stabilized"`           // 是否进入稳定状态
	DeathSaves      *DeathSaveUpdate   `json:"death_saves,omitempty"`   // 死亡豁免状态更新
	Message         string             `json:"message"`                 // 人类可读消息
}

// DeathSaveUpdate 死亡豁免更新
// 描述角色死亡豁免状态的变化
type DeathSaveUpdate struct {
	Successes int  `json:"successes"` // 成功次数
	Failures  int  `json:"failures"`  // 失败次数
	IsStable  bool `json:"is_stable"` // 是否稳定
}

// HealResult 治疗结果
// 描述一次治疗的效果
type HealResult struct {
	Amount    int    `json:"amount"`     // 治疗量
	HPBefore  int    `json:"hp_before"`  // 治疗前HP
	HPAfter   int    `json:"hp_after"`   // 治疗后HP
	WasStable bool   `json:"was_stable"` // 治疗前是否处于稳定状态
	Message   string `json:"message"`    // 人类可读消息
}

// MoveResult 移动结果
// 描述一次移动操作的结果
type MoveResult struct {
	Success       bool   `json:"success"`        // 是否成功
	DistanceMoved int    `json:"distance_moved"` // 实际移动距离
	RemainingMove int    `json:"remaining_move"` // 剩余移动距离
	Message       string `json:"message"`        // 人类可读消息
}

// TurnInfo 回合信息
// 描述当前回合的详细信息
type TurnInfo struct {
	ActorID              model.ID `json:"actor_id"`               // 角色ID
	ActorName            string   `json:"actor_name"`             // 角色名称
	Round                int      `json:"round"`                  // 回合数
	InitiativePos        int      `json:"initiative_position"`    // 先攻序列中的位置（从1开始）
	MovementLeft         int      `json:"movement_left"`          // 剩余移动距离
	ActionAvailable      bool     `json:"action_available"`       // 动作是否可用
	BonusActionAvailable bool     `json:"bonus_action_available"` // 奖励动作是否可用
	ReactionAvailable    bool     `json:"reaction_available"`     // 反应是否可用
}

// ============================================================================
// 战斗API
// ============================================================================

// StartCombat 开始一场战斗遭遇
// 在指定场景中启动一场标准战斗，所有参与者投先攻骰排序
func (e *Engine) StartCombat(ctx context.Context, req StartCombatRequest) (*StartCombatResult, error) {
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

	// 如果没有指定场景ID，自动创建一个战斗场景
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

	// 验证参与者
	for _, pid := range req.ParticipantIDs {
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

	// 掷先攻
	combat.Initiative = make([]model.CombatantEntry, 0, len(req.ParticipantIDs))
	for _, actorID := range req.ParticipantIDs {
		entry, err := e.rollInitiative(game, actorID)
		if err != nil {
			return nil, err
		}
		combat.Initiative = append(combat.Initiative, entry)
	}

	// 按先攻值排序
	sort.Slice(combat.Initiative, func(i, j int) bool {
		return combat.Initiative[i].InitiativeTotal > combat.Initiative[j].InitiativeTotal
	})

	combat.CurrentIndex = 0
	combat.CurrentTurn = &model.TurnState{
		ActorID:         combat.Initiative[0].ActorID,
		Round:           1,
		ActionUsed:      false,
		BonusActionUsed: false,
		ReactionUsed:    false,
	}

	game.Combat = combat

	// 切换到战斗阶段
	game.Phase = model.PhaseCombat

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &StartCombatResult{
		Combat: combatStateToInfo(combat),
	}, nil
}

// StartCombatWithSurprise 开始带突袭判定的战斗
// 启动战斗并区分潜行方和观察方，被突袭方第一回合无法行动
func (e *Engine) StartCombatWithSurprise(ctx context.Context, req StartCombatWithSurpriseRequest) (*StartCombatWithSurpriseResult, error) {
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

	// 如果没有指定场景ID，自动创建一个战斗场景
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

	// 隐秘方进行隐匿检定，观察方进行察觉检定
	// 简化实现：隐秘方有优势
	stealthMap := make(map[model.ID]bool)
	for _, id := range req.StealthySide {
		stealthMap[id] = true
	}

	// 创建战斗状态
	combat := &model.CombatState{
		ID:      model.NewID(),
		SceneID: req.SceneID,
		Status:  model.CombatStatusActive,
		Round:   1,
		Log:     make([]model.CombatLogEntry, 0),
	}

	allParticipants := append(req.StealthySide, req.Observers...)
	combat.Initiative = make([]model.CombatantEntry, 0, len(allParticipants))

	for _, actorID := range allParticipants {
		entry, err := e.rollInitiative(game, actorID)
		if err != nil {
			return nil, err
		}

		// 被突袭的角色在第一回合无法行动
		if stealthMap[actorID] {
			entry.IsSurprised = false // 隐秘方不被突袭
		} else {
			entry.IsSurprised = true // 观察方被突袭
		}

		combat.Initiative = append(combat.Initiative, entry)
	}

	// 按先攻值排序
	sort.Slice(combat.Initiative, func(i, j int) bool {
		return combat.Initiative[i].InitiativeTotal > combat.Initiative[j].InitiativeTotal
	})

	combat.CurrentIndex = 0
	combat.CurrentTurn = &model.TurnState{
		ActorID:         combat.Initiative[0].ActorID,
		Round:           1,
		ActionUsed:      false,
		BonusActionUsed: false,
		ReactionUsed:    false,
	}

	game.Combat = combat
	game.Phase = model.PhaseCombat

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &StartCombatWithSurpriseResult{
		Combat: combatStateToInfo(combat),
	}, nil
}

// EndCombat 结束当前战斗
// 主动结束一场活跃的战斗，将游戏阶段切换回探索阶段
func (e *Engine) EndCombat(ctx context.Context, req EndCombatRequest) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return err
	}

	if err := e.checkPermission(game.Phase, OpEndCombat); err != nil {
		return err
	}

	if game.Combat == nil || game.Combat.Status != model.CombatStatusActive {
		return ErrCombatNotActive
	}

	game.Combat.Status = model.CombatStatusFinished
	game.Combat = nil
	game.Phase = model.PhaseExploration

	if err := e.saveGame(ctx, game); err != nil {
		return err
	}

	return nil
}

// GetCurrentCombat 获取当前活跃的战斗状态
// 返回当前战斗的信息，如果无活跃战斗则返回错误
func (e *Engine) GetCurrentCombat(ctx context.Context, req GetCurrentCombatRequest) (*GetCurrentCombatResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if game.Combat == nil || game.Combat.Status != model.CombatStatusActive {
		return nil, ErrCombatNotActive
	}

	combatCopy := *game.Combat
	return &GetCurrentCombatResult{
		Combat: combatStateToInfo(&combatCopy),
	}, nil
}

// NextTurn 推进到下一个角色的回合
// 将战斗回合推进到先攻序列中的下一个角色
func (e *Engine) NextTurn(ctx context.Context, req NextTurnRequest) (*NextTurnResult, error) {
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

	// 推进回合
	game.Combat.AdvanceTurn()

	// 处理突袭：如果被突袭，跳过第一回合的动作
	currentCombatant := game.Combat.GetCurrentCombatant()
	if currentCombatant != nil && currentCombatant.IsSurprised && game.Combat.Round == 1 {
		// 被突袭的角色跳过第一回合
		game.Combat.CurrentTurn.ActionUsed = true
		game.Combat.CurrentTurn.BonusActionUsed = true
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	combatCopy := *game.Combat
	return &NextTurnResult{
		Combat: combatStateToInfo(&combatCopy),
	}, nil
}

// GetCurrentTurn 获取当前回合的信息
// 返回当前行动角色的回合状态详情
func (e *Engine) GetCurrentTurn(ctx context.Context, req GetCurrentTurnRequest) (*TurnInfo, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if game.Combat == nil || game.Combat.Status != model.CombatStatusActive {
		return nil, ErrCombatNotActive
	}

	if game.Combat.CurrentTurn == nil {
		return nil, fmt.Errorf("no active turn")
	}

	actor, ok := game.GetActor(game.Combat.CurrentTurn.ActorID)
	if !ok {
		return nil, ErrNotFound
	}

	var name string
	var baseActor *model.Actor
	switch a := actor.(type) {
	case *model.PlayerCharacter:
		name = a.Name
		baseActor = &a.Actor
	case *model.NPC:
		name = a.Name
		baseActor = &a.Actor
	case *model.Enemy:
		name = a.Name
		baseActor = &a.Actor
	case *model.Companion:
		name = a.Name
		baseActor = &a.Actor
	}

	return &TurnInfo{
		ActorID:              game.Combat.CurrentTurn.ActorID,
		ActorName:            name,
		Round:                game.Combat.CurrentTurn.Round,
		InitiativePos:        game.Combat.CurrentIndex + 1,
		MovementLeft:         baseActor.Speed - game.Combat.CurrentTurn.MovementUsed,
		ActionAvailable:      !game.Combat.CurrentTurn.ActionUsed,
		BonusActionAvailable: !game.Combat.CurrentTurn.BonusActionUsed,
		ReactionAvailable:    !game.Combat.CurrentTurn.ReactionUsed,
	}, nil
}

// ExecuteAction 在当前回合执行一个动作
// 执行角色当前回合的动作（冲刺、脱离、闪避、帮助等）
func (e *Engine) ExecuteAction(ctx context.Context, req ExecuteActionRequest) (*ExecuteActionResult, error) {
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

	// 检查是否是该角色的回合
	if !game.Combat.IsActorTurn(req.ActorID) {
		return nil, ErrNotYourTurn
	}

	// 检查角色是否失能
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

	if baseActor.IsIncapacitated() {
		return nil, ErrActorIncapacitated
	}

	result := &ActionResult{
		Success: true,
		Effects: make([]EffectDetail, 0),
	}

	// 根据动作类型处理
	switch req.Action.Type {
	case model.ActionDash:
		result.Message = fmt.Sprintf("%s 执行冲刺动作", baseActor.Name)
		result.Effects = append(result.Effects, EffectDetail{
			Type:        "movement",
			Description: "额外移动速度",
			Value:       baseActor.Speed,
		})
		game.Combat.CurrentTurn.ActionUsed = true

	case model.ActionDisengage:
		result.Message = fmt.Sprintf("%s 执行脱离动作", baseActor.Name)
		result.Effects = append(result.Effects, EffectDetail{
			Type:        "disengage",
			Description: "本回合移动不会引发借机攻击",
		})
		game.Combat.CurrentTurn.ActionUsed = true

	case model.ActionDodge:
		result.Message = fmt.Sprintf("%s 执行闪避动作", baseActor.Name)
		result.Effects = append(result.Effects, EffectDetail{
			Type:        "dodge",
			Description: "本回合内攻击者对该角色有劣势",
		})
		game.Combat.CurrentTurn.ActionUsed = true

	case model.ActionHelp:
		result.Message = fmt.Sprintf("%s 执行协助动作", baseActor.Name)
		game.Combat.CurrentTurn.ActionUsed = true

	case model.ActionHide:
		result.Message = fmt.Sprintf("%s 尝试躲藏", baseActor.Name)
		game.Combat.CurrentTurn.ActionUsed = true

	case model.ActionReady:
		result.Message = fmt.Sprintf("%s 执行准备动作", baseActor.Name)
		game.Combat.CurrentTurn.ActionUsed = true

	case model.ActionSearch:
		result.Message = fmt.Sprintf("%s 进行搜索", baseActor.Name)
		game.Combat.CurrentTurn.ActionUsed = true

	default:
		result.Message = fmt.Sprintf("%s 执行动作: %s", baseActor.Name, req.Action.Type)
		game.Combat.CurrentTurn.ActionUsed = true
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	combatCopy := *game.Combat
	return &ExecuteActionResult{
		ActionResult: result,
		Combat:       combatStateToInfo(&combatCopy),
	}, nil
}

// ExecuteAttack 执行攻击动作
// 对目标执行一次完整的攻击，包括攻击掷骰（支持优势/劣势）、命中判定、伤害计算、
// 暴击/大失败判定、武器掌控效果应用，以及未命中时的擦伤处理
// 参数:
//
//	ctx - 上下文
//	req - 攻击请求，包含游戏会话ID、攻击者ID、目标ID和攻击输入（武器、优劣势、额外伤害等）
//
// 返回:
//
//	*ExecuteAttackResult - 攻击结果，包含攻击掷骰详情、命中状态、伤害结果、武器掌控效果和当前战斗状态
//	error - 战斗未激活、角色不存在或保存失败时返回错误
func (e *Engine) ExecuteAttack(ctx context.Context, req ExecuteAttackRequest) (*ExecuteAttackResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if game.Combat == nil || game.Combat.Status != model.CombatStatusActive {
		return nil, ErrCombatNotActive
	}

	// 获取攻击者
	attacker, ok := game.GetActor(req.AttackerID)
	if !ok {
		return nil, ErrNotFound
	}

	// 获取目标
	target, ok := game.GetActor(req.TargetID)
	if !ok {
		return nil, ErrNotFound
	}

	var attackerActor *model.Actor
	var targetActor *model.Actor
	switch a := attacker.(type) {
	case *model.PlayerCharacter:
		attackerActor = &a.Actor
	case *model.NPC:
		attackerActor = &a.Actor
	case *model.Enemy:
		attackerActor = &a.Actor
	case *model.Companion:
		attackerActor = &a.Actor
	}
	switch a := target.(type) {
	case *model.PlayerCharacter:
		targetActor = &a.Actor
	case *model.NPC:
		targetActor = &a.Actor
	case *model.Enemy:
		targetActor = &a.Actor
	case *model.Companion:
		targetActor = &a.Actor
	}

	// 计算攻击加值
	attackBonus := rules.CalcAttachBonus(attacker, attackerActor.AbilityScores.Strength)

	// 掷攻击骰
	var rollValue int
	if req.Attack.Advantage.Advantage {
		rollResult, _ := e.roller.RollAdvantage(0)
		rollValue = rollResult.Rolls[0].Value
	} else if req.Attack.Advantage.Disadvantage {
		rollResult, _ := e.roller.RollDisadvantage(0)
		rollValue = rollResult.Rolls[0].Value
	} else {
		rollResult, _ := e.roller.Roll("1d20")
		rollValue = rollResult.Rolls[0].Value
	}

	// 使用 rules.PerformAttackRoll 执行攻击检定
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
		Message:     fmt.Sprintf("攻击掷骰 %d (总计 %d) vs AC %d", rollValue, attackCheck.Total, targetActor.ArmorClass),
	}

	// 如果命中，计算伤害
	if attackCheck.Hit {
		damageResult, err := e.calculateAndApplyDamage(game, req.AttackerID, req.TargetID, req.Attack, attackCheck.IsCritical)
		if err != nil {
			return nil, err
		}
		attackResult.Damage = damageResult
		attackResult.Message += fmt.Sprintf(" - 命中！造成 %d 点伤害", damageResult.FinalDamage)

		// 应用武器掌控效果
		if req.Attack.WeaponMastery != "" {
			attackResult = e.applyWeaponMastery(attackResult, req.Attack.WeaponMastery, target)
		}
	} else {
		attackResult.Message += " - 未命中"

		// 未命中时应用擦伤效果
		if req.Attack.WeaponMastery == string(model.MasteryGraze) {
			// 擦伤：未命中时仍造成属性修正伤害
			abilityMod := 0
			if attacker, ok := game.GetActor(req.AttackerID); ok {
				var actor *model.Actor
				switch a := attacker.(type) {
				case *model.PlayerCharacter:
					actor = &a.Actor
				case *model.Enemy:
					actor = &a.Actor
				default:
					actor = &model.Actor{}
				}
				abilityMod = rules.AbilityModifier(actor.AbilityScores.Strength)
			}
			attackResult.GrazeDamage = abilityMod
			attackResult.Message += fmt.Sprintf("，擦伤造成 %d 点伤害", abilityMod)

			// 应用擦伤伤害
			_, _ = e.applyDamageToTarget(game, req.AttackerID, req.TargetID, abilityMod, model.DamageTypeSlashing, false)
		}
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	combatCopy := *game.Combat
	return &ExecuteAttackResult{
		AttackResult: attackResult,
		Combat:       combatStateToInfo(&combatCopy),
	}, nil
}

// ExecuteDamage 直接对角色施加伤害
// 用于处理非攻击来源的伤害（如陷阱、环境伤害、法术伤害等），
// 自动处理伤害抗性、易伤、临时HP扣除、专注检定和死亡判定
// 参数:
//
//	ctx - 上下文
//	req - 伤害请求，包含游戏会话ID、目标ID和伤害输入（伤害量、类型、来源等）
//
// 返回:
//
//	*ExecuteDamageResult - 伤害结果，包含原始伤害、抗性/易伤应用、最终伤害、目标HP变化和死亡豁免状态
//	error - 角色不存在、伤害计算失败或保存失败时返回错误
func (e *Engine) ExecuteDamage(ctx context.Context, req ExecuteDamageRequest) (*ExecuteDamageResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	damageResult, err := e.applyDamageToTarget(game, req.Damage.Source, req.TargetID, req.Damage.Amount, req.Damage.Type, false)
	if err != nil {
		return nil, err
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	result := &ExecuteDamageResult{
		DamageResult: damageResult,
	}
	if game.Combat != nil {
		combatCopy := *game.Combat
		result.Combat = combatStateToInfo(&combatCopy)
	}
	return result, nil
}

// ExecuteHealing 对角色进行治疗
// 为目标角色恢复生命值，治疗量不会超过角色最大HP。
// 如果角色处于稳定状态（Stabilized）且治疗后HP大于0，则自动移除稳定状态
// 参数:
//
//	ctx - 上下文
//	req - 治疗请求，包含游戏会话ID、目标ID和治疗量
//
// 返回:
//
//	*ExecuteHealingResult - 治疗结果，包含目标ID、实际治疗量、治疗后当前HP和消息
//	error - 角色不存在或保存失败时返回错误
func (e *Engine) ExecuteHealing(ctx context.Context, req ExecuteHealingRequest) (*ExecuteHealingResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	actor, ok := game.GetActor(req.TargetID)
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

	wasStable := baseActor.HasCondition(model.ConditionStabilized)

	// 应用治疗
	baseActor.HitPoints.Current += req.Amount
	if baseActor.HitPoints.Current > baseActor.HitPoints.Maximum {
		baseActor.HitPoints.Current = baseActor.HitPoints.Maximum
	}

	// 如果角色稳定但HP>0，移除稳定状态
	if baseActor.HitPoints.Current > 0 && wasStable {
		newConditions := make([]model.ConditionInstance, 0)
		for _, c := range baseActor.Conditions {
			if c.Type != model.ConditionStabilized {
				newConditions = append(newConditions, c)
			}
		}
		baseActor.Conditions = newConditions
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &ExecuteHealingResult{
		TargetID:  req.TargetID,
		Healed:    req.Amount,
		CurrentHP: baseActor.HitPoints.Current,
		Message:   fmt.Sprintf("恢复 %d 点HP", req.Amount),
	}, nil
}

// MoveActor 在场景中移动角色
// 将角色从当前位置移动到目标位置，处理战斗中的移动消耗
func (e *Engine) MoveActor(ctx context.Context, req MoveActorRequest) (*MoveActorResult, error) {
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

	// 计算移动距离
	from := baseActor.Position
	if from == nil {
		from = &model.Point{X: 0, Y: 0}
	}

	distance := calculateDistance(from, &req.To)
	speedRemaining := baseActor.Speed

	if game.Combat != nil && game.Combat.Status == model.CombatStatusActive {
		speedRemaining -= game.Combat.CurrentTurn.MovementUsed
	}

	if distance > speedRemaining {
		combatCopy := *game.Combat
		return &MoveActorResult{
			MoveResult: &MoveResult{
				Success:       false,
				DistanceMoved: 0,
				RemainingMove: speedRemaining,
				Message:       "移动距离不足",
			},
			Combat: combatStateToInfo(&combatCopy),
		}, nil
	}

	baseActor.Position = &req.To

	if game.Combat != nil && game.Combat.Status == model.CombatStatusActive {
		game.Combat.CurrentTurn.MovementUsed += distance
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	var combatInfo *CombatInfo
	if game.Combat != nil {
		combatCopy := *game.Combat
		combatInfo = combatStateToInfo(&combatCopy)
	}

	return &MoveActorResult{
		MoveResult: &MoveResult{
			Success:       true,
			DistanceMoved: distance,
			RemainingMove: speedRemaining - distance,
			Message:       fmt.Sprintf("移动到 (%d, %d)", req.To.X, req.To.Y),
		},
		Combat: combatInfo,
	}, nil
}

// ============================================================================
// 内部辅助函数
// ============================================================================

// rollInitiative 掷先攻
// 为指定角色投掷先攻骰（1d20+敏捷修正+先攻加值）
func (e *Engine) rollInitiative(game *model.GameState, actorID model.ID) (model.CombatantEntry, error) {
	actor, ok := game.GetActor(actorID)
	if !ok {
		return model.CombatantEntry{}, ErrNotFound
	}

	var name string
	var baseActor *model.Actor
	switch a := actor.(type) {
	case *model.PlayerCharacter:
		name = a.Name
		baseActor = &a.Actor
	case *model.NPC:
		name = a.Name
		baseActor = &a.Actor
	case *model.Enemy:
		name = a.Name
		baseActor = &a.Actor
	case *model.Companion:
		name = a.Name
		baseActor = &a.Actor
	}

	// 先攻 = 1d20 + 敏捷修正
	roll, _ := e.roller.Roll("1d20")
	dexMod := rules.AbilityModifier(baseActor.AbilityScores.Dexterity)
	initiativeBonus := dexMod + baseActor.InitiativeBonus
	total := roll.Total + initiativeBonus

	return model.CombatantEntry{
		ActorID:         actorID,
		ActorName:       name,
		InitiativeRoll:  roll.Rolls[0].Value,
		InitiativeTotal: total,
	}, nil
}

// calculateAndApplyDamage 计算并应用伤害
// 计算攻击伤害并将其应用到目标角色
func (e *Engine) calculateAndApplyDamage(game *model.GameState, attackerID, targetID model.ID, attack AttackInput, isCritical bool) (*DamageResult, error) {
	_, ok := game.GetActor(targetID)
	if !ok {
		return nil, ErrNotFound
	}

	// 获取攻击者信息
	attacker, ok := game.GetActor(attackerID)
	if !ok {
		return nil, ErrNotFound
	}

	var attackerActor *model.Actor
	var attackerPC *model.PlayerCharacter
	switch a := attacker.(type) {
	case *model.PlayerCharacter:
		attackerActor = &a.Actor
		attackerPC = a
	case *model.NPC:
		attackerActor = &a.Actor
	case *model.Enemy:
		attackerActor = &a.Actor
	case *model.Companion:
		attackerActor = &a.Actor
	}

	// 基础伤害掷骰
	roll, _ := e.roller.Roll("2d6")
	strMod := rules.AbilityModifier(attackerActor.AbilityScores.Strength)
	baseDamage := roll.Total + strMod

	// 应用职业特性伤害钩子
	if attackerPC != nil && attackerPC.FeatureHooks != nil {
		ctx := &model.DamageContext{
			BaseDamage: roll.Total,
			Bonus:      strMod,
			DamageType: model.DamageTypeSlashing, // TODO: 从武器获取
			IsMelee:    true,
			IsRanged:   false,
		}
		for _, hook := range attackerPC.FeatureHooks {
			hook.OnDamageCalc(ctx)
		}
		baseDamage = roll.Total + ctx.Bonus
	}

	if isCritical {
		baseDamage *= 2 // 暴击伤害翻倍
	}

	// 创建伤害抗性
	resistances := model.NewDamageResistances()

	// 计算最终伤害
	calc := rules.CalculateDamage(baseDamage, 0, model.DamageTypeSlashing, resistances, isCritical)

	// 应用伤害（包含专注检查）
	result, err := e.applyDamageToTarget(game, attackerID, targetID, calc.FinalDamage, model.DamageTypeSlashing, false)
	if err != nil {
		return nil, err
	}

	result.RawDamage = baseDamage
	if isCritical {
		result.RawDamage = baseDamage * 2
	}

	return result, nil
}

// applyDamageToTarget 对目标应用伤害
// 处理伤害的应用，包括抗性、易伤、临时HP和死亡判定
func (e *Engine) applyDamageToTarget(game *model.GameState, sourceID, targetID model.ID, amount int, damageType model.DamageType, isCritical bool) (*DamageResult, error) {
	target, ok := game.GetActor(targetID)
	if !ok {
		return nil, ErrNotFound
	}

	var targetActor *model.Actor
	var pc *model.PlayerCharacter
	switch a := target.(type) {
	case *model.PlayerCharacter:
		targetActor = &a.Actor
		pc = a
	case *model.NPC:
		targetActor = &a.Actor
	case *model.Enemy:
		targetActor = &a.Actor
	case *model.Companion:
		targetActor = &a.Actor
	}

	hpBefore := targetActor.HitPoints.Current

	// 创建伤害抗性
	resistances := model.NewDamageResistances()
	// 从敌人数据加载抗性
	if enemy, ok := target.(*model.Enemy); ok {
		for _, immunity := range enemy.DamageResistances {
			for _, dt := range immunity.DamageTypes {
				resistances.AddResistance(dt)
			}
		}
		for _, immunity := range enemy.DamageImmunities {
			for _, dt := range immunity.DamageTypes {
				resistances.AddImmunity(dt)
			}
		}
		for _, vuln := range enemy.DamageVulnerabilities {
			for _, dt := range vuln.DamageTypes {
				resistances.AddVulnerability(dt)
			}
		}
	}

	// 计算伤害
	calc := rules.CalculateDamage(amount, 0, damageType, resistances, isCritical)

	// 应用伤害到HP
	tempHP := targetActor.TempHitPoints
	newHP, newTempHP, _ := rules.ApplyDamage(hpBefore, tempHP, calc.FinalDamage)
	targetActor.HitPoints.Current = newHP
	targetActor.TempHitPoints = newTempHP

	result := &DamageResult{
		RawDamage:      amount,
		FinalDamage:    calc.FinalDamage,
		TargetHPBefore: hpBefore,
		TargetHPAfter:  newHP,
		Message:        fmt.Sprintf("造成 %d 点伤害", calc.FinalDamage),
	}

	// 专注检查：如果目标正在专注，受伤时需要进行专注检定
	if pc != nil && pc.Spellcasting != nil && pc.Spellcasting.IsConcentrating() {
		concResult, err := e.ConcentrationCheck(context.Background(), ConcentrationCheckRequest{
			GameID:      game.ID,
			CasterID:    pc.ID,
			DamageTaken: calc.FinalDamage,
		})
		if err == nil && !concResult.Success {
			// 专注失败，结束专注法术
			_ = e.EndConcentration(context.Background(), EndConcentrationRequest{
				GameID:   game.ID,
				CasterID: pc.ID,
			})
			result.Message += "，专注被打断"
		} else if err == nil {
			result.Message += fmt.Sprintf("，专注检定成功 (DC %d)", concResult.DC)
		}
	}

	// 检查是否死亡
	if newHP <= 0 {
		// PC需要死亡豁免
		if pc != nil {
			if amount >= targetActor.HitPoints.Maximum {
				// 即死规则：伤害超过HP最大值，立即死亡
				result.IsDead = true
				result.Message = "即死！伤害超过HP最大值"
			} else {
				// 进入濒死状态
				targetActor.HitPoints.Current = 0
				result.IsDead = false
				result.DeathSaves = &DeathSaveUpdate{
					Successes: pc.DeathSaveSuccesses,
					Failures:  pc.DeathSaveFailures,
				}
			}
		} else {
			// NPC/敌人直接死亡
			result.IsDead = true
			// 标记为战败
			if game.Combat != nil {
				combatant := game.Combat.GetCombatantByActorID(targetID)
				if combatant != nil {
					combatant.IsDefeated = true
				}
			}
		}
	}

	return result, nil
}

// applyWeaponMastery 应用武器掌控效果
func (e *Engine) applyWeaponMastery(attackResult *AttackResult, masteryType string, target any) *AttackResult {
	mastery := model.WeaponMasteryType(masteryType)
	masteryEffect := model.GetMasteryEffect(mastery)

	switch mastery {
	case model.MasteryTopple:
		// 击倒：目标进行STR或DEX豁免，失败则倒地
		attackResult.Effects = append(attackResult.Effects, AttackEffect{
			Type:        "topple",
			Description: masteryEffect.Description,
		})
		// TODO: 实现豁免检定和倒地状态应用

	case model.MasteryPush:
		// 推击：将目标推离5尺
		attackResult.Effects = append(attackResult.Effects, AttackEffect{
			Type:        "push",
			Description: masteryEffect.Description,
		})
		// TODO: 实现推离逻辑

	case model.MasteryVex:
		// 烦扰：对目标下次攻击有优势
		attackResult.Effects = append(attackResult.Effects, AttackEffect{
			Type:        "vex",
			Description: masteryEffect.Description,
		})
		// TODO: 追踪 vex 状态

	case model.MasterySlow:
		// 减缓：目标速度降低10尺
		attackResult.Effects = append(attackResult.Effects, AttackEffect{
			Type:        "slow",
			Description: masteryEffect.Description,
		})
		// TODO: 应用速度降低效果

	case model.MasterySap:
		// 钝击：目标下次攻击有劣势
		attackResult.Effects = append(attackResult.Effects, AttackEffect{
			Type:        "sap",
			Description: masteryEffect.Description,
		})
		// TODO: 应用攻击劣势状态

	case model.MasteryCleave:
		// 劈砍：击杀后可攻击邻近生物
		if attackResult.Damage != nil && attackResult.Damage.IsDead {
			attackResult.Effects = append(attackResult.Effects, AttackEffect{
				Type:        "cleave",
				Description: masteryEffect.Description,
			})
		}

	case model.MasteryNick:
		// 轻捷：额外攻击可作为附赠动作
		// 由战斗系统特殊处理，不添加效果标记
	}

	return attackResult
}

// calculateDistance 计算两点间距离（网格移动）
// 使用曼哈顿距离计算（简化版5-10-5规则）
func calculateDistance(from, to *model.Point) int {
	dx := to.X - from.X
	dy := to.Y - from.Y
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	// 简化的5-10-5规则
	return dx*5 + dy*5
}

// combatStateToInfo 将 CombatState 转换为 CombatInfo
// 用于将内部模型转换为API返回格式
func combatStateToInfo(combat *model.CombatState) *CombatInfo {
	if combat == nil {
		return nil
	}

	info := &CombatInfo{
		ID:           combat.ID,
		SceneID:      combat.SceneID,
		Status:       combat.Status,
		Round:        combat.Round,
		CurrentIndex: combat.CurrentIndex,
		Initiative:   make([]CombatantEntryInfo, 0, len(combat.Initiative)),
	}

	// 转换先攻列表
	for _, entry := range combat.Initiative {
		info.Initiative = append(info.Initiative, CombatantEntryInfo{
			ActorID:         entry.ActorID,
			ActorName:       entry.ActorName,
			InitiativeRoll:  entry.InitiativeRoll,
			InitiativeTotal: entry.InitiativeTotal,
			IsSurprised:     entry.IsSurprised,
			IsDefeated:      entry.IsDefeated,
		})
	}

	// 转换当前回合信息
	if combat.CurrentTurn != nil {
		info.CurrentTurn = &TurnInfo{
			ActorID:              combat.CurrentTurn.ActorID,
			Round:                combat.CurrentTurn.Round,
			InitiativePos:        combat.CurrentIndex + 1,
			MovementLeft:         0, // 需要从game state获取
			ActionAvailable:      !combat.CurrentTurn.ActionUsed,
			BonusActionAvailable: !combat.CurrentTurn.BonusActionUsed,
			ReactionAvailable:    !combat.CurrentTurn.ReactionUsed,
		}
	}

	return info
}
