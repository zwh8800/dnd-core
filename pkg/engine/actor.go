package engine

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/zwh8800/dnd-core/pkg/data"
	"github.com/zwh8800/dnd-core/pkg/model"
	"github.com/zwh8800/dnd-core/pkg/rules"
)

// AbilityScoresInput 属性值输入
type AbilityScoresInput struct {
	Strength     int `json:"strength"`     // 力量
	Dexterity    int `json:"dexterity"`    // 敏捷
	Constitution int `json:"constitution"` // 体质
	Intelligence int `json:"intelligence"` // 智力
	Wisdom       int `json:"wisdom"`       // 感知
	Charisma     int `json:"charisma"`     // 魅力
}

// PersonalityTraitsInput 个性特征输入
type PersonalityTraitsInput struct {
	Traits string `json:"traits"` // 个性特征
	Ideals string `json:"ideals"` // 理想
	Bonds  string `json:"bonds"`  // 羁绊
	Flaws  string `json:"flaws"`  // 缺陷
}

// FeatInstanceInput 专长实例输入
type FeatInstanceInput struct {
	FeatID        string `json:"feat_id"`        // 专长ID
	Source        string `json:"source"`         // 来源（background/level_up/variant）
	AcquiredLevel int    `json:"acquired_level"` // 获得时的等级
}

// ClassLevelInput 职业等级输入
type ClassLevelInput struct {
	Class string `json:"class"` // 职业名称
	Level int    `json:"level"` // 职业等级
}

// SpellcasterStateInput 施法者状态输入
type SpellcasterStateInput struct {
	SpellcastingAbility string   `json:"spellcasting_ability"`      // 施法属性
	PreparedSpells      []string `json:"prepared_spells,omitempty"` // 已准备的法术
	KnownSpells         []string `json:"known_spells,omitempty"`    // 已知的法术
	PreparationType     string   `json:"preparation_type"`          // 准备类型（prepared/known）
}

// CraftingProgressInput 制作进度输入
type CraftingProgressInput struct {
	RecipeID   string `json:"recipe_id"`   // 配方ID
	DaysWorked int    `json:"days_worked"` // 已工作天数
	MoneySpent int    `json:"money_spent"` // 已花费金额
}

// FighterStateInput 战士特性状态输入
type FighterStateInput struct {
	FightingStyle string `json:"fighting_style"` // 战斗风格
	Archetype     string `json:"archetype"`      // 武术范型
}

// PlayerCharacterInput 玩家角色创建输入
type PlayerCharacterInput struct {
	// 基础信息
	Name          string             `json:"name"`              // 角色名称
	Race          string             `json:"race"`              // 种族
	Subrace       string             `json:"subrace,omitempty"` // 子种族（可选）
	Background    string             `json:"background"`        // 背景ID
	Class         string             `json:"class"`             // 主职业
	Level         int                `json:"level"`             // 等级（默认1）
	Alignment     string             `json:"alignment"`         // 阵营
	AbilityScores AbilityScoresInput `json:"ability_scores"`    // 属性值

	// 可选配置
	Description string `json:"description,omitempty"` // 角色描述
	Size        string `json:"size,omitempty"`        // 体型（默认继承种族）
	Speed       *int   `json:"speed,omitempty"`       // 速度（可选，不填则使用种族默认值）

	// 个性特征（可选）
	Personality *PersonalityTraitsInput `json:"personality,omitempty"` // 个性特征

	// 专长（可选，通常来自变体人类或背景）
	Feats []FeatInstanceInput `json:"feats,omitempty"` // 初始专长列表

	// 施法者配置（可选，仅当职业是施法者时需要）
	Spellcasting *SpellcasterStateInput `json:"spellcasting,omitempty"` // 施法者状态

	// 经济系统
	Gold int `json:"gold,omitempty"` // 初始金币（铜币单位，默认使用背景起始财富）

	// 制作进度（可选，用于继承已有角色）
	CraftingProgress map[string]*CraftingProgressInput `json:"crafting_progress,omitempty"` // 制作进度

	// 职业特性配置
	FighterState *FighterStateInput `json:"fighter_state,omitempty"` // 战士特性（仅战士职业需要）

	// 生命值（可选，不填则自动计算）
	HitPoints int `json:"hit_points,omitempty"` // 初始HP
}

// NPCInput NPC创建输入
type NPCInput struct {
	// 基础信息
	Name          string             `json:"name"`           // NPC名称
	Description   string             `json:"description"`    // 描述
	Size          string             `json:"size"`           // 体型（Tiny/Small/Medium/Large/Huge/Gargantuan）
	CreatureType  string             `json:"creature_type"`  // 生物类型（Aberration/Beast/Celestial等）
	Speed         int                `json:"speed"`          // 基础速度（英尺）
	AbilityScores AbilityScoresInput `json:"ability_scores"` // 属性值

	// NPC特有字段
	Occupation  string   `json:"occupation"`  // 职业/身份
	Faction     string   `json:"faction"`     // 所属组织
	Attitude    string   `json:"attitude"`    // 对玩家的态度（friendly/indifferent/hostile）
	Disposition string   `json:"disposition"` // 性格倾向（helpful/friendly/indifferent/suspicious/hostile）
	KnownInfo   []string `json:"known_info"`  // 知道的信息
	QuestGiver  bool     `json:"quest_giver"` // 是否能给予任务
	Merchant    bool     `json:"merchant"`    // 是否是商人

	// 可选配置
	ArmorClass      *int `json:"armor_class,omitempty"`      // 护甲等级（不填则自动计算）
	HitPoints       *int `json:"hit_points,omitempty"`       // 生命值（不填则使用默认值10）
	InitiativeBonus *int `json:"initiative_bonus,omitempty"` // 先攻修正
	Inspiration     bool `json:"inspiration"`                // 是否拥有灵感点
}

// DamageImmunityInput 伤害免疫/抗性/易伤输入
type DamageImmunityInput struct {
	DamageTypes []string `json:"damage_types"` // 伤害类型
	NonMagical  bool     `json:"non_magical"`  // 是否仅对非魔法攻击有效
}

// EnemyInput 敌人创建输入
type EnemyInput struct {
	// 基础信息
	Name            string             `json:"name"`             // 敌人名称
	Description     string             `json:"description"`      // 描述
	Size            string             `json:"size"`             // 体型
	CreatureType    string             `json:"creature_type"`    // 生物类型
	Speed           int                `json:"speed"`            // 基础速度（英尺）
	AbilityScores   AbilityScoresInput `json:"ability_scores"`   // 属性值
	ChallengeRating string             `json:"challenge_rating"` // 挑战等级（如 "1/4", "2", "10"）
	HitPoints       int                `json:"hit_points"`       // 生命值
	ArmorClass      int                `json:"armor_class"`      // 护甲等级

	// 战斗属性
	XPValue               int                   `json:"xp_value"`               // 经验值
	AttackBonus           int                   `json:"attack_bonus"`           // 攻击加值
	DamagePerRound        int                   `json:"damage_per_round"`       // 每回合伤害
	Senses                []string              `json:"senses"`                 // 感官（黑暗视觉等）
	DamageImmunities      []DamageImmunityInput `json:"damage_immunities"`      // 伤害免疫
	DamageResistances     []DamageImmunityInput `json:"damage_resistances"`     // 伤害抗性
	DamageVulnerabilities []DamageImmunityInput `json:"damage_vulnerabilities"` // 伤害易伤
	ConditionImmunities   []string              `json:"condition_immunities"`   // 状态免疫

	// 传说动作
	LegendaryActionsRemaining int `json:"legendary_actions_remaining"` // 传说动作剩余次数

	// 可选配置
	InitiativeBonus *int `json:"initiative_bonus,omitempty"` // 先攻修正
}

