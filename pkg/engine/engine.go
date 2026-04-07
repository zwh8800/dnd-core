package engine

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/zwh8800/dnd-core/pkg/dice"
	"github.com/zwh8800/dnd-core/pkg/model"
	"github.com/zwh8800/dnd-core/pkg/storage"
)

// Engine 是D&D 5e游戏引擎的核心控制器，提供对所有游戏系统的统一访问入口。
// Engine是并发安全的，可以在多个goroutine中同时使用。
// 所有对游戏状态的修改都会自动进行阶段权限验证和状态一致性检查。
type Engine struct {
	mu     sync.RWMutex
	store  storage.Store
	roller *dice.Roller
	config Config
	closed bool
}

// New 创建并初始化一个新的引擎实例
// 参数:
//
//	cfg - 引擎配置，可以传入DefaultConfig()使用默认配置
//
// 返回:
//
//	*Engine - 初始化完成的引擎实例
//	error - 初始化过程中可能发生的错误（如存储初始化失败）
//
// 使用场景: 应用程序启动时调用一次
func New(cfg Config) (*Engine, error) {
	// 使用默认存储（如果未指定）
	if cfg.Storage == nil {
		cfg.Storage = storage.NewMemoryStore()
	}

	// 初始化存储
	if err := cfg.Storage.Init(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}

	// 创建骰子投掷器
	roller := dice.New(cfg.DiceSeed)

	return &Engine{
		store:  cfg.Storage,
		roller: roller,
		config: cfg,
	}, nil
}

// DefaultConfig 返回使用内存存储的默认配置
// 返回:
//
//	Config - 可直接用于New()的默认配置
//
// 使用场景: 快速启动或测试时调用
func DefaultConfig() Config {
	return Config{
		Storage:  storage.NewMemoryStore(),
		DiceSeed: 0,
		DataPath: "",
	}
}

// Close 释放引擎占用的所有资源，包括存储后端连接
// 返回:
//
//	error - 关闭过程中可能发生的错误
//
// 使用场景: 应用程序关闭时调用
func (e *Engine) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.closed {
		return nil
	}

	e.closed = true
	return e.store.Close(context.Background())
}

// getStore 获取存储后端
func (e *Engine) getStore() storage.Store {
	return e.store
}

// getRoller 获取骰子投掷器
func (e *Engine) getRoller() *dice.Roller {
	return e.roller
}

// loadGame 加载游戏状态（内部方法，不获取锁）
func (e *Engine) loadGame(ctx context.Context, gameID model.ID) (*model.GameState, error) {
	game, err := e.store.LoadGame(ctx, gameID)
	if err != nil {
		var notFound *storage.ErrGameNotFound
		if errors.As(err, &notFound) {
			return nil, ErrNotFound
		}
		return nil, &EngineError{
			Op:  "loadGame",
			Err: err,
		}
	}
	return game, nil
}

// saveGame 保存游戏状态（内部方法，不获取锁）
func (e *Engine) saveGame(ctx context.Context, game *model.GameState) error {
	if err := e.store.SaveGame(ctx, game); err != nil {
		return &EngineError{
			Op:    "saveGame",
			Err:   err,
			Phase: game.Phase,
		}
	}
	return nil
}

// NewTestEngine creates an engine instance for testing
// Uses a fixed seed to ensure reproducible tests
func NewTestEngine(t *testing.T) *Engine {
	cfg := DefaultConfig()
	cfg.DiceSeed = 42
	e, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		e.Close()
	})
	return e
}
