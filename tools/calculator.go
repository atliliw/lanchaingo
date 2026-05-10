package tools

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

// Calculator evaluates simple math expressions.
type Calculator struct{}

func NewCalculator() *Calculator { return &Calculator{} }

func (c *Calculator) Name() string { return "calculator" }

func (c *Calculator) Description() string {
	return `璁＄畻鏁板琛ㄨ揪寮忋€傛敮鎸佸熀鏈繍绠楋紙鍔犲噺涔橀櫎锛夈€?绀轰緥:
- '2 + 2' 鈫?4
- '10 - 3' 鈫?7
- '3.14 * 10' 鈫?31.4
- '100 / 4' 鈫?25`
}

func (c *Calculator) Run(_ context.Context, input string) (string, error) {
	expr := strings.TrimSpace(input)

	if v, err := strconv.ParseFloat(expr, 64); err == nil {
		return fmt.Sprintf("%v = %v", expr, v), nil
	}

	if strings.Contains(expr, "+") {
		parts := strings.SplitN(expr, "+", 2)
		if len(parts) == 2 {
			a, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			b, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			if err1 == nil && err2 == nil {
				return fmt.Sprintf("%v = %v", expr, a+b), nil
			}
		}
	}
	if strings.Contains(expr, "-") {
		parts := strings.SplitN(expr, "-", 2)
		if len(parts) == 2 {
			a, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			b, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			if err1 == nil && err2 == nil {
				return fmt.Sprintf("%v = %v", expr, a-b), nil
			}
		}
	}
	if strings.Contains(expr, "*") {
		parts := strings.SplitN(expr, "*", 2)
		if len(parts) == 2 {
			a, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			b, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			if err1 == nil && err2 == nil {
				return fmt.Sprintf("%v = %v", expr, a*b), nil
			}
		}
	}
	if strings.Contains(expr, "/") {
		parts := strings.SplitN(expr, "/", 2)
		if len(parts) == 2 {
			a, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			b, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			if err1 == nil && err2 == nil {
				if b == 0 {
					return "", fmt.Errorf("division by zero")
				}
				return fmt.Sprintf("%v = %v", expr, a/b), nil
			}
		}
	}

	return "", fmt.Errorf("cannot evaluate expression: %s", expr)
}