// CompanionInput 同伴创建输入
type CompanionInput struct {
	// 基础信息
	Name          string             `json:"name"`           // 同伴名称
	Description   string             `json:"description"`    // 描述
	Size          string             `json:"size"`           // 体型
	CreatureType  string             `json:"creature_type"`  // 生物类型
	Speed         int                `json:"speed"`          // 速度（英尺）
	AbilityScores AbilityScoresInput `json:"ability_scores"` // 属性值

	// 同伴特有字段
	LeaderID     string   `json:"leader_id"`     // 领导者（玩家）ID
	Loyalty      int      `json:"loyalty"`       // 忠诚度
	BehaviorMode string   `json:"behavior_mode"` // 行为模式（攻击/防御/辅助）
	Commands     []string `json:"commands"`      // 可接受的命令

	// 可选配置
	ArmorClass      *int `json:"armor_class,omitempty"`      // 护甲等级（不填则自动计算）
	HitPoints       *int `json:"hit_points,omitempty"`       // 生命值（不填则自动计算）
	InitiativeBonus *int `json:"initiative_bonus,omitempty"` // 先攻修正
}

// ActorInfo 角色基本信息
type ActorInfo struct {
	ID         model.ID        `json:"id"`          // 角色唯一标识
	Type       model.ActorType `json:"type"`        // 角色类型
	Name       string          `json:"name"`        // 角色名称
	HitPoints  model.HitPoints `json:"hit_points"`  // 生命值
	TempHP     int             `json:"temp_hp"`     // 临时HP
	ArmorClass int             `json:"armor_class"` // 护甲等级
	Speed      int             `json:"speed"`       // 移动速度
	Conditions []string        `json:"conditions"`  // 状态效果列表
	Exhaustion int             `json:"exhaustion"`  // 力竭等级
	SceneID    model.ID        `json:"scene_id"`    // 所在场景ID
	Position   *model.Point    `json:"position"`    // 位置坐标
}

// PlayerCharacterInfo 玩家角色完整信息
type PlayerCharacterInfo struct {
	ID               model.ID           `json:"id"`                // 角色唯一标识
	Name             string             `json:"name"`              // 角色名称
	Race             string             `json:"race"`              // 种族
	Background       string             `json:"background"`        // 背景
	Classes          []ClassInfo        `json:"classes"`           // 职业信息
	TotalLevel       int                `json:"total_level"`       // 总等级
	Experience       int                `json:"experience"`        // 经验值
	AbilityScores    AbilityScoresInput `json:"ability_scores"`    // 属性值
	HitPoints        model.HitPoints    `json:"hit_points"`        // 生命值
	ArmorClass       int                `json:"armor_class"`       // 护甲等级
	Speed            int                `json:"speed"`             // 移动速度
	ProficiencyBonus int                `json:"proficiency_bonus"` // 熟练加值
	Features         []string           `json:"features"`          // 特性列表
	RacialTraits     []string           `json:"racial_traits"`     // 种族特性
}

// ClassInfo 职业信息
type ClassInfo struct {
	Class      model.ClassID `json:"class"`              // 职业ID
	ClassLevel int           `json:"class_level"`        // 职业等级
	Features   []string      `json:"features,omitempty"` // 职业特性
}

// ActorFilter 角色过滤条件
type ActorFilter struct {
	Types   []model.ActorType `json:"types,omitempty"`    // 角色类型过滤列表
	SceneID *model.ID         `json:"scene_id,omitempty"` // 场景ID过滤
	Alive   *bool             `json:"alive,omitempty"`    // 存活状态过滤
}

// ActorUpdate 角色更新内容
type ActorUpdate struct {
	AbilityScores *AbilityScoresInput `json:"ability_scores,omitempty"` // 属性值更新
	HitPoints     *HitPointUpdate     `json:"hit_points,omitempty"`     // HP更新
	Conditions    *ConditionUpdate    `json:"conditions,omitempty"`     // 状态效果更新
	Position      *model.Point        `json:"position,omitempty"`       // 位置更新
	SceneID       *model.ID           `json:"scene_id,omitempty"`       // 场景ID更新
	Custom        map[string]any      `json:"custom,omitempty"`         // 自定义字段
}

// HitPointUpdate HP更新
type HitPointUpdate struct {
	Current       *int `json:"current,omitempty"`         // 当前HP
	TempHitPoints *int `json:"temp_hit_points,omitempty"` // 临时HP
}

// ConditionUpdate 状态效果更新
type ConditionUpdate struct {
	Add    []model.ConditionInstance `json:"add,omitempty"`    // 添加的状态效果
	Remove []model.ConditionType     `json:"remove,omitempty"` // 移除的状态效果类型
}

// LevelUpResult 升级结果
type LevelUpResult struct {
	OldLevel             int      `json:"old_level"`             // 原等级
	NewLevel             int      `json:"new_level"`             // 新等级
	HPGain               int      `json:"hp_gain"`               // HP增长值
	NewFeatures          []string `json:"new_features"`          // 新获得特性
	SpellSlotsUpdated    bool     `json:"spell_slots_updated"`   // 法术位是否更新
	ProficiencyIncreased bool     `json:"proficiency_increased"` // 熟练加值是否增加
	Message              string   `json:"message"`               // 人类可读消息
}

// RestResult 休息结果
type RestResult struct {
	ActorResults []ActorRestResult `json:"actor_results"` // 各角色休息结果
	Message      string            `json:"message"`       // 人类可读消息
}

// ActorRestResult 角色休息结果
type ActorRestResult struct {
	ActorID            model.ID              `json:"actor_id"`             // 角色ID
	HPRecovered        int                   `json:"hp_recovered"`         // 恢复的HP
	HitDiceUsed        int                   `json:"hit_dice_used"`        // 使用的生命骰数量
	SpellSlotsRestored bool                  `json:"spell_slots_restored"` // 法术位是否恢复
	ConditionsRemoved  []model.ConditionType `json:"conditions_removed"`   // 移除的状态效果
	ExhaustionReduced  bool                  `json:"exhaustion_reduced"`   // 力竭是否减少
	AbilitiesRestored  bool                  `json:"abilities_restored"`   // 能力是否恢复
}

// Request 结构体定义

// CreatePCRequest 创建玩家角色请求
type CreatePCRequest struct {
	GameID model.ID              `json:"game_id"` // 游戏会话ID
	PC     *PlayerCharacterInput `json:"pc"`      // 玩家角色创建参数
}

// CreatePCResult 创建玩家角色结果
type CreatePCResult struct {
	Actor *ActorInfo `json:"actor"` // 创建的角色信息
}

// CreateNPCRequest 创建NPC请求
type CreateNPCRequest struct {
	GameID model.ID  `json:"game_id"` // 游戏会话ID
	NPC    *NPCInput `json:"npc"`     // NPC创建参数
}

// CreateNPCResult 创建NPC结果
type CreateNPCResult struct {
	Actor *ActorInfo `json:"actor"` // 创建的角色信息
}

// CreateEnemyRequest 创建敌人请求
type CreateEnemyRequest struct {
	GameID model.ID    `json:"game_id"` // 游戏会话ID
	Enemy  *EnemyInput `json:"enemy"`   // 敌人创建参数
}

// CreateEnemyResult 创建敌人结果
type CreateEnemyResult struct {
	Actor *ActorInfo `json:"actor"` // 创建的角色信息
}

// CreateCompanionRequest 创建同伴请求
type CreateCompanionRequest struct {
	GameID    model.ID        `json:"game_id"`   // 游戏会话ID
	Companion *CompanionInput `json:"companion"` // 同伴创建参数
}

// CreateCompanionResult 创建同伴结果
type CreateCompanionResult struct {
	Actor *ActorInfo `json:"actor"` // 创建的角色信息
}

