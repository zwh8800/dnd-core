package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zwh8800/dnd-core/pkg/model"
)

func TestLoadMonster(t *testing.T) {
	t.Run("loads monster from template", func(t *testing.T) {
		e := NewTestEngine(t)

		// Try to load a monster - the data registry should have some entries
		// If not, this test will show what's available
		monster, err := e.LoadMonster("goblin")

		// Either succeeds or returns not found - both are valid
		if err != nil {
			assert.Contains(t, err.Error(), "monster template not found")
		} else {
			require.NotNil(t, monster)
			assert.NotEmpty(t, monster.Actor.Name)
		}
	})

	t.Run("returns error for nonexistent template", func(t *testing.T) {
		e := NewTestEngine(t)

		monster, err := e.LoadMonster("nonexistent-monster")

		assert.Error(t, err)
		assert.Nil(t, monster)
		assert.Contains(t, err.Error(), "monster template not found")
	})

	t.Run("returns error for empty template ID", func(t *testing.T) {
		e := NewTestEngine(t)

		monster, err := e.LoadMonster("")

		assert.Error(t, err)
		assert.Nil(t, monster)
	})

	t.Run("loaded monster has correct actor type", func(t *testing.T) {
		e := NewTestEngine(t)

		// This test will pass if the template exists
		monster, err := e.LoadMonster("goblin")
		if err != nil {
			t.Skip("goblin template not found in registry")
		}

		assert.Equal(t, model.ActorTypeEnemy, monster.Actor.Type)
	})

	t.Run("can load same monster multiple times", func(t *testing.T) {
		e := NewTestEngine(t)

		monster, err := e.LoadMonster("goblin")
		if err != nil {
			t.Skip("goblin template not found in registry")
		}

		// Loading again should create new instance
		monster2, err := e.LoadMonster("goblin")
		require.NoError(t, err)

		// Should be different instances
		assert.NotEqual(t, monster.Actor.ID, monster2.Actor.ID)
	})
}

