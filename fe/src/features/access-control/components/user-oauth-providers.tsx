'use client'

import { useTranslation } from 'react-i18next'

import { GoogleIcon } from '@/components/icons/google-icon'
import { cn } from '@/lib/utils'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip'

type UserOAuthProvidersProps = {
  providers?: string[] | null
  /** Smaller, borderless marks for inline use next to text. */
  compact?: boolean
}

function providerLabel(provider: string): string {
  switch (provider.toLowerCase()) {
    case 'google':
      return 'Google'
    default:
      return provider
  }
}

export function UserOAuthProviders({ providers, compact = false }: UserOAuthProvidersProps) {
  const { t } = useTranslation('admin')
  const list = (providers ?? []).filter(Boolean)

  if (list.length === 0) {
    return <span className="text-muted-foreground">—</span>
  }

  return (
    <div className={cn('flex items-center', compact ? 'gap-1' : 'gap-1.5')}>
      {list.map((provider) => {
        const key = provider.toLowerCase()
        const label = providerLabel(provider)
        if (key === 'google') {
          return (
            <Tooltip key={provider}>
              <TooltipTrigger asChild>
                <span
                  className={cn(
                    'inline-flex items-center justify-center',
                    compact ? 'size-4' : 'size-6 rounded-md border border-border bg-background',
                  )}
                  aria-label={t('access.users.oauth.google', {
                    defaultValue: 'Google',
                  })}
                >
                  <GoogleIcon className={compact ? 'size-3' : 'size-3.5'} />
                </span>
              </TooltipTrigger>
              <TooltipContent side="top">{label}</TooltipContent>
            </Tooltip>
          )
        }
        return (
          <Tooltip key={provider}>
            <TooltipTrigger asChild>
              <span
                className="inline-flex h-6 items-center rounded-md border border-border bg-muted px-1.5 text-xs font-medium uppercase tracking-wide text-muted-foreground"
                aria-label={label}
              >
                {label.slice(0, 3)}
              </span>
            </TooltipTrigger>
            <TooltipContent side="top">{label}</TooltipContent>
          </Tooltip>
        )
      })}
    </div>
  )
}
