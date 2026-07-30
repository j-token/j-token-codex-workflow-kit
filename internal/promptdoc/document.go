package promptdoc

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	promptMarker         = "<!-- codex-workflow:prompt:start -->"
	koreanPromptHeading  = "## Codex 실행 프롬프트"
	englishPromptHeading = "## Codex Execution Prompt"
)

type Document struct {
	Path            string `json:"path"`
	Raw             string `json:"raw"`
	Model           string `json:"model"`
	ReasoningEffort string `json:"reasoningEffort"`
	Prompt          string `json:"prompt"`
}

func Load(path string) (Document, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return Document{}, fmt.Errorf("프롬프트 문서를 읽을 수 없습니다: %w", err)
	}

	doc, err := Parse(string(content))
	if err != nil {
		return Document{}, err
	}
	doc.Path = path
	return doc, nil
}

func Parse(content string) (Document, error) {
	normalized := strings.TrimPrefix(strings.ReplaceAll(content, "\r\n", "\n"), "\uFEFF")
	doc := Document{Raw: normalized}
	lines := strings.Split(normalized, "\n")

	if len(lines) < 3 || strings.TrimSpace(lines[0]) != "---" {
		return Document{}, errors.New("프롬프트 문서에 YAML frontmatter가 없습니다")
	}

	frontmatterEnd := -1
	for index := 1; index < len(lines); index++ {
		line := strings.TrimSpace(lines[index])
		if line == "---" {
			frontmatterEnd = index
			break
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		switch strings.TrimSpace(key) {
		case "model":
			doc.Model = trimYAMLScalar(value)
		case "reasoning_effort":
			doc.ReasoningEffort = trimYAMLScalar(value)
		}
	}
	if frontmatterEnd < 0 {
		return Document{}, errors.New("프롬프트 문서의 YAML frontmatter가 닫히지 않았습니다")
	}

	promptStart := -1
	for index := frontmatterEnd + 1; index < len(lines); index++ {
		line := strings.TrimSpace(lines[index])
		if line == promptMarker || line == koreanPromptHeading || line == englishPromptHeading {
			promptStart = index + 1
			break
		}
	}
	if promptStart < 0 {
		return Document{}, fmt.Errorf("프롬프트 문서에 `%s` 표시 또는 Codex 실행 프롬프트 섹션이 없습니다", promptMarker)
	}

	fenceStart := -1
	fence := ""
	for index := promptStart; index < len(lines); index++ {
		if marker, ok := openingFence(lines[index]); ok {
			fenceStart = index + 1
			fence = marker
			break
		}
	}
	if fenceStart < 0 {
		return Document{}, errors.New("Codex 실행 프롬프트 코드 블록이 없습니다")
	}

	fenceEnd := -1
	for index := fenceStart; index < len(lines); index++ {
		if isClosingFence(lines[index], fence) {
			fenceEnd = index
			break
		}
	}
	if fenceEnd < 0 {
		return Document{}, errors.New("Codex 실행 프롬프트 코드 블록이 닫히지 않았습니다")
	}

	doc.Prompt = strings.Join(lines[fenceStart:fenceEnd], "\n")
	if doc.Model == "" {
		return Document{}, errors.New("frontmatter의 `model` 값이 비어 있습니다")
	}
	if doc.ReasoningEffort == "" {
		return Document{}, errors.New("frontmatter의 `reasoning_effort` 값이 비어 있습니다")
	}
	if strings.TrimSpace(doc.Prompt) == "" {
		return Document{}, errors.New("Codex 실행 프롬프트가 비어 있습니다")
	}

	return doc, nil
}

func openingFence(line string) (string, bool) {
	trimmed := strings.TrimSpace(line)
	if len(trimmed) < 3 || (trimmed[0] != '`' && trimmed[0] != '~') {
		return "", false
	}
	character := trimmed[0]
	count := 0
	for count < len(trimmed) && trimmed[count] == character {
		count++
	}
	if count < 3 {
		return "", false
	}
	return strings.Repeat(string(character), count), true
}

func isClosingFence(line, opening string) bool {
	trimmed := strings.TrimSpace(line)
	if opening == "" || len(trimmed) < len(opening) || trimmed[0] != opening[0] {
		return false
	}
	for index := range trimmed {
		if trimmed[index] != opening[0] {
			return false
		}
	}
	return len(trimmed) >= len(opening)
}

func trimYAMLScalar(value string) string {
	return strings.Trim(strings.TrimSpace(value), "\"'")
}
