# go-myanimelist

[![Go Reference](https://pkg.go.dev/badge/github.com/dango-service/go-myanimelist.svg)](https://pkg.go.dev/github.com/dango-service/go-myanimelist)
[![CI](https://github.com/dango-service/go-myanimelist/actions/workflows/ci.yml/badge.svg)](https://github.com/dango-service/go-myanimelist/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A lightweight, standalone Go client library for the [MyAnimeList](https://myanimelist.net) API v2 (`https://api.myanimelist.net/v2`).
Designed with **zero external dependencies** (Go standard library only), direct streaming deserialization bounded by memory limits, and strict efficiency for background sync daemons like [Dango](https://github.com/dango-service/dango).

---

## Features

- **Zero External Dependencies**: Built entirely upon the Go standard library (`net/http`, `encoding/json`, `crypto/rand`, `context`, `net/url`).
- **Streaming JSON**: Direct streaming decoding with `json.NewDecoder` bounded by `io.LimitReader` (10 MiB ceiling) to prevent unbounded memory growth and memory exhaustion attacks.
- **Dual Authentication**: Automatic header negotiation - uses `X-MAL-CLIENT-ID` for public catalog lookups and `Authorization: Bearer <token>` for authenticated user actions.
- **Robust Rate Limit Handling**: Automatic parsing of MyAnimeList `Retry-After` headers (both integer seconds and RFC1123 HTTP-Date) on HTTP 429 with typed sentinel `ErrRateLimited` errors.
- **Credential Protection**: Redacted string formatting (`[REDACTED]`) on `Token` to prevent secret leakage in logs, traces, or debug outputs.
- **Native RFC 7636 PKCE**: Built-in cryptographic code verifier generation and OAuth2 authorization helpers without needing `golang.org/x/oauth2`.
- **Modular Services**:
  - `Anime`: Search, details, rankings (`all`, `airing`, `upcoming`, `tv`, `movie`, `bypopularity`), seasonal listings (`winter`, `spring`, `summer`, `fall`), and user anime list tracking.
  - `Manga`: Search, details, and user reading list tracking.
  - `User`: Authenticated user profile and anime statistics (`/users/@me`).
  - `Auth`: PKCE code verifier generator, authorization URL formatter, authorization code exchange, and token refresh.
- **100% Deterministic Testing**: Complete test suite using `net/http/httptest` with zero real-network requests.

---

## Installation

```bash
go get github.com/dango-service/go-myanimelist
```

---

## Usage Examples

### 1. Client Initialization (Public vs Authenticated)

```go
package main

import (
    "context"
    "fmt"
    "log"

    mal "github.com/dango-service/go-myanimelist"
)

func main() {
    // Public queries using Client ID
    publicClient := mal.NewClient("YOUR_MAL_CLIENT_ID")

    // Or authenticated user queries using Bearer Token
    userClient := mal.NewClient("YOUR_MAL_CLIENT_ID",
        mal.WithToken("USER_ACCESS_TOKEN"),
    )

    // Check authenticated user profile
    profile, _, err := userClient.User.GetMyProfile(context.Background())
    if err != nil {
        log.Fatalf("failed to fetch profile: %v", err)
    }

    fmt.Printf("Logged in as: %s (ID: %d)\n", profile.Name, profile.ID)
}
```

### 2. Anime Search & Details

```go
// Search anime with optional field expansion
results, _, err := client.Anime.Search(ctx, "Frieren", mal.WithLimit(5))
for _, edge := range results.Data {
    fmt.Printf("[%d] %s\n", edge.Node.ID, edge.Node.Title)
}

// Fetch detailed metadata
anime, _, err := client.Anime.GetByID(ctx, 52991, "alternative_titles", "genres", "num_episodes")
fmt.Printf("English title: %s, Episodes: %d\n", anime.AlternativeTitles.En, anime.NumEpisodes)
```

### 3. User Anime List (Sync & Mutations)

```go
// Fetch user's watching list (pass "@me" for current user)
animelist, _, err := client.Anime.GetUserAnimeList(ctx, "@me", mal.WithLimit(100))

// Update anime status and progress
status := mal.AnimeStatusCompleted
score := 10
episodes := 28
comments := "A modern fantasy masterpiece"

updated, _, err := client.Anime.UpdateStatus(ctx, 52991, &mal.UpdateStatusInput{
    Status:             &status,
    Score:              &score,
    NumWatchedEpisodes: &episodes,
    Comments:           &comments,
    Tags:               []string{"masterpiece", "fantasy"},
})

// Delete anime from list
_, err = client.Anime.DeleteStatus(ctx, 52991)
```

### 4. Native PKCE OAuth2 Authorization

```go
// 1. Generate high-entropy 128-char code verifier
codeVerifier, err := client.Auth.GenerateCodeVerifier()
if err != nil {
    log.Fatal(err)
}

// 2. Generate authorization URL (code_challenge = codeVerifier in 'plain' method)
authURL := client.Auth.BuildAuthorizeURL("https://dango.app/callback", "csrf_state_token", codeVerifier)
fmt.Printf("Open browser to: %s\n", authURL)

// 3. Exchange code for access & refresh tokens
token, _, err := client.Auth.ExchangeCode(ctx, clientSecret, authCode, codeVerifier, "https://dango.app/callback")
if err != nil {
    log.Fatal(err)
}

// 4. Refresh token when expired
newToken, _, err := client.Auth.RefreshToken(ctx, clientSecret, token.RefreshToken)
```

---

## License

MIT License. See [LICENSE](LICENSE) for details.
