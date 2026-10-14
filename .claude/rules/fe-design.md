---
description: "FE visual design — theme, page spacing, mobile, clickable table links. Full rules in fe/DESIGN.md."
paths:
  - "fe/**/*"
---

# FE design

Authoritative visual rules: **[`fe/DESIGN.md`](../../fe/DESIGN.md)**.

Must follow when building or changing UI:

1. **Theme consistency** — CSS variables in `src/styles/index.css` (incl. `success` / `warning` and the `--z-*` stacking scale); surfaces `rounded-xl border`, controls `rounded-lg`. `pnpm lint` rejects raw palette colors, `z-[NNNN]`, arbitrary type sizes.
2. **Page containers** — `PageContainer` + `PageHeader` + `Panel` from `@/components/common/page-layout` (admin: `p-4` + `gap-4`).
3. **Mobile-first** — responsive web (usable ~375px); not a native mobile-app shell.
4. **Clickable table cells** (id, name, …) — use **link color** `text-primary` so they look clickable before hover.
5. **shadcn first** — see also `fe-ui-library-first.md`; forms: `fe-form-layout.md`.

Architecture / structure stays in `fe/CLAUDE.md`.
