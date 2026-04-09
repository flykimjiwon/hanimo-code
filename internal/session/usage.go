package session

import (
	"fmt"
)

// UsageSummary holds aggregated token/cost totals.
type UsageSummary struct {
	TokensIn  int64
	TokensOut int64
	CostUSD   float64
}

// pricing per 1M tokens (input, output) in USD.
var pricing = map[string][2]float64{
	// OpenAI
	"gpt-4o":            {2.50, 10.00},
	"gpt-4o-mini":       {0.15, 0.60},
	"gpt-4.1":           {2.00, 8.00},
	"gpt-4.1-mini":      {0.40, 1.60},
	"gpt-4.1-nano":      {0.10, 0.40},
	"o3":                {2.00, 8.00},
	"o3-mini":           {1.10, 4.40},
	"o4-mini":           {1.10, 4.40},
	// Anthropic
	"claude-sonnet-4-20250514":    {3.00, 15.00},
	"claude-haiku-4-20250414":     {0.80, 4.00},
	"claude-opus-4-20250514":      {15.00, 75.00},
	"claude-3.5-sonnet-20241022":  {3.00, 15.00},
	"claude-3.5-haiku-20241022":   {0.80, 4.00},
	// Google
	"gemini-2.5-pro":   {1.25, 10.00},
	"gemini-2.5-flash": {0.15, 0.60},
	"gemini-2.0-flash": {0.10, 0.40},
}

// estimateCost calculates the cost in USD for a given model and token counts.
func estimateCost(model string, tokensIn, tokensOut int) float64 {
	p, ok := pricing[model]
	if !ok {
		return 0
	}
	return (float64(tokensIn)*p[0] + float64(tokensOut)*p[1]) / 1_000_000
}

// LogUsage records a usage entry for a session.
func LogUsage(sessionID, provider, model string, tokensIn, tokensOut int) error {
	cost := estimateCost(model, tokensIn, tokensOut)
	_, err := DB().Exec(
		`INSERT INTO usage_log (session_id, provider, model, tokens_in, tokens_out, cost_usd) VALUES (?, ?, ?, ?, ?, ?)`,
		sessionID, provider, model, tokensIn, tokensOut, cost,
	)
	if err != nil {
		return fmt.Errorf("log usage: %w", err)
	}
	return nil
}

// GetSessionUsage returns the aggregated usage for a single session.
func GetSessionUsage(sessionID string) (UsageSummary, error) {
	var s UsageSummary
	err := DB().QueryRow(
		`SELECT COALESCE(SUM(tokens_in),0), COALESCE(SUM(tokens_out),0), COALESCE(SUM(cost_usd),0)
		 FROM usage_log WHERE session_id = ?`, sessionID,
	).Scan(&s.TokensIn, &s.TokensOut, &s.CostUSD)
	if err != nil {
		return s, fmt.Errorf("get session usage: %w", err)
	}
	return s, nil
}

// GetTotalUsage returns the aggregated usage across all sessions.
func GetTotalUsage() (UsageSummary, error) {
	var s UsageSummary
	err := DB().QueryRow(
		`SELECT COALESCE(SUM(tokens_in),0), COALESCE(SUM(tokens_out),0), COALESCE(SUM(cost_usd),0)
		 FROM usage_log`,
	).Scan(&s.TokensIn, &s.TokensOut, &s.CostUSD)
	if err != nil {
		return s, fmt.Errorf("get total usage: %w", err)
	}
	return s, nil
}