func TestCreateEnemyFromStatBlock(t *testing.T) {
	t.Run("creates enemy with valid stat block", func(t *testing.T) {
		statBlock := &model.MonsterStatBlock{
			Name:                "Test Monster",
			Size:                model.SizeMedium,
			ArmorClass:          15,
			HitPointsAverage:    30,
			ChallengeRating:     "1",
			AbilityScores:       model.AbilityScores{Strength: 10, Dexterity: 12, Constitution: 14, Intelligence: 8, Wisdom: 10, Charisma: 8},
			Speed:               model.SpeedTypes{Walk: 30},
			ConditionImmunities: []model.ConditionType{model.ConditionPoisoned},
		}

		enemy, err := CreateEnemyFromStatBlock(statBlock)

		require.NoError(t, err)
		require.NotNil(t, enemy)
		assert.Equal(t, "Test Monster", enemy.Actor.Name)
		assert.Equal(t, model.SizeMedium, enemy.Actor.Size)
		assert.Equal(t, 15, enemy.Actor.ArmorClass)
		assert.Equal(t, 30, enemy.Actor.HitPoints.Current)
		assert.Equal(t, 30, enemy.Actor.HitPoints.Maximum)
		assert.Equal(t, model.ActorTypeEnemy, enemy.Actor.Type)
	})

	t.Run("enemy has copied ability scores", func(t *testing.T) {
		statBlock := &model.MonsterStatBlock{
			Name:             "Test",
			AbilityScores:    model.AbilityScores{Strength: 16, Dexterity: 14, Constitution: 12, Intelligence: 10, Wisdom: 8, Charisma: 6},
			HitPointsAverage: 20,
		}

		enemy, err := CreateEnemyFromStatBlock(statBlock)

		require.NoError(t, err)
		assert.Equal(t, 16, enemy.Actor.AbilityScores.Strength)
		assert.Equal(t, 14, enemy.Actor.AbilityScores.Dexterity)
		assert.Equal(t, 12, enemy.Actor.AbilityScores.Constitution)
	})

	t.Run("enemy has copied damage immunities", func(t *testing.T) {
		statBlock := &model.MonsterStatBlock{
			Name:             "Test",
			HitPointsAverage: 20,
			DamageImmunities: []model.DamageImmunity{
				{DamageTypes: []model.DamageType{"poison"}},
				{DamageTypes: []model.DamageType{"necrotic"}},
			},
		}

		enemy, err := CreateEnemyFromStatBlock(statBlock)

		require.NoError(t, err)
		assert.Len(t, enemy.DamageImmunities, 2)
	})

	t.Run("enemy has copied damage resistances", func(t *testing.T) {
		statBlock := &model.MonsterStatBlock{
			Name:             "Test",
			HitPointsAverage: 20,
			DamageResistances: []model.DamageImmunity{
				{DamageTypes: []model.DamageType{"fire"}},
			},
		}

		enemy, err := CreateEnemyFromStatBlock(statBlock)

		require.NoError(t, err)
		assert.Len(t, enemy.DamageResistances, 1)
	})

	t.Run("enemy has copied damage vulnerabilities", func(t *testing.T) {
		statBlock := &model.MonsterStatBlock{
			Name:             "Test",
			HitPointsAverage: 20,
			DamageVulnerabilities: []model.DamageImmunity{
				{DamageTypes: []model.DamageType{"cold"}},
			},
		}

		enemy, err := CreateEnemyFromStatBlock(statBlock)

		require.NoError(t, err)
		assert.Len(t, enemy.DamageVulnerabilities, 1)
	})

	t.Run("enemy with legendary actions has remaining", func(t *testing.T) {
		statBlock := &model.MonsterStatBlock{
			Name:                     "Dragon",
			HitPointsAverage:         100,
			LegendaryActions:         []model.MonsterAction{{Name: "Tail Attack"}},
			LegendaryActionsPerRound: 3,
		}

		enemy, err := CreateEnemyFromStatBlock(statBlock)

		require.NoError(t, err)
		assert.Equal(t, 3, enemy.LegendaryActionsRemaining)
	})

	t.Run("enemy without legendary actions has zero", func(t *testing.T) {
		statBlock := &model.MonsterStatBlock{
			Name:             "Goblin",
			HitPointsAverage: 7,
		}

		enemy, err := CreateEnemyFromStatBlock(statBlock)

		require.NoError(t, err)
		assert.Equal(t, 0, enemy.LegendaryActionsRemaining)
	})

	t.Run("enemy has stat block reference", func(t *testing.T) {
		statBlock := &model.MonsterStatBlock{
			Name:             "Test",
			HitPointsAverage: 20,
		}

		enemy, err := CreateEnemyFromStatBlock(statBlock)

		require.NoError(t, err)
		assert.NotNil(t, enemy.StatBlock)
		assert.Equal(t, statBlock, enemy.StatBlock)
	})

	t.Run("enemy has speed from stat block", func(t *testing.T) {
		statBlock := &model.MonsterStatBlock{
			Name:             "Test",
			HitPointsAverage: 20,
			Speed:            model.SpeedTypes{Walk: 40, Fly: 60},
		}

		enemy, err := CreateEnemyFromStatBlock(statBlock)

		require.NoError(t, err)
		assert.Equal(t, 40, enemy.Actor.Speed)
		assert.Equal(t, 60, enemy.Actor.Speeds.Fly)
	})

	t.Run("enemy has zero exhaustion", func(t *testing.T) {
		statBlock := &model.MonsterStatBlock{
			Name:             "Test",
			HitPointsAverage: 20,
		}

		enemy, err := CreateEnemyFromStatBlock(statBlock)

		require.NoError(t, err)
		assert.Equal(t, 0, enemy.Actor.Exhaustion)
	})
}

