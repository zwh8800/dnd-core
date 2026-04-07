package model

// SpellSchool 代表法术学派
type SpellSchool string

const (
	SpellSchoolAbjuration    SpellSchool = "Abjuration"    // 防护系
	SpellSchoolConjuration   SpellSchool = "Conjuration"   // 咒法系
	SpellSchoolDivination    SpellSchool = "Divination"    // 预言系
	SpellSchoolEnchantment   SpellSchool = "Enchantment"   // 附魔系
	SpellSchoolEvocation     SpellSchool = "Evocation"     // 塑能系
	SpellSchoolIllusion      SpellSchool = "Illusion"      // 幻术系
	SpellSchoolNecromancy    SpellSchool = "Necromancy"    // 死灵系
	SpellSchoolTransmutation SpellSchool = "Transmutation" // 变化系
)

// SpellComponent 代表法术成分
type SpellComponent string

const (
	SpellComponentVerbal   SpellComponent = "V" // 言语
	SpellComponentSomatic  SpellComponent = "S" // 姿势
	SpellComponentMaterial SpellComponent = "M" // 材料
)

// SpellCastTime 代表施法时间
type SpellCastTime struct {
	Value int    `json:"value"`
	Unit  string `json:"unit"` // "action", "bonus_action", "reaction", "minute"
}

// Spell 代表一个法术定义
type Spell struct {
	ID             string           `json:"id"`
	Name           string           `json:"name"`
	Level          int              `json:"level"` // 法术等级（0-9）
	School         SpellSchool      `json:"school"`
	CastTime       SpellCastTime    `json:"cast_time"`
	Range          string           `json:"range"`
	Components     []SpellComponent `json:"components"`
	Materials      string           `json:"materials,omitempty"` // 材料描述
	Duration       string           `json:"duration"`
	Concentration  bool             `json:"concentration"` // 是否需要专注
	Ritual         bool             `json:"ritual"`        // 是否是仪式法术
	Description    string           `json:"description"`
	AtHigherLevels string           `json:"at_higher_levels,omitempty"` // 高环效果

	// 效果相关
	DamageDice  string     `json:"damage_dice,omitempty"`
	DamageType  DamageType `json:"damage_type,omitempty"`
	SaveDC      Ability    `json:"save_dc,omitempty"` // 豁免属性
	HealingDice string     `json:"healing_dice,omitempty"`

	// 施法职业
	Classes []string `json:"classes"`
}

// SpellSlotTracker 追踪法术位的使用情况
type SpellSlotTracker struct {
	// Slots 存储每个环级的法术位 [总数量, 已使用数量]
	// 索引0对应戏法（无限），索引1-9对应1-9环
	Slots [10][2]int `json:"slots"`
}

// NewSpellSlotTracker 创建新的法术位追踪器
func NewSpellSlotTracker(slots [10]int) *SpellSlotTracker {
	tracker := &SpellSlotTracker{}
	for i := 0; i < 10; i++ {
		tracker.Slots[i] = [2]int{slots[i], 0}
	}
	return tracker
}

// GetAvailableSlots 获取指定环级的可用法术位数量
func (st *SpellSlotTracker) GetAvailableSlots(level int) int {
	if level < 0 || level > 9 {
		return 0
	}
	total := st.Slots[level][0]
	used := st.Slots[level][1]
	return total - used
}

// UseSlot 使用一个指定环级的法术位
func (st *SpellSlotTracker) UseSlot(level int) bool {
	if level < 1 || level > 9 {
		return false // 戏法不需要法术位
	}
	available := st.GetAvailableSlots(level)
	if available <= 0 {
		return false
	}
	st.Slots[level][1]++
	return true
}

// RestoreSlot 恢复一个指定环级的法术位
func (st *SpellSlotTracker) RestoreSlot(level int) {
	if level < 1 || level > 9 {
		return
	}
	if st.Slots[level][1] > 0 {
		st.Slots[level][1]--
	}
}

// RestoreAll 恢复所有法术位
func (st *SpellSlotTracker) RestoreAll() {
	for i := 1; i < 10; i++ {
		st.Slots[i][1] = 0
	}
}

// HasSlotAvailable 检查是否有指定环级或更高环级的可用法术位
func (st *SpellSlotTracker) HasSlotAvailable(minLevel int) bool {
	for level := minLevel; level <= 9; level++ {
		if st.GetAvailableSlots(level) > 0 {
			return true
		}
	}
	return false
}

// GetLowestAvailableSlot 获取最低可用法术位的环级
func (st *SpellSlotTracker) GetLowestAvailableSlot(minLevel int) int {
	for level := minLevel; level <= 9; level++ {
		if st.GetAvailableSlots(level) > 0 {
			return level
		}
	}
	return 0
}

// SpellcasterState 代表施法者的状态
type SpellcasterState struct {
	SpellcastingAbility Ability           `json:"spellcasting_ability"`
	SpellSaveDC         int               `json:"spell_save_dc"`
	SpellAttackBonus    int               `json:"spell_attack_bonus"`
	Slots               *SpellSlotTracker `json:"slots"`
	PreparedSpells      []string          `json:"prepared_spells"`     // 已准备的法术
	KnownSpells         []string          `json:"known_spells"`        // 已知的法术
	PreparationType     string            `json:"preparation_type"`    // "prepared" 或 "known"
	ConcentrationSpell  string            `json:"concentration_spell"` // 当前专注的法术
}

// CanPrepareSpell 检查是否可以准备此法术
func (sc *SpellcasterState) CanPrepareSpell(spellID string) bool {
	if sc.PreparationType == "known" {
		return false // 已知型施法者不需要准备
	}
	// 检查是否在已知法术列表中
	for _, s := range sc.KnownSpells {
		if s == spellID {
			return true
		}
	}
	return false
}

// IsConcentrating 检查是否正在专注某个法术
func (sc *SpellcasterState) IsConcentrating() bool {
	return sc.ConcentrationSpell != ""
}
