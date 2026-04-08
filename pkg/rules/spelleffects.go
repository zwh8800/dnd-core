package rules

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/zwh8800/dnd-core/pkg/dice"
	"github.com/zwh8800/dnd-core/pkg/model"
)

// ExecuteSpellResult 执行法术的结果
type ExecuteSpellResult struct {
	SpellName     string              `json:"spell_name"`
	Level         int                 `json:"level"`
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
	Damage      int               `json:"damage,omitempty"`
	Healing     int               `json:"healing,omitempty"`
	Effect      string            `json:"effect,omitempty"`
}

// ExecuteSpell 执行法术效果
func ExecuteSpell(spellDef *model.SpellDefinition, casterAbilityScore int, casterLevel int, targets []model.ID, gameState *model.GameState) (*ExecuteSpellResult, error) {
	result := &ExecuteSpellResult{
		SpellName:     spellDef.Name,
		Level:         spellDef.Level,
		Targets:       make([]SpellTargetResult, 0),
		Concentration: spellDef.Concentration,
	}

	// 计算施法豁免DC
	profBonus := ProficiencyBonus(casterLevel)
	result.CasterSaveDC = 8 + profBonus + AbilityModifier(casterAbilityScore)

	// 对每个效果进行处理
	for _, effect := range spellDef.Effects {
		switch effect.Type {
		case model.SpellEffectDamage:
			handleDamageEffect(spellDef, effect, casterLevel, targets, gameState, result)
		case model.SpellEffectHealing:
			handleHealingEffect(spellDef, effect, casterAbilityScore, targets, gameState, result)
		case model.SpellEffectCondition:
			handleConditionEffect(spellDef, effect, targets, gameState, result)
		case model.SpellEffectBuff, model.SpellEffectDebuff:
			handleBuffDebuffEffect(spellDef, effect, targets, gameState, result)
		case model.SpellEffectTeleport:
			handleTeleportEffect(spellDef, effect, targets, gameState, result)
		case model.SpellEffectUtility:
			handleUtilityEffect(spellDef, effect, targets, gameState, result)
		}
	}

	result.Message = fmt.Sprintf("成功施放 %s", spellDef.Name)
	return result, nil
}

// handleDamageEffect 处理伤害类法术
func handleDamageEffect(spellDef *model.SpellDefinition, effect model.SpellEffect, casterLevel int, targets []model.ID, gameState *model.GameState, result *ExecuteSpellResult) {
	if effect.Damage == nil {
		return
	}

	damageDice := effect.Damage.BaseDice

	// 处理升环伤害
	if spellDef.Level > 0 && effect.Damage.UpcastDicePerLevel != "" && effect.Damage.UpcastStartLevel > 0 {
		// 计算升环增加的骰子数
		// 例如火球术: 3环8d6, 每升一环增加1d6
		// 如果 spellDef.Level 是当前施放环级
		upcastLevels := spellDef.Level - effect.Damage.UpcastStartLevel
		if upcastLevels > 0 {
			// 解析基础骰子数量
			baseDiceCount := parseDiceCount(damageDice)
			// 每升一级增加的骰子数
			upcastDiceCount := parseDiceCount(effect.Damage.UpcastDicePerLevel)
			// 总骰子数 = 基础 + 升环增加
			totalDiceCount := baseDiceCount + (upcastDiceCount * upcastLevels)
			// 重新构建骰子表达式
			damageDice = fmt.Sprintf("%dd6", totalDiceCount)
		}
	}

	roller := dice.New(0)

	for _, targetID := range targets {
		targetResult := SpellTargetResult{
			ActorID: targetID,
		}

		// 如果需要豁免
		if effect.SaveAbility != "" {
			// 找到目标获取其属性
			actor, ok := gameState.GetActor(targetID)
			if !ok {
				continue
			}

			var abilityScore int
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
			default:
				continue
			}

			switch effect.SaveAbility {
			case model.AbilityStrength:
				abilityScore = baseActor.AbilityScores.Strength
			case model.AbilityDexterity:
				abilityScore = baseActor.AbilityScores.Dexterity
			case model.AbilityConstitution:
				abilityScore = baseActor.AbilityScores.Constitution
			case model.AbilityIntelligence:
				abilityScore = baseActor.AbilityScores.Intelligence
			case model.AbilityWisdom:
				abilityScore = baseActor.AbilityScores.Wisdom
			case model.AbilityCharisma:
				abilityScore = baseActor.AbilityScores.Charisma
			}

			saveMod := AbilityModifier(abilityScore)
			saveRoll, _ := roller.Roll("d20")
			saveTotal := saveRoll.Total + saveMod

			targetResult.SaveRoll = saveRoll
			targetResult.SaveTotal = saveTotal
			targetResult.SaveSuccess = saveTotal >= result.CasterSaveDC

			// 根据豁免结果计算伤害
			if targetResult.SaveSuccess {
				if effect.SaveSuccessEffect == "half" {
					// 豁免成功，伤害减半
					baseDmg := parseDice(damageDice)
					targetResult.Damage = baseDmg / 2
				} else if effect.SaveSuccessEffect == "none" {
					targetResult.Damage = 0
				}
			} else {
				// 豁免失败，全额伤害
				targetResult.Damage = parseDice(damageDice)
			}
		} else if spellDef.RequiresAttackRoll {
			// 需要攻击掷骰
			// 这里简化处理，实际应该在 CastSpell 中处理
		} else {
			// 自动命中
			targetResult.Damage = parseDice(damageDice)
		}

		result.Targets = append(result.Targets, targetResult)
	}
}

