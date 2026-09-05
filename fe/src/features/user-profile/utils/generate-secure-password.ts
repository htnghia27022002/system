const UPPER = 'ABCDEFGHJKLMNPQRSTUVWXYZ'
const LOWER = 'abcdefghijkmnopqrstuvwxyz'
const DIGITS = '23456789'
const SYMBOLS = '!@#$%^&*-_=+'
const ALL = UPPER + LOWER + DIGITS + SYMBOLS

function randomIndex(max: number): number {
  const bytes = new Uint8Array(1)
  crypto.getRandomValues(bytes)
  return bytes[0] % max
}

function pick(set: string): string {
  return set[randomIndex(set.length)] ?? set[0]
}

/** Cryptographically random password with mixed character classes. */
export function generateSecurePassword(length = 16): string {
  const size = Math.max(length, 12)
  const chars = Array.from({ length: size }, () => pick(ALL))
  chars[0] = pick(UPPER)
  chars[1] = pick(LOWER)
  chars[2] = pick(DIGITS)
  chars[3] = pick(SYMBOLS)
  for (let i = chars.length - 1; i > 0; i -= 1) {
    const j = randomIndex(i + 1)
    const current = chars[i]
    chars[i] = chars[j] ?? current
    chars[j] = current
  }
  return chars.join('')
}
