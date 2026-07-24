# Implementation Plan: Public Landing Rebuild + Tools IA + Shared Chrome

**Feature**: `003-landing-page` | **Date**: 2026-07-24 | **Spec**: [spec.md](./spec.md) | **Tasks**: [tasks.md](./tasks.md)

**Input**: Feature specification from `/docs/features/003-landing-page/spec.md` (greenfield rebuild revision)

**Stale note**: This plan **replaces** the prior evolve-tokenomics plan. Prior [`fe-tasks-verify.md`](./fe-tasks-verify.md) is **obsolete** — `@fe` must re-verify after re-implement. Do not treat prior FE verify as evidence for this revision.

## Summary

Ship a **brand-new** public marketing home at `/` (GSAP scrollytelling + reduced-motion fallback), an extensible Tools IA (`/tools` hub + first nested page `/tools/webhooks`), and **shared public chrome** (Navigate, auth-aware Login/Logout, Theme, Locale **Select with flag icons** for EN/VI). Admin chrome gains a **home icon** to public `/`. Catalog stays FE-only and drives landing teaser + hub. Discard the retired landing prototype composition (including tokenomics as the primary story). No backend/CMS in Phase 1.

## Technical Context

**Language/Version**: TypeScript (strict) on Next.js 15 App Router + React 19  

**Primary Dependencies**: GSAP + ScrollTrigger (existing), Framer Motion (optional chrome motion), Tailwind CSS v4, shadcn/ui (`Select`, `Button`, `Card`, `Badge`, sidebar primitives), `next/link`, `next-themes`, `react-i18next`  

**Storage**: N/A — static TypeScript catalog/content modules  

**Testing**: Vitest + Testing Library (`fe/`); update landing/tools tests; QA manual acceptance  

**Target Platform**: Modern browsers (desktop + mobile); public visitors + signed-in users on public chrome; admin users for home escape  

**Project Type**: Frontend package within monorepo (`fe/` standalone)  

**Performance Goals**: Smooth scroll-linked storytelling when motion enabled; usable reduced-motion path; no horizontal overflow of primary hero content at 375px  

**Constraints**: English feature docs; FE package independence (no `be/` imports); feature import boundaries (`@/features/tools` public barrel); Locale Select (not toggle) on public P1 routes; reduced-motion required  

**Scale/Scope**: Three public routes (`/`, `/tools`, `/tools/webhooks`); catalog N small for P1 (append-extensible); full webhook product deferred (US12)

## Constitution Check

*GATE: Must pass before design lock. Re-checked after design below.*

| Principle | Status | Notes |
|-----------|--------|--------|
| I. Package independence | **Pass** | FE-only; catalog/content in `fe/`; no BE coupling |
| II. Role-owned artifacts | **Pass** | `tasks.md` + `plan.md` owned by technical-architect; no app code; verify remains `@fe` |
| III. Spec before code | **Pass** | Regenerated `spec.md` → `tasks.md` → `plan.md` before `@fe` re-implement |
| IV. API contract alignment | **N/A (P1)** | No new HTTP APIs |
| V. English documentation | **Pass** | All feature docs English; UI EN/VI via i18n |

**Post-design re-check**: Still pass. Shared chrome lives in `components/common/` (cross-feature); domain pages stay in `features/landing` and `features/tools`. Layout split preserves independent scroll behavior without a second home route.

## Complexity Tracking

> No constitution violations requiring justification.

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| — | — | — |

---

## Research Decisions (Phase 0)

### R1 — Greenfield landing (not evolve tokenomics)

- **Decision**: Rebuild `/` composition and visual story; retire prototype sections (notably tokenomics) from the shipped primary experience. Reuse GSAP/scrolly primitives (`useGsapContext`, `ScrollyShell`, reduced-motion hook) where they accelerate delivery.
- **Rationale**: Spec clarification — rebuild, not evolve. Stakeholders must not see the old prototype as the product story.
- **Alternatives**: Incremental restyle of tokenomics prototype — **superseded**.

### R2 — FE-only shared catalog

