# Tasks: Public Landing Rebuild + Tools IA + Shared Chrome (003-landing-page)

**Input**: [spec.md](./spec.md)  
**Plan**: [plan.md](./plan.md)  
**Prerequisites**: `spec.md` (Status: Draft — ready for tasks; **greenfield rebuild**)

**Stale artifacts**: Prior `tasks.md` / `plan.md` are **replaced** by this regeneration. Prior [`fe-tasks-verify.md`](./fe-tasks-verify.md) is **obsolete** — `@fe` must re-implement against these tasks and rewrite verify after completion. Do not treat prior checkboxes or verify claims as done for this revision.

**Tests**: Spec does not require TDD. Update/extend Vitest where routes, catalog, chrome, or landing sections change. QA owns manual acceptance against Independent Tests.

**Organization**: Phases follow user stories (US1–US11 in scope; US12 deferred). Task labels use `[FE]` / `[QA]` per constitution; `[USn]` maps to spec user stories. No `[BE]` for Phase 1.

## Format: `[ID] [Prefix] [P?] [Story?] Description`

- **`[FE]` / `[QA]`**: Role ownership (constitution)
- **`[P]`**: Parallelizable (different files, no incomplete dependency)
- **`[USn]`**: User story from spec (story phases only)
- Include exact file paths in descriptions

## Path Conventions

Frontend package: `fe/src/` (App Router + `features/`). No backend work for Phase 1.

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Reset delivery baseline for the greenfield revision; confirm stack and route targets.

- [x] T001 [FE] Confirm GSAP + ScrollTrigger remain available in `fe/package.json` and document that `fe/src/features/landing/` is a **greenfield rebuild target** (reuse motion primitives only where useful; discard retired prototype composition such as tokenomics as the primary story)
- [x] T002 [FE] [P] Confirm existing tools scaffolding under `fe/src/features/tools/` (`types.ts`, `catalog.ts`, `ToolsHubPage.tsx`, `index.ts`) will be **aligned** to this spec — not treated as finished
- [x] T003 [FE] [P] Confirm App Router targets: `fe/src/app/(public)/page.tsx` (`/`), `fe/src/app/(public)/tools/page.tsx` (`/tools`), and **new** `fe/src/app/(public)/tools/webhooks/page.tsx` (`/tools/webhooks`)

**Checkpoint**: Team agrees Phase 1 routes and greenfield intent; prior verify is ignored.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Shared public chrome primitives, catalog readiness for Webhooks, layout scroll split, and route shells. **BLOCKS all user stories.**

**⚠️ CRITICAL**: Keep landing viewport lock off `/tools` and `/tools/webhooks`. Do not wrap all `(public)` children in `LandingLayout`.

- [x] T004 [FE] Create shared locale control `fe/src/components/common/locale-select.tsx` using shadcn `Select` from `fe/src/components/ui/select.tsx` — options EN/VI with **flag icon + language label** (accessible text; not flag-only); wire `i18n.changeLanguage`; no binary EN↔VI toggle
- [x] T005 [FE] Create shared public chrome `fe/src/components/common/public-site-header.tsx` composing: brand/home link, Navigate (Home `/`, Tools `/tools`), auth slot (Login/Register or Logout), `LocaleSelect`, existing `ThemeToggle` from `fe/src/components/common/theme-toggle.tsx`; support optional mobile nav menu; accept optional `variant` / children for landing section links if needed
- [x] T006 [FE] Align catalog model in `fe/src/features/tools/types.ts` + `fe/src/features/tools/catalog.ts` with FR-009 (`id`, `name`, `description`, `order`, `status`, optional `href`, optional `category`); ensure helpers `getAllTools` / `getToolsTeaser` remain empty-safe
- [x] T007 [FE] Seed/update `fe/src/features/tools/catalog.ts` so **Webhooks** is a real `available` entry with `href: '/tools/webhooks'`; keep ≥1 `comingSoon` entry for extensibility demos; fix any misleading `href` values that point at `/tools` as if they were tool apps
- [x] T008 [FE] Add thin route `fe/src/app/(public)/tools/webhooks/page.tsx` (metadata + stub page from `@/features/tools`); optionally add `fe/src/app/(public)/tools/layout.tsx` for shared tools chrome if it keeps document scroll (no `100dvh` overflow lock)
- [x] T009 [FE] Confirm/retain layout split: `fe/src/app/(public)/layout.tsx` stays thin; `LandingLayout` in `fe/src/layouts/landing-layout.tsx` is applied **only** inside landing (`LandingPage` or equivalent) — `/tools` and `/tools/webhooks` use normal document scroll
- [x] T010 [FE] [P] Export tools public API from `fe/src/features/tools/index.ts` (catalog getters, types, hub page, webhooks/detail shell) so landing imports `@/features/tools` only
- [x] T011 [FE] [P] Add `/tools/webhooks` to `fe/src/app/sitemap.ts` alongside `/` and `/tools`

