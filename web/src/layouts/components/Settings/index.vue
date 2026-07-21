<script lang="ts" setup>
import { watchEffect } from "vue"
import { useI18n } from "vue-i18n"
import { storeToRefs } from "pinia"
import { useSettingsStore } from "@/store/modules/settings"
import { useLayoutMode } from "@/hooks/useLayoutMode"
import { resetConfigLayout } from "@/utils"
import SelectLayoutMode from "./SelectLayoutMode.vue"
import { Refresh } from "@element-plus/icons-vue"

const { t } = useI18n()

const { isLeft } = useLayoutMode()
const settingsStore = useSettingsStore()

/** 使用 storeToRefs 将提取的属性保持其响应性 */
const {
  showTagsView,
  showLogo,
  fixedHeader,
  showFooter,
  showNotify,
  showThemeSwitch,
  showScreenfull,
  showSearchMenu,
  showFontSelector,
  cacheTagsView,
  showWatermark,
  showGreyMode,
  showColorWeakness
} = storeToRefs(settingsStore)

/** 定义 switch 设置项 */
const switchSettings = {
  [t('settings.layout.showTagsView')]: showTagsView,
  [t('settings.layout.showLogo')]: showLogo,
  [t('settings.layout.fixedHeader')]: fixedHeader,
  [t('settings.layout.showFooter')]: showFooter,
  [t('settings.layout.showNotify')]: showNotify,
  [t('settings.layout.showThemeSwitch')]: showThemeSwitch,
  [t('settings.layout.showScreenfull')]: showScreenfull,
  [t('settings.layout.showSearchMenu')]: showSearchMenu,
  [t('settings.layout.showFontSelector')]: showFontSelector,
  [t('settings.layout.cacheTagsView')]: cacheTagsView,
  [t('settings.layout.showWatermark')]: showWatermark,
  [t('settings.layout.showGreyMode')]: showGreyMode,
  [t('settings.layout.showColorWeakness')]: showColorWeakness
}

/** 非左侧模式时，Header 都是 fixed 布局 */
watchEffect(() => {
  !isLeft.value && (fixedHeader.value = true)
})
</script>

<template>
  <div class="setting-container">
    <h4>{{ t('settings.layout.title') }}</h4>
    <SelectLayoutMode />
    <el-divider />
    <h4>{{ t('settings.features.title') }}</h4>
    <div class="setting-item" v-for="(settingValue, settingName, index) in switchSettings" :key="index">
      <span class="setting-name">{{ settingName }}</span>
      <el-switch v-model="settingValue.value" :disabled="!isLeft && settingName === t('settings.layout.fixedHeader')" />
    </div>
    <el-button type="danger" :icon="Refresh" @click="resetConfigLayout">{{ t('settings.reset') }}</el-button>
  </div>
</template>

<style lang="scss" scoped>
@import "@/styles/mixins.scss";

.setting-container {
  padding: 20px;
  .setting-item {
    font-size: 14px;
    color: var(--el-text-color-regular);
    padding: 5px 0;
    display: flex;
    justify-content: space-between;
    align-items: center;
    .setting-name {
      @extend %ellipsis;
    }
  }
  .el-button {
    margin-top: 40px;
    width: 100%;
  }
}
</style>
