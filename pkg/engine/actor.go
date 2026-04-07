package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/zwh8800/dnd-core/internal/model"
	"github.com/zwh8800/dnd-core/internal/rules"
)

// ActorFilter 角色过滤条件
type ActorFilter struct {
	Types   []model.ActorType
	SceneID *model.ID
	Alive   *bool
}

// ActorUpdate 角色更新内容
type ActorUpdate struct {
	AbilityScores *model.AbilityScores
	HitPoints     *HitPointUpdate
	Conditions    *ConditionUpdate
	Position      *model.Point
	SceneID       *model.ID
	Custom        map[string]any
}

// HitPointUpdate HP更新
type HitPointUpdate struct {
	Current       *int
	TempHitPoints *int
}

// ConditionUpdate 状态更新
type ConditionUpdate struct {
	Add    []model.ConditionInstance
	Remove []model.ConditionType
}

// LevelUpResult 升级结果
type LevelUpResult struct {
	OldLevel             int      `json:"old_level"`
	NewLevel             int      `json:"new_level"`
	HPGain               int      `json:"hp_gain"`
	NewFeatures          []string `json:"new_features"`
	SpellSlotsUpdated    bool     `json:"spell_slots_updated"`
	ProficiencyIncreased bool     `json:"proficiency_increased"`
	Message              string   `json:"message"`
}

// RestResult 休息结果
type RestResult struct {
	ActorResults []ActorRestResult `json:"actor_results"`
	Message      string            `json:"message"`
}

// ActorRestResult 角色休息结果
type ActorRestResult struct {
	ActorID            model.ID              `json:"actor_id"`
	HPRecovered        int                   `json:"hp_recovered"`
	HitDiceUsed        int                   `json:"hit_dice_used"`
	SpellSlotsRestored bool                  `json:"spell_slots_restored"`
	ConditionsRemoved  []model.ConditionType `json:"conditions_removed"`
	ExhaustionReduced  bool                  `json:"exhaustion_reduced"`
	AbilitiesRestored  bool                  `json:"abilities_restored"`
}

// CreatePC 创建一个新的玩家角色
func (e *Engine) CreatePC(ctx context.Context, gameID model.ID, pc *model.PlayerCharacter) (*model.PlayerCharacter, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	// 检查权限
	if err := e.checkPermission(game.Phase, OpCreatePC); err != nil {
		return nil, err
	}

	// 生成ID
	if pc.ID == "" {
		pc.ID = model.NewID()
	}

	// 计算派生值
	pc.ArmorClass = calculateArmorClass(pc)
	if pc.HitPoints.Maximum <= 0 {
		pc.HitPoints.Maximum = calculateMaxHP(pc)
		pc.HitPoints.Current = pc.HitPoints.Maximum
	}

	// 创建库存
	inventory := model.NewInventory(pc.ID)
	game.Inventories[inventory.ID] = inventory
	pc.InventoryID = inventory.ID

	// 添加到游戏
	game.PCs[pc.ID] = pc
	game.UpdatedAt = time.Now()

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	pcCopy := *pc
	return &pcCopy, nil
}

// CreateNPC 创建一个新的非玩家角色
func (e *Engine) CreateNPC(ctx context.Context, gameID model.ID, npc *model.NPC) (*model.NPC, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpCreateNPC); err != nil {
		return nil, err
	}

	if npc.ID == "" {
		npc.ID = model.NewID()
	}

	npc.ArmorClass = calculateArmorClassFromActor(&npc.Actor)

	game.NPCs[npc.ID] = npc

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	npcCopy := *npc
	return &npcCopy, nil
}

// CreateEnemy 创建一个新的敌人/怪物
func (e *Engine) CreateEnemy(ctx context.Context, gameID model.ID, enemy *model.Enemy) (*model.Enemy, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpCreateEnemy); err != nil {
		return nil, err
	}

	if enemy.ID == "" {
		enemy.ID = model.NewID()
	}

	game.Enemies[enemy.ID] = enemy

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	enemyCopy := *enemy
	return &enemyCopy, nil
}

