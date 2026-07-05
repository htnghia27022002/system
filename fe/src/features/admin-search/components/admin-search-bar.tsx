'use client'

import Link from 'next/link'
import { Loader2Icon, SearchIcon, XIcon } from 'lucide-react'
import { useEffect, useRef, useState, type KeyboardEvent } from 'react'
import { useTranslation } from 'react-i18next'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useDebouncedValue } from '@/hooks/use-debounced-value'
import { cn } from '@/lib/utils'

import { useAdminSearch } from '../hooks/use-admin-search'
import { useSearchNavigation } from '../hooks/use-search-navigation'
import { entityTypeLabel } from '../lib/search-routes'
import { SearchEntityBadge } from './search-entity-badge'

const DEBOUNCE_MS = 200

export function AdminSearchBar() {
  const { t } = useTranslation('admin')
  const rootRef = useRef<HTMLDivElement>(null)
  const [query, setQuery] = useState('')
  const [activeQuery, setActiveQuery] = useState('')
  const [open, setOpen] = useState(false)

  const debouncedQuery = useDebouncedValue(query, DEBOUNCE_MS)

  const { canUseSearch, resolveHitHref } = useSearchNavigation()

  const { data, isLoading, isFetching, isError } = useAdminSearch({
    q: activeQuery,
    pageSize: 8,
  })

  useEffect(() => {
    function handlePointerDown(event: MouseEvent) {
      if (!rootRef.current?.contains(event.target as Node)) {
        setOpen(false)
      }
    }

    document.addEventListener('mousedown', handlePointerDown)
    return () => document.removeEventListener('mousedown', handlePointerDown)
  }, [])

  useEffect(() => {
    const trimmedDebounced = debouncedQuery.trim()
    const trimmedQuery = query.trim()

    // Wait until debounce matches the current input to avoid stale search/clear loops.
    if (trimmedDebounced !== trimmedQuery) {
      return
    }

    if (trimmedDebounced.length === 0) {
      setActiveQuery('')
      setOpen(false)
      return
    }

    setActiveQuery(trimmedDebounced)
    setOpen(true)
  }, [debouncedQuery, query])

  function runSearchImmediate(nextQuery: string) {
    const trimmed = nextQuery.trim()
    setActiveQuery(trimmed)
    setOpen(trimmed.length > 0)
  }

  function handleKeyDown(event: KeyboardEvent<HTMLInputElement>) {
    if (event.key === 'Enter') {
      event.preventDefault()
      runSearchImmediate(query)
    }
    if (event.key === 'Escape') {
      setOpen(false)
    }
  }

  function handleClear() {
    setQuery('')
    setActiveQuery('')
    setOpen(false)
  }

  const hits = data?.hits ?? []
  const showPanel = open && activeQuery.length > 0
  const showClear = Boolean(query.trim())
  const showSpinner = isFetching && activeQuery.length > 0
  const viewAllHref = canUseSearch
    ? `/admin/search?q=${encodeURIComponent(activeQuery)}`
    : null

  return (
    <div ref={rootRef} className="relative w-full min-w-0 max-w-sm md:max-w-lg">
      <div className="relative">
        <SearchIcon
          className="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-muted-foreground"
          aria-hidden
        />
        <Input
          value={query}
          onChange={(event) => setQuery(event.target.value)}
          onKeyDown={handleKeyDown}
          onFocus={() => {
            if (activeQuery) setOpen(true)
          }}
          placeholder={t('search.placeholder')}
          className={cn('h-9 pl-9', showClear || showSpinner ? 'pr-16' : 'pr-3')}
          aria-label={t('search.barLabel')}
          aria-expanded={showPanel}
          aria-controls="admin-search-results"
          aria-autocomplete="list"
          role="combobox"
        />
        <div className="absolute top-1/2 right-1 flex -translate-y-1/2 items-center gap-0.5">
          {showSpinner ? (
            <Loader2Icon
              className="size-4 animate-spin text-muted-foreground"
              aria-hidden
            />
          ) : null}
          {showClear ? (
            <Button
              type="button"
              variant="ghost"
              size="icon"
              className="size-7 shrink-0"
              onClick={handleClear}
              aria-label={t('search.clearSearch')}
            >
              <XIcon className="size-3.5" />
            </Button>
          ) : null}
        </div>
      </div>

      {showPanel ? (
        <div
          id="admin-search-results"
          className="absolute top-[calc(100%+0.5rem)] left-0 z-50 w-full min-w-[min(100%,20rem)] overflow-hidden rounded-xl border bg-popover text-popover-foreground shadow-lg"
          role="listbox"
        >
          {isLoading ? (
            <p className="flex items-center gap-2 px-3 py-4 text-sm text-muted-foreground">
              <Loader2Icon className="size-4 shrink-0 animate-spin" aria-hidden />
              {t('search.loading')}
            </p>
          ) : isError ? (
            <p className="px-3 py-4 text-sm text-destructive">{t('search.error')}</p>
          ) : hits.length === 0 ? (
            <p className="px-3 py-4 text-sm text-muted-foreground">{t('search.noResults')}</p>
          ) : (
            <ul className="max-h-80 overflow-y-auto py-1">
              {hits.map((hit) => {
                const href = resolveHitHref(hit)
                const meta = hit.metadata?.email ?? hit.metadata?.key ?? hit.metadata?.slug

                return (
                  <li key={`${hit.entityType}:${hit.entityId}`} className="border-b last:border-b-0">
                    {href ? (
                      <Link
                        href={href}
                        className="block px-3 py-2.5 transition-colors hover:bg-muted/60"
                        role="option"
                        onClick={() => setOpen(false)}
                      >
                        <SearchResultRow
                          entityType={hit.entityType}
                          title={hit.title}
                          meta={meta}
                          snippet={hit.snippet}
                        />
                      </Link>
                    ) : (
                      <div className="px-3 py-2.5" role="option">
                        <SearchResultRow
                          entityType={hit.entityType}
                          title={hit.title}
                          meta={meta}
                          snippet={hit.snippet}
                        />
                      </div>
                    )}
                  </li>
                )
              })}
            </ul>
          )}

          {hits.length > 0 && viewAllHref ? (
            <div className="border-t bg-muted/30 px-3 py-2">
              <Link
                href={viewAllHref}
                className="text-sm font-medium text-primary hover:underline"
                onClick={() => setOpen(false)}
              >
                {t('search.viewAllResults')}
              </Link>
            </div>
          ) : null}

          {data?.degraded ? (
            <p className="border-t px-3 py-2 text-xs text-amber-600 dark:text-amber-400">
              {t('search.degraded')}
            </p>
          ) : null}
        </div>
      ) : null}
    </div>
  )
}

function SearchResultRow({
  entityType,
  title,
  meta,
  snippet,
}: {
  entityType: Parameters<typeof SearchEntityBadge>[0]['entityType']
  title: string
  meta?: string
  snippet?: string
}) {
  return (
    <div className="flex flex-col gap-0.5">
      <div className="flex items-center gap-2">
        <SearchEntityBadge entityType={entityType} />
        <span className="min-w-0 truncate text-sm font-medium">{title}</span>
      </div>
      {meta ? (
        <p className="truncate text-xs text-muted-foreground">{meta}</p>
      ) : null}
      {snippet ? (
        <p className={cn('line-clamp-2 text-xs text-muted-foreground')}>{snippet}</p>
      ) : null}
      <span className="sr-only">{entityTypeLabel(entityType)}</span>
    </div>
  )
}
