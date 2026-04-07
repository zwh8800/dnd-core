package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/zwh8800/dnd-core/pkg/data"
	"github.com/zwh8800/dnd-core/pkg/model"
	"github.com/zwh8800/dnd-core/pkg/rules"
)

// AbilityScoresInput 属性值输入
type AbilityScoresInput struct {
	Strength     int `json:"strength"`     // 力量
	Dexterity    int `json:"dexterity"`    // 敏捷
	Constitution int `json:"constitution"` // 体质
	Intelligence int `json:"intelligence"` // 智力
	Wisdom       int `json:"wisdom"`       // 感知
	Charisma     int `json:"charisma"`     // 魅力
}

// PlayerCharacterInput 玩家角色创建输入
type PlayerCharacterInput struct {
	Name          string             `json:"name"`                 // 角色名称
	Race          string             `json:"race"`                 // 种族
	Background    string             `json:"background"`           // 背景
	Class         string             `json:"class"`                // 职业
	Level         int                `json:"level"`                // 等级
	Alignment     string             `json:"alignment"`            // 阵营
	AbilityScores AbilityScoresInput `json:"ability_scores"`       // 属性值
	HitPoints     int                `json:"hit_points,omitempty"` // 初始HP（可选，不填则自动计算）
}

// NPCInput NPC创建输入
type NPCInput struct {
	Name          string             `json:"name"`           // NPC名称
	Description   string             `json:"description"`    // 描述
	Size          model.Size         `json:"size"`           // 体型
	Speed         int                `json:"speed"`          // 速度（尺）
	AbilityScores AbilityScoresInput `json:"ability_scores"` // 属性值
}

// EnemyInput 敌人创建输入
type EnemyInput struct {
	Name            string             `json:"name"`             // 敌人名称
	Description     string             `json:"description"`      // 描述
	Size            model.Size         `json:"size"`             // 体型
	Speed           int                `json:"speed"`            // 速度（尺）
	AbilityScores   AbilityScoresInput `json:"ability_scores"`   // 属性值
	ChallengeRating float64            `json:"challenge_rating"` // 挑战等级
	HitPoints       int                `json:"hit_points"`       // 生命值
	ArmorClass      int                `json:"armor_class"`      // 护甲等级
}

// CompanionInput 同伴创建输入
type CompanionInput struct {
	Name          string             `json:"name"`           // 同伴名称
	Description   string             `json:"description"`    // 描述
	Size          model.Size         `json:"size"`           // 体型
	Speed         int                `json:"speed"`          // 速度（尺）
	AbilityScores AbilityScoresInput `json:"ability_scores"` // 属性值
	LeaderID      model.ID           `json:"leader_id"`      // 领导者ID
}

// ActorInfo 角色基本信息
type ActorInfo struct {
	ID         model.ID        `json:"id"`          // 角色唯一标识
	Type       model.ActorType `json:"type"`        // 角色类型
	Name       string          `json:"name"`        // 角色名称
	HitPoints  model.HitPoints `json:"hit_points"`  // 生命值
	TempHP     int             `json:"temp_hp"`     // 临时HP
	ArmorClass int             `json:"armor_class"` // 护甲等级
	Speed      int             `json:"speed"`       // 移动速度
	Conditions []string        `json:"conditions"`  // 状态效果列表
	Exhaustion int             `json:"exhaustion"`  // 力竭等级
	SceneID    model.ID        `json:"scene_id"`    // 所在场景ID
	Position   *model.Point    `json:"position"`    // 位置坐标
}

// PlayerCharacterInfo 玩家角色完整信息
type PlayerCharacterInfo struct {
	ID               model.ID           `json:"id"`                // 角色唯一标识
	Name             string             `json:"name"`              // 角色名称
	Race             string             `json:"race"`              // 种族
	Background       string             `json:"background"`        // 背景
	Classes          []ClassInfo        `json:"classes"`           // 职业信息
	TotalLevel       int                `json:"total_level"`       // 总等级
	Experience       int                `json:"experience"`        // 经验值
	AbilityScores    AbilityScoresInput `json:"ability_scores"`    // 属性值
	HitPoints        model.HitPoints    `json:"hit_points"`        // 生命值
	ArmorClass       int                `json:"armor_class"`       // 护甲等级
	Speed            int                `json:"speed"`             // 移动速度
	ProficiencyBonus int                `json:"proficiency_bonus"` // 熟练加值
	Features         []string           `json:"features"`          // 特性列表
	RacialTraits     []string           `json:"racial_traits"`     // 种族特性
}

