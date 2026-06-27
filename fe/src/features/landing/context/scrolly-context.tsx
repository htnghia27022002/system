'use client'

import {
  createContext,
  useCallback,
  useContext,
  useMemo,
  useRef,
  type ReactNode,
  type RefObject,
} from 'react'

type ScrollyContextValue = {
  scrollerRef: RefObject<HTMLDivElement | null>
  scrollToSection: (id: string) => void
}

const ScrollyContext = createContext<ScrollyContextValue | null>(null)

type ScrollyProviderProps = {
  children: ReactNode
}

export function ScrollyProvider({ children }: ScrollyProviderProps) {
  const scrollerRef = useRef<HTMLDivElement | null>(null)

  const scrollToSection = useCallback((id: string) => {
    const scroller = scrollerRef.current
    const target = document.getElementById(id)
    if (!scroller || !target) return

    const scrollerTop = scroller.getBoundingClientRect().top
    const targetTop = target.getBoundingClientRect().top
    const nextScrollTop = scroller.scrollTop + (targetTop - scrollerTop)

    scroller.scrollTo({ top: nextScrollTop, behavior: 'auto' })
  }, [])

  const value = useMemo(
    () => ({ scrollerRef, scrollToSection }),
    [scrollToSection],
  )

  return <ScrollyContext.Provider value={value}>{children}</ScrollyContext.Provider>
}

export function useScrolly(): ScrollyContextValue {
  const ctx = useContext(ScrollyContext)
  if (!ctx) {
    throw new Error('useScrolly must be used within ScrollyProvider')
  }
  return ctx
}
