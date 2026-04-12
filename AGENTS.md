# AGENTS.md

This file provides guidance to Qoder (qoder.com) when working with code in this repository.

## Project Overview

**dnd-core** is a D&D 5e game engine library written in Go 1.24.2. It provides a complete implementation of D&D 5e rules including character creation, combat, spells, exploration, and session management.

Repository: `github.com/zwh8800/dnd-core`

## Development Commands

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests for specific package
go test ./pkg/engine
go test ./pkg/rules
go test ./pkg/model

# Run a single test
go test ./pkg/engine -run TestNew

# Run tests with verbose output
go test -v ./pkg/engine

# Run tests with coverage
go test -cover ./...
```

### Building

This is a library project with no main package. Build and verify with:

```bash
# Build all packages
go build ./...

# Vet code for common issues
go vet ./...

# Format code
go fmt ./...
```

### Dependencies

```bash
# Add dependency
go get <package>

# Tidy dependencies
go mod tidy
```

## Architecture

### Layered Design

```
pkg/engine/          - Engine Layer: Public API, game state orchestration
    ↓
pkg/rules/           - Rules Layer: D&D 5e rule implementations
pkg/data/            - Data Layer: Game data registry (races, classes, spells, etc.)
    ↓
pkg/model/           - Model Layer: Core data structures and type definitions
    ↓
pkg/storage/         - Storage Layer: Persistence abstraction
pkg/dice/            - Dice Layer: Randomization and dice rolling
```

### Key Packages

**pkg/engine** - Core game engine providing:
- Thread-safe `Engine` struct as the single entry point
- 50+ public methods for all game operations
- Phase-based permission system (CharacterCreation, Exploration, Combat, Rest)
- Request/Result pattern for all API calls
- Automatic state persistence via storage layer
- File organization: organized by game system (actor.go, combat.go, spell.go, etc.)

**pkg/model** - Data models only (no business logic):
- Core types: PlayerCharacter, NPC, Enemy, Companion, Actor
- Game types: GameState, Combat, Scene, Quest
- D&D types: Ability, Skill, Feat, Spell, Item, Condition
- Class-specific hooks in `*_hooks.go` files (12 classes)
- Helper methods for convenience (IsAlive, IsDead, etc.)

**pkg/rules** - D&D 5e rule implementations:
- Attack rolls and damage calculation
- Check and saving throw mechanics
- Combat rules (cover, grapple, shove, opportunity attacks)
- Death saves, exhaustion, rest mechanics
- Spell effects and multiclass calculations
- Social interaction and exploration rules

**pkg/data** - Global game data registry:
- Thread-safe singleton: `data.GlobalRegistry`
- Register/Get/List pattern for all data types
- Contains: 12 classes, races, backgrounds, 100+ feats, 400+ spells, 100+ magic items, monsters, etc.
- Uses `data.GlobalRegistry.GetRace()`, `GetFeat()`, etc. for lookups
- Never modify - read-only access only

**pkg/storage** - Persistence abstraction:
- `Store` interface for pluggable backends
- Current implementation: in-memory store
- Methods: SaveGame, LoadGame, DeleteGame, ListGames, UpdateGame
- Deep copies for isolation between operations

**pkg/dice** - Dice rolling utility:
- D&D notation parser: "4d6kh3", "2d20", "adv", "dis"
- Supports keep-high, keep-low, drop modifiers
- Seeded random for reproducible tests
- Advantage/disadvantage mechanics

### Game Phase System

The engine uses a **phase-based permission model** to control what operations are allowed when.

**Phases:**
- `PhaseCharacterCreation` - Building characters (initial phase)
- `PhaseExploration` - Traveling, interaction, general gameplay
- `PhaseCombat` - Active combat encounters
- `PhaseRest` - Long/short rests

**Phase Transitions:**
```
CharacterCreation → Exploration ↔ Combat
                        ↓↑
                      Rest
