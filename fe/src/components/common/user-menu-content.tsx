'use client'

import {
  CheckIcon,
  LanguagesIcon,
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
import { UserInfo } from '@/components/common/user-info'
import { useAuthStore } from '@/store/auth-store'
import type { AuthUser } from '@/types/auth'

type UserMenuContentProps = {
  user: AuthUser
  onNavigate?: () => void
}

export function UserMenuContent({ user, onNavigate }: UserMenuContentProps) {
  const { t, i18n } = useTranslation('admin')
  const signOut = useAuthStore((state) => state.signOut)
  const { theme, setTheme } = useTheme()

  const themes = [
    { value: 'light', icon: SunIcon, label: t('userMenu.themeLight') },
    { value: 'dark', icon: MoonIcon, label: t('userMenu.themeDark') },
    { value: 'system', icon: MonitorIcon, label: t('userMenu.themeSystem') },
  ] as const

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

      <DropdownMenuItem
        onClick={() => void i18n.changeLanguage(i18n.language === 'en' ? 'vi' : 'en')}
      >
        <LanguagesIcon className="size-4" />
        {i18n.language === 'en'
          ? t('userMenu.switchToVi')
          : t('userMenu.switchToEn')}
      </DropdownMenuItem>

      <DropdownMenuSeparator />

      <DropdownMenuItem disabled>
        <SettingsIcon className="size-4" />
        {t('userMenu.settings')}
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
