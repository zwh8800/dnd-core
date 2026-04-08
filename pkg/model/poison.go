package model

// PoisonDeliveryType 毒药传递方式
type PoisonDeliveryType string

const (
	PoisonContact  PoisonDeliveryType = "contact"  // 接触
	PoisonIngested PoisonDeliveryType = "ingested" // 摄入
	PoisonInhaled  PoisonDeliveryType = "inhaled"  // 吸入
	PoisonInjury   PoisonDeliveryType = "injury"   // 伤口
)

// PoisonEffect 毒药效果
type PoisonEffect struct {
	DamageDice   string `json:"damage_dice,omitempty"`   // 伤害骰
	DamageType   string `json:"damage_type,omitempty"`   // 伤害类型
	SaveDC       int    `json:"save_dc"`                 // 豁免DC
	Duration     string `json:"duration"`                // 持续时间
	StatusEffect string `json:"status_effect,omitempty"` // 状态效果
	Description  string `json:"description"`             // 效果描述
}

// PoisonDefinition 毒药定义
type PoisonDefinition struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Type        PoisonDeliveryType `json:"type"`
	Description string             `json:"description"`
	Effect      PoisonEffect       `json:"effect"`
	Price       int                `json:"price"`  // 价格（铜币）
	Rarity      string             `json:"rarity"` // 稀有度
}

// PoisonInstance 毒药实例
type PoisonInstance struct {
	PoisonID        string `json:"poison_id"`
	RemainingUses   int    `json:"remaining_uses"`             // 剩余使用次数
	AppliedTo       string `json:"applied_to,omitempty"`       // 应用到的武器/弹药ID
	ApplicationTime string `json:"application_time,omitempty"` // 涂抹时间
	ExpiresAfter    string `json:"expires_after,omitempty"`    // 过期时间（1分钟）
}
