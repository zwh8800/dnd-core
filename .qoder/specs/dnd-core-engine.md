# D&D 5e Core Game Engine 实现计划

## Context

项目需要开发一个纯Go语言的D&D 5e核心游戏引擎，作为LLM扮演DM的基础设施。引擎负责管理游戏状态、执行D&D 5e规则、提供丰富的API供LLM调用，从而约束LLM的自由发挥并确保规则执行的准确性。

**技术决策：**
- API风格：纯Go Library（直接import使用）
- 规则深度：完整实现D&D 5e规则引擎
- 存储方案：可插拔存储后端（接口抽象）
- 并发模型：单例并发安全

**重要开发原则：**
- 在编写具体实现时，如果遇到D&D 5e规则相关的疑问，必须查阅 `docs/rules-md/` 目录中的对应章节内容获取准确信息
- 禁止基于个人理解进行猜测或推断规则细节
- 规则书章节对应关系：
  - 角色创建: `docs/rules-md/01-创建角色.md`
  - 种族: `docs/rules-md/02-种族/`
  - 职业: `docs/rules-md/03-职业/`
  - 装备: `docs/rules-md/05-装备.md`
  - 属性值应用: `docs/rules-md/07-属性值应用.md`
  - 冒险规则: `docs/rules-md/08-冒险规则.md`
  - 战斗规则: `docs/rules-md/09-战斗规则.md`
  - 施法系统: `docs/rules-md/10-施法系统.md`
  - 法术: `docs/rules-md/11-法术/`
  - 状态效果: `docs/rules-md/附录A-状态.md`

---

## 包结构设计

```
github.com/zwh8800/dnd-core/
├── pkg/engine/                    # 公共API - 唯一的外部导入路径
│   ├── engine.go                  # 引擎单例、生命周期
│   ├── phase.go                   # 游戏阶段管理与权限控制
│   ├── game.go                    # 游戏会话管理
│   ├── actor.go                   # 角色操作API
│   ├── combat.go                  # 战斗控制API
│   ├── dice.go                    # 骰子API
│   ├── check.go                   # 检定API
│   ├── spell.go                   # 法术API
│   ├── inventory.go               # 物品/装备API
│   ├── quest.go                   # 任务API
│   ├── scene.go                   # 场景API
│   └── state.go                   # 状态查询API
│
├── internal/                      # 私有实现
│   ├── model/                     # 核心数据模型
│   │   ├── actor.go               # Actor基类、PC、NPC、Companion、Enemy
│   │   ├── ability.go             # 属性值、修正值、技能
│   │   ├── combat.go              # 战斗状态、先攻、回合
│   │   ├── character.go           # 角色卡（职业、种族、等级、HP等）
│   │   ├── equipment.go           # 物品、武器、护甲、库存、槽位
│   │   ├── spell.go               # 法术、法术位、法术列表
│   │   ├── condition.go           # 状态效果（目盲、魅惑等）
│   │   ├── dice.go                # 骰子类型、结果
│   │   ├── scene.go               # 场景、地点、连接
│   │   ├── quest.go               # 任务、目标、奖励
│   │   ├── damage.go              # 伤害类型、抗性/免疫/弱点
│   │   ├── action.go              # 动作（攻击、施法、 Dash等）
│   │   ├── rest.go                # 休息机制（短休/长休）
│   │   └── id.go                  # ID生成（ULID）
│   │
│   ├── rules/                     # D&D 5e规则引擎（纯函数）
│   │   ├── calculator.go          # 修正值计算、熟练加值、AC、HP
│   │   ├── check.go               # 检定解析逻辑
│   │   ├── attack.go              # 攻击解析、暴击逻辑
│   │   ├── damage.go              # 伤害计算（含抗性）
│   │   ├── spell.go               # 法术DC、攻击加值、法术位
│   │   ├── condition.go           # 状态效果应用
│   │   ├── level.go               # 升级计算、XP阈值
│   │   ├── rest.go                # 休息恢复计算
│   │   └── constants.go           # D&D常量（DC表、XP表、熟练加值表）
│   │
│   ├── dice/                      # 骰子引擎
│   │   ├── roller.go              # 骰子表达式解析器和投掷器
│   │   ├── expression.go          # "2d6+3", "4d6kh3", 优势/劣势
│   │   └── result.go              # 投掷结果详情
│   │
│   ├── combat/                    # 战斗状态机
│   │   ├── manager.go             # 战斗生命周期
│   │   ├── initiative.go          # 先攻排序
│   │   ├── turn.go                # 回合/动作管理
│   │   ├── opportunity.go         # 借机攻击追踪
│   │   └── movement.go            # 移动追踪
│   │
│   ├── spell/                     # 法术引擎
│   │   ├── caster.go              # 施法者状态管理
│   │   ├── slots.go               # 法术位追踪
│   │   ├── concentration.go       # 专注管理
│   │   └── resolver.go            # 法术效果解析
│   │
│   ├── inventory/                 # 库存管理
│   │   ├── manager.go             # 库存CRUD
│   │   ├── equipment.go           # 装备槽位管理
│   │   ├── encumbrance.go         # 负重计算
│   │   └── currency.go            # 货币管理
│   │
│   ├── quest/                     # 任务管理
│   │   ├── tracker.go             # 任务生命周期
│   │   └── objective.go           # 目标追踪
│   │
│   ├── scene/                     # 场景管理
│   │   ├── manager.go             # 场景CRUD和转换
│   │   └── topology.go            # 场景图（地点连接）
│   │
│   ├── storage/                   # 存储抽象
│   │   ├── store.go               # 存储接口定义
│   │   ├── memory.go              # 内存存储（默认）
│   │   └── json.go                # JSON文件存储
│   │
│   └── data/                      # 静态D&D 5e数据（embedded）
│       ├── races.go               # 种族定义
│       ├── classes.go             # 职业定义
│       ├── skills.go              # 技能定义
│       ├── conditions.go          # 状态定义
│       └── items.go               # 装备/物品数据库
│
└── testutil/                      # 测试工具
    ├── factory.go                 # 测试工厂
    └── fixture.go                 # 预构建测试夹具
```

---

## 核心数据模型

### 1. Actor基类 (internal/model/actor.go)

```go
type Actor struct {
    ID          ID                   `json:"id"`
    Type        ActorType            `json:"type"`
    Name        string               `json:"name"`
    Size        Size                 `json:"size"`
    Speed       int                  `json:"speed"`
    AbilityScores AbilityScores      `json:"ability_scores"`
    Proficiencies Proficiencies      `json:"proficiencies"`
    HitPoints     HitPoints          `json:"hit_points"`
    TempHitPoints int                `json:"temp_hit_points"`
    ArmorClass    int                `json:"armor_class"`
    Conditions    []ConditionInstance `json:"conditions"`
    Exhaustion    int                `json:"exhaustion"`
    SceneID       ID                 `json:"scene_id"`
}
```

### 2. 玩家角色 (internal/model/character.go)

```go
type PlayerCharacter struct {
    Actor             
    Race              RaceReference    `json:"race"`
    Classes           []ClassLevel     `json:"classes"`
    TotalLevel        int              `json:"total_level"`
    Experience        int              `json:"experience"`
    HitDice           []HitDiceEntry   `json:"hit_dice"`
    Inspiration       bool             `json:"inspiration"`
    DeathSaveSuccesses int             `json:"death_save_successes"`
    DeathSaveFailures  int             `json:"death_save_failures"`
    InventoryID       ID               `json:"inventory_id"`
    Spellcasting      *SpellcasterState `json:"spellcasting,omitempty"`
}
```

### 3. 战斗状态 (internal/model/combat.go)

```go
type CombatState struct {
    ID           ID               `json:"id"`
    SceneID      ID               `json:"scene_id"`
    Status       CombatStatus     `json:"status"`
    Round        int              `json:"round"`
    Initiative   []CombatantEntry `json:"initiative"`
    CurrentIndex int              `json:"current_index"`
    CurrentTurn  *TurnState       `json:"current_turn"`
    Log          []CombatLogEntry `json:"log"`
}
```

### 4. 游戏状态 (internal/model/state.go)

