import { ref, watch } from 'vue'
import { defineStore } from 'pinia'
import { useI18n } from 'vue-i18n'
import { setLanguage, getLanguage } from '@/utils/cache/local-storage'

export type Language = 'en' | 'zh'

export const useLocaleStore = defineStore('locale', () => {
  const { locale } = useI18n()
  
  // 当前语言
  const currentLanguage = ref<Language>(getLanguage() || 'en')
  
  // 设置语言
  const setLocale = (lang: Language) => {
    currentLanguage.value = lang
    locale.value = lang
    setLanguage(lang)
  }
  
  // 监听语言变化
  watch(currentLanguage, (newLang) => {
    locale.value = newLang
    setLanguage(newLang)
  })
  
  // 初始化语言
  const initLocale = () => {
    const savedLang = getLanguage()
    if (savedLang) {
      setLocale(savedLang as Language)
    }
  }
  
  return {
    currentLanguage,
    setLocale,
    initLocale
  }
})