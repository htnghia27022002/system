'use client'

import Image from 'next/image'

import logoLockup from '@/assets/favicon/logo-lockup.svg'

export function AdminAppLogo() {
  return (
    <Image
      src={logoLockup}
      alt="htnghia"
      width={152}
      height={48}
      className="h-8 w-auto max-w-[9.5rem]"
      priority
    />
  )
}