```go
type GameState struct {
    ID           ID                      `json:"id"`
    Name         string                  `json:"name"`
    Phase        Phase                   `json:"phase"`           // 当前游戏阶段
    PCs          map[ID]*PlayerCharacter `json:"pcs"`
    Companions   map[ID]*Companion       `json:"companions"`
    NPCs         map[ID]*NPC             `json:"npcs"`
    Enemies      map[ID]*Enemy           `json:"enemies"`
    Scenes       map[ID]*Scene           `json:"scenes"`
    CurrentScene *ID                     `json:"current_scene,omitempty"`
    Items        map[ID]*Item            `json:"items"`
    Inventories  map[ID]*Inventory       `json:"inventories"`
    Combat       *CombatState            `json:"combat,omitempty"`
    Quests       map[ID]*Quest           `json:"quests"`
}
```

---

## 公共API设计

### 引擎初始化

```go
// pkg/engine/engine.go

// Engine 是D&D 5e游戏引擎的核心控制器，提供对所有游戏系统的统一访问入口。
// Engine是并发安全的，可以在多个goroutine中同时使用。
// 所有对游戏状态的修改都会自动进行阶段权限验证和状态一致性检查。
type Engine struct {
    // ... 内部字段
}

// Config 包含引擎的配置选项
type Config struct {
    // Storage 指定存储后端，用于游戏状态的持久化
    // 如果为nil，将使用默认的内存存储
    Storage storage.Store

    // DiceSeed 指定骰子随机数生成器的种子
    // 如果为0，将使用系统时间作为种子
    // 设置固定种子可用于测试或可重现的游戏
    DiceSeed int64

    // DataPath 指定自定义数据文件的路径
    // 用于覆盖内置的种族、职业、法术等数据
    // 如果为空，将仅使用内置数据
    DataPath string
}

// New 创建并初始化一个新的引擎实例
// 参数:
//   cfg - 引擎配置，可以传入DefaultConfig()使用默认配置
// 返回:
//   *Engine - 初始化完成的引擎实例
//   error - 初始化过程中可能发生的错误（如存储初始化失败）
// 使用场景: 应用程序启动时调用一次
func New(cfg Config) (*Engine, error)

// DefaultConfig 返回使用内存存储的默认配置
// 返回:
//   Config - 可直接用于New()的默认配置
// 使用场景: 快速启动或测试时调用
func DefaultConfig() Config

// Close 释放引擎占用的所有资源，包括存储后端连接
// 返回:
//   error - 关闭过程中可能发生的错误
// 使用场景: 应用程序关闭时调用
func (e *Engine) Close() error
```

### 游戏阶段管理

```go
// pkg/engine/phase.go

// Phase 定义游戏进行的不同阶段，每个阶段限制可用的操作集合
type Phase string

const (
    // PhaseCharacterCreation 角色创建阶段
    // 可用操作: CreatePC, CreateNPC, CreateEnemy, CreateCompanion
    //           GetActor, UpdateActor, RemoveActor
    //           GetSpellList, GetSpellInfo, Roll, LevelUp
    PhaseCharacterCreation Phase = "character_creation"

    // PhaseExploration 探索阶段 - 默认阶段
    // 可用操作: 所有角色操作、场景操作、物品操作、任务操作
    //           骰子、检定、法术施放、ShortRest
    // 不可用操作: StartCombat, ExecuteAttack, ExecuteDamage,
    //           ExecuteAction(战斗动作), NextTurn
    PhaseExploration Phase = "exploration"

    // PhaseCombat 战斗阶段
    // 可用操作: 所有战斗操作、角色状态查询、骰子、检定
    //           GetActor, GetInventory, CastSpell, ExecuteHealing
    // 不可用操作: CreateScene, CreateQuest, StartLongRest,
    //           TransferItem(跨场景)
    PhaseCombat Phase = "combat"

    // PhaseRest 休息阶段
    // 可用操作: ShortRest, StartLongRest, EndLongRest
    //           GetActor, UpdateActor, ExecuteHealing
    //           角色状态查询
    // 不可用操作: 战斗操作、场景移动、任务操作
    PhaseRest Phase = "rest"
)

// SetPhase 切换游戏当前阶段
// 参数:
//   ctx - 上下文，支持取消和超时
//   gameID - 目标游戏会话ID
//   phase - 要切换到的新阶段
//   reason - 切换原因的描述（用于日志和LLM上下文）
// 返回:
//   *PhaseTransitionResult - 包含阶段转换详情
//   error - 转换失败时返回错误（如无效阶段）
// 使用场景: DM决定改变游戏流程时调用
// 注意: 从Combat阶段退出会自动结束战斗
func (e *Engine) SetPhase(ctx context.Context, gameID model.ID, phase Phase, reason string) (*PhaseTransitionResult, error)

// GetPhase 获取游戏当前阶段
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
// 返回:
//   Phase - 当前游戏阶段
//   error - 游戏不存在时返回错误
// 使用场景: 在执行操作前检查当前阶段
func (e *Engine) GetPhase(ctx context.Context, gameID model.ID) (Phase, error)

// GetPhaseInfo 获取指定阶段的详细信息
// 参数:
//   phase - 要查询的阶段
// 返回:
//   *PhaseInfo - 包含阶段名称、描述、可用操作列表
// 使用场景: 向LLM展示当前阶段可执行的操作
func (e *Engine) GetPhaseInfo(phase Phase) *PhaseInfo

// GetAllowedOperations 获取当前阶段允许的所有操作
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
// 返回:
//   []Operation - 允许的操作列表
//   error - 游戏不存在时返回错误
// 使用场景: 构建LLM的可用操作菜单
func (e *Engine) GetAllowedOperations(ctx context.Context, gameID model.ID) ([]Operation, error)

// PhaseTransitionResult 阶段转换的结果
type PhaseTransitionResult struct {
    OldPhase       Phase      `json:"old_phase"`
    NewPhase       Phase      `json:"new_phase"`
    Reason         string     `json:"reason"`
    Timestamp      time.Time  `json:"timestamp"`
    AutoActions    []string   `json:"auto_actions"`
    Message        string     `json:"message"`
}

// PhaseInfo 阶段的详细信息
type PhaseInfo struct {
    Phase          Phase      `json:"phase"`
    Name           string     `json:"name"`
    Description    string     `json:"description"`
    AllowedOps     []Operation `json:"allowed_ops"`
    DeniedOps      []Operation `json:"denied_ops"`
    TypicalActions []string   `json:"typical_actions"`
}

// Operation 定义引擎支持的所有操作类型
type Operation string

const (
    // 角色操作
    OpCreatePC        Operation = "create_pc"
    OpCreateNPC       Operation = "create_npc"
    OpCreateEnemy     Operation = "create_enemy"
    OpCreateCompanion Operation = "create_companion"
    OpGetActor        Operation = "get_actor"
    OpUpdateActor     Operation = "update_actor"
    OpRemoveActor     Operation = "remove_actor"

    // 骰子与检定
    OpRoll            Operation = "roll"
    OpAbilityCheck    Operation = "ability_check"
    OpSkillCheck      Operation = "skill_check"
    OpSavingThrow     Operation = "saving_throw"

    // 战斗操作
    OpStartCombat     Operation = "start_combat"
    OpEndCombat       Operation = "end_combat"
    OpNextTurn        Operation = "next_turn"
    OpExecuteAction   Operation = "execute_action"
    OpExecuteAttack   Operation = "execute_attack"
    OpExecuteDamage   Operation = "execute_damage"
    OpExecuteHealing  Operation = "execute_healing"

    // 法术操作
    OpCastSpell       Operation = "cast_spell"
    OpGetSpellSlots   Operation = "get_spell_slots"
    OpPrepareSpells   Operation = "prepare_spells"

    // 物品操作
    OpAddItem         Operation = "add_item"
    OpEquipItem       Operation = "equip_item"
    OpGetInventory    Operation = "get_inventory"
    OpTransferItem    Operation = "transfer_item"

    // 场景操作
    OpCreateScene     Operation = "create_scene"
    OpMoveActorToScene Operation = "move_actor_to_scene"
    OpSetCurrentScene Operation = "set_current_scene"

    // 任务操作
    OpCreateQuest     Operation = "create_quest"
    OpAcceptQuest     Operation = "accept_quest"
    OpUpdateObjective Operation = "update_objective"
    OpCompleteQuest   Operation = "complete_quest"

    // 休息操作
    OpShortRest       Operation = "short_rest"
    OpStartLongRest   Operation = "start_long_rest"
    OpEndLongRest     Operation = "end_long_rest"

    // 经验与升级
    OpAddExperience   Operation = "add_experience"
    OpLevelUp         Operation = "level_up"

    // 状态查询（所有阶段都允许）
    OpGetStateSummary Operation = "get_state_summary"
    OpGetActorSheet   Operation = "get_actor_sheet"
    OpGetPhase        Operation = "get_phase"
)
```

