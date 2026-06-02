import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'

import enAdmin from '@/locales/en/admin.json'
import en from '@/locales/en/common.json'
import viAdmin from '@/locales/vi/admin.json'
import vi from '@/locales/vi/common.json'

void i18n.use(initReactI18next).init({
  resources: {
    en: { common: en, admin: enAdmin },
    vi: { common: vi, admin: viAdmin },
  },
  lng: 'en',
  fallbackLng: 'en',
  defaultNS: 'common',
  interpolation: {
    escapeValue: false,
  },
})

export default i18n
