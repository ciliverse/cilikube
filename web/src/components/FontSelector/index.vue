<template>
  <el-dropdown trigger="click" @command="handleFontChange">
    <el-button text>
      <el-icon><Document /></el-icon>
      {{ currentFontLabel }}
    </el-button>
    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item 
          v-for="font in availableFonts" 
          :key="font"
          :command="font"
          :class="{ 'is-active': currentFont === font }"
        >
          {{ $t(`font.${font}`) }}
        </el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useSettingsStore } from '@/store/modules/settings'
import { useLocaleStore } from '@/store/modules/locale'
import { getAvailableFonts } from '@/hooks/useFont'
import { Document } from '@element-plus/icons-vue'

const { t } = useI18n()
const settingsStore = useSettingsStore()
const localeStore = useLocaleStore()

const currentFont = computed(() => settingsStore.fontFamily)
const currentLanguage = computed(() => localeStore.currentLanguage)

const availableFonts = computed(() => {
  return getAvailableFonts(currentLanguage.value)
})

const currentFontLabel = computed(() => {
  // 直接显示当前字体的标签，字体切换逻辑由 useFont 处理
  return t(`font.${currentFont.value}`)
})

const handleFontChange = (font: string) => {
  settingsStore.fontFamily = font
}
</script>

<style scoped>
.is-active {
  color: var(--el-color-primary);
  font-weight: bold;
}
</style>