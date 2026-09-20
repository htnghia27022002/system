'use client'

import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Spinner } from '@/components/ui/spinner'

import { P1_PLACE_CATEGORIES } from '../category-filters'
import type { PlaceCategoryKey, PlaceWriteInput } from '../types'

type PlaceQuickCreateProps = {
  locationId: string
  locationName: string
  lat?: number | null
  lng?: number | null
  isPending: boolean
  onSubmit: (input: PlaceWriteInput) => void
  onCancel: () => void
}

export function PlaceQuickCreate({
  locationId,
  locationName,
  lat,
  lng,
  isPending,
  onSubmit,
  onCancel,
}: PlaceQuickCreateProps) {
  const { t } = useTranslation('admin')
  const [name, setName] = useState('')
  const [category, setCategory] = useState<PlaceCategoryKey>('eatery')
  const hasCoords = lat != null && lng != null

  useEffect(() => {
    setName('')
    setCategory('eatery')
  }, [locationId])

  return (
    <form
      className="grid gap-2 rounded-xl border border-border bg-muted/40 p-3"
      onSubmit={(event) => {
        event.preventDefault()
        const trimmed = name.trim()
        if (!trimmed) return
        onSubmit({
          locationId,
          name: trimmed,
          category,
          status: hasCoords ? 'active' : 'pending',
          lat: lat ?? undefined,
          lng: lng ?? undefined,
        })
      }}
    >
      <p className="text-xs text-muted-foreground">
        {t('maps.create.parentLocation', { name: locationName })}
      </p>
      <Input
        value={name}
        onChange={(event) => setName(event.target.value)}
        placeholder={t('maps.places.quickName')}
        aria-label={t('maps.places.fields.name')}
        required
        autoFocus
      />
      <Select
        value={category}
        onValueChange={(value) => setCategory(value as PlaceCategoryKey)}
      >
        <SelectTrigger aria-label={t('maps.create.category')}>
          <SelectValue />
        </SelectTrigger>
        <SelectContent className="z-[1400]">
          {P1_PLACE_CATEGORIES.map((key) => (
            <SelectItem key={key} value={key}>
              {t(`maps.categories.${key}`)}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
      <div className="flex justify-end gap-2">
        <Button type="button" variant="outline" size="sm" onClick={onCancel}>
          {t('access.actions.cancel')}
        </Button>
        <Button type="submit" size="sm" disabled={isPending || !name.trim()}>
          {isPending ? <Spinner className="size-4" /> : null}
          {t('maps.places.createAction')}
        </Button>
      </div>
    </form>
  )
}
