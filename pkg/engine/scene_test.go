package engine

import (
	"context"
	"testing"

	"github.com/zwh8800/dnd-core/internal/model"
)

// TestCreateScene 测试创建场景
func TestCreateScene(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	result, err := engine.CreateScene(ctx, gameID, "Tavern", "A cozy tavern", model.SceneTypeIndoor)
	if err != nil {
		t.Fatalf("Failed to create scene: %v", err)
	}

	if result.Scene == nil {
		t.Fatal("Result scene is nil")
	}
	if result.Scene.Name != "Tavern" {
		t.Errorf("Expected scene name 'Tavern', got %s", result.Scene.Name)
	}
	if result.Scene.Type != model.SceneTypeIndoor {
		t.Errorf("Expected scene type %s, got %s", model.SceneTypeIndoor, result.Scene.Type)
	}
}

// TestGetScene 测试获取场景
func TestGetScene(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	// 创建场景
	result, err := engine.CreateScene(ctx, gameID, "Forest", "A dense forest", model.SceneTypeOutdoor)
	if err != nil {
		t.Fatalf("Failed to create scene: %v", err)
	}

	// 获取场景
	scene, err := engine.GetScene(ctx, gameID, result.Scene.ID)
	if err != nil {
		t.Fatalf("Failed to get scene: %v", err)
	}

	if scene.ID != result.Scene.ID {
		t.Errorf("Expected scene ID %s, got %s", result.Scene.ID, scene.ID)
	}
	if scene.Name != "Forest" {
		t.Errorf("Expected scene name 'Forest', got %s", scene.Name)
	}
}

// TestUpdateScene 测试更新场景
func TestUpdateScene(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	// 创建场景
	result, err := engine.CreateScene(ctx, gameID, "Dungeon", "A dark dungeon", model.SceneTypeIndoor)
	if err != nil {
		t.Fatalf("Failed to create scene: %v", err)
	}

	// 更新场景
	updates := SceneUpdate{
		Description: "A very dark and dangerous dungeon",
		IsDark:      true,
		IsDarkSet:   true,
		LightLevel:  "dim",
	}

	err = engine.UpdateScene(ctx, gameID, result.Scene.ID, updates)
	if err != nil {
		t.Fatalf("Failed to update scene: %v", err)
	}

	// 验证更新
	scene, err := engine.GetScene(ctx, gameID, result.Scene.ID)
	if err != nil {
		t.Fatalf("Failed to get scene: %v", err)
	}

	if scene.Description != "A very dark and dangerous dungeon" {
		t.Errorf("Expected description 'A very dark and dangerous dungeon', got %s", scene.Description)
	}
	if !scene.IsDark {
		t.Error("Expected scene to be dark")
	}
	if scene.LightLevel != "dim" {
		t.Errorf("Expected light level 'dim', got %s", scene.LightLevel)
	}
}

