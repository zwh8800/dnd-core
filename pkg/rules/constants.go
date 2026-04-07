package rules

// D&D 5e 核心常量

// 标准难度等级
const (
	DCEasy             = 10
	DCMedium           = 15
	DCHard             = 20
	DCVeryHard         = 25
	DCNearlyImpossible = 30
)

// 升级经验值阈值（D&D 5e标准）
var XPThresholds = map[int]int{
	1:  0,
	2:  300,
	3:  900,
	4:  2700,
	5:  6500,
	6:  14000,
	7:  23000,
	8:  34000,
	9:  48000,
	10: 64000,
	11: 85000,
	12: 100000,
	13: 120000,
	14: 140000,
	15: 165000,
	16: 195000,
	17: 225000,
	18: 265000,
	19: 305000,
	20: 355000,
}

// GetLevelByXP 根据经验值获取等级
func GetLevelByXP(xp int) int {
	for level := 20; level >= 1; level-- {
		if xp >= XPThresholds[level] {
			return level
		}
	}
	return 1
}

// GetXPForLevel 获取指定等级所需的经验值
func GetXPForLevel(level int) int {
	if xp, ok := XPThresholds[level]; ok {
		return xp
	}
	return 0
}

// 熟练加值表
var ProficiencyBonusTable = map[int]int{
	1:  2,
	2:  2,
	3:  2,
	4:  2,
	5:  3,
	6:  3,
	7:  3,
	8:  3,
	9:  4,
	10: 4,
	11: 4,
	12: 4,
	13: 5,
	14: 5,
	15: 5,
	16: 5,
	17: 6,
	18: 6,
	19: 6,
	20: 6,
}

// 休息持续时间（小时）
const (
	ShortRestDuration = 1 // 短休1小时
	LongRestDuration  = 8 // 长休8小时
)

// 力竭等级效果
var ExhaustionEffects = map[int]string{
	1: "属性检定劣势",
	2: "速度减半",
	3: "攻击检定和豁免检定劣势",
	4: "HP最大值减半",
	5: "速度降为0",
	6: "死亡",
}

// 体型对应的空间占用（英尺）
var SizeSpaceMap = map[string]int{
	"Tiny":       2,
	"Small":      5,
	"Medium":     5,
	"Large":      10,
	"Huge":       15,
	"Gargantuan": 20,
}

// 标准移动速度（英尺）
const (
	StandardSpeed  = 30 // 中等体型通常速度
	FlySpeedBase   = 30 // 基础飞行速度
	SwimSpeedBase  = 30 // 基础游泳速度
	ClimbSpeedBase = 30 // 基础攀爬速度
)

// 负重计算倍数
const (
	CarryingCapacityMultiplier = 15 // 负重能力 = 力量值 × 15
	PushDragLiftMultiplier     = 30 // 推/拖/举 = 力量值 × 30
)

// 休息恢复规则
const (
	// 长休恢复的生命骰比例（总等级的一半，至少1个）
	LongRestHitDiceRecoveryRatio = 2
	// 长休减少的力竭等级
	LongRestExhaustionReduction = 1
)
