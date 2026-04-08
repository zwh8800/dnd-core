package model

import (
	"errors"
	"fmt"
	"time"
)

// ErrInvalidActorType 代表无效的角色类型
var ErrInvalidActorType = errors.New("invalid actor type")

// Phase 定义游戏进行的不同阶段
type Phase string

const (
	// PhaseCharacterCreation 角色创建阶段
	PhaseCharacterCreation Phase = "character_creation"

	// PhaseExploration 探索阶段 - 默认阶段
	PhaseExploration Phase = "exploration"

	// PhaseCombat 战斗阶段
	PhaseCombat Phase = "combat"

	// PhaseRest 休息阶段
	PhaseRest Phase = "rest"
)

// GameState 代表一个完整的游戏状态
type GameState struct {
	ID          ID        `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Phase       Phase     `json:"phase"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 角色
	PCs        map[ID]*PlayerCharacter `json:"pcs"`
	Companions map[ID]*Companion       `json:"companions"`
	NPCs       map[ID]*NPC             `json:"npcs"`
	Enemies    map[ID]*Enemy           `json:"enemies"`

	// 场景
	Scenes       map[ID]*Scene `json:"scenes"`
	CurrentScene *ID           `json:"current_scene,omitempty"`

	// 物品（场景中的物品，不在任何角色库存中）
	Items map[ID]*Item `json:"items"`

	// 库存
	Inventories map[ID]*Inventory `json:"inventories"`

	// 战斗
	Combat *CombatState `json:"combat,omitempty"`

	// 任务
	Quests map[ID]*Quest `json:"quests"`

	// 休息
	ActiveRest *RestState `json:"active_rest,omitempty"`

	// 旅行
	CurrentTravel *TravelState `json:"current_travel,omitempty"`

	// 生活方式
	Lifestyle *LifestyleState `json:"lifestyle,omitempty"`

	// 游戏时间
	GameTime GameTime `json:"game_time"`
}

// GameTime 代表游戏内时间
type GameTime struct {
	Year   int `json:"year"`
	Month  int `json:"month"`
	Day    int `json:"day"`
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
}

// NewGameState 创建新的游戏状态
func NewGameState(name, description string) *GameState {
	return &GameState{
		ID:          NewID(),
		Name:        name,
		Description: description,
		Phase:       PhaseCharacterCreation,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		PCs:         make(map[ID]*PlayerCharacter),
		Companions:  make(map[ID]*Companion),
		NPCs:        make(map[ID]*NPC),
		Enemies:     make(map[ID]*Enemy),
		Scenes:      make(map[ID]*Scene),
		Items:       make(map[ID]*Item),
		Inventories: make(map[ID]*Inventory),
		Quests:      make(map[ID]*Quest),
		GameTime: GameTime{
			Year:   1,
			Month:  1,
			Day:    1,
			Hour:   8,
			Minute: 0,
		},
	}
}

// AddActor 添加角色到游戏
func (gs *GameState) AddActor(actor any) error {
	switch a := actor.(type) {
	case *PlayerCharacter:
		gs.PCs[a.ID] = a
	case *Companion:
		gs.Companions[a.ID] = a
	case *NPC:
		gs.NPCs[a.ID] = a
	case *Enemy:
		gs.Enemies[a.ID] = a
	default:
		return ErrInvalidActorType
	}
	gs.UpdatedAt = time.Now()
	return nil
}

// RemoveActor 从游戏移除角色
func (gs *GameState) RemoveActor(actorID ID) bool {
	if _, ok := gs.PCs[actorID]; ok {
		delete(gs.PCs, actorID)
		gs.UpdatedAt = time.Now()
		return true
	}
	if _, ok := gs.Companions[actorID]; ok {
		delete(gs.Companions, actorID)
		gs.UpdatedAt = time.Now()
		return true
	}
	if _, ok := gs.NPCs[actorID]; ok {
		delete(gs.NPCs, actorID)
		gs.UpdatedAt = time.Now()
		return true
	}
	if _, ok := gs.Enemies[actorID]; ok {
		delete(gs.Enemies, actorID)
		gs.UpdatedAt = time.Now()
		return true
	}
	return false
}

// GetActor 获取角色
func (gs *GameState) GetActor(actorID ID) (any, bool) {
	if pc, ok := gs.PCs[actorID]; ok {
		return pc, true
	}
	if c, ok := gs.Companions[actorID]; ok {
		return c, true
	}
	if n, ok := gs.NPCs[actorID]; ok {
		return n, true
	}
	if e, ok := gs.Enemies[actorID]; ok {
		return e, true
	}
	return nil, false
}

// GetAllActors 获取所有角色
func (gs *GameState) GetAllActors() []ActorSnapshot {
	snapshots := make([]ActorSnapshot, 0)

	for _, pc := range gs.PCs {
		snapshots = append(snapshots, ActorToSnapshot(&pc.Actor, ActorTypePC, pc.Name))
	}
	for _, c := range gs.Companions {
		snapshots = append(snapshots, ActorToSnapshot(&c.Actor, ActorTypeCompanion, c.Name))
	}
	for _, n := range gs.NPCs {
		snapshots = append(snapshots, ActorToSnapshot(&n.Actor, ActorTypeNPC, n.Name))
	}
	for _, e := range gs.Enemies {
		snapshots = append(snapshots, ActorToSnapshot(&e.Actor, ActorTypeEnemy, e.Name))
	}

	return snapshots
}

// ActorSnapshot 角色快照（简化版）
type ActorSnapshot struct {
	ID         ID        `json:"id"`
	Type       ActorType `json:"type"`
	Name       string    `json:"name"`
	HitPoints  HitPoints `json:"hit_points"`
	ArmorClass int       `json:"armor_class"`
	Conditions []string  `json:"conditions"`
	SceneID    ID        `json:"scene_id"`
	Summary    string    `json:"summary"`
}

// ActorToSnapshot 将Actor转换为快照
func ActorToSnapshot(actor *Actor, actorType ActorType, name string) ActorSnapshot {
	conditions := make([]string, len(actor.Conditions))
	for i, c := range actor.Conditions {
		conditions[i] = string(c.Type)
	}

	return ActorSnapshot{
		ID:         actor.ID,
		Type:       actorType,
		Name:       name,
		HitPoints:  actor.HitPoints,
		ArmorClass: actor.ArmorClass,
		Conditions: conditions,
		SceneID:    actor.SceneID,
		Summary:    generateActorSummary(actor, actorType),
	}
}

// generateActorSummary 生成角色摘要
func generateActorSummary(actor *Actor, actorType ActorType) string {
	status := "Active"
	if !actor.IsAlive() {
		status = "Dead"
	} else if actor.IsIncapacitated() {
		status = "Incapacitated"
	}

	return fmt.Sprintf("%s - HP: %d/%d - AC: %d - %s",
		actorType,
		actor.HitPoints.Current,
		actor.HitPoints.Maximum,
		actor.ArmorClass,
		status)
}
