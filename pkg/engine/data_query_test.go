package engine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zwh8800/dnd-core/pkg/model"
)

// ============================================================================
// 测试：种族查询
// ============================================================================

func TestListRaces(t *testing.T) {
	t.Run("lists all races with pagination", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.ListRaces(ctx, ListRacesRequest{})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotEmpty(t, result.Races)
		assert.Greater(t, result.Pagination.TotalCount, 0)
		assert.Equal(t, 1, result.Pagination.Page)
		assert.Equal(t, 20, result.Pagination.PageSize)
	})

	t.Run("pagination with custom page size", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.ListRaces(ctx, ListRacesRequest{
			Pagination: &PaginationRequest{
				Page:     1,
				PageSize: 3,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Races, 3)
		assert.True(t, result.Pagination.HasNext)
		assert.False(t, result.Pagination.HasPrev)
	})

	t.Run("pagination second page", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.ListRaces(ctx, ListRacesRequest{
			Pagination: &PaginationRequest{
				Page:     2,
				PageSize: 3,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.Pagination.HasPrev)
	})

	t.Run("page beyond available data returns empty", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.ListRaces(ctx, ListRacesRequest{
			Pagination: &PaginationRequest{
				Page:     100,
				PageSize: 10,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Empty(t, result.Races)
	})
}

func TestGetRace(t *testing.T) {
	t.Run("gets race by name", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.GetRace(ctx, GetRaceRequest{
			Name: "人类",
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "人类", result.Race.Name)
		assert.Equal(t, 30, result.Race.Speed)
		assert.NotEmpty(t, result.Race.AbilityBonuses)
	})

	t.Run("returns error for non-existent race", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		_, err := e.GetRace(ctx, GetRaceRequest{
			Name: "不存在的种族",
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "不存在的种族")
	})
}

// ============================================================================
// 测试：职业查询
// ============================================================================

func TestListClasses(t *testing.T) {
	t.Run("lists all classes with pagination", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.ListClasses(ctx, ListClassesRequest{})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotEmpty(t, result.Classes)
		assert.Greater(t, result.Pagination.TotalCount, 0)
	})

	t.Run("pagination respects page size", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.ListClasses(ctx, ListClassesRequest{
			Pagination: &PaginationRequest{
				Page:     1,
				PageSize: 5,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.LessOrEqual(t, len(result.Classes), 5)
	})
}

func TestGetClass(t *testing.T) {
	t.Run("gets class by ID", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.GetClass(ctx, GetClassRequest{
			ID: "战士",
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "战士", result.Class.Name)
		assert.NotEmpty(t, result.Class.PrimaryAbilities)
	})

	t.Run("returns error for non-existent class", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		_, err := e.GetClass(ctx, GetClassRequest{
			ID: "不存在的职业",
		})
		assert.Error(t, err)
	})
}

// ============================================================================
// 测试：背景查询
// ============================================================================

func TestListBackgrounds(t *testing.T) {
	t.Run("lists all backgrounds with pagination", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.ListBackgrounds(ctx, ListBackgroundsRequest{})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotEmpty(t, result.Backgrounds)
		assert.Greater(t, result.Pagination.TotalCount, 0)
	})

	t.Run("pagination works correctly", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.ListBackgrounds(ctx, ListBackgroundsRequest{
			Pagination: &PaginationRequest{
				Page:     1,
				PageSize: 2,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.LessOrEqual(t, len(result.Backgrounds), 2)
	})
}

func TestGetBackground(t *testing.T) {
	t.Run("gets background by name", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.GetBackground(ctx, GetBackgroundRequest{
			ID: "acolyte",
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "侍僧", result.Background.Name)
	})

	t.Run("returns error for non-existent background", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		_, err := e.GetBackground(ctx, GetBackgroundRequest{
			ID: "nonexistent-background",
		})
		assert.Error(t, err)
	})
}

// ============================================================================
// 测试：专长查询
// ============================================================================

func TestListFeatsData(t *testing.T) {
	t.Run("lists all feats with pagination", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.ListFeatsData(ctx, ListFeatsDataRequest{})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotEmpty(t, result.Feats)
		assert.Greater(t, result.Pagination.TotalCount, 0)
	})

	t.Run("pagination with small page size", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.ListFeatsData(ctx, ListFeatsDataRequest{
			Pagination: &PaginationRequest{
				Page:     1,
				PageSize: 2,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Feats, 2)
		assert.True(t, result.Pagination.HasNext)
	})

	t.Run("default pagination applied", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.ListFeatsData(ctx, ListFeatsDataRequest{})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 1, result.Pagination.Page)
		assert.Equal(t, 20, result.Pagination.PageSize)
	})
}

func TestGetFeatData(t *testing.T) {
	t.Run("gets feat by ID", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.GetFeatData(ctx, GetFeatDataRequest{
			ID: "alert",
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "警觉", result.Feat.Name)
		assert.Equal(t, "alert", result.Feat.ID)
		assert.NotEmpty(t, result.Feat.Description)
	})

	t.Run("returns error for non-existent feat", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		_, err := e.GetFeatData(ctx, GetFeatDataRequest{
			ID: "nonexistent-feat",
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "nonexistent-feat")
	})
}

// ============================================================================
// 测试：怪物查询
// ============================================================================

func TestListMonsters(t *testing.T) {
	t.Run("lists all monsters with pagination", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.ListMonsters(ctx, ListMonstersRequest{})
		require.NoError(t, err)
		require.NotNil(t, result)
		// 怪物数据可能为空，所以不强制要求非空
		assert.GreaterOrEqual(t, result.Pagination.TotalCount, 0)
	})

	t.Run("pagination works with page 2", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		// 先获取总数
		allResult, err := e.ListMonsters(ctx, ListMonstersRequest{})
		require.NoError(t, err)

		result, err := e.ListMonsters(ctx, ListMonstersRequest{
			Pagination: &PaginationRequest{
				Page:     2,
				PageSize: 5,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, result)

		// 验证分页信息
		if allResult.Pagination.TotalCount > 5 {
			assert.True(t, result.Pagination.HasPrev)
		}
	})
}

func TestGetMonster(t *testing.T) {
	t.Run("gets monster by ID", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		// 先获取怪物列表
		listResult, err := e.ListMonsters(ctx, ListMonstersRequest{})
		require.NoError(t, err)

		// 如果没有怪物数据，跳过测试
		if len(listResult.Monsters) == 0 {
			t.Skip("no monster data available")
		}

		monsterID := listResult.Monsters[0].ID

		result, err := e.GetMonster(ctx, GetMonsterRequest{
			ID: monsterID,
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, monsterID, result.Monster.ID)
		assert.NotEmpty(t, result.Monster.Name)
	})

	t.Run("returns error for non-existent monster", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		_, err := e.GetMonster(ctx, GetMonsterRequest{
			ID: "nonexistent-monster",
		})
		assert.Error(t, err)
	})
}

