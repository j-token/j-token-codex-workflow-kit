package review

import (
	"bytes"
	"context"
	"crypto/rand"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/j-token/j-token-codex-workflow-kit/internal/promptdoc"
)

//go:embed web/index.html web/assets/*
var assets embed.FS

type Result struct {
	Status          string          `json:"status"`
	Path            string          `json:"path"`
	Model           string          `json:"model"`
	ModelReason     string          `json:"modelReason"`
	ReasoningEffort string          `json:"reasoningEffort"`
	Prompt          string          `json:"prompt"`
	Facts           []ReviewedFact  `json:"facts,omitempty"`
	Selections      []Selection     `json:"selections,omitempty"`
	Comments        []ReviewComment `json:"comments,omitempty"`
	RestartOptions  []string        `json:"restartOptions,omitempty"`
}

type Options struct {
	OpenBrowser bool
	Writer      io.Writer
	StartAt     string
}

type submission struct {
	Action          string          `json:"action"`
	Model           string          `json:"model"`
	ModelReason     string          `json:"modelReason"`
	ReasoningEffort string          `json:"reasoningEffort"`
	Prompt          string          `json:"prompt"`
	Facts           []ReviewedFact  `json:"facts"`
	Selections      []Selection     `json:"selections"`
	Comments        []ReviewComment `json:"comments"`
}

type ReviewedFact struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Content  string `json:"content"`
	Evidence string `json:"evidence"`
}

type Selection struct {
	QuestionID string   `json:"questionId"`
	Question   string   `json:"question"`
	Option     string   `json:"option,omitempty"`
	Options    []string `json:"options,omitempty"`
	Comment    string   `json:"comment,omitempty"`
	Custom     bool     `json:"custom"`
}

type ReviewComment struct {
	ID      string `json:"id"`
	Kind    string `json:"kind"`
	Section string `json:"section,omitempty"`
	Quote   string `json:"quote,omitempty"`
	Comment string `json:"comment"`
}

type session struct {
	document promptdoc.Document
	startAt  string
	result   chan Result
	mu       sync.Mutex
	final    *Result
}

func Run(ctx context.Context, document promptdoc.Document, options Options) (Result, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return Result{}, fmt.Errorf("로컬 검토 서버를 시작할 수 없습니다: %w", err)
	}
	defer listener.Close()

	token, err := randomToken()
	if err != nil {
		return Result{}, err
	}

	s := &session{document: document, startAt: normalizeStartAt(options.StartAt), result: make(chan Result, 1)}
	basePath := "/" + token
	server := &http.Server{
		Handler:           s.handler(basePath),
		ReadHeaderTimeout: 5 * time.Second,
	}

	serveErrors := make(chan error, 1)
	go func() {
		err := server.Serve(listener)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serveErrors <- err
		}
	}()

	url := "http://" + listener.Addr().String() + basePath + "/"
	if options.Writer != nil {
		fmt.Fprintf(options.Writer, "Codex 프롬프트 검토: %s\n", url)
	}
	if options.OpenBrowser {
		if err := openBrowser(url); err != nil && options.Writer != nil {
			fmt.Fprintf(options.Writer, "브라우저를 자동으로 열지 못했습니다: %v\n", err)
		}
	}

	select {
	case result := <-s.result:
		shutdownContext, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownContext)
		return result, nil
	case err := <-serveErrors:
		return Result{}, fmt.Errorf("로컬 검토 서버 오류: %w", err)
	case <-ctx.Done():
		shutdownContext, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownContext)
		return Result{}, ctx.Err()
	}
}

