import { Link, Outlet } from 'react-router-dom'
import { useTranslation } from 'react-i18next'

import { NavLoadingBar } from '@/components/common/nav-loading-bar'
import { ThemeToggle } from '@/components/common/theme-toggle'
import { Button } from '@/components/ui/button'

export function AuthLayout() {
  const { t } = useTranslation('common')

  return (
    <div className="flex min-h-[100dvh] flex-col bg-secondary">
      <header className="relative flex h-11 items-center bg-[var(--ds-surface-black)] text-[var(--ds-body-on-dark)]">
        <NavLoadingBar />
        <div className="mx-auto flex h-full w-full max-w-7xl items-center justify-between px-4 sm:px-6">
          <Button
            variant="link"
            className="h-auto p-0 text-[12px] font-semibold tracking-[-0.12px] text-white hover:text-white/90"
            asChild
          >
            <Link to="/">{t('appName')}</Link>
          </Button>
          <nav
            className="flex items-center gap-2 sm:gap-4"
            aria-label={t('auth.nav.label')}
          >
            <Button
              variant="link"
              className="h-auto p-0 text-[12px] tracking-[-0.12px] text-white/80 hover:text-white"
              asChild
            >
              <Link to="/login">{t('auth.actions.signIn')}</Link>
            </Button>
            <Button
              variant="link"
              className="h-auto p-0 text-[12px] tracking-[-0.12px] text-white/80 hover:text-white"
              asChild
            >
              <Link to="/register">{t('auth.actions.register')}</Link>
            </Button>
            <ThemeToggle className="border-white/25 bg-transparent text-white hover:bg-white/10 hover:text-white" />
          </nav>
        </div>
      </header>

      <div className="flex h-[52px] items-center border-b border-black/8 bg-secondary/80 backdrop-blur-xl backdrop-saturate-150 dark:border-white/10">
        <div className="mx-auto flex h-full w-full max-w-7xl items-center px-4 sm:px-6">
          <p className="text-[12px] leading-none tracking-[-0.12px] text-muted-foreground">
            {t('auth.layout.subtitle')}
          </p>
        </div>
      </div>

      <main className="flex flex-1 items-center justify-center px-4 py-10 sm:py-14">
        <Outlet />
      </main>
    </div>
  )
}
