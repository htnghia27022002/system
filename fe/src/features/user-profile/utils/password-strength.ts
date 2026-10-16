export type PasswordChecks = {
  length: boolean
  mixedCase: boolean
  number: boolean
  symbol: boolean
}

export type PasswordStrength = {
  /** 0 = empty, 1 = weak, 2 = fair, 3 = good, 4 = strong */
  score: 0 | 1 | 2 | 3 | 4
  checks: PasswordChecks
}

export const PASSWORD_MIN_LENGTH = 8

/** Client-side hint only; the server enforces the real policy. */
export function evaluatePassword(password: string): PasswordStrength {
  const checks: PasswordChecks = {
    length: password.length >= PASSWORD_MIN_LENGTH,
    mixedCase: /[a-z]/.test(password) && /[A-Z]/.test(password),
    number: /\d/.test(password),
    symbol: /[^A-Za-z0-9]/.test(password),
  }

  if (password.length === 0) return { score: 0, checks }
  if (!checks.length) return { score: 1, checks }

  const variety = [checks.mixedCase, checks.number, checks.symbol].filter(Boolean).length
  const lengthBonus = password.length >= 12 ? 1 : 0
  const score = Math.min(4, Math.max(1, variety + lengthBonus)) as PasswordStrength['score']
  return { score, checks }
}
