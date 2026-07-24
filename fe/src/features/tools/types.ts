export type ToolStatus = 'available' | 'comingSoon'

export type ToolEntry = {
  /** Stable catalog key */
  id: string
  /** Display name */
  name: string
  /** Short blurb for teaser and hub */
  description: string
  /** Ascending sort order for hub listing */
  order: number
  /** Drives link vs unavailable UI */
  status: ToolStatus
  /** Internal path or absolute URL; navigate only when available + present */
  href?: string
  /** Optional label/group on hub (no P1 filters) */
  category?: string
}

export type ToolCatalog = readonly ToolEntry[]