// ============================================================================
// 测试：法术查询
// ============================================================================

func TestListSpells(t *testing.T) {
	t.Run("lists all spells with pagination", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.ListSpells(ctx, ListSpellsRequest{})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotEmpty(t, result.Spells)
		assert.Greater(t, result.Pagination.TotalCount, 0)
	})

	t.Run("pagination with page size 10", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.ListSpells(ctx, ListSpellsRequest{
			Pagination: &PaginationRequest{
				Page:     1,
				PageSize: 10,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.LessOrEqual(t, len(result.Spells), 10)
	})
}

func TestGetSpell(t *testing.T) {
	t.Run("gets spell by ID", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.GetSpell(ctx, GetSpellRequest{
			ID: "fire-bolt",
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "火焰箭", result.Spell.Name)
		assert.Equal(t, 0, result.Spell.Level)
	})

	t.Run("returns error for non-existent spell", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		_, err := e.GetSpell(ctx, GetSpellRequest{
			ID: "nonexistent-spell",
		})
		assert.Error(t, err)
	})
}

// ============================================================================
// 测试：武器查询
// ============================================================================

func TestListWeapons(t *testing.T) {
	t.Run("lists all weapons with pagination", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.ListWeapons(ctx, ListWeaponsRequest{})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotEmpty(t, result.Weapons)
		assert.Greater(t, result.Pagination.TotalCount, 0)
	})
}

