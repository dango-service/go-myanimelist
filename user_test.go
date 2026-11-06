package mal

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestUserService_GetMyProfile(t *testing.T) {
	fixture, err := os.ReadFile("testdata/user_profile_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/@me" {
			t.Errorf("expected path '/users/@me', got %q", r.URL.Path)
		}
		if r.URL.Query().Get("fields") != "anime_statistics,time_zone" {
			t.Errorf("expected fields, got %q", r.URL.Query().Get("fields"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("test-cid", WithToken("user-token"), WithBaseURL(server.URL))
	profile, resp, err := client.User.GetMyProfile(context.Background(), "anime_statistics", "time_zone")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if profile.ID != 1234567 || profile.Name != "hicham" {
		t.Errorf("unexpected profile data: %+v", profile)
	}
	if profile.AnimeStatistics == nil || profile.AnimeStatistics.NumEpisodes != 2650 {
		t.Errorf("unexpected anime statistics: %+v", profile.AnimeStatistics)
	}
}
