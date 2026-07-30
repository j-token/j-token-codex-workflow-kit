package main

import (
	"context"
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
	if err := flags.Parse(normalizeReviewArgs(args)); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return errors.New("사용법: codex-workflow review <prompt.md> [--json] [--no-open]")
	}

	path, err := filepath.Abs(flags.Arg(0))
	if err != nil {
		return fmt.Errorf("프롬프트 문서 경로를 확인할 수 없습니다: %w", err)
	}
	document, err := promptdoc.Load(path)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	result, err := review.Run(ctx, document, review.Options{OpenBrowser: !*noOpen, Writer: os.Stderr})
	if err != nil {
		return err
	}

	if *jsonOutput {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetEscapeHTML(false)
		return encoder.Encode(result)
	}
	if result.Status == "approved" {
		fmt.Printf("승인됨: model=%s reasoning_effort=%s\n", result.Model, result.ReasoningEffort)
	} else {
		fmt.Println("취소됨")
	}
	return nil
}

func normalizeReviewArgs(args []string) []string {
	flags := make([]string, 0, len(args))
	positional := make([]string, 0, 1)
	for _, argument := range args {
		if len(argument) > 0 && argument[0] == '-' {
			flags = append(flags, argument)
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
  codex-workflow review <prompt.md> [--json] [--no-open]
  codex-workflow version

명령:
  review   프롬프트 문서의 계획·모델·추론 강도·프롬프트 검토
  version  바이너리 버전 출력
`, version)
}
