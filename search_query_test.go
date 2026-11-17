package mal

import (
	"strings"
	"testing"
)

func TestSearchQueryLimitsLongTitlesAtWordBoundary(t *testing.T) {
	title := "Reborn as a Space Mercenary: I Woke Up Piloting the Strongest Starship"
	got := searchQuery(title)
	if got != "Reborn as a Space Mercenary: I Woke Up Piloting the Strongest" {
		t.Fatalf("search query = %q", got)
	}
	if len([]rune(got)) > 64 || !strings.HasPrefix(title, got) {
		t.Fatalf("invalid shortened query %q", got)
	}
}
