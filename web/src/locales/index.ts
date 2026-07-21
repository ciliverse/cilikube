import { createI18n } from 'vue-i18n'
import { getLanguage } from '@/utils/cache/local-storage'
import en from './lang/en'
import zh from './lang/zh'

const messages = {
  en,
  zh
}

const i18n = createI18n({
  legacy: false,
  locale: getLanguage() || 'en',
  fallbackLocale: 'en',
  messages,
  globalInjection: true
})

export default i18n