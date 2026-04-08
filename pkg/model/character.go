package model

// RaceReference 引用一个种族定义
type RaceReference struct {
	Name string `json:"name"`
	// Subrace 可选的子种族
	Subrace string `json:"subrace,omitempty"`
}

// HitDiceEntry 代表一种生命骰
type HitDiceEntry struct {
	DiceType int `json:"dice_type"` // d6, d8, d10, d12
	Total    int `json:"total"`     // 总共有多少颗
	Used     int `json:"used"`      // 已使用的数量
}

// Alignment 代表阵营
type Alignment string

const (
	AlignmentLawfulGood     Alignment = "Lawful Good"     // 守序善良
	AlignmentNeutralGood    Alignment = "Neutral Good"    // 中立善良
	AlignmentChaoticGood    Alignment = "Chaotic Good"    // 混乱善良
	AlignmentLawfulNeutral  Alignment = "Lawful Neutral"  // 守序中立
	AlignmentTrueNeutral    Alignment = "True Neutral"    // 绝对中立
	AlignmentChaoticNeutral Alignment = "Chaotic Neutral" // 混乱中立
	AlignmentLawfulEvil     Alignment = "Lawful Evil"     // 守序邪恶
	AlignmentNeutralEvil    Alignment = "Neutral Evil"    // 中立邪恶
	AlignmentChaoticEvil    Alignment = "Chaotic Evil"    // 混乱邪恶
)

// PersonalityTraits 代表个性特征
type PersonalityTraits struct {
	Traits     string `json:"traits"`     // 个性特征
	Ideals     string `json:"ideals"`     // 理想
	Bonds      string `json:"bonds"`      // 羁绊
	Flaws      string `json:"flaws"`      // 缺陷
	Background string `json:"background"` // 背景
}

// PlayerCharacter 代表玩家角色
type PlayerCharacter struct {
	Actor

	// 种族与职业
	Race       RaceReference `json:"race"`
	Classes    []ClassLevel  `json:"classes"`
	TotalLevel int           `json:"total_level"`
	Alignment  Alignment     `json:"alignment"`

	// 背景系统
	BackgroundID string `json:"background_id,omitempty"` // 背景 ID

	// 经验与升级
	Experience         int `json:"experience"`
	DeathSaveSuccesses int `json:"death_save_successes"`
	DeathSaveFailures  int `json:"death_save_failures"`

	// 激励与生命骰
	Inspiration bool           `json:"inspiration"`
	HitDice     []HitDiceEntry `json:"hit_dice"`

	// 库存与法术
	InventoryID  ID                `json:"inventory_id"`
	Spellcasting *SpellcasterState `json:"spellcasting,omitempty"`

	// 个性
	Personality *PersonalityTraits `json:"personality,omitempty"`

	// 特性与能力
	Features     []string `json:"features"`
	RacialTraits []string `json:"racial_traits"`

	// 专长系统
	Feats []FeatInstance `json:"feats,omitempty"` // 角色已获得的专长

	// 制作系统
	CraftingProgress map[string]*CraftingProgress `json:"crafting_progress,omitempty"` // 制作进度
	Gold             int                          `json:"gold"`                        // 金币数量（铜币单位）

	// 职业特性系统
	FeatureHooks map[ClassID]FeatureHook `json:"-"`                       // 运行时特性钩子(不序列化)
	FighterState *FighterFeatures        `json:"fighter_state,omitempty"` // 战士特性状态
}

// GetSpellcastingAbility 获取施法属性（基于主要施法职业）
func (pc *PlayerCharacter) GetSpellcastingAbility() Ability {
	if pc.Spellcasting == nil {
		return ""
	}
	return pc.Spellcasting.SpellcastingAbility
}

// IsSpellcaster 检查是否是施法者
func (pc *PlayerCharacter) IsSpellcaster() bool {
	return pc.Spellcasting != nil
}

// NPC 代表非玩家角色
type NPC struct {
	Actor

	// NPC特有字段
	Occupation  string   `json:"occupation"`  // 职业/身份
	Faction     string   `json:"faction"`     // 所属组织
	Attitude    string   `json:"attitude"`    // 对玩家的态度
	Disposition string   `json:"disposition"` // 性格倾向
	KnownInfo   []string `json:"known_info"`  // 知道的信息
	QuestGiver  bool     `json:"quest_giver"` // 是否能给予任务
	Merchant    bool     `json:"merchant"`    // 是否是商人

	// 社交互动状态
	SocialState *SocialInteractionState `json:"social_state,omitempty"`
}

// Enemy 代表敌人/怪物
type Enemy struct {
	Actor

	// 敌人特有字段
	ChallengeRating       string           `json:"challenge_rating"` // 挑战等级（如 "1/4", "2", "10"）
	XPValue               int              `json:"xp_value"`         // 经验值
	AttackBonus           int              `json:"attack_bonus"`     // 攻击加值
	DamagePerRound        int              `json:"damage_per_round"` // 每回合伤害
	Senses                []string         `json:"senses"`           // 感官（黑暗视觉等）
	DamageImmunities      []DamageImmunity `json:"damage_immunities"`
	DamageResistances     []DamageImmunity `json:"damage_resistances"`
	DamageVulnerabilities []DamageImmunity `json:"damage_vulnerabilities"`
	ConditionImmunities   []ConditionType  `json:"condition_immunities"`

	// 怪物模板引用（如果使用模板创建）
	StatBlock *MonsterStatBlock `json:"-"`

	// 传说动作追踪
	LegendaryActionsRemaining int `json:"legendary_actions_remaining"`

	// 充能动作追踪（动作索引 -> 剩余次数）
	ActionRecharges map[int]int `json:"action_recharges"`
}

// Companion 代表同伴/盟友（AI控制）
type Companion struct {
	Actor

	// 同伴特有字段
	LeaderID     ID       `json:"leader_id"`     // 领导者（玩家）ID
	Loyalty      int      `json:"loyalty"`       // 忠诚度
	BehaviorMode string   `json:"behavior_mode"` // 行为模式（攻击/防御/辅助）
	Commands     []string `json:"commands"`      // 可接受的命令
}
