package engine

import (
	"context"
	"testing"

	"github.com/zwh8800/dnd-core/internal/model"
)

// TestPerformAbilityCheck 测试属性检定
func TestPerformAbilityCheck(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	// 创建PC
	pc := &model.PlayerCharacter{
		Actor: model.Actor{
			Name:  "Hero",
			Size:  model.SizeMedium,
			Speed: 30,
			AbilityScores: model.AbilityScores{
				Strength: 16, Dexterity: 12, Constitution: 14,
				Intelligence: 10, Wisdom: 8, Charisma: 13,
			},
		},
		Race: model.RaceReference{Name: "Human"},
		Classes: []model.ClassLevel{
			{ClassName: "Fighter", Level: 1},
		},
		TotalLevel: 1,
	}
	pcResult, err := engine.CreatePC(ctx, gameID, pc)
	if err != nil {
		t.Fatalf("Failed to create PC: %v", err)
	}

	// 执行力量检定
	result, err := engine.PerformAbilityCheck(ctx, gameID, AbilityCheckRequest{
		ActorID: pcResult.ID,
		Ability: model.AbilityStrength,
		DC:      15,
		Reason:  "Test strength check",
	})
	if err != nil {
		t.Fatalf("Failed to perform ability check: %v", err)
	}

	if result.ActorID != pcResult.ID {
		t.Errorf("Expected actor ID %s, got %s", pcResult.ID, result.ActorID)
	}
	if result.Ability != model.AbilityStrength {
		t.Errorf("Expected ability %s, got %s", model.AbilityStrength, result.Ability)
	}
	if result.AbilityScore != 16 {
		t.Errorf("Expected ability score 16, got %d", result.AbilityScore)
	}
	if result.AbilityMod != 3 {
		t.Errorf("Expected ability modifier 3, got %d", result.AbilityMod)
	}
}

// TestPerformSkillCheck 测试技能检定
func TestPerformSkillCheck(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	// 创建PC
	pc := &model.PlayerCharacter{
		Actor: model.Actor{
			Name:  "Rogue",
			Size:  model.SizeMedium,
			Speed: 30,
			AbilityScores: model.AbilityScores{
				Strength: 10, Dexterity: 16, Constitution: 12,
				Intelligence: 14, Wisdom: 10, Charisma: 8,
			},
		},
		Race: model.RaceReference{Name: "Elf"},
		Classes: []model.ClassLevel{
			{ClassName: "Rogue", Level: 1},
		},
		TotalLevel: 1,
	}
	pcResult, err := engine.CreatePC(ctx, gameID, pc)
	if err != nil {
		t.Fatalf("Failed to create PC: %v", err)
	}

	// 执行隐匿检定（基于敏捷）
	result, err := engine.PerformSkillCheck(ctx, gameID, SkillCheckRequest{
		ActorID: pcResult.ID,
		Skill:   model.SkillStealth,
		DC:      15,
		Reason:  "Test stealth check",
	})
	if err != nil {
		t.Fatalf("Failed to perform skill check: %v", err)
	}

	if result.Skill != model.SkillStealth {
		t.Errorf("Expected skill %s, got %s", model.SkillStealth, result.Skill)
	}
	if result.Ability != model.AbilityDexterity {
		t.Errorf("Expected ability %s, got %s", model.AbilityDexterity, result.Ability)
	}
}

// TestPerformSavingThrow 测试豁免检定
func TestPerformSavingThrow(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	// 创建PC
	pc := &model.PlayerCharacter{
		Actor: model.Actor{
			Name:  "Hero",
			Size:  model.SizeMedium,
			Speed: 30,
			AbilityScores: model.AbilityScores{
				Strength: 16, Dexterity: 12, Constitution: 14,
				Intelligence: 10, Wisdom: 8, Charisma: 13,
			},
		},
		Race: model.RaceReference{Name: "Human"},
		Classes: []model.ClassLevel{
			{ClassName: "Fighter", Level: 1},
		},
		TotalLevel: 1,
	}
	pcResult, err := engine.CreatePC(ctx, gameID, pc)
	if err != nil {
		t.Fatalf("Failed to create PC: %v", err)
	}

	// 执行体质豁免
	result, err := engine.PerformSavingThrow(ctx, gameID, SavingThrowRequest{
		ActorID: pcResult.ID,
		Ability: model.AbilityConstitution,
		DC:      15,
		Reason:  "Test constitution save",
	})
	if err != nil {
		t.Fatalf("Failed to perform saving throw: %v", err)
	}

	if result.ActorID != pcResult.ID {
		t.Errorf("Expected actor ID %s, got %s", pcResult.ID, result.ActorID)
	}
	if result.Ability != model.AbilityConstitution {
		t.Errorf("Expected ability %s, got %s", model.AbilityConstitution, result.Ability)
	}
	if result.DC != 15 {
		t.Errorf("Expected DC 15, got %d", result.DC)
	}
}

