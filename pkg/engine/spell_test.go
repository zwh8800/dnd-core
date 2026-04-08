package engine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zwh8800/dnd-core/pkg/model"
)

func TestCastSpell(t *testing.T) {
	t.Run("fails for non-spellcaster", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for spellcasting",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Fighter",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     16,
					Dexterity:    12,
					Constitution: 14,
					Intelligence: 10,
					Wisdom:       12,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		_, err = e.CastSpell(ctx, CastSpellRequest{
			GameID:   gameID,
			CasterID: actorID,
			Spell: SpellInput{
				SpellID: "firebolt",
			},
		})

		assert.Error(t, err)
	})

	t.Run("fails for invalid caster", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for spellcasting",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		_, err = e.CastSpell(ctx, CastSpellRequest{
			GameID:   gameID,
			CasterID: model.ID("invalid-caster"),
			Spell: SpellInput{
				SpellID: "firebolt",
			},
		})

		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("fails for unknown spell", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for spellcasting",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Wizard",
				Race:  "Human",
				Class: "Wizard",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     8,
					Dexterity:    12,
					Constitution: 10,
					Intelligence: 18,
					Wisdom:       12,
					Charisma:     10,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		_, err = e.CastSpell(ctx, CastSpellRequest{
			GameID:   gameID,
			CasterID: actorID,
			Spell: SpellInput{
				SpellID: "nonexistent_spell",
			},
		})

		assert.Error(t, err)
	})

	t.Run("fails with invalid game", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		_, err := e.CastSpell(ctx, CastSpellRequest{
			GameID:   model.ID("invalid-game"),
			CasterID: model.ID("caster-id"),
			Spell: SpellInput{
				SpellID: "firebolt",
			},
		})

		assert.Error(t, err)
	})

	t.Run("fails for npc attempting to cast", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for spellcasting",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreateNPC(ctx, CreateNPCRequest{
			GameID: gameID,
			NPC: &NPCInput{
				Name: "Mage NPC",
				AbilityScores: AbilityScoresInput{
					Strength:     10,
					Dexterity:    12,
					Constitution: 10,
					Intelligence: 16,
					Wisdom:       12,
					Charisma:     10,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		_, err = e.CastSpell(ctx, CastSpellRequest{
			GameID:   gameID,
			CasterID: actorID,
			Spell: SpellInput{
				SpellID: "firebolt",
			},
		})

		assert.Error(t, err)
	})
}

func TestGetSpellSlots(t *testing.T) {
	t.Run("fails for non-spellcaster", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for spell slots",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Fighter",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     16,
					Dexterity:    12,
					Constitution: 14,
					Intelligence: 10,
					Wisdom:       12,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		_, err = e.GetSpellSlots(ctx, GetSpellSlotsRequest{
			GameID:   gameID,
			CasterID: actorID,
		})

		assert.Error(t, err)
	})

	t.Run("fails for invalid caster", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for spell slots",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		_, err = e.GetSpellSlots(ctx, GetSpellSlotsRequest{
			GameID:   gameID,
			CasterID: model.ID("invalid-caster"),
		})

		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("fails with invalid game", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		_, err := e.GetSpellSlots(ctx, GetSpellSlotsRequest{
			GameID:   model.ID("invalid-game"),
			CasterID: model.ID("caster-id"),
		})

		assert.Error(t, err)
	})

	t.Run("fails for npc", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for spell slots",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreateNPC(ctx, CreateNPCRequest{
			GameID: gameID,
			NPC: &NPCInput{
				Name: "Test NPC",
				AbilityScores: AbilityScoresInput{
					Strength:     10,
					Dexterity:    10,
					Constitution: 10,
					Intelligence: 10,
					Wisdom:       10,
					Charisma:     10,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		_, err = e.GetSpellSlots(ctx, GetSpellSlotsRequest{
			GameID:   gameID,
			CasterID: actorID,
		})

		assert.Error(t, err)
	})

	t.Run("fails for enemy", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for spell slots",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreateEnemy(ctx, CreateEnemyRequest{
			GameID: gameID,
			Enemy: &EnemyInput{
				Name:            "Dragon",
				AbilityScores:   AbilityScoresInput{Strength: 20, Dexterity: 14, Constitution: 18, Intelligence: 16, Wisdom: 14, Charisma: 18},
				ChallengeRating: 10,
				HitPoints:       200,
				ArmorClass:      19,
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		_, err = e.GetSpellSlots(ctx, GetSpellSlotsRequest{
			GameID:   gameID,
			CasterID: actorID,
		})

		assert.Error(t, err)
	})
}

