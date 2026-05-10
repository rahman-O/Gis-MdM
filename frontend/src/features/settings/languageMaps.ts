/** Form uses short codes from the spec; backend (legacy Angular) uses locale codes. */
const FORM_TO_API: Record<string, string> = {
  en: 'en_US',
  ru: 'ru_RU',
  de: 'de_DE',
  fr: 'fr_FR',
  es: 'es_ES',
  pt: 'pt_PT',
  zh: 'zh_CN',
}

const API_TO_FORM: Record<string, string> = {
  en_US: 'en',
  ru_RU: 'ru',
  de_DE: 'de',
  fr_FR: 'fr',
  es_ES: 'es',
  pt_PT: 'pt',
  zh_CN: 'zh',
}

export const LANGUAGE_OPTIONS: { value: string; label: string }[] = [
  { value: 'en', label: 'English' },
  { value: 'ru', label: 'Russian' },
  { value: 'de', label: 'German' },
  { value: 'fr', label: 'French' },
  { value: 'es', label: 'Spanish' },
  { value: 'pt', label: 'Portuguese' },
  { value: 'zh', label: 'Chinese' },
]

export function formLanguageToApi(code: string): string {
  return FORM_TO_API[code] ?? 'en_US'
}

export function apiLanguageToForm(api: string | null | undefined): string {
  if (!api) return 'en'
  if (API_TO_FORM[api]) return API_TO_FORM[api]
  const short = api.slice(0, 2).toLowerCase()
  if (LANGUAGE_OPTIONS.some((o) => o.value === short)) return short
  return 'en'
}
