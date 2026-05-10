package output_parsers

import "testing"

// 娴嬭瘯 CommaSeparatedListOutputParser.Parse锛氶€楀彿鍒嗛殧鐨勫瓧绗︿覆琚В鏋愪负瀛楃涓插垏鐗?// 姣忎釜鍏冪礌鑷姩鍘婚櫎鍓嶅悗绌烘牸
func TestCommaSeparatedListOutputParser(t *testing.T) {
	p := NewCommaSeparatedListOutputParser()

	result, err := p.Parse("a, b, c")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	list, ok := result.([]string)
	if !ok {
		t.Fatalf("expected []string, got %T", result)
	}
	if len(list) != 3 {
		t.Fatalf("expected 3 items, got %d", len(list))
	}
	if list[0] != "a" || list[1] != "b" || list[2] != "c" {
		t.Errorf("unexpected list content: %v", list)
	}
}

// 娴嬭瘯绌哄瓧绗︿覆瑙ｆ瀽锛氳繑鍥炵┖鍒囩墖
func TestCommaSeparatedListEmpty(t *testing.T) {
	p := NewCommaSeparatedListOutputParser()

	result, err := p.Parse("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	list, ok := result.([]string)
	if !ok {
		t.Fatalf("expected []string, got %T", result)
	}
	if len(list) != 0 {
		t.Errorf("expected empty list, got %v", list)
	}
}

// 娴嬭瘯鍗曚釜鍏冪礌锛堟棤閫楀彿锛夛細杩斿洖鍙惈涓€涓厓绱犵殑鍒囩墖
func TestCommaSeparatedListSingle(t *testing.T) {
	p := NewCommaSeparatedListOutputParser()

	result, _ := p.Parse("only")
	list := result.([]string)
	if len(list) != 1 || list[0] != "only" {
		t.Errorf("unexpected: %v", list)
	}
}

// 娴嬭瘯 CommaSeparatedListOutputParser.Type 杩斿洖姝ｇ‘鐨勮В鏋愬櫒鏍囪瘑
func TestCommaSeparatedListType(t *testing.T) {
	p := NewCommaSeparatedListOutputParser()
	if p.Type() != "comma_separated_list" {
		t.Errorf("expected comma_separated_list, got %s", p.Type())
	}
}
