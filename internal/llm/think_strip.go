package llm

import "strings"

// thinkTags는 처리할 thinking 태그 쌍 목록.
// 알려진 모델들의 thinking 출력 형식을 모두 커버:
//   - <think>...</think>           : Qwen3.6 표준 (HF 공식 문서 명시)
//   - <thinking>...</thinking>     : Anthropic 스타일
//   - <reasoning>...</reasoning>   : 일반 reasoning 모델
//   - <|think|>...<|/think|>       : Qwen3 alternate (파이프 변형)
//
// 새 태그 추가 시 본 슬라이스에 한 줄 append하면 자동 지원.
var thinkTags = []struct {
	open  string
	close string
}{
	{"<think>", "</think>"},
	{"<thinking>", "</thinking>"},
	{"<reasoning>", "</reasoning>"},
	{"<|think|>", "<|/think|>"},
}

// StripThinking은 입력 문자열에서 모든 thinking 블록을 제거한다.
// 다중 블록, 닫힘 누락, 빈 thinking 모두 안전 처리.
//
// **닫힘 누락 정책**: 닫는 태그가 없으면 시작 태그부터 문자열 끝까지를
// thinking으로 간주하고 모두 제거 (사용자 화면 노출 0). LLM이 thinking
// 도중 끊겨도 raw text가 새지 않음.
//
// **본문 노출 안전 보장**: 시작 태그가 발견되어 닫는 태그 검색 중인 동안
// 그 사이의 어떤 텍스트도 반환되지 않는다.
//
// 본 함수는 stateless (string in → string out). Streaming용으로는
// ThinkingStripper를 사용.
func StripThinking(s string) string {
	if s == "" {
		return s
	}
	// 어떤 thinking 태그도 포함 안 하면 빠른 경로
	hasAny := false
	for _, tag := range thinkTags {
		if strings.Contains(s, tag.open) {
			hasAny = true
			break
		}
	}
	if !hasAny {
		return s
	}

	var out strings.Builder
	out.Grow(len(s))
	rest := s

	for {
		// 가장 먼저 등장하는 시작 태그 찾기
		earliestIdx := -1
		var earliestTag struct{ open, close string }
		for _, tag := range thinkTags {
			idx := strings.Index(rest, tag.open)
			if idx >= 0 && (earliestIdx == -1 || idx < earliestIdx) {
				earliestIdx = idx
				earliestTag = struct{ open, close string }{tag.open, tag.close}
			}
		}
		if earliestIdx == -1 {
			// 더 이상 thinking 시작 태그 없음 → 나머지 그대로 출력
			out.WriteString(rest)
			break
		}

		// 시작 태그 앞 부분 출력
		out.WriteString(rest[:earliestIdx])

		// 시작 태그 이후에서 닫는 태그 찾기
		afterOpen := rest[earliestIdx+len(earliestTag.open):]
		closeIdx := strings.Index(afterOpen, earliestTag.close)
		if closeIdx == -1 {
			// 닫힘 누락 — 정책에 따라 시작 태그부터 끝까지 모두 버림.
			// 본문 노출 0 보장.
			break
		}

		// 닫는 태그까지 모두 버리고 그 뒤만 다음 iteration으로
		rest = afterOpen[closeIdx+len(earliestTag.close):]
	}

	// 연속된 빈 줄 정리 — thinking 블록만 있던 라인은 빈 줄로 남기
	// 사용자에게 거슬리는 빈 줄 다중 발생 방지.
	result := out.String()
	for strings.Contains(result, "\n\n\n") {
		result = strings.ReplaceAll(result, "\n\n\n", "\n\n")
	}
	return strings.TrimLeft(result, "\n")
}

// ThinkingStripper는 streaming 청크 단위로 thinking을 제거하는 stateful
// 필터. 청크가 태그 중간에서 분할 도착해도 부분 태그가 사용자 화면에
// 절대 노출되지 않도록 buffering 처리한다.
//
// 사용 패턴:
//
//	s := NewThinkingStripper()
//	for chunk := range stream {
//	    visibleText := s.Feed(chunk.Content)
//	    if visibleText != "" {
//	        // 화면에 출력
//	    }
//	}
//	tail := s.Flush()  // stream 끝
//	if tail != "" { /* 출력 */ }
type ThinkingStripper struct {
	// pending은 아직 출력하지 않은 buffer (partial 태그 또는 thinking 본문).
	pending strings.Builder
	// inside는 현재 thinking 블록 내부면 true.
	inside bool
	// closeTag는 inside일 때 찾고 있는 닫는 태그.
	closeTag string
}

