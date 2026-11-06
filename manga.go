package mal

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// MangaStatus represents valid manga reading status values.
const (
	MangaStatusReading    = "reading"
	MangaStatusCompleted  = "completed"
	MangaStatusOnHold     = "on_hold"
	MangaStatusDropped    = "dropped"
	MangaStatusPlanToRead = "plan_to_read"
)

// Manga represents a MyAnimeList manga record.
type Manga struct {
	ID                int                `json:"id"`
	Title             string             `json:"title"`
	MainPicture       *Picture           `json:"main_picture,omitempty"`
	AlternativeTitles *AlternativeTitles `json:"alternative_titles,omitempty"`
	StartDate         string             `json:"start_date,omitempty"`
	EndDate           string             `json:"end_date,omitempty"`
	Synopsis          string             `json:"synopsis,omitempty"`
	Mean              float64            `json:"mean,omitempty"`
	Rank              int                `json:"rank,omitempty"`
	Popularity        int                `json:"popularity,omitempty"`
	NumListUsers      int                `json:"num_list_users,omitempty"`
	NumScoringUsers   int                `json:"num_scoring_users,omitempty"`
	NSFW              string             `json:"nsfw,omitempty"`
	CreatedAt         string             `json:"created_at,omitempty"`
	UpdatedAt         string             `json:"updated_at,omitempty"`
	MediaType         string             `json:"media_type,omitempty"`
	Status            string             `json:"status,omitempty"`
	Genres            []Genre            `json:"genres,omitempty"`
	MyListStatus      *MangaListStatus   `json:"my_list_status,omitempty"`
	NumVolumes        int                `json:"num_volumes,omitempty"`
	NumChapters       int                `json:"num_chapters,omitempty"`
	Pictures          []Picture          `json:"pictures,omitempty"`
}

// MangaListStatus reflects a user's reading status for a manga.
type MangaListStatus struct {
	Status          string    `json:"status"`
	Score           int       `json:"score"`
	NumVolumesRead  int       `json:"num_volumes_read"`
	NumChaptersRead int       `json:"num_chapters_read"`
	IsRereading     bool      `json:"is_rereading"`
	UpdatedAt       time.Time `json:"updated_at,omitempty"`
	Priority        int       `json:"priority,omitempty"`
	Comments        string    `json:"comments,omitempty"`
	Tags            []string  `json:"tags,omitempty"`
	StartDate       string    `json:"start_date,omitempty"`
	FinishDate      string    `json:"finish_date,omitempty"`
}

// MangaEdge wraps a Manga node in paginated list results.
type MangaEdge struct {
	Node Manga `json:"node"`
}

// MangaList represents paginated manga search results.
type MangaList struct {
	Data   []MangaEdge `json:"data"`
	Paging Paging      `json:"paging"`
}

// UserMangaEdge pairs a Manga record with user list status.
type UserMangaEdge struct {
	Node       Manga           `json:"node"`
	ListStatus MangaListStatus `json:"list_status"`
}

// UserMangaList represents paginated user manga reading list.
type UserMangaList struct {
	Data   []UserMangaEdge `json:"data"`
	Paging Paging          `json:"paging"`
}

// UpdateMangaStatusInput specifies mutation fields for a manga list entry.
type UpdateMangaStatusInput struct {
	Status          *string
	Score           *int
	NumVolumesRead  *int
	NumChaptersRead *int
	IsRereading     *bool
	Priority        *int
	Comments        *string
	Tags            []string
	StartDate       *string
	FinishDate      *string
}

// ToURLValues serializes manga status mutation fields into url.Values.
func (in *UpdateMangaStatusInput) ToURLValues() url.Values {
	values := url.Values{}
	if in == nil {
		return values
	}

	if in.Status != nil && *in.Status != "" {
		values.Set("status", *in.Status)
	}
	if in.Score != nil {
		values.Set("score", strconv.Itoa(*in.Score))
	}
	if in.NumVolumesRead != nil {
		values.Set("num_volumes_read", strconv.Itoa(*in.NumVolumesRead))
	}
	if in.NumChaptersRead != nil {
		values.Set("num_chapters_read", strconv.Itoa(*in.NumChaptersRead))
	}
	if in.IsRereading != nil {
		if *in.IsRereading {
			values.Set("is_rereading", "true")
		} else {
			values.Set("is_rereading", "false")
		}
	}
	if in.Priority != nil {
		values.Set("priority", strconv.Itoa(*in.Priority))
	}
	if in.Comments != nil {
		values.Set("comments", *in.Comments)
	}
	if len(in.Tags) > 0 {
		values.Set("tags", strings.Join(in.Tags, ","))
	}
	if in.StartDate != nil && *in.StartDate != "" {
		values.Set("start_date", *in.StartDate)
	}
	if in.FinishDate != nil && *in.FinishDate != "" {
		values.Set("finish_date", *in.FinishDate)
	}

	return values
}

