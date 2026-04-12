# D&D 5e 游戏引擎 API 文档

本文档分类总结了 `pkg/engine` 包中所有对外提供的公共 API。

---

## 1. 引擎初始化

### New
```go
func New(cfg Config) (*Engine, error)
```
创建并初始化一个新的引擎实例

### DefaultConfig
```go
func DefaultConfig() Config
```
返回使用内存存储的默认配置

### Close
```go
func (e *Engine) Close() error
```
释放引擎占用的所有资源

### NewTestEngine
```go
func NewTestEngine(t *testing.T) *Engine
```
创建用于测试的引擎实例（内部测试使用）

---

## 2. 游戏会话管理

### NewGame
```go
func (e *Engine) NewGame(ctx context.Context, req NewGameRequest) (*NewGameResult, error)
```
创建一个新的游戏会话

### LoadGame
```go
func (e *Engine) LoadGame(ctx context.Context, req LoadGameRequest) (*LoadGameResult, error)
```
从存储加载一个已存在的游戏会话

### SaveGame
```go
func (e *Engine) SaveGame(ctx context.Context, req SaveGameRequest) error
```
将当前游戏状态持久化到存储后端

### DeleteGame
```go
func (e *Engine) DeleteGame(ctx context.Context, req DeleteGameRequest) error
```
从存储中删除一个游戏会话

### ListGames
```go
func (e *Engine) ListGames(ctx context.Context, req ListGamesRequest) ([]GameSummary, error)
```
列出所有可用的游戏会话

---

## 3. 角色管理

### CreatePC
```go
func (e *Engine) CreatePC(ctx context.Context, req CreatePCRequest) (*CreatePCResult, error)
```
创建一个新的玩家角色（PC）

### CreateNPC
```go
func (e *Engine) CreateNPC(ctx context.Context, req CreateNPCRequest) (*CreateNPCResult, error)
```
创建一个新的非玩家角色（NPC）

### CreateEnemy
```go
func (e *Engine) CreateEnemy(ctx context.Context, req CreateEnemyRequest) (*CreateEnemyResult, error)
```
创建一个新的敌人/怪物

### CreateCompanion
```go
func (e *Engine) CreateCompanion(ctx context.Context, req CreateCompanionRequest) (*CreateCompanionResult, error)
```
创建一个同伴角色

### GetActor
```go
func (e *Engine) GetActor(ctx context.Context, req GetActorRequest) (*GetActorResult, error)
```
获取任意类型角色的基本信息

### GetPC
```go
func (e *Engine) GetPC(ctx context.Context, req GetPCRequest) (*GetPCResult, error)
```
获取玩家角色的完整数据

### UpdateActor
```go
func (e *Engine) UpdateActor(ctx context.Context, req UpdateActorRequest) error
```
更新角色的部分状态

### RemoveActor
```go
func (e *Engine) RemoveActor(ctx context.Context, req RemoveActorRequest) error
```
从游戏中移除一个角色

### ListActors
```go
func (e *Engine) ListActors(ctx context.Context, req ListActorsRequest) (*ListActorsResult, error)
```
列出游戏中的所有角色，支持按条件过滤

---

## 4. 角色升级与经验

### AddExperience
```go
func (e *Engine) AddExperience(ctx context.Context, req AddExperienceRequest) (*AddExperienceResult, error)
```
为玩家角色添加经验值

### LevelUp
```go
func (e *Engine) LevelUp(ctx context.Context, req LevelUpRequest) (*LevelUpResult, error)
```
手动触发玩家角色升级

---

## 5. 休息系统

### ShortRest
```go
func (e *Engine) ShortRest(ctx context.Context, req ShortRestRequest) (*RestResult, error)
```
为指定角色执行短休

### StartLongRest
```go
func (e *Engine) StartLongRest(ctx context.Context, req StartLongRestRequest) (*RestResult, error)
```
开始长休过程

### EndLongRest
```go
func (e *Engine) EndLongRest(ctx context.Context, req EndLongRestRequest) (*RestResult, error)
```
结束长休并应用恢复效果

---

## 6. 战斗系统

### StartCombat
```go
func (e *Engine) StartCombat(ctx context.Context, req StartCombatRequest) (*StartCombatResult, error)
```
开始一场战斗遭遇

### StartCombatWithSurprise
```go
func (e *Engine) StartCombatWithSurprise(ctx context.Context, req StartCombatWithSurpriseRequest) (*StartCombatWithSurpriseResult, error)
```
开始带突袭判定的战斗

### EndCombat
```go
func (e *Engine) EndCombat(ctx context.Context, req EndCombatRequest) error
```
结束当前战斗

