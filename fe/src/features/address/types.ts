export type Country = {
  id: string
  code: string
  code3: string
  name: string
  nameLocal?: string
}

export type AdminDivision = {
  id: string
  countryCode: string
  parentId?: string | null
  level: number
  code: string
  name: string
  nameEn?: string
  fullName?: string
  type: string
  path: string
  lat?: number | null
  lng?: number | null
}

export type CatalogList<T> = {
  items: T[]
}

export type AddressValue = {
  countryCode: string
  adminDivisionId: string
}
