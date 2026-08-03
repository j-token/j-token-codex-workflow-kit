package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/j-token/j-token-codex-workflow-kit/internal/promptdoc"
	"github.com/j-token/j-token-codex-workflow-kit/internal/review"
)

var version = "dev"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "오류:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		printHelp()
		return nil
	}

	switch args[0] {
	case "version", "--version", "-v":
		fmt.Println(version)
		return nil
	case "help", "--help", "-h":
		printHelp()
		return nil
	case "review":
		return runReview(args[1:])
	default:
		return fmt.Errorf("알 수 없는 명령 `%s`입니다. `codex-workflow --help`를 확인하세요", args[0])
	}
}

func runReview(args []string) error {
	flags := flag.NewFlagSet("review", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	jsonOutput := flags.Bool("json", false, "승인 결과를 JSON으로 출력")
	noOpen := flags.Bool("no-open", false, "브라우저를 자동으로 열지 않음")
	startAt := flags.String("start-at", "facts", "검토 시작 단계: facts, choices, plan")
	if err := flags.Parse(normalizeReviewArgs(args)); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return errors.New("사용법: codex-workflow review <prompt.md> [--json] [--no-open] [--start-at=facts|choices|plan]")
	}
	if *startAt != "facts" && *startAt != "choices" && *startAt != "plan" {
		return fmt.Errorf("지원하지 않는 시작 단계 `%s`입니다: facts, choices, plan 중 하나를 사용하세요", *startAt)
	}

	path, err := filepath.Abs(flags.Arg(0))
	if err != nil {
		return fmt.Errorf("프롬프트 문서 경로를 확인할 수 없습니다: %w", err)
	}
	document, err := promptdoc.Load(path)
	if err != nil {
		return err
	}
	if previous, err := loadReviewBaseline(path); err == nil && previous != "" {
		if previousDocument, parseErr := promptdoc.Parse(previous); parseErr == nil {
			document.PreviousReview = previousDocument.ReviewMarkdown
		}
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	result, err := review.Run(ctx, document, review.Options{OpenBrowser: !*noOpen, Writer: os.Stderr, StartAt: *startAt})
	if err != nil {
		return err
	}
	if result.Status == "feedback" {
		if err := saveReviewBaseline(path, document.Raw); err != nil {
			return fmt.Errorf("이전 검토본을 저장할 수 없습니다: %w", err)
		}
	} else {
		_ = clearReviewBaseline(path)
	}

	if *jsonOutput {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetEscapeHTML(false)
		return encoder.Encode(result)
	}
	if result.Status == "approved" {
		fmt.Printf("승인됨: model=%s reasoning_effort=%s\n", result.Model, result.ReasoningEffort)
	} else if result.Status == "feedback" {
		fmt.Printf("피드백 요청됨: comments=%d\n", len(result.Comments))
	} else {
		fmt.Println("취소됨")
	}
	return nil
}

func reviewBaselinePath(documentPath string) (string, error) {
	cacheDirectory, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256([]byte(filepath.Clean(documentPath)))
	return filepath.Join(cacheDirectory, "j-token-codex-workflow", "review-baselines", fmt.Sprintf("%x.md", digest)), nil
}

func loadReviewBaseline(documentPath string) (string, error) {
	baselinePath, err := reviewBaselinePath(documentPath)
	if err != nil {
		return "", err
	}
	content, err := os.ReadFile(baselinePath)
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func saveReviewBaseline(documentPath, content string) error {
	baselinePath, err := reviewBaselinePath(documentPath)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(baselinePath), 0o700); err != nil {
		return err
	}
	return os.WriteFile(baselinePath, []byte(content), 0o600)
}

func clearReviewBaseline(documentPath string) error {
	baselinePath, err := reviewBaselinePath(documentPath)
	if err != nil {
		return err
	}
	err = os.Remove(baselinePath)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func normalizeReviewArgs(args []string) []string {
	flags := make([]string, 0, len(args))
	positional := make([]string, 0, 1)
	for index := 0; index < len(args); index++ {
		argument := args[index]
		if len(argument) > 0 && argument[0] == '-' {
			flags = append(flags, argument)
			if argument == "--start-at" && index+1 < len(args) {
				index++
				flags = append(flags, args[index])
			}
			continue
		}
		positional = append(positional, argument)
	}
	return append(flags, positional...)
}

func printHelp() {
	fmt.Printf(`codex-workflow %s

Codex 실행 프롬프트를 브라우저에서 검토합니다.

사용법:
  codex-workflow review <prompt.md> [--json] [--no-open] [--start-at=facts|choices|plan]
  codex-workflow version

명령:
  review   프롬프트 문서의 계획·모델·추론 강도·프롬프트 검토
  version  바이너리 버전 출력
`, version)
}