// ClassInfo 职业信息
type ClassInfo struct {
	Class      model.ClassID `json:"class"`              // 职业ID
	ClassLevel int           `json:"class_level"`        // 职业等级
	Features   []string      `json:"features,omitempty"` // 职业特性
}

// ActorFilter 角色过滤条件
type ActorFilter struct {
	Types   []model.ActorType `json:"types,omitempty"`    // 角色类型过滤列表
	SceneID *model.ID         `json:"scene_id,omitempty"` // 场景ID过滤
	Alive   *bool             `json:"alive,omitempty"`    // 存活状态过滤
}

// ActorUpdate 角色更新内容
type ActorUpdate struct {
	AbilityScores *AbilityScoresInput `json:"ability_scores,omitempty"` // 属性值更新
	HitPoints     *HitPointUpdate     `json:"hit_points,omitempty"`     // HP更新
	Conditions    *ConditionUpdate    `json:"conditions,omitempty"`     // 状态效果更新
	Position      *model.Point        `json:"position,omitempty"`       // 位置更新
	SceneID       *model.ID           `json:"scene_id,omitempty"`       // 场景ID更新
	Custom        map[string]any      `json:"custom,omitempty"`         // 自定义字段
}

// HitPointUpdate HP更新
type HitPointUpdate struct {
	Current       *int `json:"current,omitempty"`         // 当前HP
	TempHitPoints *int `json:"temp_hit_points,omitempty"` // 临时HP
}

// ConditionUpdate 状态效果更新
type ConditionUpdate struct {
	Add    []model.ConditionInstance `json:"add,omitempty"`    // 添加的状态效果
	Remove []model.ConditionType     `json:"remove,omitempty"` // 移除的状态效果类型
}

// LevelUpResult 升级结果
type LevelUpResult struct {
	OldLevel             int      `json:"old_level"`             // 原等级
	NewLevel             int      `json:"new_level"`             // 新等级
	HPGain               int      `json:"hp_gain"`               // HP增长值
	NewFeatures          []string `json:"new_features"`          // 新获得特性
	SpellSlotsUpdated    bool     `json:"spell_slots_updated"`   // 法术位是否更新
	ProficiencyIncreased bool     `json:"proficiency_increased"` // 熟练加值是否增加
	Message              string   `json:"message"`               // 人类可读消息
}

// RestResult 休息结果
type RestResult struct {
	ActorResults []ActorRestResult `json:"actor_results"` // 各角色休息结果
	Message      string            `json:"message"`       // 人类可读消息
}

// ActorRestResult 角色休息结果
type ActorRestResult struct {
	ActorID            model.ID              `json:"actor_id"`             // 角色ID
	HPRecovered        int                   `json:"hp_recovered"`         // 恢复的HP
	HitDiceUsed        int                   `json:"hit_dice_used"`        // 使用的生命骰数量
	SpellSlotsRestored bool                  `json:"spell_slots_restored"` // 法术位是否恢复
	ConditionsRemoved  []model.ConditionType `json:"conditions_removed"`   // 移除的状态效果
	ExhaustionReduced  bool                  `json:"exhaustion_reduced"`   // 力竭是否减少
	AbilitiesRestored  bool                  `json:"abilities_restored"`   // 能力是否恢复
}

// Request 结构体定义

// CreatePCRequest 创建玩家角色请求
type CreatePCRequest struct {
	GameID model.ID              `json:"game_id"` // 游戏会话ID
	PC     *PlayerCharacterInput `json:"pc"`      // 玩家角色创建参数
}

// CreatePCResult 创建玩家角色结果
type CreatePCResult struct {
	Actor *ActorInfo `json:"actor"` // 创建的角色信息
}

// CreateNPCRequest 创建NPC请求
type CreateNPCRequest struct {
	GameID model.ID  `json:"game_id"` // 游戏会话ID
	NPC    *NPCInput `json:"npc"`     // NPC创建参数
}

