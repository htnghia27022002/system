'use client'

import { useTranslation } from 'react-i18next'

import { StatusBadge } from '@/components/ui/status-badge'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import { cn } from '@/lib/utils'

import { isIngestInProgress } from '../hooks/use-map-ingest'
import type { IngestRun } from '../types'

type IngestStatusProps = {
  run: IngestRun | null | undefined
  className?: string
}

function statusVariant(status: IngestRun['status']) {
  if (status === 'completed') return 'success' as const
  if (status === 'failed') return 'error' as const
  if (status === 'running' || status === 'queued') return 'info' as const
  return 'neutral' as const
}

export function IngestStatus({ run, className }: IngestStatusProps) {
  const { t } = useTranslation('admin')

  if (!run) return null

  const inProgress = isIngestInProgress(run.status)
  const title =
    run.status === 'queued'
      ? t('maps.ingest.queued')
      : run.status === 'running'
        ? t('maps.ingest.running')
        : run.status === 'failed'
          ? t('maps.ingest.failed')
          : t('maps.ingest.completed')

  const detail = [
    t('maps.ingest.counters', {
      success: run.success,
      total: run.sourcesTotal,
      itemsSuccess: run.itemsSuccess,
      itemsError: run.itemsError,
    }),
    run.errorMessage,
    ...(run.sources ?? [])
      .filter((source) => source.status === 'failed')
      .map((source) =>
        source.errorMessage
          ? `${t('maps.ingest.sourceFailed', { name: source.sourceName })}: ${source.errorMessage}`
          : t('maps.ingest.sourceFailed', { name: source.sourceName }),
      ),
  ]
    .filter(Boolean)
    .join(' ')

  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <div
          className={cn(
            'rounded-lg border border-border bg-card/95 p-1.5 shadow-md backdrop-blur-sm',
            className,
          )}
        >
          <StatusBadge
            variant={statusVariant(run.status)}
            label={title}
            pulse={inProgress}
          />
        </div>
      </TooltipTrigger>
      <TooltipContent side="bottom" sideOffset={8} className="z-[1200] max-w-xs">
        {inProgress ? `${t('maps.ingest.inProgress')} ` : null}
        {detail}
      </TooltipContent>
    </Tooltip>
  )
}