func TestPrepareSpells(t *testing.T) {
	t.Run("fails for non-spellcaster", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for spell preparation",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Fighter",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     16,
					Dexterity:    12,
					Constitution: 14,
					Intelligence: 10,
					Wisdom:       12,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		err = e.PrepareSpells(ctx, PrepareSpellsRequest{
			GameID:   gameID,
			CasterID: actorID,
			SpellIDs: []string{"shield"},
		})

		assert.Error(t, err)
	})

	t.Run("fails for invalid caster", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for spell preparation",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		err = e.PrepareSpells(ctx, PrepareSpellsRequest{
			GameID:   gameID,
			CasterID: model.ID("invalid-caster"),
			SpellIDs: []string{"shield"},
		})

		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("fails with invalid game", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		err := e.PrepareSpells(ctx, PrepareSpellsRequest{
			GameID:   model.ID("invalid-game"),
			CasterID: model.ID("caster-id"),
			SpellIDs: []string{"shield"},
		})

		assert.Error(t, err)
	})

	t.Run("fails for npc", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for spell preparation",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreateNPC(ctx, CreateNPCRequest{
			GameID: gameID,
			NPC: &NPCInput{
				Name: "Test NPC",
				AbilityScores: AbilityScoresInput{
					Strength:     10,
					Dexterity:    10,
					Constitution: 10,
					Intelligence: 10,
					Wisdom:       10,
					Charisma:     10,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		err = e.PrepareSpells(ctx, PrepareSpellsRequest{
			GameID:   gameID,
			CasterID: actorID,
			SpellIDs: []string{"shield"},
		})

		assert.Error(t, err)
	})

	t.Run("fails for known-type caster", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for spell preparation",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Sorcerer",
				Race:  "Human",
				Class: "Sorcerer",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     8,
					Dexterity:    12,
					Constitution: 14,
					Intelligence: 10,
					Wisdom:       10,
					Charisma:     18,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		err = e.PrepareSpells(ctx, PrepareSpellsRequest{
			GameID:   gameID,
			CasterID: actorID,
			SpellIDs: []string{},
		})

		assert.Error(t, err)
	})
}

func TestLearnSpell(t *testing.T) {
	t.Run("fails for non-spellcaster", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for learning spells",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Fighter",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     16,
					Dexterity:    12,
					Constitution: 14,
					Intelligence: 10,
					Wisdom:       12,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		err = e.LearnSpell(ctx, LearnSpellRequest{
			GameID:   gameID,
			CasterID: actorID,
			SpellID:  "firebolt",
		})

		assert.Error(t, err)
	})

	t.Run("fails for invalid caster", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for learning spells",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		err = e.LearnSpell(ctx, LearnSpellRequest{
			GameID:   gameID,
			CasterID: model.ID("invalid-caster"),
			SpellID:  "firebolt",
		})

		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("fails with invalid game", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		err := e.LearnSpell(ctx, LearnSpellRequest{
			GameID:   model.ID("invalid-game"),
			CasterID: model.ID("caster-id"),
			SpellID:  "firebolt",
		})

		assert.Error(t, err)
	})

	t.Run("fails for npc", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for learning spells",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreateNPC(ctx, CreateNPCRequest{
			GameID: gameID,
			NPC: &NPCInput{
				Name: "Test NPC",
				AbilityScores: AbilityScoresInput{
					Strength:     10,
					Dexterity:    10,
					Constitution: 10,
					Intelligence: 10,
					Wisdom:       10,
					Charisma:     10,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		err = e.LearnSpell(ctx, LearnSpellRequest{
			GameID:   gameID,
			CasterID: actorID,
			SpellID:  "firebolt",
		})

		assert.Error(t, err)
	})

	t.Run("fails for enemy", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for learning spells",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreateEnemy(ctx, CreateEnemyRequest{
			GameID: gameID,
			Enemy: &EnemyInput{
				Name:            "Goblin",
				AbilityScores:   AbilityScoresInput{Strength: 10, Dexterity: 12, Constitution: 10, Intelligence: 8, Wisdom: 10, Charisma: 8},
				ChallengeRating: 0.25,
				HitPoints:       7,
				ArmorClass:      15,
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		err = e.LearnSpell(ctx, LearnSpellRequest{
			GameID:   gameID,
			CasterID: actorID,
			SpellID:  "firebolt",
		})

		assert.Error(t, err)
	})
}

