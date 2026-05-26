package llm

import (
	"strings"
	"testing"
)

// ── StripThinking (stateless) 테스트 ──

func TestStripThinking_Empty(t *testing.T) {
	if got := StripThinking(""); got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

func TestStripThinking_NoThinking(t *testing.T) {
	in := "안녕하세요. 일반 답변입니다."
	if got := StripThinking(in); got != in {
		t.Errorf("got %q, want %q", got, in)
	}
}

func TestStripThinking_SimpleThinkTag(t *testing.T) {
	in := "<think>The user wants help.</think>안녕하세요!"
	want := "안녕하세요!"
	if got := StripThinking(in); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestStripThinking_ThinkingTag_Long(t *testing.T) {
	// Anthropic 스타일 <thinking>
	in := "<thinking>\nLet me consider...\n</thinking>\n답변입니다"
	got := StripThinking(in)
	if strings.Contains(got, "thinking") || strings.Contains(got, "consider") {
		t.Errorf("thinking 누수: %q", got)
	}
	if !strings.Contains(got, "답변") {
		t.Errorf("본문 손실: %q", got)
	}
}

func TestStripThinking_ReasoningTag(t *testing.T) {
	in := "<reasoning>Step 1: ...\nStep 2: ...</reasoning>최종 답"
	want := "최종 답"
	if got := StripThinking(in); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestStripThinking_PipeTag(t *testing.T) {
	// Qwen3 alternate 형식
	in := "<|think|>thinking content<|/think|>visible"
	want := "visible"
	if got := StripThinking(in); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestStripThinking_MultipleBlocks(t *testing.T) {
	in := "<think>a</think>X<think>b</think>Y<think>c</think>Z"
	want := "XYZ"
	if got := StripThinking(in); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestStripThinking_EmptyThinking(t *testing.T) {
	in := "<think></think>본문"
	want := "본문"
	if got := StripThinking(in); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestStripThinking_UnclosedThinking_NoLeak(t *testing.T) {
	// 닫힘 누락 — 시작 태그부터 끝까지 모두 버림 (본문 노출 0)
	in := "안녕<think>thinking but never closes\nmore thinking content"
	got := StripThinking(in)
	if strings.Contains(got, "thinking") {
		t.Errorf("닫힘 누락 시 누수 — got %q", got)
	}
	// "안녕"은 보존
	if !strings.Contains(got, "안녕") {
		t.Errorf("시작 태그 앞 본문 손실: %q", got)
	}
}

func TestStripThinking_OrphanCloseTag(t *testing.T) {
	// 닫는 태그만 있고 시작 태그 없음 — literal로 처리 (드물지만 안전)
	in := "no opening</think>foo"
	got := StripThinking(in)
	// 본 구현은 시작 태그 없으면 그대로 둠 (일반 텍스트로 취급)
	if !strings.Contains(got, "no opening") || !strings.Contains(got, "foo") {
		t.Errorf("orphan close 처리 오류: %q", got)
	}
}

func TestStripThinking_ThinkingBetweenContent(t *testing.T) {
	in := "Hello <think>internal</think>world"
	want := "Hello world"
	if got := StripThinking(in); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestStripThinking_LargeThinking(t *testing.T) {
	// 매우 긴 thinking (10K chars)
	thinkContent := strings.Repeat("thinking ", 1000)
	in := "<think>" + thinkContent + "</think>짧은 답"
	want := "짧은 답"
	if got := StripThinking(in); got != want {
		t.Errorf("긴 thinking 처리 실패: got %d chars, want %q", len(got), want)
	}
}

// ── ThinkingStripper (stateful streaming) 테스트 ──

func TestStripper_FullChunkAtOnce(t *testing.T) {
	s := NewThinkingStripper()
	got := s.Feed("<think>internal</think>visible")
	got += s.Flush()
	if got != "visible" {
		t.Errorf("got %q, want %q", got, "visible")
	}
}

func TestStripper_OpeningTagSplit(t *testing.T) {
	// <th + ink>...</think>
	s := NewThinkingStripper()
	out := ""
	out += s.Feed("hello <th")
	// 이 시점에 "<th"는 partial 시작 태그 — 출력 안 됨
	if strings.Contains(out, "<th") {
		t.Errorf("partial 시작 태그 누수: %q", out)
	}
	out += s.Feed("ink>internal</think>world")
	out += s.Flush()
	want := "hello world"
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

func TestStripper_ClosingTagSplit(t *testing.T) {
	// <think>x</thi + nk>foo
	s := NewThinkingStripper()
	out := ""
	out += s.Feed("<think>x</thi")
	out += s.Feed("nk>foo")
	out += s.Flush()
	if out != "foo" {
		t.Errorf("got %q, want %q", out, "foo")
	}
}

func TestStripper_PartialTagAtEnd_NoLeak(t *testing.T) {
	// "abc<th" 까지만 도착하고 stream 끝난 경우
	// Flush 시점에 <th는 일반 텍스트로 취급되어 출력 (false positive 허용)
	s := NewThinkingStripper()
	out := s.Feed("abc<th")
	out += s.Flush()
	// Flush는 partial을 그대로 출력 (정확한 시작 태그가 아니면 본문)
	if !strings.HasPrefix(out, "abc") {
		t.Errorf("정상 본문 'abc' 누락: %q", out)
	}
}

func TestStripper_MultipleChunksOneBlock(t *testing.T) {
	s := NewThinkingStripper()
	out := ""
	chunks := []string{
		"<think>",
		"part1 ",
		"part2 ",
		"part3</think>",
		"answer",
	}
	for _, c := range chunks {
		out += s.Feed(c)
	}
	out += s.Flush()
	if out != "answer" {
		t.Errorf("got %q, want %q", out, "answer")
	}
}

func TestStripper_MultipleBlocksAcrossChunks(t *testing.T) {
	s := NewThinkingStripper()
	out := ""
	out += s.Feed("<think>a")
	out += s.Feed("</think>X")
	out += s.Feed("<think>b</think>Y")
	out += s.Flush()
	if out != "XY" {
		t.Errorf("got %q, want %q", out, "XY")
	}
}

func TestStripper_UnclosedAtFlush_DropsContent(t *testing.T) {
	s := NewThinkingStripper()
	out := s.Feed("visible <think>started but no end")
	out += s.Flush()
	// "visible "는 보존, thinking 부분은 모두 버림
	if !strings.Contains(out, "visible") {
		t.Errorf("정상 본문 손실: %q", out)
	}
	if strings.Contains(out, "started") || strings.Contains(out, "end") {
		t.Errorf("미닫힌 thinking 누수: %q", out)
	}
}

func TestStripper_RealWorldQwen36Pattern(t *testing.T) {
	// Qwen3.6의 실제 응답 형태 (HF 문서 기준):
	// <think>\nreasoning here\n</think>\n\nactual answer
	s := NewThinkingStripper()
	out := s.Feed("<think>\nThe user")
	out += s.Feed(" wants me")
	out += s.Feed(" to count to 3.\n</think>\n\n")
	out += s.Feed("1, 2, 3!")
	out += s.Flush()

	if strings.Contains(out, "user wants") || strings.Contains(out, "count to") {
		t.Errorf("thinking 본문 누수: %q", out)
	}
	if !strings.Contains(out, "1, 2, 3!") {
		t.Errorf("실제 답변 손실: %q", out)
	}
}

func TestStripper_AnswerThenThinking(t *testing.T) {
	// 답변 먼저 + thinking (드문 케이스, 그래도 처리)
	s := NewThinkingStripper()
	out := s.Feed("Quick answer. ")
	out += s.Feed("<think>self-check</think>")
	out += s.Feed(" Done.")
	out += s.Flush()
	if !strings.Contains(out, "Quick answer") {
		t.Errorf("앞 본문 손실: %q", out)
	}
	if strings.Contains(out, "self-check") {
		t.Errorf("thinking 누수: %q", out)
	}
	if !strings.Contains(out, "Done") {
		t.Errorf("뒤 본문 손실: %q", out)
	}
}

func TestStripper_OnlyThinking_NoAnswer(t *testing.T) {
	// 응답이 thinking만 있고 답변 0 — 드물지만 가능
	s := NewThinkingStripper()
	out := s.Feed("<think>only thinking</think>")
	out += s.Flush()
	if out != "" {
		t.Errorf("thinking-only 응답이 비어야 함: %q", out)
	}
}

func TestStripper_PipeVariant_Streaming(t *testing.T) {
	// <|think|> 변형이 청크 분할로 와도 정상 처리
	s := NewThinkingStripper()
	out := ""
	out += s.Feed("<|thi")
	out += s.Feed("nk|>internal<|/think|>")
	out += s.Feed("answer")
	out += s.Flush()
	if out != "answer" {
		t.Errorf("got %q, want %q", out, "answer")
	}
}

func TestStripper_NoThinking_PassThrough(t *testing.T) {
	// thinking 없는 일반 응답이 분할 도착해도 그대로
	s := NewThinkingStripper()
	out := ""
	out += s.Feed("Hello, ")
	out += s.Feed("how can ")
	out += s.Feed("I help?")
	out += s.Flush()
	if out != "Hello, how can I help?" {
		t.Errorf("got %q", out)
	}
}

// ── safeOutputLen 단위 ──

func TestSafeOutputLen(t *testing.T) {
	cases := []struct {
		buf     string
		wantOut int
	}{
		{"hello", 5},        // 일반 텍스트 — 전체 안전
		{"hello <", 6},      // "<"는 partial 시작 가능성
		{"hello <t", 6},     // "<t"는 partial
		{"hello <th", 6},    // "<th"는 partial
		{"hello <thi", 6},   // "<thi"는 partial
		{"hello <thin", 6},  // "<thin"는 partial
		{"hello <think", 6}, // "<think"는 partial (앞 6자 안전)
		{"<think", 0},       // 전체가 partial
		{"", 0},
		{"safe text", 9},
	}
	for _, c := range cases {
		got := safeOutputLen(c.buf)
		if got != c.wantOut {
			t.Errorf("safeOutputLen(%q) = %d, want %d", c.buf, got, c.wantOut)
		}
	}
}