func TestGetWeapon(t *testing.T) {
	t.Run("gets weapon by ID", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		// 先获取武器列表
		listResult, err := e.ListWeapons(ctx, ListWeaponsRequest{})
		require.NoError(t, err)
		require.NotEmpty(t, listResult.Weapons)

		weaponID := listResult.Weapons[0].ID

		result, err := e.GetWeapon(ctx, GetWeaponRequest{
			ID: weaponID,
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, weaponID, result.Weapon.ID)
	})

	t.Run("returns error for non-existent weapon", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		_, err := e.GetWeapon(ctx, GetWeaponRequest{
			ID: "nonexistent-weapon",
		})
		assert.Error(t, err)
	})
}

// ============================================================================
// 测试：护甲查询
// ============================================================================

func TestListArmors(t *testing.T) {
	t.Run("lists all armors with pagination", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.ListArmors(ctx, ListArmorsRequest{})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotEmpty(t, result.Armors)
		assert.Greater(t, result.Pagination.TotalCount, 0)
	})
}

func TestGetArmor(t *testing.T) {
	t.Run("gets armor by ID", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		// 先获取护甲列表
		listResult, err := e.ListArmors(ctx, ListArmorsRequest{})
		require.NoError(t, err)
		require.NotEmpty(t, listResult.Armors)

		armorID := listResult.Armors[0].ID

		result, err := e.GetArmor(ctx, GetArmorRequest{
			ID: armorID,
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, armorID, result.Armor.ID)
	})

	t.Run("returns error for non-existent armor", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		_, err := e.GetArmor(ctx, GetArmorRequest{
			ID: "nonexistent-armor",
		})
		assert.Error(t, err)
	})
}

// ============================================================================
// 测试：魔法物品查询
// ============================================================================

func TestListMagicItems(t *testing.T) {
	t.Run("lists all magic items with pagination", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.ListMagicItems(ctx, ListMagicItemsRequest{})
		require.NoError(t, err)
		require.NotNil(t, result)
		// 魔法物品可能为空，所以不强制要求非空
		assert.GreaterOrEqual(t, result.Pagination.TotalCount, 0)
	})
}

func TestGetMagicItem(t *testing.T) {
	t.Run("returns error for non-existent magic item", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		_, err := e.GetMagicItem(ctx, GetMagicItemRequest{
			ID: "nonexistent-magic-item",
		})
		assert.Error(t, err)
	})
}

// ============================================================================
// 测试：装备查询
// ============================================================================

func TestListGears(t *testing.T) {
	t.Run("lists all gears with pagination", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.ListGears(ctx, ListGearsRequest{})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, result.Pagination.TotalCount, 0)
	})
}

func TestGetGear(t *testing.T) {
	t.Run("returns error for non-existent gear", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		_, err := e.GetGear(ctx, GetGearRequest{
			ID: "nonexistent-gear",
		})
		assert.Error(t, err)
	})
}

// ============================================================================
// 测试：工具查询
// ============================================================================

func TestListTools(t *testing.T) {
	t.Run("lists all tools with pagination", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.ListTools(ctx, ListToolsRequest{})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, result.Pagination.TotalCount, 0)
	})
}

func TestGetTool(t *testing.T) {
	t.Run("returns error for non-existent tool", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		_, err := e.GetTool(ctx, GetToolRequest{
			ID: "nonexistent-tool",
		})
		assert.Error(t, err)
	})
}

// ============================================================================
// 测试：配方查询
// ============================================================================

func TestListRecipes(t *testing.T) {
	t.Run("lists all recipes with pagination", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.ListRecipes(ctx, ListRecipesRequest{})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, result.Pagination.TotalCount, 0)
	})

	t.Run("pagination works correctly", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.ListRecipes(ctx, ListRecipesRequest{
			Pagination: &PaginationRequest{
				Page:     1,
				PageSize: 5,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.LessOrEqual(t, len(result.Recipes), 5)
	})
}

func TestGetRecipe(t *testing.T) {
	t.Run("returns error for non-existent recipe", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		_, err := e.GetRecipe(ctx, GetRecipeRequest{
			ID: "nonexistent-recipe",
		})
		assert.Error(t, err)
	})
}