### 阶段权限矩阵

| 操作 \ 阶段 | Creation | Exploration | Combat | Rest |
|------------|----------|-------------|--------|------|
| CreatePC | ✓ | ✓ | ✗ | ✗ |
| CreateNPC | ✓ | ✓ | ✗ | ✗ |
| CreateEnemy | ✓ | ✓ | ✓* | ✗ |
| GetActor | ✓ | ✓ | ✓ | ✓ |
| UpdateActor | ✓ | ✓ | ✓ | ✓ |
| RemoveActor | ✓ | ✓ | ✗ | ✗ |
| Roll | ✓ | ✓ | ✓ | ✓ |
| AbilityCheck | ✓ | ✓ | ✓ | ✗ |
| SkillCheck | ✓ | ✓ | ✓ | ✗ |
| SavingThrow | ✗ | ✓ | ✓ | ✗ |
| StartCombat | ✗ | ✓ | ✗ | ✗ |
| EndCombat | ✗ | ✗ | ✓ | ✗ |
| NextTurn | ✗ | ✗ | ✓ | ✗ |
| ExecuteAction | ✗ | ✗ | ✓ | ✗ |
| ExecuteAttack | ✗ | ✗ | ✓ | ✗ |
| ExecuteDamage | ✗ | ✗ | ✓ | ✗ |
| ExecuteHealing | ✓ | ✓ | ✓ | ✓ |
| CastSpell | ✗ | ✓ | ✓ | ✗ |
| GetSpellSlots | ✓ | ✓ | ✓ | ✓ |
| PrepareSpells | ✓ | ✓ | ✗ | ✗ |
| AddItem | ✓ | ✓ | ✗ | ✗ |
| EquipItem | ✓ | ✓ | ✗ | ✗ |
| GetInventory | ✓ | ✓ | ✓ | ✓ |
| TransferItem | ✓ | ✓ | ✗ | ✗ |
| CreateScene | ✓ | ✓ | ✗ | ✗ |
| MoveActorToScene | ✗ | ✓ | ✗ | ✗ |
| SetCurrentScene | ✗ | ✓ | ✗ | ✗ |
| CreateQuest | ✗ | ✓ | ✗ | ✗ |
| AcceptQuest | ✗ | ✓ | ✗ | ✗ |
| UpdateObjective | ✗ | ✓ | ✗ | ✗ |
| CompleteQuest | ✗ | ✓ | ✗ | ✗ |
| ShortRest | ✗ | ✓ | ✗ | ✓ |
| StartLongRest | ✗ | ✓ | ✗ | ✓ |
| EndLongRest | ✗ | ✗ | ✗ | ✓ |
| AddExperience | ✗ | ✓ | ✗ | ✗ |
| LevelUp | ✓ | ✗ | ✗ | ✗ |

注：✓* = 仅允许在Combat阶段开始时自动创建（遭遇战）

### 阶段自动转换规则

```go
// 以下操作会触发自动阶段转换:
// StartCombat → 自动从 Exploration 转换到 Combat
// EndCombat → 自动从 Combat 转换到 Exploration
// StartLongRest → 从 Exploration 手动转换到 Rest
// EndLongRest → 从 Rest 自动转换回 Exploration
// Character Creation完成 → 从 Creation 转换到 Exploration
```

---

### 游戏会话管理

```go
// pkg/engine/game.go

// NewGame 创建一个新的游戏会话
// 参数:
//   ctx - 上下文，支持取消和超时
//   name - 游戏名称，用于标识和显示
//   description - 游戏描述，记录游戏背景信息
// 返回:
//   *model.GameState - 新创建的游戏状态，初始阶段为PhaseCharacterCreation
//   error - 创建失败时返回错误
// 使用场景: 开始新游戏时调用
// 注意: 新游戏默认处于角色创建阶段
func (e *Engine) NewGame(ctx context.Context, name, description string) (*model.GameState, error)

// LoadGame 从存储加载一个已存在的游戏会话
// 参数:
//   ctx - 上下文
//   gameID - 要加载的游戏会话ID
// 返回:
//   *model.GameState - 加载的游戏状态（深拷贝）
//   error - 游戏不存在或加载失败时返回错误
// 使用场景: 继续之前保存的游戏
func (e *Engine) LoadGame(ctx context.Context, gameID model.ID) (*model.GameState, error)

// SaveGame 将当前游戏状态持久化到存储后端
// 参数:
//   ctx - 上下文
//   gameID - 要保存的游戏会话ID
// 返回:
//   error - 保存失败时返回错误
// 使用场景: 定期保存或用户主动保存
// 注意: 保存前会自动验证状态一致性
func (e *Engine) SaveGame(ctx context.Context, gameID model.ID) error

// DeleteGame 从存储中删除一个游戏会话
// 参数:
//   ctx - 上下文
//   gameID - 要删除的游戏会话ID
// 返回:
//   error - 删除失败时返回错误
// 使用场景: 清理不再需要的游戏
func (e *Engine) DeleteGame(ctx context.Context, gameID model.ID) error

// ListGames 列出所有可用的游戏会话
// 参数:
//   ctx - 上下文
// 返回:
//   []GameSummary - 游戏摘要列表，按更新时间降序
//   error - 列出失败时返回错误
// 使用场景: 显示游戏选择列表
func (e *Engine) ListGames(ctx context.Context) ([]GameSummary, error)

type GameSummary struct {
    ID             model.ID  `json:"id"`
    Name           string    `json:"name"`
    Description    string    `json:"description"`
    Phase          Phase     `json:"phase"`
    UpdatedAt      time.Time `json:"updated_at"`
    PCCount        int       `json:"pc_count"`
    CurrentLevel   int       `json:"current_level"`
}
```

### 角色管理

