package mal

import (
	"net/http"
	"testing"
	"time"
)

func TestNewClient_Defaults(t *testing.T) {
	c := NewClient("test-client-id")

	if c.clientID != "test-client-id" {
		t.Fatalf("expected clientID 'test-client-id', got %q", c.clientID)
	}
	if c.baseURL != DefaultBaseURL {
		t.Fatalf("expected baseURL %q, got %q", DefaultBaseURL, c.baseURL)
	}
	if c.userAgent != DefaultUserAgent {
		t.Fatalf("expected userAgent %q, got %q", DefaultUserAgent, c.userAgent)
	}
	if c.maxBodyBytes != DefaultMaxBodyBytes {
		t.Fatalf("expected maxBodyBytes %d, got %d", DefaultMaxBodyBytes, c.maxBodyBytes)
	}
	if c.Anime == nil || c.Manga == nil || c.User == nil || c.Auth == nil {
		t.Fatal("expected all sub-services to be initialized")
	}
}

func TestNewClient_Options(t *testing.T) {
	customHTTP := &http.Client{Timeout: 5 * time.Second}
	c := NewClient("init-id",
		WithBaseURL("https://custom.mal.net/api/"),
		WithClientID("override-id"),
		WithToken("user-token-123"),
		WithUserAgent("custom-agent/2.0"),
		WithMaxBodyBytes(2*1024*1024),
		WithHTTPClient(customHTTP),
	)

	if c.clientID != "override-id" {
		t.Fatalf("expected clientID 'override-id', got %q", c.clientID)
	}
	if c.baseURL != "https://custom.mal.net/api" {
		t.Fatalf("expected trimmed baseURL 'https://custom.mal.net/api', got %q", c.baseURL)
	}
	if c.token != "user-token-123" {
		t.Fatalf("expected token 'user-token-123', got %q", c.token)
	}
	if c.userAgent != "custom-agent/2.0" {
		t.Fatalf("expected userAgent 'custom-agent/2.0', got %q", c.userAgent)
	}
	if c.maxBodyBytes != 2*1024*1024 {
		t.Fatalf("expected maxBodyBytes 2097152, got %d", c.maxBodyBytes)
	}
	if c.httpClient != customHTTP {
		t.Fatal("expected custom http client to be set")
	}
}
