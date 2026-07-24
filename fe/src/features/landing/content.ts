/**
 * Greenfield landing marketing copy (English).
 * Tool entries are never duplicated here — always from `@/features/tools`.
 */
export const landingContent = {
  hero: {
    brand: 'System',
    headline: 'Tools your team actually opens.',
    subline:
      'A focused workspace home: browse the catalog, jump into Webhooks, and grow the hub without redesigning the product shell.',
    primaryCta: 'Explore tools',
    primaryCtaHref: '/tools',
    secondaryCta: 'Create account',
    secondaryCtaHref: '/register',
    imageAlt: 'Stylized 3D avatar based on the site owner',
  },
  features: {
    title: 'Built for daily use',
    subtitle: 'Less ceremony, more paths into real work.',
    items: [
      {
        title: 'One catalog',
        description: 'Landing and /tools read the same entries. Add a tool once — it shows up everywhere.',
      },
      {
        title: 'Nested workspaces',
        description: 'Each tool can own /tools/{id}. Webhooks is live; more rooms ship the same way.',
      },
      {
        title: 'Ready chrome',
        description: 'Navigate, sign in, theme, and language stay consistent across home and tools.',
      },
    ],
  },
  tools: {
    title: 'Start with Webhooks',
    subtitle: 'A short preview from the shared catalog. Open the hub for the full list.',
    primaryCta: 'Browse all tools',
    primaryCtaHref: '/tools',
    emptyTitle: 'Tools coming soon',
    emptyBody: 'The catalog is empty for now. Check back after the next release.',
  },
  cta: {
    title: 'Create an account',
    subtitle: 'Sign up for admin access, or sign in if you already have a seat.',
    primaryCta: 'Create account',
    primaryCtaHref: '/register',
    secondaryCta: 'Sign in',
    secondaryCtaHref: '/login',
    footer: '© {year} System',
  },
  nav: {
    hero: 'Home',
    features: 'Value',
    tools: 'Preview',
    cta: 'Start',
  },
} as const

export type LandingContent = typeof landingContent

/** In-page section ids for nav highlighting. */
export const landingSectionIds = ['hero', 'features', 'tools', 'cta'] as const
