package engine

import (
	"context"
	"fmt"

	"github.com/zwh8800/dnd-core/internal/model"
)

// CreateSceneResult 创建场景结果
type CreateSceneResult struct {
	Scene   *model.Scene `json:"scene"`
	Message string       `json:"message"`
}

// SceneInfo 场景信息
type SceneInfo struct {
	ID          model.ID         `json:"id"`
	Name        string           `json:"name"`
	Type        model.SceneType  `json:"type"`
	Description string           `json:"description"`
	Connections []ConnectionInfo `json:"connections"`
	RegionCount int              `json:"region_count"`
	ItemCount   int              `json:"item_count"`
	IsDark      bool             `json:"is_dark"`
	LightLevel  string           `json:"light_level"`
	CustomData  map[string]any   `json:"custom_data,omitempty"`
}

// ConnectionInfo 连接信息
type ConnectionInfo struct {
	TargetSceneID model.ID `json:"target_scene_id"`
	Description   string   `json:"description"`
	Locked        bool     `json:"locked"`
	DC            int      `json:"dc,omitempty"`
	Hidden        bool     `json:"hidden"`
}

// MoveActorResult 移动角色结果
type MoveActorResult struct {
	Success   bool     `json:"success"`
	ActorID   model.ID `json:"actor_id"`
	FromScene model.ID `json:"from_scene"`
	ToScene   model.ID `json:"to_scene"`
	Message   string   `json:"message"`
}

// SceneActorInfo 场景中的角色信息
type SceneActorInfo struct {
	ActorID   model.ID        `json:"actor_id"`
	ActorName string          `json:"actor_name"`
	ActorType model.ActorType `json:"actor_type"`
	Position  *model.Point    `json:"position,omitempty"`
}

// CreateScene 创建新场景
func (e *Engine) CreateScene(ctx context.Context, gameID model.ID, name, description string, sceneType model.SceneType) (*CreateSceneResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	scene := model.NewScene(name, description, sceneType)
	game.Scenes[scene.ID] = scene

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &CreateSceneResult{
		Scene:   scene,
		Message: fmt.Sprintf("创建了场景: %s", name),
	}, nil
}

// GetScene 获取场景信息
func (e *Engine) GetScene(ctx context.Context, gameID model.ID, sceneID model.ID) (*SceneInfo, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	scene, ok := game.Scenes[sceneID]
	if !ok {
		return nil, ErrNotFound
	}

	connections := make([]ConnectionInfo, 0, len(scene.Connections))
	for _, conn := range scene.Connections {
		connections = append(connections, ConnectionInfo{
			TargetSceneID: conn.TargetSceneID,
			Description:   conn.Description,
			Locked:        conn.Locked,
			DC:            conn.DC,
			Hidden:        conn.Hidden,
		})
	}

	return &SceneInfo{
		ID:          scene.ID,
		Name:        scene.Name,
		Type:        scene.Type,
		Description: scene.Description,
		Connections: connections,
		RegionCount: len(scene.Regions),
		ItemCount:   len(scene.Items),
		IsDark:      scene.IsDark,
		LightLevel:  scene.LightLevel,
		CustomData:  scene.CustomData,
	}, nil
}

// UpdateScene 更新场景信息
func (e *Engine) UpdateScene(ctx context.Context, gameID model.ID, sceneID model.ID, updates SceneUpdate) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return err
	}

	scene, ok := game.Scenes[sceneID]
	if !ok {
		return ErrNotFound
	}

	if updates.Name != "" {
		scene.Name = updates.Name
	}
	if updates.Description != "" {
		scene.Description = updates.Description
	}
	if updates.Details != "" {
		scene.Details = updates.Details
	}
	if updates.IsDarkSet {
		scene.IsDark = updates.IsDark
	}
	if updates.LightLevel != "" {
		scene.LightLevel = updates.LightLevel
	}
	if updates.Weather != "" {
		scene.Weather = updates.Weather
	}
	if updates.Terrain != "" {
		scene.Terrain = updates.Terrain
	}

	if err := e.saveGame(ctx, game); err != nil {
		return err
	}

	return nil
}

// SceneUpdate 场景更新
type SceneUpdate struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Details     string `json:"details,omitempty"`
	IsDark      bool   `json:"is_dark"`
	IsDarkSet   bool   `json:"is_dark_set"`
	LightLevel  string `json:"light_level,omitempty"`
	Weather     string `json:"weather,omitempty"`
	Terrain     string `json:"terrain,omitempty"`
}

