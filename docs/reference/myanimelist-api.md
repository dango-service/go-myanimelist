# MyAnimeList API v2 Technical Reference

Technical specification for MyAnimeList REST API v2 (`https://api.myanimelist.net/v2`).

---

## 1. Authentication

- **Public Endpoints**: Send Client ID in header `X-MAL-CLIENT-ID: <client_id>`.
- **User Authenticated Endpoints**: Send bearer token in header `Authorization: Bearer <access_token>`.
- **OAuth2 PKCE**:
  - Authorization URL: `https://myanimelist.net/v1/oauth2/authorize`
  - Token URL: `https://myanimelist.net/v1/oauth2/token`
  - Code challenge method: `plain` (or RFC 7636 verifier 43-128 chars).

---

## 2. Rate Limits

- MAL enforces strict rate limits (~1-2 req/sec sustained).
- Status code: `429 Too Many Requests`.
- Headers: `Retry-After: <seconds>` or RFC1123 HTTP-Date.

---

## 3. Endpoints Overview

- `GET /anime?q=...&limit=...&offset=...&fields=...`: Anime text search.
- `GET /anime/{id}?fields=...`: Detailed anime record.
- `GET /anime/ranking?ranking_type=...`: Top anime rankings (`all`, `airing`, `upcoming`, `tv`, `movie`, `special`, `bypopularity`, `favorite`).
- `GET /anime/season/{year}/{season}`: Seasonal anime (`winter`, `spring`, `summer`, `fall`).
- `GET /users/@me/animelist?status=...&sort=...&limit=...&offset=...&fields=list_status`: User anime list.
- `PATCH /anime/{id}/my_list_status`: Update anime list status (`status`, `is_rewatching`, `score`, `num_watched_episodes`, `priority`, `num_times_rewatched`, `rewatch_value`, `tags`, `comments`, `start_date`, `finish_date`).
- `DELETE /anime/{id}/my_list_status`: Delete anime from user list.
- `GET /manga?q=...`: Manga search.
- `GET /manga/{id}`: Detailed manga record.
- `GET /users/@me/mangalist`: User manga list.
- `PATCH /manga/{id}/my_list_status`: Update manga list status.
- `DELETE /manga/{id}/my_list_status`: Delete manga from user list.
- `GET /users/@me`: Authenticated user profile and settings.