**Checkpoint**: LocaleSelect + PublicSiteHeader exist; Webhooks catalog + route stub resolve; layout split preserved.

---

## Phase 3: User Story 1 — New scrollytelling home (Priority: P1) 🎯 MVP motion

**Goal**: `/` is a **new** landing composition with GSAP scroll-driven storytelling and reduced-motion fallback; retired prototype sections (e.g. tokenomics) are not the primary story.

**Independent Test**: Load `/`; confirm new structure; scroll with and without reduced motion.

### Implementation

- [x] T012 [FE] [US1] Rebuild landing composition in `fe/src/features/landing/LandingPage.tsx` + `fe/src/features/landing/content.ts` for a greenfield section set (minimum: hero, value/features, tools teaser, cta — or equivalent cinematic set with ≥3 motion boundaries); remove retired tokenomics from `scrollySectionIds` / composition
- [x] T013 [FE] [US1] Delete or fully retire unused prototype section `fe/src/features/landing/components/sections/tokenomics-snap-section.tsx` (and related dead copy) so it cannot ship as the primary story
- [x] T014 [FE] [US1] Re-wire GSAP/ScrollTrigger storytelling via `fe/src/features/landing/hooks/use-gsap-context.ts` + section components under `fe/src/features/landing/components/sections/`; keep `useReducedMotion` in `fe/src/features/landing/hooks/use-reduced-motion.ts` disabling/simplifying scrubbed/pinned motion when preference is on
- [x] T015 [FE] [US1] Ensure scroll-linked transitions across ≥3 major section boundaries when motion is enabled (SC-002); polish so motion is intentional, not only static stacks
- [x] T016 [FE] [US1] Update `fe/src/features/landing/LandingPage.test.tsx` for the new section set and absence of tokenomics as required content

**Checkpoint**: New `/` storytelling works with motion on and reduced-motion path usable.

---

## Phase 4: User Story 2 — First-viewport hierarchy (Priority: P1)

**Goal**: First viewport shows brand/product, one headline, short support, ≥1 primary CTA; readable at desktop and 375px.

**Independent Test**: Open `/` at desktop and 375px; verify hierarchy without horizontal overflow of primary text/CTAs.

### Implementation

- [x] T017 [FE] [P] [US2] Implement hero hierarchy in `fe/src/features/landing/content.ts` + hero section component under `fe/src/features/landing/components/sections/` (brand signal, headline, support, primary CTA) — placeholder English OK
- [x] T018 [FE] [US2] Verify responsive first viewport (no horizontal page scroll for primary content/CTAs) in hero styles/layout; keep `/` wired solely via `fe/src/app/(public)/page.tsx` → `LandingPage`

**Checkpoint**: Hero hierarchy meets SC-001 / FR first-viewport intent independent of hub polish.

---

## Phase 5: User Story 3 — Shared chrome: Navigate (Priority: P1)

**Goal**: On `/`, `/tools`, `/tools/webhooks`, visitors can reach Home, Tools, and auth destinations via shared chrome.

**Independent Test**: From each P1 route, use nav to `/`, `/tools`, `/login` (or register).

### Implementation

