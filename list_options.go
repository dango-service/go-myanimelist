package mal

import (
	"net/url"
	"strconv"
	"strings"
)

// ListOption configures query parameters for list endpoints.
type ListOption func(*url.Values)

// WithLimit sets the maximum number of results to return (max 100 for most endpoints).
func WithLimit(limit int) ListOption {
	return func(v *url.Values) {
		if limit > 0 {
			v.Set("limit", strconv.Itoa(limit))
		}
	}
}

// WithOffset sets the pagination offset for results.
func WithOffset(offset int) ListOption {
	return func(v *url.Values) {
		if offset >= 0 {
			v.Set("offset", strconv.Itoa(offset))
		}
	}
}

// WithFields specifies which optional fields to include in the response.
func WithFields(fields ...string) ListOption {
	return func(v *url.Values) {
		if len(fields) > 0 {
			v.Set("fields", strings.Join(fields, ","))
		}
	}
}

// WithSort sets the sorting order for list endpoints that support it.
func WithSort(sort string) ListOption {
	return func(v *url.Values) {
		trimmed := strings.TrimSpace(sort)
		if trimmed != "" {
			v.Set("sort", trimmed)
		}
	}
}

// WithNSFW specifies whether to include NSFW results.
func WithNSFW(includeNSFW bool) ListOption {
	return func(v *url.Values) {
		if includeNSFW {
			v.Set("nsfw", "true")
		} else {
			v.Set("nsfw", "false")
		}
	}
}
