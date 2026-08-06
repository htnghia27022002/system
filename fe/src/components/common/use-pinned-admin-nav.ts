'use client'

import { useCallback, useEffect, useMemo, useState } from 'react'

import {
  ADMIN_HOME_HREF,
  useAdminNavLeaves,
} from '@/components/common/use-admin-nav-items'
import type { NavItem } from '@/types/navigation'

const STORAGE_KEY = 'admin.mobile.pinnedNav'
/** Pinned shortcuts in the bottom bar (Home + Menu are fixed). */
export const MAX_MOBILE_PINNED_NAV = 3

function readStoredPins(): string[] {
  if (typeof window === 'undefined') return []
  try {
    const raw = window.localStorage.getItem(STORAGE_KEY)
    if (!raw) return []
    const parsed = JSON.parse(raw) as unknown
    if (!Array.isArray(parsed)) return []
    return parsed.filter((value): value is string => typeof value === 'string')
  } catch {
    return []
  }
}

function writeStoredPins(hrefs: string[]) {
  if (typeof window === 'undefined') return
  window.localStorage.setItem(STORAGE_KEY, JSON.stringify(hrefs))
}

/**
 * User-defined mobile bottom-bar pins (Facebook-style shortcuts).
 * Home (`/admin`) is always fixed and cannot be pinned/unpinned here.
 */
export function usePinnedAdminNav() {
  const leaves = useAdminNavLeaves()
  const [pinnedHrefs, setPinnedHrefs] = useState<string[]>([])
  const [hydrated, setHydrated] = useState(false)

  useEffect(() => {
    setPinnedHrefs(readStoredPins())
    setHydrated(true)
  }, [])

  const availableByHref = useMemo(() => {
    const map = new Map<string, NavItem>()
    for (const item of leaves) {
      map.set(item.href, item)
    }
    return map
  }, [leaves])

  const pinnedItems = useMemo(() => {
    return pinnedHrefs
      .filter((href) => href !== ADMIN_HOME_HREF)
      .map((href) => availableByHref.get(href))
      .filter((item): item is NavItem => Boolean(item))
      .slice(0, MAX_MOBILE_PINNED_NAV)
  }, [pinnedHrefs, availableByHref])

  const persist = useCallback((next: string[]) => {
    const cleaned = next
      .filter((href) => href !== ADMIN_HOME_HREF)
      .filter((href, index, arr) => arr.indexOf(href) === index)
      .slice(0, MAX_MOBILE_PINNED_NAV)
    setPinnedHrefs(cleaned)
    writeStoredPins(cleaned)
  }, [])

  const isPinned = useCallback(
    (href: string) =>
      href !== ADMIN_HOME_HREF && pinnedHrefs.includes(href),
    [pinnedHrefs],
  )

  const canPin = useCallback(
    (href: string) => {
      if (href === ADMIN_HOME_HREF) return false
      if (!availableByHref.has(href)) return false
      if (isPinned(href)) return true
      return pinnedItems.length < MAX_MOBILE_PINNED_NAV
    },
    [availableByHref, isPinned, pinnedItems.length],
  )

  const pin = useCallback(
    (href: string) => {
      if (href === ADMIN_HOME_HREF || !availableByHref.has(href)) return
      if (pinnedHrefs.includes(href)) return
      if (pinnedHrefs.length >= MAX_MOBILE_PINNED_NAV) return
      persist([...pinnedHrefs, href])
    },
    [availableByHref, persist, pinnedHrefs],
  )

  const unpin = useCallback(
    (href: string) => {
      persist(pinnedHrefs.filter((value) => value !== href))
    },
    [persist, pinnedHrefs],
  )

  const togglePin = useCallback(
    (href: string) => {
      if (isPinned(href)) {
        unpin(href)
        return
      }
      pin(href)
    },
    [isPinned, pin, unpin],
  )

  return {
    hydrated,
    leaves,
    pinnedItems,
    isPinned,
    canPin,
    pin,
    unpin,
    togglePin,
    maxPins: MAX_MOBILE_PINNED_NAV,
  }
}
