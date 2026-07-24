import { render, screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'

vi.mock('framer-motion', () => ({
  motion: {
    div: ({ children, ...props }: React.ComponentProps<'div'>) => {
      const {
        initial,
        animate,
        transition,
        whileInView,
        viewport,
        variants,
        ...rest
      } = props as Record<string, unknown>
      void initial
      void animate
      void transition
      void whileInView
      void viewport
      void variants
      return <div {...rest}>{children}</div>
    },
    header: ({ children, ...props }: React.ComponentProps<'header'>) => {
      const { initial, animate, transition, whileHover, ...rest } = props as Record<string, unknown>
      void initial
      void animate
      void transition
      void whileHover
      return <header {...rest}>{children}</header>
    },
  },
}))

vi.mock('next/dynamic', () => ({
  default: () => {
    function MockHeroStage() {
      return <div data-testid="hero-character-stage" />
    }
    return MockHeroStage
  },
}))

vi.mock('next/image', () => ({
  default: ({ alt, priority: _priority, fill: _fill, sizes: _sizes, ...props }: {
    alt: string
    priority?: boolean
    fill?: boolean
    sizes?: string
    src?: string
  }) => <img alt={alt} src={typeof props.src === 'string' ? props.src : ''} />,
}))

vi.mock('next/link', () => ({
  default: ({ children, href }: { children: React.ReactNode; href: string }) => (
    <a href={href}>{children}</a>
  ),
}))

vi.mock('@/store/auth-store', () => ({
  useAuthStore: (selector: (state: unknown) => unknown) =>
    selector({ accessToken: null }),
}))

vi.mock('@/features/auth/hooks/use-sign-out', () => ({
  useSignOut: () => ({ mutate: vi.fn(), isPending: false }),
}))

vi.mock('react-i18next', () => ({
  useTranslation: () => ({
    t: (key: string) => {
      const map: Record<string, string> = {
        'nav.home': 'Home',
        'nav.tools': 'Tools',
        'nav.admin': 'Admin',
        'nav.siteLabel': 'Site',
        'nav.openMenu': 'Open navigation menu',
        'locale.label': 'Language',
        'auth.actions.signIn': 'Sign in',
        'auth.actions.register': 'Register',
        'auth.actions.signOut': 'Sign out',
        'theme.toggle': 'Toggle theme',
      }
      return map[key] ?? key
    },
    i18n: { language: 'en', changeLanguage: vi.fn() },
  }),
}))

vi.mock('next-themes', () => ({
  useTheme: () => ({ theme: 'system', setTheme: vi.fn(), resolvedTheme: 'light' }),
}))

import { LandingPage } from './LandingPage'

describe('LandingPage', () => {
  it('renders padded landing sections without scrolly shell or tokenomics', () => {
    render(<LandingPage />)

    expect(document.getElementById('hero')).toHaveTextContent(/Tools your team actually opens/i)
    expect(document.getElementById('hero')).toBeInTheDocument()
    expect(document.getElementById('features')).toBeInTheDocument()
    expect(document.getElementById('tools')).toBeInTheDocument()
    expect(document.getElementById('tokenomics')).not.toBeInTheDocument()
    expect(document.getElementById('cta')).toBeInTheDocument()
    expect(document.getElementById('scrolly-container')).not.toBeInTheDocument()
    expect(screen.getByRole('link', { name: /Browse all tools/i })).toHaveAttribute('href', '/tools')
    expect(screen.getAllByRole('link', { name: /^Tools$/i }).length).toBeGreaterThan(0)
    expect(screen.getByRole('combobox', { name: /Language/i })).toBeInTheDocument()
  })
})