// CreateNPCResult 创建NPC结果
type CreateNPCResult struct {
	Actor *ActorInfo `json:"actor"` // 创建的角色信息
}

// CreateEnemyRequest 创建敌人请求
type CreateEnemyRequest struct {
	GameID model.ID    `json:"game_id"` // 游戏会话ID
	Enemy  *EnemyInput `json:"enemy"`   // 敌人创建参数
}

// CreateEnemyResult 创建敌人结果
type CreateEnemyResult struct {
	Actor *ActorInfo `json:"actor"` // 创建的角色信息
}

// CreateCompanionRequest 创建同伴请求
type CreateCompanionRequest struct {
	GameID    model.ID        `json:"game_id"`   // 游戏会话ID
	Companion *CompanionInput `json:"companion"` // 同伴创建参数
}

// CreateCompanionResult 创建同伴结果
type CreateCompanionResult struct {
	Actor *ActorInfo `json:"actor"` // 创建的角色信息
}

// GetActorRequest 获取角色请求
type GetActorRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID
	ActorID model.ID `json:"actor_id"` // 角色ID
}

// GetActorResult 获取角色结果
type GetActorResult struct {
	Actor *ActorInfo `json:"actor"` // 角色信息
}

// GetPCRequest 获取玩家角色请求
type GetPCRequest struct {
	GameID model.ID `json:"game_id"` // 游戏会话ID
	PCID   model.ID `json:"pc_id"`   // 玩家角色ID
}

// GetPCResult 获取玩家角色结果
type GetPCResult struct {
	PC *PlayerCharacterInfo `json:"pc"` // 玩家角色完整信息
}

// UpdateActorRequest 更新角色请求
type UpdateActorRequest struct {
	GameID  model.ID    `json:"game_id"`  // 游戏会话ID
	ActorID model.ID    `json:"actor_id"` // 角色ID
	Update  ActorUpdate `json:"update"`   // 更新内容
}

// RemoveActorRequest 移除角色请求
type RemoveActorRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID
	ActorID model.ID `json:"actor_id"` // 角色ID
}

// ListActorsRequest 列出角色请求
type ListActorsRequest struct {
	GameID model.ID     `json:"game_id"` // 游戏会话ID
	Filter *ActorFilter `json:"filter"`  // 过滤条件（可选）
}

// ListActorsResult 列出角色结果
type ListActorsResult struct {
	Actors []ActorInfo `json:"actors"` // 角色列表
}

// AddExperienceRequest 添加经验值请求
type AddExperienceRequest struct {
	GameID model.ID `json:"game_id"` // 游戏会话ID
	PCID   model.ID `json:"pc_id"`   // 玩家角色ID
	XP     int      `json:"xp"`      // 添加的经验值
}

// AddExperienceResult 添加经验值结果
type AddExperienceResult struct {
	LeveledUp bool `json:"leveled_up"` // 是否升级
	OldLevel  int  `json:"old_level"`  // 原等级
	NewLevel  int  `json:"new_level"`  // 新等级
}

// LevelUpRequest 升级请求
type LevelUpRequest struct {
	GameID      model.ID `json:"game_id"`      // 游戏会话ID
	PCID        model.ID `json:"pc_id"`        // 玩家角色ID
	ClassChoice string   `json:"class_choice"` // 升级职业选择
}

// ShortRestRequest 短休请求
type ShortRestRequest struct {
	GameID   model.ID   `json:"game_id"`   // 游戏会话ID
	ActorIDs []model.ID `json:"actor_ids"` // 参与短休的角色ID列表
}

// StartLongRestRequest 开始长休请求
type StartLongRestRequest struct {
	GameID   model.ID   `json:"game_id"`   // 游戏会话ID
	ActorIDs []model.ID `json:"actor_ids"` // 参与长休的角色ID列表
}

// EndLongRestRequest 结束长休请求
type EndLongRestRequest struct {
	GameID model.ID `json:"game_id"` // 游戏会话ID
}

