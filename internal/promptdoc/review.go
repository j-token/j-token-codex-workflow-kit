package promptdoc

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

type Fact struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Content  string `json:"content"`
	Evidence string `json:"evidence"`
}

type Question struct {
	ID          string   `json:"id"`
	Prompt      string   `json:"prompt"`
	Recommended string   `json:"recommended"`
	Multiple    bool     `json:"multiple"`
	Options     []Option `json:"options"`
}

type Option struct {
	Label       string `json:"label"`
	Description string `json:"description"`
}

var markdownRenderer = goldmark.New(
	goldmark.WithExtensions(extension.GFM, extension.CJK),
	goldmark.WithParserOptions(parser.WithAutoHeadingID()),
)

func hydrateReview(document *Document, lines []string, frontmatterEnd int) error {
	body := strings.Join(lines[frontmatterEnd+1:], "\n")
	document.Title = firstTitle(body)
	document.Facts = parseFacts(section(body, "팩트와 근거"))
	document.Questions = parseQuestions(section(body, "선택이 필요한 항목"))

	reviewMarkdown := beforeSection(body, "Codex 실행 설정", "Codex 실행 프롬프트")
	reviewMarkdown = withoutSection(reviewMarkdown, "선택이 필요한 항목")
	var rendered bytes.Buffer
	if err := markdownRenderer.Convert([]byte(reviewMarkdown), &rendered); err != nil {
		return fmt.Errorf("검토용 Markdown을 렌더링할 수 없습니다: %w", err)
	}
	document.ReviewHTML = rendered.String()
	return nil
}

func firstTitle(markdown string) string {
	for _, line := range strings.Split(markdown, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(trimmed, "# "))
		}
	}
	return "Codex 실행 계획"
}

func section(markdown, name string) string {
	lines := strings.Split(markdown, "\n")
	start := -1
	for index, line := range lines {
		if headingName(line, 2) == name {
			start = index + 1
			continue
		}
		if start >= 0 && headingLevel(line) == 2 {
			return strings.Join(lines[start:index], "\n")
		}
	}
	if start >= 0 {
		return strings.Join(lines[start:], "\n")
	}
	return ""
}

func beforeSection(markdown string, names ...string) string {
	wanted := make(map[string]bool, len(names))
	for _, name := range names {
		wanted[name] = true
	}
	lines := strings.Split(markdown, "\n")
	for index, line := range lines {
		if wanted[headingName(line, 2)] {
			return strings.TrimSpace(strings.Join(lines[:index], "\n"))
		}
	}
	return strings.TrimSpace(markdown)
}

func withoutSection(markdown, name string) string {
	lines := strings.Split(markdown, "\n")
	result := make([]string, 0, len(lines))
	skipping := false
	for _, line := range lines {
		if headingName(line, 2) == name {
			skipping = true
			continue
		}
		if skipping && headingLevel(line) == 2 {
			skipping = false
		}
		if !skipping {
			result = append(result, line)
		}
	}
	return strings.TrimSpace(strings.Join(result, "\n"))
}

func headingLevel(line string) int {
	trimmed := strings.TrimSpace(line)
	count := 0
	for count < len(trimmed) && trimmed[count] == '#' {
		count++
	}
	if count == 0 || count >= len(trimmed) || trimmed[count] != ' ' {
		return 0
	}
	return count
}

func headingName(line string, level int) string {
	if headingLevel(line) != level {
		return ""
	}
	return strings.TrimSpace(strings.TrimSpace(line)[level+1:])
}

func parseFacts(markdown string) []Fact {
	lines := strings.Split(markdown, "\n")
	headerIndex := -1
	for index, line := range lines {
		cells := splitTableRow(line)
		if len(cells) >= 4 && strings.EqualFold(strings.TrimSpace(cells[0]), "ID") {
			headerIndex = index
			break
		}
	}
	if headerIndex < 0 || headerIndex+2 >= len(lines) {
		return []Fact{}
	}

	facts := make([]Fact, 0)
	for _, line := range lines[headerIndex+2:] {
		cells := splitTableRow(line)
		if len(cells) < 4 {
			if strings.TrimSpace(line) != "" {
				break
			}
			continue
		}
		id := plainInline(cells[0])
		if id == "" {
			continue
		}
		facts = append(facts, Fact{
			ID:       id,
			Status:   plainInline(cells[1]),
			Content:  plainInline(cells[2]),
			Evidence: plainInline(cells[3]),
		})
	}
	return facts
}

