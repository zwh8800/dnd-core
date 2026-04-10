package rules

import (
	"fmt"
	"math"

	"github.com/zwh8800/dnd-core/pkg/model"
)

// RestResult 休息结果
type RestResult struct {
	Success           bool        `json:"success"`
	Message           string      `json:"message"`
	HPRestored        int         `json:"hp_restored"`
	HitDiceUsed       int         `json:"hit_dice_used"`
	SlotsRestored     map[int]int `json:"slots_restored"` // 环级 -> 恢复数量
	FeaturesReset     []string    `json:"features_reset"` // 重置的特性列表
	DeathSavesCleared bool        `json:"death_saves_cleared"`
	ExhaustionReduced bool        `json:"exhaustion_reduced"`
}

// CalculateShortRest 计算短休效果
// PHB 第8章 - 短休 (Short Rest)
// 短休是指至少1小时的相对休息时间,角色在此期间不能进行除了吃饭、喝水和阅读之外的剧烈活动。
func CalculateShortRest(
	currentHP int,
	maxHP int,
	hitDice []model.HitDiceEntry,
	totalLevel int,
	canUseHitDice bool,
) *RestResult {
	result := &RestResult{
		Success:       true,
		SlotsRestored: make(map[int]int),
		FeaturesReset: make([]string, 0),
		Message:       "完成短休 (1小时)",
	}

	hpRestored := 0
	hitDiceUsed := 0

	// 短休:可以使用生命骰恢复HP
	if canUseHitDice && len(hitDice) > 0 {
		// 注意:实际掷骰由引擎执行,这里只返回可以使用的生命骰信息
		result.Message += " - 可以使用生命骰恢复HP"
		result.FeaturesReset = append(result.FeaturesReset, "可消耗生命骰恢复HP")
	}

	// 短休恢复某些职业特性 (具体由职业特性处理)
	result.FeaturesReset = append(result.FeaturesReset,
		"战士: 气 (Second Wind)",
		"武僧: 气点 (Ki Points)",
		"邪术师: 魔契法术位 (Pact Magic)",
	)

	result.HPRestored = hpRestored
	result.HitDiceUsed = hitDiceUsed

	return result
}

// CalculateLongRest 计算长休效果
// PHB 第8章 - 长休 (Long Rest)
// 长休是指至少8小时的休息时间,角色在此期间不能进行除了睡觉或不超过1小时的其他活动之外的剧烈活动。
func CalculateLongRest(
	currentHP int,
	maxHP int,
	hitDice []model.HitDiceEntry,
	totalLevel int,
	spellcaster *model.SpellcasterState,
	deathSaveSuccesses int,
	deathSaveFailures int,
	exhaustionLevel int,
) *RestResult {
	result := &RestResult{
		Success:       true,
		SlotsRestored: make(map[int]int),
		FeaturesReset: make([]string, 0),
		Message:       "完成长休 (8小时)",
	}

	// 1. 恢复所有HP
	hpRestored := maxHP - currentHP
	if hpRestored < 0 {
		hpRestored = 0
	}
	result.HPRestored = hpRestored
	result.Message += fmt.Sprintf(" - 恢复%d点HP", hpRestored)

	// 2. 恢复生命骰
	// PHB规则:长休结束时,你可以恢复一些已耗尽的生命骰
	// 你能恢复的总数最多等于你总等级的一半(至少一个)
	maxRecoveryDice := int(math.Max(1, float64(totalLevel/2)))
	result.FeaturesReset = append(result.FeaturesReset,
		fmt.Sprintf("可恢复最多%d个生命骰", maxRecoveryDice),
	)

	// 3. 恢复所有法术位 (如果是施法者)
	if spellcaster != nil && spellcaster.Slots != nil {
		slotsRestored := 0
		for level := 1; level <= 9; level++ {
			if spellcaster.Slots.Slots[level][0] > 0 {
				used := spellcaster.Slots.Slots[level][1]
				result.SlotsRestored[level] = used
				slotsRestored += used
			}
		}
		if slotsRestored > 0 {
			result.Message += fmt.Sprintf(" - 恢复%d个法术位", slotsRestored)
		}
	}

	// 4. 清除死亡豁免状态
	if deathSaveSuccesses > 0 || deathSaveFailures > 0 {
		result.DeathSavesCleared = true
		result.Message += " - 清除死亡豁免状态"
	}

	// 5. 减少1级力竭 (如果存在)
	if exhaustionLevel > 0 {
		result.ExhaustionReduced = true
		result.Message += " - 减少1级力竭"
	}

	// 6. 重置所有职业特性使用次数
	result.FeaturesReset = append(result.FeaturesReset,
		"所有每日/长休恢复的特性已重置",
		"野蛮人: 狂暴次数",
		"战士: 动作如潮, 不屈",
		"圣武士: 至圣斩, 神恩",
		"游侠: 宿敌, 自然探索者",
	)

	return result
}