### GetCurrentCombat
```go
func (e *Engine) GetCurrentCombat(ctx context.Context, req GetCurrentCombatRequest) (*GetCurrentCombatResult, error)
```
获取当前活跃的战斗状态

### NextTurn
```go
func (e *Engine) NextTurn(ctx context.Context, req NextTurnRequest) (*NextTurnResult, error)
```
推进到下一个角色的回合

### GetCurrentTurn
```go
func (e *Engine) GetCurrentTurn(ctx context.Context, req GetCurrentTurnRequest) (*TurnInfo, error)
```
获取当前回合的信息

### ExecuteAction
```go
func (e *Engine) ExecuteAction(ctx context.Context, req ExecuteActionRequest) (*ExecuteActionResult, error)
```
执行一个动作（冲刺、脱离、闪避等）

### ExecuteAttack
```go
func (e *Engine) ExecuteAttack(ctx context.Context, req ExecuteAttackRequest) (*ExecuteAttackResult, error)
```
执行攻击动作

### ExecuteDamage
```go
func (e *Engine) ExecuteDamage(ctx context.Context, req ExecuteDamageRequest) (*ExecuteDamageResult, error)
```
直接对角色施加伤害

### ExecuteHealing
```go
func (e *Engine) ExecuteHealing(ctx context.Context, req ExecuteHealingRequest) (*ExecuteHealingResult, error)
```
对角色进行治疗

### MoveActor
```go
func (e *Engine) MoveActor(ctx context.Context, req MoveActorRequest) (*MoveActorResult, error)
```
在场景中移动角色

---

## 7. 法术系统

### CastSpell
```go
func (e *Engine) CastSpell(ctx context.Context, req CastSpellRequest) (*SpellResult, error)
```
执行施法动作

### GetSpellSlots
```go
func (e *Engine) GetSpellSlots(ctx context.Context, req GetSpellSlotsRequest) (*GetSpellSlotsResult, error)
```
获取施法者的法术位状态

### PrepareSpells
```go
func (e *Engine) PrepareSpells(ctx context.Context, req PrepareSpellsRequest) error
```
准备法术（针对准备型施法者）

### LearnSpell
```go
func (e *Engine) LearnSpell(ctx context.Context, req LearnSpellRequest) error
```
学习新法术（针对已知型施法者）

### ConcentrationCheck
```go
func (e *Engine) ConcentrationCheck(ctx context.Context, req ConcentrationCheckRequest) (*ConcentrationResult, error)
```
进行专注检定

### EndConcentration
```go
func (e *Engine) EndConcentration(ctx context.Context, req EndConcentrationRequest) error
```
主动结束专注

### CastSpellRitual
```go
func (e *Engine) CastSpellRitual(ctx context.Context, req CastSpellRitualRequest) (*SpellResult, error)
```
仪式施法

### GetPactMagicSlots
```go
func (e *Engine) GetPactMagicSlots(ctx context.Context, req GetPactMagicSlotsRequest) (*GetSpellSlotsResult, error)
```
获取魔契师 Pact Magic 法术位

### RestorePactMagicSlots
```go
func (e *Engine) RestorePactMagicSlots(ctx context.Context, req RestorePactMagicSlotsRequest) error
```
恢复魔契师法术位

### IsConcentrating
```go
func (e *Engine) IsConcentrating(ctx context.Context, req IsConcentratingRequest) (*IsConcentratingResult, error)
```
检查是否正在专注

### GetConcentrationSpell
```go
func (e *Engine) GetConcentrationSpell(ctx context.Context, req GetConcentrationSpellRequest) (*SpellResult, error)
```
获取当前专注的法术详情

---

## 8. 检定系统

### PerformAbilityCheck
```go
func (e *Engine) PerformAbilityCheck(ctx context.Context, req AbilityCheckRequest) (*AbilityCheckResult, error)
```
执行属性检定

### PerformSkillCheck
```go
func (e *Engine) PerformSkillCheck(ctx context.Context, req SkillCheckRequest) (*SkillCheckResult, error)
```
执行技能检定

### PerformSavingThrow
```go
func (e *Engine) PerformSavingThrow(ctx context.Context, req SavingThrowRequest) (*SavingThrowResult, error)
```
执行豁免检定

### GetSkillAbility
```go
func (e *Engine) GetSkillAbility(skill model.Skill) model.Ability
```
获取技能对应的属性

### GetPassivePerception
```go
func (e *Engine) GetPassivePerception(ctx context.Context, req GetPassivePerceptionRequest) (*GetPassivePerceptionResult, error)
```
获取被动感知（察觉）

---

## 9. 库存管理

