import {
  CheckIcon,
  ChevronRightIcon,
  LanguagesIcon,
  LayoutDashboardIcon,
  LogOutIcon,
  MenuIcon,
  MonitorIcon,
  MoonIcon,
  PanelLeftCloseIcon,
  PanelLeftOpenIcon,
  SettingsIcon,
  SunIcon,
} from 'lucide-react'
import { useTheme } from 'next-themes'
import { useTranslation } from 'react-i18next'
import { NavLink, Outlet } from 'react-router-dom'

import { NavLoadingBar } from '@/components/common/nav-loading-bar'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Separator } from '@/components/ui/separator'
import { Sheet, SheetContent } from '@/components/ui/sheet'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'
import { useIsMobile } from '@/hooks/use-is-mobile'
import { SIDEBAR_COLLAPSED_WIDTH, useAdminSidebar } from '@/hooks/use-admin-sidebar'
import { cn } from '@/lib/utils'
import { useAuthStore } from '@/store/auth-store'

// ---------------------------------------------------------------------------
// Nav item definitions — extend here to add more routes
// ---------------------------------------------------------------------------
const NAV_ITEMS = [
  { to: '/admin', icon: LayoutDashboardIcon, labelKey: 'nav.overview', end: true },
] as const

// ---------------------------------------------------------------------------
// User avatar (initials circle)
// ---------------------------------------------------------------------------
function UserAvatar({ name, email }: { name: string; email: string }) {
  const initials = name
    .split(' ')
    .slice(0, 2)
    .map((w) => w[0]?.toUpperCase() ?? '')
    .join('')

  return (
    <span
      className="flex size-7 shrink-0 items-center justify-center rounded-full bg-primary text-[11px] font-semibold text-primary-foreground"
      aria-hidden="true"
    >
      {initials || email[0]?.toUpperCase()}
    </span>
  )
}

