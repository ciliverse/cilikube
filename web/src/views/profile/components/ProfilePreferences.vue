<template>
  <div class="profile-preferences">
    <!-- 界面设置 -->
    <el-card class="preference-card">
      <template #header>
        <span>界面设置</span>
      </template>

      <el-form label-width="120px">
        <el-form-item label="主题模式">
          <el-radio-group v-model="preferences.theme" @change="updateTheme">
            <el-radio label="light">浅色主题</el-radio>
            <el-radio label="dark">深色主题</el-radio>
            <el-radio label="auto">跟随系统</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="语言设置">
          <el-select 
            v-model="preferences.language" 
            @change="updateLanguage"
            style="width: 200px"
          >
            <el-option label="简体中文" value="zh-CN" />
            <el-option label="English" value="en-US" />
            <el-option label="繁體中文" value="zh-TW" />
          </el-select>
        </el-form-item>

        <el-form-item label="侧边栏">
          <el-switch
            v-model="preferences.sidebarCollapsed"
            @change="updateSidebarCollapsed"
            active-text="默认收起"
            inactive-text="默认展开"
          />
        </el-form-item>

        <el-form-item label="面包屑导航">
          <el-switch
            v-model="preferences.showBreadcrumb"
            @change="updateShowBreadcrumb"
            active-text="显示"
            inactive-text="隐藏"
          />
        </el-form-item>

        <el-form-item label="标签页">
          <el-switch
            v-model="preferences.showTabs"
            @change="updateShowTabs"
            active-text="显示"
            inactive-text="隐藏"
          />
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 通知设置 -->
    <el-card class="preference-card">
      <template #header>
        <span>通知设置</span>
      </template>

      <el-form label-width="120px">
        <el-form-item label="邮件通知">
          <div class="notification-options">
            <el-checkbox 
              v-model="preferences.notifications.email.security"
              @change="updateNotifications"
            >
              安全相关通知
            </el-checkbox>
            <el-checkbox 
              v-model="preferences.notifications.email.system"
              @change="updateNotifications"
            >
              系统更新通知
            </el-checkbox>
            <el-checkbox 
              v-model="preferences.notifications.email.marketing"
              @change="updateNotifications"
            >
              产品推广信息
            </el-checkbox>
          </div>
        </el-form-item>

        <el-form-item label="浏览器通知">
          <div class="notification-options">
            <el-checkbox 
              v-model="preferences.notifications.browser.enabled"
              @change="updateNotifications"
            >
              启用浏览器通知
            </el-checkbox>
            <el-checkbox 
              v-model="preferences.notifications.browser.sound"
              @change="updateNotifications"
              :disabled="!preferences.notifications.browser.enabled"
            >
              通知声音
            </el-checkbox>
          </div>
        </el-form-item>

        <el-form-item label="通知频率">
          <el-select 
            v-model="preferences.notifications.frequency" 
            @change="updateNotifications"
            style="width: 200px"
          >
            <el-option label="实时通知" value="realtime" />
            <el-option label="每小时汇总" value="hourly" />
            <el-option label="每日汇总" value="daily" />
            <el-option label="每周汇总" value="weekly" />
          </el-select>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据和隐私 -->
    <el-card class="preference-card">
      <template #header>
        <span>数据和隐私</span>
      </template>

      <el-form label-width="120px">
        <el-form-item label="数据收集">
          <el-switch
            v-model="preferences.privacy.analytics"
            @change="updatePrivacy"
            active-text="允许"
            inactive-text="禁止"
          />
          <div class="form-help">
            帮助我们改进产品体验，不会收集个人敏感信息
          </div>
        </el-form-item>

        <el-form-item label="错误报告">
          <el-switch
            v-model="preferences.privacy.errorReporting"
            @change="updatePrivacy"
            active-text="自动发送"
            inactive-text="手动发送"
          />
          <div class="form-help">
            自动发送错误报告帮助我们快速修复问题
          </div>
        </el-form-item>

        <el-form-item label="会话保持">
          <el-select 
            v-model="preferences.privacy.sessionDuration" 
            @change="updatePrivacy"
            style="width: 200px"
          >
            <el-option label="1小时" value="1h" />
            <el-option label="8小时" value="8h" />
            <el-option label="1天" value="1d" />
            <el-option label="7天" value="7d" />
            <el-option label="30天" value="30d" />
          </el-select>
          <div class="form-help">
            设置自动登出时间，提高账户安全性
          </div>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 高级设置 -->
    <el-card class="preference-card">
      <template #header>
        <span>高级设置</span>
      </template>

      <el-form label-width="120px">
        <el-form-item label="开发者模式">
          <el-switch
            v-model="preferences.advanced.developerMode"
            @change="updateAdvanced"
            active-text="启用"
            inactive-text="禁用"
          />
          <div class="form-help">
            启用后将显示更多调试信息和开发工具
          </div>
        </el-form-item>

        <el-form-item label="实验性功能">
          <el-switch
            v-model="preferences.advanced.experimentalFeatures"
            @change="updateAdvanced"
            active-text="启用"
            inactive-text="禁用"
          />
          <div class="form-help">
            体验最新功能，可能存在不稳定因素
          </div>
        </el-form-item>

        <el-form-item label="API调试">
          <el-switch
            v-model="preferences.advanced.apiDebug"
            @change="updateAdvanced"
            active-text="启用"
            inactive-text="禁用"
          />
          <div class="form-help">
            在控制台显示API请求和响应信息
          </div>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据管理 -->
    <el-card class="preference-card">
      <template #header>
        <span>数据管理</span>
      </template>

      <div class="data-management">
        <div class="data-item">
          <div class="data-info">
            <h4>导出个人数据</h4>
            <p>下载您在平台上的所有个人数据</p>
          </div>
          <el-button 
            @click="exportUserData"
            :loading="exportLoading"
          >
            导出数据
          </el-button>
        </div>

        <div class="data-item">
          <div class="data-info">
            <h4>清除缓存数据</h4>
            <p>清除浏览器中存储的缓存和临时数据</p>
          </div>
          <el-button 
            @click="clearCacheData"
            :loading="clearCacheLoading"
          >
            清除缓存
          </el-button>
        </div>

        <div class="data-item danger">
          <div class="data-info">
            <h4>删除账户</h4>
            <p>永久删除您的账户和所有相关数据，此操作不可恢复</p>
          </div>
          <el-button 
            type="danger" 
            @click="deleteAccount"
            :loading="deleteLoading"
          >
            删除账户
          </el-button>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted } from "vue"
