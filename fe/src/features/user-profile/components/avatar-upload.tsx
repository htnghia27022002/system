'use client'

import { useRef, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import axios from 'axios'

import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { InputError } from '@/components/common/input-error'
import { Spinner } from '@/components/ui/spinner'
import { useInitials } from '@/hooks/use-initials'
import { resolveMediaUrl } from '@/utils/resolve-media-url'

import { validateAvatarFile } from '../schemas/profile-schemas'
import { profileApi } from '../services/profile-api'

type AvatarUploadProps = {
  name: string
  avatarUrl?: string | null
  onUploaded: (avatarUrl: string) => void
}

function uploadErrorMessage(error: unknown, fallback: string): string {
  if (axios.isAxiosError(error)) {
    const data = error.response?.data as { message?: string; error?: string } | undefined
    return data?.message || data?.error || fallback
  }
  if (error instanceof Error && error.message) return error.message
  return fallback
}

export function AvatarUpload({ name, avatarUrl, onUploaded }: AvatarUploadProps) {
  const { t } = useTranslation('admin')
  const getInitials = useInitials()
  const inputRef = useRef<HTMLInputElement>(null)
  const [pending, setPending] = useState(false)
  const [localError, setLocalError] = useState<string | null>(null)

  const src = resolveMediaUrl(avatarUrl ?? undefined)

  const onPick = async (file: File | undefined) => {
    if (!file) return
    const clientError = validateAvatarFile(file)
    if (clientError) {
      setLocalError(clientError)
      return
    }
    setLocalError(null)
    setPending(true)
    try {
      const user = await profileApi.uploadAvatar(file)
      onUploaded(user.avatarUrl ?? '')
      toast.success(t('profile.toasts.avatarUpdated'))
    } catch (error) {
      const message = uploadErrorMessage(error, t('profile.errors.avatarUpload'))
      setLocalError(message)
      toast.error(message)
    } finally {
      setPending(false)
      if (inputRef.current) inputRef.current.value = ''
    }
  }

  return (
    <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
      <Avatar className="h-20 w-20 overflow-hidden rounded-full">
        <AvatarImage src={src} alt={name} />
        <AvatarFallback className="rounded-lg bg-neutral-200 text-lg text-black dark:bg-neutral-700 dark:text-white">
          {getInitials(name)}
        </AvatarFallback>
      </Avatar>
      <div className="space-y-2">
        <p className="text-sm text-muted-foreground">{t('profile.avatar.hint')}</p>
        <div className="flex flex-wrap gap-2">
          <Button
            type="button"
            variant="outline"
            size="sm"
            disabled={pending}
            onClick={() => inputRef.current?.click()}
          >
            {pending ? <Spinner className="size-4" /> : null}
            {t('profile.actions.changeAvatar')}
          </Button>
          <input
            ref={inputRef}
            type="file"
            accept="image/jpeg,image/png,image/webp"
            className="sr-only"
            aria-label={t('profile.actions.changeAvatar')}
            onChange={(event) => void onPick(event.target.files?.[0])}
          />
        </div>
        <InputError message={localError ?? undefined} />
      </div>
    </div>
  )
}
