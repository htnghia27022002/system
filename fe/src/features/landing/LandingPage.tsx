'use client'

import { LandingNav } from './components/landing-nav'
import { CtaSection } from './components/sections/cta-section'
import { FeaturesSection } from './components/sections/features-section'
import { HeroSection } from './components/sections/hero-section'
import { ToolsSection } from './components/sections/tools-section'

/**
 * Public home — normal document scroll with padded sections and light reveal motion.
 * Tools hub pages stay outside this shell so they share normal page scroll.
 */
export function LandingPage() {
  return (
    <div className="min-h-[100dvh] w-full bg-background text-foreground">
      <LandingNav />
      <main>
        <HeroSection />
        <FeaturesSection />
        <ToolsSection />
        <CtaSection />
      </main>
    </div>
  )
}
