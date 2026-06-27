'use client'

import Link from 'next/link'

import { Button } from '@/components/ui/button'

import { landingContent } from '../../content'
import { useScrolly } from '../../context/scrolly-context'
import { SplitWords } from '../split-words'
import { SectionSnapBase } from './section-snap-base'

export function HeroSnapSection() {
  const { hero } = landingContent
  const { scrollToSection } = useScrolly()

  return (
    <SectionSnapBase id="hero" className="landing-section-hero landing-grid-bg">
      <div className="pointer-events-none absolute inset-0 landing-glow" aria-hidden />
      <p className="text-xs font-semibold uppercase tracking-[0.25em] text-primary">{hero.label}</p>
      <p className="mt-4 text-sm font-medium uppercase tracking-widest text-muted-foreground">
        {hero.eyebrow}
      </p>
      <SplitWords
        as="h1"
        text={hero.headline}
        className="mt-6 max-w-3xl font-heading text-4xl font-semibold leading-[1.05] tracking-tight text-foreground sm:text-5xl lg:text-6xl"
      />
      <p className="mt-6 max-w-2xl text-base leading-relaxed text-muted-foreground sm:text-lg">
        {hero.subline}
      </p>
      <div className="mt-10 flex flex-wrap gap-3">
        <Button size="lg" onClick={() => scrollToSection('features')}>
          {hero.primaryCta}
        </Button>
        <Button variant="outline" size="lg" asChild>
          <Link href="/register">{hero.secondaryCta}</Link>
        </Button>
      </div>
    </SectionSnapBase>
  )
}
