---
name: be-test
description: Backend test engineer for be/test/. Use when writing or updating Go unit, integration, or e2e tests for the backend.
---

# BE Test Agent

## Invocation

```text
@agent-be/test Add unit tests for OAuth registry
@agent-be/test Add e2e test for register flow
```

Also read: [`be/test/README.md`](../../be/test/README.md), [`be/CLAUDE.md`](../../be/CLAUDE.md) §11, `be-develop` skill.

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

For queue-driven module features: start the `queue` profile (`make up-d PROFILES="dev queue"`, `NATS_ENABLED=true`) to verify NATS processing end-to-end; unit tests call the module service directly. Module tests live in `test/unit/<module>/`; registry behaviour is covered by `test/unit/module/`.

See [`.claude/rules/environment.md`](../../.claude/rules/environment.md).

## Language

Test names, comments, and docs in **English only**.
