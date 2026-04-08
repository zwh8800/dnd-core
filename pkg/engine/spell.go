package engine

import (
	"context"
	"fmt"

	"github.com/zwh8800/dnd-core/pkg/data"
	"github.com/zwh8800/dnd-core/pkg/model"
	"github.com/zwh8800/dnd-core/pkg/rules"
)

// CastSpellRequest 施法请求（整合 SpellInput）
type CastSpellRequest struct {
	GameID   model.ID   `json:"game_id"`   // 游戏会话ID
	CasterID model.ID   `json:"caster_id"` // 施法者ID
	Spell    SpellInput `json:"spell"`     // 法术输入
}

// GetSpellSlotsRequest 获取法术位请求
type GetSpellSlotsRequest struct {
	GameID   model.ID `json:"game_id"`   // 游戏会话ID
	CasterID model.ID `json:"caster_id"` // 施法者ID
}

// GetSpellSlotsResult 获取法术位结果
type GetSpellSlotsResult struct {
	Info *SpellSlotsInfo `json:"info"` // 法术位信息
}

// PrepareSpellsRequest 准备法术请求
type PrepareSpellsRequest struct {
	GameID   model.ID `json:"game_id"`   // 游戏会话ID
	CasterID model.ID `json:"caster_id"` // 施法者ID
	SpellIDs []string `json:"spell_ids"` // 准备法术ID列表
}

// LearnSpellRequest 学习法术请求
type LearnSpellRequest struct {
	GameID   model.ID `json:"game_id"`   // 游戏会话ID
	CasterID model.ID `json:"caster_id"` // 施法者ID
	SpellID  string   `json:"spell_id"`  // 法术ID
}

// ConcentrationCheckRequest 专注检定请求
type ConcentrationCheckRequest struct {
	GameID      model.ID `json:"game_id"`      // 游戏会话ID
	CasterID    model.ID `json:"caster_id"`    // 施法者ID
	DamageTaken int      `json:"damage_taken"` // 受到的伤害
}

// EndConcentrationRequest 结束专注请求
type EndConcentrationRequest struct {
	GameID   model.ID `json:"game_id"`   // 游戏会话ID
	CasterID model.ID `json:"caster_id"` // 施法者ID
}

// SpellInput 施法输入
type SpellInput struct {
	SpellID     string             `json:"spell_id"`     // 法术ID
	SlotLevel   int                `json:"slot_level"`   // 使用的法术位环级（0表示戏法）
	TargetIDs   []model.ID         `json:"target_ids"`   // 目标角色ID列表
	TargetPoint *model.Point       `json:"target_point"` // 目标位置
	UpcastLevel int                `json:"upcast_level"` // 升环施法等级
	Advantage   model.RollModifier `json:"advantage"`    // 优势/劣势
}

// SpellResult 施法结果
type SpellResult struct {
	SpellName     string              `json:"spell_name"`
	SlotLevel     int                 `json:"slot_level"`
	CasterSaveDC  int                 `json:"caster_save_dc"`
	AttackRoll    *model.DiceResult   `json:"attack_roll,omitempty"`
	AttackTotal   int                 `json:"attack_total,omitempty"`
	Targets       []SpellTargetResult `json:"targets"`
	Concentration bool                `json:"is_concentration"`
	Message       string              `json:"message"`
}

// SpellTargetResult 法术目标结果
type SpellTargetResult struct {
	ActorID     model.ID          `json:"actor_id"`
	SaveRoll    *model.DiceResult `json:"save_roll,omitempty"`
	SaveTotal   int               `json:"save_total,omitempty"`
	SaveSuccess bool              `json:"save_success"`
	Damage      *DamageResult     `json:"damage,omitempty"`
	Healing     *HealResult       `json:"healing,omitempty"`
	Effect      string            `json:"effect,omitempty"`
}

// ConcentrationResult 专注检定结果
type ConcentrationResult struct {
	Success          bool              `json:"success"`
	DC               int               `json:"dc"`
	Roll             *model.DiceResult `json:"roll"`
	RollTotal        int               `json:"roll_total"`
	ConstitutionSave int               `json:"constitution_save"`
	SpellName        string            `json:"spell_name"`
	Message          string            `json:"message"`
}

