package mal

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAnimeService_Search(t *testing.T) {
	fixture, err := os.ReadFile("testdata/anime_search_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/anime" {
			t.Errorf("expected path '/anime', got %q", r.URL.Path)
		}
		if r.URL.Query().Get("q") != "Frieren" {
			t.Errorf("expected query 'Frieren', got %q", r.URL.Query().Get("q"))
		}
		if r.URL.Query().Get("limit") != "1" {
			t.Errorf("expected limit '1', got %q", r.URL.Query().Get("limit"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("test-cid", WithBaseURL(server.URL))
	results, resp, err := client.Anime.Search(context.Background(), "Frieren", WithLimit(1))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if len(results.Data) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results.Data))
	}
	if results.Data[0].Node.ID != 52991 || results.Data[0].Node.Title != "Sousou no Frieren" {
		t.Errorf("unexpected anime node: %+v", results.Data[0].Node)
	}
}

func TestAnimeService_Search_EmptyQuery(t *testing.T) {
	client := NewClient("test-cid")
	_, _, err := client.Anime.Search(context.Background(), "   ")
	if err == nil {
		t.Fatal("expected error on empty query, got nil")
	}
}

func TestAnimeService_GetByID(t *testing.T) {
	fixture, err := os.ReadFile("testdata/anime_details_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/anime/52991" {
			t.Errorf("expected path '/anime/52991', got %q", r.URL.Path)
		}
		if r.URL.Query().Get("fields") != "alternative_titles,genres,num_episodes" {
			t.Errorf("expected fields, got %q", r.URL.Query().Get("fields"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("test-cid", WithBaseURL(server.URL))
	anime, resp, err := client.Anime.GetByID(context.Background(), 52991, "alternative_titles", "genres", "num_episodes")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if anime.ID != 52991 || anime.Title != "Sousou no Frieren" {
		t.Errorf("unexpected anime details: %+v", anime)
	}
	if anime.NumEpisodes != 28 {
		t.Errorf("expected 28 episodes, got %d", anime.NumEpisodes)
	}
	if anime.AlternativeTitles == nil || anime.AlternativeTitles.En != "Frieren: Beyond Journey's End" {
		t.Errorf("expected alternative titles, got %+v", anime.AlternativeTitles)
	}
}

func TestAnimeService_GetByID_InvalidID(t *testing.T) {
	client := NewClient("test-cid")
	_, _, err := client.Anime.GetByID(context.Background(), 0)
	if err == nil {
		t.Fatal("expected error for id 0, got nil")
	}
}

func TestAnimeService_GetRanking(t *testing.T) {
	fixture, err := os.ReadFile("testdata/anime_ranking_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/anime/ranking" {
			t.Errorf("expected path '/anime/ranking', got %q", r.URL.Path)
		}
		if r.URL.Query().Get("ranking_type") != "all" {
			t.Errorf("expected ranking_type 'all', got %q", r.URL.Query().Get("ranking_type"))
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("test-cid", WithBaseURL(server.URL))
	rankings, _, err := client.Anime.GetRanking(context.Background(), "all", WithLimit(1))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(rankings.Data) != 1 || rankings.Data[0].Ranking.Rank != 1 {
		t.Fatalf("unexpected ranking result: %+v", rankings.Data)
	}
}

func TestAnimeService_GetSeasonal(t *testing.T) {
	fixture, err := os.ReadFile("testdata/anime_seasonal_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/anime/season/2023/fall" {
			t.Errorf("expected path '/anime/season/2023/fall', got %q", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("test-cid", WithBaseURL(server.URL))
	seasonal, _, err := client.Anime.GetSeasonal(context.Background(), 2023, "fall")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(seasonal.Data) != 1 || seasonal.Data[0].Node.ID != 52991 {
		t.Fatalf("unexpected seasonal result: %+v", seasonal.Data)
	}
}

func TestAnimeService_GetSeasonal_EmptySeason(t *testing.T) {
	client := NewClient("test-cid")
	_, _, err := client.Anime.GetSeasonal(context.Background(), 2023, "")
	if err == nil {
		t.Fatal("expected error on empty season, got nil")
	}
}