// GetActorRequest 获取角色请求
type GetActorRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID
	ActorID model.ID `json:"actor_id"` // 角色ID
}

// GetActorResult 获取角色结果
type GetActorResult struct {
	Actor *ActorInfo `json:"actor"` // 角色信息
}

// GetPCRequest 获取玩家角色请求
type GetPCRequest struct {
	GameID model.ID `json:"game_id"` // 游戏会话ID
	PCID   model.ID `json:"pc_id"`   // 玩家角色ID
}

// GetPCResult 获取玩家角色结果
type GetPCResult struct {
	PC *PlayerCharacterInfo `json:"pc"` // 玩家角色完整信息
}

// UpdateActorRequest 更新角色请求
type UpdateActorRequest struct {
	GameID  model.ID    `json:"game_id"`  // 游戏会话ID
	ActorID model.ID    `json:"actor_id"` // 角色ID
	Update  ActorUpdate `json:"update"`   // 更新内容
}

// RemoveActorRequest 移除角色请求
type RemoveActorRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID
	ActorID model.ID `json:"actor_id"` // 角色ID
}

// ListActorsRequest 列出角色请求
type ListActorsRequest struct {
	GameID model.ID     `json:"game_id"` // 游戏会话ID
	Filter *ActorFilter `json:"filter"`  // 过滤条件（可选）
}

// ListActorsResult 列出角色结果
type ListActorsResult struct {
	Actors []ActorInfo `json:"actors"` // 角色列表
}

// AddExperienceRequest 添加经验值请求
type AddExperienceRequest struct {
	GameID model.ID `json:"game_id"` // 游戏会话ID
	PCID   model.ID `json:"pc_id"`   // 玩家角色ID
	XP     int      `json:"xp"`      // 添加的经验值
}

// AddExperienceResult 添加经验值结果
type AddExperienceResult struct {
	LeveledUp bool `json:"leveled_up"` // 是否升级
	OldLevel  int  `json:"old_level"`  // 原等级
	NewLevel  int  `json:"new_level"`  // 新等级
}

// LevelUpRequest 升级请求
type LevelUpRequest struct {
	GameID      model.ID `json:"game_id"`      // 游戏会话ID
	PCID        model.ID `json:"pc_id"`        // 玩家角色ID
	ClassChoice string   `json:"class_choice"` // 升级职业选择
}

// ShortRestRequest 短休请求
type ShortRestRequest struct {
	GameID       model.ID         `json:"game_id"`        // 游戏会话ID
	ActorIDs     []model.ID       `json:"actor_ids"`      // 参与短休的角色ID列表
	HitDiceCount map[model.ID]int `json:"hit_dice_count"` // 每个角色想要使用的生命骰数量(可选,默认为1)
}

// StartLongRestRequest 开始长休请求
type StartLongRestRequest struct {
	GameID   model.ID   `json:"game_id"`   // 游戏会话ID
	ActorIDs []model.ID `json:"actor_ids"` // 参与长休的角色ID列表
}

// EndLongRestRequest 结束长休请求
type EndLongRestRequest struct {
	GameID model.ID `json:"game_id"` // 游戏会话ID
}

// CreatePC 创建一个新的玩家角色（Player Character）
// 根据提供的角色配置创建PC，包括种族、职业、背景、属性、专长、法术位等，并自动计算HP、AC等派生值。
// 参数:
//
//	ctx - 上下文
//	req - 创建请求，包含游戏会话ID和玩家角色配置
//
// 返回:
//
//	*CreatePCResult - 包含创建的角色信息
//	error - 可能返回 ErrNotFound（游戏不存在）、权限错误或职业无效等错误
func (e *Engine) CreatePC(ctx context.Context, req CreatePCRequest) (*CreatePCResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	// 检查权限
	if err := e.checkPermission(game.Phase, OpCreatePC); err != nil {
		return nil, err
	}

	if req.PC == nil {
		return nil, fmt.Errorf("pc input is required")
	}

	// 校验种族有效性
	if req.PC.Race != "" {
		if data.GetRace(req.PC.Race) == nil {
			return nil, fmt.Errorf("无效的种族: %s，可用种族: %s", req.PC.Race, strings.Join(data.GetRaceNames(), ", "))
		}
	}

	// 构建 PlayerCharacter
	pc := &model.PlayerCharacter{
		Actor: model.Actor{
			ID:            model.NewID(),
			Name:          req.PC.Name,
			Description:   req.PC.Description,
			AbilityScores: abilityScoresInputToModel(req.PC.AbilityScores),
			HitPoints:     model.HitPoints{},
			Conditions:    []model.ConditionInstance{},
			Exhaustion:    0,
		},
		Race: model.RaceReference{
			Name:    req.PC.Race,
			Subrace: req.PC.Subrace,
		},
		Alignment:    model.Alignment(req.PC.Alignment),
		Classes:      []model.ClassLevel{},
		Experience:   0,
		Features:     []string{},
		RacialTraits: []string{},
		Feats:        []model.FeatInstance{},
		Gold:         req.PC.Gold,
	}

	// 设置体型和速度（如果提供）
	if req.PC.Size != "" {
		pc.Size = model.Size(req.PC.Size)
	}
	if req.PC.Speed != nil {
		pc.Speed = *req.PC.Speed
	}

	// 设置背景
	if req.PC.Background != "" {
		bgID, ok := data.ResolveBackgroundID(req.PC.Background)
		if !ok {
			return nil, fmt.Errorf("无效的背景: %s，可用背景: %s", req.PC.Background, strings.Join(data.GetBackgroundNames(), ", "))
		}
		pc.BackgroundID = string(bgID)
	}

	// 设置个性特征
	if req.PC.Personality != nil {
		pc.Personality = &model.PersonalityTraits{
			Traits:     req.PC.Personality.Traits,
			Ideals:     req.PC.Personality.Ideals,
			Bonds:      req.PC.Personality.Bonds,
			Flaws:      req.PC.Personality.Flaws,
			Background: req.PC.Background,
		}
	}

	// 设置初始专长
	if len(req.PC.Feats) > 0 {
		for _, featInput := range req.PC.Feats {
			acquiredLevel := featInput.AcquiredLevel
			if acquiredLevel < 1 {
				acquiredLevel = 1
			}
			pc.Feats = append(pc.Feats, model.FeatInstance{
				FeatID:        featInput.FeatID,
				Source:        model.FeatSource(featInput.Source),
				AcquiredLevel: acquiredLevel,
			})
		}
	}

	// 添加职业
	if req.PC.Class != "" {
		classID, err := data.GetClassID(req.PC.Class)
		if err != nil {
			return nil, fmt.Errorf("无效的职业: %s", req.PC.Class)
		}

		level := req.PC.Level
		if level < 1 {
			level = 1
		}

		// 获取该职业的特性列表
		features := getClassFeatures(classID, level)

		pc.Classes = append(pc.Classes, model.ClassLevel{
			Class:    classID,
			Level:    level,
			Features: features,
		})
		pc.TotalLevel = level

		// 初始化职业特性系统
		pc.FeatureHooks = make(map[model.ClassID]model.FeatureHook)
		if classID == model.ClassFighter {
			pc.FighterState = &model.FighterFeatures{}
			if req.PC.FighterState != nil {
				// 应用用户指定的战斗风格和范型
				pc.FighterState.SelectedFightingStyle = model.FightingStyle(req.PC.FighterState.FightingStyle)
				pc.FighterState.SelectedArchetype = model.MartialArchetype(req.PC.FighterState.Archetype)
			}
			model.UpdateFighterFeatures(pc.FighterState, level)
			pc.FeatureHooks[classID] = &model.FighterFeatureHooks{
				Features: pc.FighterState,
				Level:    level,
			}
		}
	}

	// 设置施法者状态
	if req.PC.Spellcasting != nil {
		pc.Spellcasting = &model.SpellcasterState{
			SpellcastingAbility: model.Ability(req.PC.Spellcasting.SpellcastingAbility),
			PreparedSpells:      req.PC.Spellcasting.PreparedSpells,
			KnownSpells:         req.PC.Spellcasting.KnownSpells,
			PreparationType:     req.PC.Spellcasting.PreparationType,
		}

		// 计算法术豁免DC和攻击加值
		if pc.Spellcasting.SpellcastingAbility != "" {
			abilityScore := pc.AbilityScores.Get(pc.Spellcasting.SpellcastingAbility)
			profBonus := rules.ProficiencyBonus(pc.TotalLevel)
			pc.Spellcasting.SpellSaveDC = 8 + profBonus + rules.AbilityModifier(abilityScore)
			pc.Spellcasting.SpellAttackBonus = profBonus + rules.AbilityModifier(abilityScore)
		}
	}

	// 设置制作进度
	if len(req.PC.CraftingProgress) > 0 {
		pc.CraftingProgress = make(map[string]*model.CraftingProgress)
		for recipeID, progress := range req.PC.CraftingProgress {
			pc.CraftingProgress[recipeID] = &model.CraftingProgress{
				RecipeID:   progress.RecipeID,
				DaysWorked: progress.DaysWorked,
				MoneySpent: progress.MoneySpent,
			}
		}
	}

	// 计算派生值
	pc.ArmorClass = calculateArmorClass(pc)
	if req.PC.HitPoints > 0 {
		pc.HitPoints.Maximum = req.PC.HitPoints
		pc.HitPoints.Current = req.PC.HitPoints
	} else {
		pc.HitPoints.Maximum = calculateMaxHP(pc)
		pc.HitPoints.Current = pc.HitPoints.Maximum
	}

	// 创建库存
	inventory := model.NewInventory(pc.ID)
	game.Inventories[inventory.ID] = inventory
	pc.InventoryID = inventory.ID

	// 添加到游戏
	game.PCs[pc.ID] = pc
	game.UpdatedAt = time.Now()

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &CreatePCResult{
		Actor: actorToInfo(&pc.Actor, model.ActorTypePC, pc.Name),
	}, nil
}

