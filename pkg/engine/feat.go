package engine

import (
	"context"
	"fmt"

	"github.com/zwh8800/dnd-core/pkg/data"
	"github.com/zwh8800/dnd-core/pkg/model"
	"github.com/zwh8800/dnd-core/pkg/rules"
)

// SelectFeatRequest 选择专长请求
type SelectFeatRequest struct {
	GameID model.ID `json:"game_id"` // 游戏会话ID
	PCID   model.ID `json:"pc_id"`   // 玩家角色ID
	FeatID string   `json:"feat_id"` // 专长ID
}

// SelectFeatResult 选择专长结果
type SelectFeatResult struct {
	Feats []FeatInfo `json:"feats"` // 角色当前专长列表
}

// ListFeatsRequest 列出专长请求
type ListFeatsRequest struct {
	FilterType *model.FeatType `json:"filter_type,omitempty"` // 专长类型过滤
}

// ListFeatsResult 列出专长结果
type ListFeatsResult struct {
	Feats []FeatInfo `json:"feats"` // 专长列表
}

// GetFeatDetailsRequest 获取专长详情请求
type GetFeatDetailsRequest struct {
	FeatID string `json:"feat_id"` // 专长ID
}

// GetFeatDetailsResult 获取专长详情结果
type GetFeatDetailsResult struct {
	Feat *FeatInfo `json:"feat"` // 专长详情
}

// RemoveFeatRequest 移除专长请求
type RemoveFeatRequest struct {
	GameID model.ID `json:"game_id"` // 游戏会话ID
	PCID   model.ID `json:"pc_id"`   // 玩家角色ID
	FeatID string   `json:"feat_id"` // 专长ID
}

// FeatInfo 专长信息
type FeatInfo struct {
	ID           string `json:"id"`                     // 专长ID
	Name         string `json:"name"`                   // 专长名称
	Type         string `json:"type"`                   // 专长类型
	Description  string `json:"description"`            // 专长描述
	Repeatable   bool   `json:"repeatable"`             // 是否可重复选择
	Prerequisite string `json:"prerequisite,omitempty"` // 先决条件
}

// SelectFeat 为角色选择并获得一个专长
func (e *Engine) SelectFeat(ctx context.Context, req SelectFeatRequest) (*SelectFeatResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if err := e.checkPermission(game.Phase, OpSelectFeat); err != nil {
		return nil, err
	}

	pc, ok := game.PCs[req.PCID]
	if !ok {
		return nil, ErrNotFound
	}

	// 验证先决条件
	if !rules.CheckFeatPrerequisites(pc, req.FeatID) {
		return nil, fmt.Errorf("feat prerequisites not met: %s", req.FeatID)
	}

	// 检查是否已拥有该专长（且不可重复）
	feat, exists := data.GlobalRegistry.GetFeat(req.FeatID)
	if !exists {
		return nil, fmt.Errorf("feat not found: %s", req.FeatID)
	}

	if !feat.Repeatable {
		for _, featInstance := range pc.Feats {
			if featInstance.FeatID == req.FeatID {
				return nil, fmt.Errorf("feat already possessed: %s", req.FeatID)
			}
		}
	}

	// 添加专长实例
	featInstance := model.FeatInstance{
		FeatID:        req.FeatID,
		Source:        model.FeatSourceLevelUp,
		AcquiredLevel: pc.TotalLevel,
	}
	pc.Feats = append(pc.Feats, featInstance)

	// 应用专长效果
	rules.ApplyFeatEffects(pc, req.FeatID)

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	// 返回当前专长列表
	feats := make([]FeatInfo, len(pc.Feats))
	for i, instance := range pc.Feats {
		if featDef, exists := data.GlobalRegistry.GetFeat(instance.FeatID); exists {
			feats[i] = featToInfo(featDef)
		}
	}

	return &SelectFeatResult{
		Feats: feats,
	}, nil
}

