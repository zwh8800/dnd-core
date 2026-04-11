package engine

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/zwh8800/dnd-core/pkg/data"
	"github.com/zwh8800/dnd-core/pkg/model"
)

// ============================================================================
// 分页相关结构体
// ============================================================================

// PaginationRequest 分页请求参数
type PaginationRequest struct {
	Page     int `json:"page"`      // 页码，从 1 开始
	PageSize int `json:"page_size"` // 每页数量，默认 20，最大 100
}

// PaginationInfo 分页元信息
type PaginationInfo struct {
	Page       int  `json:"page"`        // 当前页码
	PageSize   int  `json:"page_size"`   // 每页数量
	TotalCount int  `json:"total_count"` // 总记录数
	TotalPages int  `json:"total_pages"` // 总页数
	HasNext    bool `json:"has_next"`    // 是否有下一页
	HasPrev    bool `json:"has_prev"`    // 是否有上一页
}

// applyDefaults 应用分页默认值并验证
func (p *PaginationRequest) applyDefaults() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 20
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
}

// calculatePagination 计算分页信息
func calculatePagination(totalCount int, req *PaginationRequest) PaginationInfo {
	req.applyDefaults()

	totalPages := (totalCount + req.PageSize - 1) / req.PageSize
	if totalPages == 0 {
		totalPages = 1
	}

	return PaginationInfo{
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalCount: totalCount,
		TotalPages: totalPages,
		HasNext:    req.Page < totalPages,
		HasPrev:    req.Page > 1,
	}
}

// paginateSlice 对切片进行分页
func paginateSlice[T any](items []T, req *PaginationRequest) []T {
	if req == nil {
		req = &PaginationRequest{}
	}
	req.applyDefaults()

	start := (req.Page - 1) * req.PageSize
	if start >= len(items) {
		return []T{}
	}

	end := start + req.PageSize
	if end > len(items) {
		end = len(items)
	}

	return items[start:end]
}

// ============================================================================
// 权限常量（在 phase.go 中添加）
// ============================================================================

const (
	OpListRaces        Operation = "list_races"
	OpGetRace          Operation = "get_race"
	OpListClasses      Operation = "list_classes"
	OpGetClass         Operation = "get_class"
	OpListBackgrounds  Operation = "list_backgrounds"
	OpGetBackground    Operation = "get_background"
	OpListFeatsData    Operation = "list_feats_data"
	OpGetFeatData      Operation = "get_feat_data"
	OpListMonsters     Operation = "list_monsters"
	OpGetMonster       Operation = "get_monster"
	OpListSpells       Operation = "list_spells"
	OpGetSpell         Operation = "get_spell"
	OpListWeapons      Operation = "list_weapons"
	OpGetWeapon        Operation = "get_weapon"
	OpListArmors       Operation = "list_armors"
	OpGetArmor         Operation = "get_armor"
	OpListMagicItems   Operation = "list_magic_items"
	OpGetMagicItem     Operation = "get_magic_item"
	OpListGears        Operation = "list_gears"
	OpGetGear          Operation = "get_gear"
	OpListTools        Operation = "list_tools"
	OpGetTool          Operation = "get_tool"
	OpListRecipes      Operation = "list_recipes"
	OpGetRecipe        Operation = "get_recipe"
	OpListLifestyles   Operation = "list_lifestyles"
	OpGetLifestyleData Operation = "get_lifestyle_data"
	OpListMounts       Operation = "list_mounts"
	OpGetMount         Operation = "get_mount"
	OpListPoisons      Operation = "list_poisons"
	OpGetPoison        Operation = "get_poison"
	OpListTraps        Operation = "list_traps"
	OpGetTrap          Operation = "get_trap"
)

// ============================================================================
// Info 结构体定义
// ============================================================================

// RaceInfo 种族信息
type RaceInfo struct {
	Name           string         `json:"name"`            // 种族名称
	Subraces       []string       `json:"subraces"`        // 子种族列表
	AbilityBonuses map[string]int `json:"ability_bonuses"` // 属性加值
	Speed          int            `json:"speed"`           // 速度
	Size           model.Size     `json:"size"`            // 体型
	Languages      []string       `json:"languages"`       // 语言
	Traits         []string       `json:"traits"`          // 特性
	Description    string         `json:"description"`     // 描述
}

// ClassInfo 职业信息（已在 actor.go 中定义，这里复用）
// 注意：actor.go 中的 ClassInfo 用于角色职业信息，这里需要完整的职业定义

// ClassDefinitionInfo 完整职业定义信息
type ClassDefinitionInfo struct {
	ID                  model.ClassID `json:"id"`                   // 职业ID
	Name                string        `json:"name"`                 // 职业名称
	HitDie              int           `json:"hit_die"`              // 生命骰
	PrimaryAbilities    []string      `json:"primary_abilities"`    // 主要属性
	SavingThrows        []string      `json:"saving_throws"`        // 豁免熟练
	SkillChoices        []string      `json:"skill_choices"`        // 可选技能
	NumberOfSkills      int           `json:"number_of_skills"`     // 可选技能数量
	ArmorProficiencies  []string      `json:"armor_proficiencies"`  // 护甲熟练
	WeaponProficiencies []string      `json:"weapon_proficiencies"` // 武器熟练
	ToolProficiencies   []string      `json:"tool_proficiencies"`   // 工具熟练
	SpellcastingAbility string        `json:"spellcasting_ability"` // 施法属性
	CasterType          string        `json:"caster_type"`          // 施法者类型
	Description         string        `json:"description"`          // 描述
}

// BackgroundInfo 背景信息
type BackgroundInfo struct {
	ID                string   `json:"id"`                  // 背景ID
	Name              string   `json:"name"`                // 背景名称
	SkillProficienies []string `json:"skill_proficiencies"` // 技能熟练
	ToolProficiencies []string `json:"tool_proficiencies"`  // 工具熟练
	Languages         []string `json:"languages"`           // 语言
	Equipment         []string `json:"equipment"`           // 起始装备
	FeatureName       string   `json:"feature_name"`        // 特性名称
	FeatureDesc       string   `json:"feature_desc"`        // 特性描述
	Description       string   `json:"description"`         // 描述
}

// MonsterInfo 怪物信息
type MonsterInfo struct {
	ID              string `json:"id"`               // 怪物ID
	Name            string `json:"name"`             // 怪物名称
	ChallengeRating string `json:"challenge_rating"` // 挑战等级
	ArmorClass      int    `json:"armor_class"`      // 护甲等级
	HitPoints       int    `json:"hit_points"`       // 生命值
	Speed           int    `json:"speed"`            // 速度
	Size            string `json:"size"`             // 体型
	Type            string `json:"type"`             // 类型
	Alignment       string `json:"alignment"`        // 阵营
}

// SpellInfo 法术信息
type SpellInfo struct {
	ID          string   `json:"id"`           // 法术ID
	Name        string   `json:"name"`         // 法术名称
	Level       int      `json:"level"`        // 法术等级
	School      string   `json:"school"`       // 法术学派
	CastingTime string   `json:"casting_time"` // 施法时间
	Range       string   `json:"range"`        // 范围
	Duration    string   `json:"duration"`     // 持续时间
	Components  []string `json:"components"`   // 组件
	Classes     []string `json:"classes"`      // 可用职业
	Description string   `json:"description"`  // 描述
}