// CreateNPC 创建一个新的非玩家角色（NPC）
// 用于创建中立的NPC，如商人、村民、任务发布者等。
// 参数:
//
//	ctx - 上下文
//	req - 创建请求，包含游戏会话ID和NPC配置
//
// 返回:
//
//	*CreateNPCResult - 包含创建的角色信息
//	error - 可能返回游戏不存在或权限错误
func (e *Engine) CreateNPC(ctx context.Context, req CreateNPCRequest) (*CreateNPCResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpCreateNPC); err != nil {
		return nil, err
	}

	if req.NPC == nil {
		return nil, fmt.Errorf("npc input is required")
	}

	npc := &model.NPC{
		Actor: model.Actor{
			ID:            model.NewID(),
			Name:          req.NPC.Name,
			Description:   req.NPC.Description,
			AbilityScores: abilityScoresInputToModel(req.NPC.AbilityScores),
			Size:          model.Size(req.NPC.Size),
			CreatureType:  model.CreatureType(req.NPC.CreatureType),
			Speed:         req.NPC.Speed,
			HitPoints:     model.HitPoints{},
			Conditions:    []model.ConditionInstance{},
			Exhaustion:    0,
		},
		Occupation:  req.NPC.Occupation,
		Faction:     req.NPC.Faction,
		Attitude:    req.NPC.Attitude,
		Disposition: req.NPC.Disposition,
		KnownInfo:   req.NPC.KnownInfo,
		QuestGiver:  req.NPC.QuestGiver,
		Merchant:    req.NPC.Merchant,
	}

	// 设置可选字段
	if req.NPC.InitiativeBonus != nil {
		npc.InitiativeBonus = *req.NPC.InitiativeBonus
	}
	if req.NPC.Inspiration {
		npc.Inspiration = true
	}

	// 设置生命值
	if req.NPC.HitPoints != nil {
		npc.HitPoints.Maximum = *req.NPC.HitPoints
		npc.HitPoints.Current = *req.NPC.HitPoints
	} else {
		npc.HitPoints.Maximum = 10
		npc.HitPoints.Current = 10
	}

	// 设置护甲等级
	if req.NPC.ArmorClass != nil {
		npc.ArmorClass = *req.NPC.ArmorClass
	} else {
		npc.ArmorClass = calculateArmorClassFromActor(&npc.Actor)
	}

	game.NPCs[npc.ID] = npc

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &CreateNPCResult{
		Actor: actorToInfo(&npc.Actor, model.ActorTypeNPC, req.NPC.Name),
	}, nil
}

// CreateEnemy 创建一个新的敌人/怪物（Enemy）
// 用于创建战斗中的敌对角色，包含挑战等级、生命值、护甲等级等战斗属性。
// 参数:
//
//	ctx - 上下文
//	req - 创建请求，包含游戏会话ID和敌人配置（名称、描述、体型、速度、属性值、挑战等级、生命值、护甲等级）
//
// 返回:
//
//	*CreateEnemyResult - 包含创建的敌人角色信息
//	error - 可能返回游戏不存在或权限错误
func (e *Engine) CreateEnemy(ctx context.Context, req CreateEnemyRequest) (*CreateEnemyResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpCreateEnemy); err != nil {
		return nil, err
	}

	if req.Enemy == nil {
		return nil, fmt.Errorf("enemy input is required")
	}

	enemy := &model.Enemy{
		Actor: model.Actor{
			ID:            model.NewID(),
			Name:          req.Enemy.Name,
			Description:   req.Enemy.Description,
			AbilityScores: abilityScoresInputToModel(req.Enemy.AbilityScores),
			Size:          model.Size(req.Enemy.Size),
			CreatureType:  model.CreatureType(req.Enemy.CreatureType),
			Speed:         req.Enemy.Speed,
			HitPoints:     model.HitPoints{Current: req.Enemy.HitPoints, Maximum: req.Enemy.HitPoints},
			ArmorClass:    req.Enemy.ArmorClass,
			Conditions:    []model.ConditionInstance{},
			Exhaustion:    0,
		},
		ChallengeRating:           req.Enemy.ChallengeRating,
		XPValue:                   req.Enemy.XPValue,
		AttackBonus:               req.Enemy.AttackBonus,
		DamagePerRound:            req.Enemy.DamagePerRound,
		Senses:                    req.Enemy.Senses,
		LegendaryActionsRemaining: req.Enemy.LegendaryActionsRemaining,
	}

	// 设置先攻修正
	if req.Enemy.InitiativeBonus != nil {
		enemy.InitiativeBonus = *req.Enemy.InitiativeBonus
	}

	// 转换伤害免疫
	if len(req.Enemy.DamageImmunities) > 0 {
		enemy.DamageImmunities = make([]model.DamageImmunity, len(req.Enemy.DamageImmunities))
		for i, imm := range req.Enemy.DamageImmunities {
			damageTypes := make([]model.DamageType, len(imm.DamageTypes))
			for j, dt := range imm.DamageTypes {
				damageTypes[j] = model.DamageType(dt)
			}
			enemy.DamageImmunities[i] = model.DamageImmunity{
				DamageTypes: damageTypes,
				NonMagical:  imm.NonMagical,
			}
		}
	}

	// 转换伤害抗性
	if len(req.Enemy.DamageResistances) > 0 {
		enemy.DamageResistances = make([]model.DamageImmunity, len(req.Enemy.DamageResistances))
		for i, res := range req.Enemy.DamageResistances {
			damageTypes := make([]model.DamageType, len(res.DamageTypes))
			for j, dt := range res.DamageTypes {
				damageTypes[j] = model.DamageType(dt)
			}
			enemy.DamageResistances[i] = model.DamageImmunity{
				DamageTypes: damageTypes,
				NonMagical:  res.NonMagical,
			}
		}
	}

	// 转换伤害易伤
	if len(req.Enemy.DamageVulnerabilities) > 0 {
		enemy.DamageVulnerabilities = make([]model.DamageImmunity, len(req.Enemy.DamageVulnerabilities))
		for i, vul := range req.Enemy.DamageVulnerabilities {
			damageTypes := make([]model.DamageType, len(vul.DamageTypes))
			for j, dt := range vul.DamageTypes {
				damageTypes[j] = model.DamageType(dt)
			}
			enemy.DamageVulnerabilities[i] = model.DamageImmunity{
				DamageTypes: damageTypes,
				NonMagical:  vul.NonMagical,
			}
		}
	}

	// 转换状态免疫
	if len(req.Enemy.ConditionImmunities) > 0 {
		enemy.ConditionImmunities = make([]model.ConditionType, len(req.Enemy.ConditionImmunities))
		for i, cond := range req.Enemy.ConditionImmunities {
			enemy.ConditionImmunities[i] = model.ConditionType(cond)
		}
	}

	game.Enemies[enemy.ID] = enemy

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &CreateEnemyResult{
		Actor: actorToInfo(&enemy.Actor, model.ActorTypeEnemy, req.Enemy.Name),
	}, nil
}