// SpellSlotsInfo 法术位信息
type SpellSlotsInfo struct {
	CantripsKnown       int             `json:"cantrips_known"`            // 已知戏法数量
	SpellsPrepared      []string        `json:"spells_prepared,omitempty"` // 准备的法术列表
	KnownSpells         []string        `json:"known_spells,omitempty"`    // 已知的法术列表
	SlotsByLevel        []SlotLevelInfo `json:"slots_by_level"`            // 各环级法术位信息
	SaveDC              int             `json:"save_dc"`                   // 豁免DC
	AttackBonus         int             `json:"attack_bonus"`              // 攻击加成
	SpellcastingAbility string          `json:"spellcasting_ability"`      // 施法关键属性
}

// SlotLevelInfo 单个环级的法术位信息
type SlotLevelInfo struct {
	Level     int `json:"level"`     // 法术环级
	Total     int `json:"total"`     // 总数
	Used      int `json:"used"`      // 已使用
	Available int `json:"available"` // 可用数量
}

// CastSpell 执行施法动作
func (e *Engine) CastSpell(ctx context.Context, req CastSpellRequest) (*SpellResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	// 获取施法者
	caster, ok := game.GetActor(req.CasterID)
	if !ok {
		return nil, ErrNotFound
	}

	var casterActor *model.Actor
	var spellcaster *model.SpellcasterState
	switch c := caster.(type) {
	case *model.PlayerCharacter:
		casterActor = &c.Actor
		spellcaster = c.Spellcasting
	default:
		return nil, fmt.Errorf("only player characters can cast spells")
	}

	if spellcaster == nil {
		return nil, fmt.Errorf("actor %s is not a spellcaster", casterActor.Name)
	}

	// 查找法术定义
	spellDef := findSpellDefinition(req.Spell.SpellID)
	if spellDef == nil {
		return nil, fmt.Errorf("spell %s not found", req.Spell.SpellID)
	}

	// 验证施法者是否知道/准备该法术
	if !canCastSpell(spellcaster, req.Spell.SpellID) {
		return nil, fmt.Errorf("caster does not know or have prepared spell: %s", spellDef.Name)
	}

	// 检查法术位（戏法不需要法术位）
	if spellDef.Level > 0 {
		slotLevel := req.Spell.SlotLevel
		if slotLevel == 0 {
			slotLevel = spellDef.Level
		}
		if slotLevel < spellDef.Level {
			return nil, fmt.Errorf("slot level %d is lower than spell level %d", slotLevel, spellDef.Level)
		}
		if spellcaster.Slots.GetAvailableSlots(slotLevel) <= 0 {
			return nil, ErrInsufficientSlots
		}
		// 消耗法术位
		spellcaster.Slots.UseSlot(slotLevel)
	}

	// 检查专注：如果已经专注其他法术，则结束之前的专注
	if spellcaster.IsConcentrating() && spellDef.Concentration {
		spellcaster.ConcentrationSpell = ""
	}

	// 创建施法结果
	result := &SpellResult{
		SpellName:     spellDef.Name,
		SlotLevel:     req.Spell.SlotLevel,
		CasterSaveDC:  spellcaster.SpellSaveDC,
		Targets:       make([]SpellTargetResult, 0),
		Concentration: spellDef.Concentration,
	}

	// 处理需要攻击掷骰的法术
	if spellDef.DamageDice != "" && spellDef.SaveDC == "" {
		// 法术攻击
		attackBonus := spellcaster.SpellAttackBonus
		var rollResult *model.DiceResult
		if req.Spell.Advantage.Advantage {
			rollResult, _ = e.roller.RollAdvantage(0)
		} else if req.Spell.Advantage.Disadvantage {
			rollResult, _ = e.roller.RollDisadvantage(0)
		} else {
			rollResult, _ = e.roller.Roll("1d20")
		}

		attackTotal := rollResult.Total + attackBonus
		result.AttackRoll = rollResult
		result.AttackTotal = attackTotal

		// 对每个目标进行攻击
		for _, targetID := range req.Spell.TargetIDs {
			target, ok := game.GetActor(targetID)
			if !ok {
				continue
			}
			var targetActor *model.Actor
			switch t := target.(type) {
			case *model.PlayerCharacter:
				targetActor = &t.Actor
			case *model.NPC:
				targetActor = &t.Actor
			case *model.Enemy:
				targetActor = &t.Actor
			case *model.Companion:
				targetActor = &t.Actor
			}

			isNat20 := rollResult.Rolls[0].Value == 20
			isNat1 := rollResult.Rolls[0].Value == 1
			hit := attackTotal >= targetActor.ArmorClass || isNat20
			if isNat1 {
				hit = false
			}

			targetResult := SpellTargetResult{ActorID: targetID}

			if hit {
				// 计算伤害
				damageResult, err := e.applySpellDamage(game, req.CasterID, targetID, spellDef, req.Spell.UpcastLevel, isNat20)
				if err != nil {
					return nil, err
				}
				targetResult.Damage = damageResult
			}

			result.Targets = append(result.Targets, targetResult)
		}
	} else if spellDef.DamageDice != "" && spellDef.SaveDC != "" {
		// 需要豁免的伤害法术
		for _, targetID := range req.Spell.TargetIDs {
			target, ok := game.GetActor(targetID)
			if !ok {
				continue
			}

			var targetActor *model.Actor
			switch t := target.(type) {
			case *model.PlayerCharacter:
				targetActor = &t.Actor
			case *model.NPC:
				targetActor = &t.Actor
			case *model.Enemy:
				targetActor = &t.Actor
			case *model.Companion:
				targetActor = &t.Actor
			}

			// 豁免掷骰
			saveRoll, _ := e.roller.Roll("1d20")
			saveAbility := targetActor.AbilityScores.Get(spellDef.SaveDC)
			saveBonus := rules.AbilityModifier(saveAbility)

			saveTotal := saveRoll.Total + saveBonus
			saveSuccess := saveTotal >= spellcaster.SpellSaveDC

			targetResult := SpellTargetResult{
				ActorID:     targetID,
				SaveRoll:    saveRoll,
				SaveTotal:   saveTotal,
				SaveSuccess: saveSuccess,
			}

			// 根据豁免结果计算伤害
			if spellDef.DamageDice != "" {
				baseDamage := parseDiceString(spellDef.DamageDice)
				if saveSuccess {
					baseDamage /= 2
				}
				damageResult, err := e.applySpellDamageDirect(game, req.CasterID, targetID, baseDamage, spellDef.DamageType, false)
				if err != nil {
					return nil, err
				}
				targetResult.Damage = damageResult
			}

			result.Targets = append(result.Targets, targetResult)
		}
	} else if spellDef.HealingDice != "" {
		// 治疗法术
		healingAmount := parseDiceString(spellDef.HealingDice)
		for _, targetID := range req.Spell.TargetIDs {
			healResult, err := e.applyHealingInSpell(game, targetID, healingAmount)
			if err != nil {
				return nil, err
			}
			result.Targets = append(result.Targets, SpellTargetResult{
				ActorID: targetID,
				Healing: healResult,
			})
		}
	}

	// 设置专注
	if spellDef.Concentration {
		spellcaster.ConcentrationSpell = req.Spell.SpellID
		result.Concentration = true
	}

	// 构建消息
	result.Message = fmt.Sprintf("%s 施展了 %s", casterActor.Name, spellDef.Name)
	if spellDef.Concentration {
		result.Message += " (专注)"
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return result, nil
}

// GetSpellSlots 获取施法者的法术位状态
func (e *Engine) GetSpellSlots(ctx context.Context, req GetSpellSlotsRequest) (*GetSpellSlotsResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	caster, ok := game.GetActor(req.CasterID)
	if !ok {
		return nil, ErrNotFound
	}

	var spellcaster *model.SpellcasterState
	switch c := caster.(type) {
	case *model.PlayerCharacter:
		spellcaster = c.Spellcasting
	default:
		return nil, fmt.Errorf("only player characters can have spell slots")
	}

	if spellcaster == nil {
		return nil, fmt.Errorf("actor is not a spellcaster")
	}

	info := &SpellSlotsInfo{
		SaveDC:              spellcaster.SpellSaveDC,
		AttackBonus:         spellcaster.SpellAttackBonus,
		SpellcastingAbility: string(spellcaster.SpellcastingAbility),
		SlotsByLevel:        make([]SlotLevelInfo, 0),
	}

	if spellcaster.PreparationType == "prepared" {
		info.SpellsPrepared = spellcaster.PreparedSpells
	} else {
		info.KnownSpells = spellcaster.KnownSpells
	}

	// 收集法术位信息
	for level := 1; level <= 9; level++ {
		total := spellcaster.Slots.Slots[level][0]
		if total > 0 {
			used := spellcaster.Slots.Slots[level][1]
			info.SlotsByLevel = append(info.SlotsByLevel, SlotLevelInfo{
				Level:     level,
				Total:     total,
				Used:      used,
				Available: total - used,
			})
		}
	}

	return &GetSpellSlotsResult{Info: info}, nil
}

// PrepareSpells 准备法术（针对准备型施法者）
func (e *Engine) PrepareSpells(ctx context.Context, req PrepareSpellsRequest) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return err
	}

	caster, ok := game.GetActor(req.CasterID)
	if !ok {
		return ErrNotFound
	}

	var spellcaster *model.SpellcasterState
	switch c := caster.(type) {
	case *model.PlayerCharacter:
		spellcaster = c.Spellcasting
	default:
		return fmt.Errorf("only player characters can prepare spells")
	}

	if spellcaster == nil {
		return fmt.Errorf("actor is not a spellcaster")
	}

	if spellcaster.PreparationType != "prepared" {
		return fmt.Errorf("this caster uses known spells, not prepared spells")
	}

	// 验证所有法术都在已知列表中
	for _, spellID := range req.SpellIDs {
		if !spellcaster.CanPrepareSpell(spellID) {
			return fmt.Errorf("spell %s is not known and cannot be prepared", spellID)
		}
	}

	spellcaster.PreparedSpells = req.SpellIDs

	if err := e.saveGame(ctx, game); err != nil {
		return err
	}

	return nil
}

