# pkg/engine API 统一改造计划

## Context

当前 pkg/engine 目录下的 84 个公开 API 方法存在参数和返回值不统一的问题：
- 部分方法直接使用 `*model.PlayerCharacter`、`*model.GameState` 等 internal/model 类型作为参数或返回值
- 部分方法已经将参数封装为 Request 结构体（如 AbilityCheckRequest），但仍有大量方法未封装
- 违反了"不在公开 API 中暴露 internal 包类型"的设计原则

本次改造将所有 API 方法的入参封装为 XXXRequest 结构体，返回值封装为 XXXResult 结构体，确保公开 API 不直接暴露 internal/model 类型。

## 改造规则（用户确认）

1. **全部封装**：所有 API 的入参和返回值都封装为 Request/Result 结构体
2. **不保留旧签名**：删除旧方法签名，只保留新签名（破坏性变更）
3. **统一命名**：所有 Input 后缀改为 Request（如 QuestInput → CreateQuestRequest）
4. **辅助结构体改造**：ActorFilter、ActorUpdate、SceneUpdate 等也改为 XXXRequest 格式
5. **中文注释**：每个结构体和字段都必须有详细的中文注释
6. **JSON 标签**：所有字段添加 snake_case 的 json 标签

## 改造清单（按文件）

### 1. game.go（5个方法）

**新增结构体：**
- `GameInfo` - 游戏信息摘要（替代 *model.GameState）
- `NewGameRequest/NewGameResult`
- `LoadGameRequest/LoadGameResult`
- `SaveGameRequest`
- `DeleteGameRequest`
- `ListGamesRequest`

**改造后签名示例：**
```go
func (e *Engine) NewGame(ctx context.Context, req NewGameRequest) (*NewGameResult, error)
func (e *Engine) LoadGame(ctx context.Context, req LoadGameRequest) (*LoadGameResult, error)
```

### 2. actor.go（13个方法）

**新增结构体：**
- `PlayerCharacterInput` - 替代 *model.PlayerCharacter 入参
- `NPCInput` - 替代 *model.NPC 入参
- `EnemyInput` - 替代 *model.Enemy 入参
- `CompanionInput` - 替代 *model.Companion 入参
- `AbilityScoresInput` - 属性值输入
- `ActorInfo` - 角色基本信息（替代 model.ActorSnapshot）
- `PlayerCharacterInfo` - 玩家角色完整信息
- `CreatePCRequest/CreatePCResult`
- `CreateNPCRequest/CreateNPCResult`
- `CreateEnemyRequest/CreateEnemyResult`
- `CreateCompanionRequest/CreateCompanionResult`
- `GetActorRequest/GetActorResult`
- `GetPCRequest/GetPCResult`
- `UpdateActorRequest`（改造 ActorUpdate）
- `RemoveActorRequest`
- `ListActorsRequest/ListActorsResult`（改造 ActorFilter）
- `AddExperienceRequest/AddExperienceResult`
- `LevelUpRequest`
- `ShortRestRequest`
- `StartLongRestRequest`
- `EndLongRestRequest`

### 3. check.go（5个方法）

**已有结构体（无需改造）：**
- AbilityCheckRequest/Result ✓
- SkillCheckRequest/Result ✓
- SavingThrowRequest/Result ✓

**新增结构体：**
- `GetPassivePerceptionRequest/GetPassivePerceptionResult`

### 4. combat.go（11个方法）

**新增结构体：**
- `CombatInfo` - 替代 *model.CombatState
- `CombatantEntryInfo` - 战斗者条目信息
- `StartCombatRequest/StartCombatResult`
- `StartCombatWithSurpriseRequest/StartCombatWithSurpriseResult`
- `EndCombatRequest`
- `GetCurrentCombatRequest/GetCurrentCombatResult`
- `NextTurnRequest/NextTurnResult`
- `GetCurrentTurnRequest`
- `ExecuteActionRequest`（整合 ActionInput）
- `ExecuteAttackRequest`（整合 AttackInput）
- `ExecuteDamageRequest`（整合 DamageInput）
- `ExecuteHealingRequest`
- `MoveActorRequest`

**注意：** ActionInput、AttackInput、DamageInput 保留作为 Request 的子结构

### 5. dice.go（5个方法）

**已有结构体：**
- RollRequest/Result ✓
- RollAdvantageRequest ✓
- RollDisadvantageRequest ✓

**新增结构体：**
- `RollAbilityRequest`
- `RollHitDiceRequest`

### 6. inventory.go（9个方法）