func TestConcentrationCheck(t *testing.T) {
	t.Run("fails for invalid caster", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for concentration",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		_, err = e.ConcentrationCheck(ctx, ConcentrationCheckRequest{
			GameID:      gameID,
			CasterID:    model.ID("invalid-caster"),
			DamageTaken: 10,
		})

		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("fails with invalid game", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		_, err := e.ConcentrationCheck(ctx, ConcentrationCheckRequest{
			GameID:      model.ID("invalid-game"),
			CasterID:    model.ID("caster-id"),
			DamageTaken: 10,
		})

		assert.Error(t, err)
	})

	t.Run("fails for non-concentrating caster", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for concentration",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Wizard",
				Race:  "Human",
				Class: "Wizard",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     8,
					Dexterity:    12,
					Constitution: 14,
					Intelligence: 18,
					Wisdom:       12,
					Charisma:     10,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		_, err = e.ConcentrationCheck(ctx, ConcentrationCheckRequest{
			GameID:      gameID,
			CasterID:    actorID,
			DamageTaken: 10,
		})

		assert.Error(t, err)
	})

	t.Run("fails for npc", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for concentration",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreateNPC(ctx, CreateNPCRequest{
			GameID: gameID,
			NPC: &NPCInput{
				Name: "Test NPC",
				AbilityScores: AbilityScoresInput{
					Strength:     10,
					Dexterity:    10,
					Constitution: 10,
					Intelligence: 10,
					Wisdom:       10,
					Charisma:     10,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		_, err = e.ConcentrationCheck(ctx, ConcentrationCheckRequest{
			GameID:      gameID,
			CasterID:    actorID,
			DamageTaken: 10,
		})

		assert.Error(t, err)
	})

	t.Run("fails for enemy", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for concentration",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreateEnemy(ctx, CreateEnemyRequest{
			GameID: gameID,
			Enemy: &EnemyInput{
				Name:            "Goblin",
				AbilityScores:   AbilityScoresInput{Strength: 10, Dexterity: 12, Constitution: 10, Intelligence: 8, Wisdom: 10, Charisma: 8},
				ChallengeRating: 0.25,
				HitPoints:       7,
				ArmorClass:      15,
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		_, err = e.ConcentrationCheck(ctx, ConcentrationCheckRequest{
			GameID:      gameID,
			CasterID:    actorID,
			DamageTaken: 10,
		})

		assert.Error(t, err)
	})
}

func TestEndConcentration(t *testing.T) {
	t.Run("fails for invalid caster", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for concentration",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		err = e.EndConcentration(ctx, EndConcentrationRequest{
			GameID:   gameID,
			CasterID: model.ID("invalid-caster"),
		})

		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("fails with invalid game", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		err := e.EndConcentration(ctx, EndConcentrationRequest{
			GameID:   model.ID("invalid-game"),
			CasterID: model.ID("caster-id"),
		})

		assert.Error(t, err)
	})

	t.Run("fails for non-concentrating caster", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for concentration",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Wizard",
				Race:  "Human",
				Class: "Wizard",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     8,
					Dexterity:    12,
					Constitution: 14,
					Intelligence: 18,
					Wisdom:       12,
					Charisma:     10,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		err = e.EndConcentration(ctx, EndConcentrationRequest{
			GameID:   gameID,
			CasterID: actorID,
		})

		assert.Error(t, err)
	})

	t.Run("fails for npc", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for concentration",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreateNPC(ctx, CreateNPCRequest{
			GameID: gameID,
			NPC: &NPCInput{
				Name: "Test NPC",
				AbilityScores: AbilityScoresInput{
					Strength:     10,
					Dexterity:    10,
					Constitution: 10,
					Intelligence: 10,
					Wisdom:       10,
					Charisma:     10,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		err = e.EndConcentration(ctx, EndConcentrationRequest{
			GameID:   gameID,
			CasterID: actorID,
		})

		assert.Error(t, err)
	})

	t.Run("fails for enemy", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for concentration",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreateEnemy(ctx, CreateEnemyRequest{
			GameID: gameID,
			Enemy: &EnemyInput{
				Name:            "Goblin",
				AbilityScores:   AbilityScoresInput{Strength: 10, Dexterity: 12, Constitution: 10, Intelligence: 8, Wisdom: 10, Charisma: 8},
				ChallengeRating: 0.25,
				HitPoints:       7,
				ArmorClass:      15,
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		err = e.EndConcentration(ctx, EndConcentrationRequest{
			GameID:   gameID,
			CasterID: actorID,
		})

		assert.Error(t, err)
	})
}