// ============================================================================
// 测试：生活方式查询
// ============================================================================

func TestListLifestylesData(t *testing.T) {
	t.Run("lists all lifestyles with pagination", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.ListLifestylesData(ctx, ListLifestylesRequest{})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotEmpty(t, result.Lifestyles)
		assert.Greater(t, result.Pagination.TotalCount, 0)
	})
}

func TestGetLifestyleData(t *testing.T) {
	t.Run("gets lifestyle by tier", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.GetLifestyleData(ctx, GetLifestyleDataRequest{
			Tier: "modest",
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, model.LifestyleTier("modest"), result.Lifestyle.Tier)
	})

	t.Run("returns error for non-existent lifestyle", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		_, err := e.GetLifestyleData(ctx, GetLifestyleDataRequest{
			Tier: "nonexistent-tier",
		})
		assert.Error(t, err)
	})
}

// ============================================================================
// 测试：坐骑查询
// ============================================================================

func TestListMounts(t *testing.T) {
	t.Run("lists all mounts with pagination", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.ListMounts(ctx, ListMountsRequest{})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotEmpty(t, result.Mounts)
		assert.Greater(t, result.Pagination.TotalCount, 0)
	})
}

func TestGetMount(t *testing.T) {
	t.Run("gets mount by ID", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		// 先获取坐骑列表
		listResult, err := e.ListMounts(ctx, ListMountsRequest{})
		require.NoError(t, err)
		require.NotEmpty(t, listResult.Mounts)

		mountID := listResult.Mounts[0].ID

		result, err := e.GetMount(ctx, GetMountRequest{
			ID: mountID,
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, mountID, result.Mount.ID)
	})

	t.Run("returns error for non-existent mount", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		_, err := e.GetMount(ctx, GetMountRequest{
			ID: "nonexistent-mount",
		})
		assert.Error(t, err)
	})
}

// ============================================================================
// 测试：毒药查询
// ============================================================================

func TestListPoisons(t *testing.T) {
	t.Run("lists all poisons with pagination", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.ListPoisons(ctx, ListPoisonsRequest{})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotEmpty(t, result.Poisons)
		assert.Greater(t, result.Pagination.TotalCount, 0)
	})

	t.Run("pagination with page size 2", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.ListPoisons(ctx, ListPoisonsRequest{
			Pagination: &PaginationRequest{
				Page:     1,
				PageSize: 2,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.LessOrEqual(t, len(result.Poisons), 2)
	})
}

func TestGetPoison(t *testing.T) {
	t.Run("gets poison by ID", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		// 先获取毒药列表
		listResult, err := e.ListPoisons(ctx, ListPoisonsRequest{})
		require.NoError(t, err)
		require.NotEmpty(t, listResult.Poisons)

		poisonID := listResult.Poisons[0].ID

		result, err := e.GetPoison(ctx, GetPoisonRequest{
			ID: poisonID,
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, poisonID, result.Poison.ID)
		assert.NotEmpty(t, result.Poison.Name)
	})

	t.Run("returns error for non-existent poison", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		_, err := e.GetPoison(ctx, GetPoisonRequest{
			ID: "nonexistent-poison",
		})
		assert.Error(t, err)
	})
}

// ============================================================================
// 测试：陷阱查询
// ============================================================================

func TestListTraps(t *testing.T) {
	t.Run("lists all traps with pagination", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.ListTraps(ctx, ListTrapsRequest{})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotEmpty(t, result.Traps)
		assert.Greater(t, result.Pagination.TotalCount, 0)
	})
}

func TestGetTrap(t *testing.T) {
	t.Run("gets trap by ID", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		// 先获取陷阱列表
		listResult, err := e.ListTraps(ctx, ListTrapsRequest{})
		require.NoError(t, err)
		require.NotEmpty(t, listResult.Traps)

		trapID := listResult.Traps[0].ID

		result, err := e.GetTrap(ctx, GetTrapRequest{
			ID: trapID,
		})
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, trapID, result.Trap.ID)
		assert.NotEmpty(t, result.Trap.Name)
	})

	t.Run("returns error for non-existent trap", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		_, err := e.GetTrap(ctx, GetTrapRequest{
			ID: "nonexistent-trap",
		})
		assert.Error(t, err)
	})
}

