'use client'

import { landingContent } from '../../content'
import { SectionReveal } from '../section-reveal'
import { SectionBase } from './section-base'

export function FeaturesSection() {
  const { features } = landingContent

  return (
    <SectionBase id="features" className="landing-section-features">
      <SectionReveal>
        <h2 className="max-w-xl font-heading text-3xl font-semibold tracking-tight sm:text-4xl">
          {features.title}
        </h2>
        <p className="mt-3 max-w-xl text-muted-foreground">{features.subtitle}</p>
      </SectionReveal>

      <ol className="mt-12 grid gap-8 sm:grid-cols-3">
        {features.items.map((item, index) => (
          <li key={item.title}>
            <SectionReveal delay={0.05 * (index + 1)}>
              <p className="font-mono text-xs font-medium text-primary/80">
                {String(index + 1).padStart(2, '0')}
              </p>
              <h3 className="mt-3 font-heading text-lg font-semibold text-foreground">
                {item.title}
              </h3>
              <p className="mt-2 text-sm leading-relaxed text-muted-foreground">
                {item.description}
              </p>
            </SectionReveal>
          </li>
        ))}
      </ol>
    </SectionBase>
  )
}
