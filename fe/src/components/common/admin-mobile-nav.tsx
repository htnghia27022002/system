'use client'

import Link from 'next/link'
import {
  LayoutGridIcon,
  MenuIcon,
  PinIcon,
  PinOffIcon,
} from 'lucide-react'
import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'

import { UserInfo } from '@/components/common/user-info'
import { UserMenuContent } from '@/components/common/user-menu-content'
import {
  ADMIN_HOME_HREF,
  useAdminNavItems,
  useAdminNavLeaves,
} from '@/components/common/use-admin-nav-items'
import { usePinnedAdminNav } from '@/components/common/use-pinned-admin-nav'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { useSidebar } from '@/components/ui/sidebar'
import { useCurrentPath } from '@/hooks/use-current-path'
import { cn } from '@/lib/utils'
import { useAuthStore } from '@/store/auth-store'
import type { NavItem } from '@/types/navigation'

/** Inner tab row height; total bar = this + safe-area-inset-bottom. */
const BOTTOM_NAV_ROW_HEIGHT = '3rem'

function BottomTab({
  href,
  title,
  icon: Icon,
  active,
  onNavigate,
}: {
  href: string
  title: string
  icon?: NavItem['icon']
  active: boolean
  onNavigate?: () => void
}) {
  const TabIcon = Icon ?? LayoutGridIcon

  return (
    <Link
      href={href}
      onClick={onNavigate}
      aria-current={active ? 'page' : undefined}
      className={cn(
        'relative flex min-w-0 flex-1 flex-col items-center justify-center gap-0.5 px-1 py-0.5 text-[10px] font-medium leading-none',
        active ? 'text-foreground' : 'text-muted-foreground',
      )}
    >
      <TabIcon
        className={cn('size-[1.125rem] shrink-0', active && 'text-primary')}
        aria-hidden
      />
      <span className="max-w-full truncate">{title}</span>
      {active ? (
        <span
          className="absolute inset-x-0 bottom-1 mx-auto h-0.5 w-4 rounded-full bg-primary"
          aria-hidden
        />
      ) : null}
    </Link>
  )
}

function MenuNavRow({
  item,
  onNavigate,
  showPin,
  pinned,
  allowPin,
  onTogglePin,
  pinLabel,
  unpinLabel,
}: {
  item: NavItem
  onNavigate: () => void
  showPin: boolean
  pinned?: boolean
  allowPin?: boolean
  onTogglePin?: () => void
  pinLabel?: string
  unpinLabel?: string
}) {
  const Icon = item.icon ?? LayoutGridIcon

  return (
    <div className="relative flex w-full items-center rounded-2xl border bg-card shadow-sm">
      <Link
        href={item.href}
        onClick={onNavigate}
        className={cn(
          'flex min-w-0 flex-1 items-center gap-3 px-3 py-3',
          showPin && 'pr-12',
        )}
      >
        <span className="flex size-10 shrink-0 items-center justify-center rounded-xl bg-muted">
          <Icon className="size-5" aria-hidden />
        </span>
        <span className="min-w-0 truncate text-sm font-medium">{item.title}</span>
      </Link>

      {showPin ? (
        <Button
          type="button"
          variant="ghost"
          size="icon-sm"
          className="absolute top-1/2 right-2 size-8 -translate-y-1/2 shrink-0"
          disabled={!allowPin && !pinned}
          aria-label={pinned ? unpinLabel : pinLabel}
          aria-pressed={pinned}
          onClick={onTogglePin}
        >
          {pinned ? (
            <PinOffIcon className="size-3.5 text-primary" />
          ) : (
            <PinIcon className="size-3.5 text-muted-foreground" />
          )}
        </Button>
      ) : null}
    </div>
  )
}

/**
 * Facebook-style mobile admin chrome:
 * - Bottom bar: Home + user-pinned shortcuts + Menu (hamburger)
 * - Menu panel: profile, shortcuts, hierarchical full-width rows with corner pin
 * Desktop (≥768) renders nothing.
 */
