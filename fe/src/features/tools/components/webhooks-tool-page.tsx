'use client'

import { useTranslation } from 'react-i18next'

import { PageContainer, PageHeader } from '@/components/common/page-layout'

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
    <PageContainer>
      <PageHeader title={title} description={subtitle} />

      <WebhooksInbox />
    </PageContainer>
  )
}
