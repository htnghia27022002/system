'use client'

import { FilterIcon, Trash2Icon, XIcon } from 'lucide-react'
import { useTranslation } from 'react-i18next'

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
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
import { PermissionGate, PermissionKeys } from '@/features/access-control'

import {
  WEBHOOK_METHOD_FILTERS,
  WEBHOOK_READ_FILTERS,
  type WebhookMethodFilter,
  type WebhookReadFilter,
} from '../types/webhooks'

type WebhooksListToolbarProps = {
  search: string
  onSearchChange: (value: string) => void
  method: WebhookMethodFilter
  onMethodChange: (value: WebhookMethodFilter) => void
  readFilter: WebhookReadFilter
  onReadFilterChange: (value: WebhookReadFilter) => void
  shown: number
  total: number
  page: number
  activeCount: number
  onClearAll: () => void
  clearAllPending?: boolean
}

export function WebhooksListToolbar({
  search,
  onSearchChange,
  method,
  onMethodChange,
  readFilter,
  onReadFilterChange,
  shown,
  total,
  page,
  activeCount,
  onClearAll,
  clearAllPending,
}: WebhooksListToolbarProps) {
  const { t } = useTranslation('common')

  const activeFilterCount =
    (method !== 'ALL' ? 1 : 0) +
    (readFilter !== 'all' ? 1 : 0) +
    (search.trim() ? 1 : 0)

  function clearFilters() {
    onSearchChange('')
    onMethodChange('ALL')
    onReadFilterChange('all')
  }

  return (
    <div className="flex items-center gap-1.5 border-b border-border px-2 py-1.5">
      <Popover>
        <PopoverTrigger asChild>
          <Button
            type="button"
            variant="outline"
            size="sm"
            className="h-8 gap-1.5 px-2"
            aria-label={t('tools.webhooks.filters.open')}
          >
            <FilterIcon className="size-3.5" />
            <span className="text-xs">{t('tools.webhooks.filters.open')}</span>
            {activeFilterCount > 0 ? (
              <Badge
                variant="secondary"
                className="h-4 min-w-4 px-1 text-[10px] tabular-nums"
              >
                {activeFilterCount}
              </Badge>
            ) : null}
          </Button>
        </PopoverTrigger>
        <PopoverContent align="start" className="w-72 space-y-3">
          <div className="flex items-center justify-between gap-2">
            <p className="text-sm font-medium">
              {t('tools.webhooks.filters.title')}
            </p>
            {activeFilterCount > 0 ? (
              <Button
                type="button"
                variant="ghost"
                size="sm"
                className="h-7 px-2 text-xs"
                onClick={clearFilters}
              >
                <XIcon className="size-3.5" />
                {t('tools.webhooks.filters.clear')}
              </Button>
            ) : null}
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="webhooks-search-popover">
              {t('tools.webhooks.filters.search')}
            </Label>
            <Input
              id="webhooks-search-popover"
              value={search}
              onChange={(e) => onSearchChange(e.target.value)}
              placeholder={t('tools.webhooks.filters.searchPlaceholder')}
              className="h-8"
            />
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="webhooks-method-popover">
              {t('tools.webhooks.filters.method')}
            </Label>
            <Select
              value={method}
              onValueChange={(v) => onMethodChange(v as WebhookMethodFilter)}
            >
              <SelectTrigger id="webhooks-method-popover" className="h-8 w-full">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {WEBHOOK_METHOD_FILTERS.map((m) => (
                  <SelectItem key={m} value={m}>
                    {m === 'ALL'
                      ? t('tools.webhooks.filters.allMethods')
                      : m}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="webhooks-read-popover">
              {t('tools.webhooks.filters.readStatus')}
            </Label>
            <Select
              value={readFilter}
              onValueChange={(v) => onReadFilterChange(v as WebhookReadFilter)}
            >
              <SelectTrigger id="webhooks-read-popover" className="h-8 w-full">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {WEBHOOK_READ_FILTERS.map((r) => (
                  <SelectItem key={r} value={r}>
                    {t(`tools.webhooks.filters.read.${r}`)}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </PopoverContent>
      </Popover>

      <p className="min-w-0 flex-1 truncate text-[11px] text-muted-foreground">
        {t('tools.webhooks.list.showing', { shown, total, page })}
      </p>

      <PermissionGate permission={PermissionKeys.webhooks.modify}>
        <AlertDialog>
          <Tooltip>
            <TooltipTrigger asChild>
              <AlertDialogTrigger asChild>
                <Button
                  type="button"
                  variant="ghost"
                  size="icon-sm"
                  className="text-destructive"
                  disabled={activeCount === 0 || clearAllPending}
                  aria-label={t('tools.webhooks.actions.clearAll')}
                >
                  <Trash2Icon className="size-3.5" />
                </Button>
              </AlertDialogTrigger>
            </TooltipTrigger>
            <TooltipContent side="bottom">
              {t('tools.webhooks.actions.clearAll')}
            </TooltipContent>
          </Tooltip>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>
                {t('tools.webhooks.clearAll.title')}
              </AlertDialogTitle>
              <AlertDialogDescription>
                {t('tools.webhooks.clearAll.description', {
                  count: activeCount,
                })}
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>
                {t('tools.webhooks.actions.cancel')}
              </AlertDialogCancel>
              <AlertDialogAction
                onClick={onClearAll}
                disabled={clearAllPending}
              >
                {t('tools.webhooks.actions.confirmClearAll')}
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      </PermissionGate>
    </div>
  )
}