// CreateCompanion 创建一个同伴角色
func (e *Engine) CreateCompanion(ctx context.Context, gameID model.ID, companion *model.Companion) (*model.Companion, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpCreateCompanion); err != nil {
		return nil, err
	}

	if companion.ID == "" {
		companion.ID = model.NewID()
	}

	companion.ArmorClass = calculateArmorClassFromActor(&companion.Actor)

	game.Companions[companion.ID] = companion

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	companionCopy := *companion
	return &companionCopy, nil
}

// GetActor 获取任意类型的角色信息
func (e *Engine) GetActor(ctx context.Context, gameID model.ID, actorID model.ID) (model.ActorSnapshot, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return model.ActorSnapshot{}, err
	}

	if err := e.checkPermission(game.Phase, OpGetActor); err != nil {
		return model.ActorSnapshot{}, err
	}

	actor, ok := game.GetActor(actorID)
	if !ok {
		return model.ActorSnapshot{}, ErrNotFound
	}

	switch a := actor.(type) {
	case *model.PlayerCharacter:
		return model.ActorToSnapshot(&a.Actor, model.ActorTypePC, a.Name), nil
	case *model.NPC:
		return model.ActorToSnapshot(&a.Actor, model.ActorTypeNPC, a.Name), nil
	case *model.Enemy:
		return model.ActorToSnapshot(&a.Actor, model.ActorTypeEnemy, a.Name), nil
	case *model.Companion:
		return model.ActorToSnapshot(&a.Actor, model.ActorTypeCompanion, a.Name), nil
	default:
		return model.ActorSnapshot{}, fmt.Errorf("unknown actor type")
	}
}

// GetPC 获取玩家角色的完整数据
func (e *Engine) GetPC(ctx context.Context, gameID model.ID, pcID model.ID) (*model.PlayerCharacter, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	pc, ok := game.PCs[pcID]
	if !ok {
		return nil, ErrNotFound
	}

	pcCopy := *pc
	return &pcCopy, nil
}

// UpdateActor 更新角色的部分状态
func (e *Engine) UpdateActor(ctx context.Context, gameID model.ID, actorID model.ID, update ActorUpdate) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return err
	}

	if err := e.checkPermission(game.Phase, OpUpdateActor); err != nil {
		return err
	}

	actor, ok := game.GetActor(actorID)
	if !ok {
		return ErrNotFound
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

	// 应用更新
	if update.AbilityScores != nil {
		baseActor.AbilityScores = *update.AbilityScores
	}
	if update.HitPoints != nil {
		if update.HitPoints.Current != nil {
			baseActor.HitPoints.Current = *update.HitPoints.Current
		}
		if update.HitPoints.TempHitPoints != nil {
			baseActor.TempHitPoints = *update.HitPoints.TempHitPoints
		}
	}
	if update.Conditions != nil {
		// 添加新状态
		baseActor.Conditions = append(baseActor.Conditions, update.Conditions.Add...)
		// 移除指定状态
		if len(update.Conditions.Remove) > 0 {
			newConditions := make([]model.ConditionInstance, 0)
			for _, c := range baseActor.Conditions {
				shouldRemove := false
				for _, rem := range update.Conditions.Remove {
					if c.Type == rem {
						shouldRemove = true
						break
					}
				}
				if !shouldRemove {
					newConditions = append(newConditions, c)
				}
			}
			baseActor.Conditions = newConditions
		}
	}
	if update.Position != nil {
		baseActor.Position = update.Position
	}
	if update.SceneID != nil {
		baseActor.SceneID = *update.SceneID
	}

	if err := e.saveGame(ctx, game); err != nil {
		return err
	}

	return nil
}

