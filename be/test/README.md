# BE tests

All backend tests live under `be/test/`. Do **not** add `*_test.go` files under `internal/` or `public/`.

## Layout

```text
test/
├── testutil/          # Shared mocks, DB helpers, HTTP router setup
├── unit/              # Fast tests, no external deps (mocks only)
├── integration/       # Service + repository against real Postgres
└── e2e/               # Full HTTP stack (Gin + handlers + DB)
```

| Tier | Path | Build tag | Requires Postgres |
|------|------|-----------|-------------------|
| Unit | `test/unit/` | none | No |
| Integration | `test/integration/` | `integration` | Yes (skips if unavailable) |
| E2E | `test/e2e/` | `e2e` | Yes (skips if unavailable) |

## Commands (from repo root)

```bash
make test-be              # unit only (default CI path)
make test-be-integration  # integration tier
make test-be-e2e          # e2e tier
make test-be-all          # all tiers
```

Inside `be/` container or with local Go:

```bash
make test                 # unit
make test-integration
make test-e2e
make test-all
```

Use `-short` to skip integration/e2e tiers that call `testutil.SkipIfShort`.

## Adding tests

1. **Unit** — `test/unit/<feature>/<name>_test.go`, package `<feature>_test`, use `test/testutil` mocks.
2. **Integration** — `test/integration/`, file starts with `//go:build integration`, use `testutil.ConnectPostgres`.
3. **E2E** — `test/e2e/`, file starts with `//go:build e2e`, use `testutil.NewTestRouter`.

Before editing shared test helpers, run GitNexus: `npx gitnexus query "test auth login"`.

## Environment (integration / e2e)

Uses the same env as the running BE container (`DB_HOST=postgres`, etc.). Stack must be up:

```bash
make up-d && make test-be-all
```
