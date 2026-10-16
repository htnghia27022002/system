'use client'

import Link from 'next/link'
import {
  LanguagesIcon,
  LogOutIcon,
  MonitorIcon,
  MoonIcon,
  SunIcon,
  UserRoundCogIcon,
} from 'lucide-react'
import { useTheme } from 'next-themes'
import { useTranslation } from 'react-i18next'

import {
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuPortal,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuSeparator,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
} from '@/components/ui/dropdown-menu'
import {
  LOCALE_OPTIONS,
  normalizeLocale,
  type LocaleCode,
} from '@/components/common/locale-flags'
import { UserInfo } from '@/components/common/user-info'
import { useSignOut } from '@/features/auth'
import type { AuthUser } from '@/types/auth'

type UserMenuContentProps = {
  user: AuthUser
  onNavigate?: () => void
}

export function UserMenuContent({ user, onNavigate }: UserMenuContentProps) {
  const { t, i18n } = useTranslation('admin')
  const signOut = useSignOut()
  const { theme = 'system', setTheme } = useTheme()
  const activeLocale = normalizeLocale(i18n.language)

  const themes = [
    { value: 'light', icon: SunIcon, label: t('userMenu.themeLight') },
    { value: 'dark', icon: MoonIcon, label: t('userMenu.themeDark') },
    { value: 'system', icon: MonitorIcon, label: t('userMenu.themeSystem') },
  ] as const
  const activeTheme = themes.find((item) => item.value === theme) ?? themes[2]
  const ThemeIcon = activeTheme.icon
  const activeLocaleOption =
    LOCALE_OPTIONS.find((option) => option.code === activeLocale) ?? LOCALE_OPTIONS[0]

  return (
    <>
      <DropdownMenuLabel className="p-0 font-normal">
        <div className="flex items-center gap-2 px-1.5 py-1.5 text-left text-sm">
          <UserInfo user={user} showEmail />
        </div>
      </DropdownMenuLabel>
      <DropdownMenuSeparator />

      <DropdownMenuGroup>
        <DropdownMenuItem asChild>
          <Link href="/admin/profile" onClick={() => onNavigate?.()}>
            <UserRoundCogIcon />
            {t('userMenu.settings')}
          </Link>
        </DropdownMenuItem>

        <DropdownMenuSub>
          <DropdownMenuSubTrigger>
            <ThemeIcon />
            {t('userMenu.theme')}
            <span className="ml-auto text-xs text-muted-foreground">{activeTheme.label}</span>
          </DropdownMenuSubTrigger>
          <DropdownMenuPortal>
            <DropdownMenuSubContent className="min-w-36">
              <DropdownMenuRadioGroup value={theme} onValueChange={setTheme}>
                {themes.map(({ value, icon: Icon, label }) => (
                  <DropdownMenuRadioItem key={value} value={value}>
                    <Icon />
                    {label}
                  </DropdownMenuRadioItem>
                ))}
              </DropdownMenuRadioGroup>
            </DropdownMenuSubContent>
          </DropdownMenuPortal>
        </DropdownMenuSub>

        <DropdownMenuSub>
          <DropdownMenuSubTrigger>
            <LanguagesIcon />
            {t('userMenu.language')}
            <span className="ml-auto flex items-center gap-1.5 text-xs text-muted-foreground">
              <activeLocaleOption.Flag className="h-3 w-4" />
              {activeLocaleOption.code.toUpperCase()}
            </span>
          </DropdownMenuSubTrigger>
          <DropdownMenuPortal>
            <DropdownMenuSubContent className="min-w-40">
              <DropdownMenuRadioGroup
                value={activeLocale}
                onValueChange={(code) => {
                  if (code !== activeLocale) void i18n.changeLanguage(code as LocaleCode)
                }}
              >
                {LOCALE_OPTIONS.map(({ code, label, Flag }) => (
                  <DropdownMenuRadioItem key={code} value={code}>
                    <Flag />
                    {label}
                  </DropdownMenuRadioItem>
                ))}
              </DropdownMenuRadioGroup>
            </DropdownMenuSubContent>
          </DropdownMenuPortal>
        </DropdownMenuSub>
      </DropdownMenuGroup>

      <DropdownMenuSeparator />

      <DropdownMenuItem
        variant="destructive"
        onClick={() => {
          onNavigate?.()
          signOut.mutate()
        }}
      >
        <LogOutIcon />
        {t('userMenu.signOut')}
      </DropdownMenuItem>
    </>
  )
}