func splitTableRow(line string) []string {
	trimmed := strings.TrimSpace(line)
	if !strings.Contains(trimmed, "|") {
		return nil
	}
	trimmed = strings.TrimPrefix(strings.TrimSuffix(trimmed, "|"), "|")
	var cells []string
	var current strings.Builder
	escaped := false
	for _, character := range trimmed {
		switch {
		case escaped:
			current.WriteRune(character)
			escaped = false
		case character == '\\':
			escaped = true
		case character == '|':
			cells = append(cells, strings.TrimSpace(current.String()))
			current.Reset()
		default:
			current.WriteRune(character)
		}
	}
	if escaped {
		current.WriteRune('\\')
	}
	cells = append(cells, strings.TrimSpace(current.String()))
	return cells
}

func parseQuestions(markdown string) []Question {
	lines := strings.Split(markdown, "\n")
	questions := make([]Question, 0)
	var current *Question
	flush := func() {
		if current == nil || strings.TrimSpace(current.Prompt) == "" {
			return
		}
		if len(current.Options) == 0 && current.Recommended != "" {
			current.Options = append(current.Options, Option{Label: current.Recommended, Description: "AI 권장안"})
		}
		questions = append(questions, *current)
	}

	for _, line := range lines {
		if headingLevel(line) == 3 {
			flush()
			id, prompt := splitQuestionHeading(headingName(line, 3), len(questions)+1)
			if strings.HasPrefix(prompt, "없음") || strings.HasPrefix(id, "없음") {
				current = nil
				continue
			}
			current = &Question{ID: id, Prompt: prompt, Options: []Option{}}
			continue
		}
		if current == nil {
			continue
		}
		trimmed := strings.TrimSpace(line)
		if value, ok := listValue(trimmed, "권장"); ok {
			current.Recommended = plainInline(value)
			continue
		}
		if value, ok := listValue(trimmed, "선택 방식"); ok {
			current.Multiple = isMultipleChoice(value)
			continue
		}
		if value, ok := listValue(trimmed, "옵션"); ok {
			label, description := splitOption(value)
			if label != "" {
				current.Options = append(current.Options, Option{Label: label, Description: description})
			}
		}
	}
	flush()
	return questions
}

func isMultipleChoice(value string) bool {
	normalized := strings.ToLower(strings.TrimSpace(plainInline(value)))
	switch normalized {
	case "multiple", "multi", "다중", "복수", "다중 선택", "복수 선택":
		return true
	default:
		return false
	}
}

func splitQuestionHeading(heading string, fallback int) (string, string) {
	parts := strings.SplitN(strings.TrimSpace(heading), ".", 2)
	if len(parts) == 2 && strings.TrimSpace(parts[0]) != "" {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	}
	return fmt.Sprintf("Q%d", fallback), strings.TrimSpace(heading)
}

func listValue(line, label string) (string, bool) {
	prefix := "- " + label + ":"
	if !strings.HasPrefix(line, prefix) {
		return "", false
	}
	return strings.TrimSpace(strings.TrimPrefix(line, prefix)), true
}

func splitOption(value string) (string, string) {
	plain := plainInline(value)
	for _, separator := range []string{" — ", " – ", " - ", ": "} {
		if left, right, ok := strings.Cut(plain, separator); ok {
			return strings.TrimSpace(left), strings.TrimSpace(right)
		}
	}
	return strings.TrimSpace(plain), ""
}

var (
	markdownLinkPattern = regexp.MustCompile(`\[([^\]]+)\]\([^\)]+\)`)
	markdownMarkPattern = regexp.MustCompile("[`*_~]+")
)

func plainInline(value string) string {
	value = markdownLinkPattern.ReplaceAllString(value, "$1")
	value = markdownMarkPattern.ReplaceAllString(value, "")
	return strings.TrimSpace(value)
}
