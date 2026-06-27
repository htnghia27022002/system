'use client'

import { useEffect, useRef } from 'react'
import gsap from 'gsap'
import { ScrollTrigger } from 'gsap/ScrollTrigger'

import { landingContent } from '../../content'
import { useScrolly } from '../../context/scrolly-context'
import { useReducedMotion } from '../../hooks/use-reduced-motion'
import { devScrollLog } from '../../lib/dev-scroll-log'
import { SplitWords } from '../split-words'
import { Web3Artifact } from '../web3-artifact'
import { SectionSnapBase } from './section-snap-base'

gsap.registerPlugin(ScrollTrigger)

export function FeaturesSnapSection() {
  const sectionRef = useRef<HTMLElement>(null)
  const artifactsRef = useRef<HTMLDivElement>(null)
  const { scrollerRef } = useScrolly()
  const reducedMotion = useReducedMotion()
  const { features } = landingContent

  useEffect(() => {
    if (reducedMotion) return

    const section = sectionRef.current
    const scroller = scrollerRef.current
    const artifacts = artifactsRef.current
    if (!section || !scroller || !artifacts) return

    const artifactA = artifacts.querySelector('[data-artifact-id="a"]')
    const artifactB = artifacts.querySelector('[data-artifact-id="b"]')
    const artifactC = artifacts.querySelector('[data-artifact-id="c"]')
    if (!artifactA || !artifactB || !artifactC) return

    const mm = gsap.matchMedia()
    const cleanupFns: Array<() => void> = []

    mm.add('(min-width: 769px)', () => {
      const pinTrigger = ScrollTrigger.create({
        trigger: section,
        scroller,
        start: 'top top',
        end: '+=150%',
        pin: true,
        pinSpacing: true,
        anticipatePin: 1,
        invalidateOnRefresh: true,
      })
      cleanupFns.push(() => pinTrigger.kill())

      const tl = gsap.timeline({
        scrollTrigger: {
          trigger: section,
          scroller,
          start: 'top top',
          end: '+=150%',
          scrub: 1,
          invalidateOnRefresh: true,
          onUpdate: (self) => {
            if (process.env.NODE_ENV === 'development') {
              section.dataset.scrollProgress = self.progress.toFixed(3)
            }
            devScrollLog('features', self.progress)
          },
        },
      })

      tl.fromTo(artifactA, { rotation: 0 }, { rotation: 360, ease: 'none' }, 0)
        .fromTo(artifactB, { scale: 0, opacity: 0 }, { scale: 1, opacity: 1, ease: 'none' }, 0)
        .fromTo(artifactC, { x: '-40vw', opacity: 0 }, { x: '0vw', opacity: 1, ease: 'none' }, 0)

      if (tl.scrollTrigger) cleanupFns.push(() => tl.scrollTrigger?.kill())

      const wordTween = gsap.fromTo(
        section.querySelectorAll('.split-word'),
        { y: 40, opacity: 0 },
        {
          y: 0,
          opacity: 1,
          stagger: 0.04,
          ease: 'power2.out',
          scrollTrigger: {
            trigger: section,
            scroller,
            start: 'top 80%',
            end: 'top 40%',
            scrub: 1,
          },
        },
      )
      if (wordTween.scrollTrigger) cleanupFns.push(() => wordTween.scrollTrigger?.kill())
    })

    return () => {
      cleanupFns.forEach((fn) => fn())
      mm.revert()
    }
  }, [reducedMotion, scrollerRef])

  return (
    <SectionSnapBase ref={sectionRef} id="features" className="landing-section-features">
      <p className="text-xs font-semibold uppercase tracking-[0.25em] text-primary">
        {features.label}
      </p>
      <SplitWords
        as="h2"
        text={features.title}
        className="mt-4 font-heading text-3xl font-semibold tracking-tight sm:text-4xl lg:text-5xl"
      />
      <p className="mt-4 max-w-2xl text-muted-foreground">{features.subtitle}</p>

      <div
        ref={artifactsRef}
        className="relative mt-12 grid min-h-[280px] place-items-center md:min-h-[320px]"
      >
        <Web3Artifact
          artifactId="a"
          variant="hex"
          className="absolute left-[8%] top-[10%] md:left-[12%]"
        />
        <Web3Artifact
          artifactId="b"
          variant="diamond"
          className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2"
        />
        <Web3Artifact
          artifactId="c"
          variant="ring"
          className="absolute right-[8%] bottom-[10%] md:right-[12%]"
        />
      </div>

      <ul className="mt-10 grid gap-4 sm:grid-cols-3">
        {features.items.map((item) => (
          <li
            key={item.title}
            className="rounded-xl border border-border/60 bg-card/40 p-4 backdrop-blur-sm"
          >
            <h3 className="font-medium text-foreground">{item.title}</h3>
            <p className="mt-2 text-sm text-muted-foreground">{item.description}</p>
          </li>
        ))}
      </ul>
    </SectionSnapBase>
  )
}
