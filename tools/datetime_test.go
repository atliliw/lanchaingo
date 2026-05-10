package tools

import (
	"context"
	"testing"
)

// 娴嬭瘯 DateTimeTool 杩斿洖闈炵┖鏃堕棿瀛楃涓?func TestDateTimeTool(t *testing.T) {
	d := NewDateTimeTool()
	result, err := d.Run(context.Background(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == "" {
		t.Error("expected non-empty datetime string")
	}
}

// 娴嬭瘯 DateTimeTool 鍏冩暟鎹?func TestDateTimeMeta(t *testing.T) {
	d := NewDateTimeTool()
	if d.Name() != "datetime" {
		t.Errorf("expected datetime, got %s", d.Name())
	}
	if d.Description() == "" {
		t.Error("expected non-empty description")
	}
}
