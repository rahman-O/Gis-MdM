import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'
import en from '@/i18n/locales/en.json'
import ar from '@/i18n/locales/ar.json'

const stored = typeof window !== 'undefined' ? window.localStorage.getItem('hmdm_ui_lang') : null
const lng = stored && ['en', 'ar'].includes(stored) ? stored : 'en'

void i18n.use(initReactI18next).init({
  lng,
  fallbackLng: 'en',
  resources: {
    en: { translation: en },
    ar: { translation: ar },
  },
  interpolation: { escapeValue: false },
})

export function setUiLanguage(code: 'en' | 'ar'): void {
  if (typeof window !== 'undefined') {
    window.localStorage.setItem('hmdm_ui_lang', code)
  }
  void i18n.changeLanguage(code)
}

export default i18n
