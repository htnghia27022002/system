'use client'

import { CircleAlertIcon } from 'lucide-react'
import { useTranslation } from 'react-i18next'

import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import { useInitials } from '@/hooks/use-initials'
import { resolveMediaUrl } from '@/utils/resolve-media-url'

import { UserOAuthProviders } from './user-oauth-providers'
import type { ManagedUser } from '../types'

type UserIdentityProps = {
  user: ManagedUser
  isCurrentUser?: boolean
  /** Makes the name a link-styled button (opens the edit dialog). */
  onOpen?: () => void
}

/** Avatar + name + email, with linked sign-in providers and verification state. */
export function UserIdentity({ user, isCurrentUser, onOpen }: UserIdentityProps) {
  const { t } = useTranslation('admin')
  const getInitials = useInitials()
  const name = user.name || user.email

  return (
    <div className="flex min-w-0 items-center gap-3">
      <Avatar className="size-9 shrink-0 overflow-hidden rounded-full">
        <AvatarImage src={resolveMediaUrl(user.avatarUrl)} alt="" className="object-cover" />
        <AvatarFallback className="rounded-full bg-primary/10 text-xs font-medium text-primary">
          {getInitials(name)}
        </AvatarFallback>
      </Avatar>
      <div className="min-w-0">
        <div className="flex min-w-0 items-center gap-2">
          {onOpen ? (
            <button
              type="button"
              className="truncate text-left font-medium text-primary underline-offset-4 outline-none hover:underline focus-visible:underline"
              onClick={onOpen}
            >
              {name}
            </button>
          ) : (
            <span className="truncate font-medium">{name}</span>
          )}
          {isCurrentUser ? (
            <Badge variant="outline" className="shrink-0">
              {t('access.users.you')}
            </Badge>
          ) : null}
        </div>
        <div className="flex min-w-0 items-center gap-1.5 text-sm font-normal text-muted-foreground">
          <span className="truncate">{user.email}</span>
          {user.emailVerified === false ? (
            <Tooltip>
              <TooltipTrigger asChild>
                <CircleAlertIcon
                  className="size-3.5 shrink-0 text-warning"
                  aria-label={t('access.users.unverified')}
                />
              </TooltipTrigger>
              <TooltipContent side="top">{t('access.users.unverified')}</TooltipContent>
            </Tooltip>
          ) : null}
          {user.oauthProviders?.length ? (
            <UserOAuthProviders providers={user.oauthProviders} compact />
          ) : null}
        </div>
      </div>
    </div>
  )
}
