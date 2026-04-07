package engine

import (
	"context"
	"testing"

	"github.com/zwh8800/dnd-core/internal/model"
)

// TestNewGame 测试创建新游戏
func TestNewGame(t *testing.T) {
	engine, _ := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	game, err := engine.NewGame(ctx, "New Adventure", "A new adventure begins")
	if err != nil {
		t.Fatalf("Failed to create new game: %v", err)
	}

	if game.Name != "New Adventure" {
		t.Errorf("Expected name 'New Adventure', got %s", game.Name)
	}
	if game.Description != "A new adventure begins" {
		t.Errorf("Expected description 'A new adventure begins', got %s", game.Description)
	}
	if game.Phase != model.PhaseCharacterCreation {
		t.Errorf("Expected phase %s, got %s", model.PhaseCharacterCreation, game.Phase)
	}
	if game.ID == "" {
		t.Error("Expected game ID to be set")
	}
}

// TestLoadGame 测试加载游戏
func TestLoadGame(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	game, err := engine.LoadGame(ctx, gameID)
	if err != nil {
		t.Fatalf("Failed to load game: %v", err)
	}

	if game.ID != gameID {
		t.Errorf("Expected game ID %s, got %s", gameID, game.ID)
	}
	if game.Name != "Test Game" {
		t.Errorf("Expected name 'Test Game', got %s", game.Name)
	}
}

// TestLoadGameNotFound 测试加载不存在的游戏
func TestLoadGameNotFound(t *testing.T) {
	engine, _ := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()
	fakeID := model.ID("00000000000000000000000000")

	_, err := engine.LoadGame(ctx, fakeID)
	if err == nil {
		t.Fatal("Expected error when loading non-existent game, got nil")
	}
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

// TestSaveGame 测试保存游戏
func TestSaveGame(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	err := engine.SaveGame(ctx, gameID)
	if err != nil {
		t.Fatalf("Failed to save game: %v", err)
	}
}

// TestSaveGameNotFound 测试保存不存在的游戏
func TestSaveGameNotFound(t *testing.T) {
	engine, _ := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()
	fakeID := model.ID("00000000000000000000000000")

	err := engine.SaveGame(ctx, fakeID)
	if err == nil {
		t.Fatal("Expected error when saving non-existent game, got nil")
	}
}

// TestDeleteGame 测试删除游戏
func TestDeleteGame(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	err := engine.DeleteGame(ctx, gameID)
	if err != nil {
		t.Fatalf("Failed to delete game: %v", err)
	}

	// 验证游戏已删除
	_, err = engine.LoadGame(ctx, gameID)
	if err == nil {
		t.Fatal("Expected error when loading deleted game, got nil")
	}
}

// TestDeleteGameNotFound 测试删除不存在的游戏
func TestDeleteGameNotFound(t *testing.T) {
	engine, _ := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()
	fakeID := model.ID("00000000000000000000000000")

	err := engine.DeleteGame(ctx, fakeID)
	if err == nil {
		t.Fatal("Expected error when deleting non-existent game, got nil")
	}
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

// TestListGames 测试列出游戏
func TestListGames(t *testing.T) {
	engine, _ := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	// 创建多个游戏
	for i := 0; i < 3; i++ {
		_, err := engine.NewGame(ctx, "Game", "Description")
		if err != nil {
			t.Fatalf("Failed to create game %d: %v", i, err)
		}
	}

	// 列出游戏
	games, err := engine.ListGames(ctx)
	if err != nil {
		t.Fatalf("Failed to list games: %v", err)
	}

	// 应该至少有4个游戏（1个createTestGame创建的 + 3个新创建的）
	if len(games) < 4 {
		t.Errorf("Expected at least 4 games, got %d", len(games))
	}

	// 验证按更新时间排序
	for i := 1; i < len(games); i++ {
		if games[i-1].UpdatedAt.Before(games[i].UpdatedAt) {
			t.Error("Games should be sorted by updated time descending")
			break
		}
	}
}
