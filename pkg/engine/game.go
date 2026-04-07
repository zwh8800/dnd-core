package engine

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/zwh8800/dnd-core/internal/model"
	"github.com/zwh8800/dnd-core/internal/storage"
)

// GameSummary 游戏摘要
type GameSummary struct {
	ID           model.ID    `json:"id"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	Phase        model.Phase `json:"phase"`
	UpdatedAt    time.Time   `json:"updated_at"`
	PCCount      int         `json:"pc_count"`
	CurrentLevel int         `json:"current_level"`
}

// NewGame 创建一个新的游戏会话
func (e *Engine) NewGame(ctx context.Context, name, description string) (*model.GameState, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game := model.NewGameState(name, description)

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	gameCopy := *game
	return &gameCopy, nil
}

// LoadGame 从存储加载一个已存在的游戏会话
func (e *Engine) LoadGame(ctx context.Context, gameID model.ID) (*model.GameState, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	gameCopy := *game
	return &gameCopy, nil
}

// SaveGame 将当前游戏状态持久化到存储后端
func (e *Engine) SaveGame(ctx context.Context, gameID model.ID) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, gameID)
	if err != nil {
		return err
	}

	game.UpdatedAt = time.Now()
	return e.saveGame(ctx, game)
}

// DeleteGame 从存储中删除一个游戏会话
func (e *Engine) DeleteGame(ctx context.Context, gameID model.ID) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.store.DeleteGame(ctx, gameID); err != nil {
		var notFound *storage.ErrGameNotFound
		if errors.As(err, &notFound) {
			return ErrNotFound
		}
		return &EngineError{
			Op:  "deleteGame",
			Err: err,
		}
	}
	return nil
}

// ListGames 列出所有可用的游戏会话
func (e *Engine) ListGames(ctx context.Context) ([]GameSummary, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	metas, err := e.store.ListGames(ctx)
	if err != nil {
		return nil, &EngineError{
			Op:  "listGames",
			Err: err,
		}
	}

	result := make([]GameSummary, 0, len(metas))
	for _, m := range metas {
		result = append(result, GameSummary{
			ID:           m.ID,
			Name:         m.Name,
			Description:  m.Description,
			Phase:        m.Phase,
			UpdatedAt:    m.UpdatedAt,
			PCCount:      m.PCCount,
			CurrentLevel: m.CurrentLevel,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].UpdatedAt.After(result[j].UpdatedAt)
	})

	return result, nil
}
