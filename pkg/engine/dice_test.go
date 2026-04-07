package engine

import (
	"context"
	"testing"
)

// TestRoll 测试基本掷骰
func TestRoll(t *testing.T) {
	engine, _ := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	result, err := engine.Roll(ctx, RollRequest{
		Expression: "1d20",
		Modifier:   0,
		Reason:     "Test roll",
	})
	if err != nil {
		t.Fatalf("Failed to roll: %v", err)
	}

	if result.Total < 1 || result.Total > 20 {
		t.Errorf("Expected roll result between 1 and 20, got %d", result.Total)
	}
	if result.Expression != "1d20" {
		t.Errorf("Expected expression '1d20', got %s", result.Expression)
	}
	if result.Message == "" {
		t.Error("Expected message to be set")
	}
}

// TestRollWithModifier 测试带修正的掷骰
func TestRollWithModifier(t *testing.T) {
	engine, _ := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	result, err := engine.Roll(ctx, RollRequest{
		Expression: "1d20",
		Modifier:   5,
		Reason:     "Test roll with modifier",
	})
	if err != nil {
		t.Fatalf("Failed to roll: %v", err)
	}

	// 结果应该在6-25之间（1-20 + 5）
	if result.Total < 6 || result.Total > 25 {
		t.Errorf("Expected roll result between 6 and 25, got %d", result.Total)
	}
	if result.Modifier != 5 {
		t.Errorf("Expected modifier 5, got %d", result.Modifier)
	}
}

// TestRollMultipleDice 测试多骰子掷骰
func TestRollMultipleDice(t *testing.T) {
	engine, _ := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	result, err := engine.Roll(ctx, RollRequest{
		Expression: "2d6",
		Modifier:   0,
		Reason:     "Test multiple dice",
	})
	if err != nil {
		t.Fatalf("Failed to roll: %v", err)
	}

	// 结果应该在2-12之间
	if result.Total < 2 || result.Total > 12 {
		t.Errorf("Expected roll result between 2 and 12, got %d", result.Total)
	}
	if len(result.Rolls) != 2 {
		t.Errorf("Expected 2 rolls, got %d", len(result.Rolls))
	}
}

// TestRollInvalidExpression 测试无效表达式
func TestRollInvalidExpression(t *testing.T) {
	engine, _ := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	_, err := engine.Roll(ctx, RollRequest{
		Expression: "invalid",
		Modifier:   0,
		Reason:     "Test invalid expression",
	})
	if err == nil {
		t.Fatal("Expected error for invalid expression, got nil")
	}
}

// TestRollAdvantage 测试优势掷骰
func TestRollAdvantage(t *testing.T) {
	engine, _ := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	result, err := engine.RollAdvantage(ctx, RollAdvantageRequest{
		Modifier: 0,
		Reason:   "Test advantage",
	})
	if err != nil {
		t.Fatalf("Failed to roll advantage: %v", err)
	}

	// 优势掷骰应该在1-20之间
	if result.Total < 1 || result.Total > 20 {
		t.Errorf("Expected advantage roll between 1 and 20, got %d", result.Total)
	}
}

// TestRollAdvantageWithModifier 测试带修正的优势掷骰
func TestRollAdvantageWithModifier(t *testing.T) {
	engine, _ := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	result, err := engine.RollAdvantage(ctx, RollAdvantageRequest{
		Modifier: 3,
		Reason:   "Test advantage with modifier",
	})
	if err != nil {
		t.Fatalf("Failed to roll advantage: %v", err)
	}

	// 结果应该在4-23之间（1-20 + 3）
	if result.Total < 4 || result.Total > 23 {
		t.Errorf("Expected advantage roll between 4 and 23, got %d", result.Total)
	}
}

// TestRollDisadvantage 测试劣势掷骰
func TestRollDisadvantage(t *testing.T) {
	engine, _ := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	result, err := engine.RollDisadvantage(ctx, RollDisadvantageRequest{
		Modifier: 0,
		Reason:   "Test disadvantage",
	})
	if err != nil {
		t.Fatalf("Failed to roll disadvantage: %v", err)
	}

	// 劣势掷骰应该在1-20之间
	if result.Total < 1 || result.Total > 20 {
		t.Errorf("Expected disadvantage roll between 1 and 20, got %d", result.Total)
	}
}

// TestRollAbility 测试属性掷骰（4d6去掉最低）
func TestRollAbility(t *testing.T) {
	engine, _ := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	result, err := engine.RollAbility(ctx)
	if err != nil {
		t.Fatalf("Failed to roll ability: %v", err)
	}

	// 属性掷骰结果应该在3-18之间（3d6）
	if result.Total < 3 || result.Total > 18 {
		t.Errorf("Expected ability roll between 3 and 18, got %d", result.Total)
	}

	// 应该有4个掷骰，其中1个被丢弃
	if len(result.Rolls) != 4 {
		t.Errorf("Expected 4 rolls, got %d", len(result.Rolls))
	}

	droppedCount := 0
	for _, roll := range result.Rolls {
		if roll.Dropped {
			droppedCount++
		}
	}
	if droppedCount != 1 {
		t.Errorf("Expected 1 dropped roll, got %d", droppedCount)
	}
}

// TestRollHitDice 测试生命骰掷骰
func TestRollHitDice(t *testing.T) {
	engine, _ := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	result, err := engine.RollHitDice(ctx, 8, 2)
	if err != nil {
		t.Fatalf("Failed to roll hit dice: %v", err)
	}

	// 结果应该在3-10之间（1d8 + 2）
	if result.Total < 3 || result.Total > 10 {
		t.Errorf("Expected hit dice roll between 3 and 10, got %d", result.Total)
	}
	if result.Modifier != 2 {
		t.Errorf("Expected modifier 2, got %d", result.Modifier)
	}
}

// TestRollHitDiceNoModifier 测试无修正的生命骰掷骰
func TestRollHitDiceNoModifier(t *testing.T) {
	engine, _ := createTestGame(t)
	defer engine.Close()

	ctx := context.Background()

	result, err := engine.RollHitDice(ctx, 6, 0)
	if err != nil {
		t.Fatalf("Failed to roll hit dice: %v", err)
	}

	// 结果应该在1-6之间
	if result.Total < 1 || result.Total > 6 {
		t.Errorf("Expected hit dice roll between 1 and 6, got %d", result.Total)
	}
}