```go
// pkg/engine/actor.go

// CreatePC 创建一个新的玩家角色
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   pc - 玩家角色数据，包含种族、职业、属性等
// 返回:
//   *model.PlayerCharacter - 创建完成的角色（含自动生成的ID和计算的派生值）
//   error - 创建失败时返回错误（如阶段不允许、数据无效）
// 使用场景: 玩家创建新角色时调用
// 权限: 需要 Creation/Exploration 阶段
func (e *Engine) CreatePC(ctx context.Context, gameID model.ID, pc *model.PlayerCharacter) (*model.PlayerCharacter, error)

// CreateNPC 创建一个新的非玩家角色
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   npc - NPC数据
// 返回:
//   *model.NPC - 创建的NPC
//   error - 创建失败时返回错误
// 使用场景: DM添加NPC到游戏世界
func (e *Engine) CreateNPC(ctx context.Context, gameID model.ID, npc *model.NPC) (*model.NPC, error)

// CreateEnemy 创建一个新的敌人/怪物
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   enemy - 敌人数据
// 返回:
//   *model.Enemy - 创建的敌人
//   error - 创建失败时返回错误
// 使用场景: 战斗前放置敌人
func (e *Engine) CreateEnemy(ctx context.Context, gameID model.ID, enemy *model.Enemy) (*model.Enemy, error)

// CreateCompanion 创建一个同伴角色（由AI控制的盟友）
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   companion - 同伴数据
// 返回:
//   *model.Companion - 创建的同伴
//   error - 创建失败时返回错误
// 使用场景: 玩家获得同伴时调用
func (e *Engine) CreateCompanion(ctx context.Context, gameID model.ID, companion *model.Companion) (*model.Companion, error)

// GetActor 获取任意类型的角色信息
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   actorID - 角色ID
// 返回:
//   ActorSnapshot - 角色快照（深拷贝），包含基本信息
//   error - 角色不存在时返回错误
// 使用场景: 查询角色状态
// 权限: 所有阶段都允许
func (e *Engine) GetActor(ctx context.Context, gameID model.ID, actorID model.ID) (ActorSnapshot, error)

// GetPC 获取玩家角色的完整数据
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   pcID - 玩家角色ID
// 返回:
//   *model.PlayerCharacter - 完整角色数据
//   error - 角色不存在或不是PC时返回错误
func (e *Engine) GetPC(ctx context.Context, gameID model.ID, pcID model.ID) (*model.PlayerCharacter, error)

// UpdateActor 更新角色的部分状态
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   actorID - 要更新的角色ID
//   update - 更新内容，只包含需要修改的字段
// 返回:
//   error - 更新失败时返回错误
// 使用场景: 修改角色属性、状态、位置等
// 注意: update中为nil的字段不会被修改
func (e *Engine) UpdateActor(ctx context.Context, gameID model.ID, actorID model.ID, update ActorUpdate) error

// RemoveActor 从游戏中移除一个角色
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   actorID - 要移除的角色ID
// 返回:
//   error - 移除失败时返回错误
// 使用场景: 角色死亡或离开队伍
// 权限: 不能在Combat阶段移除参与战斗的角色
func (e *Engine) RemoveActor(ctx context.Context, gameID model.ID, actorID model.ID) error

// ListActors 列出游戏中的角色，可按条件过滤
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   filter - 过滤条件，为nil时返回所有角色
// 返回:
//   []ActorSnapshot - 符合条件的角色列表
//   error - 查询失败时返回错误
func (e *Engine) ListActors(ctx context.Context, gameID model.ID, filter *ActorFilter) ([]ActorSnapshot, error)

// AddExperience 为玩家角色添加经验值
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   pcID - 玩家角色ID
//   xp - 要添加的经验值数量
// 返回:
//   *LevelUpResult - 如果升级则包含升级详情，否则为nil
//   error - 添加失败时返回错误
// 使用场景: 完成任务或击败敌人后奖励XP
// 权限: 需要 Exploration 阶段
// 注意: 如果经验值达到升级阈值，会自动触发升级
func (e *Engine) AddExperience(ctx context.Context, gameID model.ID, pcID model.ID, xp int) (*LevelUpResult, error)

// LevelUp 手动触发玩家角色升级
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   pcID - 玩家角色ID
//   classChoice - 选择要升级的职业（多职业时）
// 返回:
//   *LevelUpResult - 升级结果，包含HP增长、新特性等
//   error - 升级失败时返回错误（如未达到XP要求）
// 使用场景: DM手动调整角色等级
// 权限: 需要 Creation 阶段
func (e *Engine) LevelUp(ctx context.Context, gameID model.ID, pcID model.ID, classChoice string) (*LevelUpResult, error)

// ShortRest 为指定角色执行短休
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   actorIDs - 参与短休的角色ID列表
// 返回:
//   *RestResult - 休息结果，包含每个角色的恢复情况
//   error - 休息失败时返回错误
// 使用场景: 角色花1小时休息恢复
// 权限: 需要 Exploration/Rest 阶段
func (e *Engine) ShortRest(ctx context.Context, gameID model.ID, actorIDs []model.ID) (*RestResult, error)

// StartLongRest 开始长休过程
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   actorIDs - 参与长休的角色ID列表
// 返回:
//   *RestResult - 长休开始结果
//   error - 开始失败时返回错误
// 使用场景: 角色开始8小时休息
// 权限: 需要 Exploration/Rest 阶段
// 注意: 调用后游戏阶段自动切换到Rest
// 规则参考: docs/rules-md/08-冒险规则.md (长休规则)
func (e *Engine) StartLongRest(ctx context.Context, gameID model.ID, actorIDs []model.ID) (*RestResult, error)

// EndLongRest 结束长休并应用恢复效果
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
// 返回:
//   *RestResult - 长休结束结果，包含所有角色的恢复情况
//   error - 结束失败时返回错误
// 使用场景: 长休完成后调用
// 权限: 仅 Rest 阶段允许
// 注意: 调用后游戏阶段自动切换回Exploration
// 长休效果:
//   - 恢复所有HP
//   - 恢复所有已消耗的法术位
//   - 恢复所有Hit Dice（最多恢复一半，至少一个）
//   - 减少一级力竭（如果适用）
// 规则参考: docs/rules-md/08-冒险规则.md (长休规则)
func (e *Engine) EndLongRest(ctx context.Context, gameID model.ID) (*RestResult, error)

type ActorSnapshot struct {
    ID            model.ID        `json:"id"`
    Type          model.ActorType `json:"type"`
    Name          string          `json:"name"`
    HitPoints     model.HitPoints `json:"hit_points"`
    ArmorClass    int             `json:"armor_class"`
    Conditions    []string        `json:"conditions"`
    SceneID       model.ID        `json:"scene_id"`
    Summary       string          `json:"summary"`
}

type ActorFilter struct {
    Types   []model.ActorType
    SceneID *model.ID
    Alive   *bool
}

type ActorUpdate struct {
    AbilityScores   *model.AbilityScores
    HitPoints       *HitPointUpdate
    Conditions      *ConditionUpdate
    Position        *model.Point
    SceneID         *model.ID
    Custom          map[string]any
}

type HitPointUpdate struct {
    Current       *int
    TempHitPoints *int
}

type ConditionUpdate struct {
    Add    []model.ConditionInstance
    Remove []model.ConditionType
}

type LevelUpResult struct {
    OldLevel              int      `json:"old_level"`
    NewLevel              int      `json:"new_level"`
    HPGain                int      `json:"hp_gain"`
    NewFeatures           []string `json:"new_features"`
    SpellSlotsUpdated     bool     `json:"spell_slots_updated"`
    ProficiencyIncreased  bool     `json:"proficiency_increased"`
    Message               string   `json:"message"`
}

type RestResult struct {
    ActorResults []ActorRestResult `json:"actor_results"`
    Message      string            `json:"message"`
}

type ActorRestResult struct {
    ActorID           model.ID              `json:"actor_id"`
    HPRecovered       int                   `json:"hp_recovered"`
    HitDiceUsed       int                   `json:"hit_dice_used"`
    SpellSlotsRestored bool                 `json:"spell_slots_restored"`
    ConditionsRemoved []model.ConditionType `json:"conditions_removed"`
    ExhaustionReduced bool                  `json:"exhaustion_reduced"`
    AbilitiesRestored bool                  `json:"abilities_restored"`
}
```

### 骰子系统

```go
// pkg/engine/dice.go

// Roll 解析并投掷骰子表达式
// 参数:
//   ctx - 上下文
//   expression - 骰子表达式，如 "2d6+3", "1d20", "4d6kh3"
// 返回:
//   *model.DiceResult - 投掷结果，包含每个骰子的点数和总计
//   error - 表达式无效时返回错误
// 使用场景: 通用的骰子投掷
// 权限: 所有阶段都允许
func (e *Engine) Roll(ctx context.Context, expression string) (*model.DiceResult, error)

// RollWithModifier 投掷骰子并添加固定修正值
// 参数:
//   ctx - 上下文
//   expression - 骰子表达式
//   modifier - 要添加的修正值（可正可负）
// 返回:
//   *model.DiceResult - 投掷结果
//   error - 表达式无效时返回错误
func (e *Engine) RollWithModifier(ctx context.Context, expression string, modifier int) (*model.DiceResult, error)

// RollAdvantage 投掷带优势的d20（掷2取高）
// 参数:
//   ctx - 上下文
//   modifier - 添加到结果的修正值
// 返回:
//   *model.DiceResult - 投掷结果
// 使用场景: 角色有优势时的攻击掷骰或检定
func (e *Engine) RollAdvantage(ctx context.Context, modifier int) (*model.DiceResult, error)

// RollDisadvantage 投掷带劣势的d20（掷2取低）
// 参数:
//   ctx - 上下文
//   modifier - 添加到结果的修正值
// 返回:
//   *model.DiceResult - 投掷结果
func (e *Engine) RollDisadvantage(ctx context.Context, modifier int) (*model.DiceResult, error)

// RollSecret 投掷隐藏骰子（结果对玩家不可见）
// 参数:
//   ctx - 上下文
//   expression - 骰子表达式
// 返回:
//   *model.DiceResult - 投掷结果
//   error - 表达式无效时返回错误
// 使用场景: DM的秘密检定
func (e *Engine) RollSecret(ctx context.Context, expression string) (*model.DiceResult, error)
```