### AddItem
```go
func (e *Engine) AddItem(ctx context.Context, req AddItemRequest) (*InventoryResult, error)
```
添加物品到角色库存

### RemoveItem
```go
func (e *Engine) RemoveItem(ctx context.Context, req RemoveItemRequest) (*InventoryResult, error)
```
从角色库存移除物品

### GetInventory
```go
func (e *Engine) GetInventory(ctx context.Context, req GetInventoryRequest) (*InventoryInfo, error)
```
获取角色库存信息

### EquipItem
```go
func (e *Engine) EquipItem(ctx context.Context, req EquipItemRequest) (*EquipResult, error)
```
装备物品到指定槽位

### UnequipItem
```go
func (e *Engine) UnequipItem(ctx context.Context, req UnequipItemRequest) (*EquipResult, error)
```
卸下装备

### GetEquipment
```go
func (e *Engine) GetEquipment(ctx context.Context, req GetEquipmentRequest) (*EquipmentInfo, error)
```
获取角色当前装备信息

### AttuneItem
```go
func (e *Engine) AttuneItem(ctx context.Context, req AttuneItemRequest) (*AttuneResult, error)
```
调谐或解除调谐魔法物品

### TransferItem
```go
func (e *Engine) TransferItem(ctx context.Context, req TransferItemRequest) (*TransferResult, error)
```
将物品从一个角色转移给另一个角色

### AddCurrency
```go
func (e *Engine) AddCurrency(ctx context.Context, req AddCurrencyRequest) (*InventoryResult, error)
```
添加货币

---

## 10. 专长系统

### SelectFeat
```go
func (e *Engine) SelectFeat(ctx context.Context, req SelectFeatRequest) (*SelectFeatResult, error)
```
为角色选择并获得一个专长

### ListFeats
```go
func (e *Engine) ListFeats(ctx context.Context, req ListFeatsRequest) (*ListFeatsResult, error)
```
列出可选专长

### GetFeatDetails
```go
func (e *Engine) GetFeatDetails(ctx context.Context, req GetFeatDetailsRequest) (*GetFeatDetailsResult, error)
```
获取专长详情

### RemoveFeat
```go
func (e *Engine) RemoveFeat(ctx context.Context, req RemoveFeatRequest) error
```
从角色移除专长

### GetActorFeats
```go
func (e *Engine) GetActorFeats(ctx context.Context, req GetActorRequest) (*ListFeatsResult, error)
```
获取角色的专长列表

---

## 11. 场景管理

### CreateScene
```go
func (e *Engine) CreateScene(ctx context.Context, req CreateSceneRequest) (*CreateSceneResult, error)
```
创建新场景

### GetScene
```go
func (e *Engine) GetScene(ctx context.Context, req GetSceneRequest) (*SceneInfo, error)
```
获取场景信息

### UpdateScene
```go
func (e *Engine) UpdateScene(ctx context.Context, req UpdateSceneRequest) error
```
更新场景信息

### DeleteScene
```go
func (e *Engine) DeleteScene(ctx context.Context, req DeleteSceneRequest) error
```
删除场景

### ListScenes
```go
func (e *Engine) ListScenes(ctx context.Context, req ListScenesRequest) (*ListScenesResult, error)
```
列出所有场景

### SetCurrentScene
```go
func (e *Engine) SetCurrentScene(ctx context.Context, req SetCurrentSceneRequest) error
```
设置当前场景

### GetCurrentScene
```go
func (e *Engine) GetCurrentScene(ctx context.Context, req GetCurrentSceneRequest) (*SceneInfo, error)
```
获取当前场景

### AddSceneConnection
```go
func (e *Engine) AddSceneConnection(ctx context.Context, req AddSceneConnectionRequest) error
```
添加场景连接

### RemoveSceneConnection
```go
func (e *Engine) RemoveSceneConnection(ctx context.Context, req RemoveSceneConnectionRequest) error
```
移除场景连接

### MoveActorToScene
```go
func (e *Engine) MoveActorToScene(ctx context.Context, req MoveActorToSceneRequest) (*MoveActorResult, error)
```
移动角色到另一个场景

### GetSceneActors
```go
func (e *Engine) GetSceneActors(ctx context.Context, req GetSceneActorsRequest) (*GetSceneActorsResult, error)
```
获取场景中的所有角色

### AddItemToScene
```go
func (e *Engine) AddItemToScene(ctx context.Context, req AddItemToSceneRequest) error
```
添加物品到场景

### RemoveItemFromScene
```go
func (e *Engine) RemoveItemFromScene(ctx context.Context, req RemoveItemFromSceneRequest) error
```
从场景移除物品

