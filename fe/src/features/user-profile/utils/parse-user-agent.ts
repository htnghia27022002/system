export type DeviceType = 'desktop' | 'mobile' | 'tablet'

export type ParsedUserAgent = {
  browser: string | null
  os: string | null
  deviceType: DeviceType
}

const BROWSERS: Array<[RegExp, string]> = [
  [/Edg(e|A|iOS)?\//, 'Edge'],
  [/(OPR|Opera)\//, 'Opera'],
  [/SamsungBrowser\//, 'Samsung Internet'],
  [/(Firefox|FxiOS)\//, 'Firefox'],
  [/(Chrome|CriOS)\//, 'Chrome'],
  [/Version\/[\d.]+.*Safari\//, 'Safari'],
]

const OPERATING_SYSTEMS: Array<[RegExp, string]> = [
  [/iPad/, 'iPadOS'],
  [/(iPhone|iPod)/, 'iOS'],
  [/Android/, 'Android'],
  [/CrOS/, 'ChromeOS'],
  [/Windows/, 'Windows'],
  [/Mac OS X|Macintosh/, 'macOS'],
  [/Linux/, 'Linux'],
]

/** Reduces a raw User-Agent header to a readable browser / OS / device type. */
export function parseUserAgent(userAgent: string | null | undefined): ParsedUserAgent {
  const ua = userAgent?.trim() ?? ''
  const browser = BROWSERS.find(([pattern]) => pattern.test(ua))?.[1] ?? null
  const os = OPERATING_SYSTEMS.find(([pattern]) => pattern.test(ua))?.[1] ?? null

  let deviceType: DeviceType = 'desktop'
  if (/iPad|Tablet/i.test(ua) || (/Android/.test(ua) && !/Mobi/.test(ua))) {
    deviceType = 'tablet'
  } else if (/Mobi|iPhone|iPod/.test(ua)) {
    deviceType = 'mobile'
  }

  return { browser, os, deviceType }
}