// ---------------------------------------------------------------------------
// User settings dropdown
// ---------------------------------------------------------------------------
function UserSettingsMenu({
  collapsed,
  onNavigate,
}: {
  collapsed: boolean
  onNavigate?: () => void
}) {
  const { t, i18n } = useTranslation('admin')
  const user = useAuthStore((state) => state.user)
  const signOut = useAuthStore((state) => state.signOut)
  const { theme, setTheme } = useTheme()

  const name = user?.name ?? user?.email ?? ''
  const email = user?.email ?? ''

  const themes = [
    { value: 'light', icon: SunIcon, label: t('userMenu.themeLight') },
    { value: 'dark', icon: MoonIcon, label: t('userMenu.themeDark') },
    { value: 'system', icon: MonitorIcon, label: t('userMenu.themeSystem') },
  ] as const

  return (
    <DropdownMenu>
      <Tooltip>
        <TooltipTrigger asChild>
          <DropdownMenuTrigger asChild>
            <button
              className={cn(
                'flex w-full items-center gap-2.5 rounded-lg px-2 py-2 text-left',
                'text-sm text-sidebar-foreground',
                'hover:bg-sidebar-accent hover:text-sidebar-accent-foreground',
                'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-sidebar-ring',
                'transition-colors duration-150',
                collapsed && 'justify-center px-0',
              )}
            >
              <UserAvatar name={name} email={email} />
              {!collapsed && (
                <span className="flex min-w-0 flex-1 flex-col overflow-hidden">
                  <span className="truncate text-[13px] font-semibold leading-tight">{name}</span>
                  <span className="truncate text-[11px] text-muted-foreground">{email}</span>
                </span>
              )}
              {!collapsed && <ChevronRightIcon className="ml-auto size-3.5 shrink-0 text-muted-foreground" />}
            </button>
          </DropdownMenuTrigger>
        </TooltipTrigger>
        {collapsed && (
          <TooltipContent side="right">
            <p className="font-semibold">{name}</p>
            <p className="text-xs text-muted-foreground">{email}</p>
          </TooltipContent>
        )}
      </Tooltip>

      <DropdownMenuContent
        side={collapsed ? 'right' : 'top'}
        align="end"
        sideOffset={8}
        className="w-64"
      >
        <DropdownMenuLabel className="pb-1">
          <p className="text-[13px] font-semibold leading-tight">{name}</p>
          <p className="text-xs font-normal text-muted-foreground">{email}</p>
        </DropdownMenuLabel>
        <DropdownMenuSeparator />

        {/* Theme */}
        <DropdownMenuGroup>
          {themes.map(({ value, icon: Icon, label }) => (
            <DropdownMenuItem
              key={value}
              onClick={() => setTheme(value)}
              className="gap-2"
            >
              <Icon className="size-4" />
              {label}
              {theme === value && <CheckIcon className="ml-auto size-3.5" />}
            </DropdownMenuItem>
          ))}
        </DropdownMenuGroup>

        <DropdownMenuSeparator />

        {/* Language */}
        <DropdownMenuItem
          onClick={() => void i18n.changeLanguage(i18n.language === 'en' ? 'vi' : 'en')}
          className="gap-2"
        >
          <LanguagesIcon className="size-4" />
          {i18n.language === 'en' ? t('userMenu.switchToVi') : t('userMenu.switchToEn')}
        </DropdownMenuItem>

        <DropdownMenuSeparator />

        {/* Settings (placeholder) */}
        <DropdownMenuItem className="gap-2" disabled>
          <SettingsIcon className="size-4" />
          {t('userMenu.settings')}
        </DropdownMenuItem>

        <DropdownMenuSeparator />

        {/* Sign out */}
        <DropdownMenuItem
          className="gap-2 text-destructive focus:text-destructive"
          onClick={() => {
            onNavigate?.()
            signOut()
          }}
        >
          <LogOutIcon className="size-4" />
          {t('userMenu.signOut')}
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}

// ---------------------------------------------------------------------------
// Sidebar inner content (used in both desktop aside and mobile Sheet)
// ---------------------------------------------------------------------------
function SidebarContent({
  collapsed,
  onNavigate,
}: {
  collapsed: boolean
  onNavigate?: () => void
}) {
  const { t } = useTranslation('admin')

  return (
    <div className="flex flex-1 flex-col overflow-hidden">
      {/* Brand header */}
      <div
        className={cn(
          'flex shrink-0 flex-col px-3 pb-3 pt-4',
          collapsed && 'items-center px-0',
        )}
      >
        {collapsed ? (
          <span
            className="flex size-7 items-center justify-center rounded-md bg-primary text-[11px] font-bold text-primary-foreground"
            aria-label={t('shell.badge')}
          >
            BO
          </span>
        ) : (
          <>
            <p className="text-[11px] font-semibold uppercase tracking-widest text-muted-foreground">
              {t('shell.badge')}
            </p>
            <p className="text-[17px] font-semibold leading-tight tracking-[-0.2px] text-foreground">
              {t('shell.title')}
            </p>
          </>
        )}
      </div>

      <Separator />

      {/* Navigation */}
      <nav className={cn('flex flex-1 flex-col gap-0.5 overflow-y-auto p-2', collapsed && 'items-center')}>
        {NAV_ITEMS.map(({ to, icon: Icon, labelKey, end }) => (
          <NavLink key={to} to={to} end={end} onClick={onNavigate}>
            {({ isActive }) => (
              <Tooltip>
                <TooltipTrigger asChild>
                  <span
                    role="menuitem"
                    className={cn(
                      'flex w-full cursor-pointer items-center gap-2.5 rounded-lg px-2.5 py-2',
                      'text-[13px] font-medium text-sidebar-foreground',
                      'hover:bg-sidebar-accent hover:text-sidebar-accent-foreground',
                      'transition-colors duration-150',
                      isActive && 'bg-sidebar-accent text-sidebar-accent-foreground',
                      collapsed && 'w-9 justify-center px-0',
                    )}
                  >
                    <Icon className="size-4 shrink-0" />
                    {!collapsed && (
                      <span className="truncate">{t(labelKey)}</span>
                    )}
                  </span>
                </TooltipTrigger>
                {collapsed && (
                  <TooltipContent side="right">{t(labelKey)}</TooltipContent>
                )}
              </Tooltip>
            )}
          </NavLink>
        ))}
      </nav>

      <Separator />

      {/* User settings */}
      <div className={cn('shrink-0 p-2', collapsed && 'flex justify-center')}>
        <UserSettingsMenu collapsed={collapsed} onNavigate={onNavigate} />
      </div>
    </div>
  )
}

// ---------------------------------------------------------------------------
// Desktop sidebar (resizable)
// ---------------------------------------------------------------------------
type DesktopSidebarProps = {
  width: number
  isCollapsed: boolean
  isResizing: boolean
  toggleCollapsed: () => void
  handleResizeStart: (e: React.MouseEvent) => void
}

function DesktopSidebar({
  width,
  isCollapsed,
  isResizing,
  toggleCollapsed,
  handleResizeStart,
}: DesktopSidebarProps) {
  const { t } = useTranslation('admin')

  return (
    <aside
      style={{ width }}
      className={cn(
        'relative hidden shrink-0 flex-col border-r border-sidebar-border bg-sidebar',
        'md:flex',
        !isResizing && 'transition-[width] duration-200 ease-in-out',
      )}
    >
      <SidebarContent collapsed={isCollapsed} />

      {/* Collapse toggle button — pinned to top of sidebar */}
      <div className="absolute -right-3.5 top-4 z-10">
        <Tooltip>
          <TooltipTrigger asChild>
            <Button
              variant="outline"
              size="icon"
              className="size-7 rounded-full border-border bg-background shadow-sm hover:bg-sidebar-accent"
              onClick={toggleCollapsed}
              aria-label={isCollapsed ? t('shell.expandMenu') : t('shell.collapseMenu')}
              aria-expanded={!isCollapsed}
            >
              {isCollapsed ? (
                <PanelLeftOpenIcon className="size-3.5" />
              ) : (
                <PanelLeftCloseIcon className="size-3.5" />
              )}
            </Button>
          </TooltipTrigger>
          <TooltipContent side="right">
            {isCollapsed ? t('shell.expandMenu') : t('shell.collapseMenu')}
          </TooltipContent>
        </Tooltip>
      </div>

      {/* Resize handle — right edge, only visible when not collapsed */}
      {!isCollapsed && (
        <div
          className={cn(
            'absolute inset-y-0 -right-1 z-20 w-2 cursor-col-resize',
            'group/handle flex items-center justify-center',
          )}
          onMouseDown={handleResizeStart}
          aria-hidden="true"
        >
          <div
            className={cn(
              'h-8 w-0.5 rounded-full',
              'bg-transparent transition-colors duration-150',
              'group-hover/handle:bg-primary/50',
              isResizing && 'bg-primary',
            )}
          />
        </div>
      )}
    </aside>
  )
}

// ---------------------------------------------------------------------------
// Mobile sidebar (Sheet)
// ---------------------------------------------------------------------------
function MobileSidebar({
  isOpen,
  onClose,
}: {
  isOpen: boolean
  onClose: () => void
}) {
  const { t } = useTranslation('admin')

  return (
    <Sheet open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <SheetContent
        side="left"
        className="flex flex-col p-0"
        style={{ width: 260 }}
        title={t('shell.title')}
      >
          <SidebarContent collapsed={false} onNavigate={onClose} />
      </SheetContent>
    </Sheet>
  )
}

// ---------------------------------------------------------------------------
// Root layout
// ---------------------------------------------------------------------------
export function AdminLayout() {
  const { t } = useTranslation('admin')
  const isMobile = useIsMobile()
  const { width, isCollapsed, isMobileOpen, isResizing, toggleCollapsed, openMobile, closeMobile, handleResizeStart } =
    useAdminSidebar()

  const effectiveCollapseWidth = isCollapsed ? SIDEBAR_COLLAPSED_WIDTH : width

  return (
    <TooltipProvider delayDuration={300}>
    <div className="flex min-h-[100dvh] bg-background">
      {/* Desktop sidebar */}
      <DesktopSidebar
        width={effectiveCollapseWidth}
        isCollapsed={isCollapsed}
        isResizing={isResizing}
        toggleCollapsed={toggleCollapsed}
        handleResizeStart={handleResizeStart}
      />

      {/* Mobile sheet */}
      {isMobile && (
        <MobileSidebar isOpen={isMobileOpen} onClose={closeMobile} />
      )}

      {/* Main content area */}
      <div className="flex min-w-0 flex-1 flex-col">
        <header className="relative border-b border-border bg-background">
          <NavLoadingBar />
          <div className="flex h-14 items-center gap-2 px-4">
            {/* Mobile: hamburger */}
            <Button
              variant="ghost"
              size="icon"
              className="md:hidden"
              onClick={openMobile}
              aria-label={t('shell.openMenu')}
            >
              <MenuIcon className="size-4" />
            </Button>

            {/* Mobile: brand name */}
            <span className="text-[15px] font-semibold md:hidden">{t('shell.title')}</span>

            {/* Spacer */}
            <div className="flex-1" />
          </div>
        </header>

        <main className="flex-1 bg-secondary/40 p-4 md:p-6">
          <Outlet />
        </main>
      </div>
    </div>
    </TooltipProvider>
  )
}
