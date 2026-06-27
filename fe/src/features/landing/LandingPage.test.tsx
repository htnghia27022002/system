import { render } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'

vi.mock('gsap', () => ({
  default: {
    registerPlugin: vi.fn(),
    matchMedia: () => ({ add: vi.fn(), revert: vi.fn() }),
    context: vi.fn(() => ({ revert: vi.fn() })),
    timeline: vi.fn(() => ({ scrollTrigger: null })),
    fromTo: vi.fn(() => ({ scrollTrigger: null })),
    set: vi.fn(),
    utils: { toArray: vi.fn(() => []) },
  },
}))

vi.mock('gsap/ScrollTrigger', () => ({
  ScrollTrigger: {
    create: vi.fn(),
    refresh: vi.fn(),
    defaults: vi.fn(),
  },
}))

vi.mock('framer-motion', () => ({
  motion: {
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

vi.mock('next/image', () => ({
  default: ({ alt, priority: _priority, ...props }: { alt: string; priority?: boolean }) => (
    <img alt={alt} {...props} />
  ),
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
    t: (key: string) => key,
    i18n: { language: 'en', changeLanguage: vi.fn() },
  }),
}))

import { LandingPage } from './LandingPage'

describe('LandingPage', () => {
  it('renders scrollytelling sections', () => {
    render(<LandingPage />)

    expect(document.getElementById('hero')).toHaveTextContent(/Scroll-driven experiences/i)
    expect(document.getElementById('hero')).toBeInTheDocument()
    expect(document.getElementById('features')).toBeInTheDocument()
    expect(document.getElementById('tokenomics')).toBeInTheDocument()
    expect(document.getElementById('cta')).toBeInTheDocument()
    expect(document.getElementById('scrolly-container')).toBeInTheDocument()
  })
})