// DeleteScene 删除场景
func (e *Engine) DeleteScene(ctx context.Context, gameID model.ID, sceneID model.ID) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return err
	}

	if _, ok := game.Scenes[sceneID]; !ok {
		return ErrNotFound
	}

	// 检查是否有角色在该场景
	for _, pc := range game.PCs {
		if pc.SceneID == sceneID {
			return fmt.Errorf("scene %s still has actors in it", sceneID)
		}
	}
	for _, npc := range game.NPCs {
		if npc.SceneID == sceneID {
			return fmt.Errorf("scene %s still has actors in it", sceneID)
		}
	}
	for _, enemy := range game.Enemies {
		if enemy.SceneID == sceneID {
			return fmt.Errorf("scene %s still has actors in it", sceneID)
		}
	}
	for _, companion := range game.Companions {
		if companion.SceneID == sceneID {
			return fmt.Errorf("scene %s still has actors in it", sceneID)
		}
	}

	// 删除场景
	delete(game.Scenes, sceneID)

	if err := e.saveGame(ctx, game); err != nil {
		return err
	}

	return nil
}

// ListScenes 列出所有场景
func (e *Engine) ListScenes(ctx context.Context, gameID model.ID) ([]SceneInfo, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	result := make([]SceneInfo, 0, len(game.Scenes))
	for _, scene := range game.Scenes {
		connections := make([]ConnectionInfo, 0, len(scene.Connections))
		for _, conn := range scene.Connections {
			connections = append(connections, ConnectionInfo{
				TargetSceneID: conn.TargetSceneID,
				Description:   conn.Description,
				Locked:        conn.Locked,
				DC:            conn.DC,
				Hidden:        conn.Hidden,
			})
		}

		result = append(result, SceneInfo{
			ID:          scene.ID,
			Name:        scene.Name,
			Type:        scene.Type,
			Description: scene.Description,
			Connections: connections,
			RegionCount: len(scene.Regions),
			ItemCount:   len(scene.Items),
			IsDark:      scene.IsDark,
			LightLevel:  scene.LightLevel,
		})
	}

	return result, nil
}

// SetCurrentScene 设置当前场景
func (e *Engine) SetCurrentScene(ctx context.Context, gameID model.ID, sceneID model.ID) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return err
	}

	if _, ok := game.Scenes[sceneID]; !ok {
		return ErrNotFound
	}

	game.CurrentScene = &sceneID

	if err := e.saveGame(ctx, game); err != nil {
		return err
	}

	return nil
}

// GetCurrentScene 获取当前场景
func (e *Engine) GetCurrentScene(ctx context.Context, gameID model.ID) (*SceneInfo, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	if game.CurrentScene == nil {
		return nil, fmt.Errorf("no current scene set")
	}

	scene, ok := game.Scenes[*game.CurrentScene]
	if !ok {
		return nil, ErrNotFound
	}

	connections := make([]ConnectionInfo, 0, len(scene.Connections))
	for _, conn := range scene.Connections {
		connections = append(connections, ConnectionInfo{
			TargetSceneID: conn.TargetSceneID,
			Description:   conn.Description,
			Locked:        conn.Locked,
			DC:            conn.DC,
			Hidden:        conn.Hidden,
		})
	}

	return &SceneInfo{
		ID:          scene.ID,
		Name:        scene.Name,
		Type:        scene.Type,
		Description: scene.Description,
		Connections: connections,
		RegionCount: len(scene.Regions),
		ItemCount:   len(scene.Items),
		IsDark:      scene.IsDark,
		LightLevel:  scene.LightLevel,
	}, nil
}

// AddSceneConnection 添加场景连接
func (e *Engine) AddSceneConnection(ctx context.Context, gameID model.ID, sceneID model.ID, targetSceneID model.ID, description string, locked bool, dc int, hidden bool) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return err
	}

	scene, ok := game.Scenes[sceneID]
	if !ok {
		return ErrNotFound
	}

	if _, ok := game.Scenes[targetSceneID]; !ok {
		return ErrNotFound
	}

	scene.Connections[targetSceneID] = &model.SceneConnection{
		TargetSceneID: targetSceneID,
		Description:   description,
		Locked:        locked,
		DC:            dc,
		Hidden:        hidden,
	}

	if err := e.saveGame(ctx, game); err != nil {
		return err
	}

	return nil
}

// RemoveSceneConnection 移除场景连接
func (e *Engine) RemoveSceneConnection(ctx context.Context, gameID model.ID, sceneID model.ID, targetSceneID model.ID) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return err
	}

	scene, ok := game.Scenes[sceneID]
	if !ok {
		return ErrNotFound
	}

	delete(scene.Connections, targetSceneID)

	if err := e.saveGame(ctx, game); err != nil {
		return err
	}

	return nil
}

