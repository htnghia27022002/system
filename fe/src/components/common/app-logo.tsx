import Image from 'next/image'
import Link from 'next/link'

import logoIcon from '@/assets/favicon/logo-icon.svg'
import { cn } from '@/lib/utils'

type AppLogoProps = {
  className?: string
}

export function AppLogo({ className }: AppLogoProps) {
  return (
    <Link
      href="/"
      className={cn('inline-flex size-10 items-center justify-center', className)}
      aria-label="Home"
    >
      <Image
        src={logoIcon}
        alt=""
        width={40}
        height={36}
        className="h-9 w-auto"
        priority
      />
    </Link>
  )
}
