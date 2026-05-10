package tools

import (
	"context"
	"testing"
)

func TestCalculatorAdd(t *testing.T) {
	c := NewCalculator()
	r, _ := c.Run(context.Background(), "2 + 3")
	if r != "2 + 3 = 5" { t.Errorf("got %s", r) }
}
func TestCalculatorSubtract(t *testing.T) {
	c := NewCalculator()
	r, _ := c.Run(context.Background(), "10 - 3")
	if r != "10 - 3 = 7" { t.Errorf("got %s", r) }
}
func TestCalculatorMultiply(t *testing.T) {
	c := NewCalculator()
	r, _ := c.Run(context.Background(), "4 * 5")
	if r != "4 * 5 = 20" { t.Errorf("got %s", r) }
}
func TestCalculatorDivide(t *testing.T) {
	c := NewCalculator()
	r, _ := c.Run(context.Background(), "100 / 4")
	if r != "100 / 4 = 25" { t.Errorf("got %s", r) }
}
func TestCalculatorDivideByZero(t *testing.T) {
	c := NewCalculator()
	_, err := c.Run(context.Background(), "1 / 0")
	if err == nil { t.Error("expected error") }
}
func TestCalculatorNumber(t *testing.T) {
	c := NewCalculator()
	r, _ := c.Run(context.Background(), "42")
	if r != "42 = 42" { t.Errorf("got %s", r) }
}
func TestCalculatorMeta(t *testing.T) {
	c := NewCalculator()
	if c.Name() != "calculator" { t.Errorf("expected calculator") }
}
func TestDateTimeTool(t *testing.T) {
	d := NewDateTimeTool()
	r, _ := d.Run(context.Background(), "")
	if r == "" { t.Error("expected non-empty datetime") }
}
func TestSimpleMathSqrt(t *testing.T) {
	m := NewSimpleMathTool()
	r, _ := m.Run(context.Background(), "sqrt(16)")
	if r != "sqrt(16) = 4" { t.Errorf("got %s", r) }
}
func TestSimpleMathPow(t *testing.T) {
	m := NewSimpleMathTool()
	r, _ := m.Run(context.Background(), "pow(2,10)")
	if r != "pow(2,10) = 1024" { t.Errorf("got %s", r) }
}
func TestSimpleMathAbs(t *testing.T) {
	m := NewSimpleMathTool()
	r, _ := m.Run(context.Background(), "abs(-5)")
	if r != "abs(-5) = 5" { t.Errorf("got %s", r) }
}
func TestSearchMeta(t *testing.T) {
	s := NewDuckDuckGoSearchTool()
	if s.Name() != "search" { t.Errorf("expected search") }
}
func TestURLFetchMeta(t *testing.T) {
	u := NewURLFetchTool()
	if u.Name() != "url_fetch" { t.Errorf("expected url_fetch") }
}
func TestWikipediaMeta(t *testing.T) {
	w := NewWikipediaTool()
	if w.Name() != "wikipedia" { t.Errorf("expected wikipedia") }
}
func TestPythonREPLMeta(t *testing.T) {
	p := NewPythonREPLTool()
	if p.Name() != "python_repl" { t.Errorf("expected python_repl") }
}