```

**Permission System:**
- Each phase has an allowed operations map in `phase.go`
- Data queries allowed in ALL phases
- Engine methods must call `e.checkPermission(game.Phase, OpYourOperation)`
- Phase transitions trigger automatic actions (e.g., combat cleanup, XP rewards)

## API Design Rules (CRITICAL)

All new APIs in `pkg/engine` MUST follow these rules (see `pkg/engine/api-guidelines.md` for full spec):

### 1. Engine Method Pattern

```go
// ✅ Correct: Engine method with Request/Result
func (e *Engine) OperationName(ctx context.Context, req OperationRequest) (*OperationResult, error)

// ❌ Wrong: Standalone function
func SelectFeat(pc *model.PlayerCharacter, featID string) error
```

### 2. Standard Method Structure

Every engine method must:
1. Acquire lock (`e.mu.Lock()` for writes, `e.mu.RLock()` for reads)
2. Load game state via `e.loadGame(ctx, req.GameID)`
3. Check permission via `e.checkPermission(game.Phase, OpYourOperation)`
4. Execute business logic
5. Save game state via `e.saveGame(ctx, game)` (if modified)
6. Return result

### 3. Request/Result Pattern

```go
type OperationRequest struct {
    GameID  model.ID `json:"game_id"`  // 游戏会话ID（必填）
    ActorID model.ID `json:"actor_id"` // 角色ID
    // 其他字段...
}

