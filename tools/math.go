package tools

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// SimpleMathTool handles advanced math operations.
type SimpleMathTool struct{}

func NewSimpleMathTool() *SimpleMathTool { return &SimpleMathTool{} }

func (m *SimpleMathTool) Name() string { return "simple_math" }

func (m *SimpleMathTool) Description() string {
	return `楂樼骇鏁板璁＄畻銆傛敮鎸佺殑鍑芥暟: sqrt, pow, sin, cos, abs銆?绀轰緥:
- 'sqrt(16)' 鈫?4
- 'pow(2,10)' 鈫?1024
- 'abs(-5)' 鈫?5`
}

func (m *SimpleMathTool) Run(_ context.Context, input string) (string, error) {
	input = strings.TrimSpace(input)

	if strings.HasPrefix(input, "sqrt(") && strings.HasSuffix(input, ")") {
		inner := strings.TrimSpace(input[5 : len(input)-1])
		v, err := strconv.ParseFloat(inner, 64)
		if err != nil {
			return "", fmt.Errorf("invalid input: %s", inner)
		}
		return fmt.Sprintf("sqrt(%v) = %v", inner, math.Sqrt(v)), nil
	}

	if strings.HasPrefix(input, "pow(") && strings.HasSuffix(input, ")") {
		inner := strings.TrimSpace(input[4 : len(input)-1])
		parts := strings.Split(inner, ",")
		if len(parts) == 2 {
			a, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			b, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			if err1 == nil && err2 == nil {
				return fmt.Sprintf("pow(%v,%v) = %v", a, b, math.Pow(a, b)), nil
			}
		}
	}

	if strings.HasPrefix(input, "abs(") && strings.HasSuffix(input, ")") {
		inner := strings.TrimSpace(input[4 : len(input)-1])
		v, err := strconv.ParseFloat(inner, 64)
		if err != nil {
			return "", fmt.Errorf("invalid input: %s", inner)
		}
		return fmt.Sprintf("abs(%v) = %v", inner, math.Abs(v)), nil
	}

	return "", fmt.Errorf("cannot evaluate: %s", input)
}