// MoveActorToScene 移动角色到另一个场景
func (e *Engine) MoveActorToScene(ctx context.Context, gameID model.ID, actorID model.ID, sceneID model.ID) (*MoveActorResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	// 验证目标场景存在
	if _, ok := game.Scenes[sceneID]; !ok {
		return nil, ErrNotFound
	}

	// 获取角色
	actor, ok := game.GetActor(actorID)
	if !ok {
		return nil, ErrNotFound
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
	}

	fromScene := baseActor.SceneID
	baseActor.SceneID = sceneID

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &MoveActorResult{
		Success:   true,
		ActorID:   actorID,
		FromScene: fromScene,
		ToScene:   sceneID,
		Message:   fmt.Sprintf("将 %s 移动到场景 %s", baseActor.Name, sceneID),
	}, nil
}

// GetSceneActors 获取场景中的所有角色
func (e *Engine) GetSceneActors(ctx context.Context, gameID model.ID, sceneID model.ID) ([]SceneActorInfo, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	result := make([]SceneActorInfo, 0)

	for _, pc := range game.PCs {
		if pc.SceneID == sceneID {
			result = append(result, SceneActorInfo{
				ActorID:   pc.ID,
				ActorName: pc.Name,
				ActorType: model.ActorTypePC,
				Position:  pc.Position,
			})
		}
	}
	for _, npc := range game.NPCs {
		if npc.SceneID == sceneID {
			result = append(result, SceneActorInfo{
				ActorID:   npc.ID,
				ActorName: npc.Name,
				ActorType: model.ActorTypeNPC,
				Position:  npc.Position,
			})
		}
	}
	for _, enemy := range game.Enemies {
		if enemy.SceneID == sceneID {
			result = append(result, SceneActorInfo{
				ActorID:   enemy.ID,
				ActorName: enemy.Name,
				ActorType: model.ActorTypeEnemy,
				Position:  enemy.Position,
			})
		}
	}
	for _, companion := range game.Companions {
		if companion.SceneID == sceneID {
			result = append(result, SceneActorInfo{
				ActorID:   companion.ID,
				ActorName: companion.Name,
				ActorType: model.ActorTypeCompanion,
				Position:  companion.Position,
			})
		}
	}

	return result, nil
}

// AddItemToScene 添加物品到场景
func (e *Engine) AddItemToScene(ctx context.Context, gameID model.ID, sceneID model.ID, itemID model.ID) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return err
	}

	scene, ok := game.Scenes[sceneID]
	if !ok {
		return ErrNotFound
	}

	scene.Items = append(scene.Items, itemID)

	if err := e.saveGame(ctx, game); err != nil {
		return err
	}

	return nil
}

// RemoveItemFromScene 从场景移除物品
func (e *Engine) RemoveItemFromScene(ctx context.Context, gameID model.ID, sceneID model.ID, itemID model.ID) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return err
	}

	scene, ok := game.Scenes[sceneID]
	if !ok {
		return ErrNotFound
	}

	for i, id := range scene.Items {
		if id == itemID {
			scene.Items = append(scene.Items[:i], scene.Items[i+1:]...)
			break
		}
	}

	if err := e.saveGame(ctx, game); err != nil {
		return err
	}

	return nil
}

// GetSceneItems 获取场景中的物品
func (e *Engine) GetSceneItems(ctx context.Context, gameID model.ID, sceneID model.ID) ([]model.ID, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	scene, ok := game.Scenes[sceneID]
	if !ok {
		return nil, ErrNotFound
	}

	return scene.Items, nil
}

// AddRegion 添加场景区域
func (e *Engine) AddRegion(ctx context.Context, gameID model.ID, sceneID model.ID, name, description string, bounds *model.Rect) (*model.SceneRegion, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	scene, ok := game.Scenes[sceneID]
	if !ok {
		return nil, ErrNotFound
	}

	region := &model.SceneRegion{
		ID:          model.NewID(),
		Name:        name,
		Description: description,
		Bounds:      bounds,
	}

	scene.Regions[region.ID] = region

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return region, nil
}

// RemoveRegion 移除场景区域
func (e *Engine) RemoveRegion(ctx context.Context, gameID model.ID, sceneID model.ID, regionID model.ID) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return err
	}

	scene, ok := game.Scenes[sceneID]
	if !ok {
		return ErrNotFound
	}

	delete(scene.Regions, regionID)

	if err := e.saveGame(ctx, game); err != nil {
		return err
	}

	return nil
}
