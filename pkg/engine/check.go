package engine

import (
	"context"
	"fmt"

	"github.com/zwh8800/dnd-core/internal/model"
	"github.com/zwh8800/dnd-core/internal/rules"
)

// AbilityCheckRequest 属性检定请求
type AbilityCheckRequest struct {
	ActorID   model.ID           `json:"actor_id"`
	Ability   model.Ability      `json:"ability"`
	DC        int                `json:"dc,omitempty"`        // 难度等级（可选）
	Advantage model.RollModifier `json:"advantage,omitempty"` // 优势/劣势
	Reason    string             `json:"reason,omitempty"`
}

// AbilityCheckResult 属性检定结果
type AbilityCheckResult struct {
	ActorID      model.ID          `json:"actor_id"`
	ActorName    string            `json:"actor_name"`
	Ability      model.Ability     `json:"ability"`
	AbilityScore int               `json:"ability_score"`
	AbilityMod   int               `json:"ability_modifier"`
	Roll         *model.DiceResult `json:"roll"`
	RollTotal    int               `json:"roll_total"`
	DC           int               `json:"dc,omitempty"`
	Success      bool              `json:"success"`
	Margin       int               `json:"margin"` // 成功/失败的幅度
	Message      string            `json:"message"`
}

// SkillCheckRequest 技能检定请求
type SkillCheckRequest struct {
	ActorID   model.ID           `json:"actor_id"`
	Skill     model.Skill        `json:"skill"`
	DC        int                `json:"dc,omitempty"`
	Advantage model.RollModifier `json:"advantage,omitempty"`
	Reason    string             `json:"reason,omitempty"`
}

// SkillCheckResult 技能检定结果
type SkillCheckResult struct {
	ActorID    model.ID          `json:"actor_id"`
	ActorName  string            `json:"actor_name"`
	Skill      model.Skill       `json:"skill"`
	Ability    model.Ability     `json:"ability"`
	Proficient bool              `json:"proficient"`
	Expertise  bool              `json:"expertise"`
	Roll       *model.DiceResult `json:"roll"`
	RollTotal  int               `json:"roll_total"`
	DC         int               `json:"dc,omitempty"`
	Success    bool              `json:"success"`
	Margin     int               `json:"margin"`
	Message    string            `json:"message"`
}

// SavingThrowRequest 豁免检定请求
type SavingThrowRequest struct {
	ActorID   model.ID           `json:"actor_id"`
	Ability   model.Ability      `json:"ability"`
	DC        int                `json:"dc"`
	Advantage model.RollModifier `json:"advantage,omitempty"`
	Reason    string             `json:"reason,omitempty"`
}

// SavingThrowResult 豁免检定结果
type SavingThrowResult struct {
	ActorID    model.ID          `json:"actor_id"`
	ActorName  string            `json:"actor_name"`
	Ability    model.Ability     `json:"ability"`
	Proficient bool              `json:"proficient"`
	Roll       *model.DiceResult `json:"roll"`
	RollTotal  int               `json:"roll_total"`
	DC         int               `json:"dc"`
	Success    bool              `json:"success"`
	Margin     int               `json:"margin"`
	Message    string            `json:"message"`
}

