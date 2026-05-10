package prompts

import "testing"

func TestNewPromptTemplate(t *testing.T) {
	pt := NewPromptTemplate("Hello {name}", []string{"name"})
	if pt.GetTemplate() != "Hello {name}" { t.Error("template mismatch") }
}
func TestPromptTemplateFormat(t *testing.T) {
	pt := NewPromptTemplate("Hello {name}", []string{"name"})
	r, _ := pt.Format(map[string]string{"name": "Alice"})
	if r != "Hello Alice" { t.Errorf("got %s", r) }
}
func TestPromptTemplateMissingVar(t *testing.T) {
	pt := NewPromptTemplate("Hello {name}", []string{"name"})
	_, err := pt.Format(map[string]string{})
	if err == nil { t.Error("expected error") }
}
func TestChatPromptTemplate(t *testing.T) {
	cpt := NewChatPromptTemplate()
	cpt.AddSystemMessage("You are {role}")
	cpt.AddHumanMessage("Hi {name}")
	msgs, _ := cpt.FormatMessages(map[string]string{"role":"assistant","name":"Bob"})
	if len(msgs) != 2 { t.Errorf("expected 2 messages") }
	if msgs[0].Content != "You are assistant" { t.Error("content mismatch") }
}
func TestFewShotPromptTemplate(t *testing.T) {
	ep := NewPromptTemplate("{input}={output}", []string{"input","output"})
	fpt := NewFewShotPromptTemplate(
		[]map[string]string{{"input":"x","output":"y"}},
		ep, "end", nil, "\n", "prefix")
	r, _ := fpt.Format(nil)
	if r != "prefix\nx=y\nend" { t.Errorf("got %s", r) }
}
