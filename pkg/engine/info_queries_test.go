package engine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zwh8800/dnd-core/pkg/model"
)

func TestGetLifestyleInfo(t *testing.T) {
	t.Run("gets lifestyle info for modest lifestyle", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		// Create a game
		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test lifestyle",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Get lifestyle info for modest
		result, err := e.GetLifestyleInfo(ctx, GetLifestyleRequest{
			GameID:    gameID,
			Lifestyle: model.LifestyleModest,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, model.LifestyleModest, result.Lifestyle)
		assert.NotEmpty(t, result.Description)
	})

	t.Run("gets lifestyle info for wealthy lifestyle", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test lifestyle",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		result, err := e.GetLifestyleInfo(ctx, GetLifestyleRequest{
			GameID:    gameID,
			Lifestyle: model.LifestyleWealthy,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, model.LifestyleWealthy, result.Lifestyle)
		assert.NotEmpty(t, result.Description)
	})

	t.Run("returns default modest lifestyle when not specified", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test lifestyle",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		result, err := e.GetLifestyleInfo(ctx, GetLifestyleRequest{
			GameID: gameID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, model.LifestyleModest, result.Lifestyle)
	})

	t.Run("returns error for non-existent game", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		_, err := e.GetLifestyleInfo(ctx, GetLifestyleRequest{
			GameID:    "non-existent",
			Lifestyle: model.LifestyleModest,
		})

		assert.Error(t, err)
	})

	t.Run("all lifestyle tiers are queryable", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test lifestyle",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		lifestyles := []model.LifestyleTier{
			model.LifestyleWretched,
			model.LifestyleSqualid,
			model.LifestylePoor,
			model.LifestyleModest,
			model.LifestyleComfortable,
			model.LifestyleWealthy,
			model.LifestyleAristocratic,
		}

		for _, lifestyle := range lifestyles {
			t.Run(string(lifestyle), func(t *testing.T) {
				result, err := e.GetLifestyleInfo(ctx, GetLifestyleRequest{
					GameID:    gameID,
					Lifestyle: lifestyle,
				})

				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, lifestyle, result.Lifestyle)
				assert.NotEmpty(t, result.Description)
			})
		}
	})
}

func TestGetCraftingInfo(t *testing.T) {
	t.Run("gets crafting info for an item", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		result, err := e.GetCraftingInfo(ctx, GetCraftingInfoRequest{
			ItemName: "Longsword",
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "Longsword", result.ItemName)
		assert.NotEmpty(t, result.Description)
	})

	t.Run("gets crafting info for different items", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		items := []string{"Shield", "Chain Mail", "Shortbow", "Leather Armor"}

		for _, itemName := range items {
			t.Run(itemName, func(t *testing.T) {
				result, err := e.GetCraftingInfo(ctx, GetCraftingInfoRequest{
					ItemName: itemName,
				})

				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, itemName, result.ItemName)
			})
		}
	})
}

func TestGetCarryingCapacity(t *testing.T) {
	t.Run("gets carrying capacity for PC", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		// Create a game
		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test carrying capacity",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Create a PC with known strength
		pcResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Strong Character",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     20, // STR 20 = 300 lbs carrying capacity
					Dexterity:    14,
					Constitution: 13,
					Intelligence: 12,
					Wisdom:       10,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)

		result, err := e.GetCarryingCapacity(ctx, GetCarryingCapacityRequest{
			GameID:  gameID,
			ActorID: pcResult.Actor.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, pcResult.Actor.ID, result.ActorID)
		assert.Equal(t, 20, result.Strength)
		assert.Equal(t, 300, result.CarryingCapacity) // 20 * 15
		assert.Equal(t, 600, result.PushDragLift)     // 20 * 30
		assert.Equal(t, 0, result.CurrentWeight)
		assert.False(t, result.IsEncumbered)
		assert.False(t, result.IsHeavilyEncumbered)
	})

	t.Run("gets carrying capacity for different strength values", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test carrying capacity",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		tests := []struct {
			name             string
			strength         int
			expectedCapacity int
			expectedPushLift int
		}{
			{"STR 10", 10, 150, 300},
			{"STR 15", 15, 225, 450},
			{"STR 18", 18, 270, 540},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				pcResult, err := e.CreatePC(ctx, CreatePCRequest{
					GameID: gameID,
					PC: &PlayerCharacterInput{
						Name:  "Test Character",
						Race:  "Human",
						Class: "Fighter",
						Level: 1,
						AbilityScores: AbilityScoresInput{
							Strength:     tt.strength,
							Dexterity:    14,
							Constitution: 13,
							Intelligence: 12,
							Wisdom:       10,
							Charisma:     8,
						},
					},
				})
				require.NoError(t, err)

				result, err := e.GetCarryingCapacity(ctx, GetCarryingCapacityRequest{
					GameID:  gameID,
					ActorID: pcResult.Actor.ID,
				})

				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.strength, result.Strength)
				assert.Equal(t, tt.expectedCapacity, result.CarryingCapacity)
				assert.Equal(t, tt.expectedPushLift, result.PushDragLift)
			})
		}
	})

	t.Run("gets carrying capacity for NPC", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test carrying capacity",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		npcResult, err := e.CreateNPC(ctx, CreateNPCRequest{
			GameID: gameID,
			NPC: &NPCInput{
				Name: "Test NPC",
				AbilityScores: AbilityScoresInput{
					Strength:     16,
					Dexterity:    12,
					Constitution: 14,
					Intelligence: 10,
					Wisdom:       13,
					Charisma:     11,
				},
			},
		})
		require.NoError(t, err)

		result, err := e.GetCarryingCapacity(ctx, GetCarryingCapacityRequest{
			GameID:  gameID,
			ActorID: npcResult.Actor.ID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 16, result.Strength)
		assert.Equal(t, 240, result.CarryingCapacity) // 16 * 15
		assert.Equal(t, 480, result.PushDragLift)     // 16 * 30
	})

	t.Run("returns error for non-existent actor", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test carrying capacity",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		_, err = e.GetCarryingCapacity(ctx, GetCarryingCapacityRequest{
			GameID:  gameID,
			ActorID: "non-existent",
		})

		assert.Error(t, err)
	})

	t.Run("returns error for non-existent game", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		_, err := e.GetCarryingCapacity(ctx, GetCarryingCapacityRequest{
			GameID:  "non-existent",
			ActorID: "some-actor",
		})

		assert.Error(t, err)
	})
}

