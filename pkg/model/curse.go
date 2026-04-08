package model

// CurseDefinition 诅咒定义
type CurseDefinition struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Effect       string `json:"effect"`        // 效果描述
	RemoveMethod string `json:"remove_method"` // 移除方法（如 Remove Curse 法术）
	DC           int    `json:"dc,omitempty"`  // 豁免DC（如果有）
}

// CurseInstance 诅咒实例
type CurseInstance struct {
	CurseID           string `json:"curse_id"`
	Source            string `json:"source"`             // 诅咒来源
	AppliedAt         string `json:"applied_at"`         // 施加时间
	RemainingDuration string `json:"remaining_duration"` // 剩余持续时间
	IsPermanent       bool   `json:"is_permanent"`       // 是否永久
}