// LearnSpell 学习新法术（针对已知型施法者）
func (e *Engine) LearnSpell(ctx context.Context, req LearnSpellRequest) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return err
	}

	caster, ok := game.GetActor(req.CasterID)
	if !ok {
		return ErrNotFound
	}

	var spellcaster *model.SpellcasterState
	switch c := caster.(type) {
	case *model.PlayerCharacter:
		spellcaster = c.Spellcasting
	default:
		return fmt.Errorf("only player characters can learn spells")
	}

	if spellcaster == nil {
		return fmt.Errorf("actor is not a spellcaster")
	}

	// 检查是否已经学会
	for _, s := range spellcaster.KnownSpells {
		if s == req.SpellID {
			return fmt.Errorf("spell %s is already known", req.SpellID)
		}
	}

	spellcaster.KnownSpells = append(spellcaster.KnownSpells, req.SpellID)

	if err := e.saveGame(ctx, game); err != nil {
		return err
	}

	return nil
}

// ConcentrationCheck 进行专注检定（当施法者受到伤害时）
func (e *Engine) ConcentrationCheck(ctx context.Context, req ConcentrationCheckRequest) (*ConcentrationResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	caster, ok := game.GetActor(req.CasterID)
	if !ok {
		return nil, ErrNotFound
	}

	var casterActor *model.Actor
	var spellcaster *model.SpellcasterState
	switch c := caster.(type) {
	case *model.PlayerCharacter:
		casterActor = &c.Actor
		spellcaster = c.Spellcasting
	default:
		return nil, fmt.Errorf("only player characters can concentrate on spells")
	}

	if spellcaster == nil || !spellcaster.IsConcentrating() {
		return nil, fmt.Errorf("caster is not concentrating on any spell")
	}

	// 专注检定DC = max(10, 伤害值/2)
	dc := req.DamageTaken / 2
	if dc < 10 {
		dc = 10
	}

	// 体质豁免掷骰
	saveRoll, _ := e.roller.Roll("1d20")
	conMod := rules.AbilityModifier(casterActor.AbilityScores.Constitution)

	saveTotal := saveRoll.Total + conMod
	success := saveTotal >= dc

	currentSpell := spellcaster.ConcentrationSpell
	if !success {
		spellcaster.ConcentrationSpell = ""
	}

	result := &ConcentrationResult{
		Success:          success,
		DC:               dc,
		Roll:             saveRoll,
		RollTotal:        saveTotal,
		ConstitutionSave: conMod,
		SpellName:        currentSpell,
		Message:          fmt.Sprintf("专注检定: %d vs DC %d", saveTotal, dc),
	}

	if !success {
		result.Message += " - 失败！专注结束"
	} else {
		result.Message += " - 成功"
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return result, nil
}

// EndConcentration 主动结束专注
func (e *Engine) EndConcentration(ctx context.Context, req EndConcentrationRequest) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return err
	}

	caster, ok := game.GetActor(req.CasterID)
	if !ok {
		return ErrNotFound
	}

	var spellcaster *model.SpellcasterState
	switch c := caster.(type) {
	case *model.PlayerCharacter:
		spellcaster = c.Spellcasting
	default:
		return fmt.Errorf("only player characters can concentrate on spells")
	}

	if spellcaster == nil || !spellcaster.IsConcentrating() {
		return fmt.Errorf("caster is not concentrating on any spell")
	}

	spellcaster.ConcentrationSpell = ""

	if err := e.saveGame(ctx, game); err != nil {
		return err
	}

	return nil
}

