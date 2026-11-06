package mal

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strings"
)

const (
	// DefaultAuthEndpoint is the standard OAuth2 authorization endpoint for MyAnimeList.
	DefaultAuthEndpoint = "https://myanimelist.net/v1/oauth2/authorize"

	// DefaultTokenEndpoint is the standard OAuth2 token endpoint for MyAnimeList.
	DefaultTokenEndpoint = "https://myanimelist.net/v1/oauth2/token"

	// pkceCharset contains valid characters for PKCE code verifier (RFC 7636).
	pkceCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._~"
)

// Token represents OAuth2 credentials returned by MyAnimeList.
type Token struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// String masks sensitive token credentials to prevent accidental logging.
func (t Token) String() string {
	return fmt.Sprintf("mal.Token{TokenType:%q, ExpiresIn:%d, AccessToken:[REDACTED], RefreshToken:[REDACTED]}", t.TokenType, t.ExpiresIn)
}

// GoString masks sensitive token credentials when formatted with %#v.
func (t Token) GoString() string {
	return t.String()
}

// GenerateCodeVerifier generates a cryptographically secure 128-character PKCE code verifier.
func (s *AuthService) GenerateCodeVerifier() (string, error) {
	bytes := make([]byte, 128)
	charsetLen := big.NewInt(int64(len(pkceCharset)))

	for i := range bytes {
		num, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", fmt.Errorf("mal: generate code verifier: %w", err)
		}
		bytes[i] = pkceCharset[num.Int64()]
	}

	return string(bytes), nil
}

// BuildAuthorizeURL generates the MyAnimeList OAuth2 authorization URL with PKCE parameters.
func (s *AuthService) BuildAuthorizeURL(redirectURI, state, codeChallenge string) string {
	v := url.Values{}
	v.Set("response_type", "code")
	v.Set("client_id", s.client.clientID)
	v.Set("code_challenge", codeChallenge)
	v.Set("code_challenge_method", "plain")

	if redirectURI != "" {
		v.Set("redirect_uri", redirectURI)
	}
	if state != "" {
		v.Set("state", state)
	}

	return DefaultAuthEndpoint + "?" + v.Encode()
}

// ExchangeCode exchanges an authorization code and PKCE code_verifier for access and refresh tokens.
func (s *AuthService) ExchangeCode(ctx context.Context, clientSecret, code, codeVerifier, redirectURI string) (*Token, *Response, error) {
	if code == "" {
		return nil, nil, fmt.Errorf("mal: authorization code cannot be empty")
	}
	if codeVerifier == "" {
		return nil, nil, fmt.Errorf("mal: code verifier cannot be empty")
	}

	values := url.Values{}
	values.Set("client_id", s.client.clientID)
	if clientSecret != "" {
		values.Set("client_secret", clientSecret)
	}
	values.Set("code", code)
	values.Set("code_verifier", codeVerifier)
	values.Set("grant_type", "authorization_code")
	if redirectURI != "" {
		values.Set("redirect_uri", redirectURI)
	}

	body := strings.NewReader(values.Encode())
	var token Token
	resp, err := s.client.doRequest(ctx, http.MethodPost, s.client.tokenURL, body, "application/x-www-form-urlencoded", &token)
	if err != nil {
		return nil, resp, fmt.Errorf("mal: exchange code: %w", err)
	}

	return &token, resp, nil
}

// RefreshToken obtains a new access token using an existing refresh token.
func (s *AuthService) RefreshToken(ctx context.Context, clientSecret, refreshToken string) (*Token, *Response, error) {
	if refreshToken == "" {
		return nil, nil, fmt.Errorf("mal: refresh token cannot be empty")
	}

	values := url.Values{}
	values.Set("client_id", s.client.clientID)
	if clientSecret != "" {
		values.Set("client_secret", clientSecret)
	}
	values.Set("grant_type", "refresh_token")
	values.Set("refresh_token", refreshToken)

	body := strings.NewReader(values.Encode())
	var token Token
	resp, err := s.client.doRequest(ctx, http.MethodPost, s.client.tokenURL, body, "application/x-www-form-urlencoded", &token)
	if err != nil {
		return nil, resp, fmt.Errorf("mal: refresh token: %w", err)
	}

	return &token, resp, nil
}
