package engine

import (
	"context"
	"fmt"

	"github.com/zwh8800/dnd-core/internal/model"
)

// CreateSceneRequest 创建场景请求
type CreateSceneRequest struct {
	GameID      model.ID        `json:"game_id"`     // 游戏会话ID
	Name        string          `json:"name"`        // 场景名称
	Description string          `json:"description"` // 场景描述
	SceneType   model.SceneType `json:"scene_type"`  // 场景类型
}

// GetSceneRequest 获取场景请求
type GetSceneRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID
	SceneID model.ID `json:"scene_id"` // 场景ID
}

// UpdateSceneRequest 更新场景请求
type UpdateSceneRequest struct {
	GameID  model.ID    `json:"game_id"`  // 游戏会话ID
	SceneID model.ID    `json:"scene_id"` // 场景ID
	Updates SceneUpdate `json:"updates"`  // 更新内容
}

// DeleteSceneRequest 删除场景请求
type DeleteSceneRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID
	SceneID model.ID `json:"scene_id"` // 场景ID
}

// ListScenesRequest 列出场景请求
type ListScenesRequest struct {
	GameID model.ID `json:"game_id"` // 游戏会话ID
}

// ListScenesResult 列出场景结果
type ListScenesResult struct {
	Scenes []SceneInfo `json:"scenes"` // 场景列表
}

// SetCurrentSceneRequest 设置当前场景请求
type SetCurrentSceneRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID
	SceneID model.ID `json:"scene_id"` // 场景ID
}

// GetCurrentSceneRequest 获取当前场景请求
type GetCurrentSceneRequest struct {
	GameID model.ID `json:"game_id"` // 游戏会话ID
}

// AddSceneConnectionRequest 添加场景连接请求
type AddSceneConnectionRequest struct {
	GameID        model.ID `json:"game_id"`         // 游戏会话ID
	SceneID       model.ID `json:"scene_id"`        // 源场景ID
	TargetSceneID model.ID `json:"target_scene_id"` // 目标场景ID
	Description   string   `json:"description"`     // 连接描述
	Locked        bool     `json:"locked"`          // 是否锁定
	DC            int      `json:"dc"`              // 解锁难度等级
	Hidden        bool     `json:"hidden"`          // 是否隐藏
}

// RemoveSceneConnectionRequest 移除场景连接请求
type RemoveSceneConnectionRequest struct {
	GameID        model.ID `json:"game_id"`         // 游戏会话ID
	SceneID       model.ID `json:"scene_id"`        // 源场景ID
	TargetSceneID model.ID `json:"target_scene_id"` // 目标场景ID
}

// MoveActorToSceneRequest 移动角色到场景请求
type MoveActorToSceneRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID
	ActorID model.ID `json:"actor_id"` // 角色ID
	SceneID model.ID `json:"scene_id"` // 目标场景ID
}

// GetSceneActorsRequest 获取场景角色请求
type GetSceneActorsRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID
	SceneID model.ID `json:"scene_id"` // 场景ID
}

// GetSceneActorsResult 获取场景角色结果
type GetSceneActorsResult struct {
	Actors []SceneActorInfo `json:"actors"` // 场景中的角色列表
}

// AddItemToSceneRequest 添加物品到场景请求
type AddItemToSceneRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID
	SceneID model.ID `json:"scene_id"` // 场景ID
	ItemID  model.ID `json:"item_id"`  // 物品ID
}

// RemoveItemFromSceneRequest 从场景移除物品请求
type RemoveItemFromSceneRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID
	SceneID model.ID `json:"scene_id"` // 场景ID
	ItemID  model.ID `json:"item_id"`  // 物品ID
}

// GetSceneItemsRequest 获取场景物品请求
type GetSceneItemsRequest struct {
	GameID  model.ID `json:"game_id"`  // 游戏会话ID
	SceneID model.ID `json:"scene_id"` // 场景ID
}

// GetSceneItemsResult 获取场景物品结果
type GetSceneItemsResult struct {
	Items []model.ID `json:"items"` // 物品ID列表
}

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

