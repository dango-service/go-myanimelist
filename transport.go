package mal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (c *Client) doRequest(ctx context.Context, method, path string, body io.Reader, contentType string, v any) (*Response, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("mal: context error: %w", err)
	}

	urlStr := path
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		urlStr = c.baseURL + "/" + strings.TrimLeft(path, "/")
	}

	req, err := http.NewRequestWithContext(ctx, method, urlStr, body)
	if err != nil {
		return nil, fmt.Errorf("mal: create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	// Dual authentication strategy:
	// Prioritize Bearer token for user actions, fall back to X-MAL-CLIENT-ID for public queries.
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	} else if c.clientID != "" {
		req.Header.Set("X-MAL-CLIENT-ID", c.clientID)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, fmt.Errorf("mal: request canceled: %w", ctxErr)
		}
		return nil, fmt.Errorf("mal: execute request: %w", err)
	}
	defer resp.Body.Close()

	rateLimit := parseRateLimit(resp.Header)
	response := &Response{
		Response:  resp,
		RateLimit: rateLimit,
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return response, c.handleErrorResponse(resp, rateLimit)
	}

	if v != nil {
		limitedReader := io.LimitReader(resp.Body, c.maxBodyBytes)
		if err := json.NewDecoder(limitedReader).Decode(v); err != nil && err != io.EOF {
			return response, fmt.Errorf("mal: decode response: %w", err)
		}
	}

	return response, nil
}

func (c *Client) handleErrorResponse(resp *http.Response, rateLimit RateLimitInfo) error {
	limitedReader := io.LimitReader(resp.Body, c.maxBodyBytes)
	errResp := &ErrorResponse{
		StatusCode: resp.StatusCode,
		RateLimit:  rateLimit,
	}

	_ = json.NewDecoder(limitedReader).Decode(errResp)

	switch resp.StatusCode {
	case http.StatusTooManyRequests:
		errResp.Err = ErrRateLimited
	case http.StatusUnauthorized:
		errResp.Err = ErrUnauthorized
	case http.StatusForbidden:
		errResp.Err = ErrForbidden
	case http.StatusNotFound:
		errResp.Err = ErrNotFound
	case http.StatusBadRequest:
		errResp.Err = ErrBadRequest
	case http.StatusConflict:
		errResp.Err = ErrConflict
	default:
		errResp.Err = fmt.Errorf("mal: HTTP %d", resp.StatusCode)
	}

	return errResp
}
