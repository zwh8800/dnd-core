package engine

import (
	"context"
	"fmt"
	"sort"

	"github.com/zwh8800/dnd-core/internal/model"
	"github.com/zwh8800/dnd-core/internal/rules"
)

// ActionInput 动作输入
type ActionInput struct {
	Type    model.ActionType `json:"type"`
	Details map[string]any   `json:"details,omitempty"`
}

// ActionResult 动作执行结果
type ActionResult struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Roll    *model.DiceResult `json:"roll,omitempty"`
	Effects []EffectDetail    `json:"effects"`
}

// EffectDetail 效果详情
type EffectDetail struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Value       int    `json:"value,omitempty"`
}

// AttackInput 攻击输入
type AttackInput struct {
	WeaponID    *model.ID          `json:"weapon_id,omitempty"`
	SpellID     *string            `json:"spell_id,omitempty"`
	IsUnarmed   bool               `json:"is_unarmed"`
	IsOffHand   bool               `json:"is_off_hand"`
	Advantage   model.RollModifier `json:"advantage"`
	ExtraDamage []DamageInput      `json:"extra_damage,omitempty"`
}

// AttackResult 攻击结果
type AttackResult struct {
	Roll        *model.DiceResult `json:"roll"`
	AttackTotal int               `json:"attack_total"`
	TargetAC    int               `json:"target_ac"`
	Hit         bool              `json:"hit"`
	IsCritical  bool              `json:"is_critical"`
	IsFumble    bool              `json:"is_fumble"`
	Damage      *DamageResult     `json:"damage,omitempty"`
	Message     string            `json:"message"`
}

// DamageInput 伤害输入
type DamageInput struct {
	Amount int              `json:"amount"`
	Type   model.DamageType `json:"type"`
	Dice   string           `json:"dice,omitempty"`
	Source model.ID         `json:"source"`
}

// DamageResult 伤害结果
type DamageResult struct {
	RawDamage       int                `json:"raw_damage"`
	Resistances     []model.DamageType `json:"resistances_applied"`
	Vulnerabilities []model.DamageType `json:"vulnerabilities_applied"`
	FinalDamage     int                `json:"final_damage"`
	TargetHPBefore  int                `json:"target_hp_before"`
	TargetHPAfter   int                `json:"target_hp_after"`
	IsDead          bool               `json:"is_dead"`
	IsStabilized    bool               `json:"is_stabilized"`
	DeathSaves      *DeathSaveUpdate   `json:"death_saves,omitempty"`
	Message         string             `json:"message"`
}

// DeathSaveUpdate 死亡豁免更新
type DeathSaveUpdate struct {
	Successes int  `json:"successes"`
	Failures  int  `json:"failures"`
	IsStable  bool `json:"is_stable"`
}

// HealResult 治疗结果
type HealResult struct {
	Amount    int    `json:"amount"`
	HPBefore  int    `json:"hp_before"`
	HPAfter   int    `json:"hp_after"`
	WasStable bool   `json:"was_stable"`
	Message   string `json:"message"`
}

// MoveResult 移动结果
type MoveResult struct {
	Success       bool   `json:"success"`
	DistanceMoved int    `json:"distance_moved"`
	RemainingMove int    `json:"remaining_move"`
	Message       string `json:"message"`
}

// TurnInfo 回合信息
type TurnInfo struct {
	ActorID              model.ID `json:"actor_id"`
	ActorName            string   `json:"actor_name"`
	Round                int      `json:"round"`
	InitiativePos        int      `json:"initiative_position"`
	MovementLeft         int      `json:"movement_left"`
	ActionAvailable      bool     `json:"action_available"`
	BonusActionAvailable bool     `json:"bonus_action_available"`
	ReactionAvailable    bool     `json:"reaction_available"`
}

// StartCombat 开始一场战斗遭遇
func (e *Engine) StartCombat(ctx context.Context, gameID model.ID, sceneID model.ID, participantIDs []model.ID) (*model.CombatState, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpStartCombat); err != nil {
		return nil, err
	}

	if game.Combat != nil && game.Combat.Status == model.CombatStatusActive {
		return nil, ErrCombatAlreadyActive
	}

	// 验证参与者
	for _, pid := range participantIDs {
		if _, ok := game.GetActor(pid); !ok {
			return nil, fmt.Errorf("actor %s not found", pid)
		}
	}

	// 创建战斗状态
	combat := &model.CombatState{
		ID:      model.NewID(),
		SceneID: sceneID,
		Status:  model.CombatStatusActive,
		Round:   1,
		Log:     make([]model.CombatLogEntry, 0),
	}

	// 掷先攻
	combat.Initiative = make([]model.CombatantEntry, 0, len(participantIDs))
	for _, actorID := range participantIDs {
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

	combatCopy := *combat
	return &combatCopy, nil
}

