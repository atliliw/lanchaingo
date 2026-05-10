package prompts

import (
	"testing"
)

// 娴嬭瘯 NewPromptTemplate锛氭ā鏉垮瓧绗︿覆鍜屽彉閲忓垪琛ㄦ纭垵濮嬪寲
func TestNewPromptTemplate(t *testing.T) {
	pt := NewPromptTemplate("Hello {name}", []string{"name"})
	if pt == nil {
		t.Fatal("expected non-nil PromptTemplate")
	}
	if pt.GetTemplate() != "Hello {name}" {
		t.Errorf("unexpected template: %s", pt.GetTemplate())
	}
}

// 娴嬭瘯 Format锛氬涓彉閲忓悓鏃舵浛鎹紝缁撴灉姝ｇ‘
func TestPromptTemplateFormat(t *testing.T) {
	pt := NewPromptTemplate("Hello {name}, you are {age} years old", []string{"name", "age"})

	result, err := pt.Format(map[string]string{
		"name": "Alice",
		"age":  "30",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Hello Alice, you are 30 years old" {
		t.Errorf("unexpected result: %s", result)
	}
}

// 娴嬭瘯 Format 缂哄皯蹇呭～鍙橀噺 鈫?杩斿洖 error
func TestPromptTemplateFormatMissingVariable(t *testing.T) {
	pt := NewPromptTemplate("Hello {name}", []string{"name"})

	_, err := pt.Format(map[string]string{})
	if err == nil {
		t.Error("expected error for missing variable")
	}
}

// 娴嬭瘯 Format 浼犲叆 nil 鈫?杩斿洖 error
func TestPromptTemplateFormatNilMap(t *testing.T) {
	pt := NewPromptTemplate("Hello {name}", []string{"name"})

	_, err := pt.Format(nil)
	if err == nil {
		t.Error("expected error for nil values map")
	}
}

// 娴嬭瘯 Format 澶氫釜鍗犱綅绗﹂噸澶嶅嚭鐜帮細姣忎釜 {var} 閮借姝ｇ‘鏇挎崲
func TestPromptTemplateMultiplePlaceholders(t *testing.T) {
	pt := NewPromptTemplate("{a} + {b} = {c}", []string{"a", "b", "c"})

	result, err := pt.Format(map[string]string{
		"a": "1",
		"b": "2",
		"c": "3",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "1 + 2 = 3" {
		t.Errorf("unexpected result: %s", result)
	}
}

// 娴嬭瘯 GetInputVariables锛氳繑鍥炵殑鍒囩墖鏄嫹璐濓紝澶栭儴淇敼涓嶅奖鍝嶅唴閮ㄧ姸鎬?func TestPromptTemplateGetInputVariables(t *testing.T) {
	pt := NewPromptTemplate("Test", []string{"x", "y"})
	vars := pt.GetInputVariables()

	if len(vars) != 2 {
		t.Fatalf("expected 2 vars, got %d", len(vars))
	}
	if vars[0] != "x" || vars[1] != "y" {
		t.Errorf("unexpected vars: %v", vars)
	}

	vars[0] = "modified"
	if pt.inputVariables[0] == "modified" {
		t.Error("GetInputVariables should return a copy")
	}
}
