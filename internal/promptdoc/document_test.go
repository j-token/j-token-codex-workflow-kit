package promptdoc

import (
	"strings"
	"testing"
)

const validDocument = `---
prompt: binary-review
status: ready
type: plan
risk: standard
model: gpt-5.6-terra
reasoning_effort: medium
---

# 바이너리 검토

## Codex 실행 프롬프트

` + "```text" + `
# 지시

변경을 구현하고 검증하세요.
` + "```" + `
`

func TestParse(t *testing.T) {
	doc, err := Parse(validDocument)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if doc.Model != "gpt-5.6-terra" {
		t.Fatalf("Model = %q", doc.Model)
	}
	if doc.ReasoningEffort != "medium" {
		t.Fatalf("ReasoningEffort = %q", doc.ReasoningEffort)
	}
	if !strings.Contains(doc.Prompt, "변경을 구현") {
		t.Fatalf("Prompt = %q", doc.Prompt)
	}
}

func TestParseRequiresPromptSection(t *testing.T) {
	_, err := Parse("---\nmodel: test\nreasoning_effort: low\n---\n")
	if err == nil || !strings.Contains(err.Error(), "Codex 실행 프롬프트") {
		t.Fatalf("error = %v", err)
	}
}

func TestParseSupportsBOMCRLFAndLanguageIndependentMarker(t *testing.T) {
	content := "\uFEFF---\r\nmodel: test-model\r\nreasoning_effort: low\r\n---\r\n\r\n## 任意の見出し\r\n\r\n<!-- codex-workflow:prompt:start -->\r\n~~~text\r\nkeep this prompt\r\n~~~\r\n"
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if doc.Prompt != "keep this prompt" {
		t.Fatalf("Prompt = %q", doc.Prompt)
	}
}

func TestParsePreservesNestedFenceAndBoundaryWhitespace(t *testing.T) {
	content := `---
model: test-model
reasoning_effort: medium
---

## Codex Execution Prompt

` + "````text" + `

first line
` + "```go" + `
fmt.Println("nested")
` + "```" + `
last line

` + "````" + `
`
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if !strings.HasPrefix(doc.Prompt, "\nfirst line") || !strings.Contains(doc.Prompt, "```go") || !strings.HasSuffix(doc.Prompt, "last line\n") {
		t.Fatalf("Prompt boundary or nested fence was changed: %q", doc.Prompt)
	}
}