// StartCombatWithSurprise 开始带突袭判定的战斗
func (e *Engine) StartCombatWithSurprise(ctx context.Context, gameID model.ID, sceneID model.ID, stealthySide []model.ID, observers []model.ID) (*model.CombatState, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpStartCombat); err != nil {
		return nil, err
	}

	if game.Combat != nil && game.Combat.Status == model.CombatStatusActive {
		return nil, ErrCombatAlreadyActive
	}

	// 隐秘方进行隐匿检定，观察方进行察觉检定
	// 简化实现：隐秘方有优势
	stealthMap := make(map[model.ID]bool)
	for _, id := range stealthySide {
		stealthMap[id] = true
	}

	// 创建战斗状态
	combat := &model.CombatState{
		ID:      model.NewID(),
		SceneID: sceneID,
		Status:  model.CombatStatusActive,
		Round:   1,
		Log:     make([]model.CombatLogEntry, 0),
	}

	allParticipants := append(stealthySide, observers...)
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

	combatCopy := *combat
	return &combatCopy, nil
}

// EndCombat 结束当前战斗
func (e *Engine) EndCombat(ctx context.Context, gameID model.ID) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
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
func (e *Engine) GetCurrentCombat(ctx context.Context, gameID model.ID) (*model.CombatState, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	if game.Combat == nil || game.Combat.Status != model.CombatStatusActive {
		return nil, ErrCombatNotActive
	}

	combatCopy := *game.Combat
	return &combatCopy, nil
}

// NextTurn 推进到下一个角色的回合
func (e *Engine) NextTurn(ctx context.Context, gameID model.ID) (*model.CombatState, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
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
	return &combatCopy, nil
}

