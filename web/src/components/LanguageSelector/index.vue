<template>
  <el-dropdown trigger="click" @command="handleCommand">
    <el-button text>
      <el-icon><Operation /></el-icon>
      {{ currentLanguageLabel }}
    </el-button>
    <template #dropdown>
      <el-dropdown-menu>
        <!-- Language Section -->
        <div class="dropdown-section">
          <div class="section-title">{{ $t('settings.language') }}</div>
          <el-dropdown-item command="lang:en" :class="{ 'is-active': currentLanguage === 'en' }">
            {{ $t('language.english') }}
          </el-dropdown-item>
          <el-dropdown-item command="lang:zh" :class="{ 'is-active': currentLanguage === 'zh' }">
            {{ $t('language.chinese') }}
          </el-dropdown-item>
        </div>
        
        <!-- Divider -->
        <el-divider class="dropdown-divider" />
        
        <!-- Font Section -->
        <div class="dropdown-section">
          <div class="section-title">{{ $t('settings.font') }}</div>
          <el-dropdown-item 
            v-for="font in availableFonts" 
            :key="font"
            :command="`font:${font}`"
            :class="{ 'is-active': currentFont === font }"
          >
            {{ $t(`font.${font}`) }}
          </el-dropdown-item>
        </div>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useLocaleStore, type Language } from '@/store/modules/locale'
import { useSettingsStore } from '@/store/modules/settings'
import { getAvailableFonts } from '@/hooks/useFont'
import { Operation } from '@element-plus/icons-vue'

const { t } = useI18n()
const localeStore = useLocaleStore()
const settingsStore = useSettingsStore()

const currentLanguage = computed(() => localeStore.currentLanguage)
const currentFont = computed(() => settingsStore.fontFamily)

const availableFonts = computed(() => {
  return getAvailableFonts(currentLanguage.value)
})

const currentLanguageLabel = computed(() => {
  return currentLanguage.value === 'en' ? t('language.english') : t('language.chinese')
})

const handleCommand = (command: string) => {
  if (command.startsWith('lang:')) {
    const language = command.replace('lang:', '') as Language
    localeStore.setLocale(language)
  } else if (command.startsWith('font:')) {
    const font = command.replace('font:', '')
    settingsStore.fontFamily = font
  }
}
</script>

<style scoped>
.is-active {
  color: var(--el-color-primary);
  font-weight: bold;
}

.dropdown-section {
  padding: 4px 0;
}

.section-title {
  padding: 8px 16px 4px;
  font-size: 12px;
  font-weight: 600;
  color: var(--el-text-color-regular);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.dropdown-divider {
  margin: 4px 0;
}

:deep(.el-dropdown-menu) {
  min-width: 160px;
}

:deep(.el-dropdown-menu__item) {
  padding: 8px 16px;
  font-size: 14px;
}

:deep(.el-dropdown-menu__item:hover) {
  background-color: var(--el-fill-color-light);
}
</style>