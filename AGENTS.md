# AGENTS.md

This file provides guidance to Qoder (qoder.com) when working with code in this repository.

## Project Overview

DND-Core is a comprehensive **Dungeons & Dragons 5e Game Engine** written in Go. It serves as a headless game rule engine and state management system for LLM-based DM (Dungeon Master) applications, enforcing D&D 5e rules accurately with 178+ public APIs.

**Module:** `github.com/zwh8800/dnd-core`  
**Go Version:** 1.24.2

## Development Commands

### Build
```bash
go build ./...
```

### Run All Tests
```bash
go test ./... -v
```

### Run Single Test
```bash
go test ./pkg/engine -run TestCombatEncounters -v
```

### Run Tests for Specific Package
```bash
go test ./pkg/engine/... -v
go test ./pkg/engine/testsuite/... -v
```

### Check for Compilation Errors
```bash
go vet ./...
```

### Run with Race Detection
```bash
go test ./... -race -v
```

## Code Architecture

### Package Structure

The codebase follows a clean separation of concerns with 5 main packages:

#### `pkg/engine/` - Public API (70 files, 30.5K LOC)
The single entry point for all external consumers. Contains the `Engine` struct which provides all 178+ public APIs organized by feature:

**Key Files by Feature:**
- `engine.go` - Core Engine struct, initialization, and lifecycle
- `game.go` - Game session management (NewGame, LoadGame, SaveGame, DeleteGame)
- `actor.go` - Character/actor management (CreatePC, CreateNPC, CreateEnemy, GetActor, UpdateActor)
- `combat.go` - Combat system (StartCombat, EndCombat, NextTurn, ExecuteAction, ExecuteAttack, ExecuteDamage)
- `spell.go` - Spell system (CastSpell, GetSpellSlots, PrepareSpells, ConcentrationCheck)
- `check.go` - Ability/skill/saving throw checks
- `inventory.go` - Inventory and equipment management
- `movement.go` - Actor movement and positioning
- `scene.go` - Scene/location management
- `quest.go` - Quest tracking system
- `death_saves.go` - Death saving throws
- `exhaustion.go` - Exhaustion levels
- `feat.go` - Feat system
- `multiclass.go` - Multiclass support
- `magic_items.go` - Magic item handling
- `crafting.go` - Crafting system
- `exploration.go` - Travel and exploration
- `environment.go` - Environmental effects
- `trap.go`, `poison.go`, `curse.go` - Hazards
- `data_query.go` - Data lookup queries (ListRaces, GetClass, ListSpells, etc.)
- `state.go` - State summary and snapshots

**Design Pattern:** All public methods follow the pattern:
```go
func (e *Engine) MethodName(ctx context.Context, params...) (Result, error)
```

The engine uses `sync.RWMutex` for thread-safe concurrent access.

#### `pkg/model/` - Data Models (42 files, 4.8K LOC)
Defines all D&D 5e data structures:

**Core Types:**
- `Actor` - Base class for all game entities (PC, NPC, Enemy, Companion)
- `Creature` - Creature type definitions
- `Combat` - Combat state, initiative, turn management
- `Spell`, `SpellSlot` - Spell mechanics
- `Equipment`, `Weapon`, `Armor` - Item system
- `Condition` - Status effects
- `Ability`, `Skill` - Character abilities and skills
- `Scene`, `Quest` - World state
- `Rest` - Rest mechanics

**Class Hooks:** Each D&D class has dedicated hooks files (e.g., `barbarian_hooks.go`, `wizard_hooks.go`) for class-specific abilities.

#### `pkg/rules/` - Rule Engine (20 files)
Pure D&D 5e rule implementations, separated from state management:

- `combat_rules.go` - Combat mechanics, opportunity attacks, grappling
- `attack.go` - Attack roll calculations, critical hits
- `calculator.go` - Modifier calculations, proficiency bonuses, AC, HP
- `check.go` - Check calculations
- `spelleffects.go` - Spell effects
- `death.go`, `exhaustion.go` - Death and exhaustion rules
- `rest.go` - Rest recovery calculations
- `multiclass.go` - Multiclass spell slot calculations
- `constants.go` - D&D constants (DC tables, XP thresholds, proficiency bonuses)

#### `pkg/data/` - Game Content (17 files)
Pre-loaded D&D 5e content database:

- `classes.go`, `races.go` - Character options
- `spells.go` - Complete spell database (52K LOC)
- `feats.go` - Feat database (20K LOC)
- `magicitems.go` - Magic items (24K LOC)
- `monsters.go`, `weapons.go`, `armors.go` - Game data
- `registry.go` - Data registry and initialization (14.8K LOC)

#### `pkg/dice/` & `pkg/storage/`
- `dice/roller.go` - Dice expression parser (d20, 2d6+3, 4d6kh3, advantage/disadvantage)
- `storage/store.go`, `storage/memory.go` - Storage abstraction layer with in-memory default

### Key Architectural Principles

1. **Single Entry Point**: All external interactions go through `pkg/engine.Engine`
2. **Rule Enforcement**: `pkg/rules/` contains pure functions for rule calculations
3. **Pluggable Storage**: Storage interface allows future database backends
4. **Concurrent Safety**: Engine uses `sync.RWMutex` for thread-safe access
5. **Context Support**: All public methods accept `context.Context` for cancellation/timeout
6. **Error Handling**: Custom error types in `pkg/engine/error.go`
7. **Mixed Language**: Code comments are in Chinese (中文), variable names in English

## Testing Strategy

- Tests are co-located with source files (`*_test.go`)
- 35 test files with 17.1K LOC of tests
- Integration tests in `pkg/engine/testsuite/`
- Use `testify` for assertions
- Test naming follows pattern: `Test{Feature}` with subtests using `t.Run()`

## Important Notes

- **No Makefile**: Use direct `go` commands
- **Chinese Documentation**: Many comments and some docs are in Chinese
- **SRD Reference**: Complete D&D 5e SRD documentation in `docs/SRD-md/` (90+ markdown files)
- **API Documentation**: See `docs/engine-api.md` for complete API reference (1,280 lines)
- **Large Data Files**: Spell/feat/magic item data files are 20K-52K LOC each
- **No Linting Configured**: No golangci-lint or similar tools configured yet
