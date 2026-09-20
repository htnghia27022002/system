'use client'

import type { ReactNode } from 'react'
import { DatabaseIcon, DownloadIcon } from 'lucide-react'
import { useTranslation } from 'react-i18next'

import { PermissionGate, PermissionKeys } from '@/features/access-control'
import { Button } from '@/components/ui/button'
import { Spinner } from '@/components/ui/spinner'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import { cn } from '@/lib/utils'

type MapActionsProps = {
  ingesting: boolean
  onLoadData: () => void
  onOpenSources: () => void
}

export function MapActions({
  ingesting,
  onLoadData,
  onOpenSources,
}: MapActionsProps) {
  const { t } = useTranslation('admin')

  return (
    <div className="pointer-events-none absolute inset-x-3 top-3 z-[1100] flex justify-end">
      <div className="pointer-events-auto flex flex-col overflow-hidden rounded-lg border border-border bg-card/95 shadow-md backdrop-blur-sm">
        <PermissionGate permission={PermissionKeys.maps.modify}>
          <IconTipButton
            label={t('maps.ingest.loadData')}
            onClick={onLoadData}
            disabled={ingesting}
          >
            {ingesting ? <Spinner className="size-4" /> : <DownloadIcon />}
          </IconTipButton>
        </PermissionGate>
        <IconTipButton label={t('maps.tabs.sources')} onClick={onOpenSources}>
          <DatabaseIcon />
        </IconTipButton>
      </div>
    </div>
  )
}

function IconTipButton({
  label,
  onClick,
  disabled,
  children,
}: {
  label: string
  onClick: () => void
  disabled?: boolean
  children: ReactNode
}) {
  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <Button
          type="button"
          size="icon-lg"
          variant="ghost"
          aria-label={label}
          disabled={disabled}
          className={cn('rounded-none')}
          onClick={onClick}
        >
          {children}
        </Button>
      </TooltipTrigger>
      <TooltipContent side="left" sideOffset={8} className="z-[1200]">
        {label}
      </TooltipContent>
    </Tooltip>
  )
}