// RemoveActor 从游戏中移除一个角色
func (e *Engine) RemoveActor(ctx context.Context, gameID model.ID, actorID model.ID) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return err
	}

	if err := e.checkPermission(game.Phase, OpRemoveActor); err != nil {
		return err
	}

	// 检查是否在战斗中
	if game.Phase == model.PhaseCombat && game.Combat != nil {
		combatant := game.Combat.GetCombatantByActorID(actorID)
		if combatant != nil && !combatant.IsDefeated {
			return ErrInvalidState
		}
	}

	if !game.RemoveActor(actorID) {
		return ErrNotFound
	}

	if err := e.saveGame(ctx, game); err != nil {
		return err
	}

	return nil
}

// ListActors 列出游戏中的角色，可按条件过滤
func (e *Engine) ListActors(ctx context.Context, gameID model.ID, filter *ActorFilter) ([]model.ActorSnapshot, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	allActors := game.GetAllActors()

	if filter == nil {
		return allActors, nil
	}

	result := make([]model.ActorSnapshot, 0)
	for _, actor := range allActors {
		// 按类型过滤
		if len(filter.Types) > 0 {
			typeMatch := false
			for _, t := range filter.Types {
				if actor.Type == t {
					typeMatch = true
					break
				}
			}
			if !typeMatch {
				continue
			}
		}

		// 按场景ID过滤
		if filter.SceneID != nil && actor.SceneID != *filter.SceneID {
			continue
		}

		// 按存活状态过滤
		if filter.Alive != nil {
			isAlive := actor.HitPoints.Current > 0
			if isAlive != *filter.Alive {
				continue
			}
		}

		result = append(result, actor)
	}

	return result, nil
}

// AddExperience 为玩家角色添加经验值
func (e *Engine) AddExperience(ctx context.Context, gameID model.ID, pcID model.ID, xp int) (*LevelUpResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpAddExperience); err != nil {
		return nil, err
	}

	pc, ok := game.PCs[pcID]
	if !ok {
		return nil, ErrNotFound
	}

	oldLevel := pc.TotalLevel
	pc.Experience += xp

	// 检查是否升级
	newLevel := rules.GetLevelByXP(pc.Experience)
	if newLevel > oldLevel {
		pc.TotalLevel = newLevel
		// 升级逻辑在 LevelUp 方法中处理
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return nil, nil
}

// LevelUp 手动触发玩家角色升级
func (e *Engine) LevelUp(ctx context.Context, gameID model.ID, pcID model.ID, classChoice string) (*LevelUpResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpLevelUp); err != nil {
		return nil, err
	}

	pc, ok := game.PCs[pcID]
	if !ok {
		return nil, ErrNotFound
	}

	oldLevel := pc.TotalLevel
	newLevel := oldLevel + 1

	// 检查XP是否足够
	requiredXP := rules.GetXPForLevel(newLevel)
	if pc.Experience < requiredXP {
		return nil, fmt.Errorf("insufficient experience: need %d, have %d", requiredXP, pc.Experience)
	}

	// 计算HP增长
	hpGain := 0
	classToLevel := classChoice
	if classToLevel == "" && len(pc.Classes) > 0 {
		classToLevel = pc.Classes[0].ClassName
	}

	hitDiceType := rules.HitDiceByClass[classToLevel]
	if hitDiceType == 0 {
		hitDiceType = 8 // 默认d8
	}

	// 简化的HP计算：取平均值+CON修正
	conMod := rules.AbilityModifier(pc.AbilityScores.Constitution)
	hpGain = (hitDiceType / 2) + 1 + conMod
	if hpGain < 1 {
		hpGain = 1
	}

	pc.HitPoints.Maximum += hpGain
	pc.HitPoints.Current += hpGain
	pc.TotalLevel = newLevel

	// 更新熟练加值检查
	oldProfBonus := rules.ProficiencyBonus(oldLevel)
	newProfBonus := rules.ProficiencyBonus(newLevel)
	proficiencyIncreased := newProfBonus > oldProfBonus

	result := &LevelUpResult{
		OldLevel:             oldLevel,
		NewLevel:             newLevel,
		HPGain:               hpGain,
		SpellSlotsUpdated:    pc.Spellcasting != nil,
		ProficiencyIncreased: proficiencyIncreased,
		Message:              fmt.Sprintf("升级到等级 %d！HP +%d", newLevel, hpGain),
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return result, nil
}

