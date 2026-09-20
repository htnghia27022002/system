'use client'

import type { ReactNode } from 'react'
import {
  FilterXIcon,
  LayersIcon,
  LocateFixedIcon,
  MapPinnedIcon,
  TagsIcon,
} from 'lucide-react'
import { useTranslation } from 'react-i18next'

import { Button } from '@/components/ui/button'
import { AddressSelect } from '@/features/address'
import { Field, FieldLabel } from '@/components/ui/field'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import { Spinner } from '@/components/ui/spinner'
import { cn } from '@/lib/utils'

import { P1_PLACE_CATEGORIES } from '../category-filters'
import type { MapSearchHit } from '../map-search'
import { mapBasemapPreviewUrl } from '../providers/create-map-provider'
import type { MapBasemap } from '../providers/map-provider'
import { MapSearch } from './map-search'

export type MapFilterValues = {
  q: string
  category: string
  countryCode: string
  adminDivisionId: string
}

type MapFiltersProps = {
  values: MapFilterValues
  locating?: boolean
  onLocate?: () => void
  searchHits?: MapSearchHit[]
  searchLoading?: boolean
  onSearchSelect?: (hit: MapSearchHit) => void
  basemap: MapBasemap
  onBasemapChange: (basemap: MapBasemap) => void
  onChange: (next: MapFilterValues) => void
}

const ALL = 'all'

export function MapFilters({
  values,
  locating = false,
  onLocate,
  searchHits = [],
  searchLoading = false,
  onSearchSelect,
  basemap,
  onBasemapChange,
  onChange,
}: MapFiltersProps) {
  const { t } = useTranslation('admin')

  const categoryActive = Boolean(values.category && values.category !== ALL)
  const locationActive = Boolean(values.adminDivisionId)
  const hasFilters = Boolean(
    values.q.trim() || categoryActive || locationActive,
  )

  return (
    <div className="pointer-events-none absolute inset-x-3 top-3 z-[1100] flex flex-col items-start gap-2">
      <MapSearch
        value={values.q}
        hits={searchHits}
        loading={searchLoading}
        onChange={(q) => onChange({ ...values, q })}
        onSelect={(hit) => onSearchSelect?.(hit)}
      />

      <div className="pointer-events-auto flex flex-col overflow-hidden rounded-lg border border-border bg-card/95 shadow-md backdrop-blur-sm">
        <FilterPopoverButton
          label={t('maps.filters.category')}
          active={categoryActive}
          icon={<TagsIcon />}
        >
          <Field>
            <FieldLabel>{t('maps.filters.category')}</FieldLabel>
            <Select
              value={values.category || ALL}
              onValueChange={(category) => onChange({ ...values, category })}
            >
              <SelectTrigger aria-label={t('maps.filters.category')}>
                <SelectValue placeholder={t('maps.filters.allCategories')} />
              </SelectTrigger>
              <SelectContent className="z-[1400]">
                <SelectItem value={ALL}>
                  {t('maps.filters.allCategories')}
                </SelectItem>
                {P1_PLACE_CATEGORIES.map((key) => (
                  <SelectItem key={key} value={key}>
                    {t(`maps.categories.${key}`)}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </Field>
        </FilterPopoverButton>

        <FilterPopoverButton
          label={t('maps.filters.division')}
          active={locationActive}
          icon={<MapPinnedIcon />}
        >
          <AddressSelect
            allowEmpty
            value={{
              countryCode: values.countryCode,
              adminDivisionId: values.adminDivisionId,
            }}
            onChange={(next) =>
              onChange({
                ...values,
                countryCode: next.countryCode,
                adminDivisionId: next.adminDivisionId,
              })
            }
            contentClassName="z-[1400]"
          />
        </FilterPopoverButton>

        <FilterPopoverButton
          label={t('maps.basemap.title')}
          active={basemap !== 'street'}
          icon={<LayersIcon />}
          contentClassName="w-80"
        >
          <div className="grid gap-2">
            <p className="text-sm font-medium">{t('maps.basemap.title')}</p>
            <div className="grid grid-cols-2 gap-2">
              {(['street', 'satellite'] as const).map((mode) => {
                const selected = basemap === mode
                const preview = mapBasemapPreviewUrl(mode)
                return (
                  <button
                    key={mode}
                    type="button"
                    onClick={() => onBasemapChange(mode)}
                    aria-pressed={selected}
                    className={cn(
                      'overflow-hidden rounded-lg border text-left transition-colors',
                      selected
                        ? 'border-primary ring-2 ring-primary/30'
                        : 'border-border hover:border-foreground/30',
                    )}
                  >
                    {preview ? (
                      // Preview tiles already used by the map adapter.
                      // eslint-disable-next-line @next/next/no-img-element
                      <img
                        src={preview}
                        alt=""
                        className="h-20 w-full object-cover"
                        draggable={false}
                      />
                    ) : (
                      <div
                        className={cn(
                          'h-20 w-full',
                          mode === 'satellite' ? 'bg-foreground/80' : 'bg-muted',
                        )}
                      />
                    )}
                    <span className="block px-2 py-1.5 text-xs font-medium">
                      {t(`maps.basemap.${mode}`)}
                    </span>
                  </button>
                )
              })}
            </div>
          </div>
        </FilterPopoverButton>

        {onLocate ? (
          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                type="button"
                size="icon-lg"
                variant="ghost"
                aria-label={t('maps.locate.action')}
                disabled={locating}
                className="rounded-none"
                onClick={onLocate}
              >
                {locating ? <Spinner className="size-4" /> : <LocateFixedIcon />}
              </Button>
            </TooltipTrigger>
            <TooltipContent side="right" sideOffset={8} className="z-[1200]">
              {t('maps.locate.action')}
            </TooltipContent>
          </Tooltip>
        ) : null}

        {hasFilters ? (
          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                type="button"
                size="icon-lg"
                variant="ghost"
                aria-label={t('maps.filters.clear')}
                className="rounded-none"
                onClick={() =>
                  onChange({
                    q: '',
                    category: ALL,
                    countryCode: values.countryCode,
                    adminDivisionId: '',
                  })
                }
              >
                <FilterXIcon />
              </Button>
            </TooltipTrigger>
            <TooltipContent side="right" sideOffset={8} className="z-[1200]">
              {t('maps.filters.clear')}
            </TooltipContent>
          </Tooltip>
        ) : null}
      </div>
    </div>
  )
}

function FilterPopoverButton({
  label,
  active,
  icon,
  contentClassName,
  children,
}: {
  label: string
  active: boolean
  icon: ReactNode
  contentClassName?: string
  children: ReactNode
}) {
  return (
    <Popover>
      <Tooltip>
        <TooltipTrigger asChild>
          <PopoverTrigger asChild>
            <Button
              type="button"
              size="icon-lg"
              variant="ghost"
              aria-label={label}
              aria-pressed={active}
              className={cn(
                'rounded-none',
                active && 'bg-primary/10 text-primary',
              )}
            >
              {icon}
            </Button>
          </PopoverTrigger>
        </TooltipTrigger>
        <TooltipContent side="right" sideOffset={8} className="z-[1200]">
          {label}
        </TooltipContent>
      </Tooltip>
      <PopoverContent
        align="start"
        side="right"
        className={cn('z-[1200] w-72', contentClassName)}
        onOpenAutoFocus={(event) => event.preventDefault()}
        onInteractOutside={(event) => {
          const target = event.target as HTMLElement | null
          if (target?.closest('[data-slot="select-content"]')) {
            event.preventDefault()
          }
        }}
      >
        {children}
      </PopoverContent>
    </Popover>
  )
}
