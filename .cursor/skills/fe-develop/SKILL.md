---
name: fe-develop
description: General frontend implementation skill for React + TypeScript in this repository. Use for feature/page/layout implementation, architecture-safe file placement, import boundaries, accessibility baseline, and quality checks. Theme via fe/src/styles/index.css. For landing/marketing polish, also read design-taste-frontend.
disable-model-invocation: false
---

# FE Develop Skill

## Skill Routing (Indexing Optimization)

This skill is the default FE implementation skill.  
For best indexing and best output quality, route tasks as follows:

- Use `fe-develop` for architecture, feature coding, refactoring safety, and repository conventions.
- Use `design-taste-frontend` for design-direction-heavy tasks:
  - landing page visual direction
  - portfolio aesthetics
  - marketing page composition
  - high-fidelity redesign polish

When a request matches the design-heavy cases above, also read and apply:
- `.cursor/skills/design-taste-frontend/SKILL.md`

**Project design authority:** For UI styling, use `fe/src/styles/index.css` CSS variables and shadcn components in `fe/src/components/ui/`. For marketing/landing polish, also read `.cursor/skills/design-taste-frontend/SKILL.md`.

Do not duplicate design rules here. Keep this file lean and implementation-focused.

## Scope

Use this skill for:
- React + TypeScript UI implementation in `fe/`
- Feature/page/layout development
- Product-facing frontend refinement
- Incremental redesigns where existing UI must be improved without random rewrites

Do not use this skill as the primary pattern for:
- Backend logic
- Database design
- Dense enterprise data-table workflows as the core task
- Full aesthetic direction for landing/portfolio pages (use `design-taste-frontend`)

## Mandatory Context First

Before editing code:
1. Read the user brief and infer page/product intent.
2. Confirm the target module and route surface.
3. Respect `fe/README.md` and `fe/AGENTS.md` as architectural policy.
4. Theme and tokens: `fe/src/styles/index.css` (`:root` / `.dark`).
5. If the request is ambiguous, ask one focused clarifying question.

## Architecture Rules (Repository-Specific)

- Follow the `src` structure defined in `fe/README.md` and `fe/AGENTS.md`.
- Keep business logic inside `src/features/<feature-name>`.
- Keep route wiring in `src/app/` (thin pages + metadata only).
- Keep shared reusable UI in `src/components/`.
- Keep global API plumbing in `src/services/`.
- Avoid creating new top-level folders unless explicitly requested.

## Feature Boundary Rules

- Do not deep-import another feature's internals.
- Expose feature public API via `index.ts` when appropriate.
- Keep feature-specific types/services/hooks inside the feature folder.
- Keep shared utilities generic and business-agnostic.

## Strict Import Matrix (Enforced)

Use this matrix for every import decision:

- `src/app/**/*` can import from: `features`, `components`, `layouts`, `hooks`, `services`, `store`, `utils`, `types`, `config`.
- `src/layouts/*` can import from: `components`, `hooks`, `services`, `utils`, `types`, `config`.
- `src/features/<feature>/*` can import from: same feature, `components`, `hooks`, `services`, `utils`, `types`, `config`.
- `src/components/*` can import from: `hooks`, `utils`, `types`, `config`, `services` (only if truly generic).
- `src/hooks/*` can import from: `services`, `utils`, `types`, `config`.
- `src/services/*` can import from: `utils`, `types`, `config`.
- `src/store/*` can import from: `features/*/index.ts`, `services`, `utils`, `types`, `config`.

If an import direction is not listed, treat it as forbidden unless the user explicitly asks for an exception.

## Forbidden Import Patterns (Hard Ban)

- Do not import `src/features/<other-feature>/<internal-file>` from outside that feature.
- Do not import from `pages` into `features`, `components`, `hooks`, `services`, or `store`.
- Do not import from `layouts` into `features` internals unless explicitly approved.
- Do not let `components` depend on route-specific files from `pages`.
- Do not let `services` import from `pages`, `layouts`, `components`, or feature UI files.
- Do not create circular dependencies across `features`, `services`, and `store`.

Bad examples:
- `import { X } from "@/features/auth/components/LoginForm"` from `src/features/user/*`
- `import { DashboardPage } from "@/pages/dashboard-page"` from `src/features/auth/*`
- `import { Header } from "@/layouts/main-layout"` from `src/services/api-client.ts`

Good examples:
- `import { useAuth } from "@/features/auth"` from `src/pages/login-page.tsx`
- `import { apiClient } from "@/services/api-client"` from `src/features/auth/services/auth-api.ts`
- `import { formatDate } from "@/utils/formatDate"` from `src/components/common/DateBadge.tsx`

## UI Library First (Required)

