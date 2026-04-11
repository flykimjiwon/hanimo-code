package providers

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	openai "github.com/sashabaranov/go-openai"
)

func TestFriendlyError(t *testing.T) {
	cases := []struct {
		name   string
		err    error
		wants  string // substring that must appear in the localized message
	}{
		{"nil", nil, ""},
		{"401", &openai.APIError{HTTPStatusCode: 401, Message: "unauthorized"}, "인증 실패"},
		{"403", &openai.APIError{HTTPStatusCode: 403, Message: "forbidden"}, "권한 거부"},
		{"404", &openai.APIError{HTTPStatusCode: 404, Message: "not found"}, "모델을 찾을 수 없음"},
		{"429", &openai.APIError{HTTPStatusCode: 429, Message: "rate"}, "요청 한도 초과"},
		{"500", &openai.APIError{HTTPStatusCode: 500, Message: "boom"}, "서버 오류"},
		{"timeout", fmt.Errorf("context deadline exceeded"), "연결 타임아웃"},
		{"refused", fmt.Errorf("dial tcp: connection refused"), "연결 거부"},
		{"dns", fmt.Errorf("lookup api.x.example: no such host"), "호스트를 찾을 수 없음"},
	}
	for _, c := range cases {
		got := friendlyError(c.err)
		if c.err == nil {
			if got != nil {
				t.Errorf("%s: want nil, got %v", c.name, got)
			}
			continue
		}
		if got == nil || !strings.Contains(got.Error(), c.wants) {
			t.Errorf("%s: want %q in message, got %v", c.name, c.wants, got)
		}
		if !errors.Is(got, c.err) {
			t.Errorf("%s: errors.Is chain broken", c.name)
		}
	}
}

func TestRetryableStatus(t *testing.T) {
	retryable := []int{408, 429, 500, 502, 503, 504, 599}
	notRetryable := []int{200, 400, 401, 403, 404}
	for _, s := range retryable {
		if !retryableStatus(s) {
			t.Errorf("status %d should be retryable", s)
		}
	}
	for _, s := range notRetryable {
		if retryableStatus(s) {
			t.Errorf("status %d should NOT be retryable", s)
		}
	}
}
