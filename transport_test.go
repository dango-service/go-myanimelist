package mal

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTransport_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-MAL-CLIENT-ID") != "test-cid" {
			t.Errorf("expected X-MAL-CLIENT-ID header 'test-cid', got %q", r.Header.Get("X-MAL-CLIENT-ID"))
		}
		if r.Header.Get("User-Agent") != DefaultUserAgent {
			t.Errorf("expected User-Agent %q, got %q", DefaultUserAgent, r.Header.Get("User-Agent"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"hello world"}`))
	}))
	defer server.Close()

	client := NewClient("test-cid", WithBaseURL(server.URL))
	var out struct {
		Message string `json:"message"`
	}

	resp, err := client.doRequest(context.Background(), http.MethodGet, "/test", nil, "", &out)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if out.Message != "hello world" {
		t.Fatalf("expected output 'hello world', got %q", out.Message)
	}
}

func TestTransport_BearerTokenPrecedence(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer secret-token" {
			t.Errorf("expected Bearer token, got %q", r.Header.Get("Authorization"))
		}
		if r.Header.Get("X-MAL-CLIENT-ID") != "" {
			t.Errorf("expected empty X-MAL-CLIENT-ID when bearer token is present, got %q", r.Header.Get("X-MAL-CLIENT-ID"))
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := NewClient("test-cid", WithToken("secret-token"), WithBaseURL(server.URL))
	_, err := client.doRequest(context.Background(), http.MethodGet, "/test", nil, "", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTransport_RateLimitError_Seconds(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "12")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"error":"too_many_requests","message":"Rate limit exceeded"}`))
	}))
	defer server.Close()

	client := NewClient("test-cid", WithBaseURL(server.URL))
	resp, err := client.doRequest(context.Background(), http.MethodGet, "/limited", nil, "", nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsRateLimited(err) {
		t.Fatalf("expected IsRateLimited(err) to be true, got error: %v", err)
	}
	if resp == nil || resp.RateLimit.RetryAfter != 12*time.Second {
		t.Fatalf("expected RetryAfter 12s, got %v", resp.RateLimit.RetryAfter)
	}
}

func TestTransport_RateLimitError_HTTPDate(t *testing.T) {
	targetTime := time.Now().Add(30 * time.Second).UTC().Truncate(time.Second)
	dateStr := targetTime.Format(http.TimeFormat)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", dateStr)
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"error":"rate_limited"}`))
	}))
	defer server.Close()

	client := NewClient("test-cid", WithBaseURL(server.URL))
	resp, err := client.doRequest(context.Background(), http.MethodGet, "/limited", nil, "", nil)

	if !IsRateLimited(err) {
		t.Fatalf("expected IsRateLimited(err) to be true, got %v", err)
	}
	if resp.RateLimit.RetryAfter <= 0 || resp.RateLimit.RetryAfter > 35*time.Second {
		t.Fatalf("expected RetryAfter around 30s, got %v", resp.RateLimit.RetryAfter)
	}
}

func TestTransport_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"unauthorized","message":"Invalid access token"}`))
	}))
	defer server.Close()

	client := NewClient("test-cid", WithBaseURL(server.URL))
	_, err := client.doRequest(context.Background(), http.MethodGet, "/secure", nil, "", nil)

	if !IsUnauthorized(err) {
		t.Fatalf("expected IsUnauthorized(err) to be true, got %v", err)
	}
}

func TestTransport_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"not_found","message":"Anime not found"}`))
	}))
	defer server.Close()

	client := NewClient("test-cid", WithBaseURL(server.URL))
	_, err := client.doRequest(context.Background(), http.MethodGet, "/missing", nil, "", nil)

	if !IsNotFound(err) {
		t.Fatalf("expected IsNotFound(err) to be true, got %v", err)
	}
}

func TestTransport_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-cid", WithBaseURL(server.URL))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.doRequest(ctx, http.MethodGet, "/test", nil, "", nil)
	if err == nil {
		t.Fatal("expected context error, got nil")
	}
}
