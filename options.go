package mal

import (
	"net/http"
	"strings"
)

// ClientOption represents a functional option for configuring a Client.
type ClientOption func(*Client)

// WithHTTPClient overrides the default HTTP client used for requests.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		if httpClient != nil {
			c.httpClient = httpClient
		}
	}
}

// WithBaseURL overrides the default base API URL (default: "https://api.myanimelist.net/v2").
// If a custom base URL is supplied (e.g. in test servers), tokenURL is automatically adjusted.
func WithBaseURL(rawURL string) ClientOption {
	return func(c *Client) {
		trimmed := strings.TrimRight(rawURL, "/")
		if trimmed != "" {
			c.baseURL = trimmed
			c.tokenURL = trimmed + "/v1/oauth2/token"
		}
	}
}

// WithTokenURL overrides the OAuth2 token endpoint (default: "https://myanimelist.net/v1/oauth2/token").
func WithTokenURL(tokenURL string) ClientOption {
	return func(c *Client) {
		trimmed := strings.TrimRight(tokenURL, "/")
		if trimmed != "" {
			c.tokenURL = trimmed
		}
	}
}

// WithClientID sets or overrides the public MyAnimeList client ID.
func WithClientID(clientID string) ClientOption {
	return func(c *Client) {
		c.clientID = clientID
	}
}

// WithToken sets the user OAuth2 access token for authenticated requests.
func WithToken(token string) ClientOption {
	return func(c *Client) {
		c.token = strings.TrimSpace(token)
	}
}

// WithUserAgent sets a custom User-Agent header for all outgoing requests.
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) {
		c.userAgent = strings.TrimSpace(userAgent)
	}
}

// WithMaxBodyBytes overrides the maximum allowed response body size (default: 10 MiB).
func WithMaxBodyBytes(maxBytes int64) ClientOption {
	return func(c *Client) {
		if maxBytes > 0 {
			c.maxBodyBytes = maxBytes
		}
	}
}