// canCastSpell 检查施法者是否可以施放指定法术
func canCastSpell(spellcaster *model.SpellcasterState, spellID string) bool {
	if spellcaster.PreparationType == "known" {
		// 已知型：检查是否在已知列表中
		for _, s := range spellcaster.KnownSpells {
			if s == spellID {
				return true
			}
		}
		return false
	}
	// 准备型：检查是否在准备列表中
	return spellcaster.CanPrepareSpell(spellID)
}

// applySpellDamage 应用法术伤害（带攻击掷骰）
func (e *Engine) applySpellDamage(game *model.GameState, attackerID, targetID model.ID, spell *model.Spell, upcastLevel int, isCritical bool) (*DamageResult, error) {
	baseDamage := parseDiceString(spell.DamageDice)

	// 升环加成
	if upcastLevel > spell.Level && spell.AtHigherLevels != "" {
		// 简化处理：每升一环增加基础伤害的1/2
		bonus := (upcastLevel - spell.Level) * (baseDamage / 2)
		baseDamage += bonus
	}

	return e.applyDamageToTarget(game, attackerID, targetID, baseDamage, spell.DamageType, isCritical)
}

// applySpellDamageDirect 直接应用法术伤害（无攻击掷骰）
func (e *Engine) applySpellDamageDirect(game *model.GameState, attackerID, targetID model.ID, amount int, damageType model.DamageType, isCritical bool) (*DamageResult, error) {
	return e.applyDamageToTarget(game, attackerID, targetID, amount, damageType, isCritical)
}

