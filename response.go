package mal

import (
	"net/http"
	"strconv"
	"time"
)

// RateLimitInfo contains rate limit information extracted from HTTP response headers.
type RateLimitInfo struct {
	// RetryAfter represents the duration the client should wait before making subsequent requests.
	RetryAfter time.Duration
}

// Response wraps http.Response and includes rate limit details.
type Response struct {
	*http.Response

	// RateLimit holds parsed rate limit metadata, populated on all responses.
	RateLimit RateLimitInfo
}

// parseRateLimit extracts rate limit details from HTTP response headers.
func parseRateLimit(h http.Header) RateLimitInfo {
	info := RateLimitInfo{}

	if val := h.Get("Retry-After"); val != "" {
		if sec, err := strconv.Atoi(val); err == nil && sec >= 0 {
			info.RetryAfter = time.Duration(sec) * time.Second
		} else if parsedTime, err := http.ParseTime(val); err == nil {
			dur := time.Until(parsedTime)
			if dur > 0 {
				info.RetryAfter = dur
			}
		}
	}

	return info
}
