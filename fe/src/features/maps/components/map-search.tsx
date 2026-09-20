'use client'

import { MapPinnedIcon, SearchIcon, StoreIcon, GlobeIcon } from 'lucide-react'
import { useEffect, useId, useRef, useState } from 'react'
import { useTranslation } from 'react-i18next'

import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Spinner } from '@/components/ui/spinner'
import { cn } from '@/lib/utils'

import type { MapSearchHit, MapSearchKind } from '../map-search'

type MapSearchProps = {
  value: string
  hits: MapSearchHit[]
  loading?: boolean
  onChange: (value: string) => void
  onSelect: (hit: MapSearchHit) => void
}

const KIND_ICON: Record<MapSearchKind, typeof StoreIcon> = {
  place: StoreIcon,
  location: MapPinnedIcon,
  geocode: GlobeIcon,
}

export function MapSearch({
  value,
  hits,
  loading = false,
  onChange,
  onSelect,
}: MapSearchProps) {
  const { t } = useTranslation('admin')
  const listId = useId()
  const rootRef = useRef<HTMLDivElement | null>(null)
  const [open, setOpen] = useState(false)
  const [activeIndex, setActiveIndex] = useState(0)
  const showList = open && value.trim().length > 0

  useEffect(() => {
    setActiveIndex(0)
  }, [value, hits.length])

  useEffect(() => {
    if (!showList) return
    const onPointerDown = (event: PointerEvent) => {
      if (!rootRef.current?.contains(event.target as Node)) {
        setOpen(false)
      }
    }
    document.addEventListener('pointerdown', onPointerDown)
    return () => document.removeEventListener('pointerdown', onPointerDown)
  }, [showList])

  const select = (hit: MapSearchHit) => {
    onSelect(hit)
    onChange(hit.title)
    setOpen(false)
  }

  return (
    <div ref={rootRef} className="pointer-events-auto relative w-full max-w-sm">
      <div className="pointer-events-auto flex w-full items-center gap-2 rounded-lg border border-border bg-card/95 px-2.5 shadow-md backdrop-blur-sm">
        <SearchIcon className="size-4 shrink-0 text-muted-foreground" />
        <Label htmlFor="maps-filter-q" className="sr-only">
          {t('maps.filters.search')}
        </Label>
        <Input
          id="maps-filter-q"
          role="combobox"
          aria-expanded={showList}
          aria-controls={listId}
          aria-autocomplete="list"
          aria-activedescendant={
            showList && hits[activeIndex]
              ? `${listId}-${hits[activeIndex].id}`
              : undefined
          }
          value={value}
          onChange={(event) => {
            onChange(event.target.value)
            setOpen(true)
          }}
          onFocus={() => setOpen(true)}
          onKeyDown={(event) => {
            if (!showList) return
            if (event.key === 'ArrowDown') {
              event.preventDefault()
              setActiveIndex((index) =>
                hits.length === 0 ? 0 : Math.min(index + 1, hits.length - 1),
              )
            }
            if (event.key === 'ArrowUp') {
              event.preventDefault()
              setActiveIndex((index) => Math.max(index - 1, 0))
            }
            if (event.key === 'Enter' && hits[activeIndex]) {
              event.preventDefault()
              select(hits[activeIndex])
            }
            if (event.key === 'Escape') {
              event.preventDefault()
              setOpen(false)
            }
          }}
          placeholder={t('maps.filters.searchPlaceholder')}
          className="h-10 border-0 bg-transparent shadow-none focus-visible:border-transparent focus-visible:ring-0 dark:bg-transparent"
        />
        {loading ? <Spinner className="size-4 shrink-0" /> : null}
      </div>

      {showList ? (
        <ul
          id={listId}
          role="listbox"
          className="absolute top-[calc(100%+0.35rem)] z-[1200] max-h-80 w-full overflow-auto rounded-lg border border-border bg-card py-1 shadow-md"
        >
          {hits.length === 0 && !loading ? (
            <li className="px-3 py-2 text-sm text-muted-foreground">
              {t('maps.search.empty')}
            </li>
          ) : null}
          {hits.map((hit, index) => {
            const Icon = KIND_ICON[hit.kind]
            const active = index === activeIndex
            return (
              <li key={hit.id} role="presentation">
                <button
                  id={`${listId}-${hit.id}`}
                  type="button"
                  role="option"
                  aria-selected={active}
                  className={cn(
                    'flex w-full items-start gap-2 px-3 py-2 text-left text-sm',
                    active ? 'bg-accent text-accent-foreground' : 'hover:bg-muted/70',
                  )}
                  onMouseEnter={() => setActiveIndex(index)}
                  onClick={() => select(hit)}
                >
                  <Icon className="mt-0.5 size-4 shrink-0 text-muted-foreground" />
                  <span className="min-w-0">
                    <span className="block truncate font-medium">{hit.title}</span>
                    <span className="block truncate text-xs text-muted-foreground">
                      {hit.subtitle || t(`maps.search.kinds.${hit.kind}`)}
                    </span>
                  </span>
                </button>
              </li>
            )
          })}
        </ul>
      ) : null}
    </div>
  )
}
