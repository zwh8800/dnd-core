package model

// LifestyleTier 代表生活方式等级
type LifestyleTier string

const (
	LifestyleWretched     LifestyleTier = "wretched"     // 悲惨
	LifestyleSqualid      LifestyleTier = "squalid"      // 肮脏
	LifestylePoor         LifestyleTier = "poor"         // 贫穷
	LifestyleModest       LifestyleTier = "modest"       // 朴素
	LifestyleComfortable  LifestyleTier = "comfortable"  // 舒适
	LifestyleWealthy      LifestyleTier = "wealthy"      // 富裕
	LifestyleAristocratic LifestyleTier = "aristocratic" // 贵族
)

// LifestyleState 代表生活方式状态
type LifestyleState struct {
	CurrentTier LifestyleTier `json:"current_tier"` // 当前生活方式
	DailyCost   int           `json:"daily_cost"`   // 每日花费（铜币）
	DaysElapsed int           `json:"days_elapsed"` // 已过去的天数
	TotalSpent  int           `json:"total_spent"`  // 总花费（铜币）
	LastPayment int           `json:"last_payment"` // 上次支付
}

// LifestyleCost 生活方式费用结果
type LifestyleCost struct {
	DailyCost   int `json:"daily_cost"`
	MonthlyCost int `json:"monthly_cost"`
}
