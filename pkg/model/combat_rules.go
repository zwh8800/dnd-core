package model

// CoverType 代表掩护类型
type CoverType string

const (
	CoverNone          CoverType = "none"           // 无掩护
	CoverHalf          CoverType = "half"           // 半身掩护 (+2 AC和DEX豁免)
	CoverThreeQuarters CoverType = "three_quarters" // 四分之三掩护 (+5 AC和DEX豁免)
	CoverTotal         CoverType = "total"          // 全身掩护 (无法被直接targeting)
)

// CoverBonus 掩护的AC和豁免加值
type CoverBonus struct {
	ACBonus      int `json:"ac_bonus"`       // AC加值
	DexSaveBonus int `json:"dex_save_bonus"` // 敏捷豁免加值
}

// GetCoverBonus 获取掩护加值
func (c CoverType) GetCoverBonus() CoverBonus {
	switch c {
	case CoverHalf:
		return CoverBonus{ACBonus: 2, DexSaveBonus: 2}
	case CoverThreeQuarters:
		return CoverBonus{ACBonus: 5, DexSaveBonus: 5}
	case CoverTotal:
		return CoverBonus{ACBonus: 999, DexSaveBonus: 999} // 实际上无法被targeting
	default:
		return CoverBonus{ACBonus: 0, DexSaveBonus: 0}
	}
}

// GrappleState 擒抱状态
type GrappleState struct {
	GrapplerID ID  `json:"grappler_id"` // 擒抱者ID
	GrappeeID  ID  `json:"grappee_id"`  // 被擒抱者ID
	EscapeDC   int `json:"escape_dc"`   // 逃脱DC
}

// JumpType 跳跃类型
type JumpType string

const (
	JumpTypeLong JumpType = "long" // 跳远
	JumpTypeHigh JumpType = "high" // 跳高
)

// SuffocationState 窒息状态
type SuffocationState struct {
	CanHoldBreath       bool `json:"can_hold_breath"`       // 是否能继续闭气
	RoundsUntilDrowning int  `json:"rounds_until_drowning"` // 窒息轮数
}
