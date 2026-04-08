package model

// CraftingType 制作类型
type CraftingType string

const (
	CraftingTypeNonMagical CraftingType = "non_magical" // 非魔法物品
	CraftingTypePotion     CraftingType = "potion"      // 药水酿造
	CraftingTypeScroll     CraftingType = "scroll"      // 法术卷轴
	CraftingTypeMagicItem  CraftingType = "magic_item"  // 魔法物品
)

// CraftingMaterial 制作材料
type CraftingMaterial struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Cost     int    `json:"cost"` // 成本（铜币）
}

// CraftingRecipe 制作配方
type CraftingRecipe struct {
	ID            string             `json:"id"`
	Name          string             `json:"name"`
	Type          CraftingType       `json:"type"`
	Description   string             `json:"description"`
	Materials     []CraftingMaterial `json:"materials"`      // 所需材料
	ToolsRequired []string           `json:"tools_required"` // 所需工具熟练
	SkillRequired string             `json:"skill_required"` // 所需技能（如Arcana）
	TimeDays      int                `json:"time_days"`      // 制作时间（天）
	DC            int                `json:"dc"`             // 制作DC
	MinLevel      int                `json:"min_level"`      // 最低等级
	SpellRequired string             `json:"spell_required"` // 需要能施放的法术
	Cost          int                `json:"cost"`           // 总成本（铜币）
}

// CraftingProgress 制作进度
type CraftingProgress struct {
	RecipeID    string `json:"recipe_id"`
	DaysWorked  int    `json:"days_worked"`   // 已工作天数
	TotalDays   int    `json:"total_days"`    // 总需要天数
	MoneySpent  int    `json:"money_spent"`   // 已花费金额
	LastWorkDay string `json:"last_work_day"` // 最后工作日
	IsComplete  bool   `json:"is_complete"`   // 是否完成
}