// ItemInfo 物品信息（通用）
type ItemInfo struct {
	ID          string `json:"id"`          // 物品ID
	Name        string `json:"name"`        // 物品名称
	Type        string `json:"type"`        // 物品类型
	Description string `json:"description"` // 描述
	Cost        int    `json:"cost"`        // 价格（CP）
	Weight      int    `json:"weight"`      // 重量（磅）
}

// WeaponInfo 武器信息
type WeaponInfo struct {
	ItemInfo
	Damage     string `json:"damage"`      // 伤害
	DamageType string `json:"damage_type"` // 伤害类型
	Property   string `json:"property"`    // 属性
	Range      string `json:"range"`       // 射程
}

// ArmorInfo 护甲信息
type ArmorInfo struct {
	ItemInfo
	ArmorClass          string `json:"armor_class"`          // 护甲等级
	StrRequired         int    `json:"str_required"`         // 所需力量
	StealthDisadvantage bool   `json:"stealth_disadvantage"` // 隐匿劣势
}

// LifestyleInfo 生活方式信息
type LifestyleInfo struct {
	Tier        model.LifestyleTier `json:"tier"`         // 生活方式等级
	Name        string              `json:"name"`         // 名称
	DailyCost   int                 `json:"daily_cost"`   // 每日花费（GP）
	MonthlyCost int                 `json:"monthly_cost"` // 每月花费（GP）
	Description string              `json:"description"`  // 描述
}

// MountInfo 坐骑信息
type MountInfo struct {
	ID               string `json:"id"`                // 坐骑ID
	Name             string `json:"name"`              // 坐骑名称
	Type             string `json:"type"`              // 类型
	Speed            int    `json:"speed"`             // 速度
	CarryingCapacity int    `json:"carrying_capacity"` // 载重能力
	Cost             int    `json:"cost"`              // 价格（GP）
	Description      string `json:"description"`       // 描述
}

// PoisonInfo 毒药信息
type PoisonInfo struct {
	ID          string `json:"id"`          // 毒药ID
	Name        string `json:"name"`        // 毒药名称
	Type        string `json:"type"`        // 类型
	DC          int    `json:"dc"`          // 豁免DC
	Cost        int    `json:"cost"`        // 价格（GP）
	Description string `json:"description"` // 描述
}

// TrapInfo 陷阱信息
type TrapInfo struct {
	ID          string `json:"id"`          // 陷阱ID
	Name        string `json:"name"`        // 陷阱名称
	Type        string `json:"type"`        // 类型
	DC          int    `json:"dc"`          // 侦测/解除DC
	Damage      string `json:"damage"`      // 伤害
	Description string `json:"description"` // 描述
}

// ============================================================================
// 请求/结果结构体定义
// ============================================================================

// ListRacesRequest 列出种族请求
type ListRacesRequest struct {
	Pagination *PaginationRequest `json:"pagination"` // 分页参数（可选）
}

// ListRacesResult 列出种族结果
type ListRacesResult struct {
	Races      []RaceInfo     `json:"races"`      // 种族列表
	Pagination PaginationInfo `json:"pagination"` // 分页信息
}

// GetRaceRequest 获取种族请求
type GetRaceRequest struct {
	Name string `json:"name"` // 种族名称（必填）
}

// GetRaceResult 获取种族结果
type GetRaceResult struct {
	Race RaceInfo `json:"race"` // 种族信息
}

// ListClassesRequest 列出职业请求
type ListClassesRequest struct {
	Pagination *PaginationRequest `json:"pagination"` // 分页参数（可选）
}

// ListClassesResult 列出职业结果
type ListClassesResult struct {
	Classes    []ClassDefinitionInfo `json:"classes"`    // 职业列表
	Pagination PaginationInfo        `json:"pagination"` // 分页信息
}

// GetClassRequest 获取职业请求
type GetClassRequest struct {
	ID model.ClassID `json:"id"` // 职业ID（必填）
}

// GetClassResult 获取职业结果
type GetClassResult struct {
	Class ClassDefinitionInfo `json:"class"` // 职业信息
}

// ListBackgroundsRequest 列出背景请求
type ListBackgroundsRequest struct {
	Pagination *PaginationRequest `json:"pagination"` // 分页参数（可选）
}

// ListBackgroundsResult 列出背景结果
type ListBackgroundsResult struct {
	Backgrounds []BackgroundInfo `json:"backgrounds"` // 背景列表
	Pagination  PaginationInfo   `json:"pagination"`  // 分页信息
}

// GetBackgroundRequest 获取背景请求
type GetBackgroundRequest struct {
	ID string `json:"id"` // 背景ID（必填）
}

// GetBackgroundResult 获取背景结果
type GetBackgroundResult struct {
	Background BackgroundInfo `json:"background"` // 背景信息
}

// ListFeatsDataRequest 列出专长请求（用于数据查询）
type ListFeatsDataRequest struct {
	Pagination *PaginationRequest `json:"pagination"` // 分页参数（可选）
}

// ListFeatsDataResult 列出专长结果
type ListFeatsDataResult struct {
	Feats      []FeatInfo     `json:"feats"`      // 专长列表
	Pagination PaginationInfo `json:"pagination"` // 分页信息
}

// GetFeatDataRequest 获取专长请求（用于数据查询，与 feat.go 中的不同）
type GetFeatDataRequest struct {
	ID string `json:"id"` // 专长ID（必填）
}

// GetFeatDataResult 获取专长结果
type GetFeatDataResult struct {
	Feat FeatInfo `json:"feat"` // 专长信息
}

// ListMonstersRequest 列出怪物请求
type ListMonstersRequest struct {
	Pagination *PaginationRequest `json:"pagination"` // 分页参数（可选）
}

// ListMonstersResult 列出怪物结果
type ListMonstersResult struct {
	Monsters   []MonsterInfo  `json:"monsters"`   // 怪物列表
	Pagination PaginationInfo `json:"pagination"` // 分页信息
}

// GetMonsterRequest 获取怪物请求
type GetMonsterRequest struct {
	ID string `json:"id"` // 怪物ID（必填）
}

// GetMonsterResult 获取怪物结果
type GetMonsterResult struct {
	Monster MonsterInfo `json:"monster"` // 怪物信息
}

// ListSpellsRequest 列出法术请求
type ListSpellsRequest struct {
	Pagination *PaginationRequest `json:"pagination"` // 分页参数（可选）
}

// ListSpellsResult 列出法术结果
type ListSpellsResult struct {
	Spells     []SpellInfo    `json:"spells"`     // 法术列表
	Pagination PaginationInfo `json:"pagination"` // 分页信息
}

// GetSpellRequest 获取法术请求
type GetSpellRequest struct {
	ID string `json:"id"` // 法术ID（必填）
}

// GetSpellResult 获取法术结果
type GetSpellResult struct {
	Spell SpellInfo `json:"spell"` // 法术信息
}

// ListWeaponsRequest 列出武器请求
type ListWeaponsRequest struct {
	Pagination *PaginationRequest `json:"pagination"` // 分页参数（可选）
}

// ListWeaponsResult 列出武器结果
type ListWeaponsResult struct {
	Weapons    []WeaponInfo   `json:"weapons"`    // 武器列表
	Pagination PaginationInfo `json:"pagination"` // 分页信息
}

// GetWeaponRequest 获取武器请求
type GetWeaponRequest struct {
	ID string `json:"id"` // 武器ID（必填）
}

// GetWeaponResult 获取武器结果
type GetWeaponResult struct {
	Weapon WeaponInfo `json:"weapon"` // 武器信息
}

