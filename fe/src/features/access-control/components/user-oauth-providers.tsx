'use client'

import { useTranslation } from 'react-i18next'

import { GoogleIcon } from '@/components/icons/google-icon'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip'

type UserOAuthProvidersProps = {
  providers?: string[] | null
}

function providerLabel(provider: string): string {
  switch (provider.toLowerCase()) {
    case 'google':
      return 'Google'
    default:
      return provider
  }
}

export function UserOAuthProviders({ providers }: UserOAuthProvidersProps) {
  const { t } = useTranslation('admin')
  const list = (providers ?? []).filter(Boolean)

  if (list.length === 0) {
    return <span className="text-muted-foreground">—</span>
  }

  return (
    <div className="flex items-center gap-1.5">
      {list.map((provider) => {
        const key = provider.toLowerCase()
        const label = providerLabel(provider)
        if (key === 'google') {
          return (
            <Tooltip key={provider}>
              <TooltipTrigger asChild>
                <span
                  className="inline-flex size-6 items-center justify-center rounded-md border border-border bg-background"
                  aria-label={t('access.users.oauth.google', {
                    defaultValue: 'Google',
                  })}
                >
                  <GoogleIcon className="size-3.5" />
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
                className="inline-flex h-6 items-center rounded-md border border-border bg-muted px-1.5 text-[10px] font-medium uppercase tracking-wide text-muted-foreground"
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