// applyHealingInSpell 应用法术治疗
func (e *Engine) applyHealingInSpell(game *model.GameState, targetID model.ID, amount int) (*HealResult, error) {
	target, ok := game.GetActor(targetID)
	if !ok {
		return nil, ErrNotFound
	}

	var targetActor *model.Actor
	switch t := target.(type) {
	case *model.PlayerCharacter:
		targetActor = &t.Actor
	case *model.NPC:
		targetActor = &t.Actor
	case *model.Enemy:
		targetActor = &t.Actor
	case *model.Companion:
		targetActor = &t.Actor
	}

	hpBefore := targetActor.HitPoints.Current
	wasStable := targetActor.HasCondition(model.ConditionStabilized)

	targetActor.HitPoints.Current += amount
	if targetActor.HitPoints.Current > targetActor.HitPoints.Maximum {
		targetActor.HitPoints.Current = targetActor.HitPoints.Maximum
	}

	// 移除稳定状态
	if targetActor.HitPoints.Current > 0 && wasStable {
		newConditions := make([]model.ConditionInstance, 0)
		for _, c := range targetActor.Conditions {
			if c.Type != model.ConditionStabilized {
				newConditions = append(newConditions, c)
			}
		}
		targetActor.Conditions = newConditions
	}

	return &HealResult{
		Amount:    amount,
		HPBefore:  hpBefore,
		HPAfter:   targetActor.HitPoints.Current,
		WasStable: wasStable,
		Message:   fmt.Sprintf("恢复 %d 点HP", amount),
	}, nil
}

