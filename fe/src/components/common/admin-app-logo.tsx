import { useTranslation } from 'react-i18next'

export function AdminAppLogo() {
  const { t } = useTranslation('common')

  return (
    <>
      <div className="flex aspect-square size-8 items-center justify-center rounded-md bg-primary text-sm font-semibold text-primary-foreground">
        {t('appName').charAt(0).toUpperCase()}
      </div>
      <div className="ml-1 grid flex-1 text-left text-sm">
        <span className="mb-0.5 truncate leading-tight font-semibold">
          {t('appName')}
        </span>
      </div>
    </>
  )
}
