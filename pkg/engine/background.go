package engine

import (
	"context"
	"fmt"

	"github.com/zwh8800/dnd-core/pkg/data"
	"github.com/zwh8800/dnd-core/pkg/model"
	"github.com/zwh8800/dnd-core/pkg/rules"
)

// ApplyBackgroundRequest 应用背景请求
type ApplyBackgroundRequest struct {
	GameID       model.ID           `json:"game_id"`       // 游戏会话ID（必填）
	PCID         model.ID           `json:"pc_id"`         // 玩家角色ID（必填）
	BackgroundID model.BackgroundID `json:"background_id"` // 背景ID（必填）
}

// ApplyBackgroundResult 应用背景结果
type ApplyBackgroundResult struct {
	BackgroundID          string   `json:"background_id"`             // 背景ID
	BackgroundName        string   `json:"background_name"`           // 背景名称
	SkillProficiencies    []string `json:"skill_proficiencies"`       // 获得的技能熟练
	ToolProficiencies     []string `json:"tool_proficiencies"`        // 获得的工具熟练
	LanguageProficiencies []string `json:"language_proficiencies"`    // 获得的语言熟练
	Features              []string `json:"features"`                  // 获得的特性
	AssociatedFeat        string   `json:"associated_feat,omitempty"` // 关联专长（如果有）
	Message               string   `json:"message"`                   // 人类可读消息
}

// GetBackgroundFeaturesRequest 获取背景特性请求
type GetBackgroundFeaturesRequest struct {
	GameID model.ID `json:"game_id"` // 游戏会话ID（必填）
	PCID   model.ID `json:"pc_id"`   // 玩家角色ID（必填）
}

// GetBackgroundFeaturesResult 获取背景特性结果
type GetBackgroundFeaturesResult struct {
	BackgroundID   string   `json:"background_id"`   // 背景ID
	BackgroundName string   `json:"background_name"` // 背景名称
	Features       []string `json:"features"`        // 背景特性列表
	Description    string   `json:"description"`     // 背景描述
}

// ApplyBackground 为背景应用到角色
// 应用背景会给予角色技能熟练、工具熟练、语言熟练、特性以及可能的关联专长。
// 如果角色已有背景，会被新背景覆盖。
// 参数:
//
//	ctx - 上下文
//	req - 应用背景请求，包含游戏会话ID、角色ID和背景ID
//
// 返回:
//
//	*ApplyBackgroundResult - 应用结果，包含给予的熟练和特性
//	error - 可能返回 ErrNotFound（角色或背景不存在）、游戏不存在或权限错误
func (e *Engine) ApplyBackground(ctx context.Context, req ApplyBackgroundRequest) (*ApplyBackgroundResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpApplyBackground); err != nil {
		return nil, err
	}

	pc, ok := game.PCs[req.PCID]
	if !ok {
		return nil, ErrNotFound
	}

	// 验证背景是否存在
	bg, exists := data.GlobalRegistry.GetBackground(string(req.BackgroundID))
	if !exists {
		return nil, fmt.Errorf("background not found: %s", req.BackgroundID)
	}

	// 调用规则层应用背景
	err = rules.ApplyBackground(pc, req.BackgroundID)
	if err != nil {
		return nil, fmt.Errorf("failed to apply background: %w", err)
	}

	// 收集应用的效果信息
	skillProfs := make([]string, 0)
	for skill := range pc.Proficiencies.ProficientSkills {
		// 简化：返回所有技能熟练，实际应只返回新获得的
		skillProfs = append(skillProfs, string(skill))
	}

	toolProfs := make([]string, 0)
	for tool := range pc.Proficiencies.ToolProficiencies {
		toolProfs = append(toolProfs, tool)
	}

	langProfs := make([]string, 0)
	for lang := range pc.Proficiencies.LanguageProficiencies {
		langProfs = append(langProfs, lang)
	}

	features := make([]string, len(pc.Features))
	copy(features, pc.Features)

	associatedFeat := ""
	if len(pc.Feats) > 0 {
		// 获取最新添加的专长（如果背景给予了专长）
		latestFeat := pc.Feats[len(pc.Feats)-1]
		if latestFeat.Source == model.FeatSourceBackground {
			associatedFeat = latestFeat.FeatID
		}
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	message := fmt.Sprintf("已应用背景：%s", bg.Name)
	if associatedFeat != "" {
		message += fmt.Sprintf("，获得专长：%s", associatedFeat)
	}

	return &ApplyBackgroundResult{
		BackgroundID:          string(req.BackgroundID),
		BackgroundName:        bg.Name,
		SkillProficiencies:    skillProfs,
		ToolProficiencies:     toolProfs,
		LanguageProficiencies: langProfs,
		Features:              features,
		AssociatedFeat:        associatedFeat,
		Message:               message,
	}, nil
}

// GetBackgroundFeatures 获取角色背景的特性
// 返回角色当前背景的特性列表和描述信息。
// 参数:
//
//	ctx - 上下文
//	req - 获取请求，包含游戏会话ID和角色ID
//
// 返回:
//
//	*GetBackgroundFeaturesResult - 背景特性信息
//	error - 可能返回 ErrNotFound（角色不存在）、游戏不存在或权限错误
func (e *Engine) GetBackgroundFeatures(ctx context.Context, req GetBackgroundFeaturesRequest) (*GetBackgroundFeaturesResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	// 背景特性查询不需要特殊权限，所有阶段允许
	// 但需要验证游戏存在

	pc, ok := game.PCs[req.PCID]
	if !ok {
		return nil, ErrNotFound
	}

	if pc.BackgroundID == "" {
		return &GetBackgroundFeaturesResult{
			Features: nil,
		}, nil
	}

	bg, exists := data.GlobalRegistry.GetBackground(pc.BackgroundID)
	if !exists {
		return &GetBackgroundFeaturesResult{
			BackgroundID:   pc.BackgroundID,
			BackgroundName: "未知背景",
			Features:       nil,
		}, nil
	}

	features := rules.GetBackgroundFeatures(pc)
	featureNames := make([]string, len(features))
	copy(featureNames, features)

	return &GetBackgroundFeaturesResult{
		BackgroundID:   pc.BackgroundID,
		BackgroundName: bg.Name,
		Features:       featureNames,
		Description:    bg.FeatureDescription,
	}, nil
}
