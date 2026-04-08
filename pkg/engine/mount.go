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
