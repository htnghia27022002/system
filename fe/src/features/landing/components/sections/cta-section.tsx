'use client'

import Link from 'next/link'

import { Button } from '@/components/ui/button'

import { landingContent } from '../../content'
import { SectionReveal } from '../section-reveal'
import { SectionBase } from './section-base'

export function CtaSection() {
  const { cta } = landingContent
  const year = new Date().getFullYear()
  const footer = cta.footer.replace('{year}', String(year))

  return (
    <SectionBase id="cta" className="landing-section-cta pb-16 sm:pb-20">
      <SectionReveal>
        <div className="rounded-2xl border border-border bg-card px-6 py-10 sm:px-10 sm:py-12">
          <h2 className="font-heading text-3xl font-semibold tracking-tight sm:text-4xl">
            {cta.title}
          </h2>
          <p className="mt-3 max-w-xl text-muted-foreground">{cta.subtitle}</p>
          <div className="mt-8 flex flex-wrap gap-3">
            <Button size="lg" asChild>
              <Link href={cta.primaryCtaHref}>{cta.primaryCta}</Link>
            </Button>
            <Button variant="outline" size="lg" asChild>
              <Link href={cta.secondaryCtaHref}>{cta.secondaryCta}</Link>
            </Button>
          </div>
        </div>
      </SectionReveal>
      <footer className="mt-10 text-center text-sm text-muted-foreground">{footer}</footer>
    </SectionBase>
  )
}
