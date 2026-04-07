# pkg/engine API 设计规范

本文档定义了 `pkg/engine` 包中所有对外 API 的设计规范和实现模式。

## 核心设计原则

### 1. Engine 方法模式

所有对外 API 必须是 `Engine` 结构体的方法，遵循以下签名：

```go
// 成功返回 Result，失败返回 error
func (e *Engine) OperationName(ctx context.Context, req OperationRequest) (*OperationResult, error)

// 或者仅返回 error（简单操作）
func (e *Engine) OperationName(ctx context.Context, req OperationRequest) error
```

**错误示例**：
```go
// ❌ 不要使用独立函数
func SelectFeat(pc *model.PlayerCharacter, featID string) error
```

**正确示例**：
```go
// ✅ 使用 Engine 方法
func (e *Engine) SelectFeat(ctx context.Context, req SelectFeatRequest) (*SelectFeatResult, error)
```

### 2. Request/Result 结构体

每个 API 操作都需要定义对应的 Request 和 Result 结构体。

#### Request 结构体规范

```go
// OperationNameRequest 操作名称请求
// 用于描述该操作的用途（单行注释）
type OperationNameRequest struct {
    GameID  model.ID `json:"game_id"`  // 游戏会话ID（必填）
    ActorID model.ID `json:"actor_id"` // 角色ID（如果适用）
    // 其他操作特定参数...
}
```

**命名规则**：
- 结构体名：`<操作名>Request`，如 `CreatePCRequest`、`GetActorRequest`
- 所有字段必须导出（大写开头）
- 必须包含 `json` tag
- 字段注释使用中文，说明字段用途和是否必填

**通用字段**：
- `GameID model.ID` - 所有请求都必须包含游戏会话ID

#### Result 结构体规范

```go
// OperationNameResult 操作名称结果
// 用于描述该操作返回的数据
type OperationNameResult struct {
    Actor    *ActorInfo  `json:"actor"`              // 创建的角色信息
    LeveledUp bool       `json:"leveled_up"`         // 是否升级
    // 其他操作特定结果...
}
```

**命名规则**：
- 结构体名：`<操作名>Result`
- 返回的数据应该封装为 Info 结构体（不直接暴露 model 层）
- 包含人类可读的消息字段（如适用）

### 3. 并发控制

所有 API 方法必须实现正确的并发控制：

```go
func (e *Engine) OperationName(ctx context.Context, req OperationRequest) (*OperationResult, error) {
    // 写操作使用写锁
    e.mu.Lock()
    defer e.mu.Unlock()

    // 读操作使用读锁
    e.mu.RLock()
    defer e.mu.RUnlock()
}
```

**规则**：
- 修改游戏状态的操作：`e.mu.Lock()`
- 仅查询状态的操作：`e.mu.RLock()`

### 4. 状态管理

所有状态操作必须通过引擎的 load/save 方法：

```go
func (e *Engine) OperationName(ctx context.Context, req OperationRequest) (*OperationResult, error) {
    e.mu.Lock()
    defer e.mu.Unlock()

    // 1. 加载游戏状态
    game, err := e.loadGame(ctx, req.GameID)
    if err != nil {
        return nil, err
    }

    // 2. 验证权限
    if err := e.checkPermission(game.Phase, OpOperationName); err != nil {
        return nil, err
    }

    // 3. 执行业务逻辑
    // ...

    // 4. 保存游戏状态
    if err := e.saveGame(ctx, game); err != nil {
        return nil, err
    }

    // 5. 返回结果
    return &OperationResult{...}, nil
}
```

### 5. 权限验证

每个修改操作必须在加载游戏后立即验证权限：

```go
if err := e.checkPermission(game.Phase, OpYourOperation); err != nil {
    return nil, err
}
```

**添加新权限**：
1. 在 `phase.go` 中定义操作常量：
   ```go
   const (
       OpYourOperation Operation = "your_operation"
   )
   ```

2. 在 `phasePermissions` map 中添加权限：
   ```go
   model.PhaseExploration: {
       // ... 其他权限
       OpYourOperation: true,
   },
   ```

### 6. Info 结构体封装

不直接返回 `model` 层的结构体，使用 Info 结构体封装：

```go
// ActorInfo 角色基本信息
type ActorInfo struct {
    ID         model.ID        `json:"id"`
    Type       model.ActorType `json:"type"`
    Name       string          `json:"name"`
    HitPoints  model.HitPoints `json:"hit_points"`
    // ...
}

// 辅助转换函数
func actorToInfo(actor *model.Actor, actorType model.ActorType, name string) *ActorInfo {
    return &ActorInfo{
        ID:         actor.ID,
        Type:       actorType,
        Name:       name,
        HitPoints:  actor.HitPoints,
        // ...
    }
}
```

### 7. 错误处理

遵循统一的错误处理模式：

```go
// 资源不存在
if !ok {
    return nil, ErrNotFound
}

// 参数验证失败
if req.SomeField == "" {
    return nil, fmt.Errorf("some_field is required")
}

// 业务逻辑错误
if !condition {
    return nil, fmt.Errorf("specific business error message")
}
```

