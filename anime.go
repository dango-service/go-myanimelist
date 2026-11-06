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

// Anime represents a MyAnimeList anime record.
type Anime struct {
	ID                     int                `json:"id"`
	Title                  string             `json:"title"`
	MainPicture            *Picture           `json:"main_picture,omitempty"`
	AlternativeTitles      *AlternativeTitles `json:"alternative_titles,omitempty"`
	StartDate              string             `json:"start_date,omitempty"`
	EndDate                string             `json:"end_date,omitempty"`
	Synopsis               string             `json:"synopsis,omitempty"`
	Mean                   float64            `json:"mean,omitempty"`
	Rank                   int                `json:"rank,omitempty"`
	Popularity             int                `json:"popularity,omitempty"`
	NumListUsers           int                `json:"num_list_users,omitempty"`
	NumScoringUsers        int                `json:"num_scoring_users,omitempty"`
	NSFW                   string             `json:"nsfw,omitempty"`
	CreatedAt              string             `json:"created_at,omitempty"`
	UpdatedAt              string             `json:"updated_at,omitempty"`
	MediaType              string             `json:"media_type,omitempty"`
	Status                 string             `json:"status,omitempty"`
	Genres                 []Genre            `json:"genres,omitempty"`
	MyListStatus           *AnimeListStatus   `json:"my_list_status,omitempty"`
	NumEpisodes            int                `json:"num_episodes,omitempty"`
	StartSeason            *SeasonInfo        `json:"start_season,omitempty"`
	Broadcast              *Broadcast         `json:"broadcast,omitempty"`
	Source                 string             `json:"source,omitempty"`
	AverageEpisodeDuration int                `json:"average_episode_duration,omitempty"`
	Rating                 string             `json:"rating,omitempty"`
	Pictures               []Picture          `json:"pictures,omitempty"`
	Studios                []Studio           `json:"studios,omitempty"`
}

// AnimeListStatus reflects a user's watch list status for an anime.
type AnimeListStatus struct {
	Status             string    `json:"status"`
	Score              int       `json:"score"`
	NumEpisodesWatched int       `json:"num_episodes_watched"`
	IsRewatching       bool      `json:"is_rewatching"`
	UpdatedAt          time.Time `json:"updated_at,omitempty"`
	Priority           int       `json:"priority,omitempty"`
	NumTimesRewatched  int       `json:"num_times_rewatched,omitempty"`
	RewatchValue       int       `json:"rewatch_value,omitempty"`
	Tags               []string  `json:"tags,omitempty"`
	Comments           string    `json:"comments,omitempty"`
	StartDate          string    `json:"start_date,omitempty"`
	FinishDate         string    `json:"finish_date,omitempty"`
}

// AnimeEdge wraps an Anime node in paginated list results.
type AnimeEdge struct {
	Node Anime `json:"node"`
}

// AnimeList represents paginated anime search and seasonal results.
type AnimeList struct {
	Data   []AnimeEdge `json:"data"`
	Paging Paging      `json:"paging"`
}

// RankingInfo contains the ranking position of a ranked node.
type RankingInfo struct {
	Rank int `json:"rank"`
}

// RankingEdge wraps an Anime node with ranking details.
type RankingEdge struct {
	Node    Anime       `json:"node"`
	Ranking RankingInfo `json:"ranking"`
}

// AnimeRankingList represents paginated ranking results.
type AnimeRankingList struct {
	Data   []RankingEdge `json:"data"`
	Paging Paging        `json:"paging"`
}

// Search searches for anime by title keywords.
func (s *AnimeService) Search(ctx context.Context, query string, opts ...ListOption) (*AnimeList, *Response, error) {
	trimmed := strings.TrimSpace(query)
	if trimmed == "" {
		return nil, nil, fmt.Errorf("mal: search query cannot be empty")
	}

	values := url.Values{}
	values.Set("q", trimmed)
	for _, opt := range opts {
		if opt != nil {
			opt(&values)
		}
	}

	path := "/anime?" + values.Encode()
	var result AnimeList
	resp, err := s.client.doRequest(ctx, http.MethodGet, path, nil, "", &result)
	if err != nil {
		return nil, resp, fmt.Errorf("mal: search anime: %w", err)
	}

	return &result, resp, nil
}

// GetByID retrieves anime details by its MyAnimeList ID.
func (s *AnimeService) GetByID(ctx context.Context, id int, fields ...string) (*Anime, *Response, error) {
	if id <= 0 {
		return nil, nil, fmt.Errorf("mal: anime id must be greater than zero")
	}

	path := fmt.Sprintf("/anime/%d", id)
	if len(fields) > 0 {
		values := url.Values{}
		values.Set("fields", strings.Join(fields, ","))
		path += "?" + values.Encode()
	}

	var result Anime
	resp, err := s.client.doRequest(ctx, http.MethodGet, path, nil, "", &result)
	if err != nil {
		return nil, resp, fmt.Errorf("mal: get anime by id: %w", err)
	}

	return &result, resp, nil
}

// GetRanking retrieves anime rankings (e.g. "all", "airing", "upcoming", "tv", "movie", "bypopularity", "favorite").
func (s *AnimeService) GetRanking(ctx context.Context, rankingType string, opts ...ListOption) (*AnimeRankingList, *Response, error) {
	trimmed := strings.TrimSpace(rankingType)
	if trimmed == "" {
		trimmed = "all"
	}

	values := url.Values{}
	values.Set("ranking_type", trimmed)
	for _, opt := range opts {
		if opt != nil {
			opt(&values)
		}
	}

	path := "/anime/ranking?" + values.Encode()
	var result AnimeRankingList
	resp, err := s.client.doRequest(ctx, http.MethodGet, path, nil, "", &result)
	if err != nil {
		return nil, resp, fmt.Errorf("mal: get anime ranking: %w", err)
	}

	return &result, resp, nil
}

// GetSeasonal retrieves seasonal anime for a specific year and season ("winter", "spring", "summer", "fall").
func (s *AnimeService) GetSeasonal(ctx context.Context, year int, season string, opts ...ListOption) (*AnimeList, *Response, error) {
	trimmedSeason := strings.ToLower(strings.TrimSpace(season))
	if trimmedSeason == "" {
		return nil, nil, fmt.Errorf("mal: season cannot be empty")
	}

	values := url.Values{}
	for _, opt := range opts {
		if opt != nil {
			opt(&values)
		}
	}

	path := fmt.Sprintf("/anime/season/%s/%s", strconv.Itoa(year), trimmedSeason)
	if len(values) > 0 {
		path += "?" + values.Encode()
	}

	var result AnimeList
	resp, err := s.client.doRequest(ctx, http.MethodGet, path, nil, "", &result)
	if err != nil {
		return nil, resp, fmt.Errorf("mal: get seasonal anime: %w", err)
	}

	return &result, resp, nil
}
