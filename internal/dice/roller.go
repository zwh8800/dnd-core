package dice

import (
	"fmt"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/zwh8800/dnd-core/internal/model"
)

// Roller 骰子投掷器
type Roller struct {
	rng *rand.Rand
}

// New 创建新的骰子投掷器
func New(seed int64) *Roller {
	var src rand.Source
	if seed == 0 {
		src = rand.NewSource(rand.Int63())
	} else {
		src = rand.NewSource(seed)
	}
	return &Roller{
		rng: rand.New(src),
	}
}

// 骰子表达式正则
var diceRegex = regexp.MustCompile(`^(\d+)d(\d+)(kh\d+|kl\d+|dh\d+|dl\d+|min:\d+)?([+-]\d+)?$`)

// ParseExpression 解析骰子表达式
func ParseExpression(expr string) (*model.DiceExpression, error) {
	expr = strings.TrimSpace(strings.ToLower(expr))

	result := &model.DiceExpression{
		Original: expr,
	}

	// 处理优势/劣势
	if expr == "adv" || expr == "advantage" {
		result.IsAdvantage = true
		result.DiceCount = 2
		result.DiceType = model.DiceD20
		return result, nil
	}
	if expr == "dis" || expr == "disadvantage" {
		result.IsDisadvantage = true
		result.DiceCount = 2
		result.DiceType = model.DiceD20
		return result, nil
	}

	// 简单 d20 情况
	if expr == "d20" {
		result.DiceCount = 1
		result.DiceType = model.DiceD20
		return result, nil
	}

	matches := diceRegex.FindStringSubmatch(expr)
	if matches == nil {
		return nil, fmt.Errorf("invalid dice expression: %s", expr)
	}

	count, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, fmt.Errorf("invalid dice count: %s", matches[1])
	}
	result.DiceCount = count

	diceType, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid dice type: %s", matches[2])
	}
	result.DiceType = model.DiceType(diceType)

	if !model.IsValidDiceType(result.DiceType) {
		return nil, fmt.Errorf("invalid dice type: d%d", diceType)
	}

	// 处理修饰符
	if matches[3] != "" {
		mod := matches[3]
		if strings.HasPrefix(mod, "kh") {
			n, _ := strconv.Atoi(mod[2:])
			result.KeepHigh = n
		} else if strings.HasPrefix(mod, "kl") {
			n, _ := strconv.Atoi(mod[2:])
			result.KeepLow = n
		} else if strings.HasPrefix(mod, "dh") {
			n, _ := strconv.Atoi(mod[2:])
			result.DropHigh = n
		} else if strings.HasPrefix(mod, "dl") {
			n, _ := strconv.Atoi(mod[2:])
			result.DropLow = n
		} else if strings.HasPrefix(mod, "min:") {
			n, _ := strconv.Atoi(mod[4:])
			result.MinValue = n
		}
	}

	// 处理固定修正值
	if matches[4] != "" {
		mod, _ := strconv.Atoi(matches[4])
		result.Modifier = mod
	}

	return result, nil
}

// Roll 投掷骰子表达式
func (r *Roller) Roll(expression string) (*model.DiceResult, error) {
	expr, err := ParseExpression(expression)
	if err != nil {
		return nil, err
	}

	return r.RollExpression(expr)
}

// RollExpression 使用已解析的表达式投掷
func (r *Roller) RollExpression(expr *model.DiceExpression) (*model.DiceResult, error) {
	// 处理优势/劣势
	if expr.IsAdvantage {
		return r.rollAdvantage(expr.Modifier)
	}
	if expr.IsDisadvantage {
		return r.rollDisadvantage(expr.Modifier)
	}

	result := &model.DiceResult{
		Expression: expr.Original,
		Modifier:   expr.Modifier,
		Rolls:      make([]model.DiceRoll, 0),
	}

	// 投掷骰子
	for i := 0; i < expr.DiceCount; i++ {
		value := r.rollDie(int(expr.DiceType))

		// 应用最小值
		if expr.MinValue > 0 && value < expr.MinValue {
			value = expr.MinValue
		}

		result.Rolls = append(result.Rolls, model.DiceRoll{
			DiceType: expr.DiceType,
			Value:    value,
			Kept:     true,
		})
	}

	// 应用保留/丢弃逻辑
	result.Rolls = applyKeepDrop(result.Rolls, expr)

	// 计算总计
	total := expr.Modifier
	for _, roll := range result.Rolls {
		if roll.Kept {
			total += roll.Value
		}
	}
	result.Total = total

	return result, nil
}