// ListArmorsRequest 列出护甲请求
type ListArmorsRequest struct {
	Pagination *PaginationRequest `json:"pagination"` // 分页参数（可选）
}

// ListArmorsResult 列出护甲结果
type ListArmorsResult struct {
	Armors     []ArmorInfo    `json:"armors"`     // 护甲列表
	Pagination PaginationInfo `json:"pagination"` // 分页信息
}

// GetArmorRequest 获取护甲请求
type GetArmorRequest struct {
	ID string `json:"id"` // 护甲ID（必填）
}

// GetArmorResult 获取护甲结果
type GetArmorResult struct {
	Armor ArmorInfo `json:"armor"` // 护甲信息
}

// ListMagicItemsRequest 列出魔法物品请求
type ListMagicItemsRequest struct {
	Pagination *PaginationRequest `json:"pagination"` // 分页参数（可选）
}

// ListMagicItemsResult 列出魔法物品结果
type ListMagicItemsResult struct {
	MagicItems []ItemInfo     `json:"magic_items"` // 魔法物品列表
	Pagination PaginationInfo `json:"pagination"`  // 分页信息
}

// GetMagicItemRequest 获取魔法物品请求
type GetMagicItemRequest struct {
	ID string `json:"id"` // 魔法物品ID（必填）
}

// GetMagicItemResult 获取魔法物品结果
type GetMagicItemResult struct {
	MagicItem ItemInfo `json:"magic_item"` // 魔法物品信息
}

// ListGearsRequest 列出冒险装备请求
type ListGearsRequest struct {
	Pagination *PaginationRequest `json:"pagination"` // 分页参数（可选）
}

// ListGearsResult 列出冒险装备结果
type ListGearsResult struct {
	Gears      []ItemInfo     `json:"gears"`      // 冒险装备列表
	Pagination PaginationInfo `json:"pagination"` // 分页信息
}

// GetGearRequest 获取冒险装备请求
type GetGearRequest struct {
	ID string `json:"id"` // 冒险装备ID（必填）
}

// GetGearResult 获取冒险装备结果
type GetGearResult struct {
	Gear ItemInfo `json:"gear"` // 冒险装备信息
}

// ListToolsRequest 列出工具请求
type ListToolsRequest struct {
	Pagination *PaginationRequest `json:"pagination"` // 分页参数（可选）
}

// ListToolsResult 列出工具结果
type ListToolsResult struct {
	Tools      []ItemInfo     `json:"tools"`      // 工具列表
	Pagination PaginationInfo `json:"pagination"` // 分页信息
}

// GetToolRequest 获取工具请求
type GetToolRequest struct {
	ID string `json:"id"` // 工具ID（必填）
}

// GetToolResult 获取工具结果
type GetToolResult struct {
	Tool ItemInfo `json:"tool"` // 工具信息
}

// ListRecipesRequest 列出配方请求
type ListRecipesRequest struct {
	Pagination *PaginationRequest `json:"pagination"` // 分页参数（可选）
}

// ListRecipesResult 列出配方结果
type ListRecipesResult struct {
	Recipes    []CraftingRecipeInfo `json:"recipes"`    // 配方列表
	Pagination PaginationInfo       `json:"pagination"` // 分页信息
}

// GetRecipeRequest 获取配方请求
type GetRecipeRequest struct {
	ID string `json:"id"` // 配方ID（必填）
}

// GetRecipeResult 获取配方结果
type GetRecipeResult struct {
	Recipe CraftingRecipeInfo `json:"recipe"` // 配方信息
}

// ListLifestylesRequest 列出生活方式请求
type ListLifestylesRequest struct {
	Pagination *PaginationRequest `json:"pagination"` // 分页参数（可选）
}

// ListLifestylesResult 列出生活方式结果
type ListLifestylesResult struct {
	Lifestyles []LifestyleInfo `json:"lifestyles"` // 生活方式列表
	Pagination PaginationInfo  `json:"pagination"` // 分页信息
}

// GetLifestyleDataRequest 获取生活方式数据请求
type GetLifestyleDataRequest struct {
	Tier model.LifestyleTier `json:"tier"` // 生活方式等级（必填）
}

// GetLifestyleDataResult 获取生活方式数据结果
type GetLifestyleDataResult struct {
	Lifestyle LifestyleInfo `json:"lifestyle"` // 生活方式信息
}

// ListMountsRequest 列出坐骑请求
type ListMountsRequest struct {
	Pagination *PaginationRequest `json:"pagination"` // 分页参数（可选）
}

// ListMountsResult 列出坐骑结果
type ListMountsResult struct {
	Mounts     []MountInfo    `json:"mounts"`     // 坐骑列表
	Pagination PaginationInfo `json:"pagination"` // 分页信息
}

// GetMountRequest 获取坐骑请求
type GetMountRequest struct {
	ID string `json:"id"` // 坐骑ID（必填）
}

// GetMountResult 获取坐骑结果
type GetMountResult struct {
	Mount MountInfo `json:"mount"` // 坐骑信息
}

// ListPoisonsRequest 列出毒药请求
type ListPoisonsRequest struct {
	Pagination *PaginationRequest `json:"pagination"` // 分页参数（可选）
}

// ListPoisonsResult 列出毒药结果
type ListPoisonsResult struct {
	Poisons    []PoisonInfo   `json:"poisons"`    // 毒药列表
	Pagination PaginationInfo `json:"pagination"` // 分页信息
}

// GetPoisonRequest 获取毒药请求
type GetPoisonRequest struct {
	ID string `json:"id"` // 毒药ID（必填）
}

// GetPoisonResult 获取毒药结果
type GetPoisonResult struct {
	Poison PoisonInfo `json:"poison"` // 毒药信息
}

// ListTrapsRequest 列出陷阱请求
type ListTrapsRequest struct {
	Pagination *PaginationRequest `json:"pagination"` // 分页参数（可选）
}

// ListTrapsResult 列出陷阱结果
type ListTrapsResult struct {
	Traps      []TrapInfo     `json:"traps"`      // 陷阱列表
	Pagination PaginationInfo `json:"pagination"` // 分页信息
}

// GetTrapRequest 获取陷阱请求
type GetTrapRequest struct {
	ID string `json:"id"` // 陷阱ID（必填）
}

// GetTrapResult 获取陷阱结果
type GetTrapResult struct {
	Trap TrapInfo `json:"trap"` // 陷阱信息
}

// ============================================================================
// 辅助转换函数
// ============================================================================

// raceToInfo 将 RaceDefinition 转换为 RaceInfo
func raceToInfo(race *data.RaceDefinition) RaceInfo {
	abilityBonuses := make(map[string]int)
	for ability, bonus := range race.AbilityBonuses {
		abilityBonuses[string(ability)] = bonus
	}

	return RaceInfo{
		Name:           race.Name,
		Subraces:       race.Subraces,
		AbilityBonuses: abilityBonuses,
		Speed:          race.Speed,
		Size:           race.Size,
		Languages:      race.Languages,
		Traits:         race.Traits,
		Description:    race.Description,
	}
}

