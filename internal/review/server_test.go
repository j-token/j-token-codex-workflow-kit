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
		ModelReason:     "정밀 검증이 필요함",
		ReasoningEffort: "high",
		Prompt:          "updated prompt",
		Facts: []ReviewedFact{
			{ID: "F1", Status: "확인됨", Content: "수정한 사실", Evidence: "테스트"},
		},
		Selections: []Selection{
			{QuestionID: "Q1", Question: "방식?", Option: "권장안", Options: []string{"권장안"}, Comment: "의견"},
		},
	})
	postRequest := httptest.NewRequest(http.MethodPost, "/token/api/submit", bytes.NewReader(body))
	postResponse := httptest.NewRecorder()
	handler.ServeHTTP(postResponse, postRequest)
	if postResponse.Code != http.StatusOK {
		t.Fatalf("POST status = %d, body = %s", postResponse.Code, postResponse.Body.String())
	}

	result := <-s.result
	if result.Status != "approved" || result.Model != "gpt-5.6-sol" || result.ModelReason != "정밀 검증이 필요함" || result.Prompt != "updated prompt" || len(result.Facts) != 1 || len(result.Selections) != 1 {
		t.Fatalf("result = %#v", result)
	}
}

func TestFeedbackReturnsComments(t *testing.T) {
	s := &session{document: promptdoc.Document{Path: "prompt.md"}, startAt: "plan", result: make(chan Result, 1)}
	handler := s.handler("/token")

	getRequest := httptest.NewRequest(http.MethodGet, "/token/api/document", nil)
	getResponse := httptest.NewRecorder()
	handler.ServeHTTP(getResponse, getRequest)
	var documentPayload struct {
		StartAt string `json:"startAt"`
	}
	if err := json.Unmarshal(getResponse.Body.Bytes(), &documentPayload); err != nil || documentPayload.StartAt != "plan" {
		t.Fatalf("document payload = %s, error = %v", getResponse.Body.String(), err)
	}

	body, _ := json.Marshal(submission{
		Action:          "feedback",
		Model:           "gpt-5.6-terra",
		ModelReason:     "일반 구현 작업",
		ReasoningEffort: "medium",
		Prompt:          "revise",
		Comments: []ReviewComment{
			{ID: "C1", Kind: "inline", Section: "실행 계획", Quote: "기존 단계", Comment: "검증을 먼저 수행하세요"},
			{ID: "C2", Kind: "global", Comment: "범위를 줄여 주세요"},
			{ID: "C3", Kind: "block", Section: "팩트와 근거", Quote: "표 전체", Comment: "근거 열을 보강하세요"},
		},
	})
	request := httptest.NewRequest(http.MethodPost, "/token/api/submit", bytes.NewReader(body))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	result := <-s.result
	if result.Status != "feedback" || len(result.Comments) != 3 || result.Comments[0].Quote != "기존 단계" || result.Comments[2].Kind != "block" || len(result.RestartOptions) != 3 {
		t.Fatalf("result = %#v", result)
	}
}

func TestRejectsModelOutsideGPT56AndMissingReason(t *testing.T) {
	tests := []submission{
		{Action: "approve", Model: "gpt-5.5", ModelReason: "이전 모델", ReasoningEffort: "medium", Prompt: "prompt"},
		{Action: "approve", Model: "gpt-5.6-terra", ReasoningEffort: "medium", Prompt: "prompt"},
	}
	for _, input := range tests {
		s := &session{document: promptdoc.Document{}, result: make(chan Result, 1)}
		body, _ := json.Marshal(input)
		request := httptest.NewRequest(http.MethodPost, "/token/api/submit", bytes.NewReader(body))
		response := httptest.NewRecorder()
		s.handler("/token").ServeHTTP(response, request)
		if response.Code != http.StatusBadRequest {
			t.Fatalf("input = %#v, status = %d", input, response.Code)
		}
	}
}

func TestRejectsApprovalWithComments(t *testing.T) {
	s := &session{document: promptdoc.Document{}, result: make(chan Result, 1)}
	body, _ := json.Marshal(submission{
		Action:          "approve",
		Model:           "gpt-5.6-terra",
		ModelReason:     "일반 구현 작업",
		ReasoningEffort: "medium",
		Prompt:          "prompt",
		Comments:        []ReviewComment{{ID: "C1", Kind: "global", Comment: "수정 필요"}},
	})
	request := httptest.NewRequest(http.MethodPost, "/token/api/submit", bytes.NewReader(body))
	response := httptest.NewRecorder()
	s.handler("/token").ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest || !bytes.Contains(response.Body.Bytes(), []byte("승인할 수 없습니다")) {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
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
			body, _ := json.Marshal(submission{Action: action, Model: "gpt-5.6-terra", ModelReason: "일반 구현 작업", ReasoningEffort: "low", Prompt: " prompt with spaces "})
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

func TestServesEmbeddedReviewAssets(t *testing.T) {
	s := &session{document: promptdoc.Document{}, result: make(chan Result, 1)}
	handler := s.handler("/token")

	for _, target := range []string{"/token/", "/token/assets/app.css", "/token/assets/app.js", "/token/assets/mermaid.min.js"} {
		request := httptest.NewRequest(http.MethodGet, target, nil)
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if response.Code != http.StatusOK {
			t.Fatalf("GET %s status = %d", target, response.Code)
		}
		if response.Body.Len() == 0 {
			t.Fatalf("GET %s returned an empty body", target)
		}
	}
}
