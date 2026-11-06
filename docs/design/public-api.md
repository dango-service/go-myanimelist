# Public API Contract  -  go-myanimelist

Public interface contract and signatures for `go-myanimelist`.

---

## 1. Client Constructor & Options

```go
type Client struct {
    Anime *AnimeService
    Manga *MangaService
    User  *UserService
    Auth  *AuthService
}

func NewClient(clientID string, opts ...ClientOption) *Client

type ClientOption func(*Client)

func WithHTTPClient(httpClient *http.Client) ClientOption
func WithBaseURL(baseURL string) ClientOption
func WithToken(token string) ClientOption
func WithUserAgent(userAgent string) ClientOption
func WithMaxBodyBytes(maxBytes int64) ClientOption
```

---

## 2. Service Signatures

### AnimeService
```go
func (s *AnimeService) Search(ctx context.Context, query string, opts ...ListOption) (*AnimeList, *Response, error)
func (s *AnimeService) GetByID(ctx context.Context, id int, fields ...string) (*Anime, *Response, error)
func (s *AnimeService) GetRanking(ctx context.Context, rankingType string, opts ...ListOption) (*AnimeList, *Response, error)
func (s *AnimeService) GetSeasonal(ctx context.Context, year int, season string, opts ...ListOption) (*AnimeList, *Response, error)
func (s *AnimeService) GetUserAnimeList(ctx context.Context, username string, opts ...ListOption) (*UserAnimeList, *Response, error)
func (s *AnimeService) UpdateStatus(ctx context.Context, animeID int, input *UpdateStatusInput) (*AnimeListStatus, *Response, error)
func (s *AnimeService) DeleteStatus(ctx context.Context, animeID int) (*Response, error)
```

### MangaService
```go
func (s *MangaService) Search(ctx context.Context, query string, opts ...ListOption) (*MangaList, *Response, error)
func (s *MangaService) GetByID(ctx context.Context, id int, fields ...string) (*Manga, *Response, error)
func (s *MangaService) GetUserMangaList(ctx context.Context, username string, opts ...ListOption) (*UserMangaList, *Response, error)
func (s *MangaService) UpdateStatus(ctx context.Context, mangaID int, input *UpdateStatusInput) (*MangaListStatus, *Response, error)
func (s *MangaService) DeleteStatus(ctx context.Context, mangaID int) (*Response, error)
```

### UserService
```go
func (s *UserService) GetMyProfile(ctx context.Context) (*User, *Response, error)
```

### AuthService
```go
func (s *AuthService) BuildAuthorizeURL(redirectURI, state, codeChallenge string) string
func (s *AuthService) ExchangeCode(ctx context.Context, clientSecret, code, codeVerifier, redirectURI string) (*Token, *Response, error)
func (s *AuthService) RefreshToken(ctx context.Context, clientSecret, refreshToken string) (*Token, *Response, error)
func (s *AuthService) GenerateCodeVerifier() (string, error)
```
