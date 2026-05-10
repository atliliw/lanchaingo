package tools

import (
	"context"
	"testing"
)

// 娴嬭瘯 Calculator 鍔犳硶
func TestCalculatorAdd(t *testing.T) {
	c := NewCalculator()
	result, err := c.Run(context.Background(), "2 + 3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "2 + 3 = 5" {
		t.Errorf("expected '2 + 3 = 5', got '%s'", result)
	}
}

// 娴嬭瘯 Calculator 鍑忔硶
func TestCalculatorSubtract(t *testing.T) {
	c := NewCalculator()
	result, err := c.Run(context.Background(), "10 - 3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "10 - 3 = 7" {
		t.Errorf("unexpected: %s", result)
	}
}

// 娴嬭瘯 Calculator 涔樻硶
func TestCalculatorMultiply(t *testing.T) {
	c := NewCalculator()
	result, err := c.Run(context.Background(), "4 * 5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "4 * 5 = 20" {
		t.Errorf("unexpected: %s", result)
	}
}

// 娴嬭瘯 Calculator 闄ゆ硶
func TestCalculatorDivide(t *testing.T) {
	c := NewCalculator()
	result, err := c.Run(context.Background(), "100 / 4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "100 / 4 = 25" {
		t.Errorf("unexpected: %s", result)
	}
}

// 娴嬭瘯 Calculator 闄ら浂閿欒
func TestCalculatorDivideByZero(t *testing.T) {
	c := NewCalculator()
	_, err := c.Run(context.Background(), "1 / 0")
	if err == nil {
		t.Error("expected error for division by zero")
	}
}

// 娴嬭瘯 Calculator 鏃犳晥琛ㄨ揪寮?func TestCalculatorInvalid(t *testing.T) {
	c := NewCalculator()
	_, err := c.Run(context.Background(), "invalid")
	if err == nil {
		t.Error("expected error for invalid expression")
	}
}

// 娴嬭瘯 Calculator 绾暟瀛?func TestCalculatorNumber(t *testing.T) {
	c := NewCalculator()
	result, err := c.Run(context.Background(), "42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "42 = 42" {
		t.Errorf("unexpected: %s", result)
	}
}

// 娴嬭瘯 Name 鍜?Description
func TestCalculatorMeta(t *testing.T) {
	c := NewCalculator()
	if c.Name() != "calculator" {
		t.Errorf("expected calculator, got %s", c.Name())
	}
	if c.Description() == "" {
		t.Error("expected non-empty description")
	}
}
