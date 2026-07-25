'use client'

import Link from 'next/link'
import {
  CheckIcon,
  LogOutIcon,
  MonitorIcon,
  MoonIcon,
  SettingsIcon,
  SunIcon,
} from 'lucide-react'
import { useTheme } from 'next-themes'
import { useTranslation } from 'react-i18next'

import {
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
} from '@/components/ui/dropdown-menu'
import {
  LOCALE_OPTIONS,
  normalizeLocale,
  type LocaleCode,
} from '@/components/common/locale-flags'
import { UserInfo } from '@/components/common/user-info'
import { useAuthStore } from '@/store/auth-store'
import { cn } from '@/lib/utils'
import type { AuthUser } from '@/types/auth'

type UserMenuContentProps = {
  user: AuthUser
  onNavigate?: () => void
}

export function UserMenuContent({ user, onNavigate }: UserMenuContentProps) {
  const { t, i18n } = useTranslation('admin')
  const signOut = useAuthStore((state) => state.signOut)
  const { theme, setTheme } = useTheme()
  const activeLocale = normalizeLocale(i18n.language)

  const themes = [
    { value: 'light', icon: SunIcon, label: t('userMenu.themeLight') },
    { value: 'dark', icon: MoonIcon, label: t('userMenu.themeDark') },
    { value: 'system', icon: MonitorIcon, label: t('userMenu.themeSystem') },
  ] as const

  const selectLocale = (code: LocaleCode) => {
    if (code !== activeLocale) {
      void i18n.changeLanguage(code)
    }
  }

  return (
    <>
      <DropdownMenuLabel className="p-0 font-normal">
        <div className="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
          <UserInfo user={user} showEmail />
        </div>
      </DropdownMenuLabel>
      <DropdownMenuSeparator />

      <DropdownMenuGroup>
        {themes.map(({ value, icon: Icon, label }) => (
          <DropdownMenuItem key={value} onClick={() => setTheme(value)}>
            <Icon className="size-4" />
            {label}
            {theme === value ? <CheckIcon className="ml-auto size-3.5" /> : null}
          </DropdownMenuItem>
        ))}
      </DropdownMenuGroup>

      <DropdownMenuSeparator />

      <DropdownMenuLabel className="px-2 py-1.5 text-xs font-normal text-muted-foreground">
        {t('userMenu.language')}
      </DropdownMenuLabel>
      <DropdownMenuGroup>
        {LOCALE_OPTIONS.map(({ code, label, Flag }) => {
          const isActive = activeLocale === code
          return (
            <DropdownMenuItem
              key={code}
              onSelect={(event) => {
                event.preventDefault()
                selectLocale(code)
              }}
              aria-current={isActive ? 'true' : undefined}
              className={cn(isActive && 'bg-accent')}
            >
              <Flag />
              <span>{label}</span>
              {isActive ? <CheckIcon className="ml-auto size-3.5" /> : null}
            </DropdownMenuItem>
          )
        })}
      </DropdownMenuGroup>

      <DropdownMenuSeparator />

      <DropdownMenuItem asChild>
        <Link
          href="/admin/profile"
          onClick={() => onNavigate?.()}
        >
          <SettingsIcon className="size-4" />
          {t('userMenu.settings')}
        </Link>
      </DropdownMenuItem>

      <DropdownMenuSeparator />

      <DropdownMenuItem
        className="text-destructive focus:text-destructive"
        onClick={() => {
          onNavigate?.()
          signOut()
        }}
      >
        <LogOutIcon className="size-4" />
        {t('userMenu.signOut')}
      </DropdownMenuItem>
    </>
  )
}