// CreateCompanion 创建一个同伴角色（Companion）
// 用于创建跟随玩家的同伴，如动物伙伴、构装体或追随者等。同伴归属于指定的领导者。
// 参数:
//
//	ctx - 上下文
//	req - 创建请求，包含游戏会话ID和同伴配置（名称、描述、体型、速度、属性值、领导者ID）
//
// 返回:
//
//	*CreateCompanionResult - 包含创建的同伴角色信息
//	error - 可能返回游戏不存在或权限错误
func (e *Engine) CreateCompanion(ctx context.Context, req CreateCompanionRequest) (*CreateCompanionResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpCreateCompanion); err != nil {
		return nil, err
	}

	if req.Companion == nil {
		return nil, fmt.Errorf("companion input is required")
	}

	companion := &model.Companion{
		Actor: model.Actor{
			ID:            model.NewID(),
			Name:          req.Companion.Name,
			Description:   req.Companion.Description,
			AbilityScores: abilityScoresInputToModel(req.Companion.AbilityScores),
			Size:          model.Size(req.Companion.Size),
			CreatureType:  model.CreatureType(req.Companion.CreatureType),
			Speed:         req.Companion.Speed,
			HitPoints:     model.HitPoints{},
			Conditions:    []model.ConditionInstance{},
			Exhaustion:    0,
		},
		LeaderID:     model.ID(req.Companion.LeaderID),
		Loyalty:      req.Companion.Loyalty,
		BehaviorMode: req.Companion.BehaviorMode,
		Commands:     req.Companion.Commands,
	}

	// 设置先攻修正
	if req.Companion.InitiativeBonus != nil {
		companion.InitiativeBonus = *req.Companion.InitiativeBonus
	}

	// 设置生命值
	if req.Companion.HitPoints != nil {
		companion.HitPoints.Maximum = *req.Companion.HitPoints
		companion.HitPoints.Current = *req.Companion.HitPoints
	} else {
		// 默认使用属性修正计算HP
		conMod := rules.AbilityModifier(companion.AbilityScores.Constitution)
		companion.HitPoints.Maximum = 10 + conMod
		companion.HitPoints.Current = companion.HitPoints.Maximum
	}

	// 设置护甲等级
	if req.Companion.ArmorClass != nil {
		companion.ArmorClass = *req.Companion.ArmorClass
	} else {
		companion.ArmorClass = calculateArmorClassFromActor(&companion.Actor)
	}

	game.Companions[companion.ID] = companion

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &CreateCompanionResult{
		Actor: actorToInfo(&companion.Actor, model.ActorTypeCompanion, req.Companion.Name),
	}, nil
}

// GetActor 获取任意类型角色的基本信息
// 根据角色ID查找并返回角色的基本信息，支持PC、NPC、Enemy、Companion所有类型。
// 参数:
//
//	ctx - 上下文
//	req - 获取请求，包含游戏会话ID和要查询的角色ID
//
// 返回:
//
//	*GetActorResult - 包含角色的基本信息（ID、类型、名称、HP、AC、速度、状态效果等）
//	error - 可能返回 ErrNotFound（角色不存在）、游戏不存在或权限错误
func (e *Engine) GetActor(ctx context.Context, req GetActorRequest) (*GetActorResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpGetActor); err != nil {
		return nil, err
	}

	actor, ok := game.GetActor(req.ActorID)
	if !ok {
		return nil, ErrNotFound
	}

	var info *ActorInfo
	switch a := actor.(type) {
	case *model.PlayerCharacter:
		info = actorToInfo(&a.Actor, model.ActorTypePC, a.Name)
	case *model.NPC:
		info = actorToInfo(&a.Actor, model.ActorTypeNPC, a.Name)
	case *model.Enemy:
		info = actorToInfo(&a.Actor, model.ActorTypeEnemy, a.Name)
	case *model.Companion:
		info = actorToInfo(&a.Actor, model.ActorTypeCompanion, a.Name)
	default:
		return nil, fmt.Errorf("unknown actor type")
	}

	return &GetActorResult{Actor: info}, nil
}

// GetPC 获取玩家角色的完整数据
// 返回玩家角色的详细信息，包括职业、等级、经验值、属性值、熟练加值、种族特性等完整数据。
// 参数:
//
//	ctx - 上下文
//	req - 获取请求，包含游戏会话ID和玩家角色ID
//
// 返回:
//
//	*GetPCResult - 包含玩家角色的完整信息（职业、等级、经验、属性、特性等）
//	error - 可能返回 ErrNotFound（角色不存在）或游戏不存在错误
func (e *Engine) GetPC(ctx context.Context, req GetPCRequest) (*GetPCResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	pc, ok := game.PCs[req.PCID]
	if !ok {
		return nil, ErrNotFound
	}

	return &GetPCResult{
		PC: playerCharacterToInfo(pc),
	}, nil
}

// UpdateActor 更新角色的部分状态
// 用于更新角色的属性值、生命值、状态效果、位置、所在场景等信息。支持所有类型的角色。
// 参数:
//
//	ctx - 上下文
//	req - 更新请求，包含游戏会话ID、角色ID和更新内容（属性值、HP、状态效果、位置、场景ID、自定义字段）
//
// 返回:
//
//	error - 可能返回 ErrNotFound（角色不存在）、游戏不存在或权限错误
func (e *Engine) UpdateActor(ctx context.Context, req UpdateActorRequest) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return err
	}

	if err := e.checkPermission(game.Phase, OpUpdateActor); err != nil {
		return err
	}

	actor, ok := game.GetActor(req.ActorID)
	if !ok {
		return ErrNotFound
	}

	var baseActor *model.Actor
	switch a := actor.(type) {
	case *model.PlayerCharacter:
		baseActor = &a.Actor
	case *model.NPC:
		baseActor = &a.Actor
	case *model.Enemy:
		baseActor = &a.Actor
	case *model.Companion:
		baseActor = &a.Actor
	}

	// 应用更新
	if req.Update.AbilityScores != nil {
		baseActor.AbilityScores = abilityScoresInputToModel(*req.Update.AbilityScores)
	}
	if req.Update.HitPoints != nil {
		if req.Update.HitPoints.Current != nil {
			baseActor.HitPoints.Current = *req.Update.HitPoints.Current
		}
		if req.Update.HitPoints.TempHitPoints != nil {
			baseActor.TempHitPoints = *req.Update.HitPoints.TempHitPoints
		}
	}
	if req.Update.Conditions != nil {
		// 添加新状态
		baseActor.Conditions = append(baseActor.Conditions, req.Update.Conditions.Add...)
		// 移除指定状态
		if len(req.Update.Conditions.Remove) > 0 {
			newConditions := make([]model.ConditionInstance, 0)
			for _, c := range baseActor.Conditions {
				shouldRemove := false
				for _, rem := range req.Update.Conditions.Remove {
					if c.Type == rem {
						shouldRemove = true
						break
					}
				}
				if !shouldRemove {
					newConditions = append(newConditions, c)
				}
			}
			baseActor.Conditions = newConditions
		}
	}
	if req.Update.Position != nil {
		baseActor.Position = req.Update.Position
	}
	if req.Update.SceneID != nil {
		baseActor.SceneID = *req.Update.SceneID
	}

	if err := e.saveGame(ctx, game); err != nil {
		return err
	}

	return nil
}

