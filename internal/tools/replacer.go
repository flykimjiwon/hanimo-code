package tools

import (
	"fmt"
	"math"
	"strings"
)

// SmartReplace tries multiple strategies to find and replace oldStr in content.
// Returns (newContent, strategyUsed, error).
func SmartReplace(content, oldStr, newStr string) (string, string, error) {
	strategies := []struct {
		name string
		fn   func(string, string, string) (string, bool)
	}{
		{"exact", strategyExact},
		{"trimmed_lines", strategyTrimmedLines},
		{"whitespace_normalized", strategyWhitespaceNormalized},
		{"indentation_flexible", strategyIndentationFlexible},
		{"block_anchor", strategyBlockAnchor},
	}

	for _, s := range strategies {
		result, ok := s.fn(content, oldStr, newStr)
		if ok {
			return result, s.name, nil
		}
	}

	// No strategy worked — find closest match for helpful error
	closest := findClosestMatch(content, oldStr)
	if closest != "" {
		return "", "", fmt.Errorf("no match found. Closest content:\n%s", closest)
	}
	return "", "", fmt.Errorf("no match found for the given old_string")
}

// strategyExact performs a direct string match (current behavior).
func strategyExact(content, oldStr, newStr string) (string, bool) {
	count := strings.Count(content, oldStr)
	if count == 1 {
		return strings.Replace(content, oldStr, newStr, 1), true
	}
	return "", false
}

// strategyTrimmedLines trims whitespace from each line, then matches.
func strategyTrimmedLines(content, oldStr, newStr string) (string, bool) {
	contentLines := strings.Split(content, "\n")
	oldLines := strings.Split(oldStr, "\n")

	trimmedOld := make([]string, len(oldLines))
	for i, l := range oldLines {
		trimmedOld[i] = strings.TrimSpace(l)
	}

	var matchIndices []int
	for i := 0; i <= len(contentLines)-len(oldLines); i++ {
		match := true
		for j, tl := range trimmedOld {
			if strings.TrimSpace(contentLines[i+j]) != tl {
				match = false
				break
			}
		}
		if match {
			matchIndices = append(matchIndices, i)
		}
	}

	if len(matchIndices) != 1 {
		return "", false
	}

	idx := matchIndices[0]
	var result []string
	result = append(result, contentLines[:idx]...)
	result = append(result, strings.Split(newStr, "\n")...)
	result = append(result, contentLines[idx+len(oldLines):]...)
	return strings.Join(result, "\n"), true
}

// strategyWhitespaceNormalized collapses all whitespace runs to single space.
func strategyWhitespaceNormalized(content, oldStr, newStr string) (string, bool) {
	normalize := func(s string) string {
		fields := strings.Fields(s)
		return strings.Join(fields, " ")
	}

	normContent := normalize(content)
	normOld := normalize(oldStr)

	count := strings.Count(normContent, normOld)
	if count != 1 {
		return "", false
	}

	// Find the match position in normalized space, then map back
	// We need to find the original substring boundaries
	contentLines := strings.Split(content, "\n")
	oldLines := strings.Split(oldStr, "\n")

	// Use line-based matching with normalized comparison
	normOldLines := make([]string, len(oldLines))
	for i, l := range oldLines {
		normOldLines[i] = normalize(l)
	}

	var matchIndices []int
	for i := 0; i <= len(contentLines)-len(oldLines); i++ {
		match := true
		for j, nl := range normOldLines {
			if normalize(contentLines[i+j]) != nl {
				match = false
				break
			}
		}
		if match {
			matchIndices = append(matchIndices, i)
		}
	}

	if len(matchIndices) != 1 {
		return "", false
	}

	idx := matchIndices[0]
	var result []string
	result = append(result, contentLines[:idx]...)
	result = append(result, strings.Split(newStr, "\n")...)
	result = append(result, contentLines[idx+len(oldLines):]...)
	return strings.Join(result, "\n"), true
}

