package model

import "time"

// CombatStatus 代表战斗的状态
type CombatStatus string

const (
	CombatStatusActive   CombatStatus = "active"   // 战斗进行中
	CombatStatusFinished CombatStatus = "finished" // 战斗已结束
)

// CombatantEntry 代表先攻列表中的一个条目
type CombatantEntry struct {
	ActorID         ID     `json:"actor_id"`
	ActorName       string `json:"actor_name"`
	InitiativeRoll  int    `json:"initiative_roll"`
	InitiativeTotal int    `json:"initiative_total"`
	IsSurprised     bool   `json:"is_surprised"`
	IsDefeated      bool   `json:"is_defeated"`
}

// TurnState 代表当前回合的状态
type TurnState struct {
	ActorID         ID   `json:"actor_id"`
	Round           int  `json:"round"`
	MovementUsed    int  `json:"movement_used"`
	ActionUsed      bool `json:"action_used"`
	BonusActionUsed bool `json:"bonus_action_used"`
	ReactionUsed    bool `json:"reaction_used"`
	ReactionTurn    int  `json:"reaction_turn"` // 可以在哪个回合使用反应
}

// CombatLogEntry 代表战斗日志中的一条记录
type CombatLogEntry struct {
	Timestamp   time.Time `json:"timestamp"`
	Round       int       `json:"round"`
	ActorID     ID        `json:"actor_id,omitempty"`
	Action      string    `json:"action"`
	Description string    `json:"description"`
}

// CombatState 代表一场战斗的完整状态
type CombatState struct {
	ID           ID               `json:"id"`
	SceneID      ID               `json:"scene_id"`
	Status       CombatStatus     `json:"status"`
	Round        int              `json:"round"`
	Initiative   []CombatantEntry `json:"initiative"`
	CurrentIndex int              `json:"current_index"`
	CurrentTurn  *TurnState       `json:"current_turn"`
	Log          []CombatLogEntry `json:"log"`
}

// GetCurrentCombatant 获取当前回合的战斗者
func (cs *CombatState) GetCurrentCombatant() *CombatantEntry {
	if cs.CurrentIndex < 0 || cs.CurrentIndex >= len(cs.Initiative) {
		return nil
	}
	return &cs.Initiative[cs.CurrentIndex]
}

// GetCurrentActorID 获取当前回合角色的ID
func (cs *CombatState) GetCurrentActorID() ID {
	combatant := cs.GetCurrentCombatant()
	if combatant == nil {
		return ""
	}
	return combatant.ActorID
}

// IsActorTurn 检查是否是指定角色的回合
func (cs *CombatState) IsActorTurn(actorID ID) bool {
	return cs.GetCurrentActorID() == actorID
}

// GetCombatantByActorID 根据角色ID查找战斗者条目
func (cs *CombatState) GetCombatantByActorID(actorID ID) *CombatantEntry {
	for i := range cs.Initiative {
		if cs.Initiative[i].ActorID == actorID {
			return &cs.Initiative[i]
		}
	}
	return nil
}

// AdvanceTurn 推进到下一个回合
func (cs *CombatState) AdvanceTurn() {
	if len(cs.Initiative) == 0 {
		return
	}

	cs.CurrentIndex++
	if cs.CurrentIndex >= len(cs.Initiative) {
		cs.CurrentIndex = 0
		cs.Round++
	}

	// 跳过已 defeat 的战斗者
	for cs.CurrentIndex < len(cs.Initiative) && cs.Initiative[cs.CurrentIndex].IsDefeated {
		cs.CurrentIndex++
		if cs.CurrentIndex >= len(cs.Initiative) {
			cs.CurrentIndex = 0
			cs.Round++
		}
	}

	// 创建新的回合状态
	cs.CurrentTurn = &TurnState{
		ActorID:         cs.GetCurrentActorID(),
		Round:           cs.Round,
		MovementUsed:    0,
		ActionUsed:      false,
		BonusActionUsed: false,
		ReactionUsed:    false,
		ReactionTurn:    cs.Round,
	}
}