// UseHitDice 使用生命骰恢复HP
// PHB 第8章 - 生命骰 (Hit Dice)
// 你可以在短休期间花费生命骰来恢复HP。花费一个生命骰时,掷骰并加上你的体质调整值,
// 恢复等于该数值的HP。你可以决定每个生命骰是否用于恢复HP,也可以在一个短休中花费多个生命骰。
func UseHitDice(
	currentHP int,
	maxHP int,
	hitDice []model.HitDiceEntry,
	conMod int,
	diceCount int,
) (*RestResult, error) {
	if diceCount <= 0 {
		return nil, fmt.Errorf("必须至少使用1个生命骰")
	}

	result := &RestResult{
		Success:       true,
		SlotsRestored: make(map[int]int),
		FeaturesReset: make([]string, 0),
	}

	totalHealed := 0
	diceUsed := 0

	// 按顺序使用生命骰 (优先使用小的骰子)
	remaining := diceCount
	for i := 0; i < len(hitDice) && remaining > 0; i++ {
		entry := &hitDice[i]
		available := entry.Total - entry.Used

		if available <= 0 {
			continue
		}

		useCount := remaining
		if useCount > available {
			useCount = available
		}

		// 计算恢复的HP (这里返回平均值,实际掷骰由引擎执行)
		avgRoll := (entry.DiceType + 1) / 2 // 平均值
		healPerDie := avgRoll + conMod
		if healPerDie < 1 {
			healPerDie = 1 // 至少恢复1点
		}

		totalHealed += healPerDie * useCount
		entry.Used += useCount
		diceUsed += useCount
		remaining -= useCount
	}

	if diceUsed == 0 {
		return nil, fmt.Errorf("没有可用的生命骰")
	}

	// 应用治疗 (不能超过最大HP)
	actualHealed := totalHealed
	if currentHP+actualHealed > maxHP {
		actualHealed = maxHP - currentHP
	}

	result.HPRestored = actualHealed
	result.HitDiceUsed = diceUsed
	result.Message = fmt.Sprintf("使用%d个生命骰,恢复%d点HP", diceUsed, actualHealed)

	return result, nil
}

// CanTakeRest 检查是否可以进行休息
func CanTakeRest(
	isInCombat bool,
	hasExhaustion int,
) (bool, string) {
	if isInCombat {
		return false, "战斗中无法进行休息"
	}

	if hasExhaustion >= 6 {
		return false, "6级力竭状态无法休息"
	}

	return true, "可以进行休息"
}

// InterruptRest 检查是否打断休息
// PHB 第8章: 如果休息被至少1轮的战斗或类似干扰打断,你必须重新开始休息才能获得休息的益处。
func InterruptRest(restType model.RestType, interruptionDuration string) *RestResult {
	return &RestResult{
		Success: false,
		Message: fmt.Sprintf("%s被%s打断,必须重新开始休息", restType, interruptionDuration),
	}
}