### 检定系统

```go
// pkg/engine/check.go

// AbilityCheck 执行属性检定
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   actorID - 执行检定的角色ID
//   ability - 要检定的属性（STR/DEX/CON/INT/WIS/CHA）
//   dc - 难度等级
//   modifier - 额外修正值（优势/劣势）
// 返回:
//   *CheckResult - 检定结果
//   error - 角色不存在时返回错误
// 权限: 需要 Creation/Exploration/Combat 阶段
func (e *Engine) AbilityCheck(ctx context.Context, gameID model.ID, actorID model.ID, ability model.Ability, dc int, modifier RollModifier) (*CheckResult, error)

// SkillCheck 执行技能检定
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   actorID - 角色ID
//   skill - 要检定的技能
//   dc - 难度等级
//   modifier - 额外修正值
// 返回:
//   *CheckResult - 检定结果
//   error - 角色不存在时返回错误
func (e *Engine) SkillCheck(ctx context.Context, gameID model.ID, actorID model.ID, skill model.Skill, dc int, modifier RollModifier) (*CheckResult, error)

// SavingThrow 执行豁免检定
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   actorID - 角色ID
//   ability - 要豁免的属性
//   dc - 豁免DC
//   modifier - 额外修正值
// 返回:
//   *CheckResult - 检定结果
//   error - 角色不存在时返回错误
// 权限: 需要 Exploration/Combat 阶段
func (e *Engine) SavingThrow(ctx context.Context, gameID model.ID, actorID model.ID, ability model.Ability, dc int, modifier RollModifier) (*CheckResult, error)

// PassiveCheck 计算被动检定分数
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   actorID - 角色ID
//   skill - 要计算的技能
// 返回:
//   int - 被动检定分数（10 + 所有修正值）
// 使用场景: 不主动掷骰时的默认察觉
func (e *Engine) PassiveCheck(ctx context.Context, gameID model.ID, actorID model.ID, skill model.Skill) int

// Contest 执行对抗检定
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   actorA, actorB - 对抗双方角色ID
//   skillA, skillB - 各自使用的技能
// 返回:
//   *ContestResult - 对抗结果
//   error - 角色不存在时返回错误
func (e *Engine) Contest(ctx context.Context, gameID model.ID, actorA model.ID, skillA model.Skill, actorB model.ID, skillB model.Skill) (*ContestResult, error)

// GroupCheck 执行群体检定
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   actorIDs - 参与检定的角色ID列表
//   skill - 使用的技能
//   dc - 难度等级
// 返回:
//   *GroupCheckResult - 群体检定结果
func (e *Engine) GroupCheck(ctx context.Context, gameID model.ID, actorIDs []model.ID, skill model.Skill, dc int) (*GroupCheckResult, error)

type CheckResult struct {
    Roll       *model.DiceResult `json:"roll"`
    Modifier   int               `json:"modifier"`
    Total      int               `json:"total"`
    DC         int               `json:"dc"`
    Success    bool              `json:"success"`
    ActorID    model.ID          `json:"actor_id"`
    CheckType  string            `json:"check_type"`
    Details    string            `json:"details"`
}

type ContestResult struct {
    ActorAResult *CheckResult `json:"actor_a_result"`
    ActorBResult *CheckResult `json:"actor_b_result"`
    Winner       model.ID     `json:"winner"`
}

type GroupCheckResult struct {
    Individual []*CheckResult `json:"individual"`
    Successes  int            `json:"successes"`
    Failures   int            `json:"failures"`
    Passed     bool           `json:"passed"`
}
```

### 战斗系统

```go
// pkg/engine/combat.go

// StartCombat 开始一场战斗遭遇
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   sceneID - 战斗发生的场景ID
//   participantIDs - 参与战斗的所有角色ID
// 返回:
//   *model.CombatState - 创建的战斗状态，包含先攻顺序
//   error - 开始失败时返回错误
// 使用场景: 遭遇敌人进入战斗
// 权限: 仅 Exploration 阶段允许
// 注意: 调用后自动将游戏阶段切换到Combat
func (e *Engine) StartCombat(ctx context.Context, gameID model.ID, sceneID model.ID, participantIDs []model.ID) (*model.CombatState, error)

// StartCombatWithSurprise 开始带突袭判定的战斗
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   sceneID - 战斗场景ID
//   stealthySide - 隐秘一方的角色ID
//   observers - 观察一方的角色ID
// 返回:
//   *model.CombatState - 战斗状态，包含突袭标记
//   error - 开始失败时返回错误
func (e *Engine) StartCombatWithSurprise(ctx context.Context, gameID model.ID, sceneID model.ID, stealthySide []model.ID, observers []model.ID) (*model.CombatState, error)

// EndCombat 结束当前战斗
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
// 返回:
//   error - 结束失败时返回错误
// 注意: 调用后自动将游戏阶段切换回Exploration
func (e *Engine) EndCombat(ctx context.Context, gameID model.ID) error

// GetCurrentCombat 获取当前活跃的战斗状态
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
// 返回:
//   *model.CombatState - 当前战斗状态
//   error - 没有活跃战斗时返回ErrCombatNotActive
func (e *Engine) GetCurrentCombat(ctx context.Context, gameID model.ID) (*model.CombatState, error)

// NextTurn 推进到下一个角色的回合
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
// 返回:
//   *model.CombatState - 更新后的战斗状态
//   error - 推进失败时返回错误
// 权限: 仅 Combat 阶段允许
func (e *Engine) NextTurn(ctx context.Context, gameID model.ID) (*model.CombatState, error)

// GetCurrentTurn 获取当前回合的信息
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
// 返回:
//   *TurnInfo - 当前回合详情
func (e *Engine) GetCurrentTurn(ctx context.Context, gameID model.ID) (*TurnInfo, error)

// ExecuteAction 在当前回合执行一个动作
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   actorID - 执行动作的角色ID（必须是当前回合角色）
//   action - 要执行的动作
// 返回:
//   *ActionResult - 动作执行结果
//   error - 执行失败时返回错误
// 权限: 仅 Combat 阶段允许
func (e *Engine) ExecuteAction(ctx context.Context, gameID model.ID, actorID model.ID, action ActionInput) (*ActionResult, error)

// ExecuteAttack 执行攻击动作
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   attackerID - 攻击者ID
//   targetID - 目标ID
//   attack - 攻击输入
// 返回:
//   *AttackResult - 攻击结果
func (e *Engine) ExecuteAttack(ctx context.Context, gameID model.ID, attackerID model.ID, targetID model.ID, attack AttackInput) (*AttackResult, error)

// ExecuteDamage 直接对角色施加伤害
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   targetID - 承受伤害的角色ID
//   damage - 伤害输入
// 返回:
//   *DamageResult - 伤害结果
func (e *Engine) ExecuteDamage(ctx context.Context, gameID model.ID, targetID model.ID, damage DamageInput) (*DamageResult, error)

// ExecuteHealing 对角色进行治疗
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   targetID - 治疗的目标角色ID
//   amount - 治疗的HP数量
// 返回:
//   *HealResult - 治疗结果
// 权限: 所有阶段都允许
func (e *Engine) ExecuteHealing(ctx context.Context, gameID model.ID, targetID model.ID, amount int) (*HealResult, error)

// MoveActor 在场景中移动角色
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   actorID - 移动的角色ID
//   to - 目标位置
// 返回:
//   *MoveResult - 移动结果
func (e *Engine) MoveActor(ctx context.Context, gameID model.ID, actorID model.ID, to model.Point) (*MoveResult, error)

// TriggerReaction 触发一个待决反应
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   reactionID - 要触发的反应ID
// 返回:
//   *ActionResult - 反应执行结果
func (e *Engine) TriggerReaction(ctx context.Context, gameID model.ID, reactionID model.ID) (*ActionResult, error)

type ActionInput struct {
    Type    model.ActionType `json:"type"`
    Details map[string]any   `json:"details,omitempty"`
}

type ActionResult struct {
    Success bool              `json:"success"`
    Message string            `json:"message"`
    Roll    *model.DiceResult `json:"roll,omitempty"`
    Effects []EffectDetail    `json:"effects"`
}

type AttackInput struct {
    WeaponID    *model.ID        `json:"weapon_id,omitempty"`
    SpellID     *string          `json:"spell_id,omitempty"`
    IsUnarmed   bool             `json:"is_unarmed"`
    IsOffHand   bool             `json:"is_off_hand"`
    Advantage   model.RollModifier `json:"advantage"`
    ExtraDamage []DamageInput    `json:"extra_damage,omitempty"`
}

type AttackResult struct {
    Roll        *model.DiceResult `json:"roll"`
    AttackTotal int               `json:"attack_total"`
    TargetAC    int               `json:"target_ac"`
    Hit         bool              `json:"hit"`
    IsCritical  bool              `json:"is_critical"`
    IsFumble    bool              `json:"is_fumble"`
    Damage      *DamageResult     `json:"damage,omitempty"`
    Message     string            `json:"message"`
}

type DamageInput struct {
    Amount int              `json:"amount"`
    Type   model.DamageType `json:"type"`
    Dice   string           `json:"dice,omitempty"`
    Source model.ID         `json:"source"`
}

type DamageResult struct {
    RawDamage       int                `json:"raw_damage"`
    Resistances     []model.DamageType `json:"resistances_applied"`
    Vulnerabilities []model.DamageType `json:"vulnerabilities_applied"`
    FinalDamage     int                `json:"final_damage"`
    TargetHPBefore  int                `json:"target_hp_before"`
    TargetHPAfter   int                `json:"target_hp_after"`
    IsDead          bool               `json:"is_dead"`
    IsStabilized    bool               `json:"is_stabilized"`
    DeathSaves      *DeathSaveUpdate   `json:"death_saves,omitempty"`
    Message         string             `json:"message"`
}

type DeathSaveUpdate struct {
    Successes int  `json:"successes"`
    Failures  int  `json:"failures"`
    IsStable  bool `json:"is_stable"`
}

type HealResult struct {
    Amount    int    `json:"amount"`
    HPBefore  int    `json:"hp_before"`
    HPAfter   int    `json:"hp_after"`
    WasStable bool   `json:"was_stable"`
    Message   string `json:"message"`
}

type MoveResult struct {
    Success       bool   `json:"success"`
    DistanceMoved int    `json:"distance_moved"`
    RemainingMove int    `json:"remaining_move"`
    Message       string `json:"message"`
}

type TurnInfo struct {
    ActorID            model.ID `json:"actor_id"`
    ActorName          string   `json:"actor_name"`
    Round              int      `json:"round"`
    InitiativePos      int      `json:"initiative_position"`
    MovementLeft       int      `json:"movement_left"`
    ActionAvailable    bool     `json:"action_available"`
    BonusActionAvailable bool   `json:"bonus_action_available"`
    ReactionAvailable  bool     `json:"reaction_available"`
}
```

