package engine

import (
	"context"
	"fmt"

	"github.com/zwh8800/dnd-core/pkg/model"
	"github.com/zwh8800/dnd-core/pkg/rules"
)

// StateSummary 游戏状态摘要
type StateSummary struct {
	GameName     string                `json:"game_name"`
	Phase        model.Phase           `json:"phase"`
	CurrentScene *SceneSummary         `json:"current_scene"`
	PartyMembers []model.ActorSnapshot `json:"party_members"`
	ActiveCombat *CombatSummary        `json:"active_combat"`
	ActiveQuests []QuestSummary        `json:"active_quests"`
	Time         string                `json:"time"`
}

// SceneSummary 场景摘要
type SceneSummary struct {
	ID   model.ID `json:"id"`
	Name string   `json:"name"`
}

// QuestSummary 任务摘要
type QuestSummary struct {
	ID     model.ID          `json:"id"`
	Name   string            `json:"name"`
	Status model.QuestStatus `json:"status"`
}

// CombatSummary 战斗摘要
type CombatSummary struct {
	Round        int              `json:"round"`
	TurnOrder    []TurnOrderEntry `json:"turn_order"`
	CurrentActor string           `json:"current_actor"`
	Combatants   []CombatantBrief `json:"combatants"`
}

// TurnOrderEntry 先攻顺序条目
type TurnOrderEntry struct {
	ActorName  string `json:"actor_name"`
	Initiative int    `json:"initiative"`
	IsCurrent  bool   `json:"is_current"`
}

// CombatantBrief 战斗者简要信息
type CombatantBrief struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	HP         int      `json:"hp"`
	MaxHP      int      `json:"max_hp"`
	AC         int      `json:"ac"`
	Conditions []string `json:"conditions"`
	IsDefeated bool     `json:"is_defeated"`
}

// ActorSheet 完整角色卡
type ActorSheet struct {
	BasicInfo     string           `json:"basic_info"`
	AbilityScores map[string]int   `json:"ability_scores"`
	Skills        map[string]int   `json:"skills"`
	SavingThrows  map[string]int   `json:"saving_throws"`
	Combat        CombatSheetInfo  `json:"combat"`
	Spellcasting  *SpellSheetInfo  `json:"spellcasting,omitempty"`
	Equipment     []EquipmentEntry `json:"equipment"`
	Conditions    []string         `json:"conditions"`
	Features      []string         `json:"features"`
}

// CombatSheetInfo 战斗信息
type CombatSheetInfo struct {
	HP         int    `json:"hp"`
	MaxHP      int    `json:"max_hp"`
	TempHP     int    `json:"temp_hp"`
	AC         int    `json:"ac"`
	Speed      int    `json:"speed"`
	HitDice    string `json:"hit_dice"`
	DeathSaves string `json:"death_saves"`
}

// SpellSheetInfo 法术信息
type SpellSheetInfo struct {
	Ability        string      `json:"ability"`
	SaveDC         int         `json:"save_dc"`
	AttackBonus    int         `json:"attack_bonus"`
	SlotsRemaining map[int]int `json:"slots_remaining"`
	PreparedSpells []string    `json:"prepared_spells"`
}

// EquipmentEntry 装备条目
type EquipmentEntry struct {
	Name     string `json:"name"`
	Slot     string `json:"slot"`
	Equipped bool   `json:"equipped"`
}

// GetStateSummaryRequest 获取状态摘要请求
type GetStateSummaryRequest struct {
	GameID model.ID `json:"game_id"` // 游戏会话ID
}

// GetActorSheetRequest 获取角色卡请求
type GetActorSheetRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID
	ActorID model.ID `json:"actor_id"` // 角色ID
}

// GetCombatSummaryRequest 获取战斗摘要请求
type GetCombatSummaryRequest struct {
	GameID model.ID `json:"game_id"` // 游戏会话ID
}