type OperationResult struct {
    // 返回的数据使用 Info 结构体封装
    Actor *ActorInfo `json:"actor"`
}
```

### 4. Info Struct Encapsulation

Never return raw `model` types directly. Use Info structs:

```go
type ActorInfo struct {
    ID         model.ID        `json:"id"`
    Type       model.ActorType `json:"type"`
    Name       string          `json:"name"`
    HitPoints  model.HitPoints `json:"hit_points"`
}
```

### 5. Adding Permissions

When adding a new operation:
1. Define constant in `phase.go`: `const OpYourOperation Operation = "your_operation"`
2. Add to `phasePermissions` map for allowed phases

### 6. Checklist for New APIs

- [ ] Defined `<Operation>Request` struct with json tags
- [ ] Defined `<Operation>Result` struct (if applicable)
- [ ] Method is on `(e *Engine)` receiver
- [ ] Correct mutex (Lock for writes, RLock for reads)
- [ ] Calls `e.loadGame()` and `e.checkPermission()`
- [ ] Calls `e.saveGame()` after state modifications
- [ ] All fields have json tags and Chinese comments
- [ ] Added permission constant and phase mapping

## Testing Patterns

### Test Structure

- Use table-driven tests with `t.Run()`
- Use `testify/assert` and `testify/require` for assertions
- Use `engine.NewTestEngine(t)` for test engine (auto-cleanup, deterministic dice)

### Example Test

```go
func TestCreatePC(t *testing.T) {
    e := engine.NewTestEngine(t)
    ctx := context.Background()

    t.Run("creates character successfully", func(t *testing.T) {
        gameResult, err := e.NewGame(ctx, engine.NewGameRequest{
            Name:        "Test Game",
            Description: "Test",
        })
        require.NoError(t, err)

        result, err := e.CreatePC(ctx, engine.CreatePCRequest{
            GameID: gameResult.Game.ID,
            Name:   "Test Character",
            // ...
        })
        require.NoError(t, err)
        assert.NotEmpty(t, result.PC.ID)
    })
}
```

## Important Patterns

### Thread Safety

All Engine methods are concurrent-safe. Never modify game state directly - always use Engine methods which handle:
- Locking via `sync.RWMutex`
- Permission checks
- State persistence
- Validation

### Data Access

Use `data.GlobalRegistry` for all game data lookups:
```go
feat, exists := data.GlobalRegistry.GetFeat(featID)
race, exists := data.GlobalRegistry.GetRace(raceID)
class, exists := data.GlobalRegistry.GetClass(classID)
```

Never hardcode game data values.

### Error Handling

Use predefined errors from `error.go`:
- `ErrNotFound` - Resource not found
- `ErrInvalidState` - Invalid state
- `ErrPhaseNotAllowed` - Phase permission denied

Return formatted errors with context:
```go
return nil, fmt.Errorf("feat not found: %s", req.FeatID)
```

### File Organization

Engine layer is organized by game system:
- `actor.go` - Character/NPC management
- `combat.go` - Combat system
- `spell.go` - Spell system
- `check.go` - Ability/skill checks
- `inventory.go` - Item management
- `feat.go` - Feat selection
- `phase.go` - Phase and permission definitions
- `scene.go` - Scene management
- Other files for specific systems (quest, social, exploration, etc.)

Each file contains:
1. Request/Result struct definitions
2. Info struct definitions
3. Engine method implementations
4. Helper conversion functions (e.g., `actorToInfo`)

### Character Classes

The game supports 12 D&D 5e classes, each with class-specific hook files in `pkg/model/`:

| Class | Key Features | Hook File |
|-------|--------------|-----------|
| Barbarian | Rage, Unarmored Defense, Reckless Attack | `barbarian_hooks.go` |
| Bard | Bardic Inspiration, Spellcasting (CHA) | `bard_hooks.go` |
| Cleric | Channel Divinity, Spellcasting (WIS) | `cleric_hooks.go` |
| Druid | Wild Shape, Spellcasting (WIS) | `druid_hooks.go` |
| Fighter | Action Surge, Fighting Style, Extra Attack | `fighter_hooks.go` |
| Monk | Ki Points, Flurry of Blows, Stunning Strike | `monk_hooks.go` |
| Paladin | Divine Smite, Aura of Protection | `paladin_hooks.go` |
| Ranger | Favored Enemy, Spellcasting (WIS) | `ranger_hooks.go` |
| Rogue | Sneak Attack, Cunning Action, Expertise | `rogue_hooks.go` |
| Sorcerer | Sorcery Points, Metamagic | `sorcerer_hooks.go` |
| Warlock | Pact Magic, Invocations, Eldritch Blast | `warlock_hooks.go` |
| Wizard | Spellcasting (INT), Arcane Recovery | `wizard_hooks.go` |

### Combat System Flow

Combat follows this lifecycle:
1. `StartCombat()` - Initialize combat, set phase to PhaseCombat, generate initiative
2. `NextTurn()` - Progress initiative order, reset actions
3. `ExecuteAction()` - Actor takes action (attack, spell, move, etc.)
4. `ExecuteAttack()` - Roll attack vs AC, trigger damage on hit
5. `ExecuteDamage()` - Apply damage, check for unconsciousness/death saves
6. `EndCombat()` - Apply XP rewards, return to PhaseExploration

### Configuration and Initialization

```go
type Config struct {
    Storage  storage.Store  // Defaults to in-memory if nil
    DiceSeed int64          // 0 = random, fixed = deterministic for tests
    DataPath string         // Empty = built-in data only
}

// Usage
cfg := engine.DefaultConfig()
cfg.DiceSeed = 42  // For deterministic tests
e, err := engine.New(cfg)
```

## Comment Conventions

- Use Chinese comments for business logic and API documentation
- All exported types and methods must have doc comments
- Struct field comments should include purpose and required/optional status
- Method comments should document Parameters and Returns sections

Example:
```go
// SelectFeatRequest 选择专长请求
// 用于描述该操作的具体用途（多行说明）
type SelectFeatRequest struct {
    GameID model.ID `json:"game_id"` // 游戏会话ID（必填）
    FeatID string   `json:"feat_id"` // 专长ID（必填）
}
```

## Prohibited Patterns

1. **Never modify game state directly** - Use Engine methods only
2. **Never skip permission checks** - All operations must validate phase
3. **Never return raw model types** - Use Info struct encapsulation
4. **Never use standalone functions** - All APIs must be Engine methods
5. **Never forget to save state** - Always call `e.saveGame()` after modifications
6. **Never skip locking** - All methods must acquire appropriate mutex
7. **Never hardcode game data** - Use `data.GlobalRegistry` for lookups
