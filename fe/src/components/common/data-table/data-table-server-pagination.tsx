'use client'

import {
  ChevronLeftIcon,
  ChevronRightIcon,
  ChevronsLeftIcon,
  ChevronsRightIcon,
} from 'lucide-react'
import { useTranslation } from 'react-i18next'

import { Button } from '@/components/ui/button'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'

const PAGE_SIZES = [10, 25, 50] as const

export type DataTableServerPaging = {
  page: number
  pageSize: number
  total: number
  onPageChange: (page: number) => void
  onPageSizeChange?: (pageSize: number) => void
}

type DataTableServerPaginationProps = {
  paging: DataTableServerPaging
}

export function DataTableServerPagination({ paging }: DataTableServerPaginationProps) {
  const { t } = useTranslation('common')
  const { page, pageSize, total, onPageChange, onPageSizeChange } = paging
  const pageCount = Math.max(1, Math.ceil(total / Math.max(pageSize, 1)))
  const from = total === 0 ? 0 : (page - 1) * pageSize + 1
  const to = Math.min(page * pageSize, total)
  const canPrev = page > 1
  const canNext = page < pageCount

  return (
    <div className="flex flex-col items-center justify-between gap-3 border-t px-3 py-3 sm:flex-row">
      <p className="order-2 text-sm text-muted-foreground sm:order-1">
        {total === 0
          ? t('table.noResults')
          : t('table.range', { from, to, total })}
      </p>
      <div className="order-1 flex items-center gap-4 sm:order-2">
        {onPageSizeChange ? (
          <div className="flex items-center gap-2">
            <span className="text-sm text-muted-foreground">{t('table.rows')}</span>
            <Select
              value={String(pageSize)}
              onValueChange={(value) => onPageSizeChange(Number(value))}
            >
              <SelectTrigger className="h-8 w-16 text-xs">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {PAGE_SIZES.map((size) => (
                  <SelectItem key={size} value={String(size)}>
                    {size}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        ) : null}
        <div className="flex items-center gap-1">
          <Button
            variant="ghost"
            size="icon"
            className="hidden size-8 sm:flex"
            onClick={() => onPageChange(1)}
            disabled={!canPrev}
            aria-label={t('table.firstPage')}
          >
            <ChevronsLeftIcon className="size-4" />
          </Button>
          <Button
            variant="ghost"
            size="icon"
            className="size-8"
            onClick={() => onPageChange(page - 1)}
            disabled={!canPrev}
            aria-label={t('table.previousPage')}
          >
            <ChevronLeftIcon className="size-4" />
          </Button>
          <span className="min-w-[4.5rem] text-center text-sm">
            {t('table.pageOf', { page, pageCount })}
          </span>
          <Button
            variant="ghost"
            size="icon"
            className="size-8"
            onClick={() => onPageChange(page + 1)}
            disabled={!canNext}
            aria-label={t('table.nextPage')}
          >
            <ChevronRightIcon className="size-4" />
          </Button>
          <Button
            variant="ghost"
            size="icon"
            className="hidden size-8 sm:flex"
            onClick={() => onPageChange(pageCount)}
            disabled={!canNext}
            aria-label={t('table.lastPage')}
          >
            <ChevronsRightIcon className="size-4" />
          </Button>
        </div>
      </div>
    </div>
  )
}
