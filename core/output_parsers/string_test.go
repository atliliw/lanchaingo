package output_parsers

import "testing"

// 娴嬭瘯 StrOutputParser.Parse锛氭甯告枃鏈簲鍘熸牱杩斿洖
func TestStrOutputParser(t *testing.T) {
	p := NewStrOutputParser()

	result, err := p.Parse("hello world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "hello world" {
		t.Errorf("expected hello world, got %v", result)
	}
}

// 娴嬭瘯 StrOutputParser 澶勭悊绌哄瓧绗︿覆
func TestStrOutputParserEmpty(t *testing.T) {
	p := NewStrOutputParser()

	result, err := p.Parse("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "" {
		t.Errorf("expected empty string, got %v", result)
	}
}

// 娴嬭瘯 StrOutputParser.Type 杩斿洖姝ｇ‘鐨勮В鏋愬櫒鏍囪瘑
func TestStrOutputParserType(t *testing.T) {
	p := NewStrOutputParser()
	if p.Type() != "string" {
		t.Errorf("expected string, got %s", p.Type())
	}
}

// 娴嬭瘯 StrOutputParser.GetFormatInstructions锛氭亽绛夎В鏋愬櫒涓嶉渶瑕佹牸寮忓寲鎸囦护
func TestStrOutputParserFormatInstructions(t *testing.T) {
	p := NewStrOutputParser()
	if p.GetFormatInstructions() != "" {
		t.Errorf("expected empty instructions, got %s", p.GetFormatInstructions())
	}
}
