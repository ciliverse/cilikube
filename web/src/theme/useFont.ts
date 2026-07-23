import { useSyncExternalStore } from 'react'
import {
  FONT_PACKS,
  getStoredFontId,
  resolveFont,
  setFontId,
  subscribeFont,
  type FontPack,
} from './fonts'

export function useFont(): {
  font: FontPack
  fontId: string
  fonts: FontPack[]
  setFont: (id: string) => void
} {
  const fontId = useSyncExternalStore(subscribeFont, getStoredFontId, () => 'hud')
  return {
    font: resolveFont(fontId),
    fontId,
    fonts: FONT_PACKS,
    setFont: setFontId,
  }
}
