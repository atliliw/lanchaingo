package tools

import (
	"testing"
)

// 娴嬭瘯 DuckDuckGoSearchTool 鍏冩暟鎹?func TestSearchMeta(t *testing.T) {
	s := NewDuckDuckGoSearchTool()
	if s.Name() != "search" {
		t.Errorf("expected search, got %s", s.Name())
	}
	if s.Description() == "" {
		t.Error("expected non-empty description")
	}
}

// 娴嬭瘯 URLFetchTool 鍏冩暟鎹?func TestURLFetchMeta(t *testing.T) {
	u := NewURLFetchTool()
	if u.Name() != "url_fetch" {
		t.Errorf("expected url_fetch, got %s", u.Name())
	}
	if u.Description() == "" {
		t.Error("expected non-empty description")
	}
}

// 娴嬭瘯 WikipediaTool 鍏冩暟鎹?func TestWikipediaMeta(t *testing.T) {
	w := NewWikipediaTool()
	if w.Name() != "wikipedia" {
		t.Errorf("expected wikipedia, got %s", w.Name())
	}
	if w.Description() == "" {
		t.Error("expected non-empty description")
	}
}

// 娴嬭瘯 PythonREPLTool 鍏冩暟鎹?func TestPythonREPLMeta(t *testing.T) {
	p := NewPythonREPLTool()
	if p.Name() != "python_repl" {
		t.Errorf("expected python_repl, got %s", p.Name())
	}
	if p.Description() == "" {
		t.Error("expected non-empty description")
	}
}
