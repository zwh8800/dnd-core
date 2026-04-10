package rules

import "strings"

// CalculateCarryingCapacity 计算负重能力
// PHB Ch.7: 负重能力 = 力量值 × 15（磅）
func CalculateCarryingCapacity(strength int) int {
	return strength * CarryingCapacityMultiplier
}

// CalculatePushDragLift 计算推/拖/举能力
// PHB Ch.7: 推/拖/举能力 = 力量值 × 30（磅）
func CalculatePushDragLift(strength int) int {
	return strength * PushDragLiftMultiplier
}

// JoinStrings 连接字符串切片
// 辅助函数，用于将字符串切片连接为单个字符串
func JoinStrings(strs []string, sep string) string {
	return strings.Join(strs, sep)
}
