package engine

import (
	"context"
	"testing"

	"github.com/zwh8800/dnd-core/internal/model"
)

// createTestGame 创建测试用的引擎和游戏
func createTestGame(t *testing.T) (*Engine, model.ID) {
	t.Helper()
	engine, err := New(DefaultConfig())
	if err != nil {
		t.Fatalf("Failed to create engine: %v", err)
	}

	ctx := context.Background()
	game := model.NewGameState("Test Game", "A test game")
	// 设置游戏阶段为exploration，允许大多数操作
	game.Phase = model.PhaseExploration

	err = engine.saveGame(ctx, game)
	if err != nil {
		t.Fatalf("Failed to save game: %v", err)
	}

	return engine, game.ID
}
