package tools

import (
	"context"
	"testing"
)

// 娴嬭瘯 SimpleMathTool sqrt
func TestSimpleMathSqrt(t *testing.T) {
	m := NewSimpleMathTool()
	result, err := m.Run(context.Background(), "sqrt(16)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "sqrt(16) = 4" {
		t.Errorf("unexpected: %s", result)
	}
}

// 娴嬭瘯 SimpleMathTool pow
func TestSimpleMathPow(t *testing.T) {
	m := NewSimpleMathTool()
	result, err := m.Run(context.Background(), "pow(2,10)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "pow(2,10) = 1024" {
		t.Errorf("unexpected: %s", result)
	}
}

// 娴嬭瘯 SimpleMathTool abs
func TestSimpleMathAbs(t *testing.T) {
	m := NewSimpleMathTool()
	result, err := m.Run(context.Background(), "abs(-5)")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "abs(-5) = 5" {
		t.Errorf("unexpected: %s", result)
	}
}

// 娴嬭瘯 SimpleMathTool 鏃犳晥杈撳叆
func TestSimpleMathInvalid(t *testing.T) {
	m := NewSimpleMathTool()
	_, err := m.Run(context.Background(), "unknown(x)")
	if err == nil {
		t.Error("expected error for unknown function")
	}
}

// 娴嬭瘯 SimpleMathTool 鍏冩暟鎹?func TestSimpleMathMeta(t *testing.T) {
	m := NewSimpleMathTool()
	if m.Name() != "simple_math" {
		t.Errorf("expected simple_math, got %s", m.Name())
	}
}
