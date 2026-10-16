const UNITS: Array<[Intl.RelativeTimeFormatUnit, number]> = [
  ['year', 365 * 24 * 60 * 60],
  ['month', 30 * 24 * 60 * 60],
  ['week', 7 * 24 * 60 * 60],
  ['day', 24 * 60 * 60],
  ['hour', 60 * 60],
  ['minute', 60],
]

export function formatDateTime(iso: string, locale: string): string {
  const date = new Date(iso)
  if (Number.isNaN(date.getTime())) return iso
  return date.toLocaleString(locale, { dateStyle: 'medium', timeStyle: 'short' })
}

export function formatDate(iso: string, locale: string): string {
  const date = new Date(iso)
  if (Number.isNaN(date.getTime())) return iso
  return date.toLocaleDateString(locale, { dateStyle: 'medium' })
}

/**
 * "5 minutes ago" / "in 3 days". Returns null under a minute so callers can
 * show their own "just now" copy.
 */
export function formatRelative(iso: string, locale: string, now = Date.now()): string | null {
  const time = new Date(iso).getTime()
  if (Number.isNaN(time)) return iso
  const diffSeconds = Math.round((time - now) / 1000)
  if (Math.abs(diffSeconds) < 60) return null

  const formatter = new Intl.RelativeTimeFormat(locale, { numeric: 'auto' })
  for (const [unit, seconds] of UNITS) {
    if (Math.abs(diffSeconds) >= seconds) {
      return formatter.format(Math.round(diffSeconds / seconds), unit)
    }
  }
  return null
}
