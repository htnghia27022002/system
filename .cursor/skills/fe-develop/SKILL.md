---
name: fe-develop
description: General frontend implementation skill for React + TypeScript in this repository. Use for feature/page/layout implementation, architecture-safe file placement, import boundaries, accessibility baseline, and quality checks. When styling UI, components, layouts, or marketing surfaces, follow fe/DESIGN.md as the project design system. For extra anti-slop polish on landing pages, also read design-taste-frontend.
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

**Project design authority:** For any UI styling, component design, layout composition, or visual refactor in `fe/`, `fe/DESIGN.md` is mandatory and overrides generic defaults. Do not duplicate `DESIGN.md` here; read the file and apply its tokens, components, and do/don't rules.

Do not duplicate all design-taste rules inside this file. Keep this file lean and implementation-focused.

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
4. If the task changes UI (components, pages, layouts, styles, tokens), read `fe/DESIGN.md` first and implement against it.
5. If the request is ambiguous, ask one focused clarifying question.

## Architecture Rules (Repository-Specific)

- Follow the `src` structure defined in `fe/README.md`.
- Keep business logic inside `src/features/<feature-name>`.
- Keep route-level composition in `src/pages`.
- Keep shared reusable UI in `src/components`.
- Keep global API plumbing in `src/services`.
- Avoid creating new top-level folders unless explicitly requested.

## Feature Boundary Rules

- Do not deep-import another feature's internals.
- Expose feature public API via `index.ts` when appropriate.
- Keep feature-specific types/services/hooks inside the feature folder.
- Keep shared utilities generic and business-agnostic.

## Strict Import Matrix (Enforced)

Use this matrix for every import decision:

- `src/app/**/*` can import from: `screens`, `features`, `components`, `layouts`, `hooks`, `services`, `store`, `utils`, `types`, `config`.
- `src/screens/*` can import from: `features`, `components`, `layouts`, `hooks`, `services`, `store`, `utils`, `types`, `config`.
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

## Project Design System (`fe/DESIGN.md`)

`DESIGN.md` is the **visual reference** (tokens, patterns, do/don't). Apply it **through** the library stack — not by replacing it.

Apply this whenever you design or restyle UI in `fe/`.

### When to use
- New pages, layouts, marketing tiles, auth screens, nav, footer, cards, buttons, forms.
- Visual refactors, spacing/typography/radius/color updates, component variants.
- Any task where the user says "design", "style", "UI polish", or "make it look like the spec".

### Implementation rules
- Map DESIGN patterns onto **existing** components (e.g. `store-utility-card` → `Card` + tokenized `className`; `button-primary` → `Button variant="pill"`; forms → `Input` + `Label`).
- Put tokens in `src/styles/design-tokens.css` and shadcn theme aliases; extend component **variants** when a pattern repeats — do not add parallel component CSS layers.
- Use the documented color, typography, spacing, radius, and elevation tokens. Prefer CSS variables / Tailwind theme tokens aligned with `DESIGN.md`; avoid random hex values.
- Single accent only: Action Blue (`#0066cc` / `{colors.primary}`) for interactive elements on light surfaces.
- Typography: body at 17px / weight 400; headlines at weight 600 with negative letter-spacing; no weight 500.
- No decorative gradients, no card/button shadows (product imagery shadow only), no second accent color.
- Button grammar: pill radius for primary CTAs; `active:scale-[0.95]` for press feedback; visible focus ring on keyboard focus.
- Section rhythm: alternate light/parchment/dark full-bleed tiles; use surface color change as dividers, not heavy borders.
- Font stack: `SF Pro Display`, `SF Pro Text`, with `system-ui` / `-apple-system` fallbacks; Inter only as off-Apple substitute per `DESIGN.md`.

### Conflict resolution
1. Installed libraries + their docs (shadcn/ui, Radix, next-themes, etc.)
2. `fe/DESIGN.md` (visual tokens and patterns applied **via** those components)
3. `fe/README.md` + `fe/AGENTS.md` (structure and boundaries)
4. `design-taste-frontend` skill (anti-slop guardrails, not a replacement for `DESIGN.md`)

If `design-taste-frontend` conflicts with `DESIGN.md`, follow `DESIGN.md` for this repository. Never satisfy `DESIGN.md` by dropping the UI library.

## UI and Layout Discipline

- Avoid generic templated output (for example: repetitive three-equal-card sections by default).
- Prefer clear visual hierarchy over decorative noise.
- Use consistent spacing, typography, radius, and color strategy from `fe/DESIGN.md`.
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
- [ ] UI work followed `fe/DESIGN.md` via tokens/variants on existing components.
- [ ] For extra landing/marketing polish, `design-taste-frontend` was checked when relevant.
- [ ] UI is coherent (type, spacing, color, radius, components).
- [ ] Interaction states exist (loading, empty, error, focus/active).
- [ ] Accessibility baseline is met (contrast, labels, keyboard focus, reduced motion).
- [ ] Responsive behavior is validated.
- [ ] No unrelated refactor introduced.
- [ ] Content written in English for modified/created files.
