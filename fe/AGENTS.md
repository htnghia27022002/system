# FE Agent Rules (React + TypeScript)

This file defines strict implementation rules for AI agents working in `fe/`.
Follow these rules before creating, editing, moving, or deleting files.

## 1) Source of Truth

- Use `README.md` in this folder as the architecture baseline.
- Use `DESIGN.md` in this folder as the visual design baseline for all UI work (colors, typography, spacing, components, layout rhythm).
- Do not introduce a different folder strategy unless explicitly requested by the user.
- Prefer consistency with existing patterns over personal preference.

## 2) Required `src` Structure

Agents must keep code organized under this structure:

- `src/assets`
- `src/components`
- `src/features`
- `src/hooks`
- `src/layouts`
- `src/pages`
- `src/services`
- `src/store`
- `src/styles`
- `src/types`
- `src/utils`
- `src/config`

Do not create random top-level folders in `src` without user approval.

## 3) Feature-First Rule

- Business/domain logic belongs in `src/features/<feature-name>/`.
- Each feature should own its local `components`, `hooks`, `services`, and `types`.
- Shared UI only goes to `src/components` when it is truly cross-feature.
- Do not place feature-specific logic in shared folders.

## 4) File Placement Rules

- Route screens go in `src/pages`.
- Layout shells go in `src/layouts`.
- Global API plumbing (base client/interceptors/adapters) goes in `src/services`.
- Global reusable hooks go in `src/hooks`.
- Global types go in `src/types`; feature types stay inside feature folders.
- Runtime/environment configuration goes in `src/config`.
- Generic helpers go in `src/utils`.
- App-wide styles/themes go in `src/styles`.
- Tests are colocated next to the file under test (`*.test.ts` / `*.test.tsx`), not in a separate top-level test tree.

## 5) Import and Boundary Rules

- Avoid deep imports into another feature's internals.
- Prefer exporting a feature's public surface via `src/features/<feature>/index.ts`.
- `pages` should compose from `features` and shared modules, not duplicate business logic.
- Keep dependency direction simple and predictable:
  - Shared/global layers can be used by features.
  - Features should not tightly couple to each other.
- Use the `@/*` path alias (mapped to `src/*`) instead of long relative chains (`../../../`).
- Import features through their public entry (`@/features/<feature>`), not deep internal files.

## 6) Naming Conventions

- Folder names: `kebab-case`
- React component files: `PascalCase.tsx`
- Hook files: `useXxx.ts`
- Utility files: `camelCase.ts`
- Keep names explicit and domain-oriented (e.g., `user-profile`, `auth-session`).

## 7) Editing Behavior for Agents

- Make minimal, targeted changes.
- Do not refactor unrelated modules in the same task.
- Do not move files across layers without a clear architectural reason.
- Do not add new dependencies unless the task requires them.
- Preserve backward compatibility unless breaking changes are requested.

## 8) Quality Guardrails

- Keep components focused (UI rendering first; move heavy logic into hooks/services).
- Avoid API calls directly inside deeply presentational components.
- Co-locate files that change together.
- Avoid excessive folder nesting unless complexity justifies it.
- Prefer readability and maintainability over clever abstractions.

## 9) UI Library and Design Reference

- **Library-first:** Use installed stack components (`src/components/ui` from shadcn/ui, Radix primitives, TanStack Query, next-themes, etc.) and follow each library’s official docs/API. Do not replace them with parallel hand-rolled UI (raw `<button>`, bespoke input classes, duplicate card shells).
- **Registry before markup:** For any UI primitive (progress bar, dialog, tabs, etc.), check `src/components/ui/` and [shadcn components](https://ui.shadcn.com/docs/components). If missing, run `pnpm dlx shadcn@latest add <name> -y` in `fe/` before writing custom elements. Files under `src/components/common/` must **compose** shadcn/ui, not duplicate it.
- **`fe/DESIGN.md` is a visual reference:** Map its patterns (`button-primary`, `store-utility-card`, `global-nav`, tokens) onto existing components via theme tokens, variants, and `className` — not by inventing separate `ds-*` component layers.
- **Tokens live in `src/styles/design-tokens.css`:** Bridge DESIGN colors/typography/radius to CSS variables and shadcn theme aliases. Extend shadcn variants (e.g. `Button` `variant="pill"`, `Progress` `indeterminate`) when a DESIGN pattern repeats.

## 10) Language Rule

- Any newly created or modified file content must be written in English.
- Do not mix Vietnamese and English in source text unless explicitly requested.

## 11) Before Finishing Any Task

Agents should self-check:

1. Is the file placed in the correct folder?
2. Does it respect feature boundaries?
3. Are naming conventions followed?
4. Did we avoid unrelated refactors?
5. Is the content in English?
6. For UI work: was a matching shadcn component used or added before custom markup?

If any answer is "no", fix it before finalizing.