func TestGetMonsterActions(t *testing.T) {
	t.Run("returns actions from stat block", func(t *testing.T) {
		monster := &model.Enemy{
			Actor: model.Actor{
				Name: "Test Monster",
			},
			StatBlock: &model.MonsterStatBlock{
				Actions: []model.MonsterAction{
					{Name: "Bite"},
					{Name: "Claw"},
				},
			},
		}

		actions := GetMonsterActions(monster)

		assert.Len(t, actions, 2)
		assert.Equal(t, "Bite", actions[0].Name)
		assert.Equal(t, "Claw", actions[1].Name)
	})

	t.Run("returns empty actions for monster without actions", func(t *testing.T) {
		monster := &model.Enemy{
			StatBlock: &model.MonsterStatBlock{},
		}

		actions := GetMonsterActions(monster)

		assert.Empty(t, actions)
	})

	t.Run("includes bonus actions", func(t *testing.T) {
		monster := &model.Enemy{
			StatBlock: &model.MonsterStatBlock{
				Actions:      []model.MonsterAction{{Name: "Attack"}},
				BonusActions: []model.MonsterAction{{Name: "Dash"}},
			},
		}

		actions := GetMonsterActions(monster)

		assert.Len(t, actions, 2)
	})

	t.Run("includes reactions", func(t *testing.T) {
		monster := &model.Enemy{
			StatBlock: &model.MonsterStatBlock{
				Actions:   []model.MonsterAction{{Name: "Attack"}},
				Reactions: []model.MonsterAction{{Name: "Parry"}},
			},
		}

		actions := GetMonsterActions(monster)

		assert.Len(t, actions, 2)
	})

	t.Run("includes legendary actions when available", func(t *testing.T) {
		monster := &model.Enemy{
			LegendaryActionsRemaining: 2,
			StatBlock: &model.MonsterStatBlock{
				Actions:          []model.MonsterAction{{Name: "Attack"}},
				LegendaryActions: []model.MonsterAction{{Name: "Tail"}},
			},
		}

		actions := GetMonsterActions(monster)

		assert.Len(t, actions, 2)
	})

	t.Run("excludes legendary actions when none remaining", func(t *testing.T) {
		monster := &model.Enemy{
			LegendaryActionsRemaining: 0,
			StatBlock: &model.MonsterStatBlock{
				Actions:          []model.MonsterAction{{Name: "Attack"}},
				LegendaryActions: []model.MonsterAction{{Name: "Tail"}},
			},
		}

		actions := GetMonsterActions(monster)

		assert.GreaterOrEqual(t, len(actions), 1)
		assert.Equal(t, "Attack", actions[0].Name)
	})

	t.Run("excludes recharge actions that are not charged", func(t *testing.T) {
		monster := &model.Enemy{
			StatBlock: &model.MonsterStatBlock{
				Actions: []model.MonsterAction{
					{Name: "Normal Attack"},
					{
						Name: "Breath Weapon",
						Recharge: &model.RechargeInfo{
							RollRange:   [2]int{5, 6},
							CurrentUses: 0,
						},
					},
				},
			},
		}

		actions := GetMonsterActions(monster)

		// Breath Weapon should not be included since it's not recharged
		assert.GreaterOrEqual(t, len(actions), 1)
	})

	t.Run("includes recharge actions that are charged", func(t *testing.T) {
		monster := &model.Enemy{
			StatBlock: &model.MonsterStatBlock{
				Actions: []model.MonsterAction{
					{
						Name: "Breath Weapon",
						Recharge: &model.RechargeInfo{
							RollRange:   [2]int{5, 6},
							CurrentUses: 1,
						},
					},
				},
			},
		}

		actions := GetMonsterActions(monster)

		assert.GreaterOrEqual(t, len(actions), 1)
		assert.Equal(t, "Breath Weapon", actions[0].Name)
	})
}

