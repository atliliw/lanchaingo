package output_parsers

import (
	"testing"
)

// testPerson 鏄敤浜庢祴璇?TypedOutputParser 鐨勭ず渚嬬粨鏋勪綋
// json tag 瀹氫箟浜?JSON 瀛楁鍚嶆槧灏?type testPerson struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// 娴嬭瘯 TypedOutputParser.Parse锛欽SON 鑳芥纭В鏋愪负鐩爣缁撴瀯浣撶被鍨?// 楠岃瘉瀛楁鍊煎畬鏁翠笖绫诲瀷姝ｇ‘
func TestTypedOutputParser(t *testing.T) {
	p := NewTypedOutputParser(testPerson{})

	result, err := p.Parse(`{"name": "Alice", "age": 30}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	person, ok := result.(testPerson)
	if !ok {
		t.Fatalf("expected testPerson, got %T", result)
	}
	if person.Name != "Alice" {
		t.Errorf("expected Alice, got %s", person.Name)
	}
	if person.Age != 30 {
		t.Errorf("expected 30, got %d", person.Age)
	}
}

// 娴嬭瘯 TypedOutputParser 澶勭悊闈炴硶 JSON 鈫?杩斿洖閿欒
func TestTypedOutputParserInvalidJSON(t *testing.T) {
	p := NewTypedOutputParser(testPerson{})

	_, err := p.Parse("not json")
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

// 娴嬭瘯 TypedOutputParser.Type 杩斿洖 "typed({绫诲瀷鍚峿)" 鏍煎紡
func TestTypedOutputParserTypeName(t *testing.T) {
	p := NewTypedOutputParser(testPerson{})
	if p.Type() != "typed(testPerson)" {
		t.Errorf("expected typed(testPerson), got %s", p.Type())
	}
}

// 娴嬭瘯 GetFormatInstructions锛氬簲鏈夊寘鍚ず渚?JSON 鐨勬牸寮忓寲璇存槑
func TestTypedOutputParserFormatInstructions(t *testing.T) {
	p := NewTypedOutputParser(testPerson{})
	if p.GetFormatInstructions() == "" {
		t.Error("expected non-empty format instructions")
	}
}
