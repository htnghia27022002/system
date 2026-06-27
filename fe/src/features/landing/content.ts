export const landingContent = {
  hero: {
    label: '01 — Hero',
    eyebrow: 'Web3 Prototype',
    headline: 'Scroll-driven experiences for the decentralized web.',
    subline:
      'A scrollytelling landing prototype — pinned sections, scrubbed motion, and neon cyber aesthetics. Personal details coming soon.',
    primaryCta: 'Explore features',
    secondaryCta: 'Join waitlist',
  },
  features: {
    label: '02 — Features',
    title: 'Built for momentum',
    subtitle:
      'Pinned scroll zone — artifacts animate in sync with your wheel. Scroll slowly to scrub through the sequence.',
    items: [
      { title: 'On-chain ready', description: 'Wallet-native flows and smart-contract hooks.' },
      { title: 'Realtime sync', description: 'Sub-second state across nodes and clients.' },
      { title: 'Composable UI', description: 'Glass panels, glow borders, token-driven theme.' },
    ],
  },
  tokenomics: {
    label: '03 — Tokenomics',
    title: 'Allocation that stacks',
    subtitle: 'Cards stack as you scroll — each layer reveals the next allocation tier.',
    cards: [
      { title: 'Community', share: '40%', description: 'Ecosystem grants, builders, and governance.' },
      { title: 'Treasury', share: '35%', description: 'Protocol reserves and long-term development.' },
      { title: 'Core team', share: '25%', description: 'Vested allocation for contributors.' },
    ],
  },
  cta: {
    label: '04 — CTA',
    title: 'Ready to ship your Web3 surface?',
    subtitle: 'Prototype landing — replace with your launch copy, links, and contract address.',
    primaryCta: 'Get started',
    secondaryCta: 'Read docs',
    email: 'hello@example.com',
    emailHref: 'mailto:hello@example.com',
    footer: '© {year} htnghia. GSAP scrollytelling prototype.',
  },
  nav: {
    hero: 'Hero',
    features: 'Features',
    tokenomics: 'Tokenomics',
    cta: 'CTA',
  },
} as const

export type LandingContent = typeof landingContent

export const scrollySectionIds = ['hero', 'features', 'tokenomics', 'cta'] as const
