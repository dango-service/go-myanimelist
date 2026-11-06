package mal

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func TestAnimeService_GetUserAnimeList(t *testing.T) {
	fixture, err := os.ReadFile("testdata/user_animelist_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/hicham/animelist" {
			t.Errorf("expected path '/users/hicham/animelist', got %q", r.URL.Path)
		}
		if r.URL.Query().Get("status") != AnimeStatusWatching {
			t.Errorf("expected status 'watching', got %q", r.URL.Query().Get("status"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("test-cid", WithToken("user-token"), WithBaseURL(server.URL))
	list, resp, err := client.Anime.GetUserAnimeList(context.Background(), "hicham", func(v *url.Values) {
		v.Set("status", AnimeStatusWatching)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if len(list.Data) != 1 {
		t.Fatalf("expected 1 item, got %d", len(list.Data))
	}
	if list.Data[0].ListStatus.Status != "watching" || list.Data[0].ListStatus.NumEpisodesWatched != 14 {
		t.Errorf("unexpected list status: %+v", list.Data[0].ListStatus)
	}
}

func TestAnimeService_GetUserAnimeList_DefaultMe(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/@me/animelist" {
			t.Errorf("expected path '/users/@me/animelist', got %q", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":[]}`))
	}))
	defer server.Close()

	client := NewClient("test-cid", WithBaseURL(server.URL))
	_, _, err := client.Anime.GetUserAnimeList(context.Background(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAnimeService_UpdateStatus(t *testing.T) {
	fixture, err := os.ReadFile("testdata/update_status_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected method PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/anime/52991/my_list_status" {
			t.Errorf("expected path '/anime/52991/my_list_status', got %q", r.URL.Path)
		}
		bodyBytes, _ := io.ReadAll(r.Body)
		parsed, _ := url.ParseQuery(string(bodyBytes))
		if parsed.Get("status") != AnimeStatusCompleted {
			t.Errorf("expected status 'completed', got %q", parsed.Get("status"))
		}
		if parsed.Get("score") != "9" {
			t.Errorf("expected score '9', got %q", parsed.Get("score"))
		}
		if parsed.Get("num_watched_episodes") != "28" {
			t.Errorf("expected episodes '28', got %q", parsed.Get("num_watched_episodes"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	status := AnimeStatusCompleted
	score := 9
	episodes := 28
	comments := "Incredible adaptation"
	input := &UpdateStatusInput{
		Status:             &status,
		Score:              &score,
		NumWatchedEpisodes: &episodes,
		Comments:           &comments,
		Tags:               []string{"masterpiece", "fantasy"},
	}

	client := NewClient("test-cid", WithToken("user-token"), WithBaseURL(server.URL))
	res, resp, err := client.Anime.UpdateStatus(context.Background(), 52991, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if res.Status != "completed" || res.Score != 9 || res.NumEpisodesWatched != 28 {
		t.Errorf("unexpected updated status: %+v", res)
	}
}

func TestAnimeService_UpdateStatus_NilInput(t *testing.T) {
	client := NewClient("test-cid")
	_, _, err := client.Anime.UpdateStatus(context.Background(), 52991, nil)
	if err == nil {
		t.Fatal("expected error for nil input, got nil")
	}
}

func TestAnimeService_UpdateStatus_InvalidID(t *testing.T) {
	client := NewClient("test-cid")
	status := AnimeStatusWatching
	_, _, err := client.Anime.UpdateStatus(context.Background(), 0, &UpdateStatusInput{Status: &status})
	if err == nil {
		t.Fatal("expected error for id 0, got nil")
	}
}

func TestAnimeService_DeleteStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected method DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/anime/52991/my_list_status" {
			t.Errorf("expected path '/anime/52991/my_list_status', got %q", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-cid", WithToken("user-token"), WithBaseURL(server.URL))
	resp, err := client.Anime.DeleteStatus(context.Background(), 52991)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestAnimeService_DeleteStatus_InvalidID(t *testing.T) {
	client := NewClient("test-cid")
	_, err := client.Anime.DeleteStatus(context.Background(), -1)
	if err == nil {
		t.Fatal("expected error for invalid id, got nil")
	}
}
