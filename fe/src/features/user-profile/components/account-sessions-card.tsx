'use client'

import {
  GlobeIcon,
  MonitorIcon,
  SmartphoneIcon,
  TabletIcon,
  type LucideIcon,
} from 'lucide-react'
import { useMemo, useState } from 'react'
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
import { Button } from '@/components/ui/button'
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { StatusBadge } from '@/components/ui/status-badge'

import { SettingsSection } from './settings-section'
import { useAccountSessions } from '../hooks/use-account-sessions'
import type { AccountSession } from '../types'
import { formatDate, formatDateTime, formatRelative } from '../utils/format-date'
import { parseUserAgent, type DeviceType } from '../utils/parse-user-agent'

const DEVICE_ICONS: Record<DeviceType, LucideIcon> = {
  desktop: MonitorIcon,
  mobile: SmartphoneIcon,
  tablet: TabletIcon,
}

function sortSessions(items: AccountSession[]): AccountSession[] {
  return [...items].sort((a, b) => {
    if (a.current !== b.current) return a.current ? -1 : 1
    return new Date(b.lastUsedAt).getTime() - new Date(a.lastUsedAt).getTime()
  })
}

function RevokeSessionButton({
  disabled,
  onConfirm,
}: {
  disabled: boolean
  onConfirm: () => void
}) {
  const { t } = useTranslation('admin')
  return (
    <AlertDialog>
      <AlertDialogTrigger asChild>
        <Button
          type="button"
          variant="outline"
          size="sm"
          className="text-destructive hover:text-destructive"
          disabled={disabled}
        >
          {t('profile.sessions.revoke')}
        </Button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>{t('profile.sessions.revokeTitle')}</AlertDialogTitle>
          <AlertDialogDescription>
            {t('profile.sessions.revokeDescription')}
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>{t('profile.sessions.cancel')}</AlertDialogCancel>
          <AlertDialogAction variant="destructive" onClick={onConfirm} disabled={disabled}>
            {t('profile.sessions.revokeConfirm')}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}

function SessionRow({
  session,
  revokeDisabled,
  onRevoke,
}: {
  session: AccountSession
  revokeDisabled: boolean
  onRevoke: () => void
}) {
  const { t, i18n } = useTranslation('admin')
  const locale = i18n.language
  const { browser, os, deviceType } = parseUserAgent(session.userAgent)
  const Icon = browser || os ? DEVICE_ICONS[deviceType] : GlobeIcon

  const deviceName =
    browser && os
      ? t('profile.sessions.deviceOn', { browser, os })
      : browser || os || t('profile.sessions.unknownDevice')

  const lastActive = session.current
    ? t('profile.sessions.activeNow')
    : (() => {
        const relative = formatRelative(session.lastUsedAt, locale)
        return relative
          ? t('profile.sessions.lastActive', { time: relative })
          : t('profile.sessions.activeNow')
      })()

  return (
    <li
      data-testid={`session-row-${session.id}`}
      data-current={session.current ? 'true' : 'false'}
      className="flex flex-col gap-3 px-4 py-4 sm:flex-row sm:items-center"
    >
      <div className="flex min-w-0 flex-1 items-start gap-3">
        <span className="flex size-10 shrink-0 items-center justify-center rounded-lg bg-muted text-muted-foreground">
          <Icon className="size-5" aria-hidden />
        </span>
        <div className="min-w-0 flex-1 space-y-1">
          <div className="flex flex-wrap items-center gap-2">
            <p className="truncate font-medium" title={session.userAgent ?? undefined}>
              {deviceName}
            </p>
            {session.current ? (
              <StatusBadge variant="success" label={t('profile.sessions.current')} />
            ) : null}
          </div>
          <p className="flex flex-wrap items-center gap-x-2 text-sm text-muted-foreground">
            <span>{session.ipAddress?.trim() || t('profile.sessions.unknown')}</span>
            <span aria-hidden>·</span>
            <time
              dateTime={session.lastUsedAt}
              title={formatDateTime(session.lastUsedAt, locale)}
              className={session.current ? 'text-success' : undefined}
            >
              {lastActive}
            </time>
          </p>
          <p className="text-xs text-muted-foreground">
            <time dateTime={session.createdAt} title={formatDateTime(session.createdAt, locale)}>
              {t('profile.sessions.signedIn', { date: formatDate(session.createdAt, locale) })}
            </time>
            <span aria-hidden> · </span>
            <time dateTime={session.expiresAt} title={formatDateTime(session.expiresAt, locale)}>
              {t('profile.sessions.expiresOn', { date: formatDate(session.expiresAt, locale) })}
            </time>
          </p>
        </div>
      </div>
      {session.current ? null : (
        <div className="pl-13 sm:pl-0">
          <RevokeSessionButton disabled={revokeDisabled} onConfirm={onRevoke} />
        </div>
      )}
    </li>
  )
}

function SessionListSkeleton() {
  return (
    <ul className="divide-y divide-border" aria-hidden>
      {[0, 1, 2].map((i) => (
        <li key={i} className="flex items-center gap-3 px-4 py-4">
          <Skeleton className="size-10 rounded-lg" />
          <div className="flex-1 space-y-2">
            <Skeleton className="h-4 w-40" />
            <Skeleton className="h-3 w-56" />
          </div>
        </li>
      ))}
    </ul>
  )
}

export function AccountSessionsCard() {
  const { t } = useTranslation('admin')
  const { query, items, hasOtherSessions, revokeSession, revokeOthers } =
    useAccountSessions()
  const [revokeOthersOpen, setRevokeOthersOpen] = useState(false)
  const sessions = useMemo(() => sortSessions(items), [items])

  const isInitialLoading = query.isLoading && items.length === 0
  const isInitialError = query.isError && items.length === 0

  return (
    <SettingsSection
      title={t('profile.sections.sessions')}
      description={t('profile.sessions.description')}
    >
      <Card className="gap-0 py-0">
        <CardHeader className="border-b py-4">
          <CardTitle>
            {isInitialLoading ? (
              <Skeleton className="h-5 w-32" />
            ) : (
              t('profile.sessions.activeCount', { count: items.length })
            )}
          </CardTitle>
          <CardDescription>{t('profile.sessions.listHint')}</CardDescription>
          {hasOtherSessions ? (
            <CardAction>
              <AlertDialog open={revokeOthersOpen} onOpenChange={setRevokeOthersOpen}>
                <AlertDialogTrigger asChild>
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    className="text-destructive hover:text-destructive"
                    disabled={revokeOthers.isPending}
                  >
                    {t('profile.sessions.revokeOthers')}
                  </Button>
                </AlertDialogTrigger>
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>{t('profile.sessions.revokeOthersTitle')}</AlertDialogTitle>
                    <AlertDialogDescription>
                      {t('profile.sessions.revokeOthersDescription')}
                    </AlertDialogDescription>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>{t('profile.sessions.cancel')}</AlertDialogCancel>
                    <AlertDialogAction
                      variant="destructive"
                      disabled={revokeOthers.isPending}
                      onClick={() => {
                        revokeOthers.mutate(undefined, {
                          onSuccess: () => setRevokeOthersOpen(false),
                        })
                      }}
                    >
                      {t('profile.sessions.revokeOthersConfirm')}
                    </AlertDialogAction>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>
            </CardAction>
          ) : null}
        </CardHeader>

        <CardContent className="px-0" aria-busy={query.isFetching}>
          {isInitialLoading ? (
            <SessionListSkeleton />
          ) : isInitialError ? (
            <div className="flex flex-col items-start gap-3 px-4 py-6">
              <p className="text-sm text-destructive" role="alert">
                {t('profile.sessions.loadFailed')}
              </p>
              <Button type="button" variant="outline" size="sm" onClick={() => void query.refetch()}>
                {t('profile.sessions.retry')}
              </Button>
            </div>
          ) : sessions.length === 0 ? (
            <p className="px-4 py-6 text-sm text-muted-foreground">{t('profile.sessions.empty')}</p>
          ) : (
            <ul className="divide-y divide-border">
              {sessions.map((session) => (
                <SessionRow
                  key={session.id}
                  session={session}
                  revokeDisabled={revokeSession.isPending}
                  onRevoke={() => revokeSession.mutate(session.id)}
                />
              ))}
            </ul>
          )}
        </CardContent>
      </Card>
    </SettingsSection>
  )
}
