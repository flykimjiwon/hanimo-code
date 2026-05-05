//go:build onprem

// Onprem (TECHAI 사내망) profile. Build with `go build -tags=onprem .` or
// `wails build -tags=onprem`. See default.go for the rationale behind the
// profile split.
//
// TECHAI / Shinhan 사내망 vLLM 운영자는 모델을 vendor-prefix 없이 등록한다
// (예: `gpt-oss-120b`). 외부망 default 빌드와 같은 카탈로그를 쓰면 endpoint
// 컨벤션 불일치로 404가 난다. 이 프로파일이 그 분기를 컴파일 타임에 박는다.
package build

const (
	Profile              = "onprem"
	RecommendedProvider  = "vllm"
	RecommendedSuperModel = "gpt-oss-120b"
	RecommendedDevModel   = "gpt-oss-120b"
)
