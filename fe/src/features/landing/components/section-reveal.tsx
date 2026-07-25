'use client'

import { motion, type Variants } from 'framer-motion'
import { useEffect, useState, type ReactNode } from 'react'

import { useReducedMotion } from '../hooks/use-reduced-motion'

type SectionRevealProps = {
  children: ReactNode
  className?: string
  delay?: number
}

const variants: Variants = {
  hidden: { opacity: 0, y: 20 },
  visible: { opacity: 1, y: 0 },
}

/** Light enter animation — plain div until mounted to avoid hydration mismatch. */
export function SectionReveal({ children, className, delay = 0 }: SectionRevealProps) {
  const reducedMotion = useReducedMotion()
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    setMounted(true)
  }, [])

  if (!mounted || reducedMotion) {
    return <div className={className}>{children}</div>
  }

  return (
    <motion.div
      className={className}
      initial="hidden"
      whileInView="visible"
      viewport={{ once: true, margin: '-8% 0px' }}
      transition={{ duration: 0.45, ease: [0.22, 1, 0.36, 1], delay }}
      variants={variants}
    >
      {children}
    </motion.div>
  )
}
