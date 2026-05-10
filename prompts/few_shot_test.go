package prompts

import (
	"testing"
)

// 娴嬭瘯 NewLengthBasedExampleSelector锛氭湁绀轰緥鏃?SelectExamples 杩斿洖闈炵┖缁撴灉
func TestNewLengthBasedExampleSelector(t *testing.T) {
	examples := []map[string]string{
		{"input": "hello", "output": "hi"},
		{"input": "bye", "output": "goodbye"},
	}
	sel := NewLengthBasedExampleSelector(examples, 100)

	selected, err := sel.SelectExamples(map[string]string{"input": "test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(selected) == 0 {
		t.Error("expected at least one selected example")
	}
}

// 娴嬭瘯 maxLength 榛樿鍊硷細浼?0 鏃朵娇鐢?2048
func TestNewLengthBasedExampleSelectorDefaultMaxLength(t *testing.T) {
	sel := NewLengthBasedExampleSelector(nil, 0)
	if sel.maxLength != 2048 {
		t.Errorf("expected default 2048, got %d", sel.maxLength)
	}
}

// 娴嬭瘯 AddExample锛氬姩鎬佹坊鍔犵ず渚嬪悗 SelectExamples 鑳借繑鍥炴柊绀轰緥
func TestExampleSelectorAddExample(t *testing.T) {
	sel := NewLengthBasedExampleSelector(nil, 100)
	sel.AddExample(map[string]string{"input": "q", "output": "a"})

	selected, _ := sel.SelectExamples(nil)
	if len(selected) != 1 {
		t.Errorf("expected 1 example after add, got %d", len(selected))
	}
}

// 娴嬭瘯 FewShotPromptTemplate.Format锛氬畬鏁存祦绋嬧€斺€攑refix + 鏍煎紡鍖?examples + suffix
// 楠岃瘉鎵€鏈夌粍鎴愰儴鍒嗘寜姝ｇ‘椤哄簭鎷兼帴
func TestFewShotPromptTemplateFormat(t *testing.T) {
	examples := []map[string]string{
		{"input": "sunny", "output": "It's sunny!"},
	}
	examplePrompt := NewPromptTemplate("Q: {input}\nA: {output}", []string{"input", "output"})

	fpt := NewFewShotPromptTemplate(
		examples,
		examplePrompt,
		"Now answer: {question}",
		[]string{"question"},
		"\n\n",
		"Here are some examples:",
	)

	result, err := fpt.Format(map[string]string{
		"question": "What's the weather?",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Here are some examples:\n\nQ: sunny\nA: It's sunny!\n\nNow answer: What's the weather?"
	if result != expected {
		t.Errorf("unexpected result:\nexpected: %s\n\ngot: %s", expected, result)
	}
}

// 娴嬭瘯鏃?prefix 鐨勬儏鍐碉細缁撴灉鐩存帴浠?examples 寮€濮?func TestFewShotPromptTemplateNoPrefix(t *testing.T) {
	examples := []map[string]string{
		{"input": "x", "output": "y"},
	}
	ep := NewPromptTemplate("{input}={output}", []string{"input", "output"})

	fpt := NewFewShotPromptTemplate(examples, ep, "end", nil, "\n", "")

	result, err := fpt.Format(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "x=y\nend" {
		t.Errorf("unexpected result: %s", result)
	}
}

// 娴嬭瘯 exampleSeparator 榛樿鍊硷細浼犵┖瀛楃涓叉椂浣跨敤 "\n\n"
func TestFewShotPromptTemplateDefaultSeparator(t *testing.T) {
	fpt := NewFewShotPromptTemplate(nil, nil, "", nil, "", "")
	if fpt.exampleSeparator != "\n\n" {
		t.Errorf("expected default \\n\\n separator, got %q", fpt.exampleSeparator)
	}
}