// PerformAbilityCheck 执行属性检定
func (e *Engine) PerformAbilityCheck(ctx context.Context, gameID model.ID, req AbilityCheckRequest) (*AbilityCheckResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	actor, ok := game.GetActor(req.ActorID)
	if !ok {
		return nil, ErrNotFound
	}

	var baseActor *model.Actor
	var name string
	switch a := actor.(type) {
	case *model.PlayerCharacter:
		baseActor = &a.Actor
		name = a.Name
	case *model.NPC:
		baseActor = &a.Actor
		name = a.Name
	case *model.Enemy:
		baseActor = &a.Actor
		name = a.Name
	case *model.Companion:
		baseActor = &a.Actor
		name = a.Name
	}

	// 获取属性值和修正
	abilityScore := baseActor.AbilityScores.Get(req.Ability)
	abilityMod := rules.AbilityModifier(abilityScore)

	// 掷骰
	var rollResult *model.DiceResult
	if req.Advantage.Advantage {
		rollResult, _ = e.roller.RollAdvantage(abilityMod)
	} else if req.Advantage.Disadvantage {
		rollResult, _ = e.roller.RollDisadvantage(abilityMod)
	} else {
		rollResult, _ = e.roller.Roll("1d20")
		rollResult.Total += abilityMod
	}

	// 判断是否成功
	success := false
	margin := 0
	if req.DC > 0 {
		success = rollResult.Total >= req.DC
		margin = rollResult.Total - req.DC
	}

	message := fmt.Sprintf("%s 进行%s检定: %d", name, req.Ability, rollResult.Total)
	if req.DC > 0 {
		if success {
			message += fmt.Sprintf(" (成功，超出 %d)", margin)
		} else {
			message += fmt.Sprintf(" (失败，差 %d)", -margin)
		}
	}

	return &AbilityCheckResult{
		ActorID:      req.ActorID,
		ActorName:    name,
		Ability:      req.Ability,
		AbilityScore: abilityScore,
		AbilityMod:   abilityMod,
		Roll:         rollResult,
		RollTotal:    rollResult.Total,
		DC:           req.DC,
		Success:      success,
		Margin:       margin,
		Message:      message,
	}, nil
}

// PerformSkillCheck 执行技能检定
func (e *Engine) PerformSkillCheck(ctx context.Context, gameID model.ID, req SkillCheckRequest) (*SkillCheckResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	actor, ok := game.GetActor(req.ActorID)
	if !ok {
		return nil, ErrNotFound
	}

	var baseActor *model.Actor
	var name string
	var level int
	switch a := actor.(type) {
	case *model.PlayerCharacter:
		baseActor = &a.Actor
		name = a.Name
		level = a.TotalLevel
	case *model.NPC:
		baseActor = &a.Actor
		name = a.Name
		level = 1
	case *model.Enemy:
		baseActor = &a.Actor
		name = a.Name
		level = 1
	case *model.Companion:
		baseActor = &a.Actor
		name = a.Name
		level = 1
	}

	// 获取技能对应的属性
	ability := model.SkillAbilityMap[req.Skill]
	abilityScore := baseActor.AbilityScores.Get(ability)
	abilityMod := rules.AbilityModifier(abilityScore)

	// 检查是否熟练
	proficient := baseActor.Proficiencies.IsProficient(req.Skill)
	expertise := baseActor.Proficiencies.HasExpertise(req.Skill)

	// 计算总加值
	bonus := abilityMod
	if proficient {
		profBonus := rules.ProficiencyBonus(level)
		if expertise {
			profBonus *= 2
		}
		bonus += profBonus
	}

	// 掷骰
	var rollResult *model.DiceResult
	if req.Advantage.Advantage {
		rollResult, _ = e.roller.RollAdvantage(bonus)
	} else if req.Advantage.Disadvantage {
		rollResult, _ = e.roller.RollDisadvantage(bonus)
	} else {
		rollResult, _ = e.roller.Roll("1d20")
		rollResult.Total += bonus
	}

	// 判断是否成功
	success := false
	margin := 0
	if req.DC > 0 {
		success = rollResult.Total >= req.DC
		margin = rollResult.Total - req.DC
	}

	message := fmt.Sprintf("%s 进行%s检定: %d", name, req.Skill, rollResult.Total)
	if req.DC > 0 {
		if success {
			message += fmt.Sprintf(" (成功，超出 %d)", margin)
		} else {
			message += fmt.Sprintf(" (失败，差 %d)", -margin)
		}
	}

	return &SkillCheckResult{
		ActorID:    req.ActorID,
		ActorName:  name,
		Skill:      req.Skill,
		Ability:    ability,
		Proficient: proficient,
		Expertise:  expertise,
		Roll:       rollResult,
		RollTotal:  rollResult.Total,
		DC:         req.DC,
		Success:    success,
		Margin:     margin,
		Message:    message,
	}, nil
}