// ShortRest 为指定角色执行短休
func (e *Engine) ShortRest(ctx context.Context, gameID model.ID, actorIDs []model.ID) (*RestResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpShortRest); err != nil {
		return nil, err
	}

	result := &RestResult{
		ActorResults: make([]ActorRestResult, 0),
		Message:      "短休完成",
	}

	for _, actorID := range actorIDs {
		actorResult, err := e.processShortRest(game, actorID)
		if err != nil {
			return nil, fmt.Errorf("short rest failed for actor %s: %w", actorID, err)
		}
		result.ActorResults = append(result.ActorResults, actorResult)
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return result, nil
}

// StartLongRest 开始长休过程
func (e *Engine) StartLongRest(ctx context.Context, gameID model.ID, actorIDs []model.ID) (*RestResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpStartLongRest); err != nil {
		return nil, err
	}

	// 创建长休状态
	restState := model.NewLongRest(actorIDs)
	restState.Start()
	game.ActiveRest = restState

	// 切换到休息阶段
	game.Phase = model.PhaseRest

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &RestResult{
		Message: "长休开始，需要8小时才能完成",
	}, nil
}

// EndLongRest 结束长休并应用恢复效果
func (e *Engine) EndLongRest(ctx context.Context, gameID model.ID) (*RestResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpEndLongRest); err != nil {
		return nil, err
	}

	if game.ActiveRest == nil || game.ActiveRest.Type != model.RestTypeLong {
		return nil, ErrInvalidState
	}

	result := &RestResult{
		ActorResults: make([]ActorRestResult, 0),
	}

	// 对每个参与者应用长休效果
	for _, actorID := range game.ActiveRest.ParticipantIDs {
		actorResult, err := e.processLongRestRecovery(game, actorID)
		if err != nil {
			return nil, fmt.Errorf("long rest recovery failed for actor %s: %w", actorID, err)
		}
		result.ActorResults = append(result.ActorResults, actorResult)
	}

	// 完成长休
	game.ActiveRest.Complete()
	game.Phase = model.PhaseExploration
	game.ActiveRest = nil

	result.Message = "长休完成，队伍完全恢复"

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return result, nil
}

// processShortRest 处理单个角色的短休
func (e *Engine) processShortRest(game *model.GameState, actorID model.ID) (ActorRestResult, error) {
	result := ActorRestResult{
		ActorID: actorID,
	}

	// 查找角色
	actor, ok := game.GetActor(actorID)
	if !ok {
		return result, ErrNotFound
	}

	var baseActor *model.Actor
	var pc *model.PlayerCharacter
	switch a := actor.(type) {
	case *model.PlayerCharacter:
		baseActor = &a.Actor
		pc = a
	case *model.Companion:
		baseActor = &a.Actor
	case *model.NPC:
		baseActor = &a.Actor
	case *model.Enemy:
		baseActor = &a.Actor
	}

	// 短休可以掷生命骰恢复HP（仅PC）
	if pc != nil && len(pc.HitDice) > 0 {
		// 简化的短休逻辑：恢复一部分HP
		hpRecovered := 0
		for i := range pc.HitDice {
			if pc.HitDice[i].Used < pc.HitDice[i].Total {
				// 使用一个生命骰
				pc.HitDice[i].Used++
				// 简化：取平均值
				hpRecovered += (pc.HitDice[i].DiceType / 2) + 1
				result.HitDiceUsed++
				break // 简化：只用一个生命骰
			}
		}
		conMod := rules.AbilityModifier(pc.AbilityScores.Constitution)
		hpRecovered += conMod
		if hpRecovered > 0 {
			baseActor.HitPoints.Current += hpRecovered
			if baseActor.HitPoints.Current > baseActor.HitPoints.Maximum {
				baseActor.HitPoints.Current = baseActor.HitPoints.Maximum
			}
			result.HPRecovered = hpRecovered
		}
	}

	return result, nil
}

