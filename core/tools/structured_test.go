package tools

import (
	"context"
	"encoding/json"
	"testing"
)

// 娴嬭瘯 NewStructuredTool锛氶獙璇佺粨鏋勫寲宸ュ叿鐨?Name/Description 鑳芥纭祴鍊?func TestNewFuncTool(t *testing.T) {
	tool := NewStructuredTool(
		"echo",
		"echoes input",
		"string",
		map[string]any{"type": "string"},
		func(ctx context.Context, params json.RawMessage) (string, error) {
			return string(params), nil
		},
	)

	if tool.Name() != "echo" {
		t.Errorf("expected echo, got %s", tool.Name())
	}
	if tool.Description() != "echoes input" {
		t.Errorf("expected echoes input, got %s", tool.Description())
	}
}

// 娴嬭瘯 FuncTool.Run锛氶獙璇?execute 鍥炶皟鍑芥暟琚纭皟鐢紝涓斿弬鏁伴€忎紶
func TestFuncToolRun(t *testing.T) {
	tool := NewStructuredTool(
		"echo",
		"echo",
		"string",
		map[string]any{"type": "string"},
		func(ctx context.Context, params json.RawMessage) (string, error) {
			return "echo: " + string(params), nil
		},
	)

	result, err := tool.Run(context.Background(), `"hello"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "echo: \"hello\"" {
		t.Errorf("expected echo: hello, got %s", result)
	}
}

// 娴嬭瘯 FuncTool.InputSchema 鍜?InputType锛氶獙璇佺粨鏋勫寲宸ュ叿鐨?Schema 鍏冩暟鎹?func TestFuncToolInputSchema(t *testing.T) {
	schema := map[string]any{"type": "object", "properties": map[string]any{}}
	tool := NewStructuredTool("t", "d", "Input", schema, nil)

	if tool.InputSchema()["type"] != "object" {
		t.Error("schema type mismatch")
	}
	if tool.InputType() != "Input" {
		t.Errorf("expected Input, got %s", tool.InputType())
	}
}

// 娴嬭瘯 ToolError锛氶獙璇侀敊璇秷鎭牸寮忎负 "tool {name}: {message}"
func TestToolError(t *testing.T) {
	err := NewToolError("calc", "division by zero", nil)
	if err.Error() != "tool calc: division by zero" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
	if err.ToolName != "calc" {
		t.Errorf("expected calc, got %s", err.ToolName)
	}
}
