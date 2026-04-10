package engine

import (
	"context"
	"fmt"

	"github.com/zwh8800/dnd-core/pkg/data"
	"github.com/zwh8800/dnd-core/pkg/model"
	"github.com/zwh8800/dnd-core/pkg/rules"
)

// ValidateMulticlassRequest 多职业验证请求
type ValidateMulticlassRequest struct {
	GameID   model.ID      `json:"game_id"`   // 游戏会话ID（必填）
	PCID     model.ID      `json:"pc_id"`     // 玩家角色ID（必填）
	NewClass model.ClassID `json:"new_class"` // 要验证的新职业ID（必填）
}

// ValidateMulticlassResult 多职业验证结果
type ValidateMulticlassResult struct {
	Valid             bool   `json:"valid"`                     // 是否合法
	MeetsRequirements bool   `json:"meets_requirements"`        // 是否满足属性要求
	Message           string `json:"message"`                   // 验证消息
	PrimaryAbility    string `json:"primary_ability,omitempty"` // 主属性
	RequiredScore     int    `json:"required_score"`            // 需要的属性分数
	CurrentScore      int    `json:"current_score"`             // 当前属性分数
}

// GetMulticlassSpellSlotsRequest 获取多职业法术位请求
type GetMulticlassSpellSlotsRequest struct {
	GameID model.ID `json:"game_id"` // 游戏会话ID（必填）
	PCID   model.ID `json:"pc_id"`   // 玩家角色ID（必填）
}

// SpellSlotInfo 法术位信息
type SpellSlotInfo struct {
	Level int `json:"level"` // 法术环级
	Total int `json:"total"` // 总法术位数
	Used  int `json:"used"`  // 已用法术位数
}

// GetMulticlassSpellSlotsResult 获取多职业法术位结果
type GetMulticlassSpellSlotsResult struct {
	EffectiveCasterLevel int             `json:"effective_caster_level"` // 有效施法者等级
	SpellSlots           []SpellSlotInfo `json:"spell_slots"`            // 法术位表
	IsMulticlassCaster   bool            `json:"is_multiclass_caster"`   // 是否为多职业施法者
	Message              string          `json:"message"`                // 人类可读消息
}

// ValidateMulticlassChoice 验证多职业选择是否合法
// 检查角色是否满足多职业要求（主要属性至少13分），包括新职业和已有职业的要求。
// 参数:
//
//	ctx - 上下文
//	req - 验证请求，包含游戏会话ID、角色ID和要验证的新职业
//
// 返回:
//
//	*ValidateMulticlassResult - 验证结果，包含是否合法、属性要求等信息
//	error - 可能返回 ErrNotFound（角色或职业不存在）、游戏不存在或权限错误
func (e *Engine) ValidateMulticlassChoice(ctx context.Context, req ValidateMulticlassRequest) (*ValidateMulticlassResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpValidateMulticlass); err != nil {
		return nil, err
	}

	pc, ok := game.PCs[req.PCID]
	if !ok {
		return nil, ErrNotFound
	}

	// 验证职业定义是否存在
	classDef, exists := data.GlobalRegistry.GetClass(req.NewClass)
	if !exists {
		return nil, fmt.Errorf("class definition not found: %s", req.NewClass)
	}

	// 调用规则层验证
	err = rules.ValidateMulticlass(pc, req.NewClass)
	if err != nil {
		// 获取主属性信息用于返回详细消息
		primaryAbility := ""
		requiredScore := 13
		currentScore := 0

		if len(classDef.PrimaryAbilities) > 0 {
			primaryAbility = string(classDef.PrimaryAbilities[0])
			currentScore = pc.AbilityScores.Get(classDef.PrimaryAbilities[0])
		}

		return &ValidateMulticlassResult{
			Valid:             false,
			MeetsRequirements: false,
			Message:           err.Error(),
			PrimaryAbility:    primaryAbility,
			RequiredScore:     requiredScore,
			CurrentScore:      currentScore,
		}, nil
	}

	// 验证通过
	primaryAbility := ""
	currentScore := 0
	if len(classDef.PrimaryAbilities) > 0 {
		primaryAbility = string(classDef.PrimaryAbilities[0])
		currentScore = pc.AbilityScores.Get(classDef.PrimaryAbilities[0])
	}

	return &ValidateMulticlassResult{
		Valid:             true,
		MeetsRequirements: true,
		Message:           fmt.Sprintf("满足 %s 的多职业要求", classDef.Name),
		PrimaryAbility:    primaryAbility,
		RequiredScore:     13,
		CurrentScore:      currentScore,
	}, nil
}

