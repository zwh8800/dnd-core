package model

// DamageType 代表伤害类型（已在equipment.go中定义基础类型，这里添加常量和功能）
// 注意：这里重新定义是为了保持包的完整性，实际使用时以首次定义为准
// 此处添加伤害类型的辅助方法

// AllDamageTypes 返回所有伤害类型
func AllDamageTypes() []DamageType {
	return []DamageType{
		DamageTypeAcid,
		DamageTypeBludgeoning,
		DamageTypeCold,
		DamageTypeFire,
		DamageTypeForce,
		DamageTypeLightning,
		DamageTypeNecrotic,
		DamageTypePiercing,
		DamageTypePoison,
		DamageTypePsychic,
		DamageTypeRadiant,
		DamageTypeSlashing,
		DamageTypeThunder,
	}
}

// DamageResistance 代表对某种伤害类型的抗性/免疫/弱点
type DamageResistance struct {
	Type          DamageType `json:"type"`
	Resistance    bool       `json:"resistance"`    // 抗性（伤害减半）
	Immunity      bool       `json:"immunity"`      // 免疫（伤害为零）
	Vulnerability bool       `json:"vulnerability"` // 弱点（伤害翻倍）
}

// DamageResistances 存储所有伤害类型的抗性/免疫/弱点
type DamageResistances struct {
	Resistances     map[DamageType]bool `json:"resistances"`
	Immunities      map[DamageType]bool `json:"immunities"`
	Vulnerabilities map[DamageType]bool `json:"vulnerabilities"`
}

// NewDamageResistances 创建新的伤害抗性配置
func NewDamageResistances() *DamageResistances {
	return &DamageResistances{
		Resistances:     make(map[DamageType]bool),
		Immunities:      make(map[DamageType]bool),
		Vulnerabilities: make(map[DamageType]bool),
	}
}

// AddResistance 添加抗性
func (dr *DamageResistances) AddResistance(damageType DamageType) {
	dr.Resistances[damageType] = true
	// 如果有免疫，移除抗性（免疫优先）
	delete(dr.Immunities, damageType)
}

// AddImmunity 添加免疫
func (dr *DamageResistances) AddImmunity(damageType DamageType) {
	dr.Immunities[damageType] = true
	// 移除抗性（免疫优先）
	delete(dr.Resistances, damageType)
	// 移除弱点（免疫优先）
	delete(dr.Vulnerabilities, damageType)
}

// AddVulnerability 添加弱点
func (dr *DamageResistances) AddVulnerability(damageType DamageType) {
	// 如果已经有免疫，不添加弱点
	if dr.Immunities[damageType] {
		return
	}
	dr.Vulnerabilities[damageType] = true
}

// HasResistance 检查是否有抗性
func (dr *DamageResistances) HasResistance(damageType DamageType) bool {
	return dr.Resistances[damageType]
}

// HasImmunity 检查是否免疫
func (dr *DamageResistances) HasImmunity(damageType DamageType) bool {
	return dr.Immunities[damageType]
}

// HasVulnerability 检查是否有弱点
func (dr *DamageResistances) HasVulnerability(damageType DamageType) bool {
	return dr.Vulnerabilities[damageType]
}

// CalculateDamage 根据抗性/免疫/弱点计算最终伤害
func (dr *DamageResistances) CalculateDamage(baseDamage int, damageType DamageType) int {
	// 免疫：伤害为0
	if dr.HasImmunity(damageType) {
		return 0
	}

	finalDamage := baseDamage

	// 弱点：伤害翻倍
	if dr.HasVulnerability(damageType) {
		finalDamage *= 2
	}

	// 抗性：伤害减半
	if dr.HasResistance(damageType) {
		finalDamage /= 2
	}

	return finalDamage
}

// DamageInstance 代表一次伤害实例
type DamageInstance struct {
	Amount     int        `json:"amount"`
	Type       DamageType `json:"type"`
	Source     string     `json:"source"`      // 伤害来源
	SourceID   ID         `json:"source_id"`   // 伤害来源ID
	IsCritical bool       `json:"is_critical"` // 是否是暴击
}
