# Architecture Rules  -  go-myanimelist

Architectural invariants for `go-myanimelist`. Every contribution must uphold these constraints.

---

## 1. Zero External Dependencies

- `go-myanimelist` must depend ONLY on the Go standard library.
- No `golang.org/x/oauth2`, no third-party HTTP routers, loggers, or JSON libraries.
- Standard packages only: `net/http`, `encoding/json`, `crypto/rand`, `crypto/sha256`, `context`, `time`, `net/url`, `fmt`, `errors`, `io`.

---

## 2. Stateless Client & Dual Authentication

- The root `Client` holds references to modular domain services (`Anime`, `Manga`, `User`, `Auth`).
- Functional options configure credentials:
  - `WithClientID(clientID)`: Sets the public Client ID sent via `X-MAL-CLIENT-ID` header.
  - `WithToken(token)`: Sets the OAuth2 user access token sent via `Authorization: Bearer <token>` header.
- If both are set, `Authorization: Bearer <token>` takes precedence for user endpoints.

---

## 3. Streaming Deserialization & Bounded Memory

- Deserialization MUST use `json.NewDecoder` on `io.LimitReader(resp.Body, maxBodyBytes)`.
- Default body ceiling: 10 MiB (`10 * 1024 * 1024` bytes).
- Never read unbounded streams via `io.ReadAll(resp.Body)`.
- Never retain raw HTTP byte arrays on the `Response` object.

---

## 4. Rate Limit & Error Handling

- On HTTP 429 Too Many Requests, parse `Retry-After` seconds or RFC1123 HTTP-Date into `RateLimitInfo`.
- Return structured error wrapping `ErrRateLimited`.
- Provide helper predicates: `IsRateLimited(err)`, `IsNotFound(err)`, `IsUnauthorized(err)`.
- Wrap errors using `%w` (`fmt.Errorf("mal: %s: %w", op, err)`).

---

## 5. Security & Credential Protection

- `Token` structs must implement `fmt.Stringer` and `fmt.GoStringer` returning `[REDACTED]`.
- OAuth2 authorization code flow must use RFC 7636 PKCE (`code_verifier` with 128 characters and `code_challenge` / `plain` or `s256`).
- No tokens or credentials logged in stdout or tests.
