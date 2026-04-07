package storage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/zwh8800/dnd-core/pkg/model"
)

// MemoryStore 内存存储实现
type MemoryStore struct {
	mu    sync.RWMutex
	games map[model.ID]*model.GameState
}

// NewMemoryStore 创建新的内存存储
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		games: make(map[model.ID]*model.GameState),
	}
}

// Init 初始化存储（内存存储无需初始化）
func (s *MemoryStore) Init(ctx context.Context) error {
	return nil
}

// SaveGame 保存游戏状态
func (s *MemoryStore) SaveGame(ctx context.Context, game *model.GameState) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 深拷贝
	gameCopy := *game
	s.games[game.ID] = &gameCopy
	return nil
}

// LoadGame 加载游戏状态
func (s *MemoryStore) LoadGame(ctx context.Context, gameID model.ID) (*model.GameState, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	game, ok := s.games[gameID]
	if !ok {
		return nil, &ErrGameNotFound{ID: gameID}
	}

	// 深拷贝返回
	gameCopy := *game
	return &gameCopy, nil
}

// DeleteGame 删除游戏
func (s *MemoryStore) DeleteGame(ctx context.Context, gameID model.ID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.games[gameID]; !ok {
		return &ErrGameNotFound{ID: gameID}
	}

	delete(s.games, gameID)
	return nil
}

// ListGames 列出所有游戏
func (s *MemoryStore) ListGames(ctx context.Context) ([]GameMeta, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]GameMeta, 0, len(s.games))
	for _, game := range s.games {
		pcCount := len(game.PCs)
		maxLevel := 0
		for _, pc := range game.PCs {
			if pc.TotalLevel > maxLevel {
				maxLevel = pc.TotalLevel
			}
		}

		meta := GameMeta{
			ID:           game.ID,
			Name:         game.Name,
			Description:  game.Description,
			Phase:        game.Phase,
			UpdatedAt:    game.UpdatedAt,
			PCCount:      pcCount,
			CurrentLevel: maxLevel,
		}
		result = append(result, meta)
	}

	return result, nil
}

// UpdateGame 原子更新游戏状态
func (s *MemoryStore) UpdateGame(ctx context.Context, gameID model.ID, fn func(*model.GameState) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	game, ok := s.games[gameID]
	if !ok {
		return &ErrGameNotFound{ID: gameID}
	}

	// 深拷贝
	gameCopy := *game
	err := fn(&gameCopy)
	if err != nil {
		return fmt.Errorf("update function failed: %w", err)
	}

	// 更新时间戳
	gameCopy.UpdatedAt = time.Now()
	s.games[gameID] = &gameCopy
	return nil
}

// Close 关闭存储（内存存储无需关闭）
func (s *MemoryStore) Close(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.games = make(map[model.ID]*model.GameState)
	return nil
}

// GameCount 返回游戏数量（调试用）
func (s *MemoryStore) GameCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.games)
}
