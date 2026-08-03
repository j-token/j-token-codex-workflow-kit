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
model_reason: "일반 구현과 검증의 균형이 필요함"
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
	if doc.ModelReason != "일반 구현과 검증의 균형이 필요함" {
		t.Fatalf("ModelReason = %q", doc.ModelReason)
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

func TestParseBuildsThreeStepReviewData(t *testing.T) {
	content := `---
model: gpt-5.6-terra
reasoning_effort: medium
---

# UI 개편 계획

## 팩트와 근거

| ID | 상태 | 내용 | 근거·확인 방법 |
| --- | --- | --- | --- |
| F1 | 확인됨 | 현재 화면은 raw Markdown을 표시한다. | index.html의 pre 요소 |
| F2 | 추정 | Mermaid 코드 블록이 포함될 수 있다. | 프롬프트 템플릿 |

## 선택이 필요한 항목

### Q1. 계획 화면의 기본 폭을 어떻게 할까요?

- 권장: 문서 중심
- 선택 방식: multiple
- 옵션: **문서 중심** — 본문 가독성을 우선합니다.
- 옵션: **전체 폭** — 한 화면의 정보량을 우선합니다.

## 실행 계획

~~~mermaid
flowchart LR
    A[사실] --> B[선택]
~~~

## Codex 실행 설정

- 모델: gpt-5.6-terra

## Codex 실행 프롬프트

~~~text
구현하세요.
~~~
`
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if doc.Title != "UI 개편 계획" {
		t.Fatalf("Title = %q", doc.Title)
	}
	if len(doc.Facts) != 2 || doc.Facts[0].ID != "F1" || !strings.Contains(doc.Facts[0].Evidence, "index.html") {
		t.Fatalf("Facts = %#v", doc.Facts)
	}
	if len(doc.Questions) != 1 || doc.Questions[0].Recommended != "문서 중심" || !doc.Questions[0].Multiple || len(doc.Questions[0].Options) != 2 {
		t.Fatalf("Questions = %#v", doc.Questions)
	}
	if !strings.Contains(doc.ReviewHTML, `class="language-mermaid"`) {
		t.Fatalf("ReviewHTML does not contain Mermaid block: %s", doc.ReviewHTML)
	}
	if !strings.Contains(doc.ReviewMarkdown, "## 실행 계획") || strings.Contains(doc.ReviewMarkdown, "## 선택이 필요한 항목") {
		t.Fatalf("ReviewMarkdown contains the wrong sections: %s", doc.ReviewMarkdown)
	}
	if strings.Contains(doc.ReviewHTML, "새 작업용 프롬프트") || strings.Contains(doc.ReviewHTML, "선택이 필요한 항목") {
		t.Fatalf("ReviewHTML contains hidden review sections: %s", doc.ReviewHTML)
	}
}

func TestParseEscapesUnsafeRawHTML(t *testing.T) {
	content := `---
model: test-model
reasoning_effort: low
---

# 안전한 문서

<script>alert("no")</script>

## Codex 실행 프롬프트

~~~text
실행
~~~
`
	doc, err := Parse(content)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if strings.Contains(doc.ReviewHTML, "<script>") {
		t.Fatalf("unsafe HTML was rendered: %s", doc.ReviewHTML)
	}
}