import { ElMessage, ElMessageBox } from "element-plus"
import { 
  getUserPreferencesApi, 
  updateUserPreferencesApi,
  exportUserDataApi,
  clearUserCacheApi,
  deleteUserAccountApi
} from "@/api/profile"

const emit = defineEmits<{
  loading: [loading: boolean]
}>()

const exportLoading = ref(false)
const clearCacheLoading = ref(false)
const deleteLoading = ref(false)

const preferences = reactive({
  theme: "light",
  language: "zh-CN",
  sidebarCollapsed: false,
  showBreadcrumb: true,
  showTabs: true,
  notifications: {
    email: {
      security: true,
      system: true,
      marketing: false
    },
    browser: {
      enabled: true,
      sound: true
    },
    frequency: "realtime"
  },
  privacy: {
    analytics: true,
    errorReporting: true,
    sessionDuration: "8h"
  },
  advanced: {
    developerMode: false,
    experimentalFeatures: false,
    apiDebug: false
  }
})

/** 加载用户偏好设置 */
const loadUserPreferences = async () => {
  try {
    const { data } = await getUserPreferencesApi()
    Object.assign(preferences, data)
  } catch (error: any) {
    console.error("加载偏好设置失败:", error)
    ElMessage.error("加载偏好设置失败")
  }
}

/** 更新偏好设置 */
const updatePreferences = async (section: string, data: any) => {
  try {
    await updateUserPreferencesApi({ [section]: data })
    ElMessage.success("设置已保存")
  } catch (error: any) {
    console.error("保存设置失败:", error)
    ElMessage.error("保存设置失败")
  }
}

/** 更新主题 */
const updateTheme = (theme: string) => {
  updatePreferences("theme", theme)
  // 应用主题到界面
  document.documentElement.setAttribute("data-theme", theme)
}

/** 更新语言 */
const updateLanguage = (language: string) => {
  updatePreferences("language", language)
  // 这里可以集成i18n来切换语言
}

/** 更新侧边栏状态 */
const updateSidebarCollapsed = (collapsed: boolean) => {
  updatePreferences("sidebarCollapsed", collapsed)
}

/** 更新面包屑显示 */
const updateShowBreadcrumb = (show: boolean) => {
  updatePreferences("showBreadcrumb", show)
}

/** 更新标签页显示 */
const updateShowTabs = (show: boolean) => {
  updatePreferences("showTabs", show)
}

/** 更新通知设置 */
const updateNotifications = () => {
  updatePreferences("notifications", preferences.notifications)
}

