package mal

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// AnimeStatus represents valid watch list status values.
const (
	AnimeStatusWatching    = "watching"
	AnimeStatusCompleted   = "completed"
	AnimeStatusOnHold      = "on_hold"
	AnimeStatusDropped     = "dropped"
	AnimeStatusPlanToWatch = "plan_to_watch"
)

// UserAnimeEdge pairs an Anime record with the user's ListStatus.
type UserAnimeEdge struct {
	Node       Anime           `json:"node"`
	ListStatus AnimeListStatus `json:"list_status"`
}

// UserAnimeList represents a paginated user anime list.
type UserAnimeList struct {
	Data   []UserAnimeEdge `json:"data"`
	Paging Paging          `json:"paging"`
}

// UpdateStatusInput specifies mutation fields when updating an anime list entry.
type UpdateStatusInput struct {
	Status             *string
	Score              *int
	NumWatchedEpisodes *int
	IsRewatching       *bool
	Priority           *int
	NumTimesRewatched  *int
	RewatchValue       *int
	Tags               []string
	Comments           *string
	StartDate          *string
	FinishDate         *string
}

// ToURLValues serializes the update input into form-urlencoded parameters.
func (in *UpdateStatusInput) ToURLValues() url.Values {
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
	if in.NumWatchedEpisodes != nil {
		values.Set("num_watched_episodes", strconv.Itoa(*in.NumWatchedEpisodes))
	}
	if in.IsRewatching != nil {
		if *in.IsRewatching {
			values.Set("is_rewatching", "true")
		} else {
			values.Set("is_rewatching", "false")
		}
	}
	if in.Priority != nil {
		values.Set("priority", strconv.Itoa(*in.Priority))
	}
	if in.NumTimesRewatched != nil {
		values.Set("num_times_rewatched", strconv.Itoa(*in.NumTimesRewatched))
	}
	if in.RewatchValue != nil {
		values.Set("rewatch_value", strconv.Itoa(*in.RewatchValue))
	}
	if len(in.Tags) > 0 {
		values.Set("tags", strings.Join(in.Tags, ","))
	}
	if in.Comments != nil {
		values.Set("comments", *in.Comments)
	}
	if in.StartDate != nil && *in.StartDate != "" {
		values.Set("start_date", *in.StartDate)
	}
	if in.FinishDate != nil && *in.FinishDate != "" {
		values.Set("finish_date", *in.FinishDate)
	}

	return values
}

// GetUserAnimeList fetches a user's anime list. Pass "@me" for the currently authenticated user.
func (s *AnimeService) GetUserAnimeList(ctx context.Context, username string, opts ...ListOption) (*UserAnimeList, *Response, error) {
	trimmed := strings.TrimSpace(username)
	if trimmed == "" {
		trimmed = "@me"
	}

	values := url.Values{}
	// list_status is essential for user list sync
	values.Set("fields", "list_status")
	for _, opt := range opts {
		if opt != nil {
			opt(&values)
		}
	}

	path := fmt.Sprintf("/users/%s/animelist?%s", url.PathEscape(trimmed), values.Encode())
	var result UserAnimeList
	resp, err := s.client.doRequest(ctx, http.MethodGet, path, nil, "", &result)
	if err != nil {
		return nil, resp, fmt.Errorf("mal: get user animelist: %w", err)
	}

	return &result, resp, nil
}

// UpdateStatus adds or updates an anime entry in the authenticated user's list.
func (s *AnimeService) UpdateStatus(ctx context.Context, animeID int, input *UpdateStatusInput) (*AnimeListStatus, *Response, error) {
	if animeID <= 0 {
		return nil, nil, fmt.Errorf("mal: anime id must be greater than zero")
	}
	if input == nil {
		return nil, nil, fmt.Errorf("mal: update status input cannot be nil")
	}

	values := input.ToURLValues()
	body := strings.NewReader(values.Encode())
	path := fmt.Sprintf("/anime/%d/my_list_status", animeID)

	var result AnimeListStatus
	resp, err := s.client.doRequest(ctx, http.MethodPatch, path, body, "application/x-www-form-urlencoded", &result)
	if err != nil {
		return nil, resp, fmt.Errorf("mal: update anime status: %w", err)
	}

	return &result, resp, nil
}

// DeleteStatus removes an anime from the authenticated user's list.
func (s *AnimeService) DeleteStatus(ctx context.Context, animeID int) (*Response, error) {
	if animeID <= 0 {
		return nil, fmt.Errorf("mal: anime id must be greater than zero")
	}

	path := fmt.Sprintf("/anime/%d/my_list_status", animeID)
	resp, err := s.client.doRequest(ctx, http.MethodDelete, path, nil, "", nil)
	if err != nil {
		return resp, fmt.Errorf("mal: delete anime status: %w", err)
	}

	return resp, nil
}
