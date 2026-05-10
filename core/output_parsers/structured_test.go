package output_parsers

import "testing"

// 娴嬭瘯 StructuredOutputParser锛氭寜 schema 瑙ｆ瀽 JSON锛岄獙璇佸瓧娈靛€煎拰绫诲瀷
func TestStructuredOutputParser(t *testing.T) {
	schema := map[string]any{
		"name": map[string]any{"type": "string"},
		"age":  map[string]any{"type": "integer"},
	}
	p := NewStructuredOutputParser(schema)

	result, err := p.Parse(`{"name": "Alice", "age": 30}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	m, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", result)
	}
	if m["name"] != "Alice" {
		t.Errorf("expected Alice, got %v", m["name"])
	}
	if m["age"] != float64(30) {
		t.Errorf("expected float64(30), got %v", m["age"])
	}
}

// 娴嬭瘯缂哄皯 schema 涓畾涔夌殑蹇呭～瀛楁 鈫?杩斿洖閿欒
func TestStructuredOutputParserMissingFields(t *testing.T) {
	schema := map[string]any{
		"required_field": map[string]any{"type": "string"},
	}
	p := NewStructuredOutputParser(schema)

	_, err := p.Parse(`{"other": "value"}`)
	if err == nil {
		t.Error("expected error for missing required field")
	}
}

// 娴嬭瘯闈炴硶 JSON 杈撳叆 鈫?杩斿洖閿欒
func TestStructuredOutputParserInvalidJSON(t *testing.T) {
	schema := map[string]any{"x": map[string]any{"type": "string"}}
	p := NewStructuredOutputParser(schema)

	_, err := p.Parse("not json")
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

// 娴嬭瘯 NewStructuredOutputParserFromFields锛氱敤 field-name鈫抎escription 鏄犲皠鍒涘缓瑙ｆ瀽鍣?func TestStructuredOutputParserFromFields(t *testing.T) {
	fields := map[string]string{
		"name": "The name",
		"age":  "The age",
	}
	p := NewStructuredOutputParserFromFields(fields)

	result, err := p.Parse(`{"name": "Bob", "age": "25"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	m := result.(map[string]any)
	if m["name"] != "Bob" {
		t.Errorf("expected Bob, got %v", m["name"])
	}
}

// 娴嬭瘯 GetFormatInstructions锛氬簲杩斿洖鍖呭惈 schema 鐨?JSON 鏍煎紡鍖栬鏄?func TestStructuredOutputParserFormatInstructions(t *testing.T) {
	schema := map[string]any{"x": map[string]any{"type": "string"}}
	p := NewStructuredOutputParser(schema)

	instr := p.GetFormatInstructions()
	if instr == "" {
		t.Error("expected non-empty format instructions")
	}
}