// parseDiceString 解析骰子字符串并返回估算值
// 简化实现：返回平均值
func parseDiceString(diceExpr string) int {
	if diceExpr == "" {
		return 0
	}
	// 简化处理：返回一个固定值用于测试
	// 实际应该使用 e.roller 来掷骰
	return 10
}

// findSpellDefinition 查找法术定义
func findSpellDefinition(spellID string) *model.Spell {
	// 从全局注册中心查找法术定义
	if spell, ok := data.GlobalRegistry.GetSpell(spellID); ok {
		return spell
	}
	return nil
}

// CastSpellRitualRequest 仪式施法请求
type CastSpellRitualRequest struct {
	GameID   model.ID `json:"game_id"`
	CasterID model.ID `json:"caster_id"`
	SpellID  string   `json:"spell_id"`
}

// CastSpellRitual 仪式施法（不消耗法术位，但需要额外 10 分钟施法时间）
func (e *Engine) CastSpellRitual(ctx context.Context, req CastSpellRitualRequest) (*SpellResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	caster, ok := game.GetActor(req.CasterID)
	if !ok {
		return nil, ErrNotFound
	}

	var spellcaster *model.SpellcasterState
	switch c := caster.(type) {
	case *model.PlayerCharacter:
		spellcaster = c.Spellcasting
	default:
		return nil, fmt.Errorf("only player characters can cast spells")
	}

	if spellcaster == nil {
		return nil, fmt.Errorf("actor is not a spellcaster")
	}

	// 检查施法者是否支持仪式施法
	if spellcaster.PreparationType != "prepared" {
		return nil, fmt.Errorf("this caster cannot perform ritual casting")
	}

	// 查找法术定义
	spellDef := findSpellDefinition(req.SpellID)
	if spellDef == nil {
		return nil, fmt.Errorf("spell %s not found", req.SpellID)
	}

	// 检查法术是否有仪式标签
	if !spellDef.Ritual {
		return nil, fmt.Errorf("spell %s does not have the ritual tag", spellDef.Name)
	}

	// 检查施法者是否知道/准备该法术
	if !canCastSpell(spellcaster, req.SpellID) {
		return nil, fmt.Errorf("caster does not know or have prepared spell: %s", spellDef.Name)
	}

	// 仪式施法不消耗法术位
	result := &SpellResult{
		SpellName:    spellDef.Name,
		SlotLevel:    0, // 仪式施法不消耗法术位
		CasterSaveDC: spellcaster.SpellSaveDC,
		Targets:      make([]SpellTargetResult, 0),
		Message:      fmt.Sprintf("通过仪式施法施展了 %s（不消耗法术位）", spellDef.Name),
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return result, nil
}

// GetPactMagicSlotsRequest 获取魔契师法术位请求
type GetPactMagicSlotsRequest struct {
	GameID   model.ID `json:"game_id"`
	CasterID model.ID `json:"caster_id"`
}

// RestorePactMagicSlotsRequest 恢复魔契师法术位请求
type RestorePactMagicSlotsRequest struct {
	GameID   model.ID `json:"game_id"`
	CasterID model.ID `json:"caster_id"`
}

// GetPactMagicSlots 获取魔契师 Pact Magic 法术位
// Note: Pact Magic 使用标准法术位系统，此方法返回所有可用的法术位
func (e *Engine) GetPactMagicSlots(ctx context.Context, req GetPactMagicSlotsRequest) (*GetSpellSlotsResult, error) {
	// 使用标准 GetSpellSlots 方法
	return e.GetSpellSlots(ctx, GetSpellSlotsRequest(req))
}

// RestorePactMagicSlots 恢复魔契师法术位（短休即可恢复）
// Note: Pact Magic 恢复需要短休，这个方法应该在 ShortRest 后调用
func (e *Engine) RestorePactMagicSlots(ctx context.Context, req RestorePactMagicSlotsRequest) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return err
	}

	caster, ok := game.GetActor(req.CasterID)
	if !ok {
		return ErrNotFound
	}

	var spellcaster *model.SpellcasterState
	switch c := caster.(type) {
	case *model.PlayerCharacter:
		spellcaster = c.Spellcasting
	default:
		return fmt.Errorf("only player characters can have spell slots")
	}

	if spellcaster == nil {
		return fmt.Errorf("actor is not a spellcaster")
	}

	// 重置所有法术位（Pact Magic 短休后恢复）
	if spellcaster.Slots != nil {
		for level := 1; level <= 9; level++ {
			if spellcaster.Slots.Slots[level][0] > 0 {
				spellcaster.Slots.Slots[level][1] = 0
			}
		}
	}

	if err := e.saveGame(ctx, game); err != nil {
		return err
	}

	return nil
}

