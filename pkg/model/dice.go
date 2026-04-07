package model

import "time"

// DiceType 代表骰子类型
type DiceType int

const (
	DiceD4   DiceType = 4
	DiceD6   DiceType = 6
	DiceD8   DiceType = 8
	DiceD10  DiceType = 10
	DiceD12  DiceType = 12
	DiceD20  DiceType = 20
	DiceD100 DiceType = 100
)

// RollModifier 代表掷骰修正
type RollModifier struct {
	Advantage    bool `json:"advantage"`    // 优势
	Disadvantage bool `json:"disadvantage"` // 劣势
	Bonus        int  `json:"bonus"`        // 固定加值
}

// DiceRoll 代表单个骰子的投掷结果
type DiceRoll struct {
	DiceType DiceType `json:"dice_type"` // 骰子类型
	Value    int      `json:"value"`     // 投掷结果
	Kept     bool     `json:"kept"`      // 是否被保留（用于kh/kl等）
	Dropped  bool     `json:"dropped"`   // 是否被丢弃
}

// DiceResult 代表一次完整的骰子投掷结果
type DiceResult struct {
	Expression string     `json:"expression"` // 原始表达式
	Rolls      []DiceRoll `json:"rolls"`      // 所有骰子的结果
	Modifier   int        `json:"modifier"`   // 修正值
	Total      int        `json:"total"`      // 最终结果
	Timestamp  time.Time  `json:"timestamp"`  // 投掷时间
	Hidden     bool       `json:"hidden"`     // 是否是隐藏掷骰
}

// DiceExpression 解析后的骰子表达式
type DiceExpression struct {
	Original       string   `json:"original"`        // 原始表达式
	DiceCount      int      `json:"dice_count"`      // 骰子数量
	DiceType       DiceType `json:"dice_type"`       // 骰子类型
	Modifier       int      `json:"modifier"`        // 修正值
	KeepHigh       int      `json:"keep_high"`       // 保留最高的N个（kh）
	KeepLow        int      `json:"keep_low"`        // 保留最低的N个（kl）
	DropHigh       int      `json:"drop_high"`       // 丢弃最高的N个（dh）
	DropLow        int      `json:"drop_low"`        // 丢弃最低的N个（dl）
	MinValue       int      `json:"min_value"`       // 最小值（如min:8表示最小掷出8）
	IsAdvantage    bool     `json:"is_advantage"`    // 是否是优势掷骰
	IsDisadvantage bool     `json:"is_disadvantage"` // 是否是劣势掷骰
}

// IsValidDiceType 检查是否是有效的骰子类型
func IsValidDiceType(t DiceType) bool {
	switch t {
	case DiceD4, DiceD6, DiceD8, DiceD10, DiceD12, DiceD20, DiceD100:
		return true
	default:
		return false
	}
}