**新增结构体：**
- `ItemInput` - 替代 *model.Item 入参
- `ItemSummary` - 替代 *model.Item 返回值
- `AddItemRequest`
- `RemoveItemRequest`
- `GetInventoryRequest`
- `EquipItemRequest`
- `UnequipItemRequest`
- `GetEquipmentRequest`
- `AttuneItemRequest`
- `TransferItemRequest`
- `AddCurrencyRequest`

**注意：** InventoryInfo 和 EquipmentInfo 中的 *model.Item 改为 *ItemSummary

### 7. phase.go（3个方法）

**新增结构体：**
- `SetPhaseRequest`
- `GetPhaseRequest/GetPhaseResult`
- `GetAllowedOperationsRequest`

### 8. quest.go（9个方法）

**新增结构体：**
- `CreateQuestRequest`（整合 QuestInput）
- `GetQuestRequest`
- `ListQuestsRequest/ListQuestsResult`
- `AcceptQuestRequest`
- `UpdateQuestObjectiveRequest`
- `CompleteQuestRequest`
- `FailQuestRequest`
- `DeleteQuestRequest`
- `GetActorQuestsRequest/GetActorQuestsResult`
- `GetQuestGiverQuestsRequest/GetQuestGiverQuestsResult`

**注意：** QuestResult 中的 *model.Quest 改为 *QuestInfo

### 9. scene.go（13+个方法）

**新增结构体：**
- `CreateSceneRequest`
- `GetSceneRequest`
- `UpdateSceneRequest`（整合 SceneUpdate）
- `DeleteSceneRequest`
- `ListScenesRequest/ListScenesResult`
- `SetCurrentSceneRequest`
- `GetCurrentSceneRequest`
- `AddSceneConnectionRequest`
- `RemoveSceneConnectionRequest`
- `MoveActorToSceneRequest`
- `GetSceneActorsRequest/GetSceneActorsResult`
- `AddItemToSceneRequest`
- `RemoveItemFromSceneRequest`
- `GetSceneItemsRequest/GetSceneItemsResult`

### 10. spell.go（6个方法）

**新增结构体：**
- `CastSpellRequest`（整合 SpellInput）
- `GetSpellSlotsRequest`
- `PrepareSpellsRequest`
- `LearnSpellRequest`
- `ConcentrationCheckRequest`
- `EndConcentrationRequest`

### 11. state.go（3个方法）

**新增结构体：**
- `GetStateSummaryRequest`
- `GetActorSheetRequest`
- `GetCombatSummaryRequest`

## 改造步骤

### 阶段一：定义结构体（步骤 1-11）
按文件顺序添加所有 Request/Result/Input/Info 结构体及其字段注释，不改动方法签名

### 阶段二：改造方法签名（步骤 12-22）
从简单到复杂逐文件改造：
1. dice.go（最简单，5个方法）
2. check.go（3个已有Request，只改GetPassivePerception）
3. game.go（5个方法，新增 gameStateToInfo 辅助函数）
4. phase.go（3个方法）
5. state.go（3个方法）
6. actor.go（13个方法，涉及大量 Input→model 转换）
7. combat.go（11个方法，涉及 CombatInfo 转换）
8. inventory.go（9个方法）
9. quest.go（9个方法）
10. scene.go（13+个方法）
11. spell.go（6个方法）

### 阶段三：更新测试（步骤 23）
更新所有 `*_test.go` 文件以使用新的 Request 结构体

### 阶段四：最终验证（步骤 24-26）
- 运行 `go build ./...`
- 运行 `go test ./pkg/engine/...`
- 检查是否有遗漏的 model 类型暴露

## 关键文件

- `/Users/wastecat/code/go/dnd-core/pkg/engine/actor.go` - 最复杂，13个方法，大量 Input/Info 转换
- `/Users/wastecat/code/go/dnd-core/pkg/engine/combat.go` - 涉及 CombatInfo 封装
- `/Users/wastecat/code/go/dnd-core/pkg/engine/game.go` - GameInfo 是多个 Result 的基础
- `/Users/wastecat/code/go/dnd-core/pkg/engine/dice.go` - 最简单，建议首先改造验证模式

## 注意事项

1. **Input→model 转换**：方法内部需添加 Input 结构体到 model 类型的转换逻辑
2. **model→Info 转换**：返回值需将 model 类型转换为 Info 结构体
3. **model.ID 等基础类型**：可直接在 Request/Result 中使用，无需封装
4. **JSON 标签**：所有字段使用 `snake_case`，可选字段添加 `omitempty`
5. **破坏性变更**：不保留旧签名，外部调用方需同步修改

## 验证方法

1. 编译检查：`go build ./pkg/engine/...`
2. 运行测试：`go test ./pkg/engine/...`
3. 类型检查：确保公开 API 签名中不包含 `internal/model` 的复杂类型（`model.ID`、`model.Ability` 等基础枚举类型除外）
