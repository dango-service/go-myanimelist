package mal

import (
	"errors"
	"fmt"
)

var (
	// ErrRateLimited is returned when MyAnimeList throttles requests with HTTP 429.
	ErrRateLimited = errors.New("mal: rate limit exceeded (429)")

	// ErrUnauthorized is returned when authentication is missing or invalid (HTTP 401).
	ErrUnauthorized = errors.New("mal: unauthorized (401)")

	// ErrForbidden is returned when access is forbidden (HTTP 403).
	ErrForbidden = errors.New("mal: forbidden (403)")

	// ErrNotFound is returned when the requested resource does not exist (HTTP 404).
	ErrNotFound = errors.New("mal: resource not found (404)")

	// ErrBadRequest is returned when request parameters are invalid (HTTP 400).
	ErrBadRequest = errors.New("mal: bad request (400)")

	// ErrConflict is returned when a state conflict occurs (HTTP 409).
	ErrConflict = errors.New("mal: conflict (409)")
)

// ErrorResponse represents a structured error returned by the MyAnimeList API.
type ErrorResponse struct {
	StatusCode int           `json:"-"`
	ErrorType  string        `json:"error"`
	Message    string        `json:"message"`
	RateLimit  RateLimitInfo `json:"-"`
	Err        error         `json:"-"`
}

func (e *ErrorResponse) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("mal: HTTP %d: %s (error: %s)", e.StatusCode, e.Message, e.ErrorType)
	}
	if e.ErrorType != "" {
		return fmt.Sprintf("mal: HTTP %d: %s", e.StatusCode, e.ErrorType)
	}
	return fmt.Sprintf("mal: HTTP %d", e.StatusCode)
}

// Unwrap supports errors.Is and errors.As against the sentinel error.
func (e *ErrorResponse) Unwrap() error {
	return e.Err
}

// IsRateLimited checks if an error was caused by HTTP 429 Too Many Requests.
func IsRateLimited(err error) bool {
	return errors.Is(err, ErrRateLimited)
}

// IsUnauthorized checks if an error was caused by HTTP 401 Unauthorized.
func IsUnauthorized(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}

// IsNotFound checks if an error was caused by HTTP 404 Not Found.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}
