'use client'

import { useQueryClient } from '@tanstack/react-query'
import { KeyRoundIcon, MonitorSmartphoneIcon, UserRoundIcon } from 'lucide-react'
import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'

import { PageContainer, PageHeader } from '@/components/common/page-layout'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useAuthStore } from '@/store/auth-store'
import type { AuthUser } from '@/types/auth'

import { AccountSessionsCard } from './account-sessions-card'
import { AvatarUpload } from './avatar-upload'
import { ChangePasswordForm } from './change-password-form'
import { EmailVerificationNotice } from './email-verification-notice'
import { ProfileForm } from './profile-form'
import { SettingsSection } from './settings-section'
import { profileQueryKey, useProfile } from '../hooks/use-profile'

const TABS = ['profile', 'password', 'sessions'] as const
type TabValue = (typeof TABS)[number]

const TAB_ICONS = {
  profile: UserRoundIcon,
  password: KeyRoundIcon,
  sessions: MonitorSmartphoneIcon,
} satisfies Record<TabValue, unknown>

function isTab(value: string): value is TabValue {
  return (TABS as readonly string[]).includes(value)
}

/** Keeps the active tab in the URL hash so reloads and shared links land on it. */
function useHashTab(): [TabValue, (value: string) => void] {
  const [tab, setTab] = useState<TabValue>('profile')

  useEffect(() => {
    const sync = () => {
      const hash = window.location.hash.slice(1)
      if (isTab(hash)) setTab(hash)
    }
    sync()
    window.addEventListener('hashchange', sync)
    return () => window.removeEventListener('hashchange', sync)
  }, [])

  const change = (value: string) => {
    if (!isTab(value)) return
    setTab(value)
    const url = value === 'profile' ? window.location.pathname : `#${value}`
    window.history.replaceState(null, '', url)
  }

  return [tab, change]
}

export function ProfilePage() {
  const { t } = useTranslation('admin')
  const syncSession = useAuthStore((s) => s.syncSession)
  const queryClient = useQueryClient()
  const profileQuery = useProfile()
  const profile = profileQuery.data ?? null
  const [tab, setTab] = useHashTab()

  const applyUser = async (user: AuthUser) => {
    queryClient.setQueryData(profileQueryKey, user)
    await syncSession()
  }

  if (profileQuery.isLoading && !profile) {
    return (
      <PageContainer>
        <div className="space-y-2">
          <Skeleton className="h-8 w-56" />
          <Skeleton className="h-4 w-full max-w-md" />
        </div>
        <Skeleton className="h-9 w-72" />
        <Skeleton className="h-40 w-full" />
        <Skeleton className="h-72 w-full" />
      </PageContainer>
    )
  }

  if (!profile) {
    return (
      <PageContainer>
        <PageHeader title={t('profile.title')} />
        <Card>
          <CardContent className="flex flex-col items-start gap-3">
            <p className="text-sm text-destructive" role="alert">
              {t('profile.errors.loadFailed')}
            </p>
            <Button variant="outline" size="sm" onClick={() => void profileQuery.refetch()}>
              {t('access.actions.retry', { defaultValue: 'Retry' })}
            </Button>
          </CardContent>
        </Card>
      </PageContainer>
    )
  }

  return (
    <PageContainer className="gap-6">
      <PageHeader title={t('profile.title')} description={t('profile.description')} />

      {profile.emailVerified === false ? (
        <EmailVerificationNotice email={profile.email} />
      ) : null}

      <Tabs value={tab} onValueChange={setTab} className="min-w-0 gap-6">
        <div className="-mx-4 overflow-x-auto border-b border-border px-4">
          <TabsList variant="line" className="h-10 w-max justify-start gap-4 p-0">
            {TABS.map((value) => {
              const Icon = TAB_ICONS[value]
              return (
                <TabsTrigger key={value} value={value} className="flex-none px-1 pb-2">
                  <Icon aria-hidden />
                  {t(`profile.tabs.${value}`)}
                </TabsTrigger>
              )
            })}
          </TabsList>
        </div>

        <TabsContent value="profile" className="flex min-w-0 max-w-5xl flex-col gap-8">
          <SettingsSection
            title={t('profile.sections.avatar')}
            description={t('profile.sections.avatarDescription')}
          >
            <Card>
              <CardContent>
                <AvatarUpload
                  name={profile.name || profile.email}
                  subtitle={profile.name ? profile.email : undefined}
                  avatarUrl={profile.avatarUrl}
                  onUploaded={async (avatarUrl) => {
                    await applyUser({ ...profile, avatarUrl })
                  }}
                />
              </CardContent>
            </Card>
          </SettingsSection>
          <ProfileForm user={profile} onSaved={(user) => void applyUser(user)} />
        </TabsContent>

        <TabsContent value="password" className="min-w-0 max-w-5xl">
          <ChangePasswordForm hasPassword={profile.hasPassword !== false} />
        </TabsContent>

        <TabsContent value="sessions" className="min-w-0 max-w-5xl">
          <AccountSessionsCard />
        </TabsContent>
      </Tabs>
    </PageContainer>
  )
}
