'use client'

import { useTranslation } from 'react-i18next'

import { TextLink } from '@/components/common/text-link'

type InvalidLinkNoticeProps = {
  /** Where to request a fresh link, e.g. /forgot-password. */
  retryHref?: string
  retryLabel?: string
}

export function InvalidLinkNotice({ retryHref, retryLabel }: InvalidLinkNoticeProps) {
  const { t } = useTranslation('common')
  return (
    <div className="flex flex-col gap-4 text-center" role="alert">
      <p className="text-sm text-muted-foreground">{t('auth.links.invalid')}</p>
      <div className="flex flex-col gap-2 text-sm">
        {retryHref ? <TextLink href={retryHref}>{retryLabel}</TextLink> : null}
        <TextLink href="/login">{t('auth.links.backToSignIn')}</TextLink>
      </div>
    </div>
  )
}
