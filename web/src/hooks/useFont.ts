import { watchEffect } from "vue"
import { useSettingsStore } from "@/store/modules/settings"
import { useLocaleStore } from "@/store/modules/locale"
import { storeToRefs } from "pinia"

/** 字体 CDN 链接 */
const fontCDNs = {
  'victor-mono': 'https://fonts.googleapis.com/css2?family=Victor+Mono:wght@400&display=swap',
  'maple-mono': 'https://cdn.jsdelivr.net/fontsource/fonts/maple-mono@latest/latin-400-normal.woff2'
}

/** 字体映射 */
const fontMap = {
  // 英文字体 - Maple Mono 作为默认等宽字体
  en: {
    'maple-mono': '"Maple Mono", "Consolas", "Monaco", "Courier New", monospace',
    'victor-mono': '"Victor Mono", "Consolas", "Monaco", "Courier New", monospace'
  },
  // 中文字体 - Maple Mono 作为默认等宽字体
  zh: {
    'maple-mono': '"Maple Mono", "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", "微软雅黑", monospace',
    'kai': '"KaiTi", "楷体", "STKaiti", "Helvetica Neue", Helvetica, Arial, serif'
  }
}

/** 加载字体 */
const loadFont = (fontKey: string) => {
  if (fontCDNs[fontKey as keyof typeof fontCDNs]) {
    const fontUrl = fontCDNs[fontKey as keyof typeof fontCDNs]

    // 检查是否已经加载过该字体
    const existingElement = document.querySelector(`[data-font="${fontKey}"]`)
    if (existingElement) return

    if (fontKey === 'maple-mono') {
      // 对于 woff2 字体，使用 @font-face
      const style = document.createElement('style')
      style.setAttribute('data-font', fontKey)
      style.textContent = `
        @font-face {
          font-family: 'Maple Mono';
          src: url('${fontUrl}') format('woff2');
          font-weight: 400;
          font-style: normal;
          font-display: swap;
        }
      `
      document.head.appendChild(style)
    } else {
      // 对于 Google Fonts（Victor Mono 等），使用 link 标签
      const link = document.createElement('link')
      link.rel = 'stylesheet'
      link.href = fontUrl
      link.setAttribute('data-font', fontKey)
      document.head.appendChild(link)
    }
  }
}

/** 应用字体到全局 */
const applyFont = (font: string, language: 'en' | 'zh') => {
  const body = document.body

  // 智能默认字体选择：中英文都使用 Maple Mono 作为默认
  let actualFont = font
  if (font === 'default') {
    actualFont = 'maple-mono'
  }

  // 确保 Maple Mono 字体总是被加载（中英文都使用）
  if (actualFont === 'maple-mono' || actualFont === 'default') {
    loadFont('maple-mono')
  }

  // 如果需要加载其他外部字体，先加载
  if (fontCDNs[actualFont as keyof typeof fontCDNs]) {
    loadFont(actualFont)
  }

  // 根据语言和字体选择合适的字体族
  let fontFamily = ''

  if (language === 'en') {
    fontFamily = fontMap.en[actualFont as keyof typeof fontMap.en] || fontMap.en['maple-mono']
  } else {
    fontFamily = fontMap.zh[actualFont as keyof typeof fontMap.zh] || fontMap.zh['maple-mono']
  }

  // 添加调试信息
  console.log(`应用字体: 语言=${language}, 原始字体=${font}, 实际字体=${actualFont}, 字体族=${fontFamily}`)

  body.style.fontFamily = fontFamily
}

/** 初始化字体 */
export const initFont = () => {
  const settingsStore = useSettingsStore()
  const localeStore = useLocaleStore()
  const { fontFamily } = storeToRefs(settingsStore)
  const { currentLanguage } = storeToRefs(localeStore)

  // 使用 watchEffect 监听字体和语言变化
  watchEffect(() => {
    // 检查当前字体是否在当前语言的字体选项中存在
    const availableFonts = Object.keys(fontMap[currentLanguage.value])
    
    // 如果当前字体不在可用字体列表中，自动切换到 Maple Mono
    if (!availableFonts.includes(fontFamily.value) && fontFamily.value !== 'default') {
      console.log(`字体 ${fontFamily.value} 不可用，切换到 Maple Mono`)
      settingsStore.fontFamily = 'maple-mono'
      return // 等待下一次 watchEffect 触发
    }
    
    console.log(`🎨 应用字体: ${fontFamily.value} (语言: ${currentLanguage.value})`)
    applyFont(fontFamily.value, currentLanguage.value)
  })
}

/** 获取当前语言可用的字体选项 */
export const getAvailableFonts = (language: 'en' | 'zh') => {
  return Object.keys(fontMap[language])
}

/** 字体 hook */
export function useFont() {
  return { initFont, getAvailableFonts }
}