// classToInfo 将 ClassDefinition 转换为 ClassDefinitionInfo
func classToInfo(class *data.ClassDefinition) ClassDefinitionInfo {
	primaryAbilities := make([]string, len(class.PrimaryAbilities))
	for i, a := range class.PrimaryAbilities {
		primaryAbilities[i] = string(a)
	}

	savingThrows := make([]string, len(class.SavingThrows))
	for i, a := range class.SavingThrows {
		savingThrows[i] = string(a)
	}

	skillChoices := make([]string, len(class.SkillChoices))
	for i, s := range class.SkillChoices {
		skillChoices[i] = string(s)
	}

	armorProficiencies := make([]string, len(class.ArmorProficiencies))
	for i, a := range class.ArmorProficiencies {
		armorProficiencies[i] = string(a)
	}

	return ClassDefinitionInfo{
		ID:                  class.ID,
		Name:                class.Name,
		HitDie:              class.HitDie,
		PrimaryAbilities:    primaryAbilities,
		SavingThrows:        savingThrows,
		SkillChoices:        skillChoices,
		NumberOfSkills:      class.NumberOfSkills,
		ArmorProficiencies:  armorProficiencies,
		WeaponProficiencies: class.WeaponProficiencies,
		ToolProficiencies:   class.ToolProficiencies,
		SpellcastingAbility: string(class.SpellcastingAbility),
		CasterType:          string(class.CasterType),
		Description:         class.Description,
	}
}

// backgroundToInfo 将 BackgroundDefinition 转换为 BackgroundInfo
func backgroundToInfo(bg *model.BackgroundDefinition) BackgroundInfo {
	// 将 Skill 类型转换为字符串
	skillProfs := make([]string, len(bg.SkillProficiencies))
	for i, skill := range bg.SkillProficiencies {
		skillProfs[i] = string(skill)
	}

	return BackgroundInfo{
		ID:                string(bg.ID),
		Name:              bg.Name,
		SkillProficienies: skillProfs,
		ToolProficiencies: bg.ToolProficiencies,
		Languages:         bg.LanguageProficiencies,
		Equipment:         bg.StartingEquipment,
		FeatureName:       bg.FeatureName,
		FeatureDesc:       bg.FeatureDescription,
		Description:       bg.Description,
	}
}

// monsterToInfo 将 MonsterStatBlock 转换为 MonsterInfo
func monsterToInfo(monster *model.MonsterStatBlock) MonsterInfo {
	return MonsterInfo{
		ID:              monster.ID,
		Name:            monster.Name,
		ChallengeRating: monster.ChallengeRating,
		ArmorClass:      monster.ArmorClass,
		HitPoints:       monster.HitPointsAverage,
		Speed:           monster.Speed.Walk,
		Size:            string(monster.Size),
		Type:            string(monster.CreatureType),
		Alignment:       monster.Alignment,
	}
}

// spellToInfo 将 Spell 转换为 SpellInfo
func spellToInfo(spell *model.Spell) SpellInfo {
	components := make([]string, 0)
	for _, comp := range spell.Components {
		components = append(components, string(comp))
	}

	// 格式化施法时间
	castingTime := fmt.Sprintf("%d %s", spell.CastTime.Value, spell.CastTime.Unit)

	return SpellInfo{
		ID:          spell.ID,
		Name:        spell.Name,
		Level:       spell.Level,
		School:      string(spell.School),
		CastingTime: castingTime,
		Range:       spell.Range,
		Duration:    spell.Duration,
		Components:  components,
		Classes:     spell.Classes,
		Description: spell.Description,
	}
}

// itemToInfo 将 Item 转换为 ItemInfo
func itemToInfo(item *model.Item) ItemInfo {
	return ItemInfo{
		ID:          item.ID.String(),
		Name:        item.Name,
		Type:        string(item.Type),
		Description: item.Description,
		Cost:        item.Value,
		Weight:      int(item.Weight),
	}
}

// weaponToInfo 将 Item 转换为 WeaponInfo
func weaponToInfo(item *model.Item) WeaponInfo {
	info := WeaponInfo{
		ItemInfo: itemToInfo(item),
	}
	if item.WeaponProps != nil {
		info.Damage = item.WeaponProps.DamageDice
		info.DamageType = string(item.WeaponProps.DamageType)
		// 构建武器属性字符串
		props := []string{}
		if item.WeaponProps.Light {
			props = append(props, "Light")
		}
		if item.WeaponProps.Finesse {
			props = append(props, "Finesse")
		}
		if item.WeaponProps.Heavy {
			props = append(props, "Heavy")
		}
		if item.WeaponProps.TwoHanded {
			props = append(props, "Two-Handed")
		}
		if item.WeaponProps.Versatile != "" {
			props = append(props, fmt.Sprintf("Versatile (%s)", item.WeaponProps.Versatile))
		}
		if item.WeaponProps.Loading {
			props = append(props, "Loading")
		}
		if item.WeaponProps.Thrown {
			props = append(props, "Thrown")
		}
		if item.WeaponProps.Reach {
			props = append(props, "Reach")
		}
		info.Property = strings.Join(props, ", ")
		// 构建射程
		if item.WeaponProps.Range > 0 {
			info.Range = fmt.Sprintf("%d/%d", item.WeaponProps.Range, item.WeaponProps.LongRange)
		}
	}
	return info
}

// armorToInfo 将 Item 转换为 ArmorInfo
func armorToInfo(item *model.Item) ArmorInfo {
	info := ArmorInfo{
		ItemInfo: itemToInfo(item),
	}
	if item.ArmorProps != nil {
		info.ArmorClass = fmt.Sprintf("%d", item.ArmorProps.BaseAC)
		info.StrRequired = item.ArmorProps.StrengthRequirement
		info.StealthDisadvantage = item.ArmorProps.StealthDisadvantage
	}
	return info
}

// recipeToInfo 将 CraftingRecipe 转换为 CraftingRecipeInfo
func recipeToInfo(recipe *model.CraftingRecipe) CraftingRecipeInfo {
	return CraftingRecipeInfo{
		ID:          recipe.ID,
		Name:        recipe.Name,
		Type:        string(recipe.Type),
		Description: recipe.Description,
		TimeDays:    recipe.TimeDays,
		DC:          recipe.DC,
		MinLevel:    recipe.MinLevel,
		Cost:        recipe.Cost,
	}
}

// lifestyleToInfo 将 LifestyleData 转换为 LifestyleInfo
func lifestyleToInfo(lifestyle *data.LifestyleData) LifestyleInfo {
	return LifestyleInfo{
		Tier:        lifestyle.Tier,
		Name:        string(lifestyle.Tier),
		DailyCost:   lifestyle.DailyCost,
		MonthlyCost: lifestyle.MonthlyCost,
		Description: lifestyle.Description,
	}
}

// mountToInfo 将 MountData 转换为 MountInfo
func mountToInfo(mount *data.MountData) MountInfo {
	return MountInfo{
		ID:               mount.ID.String(),
		Name:             mount.Name,
		Type:             "Beast", // MountData 没有 Type 字段，使用默认值
		Speed:            mount.Speed,
		CarryingCapacity: int(mount.CarryCap),
		Cost:             mount.Value,
		Description:      mount.Description,
	}
}

// poisonToInfo 将 PoisonDefinition 转换为 PoisonInfo
func poisonToInfo(poison *model.PoisonDefinition) PoisonInfo {
	return PoisonInfo{
		ID:          poison.ID,
		Name:        poison.Name,
		Type:        string(poison.Type),
		DC:          poison.Effect.SaveDC,
		Cost:        poison.Price,
		Description: poison.Description,
	}
}