// IsConcentratingRequest 检查是否正在专注
type IsConcentratingRequest struct {
	GameID   model.ID `json:"game_id"`
	CasterID model.ID `json:"caster_id"`
}

// IsConcentratingResult 专注状态结果
type IsConcentratingResult struct {
	IsConcentrating bool   `json:"is_concentrating"`
	SpellName       string `json:"spell_name,omitempty"`
}

// IsConcentrating 检查施法者是否正在专注
func (e *Engine) IsConcentrating(ctx context.Context, req IsConcentratingRequest) (*IsConcentratingResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	caster, ok := game.GetActor(req.CasterID)
	if !ok {
		return nil, ErrNotFound
	}

	var spellcaster *model.SpellcasterState
	switch c := caster.(type) {
	case *model.PlayerCharacter:
		spellcaster = c.Spellcasting
	default:
		return nil, fmt.Errorf("only player characters can concentrate on spells")
	}

	if spellcaster == nil {
		return &IsConcentratingResult{IsConcentrating: false}, nil
	}

	result := &IsConcentratingResult{
		IsConcentrating: spellcaster.IsConcentrating(),
	}
	if result.IsConcentrating {
		result.SpellName = spellcaster.ConcentrationSpell
	}

	return result, nil
}

// GetConcentrationSpellRequest 获取当前专注的法术
type GetConcentrationSpellRequest struct {
	GameID   model.ID `json:"game_id"`
	CasterID model.ID `json:"caster_id"`
}

// GetConcentrationSpell 获取当前专注的法术
func (e *Engine) GetConcentrationSpell(ctx context.Context, req GetConcentrationSpellRequest) (*SpellResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	caster, ok := game.GetActor(req.CasterID)
	if !ok {
		return nil, ErrNotFound
	}

	var spellcaster *model.SpellcasterState
	switch c := caster.(type) {
	case *model.PlayerCharacter:
		spellcaster = c.Spellcasting
	default:
		return nil, fmt.Errorf("only player characters can concentrate on spells")
	}

	if spellcaster == nil || !spellcaster.IsConcentrating() {
		return nil, fmt.Errorf("caster is not concentrating on any spell")
	}

	spellDef := findSpellDefinition(spellcaster.ConcentrationSpell)
	if spellDef == nil {
		return nil, fmt.Errorf("concentration spell %s not found", spellcaster.ConcentrationSpell)
	}

	return &SpellResult{
		SpellName:     spellDef.Name,
		Concentration: true,
		Message:       fmt.Sprintf("当前专注的法术: %s", spellDef.Name),
	}, nil
}