- [x] T019 [FE] [US3] Mount `PublicSiteHeader` on landing (replace or slim `fe/src/features/landing/components/landing-nav.tsx` so site Navigate is shared — landing section jumps may remain as optional secondary links)
- [x] T020 [FE] [US3] Mount `PublicSiteHeader` on tools hub and webhooks shell (`fe/src/features/tools/ToolsHubPage.tsx` and webhooks page component) so Navigate destinations match FR-004
- [x] T021 [FE] [US3] Ensure mobile collapsed menu (if used) still exposes Home, Tools, Login/Register with accessible labeling in `fe/src/components/common/public-site-header.tsx`

**Checkpoint**: All three P1 routes share Navigate to Home/Tools/auth destinations.

---

## Phase 6: User Story 4 — Shared chrome: Login / Logout (Priority: P1)

**Goal**: Auth-aware chrome — guests see Login/Register; signed-in users see Logout (existing auth behavior).

**Independent Test**: Visit `/` as guest and signed-in; verify chrome states and logout clears session.

### Implementation

- [x] T022 [FE] [US4] Wire guest vs signed-in controls in `fe/src/components/common/public-site-header.tsx` using `useAuthStore` + existing `useSignOut` from `fe/src/features/auth/hooks/use-sign-out.ts`; links to `/login` and `/register`; never show Login + Logout as equal primaries
- [x] T023 [FE] [US4] Verify logout from public chrome ends session per existing auth rules and chrome returns to guest state on `/`, `/tools`, and `/tools/webhooks`

**Checkpoint**: Auth chrome matches SC-006 on all P1 public routes.

---

## Phase 7: User Story 5 — Shared chrome: Theme (Priority: P1)

**Goal**: Theme toggle in shared chrome; preference persists across P1 public routes.

**Independent Test**: Toggle theme on `/`; navigate to `/tools`; theme persists.

### Implementation

- [x] T024 [FE] [US5] Ensure `ThemeToggle` is always present in `PublicSiteHeader` and uses existing theme persistence (`next-themes` / current app pattern) so `/` → `/tools` → `/tools/webhooks` keep the same theme without losing place

**Checkpoint**: Theme persistence verified across the three P1 routes.

---

## Phase 8: User Story 6 — Locale Select with flags (Priority: P1)

**Goal**: Locale control is a Select with flag + label for EN/VI on all P1 public routes; EN↔VI toggle removed from those surfaces.

**Independent Test**: Open locale Select; choose VI then EN on `/`, `/tools`, `/tools/webhooks`.

### Implementation

- [x] T025 [FE] [US6] Integrate `LocaleSelect` into `PublicSiteHeader` and remove EN↔VI toggle buttons from landing/tools chrome (`landing-nav.tsx` if still present, `ToolsHubPage.tsx`, any duplicated public headers)
- [x] T026 [FE] [P] [US6] Add/adjust i18n strings for chrome/tools/landing as needed under `fe/src/locales/en/` and `fe/src/locales/vi/` (reuse existing language resources; do not invent new locales)
- [x] T027 [FE] [US6] Optionally align `fe/src/layouts/main-layout.tsx` locale control to `LocaleSelect` for consistency if that layout remains in use; admin user-menu language toggle may stay as-is unless touching it for consistency (public P1 routes are the hard requirement)

**Checkpoint**: Select+flags on all three P1 routes; old toggle gone on those surfaces (SC-005).

---

## Phase 9: User Story 7 — Tools hub + open Webhooks (Priority: P1)

**Goal**: `/tools` lists catalog (at least Webhooks); opening Webhooks reaches `/tools/webhooks`; hub remains expansion point.

**Independent Test**: Open `/tools`; open Webhooks; use back-to-tools.

### Implementation

- [x] T028 [FE] [US7] Implement/refresh hub listing UI in `fe/src/features/tools/ToolsHubPage.tsx` + `fe/src/features/tools/components/tool-card.tsx` from `getAllTools()` (name + description; status badge; available+href navigates)
- [x] T029 [FE] [US7] Ensure Webhooks catalog entry navigates to `/tools/webhooks` and hub chrome can return to `/` and stay on Tools IA
- [x] T030 [FE] [US7] Coming-soon entries remain visible without broken links (non-navigating control or disabled state — never empty `href`)

**Checkpoint**: Hub lists Webhooks and opens the nested route; extensible listing intact.

---

