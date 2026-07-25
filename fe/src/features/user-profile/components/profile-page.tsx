'use client'

import { useCallback, useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Skeleton } from '@/components/ui/skeleton'
import { useAuthStore } from '@/store/auth-store'
import type { AuthUser } from '@/types/auth'

import { AvatarUpload } from './avatar-upload'
import { ChangePasswordForm } from './change-password-form'
import { ProfileForm } from './profile-form'
import { profileApi } from '../services/profile-api'

export function ProfilePage() {
  const { t } = useTranslation('admin')
  const storeUser = useAuthStore((s) => s.user)
  const syncSession = useAuthStore((s) => s.syncSession)
  const [profile, setProfile] = useState<AuthUser | null>(null)
  const [loading, setLoading] = useState(true)
  const [loadError, setLoadError] = useState<string | null>(null)

  const load = useCallback(async () => {
    setLoading(true)
    setLoadError(null)
    try {
      const user = await profileApi.getProfile()
      setProfile(user)
    } catch {
      setLoadError(t('profile.errors.loadFailed'))
      if (storeUser) setProfile(storeUser)
    } finally {
      setLoading(false)
    }
  }, [storeUser, t])

  useEffect(() => {
    void load()
  }, [load])

  const applyUser = async (user: AuthUser) => {
    setProfile(user)
    await syncSession()
  }

  if (loading && !profile) {
    return (
      <div className="flex flex-col gap-4 p-4 md:p-6">
        <Skeleton className="h-8 w-48" />
        <Skeleton className="h-40 w-full" />
        <Skeleton className="h-64 w-full" />
      </div>
    )
  }

  if (!profile) {
    return (
      <div className="p-4 md:p-6">
        <p className="text-sm text-destructive">
          {loadError ?? t('profile.errors.loadFailed')}
        </p>
      </div>
    )
  }

  return (
    <div className="mx-auto flex w-full max-w-3xl flex-col gap-6 p-4 md:p-6">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight">
          {t('profile.title')}
        </h1>
        <p className="text-sm text-muted-foreground">{t('profile.description')}</p>
      </div>

      {loadError ? (
        <p className="text-sm text-destructive" role="alert">
          {loadError}
        </p>
      ) : null}

      <Card>
        <CardHeader>
          <CardTitle className="text-base">{t('profile.sections.avatar')}</CardTitle>
        </CardHeader>
        <CardContent>
          <AvatarUpload
            name={profile.name || profile.email}
            avatarUrl={profile.avatarUrl}
            onUploaded={async (avatarUrl) => {
              setProfile((prev) => (prev ? { ...prev, avatarUrl } : prev))
              await syncSession()
            }}
          />
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="text-base">
            {t('profile.sections.personalInfo')}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <ProfileForm user={profile} onSaved={(user) => void applyUser(user)} />
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle className="text-base">
            {t('profile.sections.password')}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <ChangePasswordForm hasPassword={profile.hasPassword !== false} />
        </CardContent>
      </Card>
    </div>
  )
}