// RemoveActor 从游戏中移除一个角色
// 将指定角色从游戏会话中删除。如果角色当前处于战斗状态且未被击败，则不允许移除。
// 参数:
//
//	ctx - 上下文
//	req - 移除请求，包含游戏会话ID和要移除的角色ID
//
// 返回:
//
//	error - 可能返回 ErrNotFound（角色不存在）、ErrInvalidState（战斗中未击败的角色无法移除）、游戏不存在或权限错误
func (e *Engine) RemoveActor(ctx context.Context, req RemoveActorRequest) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return err
	}

	if err := e.checkPermission(game.Phase, OpRemoveActor); err != nil {
		return err
	}

	// 检查是否在战斗中
	if game.Phase == model.PhaseCombat && game.Combat != nil {
		combatant := game.Combat.GetCombatantByActorID(req.ActorID)
		if combatant != nil && !combatant.IsDefeated {
			return ErrInvalidState
		}
	}

	if !game.RemoveActor(req.ActorID) {
		return ErrNotFound
	}

	if err := e.saveGame(ctx, game); err != nil {
		return err
	}

	return nil
}

// ListActors 列出游戏中的所有角色，支持按条件过滤
// 返回游戏会话中所有角色的基本信息，可按角色类型、所在场景、存活状态进行过滤。
// 参数:
//
//	ctx - 上下文
//	req - 列表请求，包含游戏会话ID和可选的过滤条件（角色类型、场景ID、存活状态）
//
// 返回:
//
//	*ListActorsResult - 包含符合条件的角色列表
//	error - 可能返回游戏不存在错误
func (e *Engine) ListActors(ctx context.Context, req ListActorsRequest) (*ListActorsResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	allActors := game.GetAllActors()

	if req.Filter == nil {
		result := make([]ActorInfo, len(allActors))
		for i, actor := range allActors {
			result[i] = *actorSnapshotToInfo(&actor)
		}
		return &ListActorsResult{Actors: result}, nil
	}

	result := make([]ActorInfo, 0)
	for _, actor := range allActors {
		// 按类型过滤
		if len(req.Filter.Types) > 0 {
			typeMatch := false
			for _, t := range req.Filter.Types {
				if actor.Type == t {
					typeMatch = true
					break
				}
			}
			if !typeMatch {
				continue
			}
		}

		// 按场景ID过滤
		if req.Filter.SceneID != nil && actor.SceneID != *req.Filter.SceneID {
			continue
		}

		// 按存活状态过滤
		if req.Filter.Alive != nil {
			isAlive := actor.HitPoints.Current > 0
			if isAlive != *req.Filter.Alive {
				continue
			}
		}

		result = append(result, *actorSnapshotToInfo(&actor))
	}

	return &ListActorsResult{Actors: result}, nil
}

// AddExperience 为玩家角色添加经验值
// 向指定的玩家角色添加经验值，并自动检查是否达到升级条件。如果经验值足够则自动升级。
// 参数:
//
//	ctx - 上下文
//	req - 添加经验值请求，包含游戏会话ID、玩家角色ID和要添加的经验值数量
//
// 返回:
//
//	*AddExperienceResult - 包含是否升级、原等级和新等级信息
//	error - 可能返回 ErrNotFound（角色不存在）、游戏不存在或权限错误
func (e *Engine) AddExperience(ctx context.Context, req AddExperienceRequest) (*AddExperienceResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpAddExperience); err != nil {
		return nil, err
	}

	pc, ok := game.PCs[req.PCID]
	if !ok {
		return nil, ErrNotFound
	}

	oldLevel := pc.TotalLevel
	pc.Experience += req.XP

	// 检查是否升级
	newLevel := rules.GetLevelByXP(pc.Experience)
	leveledUp := newLevel > oldLevel
	if leveledUp {
		pc.TotalLevel = newLevel
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &AddExperienceResult{
		LeveledUp: leveledUp,
		OldLevel:  oldLevel,
		NewLevel:  pc.TotalLevel,
	}, nil
}

// LevelUp 手动触发玩家角色升级
// 在经验值足够的的情况下，手动将玩家角色提升一级。会计算HP增长、更新熟练加值、检查法术位更新等。
// 参数:
//
//	ctx - 上下文
//	req - 升级请求，包含游戏会话ID、玩家角色ID和升级时选择的职业（为空则默认第一个职业）
//
// 返回:
//
//	*LevelUpResult - 包含升级详情（原等级、新等级、HP增长、新特性、熟练加值变化等）
//	error - 可能返回 ErrNotFound（角色不存在）、经验值不足、游戏不存在或权限错误
func (e *Engine) LevelUp(ctx context.Context, req LevelUpRequest) (*LevelUpResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpLevelUp); err != nil {
		return nil, err
	}

	pc, ok := game.PCs[req.PCID]
	if !ok {
		return nil, ErrNotFound
	}

	oldLevel := pc.TotalLevel
	newLevel := oldLevel + 1

	// 检查XP是否足够
	requiredXP := rules.GetXPForLevel(newLevel)
	if pc.Experience < requiredXP {
		return nil, fmt.Errorf("insufficient experience: need %d, have %d", requiredXP, pc.Experience)
	}

	// 计算HP增长
	hpGain := 0
	classToLevel := req.ClassChoice
	if classToLevel == "" && len(pc.Classes) > 0 {
		classToLevel = string(pc.Classes[0].Class)
	}

	classID, err := data.GetClassID(classToLevel)
	if err != nil {
		classID = model.ClassFighter // 默认战士
	}

	classDef := data.GetClass(classID)
	hitDiceType := 8 // 默认d8
	if classDef != nil {
		hitDiceType = classDef.HitDie
	}

	// 简化的HP计算：取平均值+CON修正
	conMod := rules.AbilityModifier(pc.AbilityScores.Constitution)
	hpGain = (hitDiceType / 2) + 1 + conMod
	if hpGain < 1 {
		hpGain = 1
	}

	pc.HitPoints.Maximum += hpGain
	pc.HitPoints.Current += hpGain
	pc.TotalLevel = newLevel

	// 更新熟练加值检查
	oldProfBonus := rules.ProficiencyBonus(oldLevel)
	newProfBonus := rules.ProficiencyBonus(newLevel)
	proficiencyIncreased := newProfBonus > oldProfBonus

	result := &LevelUpResult{
		OldLevel:             oldLevel,
		NewLevel:             newLevel,
		HPGain:               hpGain,
		SpellSlotsUpdated:    pc.Spellcasting != nil,
		ProficiencyIncreased: proficiencyIncreased,
		Message:              fmt.Sprintf("升级到等级 %d！HP +%d", newLevel, hpGain),
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return result, nil
}

