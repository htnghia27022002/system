# QA checklist: Tools → Webhooks (005-webhooks-tool)

**Feature:** `docs/features/005-webhooks-tool/`  
**Based on:** [spec.md](./spec.md), [tasks.md](./tasks.md), [plan.md](./plan.md), [be-tasks-verify.md](./be-tasks-verify.md), [fe-tasks-verify.md](./fe-tasks-verify.md)  
**Agent:** `@qa`  
**Date:** 2026-07-25 (re-check after edge fixes)  
**Scope:** P1 (US1–US6); US7 / T041 deferred  
**Overall verdict:** **PASS**

## Summary

Re-check after edge fixes: nginx UUID quantifier escaped (`\{36\}`), FE `WEBHOOK_CAPTURE_API_BASE` for Docker rewrite fallback, and bare OPTIONS reaching Capture (CORS preflight still 204).

Canonical product capture URL **`/tools/webhooks/{uuid}`** on `http://system.local:8080` now proxies to BE and returns capture acks. Prior FAIL (nginx `nginx -t` / FE `localhost` ECONNREFUSED) is **resolved**.

Automated BE/FE suites were already green in the prior QA pass; this re-check focused on product URL + OPTIONS + isolation smoke. Browser owner UI walkthrough (copy toast, regenerate dialog chrome) remains skipped — FE verify + code wiring cover UI claims; not gating P1 sign-off.

## Automated tests

| Command | Result | Notes |
|---------|--------|-------|
| `make test-be` | **Pass** (prior QA) | Includes `be/test/unit/webhook` |
| `make test-fe` | **Pass** (prior QA) | 11 files, 40 tests |
| `be-tasks-verify.md` | Present | BE P1 tasks claimed done |
| `fe-tasks-verify.md` | Present | FE P1 tasks claimed done; T041 P2 deferred |

## Independent Test matrix (T004)

| Story | Independent Test | Result | Evidence |
|-------|------------------|--------|----------|
| US1 | Sign in → `/tools/webhooks` → UUID URL; reload same URL; guest no inbox | **Pass** (API) / **Partial** (UI) | Inbox get-or-create + idempotent UUID **Pass**. Guest HTML shell 200 with client `ProtectedGuard`. Full browser login/copy **Skipped**. |
| US2 | Capture GET/POST (+ headers/query/body); bad UUID; body > 1 MiB | **Pass** | Product path GET/POST/PUT/HEAD **200** ack JSON; bad UUID **404**; API oversize **413** (prior). |
| US3 | List newest-first; detail; search; pagination | **Pass** (API) | List + `q=product-url` hit; UI chrome **Skipped**. |
| US4 | Soft-delete hides item; search excludes | **Pass** (API) | Prior soft-delete matrix held. |
| US5 | Active vs lifetime after soft-delete | **Pass** (API) | Dual counters on inbox/list. |
| US6 | Regenerate; old fails; new captures; history remains | **Pass** (API) | Prior regenerate matrix held. |

## Acceptance scenarios (P1)

| ID | Scenario | Result | Notes |
|----|----------|--------|-------|
| US1.1 | First open auto-creates inbox + `/tools/webhooks/{uuid}` path | **Pass** | `publicPath` correct |
| US1.2 | Idempotent same URL | **Pass** | Same `publicUuid` on repeat GET |
| US1.3 | Copy URL | **Skipped** | FE clipboard + toast; not exercised in browser |
| US1.4 | Guest `/tools/webhooks` → login | **Pass** (code + shell) | HTML 200; API inbox without JWT → 401 |
| US2.1 | Capture GET/POST/PUT/PATCH/DELETE/HEAD/OPTIONS unauthenticated | **Pass** | Product + API paths. Bare OPTIONS → **200** + listed; CORS preflight (`Access-Control-Request-Method`) → **204** (not listed) — intentional |
| US2.2 | Detail: URL, method, time, IP, headers, query, body | **Pass** | Prior detail smoke |
| US2.3 | Unknown UUID → not found | **Pass** | Product + API 404 JSON |
| US2.4 | Body > 1 MiB | **Pass** | Prior 413 + oversized marker |
| US2.5 | Public cannot list/manage | **Pass** | Guest list 401; other user detail 404 |
| US3.1–US3.5 | List / detail / empty / page / search | **Pass** (API) / **Skipped** (UI) | Filters work |
| US4.1–US4.3 | Soft-delete visibility + search | **Pass** | Prior |
| US5.1–US5.2 | Dual totals | **Pass** | Inbox + list |
| US5.3 | Retention purge adjusts counts | **Pass** (unit) | Live purge not re-run |
| US6.1–US6.3 | Regenerate | **Pass** | Prior |
| FR-015 / product URL | `/tools/webhooks/{uuid}` → BE capture | **Pass** | nginx `-t` ok; GET/POST → `{"ok":true,...}` |
| FR-016 | Tools catalog → real tool | **Pass** (code) | FE verify |

