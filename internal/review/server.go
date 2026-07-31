package review

import (
	"context"
	"crypto/rand"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/j-token/j-token-codex-workflow-kit/internal/promptdoc"
)

//go:embed web/index.html
var assets embed.FS

type Result struct {
	Status          string `json:"status"`
	Path            string `json:"path"`
	Model           string `json:"model"`
	ReasoningEffort string `json:"reasoningEffort"`
	Prompt          string `json:"prompt"`
}

type Options struct {
	OpenBrowser bool
	Writer      io.Writer
}

type submission struct {
	Action          string `json:"action"`
	Model           string `json:"model"`
	ReasoningEffort string `json:"reasoningEffort"`
	Prompt          string `json:"prompt"`
}

type session struct {
	document promptdoc.Document
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

	s := &session{document: document, result: make(chan Result, 1)}
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
		writer.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-inline'; connect-src 'self'; frame-ancestors 'none'; base-uri 'none'")
		_, _ = writer.Write(content)
	})
	mux.HandleFunc(basePath+"/api/document", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		writeJSON(writer, http.StatusOK, s.document)
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
		if input.Action != "approve" && input.Action != "cancel" {
			writeJSON(writer, http.StatusBadRequest, map[string]string{"error": "지원하지 않는 작업입니다"})
			return
		}
		if input.Action == "approve" && (strings.TrimSpace(input.Model) == "" || strings.TrimSpace(input.ReasoningEffort) == "" || strings.TrimSpace(input.Prompt) == "") {
			writeJSON(writer, http.StatusBadRequest, map[string]string{"error": "모델, 추론 강도와 프롬프트를 모두 입력하세요"})
			return
		}

		result := Result{
			Status:          map[bool]string{true: "approved", false: "cancelled"}[input.Action == "approve"],
			Path:            s.document.Path,
			Model:           strings.TrimSpace(input.Model),
			ReasoningEffort: strings.TrimSpace(input.ReasoningEffort),
			Prompt:          input.Prompt,
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
