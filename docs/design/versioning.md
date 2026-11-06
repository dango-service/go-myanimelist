# Versioning and Release Lifecycle  -  go-myanimelist

Semantic Versioning 2.0.0 guidelines for `go-myanimelist`.

---

## 1. SemVer 2.0.0

- `MAJOR`: Breaking changes to public interfaces or option types.
- `MINOR`: Backward-compatible new endpoints, options, or data structures.
- `PATCH`: Bug fixes, documentation updates, or internal performance enhancements.

---

## 2. Release Lifecycle

- Initial production release: `v0.1.0`.
- Release tags must be annotated git tags (`git tag -a v0.1.0 -m "Release v0.1.0"`).
- Git tags trigger automated GitHub Actions CI.