// trapToInfo 将 TrapDefinition 转换为 TrapInfo
func trapToInfo(trap *model.TrapDefinition) TrapInfo {
	// 构建伤害字符串
	damage := ""
	if len(trap.Effects) > 0 {
		effectStrs := make([]string, 0, len(trap.Effects))
		for _, effect := range trap.Effects {
			if effect.DamageDice != "" {
				effectStrs = append(effectStrs, effect.DamageDice)
			}
		}
		damage = strings.Join(effectStrs, "; ")
	}

	return TrapInfo{
		ID:          trap.ID,
		Name:        trap.Name,
		Type:        string(trap.Type),
		DC:          trap.DetectDC, // 使用侦测 DC
		Damage:      damage,
		Description: trap.Description,
	}
}

// ============================================================================
// Engine 方法实现
// ============================================================================

// ListRaces 列出所有种族，支持分页
// 返回游戏中所有可用的种族定义，可按名称搜索，支持分页查询。
// 参数:
//
//	ctx - 上下文
//	req - 列表请求，包含可选的分页参数
//
// 返回:
//
//	*ListRacesResult - 种族列表和分页信息
//	error - 查询失败时返回错误
func (e *Engine) ListRaces(ctx context.Context, req ListRacesRequest) (*ListRacesResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	races := data.GlobalRegistry.ListRaces()

	// 按名称排序确保一致性
	sort.Slice(races, func(i, j int) bool {
		return races[i].Name < races[j].Name
	})

	totalCount := len(races)
	pagination := calculatePagination(totalCount, req.Pagination)
	pagedRaces := paginateSlice(races, req.Pagination)

	result := make([]RaceInfo, len(pagedRaces))
	for i, race := range pagedRaces {
		result[i] = raceToInfo(race)
	}

	return &ListRacesResult{
		Races:      result,
		Pagination: pagination,
	}, nil
}

// GetRace 获取指定种族的详细信息
// 根据种族名称获取该种族的完整定义，包括属性加值、特性、语言等。
// 参数:
//
//	ctx - 上下文
//	req - 获取请求，包含种族名称
//
// 返回:
//
//	*GetRaceResult - 种族详细信息
//	error - 种族不存在时返回错误
func (e *Engine) GetRace(ctx context.Context, req GetRaceRequest) (*GetRaceResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	race, exists := data.GlobalRegistry.GetRace(req.Name)
	if !exists {
		return nil, fmt.Errorf("race not found: %s", req.Name)
	}

	return &GetRaceResult{
		Race: raceToInfo(race),
	}, nil
}

// ListClasses 列出所有职业，支持分页
// 返回游戏中所有可用的职业定义，支持分页查询。
// 参数:
//
//	ctx - 上下文
//	req - 列表请求，包含可选的分页参数
//
// 返回:
//
//	*ListClassesResult - 职业列表和分页信息
//	error - 查询失败时返回错误
func (e *Engine) ListClasses(ctx context.Context, req ListClassesRequest) (*ListClassesResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	classes := data.GlobalRegistry.ListClasses()

	// 按ID排序确保一致性
	sort.Slice(classes, func(i, j int) bool {
		return classes[i].ID < classes[j].ID
	})

	totalCount := len(classes)
	pagination := calculatePagination(totalCount, req.Pagination)
	pagedClasses := paginateSlice(classes, req.Pagination)

	result := make([]ClassDefinitionInfo, len(pagedClasses))
	for i, class := range pagedClasses {
		result[i] = classToInfo(class)
	}

	return &ListClassesResult{
		Classes:    result,
		Pagination: pagination,
	}, nil
}

// GetClass 获取指定职业的详细信息
// 根据职业ID获取该职业的完整定义，包括生命骰、属性要求、技能等。
// 参数:
//
//	ctx - 上下文
//	req - 获取请求，包含职业ID
//
// 返回:
//
//	*GetClassResult - 职业详细信息
//	error - 职业不存在时返回错误
func (e *Engine) GetClass(ctx context.Context, req GetClassRequest) (*GetClassResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	class, exists := data.GlobalRegistry.GetClass(req.ID)
	if !exists {
		return nil, fmt.Errorf("class not found: %s", req.ID)
	}

	return &GetClassResult{
		Class: classToInfo(class),
	}, nil
}

// ListBackgrounds 列出所有背景，支持分页
// 返回游戏中所有可用的背景定义，支持分页查询。
// 参数:
//
//	ctx - 上下文
//	req - 列表请求，包含可选的分页参数
//
// 返回:
//
//	*ListBackgroundsResult - 背景列表和分页信息
//	error - 查询失败时返回错误
func (e *Engine) ListBackgrounds(ctx context.Context, req ListBackgroundsRequest) (*ListBackgroundsResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	backgrounds := data.GlobalRegistry.ListBackgrounds()

	// 按ID排序确保一致性
	sort.Slice(backgrounds, func(i, j int) bool {
		return backgrounds[i].ID < backgrounds[j].ID
	})

	totalCount := len(backgrounds)
	pagination := calculatePagination(totalCount, req.Pagination)
	pagedBackgrounds := paginateSlice(backgrounds, req.Pagination)

	result := make([]BackgroundInfo, len(pagedBackgrounds))
	for i, bg := range pagedBackgrounds {
		result[i] = backgroundToInfo(bg)
	}

	return &ListBackgroundsResult{
		Backgrounds: result,
		Pagination:  pagination,
	}, nil
}

// GetBackground 获取指定背景的详细信息
// 根据背景ID获取该背景的完整定义，包括技能熟练、特性等。
// 参数:
//
//	ctx - 上下文
//	req - 获取请求，包含背景ID
//
// 返回:
//
//	*GetBackgroundResult - 背景详细信息
//	error - 背景不存在时返回错误
func (e *Engine) GetBackground(ctx context.Context, req GetBackgroundRequest) (*GetBackgroundResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	bg, exists := data.GlobalRegistry.GetBackground(req.ID)
	if !exists {
		return nil, fmt.Errorf("background not found: %s", req.ID)
	}

	return &GetBackgroundResult{
		Background: backgroundToInfo(bg),
	}, nil
}

// ListFeatsData 列出所有专长，支持分页
// 返回游戏中所有可用的专长定义，支持分页查询。
// 参数:
//
//	ctx - 上下文
//	req - 列表请求，包含可选的分页参数
//
// 返回:
//
//	*ListFeatsDataResult - 专长列表和分页信息
//	error - 查询失败时返回错误
func (e *Engine) ListFeatsData(ctx context.Context, req ListFeatsDataRequest) (*ListFeatsDataResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	feats := data.GlobalRegistry.ListFeats()

	// 按ID排序确保一致性
	sort.Slice(feats, func(i, j int) bool {
		return feats[i].ID < feats[j].ID
	})

	totalCount := len(feats)
	pagination := calculatePagination(totalCount, req.Pagination)
	pagedFeats := paginateSlice(feats, req.Pagination)

	result := make([]FeatInfo, len(pagedFeats))
	for i, feat := range pagedFeats {
		result[i] = featToInfoForData(feat)
	}

	return &ListFeatsDataResult{
		Feats:      result,
		Pagination: pagination,
	}, nil
}

// GetFeatData 获取指定专长的详细信息（用于数据查询）
// 根据专长ID获取该专长的完整定义，包括前置条件、效果等。
// 参数:
//
//	ctx - 上下文
//	req - 获取请求，包含专长ID
//
// 返回:
//
//	*GetFeatDataResult - 专长详细信息
//	error - 专长不存在时返回错误
func (e *Engine) GetFeatData(ctx context.Context, req GetFeatDataRequest) (*GetFeatDataResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	feat, exists := data.GlobalRegistry.GetFeat(req.ID)
	if !exists {
		return nil, fmt.Errorf("feat not found: %s", req.ID)
	}

	return &GetFeatDataResult{
		Feat: featToInfoForData(feat),
	}, nil
}

