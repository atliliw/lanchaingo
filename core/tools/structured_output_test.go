package tools

import (
	"testing"
)

// 娴嬭瘯 NewStructuredOutput锛氫粠 Go 缁撴瀯浣撶敓鎴愮粨鏋勫寲杈撳嚭瀹氫箟锛?// 楠岃瘉 Schema/DataType 琚纭～鍏?func TestNewStructuredOutput(t *testing.T) {
	type testStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	so, err := NewStructuredOutput(testStruct{Name: "test", Value: 42})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if so == nil {
		t.Fatal("expected non-nil StructuredOutput")
	}
	if so.DataType != "tools.testStruct" {
		t.Errorf("expected tools.testStruct, got %s", so.DataType)
	}
	if so.Schema == nil {
		t.Error("expected non-nil schema")
	}
}
