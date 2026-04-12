package engine

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/zwh8800/dnd-core/pkg/model"
	"github.com/zwh8800/dnd-core/pkg/storage"
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

// GameInfo 游戏信息摘要
type GameInfo struct {
	ID             model.ID    `json:"id"`                         // 游戏唯一标识
	Name           string      `json:"name"`                       // 游戏名称
	Description    string      `json:"description"`                // 游戏描述
	Phase          model.Phase `json:"phase"`                      // 当前游戏阶段
	CreatedAt      time.Time   `json:"created_at"`                 // 创建时间
	UpdatedAt      time.Time   `json:"updated_at"`                 // 更新时间
	PCCount        int         `json:"pc_count"`                   // 玩家角色数量
	NPCCount       int         `json:"npc_count"`                  // NPC数量
	EnemyCount     int         `json:"enemy_count"`                // 敌人数量
	CompanionCount int         `json:"companion_count"`            // 同伴数量
	SceneCount     int         `json:"scene_count"`                // 场景数量
	CurrentSceneID model.ID    `json:"current_scene_id,omitempty"` // 当前场景ID
	InCombat       bool        `json:"in_combat"`                  // 是否处于战斗中
}

// NewGameRequest 创建游戏请求
type NewGameRequest struct {
	Name        string `json:"name"`        // 游戏名称
	Description string `json:"description"` // 游戏描述
}

// NewGameResult 创建游戏结果
type NewGameResult struct {
	Game *GameInfo `json:"game"` // 新创建的游戏信息
}

// LoadGameRequest 加载游戏请求
type LoadGameRequest struct {
	GameID model.ID `json:"game_id"` // 游戏ID
}

// LoadGameResult 加载游戏结果
type LoadGameResult struct {
	Game *GameInfo `json:"game"` // 加载的游戏信息
}

// SaveGameRequest 保存游戏请求
type SaveGameRequest struct {
	GameID model.ID `json:"game_id"` // 游戏ID
}

// DeleteGameRequest 删除游戏请求
type DeleteGameRequest struct {
	GameID model.ID `json:"game_id"` // 游戏ID
}

// ListGamesRequest 列出游戏请求
type ListGamesRequest struct {
	// 无需参数，使用空结构体以统一API签名
}

type GetGameRequest struct {
	GameID model.ID `json:"game_id"` // 游戏ID
}

type GetGameResult struct {
	Game *GameInfo `json:"game"` // 加载的游戏信息
}

// gameStateToInfo 将 GameState 转换为 GameInfo
func gameStateToInfo(game *model.GameState) *GameInfo {
	info := &GameInfo{
		ID:             game.ID,
		Name:           game.Name,
		Description:    game.Description,
		Phase:          game.Phase,
		UpdatedAt:      game.UpdatedAt,
		PCCount:        len(game.PCs),
		NPCCount:       len(game.NPCs),
		EnemyCount:     len(game.Enemies),
		CompanionCount: len(game.Companions),
		SceneCount:     len(game.Scenes),
		InCombat:       game.Combat != nil,
	}
	if game.CurrentScene != nil && *game.CurrentScene != "" {
		info.CurrentSceneID = *game.CurrentScene
	}
	return info
}

// NewGame 创建一个新的游戏会话
// 在内存中初始化 GameState 对象，并持久化到存储后端。
// 参数:
//
//	ctx - 上下文，用于控制请求生命周期和取消操作
//	req - 创建游戏请求，包含游戏名称和描述
//
// 返回:
//
//	*NewGameResult - 包含新创建游戏的信息摘要
//	error - 创建失败时返回错误（如存储写入失败）
func (e *Engine) NewGame(ctx context.Context, req NewGameRequest) (*NewGameResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game := model.NewGameState(req.Name, req.Description)

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &NewGameResult{
		Game: gameStateToInfo(game),
	}, nil
}

// LoadGame 从存储加载一个已存在的游戏会话
// 根据游戏ID从存储后端读取游戏状态，并返回游戏信息摘要。
// 参数:
//
//	ctx - 上下文，用于控制请求生命周期和取消操作
//	req - 加载游戏请求，包含要加载的游戏ID
//
// 返回:
//
//	*LoadGameResult - 包含加载的游戏信息摘要
//	error - 加载失败时返回错误（如游戏不存在或存储读取失败）
func (e *Engine) LoadGame(ctx context.Context, req LoadGameRequest) (*LoadGameResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	return &LoadGameResult{
		Game: gameStateToInfo(game),
	}, nil
}

// SaveGame 将当前游戏状态持久化到存储后端
// 加载指定游戏状态，更新最后修改时间，并保存到存储。
// 参数:
//
//	ctx - 上下文，用于控制请求生命周期和取消操作
//	req - 保存游戏请求，包含要保存的游戏ID
//
// 返回:
//
//	error - 保存失败时返回错误（如游戏不存在或存储写入失败）
func (e *Engine) SaveGame(ctx context.Context, req SaveGameRequest) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return err
	}

	game.UpdatedAt = time.Now()
	return e.saveGame(ctx, game)
}

// DeleteGame 从存储中删除一个游戏会话
// 永久删除指定游戏及其所有关联数据。如果游戏不存在，返回 ErrNotFound。
// 参数:
//
//	ctx - 上下文，用于控制请求生命周期和取消操作
//	req - 删除游戏请求，包含要删除的游戏ID
//
// 返回:
//
//	error - 删除失败时返回错误（如游戏不存在或存储删除失败）
func (e *Engine) DeleteGame(ctx context.Context, req DeleteGameRequest) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.store.DeleteGame(ctx, req.GameID); err != nil {
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
// 从存储后端读取所有游戏状态，并返回游戏信息摘要列表。
// 参数:
//
//	ctx - 上下文，用于控制请求生命周期和取消操作
//	req - 列出游戏请求，无需参数
//
// 返回:
//
//	[]GameSummary - 包含所有游戏信息摘要的列表
//	error - 列出失败时返回错误（如存储读取失败）
func (e *Engine) ListGames(ctx context.Context, req ListGamesRequest) ([]GameSummary, error) {
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

// GetGame 获取游戏信息
// 获取指定游戏的最新状态，并返回游戏信息摘要。
// 参数:
//
//	ctx - 上下文，用于控制请求生命周期和取消操作
//	req - 获取游戏请求，包含要获取的游戏ID
//
// 返回:
//
//	*GameInfo - 包含游戏信息摘要
//	error - 获取失败时返回错误（如游戏不存在或存储读取失败）
func (e *Engine) GetGame(ctx context.Context, req GetGameRequest) (*GameInfo, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}
	return gameStateToInfo(game), nil
}