// CreatePC 创建一个新的玩家角色
func (e *Engine) CreatePC(ctx context.Context, req CreatePCRequest) (*CreatePCResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	// 检查权限
	if err := e.checkPermission(game.Phase, OpCreatePC); err != nil {
		return nil, err
	}

	if req.PC == nil {
		return nil, fmt.Errorf("pc input is required")
	}

	// 构建 PlayerCharacter
	pc := &model.PlayerCharacter{
		Actor: model.Actor{
			ID:            model.NewID(),
			Name:          req.PC.Name,
			AbilityScores: abilityScoresInputToModel(req.PC.AbilityScores),
			HitPoints:     model.HitPoints{},
			Conditions:    []model.ConditionInstance{},
			Exhaustion:    0,
		},
		Race: model.RaceReference{
			Name: req.PC.Race,
		},
		Alignment:    model.Alignment(req.PC.Alignment),
		Classes:      []model.ClassLevel{},
		Experience:   0,
		Features:     []string{},
		RacialTraits: []string{},
	}

	// 设置背景
	if req.PC.Background != "" {
		pc.Personality = &model.PersonalityTraits{
			Background: req.PC.Background,
		}
	}

	// 添加职业
	if req.PC.Class != "" {
		classID, err := data.GetClassID(req.PC.Class)
		if err != nil {
			return nil, fmt.Errorf("无效的职业: %s", req.PC.Class)
		}

		level := req.PC.Level
		if level < 1 {
			level = 1
		}

		// 获取该职业的特性列表
		features := getClassFeatures(classID, level)

		pc.Classes = append(pc.Classes, model.ClassLevel{
			Class:    classID,
			Level:    level,
			Features: features,
		})
		pc.TotalLevel = level

		// 初始化职业特性系统
		pc.FeatureHooks = make(map[model.ClassID]model.FeatureHook)
		if classID == model.ClassFighter {
			pc.FighterState = &model.FighterFeatures{}
			model.UpdateFighterFeatures(pc.FighterState, level)
			pc.FeatureHooks[classID] = &model.FighterFeatureHooks{
				Features: pc.FighterState,
				Level:    level,
			}
		}
	}

	// 计算派生值
	pc.ArmorClass = calculateArmorClass(pc)
	if req.PC.HitPoints > 0 {
		pc.HitPoints.Maximum = req.PC.HitPoints
		pc.HitPoints.Current = req.PC.HitPoints
	} else {
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

	return &CreatePCResult{
		Actor: actorToInfo(&pc.Actor, model.ActorTypePC, pc.Name),
	}, nil
}

// CreateNPC 创建一个新的非玩家角色
func (e *Engine) CreateNPC(ctx context.Context, req CreateNPCRequest) (*CreateNPCResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpCreateNPC); err != nil {
		return nil, err
	}

	if req.NPC == nil {
		return nil, fmt.Errorf("npc input is required")
	}

	npc := &model.NPC{
		Actor: model.Actor{
			ID:            model.NewID(),
			Name:          req.NPC.Name,
			Description:   req.NPC.Description,
			AbilityScores: abilityScoresInputToModel(req.NPC.AbilityScores),
			Size:          req.NPC.Size,
			Speed:         req.NPC.Speed,
			HitPoints:     model.HitPoints{},
			Conditions:    []model.ConditionInstance{},
			Exhaustion:    0,
		},
	}

	npc.ArmorClass = calculateArmorClassFromActor(&npc.Actor)

	game.NPCs[npc.ID] = npc

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &CreateNPCResult{
		Actor: actorToInfo(&npc.Actor, model.ActorTypeNPC, req.NPC.Name),
	}, nil
}

// CreateEnemy 创建一个新的敌人/怪物
func (e *Engine) CreateEnemy(ctx context.Context, req CreateEnemyRequest) (*CreateEnemyResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpCreateEnemy); err != nil {
		return nil, err
	}

	if req.Enemy == nil {
		return nil, fmt.Errorf("enemy input is required")
	}

	enemy := &model.Enemy{
		Actor: model.Actor{
			ID:            model.NewID(),
			AbilityScores: abilityScoresInputToModel(req.Enemy.AbilityScores),
			Size:          req.Enemy.Size,
			Speed:         req.Enemy.Speed,
			HitPoints:     model.HitPoints{Current: req.Enemy.HitPoints, Maximum: req.Enemy.HitPoints},
			ArmorClass:    req.Enemy.ArmorClass,
			Conditions:    []model.ConditionInstance{},
			Exhaustion:    0,
		},
		ChallengeRating: req.Enemy.ChallengeRating,
	}

	game.Enemies[enemy.ID] = enemy

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &CreateEnemyResult{
		Actor: actorToInfo(&enemy.Actor, model.ActorTypeEnemy, req.Enemy.Name),
	}, nil
}

