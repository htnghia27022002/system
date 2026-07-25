# FE tasks verify: Tools → Webhooks

**Feature:** `docs/features/005-webhooks-tool/`  
**Based on:** [spec.md](spec.md), [tasks.md](tasks.md), [plan.md](plan.md), [contracts/](contracts/)  
**Agent:** `@fe`  
**Date:** 2026-07-25 (admin RBAC relocate; URL toolbar compact polish)

## Summary

Owner Webhooks UI moved to **`/admin/tools/webhooks`** with `PermissionKeys.webhooks` (`view` / `modify`), sidebar **Tools → Webhooks**, `PermissionGuard` / `PermissionGate`, and `/tools/webhooks` → admin redirect. Public capture path `/tools/webhooks/{uuid}` unchanged. Follow-up polish: public URL block compacted to a single toolbar row so the inbox list dominates the viewport. `make test-fe` passed (40 tests). T041 (P2 mobile stack) remains deferred.

## Contract alignment (permissions.md)

| Contract surface | FE status |
|------------------|-----------|
| `PermissionKeys.webhooks.view` / `.modify` | `permission-keys.ts` |
| Admin page `/admin/tools/webhooks` + `PermissionGuard(webhooks:view)` | `app/admin/tools/webhooks/page.tsx` |
| Sidebar Tools → Webhooks filtered by `hasPermission` | `app-sidebar.tsx` + `nav.tools` / `nav.webhooks` i18n |
| Regenerate / soft-delete `PermissionGate(webhooks:modify)` | `webhooks-inbox.tsx`, `webhooks-request-detail.tsx` |
| Old owner `/tools/webhooks` → `/admin/tools/webhooks` | redirect `page.tsx` |
| Capture `/tools/webhooks/{uuid}` not an owner Next page | Confirmed (rewrite only) |
| Mock permission catalog | `access-control.mock.ts` seeds both keys |

## Tasks completed

Prior P1 `[FE]` tasks (T002–T043, T045) remain Done from earlier passes; this pass relocates owner UI per updated `contracts/permissions.md`.

| Item | Description | Status | Evidence |
|------|-------------|--------|----------|
| RBAC keys | Add `webhooks` resource | Done | `permission-keys.ts` |
| Admin route | `/admin/tools/webhooks` + guard | Done | `app/admin/tools/webhooks/page.tsx` |
| Sidebar | Tools group + Webhooks child | Done | `app-sidebar.tsx`, breadcrumbs |
| i18n | `nav.tools`, `nav.webhooks` | Done | `locales/en\|vi/admin.json` |
| Redirect | `/tools/webhooks` → admin | Done | public `tools/webhooks/page.tsx` |
| Catalog href | Owner link → admin path | Done | `catalog.ts` + test |
| Mutate gates | regenerate + remove | Done | `PermissionGate` |
| Mock catalog | seed webhooks keys | Done | `access-control.mock.ts` |
| T041 | P2 stacked mobile | Deferred | spec P2 |
| Verify | `fe-tasks-verify.md` | Done | this file |

## Verification commands

```bash
make test-fe
pnpm lint   # from fe/
```

| Command | Result | Notes |
|---------|--------|-------|
| `make test-fe` | Pass | 11 files, 40 tests (2026-07-25; re-run after URL toolbar compact) |
| `pnpm lint` | Skipped this pass | Prior unrelated failures |

## Follow-up polish (2026-07-25)

| Change | Files | Notes |
|--------|-------|-------|
| Compact public URL toolbar | `webhooks-inbox.tsx` | Dropped H2 + visible hint; icon-sm Copy/Refresh/Regenerate + inline totals; `p-3`; inbox `min-h-[40rem]` / `lg:h-[min(78vh,52rem)]` |
| Light page header tighten | `webhooks-tool-page.tsx` | Slightly smaller H1/subtitle; page shell `gap-4 p-4` kept |

## Key routes / files

| Surface | Path |
|---------|------|
| Owner UI | `/admin/tools/webhooks` → `WebhooksToolPage` + `PermissionGuard` |
| Legacy owner | `/tools/webhooks` → redirect to admin |
| Capture (local FE) | Next rewrite `/tools/webhooks/{uuid}` → BE |
| Sidebar | Tools → Webhooks (`webhooks:view`) |
| Mutate | regenerate / remove (`webhooks:modify`) |

## How to see the menu

1. Ensure BE permissions are seeded (`webhooks:view` / `webhooks:modify` on Admin role) — reseed DB if needed after BE catalog update.
2. Log in as admin (`admin@example.com` / seeded password).
3. Open `/admin` — sidebar shows **Tools → Webhooks**.
4. Open Access Control → Roles — Admin role includes the webhooks keys.
5. Direct URL `/admin/tools/webhooks` works with `webhooks:view`; regenerate/remove need `webhooks:modify`.

## Acceptance coverage (FE-relevant)

| Spec / contract scenario | Covered by | Result |
|--------------------------|------------|--------|
| Owner UI under admin + view guard | PermissionGuard page | Pass (code) |
| Sidebar gated by view | `hasPermission` filter | Pass (code) |
| Mutate gated by modify | PermissionGate | Pass (code) |
| Legacy path redirects | `redirect()` | Pass (code) |
| Capture URL unchanged | rewrite + URL helper | Pass (prior + code) |
| Soft-delete / regenerate / inbox | prior US1–US6 UI | Pass (prior) |

## Gaps / follow-ups for `@qa`

- [ ] Manual: login as admin with reseeded perms → Tools → Webhooks visible
- [ ] Manual: role without `webhooks:view` hides menu and redirects from admin page
- [ ] Manual: view-only role sees inbox but not regenerate/remove
- [ ] Manual: `/tools/webhooks` redirects to `/admin/tools/webhooks` (login if needed)
- [ ] Confirm BE seed UUIDs match `contracts/permissions.md`
- [x] FE owner-path relocate + RBAC gates done (T041 P2 open)

## Sign-off (FE)

- [x] Admin RBAC relocate complete per `contracts/permissions.md`
- [x] `make test-fe` passed
- [x] Matches updated permissions contract (deviations: none for FE)