func (s *session) handler(basePath string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(basePath+"/", func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != basePath+"/" {
			http.NotFound(writer, request)
			return
		}
		content, err := assets.ReadFile("web/index.html")
		if err != nil {
			http.Error(writer, "검토 화면을 읽을 수 없습니다", http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		writer.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; script-src 'self'; img-src 'self' data:; font-src 'self'; connect-src 'self'; object-src 'none'; frame-ancestors 'none'; base-uri 'none'; form-action 'self'")
		_, _ = writer.Write(content)
	})
	mux.HandleFunc(basePath+"/assets/", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet && request.Method != http.MethodHead {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		assetPath := strings.TrimPrefix(request.URL.Path, basePath+"/")
		if !strings.HasPrefix(assetPath, "assets/") || strings.Contains(assetPath, "..") {
			http.NotFound(writer, request)
			return
		}
		content, err := assets.ReadFile("web/" + assetPath)
		if err != nil {
			http.NotFound(writer, request)
			return
		}
		contentType := mime.TypeByExtension(path.Ext(assetPath))
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		writer.Header().Set("Content-Type", contentType)
		writer.Header().Set("Cache-Control", "no-store")
		http.ServeContent(writer, request, path.Base(assetPath), time.Time{}, bytes.NewReader(content))
	})
	mux.HandleFunc(basePath+"/api/document", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		writeJSON(writer, http.StatusOK, struct {
			promptdoc.Document
			StartAt string `json:"startAt"`
		}{Document: s.document, StartAt: s.startAt})
	})
	mux.HandleFunc(basePath+"/api/submit", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		request.Body = http.MaxBytesReader(writer, request.Body, 2<<20)
		var input submission
		if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
			writeJSON(writer, http.StatusBadRequest, map[string]string{"error": "요청 형식이 올바르지 않습니다"})
			return
		}

		input.Action = strings.ToLower(strings.TrimSpace(input.Action))
		if input.Action != "approve" && input.Action != "cancel" && input.Action != "feedback" {
			writeJSON(writer, http.StatusBadRequest, map[string]string{"error": "지원하지 않는 작업입니다"})
			return
		}
		if input.Action != "cancel" && (strings.TrimSpace(input.Model) == "" || strings.TrimSpace(input.ModelReason) == "" || strings.TrimSpace(input.ReasoningEffort) == "" || strings.TrimSpace(input.Prompt) == "") {
			writeJSON(writer, http.StatusBadRequest, map[string]string{"error": "모델, 모델 선택 이유, 추론 강도와 프롬프트를 모두 입력하세요"})
			return
		}
		if input.Action != "cancel" && !strings.HasPrefix(strings.TrimSpace(input.Model), "gpt-5.6-") {
			writeJSON(writer, http.StatusBadRequest, map[string]string{"error": "gpt-5.6 계열 모델만 사용할 수 있습니다"})
			return
		}
		if input.Action != "cancel" {
			if err := validateReview(input); err != nil {
				writeJSON(writer, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
		}

		status := "cancelled"
		if input.Action == "approve" {
			status = "approved"
		} else if input.Action == "feedback" {
			status = "feedback"
		}
		result := Result{
			Status:          status,
			Path:            s.document.Path,
			Model:           strings.TrimSpace(input.Model),
			ModelReason:     strings.TrimSpace(input.ModelReason),
			ReasoningEffort: strings.TrimSpace(input.ReasoningEffort),
			Prompt:          input.Prompt,
			Facts:           input.Facts,
			Selections:      input.Selections,
			Comments:        input.Comments,
		}
		if input.Action == "feedback" {
			result.RestartOptions = []string{"facts", "choices", "plan"}
		}
		final, accepted := s.finalize(result)
		if !accepted {
			writeJSON(writer, http.StatusConflict, map[string]string{
				"status": final.Status,
				"error":  "검토 결과가 이미 확정됐습니다: " + final.Status,
			})
			return
		}
		writeJSON(writer, http.StatusOK, map[string]string{"status": final.Status})
	})
	return mux
}

func validateReview(input submission) error {
	if len(input.Facts) > 100 || len(input.Selections) > 50 || len(input.Comments) > 100 {
		return errors.New("검토 항목이 허용 범위를 넘었습니다")
	}
	if !within(input.Model, 100) || !within(input.ModelReason, 2_000) || !within(input.ReasoningEffort, 100) {
		return errors.New("실행 설정 내용이 너무 깁니다")
	}
	if input.Action == "approve" && len(input.Comments) > 0 {
		return errors.New("코멘트가 있으면 계획을 승인할 수 없습니다")
	}
	if input.Action == "feedback" && len(input.Comments) == 0 {
		return errors.New("전송할 코멘트를 하나 이상 추가하세요")
	}
	for _, fact := range input.Facts {
		if !within(fact.ID, 100) || !within(fact.Status, 200) || !within(fact.Content, 10_000) || !within(fact.Evidence, 10_000) {
			return errors.New("사실 관계 검토 내용이 너무 깁니다")
		}
	}
	for _, selected := range input.Selections {
		if strings.TrimSpace(selected.Option) == "" && len(selected.Options) == 0 {
			return errors.New("모든 질문에서 옵션을 선택하세요")
		}
		if len(selected.Options) > 20 || !within(selected.QuestionID, 100) || !within(selected.Question, 2_000) || !within(selected.Option, 2_000) || !within(selected.Comment, 10_000) {
			return errors.New("옵션 검토 내용이 너무 깁니다")
		}
		for _, option := range selected.Options {
			if strings.TrimSpace(option) == "" || !within(option, 2_000) {
				return errors.New("옵션 검토 내용이 너무 깁니다")
			}
		}
	}
	for _, comment := range input.Comments {
		if strings.TrimSpace(comment.Comment) == "" {
			return errors.New("빈 코멘트는 전송할 수 없습니다")
		}
		if comment.Kind != "inline" && comment.Kind != "global" && comment.Kind != "block" {
			return errors.New("지원하지 않는 코멘트 형식입니다")
		}
		if !within(comment.ID, 100) || !within(comment.Section, 500) || !within(comment.Quote, 10_000) || !within(comment.Comment, 10_000) {
			return errors.New("코멘트 내용이 너무 깁니다")
		}
	}
	return nil
}

func normalizeStartAt(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "choices":
		return "choices"
	case "plan":
		return "plan"
	default:
		return "facts"
	}
}

func within(value string, maximum int) bool {
	return len([]rune(value)) <= maximum
}

func (s *session) finalize(candidate Result) (Result, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.final != nil {
		return *s.final, false
	}
	result := candidate
	s.final = &result
	s.result <- result
	return result, true
}

func writeJSON(writer http.ResponseWriter, status int, value any) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(value)
}

func randomToken() (string, error) {
	buffer := make([]byte, 24)
	if _, err := rand.Read(buffer); err != nil {
		return "", fmt.Errorf("검토 세션 토큰을 만들 수 없습니다: %w", err)
	}
	return hex.EncodeToString(buffer), nil
}

func openBrowser(url string) error {
	var command *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		command = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		command = exec.Command("open", url)
	default:
		command = exec.Command("xdg-open", url)
	}
	return command.Start()
}