// ListFeats 列出可选专长
func (e *Engine) ListFeats(ctx context.Context, req ListFeatsRequest) (*ListFeatsResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var feats []*model.FeatDefinition

	allFeats := data.GlobalRegistry.ListFeats()
	for _, feat := range allFeats {
		if req.FilterType != nil && feat.Type != *req.FilterType {
			continue
		}
		feats = append(feats, feat)
	}

	result := make([]FeatInfo, len(feats))
	for i, feat := range feats {
		result[i] = featToInfo(feat)
	}

	return &ListFeatsResult{
		Feats: result,
	}, nil
}

// GetFeatDetails 获取专长详情
func (e *Engine) GetFeatDetails(ctx context.Context, req GetFeatDetailsRequest) (*GetFeatDetailsResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	feat, exists := data.GlobalRegistry.GetFeat(req.FeatID)
	if !exists {
		return nil, fmt.Errorf("feat not found: %s", req.FeatID)
	}

	return &GetFeatDetailsResult{
		Feat: &FeatInfo{
			ID:           feat.ID,
			Name:         feat.Name,
			Type:         string(feat.Type),
			Description:  feat.Description,
			Repeatable:   feat.Repeatable,
			Prerequisite: formatPrerequisite(feat.Prerequisite),
		},
	}, nil
}

// RemoveFeat 从角色移除专长（用于角色重建）
func (e *Engine) RemoveFeat(ctx context.Context, req RemoveFeatRequest) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return err
	}

	if err := e.checkPermission(game.Phase, OpRemoveFeat); err != nil {
		return err
	}

	pc, ok := game.PCs[req.PCID]
	if !ok {
		return ErrNotFound
	}

	// 从专长列表中移除
	newFeats := make([]model.FeatInstance, 0, len(pc.Feats))
	for _, featInstance := range pc.Feats {
		if featInstance.FeatID != req.FeatID {
			newFeats = append(newFeats, featInstance)
		}
	}

	if len(newFeats) == len(pc.Feats) {
		return fmt.Errorf("feat not found on character: %s", req.FeatID)
	}

	pc.Feats = newFeats

	// 注意：这里不撤销专长效果，因为效果可能已与其他系统交织
	// 完整的撤销需要更复杂的逻辑，通常在角色重建时使用

	if err := e.saveGame(ctx, game); err != nil {
		return err
	}

	return nil
}

// GetActorFeats 获取角色的专长列表
func (e *Engine) GetActorFeats(ctx context.Context, req GetActorRequest) (*ListFeatsResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	pc, ok := game.PCs[req.ActorID]
	if !ok {
		return nil, ErrNotFound
	}

	feats := make([]FeatInfo, len(pc.Feats))
	for i, instance := range pc.Feats {
		if featDef, exists := data.GlobalRegistry.GetFeat(instance.FeatID); exists {
			feats[i] = featToInfo(featDef)
		}
	}

	return &ListFeatsResult{
		Feats: feats,
	}, nil
}

// featToInfo 将 FeatDefinition 转换为 FeatInfo
func featToInfo(feat *model.FeatDefinition) FeatInfo {
	return FeatInfo{
		ID:           feat.ID,
		Name:         feat.Name,
		Type:         string(feat.Type),
		Description:  feat.Description,
		Repeatable:   feat.Repeatable,
		Prerequisite: formatPrerequisite(feat.Prerequisite),
	}
}

// formatPrerequisite 格式化先决条件为人类可读字符串
func formatPrerequisite(prereq *model.FeatPrerequisite) string {
	if prereq == nil {
		return ""
	}

	// 简化实现，返回属性要求
	if len(prereq.MinimumAbilityScores) > 0 {
		return "属性要求"
	}
	if prereq.MinimumLevel > 0 {
		return fmt.Sprintf("等级 %d", prereq.MinimumLevel)
	}
	if prereq.RequiredClass != "" {
		return fmt.Sprintf("职业：%s", prereq.RequiredClass)
	}
	if prereq.RequiredFeat != "" {
		return fmt.Sprintf("专长：%s", prereq.RequiredFeat)
	}
	if prereq.Description != "" {
		return prereq.Description
	}
	return ""
}