### 法术系统

```go
// pkg/engine/spell.go

// CastSpell 执行法术施放
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   casterID - 施法者ID
//   input - 施法输入
// 返回:
//   *SpellCastResult - 施法结果
//   error - 施法失败时返回错误
// 权限: Exploration/Combat 阶段允许
func (e *Engine) CastSpell(ctx context.Context, gameID model.ID, casterID model.ID, input SpellCastInput) (*SpellCastResult, error)

// GetSpellInfo 获取法术的详细信息
// 参数:
//   ctx - 上下文
//   spellID - 法术ID
// 返回:
//   *model.Spell - 法术数据
func (e *Engine) GetSpellInfo(ctx context.Context, spellID string) (*model.Spell, error)

// GetSpellList 获取指定职业可用的法术列表
// 参数:
//   ctx - 上下文
//   className - 职业名称
//   level - 法术等级（0-9）
// 返回:
//   []model.Spell - 可用法术列表
func (e *Engine) GetSpellList(ctx context.Context, className string, level int) ([]model.Spell, error)

// GetSpellSlots 获取施法者的当前法术位状态
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   actorID - 施法者ID
// 返回:
//   *model.SpellSlotTracker - 法术位追踪器
func (e *Engine) GetSpellSlots(ctx context.Context, gameID model.ID, actorID model.ID) (*model.SpellSlotTracker, error)

// PrepareSpells 为准备型施法者准备法术
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   casterID - 施法者ID
//   spellIDs - 要准备的法术ID列表
// 返回:
//   error - 准备失败时返回错误
// 权限: Creation/Exploration 阶段允许
func (e *Engine) PrepareSpells(ctx context.Context, gameID model.ID, casterID model.ID, spellIDs []string) error

// LearnSpell 为已知型施法者学习新法术
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   casterID - 施法者ID
//   spellID - 要学习的法术ID
// 返回:
//   error - 学习失败时返回错误
func (e *Engine) LearnSpell(ctx context.Context, gameID model.ID, casterID model.ID, spellID string) error

// ConcentrationCheck 强制施法者进行专注检定
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   casterID - 施法者ID
//   damage - 受到的伤害量
// 返回:
//   *CheckResult - 专注检定结果
// 注意: DC = max(10, damage/2)
func (e *Engine) ConcentrationCheck(ctx context.Context, gameID model.ID, casterID model.ID, damage int) (*CheckResult, error)

// EndConcentration 结束施者的专注
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   casterID - 施法者ID
// 返回:
//   error - 结束失败时返回错误
func (e *Engine) EndConcentration(ctx context.Context, gameID model.ID, casterID model.ID) error

type SpellCastInput struct {
    SpellID     string     `json:"spell_id"`
    SlotLevel   int        `json:"slot_level"`
    Targets     []model.ID `json:"targets"`
    Point       *model.Point `json:"point,omitempty"`
    Description string     `json:"description"`
    IsRitual    bool       `json:"is_ritual"`
    IsReaction  bool       `json:"is_reaction"`
}

type SpellCastResult struct {
    Success  bool             `json:"success"`
    Spell    *model.Spell     `json:"spell"`
    SlotUsed int              `json:"slot_used"`
    Saves    []*CheckResult   `json:"saves"`
    Damage   []*DamageResult  `json:"damage"`
    Effects  []string         `json:"effects"`
    Message  string           `json:"message"`
}
```

### 物品系统

