# AI Agent Workflow  -  go-myanimelist

Step-by-step workflow for agents working on `go-myanimelist`.

---

## 1. Workflow Phases

1. **Read Specifications**: Understand API requirements in `docs/reference/myanimelist-api.md`.
2. **Deterministic TDD**: Write table-driven unit tests against `httptest.Server` with realistic JSON payloads from `testdata/`.
3. **Implementation**: Implement models, services, and option builders.
4. **Validation**: Execute `go test -v ./...`, `go vet ./...`, and `CGO_ENABLED=0 go build ./...`.
5. **Atomic Commit**: Commit the phase using Conventional Commits 1.0.0.

---

## 2. Invariants Check

- Zero real network requests during tests.
- ASCII English only in code, docs, and git commits.
- Strict bounded JSON decoding.
