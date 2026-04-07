package storage

import (
	"context"
	"time"

	"github.com/zwh8800/dnd-core/internal/model"
)

// GameMeta 游戏元数据
type GameMeta struct {
	ID           model.ID    `json:"id"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	Phase        model.Phase `json:"phase"`
	UpdatedAt    time.Time   `json:"updated_at"`
	PCCount      int         `json:"pc_count"`
	CurrentLevel int         `json:"current_level"`
}

// Store 定义游戏状态存储接口
type Store interface {
	// Init 初始化存储
	Init(ctx context.Context) error

	// SaveGame 保存游戏状态
	SaveGame(ctx context.Context, game *model.GameState) error

	// LoadGame 加载游戏状态
	LoadGame(ctx context.Context, gameID model.ID) (*model.GameState, error)

	// DeleteGame 删除游戏
	DeleteGame(ctx context.Context, gameID model.ID) error

	// ListGames 列出所有游戏
	ListGames(ctx context.Context) ([]GameMeta, error)

	// UpdateGame 原子更新游戏状态
	UpdateGame(ctx context.Context, gameID model.ID, fn func(*model.GameState) error) error

	// Close 关闭存储连接
	Close(ctx context.Context) error
}

// ErrGameNotFound 游戏未找到
type ErrGameNotFound struct {
	ID model.ID
}

func (e *ErrGameNotFound) Error() string {
	return "game not found: " + string(e.ID)
}

// ErrStorageError 存储错误
type ErrStorageError struct {
	Op  string
	Err error
}

func (e *ErrStorageError) Error() string {
	if e.Err != nil {
		return "storage error (" + e.Op + "): " + e.Err.Error()
	}
	return "storage error (" + e.Op + ")"
}

func (e *ErrStorageError) Unwrap() error {
	return e.Err
}
