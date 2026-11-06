# Testing Strategy  -  go-myanimelist

Deterministic testing strategy for `go-myanimelist`.

---

## 1. Zero Network Policy

- All unit tests must communicate exclusively with `net/http/httptest.Server`.
- Real network calls to `api.myanimelist.net` or `myanimelist.net` are strictly prohibited during tests.

---

## 2. Table-Driven Tests & Fixtures

- Tests must use `t.Run` for table-driven scenarios covering:
  - Success responses with realistic JSON fixtures in `testdata/`.
  - HTTP 429 Rate Limit errors with `Retry-After` (both integer seconds and RFC1123 HTTP-Date).
  - HTTP 401 Unauthorized, HTTP 404 Not Found, and HTTP 400 Bad Request.
  - Context cancellation and timeouts.
  - Empty or invalid input arguments.

---

## 3. Test Coverage Goals

- Core transport, rate limit parser, and auth token redaction: 100% coverage.
- All service endpoints: covered by at least one positive test and one error path test.
