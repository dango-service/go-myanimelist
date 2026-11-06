package mal

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// UserAnimeStatistics summarizes a user's anime tracking stats.
type UserAnimeStatistics struct {
	NumItemsWatching    int     `json:"num_items_watching"`
	NumItemsCompleted   int     `json:"num_items_completed"`
	NumItemsOnHold      int     `json:"num_items_on_hold"`
	NumItemsDropped     int     `json:"num_items_dropped"`
	NumItemsPlanToWatch int     `json:"num_items_plan_to_watch"`
	NumItems            int     `json:"num_items"`
	NumDaysWatched      float64 `json:"num_days_watched"`
	NumDaysWatching     float64 `json:"num_days_watching"`
	NumDaysCompleted    float64 `json:"num_days_completed"`
	NumDaysOnHold       float64 `json:"num_days_on_hold"`
	NumDaysDropped      float64 `json:"num_days_dropped"`
	NumDays             float64 `json:"num_days"`
	NumEpisodes         int     `json:"num_episodes"`
	NumTimesRewatched   int     `json:"num_times_rewatched"`
	MeanScore           float64 `json:"mean_score"`
}

// User represents a MyAnimeList user profile.
type User struct {
	ID              int                  `json:"id"`
	Name            string               `json:"name"`
	Picture         string               `json:"picture,omitempty"`
	Gender          string               `json:"gender,omitempty"`
	Birthday        string               `json:"birthday,omitempty"`
	Location        string               `json:"location,omitempty"`
	JoinedAt        string               `json:"joined_at,omitempty"`
	AnimeStatistics *UserAnimeStatistics `json:"anime_statistics,omitempty"`
	TimeZone        string               `json:"time_zone,omitempty"`
	IsSupporter     bool                 `json:"is_supporter,omitempty"`
}

// GetMyProfile retrieves the authenticated user's profile and settings.
func (s *UserService) GetMyProfile(ctx context.Context, fields ...string) (*User, *Response, error) {
	path := "/users/@me"
	if len(fields) > 0 {
		values := url.Values{}
		values.Set("fields", strings.Join(fields, ","))
		path += "?" + values.Encode()
	}

	var result User
	resp, err := s.client.doRequest(ctx, http.MethodGet, path, nil, "", &result)
	if err != nil {
		return nil, resp, fmt.Errorf("mal: get my profile: %w", err)
	}

	return &result, resp, nil
}