## Phase 10: User Story 8 — Webhooks tool page shell (Priority: P1)

**Goal**: `/tools/webhooks` is a real nested tool shell with identity, placeholder body, shared chrome, and clear back to `/tools`. Full webhook product logic out of scope.

**Independent Test**: Direct-load `/tools/webhooks`; confirm shell + back link; not 404.

### Implementation

- [x] T031 [FE] [US8] Create `fe/src/features/tools/components/webhooks-tool-page.tsx` (or `tool-detail-shell.tsx` specialized for webhooks) with title/identity, placeholder MVP body (EN/VI via i18n or content module), and explicit back-to-`/tools` control
- [x] T032 [FE] [US8] Wire `fe/src/app/(public)/tools/webhooks/page.tsx` to the shell + metadata; prove `/tools/{id}` pattern for future tools without requiring dynamic routes in P1
- [x] T033 [FE] [P] [US8] Add placeholder copy keys/content in `fe/src/features/tools/content.ts` and/or locale JSON for the webhooks shell

**Checkpoint**: Direct load of `/tools/webhooks` succeeds with back-to-hub (SC-008).

---

## Phase 11: User Story 9 — Landing Tools teaser → hub (Priority: P1)

**Goal**: Landing Tools teaser draws from shared catalog; primary CTA goes to `/tools`.

**Independent Test**: From `/` Tools section, primary CTA → `/tools`.

### Implementation

- [x] T034 [FE] [US9] Rebuild/align Tools teaser section (`fe/src/features/landing/components/sections/tools-snap-section.tsx` or replacement) to consume `getToolsTeaser()` from `@/features/tools` — no duplicated catalog entries in landing content
- [x] T035 [FE] [US9] Wire primary teaser CTA to `/tools` in `fe/src/features/landing/content.ts`; optional deep link to `/tools/webhooks` only if configured without replacing primary hub CTA

**Checkpoint**: Teaser + hub share one catalog; primary CTA hits `/tools`.

---

## Phase 12: User Story 10 — Extensible catalog (Priority: P1)

**Goal**: New tools via catalog entries without redesigning hub/landing chrome; future `/tools/{id}` pattern documented by Webhooks.

**Independent Test**: Config-add a second placeholder (coming soon OK); hub lists it; structure unchanged.

### Implementation

- [x] T036 [FE] [US10] Document append workflow in comments at top of `fe/src/features/tools/catalog.ts` (fields + hub/teaser refresh expectation)
- [x] T037 [FE] [US10] Extend/refresh `fe/src/features/tools/catalog.test.ts` to cover sort order, teaser subset, Webhooks href, and coming-soon non-link behavior
- [x] T038 [FE] [US10] Confirm landing teaser and hub both import catalog only via `@/features/tools` public barrel (no deep cross-feature internals)

**Checkpoint**: SC-010 satisfied; append-only extensibility clear for Phase 2+.

---

## Phase 13: User Story 11 — Admin → public home (Priority: P1)

**Goal**: Admin chrome includes an icon/control navigating to public home `/`.

**Independent Test**: As admin, activate home control → `/`.

### Implementation

- [x] T039 [FE] [US11] Add a home icon/control in admin chrome (prefer `fe/src/components/common/app-sidebar.tsx` header/actions and/or `fe/src/components/common/app-sidebar-header.tsx`) linking to `/` with accessible label (e.g. “Public home”)
- [x] T040 [FE] [US11] Ensure the control is one activation to `/` (SC-009) and does not replace admin dashboard navigation (`/admin` logo may remain admin-scoped)

**Checkpoint**: Admin users can escape to public `/` via a dedicated control.

---

## Phase 14: Polish & Cross-Cutting Concerns

**Purpose**: Cleanup, tests, SEO, and handoff hygiene after US1–US11.

