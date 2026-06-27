---
name: be-test
description: BE test agent — unit, integration, and e2e tests in be/test/. Use when writing or updating backend tests.
---

# BE Test Agent

## Invocation

```text
@be/test Add unit tests for OAuth registry
@be/test Add e2e test for register flow
```

Also read: [`be/test/README.md`](../../be/test/README.md), [`be/AGENTS.md`](../../be/AGENTS.md) §11, `be-develop` skill.

## Rules (mandatory)

1. **All tests** go under `be/test/` — never `internal/*_test.go` or `public/*_test.go`.
2. **Unit** → `test/unit/<feature>/`, no Postgres, use `test/testutil` mocks.
3. **Integration** → `test/integration/`, first line `//go:build integration`, use `testutil.ConnectPostgres`.
4. **E2E** → `test/e2e/`, first line `//go:build e2e`, use `testutil.NewTestRouter`.
5. Shared mocks/helpers only in `test/testutil/`.

## GitNexus

Before adding tests for existing behavior:

```bash
npx gitnexus query "auth login test"
npx gitnexus impact AuthService
```

## Verify (repo root)

```bash
make test-be                 # unit
make test-be-integration     # needs docker stack + Postgres
make test-be-e2e
make test-be-all
```

See [`.cursor/rules/environment.mdc`](../../.cursor/rules/environment.mdc).

## Language

Test names, comments, and docs in **English only**.