func TestGetExhaustionEffects(t *testing.T) {
	t.Run("gets exhaustion effects for level 0", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		// Create a game (required for loadGame)
		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test exhaustion effects",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		// Need to be in a valid phase
		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		result, err := e.GetExhaustionEffects(ctx, GetExhaustionEffectsRequest{
			ExhaustionLevel: 0,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 0, result.ExhaustionLevel)
		assert.Empty(t, result.Effects)
		assert.False(t, result.IsDead)
		assert.Contains(t, result.Description, "力竭等级 0")
	})

	t.Run("gets exhaustion effects for level 1", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test exhaustion effects",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		result, err := e.GetExhaustionEffects(ctx, GetExhaustionEffectsRequest{
			ExhaustionLevel: 1,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 1, result.ExhaustionLevel)
		assert.Len(t, result.Effects, 1)
		assert.Contains(t, result.Effects[0], "属性检定劣势")
		assert.False(t, result.IsDead)
	})

	t.Run("gets exhaustion effects for level 3", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test exhaustion effects",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		result, err := e.GetExhaustionEffects(ctx, GetExhaustionEffectsRequest{
			ExhaustionLevel: 3,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 3, result.ExhaustionLevel)
		assert.Len(t, result.Effects, 3)
		assert.Contains(t, result.Effects[0], "属性检定劣势")
		assert.Contains(t, result.Effects[1], "速度减半")
		assert.Contains(t, result.Effects[2], "攻击检定和豁免检定劣势")
		assert.False(t, result.IsDead)
	})

	t.Run("gets exhaustion effects for level 6 (death)", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test exhaustion effects",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		result, err := e.GetExhaustionEffects(ctx, GetExhaustionEffectsRequest{
			ExhaustionLevel: 6,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 6, result.ExhaustionLevel)
		assert.Len(t, result.Effects, 6)
		assert.True(t, result.IsDead)
		assert.Contains(t, result.Description, "死亡")
	})

	t.Run("returns error for invalid exhaustion level", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test exhaustion effects",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		t.Run("negative level", func(t *testing.T) {
			_, err := e.GetExhaustionEffects(ctx, GetExhaustionEffectsRequest{
				ExhaustionLevel: -1,
			})
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid exhaustion level")
		})

		t.Run("level 7", func(t *testing.T) {
			_, err := e.GetExhaustionEffects(ctx, GetExhaustionEffectsRequest{
				ExhaustionLevel: 7,
			})
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "invalid exhaustion level")
		})
	})

	t.Run("all valid exhaustion levels work", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "Test exhaustion effects",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "testing")
		require.NoError(t, err)

		for level := 0; level <= 6; level++ {
			t.Run(string(rune('0'+level)), func(t *testing.T) {
				result, err := e.GetExhaustionEffects(ctx, GetExhaustionEffectsRequest{
					ExhaustionLevel: level,
				})

				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, level, result.ExhaustionLevel)
				assert.Len(t, result.Effects, level)

				if level >= 6 {
					assert.True(t, result.IsDead)
				} else {
					assert.False(t, result.IsDead)
				}
			})
		}
	})
}