export function AdminMobileNav() {
  const { t } = useTranslation('admin')
  const { isMobile } = useSidebar()
  const { pathname, isCurrentPath } = useCurrentPath()
  const user = useAuthStore((state) => state.user)
  const navItems = useAdminNavItems()
  const leaves = useAdminNavLeaves()
  const {
    pinnedItems,
    isPinned,
    canPin,
    togglePin,
    maxPins,
  } = usePinnedAdminNav()

  const [menuOpen, setMenuOpen] = useState(false)

  useEffect(() => {
    setMenuOpen(false)
  }, [pathname])

  if (!isMobile) {
    return null
  }

  const homeLeaf =
    leaves.find((item) => item.href === ADMIN_HOME_HREF) ??
    ({
      title: t('nav.dashboard'),
      href: ADMIN_HOME_HREF,
      icon: LayoutGridIcon,
    } satisfies NavItem)

  const platformLeaves = navItems.filter((item) => !item.items?.length)
  const sectionGroups = navItems.filter(
    (item) => item.items && item.items.length > 0,
  )

  const closeMenu = () => setMenuOpen(false)

  const renderPinnableRow = (item: NavItem) => {
    const isHome = item.href === ADMIN_HOME_HREF
    const pinned = isPinned(item.href)
    return (
      <MenuNavRow
        key={item.href}
        item={item}
        onNavigate={closeMenu}
        showPin={!isHome}
        pinned={pinned}
        allowPin={canPin(item.href)}
        onTogglePin={() => togglePin(item.href)}
        pinLabel={t('shell.pinShortcut', { name: item.title })}
        unpinLabel={t('shell.unpinShortcut', { name: item.title })}
      />
    )
  }

  return (
    <>
      {menuOpen ? (
        <div
          className="safe-area-top fixed inset-x-0 top-0 z-40 flex flex-col bg-background md:hidden"
          style={{
            bottom: `calc(${BOTTOM_NAV_ROW_HEIGHT} + env(safe-area-inset-bottom, 0px))`,
          }}
          role="dialog"
          aria-modal="true"
          aria-label={t('shell.menuTitle')}
        >
          <header className="flex shrink-0 items-center justify-between gap-2 border-b px-4 py-3">
            <h1 className="text-xl font-semibold tracking-tight">
              {t('shell.menuTitle')}
            </h1>
          </header>

          <div className="min-h-0 flex-1 overflow-y-auto overscroll-contain px-4 py-4">
            {user ? (
              <div className="mb-5 rounded-2xl border bg-card p-3 shadow-sm">
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <button
                      type="button"
                      className="flex w-full items-center gap-3 rounded-xl p-1 text-left outline-none focus-visible:ring-2 focus-visible:ring-ring"
                    >
                      <UserInfo user={user} showEmail />
                    </button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent
                    className="w-[min(calc(100vw-2rem),20rem)] rounded-lg"
                    align="start"
                    side="bottom"
                    sideOffset={8}
                  >
                    <UserMenuContent user={user} onNavigate={closeMenu} />
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
            ) : null}

            {pinnedItems.length > 0 ? (
              <section className="mb-6">
                <h2 className="mb-2 text-sm font-semibold">
                  {t('shell.yourShortcuts')}
                </h2>
                <div className="-mx-1 flex gap-2 overflow-x-auto px-1 pb-1 [-webkit-overflow-scrolling:touch]">
                  {pinnedItems.map((item) => {
                    const Icon = item.icon ?? LayoutGridIcon
                    return (
                      <Link
                        key={item.href}
                        href={item.href}
                        onClick={closeMenu}
                        className="flex w-20 shrink-0 flex-col items-center gap-1.5 rounded-xl border bg-card p-2 text-center shadow-sm"
                      >
                        <span className="flex size-12 items-center justify-center rounded-xl bg-muted">
                          <Icon className="size-6" aria-hidden />
                        </span>
                        <span className="line-clamp-2 text-[11px] font-medium leading-tight">
                          {item.title}
                        </span>
                      </Link>
                    )
                  })}
                </div>
              </section>
            ) : null}

            <section className="space-y-5">
              <div className="flex items-end justify-between gap-2">
                <h2 className="text-sm font-semibold">{t('shell.allShortcuts')}</h2>
                <p className="text-[11px] text-muted-foreground">
                  {t('shell.pinHint', { count: maxPins })}
                </p>
              </div>

              {platformLeaves.length > 0 ? (
                <div className="space-y-2">
                  <h3 className="px-0.5 text-xs font-medium tracking-wide text-muted-foreground uppercase">
                    {t('nav.platform')}
                  </h3>
                  <div className="flex flex-col gap-2">
                    {platformLeaves.map(renderPinnableRow)}
                  </div>
                </div>
              ) : null}

              {sectionGroups.map((group) => (
                <div key={group.title} className="space-y-2">
                  <h3 className="px-0.5 text-xs font-medium tracking-wide text-muted-foreground uppercase">
                    {group.title}
                  </h3>
                  <div className="flex flex-col gap-2">
                    {group.items?.map(renderPinnableRow)}
                  </div>
                </div>
              ))}
            </section>
          </div>
        </div>
      ) : null}

      <nav
        className="fixed inset-x-0 bottom-0 z-50 border-t bg-background/95 pb-[env(safe-area-inset-bottom)] backdrop-blur-md md:hidden"
        aria-label={t('shell.bottomNav')}
      >
        <div className="flex h-12 items-stretch">
          <BottomTab
            href={homeLeaf.href}
            title={homeLeaf.title}
            icon={homeLeaf.icon}
            active={!menuOpen && isCurrentPath(homeLeaf.href)}
            onNavigate={closeMenu}
          />

          {pinnedItems.map((item) => (
            <BottomTab
              key={item.href}
              href={item.href}
              title={item.title}
              icon={item.icon}
              active={!menuOpen && isCurrentPath(item.href)}
              onNavigate={closeMenu}
            />
          ))}

          <button
            type="button"
            className={cn(
              'relative flex min-w-0 flex-1 flex-col items-center justify-center gap-0.5 px-1 py-0.5 text-[10px] font-medium leading-none',
              menuOpen ? 'text-foreground' : 'text-muted-foreground',
            )}
            aria-label={t('shell.openMenu')}
            aria-pressed={menuOpen}
            onClick={() => setMenuOpen((open) => !open)}
          >
            <MenuIcon
              className={cn('size-[1.125rem] shrink-0', menuOpen && 'text-primary')}
              aria-hidden
            />
            <span className="max-w-full truncate">{t('shell.menuTab')}</span>
            {menuOpen ? (
              <span
                className="absolute inset-x-0 bottom-1 mx-auto h-0.5 w-4 rounded-full bg-primary"
                aria-hidden
              />
            ) : null}
          </button>
        </div>
      </nav>
    </>
  )
}