// Search searches for manga by title keywords.
func (s *MangaService) Search(ctx context.Context, query string, opts ...ListOption) (*MangaList, *Response, error) {
	trimmed := strings.TrimSpace(query)
	if trimmed == "" {
		return nil, nil, fmt.Errorf("mal: manga search query cannot be empty")
	}

	values := url.Values{}
	values.Set("q", trimmed)
	for _, opt := range opts {
		if opt != nil {
			opt(&values)
		}
	}

	path := "/manga?" + values.Encode()
	var result MangaList
	resp, err := s.client.doRequest(ctx, http.MethodGet, path, nil, "", &result)
	if err != nil {
		return nil, resp, fmt.Errorf("mal: search manga: %w", err)
	}

	return &result, resp, nil
}

// GetByID retrieves manga details by its ID.
func (s *MangaService) GetByID(ctx context.Context, id int, fields ...string) (*Manga, *Response, error) {
	if id <= 0 {
		return nil, nil, fmt.Errorf("mal: manga id must be greater than zero")
	}

	path := fmt.Sprintf("/manga/%d", id)
	if len(fields) > 0 {
		values := url.Values{}
		values.Set("fields", strings.Join(fields, ","))
		path += "?" + values.Encode()
	}

	var result Manga
	resp, err := s.client.doRequest(ctx, http.MethodGet, path, nil, "", &result)
	if err != nil {
		return nil, resp, fmt.Errorf("mal: get manga by id: %w", err)
	}

	return &result, resp, nil
}

// GetUserMangaList retrieves a user's manga list. Pass "@me" for currently authenticated user.
func (s *MangaService) GetUserMangaList(ctx context.Context, username string, opts ...ListOption) (*UserMangaList, *Response, error) {
	trimmed := strings.TrimSpace(username)
	if trimmed == "" {
		trimmed = "@me"
	}

	values := url.Values{}
	values.Set("fields", "list_status")
	for _, opt := range opts {
		if opt != nil {
			opt(&values)
		}
	}

	path := fmt.Sprintf("/users/%s/mangalist?%s", url.PathEscape(trimmed), values.Encode())
	var result UserMangaList
	resp, err := s.client.doRequest(ctx, http.MethodGet, path, nil, "", &result)
	if err != nil {
		return nil, resp, fmt.Errorf("mal: get user mangalist: %w", err)
	}

	return &result, resp, nil
}

// UpdateStatus adds or updates a manga entry in the authenticated user's list.
func (s *MangaService) UpdateStatus(ctx context.Context, mangaID int, input *UpdateMangaStatusInput) (*MangaListStatus, *Response, error) {
	if mangaID <= 0 {
		return nil, nil, fmt.Errorf("mal: manga id must be greater than zero")
	}
	if input == nil {
		return nil, nil, fmt.Errorf("mal: update manga status input cannot be nil")
	}

	values := input.ToURLValues()
	body := strings.NewReader(values.Encode())
	path := fmt.Sprintf("/manga/%d/my_list_status", mangaID)

	var result MangaListStatus
	resp, err := s.client.doRequest(ctx, http.MethodPatch, path, body, "application/x-www-form-urlencoded", &result)
	if err != nil {
		return nil, resp, fmt.Errorf("mal: update manga status: %w", err)
	}

	return &result, resp, nil
}

// DeleteStatus removes a manga from the authenticated user's list.
func (s *MangaService) DeleteStatus(ctx context.Context, mangaID int) (*Response, error) {
	if mangaID <= 0 {
		return nil, fmt.Errorf("mal: manga id must be greater than zero")
	}

	path := fmt.Sprintf("/manga/%d/my_list_status", mangaID)
	resp, err := s.client.doRequest(ctx, http.MethodDelete, path, nil, "", nil)
	if err != nil {
		return resp, fmt.Errorf("mal: delete manga status: %w", err)
	}

	return resp, nil
}
