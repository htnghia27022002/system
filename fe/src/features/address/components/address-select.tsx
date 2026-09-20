'use client'

import { useTranslation } from 'react-i18next'

import { Field, FieldLabel } from '@/components/ui/field'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { cn } from '@/lib/utils'

import {
  useAddressCountries,
  useAddressDivisions,
} from '../hooks/use-address-catalogs'
import type { AddressValue } from '../types'

const EMPTY = 'empty'

type AddressSelectProps = {
  value: AddressValue
  onChange: (next: AddressValue) => void
  allowEmpty?: boolean
  disabled?: boolean
  countryLabel?: string
  divisionLabel?: string
  emptyCountryLabel?: string
  emptyDivisionLabel?: string
  contentClassName?: string
  className?: string
}

export function AddressSelect({
  value,
  onChange,
  allowEmpty = false,
  disabled = false,
  countryLabel,
  divisionLabel,
  emptyCountryLabel,
  emptyDivisionLabel,
  contentClassName,
  className,
}: AddressSelectProps) {
  const { t } = useTranslation()
  const countriesQuery = useAddressCountries()
  const divisionsQuery = useAddressDivisions(value.countryCode || undefined)
  const countries = countriesQuery.data?.items ?? []
  const divisions = divisionsQuery.data?.items ?? []
  const emptyCountry = emptyCountryLabel ?? t('address.allCountries')
  const emptyDivision = emptyDivisionLabel ?? t('address.allDivisions')

  return (
    <div className={cn('grid gap-3 sm:grid-cols-1', className)}>
      <Field>
        <FieldLabel>{countryLabel ?? t('address.country')}</FieldLabel>
        <Select
          value={value.countryCode || (allowEmpty ? EMPTY : '')}
          disabled={disabled}
          onValueChange={(countryCode) =>
            onChange({
              countryCode: countryCode === EMPTY ? '' : countryCode,
              adminDivisionId: '',
            })
          }
        >
          <SelectTrigger aria-label={countryLabel ?? t('address.country')}>
            <SelectValue placeholder={emptyCountry} />
          </SelectTrigger>
          <SelectContent className={contentClassName}>
            {allowEmpty ? (
              <SelectItem value={EMPTY}>{emptyCountry}</SelectItem>
            ) : null}
            {countries.map((country) => (
              <SelectItem key={country.id} value={country.code}>
                {country.nameLocal || country.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </Field>
      <Field>
        <FieldLabel>{divisionLabel ?? t('address.division')}</FieldLabel>
        <Select
          value={value.adminDivisionId || (allowEmpty ? EMPTY : '')}
          disabled={disabled || !value.countryCode}
          onValueChange={(adminDivisionId) =>
            onChange({
              ...value,
              adminDivisionId: adminDivisionId === EMPTY ? '' : adminDivisionId,
            })
          }
        >
          <SelectTrigger aria-label={divisionLabel ?? t('address.division')}>
            <SelectValue placeholder={emptyDivision} />
          </SelectTrigger>
          <SelectContent className={contentClassName}>
            {allowEmpty ? (
              <SelectItem value={EMPTY}>{emptyDivision}</SelectItem>
            ) : null}
            {divisions.map((division) => (
              <SelectItem key={division.id} value={division.id}>
                {division.fullName || division.name}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </Field>
    </div>
  )
}
