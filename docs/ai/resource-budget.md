# Resource and Memory Budget  -  go-myanimelist

Resource ceiling and memory guidelines for `go-myanimelist`.

---

## 1. Memory Budget

- Peak RSS must not exceed 5 MiB during background polling operations.
- Per-request body allocation is strictly bounded by `io.LimitReader(resp.Body, maxBodyBytes)` (default 10 MiB).
- Avoid saving byte buffers into memory when decoding JSON responses.

---

## 2. Allocation Guidelines

- Prefer reusing structs and slicing directly into slices of pointers.
- Do not make redundant copies of query parameter maps.
- Close `resp.Body` immediately after decoding via `defer resp.Body.Close()`.

---

## 3. Concurrency Safety

- All services are stateless and operate via pointers to their parent `Client`.
- Client options are applied once upon construction in `NewClient`.