// strategyIndentationFlexible strips minimum indentation from oldStr,
// then searches for the same content at any indentation level.
func strategyIndentationFlexible(content, oldStr, newStr string) (string, bool) {
	oldLines := strings.Split(oldStr, "\n")

	// Find minimum indentation in oldStr (ignoring blank lines)
	minIndent := math.MaxInt32
	for _, l := range oldLines {
		if strings.TrimSpace(l) == "" {
			continue
		}
		indent := len(l) - len(strings.TrimLeft(l, " \t"))
		if indent < minIndent {
			minIndent = indent
		}
	}
	if minIndent == math.MaxInt32 {
		minIndent = 0
	}

	// Strip minimum indentation from old lines
	strippedOld := make([]string, len(oldLines))
	for i, l := range oldLines {
		if len(l) >= minIndent {
			strippedOld[i] = l[minIndent:]
		} else {
			strippedOld[i] = strings.TrimLeft(l, " \t")
		}
	}

	contentLines := strings.Split(content, "\n")

	var matchIndices []int
	for i := 0; i <= len(contentLines)-len(oldLines); i++ {
		// Determine indentation of first non-blank content line
		contentIndent := 0
		for _, cl := range contentLines[i : i+len(oldLines)] {
			if strings.TrimSpace(cl) != "" {
				contentIndent = len(cl) - len(strings.TrimLeft(cl, " \t"))
				break
			}
		}

		match := true
		for j, sl := range strippedOld {
			cl := contentLines[i+j]
			// Strip the content's indentation prefix and compare
			var clStripped string
			if strings.TrimSpace(cl) == "" && strings.TrimSpace(sl) == "" {
				continue // both blank
			}
			if len(cl) >= contentIndent {
				clStripped = cl[contentIndent:]
			} else {
				clStripped = strings.TrimLeft(cl, " \t")
			}
			if clStripped != sl {
				match = false
				break
			}
		}
		if match {
			matchIndices = append(matchIndices, i)
		}
	}

	if len(matchIndices) != 1 {
		return "", false
	}

	idx := matchIndices[0]
	var result []string
	result = append(result, contentLines[:idx]...)
	result = append(result, strings.Split(newStr, "\n")...)
	result = append(result, contentLines[idx+len(oldLines):]...)
	return strings.Join(result, "\n"), true
}

// strategyBlockAnchor matches first and last lines exactly,
// accepts middle if >70% of lines are similar.
func strategyBlockAnchor(content, oldStr, newStr string) (string, bool) {
	oldLines := strings.Split(oldStr, "\n")
	if len(oldLines) < 3 {
		return "", false // need at least 3 lines for anchor strategy
	}

	firstLine := strings.TrimSpace(oldLines[0])
	lastLine := strings.TrimSpace(oldLines[len(oldLines)-1])
	contentLines := strings.Split(content, "\n")

	var matchIndices []int
	for i := 0; i <= len(contentLines)-len(oldLines); i++ {
		if strings.TrimSpace(contentLines[i]) != firstLine {
			continue
		}
		endIdx := i + len(oldLines) - 1
		if endIdx >= len(contentLines) {
			continue
		}
		if strings.TrimSpace(contentLines[endIdx]) != lastLine {
			continue
		}

		// Check middle lines similarity
		middleTotal := len(oldLines) - 2
		if middleTotal == 0 {
			matchIndices = append(matchIndices, i)
			continue
		}

		similarCount := 0
		for j := 1; j < len(oldLines)-1; j++ {
			if lineSimilarity(oldLines[j], contentLines[i+j]) >= 0.7 {
				similarCount++
			}
		}

		ratio := float64(similarCount) / float64(middleTotal)
		if ratio > 0.7 {
			matchIndices = append(matchIndices, i)
		}
	}

	if len(matchIndices) != 1 {
		return "", false
	}

	idx := matchIndices[0]
	var result []string
	result = append(result, contentLines[:idx]...)
	result = append(result, strings.Split(newStr, "\n")...)
	result = append(result, contentLines[idx+len(oldLines):]...)
	return strings.Join(result, "\n"), true
}

// lineSimilarity returns a simple ratio of matching characters between two lines.
func lineSimilarity(a, b string) float64 {
	a = strings.TrimSpace(a)
	b = strings.TrimSpace(b)
	if a == b {
		return 1.0
	}
	if len(a) == 0 && len(b) == 0 {
		return 1.0
	}
	if len(a) == 0 || len(b) == 0 {
		return 0.0
	}

	// Simple longest common subsequence ratio
	shorter, longer := a, b
	if len(a) > len(b) {
		shorter, longer = b, a
	}

	matchCount := 0
	longerIdx := 0
	for i := 0; i < len(shorter) && longerIdx < len(longer); i++ {
		for longerIdx < len(longer) {
			if shorter[i] == longer[longerIdx] {
				matchCount++
				longerIdx++
				break
			}
			longerIdx++
		}
	}

	return float64(matchCount) / float64(len(longer))
}

// findClosestMatch finds the most similar block in content to oldStr.
func findClosestMatch(content, oldStr string) string {
	oldLines := strings.Split(oldStr, "\n")
	contentLines := strings.Split(content, "\n")

	if len(oldLines) == 0 || len(contentLines) == 0 {
		return ""
	}

	bestScore := 0.0
	bestIdx := -1

	blockLen := len(oldLines)
	if blockLen > len(contentLines) {
		blockLen = len(contentLines)
	}

	for i := 0; i <= len(contentLines)-blockLen; i++ {
		score := 0.0
		for j := 0; j < blockLen; j++ {
			score += lineSimilarity(oldLines[j], contentLines[i+j])
		}
		score /= float64(blockLen)

		if score > bestScore {
			bestScore = score
			bestIdx = i
		}
	}

	if bestIdx < 0 || bestScore < 0.3 {
		return ""
	}

	end := bestIdx + blockLen
	if end > len(contentLines) {
		end = len(contentLines)
	}
	return strings.Join(contentLines[bestIdx:end], "\n")
}
