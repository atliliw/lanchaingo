package output_parsers

import "testing"

func TestStrOutputParser(t *testing.T) {
	p := NewStrOutputParser()
	r, _ := p.Parse("hello")
	if r != "hello" { t.Errorf("expected hello") }
}
func TestStrOutputParserEmpty(t *testing.T) {
	p := NewStrOutputParser()
	r, _ := p.Parse("")
	if r != "" { t.Errorf("expected empty") }
}
func TestStrOutputParserType(t *testing.T) {
	p := NewStrOutputParser()
	if p.Type() != "string" { t.Errorf("expected string") }
}
func TestJsonOutputParser(t *testing.T) {
	p := NewJsonOutputParser()
	r, err := p.Parse(`{"key":"value"}`)
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	m := r.(map[string]any)
	if m["key"] != "value" { t.Errorf("expected value") }
}
func TestJsonOutputParserInvalid(t *testing.T) {
	p := NewJsonOutputParser()
	_, err := p.Parse("not json")
	if err == nil { t.Error("expected error") }
}
func TestCommaSeparatedListOutputParser(t *testing.T) {
	p := NewCommaSeparatedListOutputParser()
	r, _ := p.Parse("a, b, c")
	list := r.([]string)
	if len(list) != 3 { t.Errorf("expected 3") }
}
func TestCommaSeparatedListEmpty(t *testing.T) {
	p := NewCommaSeparatedListOutputParser()
	r, _ := p.Parse("")
	list := r.([]string)
	if len(list) != 0 { t.Errorf("expected 0") }
}
func TestStructuredOutputParser(t *testing.T) {
	s := map[string]any{"name": map[string]any{"type": "string"}}
	p := NewStructuredOutputParser(s)
	r, _ := p.Parse(`{"name":"Alice"}`)
	m := r.(map[string]any)
	if m["name"] != "Alice" { t.Errorf("expected Alice") }
}
func TestStructuredOutputParserMissingFields(t *testing.T) {
	s := map[string]any{"req": map[string]any{"type": "string"}}
	p := NewStructuredOutputParser(s)
	_, err := p.Parse(`{"other": "value"}`)
	if err == nil { t.Error("expected error") }
}
func TestTypedOutputParser(t *testing.T) {
	type P struct{ Name string }
	p := NewTypedOutputParser(P{})
	r, _ := p.Parse(`{"name":"Alice"}`)
	person := r.(P)
	if person.Name != "Alice" { t.Errorf("expected Alice") }
}
func TestValidateJSON(t *testing.T) {
	if err := ValidateJSON(`{"a":1}`); err != nil { t.Errorf("unexpected error: %v", err) }
	if err := ValidateJSON("not json"); err == nil { t.Error("expected error") }
}
