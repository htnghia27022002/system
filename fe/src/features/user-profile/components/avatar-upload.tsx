'use client'

import { CameraIcon, UploadIcon } from 'lucide-react'
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
  /** Secondary line under the name (e.g. email). */
  subtitle?: string
  avatarUrl?: string | null
  onUploaded: (avatarUrl: string) => void
  /** Upload strategy; defaults to the signed-in user's own avatar endpoint. */
  upload?: (file: File) => Promise<{ avatarUrl?: string | null }>
}

function uploadErrorMessage(error: unknown, fallback: string): string {
  if (axios.isAxiosError(error)) {
    const data = error.response?.data as { message?: string; error?: string } | undefined
    return data?.message || data?.error || fallback
  }
  if (error instanceof Error && error.message) return error.message
  return fallback
}

export function AvatarUpload({
  name,
  subtitle,
  avatarUrl,
  onUploaded,
  upload = profileApi.uploadAvatar,
}: AvatarUploadProps) {
  const { t } = useTranslation('admin')
  const getInitials = useInitials()
  const inputRef = useRef<HTMLInputElement>(null)
  const [pending, setPending] = useState(false)
  const [localError, setLocalError] = useState<string | null>(null)

  const src = resolveMediaUrl(avatarUrl ?? undefined)
  const openPicker = () => inputRef.current?.click()

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
      const user = await upload(file)
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
    <div className="flex flex-col gap-4 sm:flex-row sm:items-center">
      <button
        type="button"
        className="group relative size-20 shrink-0 rounded-full outline-none focus-visible:ring-3 focus-visible:ring-ring/50 disabled:cursor-not-allowed"
        disabled={pending}
        onClick={openPicker}
        aria-label={t('profile.actions.changeAvatar')}
      >
        <Avatar className="size-20 overflow-hidden rounded-full border border-border">
          <AvatarImage src={src} alt="" className="object-cover" />
          <AvatarFallback className="rounded-full bg-primary/10 text-xl font-medium text-primary">
            {getInitials(name)}
          </AvatarFallback>
        </Avatar>
        <span className="absolute inset-0 flex items-center justify-center rounded-full bg-foreground/50 text-background opacity-0 transition-opacity group-hover:opacity-100 group-focus-visible:opacity-100 motion-reduce:transition-none">
          <CameraIcon className="size-5" aria-hidden />
        </span>
        {pending ? (
          <span className="absolute inset-0 flex items-center justify-center rounded-full bg-foreground/50 text-background">
            <Spinner className="size-5" />
          </span>
        ) : null}
      </button>

      <div className="min-w-0 flex-1 space-y-2">
        <div className="min-w-0">
          <p className="truncate font-medium">{name}</p>
          {subtitle ? (
            <p className="truncate text-sm text-muted-foreground">{subtitle}</p>
          ) : null}
        </div>
        <div className="flex flex-wrap items-center gap-x-3 gap-y-2">
          <Button
            type="button"
            variant="outline"
            size="sm"
            disabled={pending}
            onClick={openPicker}
          >
            {pending ? <Spinner className="size-4" /> : <UploadIcon />}
            {t('profile.actions.changeAvatar')}
          </Button>
          <p className="text-xs text-muted-foreground">{t('profile.avatar.hint')}</p>
        </div>
        <input
          ref={inputRef}
          type="file"
          accept="image/jpeg,image/png,image/webp"
          className="sr-only"
          tabIndex={-1}
          aria-hidden
          onChange={(event) => void onPick(event.target.files?.[0])}
        />
        <InputError message={localError ?? undefined} />
      </div>
    </div>
  )
}