// ShortRest 为指定角色执行短休
// 短休允许角色使用生命骰恢复HP。PC角色可以掷一个生命骰并加上体质修正值来恢复生命值。
// 参数:
//
//	ctx - 上下文
//	req - 短休请求，包含游戏会话ID和参与短休的角色ID列表
//
// 返回:
//
//	*RestResult - 包含各角色的休息结果（恢复的HP、使用的生命骰等）
//	error - 可能返回游戏不存在、权限错误或某个角色处理失败错误
func (e *Engine) ShortRest(ctx context.Context, req ShortRestRequest) (*RestResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpShortRest); err != nil {
		return nil, err
	}

	result := &RestResult{
		ActorResults: make([]ActorRestResult, 0),
		Message:      "短休完成",
	}

	for _, actorID := range req.ActorIDs {
		// 获取该角色要使用的生命骰数量
		hitDiceCount := 1 // 默认1个
		if req.HitDiceCount != nil {
			if count, exists := req.HitDiceCount[actorID]; exists {
				hitDiceCount = count
			}
		}

		actorResult, err := e.processShortRest(game, actorID, hitDiceCount)
		if err != nil {
			return nil, fmt.Errorf("short rest failed for actor %s: %w", actorID, err)
		}
		result.ActorResults = append(result.ActorResults, actorResult)
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return result, nil
}

// StartLongRest 开始长休过程
// 创建长休状态并切换到休息阶段。长休需要8小时才能完成，完成后需调用 EndLongRest 应用恢复效果。
// 参数:
//
//	ctx - 上下文
//	req - 开始长休请求，包含游戏会话ID和参与长休的角色ID列表
//
// 返回:
//
//	*RestResult - 包含长休开始的消息
//	error - 可能返回游戏不存在或权限错误
func (e *Engine) StartLongRest(ctx context.Context, req StartLongRestRequest) (*RestResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpStartLongRest); err != nil {
		return nil, err
	}

	// 创建长休状态
	restState := model.NewLongRest(req.ActorIDs)
	restState.Start()
	game.ActiveRest = restState

	// 切换到休息阶段
	game.Phase = model.PhaseRest

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &RestResult{
		Message: "长休开始，需要8小时才能完成",
	}, nil
}

// EndLongRest 结束长休并应用恢复效果
// 完成长休过程，为所有参与角色恢复全部HP、恢复法术位、恢复部分生命骰、减少力竭等级，并移除中毒、恐惧、魅惑等状态效果。
// 参数:
//
//	ctx - 上下文
//	req - 结束长休请求，包含游戏会话ID
//
// 返回:
//
//	*RestResult - 包含各角色的恢复结果（恢复的HP、恢复的法术位、移除的状态效果等）
//	error - 可能返回 ErrInvalidState（没有活跃的长休）、游戏不存在或权限错误
func (e *Engine) EndLongRest(ctx context.Context, req EndLongRestRequest) (*RestResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpEndLongRest); err != nil {
		return nil, err
	}

	if game.ActiveRest == nil || game.ActiveRest.Type != model.RestTypeLong {
		return nil, ErrInvalidState
	}

	result := &RestResult{
		ActorResults: make([]ActorRestResult, 0),
	}

	// 对每个参与者应用长休效果
	for _, actorID := range game.ActiveRest.ParticipantIDs {
		actorResult, err := e.processLongRestRecovery(game, actorID)
		if err != nil {
			return nil, fmt.Errorf("long rest recovery failed for actor %s: %w", actorID, err)
		}
		result.ActorResults = append(result.ActorResults, actorResult)
	}

	// 完成长休
	game.ActiveRest.Complete()
	game.Phase = model.PhaseExploration
	game.ActiveRest = nil

	result.Message = "长休完成，队伍完全恢复"

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return result, nil
}

// processShortRest 处理单个角色的短休
// PHB 第8章: 短休时可以使用一粒或多粒生命骰(最多等于角色等级)
func (e *Engine) processShortRest(game *model.GameState, actorID model.ID, hitDiceCount int) (ActorRestResult, error) {
	result := ActorRestResult{
		ActorID: actorID,
	}

	// 查找角色
	actor, ok := game.GetActor(actorID)
	if !ok {
		return result, ErrNotFound
	}

	var baseActor *model.Actor
	var pc *model.PlayerCharacter
	switch a := actor.(type) {
	case *model.PlayerCharacter:
		baseActor = &a.Actor
		pc = a
	case *model.Companion:
		baseActor = &a.Actor
	case *model.NPC:
		baseActor = &a.Actor
	case *model.Enemy:
		baseActor = &a.Actor
	}

	// 短休可以掷生命骰恢复HP(仅PC)
	if pc != nil && len(pc.HitDice) > 0 {
		// 如果未指定,默认使用1个生命骰
		if hitDiceCount <= 0 {
			hitDiceCount = 1
		}

		// 计算可用生命骰总数
		totalAvailable := 0
		for i := range pc.HitDice {
			totalAvailable += pc.HitDice[i].Total - pc.HitDice[i].Used
		}

		// 实际使用的生命骰不能超过可用数量
		if hitDiceCount > totalAvailable {
			hitDiceCount = totalAvailable
		}

		if hitDiceCount > 0 {
			// 使用rules层的UseHitDice函数计算恢复效果
			conMod := rules.AbilityModifier(pc.AbilityScores.Constitution)
			restResult, err := rules.UseHitDice(
				baseActor.HitPoints.Current,
				baseActor.HitPoints.Maximum,
				pc.HitDice,
				conMod,
				hitDiceCount,
			)

			if err != nil {
				return result, fmt.Errorf("使用生命骰失败: %w", err)
			}

			// 应用恢复的HP
			if restResult.HPRestored > 0 {
				baseActor.HitPoints.Current += restResult.HPRestored
				if baseActor.HitPoints.Current > baseActor.HitPoints.Maximum {
					baseActor.HitPoints.Current = baseActor.HitPoints.Maximum
				}
			}

			result.HPRecovered = restResult.HPRestored
			result.HitDiceUsed = restResult.HitDiceUsed
		}

		// 邪术师特性: 短休恢复Pact Magic法术位
		// D&D 5e规则: 邪术师完成短休或长休后恢复所有已消耗的法术位
		isWarlock := false
		for _, cl := range pc.Classes {
			if cl.Class == model.ClassWarlock {
				isWarlock = true
				break
			}
		}

		if isWarlock && pc.Spellcasting != nil && pc.FeatureHooks != nil {
			// 恢复所有法术位
			pc.Spellcasting.Slots.RestoreAll()
			result.SpellSlotsRestored = true

			// 调用邪术师钩子
			if warlockHooks, ok := pc.FeatureHooks[model.ClassWarlock].(*model.WarlockFeatureHooks); ok {
				warlockHooks.OnShortRest()
			}
		}
	}

	return result, nil
}

