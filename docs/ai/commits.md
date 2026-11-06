# Commit Conventions  -  go-myanimelist

Conventional Commits 1.0.0 specification for `go-myanimelist`.

---

## 1. Format

```text
<type>(<scope>)<!>: <description>
```

- **Type**: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `revert`.
- **Scope**: Lowercase, e.g. `client`, `anime`, `animelist`, `manga`, `auth`, `design`, `github`.
- **Description**: Imperative, lowercase, no trailing period, maximum 72 characters.

---

## 2. Standard Commit Progression

1. `docs(design): establish architecture rules, public api spec, and scope for phase 0`
2. `feat(client): implement core client, transport, and rate limit parsing`
3. `feat(anime): implement anime search, details, ranking, and seasonal catalog`
4. `feat(animelist): implement user anime list retrieval, status updates, and deletions`
5. `feat(manga,user,auth): implement manga catalog, user profile, and pkce oauth2 helpers`
6. `ci(github): add matrix test workflow, developer guidelines, and complete phase 5`

---

## 3. Invariants

- Atomic commits per logical change.
- Never `--amend` or `push --force`.
- Pure ASCII English descriptions.
