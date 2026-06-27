'use client'

import { CtaSnapSection } from './components/sections/cta-snap-section'
import { FeaturesSnapSection } from './components/sections/features-snap-section'
import { HeroSnapSection } from './components/sections/hero-snap-section'
import { TokenomicsSnapSection } from './components/sections/tokenomics-snap-section'
import { LandingNav } from './components/landing-nav'
import { ScrollProgressHud } from './components/scroll-progress-hud'
import { ScrollyShell } from './components/scrolly-shell'
import { ScrollyProvider, useScrolly } from './context/scrolly-context'
import { useGsapContext } from './hooks/use-gsap-context'
import { useReducedMotion } from './hooks/use-reduced-motion'

function LandingPageInner() {
  const { scrollerRef } = useScrolly()
  const reducedMotion = useReducedMotion()

  useGsapContext({ scrollerRef, enabled: !reducedMotion })

  return (
    <>
      <LandingNav />
      <ScrollyShell snapMode={reducedMotion ? 'mandatory' : 'proximity'}>
        <HeroSnapSection />
        <FeaturesSnapSection />
        <TokenomicsSnapSection />
        <CtaSnapSection />
      </ScrollyShell>
      <ScrollProgressHud />
    </>
  )
}

export function LandingPage() {
  return (
    <ScrollyProvider>
      <LandingPageInner />
    </ScrollyProvider>
  )
}