// TestPerformAbilityCheckNotFound 测试找不到角色的属性检定
func TestPerformAbilityCheckNotFound(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()
	fakeID := model.ID("00000000000000000000000000")

	_, err := engine.PerformAbilityCheck(ctx, gameID, AbilityCheckRequest{
		ActorID: fakeID,
		Ability: model.AbilityStrength,
		DC:      10,
	})
	if err == nil {
		t.Fatal("Expected error for non-existent actor, got nil")
	}
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

// TestPerformSkillCheckNotFound 测试找不到角色的技能检定
func TestPerformSkillCheckNotFound(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()
	fakeID := model.ID("00000000000000000000000000")

	_, err := engine.PerformSkillCheck(ctx, gameID, SkillCheckRequest{
		ActorID: fakeID,
		Skill:   model.SkillStealth,
		DC:      10,
	})
	if err == nil {
		t.Fatal("Expected error for non-existent actor, got nil")
	}
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

// TestPerformSavingThrowNotFound 测试找不到角色的豁免检定
func TestPerformSavingThrowNotFound(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()
	fakeID := model.ID("00000000000000000000000000")

	_, err := engine.PerformSavingThrow(ctx, gameID, SavingThrowRequest{
		ActorID: fakeID,
		Ability: model.AbilityConstitution,
		DC:      10,
	})
	if err == nil {
		t.Fatal("Expected error for non-existent actor, got nil")
	}
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

// TestGetPassivePerception 测试获取被动察觉
func TestGetPassivePerception(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	// 创建PC
	pc := &model.PlayerCharacter{
		Actor: model.Actor{
			Name:  "Ranger",
			Size:  model.SizeMedium,
			Speed: 30,
			AbilityScores: model.AbilityScores{
				Strength: 12, Dexterity: 14, Constitution: 12,
				Intelligence: 10, Wisdom: 16, Charisma: 8,
			},
		},
		Race: model.RaceReference{Name: "Elf"},
		Classes: []model.ClassLevel{
			{ClassName: "Ranger", Level: 1},
		},
		TotalLevel: 1,
	}
	pcResult, err := engine.CreatePC(ctx, gameID, pc)
	if err != nil {
		t.Fatalf("Failed to create PC: %v", err)
	}

	// 获取被动察觉
	passive, err := engine.GetPassivePerception(ctx, gameID, pcResult.ID)
	if err != nil {
		t.Fatalf("Failed to get passive perception: %v", err)
	}

	// 基础被动察觉 = 10 + 感知修正
	expectedBase := 10 + 3 // 16 WIS = +3
	if passive < expectedBase {
		t.Errorf("Expected passive perception at least %d, got %d", expectedBase, passive)
	}
}

// TestGetPassivePerceptionNotFound 测试找不到角色的被动察觉
func TestGetPassivePerceptionNotFound(t *testing.T) {
	engine, gameID := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()
	fakeID := model.ID("00000000000000000000000000")

	_, err := engine.GetPassivePerception(ctx, gameID, fakeID)
	if err == nil {
		t.Fatal("Expected error for non-existent actor, got nil")
	}
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

// TestGetSkillAbility 测试获取技能对应的属性
func TestGetSkillAbility(t *testing.T) {
	engine, _ := createTestGame(t)
	defer engine.Close()

	ability := engine.GetSkillAbility(model.SkillStealth)
	if ability != model.AbilityDexterity {
		t.Errorf("Expected Stealth to use Dexterity, got %s", ability)
	}

	ability = engine.GetSkillAbility(model.SkillPerception)
	if ability != model.AbilityWisdom {
		t.Errorf("Expected Perception to use Wisdom, got %s", ability)
	}

	ability = engine.GetSkillAbility(model.SkillArcana)
	if ability != model.AbilityIntelligence {
		t.Errorf("Expected Arcana to use Intelligence, got %s", ability)
	}
}
