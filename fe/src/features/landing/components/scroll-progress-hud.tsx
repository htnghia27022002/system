'use client'

import { useEffect, useState } from 'react'

import { scrollySectionIds } from '../content'
import { useScrolly } from '../context/scrolly-context'

export function ScrollProgressHud() {
  const { scrollerRef } = useScrolly()
  const [progress, setProgress] = useState<Record<string, string>>({})

  useEffect(() => {
    if (process.env.NODE_ENV !== 'development') return

    const scroller = scrollerRef.current
    if (!scroller) return

    const readProgress = () => {
      const next: Record<string, string> = {}
      scrollySectionIds.forEach((id) => {
        const el = document.getElementById(id)
        if (el?.dataset.scrollProgress) {
          next[id] = el.dataset.scrollProgress
        }
      })
      setProgress(next)
    }

    const observer = new MutationObserver(readProgress)
    scrollySectionIds.forEach((id) => {
      const el = document.getElementById(id)
      if (el) observer.observe(el, { attributes: true, attributeFilter: ['data-scroll-progress'] })
    })

    readProgress()
    return () => observer.disconnect()
  }, [scrollerRef])

  if (process.env.NODE_ENV !== 'development') return null

  return (
    <div
      data-scroll-progress-hud
      className="pointer-events-none fixed bottom-4 left-4 z-[100] rounded-lg border border-border/60 bg-background/90 px-3 py-2 font-mono text-[10px] text-muted-foreground backdrop-blur-sm"
      aria-hidden
    >
      <div className="mb-1 font-semibold text-primary">Scroll progress</div>
      {scrollySectionIds.map((id) => (
        <div key={id}>
          {id}: {progress[id] ?? '—'}
        </div>
      ))}
    </div>
  )
}
