'use client'

import { useTranslation } from 'react-i18next'

import { getToolById } from '../catalog'
import { toolsFeatureContent } from '../content'
import { WebhooksInbox } from './webhooks-inbox'

/**
 * Owner Webhooks tool — admin inbox UI for the account's public URL.
 * Mounted at `/admin/tools/webhooks` behind `PermissionGuard(webhooks:view)`.
 * Public capture at `/tools/webhooks/{uuid}` is BE-only (nginx / Next rewrite).
 */
export function WebhooksToolPage() {
  const { t } = useTranslation('common')
  const tool = getToolById('webhooks')
  const copy = toolsFeatureContent.webhooks
  const title = tool?.name ?? t('tools.webhooks.title', { defaultValue: copy.title })
  const subtitle =
    tool?.description ??
    t('tools.webhooks.subtitle', { defaultValue: copy.subtitle })

  return (
    <div className="flex flex-col gap-4 p-4">
      <header className="max-w-3xl">
        <h1 className="font-heading text-xl font-semibold tracking-tight sm:text-2xl">
          {title}
        </h1>
        <p className="mt-0.5 text-sm text-muted-foreground">{subtitle}</p>
      </header>

      <WebhooksInbox />
    </div>
  )
}
