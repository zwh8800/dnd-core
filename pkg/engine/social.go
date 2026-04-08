package engine

import (
	"context"
	"fmt"

	"github.com/zwh8800/dnd-core/pkg/model"
	"github.com/zwh8800/dnd-core/pkg/rules"
)

// InteractWithNPCRequest 与 NPC 互动请求
type InteractWithNPCRequest struct {
	GameID    model.ID              `json:"game_id"`
	NPCID     model.ID              `json:"npc_id"`
	CheckType model.SocialCheckType `json:"check_type"`
	Ability   int                   `json:"ability"`         // 相关属性值
	ProfBonus int                   `json:"prof_bonus"`      // 熟练加值
	HasProf   bool                  `json:"has_proficiency"` // 是否有熟练
}

// InteractWithNPCResult 与 NPC 互动结果
type InteractWithNPCResult struct {
	Result      *model.SocialInteractionResult `json:"result"`
	NewAttitude model.NPCAttitude              `json:"new_attitude"`
	Message     string                         `json:"message"`
}

// GetNPCAttitudeRequest 获取 NPC 态度请求
type GetNPCAttitudeRequest struct {
	GameID model.ID `json:"game_id"`
	NPCID  model.ID `json:"npc_id"`
}

// GetNPCAttitudeResult 获取 NPC 态度结果
type GetNPCAttitudeResult struct {
	Attitude    model.NPCAttitude             `json:"attitude"`
	Disposition model.NPCDisposition          `json:"disposition"`
	Interaction *model.SocialInteractionState `json:"interaction,omitempty"`
}

// InteractWithNPC 执行社交互动
func (e *Engine) InteractWithNPC(ctx context.Context, req InteractWithNPCRequest) (*InteractWithNPCResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	// 获取 NPC
	npc, ok := game.NPCs[req.NPCID]
	if !ok {
		return nil, fmt.Errorf("NPC %s not found", req.NPCID)
	}

	// 初始化社交状态
	if npc.SocialState == nil {
		npc.SocialState = &model.SocialInteractionState{
			CurrentAttitude: model.AttitudeIndifferent,
			Disposition:     model.DispositionIndifferent,
			Impressions:     make([]string, 0),
		}
	}

	// 执行社交检定
	socialResult, err := rules.PerformSocialCheck(
		req.Ability,
		req.ProfBonus,
		req.HasProf,
		npc.SocialState.Disposition,
		req.CheckType,
	)
	if err != nil {
		return nil, err
	}

	// 更新 NPC 态度
	oldAttitude := npc.SocialState.CurrentAttitude
	npc.SocialState.CurrentAttitude = socialResult.AttitudeChange
	npc.SocialState.InteractionCount++
	npc.SocialState.LastInteraction = string(req.CheckType)

	result := &InteractWithNPCResult{
		Result:      socialResult,
		NewAttitude: npc.SocialState.CurrentAttitude,
		Message: fmt.Sprintf("使用 %s 检定 NPC，态度从 %s 变为 %s",
			req.CheckType, oldAttitude, npc.SocialState.CurrentAttitude),
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return result, nil
}

// GetNPCAttitude 获取 NPC 当前态度
func (e *Engine) GetNPCAttitude(ctx context.Context, req GetNPCAttitudeRequest) (*GetNPCAttitudeResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	npc, ok := game.NPCs[req.NPCID]
	if !ok {
		return nil, fmt.Errorf("NPC %s not found", req.NPCID)
	}

	result := &GetNPCAttitudeResult{
		Attitude:    npc.SocialState.CurrentAttitude,
		Disposition: npc.SocialState.Disposition,
		Interaction: npc.SocialState,
	}

	return result, nil
}