// GetMulticlassSpellSlots 计算并返回多职业施法者的法术位表
// 对于单职业施法者，返回该职业的法术位；对于多职业组合，计算有效施法者等级后返回对应法术位。
// 参数:
//
//	ctx - 上下文
//	req - 请求，包含游戏会话ID和角色ID
//
// 返回:
//
//	*GetMulticlassSpellSlotsResult - 法术位信息，包含有效施法者等级和法术位表
//	error - 可能返回 ErrNotFound（角色不存在）、游戏不存在或权限错误
func (e *Engine) GetMulticlassSpellSlots(ctx context.Context, req GetMulticlassSpellSlotsRequest) (*GetMulticlassSpellSlotsResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpGetMulticlassSlots); err != nil {
		return nil, err
	}

	pc, ok := game.PCs[req.PCID]
	if !ok {
		return nil, ErrNotFound
	}

	// 检查是否有多个职业
	isMulticlass := len(pc.Classes) > 1

	var slots [][]int
	if isMulticlass {
		// 多职业：使用多职业法术位计算
		slots = rules.GetMulticlassSpellSlots(pc.Classes)
	} else if len(pc.Classes) == 1 {
		// 单职业：使用该职业的法术位
		classDef, exists := data.GlobalRegistry.GetClass(pc.Classes[0].Class)
		if !exists || classDef.CasterType == model.CasterTypeNone {
			return &GetMulticlassSpellSlotsResult{
				EffectiveCasterLevel: 0,
				SpellSlots:           nil,
				IsMulticlassCaster:   false,
				Message:              "该角色不是施法者",
			}, nil
		}
		slots = rules.GetSpellSlotTable(pc.Classes[0].Level)
	} else {
		return &GetMulticlassSpellSlotsResult{
			EffectiveCasterLevel: 0,
			SpellSlots:           nil,
			IsMulticlassCaster:   false,
			Message:              "该角色没有职业",
		}, nil
	}

	if slots == nil {
		return &GetMulticlassSpellSlotsResult{
			EffectiveCasterLevel: 0,
			SpellSlots:           nil,
			IsMulticlassCaster:   false,
			Message:              "该角色没有法术位",
		}, nil
	}

	// 转换为信息结构体
	spellSlotInfo := make([]SpellSlotInfo, len(slots))
	for i, slot := range slots {
		level := i + 1 // 索引0对应1环
		total := 0
		used := 0
		if len(slot) > 0 {
			total = slot[0]
		}
		if len(slot) > 1 {
			used = slot[1]
		}
		spellSlotInfo[i] = SpellSlotInfo{
			Level: level,
			Total: total,
			Used:  used,
		}
	}

	// 计算有效施法者等级
	effectiveCasterLevel := 0
	for _, classLevel := range pc.Classes {
		classDef, exists := data.GlobalRegistry.GetClass(classLevel.Class)
		if !exists {
			continue
		}
		switch classDef.CasterType {
		case model.CasterTypeFull:
			effectiveCasterLevel += classLevel.Level
		case model.CasterTypeHalf:
			effectiveCasterLevel += classLevel.Level / 2
		case model.CasterTypeThird:
			effectiveCasterLevel += classLevel.Level / 3
		}
	}

	message := fmt.Sprintf("有效施法者等级: %d", effectiveCasterLevel)
	if isMulticlass {
		message = fmt.Sprintf("多职业施法者，有效等级: %d", effectiveCasterLevel)
	}

	return &GetMulticlassSpellSlotsResult{
		EffectiveCasterLevel: effectiveCasterLevel,
		SpellSlots:           spellSlotInfo,
		IsMulticlassCaster:   isMulticlass,
		Message:              message,
	}, nil
}
