package mal

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMangaService_Search(t *testing.T) {
	fixture, err := os.ReadFile("testdata/manga_search_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/manga" {
			t.Errorf("expected path '/manga', got %q", r.URL.Path)
		}
		if r.URL.Query().Get("q") != "Frieren" {
			t.Errorf("expected query 'Frieren', got %q", r.URL.Query().Get("q"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("test-cid", WithBaseURL(server.URL))
	results, resp, err := client.Manga.Search(context.Background(), "Frieren")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if len(results.Data) != 1 || results.Data[0].Node.ID != 126479 {
		t.Errorf("unexpected manga results: %+v", results.Data)
	}
}

func TestMangaService_Search_EmptyQuery(t *testing.T) {
	client := NewClient("test-cid")
	_, _, err := client.Manga.Search(context.Background(), "   ")
	if err == nil {
		t.Fatal("expected error on empty query, got nil")
	}
}

func TestMangaService_GetByID(t *testing.T) {
	fixture, err := os.ReadFile("testdata/manga_details_success.json")
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/manga/126479" {
			t.Errorf("expected path '/manga/126479', got %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fixture)
	}))
	defer server.Close()

	client := NewClient("test-cid", WithBaseURL(server.URL))
	manga, resp, err := client.Manga.GetByID(context.Background(), 126479)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if manga.ID != 126479 || manga.Title != "Sousou no Frieren" {
		t.Errorf("unexpected manga details: %+v", manga)
	}
	if manga.NumChapters != 130 || manga.NumVolumes != 13 {
		t.Errorf("unexpected manga volume/chapters: %+v", manga)
	}
}

func TestMangaService_GetByID_InvalidID(t *testing.T) {
	client := NewClient("test-cid")
	_, _, err := client.Manga.GetByID(context.Background(), 0)
	if err == nil {
		t.Fatal("expected error on invalid id, got nil")
	}
}

func TestMangaService_UpdateAndDeleteStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/manga/126479/my_list_status" {
			t.Errorf("expected path '/manga/126479/my_list_status', got %q", r.URL.Path)
		}
		if r.Method == http.MethodPatch {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"reading","score":10,"num_volumes_read":5,"num_chapters_read":50}`))
			return
		}
		if r.Method == http.MethodDelete {
			w.WriteHeader(http.StatusOK)
			return
		}
		t.Errorf("unexpected method %s", r.Method)
	}))
	defer server.Close()

	client := NewClient("test-cid", WithToken("user-token"), WithBaseURL(server.URL))
	status := MangaStatusReading
	score := 10
	vols := 5
	chaps := 50
	updated, _, err := client.Manga.UpdateStatus(context.Background(), 126479, &UpdateMangaStatusInput{
		Status:          &status,
		Score:           &score,
		NumVolumesRead:  &vols,
		NumChaptersRead: &chaps,
	})
	if err != nil {
		t.Fatalf("unexpected update error: %v", err)
	}
	if updated.Status != "reading" || updated.Score != 10 {
		t.Errorf("unexpected updated status: %+v", updated)
	}

	delResp, err := client.Manga.DeleteStatus(context.Background(), 126479)
	if err != nil {
		t.Fatalf("unexpected delete error: %v", err)
	}
	if delResp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", delResp.StatusCode)
	}
}

func TestMangaService_NilInputAndInvalidIDs(t *testing.T) {
	client := NewClient("test-cid")
	if _, _, err := client.Manga.UpdateStatus(context.Background(), 126479, nil); err == nil {
		t.Fatal("expected error on nil input, got nil")
	}
	if _, _, err := client.Manga.UpdateStatus(context.Background(), 0, &UpdateMangaStatusInput{}); err == nil {
		t.Fatal("expected error on invalid id, got nil")
	}
	if _, err := client.Manga.DeleteStatus(context.Background(), 0); err == nil {
		t.Fatal("expected error on invalid delete id, got nil")
	}
}