// GetStateSummary 获取LLM友好的游戏状态摘要
func (e *Engine) GetStateSummary(ctx context.Context, gameID model.ID) (*StateSummary, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	summary := &StateSummary{
		GameName:     game.Name,
		Phase:        game.Phase,
		PartyMembers: make([]model.ActorSnapshot, 0),
		ActiveQuests: make([]QuestSummary, 0),
	}

	// 当前场景
	if game.CurrentScene != nil {
		if scene, ok := game.Scenes[*game.CurrentScene]; ok {
			summary.CurrentScene = &SceneSummary{
				ID:   scene.ID,
				Name: scene.Name,
			}
		}
	}

	// 队伍成员（PC）
	for _, pc := range game.PCs {
		snapshot := model.ActorToSnapshot(&pc.Actor, model.ActorTypePC, pc.Name)
		summary.PartyMembers = append(summary.PartyMembers, snapshot)
	}

	// 活跃战斗
	if game.Combat != nil && game.Combat.Status == model.CombatStatusActive {
		summary.ActiveCombat = buildCombatSummary(game)
	}

	// 活跃任务
	for _, quest := range game.Quests {
		if quest.Status == model.QuestStatusActive {
			summary.ActiveQuests = append(summary.ActiveQuests, QuestSummary{
				ID:     quest.ID,
				Name:   quest.Name,
				Status: quest.Status,
			})
		}
	}

	// 游戏时间
	summary.Time = fmt.Sprintf("Year %d, Month %d, Day %d, %02d:%02d",
		game.GameTime.Year, game.GameTime.Month, game.GameTime.Day,
		game.GameTime.Hour, game.GameTime.Minute)

	return summary, nil
}

// GetActorSheet 获取角色的完整角色卡
func (e *Engine) GetActorSheet(ctx context.Context, gameID model.ID, actorID model.ID) (*ActorSheet, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	actor, ok := game.GetActor(actorID)
	if !ok {
		return nil, ErrNotFound
	}

	var baseActor *model.Actor
	var name string
	var actorType model.ActorType

	switch a := actor.(type) {
	case *model.PlayerCharacter:
		baseActor = &a.Actor
		name = a.Name
		actorType = model.ActorTypePC
	case *model.NPC:
		baseActor = &a.Actor
		name = a.Name
		actorType = model.ActorTypeNPC
	case *model.Enemy:
		baseActor = &a.Actor
		name = a.Name
		actorType = model.ActorTypeEnemy
	case *model.Companion:
		baseActor = &a.Actor
		name = a.Name
		actorType = model.ActorTypeCompanion
	}

	sheet := &ActorSheet{
		BasicInfo:     fmt.Sprintf("%s (%s, Level %d)", name, actorType, getActorLevel(actor)),
		AbilityScores: make(map[string]int),
		Skills:        make(map[string]int),
		SavingThrows:  make(map[string]int),
		Conditions:    make([]string, 0),
		Features:      make([]string, 0),
	}

	// 属性值
	abilityScores := baseActor.AbilityScores
	sheet.AbilityScores["STR"] = abilityScores.Strength
	sheet.AbilityScores["DEX"] = abilityScores.Dexterity
	sheet.AbilityScores["CON"] = abilityScores.Constitution
	sheet.AbilityScores["INT"] = abilityScores.Intelligence
	sheet.AbilityScores["WIS"] = abilityScores.Wisdom
	sheet.AbilityScores["CHA"] = abilityScores.Charisma

	// 计算熟练加值
	level := getActorLevel(actor)
	profBonus := rules.ProficiencyBonus(level)

	// 技能修正值
	skills := model.AllSkills()
	for _, skill := range skills {
		ability := model.SkillAbilityMap[skill]
		abilityScore := abilityScores.Get(ability)
		mod := rules.SkillModifierWithLevel(abilityScore, level,
			baseActor.Proficiencies.IsProficient(skill),
			baseActor.Proficiencies.HasExpertise(skill),
			0)
		sheet.Skills[string(skill)] = mod
	}

	// 豁免检定
	for _, ability := range model.AllAbilities() {
		mod := rules.AbilityModifier(abilityScores.Get(ability))
		if baseActor.Proficiencies.IsSavingThrowProficient(ability) {
			mod += profBonus
		}
		sheet.SavingThrows[string(ability)] = mod
	}

	// 战斗信息
	sheet.Combat = CombatSheetInfo{
		HP:     baseActor.HitPoints.Current,
		MaxHP:  baseActor.HitPoints.Maximum,
		TempHP: baseActor.TempHitPoints,
		AC:     baseActor.ArmorClass,
		Speed:  baseActor.Speed,
	}

	// PC特有信息
	if pc, ok := actor.(*model.PlayerCharacter); ok {
		// 生命骰
		hitDiceStr := ""
		for _, hd := range pc.HitDice {
			remaining := hd.Total - hd.Used
			if hitDiceStr != "" {
				hitDiceStr += ", "
			}
			hitDiceStr += fmt.Sprintf("%d/%d d%d", remaining, hd.Total, hd.DiceType)
		}
		sheet.Combat.HitDice = hitDiceStr

		// 死亡豁免
		sheet.Combat.DeathSaves = fmt.Sprintf("Successes: %d, Failures: %d",
			pc.DeathSaveSuccesses, pc.DeathSaveFailures)

		// 法术信息
		if pc.Spellcasting != nil {
			sc := pc.Spellcasting
			sheet.Spellcasting = &SpellSheetInfo{
				Ability:        string(sc.SpellcastingAbility),
				SaveDC:         sc.SpellSaveDC,
				AttackBonus:    sc.SpellAttackBonus,
				PreparedSpells: sc.PreparedSpells,
			}
			if sc.Slots != nil {
				sheet.Spellcasting.SlotsRemaining = make(map[int]int)
				for i := 1; i <= 9; i++ {
					available := sc.Slots.GetAvailableSlots(i)
					if available > 0 {
						sheet.Spellcasting.SlotsRemaining[i] = available
					}
				}
			}
		}

		sheet.Features = append(sheet.Features, pc.Features...)
		sheet.Features = append(sheet.Features, pc.RacialTraits...)
	}

	// 状态效果
	for _, cond := range baseActor.Conditions {
		sheet.Conditions = append(sheet.Conditions, string(cond.Type))
	}

	return sheet, nil
}