- **Default stack:** shadcn/ui in `src/components/ui`, Radix primitives, TanStack Query, next-themes, lucide-react icons — use these before custom markup.
- **Follow official docs** for each library (shadcn component API, Radix behavior, next-themes `ThemeProvider` / `useTheme`, etc.).
- **Registry checklist (every new UI surface):**
  1. List what you need (e.g. loading bar, modal, tabs).
  2. Grep `fe/src/components/ui/` and [shadcn docs](https://ui.shadcn.com/docs/components).
  3. If absent: `cd fe && pnpm dlx shadcn@latest add <name> -y`, fix imports to `@/lib/utils`.
  4. Compose in `components/common/` or features using those primitives only.
- **Do not bypass the library** with parallel UI (raw `<div>` loading bars, `<button>`, bespoke `.ds-*` layers) when shadcn covers the pattern — e.g. nav loading → `Progress`, not a custom animated div.
- **`components/common` is for composition**, not a second design system.

## Project design system (`src/styles/index.css`)

Theme authority is **`fe/src/styles/index.css`** — edit `:root` and `.dark` CSS variables (`--primary`, `--background`, `--radius`, etc.).

Apply styling **through** shadcn/ui components and Tailwind classes that reference those tokens — not parallel custom layers.

### Conflict resolution

1. Installed libraries + their docs (shadcn/ui, Radix, next-themes)
2. `fe/src/styles/index.css` (theme tokens)
3. `fe/README.md` + `fe/AGENTS.md` (structure and boundaries)
4. `design-taste-frontend` skill (anti-slop for landing/marketing only)

## UI and Layout Discipline

- For admin CRUD forms in dialogs/sheets, follow `.cursor/rules/fe-form-layout.mdc` (2-column pairing, sections, full-width sensitive fields).
- Avoid generic templated output (for example: repetitive three-equal-card sections by default).
- Prefer clear visual hierarchy over decorative noise.
- Use consistent spacing, typography, radius, and color from CSS variables in `index.css`.
- Use Grid for multi-column composition when layout relationships matter.
- Keep hero content concise and visible without forced scrolling when building landing sections.

## Motion and Interaction Rules

- Animate only `transform` and `opacity`.
- Honor `prefers-reduced-motion` for non-trivial motion.
- Use meaningful motion (hierarchy, feedback, state transition), never motion for decoration only.
- Provide explicit hover, active, focus, loading, empty, and error states for interactive surfaces.

## Accessibility and Quality Baseline

- Ensure readable contrast for text, buttons, and form elements.
- Do not rely on placeholder as label in forms.
- Keep keyboard focus visible.
- Respect responsive behavior from mobile-first up to desktop.
- Avoid unstable viewport hacks for full-screen sections; prefer modern viewport-safe units.

## Design System and Dependency Rules

- If a specific design system is requested, use its official package.
- Do not mix multiple design systems in one surface unless explicitly required.
- Verify dependencies before importing new libraries.
- If a package is missing, provide the install command and then implement.

## Anti-Sloppy Output Rules

- No filler marketing copy or generic placeholders that reduce credibility.
- No inconsistent visual language between sections.
- No random style shifts between warm/cool palettes or mixed corner systems.
- No hand-rolled icon paths when a standard icon library is available.
- No fake "screenshot-like" div mosaics when real assets/placeholders are expected.

## Redesign Protocol

When task is a redesign:
1. Audit first: identify what must be preserved.
2. Keep brand semantics, improve structure and clarity.
3. Modernize progressively (typography, spacing, hierarchy, states, responsiveness).
4. Avoid destructive rewrites unless the user asks for full overhaul.

## Implementation Workflow

1. Infer intent and constraints.
2. Map change to correct folders/layers.
3. Implement minimal, focused edits.
4. Validate styling, responsiveness, accessibility, and state handling.
5. Run lint checks when changes are substantive.
6. Summarize what changed and why.

## Pre-Flight Checklist (Must Pass)

- [ ] Structure follows `fe/README.md` and `fe/AGENTS.md`.
- [ ] File placement and naming conventions are correct.
- [ ] Feature boundaries are respected (no deep cross-feature coupling).
- [ ] Imports pass strict matrix and forbidden-pattern checks.
- [ ] Imports use the `@/*` alias and feature public entries (no `../../../` chains, no deep feature imports).
- [ ] Tests are colocated next to the file under test (`*.test.ts` / `*.test.tsx`).
- [ ] Registry checked; shadcn component added or reused (no parallel hand-rolled UI where registry covers it).
- [ ] UI work uses `src/styles/index.css` tokens via shadcn/Tailwind (no random hex values).
- [ ] For extra landing/marketing polish, `design-taste-frontend` was checked when relevant.
- [ ] UI is coherent (type, spacing, color, radius, components).
- [ ] Interaction states exist (loading, empty, error, focus/active).
- [ ] Accessibility baseline is met (contrast, labels, keyboard focus, reduced motion).
- [ ] Responsive behavior is validated.
- [ ] No unrelated refactor introduced.
- [ ] Content written in English for modified/created files.