// SceneMoveResult 场景移动结果
type SceneMoveResult struct {
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

// SceneUpdate 场景更新
type SceneUpdate struct {
	Name        string `json:"name,omitempty"`        // 场景名称
	Description string `json:"description,omitempty"` // 场景描述
	Details     string `json:"details,omitempty"`     // 场景细节
	IsDark      bool   `json:"is_dark"`               // 是否黑暗
	IsDarkSet   bool   `json:"is_dark_set"`           // 是否设置了黑暗标志
	LightLevel  string `json:"light_level,omitempty"` // 光照等级
	Weather     string `json:"weather,omitempty"`     // 天气
	Terrain     string `json:"terrain,omitempty"`     // 地形
}

// CreateScene 创建新场景
func (e *Engine) CreateScene(ctx context.Context, req CreateSceneRequest) (*CreateSceneResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	scene := model.NewScene(req.Name, req.Description, req.SceneType)
	game.Scenes[scene.ID] = scene

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &CreateSceneResult{
		Scene:   scene,
		Message: fmt.Sprintf("创建了场景: %s", req.Name),
	}, nil
}

// GetScene 获取场景信息
func (e *Engine) GetScene(ctx context.Context, req GetSceneRequest) (*SceneInfo, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	scene, ok := game.Scenes[req.SceneID]
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
func (e *Engine) UpdateScene(ctx context.Context, req UpdateSceneRequest) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return err
	}

	scene, ok := game.Scenes[req.SceneID]
	if !ok {
		return ErrNotFound
	}

	if req.Updates.Name != "" {
		scene.Name = req.Updates.Name
	}
	if req.Updates.Description != "" {
		scene.Description = req.Updates.Description
	}
	if req.Updates.Details != "" {
		scene.Details = req.Updates.Details
	}
	if req.Updates.IsDarkSet {
		scene.IsDark = req.Updates.IsDark
	}
	if req.Updates.LightLevel != "" {
		scene.LightLevel = req.Updates.LightLevel
	}
	if req.Updates.Weather != "" {
		scene.Weather = req.Updates.Weather
	}
	if req.Updates.Terrain != "" {
		scene.Terrain = req.Updates.Terrain
	}

	if err := e.saveGame(ctx, game); err != nil {
		return err
	}

	return nil
}

// DeleteScene 删除场景
func (e *Engine) DeleteScene(ctx context.Context, req DeleteSceneRequest) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return err
	}

	if _, ok := game.Scenes[req.SceneID]; !ok {
		return ErrNotFound
	}

	// 检查是否有角色在该场景
	for _, pc := range game.PCs {
		if pc.SceneID == req.SceneID {
			return fmt.Errorf("scene %s still has actors in it", req.SceneID)
		}
	}
	for _, npc := range game.NPCs {
		if npc.SceneID == req.SceneID {
			return fmt.Errorf("scene %s still has actors in it", req.SceneID)
		}
	}
	for _, enemy := range game.Enemies {
		if enemy.SceneID == req.SceneID {
			return fmt.Errorf("scene %s still has actors in it", req.SceneID)
		}
	}
	for _, companion := range game.Companions {
		if companion.SceneID == req.SceneID {
			return fmt.Errorf("scene %s still has actors in it", req.SceneID)
		}
	}

	// 删除场景
	delete(game.Scenes, req.SceneID)

	if err := e.saveGame(ctx, game); err != nil {
		return err
	}

	return nil
}

// ListScenes 列出所有场景
func (e *Engine) ListScenes(ctx context.Context, req ListScenesRequest) (*ListScenesResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, req.GameID)
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

	return &ListScenesResult{Scenes: result}, nil
}

// SetCurrentScene 设置当前场景
func (e *Engine) SetCurrentScene(ctx context.Context, req SetCurrentSceneRequest) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return err
	}

	if _, ok := game.Scenes[req.SceneID]; !ok {
		return ErrNotFound
	}

	game.CurrentScene = &req.SceneID

	if err := e.saveGame(ctx, game); err != nil {
		return err
	}

	return nil
}