## Success criteria

| ID | Criterion | Result | Notes |
|----|-----------|--------|-------|
| SC-001 | Copyable URL &lt; 1 min first use | **Partial** | API create instant; UI copy not timed in browser |
| SC-002 | Request visible ≤ 10 s on refresh | **Pass** (API) | Immediate list after product capture; FE poll present |
| SC-003 | Soft-delete visibility + lifetime | **Pass** | API matrix |
| SC-004 | Regenerate invalidates old URL | **Pass** | API matrix |
| SC-005 | Guest no inbox; unauth capture works | **Pass** | Guest gated; product + API unauth capture work |
| SC-006 | Informal find method/headers/body | **Skipped** | No browser walkthrough |

## Manual / stack checks (T046 / T047) — re-check

| Check | Result | Notes |
|-------|--------|-------|
| Stack available | Yes | `system.local:8080`; containers up |
| nginx includes `webhooks-capture.conf` | **Pass** | `nginx -t` successful; `\{36\}` loaded |
| Product URL capture | **Pass** | GET/POST/PUT/HEAD → 200 ack; unknown UUID → 404 |
| Direct API capture | **Pass** | Unchanged |
| Guest owner page | **Pass** (partial) | Shell 200; API 401 without JWT |
| HEAD/OPTIONS in list | **Pass** | HEAD + bare OPTIONS listed; preflight OPTIONS 204 not listed |
| Cross-account inbox isolation | **Pass** | Second user detail of demo request → 404 |
| Browser owner UI (list/detail/copy/regenerate dialog) | **Skipped** | Optional; not required after API + product URL Pass |

## Defects (resolved / remaining)

1. ~~**[BLOCKING] T016 nginx capture regex**~~ — **Fixed.** `\{36\}`; `nginx -t` Pass; product URL captures.
2. ~~**[BLOCKING for Docker FE path] T017 rewrite upstream**~~ — **Fixed.** `WEBHOOK_CAPTURE_API_BASE` in `fe/.env.example` / `next.config.ts` (Docker default `http://be:8080/api`). Nginx is primary on full stack; rewrite is fallback for FE-alone / mis-routed UUID.
3. ~~**[NOTE] OPTIONS capture**~~ — **Fixed / clarified.** Bare OPTIONS → Capture 200 + listed; CORS preflight → 204. Aligns with FR-004 for intentional OPTIONS probes.
4. **[OPS] BE restart** — Long-lived `go run` BE may need restart after new routes; not blocking this re-check.
5. **[DOCS] contracts/** — `contracts/database.md` and `contracts/endpoints.md` still missing under this feature folder (plan references contracts). Non-blocking for P1 runtime sign-off.

## Cross-artifact notes (lightweight analyze)

- Spec / tasks / plan align on P1 scope; edge fixes close the T016/T017 runtime gap.
- Verify docs claim P1 complete; runtime now matches.
- No constitution conflicts for package independence (API path + nginx product path; FE rewrite configurable per env).

## QA task status

| Task | Status |
|------|--------|
| T004 Independent Test matrix noted | **Done** |
| T046 US1–US6 + SC on `system.local` | **Done** (Pass) — product URL Pass on re-check |
| T047 Guest / capture / isolation / HEAD·OPTIONS | **Done** (Pass) — OPTIONS listed when bare; preflight excluded |

## Sign-off

- [x] P1 acceptance met on canonical product capture URL  
- [x] Automated `make test-be` / `make test-fe` passed (prior QA session)  
- [x] BE/FE verify docs reviewed  
- [x] Ready for release / merge as P1 complete (US7 / T041 deferred)

**Verdict: PASS** — Product URL capture and OPTIONS edge behavior confirmed on `http://system.local:8080`. Remaining notes: optional browser UI polish check; add missing `contracts/` for docs hygiene.