// rollDie 投掷单个骰子
func (r *Roller) rollDie(sides int) int {
	return r.rng.Intn(sides) + 1
}

// rollAdvantage 投掷优势（2d20取高）
func (r *Roller) rollAdvantage(modifier int) (*model.DiceResult, error) {
	roll1 := r.rollDie(20)
	roll2 := r.rollDie(20)

	high := roll1
	if roll2 > roll1 {
		high = roll2
	}

	result := &model.DiceResult{
		Expression: "1d20 (advantage)",
		Modifier:   modifier,
		Total:      high + modifier,
		Rolls: []model.DiceRoll{
			{DiceType: model.DiceD20, Value: roll1, Kept: roll1 >= roll2},
			{DiceType: model.DiceD20, Value: roll2, Kept: roll2 >= roll1},
		},
	}

	return result, nil
}

// rollDisadvantage 投掷劣势（2d20取低）
func (r *Roller) rollDisadvantage(modifier int) (*model.DiceResult, error) {
	roll1 := r.rollDie(20)
	roll2 := r.rollDie(20)

	low := roll1
	if roll2 < roll1 {
		low = roll2
	}

	result := &model.DiceResult{
		Expression: "1d20 (disadvantage)",
		Modifier:   modifier,
		Total:      low + modifier,
		Rolls: []model.DiceRoll{
			{DiceType: model.DiceD20, Value: roll1, Kept: roll1 <= roll2},
			{DiceType: model.DiceD20, Value: roll2, Kept: roll2 <= roll1},
		},
	}

	return result, nil
}

// RollAdvantage 公开的优势投掷方法
func (r *Roller) RollAdvantage(modifier int) (*model.DiceResult, error) {
	return r.rollAdvantage(modifier)
}

// RollDisadvantage 公开的劣势投掷方法
func (r *Roller) RollDisadvantage(modifier int) (*model.DiceResult, error) {
	return r.rollDisadvantage(modifier)
}

// applyKeepDrop 应用保留/丢弃逻辑
func applyKeepDrop(rolls []model.DiceRoll, expr *model.DiceExpression) []model.DiceRoll {
	if len(rolls) == 0 {
		return rolls
	}

	// 按值排序
	sort.Slice(rolls, func(i, j int) bool {
		return rolls[i].Value < rolls[j].Value
	})

	// 初始化所有骰子为保留
	for i := range rolls {
		rolls[i].Kept = true
		rolls[i].Dropped = false
	}

	// 丢弃最高
	for i := 0; i < expr.DropHigh && i < len(rolls); i++ {
		rolls[len(rolls)-1-i].Kept = false
		rolls[len(rolls)-1-i].Dropped = true
	}

	// 丢弃最低
	for i := 0; i < expr.DropLow && i < len(rolls); i++ {
		rolls[i].Kept = false
		rolls[i].Dropped = true
	}

	// 保留最高
	if expr.KeepHigh > 0 && expr.KeepHigh < len(rolls) {
		keepCount := expr.KeepHigh
		for i := 0; i < len(rolls)-keepCount; i++ {
			rolls[i].Kept = false
			rolls[i].Dropped = true
		}
	}

	// 保留最低
	if expr.KeepLow > 0 && expr.KeepLow < len(rolls) {
		keepCount := expr.KeepLow
		for i := keepCount; i < len(rolls); i++ {
			rolls[i].Kept = false
			rolls[i].Dropped = true
		}
	}

	return rolls
}