// featToInfoForData 将 FeatDefinition 转换为 FeatInfo（用于数据查询）
func featToInfoForData(feat *model.FeatDefinition) FeatInfo {
	return FeatInfo{
		ID:           feat.ID,
		Name:         feat.Name,
		Type:         string(feat.Type),
		Description:  feat.Description,
		Repeatable:   feat.Repeatable,
		Prerequisite: feat.Prerequisite.Description,
	}
}

// ListMonsters 列出所有怪物，支持分页
// 返回游戏中所有可用的怪物模板，支持分页查询。
// 参数:
//
//	ctx - 上下文
//	req - 列表请求，包含可选的分页参数
//
// 返回:
//
//	*ListMonstersResult - 怪物列表和分页信息
//	error - 查询失败时返回错误
func (e *Engine) ListMonsters(ctx context.Context, req ListMonstersRequest) (*ListMonstersResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	monsters := data.GlobalRegistry.ListMonsters()

	// 按ID排序确保一致性
	sort.Slice(monsters, func(i, j int) bool {
		return monsters[i].ID < monsters[j].ID
	})

	totalCount := len(monsters)
	pagination := calculatePagination(totalCount, req.Pagination)
	pagedMonsters := paginateSlice(monsters, req.Pagination)

	result := make([]MonsterInfo, len(pagedMonsters))
	for i, monster := range pagedMonsters {
		result[i] = monsterToInfo(monster)
	}

	return &ListMonstersResult{
		Monsters:   result,
		Pagination: pagination,
	}, nil
}

// GetMonster 获取指定怪物的详细信息
// 根据怪物ID获取该怪物的完整定义，包括属性、技能、动作等。
// 参数:
//
//	ctx - 上下文
//	req - 获取请求，包含怪物ID
//
// 返回:
//
//	*GetMonsterResult - 怪物详细信息
//	error - 怪物不存在时返回错误
func (e *Engine) GetMonster(ctx context.Context, req GetMonsterRequest) (*GetMonsterResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	monster, exists := data.GlobalRegistry.GetMonster(req.ID)
	if !exists {
		return nil, fmt.Errorf("monster not found: %s", req.ID)
	}

	return &GetMonsterResult{
		Monster: monsterToInfo(monster),
	}, nil
}

// ListSpells 列出所有法术，支持分页
// 返回游戏中所有可用的法术定义，支持分页查询。
// 参数:
//
//	ctx - 上下文
//	req - 列表请求，包含可选的分页参数
//
// 返回:
//
//	*ListSpellsResult - 法术列表和分页信息
//	error - 查询失败时返回错误
func (e *Engine) ListSpells(ctx context.Context, req ListSpellsRequest) (*ListSpellsResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	spells := data.GlobalRegistry.ListSpells()

	// 按ID排序确保一致性
	sort.Slice(spells, func(i, j int) bool {
		return spells[i].ID < spells[j].ID
	})

	totalCount := len(spells)
	pagination := calculatePagination(totalCount, req.Pagination)
	pagedSpells := paginateSlice(spells, req.Pagination)

	result := make([]SpellInfo, len(pagedSpells))
	for i, spell := range pagedSpells {
		result[i] = spellToInfo(spell)
	}

	return &ListSpellsResult{
		Spells:     result,
		Pagination: pagination,
	}, nil
}

// GetSpell 获取指定法术的详细信息
// 根据法术ID获取该法术的完整定义，包括施法时间、范围、效果等。
// 参数:
//
//	ctx - 上下文
//	req - 获取请求，包含法术ID
//
// 返回:
//
//	*GetSpellResult - 法术详细信息
//	error - 法术不存在时返回错误
func (e *Engine) GetSpell(ctx context.Context, req GetSpellRequest) (*GetSpellResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	spell, exists := data.GlobalRegistry.GetSpell(req.ID)
	if !exists {
		return nil, fmt.Errorf("spell not found: %s", req.ID)
	}

	return &GetSpellResult{
		Spell: spellToInfo(spell),
	}, nil
}

// ListWeapons 列出所有武器，支持分页
// 返回游戏中所有可用的武器定义，支持分页查询。
// 参数:
//
//	ctx - 上下文
//	req - 列表请求，包含可选的分页参数
//
// 返回:
//
//	*ListWeaponsResult - 武器列表和分页信息
//	error - 查询失败时返回错误
func (e *Engine) ListWeapons(ctx context.Context, req ListWeaponsRequest) (*ListWeaponsResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	weapons := data.GlobalRegistry.ListWeapons()

	// 按ID排序确保一致性
	sort.Slice(weapons, func(i, j int) bool {
		return weapons[i].ID.String() < weapons[j].ID.String()
	})

	totalCount := len(weapons)
	pagination := calculatePagination(totalCount, req.Pagination)
	pagedWeapons := paginateSlice(weapons, req.Pagination)

	result := make([]WeaponInfo, len(pagedWeapons))
	for i, weapon := range pagedWeapons {
		result[i] = weaponToInfo(weapon)
	}

	return &ListWeaponsResult{
		Weapons:    result,
		Pagination: pagination,
	}, nil
}

// GetWeapon 获取指定武器的详细信息
// 根据武器ID获取该武器的完整定义，包括伤害、属性等。
// 参数:
//
//	ctx - 上下文
//	req - 获取请求，包含武器ID
//
// 返回:
//
//	*GetWeaponResult - 武器详细信息
//	error - 武器不存在时返回错误
func (e *Engine) GetWeapon(ctx context.Context, req GetWeaponRequest) (*GetWeaponResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	weapon, exists := data.GlobalRegistry.GetWeapon(req.ID)
	if !exists {
		return nil, fmt.Errorf("weapon not found: %s", req.ID)
	}

	return &GetWeaponResult{
		Weapon: weaponToInfo(weapon),
	}, nil
}

// ListArmors 列出所有护甲，支持分页
// 返回游戏中所有可用的护甲定义，支持分页查询。
// 参数:
//
//	ctx - 上下文
//	req - 列表请求，包含可选的分页参数
//
// 返回:
//
//	*ListArmorsResult - 护甲列表和分页信息
//	error - 查询失败时返回错误
func (e *Engine) ListArmors(ctx context.Context, req ListArmorsRequest) (*ListArmorsResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	armors := data.GlobalRegistry.ListArmors()

	// 按ID排序确保一致性
	sort.Slice(armors, func(i, j int) bool {
		return armors[i].ID.String() < armors[j].ID.String()
	})

	totalCount := len(armors)
	pagination := calculatePagination(totalCount, req.Pagination)
	pagedArmors := paginateSlice(armors, req.Pagination)

	result := make([]ArmorInfo, len(pagedArmors))
	for i, armor := range pagedArmors {
		result[i] = armorToInfo(armor)
	}

	return &ListArmorsResult{
		Armors:     result,
		Pagination: pagination,
	}, nil
}