```go
// pkg/engine/inventory.go

// AddItem 向角色的库存添加物品
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   actorID - 目标角色ID
//   item - 要添加的物品
//   quantity - 数量
// 返回:
//   error - 添加失败时返回错误
// 权限: Creation/Exploration 阶段允许
func (e *Engine) AddItem(ctx context.Context, gameID model.ID, actorID model.ID, item *model.Item, quantity int) error

// RemoveItem 从角色的库存移除物品
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   actorID - 角色ID
//   itemID - 要移除的物品ID
//   quantity - 移除数量
// 返回:
//   error - 移除失败时返回错误
func (e *Engine) RemoveItem(ctx context.Context, gameID model.ID, actorID model.ID, itemID model.ID, quantity int) error

// GetInventory 获取角色的完整库存
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   actorID - 角色ID
// 返回:
//   *model.Inventory - 库存数据
func (e *Engine) GetInventory(ctx context.Context, gameID model.ID, actorID model.ID) (*model.Inventory, error)

// EquipItem 装备物品到角色的装备槽
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   actorID - 角色ID
//   itemID - 要装备的物品ID
// 返回:
//   error - 装备失败时返回错误
// 权限: Creation/Exploration 阶段允许
func (e *Engine) EquipItem(ctx context.Context, gameID model.ID, actorID model.ID, itemID model.ID) error

// UnequipItem 卸除已装备的物品
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   actorID - 角色ID
//   itemID - 要卸除的物品ID
// 返回:
//   error - 卸除失败时返回错误
func (e *Engine) UnequipItem(ctx context.Context, gameID model.ID, actorID model.ID, itemID model.ID) error

// GetEquippedStats 获取装备提供的属性加成
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   actorID - 角色ID
// 返回:
//   *EquippedStats - 装备属性统计
func (e *Engine) GetEquippedStats(ctx context.Context, gameID model.ID, actorID model.ID) *EquippedStats

// TransferItem 在两个角色间转移物品
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   fromID - 给予方角色ID
//   toID - 接收方角色ID
//   itemID - 物品ID
//   quantity - 数量
// 返回:
//   error - 转移失败时返回错误
// 权限: Creation/Exploration 阶段允许
func (e *Engine) TransferItem(ctx context.Context, gameID model.ID, fromID model.ID, toID model.ID, itemID model.ID, quantity int) error

// AttuneItem 将魔法物品与角色调音
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   actorID - 角色ID
//   itemID - 要调音的物品ID
// 返回:
//   error - 调音失败时返回错误
func (e *Engine) AttuneItem(ctx context.Context, gameID model.ID, actorID model.ID, itemID model.ID) error

// UnattuneItem 解除魔法物品的调音
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   actorID - 角色ID
//   itemID - 要解除调音的物品ID
// 返回:
//   error - 解除失败时返回错误
func (e *Engine) UnattuneItem(ctx context.Context, gameID model.ID, actorID model.ID, itemID model.ID) error

type EquippedStats struct {
    TotalAC             int                 `json:"total_ac"`
    DexModifierApplied int                 `json:"dex_modifier_applied"`
    SpeedBonus          int                 `json:"speed_bonus"`
    AbilityBonuses      map[model.Ability]int `json:"ability_bonuses"`
    Resistances         []model.DamageType  `json:"resistances"`
    StealthDisadvantage bool                `json:"stealth_disadvantage"`
}
```

### 任务系统

```go
// pkg/engine/quest.go

// CreateQuest 创建一个新的任务
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   quest - 任务数据
// 返回:
//   *model.Quest - 创建的任务
//   error - 创建失败时返回错误
// 权限: Exploration 阶段允许
func (e *Engine) CreateQuest(ctx context.Context, gameID model.ID, quest *model.Quest) (*model.Quest, error)

// GetQuest 获取任务的详细信息
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   questID - 任务ID
// 返回:
//   *model.Quest - 任务数据
func (e *Engine) GetQuest(ctx context.Context, gameID model.ID, questID model.ID) (*model.Quest, error)

// ListQuests 列出任务，可按条件过滤
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   filter - 过滤条件
// 返回:
//   []model.Quest - 任务列表
func (e *Engine) ListQuests(ctx context.Context, gameID model.ID, filter *QuestFilter) ([]model.Quest, error)

// AcceptQuest 让角色接受一个任务
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   questID - 任务ID
//   actorID - 接受任务的角色ID
// 返回:
//   error - 接受失败时返回错误
func (e *Engine) AcceptQuest(ctx context.Context, gameID model.ID, questID model.ID, actorID model.ID) error

// UpdateObjective 更新任务目标的进度
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   questID - 任务ID
//   objectiveID - 目标ID
//   progress - 新进度
//   status - 目标状态
// 返回:
//   error - 更新失败时返回错误
func (e *Engine) UpdateObjective(ctx context.Context, gameID model.ID, questID model.ID, objectiveID string, progress int, status model.ObjectiveStatus) error

// CompleteQuest 完成任务并分发奖励
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   questID - 任务ID
//   actorIDs - 完成任务的角色ID列表
// 返回:
//   *QuestCompletionResult - 完成结果
func (e *Engine) CompleteQuest(ctx context.Context, gameID model.ID, questID model.ID, actorIDs []model.ID) (*QuestCompletionResult, error)

// FailQuest 标记任务为失败
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   questID - 任务ID
//   reason - 失败原因
// 返回:
//   error - 标记失败时返回错误
func (e *Engine) FailQuest(ctx context.Context, gameID model.ID, questID model.ID, reason string) error

type QuestFilter struct {
    Status  *model.QuestStatus
    ActorID *model.ID
}

type QuestCompletionResult struct {
    Quest       *model.Quest     `json:"quest"`
    Rewards     model.QuestRewards `json:"rewards"`
    DistributedTo []model.ID    `json:"distributed_to"`
    Message     string          `json:"message"`
}
```

### 场景系统

```go
// pkg/engine/scene.go

// CreateScene 创建一个新的场景/地点
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   scene - 场景数据
// 返回:
//   *model.Scene - 创建的场景
//   error - 创建失败时返回错误
// 权限: Creation/Exploration 阶段允许
func (e *Engine) CreateScene(ctx context.Context, gameID model.ID, scene *model.Scene) (*model.Scene, error)

// GetScene 获取场景的详细信息
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   sceneID - 场景ID
// 返回:
//   *model.Scene - 场景数据
func (e *Engine) GetScene(ctx context.Context, gameID model.ID, sceneID model.ID) (*model.Scene, error)

// ListScenes 列出所有场景
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
// 返回:
//   []model.Scene - 场景列表
func (e *Engine) ListScenes(ctx context.Context, gameID model.ID) ([]model.Scene, error)

// MoveActorToScene 将角色移动到不同的场景
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   actorID - 要移动的角色ID
//   sceneID - 目标场景ID
// 返回:
//   error - 移动失败时返回错误
// 权限: Exploration 阶段允许
func (e *Engine) MoveActorToScene(ctx context.Context, gameID model.ID, actorID model.ID, sceneID model.ID) error

// ConnectScenes 创建两个场景之间的连接
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   fromID, toID - 连接的两个场景ID
//   description - 连接描述（如"北边的门"）
// 返回:
//   error - 创建连接失败时返回错误
func (e *Engine) ConnectScenes(ctx context.Context, gameID model.ID, fromID model.ID, toID model.ID, description string) error

// SetCurrentScene 设置队伍当前所在的场景
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   sceneID - 要设置的场景ID
// 返回:
//   error - 设置失败时返回错误
func (e *Engine) SetCurrentScene(ctx context.Context, gameID model.ID, sceneID model.ID) error

// GetSceneActors 获取场景中的所有角色
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   sceneID - 场景ID
// 返回:
//   []ActorSnapshot - 场景中的角色列表
func (e *Engine) GetSceneActors(ctx context.Context, gameID model.ID, sceneID model.ID) ([]ActorSnapshot, error)

// AddItemToScene 在场景中放置物品
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   sceneID - 场景ID
//   item - 要放置的物品
// 返回:
//   error - 放置失败时返回错误
func (e *Engine) AddItemToScene(ctx context.Context, gameID model.ID, sceneID model.ID, item *model.Item) error

// GetSceneItems 获取场景地面上的物品
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   sceneID - 场景ID
// 返回:
//   []model.Item - 场景中的物品列表
func (e *Engine) GetSceneItems(ctx context.Context, gameID model.ID, sceneID model.ID) ([]model.Item, error)
```

### 状态查询

