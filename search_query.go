package mal

import "strings"

// searchQuery keeps a catalog search within MAL's accepted q length. The
// caller should still compare returned titles against the unabridged title.
func searchQuery(raw string) string {
	query := strings.Join(strings.Fields(raw), " ")
	chars := []rune(query)
	if len(chars) <= 64 {
		return query
	}
	cut := 64
	for i := 63; i >= 32; i-- {
		if chars[i] == ' ' {
			cut = i
			break
		}
	}
	return strings.TrimSpace(string(chars[:cut]))
}
