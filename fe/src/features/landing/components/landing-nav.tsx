'use client'

import { motion } from 'framer-motion'
import Image from 'next/image'
import Link from 'next/link'
import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'

import { ThemeToggle } from '@/components/common/theme-toggle'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'
import { useAuthStore } from '@/store/auth-store'
import { useSignOut } from '@/features/auth/hooks/use-sign-out'
import logoLockup from '@/assets/favicon/logo-lockup.svg'

import { landingContent, scrollySectionIds } from '../content'
import { useScrolly } from '../context/scrolly-context'

export function LandingNav() {
  const { t, i18n } = useTranslation('common')
  const accessToken = useAuthStore((state) => state.accessToken)
  const signOutMutation = useSignOut()
  const { scrollerRef, scrollToSection } = useScrolly()
  const [activeId, setActiveId] = useState<string>('hero')
  const { nav } = landingContent

  useEffect(() => {
    const scroller = scrollerRef.current
    if (!scroller) return

    const sections = scrollySectionIds
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
      { root: scroller, rootMargin: '-40% 0px -45% 0px', threshold: [0, 0.25, 0.5] },
    )

    sections.forEach((section) => observer.observe(section))
    return () => observer.disconnect()
  }, [scrollerRef])

  const links = [
    { id: 'hero', label: nav.hero },
    { id: 'features', label: nav.features },
    { id: 'tokenomics', label: nav.tokenomics },
    { id: 'cta', label: nav.cta },
  ]

  return (
    <motion.header
      initial={{ y: -16, opacity: 0 }}
      animate={{ y: 0, opacity: 1 }}
      transition={{ duration: 0.5, ease: [0.22, 1, 0.36, 1] }}
      whileHover={{ boxShadow: '0 0 0 1px color-mix(in oklch, var(--primary) 25%, transparent)' }}
      className="fixed inset-x-0 top-0 z-50 border-b border-border/50 bg-background/70 backdrop-blur-xl"
    >
      <div className="mx-auto flex h-16 max-w-7xl items-center justify-between gap-4 px-4 lg:px-8">
        <Link href="/" className="flex items-center gap-2" aria-label="Home">
          <Image
            src={logoLockup}
            alt=""
            width={120}
            height={38}
            className="h-8 w-auto brightness-110"
            priority
          />
        </Link>

        <nav className="hidden items-center gap-1 md:flex" aria-label="Section navigation">
          {links.map((link) => (
            <a
              key={link.id}
              href={`#${link.id}`}
              onClick={(event) => {
                event.preventDefault()
                scrollToSection(link.id)
              }}
              className={cn(
                'rounded-lg px-3 py-2 text-sm font-medium transition-colors',
                activeId === link.id
                  ? 'text-primary'
                  : 'text-muted-foreground hover:text-foreground',
              )}
            >
              {link.label}
            </a>
          ))}
        </nav>

        <div className="flex items-center gap-2">
          {accessToken ? (
            <>
              <Button variant="ghost" size="sm" asChild>
                <Link href="/admin">Admin</Link>
              </Button>
              <Button
                variant="outline"
                size="sm"
                disabled={signOutMutation.isPending}
                onClick={() => signOutMutation.mutate()}
              >
                {t('auth.actions.signOut')}
              </Button>
            </>
          ) : (
            <>
              <Button variant="ghost" size="sm" asChild className="hidden sm:inline-flex">
                <Link href="/login">{t('auth.actions.signIn')}</Link>
              </Button>
              <Button size="sm" asChild>
                <Link href="/register">{t('auth.actions.register')}</Link>
              </Button>
            </>
          )}
          <Button
            variant="outline"
            size="sm"
            className="hidden sm:inline-flex"
            onClick={() => void i18n.changeLanguage(i18n.language === 'en' ? 'vi' : 'en')}
          >
            {i18n.language === 'en' ? 'VI' : 'EN'}
          </Button>
          <ThemeToggle />
        </div>
      </div>
    </motion.header>
  )
}