// ============================================================================
// 测试：分页工具函数
// ============================================================================

func TestPaginationHelpers(t *testing.T) {
	t.Run("applyDefaults sets default values", func(t *testing.T) {
		req := &PaginationRequest{}
		req.applyDefaults()

		assert.Equal(t, 1, req.Page)
		assert.Equal(t, 20, req.PageSize)
	})

	t.Run("applyDefaults respects custom values", func(t *testing.T) {
		req := &PaginationRequest{
			Page:     3,
			PageSize: 10,
		}
		req.applyDefaults()

		assert.Equal(t, 3, req.Page)
		assert.Equal(t, 10, req.PageSize)
	})

	t.Run("applyDefaults caps max page size", func(t *testing.T) {
		req := &PaginationRequest{
			Page:     1,
			PageSize: 200,
		}
		req.applyDefaults()

		assert.Equal(t, 1, req.Page)
		assert.Equal(t, 100, req.PageSize)
	})

	t.Run("calculatePagination computes correctly", func(t *testing.T) {
		req := &PaginationRequest{
			Page:     1,
			PageSize: 10,
		}

		info := calculatePagination(25, req)

		assert.Equal(t, 1, info.Page)
		assert.Equal(t, 10, info.PageSize)
		assert.Equal(t, 25, info.TotalCount)
		assert.Equal(t, 3, info.TotalPages)
		assert.True(t, info.HasNext)
		assert.False(t, info.HasPrev)
	})

	t.Run("calculatePagination handles exact division", func(t *testing.T) {
		req := &PaginationRequest{
			Page:     2,
			PageSize: 10,
		}

		info := calculatePagination(20, req)

		assert.Equal(t, 2, info.Page)
		assert.Equal(t, 2, info.TotalPages)
		assert.False(t, info.HasNext)
		assert.True(t, info.HasPrev)
	})

	t.Run("paginateSlice returns correct page", func(t *testing.T) {
		items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		req := &PaginationRequest{
			Page:     1,
			PageSize: 3,
		}

		result := paginateSlice(items, req)

		assert.Equal(t, []int{1, 2, 3}, result)
	})

	t.Run("paginateSlice returns empty for out of range page", func(t *testing.T) {
		items := []int{1, 2, 3}
		req := &PaginationRequest{
			Page:     10,
			PageSize: 2,
		}

		result := paginateSlice(items, req)

		assert.Empty(t, result)
	})

	t.Run("paginateSlice handles nil request", func(t *testing.T) {
		items := []int{1, 2, 3, 4, 5}
		result := paginateSlice(items, nil)

		// 默认 pageSize=20，所以应该返回所有项
		assert.Len(t, result, 5)
	})
}

// ============================================================================
// 测试：数据查询操作在所有阶段都允许
// ============================================================================

func TestDataQueryOperationsInAllPhases(t *testing.T) {
	phases := []string{
		"character_creation",
		"exploration",
		"combat",
		"downtime",
	}

	for _, phase := range phases {
		t.Run("data queries work in "+phase, func(t *testing.T) {
			e := NewTestEngine(t)
			ctx := context.Background()

			// 创建游戏
			gameResult, err := e.NewGame(ctx, NewGameRequest{
				Name:        "Test Game",
				Description: "Test data queries",
			})
			require.NoError(t, err)

			// 设置阶段
			_, err = e.SetPhase(ctx, gameResult.Game.ID, model.Phase(phase), "testing")
			require.NoError(t, err)

			// 测试种族查询
			result, err := e.ListRaces(ctx, ListRacesRequest{})
			require.NoError(t, err)
			assert.NotEmpty(t, result.Races)

			// 测试专长查询
			featResult, err := e.ListFeatsData(ctx, ListFeatsDataRequest{})
			require.NoError(t, err)
			assert.NotEmpty(t, featResult.Feats)

			// 测试法术查询
			spellResult, err := e.ListSpells(ctx, ListSpellsRequest{})
			require.NoError(t, err)
			assert.NotEmpty(t, spellResult.Spells)
		})
	}
}