- **Decision**: Keep `fe/src/features/tools/catalog.ts` as single source of truth for teaser + hub (+ hrefs).
- **Rationale**: Spec assumptions; independent FE deploy; no CMS.
- **Alternatives**: BE CMS — deferred.

### R3 — Nested tool route pattern via concrete Webhooks page

- **Decision**: Ship explicit App Router page `fe/src/app/(public)/tools/webhooks/page.tsx` + feature shell component. P1 does **not** require a dynamic `[id]` segment; Webhooks proves the `/tools/{id}` IA for later tools.
- **Rationale**: Spec P1 includes `/tools/webhooks` shell; placeholder body OK.
- **Alternatives**: Dynamic `[id]` catch-all in P1 — unnecessary complexity until multiple real tool apps exist.

### R4 — Shared public chrome in `components/common`

- **Decision**: Introduce `PublicSiteHeader` + `LocaleSelect` under `fe/src/components/common/`. Landing and tools mount the same header. Landing may keep optional section-jump links as secondary UX, but site Navigate (Home, Tools, auth) is shared.
- **Rationale**: FR-004; avoids duplicated toggle/auth/theme across `landing-nav.tsx` and `ToolsHubPage.tsx`.
- **Alternatives**: Duplicate headers per feature — rejected (drift risk). Put chrome only in `layouts/` — rejected (`MainLayout` is not the landing/tools shell today; chrome is presentational and auth-aware).

### R5 — Locale Select + flags (replace toggle)

- **Decision**: shadcn `Select` with flag icon + “English” / “Tiếng Việt” labels; `i18n.changeLanguage('en'|'vi')`. Remove EN↔VI binary toggle from public P1 surfaces.
- **Rationale**: FR-005 / US6; accessible labels required (not flag-only).
- **Alternatives**: Keep toggle — rejected by spec. Flags-only without text — rejected (a11y).

### R6 — Public layout scroll split (retain)

- **Decision**: Keep thin `fe/src/app/(public)/layout.tsx`. Apply `LandingLayout` (`h-[100dvh] overflow-hidden`) **only** inside the landing page. `/tools` and `/tools/webhooks` use normal document scroll; optional `tools/layout.tsx` for shared tools chrome without viewport lock.
- **Rationale**: Prior good decision; still required so hub/tool pages are browsable.
- **Alternatives**: Nested `(landing)` route group — acceptable equivalent if implementer prefers clearer file structure.

### R7 — Admin home control

- **Decision**: Add an explicit home icon/control in admin sidebar chrome (`app-sidebar.tsx` and/or `app-sidebar-header.tsx`) linking to `/`. Do not overload the admin logo (currently `/admin`) as the only escape hatch unless product later unifies branding.
- **Rationale**: US11 / FR-012 — dedicated public-home control.
- **Alternatives**: Change AdminAppLogo link to `/` — rejected (breaks admin home navigation expectation).

### R8 — Motion stack

- **Decision**: GSAP ScrollTrigger remains P1 cinematic intent on `/`; gate with `useReducedMotion`.
- **Rationale**: Spec retains GSAP for the **new** landing.
- **Alternatives**: Framer-only scrolly — rejected for P1 rewrite cost.

### R9 — No BE work

- **Decision**: Zero `[BE]` tasks for Phase 1.
- **Rationale**: Presentation + IA only; auth reuse from `001-auth`.

---

## Data Model (Phase 1 design)

### ToolEntry

| Field | Type | Required | Notes |
|-------|------|----------|--------|
| `id` | `string` | yes | Stable key (e.g. `webhooks`) |
| `name` | `string` | yes | Display name |
| `description` | `string` | yes | Short blurb |
| `order` | `number` | yes | Ascending sort for hub |
| `status` | `'available' \| 'comingSoon'` | yes | Drives link vs non-link UI |
| `href` | `string` | no | e.g. `/tools/webhooks`; only navigate when `available` + present |
| `category` | `string` | no | Optional group label |

### ToolCatalog

- Ordered collection; helpers `getAllTools()`, `getToolsTeaser(limit?)`, empty-safe.
- Phase 1 must include **Webhooks** with `href: '/tools/webhooks'`.
- ≥1 `comingSoon` entry recommended to prove non-broken coming-soon UX.

