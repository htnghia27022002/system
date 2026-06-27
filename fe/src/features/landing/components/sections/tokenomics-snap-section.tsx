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

export function TokenomicsSnapSection() {
  const sectionRef = useRef<HTMLElement>(null)
  const stackRef = useRef<HTMLDivElement>(null)
  const { scrollerRef } = useScrolly()
  const reducedMotion = useReducedMotion()
  const { tokenomics } = landingContent

  useEffect(() => {
    if (reducedMotion) return

    const section = sectionRef.current
    const scroller = scrollerRef.current
    const stack = stackRef.current
    if (!section || !scroller || !stack) return

    const cards = gsap.utils.toArray<HTMLElement>('[data-tokenomics-card]', stack)
    if (cards.length === 0) return

    const mm = gsap.matchMedia()
    const cleanupFns: Array<() => void> = []

    mm.add('(min-width: 769px)', () => {
      const pinTrigger = ScrollTrigger.create({
        trigger: section,
        scroller,
        start: 'top top',
        end: `+=${120 * cards.length}%`,
        pin: true,
        pinSpacing: true,
        anticipatePin: 1,
        invalidateOnRefresh: true,
      })
      cleanupFns.push(() => pinTrigger.kill())

      cards.forEach((card, index) => {
        gsap.set(card, { zIndex: index + 1 })
      })

      const tl = gsap.timeline({
        scrollTrigger: {
          trigger: section,
          scroller,
          start: 'top top',
          end: `+=${120 * cards.length}%`,
          scrub: 1,
          invalidateOnRefresh: true,
          onUpdate: (self) => {
            if (process.env.NODE_ENV === 'development') {
              section.dataset.scrollProgress = self.progress.toFixed(3)
            }
            devScrollLog('tokenomics', self.progress)
          },
        },
      })

      cards.forEach((card, index) => {
        if (index === 0) {
          tl.fromTo(card, { y: 0, scale: 1 }, { y: 0, scale: 1, ease: 'none' }, 0)
          return
        }

        const segmentStart = index / cards.length
        tl.fromTo(
          card,
          { y: 120, scale: 0.92, opacity: 0.6 },
          { y: index * -12, scale: 1 - index * 0.03, opacity: 1, ease: 'none' },
          segmentStart,
        )
      })

      if (tl.scrollTrigger) cleanupFns.push(() => tl.scrollTrigger?.kill())

      const artifacts = section.querySelectorAll('.web3-artifact')
      const artifactTl = gsap.timeline({
        scrollTrigger: {
          trigger: section,
          scroller,
          start: 'top top',
          end: `+=${120 * cards.length}%`,
          scrub: 1,
        },
      })

      artifactTl
        .fromTo(
          artifacts[0],
          { rotation: 0, x: '-20vw' },
          { rotation: 180, x: '0vw', ease: 'none' },
          0,
        )
        .fromTo(
          artifacts[1],
          { scale: 0.5, opacity: 0 },
          { scale: 1, opacity: 1, ease: 'none' },
          0.2,
        )
        .fromTo(
          artifacts[2],
          { x: '30vw', opacity: 0 },
          { x: '0vw', opacity: 0.85, ease: 'none' },
          0.4,
        )

      if (artifactTl.scrollTrigger) cleanupFns.push(() => artifactTl.scrollTrigger?.kill())
    })

    return () => {
      cleanupFns.forEach((fn) => fn())
      mm.revert()
    }
  }, [reducedMotion, scrollerRef])

  return (
    <SectionSnapBase ref={sectionRef} id="tokenomics" className="landing-section-tokenomics">
      <p className="text-xs font-semibold uppercase tracking-[0.25em] text-primary">
        {tokenomics.label}
      </p>
      <SplitWords
        as="h2"
        text={tokenomics.title}
        className="mt-4 font-heading text-3xl font-semibold tracking-tight sm:text-4xl lg:text-5xl"
      />
      <p className="mt-4 max-w-2xl text-muted-foreground">{tokenomics.subtitle}</p>

      <div className="relative mt-8 hidden min-h-[80px] md:block">
        <Web3Artifact
          artifactId="tok-a"
          variant="ring"
          className="absolute left-0 top-0 scale-75 opacity-80"
        />
        <Web3Artifact
          artifactId="tok-b"
          variant="hex"
          className="absolute left-1/2 top-0 -translate-x-1/2 scale-75"
        />
        <Web3Artifact
          artifactId="tok-c"
          variant="diamond"
          className="absolute right-0 top-0 scale-75 opacity-80"
        />
      </div>

      <div
        ref={stackRef}
        className="relative mx-auto mt-10 flex w-full max-w-md flex-col gap-4 md:mt-14 md:min-h-[420px] md:max-w-lg"
      >
        {tokenomics.cards.map((card, index) => (
          <article
            key={card.title}
            data-tokenomics-card
            className="tokenomics-card w-full rounded-2xl border border-primary/30 bg-card/70 p-6 shadow-[0_0_32px_-12px_var(--landing-glow-primary)] backdrop-blur-md will-change-transform md:absolute md:inset-x-0 md:top-0"
            style={reducedMotion ? undefined : { zIndex: index + 1 }}
          >
            <div className="flex items-baseline justify-between gap-4">
              <h3 className="font-heading text-xl font-semibold text-foreground">{card.title}</h3>
              <span className="text-2xl font-semibold text-primary">{card.share}</span>
            </div>
            <p className="mt-3 text-sm leading-relaxed text-muted-foreground">{card.description}</p>
          </article>
        ))}
      </div>
    </SectionSnapBase>
  )
}