// CreateCompanion 创建一个同伴角色
func (e *Engine) CreateCompanion(ctx context.Context, req CreateCompanionRequest) (*CreateCompanionResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpCreateCompanion); err != nil {
		return nil, err
	}

	if req.Companion == nil {
		return nil, fmt.Errorf("companion input is required")
	}

	companion := &model.Companion{
		Actor: model.Actor{
			ID:            model.NewID(),
			AbilityScores: abilityScoresInputToModel(req.Companion.AbilityScores),
			Size:          req.Companion.Size,
			Speed:         req.Companion.Speed,
			HitPoints:     model.HitPoints{},
			Conditions:    []model.ConditionInstance{},
			Exhaustion:    0,
		},
		LeaderID: req.Companion.LeaderID,
	}

	companion.ArmorClass = calculateArmorClassFromActor(&companion.Actor)

	game.Companions[companion.ID] = companion

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &CreateCompanionResult{
		Actor: actorToInfo(&companion.Actor, model.ActorTypeCompanion, req.Companion.Name),
	}, nil
}

// GetActor 获取任意类型的角色信息
func (e *Engine) GetActor(ctx context.Context, req GetActorRequest) (*GetActorResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpGetActor); err != nil {
		return nil, err
	}

	actor, ok := game.GetActor(req.ActorID)
	if !ok {
		return nil, ErrNotFound
	}

	var info *ActorInfo
	switch a := actor.(type) {
	case *model.PlayerCharacter:
		info = actorToInfo(&a.Actor, model.ActorTypePC, a.Name)
	case *model.NPC:
		info = actorToInfo(&a.Actor, model.ActorTypeNPC, a.Name)
	case *model.Enemy:
		info = actorToInfo(&a.Actor, model.ActorTypeEnemy, a.Name)
	case *model.Companion:
		info = actorToInfo(&a.Actor, model.ActorTypeCompanion, a.Name)
	default:
		return nil, fmt.Errorf("unknown actor type")
	}

	return &GetActorResult{Actor: info}, nil
}

// GetPC 获取玩家角色的完整数据
func (e *Engine) GetPC(ctx context.Context, req GetPCRequest) (*GetPCResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	pc, ok := game.PCs[req.PCID]
	if !ok {
		return nil, ErrNotFound
	}

	return &GetPCResult{
		PC: playerCharacterToInfo(pc),
	}, nil
}

