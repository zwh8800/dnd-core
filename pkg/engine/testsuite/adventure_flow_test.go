package testsuite

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zwh8800/dnd-core/pkg/engine"
	"github.com/zwh8800/dnd-core/pkg/model"
)

func TestFullAdventureFlow(t *testing.T) {
	t.Run("complete adventure from creation to combat", func(t *testing.T) {
		e := engine.NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, engine.NewGameRequest{
			Name:        "The Lost Temple",
			Description: "An adventure to explore an ancient temple",
		})
		require.NoError(t, err)
		require.NotNil(t, gameResult)
		gameID := gameResult.Game.ID

		t.Run("Phase 1: Create player characters", func(t *testing.T) {
			pc1, err := e.CreatePC(ctx, engine.CreatePCRequest{
				GameID: gameID,
				PC: &engine.PlayerCharacterInput{
					Name:       "Aelindra",
					Race:       "Elf",
					Class:      "Ranger",
					Level:      3,
					Background: "Outlander",
					AbilityScores: engine.AbilityScoresInput{
						Strength:     12,
						Dexterity:    16,
						Constitution: 14,
						Intelligence: 10,
						Wisdom:       14,
						Charisma:     10,
					},
				},
			})
			require.NoError(t, err)
			require.NotNil(t, pc1)
			assert.NotEmpty(t, pc1.Actor.ID)

			pc2, err := e.CreatePC(ctx, engine.CreatePCRequest{
				GameID: gameID,
				PC: &engine.PlayerCharacterInput{
					Name:       "Thorin",
					Race:       "Dwarf",
					Class:      "Fighter",
					Level:      3,
					Background: "Soldier",
					AbilityScores: engine.AbilityScoresInput{
						Strength:     16,
						Dexterity:    12,
						Constitution: 15,
						Intelligence: 10,
						Wisdom:       12,
						Charisma:     8,
					},
				},
			})
			require.NoError(t, err)
			require.NotNil(t, pc2)
			assert.NotEmpty(t, pc2.Actor.ID)

			t.Logf("Created PCs: %s (HP: %d), %s (HP: %d)",
				pc1.Actor.Name, pc1.Actor.HitPoints.Current,
				pc2.Actor.Name, pc2.Actor.HitPoints.Current)
		})

		t.Run("Phase 2: Create quest giver NPC and quest", func(t *testing.T) {
			npc, err := e.CreateNPC(ctx, engine.CreateNPCRequest{
				GameID: gameID,
				NPC: &engine.NPCInput{
					Name:        "Elder Theron",
					Description: "A wise old sage from the village",
					Size:        model.SizeMedium,
					Speed:       25,
					AbilityScores: engine.AbilityScoresInput{
						Strength:     10,
						Dexterity:    10,
						Constitution: 12,
						Intelligence: 16,
						Wisdom:       14,
						Charisma:     14,
					},
				},
			})
			require.NoError(t, err)
			require.NotNil(t, npc)

			questResult, err := e.CreateQuest(ctx, engine.CreateQuestRequest{
				GameID:      gameID,
				Name:        "The Lost Temple",
				Description: "Find the ancient temple in the forest and recover the sacred artifact",
				GiverID:     npc.Actor.ID,
				GiverName:   "Elder Theron",
				Objectives: []engine.ObjectiveInput{
					{ID: "find_temple", Description: "Locate the ancient temple", Required: 1},
					{ID: "defeat_guardians", Description: "Defeat the temple guardians", Required: 3},
					{ID: "retrieve_artifact", Description: "Retrieve the sacred artifact", Required: 1},
				},
			})
			require.NoError(t, err)
			require.NotNil(t, questResult)
			assert.NotEmpty(t, questResult.Quest.ID)
			assert.Equal(t, "The Lost Temple", questResult.Quest.Name)

			t.Logf("Created quest: %s", questResult.Quest.Name)
		})

		t.Run("Phase 3: Create exploration scenes", func(t *testing.T) {
			villageScene, err := e.CreateScene(ctx, engine.CreateSceneRequest{
				GameID:      gameID,
				Name:        "Willowbrook Village",
				Description: "A peaceful farming village at the edge of the Darkwood",
				SceneType:   model.SceneTypeOutdoor,
			})
			require.NoError(t, err)
			require.NotNil(t, villageScene)

			forestScene, err := e.CreateScene(ctx, engine.CreateSceneRequest{
				GameID:      gameID,
				Name:        "Darkwood Forest",
				Description: "A dark and mysterious forest",
				SceneType:   model.SceneTypeWilderness,
			})
			require.NoError(t, err)
			require.NotNil(t, forestScene)

			templeScene, err := e.CreateScene(ctx, engine.CreateSceneRequest{
				GameID:      gameID,
				Name:        "Temple of the Ancients",
				Description: "An ancient temple covered in vines",
				SceneType:   model.SceneTypeDungeon,
			})
			require.NoError(t, err)
			require.NotNil(t, templeScene)

			err = e.AddSceneConnection(ctx, engine.AddSceneConnectionRequest{
				GameID:        gameID,
				SceneID:       villageScene.Scene.ID,
				TargetSceneID: forestScene.Scene.ID,
				Description:   "A dirt path leading into the forest",
				Locked:        false,
				DC:            0,
				Hidden:        false,
			})
			require.NoError(t, err)

			err = e.AddSceneConnection(ctx, engine.AddSceneConnectionRequest{
				GameID:        gameID,
				SceneID:       forestScene.Scene.ID,
				TargetSceneID: templeScene.Scene.ID,
				Description:   "Ancient stone steps leading underground",
				Locked:        true,
				DC:            15,
				Hidden:        true,
			})
			require.NoError(t, err)

			t.Logf("Created scenes with connections")
		})

		t.Run("Phase 4: Add enemy monsters", func(t *testing.T) {
			skeleton1, err := e.CreateEnemy(ctx, engine.CreateEnemyRequest{
				GameID: gameID,
				Enemy: &engine.EnemyInput{
					Name:        "Skeleton Warrior",
					Description: "An animated skeleton wielding a rusty sword",
					Size:        model.SizeMedium,
					Speed:       30,
					AbilityScores: engine.AbilityScoresInput{
						Strength:     14,
						Dexterity:    12,
						Constitution: 10,
						Intelligence: 6,
						Wisdom:       8,
						Charisma:     5,
					},
					ChallengeRating: 0.25,
					HitPoints:       13,
					ArmorClass:      13,
				},
			})
			require.NoError(t, err)
			require.NotNil(t, skeleton1)
			assert.NotEmpty(t, skeleton1.Actor.ID)

			skeleton2, err := e.CreateEnemy(ctx, engine.CreateEnemyRequest{
				GameID: gameID,
				Enemy: &engine.EnemyInput{
					Name:        "Skeleton Archer",
					Description: "An animated skeleton with a bow",
					Size:        model.SizeMedium,
					Speed:       30,
					AbilityScores: engine.AbilityScoresInput{
						Strength:     10,
						Dexterity:    14,
						Constitution: 10,
						Intelligence: 6,
						Wisdom:       8,
						Charisma:     5,
					},
					ChallengeRating: 0.25,
					HitPoints:       10,
					ArmorClass:      13,
				},
			})
			require.NoError(t, err)
			require.NotNil(t, skeleton2)
			assert.NotEmpty(t, skeleton2.Actor.ID)

			t.Logf("Created enemies: %s (AC: %d, HP: %d), %s (AC: %d, HP: %d)",
				skeleton1.Actor.Name, skeleton1.Actor.ArmorClass, skeleton1.Actor.HitPoints.Current,
				skeleton2.Actor.Name, skeleton2.Actor.ArmorClass, skeleton2.Actor.HitPoints.Current)
		})

		t.Run("Phase 5: List actors for combat", func(t *testing.T) {
			actorsResult, err := e.ListActors(ctx, engine.ListActorsRequest{
				GameID: gameID,
				Filter: &engine.ActorFilter{
					Types: []model.ActorType{model.ActorTypePC},
				},
			})
			require.NoError(t, err)
			require.GreaterOrEqual(t, len(actorsResult.Actors), 2)
			t.Logf("Found %d player characters", len(actorsResult.Actors))

			enemiesResult, err := e.ListActors(ctx, engine.ListActorsRequest{
				GameID: gameID,
				Filter: &engine.ActorFilter{
					Types: []model.ActorType{model.ActorTypeEnemy},
				},
			})
			require.NoError(t, err)
			require.GreaterOrEqual(t, len(enemiesResult.Actors), 2)
			t.Logf("Found %d enemies", len(enemiesResult.Actors))
		})

		t.Run("Phase 6: Transition to exploration phase", func(t *testing.T) {
			result, err := e.SetPhase(ctx, gameID, model.PhaseExploration, "Begin exploration of the forest")
			require.NoError(t, err)
			require.NotNil(t, result)

			phase, err := e.GetPhase(ctx, gameID)
			require.NoError(t, err)
			assert.Equal(t, model.PhaseExploration, phase)
			t.Logf("Phase set to: %s", phase)
		})

		t.Run("Phase 7: Start combat with surprise", func(t *testing.T) {
			actorsResult, err := e.ListActors(ctx, engine.ListActorsRequest{
				GameID: gameID,
				Filter: &engine.ActorFilter{
					Types: []model.ActorType{model.ActorTypePC, model.ActorTypeEnemy},
				},
			})
			require.NoError(t, err)
			require.GreaterOrEqual(t, len(actorsResult.Actors), 4)

			var participantIDs []model.ID
			var stealthyIDs []model.ID
			for i, actor := range actorsResult.Actors {
				participantIDs = append(participantIDs, actor.ID)
				if i < 2 {
					stealthyIDs = append(stealthyIDs, actor.ID)
				}
			}

			surpriseResult, err := e.StartCombatWithSurprise(ctx, engine.StartCombatWithSurpriseRequest{
				GameID:       gameID,
				SceneID:      actorsResult.Actors[0].SceneID,
				StealthySide: stealthyIDs,
				Observers:    participantIDs[len(stealthyIDs):],
			})
			require.NoError(t, err)
			require.NotNil(t, surpriseResult)
			assert.Equal(t, model.CombatStatusActive, surpriseResult.Combat.Status)

			t.Logf("Combat started with %d participants (surprise possible)",
				len(participantIDs))
		})
	})

	t.Run("adventure with skill checks", func(t *testing.T) {
		e := engine.NewTestEngine(t)
		ctx := context.Background()

		gameResult, err := e.NewGame(ctx, engine.NewGameRequest{
			Name:        "Diplomatic Mission",
			Description: "A social and exploration focused adventure",
		})
		require.NoError(t, err)
		gameID := gameResult.Game.ID

		pc, err := e.CreatePC(ctx, engine.CreatePCRequest{
			GameID: gameID,
			PC: &engine.PlayerCharacterInput{
				Name:       "Elara",
				Race:       "Half-Elf",
				Class:      "Bard",
				Level:      2,
				Background: "Entertainer",
				AbilityScores: engine.AbilityScoresInput{
					Strength:     10,
					Dexterity:    14,
					Constitution: 10,
					Intelligence: 12,
					Wisdom:       10,
					Charisma:     16,
				},
			},
		})
		require.NoError(t, err)
		require.NotNil(t, pc)

		_, err = e.SetPhase(ctx, gameID, model.PhaseExploration, "Start adventure")
		require.NoError(t, err)

		t.Run("perform ability check during exploration", func(t *testing.T) {
			checkResult, err := e.PerformAbilityCheck(ctx, engine.AbilityCheckRequest{
				GameID:  gameID,
				ActorID: pc.Actor.ID,
				Ability: model.AbilityCharisma,
				DC:      12,
				Reason:  "Persuade the guard",
			})
			require.NoError(t, err)
			require.NotNil(t, checkResult)
			assert.Equal(t, model.AbilityCharisma, checkResult.Ability)
			t.Logf("Charisma check: Roll=%d, Total=%d, DC=%d, Success=%v",
				checkResult.Roll.Total, checkResult.RollTotal, checkResult.DC, checkResult.Success)
		})

		t.Run("perform skill check", func(t *testing.T) {
			checkResult, err := e.PerformSkillCheck(ctx, engine.SkillCheckRequest{
				GameID:  gameID,
				ActorID: pc.Actor.ID,
				Skill:   model.SkillPersuasion,
				DC:      14,
				Reason:  "Charm the merchant",
			})
			require.NoError(t, err)
			require.NotNil(t, checkResult)
			assert.Equal(t, model.SkillPersuasion, checkResult.Skill)
			t.Logf("Persuasion check: Roll=%d, Total=%d, Success=%v",
				checkResult.Roll.Total, checkResult.RollTotal, checkResult.Success)
		})

		t.Run("perform saving throw", func(t *testing.T) {
			saveResult, err := e.PerformSavingThrow(ctx, engine.SavingThrowRequest{
				GameID:  gameID,
				ActorID: pc.Actor.ID,
				Ability: model.AbilityCharisma,
				DC:      10,
				Reason:  "Resist charm effect",
			})
			require.NoError(t, err)
			require.NotNil(t, saveResult)
			assert.Equal(t, model.AbilityCharisma, saveResult.Ability)
			t.Logf("Charisma saving throw: Roll=%d, Total=%d, Success=%v",
				saveResult.Roll.Total, saveResult.RollTotal, saveResult.Success)
		})
	})
}
