'use client'

import Image from 'next/image'
import Link from 'next/link'

import { Button } from '@/components/ui/button'

import { landingContent } from '../../content'
import { SectionBase } from './section-base'

export function HeroSection() {
  const { hero } = landingContent

  return (
    <SectionBase
      id="hero"
      className="landing-section-hero site-header-offset overflow-hidden"
      contentClassName="max-w-6xl lg:max-w-7xl"
      backdrop={
        <>
          <div className="pointer-events-none absolute inset-0 landing-grid-bg" aria-hidden />
          <div className="pointer-events-none absolute inset-0 landing-glow" aria-hidden />
        </>
      }
    >
      <div className="grid items-center gap-8 md:gap-10 lg:grid-cols-[minmax(0,1.05fr)_minmax(0,0.95fr)] lg:gap-12 xl:gap-14">
        <div className="order-2 min-w-0 lg:order-1">
          <p className="text-xs font-semibold uppercase tracking-[0.18em] text-primary">
            {hero.brand}
          </p>
          <h1 className="mt-4 max-w-xl font-heading text-4xl font-semibold leading-[1.08] tracking-tight text-foreground sm:text-5xl lg:text-[3.25rem]">
            {hero.headline}
          </h1>
          <p className="mt-5 max-w-lg text-base leading-relaxed text-muted-foreground sm:text-lg">
            {hero.subline}
          </p>
          <div className="mt-8 flex flex-wrap gap-3">
            <Button size="lg" asChild>
              <Link href={hero.primaryCtaHref}>{hero.primaryCta}</Link>
            </Button>
            <Button variant="outline" size="lg" asChild>
              <Link href={hero.secondaryCtaHref}>{hero.secondaryCta}</Link>
            </Button>
          </div>
        </div>

        <div className="order-1 w-full lg:order-2 lg:justify-self-end">
          <div className="relative mx-auto aspect-[3/2] w-full max-w-xl overflow-hidden rounded-2xl border border-border/70 shadow-sm">
            <Image
              src="/images/landing-hero.webp"
              alt=""
              fill
              priority
              sizes="(min-width: 1024px) 45vw, 100vw"
              className="object-cover"
            />
          </div>
        </div>
      </div>
    </SectionBase>
  )
}