func TestCastSpellRitual(t *testing.T) {
	t.Run("fails for invalid caster", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for ritual casting",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		_, err = e.CastSpellRitual(ctx, CastSpellRitualRequest{
			GameID:   gameID,
			CasterID: model.ID("invalid-caster"),
			SpellID:  "detect_magic",
		})

		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("fails with invalid game", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		_, err := e.CastSpellRitual(ctx, CastSpellRitualRequest{
			GameID:   model.ID("invalid-game"),
			CasterID: model.ID("caster-id"),
			SpellID:  "detect_magic",
		})

		assert.Error(t, err)
	})

	t.Run("fails for non-spellcaster", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for ritual casting",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Fighter",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     16,
					Dexterity:    12,
					Constitution: 14,
					Intelligence: 10,
					Wisdom:       12,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		_, err = e.CastSpellRitual(ctx, CastSpellRitualRequest{
			GameID:   gameID,
			CasterID: actorID,
			SpellID:  "detect_magic",
		})

		assert.Error(t, err)
	})

	t.Run("fails for npc", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for ritual casting",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreateNPC(ctx, CreateNPCRequest{
			GameID: gameID,
			NPC: &NPCInput{
				Name: "Test NPC",
				AbilityScores: AbilityScoresInput{
					Strength:     10,
					Dexterity:    10,
					Constitution: 10,
					Intelligence: 10,
					Wisdom:       10,
					Charisma:     10,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		_, err = e.CastSpellRitual(ctx, CastSpellRitualRequest{
			GameID:   gameID,
			CasterID: actorID,
			SpellID:  "detect_magic",
		})

		assert.Error(t, err)
	})

	t.Run("fails for unknown spell", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for ritual casting",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Wizard",
				Race:  "Human",
				Class: "Wizard",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     8,
					Dexterity:    12,
					Constitution: 14,
					Intelligence: 18,
					Wisdom:       12,
					Charisma:     10,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		_, err = e.CastSpellRitual(ctx, CastSpellRitualRequest{
			GameID:   gameID,
			CasterID: actorID,
			SpellID:  "nonexistent_ritual",
		})

		assert.Error(t, err)
	})
}

func TestIsConcentrating(t *testing.T) {
	t.Run("returns false for non-spellcaster", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for concentration",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Fighter",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     16,
					Dexterity:    12,
					Constitution: 14,
					Intelligence: 10,
					Wisdom:       12,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		result, err := e.IsConcentrating(ctx, IsConcentratingRequest{
			GameID:   gameID,
			CasterID: actorID,
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.False(t, result.IsConcentrating)
		assert.Empty(t, result.SpellName)
	})

	t.Run("fails for invalid caster", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for concentration",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		_, err = e.IsConcentrating(ctx, IsConcentratingRequest{
			GameID:   gameID,
			CasterID: model.ID("invalid-caster"),
		})

		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("fails with invalid game", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		_, err := e.IsConcentrating(ctx, IsConcentratingRequest{
			GameID:   model.ID("invalid-game"),
			CasterID: model.ID("caster-id"),
		})

		assert.Error(t, err)
	})

	t.Run("fails for npc", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for concentration",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreateNPC(ctx, CreateNPCRequest{
			GameID: gameID,
			NPC: &NPCInput{
				Name: "Test NPC",
				AbilityScores: AbilityScoresInput{
					Strength:     10,
					Dexterity:    10,
					Constitution: 10,
					Intelligence: 10,
					Wisdom:       10,
					Charisma:     10,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		_, err = e.IsConcentrating(ctx, IsConcentratingRequest{
			GameID:   gameID,
			CasterID: actorID,
		})

		assert.Error(t, err)
	})

	t.Run("fails for enemy", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for concentration",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreateEnemy(ctx, CreateEnemyRequest{
			GameID: gameID,
			Enemy: &EnemyInput{
				Name:            "Goblin",
				AbilityScores:   AbilityScoresInput{Strength: 10, Dexterity: 12, Constitution: 10, Intelligence: 8, Wisdom: 10, Charisma: 8},
				ChallengeRating: 0.25,
				HitPoints:       7,
				ArmorClass:      15,
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		_, err = e.IsConcentrating(ctx, IsConcentratingRequest{
			GameID:   gameID,
			CasterID: actorID,
		})

		assert.Error(t, err)
	})
}