// GetCombatSummary 获取战斗摘要
func (e *Engine) GetCombatSummary(ctx context.Context, gameID model.ID) (*CombatSummary, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	if game.Combat == nil || game.Combat.Status != model.CombatStatusActive {
		return nil, ErrCombatNotActive
	}

	return buildCombatSummary(game), nil
}

// buildCombatSummary 构建战斗摘要（内部方法）
func buildCombatSummary(game *model.GameState) *CombatSummary {
	combat := game.Combat
	summary := &CombatSummary{
		Round:      combat.Round,
		TurnOrder:  make([]TurnOrderEntry, 0),
		Combatants: make([]CombatantBrief, 0),
	}

	for i, entry := range combat.Initiative {
		isCurrent := i == combat.CurrentIndex
		summary.TurnOrder = append(summary.TurnOrder, TurnOrderEntry{
			ActorName:  entry.ActorName,
			Initiative: entry.InitiativeTotal,
			IsCurrent:  isCurrent,
		})

		if isCurrent {
			summary.CurrentActor = entry.ActorName
		}

		// 获取角色信息
		actor, ok := game.GetActor(entry.ActorID)
		if !ok {
			continue
		}

		var baseActor *model.Actor
		var name string
		var actorType model.ActorType

		switch a := actor.(type) {
		case *model.PlayerCharacter:
			baseActor = &a.Actor
			name = a.Name
			actorType = model.ActorTypePC
		case *model.NPC:
			baseActor = &a.Actor
			name = a.Name
			actorType = model.ActorTypeNPC
		case *model.Enemy:
			baseActor = &a.Actor
			name = a.Name
			actorType = model.ActorTypeEnemy
		case *model.Companion:
			baseActor = &a.Actor
			name = a.Name
			actorType = model.ActorTypeCompanion
		}

		conditions := make([]string, 0)
		for _, c := range baseActor.Conditions {
			conditions = append(conditions, string(c.Type))
		}

		summary.Combatants = append(summary.Combatants, CombatantBrief{
			Name:       name,
			Type:       string(actorType),
			HP:         baseActor.HitPoints.Current,
			MaxHP:      baseActor.HitPoints.Maximum,
			AC:         baseActor.ArmorClass,
			Conditions: conditions,
			IsDefeated: entry.IsDefeated,
		})
	}

	return summary
}

// getActorLevel 获取角色等级
func getActorLevel(actor any) int {
	switch a := actor.(type) {
	case *model.PlayerCharacter:
		return a.TotalLevel
	case *model.Companion:
		// 同伴等级通常与领导者相同
		return 1
	case *model.NPC:
		return 1
	case *model.Enemy:
		// 敌人使用挑战等级作为近似
		if e, ok := actor.(*model.Enemy); ok {
			return int(e.ChallengeRating)
		}
		return 1
	default:
		return 1
	}
}
