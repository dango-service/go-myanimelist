# Scope and Bounded Context  -  go-myanimelist

Scope of supported features and boundaries for `go-myanimelist`.

---

## 1. Supported Features

- Anime metadata search, details, ranking, and seasonal listings.
- User anime list tracking (status, episodes, score, rewatch count, dates, tags).
- Manga metadata search, details, and user manga list tracking.
- Authenticated user profile retrieval (`/users/@me`).
- Native OAuth2 with PKCE code generation, authorization URL formatting, code exchange, and token refresh.
- Rate limit detection (`Retry-After` header parsing).

---

## 2. Excluded Features

- Forum discussions and legacy MyAnimeList v1 XML endpoints.
- Direct database persistence, SQLite caching, or Plex/Jellyfin webhooks (responsibility of Dango).
- GUI or CLI interactive credential input handling.
