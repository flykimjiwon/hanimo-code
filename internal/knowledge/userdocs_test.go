package knowledge

import "testing"

func TestScanAndSearch(t *testing.T) {
	idx := ScanUserDocs()
	if idx.Count() == 0 {
		t.Skip("no .hanimo/knowledge/ dir with docs in test env")
	}
	t.Logf("indexed %d docs from %s", idx.Count(), idx.Root())

	toc := idx.TableOfContents()
	if toc == "" {
		t.Error("expected non-empty TOC")
	}
	t.Logf("TOC:\n%s", toc)

	results := idx.Search("svelte runes", 3)
	if len(results) == 0 {
		t.Error("expected at least 1 result for 'svelte runes'")
	}
	for _, r := range results {
		t.Logf("found: %s — %s", r.Path, r.Title)
	}

	results2 := idx.Search("엔드포인트 API", 3)
	if len(results2) == 0 {
		t.Error("expected result for '엔드포인트 API'")
	}

	results3 := idx.Search("xxxnonexistentxxx", 3)
	if len(results3) != 0 {
		t.Errorf("expected 0 results for nonsense query, got %d", len(results3))
	}
}
