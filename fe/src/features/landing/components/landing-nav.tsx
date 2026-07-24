'use client'

import { useEffect, useState } from 'react'

import { PublicSiteHeader } from '@/components/common/public-site-header'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'

import { landingContent, landingSectionIds } from '../content'

function scrollToSection(id: string) {
  document.getElementById(id)?.scrollIntoView({ behavior: 'smooth', block: 'start' })
}

/**
 * Landing chrome: shared PublicSiteHeader + in-page section jump links.
 */
export function LandingNav() {
  const [activeId, setActiveId] = useState<string>('hero')
  const { nav } = landingContent

  useEffect(() => {
    const sections = landingSectionIds
      .map((id) => document.getElementById(id))
      .filter((el): el is HTMLElement => el !== null)

    const observer = new IntersectionObserver(
      (entries) => {
        const visible = entries
          .filter((entry) => entry.isIntersecting)
          .sort((a, b) => b.intersectionRatio - a.intersectionRatio)[0]
        if (visible?.target.id) {
          setActiveId(visible.target.id)
        }
      },
      { root: null, rootMargin: '-35% 0px -50% 0px', threshold: [0, 0.25, 0.5] },
    )

    sections.forEach((section) => observer.observe(section))
    return () => observer.disconnect()
  }, [])

  const sectionLinks = [
    { id: 'features', label: nav.features },
    { id: 'tools', label: nav.tools },
    { id: 'cta', label: nav.cta },
  ]

  return (
    <PublicSiteHeader variant="landing">
      {sectionLinks.map((link) => (
        <Button
          key={link.id}
          variant="ghost"
          size="sm"
          className={cn(
            'justify-start md:justify-center',
            activeId === link.id ? 'text-primary' : 'text-muted-foreground',
          )}
          onClick={() => scrollToSection(link.id)}
        >
          {link.label}
        </Button>
      ))}
    </PublicSiteHeader>
  )
}
