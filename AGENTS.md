# Sandgrouse

LLM traffic compression proxy for developers on metered connections. Written in Go.

## Build & run

```bash
go build -o sg ./cmd/sg    # build
./sg start                  # start proxy on :8080
./sg status                 # check stats
go test ./...               # run tests
golangci-lint run           # lint
```

## Project structure

```
cmd/sg/          CLI entrypoint
internal/proxy/  HTTP proxy and compression logic
internal/dedup/  Content deduplication engine
internal/delta/  Delta encoding for file changes
internal/dash/   Web dashboard (localhost:8585)
```

## Code style

- Standard Go conventions (gofmt, go vet)
- No new dependencies without discussion
- Tests for all new functionality
- Keep functions small and focused

## Development workflow

All changes must be **atomic** and **methodological**. One logical change per unit of work.

### Principles

- **DRY** — Don't Repeat Yourself; extract shared logic, avoid copy-paste
- **Single Responsibility** — each function, file, and package does one thing well
- **Clean Code** — meaningful names, small functions, no dead code, no magic numbers
- **YAGNI** — don't build what isn't needed yet
- **Separation of Concerns** — keep layers distinct (transport, logic, storage)

### Workflow: code → build → test → commit

Every change follows this flow:

1. **Code** — make one atomic, focused change
2. **Build** — `go build ./...` must pass with zero errors
3. **Test** — `go test ./...` must pass; new code requires new tests
4. **Commit** — one commit per logical change (see commit conventions below)

Do not skip steps. Do not batch unrelated changes into a single commit. If the build or tests fail, fix before committing.

## Key decisions

- Single binary distribution, no runtime dependencies
- Proxy is transparent to LLM clients (no client-side changes)
- All data stays local, nothing sent except to the original API destination
- Supports Anthropic, OpenAI, Google Gemini, and any OpenAI-compatible API

## Config

User config lives at `~/.sandgrouse/config.yml`. See `config.example.yml` for reference.

## Docs

- Planning docs (PRDs, tracker): `docs/planning/`
- Issue templates: `.github/ISSUE_TEMPLATE/`
- Project manifesto: `MANIFESTO.md`
- Public roadmap: `ROADMAP.md`
- Release history: `CHANGELOG.md`

## Git commit conventions

Follow Angular conventional commit format:

- **Format:** `<type>(<scope>): <subject>`
- **Types:** `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `revert`
- **Scope:** optional, e.g. `proxy`, `dedup`, `delta`, `dash`, `cli`
- **Subject:** imperative mood, lowercase, no period at end, max 50 chars
- **Body:** use bullet points to explain what and why
  - Each bullet starts with `-`
  - Wrap at 72 characters
- **Breaking changes:** add `BREAKING CHANGE:` in the footer
- **No AI attribution:** do not include `Co-Authored-By` or any AI/Claude attribution lines

### Examples

```
feat(proxy): add gzip compression for anthropic responses

- Add gzip encoding to outbound proxy responses
- Reduce average payload size by ~60% for chat completions
- Skip compression for streaming responses under 1KB
```

```
fix(dedup): handle empty response bodies without panic

- Check for nil body before computing content hash
- Add regression test for empty-body edge case
```
