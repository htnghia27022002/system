# Frontend Design System — Reference & Rules

> **Rule #0 — Research before you build.**
> Before implementing any screen or feature, search for established patterns online (MDN, Radix, shadcn/ui docs, UX Collective, Material Design, Pencil & Paper). Only start coding after confirming the best pattern for the use case. Reuse what already exists in `src/components/ui/` and `src/components/common/` before writing new markup.

---

## 1. Pre-Implementation Checklist

Before opening any source file to add a feature:

- [ ] **Search first** — look up the UX pattern (table, form, badge, dialog, etc.) on established design resources.
- [ ] **Audit `src/components/ui/`** — check whether shadcn/ui already provides the primitive. Run `pnpm dlx shadcn@latest add <name> -y` if it is missing from the registry.
- [ ] **Audit `src/components/common/`** — check whether a shared composition exists (e.g. `DataTable`, `StatusBadge`, `AccessControlPageHeader`).
- [ ] **Audit `src/features/<feature>/`** — check for existing hooks, services, or components that already do part of the job.
- [ ] **Read `src/styles/index.css`** — all colors must come from the CSS variable tokens defined there, not hard-coded hex/oklch values.
- [ ] Only after the above steps pass: write new code.

---

## 2. Status / Enum Badge Pattern

### Rule
Every enum value that appears in a table or form **must** have a dedicated `Record<EnumValue, BadgeConfig>` map. Never `if`/`switch` inside JSX for badge colors.

### Structure

```ts
// features/<feature>/enum-maps.ts
import type { StatusVariant } from '@/components/ui/status-badge'

type BadgeConfig = {
  variant: StatusVariant
  /** Whether to show the animated pulse dot (for "live" / "active" states only) */
  pulse: boolean
}

export const USER_STATUS_MAP: Record<ManagedUserStatus, BadgeConfig> = {
  active:   { variant: 'success', pulse: true  },
  inactive: { variant: 'neutral', pulse: false },
}
```

### Animated Dot Rule
- `active` / `online` / `live` states → `pulse: true` → renders `animate-ping` outer ring + solid inner dot.
- All other states → `pulse: false` → solid static dot.
- Always wrap ping animation with `motion-reduce:hidden` to respect the OS "reduce motion" preference.

### Color Semantic Mapping

| Variant   | Meaning                       | Use cases                        |
|-----------|-------------------------------|----------------------------------|
| `success` | Healthy / active / complete   | active user, published, paid     |
| `warning` | Needs attention               | pending, expiring, low stock     |
| `error`   | Failed / blocked / danger     | suspended, error, overdue        |
| `info`    | Informational / in-progress   | in review, syncing, draft        |
| `neutral` | Inactive / archived / unknown | inactive, archived, none         |

---

## 3. DataTable Design Rules

### Header
- Table header row uses a subtle `bg-muted/50` background — never the same surface as data rows.
- Header text: `text-xs font-semibold uppercase tracking-wide text-muted-foreground` to visually separate it from cell content.
- Column headers for sortable columns render `DataTableColumnHeader` with asc/desc/none icons.

### Primary Column (Name / Title)
- The first meaningful column (usually `name`, `title`, `email`) is the **clickable primary column**.
- It renders as a `<button>` styled as an underline-on-hover link that opens the edit dialog for that row.
- This removes the need for users to discover the actions menu just to edit a record.
- Pattern: `variant="link"` Button or plain `<button className="... hover:underline">` with `onClick={() => openEdit(row)}`.

### Row Density
- Default density for admin tables: compact — `TableRow` height ~36px via `py-1.5` cells.
- Do not use `py-4` (default shadcn) in admin tables.

### Pagination
- Always show pagination if `data.length > defaultPageSize`.
- Show row-count selector (`10 / 25 / 50`) and page navigation.
- On a single page, hide the pagination bar.

### Empty State
- Centered `InboxIcon` + descriptive message — never an empty `<tbody>`.

### Mobile
- Below `md` breakpoint, swap the table for a card list using `MobileRecordCard`.
- The card list shares the same data source and filter state as the desktop table.

---

## 4. Form Design Rules

### Mobile
- Use `Sheet` (bottom, full-width) on mobile and `Dialog` on desktop.
- Use `useIsMobile()` to switch between the two.
- Add a sticky footer with action buttons inside the sheet/dialog so they are always visible above the virtual keyboard.
- Add `interactive-widget=resizes-content` to the viewport meta tag.

### Layout
- Single-column on mobile, two-column on desktop (where it makes sense).
- Group related fields under `<fieldset>` or a visual section divider.
- All required fields must have a visible label and accessible `aria-required`.

---

## 5. Component Reuse Hierarchy

```
shadcn/ui registry         ← always check here first
  └── src/components/ui/   ← extended / customized primitives (StatusBadge, etc.)
        └── src/components/common/  ← cross-feature compositions (DataTable, PageHeader, etc.)
              └── src/features/<f>/components/  ← feature-specific UI
```

Never write a one-off badge, skeleton, or empty-state inside a feature component when a shared version already exists.

---

## 6. Animation Rules

- **`animate-ping`** — use only for "live" status dots. Never for decorative UI.
- **`animate-pulse`** — use only for skeleton loading placeholders.
- **`animate-spin`** — use only for spinner / loading states on async actions.
- Always pair any animation with `motion-reduce:hidden` or `motion-safe:` guards.

---

## 7. CSS Variable Token Rules

All colors must consume semantic CSS variables from `src/styles/index.css`:

```
bg-background      text-foreground
bg-card            text-card-foreground
bg-muted           text-muted-foreground
bg-primary         text-primary-foreground
bg-destructive     text-destructive-foreground
```

Do not introduce per-component hex values. For status/semantic colors not in the token set (e.g. `emerald` for success), use Tailwind palette classes — but isolate them in the `StatusBadge` component so they are changed in one place.

---

*Last updated: June 2026 — updated after DataTable + StatusBadge implementation.*