// GetArmor 获取指定护甲的详细信息
// 根据护甲ID获取该护甲的完整定义，包括AC、力量要求等。
// 参数:
//
//	ctx - 上下文
//	req - 获取请求，包含护甲ID
//
// 返回:
//
//	*GetArmorResult - 护甲详细信息
//	error - 护甲不存在时返回错误
func (e *Engine) GetArmor(ctx context.Context, req GetArmorRequest) (*GetArmorResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	armor, exists := data.GlobalRegistry.GetArmor(req.ID)
	if !exists {
		return nil, fmt.Errorf("armor not found: %s", req.ID)
	}

	return &GetArmorResult{
		Armor: armorToInfo(armor),
	}, nil
}

// ListMagicItems 列出所有魔法物品，支持分页
// 返回游戏中所有可用的魔法物品定义，支持分页查询。
// 参数:
//
//	ctx - 上下文
//	req - 列表请求，包含可选的分页参数
//
// 返回:
//
//	*ListMagicItemsResult - 魔法物品列表和分页信息
//	error - 查询失败时返回错误
func (e *Engine) ListMagicItems(ctx context.Context, req ListMagicItemsRequest) (*ListMagicItemsResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	items := data.GlobalRegistry.ListMagicItems()

	// 按ID排序确保一致性
	sort.Slice(items, func(i, j int) bool {
		return items[i].ID.String() < items[j].ID.String()
	})

	totalCount := len(items)
	pagination := calculatePagination(totalCount, req.Pagination)
	pagedItems := paginateSlice(items, req.Pagination)

	result := make([]ItemInfo, len(pagedItems))
	for i, item := range pagedItems {
		result[i] = itemToInfo(item)
	}

	return &ListMagicItemsResult{
		MagicItems: result,
		Pagination: pagination,
	}, nil
}

// GetMagicItem 获取指定魔法物品的详细信息
// 根据魔法物品ID获取该物品的完整定义。
// 参数:
//
//	ctx - 上下文
//	req - 获取请求，包含魔法物品ID
//
// 返回:
//
//	*GetMagicItemResult - 魔法物品详细信息
//	error - 魔法物品不存在时返回错误
func (e *Engine) GetMagicItem(ctx context.Context, req GetMagicItemRequest) (*GetMagicItemResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	item, exists := data.GlobalRegistry.GetMagicItem(req.ID)
	if !exists {
		return nil, fmt.Errorf("magic item not found: %s", req.ID)
	}

	return &GetMagicItemResult{
		MagicItem: itemToInfo(item),
	}, nil
}

// ListGears 列出所有冒险装备，支持分页
// 返回游戏中所有可用的冒险装备定义，支持分页查询。
// 参数:
//
//	ctx - 上下文
//	req - 列表请求，包含可选的分页参数
//
// 返回:
//
//	*ListGearsResult - 冒险装备列表和分页信息
//	error - 查询失败时返回错误
func (e *Engine) ListGears(ctx context.Context, req ListGearsRequest) (*ListGearsResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	gears := data.GlobalRegistry.ListGears()

	// 按ID排序确保一致性
	sort.Slice(gears, func(i, j int) bool {
		return gears[i].ID.String() < gears[j].ID.String()
	})

	totalCount := len(gears)
	pagination := calculatePagination(totalCount, req.Pagination)
	pagedGears := paginateSlice(gears, req.Pagination)

	result := make([]ItemInfo, len(pagedGears))
	for i, gear := range pagedGears {
		result[i] = itemToInfo(gear)
	}

	return &ListGearsResult{
		Gears:      result,
		Pagination: pagination,
	}, nil
}

// GetGear 获取指定冒险装备的详细信息
// 根据冒险装备ID获取该装备的完整定义。
// 参数:
//
//	ctx - 上下文
//	req - 获取请求，包含冒险装备ID
//
// 返回:
//
//	*GetGearResult - 冒险装备详细信息
//	error - 冒险装备不存在时返回错误
func (e *Engine) GetGear(ctx context.Context, req GetGearRequest) (*GetGearResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	gear, exists := data.GlobalRegistry.GetGear(req.ID)
	if !exists {
		return nil, fmt.Errorf("gear not found: %s", req.ID)
	}

	return &GetGearResult{
		Gear: itemToInfo(gear),
	}, nil
}

// ListTools 列出所有工具，支持分页
// 返回游戏中所有可用的工具定义，支持分页查询。
// 参数:
//
//	ctx - 上下文
//	req - 列表请求，包含可选的分页参数
//
// 返回:
//
//	*ListToolsResult - 工具列表和分页信息
//	error - 查询失败时返回错误
func (e *Engine) ListTools(ctx context.Context, req ListToolsRequest) (*ListToolsResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	tools := data.GlobalRegistry.ListTools()

	// 按ID排序确保一致性
	sort.Slice(tools, func(i, j int) bool {
		return tools[i].ID.String() < tools[j].ID.String()
	})

	totalCount := len(tools)
	pagination := calculatePagination(totalCount, req.Pagination)
	pagedTools := paginateSlice(tools, req.Pagination)

	result := make([]ItemInfo, len(pagedTools))
	for i, tool := range pagedTools {
		result[i] = itemToInfo(tool)
	}

	return &ListToolsResult{
		Tools:      result,
		Pagination: pagination,
	}, nil
}

// GetTool 获取指定工具的详细信息
// 根据工具ID获取该工具的完整定义。
// 参数:
//
//	ctx - 上下文
//	req - 获取请求，包含工具ID
//
// 返回:
//
//	*GetToolResult - 工具详细信息
//	error - 工具不存在时返回错误
func (e *Engine) GetTool(ctx context.Context, req GetToolRequest) (*GetToolResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	tool, exists := data.GlobalRegistry.GetTool(req.ID)
	if !exists {
		return nil, fmt.Errorf("tool not found: %s", req.ID)
	}

	return &GetToolResult{
		Tool: itemToInfo(tool),
	}, nil
}

// ListRecipes 列出所有制作配方，支持分页
// 返回游戏中所有可用的制作配方定义，支持分页查询。
// 参数:
//
//	ctx - 上下文
//	req - 列表请求，包含可选的分页参数
//
// 返回:
//
//	*ListRecipesResult - 配方列表和分页信息
//	error - 查询失败时返回错误
func (e *Engine) ListRecipes(ctx context.Context, req ListRecipesRequest) (*ListRecipesResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	recipes := data.GlobalRegistry.ListCraftingRecipes()

	// 按ID排序确保一致性
	sort.Slice(recipes, func(i, j int) bool {
		return recipes[i].ID < recipes[j].ID
	})

	totalCount := len(recipes)
	pagination := calculatePagination(totalCount, req.Pagination)
	pagedRecipes := paginateSlice(recipes, req.Pagination)

	result := make([]CraftingRecipeInfo, len(pagedRecipes))
	for i, recipe := range pagedRecipes {
		result[i] = recipeToInfo(recipe)
	}

	return &ListRecipesResult{
		Recipes:    result,
		Pagination: pagination,
	}, nil
}

// GetRecipe 获取指定制作配方的详细信息
// 根据配方ID获取该配方的完整定义，包括材料、DC、时间等。
// 参数:
//
//	ctx - 上下文
//	req - 获取请求，包含配方ID
//
// 返回:
//
//	*GetRecipeResult - 配方详细信息
//	error - 配方不存在时返回错误
func (e *Engine) GetRecipe(ctx context.Context, req GetRecipeRequest) (*GetRecipeResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	recipe, exists := data.GlobalRegistry.GetCraftingRecipe(req.ID)
	if !exists {
		return nil, fmt.Errorf("recipe not found: %s", req.ID)
	}

	return &GetRecipeResult{
		Recipe: recipeToInfo(recipe),
	}, nil
}

