package rules

// CheckResult 代表一次检定的结果
type CheckResult struct {
	DieRoll     int  // 骰子投掷结果
	Modifier    int  // 修正值
	Total       int  // 总计
	DC          int  // 难度等级
	Success     bool // 是否成功
	IsNatural20 bool // 是否天然20
	IsNatural1  bool // 是否天然1
}

// PerformCheck 执行属性检定
func PerformCheck(roll int, modifier int, dc int) *CheckResult {
	isNat20 := roll == 20
	isNat1 := roll == 1
	total := roll + modifier

	result := &CheckResult{
		DieRoll:     roll,
		Modifier:    modifier,
		Total:       total,
		DC:          dc,
		Success:     total >= dc,
		IsNatural20: isNat20,
		IsNatural1:  isNat1,
	}

	// 天然20在死亡豁免中有特殊规则
	// 但在普通检定中，天然20/1没有自动成功/失败规则
	return result
}

// SavingThrowResult 豁免检定结果
type SavingThrowResult struct {
	CheckResult
	CriticalSuccess bool // 某些法术/效果在天然20时完全无效
}

// PerformSavingThrow 执行豁免检定
func PerformSavingThrow(roll int, modifier int, dc int) *SavingThrowResult {
	check := PerformCheck(roll, modifier, dc)
	return &SavingThrowResult{
		CheckResult:     *check,
		CriticalSuccess: check.IsNatural20, // 天然20通常是完全成功
	}
}

// DeathSaveResult 死亡豁免结果
type DeathSaveResult struct {
	Roll           int
	Success        bool
	IsCritical     bool // 天然20 = 立即恢复1HP
	IsCriticalFail bool // 天然1 = 2次失败
}

// PerformDeathSave 执行死亡豁免检定
func PerformDeathSave(roll int) *DeathSaveResult {
	result := &DeathSaveResult{
		Roll: roll,
	}

	if roll == 20 {
		result.Success = true
		result.IsCritical = true // 立即恢复1HP
	} else if roll == 1 {
		result.Success = false
		result.IsCriticalFail = true // 算2次失败
	} else if roll >= 10 {
		result.Success = true
	} else {
		result.Success = false
	}

	return result
}

// ContestedCheckResult 对抗检定结果
type ContestedCheckResult struct {
	ActorATotal int
	ActorBTotal int
	WinnerA     bool
	WinnerB     bool
	Tie         bool
}

// PerformContestedCheck 执行对抗检定
func PerformContestedCheck(rollA, modA, rollB, modB int) *ContestedCheckResult {
	totalA := rollA + modA
	totalB := rollB + modB

	result := &ContestedCheckResult{
		ActorATotal: totalA,
		ActorBTotal: totalB,
	}

	if totalA > totalB {
		result.WinnerA = true
	} else if totalB > totalA {
		result.WinnerB = true
	} else {
		result.Tie = true // 平局时保持现状
	}

	return result
}