func TestGetConcentrationSpell(t *testing.T) {
	t.Run("fails for invalid caster", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for concentration",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		_, err = e.GetConcentrationSpell(ctx, GetConcentrationSpellRequest{
			GameID:   gameID,
			CasterID: model.ID("invalid-caster"),
		})

		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("fails with invalid game", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		_, err := e.GetConcentrationSpell(ctx, GetConcentrationSpellRequest{
			GameID:   model.ID("invalid-game"),
			CasterID: model.ID("caster-id"),
		})

		assert.Error(t, err)
	})

	t.Run("fails for non-concentrating caster", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for concentration",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Wizard",
				Race:  "Human",
				Class: "Wizard",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     8,
					Dexterity:    12,
					Constitution: 14,
					Intelligence: 18,
					Wisdom:       12,
					Charisma:     10,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		_, err = e.GetConcentrationSpell(ctx, GetConcentrationSpellRequest{
			GameID:   gameID,
			CasterID: actorID,
		})

		assert.Error(t, err)
	})

	t.Run("fails for npc", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for concentration",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreateNPC(ctx, CreateNPCRequest{
			GameID: gameID,
			NPC: &NPCInput{
				Name: "Test NPC",
				AbilityScores: AbilityScoresInput{
					Strength:     10,
					Dexterity:    10,
					Constitution: 10,
					Intelligence: 10,
					Wisdom:       10,
					Charisma:     10,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		_, err = e.GetConcentrationSpell(ctx, GetConcentrationSpellRequest{
			GameID:   gameID,
			CasterID: actorID,
		})

		assert.Error(t, err)
	})

	t.Run("fails for enemy", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for concentration",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreateEnemy(ctx, CreateEnemyRequest{
			GameID: gameID,
			Enemy: &EnemyInput{
				Name:            "Goblin",
				AbilityScores:   AbilityScoresInput{Strength: 10, Dexterity: 12, Constitution: 10, Intelligence: 8, Wisdom: 10, Charisma: 8},
				ChallengeRating: 0.25,
				HitPoints:       7,
				ArmorClass:      15,
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		_, err = e.GetConcentrationSpell(ctx, GetConcentrationSpellRequest{
			GameID:   gameID,
			CasterID: actorID,
		})

		assert.Error(t, err)
	})
}