// ListLifestyles 列出所有生活方式，支持分页
// 返回游戏中所有可用的生活方式定义，支持分页查询。
// 参数:
//
//	ctx - 上下文
//	req - 列表请求，包含可选的分页参数
//
// 返回:
//
//	*ListLifestylesResult - 生活方式列表和分页信息
//	error - 查询失败时返回错误
func (e *Engine) ListLifestylesData(ctx context.Context, req ListLifestylesRequest) (*ListLifestylesResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	lifestyles := data.GlobalRegistry.ListLifestyles()

	// 按Tier排序确保一致性
	sort.Slice(lifestyles, func(i, j int) bool {
		return lifestyles[i].Tier < lifestyles[j].Tier
	})

	totalCount := len(lifestyles)
	pagination := calculatePagination(totalCount, req.Pagination)
	pagedLifestyles := paginateSlice(lifestyles, req.Pagination)

	result := make([]LifestyleInfo, len(pagedLifestyles))
	for i, lifestyle := range pagedLifestyles {
		result[i] = lifestyleToInfo(lifestyle)
	}

	return &ListLifestylesResult{
		Lifestyles: result,
		Pagination: pagination,
	}, nil
}

// GetLifestyleData 获取指定生活方式的详细信息
// 根据生活方式等级获取该生活方式的完整定义，包括花费、描述等。
// 参数:
//
//	ctx - 上下文
//	req - 获取请求，包含生活方式等级
//
// 返回:
//
//	*GetLifestyleDataResult - 生活方式详细信息
//	error - 生活方式不存在时返回错误
func (e *Engine) GetLifestyleData(ctx context.Context, req GetLifestyleDataRequest) (*GetLifestyleDataResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	lifestyle, exists := data.GlobalRegistry.GetLifestyle(req.Tier)
	if !exists {
		return nil, fmt.Errorf("lifestyle not found: %s", req.Tier)
	}

	return &GetLifestyleDataResult{
		Lifestyle: lifestyleToInfo(lifestyle),
	}, nil
}

// ListMounts 列出所有坐骑，支持分页
// 返回游戏中所有可用的坐骑定义，支持分页查询。
// 参数:
//
//	ctx - 上下文
//	req - 列表请求，包含可选的分页参数
//
// 返回:
//
//	*ListMountsResult - 坐骑列表和分页信息
//	error - 查询失败时返回错误
func (e *Engine) ListMounts(ctx context.Context, req ListMountsRequest) (*ListMountsResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	mounts := data.GlobalRegistry.ListMounts()

	// 按ID排序确保一致性
	sort.Slice(mounts, func(i, j int) bool {
		return mounts[i].ID.String() < mounts[j].ID.String()
	})

	totalCount := len(mounts)
	pagination := calculatePagination(totalCount, req.Pagination)
	pagedMounts := paginateSlice(mounts, req.Pagination)

	result := make([]MountInfo, len(pagedMounts))
	for i, mount := range pagedMounts {
		result[i] = mountToInfo(mount)
	}

	return &ListMountsResult{
		Mounts:     result,
		Pagination: pagination,
	}, nil
}

// GetMount 获取指定坐骑的详细信息
// 根据坐骑ID获取该坐骑的完整定义，包括速度、载重等。
// 参数:
//
//	ctx - 上下文
//	req - 获取请求，包含坐骑ID
//
// 返回:
//
//	*GetMountResult - 坐骑详细信息
//	error - 坐骑不存在时返回错误
func (e *Engine) GetMount(ctx context.Context, req GetMountRequest) (*GetMountResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	mount, exists := data.GlobalRegistry.GetMount(req.ID)
	if !exists {
		return nil, fmt.Errorf("mount not found: %s", req.ID)
	}

	return &GetMountResult{
		Mount: mountToInfo(mount),
	}, nil
}

// ListPoisons 列出所有毒药，支持分页
// 返回游戏中所有可用的毒药定义，支持分页查询。
// 参数:
//
//	ctx - 上下文
//	req - 列表请求，包含可选的分页参数
//
// 返回:
//
//	*ListPoisonsResult - 毒药列表和分页信息
//	error - 查询失败时返回错误
func (e *Engine) ListPoisons(ctx context.Context, req ListPoisonsRequest) (*ListPoisonsResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	poisons := data.GlobalRegistry.ListPoisons()

	// 按ID排序确保一致性
	sort.Slice(poisons, func(i, j int) bool {
		return poisons[i].ID < poisons[j].ID
	})

	totalCount := len(poisons)
	pagination := calculatePagination(totalCount, req.Pagination)
	pagedPoisons := paginateSlice(poisons, req.Pagination)

	result := make([]PoisonInfo, len(pagedPoisons))
	for i, poison := range pagedPoisons {
		result[i] = poisonToInfo(poison)
	}

	return &ListPoisonsResult{
		Poisons:    result,
		Pagination: pagination,
	}, nil
}

// GetPoison 获取指定毒药的详细信息
// 根据毒药ID获取该毒药的完整定义，包括类型、DC、效果等。
// 参数:
//
//	ctx - 上下文
//	req - 获取请求，包含毒药ID
//
// 返回:
//
//	*GetPoisonResult - 毒药详细信息
//	error - 毒药不存在时返回错误
func (e *Engine) GetPoison(ctx context.Context, req GetPoisonRequest) (*GetPoisonResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	poison, exists := data.GlobalRegistry.GetPoison(req.ID)
	if !exists {
		return nil, fmt.Errorf("poison not found: %s", req.ID)
	}

	return &GetPoisonResult{
		Poison: poisonToInfo(poison),
	}, nil
}

// ListTraps 列出所有陷阱，支持分页
// 返回游戏中所有可用的陷阱定义，支持分页查询。
// 参数:
//
//	ctx - 上下文
//	req - 列表请求，包含可选的分页参数
//
// 返回:
//
//	*ListTrapsResult - 陷阱列表和分页信息
//	error - 查询失败时返回错误
func (e *Engine) ListTraps(ctx context.Context, req ListTrapsRequest) (*ListTrapsResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	traps := data.GlobalRegistry.ListTraps()

	// 按ID排序确保一致性
	sort.Slice(traps, func(i, j int) bool {
		return traps[i].ID < traps[j].ID
	})

	totalCount := len(traps)
	pagination := calculatePagination(totalCount, req.Pagination)
	pagedTraps := paginateSlice(traps, req.Pagination)

	result := make([]TrapInfo, len(pagedTraps))
	for i, trap := range pagedTraps {
		result[i] = trapToInfo(trap)
	}

	return &ListTrapsResult{
		Traps:      result,
		Pagination: pagination,
	}, nil
}

// GetTrap 获取指定陷阱的详细信息
// 根据陷阱ID获取该陷阱的完整定义，包括类型、DC、伤害等。
// 参数:
//
//	ctx - 上下文
//	req - 获取请求，包含陷阱ID
//
// 返回:
//
//	*GetTrapResult - 陷阱详细信息
//	error - 陷阱不存在时返回错误
func (e *Engine) GetTrap(ctx context.Context, req GetTrapRequest) (*GetTrapResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	trap, exists := data.GlobalRegistry.GetTrap(req.ID)
	if !exists {
		return nil, fmt.Errorf("trap not found: %s", req.ID)
	}

	return &GetTrapResult{
		Trap: trapToInfo(trap),
	}, nil
}
