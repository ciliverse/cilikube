<script lang="ts" setup>
import { h, onMounted, computed } from "vue"
import { useTheme } from "@/hooks/useTheme"
import { useFont } from "@/hooks/useFont"
import { useLocaleStore } from "@/store/modules/locale"
import { useClusterStore } from "@/store/modules/clusterStore"
import { migrateClusterStorage, needsMigration } from "@/utils/cluster-migration"
import { ElNotification } from "element-plus"
import { useI18n } from "vue-i18n"
// Element Plus 语言包
import zhCn from "element-plus/es/locale/lang/zh-cn"
import en from "element-plus/es/locale/lang/en"

const { t } = useI18n()
const { initTheme } = useTheme()
const { initFont } = useFont()
const localeStore = useLocaleStore()
const clusterStore = useClusterStore()

// Element Plus 语言配置
const elementLocale = computed(() => {
  return localeStore.currentLanguage === 'zh' ? zhCn : en
})

/** 初始化主题 */
initTheme()

/** 初始化字体 */
initFont()

/** 初始化语言 */
localeStore.initLocale()

/** 初始化集群store */
onMounted(async () => {
  try {
    // 首先获取集群列表
    await clusterStore.fetchAvailableClusters()
    
    // 检查是否需要迁移旧的集群存储
    if (needsMigration()) {
      const migrationSuccess = await migrateClusterStorage()
      if (migrationSuccess) {
        ElNotification({
          title: t('common.success'),
          message: t('app.migrationSuccess'),
          type: 'success',
          duration: 3000
        })
      }
    }
  } catch (error) {
    console.warn('Failed to initialize cluster store:', error)
  }
})

/** 作者小心思 */
ElNotification({
  title: t('app.welcomeTitle'),
  type: "success",
  offset: 300,
  message: h(
    "a",
    { style: "color: darkblue", target: "_blank", href: "https://www.cillian.website" },
    t('app.welcomeMessage')
  ),
  duration: 0,
  position: "top-right"
})
</script>

<template>
  <el-config-provider :locale="elementLocale">
    <router-view />
  </el-config-provider>
</template>