### Public chrome

- Shared header model: Navigate targets (`/`, `/tools`, `/login`, `/register`), auth state (guest vs signed-in), theme control, locale Select (EN/VI).

### Landing content

- `fe/src/features/landing/content.ts` owns marketing copy + section ids.
- Tool **entries** never duplicated — always from `@/features/tools`.

### Entity relationships

```text
LandingPage (/) --teaser subset--> ToolCatalog
ToolsHubPage (/tools) --full list--> ToolCatalog
WebhooksToolPage (/tools/webhooks) --identity/copy--> tools content (+ optional catalog lookup by id)
ToolEntry.href --> /tools/{id} (P1 instance: webhooks)
Admin chrome --home icon--> /
PublicSiteHeader --mounted on--> /, /tools, /tools/webhooks
```

---

## Contracts

**No HTTP API contracts for P1.**

**UI / IA contract (for FE + QA):**

| Surface | Contract |
|---------|----------|
| `GET /` | Greenfield scrollytelling landing; shared chrome; Tools teaser → `/tools` |
| `GET /tools` | Hub listing from catalog; shared chrome; Webhooks entry → `/tools/webhooks` |
| `GET /tools/webhooks` | Tool shell (not 404); back to `/tools`; shared chrome |
| Locale | Select + flag + label; EN/VI; no binary toggle on P1 public surfaces |
| Auth chrome | Guest: Login/Register; Signed-in: Logout via existing auth |
| Theme | Persists across the three P1 routes |
| Admin | Home icon → `/` in one activation |
| Catalog append | Edit `catalog.ts` → hub (and teaser) update after reload |

Types use TypeScript camelCase field names matching the entity table.

---

## Project Structure

### Documentation (this feature)

```text
docs/features/003-landing-page/
├── spec.md
├── tasks.md
├── plan.md                 ← this file
├── fe-tasks-verify.md      ← obsolete until @fe rewrites after re-implement
└── qa-checklist.md         ← @qa later
```

### Source Code (FE — concrete paths)

```text
fe/src/
├── app/
│   ├── (public)/
│   │   ├── layout.tsx                 # thin — NO landing overflow lock for all children
│   │   ├── page.tsx                   # metadata + <LandingPage />
│   │   └── tools/
│   │       ├── layout.tsx             # optional: shared tools chrome, document scroll
│   │       ├── page.tsx               # metadata + <ToolsHubPage />
│   │       └── webhooks/
│   │           └── page.tsx           # NEW — metadata + webhooks shell
│   └── sitemap.ts                     # `/`, `/tools`, `/tools/webhooks`
├── layouts/
│   └── landing-layout.tsx             # viewport lock — landing-only
├── components/common/
│   ├── public-site-header.tsx         # NEW — Navigate, auth, theme, locale
│   ├── locale-select.tsx              # NEW — Select + flags EN/VI
│   ├── theme-toggle.tsx               # existing
│   ├── app-sidebar.tsx                # add public-home control
│   └── app-sidebar-header.tsx         # optional home action slot
├── features/
│   ├── landing/                       # GREENFIELD rebuild (reuse motion primitives)
│   │   ├── LandingPage.tsx
│   │   ├── LandingPage.test.tsx
│   │   ├── content.ts
│   │   ├── components/
│   │   │   ├── landing-nav.tsx        # slim/replace → PublicSiteHeader (+ optional section jumps)
│   │   │   ├── scrolly-shell.tsx
│   │   │   └── sections/              # new/rebuilt sections; retire tokenomics usage
│   │   └── hooks/
│   │       ├── use-gsap-context.ts
│   │       └── use-reduced-motion.ts
│   └── tools/
│       ├── index.ts
│       ├── types.ts
│       ├── catalog.ts                 # include webhooks href
│       ├── catalog.test.ts
│       ├── content.ts
│       ├── ToolsHubPage.tsx           # use PublicSiteHeader; drop locale toggle
│       └── components/
│           ├── tool-card.tsx
│           └── webhooks-tool-page.tsx # NEW shell (or tool-detail-shell.tsx)
└── locales/en|vi/                     # chrome/tools/landing strings as needed
```

