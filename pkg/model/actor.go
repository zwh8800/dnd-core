package model

// ActorType 定义角色的类型
type ActorType string

const (
	ActorTypePC        ActorType = "pc"        // 玩家角色
	ActorTypeNPC       ActorType = "npc"       // 非玩家角色
	ActorTypeEnemy     ActorType = "enemy"     // 敌人/怪物
	ActorTypeCompanion ActorType = "companion" // 同伴/盟友
)

// Size 代表生物的体型
type Size string

const (
	SizeTiny       Size = "Tiny"       // 超小型（2.5×2.5英尺）
	SizeSmall      Size = "Small"      // 小型（5×5英尺）
	SizeMedium     Size = "Medium"     // 中型（5×5英尺）
	SizeLarge      Size = "Large"      // 大型（10×10英尺）
	SizeHuge       Size = "Huge"       // 超大型（15×15英尺）
	SizeGargantuan Size = "Gargantuan" // 巨型（20×20英尺或更大）
)

// HitPoints 代表生物的生命值状态
type HitPoints struct {
	Current int `json:"current"`
	Maximum int `json:"maximum"`
}

// Actor 是所有生物类型的基类
type Actor struct {
	ID              ID                  `json:"id"`
	Type            ActorType           `json:"type"`
	Name            string              `json:"name"`
	Description     string              `json:"description"`
	Size            Size                `json:"size"`
	Speed           int                 `json:"speed"` // 基础移动速度（英尺）
	AbilityScores   AbilityScores       `json:"ability_scores"`
	Proficiencies   Proficiencies       `json:"proficiencies"`
	HitPoints       HitPoints           `json:"hit_points"`
	TempHitPoints   int                 `json:"temp_hit_points"`
	ArmorClass      int                 `json:"armor_class"`
	Conditions      []ConditionInstance `json:"conditions"`
	Exhaustion      int                 `json:"exhaustion"`         // 力竭等级（0-6）
	SceneID         ID                  `json:"scene_id"`           // 当前所在场景
	Position        *Point              `json:"position,omitempty"` // 战斗中的位置（可选）
	InitiativeBonus int                 `json:"initiative_bonus"`   // 先攻修正
}

// Point 代表二维坐标位置
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// IsAlive 检查生物是否存活
func (a *Actor) IsAlive() bool {
	return a.HitPoints.Current > 0 || a.TempHitPoints > 0
}

// IsDead 检查生物是否死亡
func (a *Actor) IsDead() bool {
	return a.HitPoints.Current <= 0 && !a.IsStabilized()
}

// IsStabilized 检查生物是否已稳定（不再进行死亡豁免）
func (a *Actor) IsStabilized() bool {
	// 检查是否有稳定状态
	for _, c := range a.Conditions {
		if c.Type == ConditionStabilized {
			return true
		}
	}
	return false
}

// IsIncapacitated 检查生物是否失去行动能力
func (a *Actor) IsIncapacitated() bool {
	for _, c := range a.Conditions {
		if c.Type == ConditionIncapacitated {
			return true
		}
	}
	//  unconscious 也会导致 incapacitated
	for _, c := range a.Conditions {
		if c.Type == ConditionUnconscious {
			return true
		}
	}
	return false
}

// HasCondition 检查是否具有特定状态
func (a *Actor) HasCondition(conditionType ConditionType) bool {
	for _, c := range a.Conditions {
		if c.Type == conditionType {
			return true
		}
	}
	return false
}