// PerformSavingThrow 执行豁免检定
func (e *Engine) PerformSavingThrow(ctx context.Context, gameID model.ID, req SavingThrowRequest) (*SavingThrowResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	actor, ok := game.GetActor(req.ActorID)
	if !ok {
		return nil, ErrNotFound
	}

	var baseActor *model.Actor
	var name string
	var level int
	switch a := actor.(type) {
	case *model.PlayerCharacter:
		baseActor = &a.Actor
		name = a.Name
		level = a.TotalLevel
	case *model.NPC:
		baseActor = &a.Actor
		name = a.Name
		level = 1
	case *model.Enemy:
		baseActor = &a.Actor
		name = a.Name
		level = 1
	case *model.Companion:
		baseActor = &a.Actor
		name = a.Name
		level = 1
	}

	// 获取属性修正
	abilityScore := baseActor.AbilityScores.Get(req.Ability)
	abilityMod := rules.AbilityModifier(abilityScore)

	// 检查是否豁免熟练
	proficient := baseActor.Proficiencies.IsSavingThrowProficient(req.Ability)

	// 计算总加值
	bonus := abilityMod
	if proficient {
		bonus += rules.ProficiencyBonus(level)
	}

	// 掷骰
	var rollResult *model.DiceResult
	if req.Advantage.Advantage {
		rollResult, _ = e.roller.RollAdvantage(bonus)
	} else if req.Advantage.Disadvantage {
		rollResult, _ = e.roller.RollDisadvantage(bonus)
	} else {
		rollResult, _ = e.roller.Roll("1d20")
		rollResult.Total += bonus
	}

	// 判断是否成功
	success := rollResult.Total >= req.DC
	margin := rollResult.Total - req.DC

	message := fmt.Sprintf("%s 进行%s豁免: %d vs DC %d", name, req.Ability, rollResult.Total, req.DC)
	if success {
		message += fmt.Sprintf(" (成功，超出 %d)", margin)
	} else {
		message += fmt.Sprintf(" (失败，差 %d)", -margin)
	}

	return &SavingThrowResult{
		ActorID:    req.ActorID,
		ActorName:  name,
		Ability:    req.Ability,
		Proficient: proficient,
		Roll:       rollResult,
		RollTotal:  rollResult.Total,
		DC:         req.DC,
		Success:    success,
		Margin:     margin,
		Message:    message,
	}, nil
}

// GetSkillAbility 获取技能对应的属性
func (e *Engine) GetSkillAbility(skill model.Skill) model.Ability {
	return model.SkillAbilityMap[skill]
}

// GetPassivePerception 获取被动感知（察觉）
func (e *Engine) GetPassivePerception(ctx context.Context, gameID model.ID, actorID model.ID) (int, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return 0, err
	}

	actor, ok := game.GetActor(actorID)
	if !ok {
		return 0, ErrNotFound
	}

	var baseActor *model.Actor
	var level int
	switch a := actor.(type) {
	case *model.PlayerCharacter:
		baseActor = &a.Actor
		level = a.TotalLevel
	case *model.NPC:
		baseActor = &a.Actor
		level = 1
	case *model.Enemy:
		baseActor = &a.Actor
		level = 1
	case *model.Companion:
		baseActor = &a.Actor
		level = 1
	}

	// 被动感知 = 10 + 感知修正 + (如果熟练则加熟练加值)
	wisMod := rules.AbilityModifier(baseActor.AbilityScores.Wisdom)
	passive := 10 + wisMod

	if baseActor.Proficiencies.IsProficient(model.SkillPerception) {
		passive += rules.ProficiencyBonus(level)
	}

	return passive, nil
}
