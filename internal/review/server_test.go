package review

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/j-token/j-token-codex-workflow-kit/internal/promptdoc"
)

func TestDocumentAndApprove(t *testing.T) {
	s := &session{
		document: promptdoc.Document{
			Path:            "prompt.md",
			Raw:             "document",
			Model:           "gpt-5.6-terra",
			ReasoningEffort: "medium",
			Prompt:          "implement",
		},
		result: make(chan Result, 1),
	}
	handler := s.handler("/token")

	getRequest := httptest.NewRequest(http.MethodGet, "/token/api/document", nil)
	getResponse := httptest.NewRecorder()
	handler.ServeHTTP(getResponse, getRequest)
	if getResponse.Code != http.StatusOK {
		t.Fatalf("GET status = %d", getResponse.Code)
	}

	body, _ := json.Marshal(submission{
		Action:          "approve",
		Model:           "gpt-5.6-sol",
		ReasoningEffort: "high",
		Prompt:          "updated prompt",
	})
	postRequest := httptest.NewRequest(http.MethodPost, "/token/api/submit", bytes.NewReader(body))
	postResponse := httptest.NewRecorder()
	handler.ServeHTTP(postResponse, postRequest)
	if postResponse.Code != http.StatusOK {
		t.Fatalf("POST status = %d, body = %s", postResponse.Code, postResponse.Body.String())
	}

	result := <-s.result
	if result.Status != "approved" || result.Model != "gpt-5.6-sol" || result.Prompt != "updated prompt" {
		t.Fatalf("result = %#v", result)
	}
}

func TestRejectsEmptyApproval(t *testing.T) {
	s := &session{document: promptdoc.Document{}, result: make(chan Result, 1)}
	request := httptest.NewRequest(http.MethodPost, "/token/api/submit", bytes.NewBufferString(`{"action":"approve"}`))
	response := httptest.NewRecorder()
	s.handler("/token").ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", response.Code)
	}
}

func TestOnlyOneConcurrentSubmissionWins(t *testing.T) {
	s := &session{document: promptdoc.Document{Path: "prompt.md"}, result: make(chan Result, 1)}
	handler := s.handler("/token")
	actions := []string{"approve", "cancel"}
	statuses := make(chan int, len(actions))
	var wait sync.WaitGroup

	for _, action := range actions {
		wait.Add(1)
		go func(action string) {
			defer wait.Done()
			body, _ := json.Marshal(submission{Action: action, Model: "model", ReasoningEffort: "low", Prompt: " prompt with spaces "})
			request := httptest.NewRequest(http.MethodPost, "/token/api/submit", bytes.NewReader(body))
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			statuses <- response.Code
		}(action)
	}
	wait.Wait()
	close(statuses)

	okCount := 0
	conflictCount := 0
	for status := range statuses {
		switch status {
		case http.StatusOK:
			okCount++
		case http.StatusConflict:
			conflictCount++
		}
	}
	if okCount != 1 || conflictCount != 1 {
		t.Fatalf("status counts: ok=%d conflict=%d", okCount, conflictCount)
	}
	result := <-s.result
	if result.Status == "approved" && result.Prompt != " prompt with spaces " {
		t.Fatalf("approved prompt whitespace changed: %q", result.Prompt)
	}
}