```go
// pkg/engine/state.go

// GetStateSummary 获取LLM友好的游戏状态摘要
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
// 返回:
//   *StateSummary - 游戏状态摘要
//   error - 获取失败时返回错误
// 使用场景: 向LLM展示当前游戏整体状态
// 权限: 所有阶段都允许
func (e *Engine) GetStateSummary(ctx context.Context, gameID model.ID) (*StateSummary, error)

// GetActorSheet 获取角色的完整角色卡
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
//   actorID - 角色ID
// 返回:
//   *ActorSheet - 完整角色卡
//   error - 获取失败时返回错误
func (e *Engine) GetActorSheet(ctx context.Context, gameID model.ID, actorID model.ID) (*ActorSheet, error)

// GetCombatSummary 获取战斗摘要
// 参数:
//   ctx - 上下文
//   gameID - 游戏会话ID
// 返回:
//   *CombatSummary - 战斗摘要
//   error - 没有活跃战斗时返回错误
func (e *Engine) GetCombatSummary(ctx context.Context, gameID model.ID) (*CombatSummary, error)

type StateSummary struct {
    GameName     string          `json:"game_name"`
    Phase        Phase           `json:"phase"`
    CurrentScene *SceneSummary   `json:"current_scene"`
    PartyMembers []ActorSnapshot `json:"party_members"`
    ActiveCombat *CombatSummary  `json:"active_combat"`
    ActiveQuests []QuestSummary  `json:"active_quests"`
    Time         string          `json:"time"`
}

type ActorSheet struct {
    BasicInfo     string              `json:"basic_info"`
    AbilityScores map[string]int      `json:"ability_scores"`
    Skills        map[string]int      `json:"skills"`
    SavingThrows  map[string]int      `json:"saving_throws"`
    Combat        CombatSheetInfo     `json:"combat"`
    Spellcasting  *SpellSheetInfo     `json:"spellcasting,omitempty"`
    Equipment     []EquipmentEntry    `json:"equipment"`
    Conditions    []string            `json:"conditions"`
    Features      []string            `json:"features"`
}

type CombatSummary struct {
    Round        int              `json:"round"`
    TurnOrder    []TurnOrderEntry `json:"turn_order"`
    CurrentActor string           `json:"current_actor"`
    Combatants   []CombatantBrief `json:"combatants"`
}

type TurnOrderEntry struct {
    ActorName  string `json:"actor_name"`
    Initiative int    `json:"initiative"`
    IsCurrent  bool   `json:"is_current"`
}

type CombatantBrief struct {
    Name        string   `json:"name"`
    Type        string   `json:"type"`
    HP          int      `json:"hp"`
    MaxHP       int      `json:"max_hp"`
    AC          int      `json:"ac"`
    Conditions  []string `json:"conditions"`
    IsDefeated  bool     `json:"is_defeated"`
}
```

### 错误定义

```go
// pkg/engine/errors.go

var (
    ErrNotFound              = errors.New("entity not found")
    ErrAlreadyExists         = errors.New("entity already exists")
    ErrInvalidState          = errors.New("invalid game state for this operation")
    ErrCombatNotActive       = errors.New("no active combat")
    ErrCombatAlreadyActive   = errors.New("combat is already active")
    ErrNotYourTurn           = errors.New("it is not this actor's turn")
    ErrActionAlreadyUsed     = errors.New("action has already been used this turn")
    ErrInsufficientSlots     = errors.New("insufficient spell slots")
    ErrInvalidTarget         = errors.New("invalid target for this action")
    ErrOutOfRange            = errors.New("target is out of range")
    ErrNoLineOfSight         = errors.New("no line of sight to target")
    ErrConcentrationBroken   = errors.New("concentration check failed")
    ErrActorIncapacitated    = errors.New("actor is incapacitated")
    ErrInvalidDiceExpression = errors.New("invalid dice expression")
    ErrStorageError          = errors.New("storage operation failed")
    ErrValidationFailed      = errors.New("validation failed")
    ErrPhaseNotAllowed       = errors.New("operation not allowed in current phase")
)

// EngineError 包装错误并附加上下文
type EngineError struct {
    Op      string         // 操作名称
    Err     error          // 原始错误
    Phase   Phase          // 当前阶段
    Details map[string]any // 额外详情
}

func (e *EngineError) Error() string
func (e *EngineError) Unwrap() error
```

---

## 规则引擎实现

### 核心计算 (internal/rules/calculator.go)

```go
// 属性修正值: floor((score - 10) / 2)
func AbilityModifier(score int) int

// 熟练加值: 1-4级+2, 5-8级+3, 9-12级+4, 13-16级+5, 17-20级+6
func ProficiencyBonus(totalLevel int) int

// 技能修正值: 属性修正 + (熟练加值 if 熟练) + (2×熟练加值 if 专家) + 其他加值
func SkillModifier(abilityScore int, isProficient bool, isExpert bool, bonus int) int

// 法术豁免DC: 8 + 熟练加值 + 施法属性修正
func SpellSaveDC(spellcastingAbilityScore int, proficiencyBonus int) int

// AC计算
func ArmorClass(armorType ArmorType, armorAC int, maxDexMod int, dexModifier int, shieldBonus int, otherBonus int) int

// 被动检定: 10 + 总修正值
func PassiveScore(totalModifier int, modifier RollModifier) int
```

### d20系统与伤害计算

严格按照D&D 5e规则实现：
- 优势/劣势处理
- 自然20/1特殊处理
- 伤害计算顺序：掷骰 → 修正 → 弱点(×2) → 抗性(÷2) → 免疫(→0)

---

## 存储机制

### 存储接口

```go
type Store interface {
    SaveGame(ctx context.Context, game *model.GameState) error
    LoadGame(ctx context.Context, gameID model.ID) (*model.GameState, error)
    DeleteGame(ctx context.Context, gameID model.ID) error
    ListGames(ctx context.Context) ([]GameMeta, error)
    UpdateGame(ctx context.Context, gameID model.ID, fn func(*model.GameState) error) error
    Init(ctx context.Context) error
    Close(ctx context.Context) error
}
```

---

## 并发模型

- **单写多读** - 使用sync.RWMutex
- **操作级锁** - 每个Engine方法获取适当的锁
- **深拷贝返回** - 防止外部修改内部状态
- **Context支持** - 所有公开API接受context.Context

---

## 模块依赖关系

```
pkg/engine/  ──────►  internal/*
internal/combat/ ───► internal/model/, internal/rules/, internal/dice/
internal/spell/  ───► internal/model/, internal/rules/
internal/rules/  ───► internal/model/（仅类型）
internal/dice/   ───►（无内部导入）
internal/storage/ ──► internal/model/
internal/model/  ───►（无 - 叶子包）
```

**严格DAG，无循环依赖**

---

## 实现步骤

### Phase 1: 基础架构
- [ ] `internal/model/` - 所有数据模型定义
- [ ] `internal/rules/` - 纯计算函数
- [ ] `internal/dice/` - 骰子引擎
- [ ] `internal/storage/` - 存储接口和实现
- [ ] `internal/data/` - 静态数据嵌入
- [ ] `pkg/engine/engine.go` - 引擎生命周期

### Phase 2: 阶段管理与角色系统
- [ ] `pkg/engine/phase.go` - 阶段定义和权限控制
- [ ] PC创建API
- [ ] NPC和Enemy管理
- [ ] 角色CRUD API
- [ ] 状态查询API

### Phase 3: 战斗系统
- [ ] 战斗状态机
- [ ] 先攻排序
- [ ] 回合管理
- [ ] 攻击与伤害
- [ ] 死亡豁免

### Phase 4: 法术与物品
- [ ] 法术位管理
- [ ] 施法解析
- [ ] 专注追踪
- [ ] 库存管理
- [ ] 装备槽系统

### Phase 5: 任务与场景
- [ ] 场景图管理
- [ ] 角色移动
- [ ] 任务创建和追踪
- [ ] 休息机制
- [ ] 升级系统

### Phase 6: 完善与测试
- [ ] LLM友好的输出格式化
- [ ] 错误处理优化
- [ ] 全面测试套件

---

## 验证方案

1. **单元测试** - 每个规则函数独立测试
2. **集成测试** - 完整游戏流程测试
3. **规则验证** - 对照D&D 5e规则书验证
4. **并发测试** - `go test -race ./...`
5. **存储测试** - 持久化验证

---

## 关键文件

- `internal/model/actor.go` - 核心数据模型
- `internal/rules/calculator.go` - 规则引擎核心
- `internal/model/combat.go` - 战斗数据结构
- `internal/storage/store.go` - 存储接口
- `pkg/engine/engine.go` - 公共API入口点
- `pkg/engine/phase.go` - 阶段管理与权限控制（新增）

---

## 设计模式

1. **Builder模式** - 角色创建
2. **状态机** - 战斗系统
3. **规约模式** - 规则验证
4. **观察者模式** - 游戏事件
5. **纯函数** - 规则计算
6. **权限矩阵** - 阶段操作控制（新增）
