package llm

type ModelInfo struct {
	ID          string
	DisplayName string
	Description string
}

var Models = map[string]ModelInfo{
	"openai/gpt-oss-120b": {
		ID:          "openai/gpt-oss-120b",
		DisplayName: "GPT-OSS 120B",
		Description: "General-purpose model — Super/Plan mode",
	},
	"qwen/qwen3-coder-30b": {
		ID:          "qwen/qwen3-coder-30b",
		DisplayName: "Qwen3-Coder-30B",
		Description: "Coding-specialized MoE model — Dev mode (256K context)",
	},
}