// UpdateActor 更新角色的部分状态
func (e *Engine) UpdateActor(ctx context.Context, req UpdateActorRequest) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return err
	}

	if err := e.checkPermission(game.Phase, OpUpdateActor); err != nil {
		return err
	}

	actor, ok := game.GetActor(req.ActorID)
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
	if req.Update.AbilityScores != nil {
		baseActor.AbilityScores = abilityScoresInputToModel(*req.Update.AbilityScores)
	}
	if req.Update.HitPoints != nil {
		if req.Update.HitPoints.Current != nil {
			baseActor.HitPoints.Current = *req.Update.HitPoints.Current
		}
		if req.Update.HitPoints.TempHitPoints != nil {
			baseActor.TempHitPoints = *req.Update.HitPoints.TempHitPoints
		}
	}
	if req.Update.Conditions != nil {
		// 添加新状态
		baseActor.Conditions = append(baseActor.Conditions, req.Update.Conditions.Add...)
		// 移除指定状态
		if len(req.Update.Conditions.Remove) > 0 {
			newConditions := make([]model.ConditionInstance, 0)
			for _, c := range baseActor.Conditions {
				shouldRemove := false
				for _, rem := range req.Update.Conditions.Remove {
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
	if req.Update.Position != nil {
		baseActor.Position = req.Update.Position
	}
	if req.Update.SceneID != nil {
		baseActor.SceneID = *req.Update.SceneID
	}

	if err := e.saveGame(ctx, game); err != nil {
		return err
	}

	return nil
}

// RemoveActor 从游戏中移除一个角色
func (e *Engine) RemoveActor(ctx context.Context, req RemoveActorRequest) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return err
	}

	if err := e.checkPermission(game.Phase, OpRemoveActor); err != nil {
		return err
	}

	// 检查是否在战斗中
	if game.Phase == model.PhaseCombat && game.Combat != nil {
		combatant := game.Combat.GetCombatantByActorID(req.ActorID)
		if combatant != nil && !combatant.IsDefeated {
			return ErrInvalidState
		}
	}

	if !game.RemoveActor(req.ActorID) {
		return ErrNotFound
	}

	if err := e.saveGame(ctx, game); err != nil {
		return err
	}

	return nil
}

// ListActors 列出游戏中的角色，可按条件过滤
func (e *Engine) ListActors(ctx context.Context, req ListActorsRequest) (*ListActorsResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	allActors := game.GetAllActors()

	if req.Filter == nil {
		result := make([]ActorInfo, len(allActors))
		for i, actor := range allActors {
			result[i] = *actorSnapshotToInfo(&actor)
		}
		return &ListActorsResult{Actors: result}, nil
	}

	result := make([]ActorInfo, 0)
	for _, actor := range allActors {
		// 按类型过滤
		if len(req.Filter.Types) > 0 {
			typeMatch := false
			for _, t := range req.Filter.Types {
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
		if req.Filter.SceneID != nil && actor.SceneID != *req.Filter.SceneID {
			continue
		}

		// 按存活状态过滤
		if req.Filter.Alive != nil {
			isAlive := actor.HitPoints.Current > 0
			if isAlive != *req.Filter.Alive {
				continue
			}
		}

		result = append(result, *actorSnapshotToInfo(&actor))
	}

	return &ListActorsResult{Actors: result}, nil
}

// AddExperience 为玩家角色添加经验值
func (e *Engine) AddExperience(ctx context.Context, req AddExperienceRequest) (*AddExperienceResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpAddExperience); err != nil {
		return nil, err
	}

	pc, ok := game.PCs[req.PCID]
	if !ok {
		return nil, ErrNotFound
	}

	oldLevel := pc.TotalLevel
	pc.Experience += req.XP

	// 检查是否升级
	newLevel := rules.GetLevelByXP(pc.Experience)
	leveledUp := newLevel > oldLevel
	if leveledUp {
		pc.TotalLevel = newLevel
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &AddExperienceResult{
		LeveledUp: leveledUp,
		OldLevel:  oldLevel,
		NewLevel:  pc.TotalLevel,
	}, nil
}

// LevelUp 手动触发玩家角色升级
func (e *Engine) LevelUp(ctx context.Context, req LevelUpRequest) (*LevelUpResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpLevelUp); err != nil {
		return nil, err
	}

	pc, ok := game.PCs[req.PCID]
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
	classToLevel := req.ClassChoice
	if classToLevel == "" && len(pc.Classes) > 0 {
		classToLevel = string(pc.Classes[0].Class)
	}

	classID, err := data.GetClassID(classToLevel)
	if err != nil {
		classID = model.ClassFighter // 默认战士
	}

	classDef := data.GetClass(classID)
	hitDiceType := 8 // 默认d8
	if classDef != nil {
		hitDiceType = classDef.HitDie
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
func (e *Engine) ShortRest(ctx context.Context, req ShortRestRequest) (*RestResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
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

	for _, actorID := range req.ActorIDs {
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
func (e *Engine) StartLongRest(ctx context.Context, req StartLongRestRequest) (*RestResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpStartLongRest); err != nil {
		return nil, err
	}

	// 创建长休状态
	restState := model.NewLongRest(req.ActorIDs)
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
func (e *Engine) EndLongRest(ctx context.Context, req EndLongRestRequest) (*RestResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
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
	conMod := rules.AbilityModifier(pc.AbilityScores.Constitution)

	for i, cl := range pc.Classes {
		classDef := data.GetClass(cl.Class)
		hitDiceType := 8 // 默认d8
		if classDef != nil {
			hitDiceType = classDef.HitDie
		}

		// 第一级取最大值，之后取平均
		if i == 0 {
			hp += hitDiceType + conMod
		} else {
			hp += (hitDiceType/2 + 1) + conMod
		}
	}
	return hp
}

// abilityScoresInputToModel 将 AbilityScoresInput 转换为 model.AbilityScores
func abilityScoresInputToModel(input AbilityScoresInput) model.AbilityScores {
	return model.AbilityScores{
		Strength:     input.Strength,
		Dexterity:    input.Dexterity,
		Constitution: input.Constitution,
		Intelligence: input.Intelligence,
		Wisdom:       input.Wisdom,
		Charisma:     input.Charisma,
	}
}

// actorToInfo 将 Actor 转换为 ActorInfo
func actorToInfo(actor *model.Actor, actorType model.ActorType, name string) *ActorInfo {
	conditions := make([]string, len(actor.Conditions))
	for i, c := range actor.Conditions {
		conditions[i] = string(c.Type)
	}
	info := &ActorInfo{
		ID:         actor.ID,
		Type:       actorType,
		Name:       name,
		HitPoints:  actor.HitPoints,
		TempHP:     actor.TempHitPoints,
		ArmorClass: actor.ArmorClass,
		Speed:      actor.Speed,
		Conditions: conditions,
		Exhaustion: actor.Exhaustion,
		SceneID:    actor.SceneID,
	}
	if actor.Position != nil {
		pos := *actor.Position
		info.Position = &pos
	}
	return info
}

// actorSnapshotToInfo 将 ActorSnapshot 转换为 ActorInfo
func actorSnapshotToInfo(snapshot *model.ActorSnapshot) *ActorInfo {
	return &ActorInfo{
		ID:         snapshot.ID,
		Type:       snapshot.Type,
		Name:       snapshot.Name,
		HitPoints:  snapshot.HitPoints,
		ArmorClass: snapshot.ArmorClass,
		Conditions: snapshot.Conditions,
		SceneID:    snapshot.SceneID,
	}
}

// playerCharacterToInfo 将 PlayerCharacter 转换为 PlayerCharacterInfo
func playerCharacterToInfo(pc *model.PlayerCharacter) *PlayerCharacterInfo {
	classes := make([]ClassInfo, len(pc.Classes))
	for i, cl := range pc.Classes {
		classes[i] = ClassInfo{
			Class:      cl.Class,
			ClassLevel: cl.Level,
			Features:   cl.Features,
		}
	}
	info := &PlayerCharacterInfo{
		ID:         pc.ID,
		Name:       pc.Actor.Name,
		Race:       pc.Race.Name,
		Classes:    classes,
		TotalLevel: pc.TotalLevel,
		Experience: pc.Experience,
		AbilityScores: AbilityScoresInput{
			Strength:     pc.AbilityScores.Strength,
			Dexterity:    pc.AbilityScores.Dexterity,
			Constitution: pc.AbilityScores.Constitution,
			Intelligence: pc.AbilityScores.Intelligence,
			Wisdom:       pc.AbilityScores.Wisdom,
			Charisma:     pc.AbilityScores.Charisma,
		},
		HitPoints:        pc.HitPoints,
		ArmorClass:       pc.ArmorClass,
		Speed:            pc.Speed,
		ProficiencyBonus: rules.ProficiencyBonus(pc.TotalLevel),
		Features:         pc.Features,
		RacialTraits:     pc.RacialTraits,
	}
	if pc.Personality != nil {
		info.Background = pc.Personality.Background
	}
	return info
}

// getClassFeatures 获取指定职业和等级的特性列表
func getClassFeatures(classID model.ClassID, level int) []string {
	var features []string

	// 战士特性
	if classID == model.ClassFighter {
		for lvl := 1; lvl <= level; lvl++ {
			if feats, ok := data.FighterFeaturesByLevel[lvl]; ok {
				features = append(features, feats...)
			}
		}
	}

	// TODO: 添加其他职业的特性

	return features
}