- [x] T041 [FE] [P] Remove leftover EN↔VI toggle UX from public P1 surfaces; grep `fe/src/features/landing/` and `fe/src/features/tools/` for `changeLanguage(i18n.language === 'en'` toggle patterns and replace with `LocaleSelect` usage
- [x] T042 [FE] [P] Refresh metadata titles/descriptions on `fe/src/app/(public)/page.tsx`, `tools/page.tsx`, and `tools/webhooks/page.tsx` (English)
- [x] T043 [FE] Run `make test-fe` from repo root; fix regressions in landing/tools tests
- [x] T044 [FE] Rewrite [`fe-tasks-verify.md`](./fe-tasks-verify.md) after re-implement — prior verify file is obsolete and must not be reused as evidence
- [ ] T045 [QA] Manual acceptance against Independent Tests for US1–US11 and Success Criteria SC-001–SC-010 (after FE verify)
- [ ] T046 [QA] [P] Matrix: guest vs signed-in chrome; theme persist; locale Select EN↔VI on all three routes; reduced-motion `/`; admin home icon; direct `/tools/webhooks` load

---

## Phase 15: User Story 12 — Additional tools / rich apps (Priority: P2 — deferred)

**Goal**: Deferred — more catalog tools, full webhook product, filters/search, richer hub motion.

- [ ] T047 [FE] [US12] Deferred: additional `/tools/{id}` pages + real webhook product features (out of Phase 1 Done)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: Immediate
- **Foundational (Phase 2)**: Depends on Setup — **BLOCKS** US1–US11
- **User stories**: Prefer sequential order below for a single `@fe` agent; parallel notes where marked `[P]`
- **Polish**: After US1–US11 implementation
- **US12**: Deferred

### User Story Dependencies

| Story | Depends on | Notes |
|-------|------------|-------|
| US1, US2 | Phase 2 | Greenfield landing; can start once layout split + catalog helpers exist |
| US3–US6 | Phase 2 (`PublicSiteHeader`, `LocaleSelect`) | Chrome stories share one header component |
| US7, US8 | Phase 2 (catalog Webhooks + route) | Hub + nested shell |
| US9 | US1 section shell + catalog | Teaser on new landing |
| US10 | Catalog + hub/teaser wiring | Validation/tests |
| US11 | Admin sidebar/header only | Independent of landing motion |
| US12 | Phase 1 Done | Deferred |

### Recommended `@fe` start order

1. **T001–T011** Foundational (chrome primitives, catalog Webhooks, webhooks route, layout/sitemap)
2. **T019–T027** Shared chrome Navigate + Auth + Theme + Locale Select (US3–US6) on all P1 routes
3. **T028–T033** Tools hub + Webhooks shell (US7–US8)
4. **T012–T018** + **T034–T035** Greenfield landing + hero + Tools teaser (US1, US2, US9)
5. **T036–T038** Catalog extensibility proof (US10)
6. **T039–T040** Admin home icon (US11)
7. **T041–T044** Polish + rewrite `fe-tasks-verify.md`
8. Hand off to `@qa` for T045–T046

### Parallel Opportunities

- After T004–T005: T006–T007 (catalog) ∥ T008–T011 (routes/sitemap/exports)
- US11 (admin home) can run parallel to landing rebuild once FE capacity allows
- T026 locale strings ∥ T033 webhooks copy
- T041–T042 polish items in parallel after chrome is unified

---

## Parallel Example: Foundational chrome + catalog

```bash
Task: "Create locale-select.tsx in fe/src/components/common/locale-select.tsx"
Task: "Align catalog + seed Webhooks in fe/src/features/tools/catalog.ts"
Task: "Add tools/webhooks/page.tsx route stub"
Task: "Add /tools/webhooks to fe/src/app/sitemap.ts"
```

---

## Implementation Strategy

### MVP slice (chrome + IA first)

1. Phase 1–2 foundational
2. US3–US6 shared chrome on existing `/` and `/tools`
3. US7–US8 Webhooks nested page
4. **STOP**: Validate Navigate/Auth/Theme/Locale + hub + `/tools/webhooks`

### Full Phase 1

5. US1–US2 greenfield landing + US9 teaser
6. US10 catalog extensibility checks
7. US11 admin home
8. Polish + rewrite FE verify → `@qa`

### Notes

- Prior `fe-tasks-verify.md` is **obsolete**; never mark Phase 1 complete using it
- No `[BE]` tasks — FE presentation + IA only
- Do not revive `docs/tools` / `docs/langding-page`
- Do not ship full webhook runtime in P1
