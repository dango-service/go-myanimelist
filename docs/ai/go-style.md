# Go Style Rules  -  go-myanimelist

Coding conventions for `go-myanimelist`.

---

## 1. Error Handling

- Always wrap errors using `%w`:
  ```go
  return nil, nil, fmt.Errorf("mal: %s: %w", op, err)
  ```
- Expose typed errors and predicates:
  ```go
  var (
      ErrRateLimited   = errors.New("mal: rate limit exceeded")
      ErrUnauthorized  = errors.New("mal: unauthorized")
      ErrNotFound      = errors.New("mal: not found")
      ErrBadRequest    = errors.New("mal: bad request")
  )
  ```

---

## 2. Context Cancellation

- `context.Context` must be the first parameter in all service methods performing network I/O.
- Check context cancellation early:
  ```go
  if err := ctx.Err(); err != nil {
      return nil, nil, fmt.Errorf("mal: %s: %w", op, err)
  }
  ```

---

## 3. Streaming Bounded Deserialization

- Always use `io.LimitReader`:
  ```go
  limitedBody := io.LimitReader(resp.Body, maxBodyBytes)
  if err := json.NewDecoder(limitedBody).Decode(v); err != nil {
      return nil, fmt.Errorf("mal: decode response: %w", err)
  }
  ```

---

## 4. Pure CGO_ENABLED=0

- No cgo dependencies.
- Library must compile on Windows, Linux, and macOS without a C compiler.