// handleHealingEffect 处理治疗类法术
func handleHealingEffect(spellDef *model.SpellDefinition, effect model.SpellEffect, casterAbilityScore int, targets []model.ID, gameState *model.GameState, result *ExecuteSpellResult) {
	if effect.HealingDice == "" {
		return
	}

	healingMod := AbilityModifier(casterAbilityScore)

	for _, targetID := range targets {
		actor, ok := gameState.GetActor(targetID)
		if !ok {
			continue
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
		default:
			continue
		}

		healing := parseDice(effect.HealingDice) + healingMod
		baseActor.HitPoints.Current += healing
		if baseActor.HitPoints.Current > baseActor.HitPoints.Maximum {
			baseActor.HitPoints.Current = baseActor.HitPoints.Maximum
		}

		result.Targets = append(result.Targets, SpellTargetResult{
			ActorID: targetID,
			Healing: healing,
			Effect:  fmt.Sprintf("恢复 %d HP", healing),
		})
	}
}

// handleConditionEffect 处理状态施加类法术
func handleConditionEffect(spellDef *model.SpellDefinition, effect model.SpellEffect, targets []model.ID, gameState *model.GameState, result *ExecuteSpellResult) {
	for _, targetID := range targets {
		actor, ok := gameState.GetActor(targetID)
		if !ok {
			continue
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
		default:
			continue
		}

		saveSuccess := false
		if effect.SaveAbility != "" {
			// 进行豁免检定
			abilityScore := getAbilityScore(baseActor, effect.SaveAbility)
			saveMod := AbilityModifier(abilityScore)
			roller := dice.New(0)
			saveRoll, _ := roller.Roll("d20")
			saveTotal := saveRoll.Total + saveMod

			saveSuccess = saveTotal >= result.CasterSaveDC
		}

		if !saveSuccess {
			// 施加状态
			condition := model.ConditionInstance{
				Type:     effect.ConditionApplied,
				Duration: effect.ConditionDuration,
			}
			baseActor.Conditions = append(baseActor.Conditions, condition)
		}

		result.Targets = append(result.Targets, SpellTargetResult{
			ActorID:     targetID,
			SaveSuccess: saveSuccess,
			Effect:      fmt.Sprintf("豁免成功：%v", saveSuccess),
		})
	}
}

// handleBuffDebuffEffect 处理增益/减益效果
func handleBuffDebuffEffect(spellDef *model.SpellDefinition, effect model.SpellEffect, targets []model.ID, gameState *model.GameState, result *ExecuteSpellResult) {
	for _, targetID := range targets {
		result.Targets = append(result.Targets, SpellTargetResult{
			ActorID: targetID,
			Effect:  effect.Description,
		})
	}
}

// handleTeleportEffect 处理传送效果
func handleTeleportEffect(spellDef *model.SpellDefinition, effect model.SpellEffect, targets []model.ID, gameState *model.GameState, result *ExecuteSpellResult) {
	for _, targetID := range targets {
		result.Targets = append(result.Targets, SpellTargetResult{
			ActorID: targetID,
			Effect:  effect.Description,
		})
	}
}

// handleUtilityEffect 处理实用效果
func handleUtilityEffect(spellDef *model.SpellDefinition, effect model.SpellEffect, targets []model.ID, gameState *model.GameState, result *ExecuteSpellResult) {
	for _, targetID := range targets {
		result.Targets = append(result.Targets, SpellTargetResult{
			ActorID: targetID,
			Effect:  effect.Description,
		})
	}
}

// ResolveSpellSave 处理法术豁免检定
func ResolveSpellSave(target *model.Actor, saveAbility model.Ability, dc int) (*model.DiceResult, int, bool) {
	abilityScore := getAbilityScore(target, saveAbility)
	saveMod := AbilityModifier(abilityScore)

	roller := dice.New(0)
	saveRoll, _ := roller.Roll("d20")
	saveTotal := saveRoll.Total + saveMod

	return saveRoll, saveTotal, saveTotal >= dc
}

// CalculateSpellDamage 计算法术伤害（含升环升级）
func CalculateSpellDamage(baseDice string, slotLevel int, spellLevel int, upcastDicePerLevel string) int {
	damage := parseDice(baseDice)

	// 如果升环施法，增加额外伤害
	if slotLevel > spellLevel && upcastDicePerLevel != "" {
		extraLevels := slotLevel - spellLevel
		extraDamage := parseDice(upcastDicePerLevel) * extraLevels
		damage += extraDamage
	}

	return damage
}

// parseDice 解析骰子表达式并掷骰（如 "2d6+3"）
func parseDice(diceExpr string) int {
	if diceExpr == "" {
		return 0
	}

	roller := dice.New(0)
	result, _ := roller.Roll(diceExpr)
	return result.Total
}

// parseDiceCount 解析骰子表达式中的骰子数量 (如 "8d6" -> 8)
func parseDiceCount(diceExpr string) int {
	if diceExpr == "" {
		return 0
	}

	parts := strings.Split(diceExpr, "d")
	if len(parts) != 2 {
		return 0
	}

	count, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0
	}

	return count
}

// getAbilityScore 获取演员的指定属性值
func getAbilityScore(actor *model.Actor, ability model.Ability) int {
	switch ability {
	case model.AbilityStrength:
		return actor.AbilityScores.Strength
	case model.AbilityDexterity:
		return actor.AbilityScores.Dexterity
	case model.AbilityConstitution:
		return actor.AbilityScores.Constitution
	case model.AbilityIntelligence:
		return actor.AbilityScores.Intelligence
	case model.AbilityWisdom:
		return actor.AbilityScores.Wisdom
	case model.AbilityCharisma:
		return actor.AbilityScores.Charisma
	default:
		return 10
	}
}

// NewRoller 创建骰子投掷器（避免导入循环）

// 确保导入 strings 被使用
func init() {
	_ = strings.Contains("", "")
	_ = strconv.Itoa(0)
}
