'use client'

import { useLayoutEffect, useRef } from 'react'
import gsap from 'gsap'
import { ScrollTrigger } from 'gsap/ScrollTrigger'

gsap.registerPlugin(ScrollTrigger)

type UseGsapContextOptions = {
  scrollerRef: React.RefObject<HTMLElement | null>
  enabled?: boolean
}

export function useGsapContext({ scrollerRef, enabled = true }: UseGsapContextOptions) {
  const ctxRef = useRef<gsap.Context | null>(null)

  useLayoutEffect(() => {
    if (!enabled) return

    const scroller = scrollerRef.current
    if (!scroller) return

    ctxRef.current = gsap.context(() => {
      ScrollTrigger.defaults({ scroller })
    }, scroller)

    const refresh = () => ScrollTrigger.refresh()
    refresh()

    window.addEventListener('resize', refresh)
    return () => {
      window.removeEventListener('resize', refresh)
      ctxRef.current?.revert()
      ctxRef.current = null
    }
  }, [scrollerRef, enabled])
}