// GetCurrentScene 获取当前场景
func (e *Engine) GetCurrentScene(ctx context.Context, req GetCurrentSceneRequest) (*SceneInfo, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, req.GameID)
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
func (e *Engine) AddSceneConnection(ctx context.Context, req AddSceneConnectionRequest) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return err
	}

	scene, ok := game.Scenes[req.SceneID]
	if !ok {
		return ErrNotFound
	}

	if _, ok := game.Scenes[req.TargetSceneID]; !ok {
		return ErrNotFound
	}

	scene.Connections[req.TargetSceneID] = &model.SceneConnection{
		TargetSceneID: req.TargetSceneID,
		Description:   req.Description,
		Locked:        req.Locked,
		DC:            req.DC,
		Hidden:        req.Hidden,
	}

	if err := e.saveGame(ctx, game); err != nil {
		return err
	}

	return nil
}

// RemoveSceneConnection 移除场景连接
func (e *Engine) RemoveSceneConnection(ctx context.Context, req RemoveSceneConnectionRequest) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return err
	}

	scene, ok := game.Scenes[req.SceneID]
	if !ok {
		return ErrNotFound
	}

	delete(scene.Connections, req.TargetSceneID)

	if err := e.saveGame(ctx, game); err != nil {
		return err
	}

	return nil
}

// MoveActorToScene 移动角色到另一个场景
func (e *Engine) MoveActorToScene(ctx context.Context, req MoveActorToSceneRequest) (*MoveActorResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	// 验证目标场景存在
	if _, ok := game.Scenes[req.SceneID]; !ok {
		return nil, ErrNotFound
	}

	// 获取角色
	actor, ok := game.GetActor(req.ActorID)
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
	baseActor.SceneID = req.SceneID

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &MoveActorResult{
		SceneMoveResult: &SceneMoveResult{
			Success:   true,
			ActorID:   req.ActorID,
			FromScene: fromScene,
			ToScene:   req.SceneID,
			Message:   fmt.Sprintf("将 %s 移动到场景 %s", baseActor.Name, req.SceneID),
		},
	}, nil
}

// GetSceneActors 获取场景中的所有角色
func (e *Engine) GetSceneActors(ctx context.Context, req GetSceneActorsRequest) (*GetSceneActorsResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	result := make([]SceneActorInfo, 0)

	for _, pc := range game.PCs {
		if pc.SceneID == req.SceneID {
			result = append(result, SceneActorInfo{
				ActorID:   pc.ID,
				ActorName: pc.Name,
				ActorType: model.ActorTypePC,
				Position:  pc.Position,
			})
		}
	}
	for _, npc := range game.NPCs {
		if npc.SceneID == req.SceneID {
			result = append(result, SceneActorInfo{
				ActorID:   npc.ID,
				ActorName: npc.Name,
				ActorType: model.ActorTypeNPC,
				Position:  npc.Position,
			})
		}
	}
	for _, enemy := range game.Enemies {
		if enemy.SceneID == req.SceneID {
			result = append(result, SceneActorInfo{
				ActorID:   enemy.ID,
				ActorName: enemy.Name,
				ActorType: model.ActorTypeEnemy,
				Position:  enemy.Position,
			})
		}
	}
	for _, companion := range game.Companions {
		if companion.SceneID == req.SceneID {
			result = append(result, SceneActorInfo{
				ActorID:   companion.ID,
				ActorName: companion.Name,
				ActorType: model.ActorTypeCompanion,
				Position:  companion.Position,
			})
		}
	}

	return &GetSceneActorsResult{Actors: result}, nil
}

// AddItemToScene 添加物品到场景
func (e *Engine) AddItemToScene(ctx context.Context, req AddItemToSceneRequest) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return err
	}

	scene, ok := game.Scenes[req.SceneID]
	if !ok {
		return ErrNotFound
	}

	scene.Items = append(scene.Items, req.ItemID)

	if err := e.saveGame(ctx, game); err != nil {
		return err
	}

	return nil
}

// RemoveItemFromScene 从场景移除物品
func (e *Engine) RemoveItemFromScene(ctx context.Context, req RemoveItemFromSceneRequest) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return err
	}

	scene, ok := game.Scenes[req.SceneID]
	if !ok {
		return ErrNotFound
	}

	for i, id := range scene.Items {
		if id == req.ItemID {
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
func (e *Engine) GetSceneItems(ctx context.Context, req GetSceneItemsRequest) (*GetSceneItemsResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	scene, ok := game.Scenes[req.SceneID]
	if !ok {
		return nil, ErrNotFound
	}

	return &GetSceneItemsResult{Items: scene.Items}, nil
}
