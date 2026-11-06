package mal

import (
	"net/http"
	"time"
)

const (
	// DefaultBaseURL is the standard MyAnimeList API v2 endpoint.
	DefaultBaseURL = "https://api.myanimelist.net/v2"

	// DefaultUserAgent is sent with every request unless overridden.
	DefaultUserAgent = "go-myanimelist/0.1.0"

	// DefaultMaxBodyBytes limits incoming JSON payloads to 10 MiB.
	DefaultMaxBodyBytes = 10 * 1024 * 1024
)

// Client is the MyAnimeList API v2 client.
type Client struct {
	httpClient   *http.Client
	baseURL      string
	clientID     string
	token        string
	userAgent    string
	maxBodyBytes int64

	Anime *AnimeService
	Manga *MangaService
	User  *UserService
	Auth  *AuthService
}

// NewClient creates a new MyAnimeList API v2 client.
func NewClient(clientID string, opts ...ClientOption) *Client {
	c := &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:      DefaultBaseURL,
		clientID:     clientID,
		userAgent:    DefaultUserAgent,
		maxBodyBytes: DefaultMaxBodyBytes,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(c)
		}
	}

	c.Anime = &AnimeService{client: c}
	c.Manga = &MangaService{client: c}
	c.User = &UserService{client: c}
	c.Auth = &AuthService{client: c}

	return c
}

// AnimeService provides access to anime search, details, ranking, and list management.
type AnimeService struct {
	client *Client
}

// MangaService provides access to manga search, details, and list management.
type MangaService struct {
	client *Client
}

// UserService provides access to user profile information.
type UserService struct {
	client *Client
}

// AuthService provides helpers for OAuth2 PKCE authorization code exchange and token refresh.
type AuthService struct {
	client *Client
}
