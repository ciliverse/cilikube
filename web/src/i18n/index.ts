import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'
import en from './locales/en'
import zh from './locales/zh'

export const LANG_STORAGE_KEY = 'cilikube_lang'

export type AppLang = 'zh' | 'en'

export function getStoredLang(): AppLang {
  try {
    const v = localStorage.getItem(LANG_STORAGE_KEY)
    if (v === 'en' || v === 'zh') return v
  } catch {
    /* ignore */
  }
  // Default English for public demo + new visitors (boot / login)
  return 'en'
}

export function setStoredLang(lang: AppLang) {
  try {
    localStorage.setItem(LANG_STORAGE_KEY, lang)
  } catch {
    /* ignore */
  }
  document.documentElement.lang = lang === 'zh' ? 'zh-CN' : 'en'
}

void i18n.use(initReactI18next).init({
  resources: {
    zh: { translation: zh },
    en: { translation: en },
  },
  lng: getStoredLang(),
  fallbackLng: 'en',
  interpolation: { escapeValue: false },
})

setStoredLang(getStoredLang())

export default i18n