// NewThinkingStripper는 새 stripper를 생성한다.
func NewThinkingStripper() *ThinkingStripper {
	return &ThinkingStripper{}
}

// Feed는 청크를 입력받아 사용자 화면에 출력해도 안전한 부분만 반환한다.
// partial 시작 태그는 buffer에 보존하여 다음 Feed에서 완성될 가능성에 대비.
func (s *ThinkingStripper) Feed(chunk string) string {
	if chunk == "" && s.pending.Len() == 0 {
		return ""
	}
	s.pending.WriteString(chunk)
	return s.drain(false)
}

// Flush는 stream 종료 시 호출. 미닫힌 thinking 본문은 모두 버리고,
// buffer에 남은 일반 텍스트만 반환한다.
func (s *ThinkingStripper) Flush() string {
	if s.inside {
		// 닫힘 누락 — buffer에 있는 thinking 본문 모두 버림
		s.pending.Reset()
		s.inside = false
		s.closeTag = ""
		return ""
	}
	return s.drain(true)
}

// drain은 pending buffer에서 출력 가능한 부분을 추출한다.
// finalize=true면 partial 시작 태그도 모두 출력 (stream 끝났으므로
// 더 이상 완성될 일 없음).
func (s *ThinkingStripper) drain(finalize bool) string {
	var out strings.Builder
	buf := s.pending.String()
	s.pending.Reset()

	for {
		if s.inside {
			// 닫는 태그 검색
			closeIdx := strings.Index(buf, s.closeTag)
			if closeIdx == -1 {
				// 닫는 태그 미발견 — buffer 전체 보유 (thinking 본문)
				// 단, 마지막 len(closeTag)-1 바이트만 보유해도 충분
				// (그 이상은 절대 닫는 태그를 형성할 수 없음). 메모리
				// 절약 위해 thinking 본문 자체를 버려도 안전.
				// 화면 노출 0 보장 위해 buffer는 닫는 태그 길이 -1만 남김
				if len(buf) > len(s.closeTag) {
					buf = buf[len(buf)-len(s.closeTag)+1:]
				}
				s.pending.WriteString(buf)
				return out.String()
			}
			// 닫는 태그 발견 — 그 다음부터 일반 텍스트
			buf = buf[closeIdx+len(s.closeTag):]
			s.inside = false
			s.closeTag = ""
			// 다음 루프에서 buf를 다시 검사
			continue
		}

		// !inside — 시작 태그 검색
		earliestIdx := -1
		var earliestTag struct{ open, close string }
		for _, tag := range thinkTags {
			idx := strings.Index(buf, tag.open)
			if idx >= 0 && (earliestIdx == -1 || idx < earliestIdx) {
				earliestIdx = idx
				earliestTag = struct{ open, close string }{tag.open, tag.close}
			}
		}

		if earliestIdx == -1 {
			// 시작 태그 없음 — 나머지 출력 가능 여부 판단
			if finalize {
				out.WriteString(buf)
				return out.String()
			}
			// streaming 중 — 마지막 부분이 partial 시작 태그일 가능성 체크
			safeLen := safeOutputLen(buf)
			if safeLen > 0 {
				out.WriteString(buf[:safeLen])
			}
			s.pending.WriteString(buf[safeLen:])
			return out.String()
		}

		// 시작 태그 발견 — 그 앞까지 출력하고 inside 진입
		out.WriteString(buf[:earliestIdx])
		buf = buf[earliestIdx+len(earliestTag.open):]
		s.inside = true
		s.closeTag = earliestTag.close
	}
}

// safeOutputLen은 buf에서 끝부분이 partial 시작 태그일 가능성을 고려해
// 안전하게 출력 가능한 prefix 길이를 반환한다.
//
// 예: buf = "hello <th" → safeLen = 6 (hello+space) — "<th"는 보유.
//
//	buf = "hello "    → safeLen = 6 — 전체 안전.
//	buf = "<thin"     → safeLen = 0 — 전체 보유 (다음 청크 대기).
func safeOutputLen(buf string) int {
	if buf == "" {
		return 0
	}
	maxKeep := 0
	// 가장 긴 시작 태그 길이만큼 끝쪽이 prefix-match 가능성 검사
	for _, tag := range thinkTags {
		// buf 끝쪽이 tag.open의 prefix와 일치하는 가장 긴 길이
		for k := 1; k < len(tag.open) && k <= len(buf); k++ {
			if strings.HasSuffix(buf, tag.open[:k]) {
				if k > maxKeep {
					maxKeep = k
				}
			}
		}
	}
	return len(buf) - maxKeep
}