// processLongRestRecovery 处理单个角色的长休恢复
func (e *Engine) processLongRestRecovery(game *model.GameState, actorID model.ID) (ActorRestResult, error) {
	result := ActorRestResult{
		ActorID: actorID,
	}

	actor, ok := game.GetActor(actorID)
	if !ok {
		return result, ErrNotFound
	}

	var baseActor *model.Actor
	var pc *model.PlayerCharacter
	switch a := actor.(type) {
	case *model.PlayerCharacter:
		baseActor = &a.Actor
		pc = a
	case *model.Companion:
		baseActor = &a.Actor
	case *model.NPC:
		baseActor = &a.Actor
	case *model.Enemy:
		baseActor = &a.Actor
	}

	// 恢复所有HP
	hpRecovered := baseActor.HitPoints.Maximum - baseActor.HitPoints.Current
	baseActor.HitPoints.Current = baseActor.HitPoints.Maximum
	baseActor.TempHitPoints = 0
	result.HPRecovered = hpRecovered

	// 恢复法术位（PC）
	if pc != nil && pc.Spellcasting != nil {
		pc.Spellcasting.Slots.RestoreAll()
		pc.Spellcasting.ConcentrationSpell = ""
		result.SpellSlotsRestored = true
	}

	// 恢复生命骰
	if pc != nil {
		recoveryAmount := rules.CalculateHitDiceRecovery(pc.TotalLevel)
		recovered := 0
		for i := range pc.HitDice {
			available := pc.HitDice[i].Total - pc.HitDice[i].Used
			canRecover := recoveryAmount - recovered
			if canRecover > available {
				canRecover = available
			}
			pc.HitDice[i].Used -= canRecover
			recovered += canRecover
		}
		result.HitDiceUsed = -recovered // 负数表示恢复
	}

	// 减少力竭
	if baseActor.Exhaustion > 0 {
		baseActor.Exhaustion = rules.CalculateExhaustionReduction(baseActor.Exhaustion)
		result.ExhaustionReduced = true
	}

	// 移除某些状态效果（长休结束时移除）
	removableConditions := []model.ConditionType{
		model.ConditionPoisoned,
		model.ConditionFrightened,
		model.ConditionCharmed,
	}
	for _, condType := range removableConditions {
		for i := len(baseActor.Conditions) - 1; i >= 0; i-- {
			if baseActor.Conditions[i].Type == condType {
				result.ConditionsRemoved = append(result.ConditionsRemoved, condType)
				baseActor.Conditions = append(baseActor.Conditions[:i], baseActor.Conditions[i+1:]...)
			}
		}
	}

	return result, nil
}

// calculateArmorClass 计算PC的护甲等级
func calculateArmorClass(pc *model.PlayerCharacter) int {
	// 简化实现：默认10+DEX修正
	dexMod := rules.AbilityModifier(pc.AbilityScores.Dexterity)
	return 10 + dexMod
}

// calculateArmorClassFromActor 从Actor计算AC
func calculateArmorClassFromActor(actor *model.Actor) int {
	dexMod := rules.AbilityModifier(actor.AbilityScores.Dexterity)
	return 10 + dexMod
}

// calculateMaxHP 计算最大HP
func calculateMaxHP(pc *model.PlayerCharacter) int {
	if len(pc.Classes) == 0 {
		return 10 // 默认值
	}

	hp := 0
	for _, cl := range pc.Classes {
		hitDiceType := rules.HitDiceByClass[cl.ClassName]
		if hitDiceType == 0 {
			hitDiceType = 8
		}
		// 第一级取最大值，之后取平均
		hp += hitDiceType + rules.AbilityModifier(pc.AbilityScores.Constitution)
	}
	return hp
}