// processLongRestRecovery 处理单个角色的长休恢复
func (e *Engine) processLongRestRecovery(game *model.GameState, actorID model.ID) (ActorRestResult, error) {
	result := ActorRestResult{
		ActorID: actorID,
	}

	actor, ok := game.GetActor(actorID)
	if !ok {
		return result, ErrNotFound
	}

	var baseActor *model.Actor
	var pc *model.PlayerCharacter
	switch a := actor.(type) {
	case *model.PlayerCharacter:
		baseActor = &a.Actor
		pc = a
	case *model.Companion:
		baseActor = &a.Actor
	case *model.NPC:
		baseActor = &a.Actor
	case *model.Enemy:
		baseActor = &a.Actor
	}

	// 恢复所有HP
	hpRecovered := baseActor.HitPoints.Maximum - baseActor.HitPoints.Current
	baseActor.HitPoints.Current = baseActor.HitPoints.Maximum
	baseActor.TempHitPoints = 0
	result.HPRecovered = hpRecovered

	// 恢复法术位(PC)
	// D&D 5e规则:
	// - 标准施法者(Spellcasting): 长休恢复所有法术位
	// - 邪术师(Pact Magic): 短休即可恢复所有法术位,长休也恢复
	if pc != nil && pc.Spellcasting != nil {
		// 检查是否是邪术师
		isWarlock := false
		for _, cl := range pc.Classes {
			if cl.Class == model.ClassWarlock {
				isWarlock = true
				break
			}
		}

		// 长休恢复所有法术位(对邪术师和标准施法者都适用)
		pc.Spellcasting.Slots.RestoreAll()
		pc.Spellcasting.ConcentrationSpell = ""
		result.SpellSlotsRestored = true

		// 邪术师特性:调用OnShortRest钩子(短休也能恢复)
		// 注意:长休时也会调用,但效果相同(都是恢复所有位)
		if isWarlock && pc.FeatureHooks != nil {
			if warlockHooks, ok := pc.FeatureHooks[model.ClassWarlock].(*model.WarlockFeatureHooks); ok {
				warlockHooks.OnLongRest()
			}
		}
	}

	// 恢复生命骰
	if pc != nil {
		recoveryAmount := rules.CalculateHitDiceRecovery(pc.TotalLevel)
		recovered := 0
		for i := range pc.HitDice {
			available := pc.HitDice[i].Total - pc.HitDice[i].Used
			canRecover := recoveryAmount - recovered
			if canRecover > available {
				canRecover = available
			}
			pc.HitDice[i].Used -= canRecover
			recovered += canRecover
		}
		result.HitDiceUsed = -recovered // 负数表示恢复
	}

	// 减少力竭
	if baseActor.Exhaustion > 0 {
		baseActor.Exhaustion = rules.CalculateExhaustionReduction(baseActor.Exhaustion)
		result.ExhaustionReduced = true
	}

	// 移除某些状态效果（长休结束时移除）
	removableConditions := []model.ConditionType{
		model.ConditionPoisoned,
		model.ConditionFrightened,
		model.ConditionCharmed,
	}
	for _, condType := range removableConditions {
		for i := len(baseActor.Conditions) - 1; i >= 0; i-- {
			if baseActor.Conditions[i].Type == condType {
				result.ConditionsRemoved = append(result.ConditionsRemoved, condType)
				baseActor.Conditions = append(baseActor.Conditions[:i], baseActor.Conditions[i+1:]...)
			}
		}
	}

	return result, nil
}

// calculateArmorClass 计算PC的护甲等级
func calculateArmorClass(pc *model.PlayerCharacter) int {
	// 简化实现：默认10+DEX修正
	dexMod := rules.AbilityModifier(pc.AbilityScores.Dexterity)
	return 10 + dexMod
}

// calculateArmorClassFromActor 从Actor计算AC
func calculateArmorClassFromActor(actor *model.Actor) int {
	dexMod := rules.AbilityModifier(actor.AbilityScores.Dexterity)
	return 10 + dexMod
}

// calculateMaxHP 计算最大HP
func calculateMaxHP(pc *model.PlayerCharacter) int {
	if len(pc.Classes) == 0 {
		return 10 // 默认值
	}

	hp := 0
	conMod := rules.AbilityModifier(pc.AbilityScores.Constitution)

	for i, cl := range pc.Classes {
		classDef := data.GetClass(cl.Class)
		hitDiceType := 8 // 默认d8
		if classDef != nil {
			hitDiceType = classDef.HitDie
		}

		// 第一级取最大值，之后取平均
		if i == 0 {
			hp += hitDiceType + conMod
		} else {
			hp += (hitDiceType/2 + 1) + conMod
		}
	}
	return hp
}

// abilityScoresInputToModel 将 AbilityScoresInput 转换为 model.AbilityScores
func abilityScoresInputToModel(input AbilityScoresInput) model.AbilityScores {
	return model.AbilityScores{
		Strength:     input.Strength,
		Dexterity:    input.Dexterity,
		Constitution: input.Constitution,
		Intelligence: input.Intelligence,
		Wisdom:       input.Wisdom,
		Charisma:     input.Charisma,
	}
}

// actorToInfo 将 Actor 转换为 ActorInfo
func actorToInfo(actor *model.Actor, actorType model.ActorType, name string) *ActorInfo {
	conditions := make([]string, len(actor.Conditions))
	for i, c := range actor.Conditions {
		conditions[i] = string(c.Type)
	}
	info := &ActorInfo{
		ID:         actor.ID,
		Type:       actorType,
		Name:       name,
		HitPoints:  actor.HitPoints,
		TempHP:     actor.TempHitPoints,
		ArmorClass: actor.ArmorClass,
		Speed:      actor.Speed,
		Conditions: conditions,
		Exhaustion: actor.Exhaustion,
		SceneID:    actor.SceneID,
	}
	if actor.Position != nil {
		pos := *actor.Position
		info.Position = &pos
	}
	return info
}

// actorSnapshotToInfo 将 ActorSnapshot 转换为 ActorInfo
func actorSnapshotToInfo(snapshot *model.ActorSnapshot) *ActorInfo {
	return &ActorInfo{
		ID:         snapshot.ID,
		Type:       snapshot.Type,
		Name:       snapshot.Name,
		HitPoints:  snapshot.HitPoints,
		ArmorClass: snapshot.ArmorClass,
		Conditions: snapshot.Conditions,
		SceneID:    snapshot.SceneID,
	}
}

// playerCharacterToInfo 将 PlayerCharacter 转换为 PlayerCharacterInfo
func playerCharacterToInfo(pc *model.PlayerCharacter) *PlayerCharacterInfo {
	classes := make([]ClassInfo, len(pc.Classes))
	for i, cl := range pc.Classes {
		classes[i] = ClassInfo{
			Class:      cl.Class,
			ClassLevel: cl.Level,
			Features:   cl.Features,
		}
	}
	info := &PlayerCharacterInfo{
		ID:         pc.ID,
		Name:       pc.Actor.Name,
		Race:       pc.Race.Name,
		Classes:    classes,
		TotalLevel: pc.TotalLevel,
		Experience: pc.Experience,
		AbilityScores: AbilityScoresInput{
			Strength:     pc.AbilityScores.Strength,
			Dexterity:    pc.AbilityScores.Dexterity,
			Constitution: pc.AbilityScores.Constitution,
			Intelligence: pc.AbilityScores.Intelligence,
			Wisdom:       pc.AbilityScores.Wisdom,
			Charisma:     pc.AbilityScores.Charisma,
		},
		HitPoints:        pc.HitPoints,
		ArmorClass:       pc.ArmorClass,
		Speed:            pc.Speed,
		ProficiencyBonus: rules.ProficiencyBonus(pc.TotalLevel),
		Features:         pc.Features,
		RacialTraits:     pc.RacialTraits,
	}
	if pc.Personality != nil && pc.Personality.Background != "" {
		info.Background = pc.Personality.Background
	} else if pc.BackgroundID != "" {
		// 通过注册表查询中文名
		if bg, ok := data.GlobalRegistry.GetBackground(string(pc.BackgroundID)); ok {
			info.Background = bg.Name
		} else {
			info.Background = string(pc.BackgroundID)
		}
	}
	return info
}

// getClassFeatures 获取指定职业和等级的特性列表
func getClassFeatures(classID model.ClassID, level int) []string {
	var features []string

	// 战士特性
	if classID == model.ClassFighter {
		for lvl := 1; lvl <= level; lvl++ {
			if feats, ok := data.FighterFeaturesByLevel[lvl]; ok {
				features = append(features, feats...)
			}
		}
	}

	// TODO: 添加其他职业的特性

	return features
}