/** 更新隐私设置 */
const updatePrivacy = () => {
  updatePreferences("privacy", preferences.privacy)
}

/** 更新高级设置 */
const updateAdvanced = () => {
  updatePreferences("advanced", preferences.advanced)
}

/** 导出用户数据 */
const exportUserData = async () => {
  try {
    await ElMessageBox.confirm(
      "将导出您的所有个人数据，包括个人信息、设置和活动记录。",
      "确认导出",
      {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "info"
      }
    )
    
    exportLoading.value = true
    emit("loading", true)
    
    const response = await exportUserDataApi()
    
    // 创建下载链接
    const blob = new Blob([JSON.stringify(response.data, null, 2)], {
      type: "application/json"
    })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement("a")
    link.href = url
    link.download = `user-data-${new Date().toISOString().split('T')[0]}.json`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    
    ElMessage.success("数据导出成功")
    
  } catch (error: any) {
    if (error !== "cancel") {
      console.error("导出数据失败:", error)
      ElMessage.error("导出数据失败")
    }
  } finally {
    exportLoading.value = false
    emit("loading", false)
  }
}

/** 清除缓存数据 */
const clearCacheData = async () => {
  try {
    await ElMessageBox.confirm(
      "将清除浏览器中的缓存数据，您可能需要重新登录。",
      "确认清除",
      {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning"
      }
    )
    
    clearCacheLoading.value = true
    emit("loading", true)
    
    await clearUserCacheApi()
    
    // 清除本地存储
    localStorage.clear()
    sessionStorage.clear()
    
    ElMessage.success("缓存已清除")
    
    // 延迟刷新页面
    setTimeout(() => {
      window.location.reload()
    }, 1000)
    
  } catch (error: any) {
    if (error !== "cancel") {
      console.error("清除缓存失败:", error)
      ElMessage.error("清除缓存失败")
    }
  } finally {
    clearCacheLoading.value = false
    emit("loading", false)
  }
}

/** 删除账户 */
const deleteAccount = async () => {
  try {
    await ElMessageBox.confirm(
      "此操作将永久删除您的账户和所有相关数据，且无法恢复。请输入您的用户名确认删除。",
      "危险操作",
      {
        confirmButtonText: "确认删除",
        cancelButtonText: "取消",
        type: "error",
        inputPattern: /.+/,
        inputPlaceholder: "请输入您的用户名",
        inputErrorMessage: "请输入用户名"
      }
    )
    
    deleteLoading.value = true
    emit("loading", true)
    
    await deleteUserAccountApi()
    
    ElMessage.success("账户删除成功")
    
    // 清除所有数据并跳转到登录页
    localStorage.clear()
    sessionStorage.clear()
    
    setTimeout(() => {
      window.location.href = "/login"
    }, 1000)
    
  } catch (error: any) {
    if (error !== "cancel") {
      console.error("删除账户失败:", error)
      ElMessage.error("删除账户失败")
    }
  } finally {
    deleteLoading.value = false
    emit("loading", false)
  }
}

/** 组件挂载时加载数据 */
onMounted(() => {
  // TODO: 暂时禁用偏好设置加载，等后端实现后再启用
  // loadUserPreferences()
})
</script>

<style lang="scss" scoped>
.profile-preferences {
  .preference-card {
    margin-bottom: 24px;
    
    :deep(.el-card__header) {
      span {
        font-weight: 600;
        font-size: 16px;
      }
    }
  }
  
  .notification-options {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }
  
  .form-help {
    margin-top: 4px;
    font-size: 12px;
    color: var(--el-text-color-secondary);
    line-height: 1.4;
  }
  
  .data-management {
    .data-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 20px;
      border: 1px solid var(--el-border-color-light);
      border-radius: 8px;
      margin-bottom: 16px;
      
      &.danger {
        border-color: var(--el-color-danger-light-7);
        background-color: var(--el-color-danger-light-9);
      }
      
      .data-info {
        flex: 1;
        
        h4 {
          margin: 0 0 8px 0;
          font-size: 16px;
          font-weight: 600;
          color: var(--el-text-color-primary);
        }
        
        p {
          margin: 0;
          font-size: 14px;
          color: var(--el-text-color-regular);
          line-height: 1.5;
        }
      }
    }
  }
  
  :deep(.el-form-item__label) {
    font-weight: 500;
  }
  
  :deep(.el-radio-group) {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }
}

@media (max-width: 768px) {
  .data-management {
    .data-item {
      flex-direction: column;
      align-items: flex-start;
      
      .data-info {
        margin-bottom: 16px;
      }
    }
  }
}
</style>