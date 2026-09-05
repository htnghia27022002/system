'use client'

import { useQueryClient } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'

import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useAuthStore } from '@/store/auth-store'
import type { AuthUser } from '@/types/auth'

import { AccountSessionsCard } from './account-sessions-card'
import { AvatarUpload } from './avatar-upload'
import { ChangePasswordForm } from './change-password-form'
import { ProfileForm } from './profile-form'
import { profileQueryKey, useProfile } from '../hooks/use-profile'

export function ProfilePage() {
  const { t } = useTranslation('admin')
  const syncSession = useAuthStore((s) => s.syncSession)
  const queryClient = useQueryClient()
  const profileQuery = useProfile()
  const profile = profileQuery.data ?? null

  const applyUser = async (user: AuthUser) => {
    queryClient.setQueryData(profileQueryKey, user)
    await syncSession()
  }

  if (profileQuery.isLoading && !profile) {
    return (
      <div className="flex flex-1 flex-col gap-4 p-4">
        <Skeleton className="h-8 w-48" />
        <Skeleton className="h-10 w-full max-w-md" />
        <Skeleton className="h-64 w-full" />
      </div>
    )
  }

  if (profileQuery.isError && !profile) {
    return (
      <div className="flex flex-col gap-3 p-4">
        <p className="text-sm text-destructive" role="alert">
          {t('profile.errors.loadFailed')}
        </p>
        <Button variant="outline" className="w-fit" onClick={() => void profileQuery.refetch()}>
          {t('access.actions.retry', { defaultValue: 'Retry' })}
        </Button>
      </div>
    )
  }

  if (!profile) {
    return (
      <div className="p-4">
        <p className="text-sm text-destructive">{t('profile.errors.loadFailed')}</p>
      </div>
    )
  }

  return (
    <div className="flex flex-1 flex-col gap-4 p-4">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight">
          {t('profile.title')}
        </h1>
        <p className="text-sm text-muted-foreground">{t('profile.description')}</p>
      </div>

      <Tabs defaultValue="profile" className="min-w-0 gap-4">
        <TabsList
          variant="line"
          className="h-auto w-full max-w-full flex-wrap justify-start"
        >
          <TabsTrigger value="profile" className="flex-none px-3">
            {t('profile.tabs.profile')}
          </TabsTrigger>
          <TabsTrigger value="password" className="flex-none px-3">
            {t('profile.tabs.password')}
          </TabsTrigger>
          <TabsTrigger value="sessions" className="flex-none px-3">
            {t('profile.tabs.sessions')}
          </TabsTrigger>
        </TabsList>

        <TabsContent value="profile" className="min-w-0">
          <div className="flex max-w-3xl flex-col gap-8">
            <section className="space-y-3">
              <h2 className="text-sm font-medium">{t('profile.sections.avatar')}</h2>
              <AvatarUpload
                name={profile.name || profile.email}
                avatarUrl={profile.avatarUrl}
                onUploaded={async (avatarUrl) => {
                  await applyUser({ ...profile, avatarUrl })
                }}
              />
            </section>
            <section className="space-y-3">
              <h2 className="text-sm font-medium">
                {t('profile.sections.personalInfo')}
              </h2>
              <ProfileForm user={profile} onSaved={(user) => void applyUser(user)} />
            </section>
          </div>
        </TabsContent>

        <TabsContent value="password" className="min-w-0">
          <div className="max-w-md">
            <ChangePasswordForm hasPassword={profile.hasPassword !== false} />
          </div>
        </TabsContent>

        <TabsContent value="sessions" className="min-w-0">
          <AccountSessionsCard />
        </TabsContent>
      </Tabs>
    </div>
  )
}
