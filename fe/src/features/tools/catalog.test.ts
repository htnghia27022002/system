import { describe, expect, it } from 'vitest'

import {
  getAllTools,
  getToolById,
  getToolsTeaser,
  selectToolsTeaser,
  sortTools,
  toolCatalog,
} from './catalog'
import type { ToolEntry } from './types'

describe('tool catalog helpers', () => {
  it('returns tools sorted by ascending order', () => {
    const tools = getAllTools()
    expect(tools.length).toBeGreaterThan(0)
    for (let i = 1; i < tools.length; i += 1) {
      expect(tools[i].order).toBeGreaterThanOrEqual(tools[i - 1].order)
    }
  })

  it('returns a teaser subset capped by limit', () => {
    const all = getAllTools()
    const teaser = getToolsTeaser(2)
    expect(teaser).toHaveLength(Math.min(2, all.length))
    expect(teaser.map((t) => t.id)).toEqual(all.slice(0, 2).map((t) => t.id))
  })

  it('is empty-safe for teaser when limit is zero', () => {
    expect(getToolsTeaser(0)).toEqual([])
  })

  it('handles an empty catalog without throwing', () => {
    const empty: ToolEntry[] = []
    expect(sortTools(empty)).toEqual([])
    expect(selectToolsTeaser(empty)).toEqual([])
    expect(selectToolsTeaser(empty, 5)).toEqual([])
  })

  it('exposes Webhooks as available with nested href', () => {
    const webhooks = getToolById('webhooks')
    expect(webhooks).toBeDefined()
    expect(webhooks?.status).toBe('available')
    expect(webhooks?.href).toBe('/admin/tools/webhooks')
    expect(webhooks?.description.toLowerCase()).toContain('http')
  })

  it('exposes at least one available and one comingSoon entry', () => {
    expect(toolCatalog.some((t) => t.status === 'available')).toBe(true)
    expect(toolCatalog.some((t) => t.status === 'comingSoon')).toBe(true)
  })

  it('does not give comingSoon entries a navigable tools-hub self href', () => {
    for (const tool of toolCatalog) {
      if (tool.status === 'comingSoon') {
        expect(tool.href).toBeUndefined()
      }
    }
  })
})