// TestDeleteScene 测试删除场景
func TestDeleteScene(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	// 创建场景
	result, err := engine.CreateScene(ctx, gameID, "Temp Scene", "A temporary scene", model.SceneTypeIndoor)
	if err != nil {
		t.Fatalf("Failed to create scene: %v", err)
	}

	// 删除场景
	err = engine.DeleteScene(ctx, gameID, result.Scene.ID)
	if err != nil {
		t.Fatalf("Failed to delete scene: %v", err)
	}

	// 验证场景已删除
	_, err = engine.GetScene(ctx, gameID, result.Scene.ID)
	if err == nil {
		t.Fatal("Expected error when getting deleted scene, got nil")
	}
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

// TestListScenes 测试列出场景
func TestListScenes(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	// 创建多个场景
	for i := 0; i < 3; i++ {
		_, err := engine.CreateScene(ctx, gameID, "Scene", "A scene", model.SceneTypeIndoor)
		if err != nil {
			t.Fatalf("Failed to create scene %d: %v", i, err)
		}
	}

	// 列出场景
	scenes, err := engine.ListScenes(ctx, gameID)
	if err != nil {
		t.Fatalf("Failed to list scenes: %v", err)
	}

	if len(scenes) != 3 {
		t.Errorf("Expected 3 scenes, got %d", len(scenes))
	}
}

// TestSceneConnections 测试场景连接
func TestSceneConnections(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	// 创建两个场景
	scene1, err := engine.CreateScene(ctx, gameID, "Room 1", "First room", model.SceneTypeIndoor)
	if err != nil {
		t.Fatalf("Failed to create scene 1: %v", err)
	}

	scene2, err := engine.CreateScene(ctx, gameID, "Room 2", "Second room", model.SceneTypeIndoor)
	if err != nil {
		t.Fatalf("Failed to create scene 2: %v", err)
	}

	// 添加连接
	err = engine.AddSceneConnection(ctx, gameID, scene1.Scene.ID, scene2.Scene.ID,
		"A door connects the rooms", false, 0, false)
	if err != nil {
		t.Fatalf("Failed to add scene connection: %v", err)
	}

	// 验证连接
	scene, err := engine.GetScene(ctx, gameID, scene1.Scene.ID)
	if err != nil {
		t.Fatalf("Failed to get scene: %v", err)
	}

	if len(scene.Connections) != 1 {
		t.Errorf("Expected 1 connection, got %d", len(scene.Connections))
	}

	// 移除连接
	err = engine.RemoveSceneConnection(ctx, gameID, scene1.Scene.ID, scene2.Scene.ID)
	if err != nil {
		t.Fatalf("Failed to remove scene connection: %v", err)
	}

	// 验证连接已移除
	scene, err = engine.GetScene(ctx, gameID, scene1.Scene.ID)
	if err != nil {
		t.Fatalf("Failed to get scene: %v", err)
	}

	if len(scene.Connections) != 0 {
		t.Errorf("Expected 0 connections after removal, got %d", len(scene.Connections))
	}
}

// TestMoveActorToScene 测试移动角色到场景
func TestMoveActorToScene(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	// 创建场景
	scene1, err := engine.CreateScene(ctx, gameID, "Scene 1", "First scene", model.SceneTypeIndoor)
	if err != nil {
		t.Fatalf("Failed to create scene 1: %v", err)
	}

	scene2, err := engine.CreateScene(ctx, gameID, "Scene 2", "Second scene", model.SceneTypeIndoor)
	if err != nil {
		t.Fatalf("Failed to create scene 2: %v", err)
	}

	// 创建PC
	pc := &model.PlayerCharacter{
		Actor: model.Actor{
			Name:  "Hero",
			Size:  model.SizeMedium,
			Speed: 30,
			AbilityScores: model.AbilityScores{
				Strength: 16, Dexterity: 12, Constitution: 14,
				Intelligence: 10, Wisdom: 8, Charisma: 13,
			},
		},
		Race: model.RaceReference{Name: "Human"},
		Classes: []model.ClassLevel{
			{ClassName: "Fighter", Level: 1},
		},
		TotalLevel: 1,
	}
	pcResult, err := engine.CreatePC(ctx, gameID, pc)
	if err != nil {
		t.Fatalf("Failed to create PC: %v", err)
	}

	// 移动到场景1
	result, err := engine.MoveActorToScene(ctx, gameID, pcResult.ID, scene1.Scene.ID)
	if err != nil {
		t.Fatalf("Failed to move actor to scene 1: %v", err)
	}

	if !result.Success {
		t.Error("Expected move to be successful")
	}
	if result.ToScene != scene1.Scene.ID {
		t.Errorf("Expected to scene %s, got %s", scene1.Scene.ID, result.ToScene)
	}

	// 移动到场景2
	result, err = engine.MoveActorToScene(ctx, gameID, pcResult.ID, scene2.Scene.ID)
	if err != nil {
		t.Fatalf("Failed to move actor to scene 2: %v", err)
	}

	if result.FromScene != scene1.Scene.ID {
		t.Errorf("Expected from scene %s, got %s", scene1.Scene.ID, result.FromScene)
	}
	if result.ToScene != scene2.Scene.ID {
		t.Errorf("Expected to scene %s, got %s", scene2.Scene.ID, result.ToScene)
	}
}

// TestGetSceneActors 测试获取场景中的角色
func TestGetSceneActors(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	// 创建场景
	scene, err := engine.CreateScene(ctx, gameID, "Tavern", "A tavern", model.SceneTypeIndoor)
	if err != nil {
		t.Fatalf("Failed to create scene: %v", err)
	}

	// 创建PC
	pc := &model.PlayerCharacter{
		Actor: model.Actor{
			Name:  "Hero",
			Size:  model.SizeMedium,
			Speed: 30,
			AbilityScores: model.AbilityScores{
				Strength: 16, Dexterity: 12, Constitution: 14,
				Intelligence: 10, Wisdom: 8, Charisma: 13,
			},
		},
		Race: model.RaceReference{Name: "Human"},
		Classes: []model.ClassLevel{
			{ClassName: "Fighter", Level: 1},
		},
		TotalLevel: 1,
	}
	pcResult, err := engine.CreatePC(ctx, gameID, pc)
	if err != nil {
		t.Fatalf("Failed to create PC: %v", err)
	}

	// 移动PC到场景
	_, err = engine.MoveActorToScene(ctx, gameID, pcResult.ID, scene.Scene.ID)
	if err != nil {
		t.Fatalf("Failed to move actor to scene: %v", err)
	}

	// 获取场景中的角色
	actors, err := engine.GetSceneActors(ctx, gameID, scene.Scene.ID)
	if err != nil {
		t.Fatalf("Failed to get scene actors: %v", err)
	}

	if len(actors) != 1 {
		t.Errorf("Expected 1 actor in scene, got %d", len(actors))
	}
	if actors[0].ActorID != pcResult.ID {
		t.Errorf("Expected actor ID %s, got %s", pcResult.ID, actors[0].ActorID)
	}
}

// TestSceneItems 测试场景物品
func TestSceneItems(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	// 创建场景
	scene, err := engine.CreateScene(ctx, gameID, "Treasure Room", "A room with treasure", model.SceneTypeIndoor)
	if err != nil {
		t.Fatalf("Failed to create scene: %v", err)
	}

	// 创建物品
	game, err := engine.loadGame(ctx, gameID)
	if err != nil {
		t.Fatalf("Failed to load game: %v", err)
	}

	item := &model.Item{
		ID:          model.NewID(),
		Name:        "Gold Coin",
		Description: "A shiny gold coin",
		Type:        model.ItemTypeTreasure,
	}
	game.Items[item.ID] = item

	err = engine.saveGame(ctx, game)
	if err != nil {
		t.Fatalf("Failed to save game: %v", err)
	}

	// 添加物品到场景
	err = engine.AddItemToScene(ctx, gameID, scene.Scene.ID, item.ID)
	if err != nil {
		t.Fatalf("Failed to add item to scene: %v", err)
	}

	// 获取场景物品
	items, err := engine.GetSceneItems(ctx, gameID, scene.Scene.ID)
	if err != nil {
		t.Fatalf("Failed to get scene items: %v", err)
	}

	if len(items) != 1 {
		t.Errorf("Expected 1 item in scene, got %d", len(items))
	}
	if items[0] != item.ID {
		t.Errorf("Expected item ID %s, got %s", item.ID, items[0])
	}

	// 移除物品
	err = engine.RemoveItemFromScene(ctx, gameID, scene.Scene.ID, item.ID)
	if err != nil {
		t.Fatalf("Failed to remove item from scene: %v", err)
	}

	// 验证物品已移除
	items, err = engine.GetSceneItems(ctx, gameID, scene.Scene.ID)
	if err != nil {
		t.Fatalf("Failed to get scene items: %v", err)
	}

	if len(items) != 0 {
		t.Errorf("Expected 0 items after removal, got %d", len(items))
	}
}

// TestSceneRegions 测试场景区域
func TestSceneRegions(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	// 创建场景
	scene, err := engine.CreateScene(ctx, gameID, "Battlefield", "A battlefield", model.SceneTypeOutdoor)
	if err != nil {
		t.Fatalf("Failed to create scene: %v", err)
	}

	// 添加区域
	bounds := &model.Rect{
		X: 0, Y: 0, Width: 100, Height: 100,
	}
	region, err := engine.AddRegion(ctx, gameID, scene.Scene.ID,
		"Combat Zone", "The main combat area", bounds)
	if err != nil {
		t.Fatalf("Failed to add region: %v", err)
	}

	if region.Name != "Combat Zone" {
		t.Errorf("Expected region name 'Combat Zone', got %s", region.Name)
	}

	// 移除区域
	err = engine.RemoveRegion(ctx, gameID, scene.Scene.ID, region.ID)
	if err != nil {
		t.Fatalf("Failed to remove region: %v", err)
	}

	// 验证区域已移除
	game, err := engine.loadGame(ctx, gameID)
	if err != nil {
		t.Fatalf("Failed to load game: %v", err)
	}

	if _, ok := game.Scenes[scene.Scene.ID].Regions[region.ID]; ok {
		t.Error("Expected region to be removed")
	}
}
