'use client'

import Image from 'next/image'
import Link from 'next/link'
import { MenuIcon } from 'lucide-react'
import { useState, type ReactNode } from 'react'
import { useTranslation } from 'react-i18next'

import { LocaleSelect } from '@/components/common/locale-select'
import { ThemeToggle } from '@/components/common/theme-toggle'
import { Button } from '@/components/ui/button'
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/components/ui/sheet'
import { useSignOut } from '@/features/auth/hooks/use-sign-out'
import { cn } from '@/lib/utils'
import { useAuthStore } from '@/store/auth-store'
import logoLockup from '@/assets/favicon/logo-lockup.svg'

type PublicSiteHeaderProps = {
  /** Visual variant for sticky marketing vs tools pages */
  variant?: 'landing' | 'tools'
  /** Optional secondary links (e.g. landing section jumps) rendered after site Navigate */
  children?: ReactNode
  className?: string
}

/**
 * Shared public chrome: brand, Navigate (Home/Tools), auth, LocaleSelect, ThemeToggle.
 * Mounted on `/` and `/tools` (owner Webhooks UI is under `/admin/tools/webhooks`).
 */
export function PublicSiteHeader({
  variant = 'tools',
  children,
  className,
}: PublicSiteHeaderProps) {
  const { t } = useTranslation('common')
  const accessToken = useAuthStore((state) => state.accessToken)
  const hasHydrated = useAuthStore((state) => state.hasHydrated)
  const signOutMutation = useSignOut()
  const [mobileOpen, setMobileOpen] = useState(false)

  const closeMobile = () => setMobileOpen(false)

  // Match SSR until auth rehydrates from localStorage (avoids Sign in ↔ Admin flash/mismatch).
  const showAuthed = hasHydrated && Boolean(accessToken)

  const authControls = showAuthed ? (
    <>
      <Button variant="ghost" size="sm" asChild>
        <Link href="/admin">{t('nav.admin')}</Link>
      </Button>
      <Button
        variant="outline"
        size="sm"
        disabled={signOutMutation.isPending}
        onClick={() => signOutMutation.mutate()}
      >
        {t('auth.actions.signOut')}
      </Button>
    </>
  ) : (
    <>
      <Button variant="ghost" size="sm" asChild className="hidden sm:inline-flex">
        <Link href="/login">{t('auth.actions.signIn')}</Link>
      </Button>
      <Button size="sm" asChild>
        <Link href="/register">{t('auth.actions.register')}</Link>
      </Button>
    </>
  )

  const mobileAuth = showAuthed ? (
    <>
      <Button variant="ghost" className="justify-start" asChild onClick={closeMobile}>
        <Link href="/admin">{t('nav.admin')}</Link>
      </Button>
      <Button
        variant="outline"
        className="justify-start"
        disabled={signOutMutation.isPending}
        onClick={() => {
          signOutMutation.mutate()
          closeMobile()
        }}
      >
        {t('auth.actions.signOut')}
      </Button>
    </>
  ) : (
    <>
      <Button variant="ghost" className="justify-start" asChild onClick={closeMobile}>
        <Link href="/login">{t('auth.actions.signIn')}</Link>
      </Button>
      <Button className="justify-start" asChild onClick={closeMobile}>
        <Link href="/register">{t('auth.actions.register')}</Link>
      </Button>
    </>
  )

  return (
    <header
      className={cn(
        'z-50 border-b border-border/60 bg-background/90 backdrop-blur-xl',
        variant === 'landing' ? 'fixed inset-x-0 top-0' : 'sticky top-0',
        className,
      )}
    >
      <div className="mx-auto flex h-16 max-w-6xl items-center justify-between gap-3 px-4 sm:px-6 lg:px-8">
        <div className="flex min-w-0 items-center gap-3">
          <Link href="/" className="flex shrink-0 items-center gap-2" aria-label={t('nav.home')}>
            <Image
              src={logoLockup}
              alt=""
              width={120}
              height={38}
              className="h-8 w-auto brightness-110"
              priority
            />
          </Link>

          <nav className="hidden items-center gap-1 md:flex" aria-label={t('nav.siteLabel')}>
            <Button variant="ghost" size="sm" asChild>
              <Link href="/">{t('nav.home')}</Link>
            </Button>
            <Button variant="ghost" size="sm" asChild>
              <Link href="/tools">{t('nav.tools')}</Link>
            </Button>
            {children}
          </nav>
        </div>

        <div className="flex items-center gap-2">
          <div className="hidden items-center gap-2 sm:flex">{authControls}</div>
          <LocaleSelect className="hidden sm:inline-flex" />
          <ThemeToggle />

          <Sheet open={mobileOpen} onOpenChange={setMobileOpen}>
            <SheetTrigger asChild>
              <Button
                variant="outline"
                size="icon-sm"
                className="md:hidden"
                aria-label={t('nav.openMenu')}
              >
                <MenuIcon className="size-4" />
              </Button>
            </SheetTrigger>
            <SheetContent side="right" className="w-[min(100%,20rem)]">
              <SheetHeader>
                <SheetTitle>{t('nav.siteLabel')}</SheetTitle>
              </SheetHeader>
              <nav className="mt-6 flex flex-col gap-2" aria-label={t('nav.siteLabel')}>
                <Button variant="ghost" className="justify-start" asChild onClick={closeMobile}>
                  <Link href="/">{t('nav.home')}</Link>
                </Button>
                <Button variant="ghost" className="justify-start" asChild onClick={closeMobile}>
                  <Link href="/tools">{t('nav.tools')}</Link>
                </Button>
                {children ? (
                  <div className="flex flex-col gap-2 border-t border-border pt-2" onClick={closeMobile}>
                    {children}
                  </div>
                ) : null}
                <div className="mt-2 flex flex-col gap-2 border-t border-border pt-2 sm:hidden">
                  {mobileAuth}
                </div>
                <div className="mt-2 border-t border-border pt-3 sm:hidden">
                  <LocaleSelect className="w-full" triggerClassName="w-full" />
                </div>
              </nav>
            </SheetContent>
          </Sheet>
        </div>
      </div>
    </header>
  )
}
