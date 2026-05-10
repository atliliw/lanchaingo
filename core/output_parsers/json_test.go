package output_parsers

import "testing"

// 娴嬭瘯 JsonOutputParser.Parse锛氭纭В鏋?JSON 瀵硅薄锛岄獙璇?key-value 鍜屾暟瀛楃被鍨?// 娉ㄦ剰锛欸o 鐨?json.Unmarshal 榛樿灏嗘暟瀛楄В鏋愪负 float64
func TestJsonOutputParser(t *testing.T) {
	p := NewJsonOutputParser()

	result, err := p.Parse(`{"key": "value", "num": 42}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	m, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", result)
	}
	if m["key"] != "value" {
		t.Errorf("expected value, got %v", m["key"])
	}
	if m["num"] != float64(42) {
		t.Errorf("expected float64(42), got %v (%T)", m["num"], m["num"])
	}
}

// 娴嬭瘯 JsonOutputParser 澶勭悊闈炴硶 JSON 搴旇繑鍥為敊璇?func TestJsonOutputParserInvalid(t *testing.T) {
	p := NewJsonOutputParser()

	_, err := p.Parse("not json")
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

// 娴嬭瘯 JsonOutputParser 瑙ｆ瀽 JSON 鏁扮粍
func TestJsonOutputParserArray(t *testing.T) {
	p := NewJsonOutputParser()

	result, err := p.Parse(`[1, 2, 3]`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	arr, ok := result.([]any)
	if !ok {
		t.Fatalf("expected []any, got %T", result)
	}
	if len(arr) != 3 {
		t.Errorf("expected 3 elements, got %d", len(arr))
	}
}

// 娴嬭瘯 JsonOutputParser.GetFormatInstructions锛氬簲鏈夐潪绌虹殑 JSON 鏍煎紡鍖栨寚浠?func TestJsonOutputParserFormatInstructions(t *testing.T) {
	p := NewJsonOutputParser()
	if p.GetFormatInstructions() == "" {
		t.Error("expected non-empty format instructions")
	}
}

// 娴嬭瘯 ValidateJSON 宸ュ叿鍑芥暟锛氬悎娉?JSON 杩斿洖 nil锛岄潪娉?JSON 杩斿洖 error
func TestValidateJSON(t *testing.T) {
	if err := ValidateJSON(`{"a":1}`); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := ValidateJSON(`not json`); err == nil {
		t.Error("expected error for invalid JSON")
	}
}
