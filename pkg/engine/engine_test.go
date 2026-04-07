package engine

import (
	"context"
	"testing"

	"github.com/zwh8800/dnd-core/internal/model"
	"github.com/zwh8800/dnd-core/internal/storage"
)

// TestNewEngine 测试引擎创建
func TestNewEngine(t *testing.T) {
	cfg := DefaultConfig()
	engine, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	if engine == nil {
		t.Fatal("Engine is nil")
	}
	if engine.store == nil {
		t.Fatal("Engine store is nil")
	}
	if engine.roller == nil {
		t.Fatal("Engine roller is nil")
	}

	// 测试关闭
	err = engine.Close()
	if err != nil {
		t.Fatalf("Failed to close engine: %v", err)
	}
}

// TestNewEngineWithCustomStorage 测试自定义存储
func TestNewEngineWithCustomStorage(t *testing.T) {
	cfg := Config{
		Storage:  storage.NewMemoryStore(),
		DiceSeed: 12345,
	}
	engine, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create engine with custom storage: %v", err)
	}
	defer engine.Close()

	// 验证存储已初始化
	if engine.store == nil {
		t.Fatal("Engine store is nil")
	}
}

// TestEngineClose 测试引擎关闭
func TestEngineClose(t *testing.T) {
	engine, err := New(DefaultConfig())
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	// 第一次关闭应该成功
	err = engine.Close()
	if err != nil {
		t.Fatalf("First close failed: %v", err)
	}

	// 第二次关闭也应该成功（幂等）
	err = engine.Close()
	if err != nil {
		t.Fatalf("Second close failed: %v", err)
	}
}

// TestEngineLoadSaveGame 测试游戏加载和保存
func TestEngineLoadSaveGame(t *testing.T) {
	engine, err := New(DefaultConfig())
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer engine.Close()

	ctx := context.Background()

	// 创建游戏
	game := model.NewGameState("Test Game", "A test game")

	// 保存游戏
	err = engine.saveGame(ctx, game)
	if err != nil {
		t.Fatalf("Failed to save game: %v", err)
	}

	// 加载游戏
	loadedGame, err := engine.loadGame(ctx, game.ID)
	if err != nil {
		t.Fatalf("Failed to load game: %v", err)
	}

	if loadedGame.ID != game.ID {
		t.Errorf("Expected game ID %s, got %s", game.ID, loadedGame.ID)
	}
	if loadedGame.Name != "Test Game" {
		t.Errorf("Expected game name 'Test Game', got %s", loadedGame.Name)
	}
}

// TestEngineLoadNotFound 测试加载不存在的游戏
func TestEngineLoadNotFound(t *testing.T) {
	engine, err := New(DefaultConfig())
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}
	defer engine.Close()

	ctx := context.Background()
	fakeID := model.ID("00000000000000000000000000")

	_, err = engine.loadGame(ctx, fakeID)
	if err == nil {
		t.Fatal("Expected error when loading non-existent game, got nil")
	}
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

// TestEngineError 测试引擎错误
func TestEngineError(t *testing.T) {
	err := &EngineError{
		Op:    "testOperation",
		Err:   ErrNotFound,
		Phase: model.PhaseExploration,
		Details: map[string]any{
			"game_id": "test123",
		},
	}

	expected := "engine error in testOperation (phase: exploration): entity not found"
	if err.Error() != expected {
		t.Errorf("Expected error message:\n%s\nGot:\n%s", expected, err.Error())
	}

	// 测试 Unwrap
	unwrapped := err.Unwrap()
	if unwrapped != ErrNotFound {
		t.Errorf("Expected unwrapped error to be ErrNotFound, got %v", unwrapped)
	}
}

// TestEngineErrorWithoutErr 测试没有原始错误的EngineError
func TestEngineErrorWithoutErr(t *testing.T) {
	err := &EngineError{
		Op:    "testOperation",
		Phase: model.PhaseCombat,
	}

	expected := "engine error in testOperation (phase: combat)"
	if err.Error() != expected {
		t.Errorf("Expected error message:\n%s\nGot:\n%s", expected, err.Error())
	}

	// Unwrap 应该返回 nil
	if err.Unwrap() != nil {
		t.Errorf("Expected unwrapped error to be nil, got %v", err.Unwrap())
	}
}
