package core

import (
	"context"
	"testing"
)

// 娴嬭瘯 RunnableFunc.Invoke锛氶獙璇佸嚱鏁伴€傞厤鍣ㄨ兘姝ｇ‘鎵ц骞惰繑鍥炵粨鏋?func TestRunnableFunc(t *testing.T) {
	fn := RunnableFunc(func(ctx context.Context, input map[string]any) (any, error) {
		return "ok", nil
	})

	result, err := fn.Invoke(context.Background(), map[string]any{"key": "val"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "ok" {
		t.Errorf("expected ok, got %v", result)
	}
}

// 娴嬭瘯 RunnableFunc.Batch锛氶獙璇佹壒閲忓鐞嗚兘姝ｇ‘閬嶅巻澶氫釜杈撳叆骞惰繑鍥炲搴旂粨鏋?func TestRunnableFuncBatch(t *testing.T) {
	fn := RunnableFunc(func(ctx context.Context, input map[string]any) (any, error) {
		return input["x"], nil
	})

	inputs := []map[string]any{
		{"x": 1},
		{"x": 2},
		{"x": 3},
	}

	results, err := fn.Batch(context.Background(), inputs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for i, r := range results {
		if r != i+1 {
			t.Errorf("result[%d] expected %d, got %v", i, i+1, r)
		}
	}
}

// 娴嬭瘯 RunnableFunc.Stream锛氶粯璁ゅ疄鐜板皢 Invoke 缁撴灉鍖呰涓哄崟鍏冪礌 channel锛?// 璇诲彇绗竴涓厓绱犲悗 channel 鑷姩鍏抽棴
func TestRunnableFuncStream(t *testing.T) {
	fn := RunnableFunc(func(ctx context.Context, input map[string]any) (any, error) {
		return "streamed", nil
	})

	ch, err := fn.Stream(context.Background(), map[string]any{"x": 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := <-ch
	if result != "streamed" {
		t.Errorf("expected streamed, got %v", result)
	}

	// 榛樿 Stream 鍙骇鐢熶竴涓厓绱狅紝channel 搴斿凡鍏抽棴
	_, ok := <-ch
	if ok {
		t.Error("channel should be closed after single element")
	}
}
