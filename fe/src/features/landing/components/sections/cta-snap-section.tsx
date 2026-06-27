'use client'

import Link from 'next/link'

import { Button } from '@/components/ui/button'

import { landingContent } from '../../content'
import { SplitWords } from '../split-words'
import { SectionSnapBase } from './section-snap-base'

export function CtaSnapSection() {
  const { cta } = landingContent
  const year = new Date().getFullYear()
  const footer = cta.footer.replace('{year}', String(year))

  return (
    <SectionSnapBase id="cta" className="landing-section-cta">
      <div className="rounded-3xl border border-primary/30 bg-card/50 p-8 sm:p-12 landing-glow">
        <p className="text-xs font-semibold uppercase tracking-[0.25em] text-primary">{cta.label}</p>
        <SplitWords
          as="h2"
          text={cta.title}
          className="mt-4 font-heading text-3xl font-semibold tracking-tight sm:text-4xl"
        />
        <p className="mt-4 max-w-xl text-muted-foreground">{cta.subtitle}</p>
        <div className="mt-10 flex flex-wrap gap-3">
          <Button size="lg" asChild>
            <Link href="/register">{cta.primaryCta}</Link>
          </Button>
          <Button variant="outline" size="lg" asChild>
            <a href={cta.emailHref}>{cta.secondaryCta}</a>
          </Button>
        </div>
        <p className="mt-6 text-sm text-muted-foreground">{cta.email}</p>
      </div>
      <footer className="mt-12 text-center text-sm text-muted-foreground">{footer}</footer>
    </SectionSnapBase>
  )
}
