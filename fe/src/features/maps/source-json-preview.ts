export function valueAtPath(root: unknown, path: string): unknown {
  const trimmed = path.trim()
  if (!trimmed) return root
  let current: unknown = root
  for (const part of trimmed.split('.')) {
    if (current == null || typeof current !== 'object' || Array.isArray(current)) {
      return undefined
    }
    current = (current as Record<string, unknown>)[part]
  }
  return current
}

export function resolveList(root: unknown, listPath: string): unknown[] {
  const value = valueAtPath(root, listPath)
  return Array.isArray(value) ? value : []
}

export function flattenFieldPaths(value: unknown, prefix = '', depth = 0): string[] {
  if (value == null || depth > 3) return prefix ? [prefix] : []
  if (Array.isArray(value)) {
    return prefix ? [prefix] : []
  }
  if (typeof value !== 'object') {
    return prefix ? [prefix] : []
  }
  const entries = Object.entries(value as Record<string, unknown>)
  if (entries.length === 0) {
    return prefix ? [prefix] : []
  }
  const paths: string[] = []
  for (const [key, child] of entries) {
    const next = prefix ? `${prefix}.${key}` : key
    if (child != null && typeof child === 'object' && !Array.isArray(child)) {
      paths.push(...flattenFieldPaths(child, next, depth + 1))
    } else {
      paths.push(next)
    }
  }
  return paths
}

export function suggestListPaths(root: unknown): string[] {
  if (Array.isArray(root)) return ['']
  const found: string[] = []
  walkArrays(root, '', found, 0)
  return found
}

function walkArrays(
  value: unknown,
  prefix: string,
  found: string[],
  depth: number,
) {
  if (value == null || depth > 4) return
  if (Array.isArray(value)) {
    found.push(prefix)
    return
  }
  if (typeof value !== 'object') return
  for (const [key, child] of Object.entries(value as Record<string, unknown>)) {
    const next = prefix ? `${prefix}.${key}` : key
    walkArrays(child, next, found, depth + 1)
  }
}

const AUTO_ALIASES: Record<string, string[]> = {
  'location.name': ['locationname', 'location_name', 'sitename', 'site'],
  'location.locationKey': ['locationkey', 'location_key', 'locationid'],
  'location.countryCode': ['countrycode', 'country_code', 'country'],
  'location.adminCode': ['admincode', 'admin_code'],
  'location.street': ['street', 'address'],
  'location.postalCode': ['postalcode', 'postal_code', 'zip'],
  'location.formatted': ['formatted', 'address', 'formattedaddress'],
  'location.lat': ['lat', 'latitude'],
  'location.lng': ['lng', 'lon', 'longitude'],
  'place.name': ['placename', 'place_name', 'name', 'title'],
  'place.category': ['category', 'type'],
  'place.placeKey': ['placekey', 'place_key', 'id', 'placeid'],
  'place.unit': ['unit'],
  'place.lat': ['lat', 'latitude'],
  'place.lng': ['lng', 'lon', 'longitude'],
  'news.title': ['title', 'headline'],
  'news.originalUrl': ['originalurl', 'original_url', 'url', 'link'],
  'details.description': ['description', 'desc'],
  'details.phone': ['phone', 'tel'],
  'details.hours': ['hours', 'openinghours'],
  'details.website': ['website'],
  'details.priceRange': ['pricerange', 'price_range', 'price'],
}

function normalizeToken(path: string) {
  return path.replace(/[_\-\s]/g, '').toLowerCase()
}

export function autoMapFields(fieldPaths: string[]): Record<string, string> {
  const tokens = fieldPaths.map((path) => ({
    path,
    token: normalizeToken(path.split('.').at(-1) ?? path),
    full: normalizeToken(path),
  }))
  const next: Record<string, string> = {}
  for (const [target, aliases] of Object.entries(AUTO_ALIASES)) {
    const match = tokens.find(
      (item) => aliases.includes(item.token) || aliases.includes(item.full),
    )
    if (match) next[target] = match.path
  }
  return next
}
