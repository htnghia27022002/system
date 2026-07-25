/**
 * Shared tool catalog — single source of truth for landing teaser and `/tools` hub.
 *
 * Append workflow (Phase 1+):
 * 1. Add a `ToolEntry` to `TOOL_CATALOG` below.
 * 2. Required fields: `id`, `name`, `description`, `order`, `status`.
 * 3. Optional: `href` (navigate only when `status === 'available'`), `category`.
 * 4. For a real tool app, add a route and set `href` (Webhooks owner UI is under
 *    admin at `/admin/tools/webhooks`; public capture stays `/tools/webhooks/{uuid}`).
 * 5. Set `order` ascending for hub sort; teaser uses the same sorted list (first N).
 * 6. Reload `/` and `/tools` — both surfaces read this module via `@/features/tools` only.
 *    No hub/landing chrome redesign is required to list a new entry.
 */

import type { ToolCatalog, ToolEntry } from './types'

const TOOL_CATALOG: ToolCatalog = [
  {
    id: 'webhooks',
    name: 'Webhooks',
    description:
      'Capture and inspect inbound HTTP webhook requests with a personal public URL.',
    order: 10,
    status: 'available',
    href: '/admin/tools/webhooks',
    category: 'Integrations',
  },
  {
    id: 'json-formatter',
    name: 'JSON Formatter',
    description: 'Pretty-print, validate, and minify JSON in the browser.',
    order: 20,
    status: 'comingSoon',
    category: 'Developer',
  },
  {
    id: 'base64-codec',
    name: 'Base64 Codec',
    description: 'Encode and decode Base64 strings without leaving the page.',
    order: 30,
    status: 'comingSoon',
    category: 'Developer',
  },
  {
    id: 'uuid-generator',
    name: 'UUID Generator',
    description: 'Generate version 4 UUIDs for prototypes and test fixtures.',
    order: 40,
    status: 'comingSoon',
    category: 'Developer',
  },
]

/** Sort ascending by `order`. Empty-safe. */
export function sortTools(entries: readonly ToolEntry[]): ToolEntry[] {
  return [...entries].sort((a, b) => a.order - b.order)
}

/** Teaser subset helper. Empty-safe for empty catalogs or non-positive limits. */
export function selectToolsTeaser(entries: readonly ToolEntry[], limit = 3): ToolEntry[] {
  if (limit <= 0 || entries.length === 0) return []
  return sortTools(entries).slice(0, limit)
}

/** Full catalog in configured order. Empty-safe (never returns null/undefined). */
export function getAllTools(): ToolEntry[] {
  return sortTools(TOOL_CATALOG)
}

/**
 * Teaser subset for the landing Tools section.
 * Empty-safe: returns [] when the catalog is empty or limit is 0.
 */
export function getToolsTeaser(limit = 3): ToolEntry[] {
  return selectToolsTeaser(TOOL_CATALOG, limit)
}

/** Lookup a catalog entry by id. */
export function getToolById(id: string): ToolEntry | undefined {
  return TOOL_CATALOG.find((tool) => tool.id === id)
}

/** Raw catalog reference for tests / advanced use. Prefer getters in UI. */
export const toolCatalog: ToolCatalog = TOOL_CATALOG