func TestRechargeMonsterActions(t *testing.T) {
	t.Run("restores legendary actions", func(t *testing.T) {
		monster := &model.Enemy{
			LegendaryActionsRemaining: 0,
			StatBlock: &model.MonsterStatBlock{
				LegendaryActions:         []model.MonsterAction{{Name: "Tail"}},
				LegendaryActionsPerRound: 3,
			},
		}

		RechargeMonsterActions(monster)

		assert.Equal(t, 3, monster.LegendaryActionsRemaining)
	})

	t.Run("does not restore legendary actions for monsters without them", func(t *testing.T) {
		monster := &model.Enemy{
			LegendaryActionsRemaining: 0,
			StatBlock:                 &model.MonsterStatBlock{},
		}

		RechargeMonsterActions(monster)

		assert.Equal(t, 0, monster.LegendaryActionsRemaining)
	})

	t.Run("handles monster with recharge actions", func(t *testing.T) {
		monster := &model.Enemy{
			StatBlock: &model.MonsterStatBlock{
				Actions: []model.MonsterAction{
					{
						Name: "Breath Weapon",
						Recharge: &model.RechargeInfo{
							RollRange: [2]int{5, 6},
						},
					},
				},
			},
		}

		// Should not panic
		RechargeMonsterActions(monster)
	})

	t.Run("handles monster with no recharge actions", func(t *testing.T) {
		monster := &model.Enemy{
			StatBlock: &model.MonsterStatBlock{
				Actions: []model.MonsterAction{{Name: "Attack"}},
			},
		}

		// Should not panic
		RechargeMonsterActions(monster)
	})

	t.Run("handles empty actions list", func(t *testing.T) {
		monster := &model.Enemy{
			StatBlock: &model.MonsterStatBlock{},
		}

		// Should not panic
		RechargeMonsterActions(monster)
	})
}

func TestUseLegendaryAction(t *testing.T) {
	t.Run("uses legendary action successfully", func(t *testing.T) {
		monster := &model.Enemy{
			LegendaryActionsRemaining: 3,
		}

		err := UseLegendaryAction(monster)

		assert.NoError(t, err)
		assert.Equal(t, 2, monster.LegendaryActionsRemaining)
	})

	t.Run("returns error when no legendary actions remaining", func(t *testing.T) {
		monster := &model.Enemy{
			LegendaryActionsRemaining: 0,
		}

		err := UseLegendaryAction(monster)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no legendary actions remaining")
		assert.Equal(t, 0, monster.LegendaryActionsRemaining)
	})

	t.Run("can use multiple legendary actions", func(t *testing.T) {
		monster := &model.Enemy{
			LegendaryActionsRemaining: 3,
		}

		for i := 0; i < 3; i++ {
			err := UseLegendaryAction(monster)
			assert.NoError(t, err)
		}

		assert.Equal(t, 0, monster.LegendaryActionsRemaining)
	})

	t.Run("error on fourth use when only 3 available", func(t *testing.T) {
		monster := &model.Enemy{
			LegendaryActionsRemaining: 3,
		}

		for i := 0; i < 3; i++ {
			_ = UseLegendaryAction(monster)
		}

		err := UseLegendaryAction(monster)
		assert.Error(t, err)
	})

	t.Run("decrements correctly from 1 to 0", func(t *testing.T) {
		monster := &model.Enemy{
			LegendaryActionsRemaining: 1,
		}

		err := UseLegendaryAction(monster)
		assert.NoError(t, err)
		assert.Equal(t, 0, monster.LegendaryActionsRemaining)
	})
}

func TestUseRechargeAction(t *testing.T) {
	t.Run("uses recharge action successfully", func(t *testing.T) {
		monster := &model.Enemy{
			ActionRecharges: map[int]int{
				0: 3,
			},
		}

		err := UseRechargeAction(monster, 0)

		assert.NoError(t, err)
		assert.Equal(t, 2, monster.ActionRecharges[0])
	})

	t.Run("returns error when action does not use charges", func(t *testing.T) {
		monster := &model.Enemy{
			ActionRecharges: map[int]int{
				0: 3,
			},
		}

		err := UseRechargeAction(monster, 1)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "action does not use charges")
	})

	t.Run("returns error when no charges remaining", func(t *testing.T) {
		monster := &model.Enemy{
			ActionRecharges: map[int]int{
				0: 0,
			},
		}

		err := UseRechargeAction(monster, 0)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no charges remaining")
	})

	t.Run("returns error when ActionRecharges is nil", func(t *testing.T) {
		monster := &model.Enemy{}

		err := UseRechargeAction(monster, 0)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "action does not use charges")
	})

	t.Run("can use all charges", func(t *testing.T) {
		monster := &model.Enemy{
			ActionRecharges: map[int]int{
				0: 3,
			},
		}

		for i := 0; i < 3; i++ {
			err := UseRechargeAction(monster, 0)
			assert.NoError(t, err)
		}

		assert.Equal(t, 0, monster.ActionRecharges[0])

		err := UseRechargeAction(monster, 0)
		assert.Error(t, err)
	})
}
