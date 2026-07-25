# FE Design Rules

Visual and interaction standards for `fe/`. Architecture, routing, and package boundaries stay in [`AGENTS.md`](./AGENTS.md). Agents implementing UI **must** follow this file.

**Theme source of truth:** `src/styles/index.css` (`:root` / `.dark` CSS variables).  
**UI primitives:** `src/components/ui/` (shadcn). Compose in `src/components/common/` or features — do not invent a parallel design system.

Related Cursor rules (detail): `.cursor/rules/fe-ui-library-first.mdc`, `.cursor/rules/fe-form-layout.mdc`.  
Marketing/landing polish only: `.cursor/skills/design-taste-frontend/SKILL.md`.

---

## 1) Theme & color consistency

- Customize appearance only via CSS variables in `src/styles/index.css` (`--primary`, `--background`, `--foreground`, `--muted`, `--destructive`, `--border`, `--radius`, …).
- All components (tables, cards, buttons, badges, dialogs, sidebars) must use the **same token set** — `bg-card`, `text-foreground`, `text-muted-foreground`, `text-primary`, `border-border`, etc.
- Do **not** hardcode one-off hex/rgb colors or invent parallel token layers (`.ds-*`, feature-local palettes).
- Light and dark modes must both remain readable; prefer semantic tokens over raw Tailwind color scales (`blue-500`, …) unless the token maps to a theme variable.
- Shared chrome (sidebar, header, page shell, cards) must look like one product: same radius, border treatment, and elevation pattern.

### Before custom markup

1. Check `src/components/ui/` for an existing shadcn component.
2. If missing: `cd fe && pnpm dlx shadcn@latest add <component> -y`.
3. Compose wrappers in `components/common/` or features — wrap primitives, do not reimplement them with raw `<div>` / `<button>`.

---

## 2) Page container & spacing (must match across pages)

Admin and product pages must share the **same outer padding and vertical rhythm** so screens do not feel like different apps.

| Surface | Standard |
|---------|----------|
| Admin page body (dashboard, users, roles, …) | Outer wrapper: `flex flex-1 flex-col gap-4 p-4` (match existing admin pages) |
| Card / table panel padding | Prefer `p-4` / `px-4` on small; align with siblings on the same page |
| Public content width (tools, marketing sections) | `mx-auto w-full max-w-6xl px-4 sm:px-6 lg:px-8` (or the layout already used by that product area — stay consistent within the area) |
| Stack spacing | Prefer `gap-4` / `space-y-4` / `space-y-6` — do not mix random `p-3` + `p-8` + `gap-2` on sibling pages |

**Rules**

- New admin pages must reuse the same outer `p-4` + `gap-4` pattern as `UsersTable` / dashboard overview.
- Do not invent a second page shell with different horizontal margin (`px-8` on one page, `px-3` on another).
- Dialogs/sheets: follow `.cursor/rules/fe-form-layout.mdc` (width, footer, grid). Keep header/body/footer padding aligned across create/edit dialogs.

---

## 3) Mobile-first (responsive web — not a native mobile app)

This product is a **responsive web app**. Optimize for phones and tablets; do **not** design as a separate native mobile app (no app-only chrome, bottom-tab-only IA, or desktop-only admin that breaks under ~375px).

- **Default layout is single column**; add multi-column only from `sm:` / `md:` upward.
- Tables: provide a usable mobile path (stacked cards, horizontal scroll with sticky primary column, or existing `MobileRecordCard` pattern) — do not leave a broken overflow-only desktop table with no mobile affordance.
- Touch targets: interactive controls should remain easy to tap (avoid tiny icon-only hit areas without padding).
- Forms: one column on mobile; 2 columns from `sm:` for paired fields (see form-layout rule).
- Sidebars/sheets: work with existing responsive sidebar / sheet patterns — do not assume permanent desktop sidebar width on small screens.
- Prefer `overflow-x-auto` / `min-w-0` where needed so content does not force horizontal page scroll.

---

## 4) Clickable table cells look like links

When a table cell value (**id**, **name**, email, title, …) opens a detail view, edit dialog, or navigates on click, it **must** look clickable at rest — not only on hover.

**Required appearance**

- Use **link color**: `text-primary` (theme token).
- Optional: `underline-offset-4 hover:underline` (or underline always if that matches nearby patterns).
- Keyboard: visible focus (`focus-visible:underline` / ring) — never `outline-none` without a replacement.
- Prefer `Link` from `next/link` for real routes; use `<button type="button">` styled the same way for in-page actions (open sheet/dialog).

```tsx
// ✅ GOOD — readable as a link before hover
<button
  type="button"
  className="text-left text-primary underline-offset-4 hover:underline focus-visible:underline"
  onClick={() => openEdit(row)}
>
  {row.name}
</button>

// ❌ BAD — looks like plain text; users cannot tell it is clickable
<button type="button" className="font-medium hover:text-primary">
  {row.name}
</button>
```

Non-clickable cells stay default body color (`text-foreground` / muted). Do not paint every column primary — only actionable values.

---

## 5) Form & dialog layout (summary)

Full detail: `.cursor/rules/fe-form-layout.mdc`.

- Mobile first; prefer 2 columns on `sm+` for paired fields; avoid uneven 3-column grids.
- Full-width: password, textarea, file upload, permission pickers, long selects.
- Visible `<Label>` on every field; hints as `text-xs text-muted-foreground`.
- Footer: Cancel (outline) + primary submit; sticky footer with `border-t bg-muted/30` when the body scrolls.

---

## 6) Motion & interaction

- Animate primarily `transform` and `opacity`; honor `prefers-reduced-motion`.
- Interactive surfaces need clear hover, active, focus, loading, empty, and error states.
- Do not use motion as decoration only.

---

## 7) Checklist before finishing UI work

1. Theme tokens only — no one-off palette drift across components.
2. Page outer padding/margin matches sibling pages in the same area (admin vs public).
3. Layout works on a ~375px-wide viewport (usable, not “desktop squished”).
4. Clickable id/name/… cells use **link color** (`text-primary`) so they read as clickable.
5. shadcn/registry checked before custom markup.
6. Copy and docs in English (UI strings via i18n when the feature already uses it).
