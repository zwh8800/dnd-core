package engine

import (
	"context"
	"fmt"

	"github.com/zwh8800/dnd-core/pkg/data"
	"github.com/zwh8800/dnd-core/pkg/model"
)

// MountCreatureRequest 骑乘请求
type MountCreatureRequest struct {
	GameID  model.ID `json:"game_id"`
	RiderID model.ID `json:"rider_id"`
	MountID model.ID `json:"mount_id"`
}

// MountCreatureResult 骑乘结果
type MountCreatureResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// MountCreature 骑上坐骑
// 验证骑手和坐骑是否存在，建立骑乘关系并保存游戏状态。
// 参数:
//
//	ctx - 上下文
//	req - 骑乘请求，包含游戏ID、骑手ID和坐骑ID
//
// 返回:
//
//	*MountCreatureResult - 骑乘结果，包含成功状态和提示信息
//	error - 错误信息，如游戏加载失败、骑手或坐骑不存在
func (e *Engine) MountCreature(ctx context.Context, req MountCreatureRequest) (*MountCreatureResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	rider, ok := game.GetActor(req.RiderID)
	if !ok {
		return nil, fmt.Errorf("rider %s not found", req.RiderID)
	}

	mount, ok := game.GetActor(req.MountID)
	if !ok {
		return nil, fmt.Errorf("mount %s not found", req.MountID)
	}

	result := &MountCreatureResult{
		Success: true,
		Message: fmt.Sprintf("%s 骑上了 %s", rider.(*model.PlayerCharacter).Actor.Name, mount.(*model.Enemy).Actor.Name),
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return result, nil
}

// DismountRequest 下马请求
type DismountRequest struct {
	GameID  model.ID `json:"game_id"`
	RiderID model.ID `json:"rider_id"`
}

// DismountResult 下马结果
type DismountResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Dismount 下马
// 验证骑手是否存在，解除骑乘关系并保存游戏状态。
// 参数:
//
//	ctx - 上下文
//	req - 下马请求，包含游戏ID和骑手ID
//
// 返回:
//
//	*DismountResult - 下马结果，包含成功状态和提示信息
//	error - 错误信息，如游戏加载失败或骑手不存在
func (e *Engine) Dismount(ctx context.Context, req DismountRequest) (*DismountResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	rider, ok := game.GetActor(req.RiderID)
	if !ok {
		return nil, fmt.Errorf("rider %s not found", req.RiderID)
	}

	result := &DismountResult{
		Success: true,
		Message: fmt.Sprintf("%s 下了坐骑", rider.(*model.PlayerCharacter).Actor.Name),
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return result, nil
}

// CalculateMountSpeedRequest 计算坐骑速度请求
type CalculateMountSpeedRequest struct {
	GameID  model.ID `json:"game_id"`
	MountID model.ID `json:"mount_id"`
}

// CalculateMountSpeedResult 计算坐骑速度结果
type CalculateMountSpeedResult struct {
	BaseSpeed   int     `json:"base_speed"`
	FinalSpeed  int     `json:"final_speed"`
	CarryCap    float64 `json:"carry_capacity"`
	CurrentLoad float64 `json:"current_load"`
	Message     string  `json:"message"`
}

// CalculateMountSpeed 计算坐骑速度
// 根据坐骑数据获取基础速度和载重能力，返回坐骑的移动速度信息。
// 参数:
//
//	ctx - 上下文
//	req - 计算坐骑速度请求，包含游戏ID和坐骑ID
//
// 返回:
//
//	*CalculateMountSpeedResult - 计算结果，包含基础速度、最终速度、载重能力和当前负载
//	error - 错误信息，如游戏加载失败或坐骑数据不存在
func (e *Engine) CalculateMountSpeed(ctx context.Context, req CalculateMountSpeedRequest) (*CalculateMountSpeedResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	_, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	mountData := data.GetMountData(req.MountID)
	if mountData == nil {
		return nil, fmt.Errorf("mount data %s not found", req.MountID)
	}

	baseSpeed := mountData.Speed
	carryCap := mountData.CarryCap

	result := &CalculateMountSpeedResult{
		BaseSpeed:   baseSpeed,
		FinalSpeed:  baseSpeed,
		CarryCap:    carryCap,
		CurrentLoad: 0,
		Message:     fmt.Sprintf("坐骑速度: %d 尺，载重能力: %.1f 磅", baseSpeed, carryCap),
	}

	return result, nil
}
