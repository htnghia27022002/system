import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { useInitials } from '@/hooks/use-initials'
import type { AuthUser } from '@/types/auth'

type UserInfoProps = {
  user: AuthUser
  showEmail?: boolean
}

export function UserInfo({ user, showEmail = false }: UserInfoProps) {
  const getInitials = useInitials()

  return (
    <>
      <Avatar className="h-8 w-8 overflow-hidden rounded-full">
        <AvatarImage src="" alt={user.name} />
        <AvatarFallback className="rounded-lg bg-neutral-200 text-black dark:bg-neutral-700 dark:text-white">
          {getInitials(user.name || user.email)}
        </AvatarFallback>
      </Avatar>
      <div className="grid flex-1 text-left text-sm leading-tight">
        <span className="truncate font-medium">{user.name || user.email}</span>
        {showEmail ? (
          <span className="truncate text-xs text-muted-foreground">
            {user.email}
          </span>
        ) : null}
      </div>
    </>
  )
}
