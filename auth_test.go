package mal

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestAuthService_GenerateCodeVerifier(t *testing.T) {
	client := NewClient("test-cid")
	verifier, err := client.Auth.GenerateCodeVerifier()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(verifier) != 128 {
		t.Fatalf("expected verifier length 128, got %d", len(verifier))
	}
	for i := 0; i < len(verifier); i++ {
		if !strings.ContainsRune(pkceCharset, rune(verifier[i])) {
			t.Fatalf("verifier contains invalid char: %c", verifier[i])
		}
	}
}

func TestAuthService_BuildAuthorizeURL(t *testing.T) {
	client := NewClient("my-client-id")
	authURL := client.Auth.BuildAuthorizeURL("https://dango.app/callback", "xyzState", "myCodeChallenge")

	u, err := url.Parse(authURL)
	if err != nil {
		t.Fatalf("failed to parse auth URL: %v", err)
	}

	q := u.Query()
	if q.Get("response_type") != "code" {
		t.Errorf("expected response_type 'code', got %q", q.Get("response_type"))
	}
	if q.Get("client_id") != "my-client-id" {
		t.Errorf("expected client_id 'my-client-id', got %q", q.Get("client_id"))
	}
	if q.Get("code_challenge") != "myCodeChallenge" {
		t.Errorf("expected code_challenge 'myCodeChallenge', got %q", q.Get("code_challenge"))
	}
	if q.Get("code_challenge_method") != "plain" {
		t.Errorf("expected code_challenge_method 'plain', got %q", q.Get("code_challenge_method"))
	}
	if q.Get("redirect_uri") != "https://dango.app/callback" {
		t.Errorf("expected redirect_uri, got %q", q.Get("redirect_uri"))
	}
	if q.Get("state") != "xyzState" {
		t.Errorf("expected state 'xyzState', got %q", q.Get("state"))
	}
}

func TestAuthService_ExchangeCode(t *testing.T) {
	fixture, err := os.ReadFile("testdata/token_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/oauth2/token" {
			t.Errorf("expected path '/v1/oauth2/token', got %q", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}
		bodyBytes, _ := io.ReadAll(r.Body)
		parsed, _ := url.ParseQuery(string(bodyBytes))
		if parsed.Get("grant_type") != "authorization_code" {
			t.Errorf("expected grant_type 'authorization_code', got %q", parsed.Get("grant_type"))
		}
		if parsed.Get("code") != "auth_code_123" {
			t.Errorf("expected code 'auth_code_123', got %q", parsed.Get("code"))
		}
		if parsed.Get("code_verifier") != "verifier_abc" {
			t.Errorf("expected verifier 'verifier_abc', got %q", parsed.Get("code_verifier"))
		}
		if parsed.Get("client_secret") != "secret_key" {
			t.Errorf("expected client_secret, got %q", parsed.Get("client_secret"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("my-client-id", WithBaseURL(server.URL))
	token, resp, err := client.Auth.ExchangeCode(context.Background(), "secret_key", "auth_code_123", "verifier_abc", "https://dango.app/callback")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if token.AccessToken != "mal_access_token_super_secret_12345" {
		t.Errorf("unexpected access token: %q", token.AccessToken)
	}
	if token.TokenType != "Bearer" || token.ExpiresIn != 2678400 {
		t.Errorf("unexpected token metadata: %+v", token)
	}
}

func TestAuthService_ExchangeCode_Validation(t *testing.T) {
	client := NewClient("my-client-id")
	if _, _, err := client.Auth.ExchangeCode(context.Background(), "", "", "verifier", ""); err == nil {
		t.Fatal("expected error on empty code, got nil")
	}
	if _, _, err := client.Auth.ExchangeCode(context.Background(), "", "code", "", ""); err == nil {
		t.Fatal("expected error on empty verifier, got nil")
	}
}

func TestAuthService_RefreshToken(t *testing.T) {
	fixture, err := os.ReadFile("testdata/token_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, _ := io.ReadAll(r.Body)
		parsed, _ := url.ParseQuery(string(bodyBytes))
		if parsed.Get("grant_type") != "refresh_token" {
			t.Errorf("expected grant_type 'refresh_token', got %q", parsed.Get("grant_type"))
		}
		if parsed.Get("refresh_token") != "existing_ref_token" {
			t.Errorf("expected refresh token, got %q", parsed.Get("refresh_token"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("my-client-id", WithBaseURL(server.URL))
	token, _, err := client.Auth.RefreshToken(context.Background(), "secret_key", "existing_ref_token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token.RefreshToken != "mal_refresh_token_super_secret_67890" {
		t.Errorf("unexpected refresh token: %q", token.RefreshToken)
	}
}

func TestAuthService_RefreshToken_Validation(t *testing.T) {
	client := NewClient("my-client-id")
	if _, _, err := client.Auth.RefreshToken(context.Background(), "", ""); err == nil {
		t.Fatal("expected error on empty refresh token, got nil")
	}
}

func TestToken_CredentialRedaction(t *testing.T) {
	tok := Token{
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		AccessToken:  "secret_access_token_123",
		RefreshToken: "secret_refresh_token_456",
	}

	str := tok.String()
	goStr := fmt.Sprintf("%#v", tok)

	if strings.Contains(str, "secret_access_token_123") || strings.Contains(str, "secret_refresh_token_456") {
		t.Fatalf("String() leaked token credentials: %s", str)
	}
	if strings.Contains(goStr, "secret_access_token_123") || strings.Contains(goStr, "secret_refresh_token_456") {
		t.Fatalf("GoString() leaked token credentials: %s", goStr)
	}
	if !strings.Contains(str, "[REDACTED]") {
		t.Fatalf("expected [REDACTED] in string representation: %s", str)
	}
}