func TestGetPactMagicSlots(t *testing.T) {
	t.Run("fails for non-spellcaster", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for pact magic",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Fighter",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     16,
					Dexterity:    12,
					Constitution: 14,
					Intelligence: 10,
					Wisdom:       12,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		_, err = e.GetPactMagicSlots(ctx, GetPactMagicSlotsRequest{
			GameID:   gameID,
			CasterID: actorID,
		})

		assert.Error(t, err)
	})

	t.Run("fails for invalid caster", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for pact magic",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		_, err = e.GetPactMagicSlots(ctx, GetPactMagicSlotsRequest{
			GameID:   gameID,
			CasterID: model.ID("invalid-caster"),
		})

		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("fails with invalid game", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		_, err := e.GetPactMagicSlots(ctx, GetPactMagicSlotsRequest{
			GameID:   model.ID("invalid-game"),
			CasterID: model.ID("caster-id"),
		})

		assert.Error(t, err)
	})

	t.Run("fails for npc", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for pact magic",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreateNPC(ctx, CreateNPCRequest{
			GameID: gameID,
			NPC: &NPCInput{
				Name: "Test NPC",
				AbilityScores: AbilityScoresInput{
					Strength:     10,
					Dexterity:    10,
					Constitution: 10,
					Intelligence: 10,
					Wisdom:       10,
					Charisma:     10,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		_, err = e.GetPactMagicSlots(ctx, GetPactMagicSlotsRequest{
			GameID:   gameID,
			CasterID: actorID,
		})

		assert.Error(t, err)
	})

	t.Run("fails for enemy", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for pact magic",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreateEnemy(ctx, CreateEnemyRequest{
			GameID: gameID,
			Enemy: &EnemyInput{
				Name:            "Goblin",
				AbilityScores:   AbilityScoresInput{Strength: 10, Dexterity: 12, Constitution: 10, Intelligence: 8, Wisdom: 10, Charisma: 8},
				ChallengeRating: 0.25,
				HitPoints:       7,
				ArmorClass:      15,
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		_, err = e.GetPactMagicSlots(ctx, GetPactMagicSlotsRequest{
			GameID:   gameID,
			CasterID: actorID,
		})

		assert.Error(t, err)
	})
}

func TestRestorePactMagicSlots(t *testing.T) {
	t.Run("fails for non-spellcaster", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for pact magic",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreatePC(ctx, CreatePCRequest{
			GameID: gameID,
			PC: &PlayerCharacterInput{
				Name:  "Fighter",
				Race:  "Human",
				Class: "Fighter",
				Level: 1,
				AbilityScores: AbilityScoresInput{
					Strength:     16,
					Dexterity:    12,
					Constitution: 14,
					Intelligence: 10,
					Wisdom:       12,
					Charisma:     8,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		err = e.RestorePactMagicSlots(ctx, RestorePactMagicSlotsRequest{
			GameID:   gameID,
			CasterID: actorID,
		})

		assert.Error(t, err)
	})

	t.Run("fails for invalid caster", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for pact magic",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		err = e.RestorePactMagicSlots(ctx, RestorePactMagicSlotsRequest{
			GameID:   gameID,
			CasterID: model.ID("invalid-caster"),
		})

		assert.Error(t, err)
		assert.Equal(t, ErrNotFound, err)
	})

	t.Run("fails with invalid game", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		err := e.RestorePactMagicSlots(ctx, RestorePactMagicSlotsRequest{
			GameID:   model.ID("invalid-game"),
			CasterID: model.ID("caster-id"),
		})

		assert.Error(t, err)
	})

	t.Run("fails for npc", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for pact magic",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreateNPC(ctx, CreateNPCRequest{
			GameID: gameID,
			NPC: &NPCInput{
				Name: "Test NPC",
				AbilityScores: AbilityScoresInput{
					Strength:     10,
					Dexterity:    10,
					Constitution: 10,
					Intelligence: 10,
					Wisdom:       10,
					Charisma:     10,
				},
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		err = e.RestorePactMagicSlots(ctx, RestorePactMagicSlotsRequest{
			GameID:   gameID,
			CasterID: actorID,
		})

		assert.Error(t, err)
	})

	t.Run("fails for enemy", func(t *testing.T) {
		e := NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, NewGameRequest{
			Name:        "Test Game",
			Description: "A test game for pact magic",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		createResult, err := e.CreateEnemy(ctx, CreateEnemyRequest{
			GameID: gameID,
			Enemy: &EnemyInput{
				Name:            "Goblin",
				AbilityScores:   AbilityScoresInput{Strength: 10, Dexterity: 12, Constitution: 10, Intelligence: 8, Wisdom: 10, Charisma: 8},
				ChallengeRating: 0.25,
				HitPoints:       7,
				ArmorClass:      15,
			},
		})
		require.NoError(t, err)
		actorID := createResult.Actor.ID

		err = e.RestorePactMagicSlots(ctx, RestorePactMagicSlotsRequest{
			GameID:   gameID,
			CasterID: actorID,
		})

		assert.Error(t, err)
	})
}