**预定义错误**：
- `ErrNotFound` - 资源不存在
- `ErrInvalidState` - 无效状态
- `ErrPhaseNotAllowed` - 当前阶段不允许该操作

### 8. 注释规范

#### 结构体注释
```go
// OperationNameRequest 操作名称请求
// 用于描述该操作的具体用途（多行说明）
type OperationNameRequest struct {
```

#### 字段注释
```go
type OperationNameRequest struct {
    GameID  model.ID `json:"game_id"`  // 游戏会话ID
    ActorID model.ID `json:"actor_id"` // 角色ID（可选）
}
```

#### 方法注释
```go
// OperationName 执行某个操作
// 参数:
//
//	ctx - 上下文
//	req - 请求参数
//
// 返回:
//
//	*OperationResult - 操作结果
//	error - 可能发生的错误
func (e *Engine) OperationName(ctx context.Context, req OperationNameRequest) (*OperationResult, error) {
```

## 完整示例

### 定义结构体

```go
// SelectFeatRequest 选择专长请求
type SelectFeatRequest struct {
    GameID model.ID `json:"game_id"` // 游戏会话ID
    PCID   model.ID `json:"pc_id"`   // 玩家角色ID
    FeatID string   `json:"feat_id"` // 专长ID
}

// SelectFeatResult 选择专长结果
type SelectFeatResult struct {
    Feats []FeatInfo `json:"feats"` // 角色当前专长列表
}
```

### 实现方法

```go
// SelectFeat 为角色选择并获得一个专长
func (e *Engine) SelectFeat(ctx context.Context, req SelectFeatRequest) (*SelectFeatResult, error) {
    e.mu.Lock()
    defer e.mu.Unlock()

    // 加载游戏
    game, err := e.loadGame(ctx, req.GameID)
    if err != nil {
        return nil, err
    }

    // 验证权限
    if err := e.checkPermission(game.Phase, OpSelectFeat); err != nil {
        return nil, err
    }

    // 查找角色
    pc, ok := game.PCs[req.PCID]
    if !ok {
        return nil, ErrNotFound
    }

    // 业务逻辑
    feat, exists := data.GlobalRegistry.GetFeat(req.FeatID)
    if !exists {
        return nil, fmt.Errorf("feat not found: %s", req.FeatID)
    }

    // 修改状态
    pc.Feats = append(pc.Feats, model.FeatInstance{
        FeatID: req.FeatID,
        Source: model.FeatSourceLevelUp,
    })

    // 保存游戏
    if err := e.saveGame(ctx, game); err != nil {
        return nil, err
    }

    // 返回结果
    return &SelectFeatResult{
        Feats: []FeatInfo{featToInfo(feat)},
    }, nil
}
```

### 添加权限

```go
// phase.go
const (
    OpSelectFeat Operation = "select_feat"
)

var phasePermissions = map[model.Phase]map[Operation]bool{
    model.PhaseCharacterCreation: {
        // ...
        OpSelectFeat: true,
    },
    model.PhaseExploration: {
        // ...
        OpSelectFeat: true,
    },
}
```

## 文件组织

- `actor.go` - 角色相关 API
- `combat.go` - 战斗相关 API
- `spell.go` - 法术相关 API
- `inventory.go` - 库存相关 API
- `feat.go` - 专长相关 API
- `phase.go` - 阶段和权限定义
- `scene.go` - 场景相关 API

每个文件包含：
1. Request/Result 结构体定义
2. Info 结构体定义
3. Engine 方法实现
4. 辅助转换函数

## 禁止的行为

1. **不要在 API 方法中直接操作 model 层结构体** - 使用 Info 结构体封装
2. **不要忘记加锁** - 所有方法必须正确实现并发控制
3. **不要忘记保存状态** - 修改状态后必须调用 `e.saveGame()`
4. **不要跳过权限验证** - 所有修改操作必须验证阶段权限
5. **不要使用独立的函数** - 所有 API 必须是 Engine 的方法
6. **不要返回裸的 model 类型** - 使用封装的 Info 结构体

## 检查清单

实现新 API 时，确保：

- [ ] 定义了 `<Operation>Request` 结构体
- [ ] 定义了 `<Operation>Result` 结构体（如适用）
- [ ] 方法是 `(e *Engine)` 的方法
- [ ] 方法签名为 `func (e *Engine) Op(ctx context.Context, req Req) (*Result, error)`
- [ ] 正确使用 `e.mu.Lock()` 或 `e.mu.RLock()`
- [ ] 使用 `e.loadGame()` 加载状态
- [ ] 使用 `e.checkPermission()` 验证权限
- [ ] 使用 `e.saveGame()` 保存修改
- [ ] 所有字段都有 `json` tag
- [ ] 所有字段都有中文注释
- [ ] 结构体和方法有完整的注释说明
- [ ] 在 `phase.go` 中添加了权限定义
