import { Link, Outlet } from 'react-router-dom'
import { useTranslation } from 'react-i18next'

import { NavLoadingBar } from '@/components/common/nav-loading-bar'
import { ThemeToggle } from '@/components/common/theme-toggle'
import { Button } from '@/components/ui/button'
import { useAuthStore } from '@/store/auth-store'

export function MainLayout() {
  const { t, i18n } = useTranslation('common')
  const accessToken = useAuthStore((state) => state.accessToken)
  const signOut = useAuthStore((state) => state.signOut)

  return (
    <div className="min-h-[100dvh] bg-background">
      <header className="relative border-b">
        <NavLoadingBar />
        <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4">
          <Link to="/" className="text-sm font-semibold">
            {t('appName')}
          </Link>
          <nav className="flex items-center gap-2">
            {accessToken ? (
              <Button variant="outline" size="sm" onClick={() => signOut()}>
                {t('auth.actions.signOut')}
              </Button>
            ) : (
              <>
                <Button variant="ghost" size="sm" asChild>
                  <Link to="/login">{t('auth.actions.signIn')}</Link>
                </Button>
                <Button size="sm" asChild>
                  <Link to="/register">{t('auth.actions.register')}</Link>
                </Button>
              </>
            )}
            <Button
              variant="outline"
              size="sm"
              onClick={() =>
                void i18n.changeLanguage(i18n.language === 'en' ? 'vi' : 'en')
              }
            >
              {i18n.language === 'en' ? 'VI' : 'EN'}
            </Button>
            <ThemeToggle />
          </nav>
        </div>
      </header>
      <main className="mx-auto max-w-7xl px-4 py-8">
        <Outlet />
      </main>
    </div>
  )
}