**Structure Decision**: App Router `(public)` for three routes; shared chrome in `components/common`; domain UI in `features/landing` + `features/tools`. No `be/` changes.

---

## Architecture Approach

### Blast radius

| Area | Risk | Notes |
|------|------|--------|
| `LandingPage` / sections / content | Medium–High local | Greenfield rebuild; tests + nav/HUD must follow new section ids |
| `landing-nav` → `PublicSiteHeader` | Medium | Auth/theme/locale consumers unify |
| `ToolsHubPage` / catalog | Medium | Webhooks entry + header swap; prior placeholder tools may be cleaned |
| New `/tools/webhooks` | Low–Medium | Additive route + shell |
| `(public)/layout` + `LandingLayout` | Medium | Must keep lock off tools routes |
| Admin sidebar home control | Low | Additive control; do not break `/admin` logo nav |
| Locale toggle removal | Medium | Grep public surfaces for binary `changeLanguage` toggles |

No HIGH/CRITICAL cross-package BE impact for Phase 1.

### Composition flow

```text
/  → page.tsx → LandingPage → LandingLayout + PublicSiteHeader (+ optional section jumps)
                            → ScrollyShell + GSAP (if !reducedMotion)
                            → Hero / Value / Tools teaser / CTA
                            → getToolsTeaser() from @/features/tools

/tools → tools/page.tsx → ToolsHubPage → PublicSiteHeader + getAllTools() → ToolCard
                                         Webhooks href → /tools/webhooks

/tools/webhooks → webhooks/page.tsx → WebhooksToolPage → PublicSiteHeader + back to /tools

Admin → AppSidebar (+ header) → Home icon → /
```

### Accessibility

- Reduced motion: disable/simplify scrubbed/pinned GSAP; CTAs remain actionable.
- Locale Select: keyboard operable; options expose text labels (flags are decorative or accompanied by text).
- Coming-soon tools: never broken empty links.
- Admin home control: accessible name (not icon-only without label/title).

### Out of scope (enforce in review)

- Full webhook product/runtime
- Additional nested tools beyond Webhooks (except coming-soon catalog placeholders)
- BE CMS
- New locales beyond EN/VI
- Reviving abandoned docs folders
- Shipping the retired tokenomics prototype as the home story
- Treating prior `fe-tasks-verify.md` as completion evidence

---

## Quickstart Validation Guide

**Prerequisites**: FE deps installed; `pnpm dev` in `fe/` or `make up-d`.

1. Open `/` — confirm **new** landing composition (not tokenomics-primary); motion enabled → scroll-linked transitions across ≥3 boundaries.
2. Enable reduced motion — reload `/` — content + CTAs usable.
3. Confirm shared chrome: Home, Tools, Login/Register (guest), Theme, **Locale Select** (not toggle).
4. Switch locale EN↔VI via Select on `/`, then confirm on `/tools` and `/tools/webhooks`.
5. Toggle theme on `/` → navigate to `/tools` → theme persists.
6. Sign in → Logout appears; logout returns guest chrome.
7. `/tools` lists Webhooks → open `/tools/webhooks` → back to `/tools`.
8. Direct-load `/tools/webhooks` — not 404.
9. Admin: home icon → `/`.
10. Append a coming-soon catalog entry → hub lists it without redesign.
11. Run `make test-fe`.

---

## Implementation Handoff

| Role | Next |
|------|------|
| `@fe` | Execute `[FE]` tasks in [tasks.md](./tasks.md) (recommended order in tasks); **rewrite** `fe-tasks-verify.md` (prior file obsolete) |
| `@qa` | After new FE verify: T045–T046 checklist + `make test` |
| `@be` | Not required for Phase 1 |

**Ready for `@fe`**: **Yes** — `tasks.md` + `plan.md` regenerated for the greenfield spec.

**Recommended start order**: Foundational chrome + catalog + webhooks route → US3–US6 shared chrome → US7–US8 hub/shell → US1/US2/US9 greenfield landing + teaser → US10 → US11 → polish + rewrite verify.