### GetSceneItems
```go
func (e *Engine) GetSceneItems(ctx context.Context, req GetSceneItemsRequest) (*GetSceneItemsResult, error)
```
获取场景中的物品

---

## 12. 探索系统

### StartTravel
```go
func (e *Engine) StartTravel(ctx context.Context, req StartTravelRequest) (*StartTravelResult, error)
```
开始一段新的旅行

### AdvanceTravel
```go
func (e *Engine) AdvanceTravel(ctx context.Context, req AdvanceTravelRequest) (*AdvanceTravelResult, error)
```
推进旅行进度

### Forage
```go
func (e *Engine) Forage(ctx context.Context, req ForageRequest) (*ForageResultEngine, error)
```
执行觅食行动

### Navigate
```go
func (e *Engine) Navigate(ctx context.Context, req NavigateRequest) (*NavigateResult, error)
```
执行导航检定

---

## 13. 社交互动

### InteractWithNPC
```go
func (e *Engine) InteractWithNPC(ctx context.Context, req InteractWithNPCRequest) (*InteractWithNPCResult, error)
```
执行与 NPC 的社交互动

### GetNPCAttitude
```go
func (e *Engine) GetNPCAttitude(ctx context.Context, req GetNPCAttitudeRequest) (*GetNPCAttitudeResult, error)
```
获取 NPC 当前态度

---

## 14. 任务系统

### CreateQuest
```go
func (e *Engine) CreateQuest(ctx context.Context, req CreateQuestRequest) (*QuestResult, error)
```
创建新任务

### GetQuest
```go
func (e *Engine) GetQuest(ctx context.Context, req GetQuestRequest) (*QuestInfo, error)
```
获取任务信息

### ListQuests
```go
func (e *Engine) ListQuests(ctx context.Context, req ListQuestsRequest) (*ListQuestsResult, error)
```
列出所有任务

### AcceptQuest
```go
func (e *Engine) AcceptQuest(ctx context.Context, req AcceptQuestRequest) (*QuestResult, error)
```
接受任务

### UpdateQuestObjective
```go
func (e *Engine) UpdateQuestObjective(ctx context.Context, req UpdateQuestObjectiveRequest) (*QuestResult, error)
```
更新任务目标进度

### CompleteQuest
```go
func (e *Engine) CompleteQuest(ctx context.Context, req CompleteQuestRequest) (*QuestResult, error)
```
完成任务并发放奖励

### FailQuest
```go
func (e *Engine) FailQuest(ctx context.Context, req FailQuestRequest) (*QuestResult, error)
```
标记任务失败

### DeleteQuest
```go
func (e *Engine) DeleteQuest(ctx context.Context, req DeleteQuestRequest) error
```
删除任务

### GetActorQuests
```go
func (e *Engine) GetActorQuests(ctx context.Context, req GetActorQuestsRequest) (*GetActorQuestsResult, error)
```
获取角色的任务列表

### GetQuestGiverQuests
```go
func (e *Engine) GetQuestGiverQuests(ctx context.Context, req GetQuestGiverQuestsRequest) (*GetQuestGiverQuestsResult, error)
```
获取 NPC 发布的任务列表

---

## 15. 死亡豁免

### PerformDeathSave
```go
func (e *Engine) PerformDeathSave(ctx context.Context, req PerformDeathSaveRequest) (*PerformDeathSaveResult, error)
```
执行死亡豁免检定

### StabilizeCreature
```go
func (e *Engine) StabilizeCreature(ctx context.Context, req StabilizeCreatureRequest) (*StabilizeCreatureResult, error)
```
稳定濒死生物

### GetDeathSaveStatus
```go
func (e *Engine) GetDeathSaveStatus(ctx context.Context, req GetDeathSaveStatusRequest) (*GetDeathSaveStatusResult, error)
```
获取死亡豁免状态

---

## API 分类索引

| 分类 | 文件 | API 数量 |
|------|------|----------|
| 引擎初始化 | engine.go | 4 |
| 游戏会话管理 | game.go | 5 |
| 角色管理 | actor.go | 10 |
| 角色升级与经验 | actor.go | 2 |
| 休息系统 | actor.go | 3 |
| 战斗系统 | combat.go | 12 |
| 法术系统 | spell.go | 10 |
| 检定系统 | check.go | 5 |
| 库存管理 | inventory.go | 10 |
| 专长系统 | feat.go | 5 |
| 场景管理 | scene.go | 16 |
| 探索系统 | exploration.go | 4 |
| 社交互动 | social.go | 2 |
| 任务系统 | quest.go | 10 |
| 死亡豁免 | death_saves.go | 3 |

**总计**: 116 个公共 API 方法