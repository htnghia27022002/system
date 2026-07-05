'use client'

import { usePathname, useRouter, useSearchParams } from 'next/navigation'
import { useCallback, useMemo } from 'react'

function normalizeQueryString(params: URLSearchParams): string {
  const entries = [...params.entries()].sort(([a], [b]) => a.localeCompare(b))
  return new URLSearchParams(entries).toString()
}

export function useListQueryParams(keys: readonly string[]) {
  const router = useRouter()
  const pathname = usePathname()
  const searchParams = useSearchParams()

  const values = useMemo(() => {
    const out: Record<string, string> = {}
    for (const key of keys) {
      out[key] = searchParams.get(key) ?? ''
    }
    const page = Number(searchParams.get('page') ?? '1')
    out.page = String(Number.isFinite(page) && page > 0 ? page : 1)
    return out
  }, [keys, searchParams])

  const setParams = useCallback(
    (updates: Record<string, string | undefined>) => {
      const params = new URLSearchParams(searchParams.toString())
      const shouldResetPage =
        !('page' in updates) &&
        Object.keys(updates).some((key) => key !== 'page')

      for (const [key, value] of Object.entries(updates)) {
        if (value === undefined || value === '') {
          params.delete(key)
        } else {
          params.set(key, value)
        }
      }

      if (shouldResetPage) {
        params.delete('page')
      }

      const nextQuery = normalizeQueryString(params)
      const currentQuery = normalizeQueryString(searchParams)

      if (nextQuery === currentQuery) {
        return
      }

      router.push(nextQuery ? `${pathname}?${nextQuery}` : pathname)
    },
    [pathname, router, searchParams],
  )

  return { values, setParams }
}