// GetCurrentTurn 获取当前回合的信息
func (e *Engine) GetCurrentTurn(ctx context.Context, gameID model.ID) (*TurnInfo, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, gameID)
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
func (e *Engine) ExecuteAction(ctx context.Context, gameID model.ID, actorID model.ID, action ActionInput) (*ActionResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
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
	if !game.Combat.IsActorTurn(actorID) {
		return nil, ErrNotYourTurn
	}

	// 检查角色是否失能
	actor, ok := game.GetActor(actorID)
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
	switch action.Type {
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
		result.Message = fmt.Sprintf("%s 执行动作: %s", baseActor.Name, action.Type)
		game.Combat.CurrentTurn.ActionUsed = true
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return result, nil
}

// ExecuteAttack 执行攻击动作
func (e *Engine) ExecuteAttack(ctx context.Context, gameID model.ID, attackerID model.ID, targetID model.ID, attack AttackInput) (*AttackResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	if game.Combat == nil || game.Combat.Status != model.CombatStatusActive {
		return nil, ErrCombatNotActive
	}

	// 获取攻击者
	attacker, ok := game.GetActor(attackerID)
	if !ok {
		return nil, ErrNotFound
	}

	// 获取目标
	target, ok := game.GetActor(targetID)
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

	// 计算攻击加值（简化：使用熟练加值+属性修正）
	level := 1 // 简化
	profBonus := rules.ProficiencyBonus(level)
	attackBonus := profBonus + rules.AbilityModifier(attackerActor.AbilityScores.Strength)

	// 掷攻击骰
	var rollResult *model.DiceResult
	if attack.Advantage.Advantage {
		rollResult, _ = e.roller.RollAdvantage(0)
	} else if attack.Advantage.Disadvantage {
		rollResult, _ = e.roller.RollDisadvantage(0)
	} else {
		rollResult, _ = e.roller.Roll("1d20")
	}

	attackTotal := rollResult.Total + attackBonus
	isNat20 := rollResult.Rolls[0].Value == 20
	isNat1 := rollResult.Rolls[0].Value == 1

	// 判断命中
	hit := attackTotal >= targetActor.ArmorClass || isNat20
	if isNat1 {
		hit = false
	}

	result := &AttackResult{
		Roll:        rollResult,
		AttackTotal: attackTotal,
		TargetAC:    targetActor.ArmorClass,
		Hit:         hit,
		IsCritical:  isNat20 && hit,
		IsFumble:    isNat1,
		Message:     fmt.Sprintf("攻击掷骰 %d (总计 %d) vs AC %d", rollResult.Rolls[0].Value, attackTotal, targetActor.ArmorClass),
	}

	// 如果命中，计算伤害
	if hit {
		damageResult, err := e.calculateAndApplyDamage(game, attackerID, targetID, attack, isNat20 && hit)
		if err != nil {
			return nil, err
		}
		result.Damage = damageResult
		result.Message += fmt.Sprintf(" - 命中！造成 %d 点伤害", damageResult.FinalDamage)
	} else {
		result.Message += " - 未命中"
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return result, nil
}

// ExecuteDamage 直接对角色施加伤害
func (e *Engine) ExecuteDamage(ctx context.Context, gameID model.ID, targetID model.ID, damage DamageInput) (*DamageResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	return e.applyDamageToTarget(game, damage.Source, targetID, damage.Amount, damage.Type, false)
}

// ExecuteHealing 对角色进行治疗
func (e *Engine) ExecuteHealing(ctx context.Context, gameID model.ID, targetID model.ID, amount int) (*HealResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	actor, ok := game.GetActor(targetID)
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

	hpBefore := baseActor.HitPoints.Current
	wasStable := baseActor.HasCondition(model.ConditionStabilized)

	// 应用治疗
	baseActor.HitPoints.Current += amount
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

	return &HealResult{
		Amount:    amount,
		HPBefore:  hpBefore,
		HPAfter:   baseActor.HitPoints.Current,
		WasStable: wasStable,
		Message:   fmt.Sprintf("恢复 %d 点HP", amount),
	}, nil
}

// MoveActor 在场景中移动角色
func (e *Engine) MoveActor(ctx context.Context, gameID model.ID, actorID model.ID, to model.Point) (*MoveResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	actor, ok := game.GetActor(actorID)
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

	distance := calculateDistance(from, &to)
	speedRemaining := baseActor.Speed

	if game.Combat != nil && game.Combat.Status == model.CombatStatusActive {
		speedRemaining -= game.Combat.CurrentTurn.MovementUsed
	}

	if distance > speedRemaining {
		return &MoveResult{
			Success:       false,
			DistanceMoved: 0,
			RemainingMove: speedRemaining,
			Message:       "移动距离不足",
		}, nil
	}

	baseActor.Position = &to

	if game.Combat != nil && game.Combat.Status == model.CombatStatusActive {
		game.Combat.CurrentTurn.MovementUsed += distance
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &MoveResult{
		Success:       true,
		DistanceMoved: distance,
		RemainingMove: speedRemaining - distance,
		Message:       fmt.Sprintf("移动到 (%d, %d)", to.X, to.Y),
	}, nil
}

// rollInitiative 掷先攻
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
func (e *Engine) calculateAndApplyDamage(game *model.GameState, attackerID, targetID model.ID, attack AttackInput, isCritical bool) (*DamageResult, error) {
	_, ok := game.GetActor(targetID)
	if !ok {
		return nil, ErrNotFound
	}

	// 简化伤害计算：2d6+属性修正
	roll, _ := e.roller.Roll("2d6")
	strMod := 0
	if attacker, ok := game.GetActor(attackerID); ok {
		var attackerActor *model.Actor
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
		strMod = rules.AbilityModifier(attackerActor.AbilityScores.Strength)
	}

	baseDamage := roll.Total + strMod
	if isCritical {
		baseDamage *= 2 // 暴击伤害翻倍
	}

	// 创建伤害抗性（简化）
	resistances := model.NewDamageResistances()

	// 计算最终伤害
	calc := rules.CalculateDamage(baseDamage, 0, model.DamageTypeSlashing, resistances, isCritical)

	// 应用伤害
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
		for _, dt := range enemy.DamageResistances {
			resistances.AddResistance(dt)
		}
		for _, dt := range enemy.DamageImmunities {
			resistances.AddImmunity(dt)
		}
		for _, dt := range enemy.DamageVulnerabilities {
			resistances.AddVulnerability(dt)
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

// calculateDistance 计算两点间距离（网格移动）
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
