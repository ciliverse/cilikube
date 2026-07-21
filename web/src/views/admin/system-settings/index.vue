<template>
  <div class="system-settings">
    <div class="page-header">
      <h1>系统设置</h1>
      <p class="subtitle">管理系统配置、安全设置和认证提供商</p>
    </div>

    <!-- 设置导航 -->
    <div class="settings-nav">
      <el-tabs v-model="activeTab" type="card" @tab-change="handleTabChange">
        <el-tab-pane label="认证设置" name="auth">
          <template #label>
            <div class="tab-label">
              <el-icon><Key /></el-icon>
              <span>认证设置</span>
            </div>
          </template>
        </el-tab-pane>
        
        <el-tab-pane label="安全配置" name="security">
          <template #label>
            <div class="tab-label">
              <el-icon><Lock /></el-icon>
              <span>安全配置</span>
            </div>
          </template>
        </el-tab-pane>
        
        <el-tab-pane label="系统首选项" name="preferences">
          <template #label>
            <div class="tab-label">
              <el-icon><Setting /></el-icon>
              <span>系统首选项</span>
            </div>
          </template>
        </el-tab-pane>
        
        <el-tab-pane label="邮件配置" name="email">
          <template #label>
            <div class="tab-label">
              <el-icon><Message /></el-icon>
              <span>邮件配置</span>
            </div>
          </template>
        </el-tab-pane>
      </el-tabs>
    </div>

    <!-- 设置内容 -->
    <div class="settings-content">
      <!-- 认证设置 -->
      <div v-show="activeTab === 'auth'" class="settings-panel">
        <AuthenticationSettings
          v-model:loading="loading.auth"
          @save="handleSaveAuth"
        />
      </div>

      <!-- 安全配置 -->
      <div v-show="activeTab === 'security'" class="settings-panel">
        <SecuritySettings
          v-model:loading="loading.security"
          @save="handleSaveSecurity"
        />
      </div>

      <!-- 系统首选项 -->
      <div v-show="activeTab === 'preferences'" class="settings-panel">
        <SystemPreferences
          v-model:loading="loading.preferences"
          @save="handleSavePreferences"
        />
      </div>

      <!-- 邮件配置 -->
      <div v-show="activeTab === 'email'" class="settings-panel">
        <EmailSettings
          v-model:loading="loading.email"
          @save="handleSaveEmail"
        />
      </div>
    </div>

    <!-- 保存提示 -->
    <div v-if="hasUnsavedChanges" class="save-reminder">
      <el-alert
        title="您有未保存的更改"
        type="warning"
        :closable="false"
        show-icon
      >
        <template #default>
          <p>请记得保存您的更改，否则将会丢失。</p>
          <div class="reminder-actions">
            <el-button size="small" @click="discardChanges">
              放弃更改
            </el-button>
            <el-button size="small" type="primary" @click="saveAllChanges">
              保存所有更改
            </el-button>
          </div>
        </template>
      </el-alert>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted, onBeforeUnmount } from "vue"
import { ElMessage, ElMessageBox } from "element-plus"
import {
  Key,
  Lock,
  Setting,
  Message
} from "@element-plus/icons-vue"
import AuthenticationSettings from "./components/AuthenticationSettings.vue"
import SecuritySettings from "./components/SecuritySettings.vue"
import SystemPreferences from "./components/SystemPreferences.vue"
import EmailSettings from "./components/EmailSettings.vue"

// 响应式数据
const activeTab = ref("auth")
const hasUnsavedChanges = ref(false)

const loading = reactive({
  auth: false,
  security: false,
  preferences: false,
  email: false
})

/** 标签页切换处理 */
const handleTabChange = async (tabName: string) => {
  if (hasUnsavedChanges.value) {
    try {
      await ElMessageBox.confirm(
        "您有未保存的更改，切换标签页将丢失这些更改。是否继续？",
        "确认切换",
        {
          confirmButtonText: "继续切换",
          cancelButtonText: "取消",
          type: "warning"
        }
      )
      hasUnsavedChanges.value = false
    } catch {
      // 用户取消，恢复原标签页
      return
    }
  }
  activeTab.value = tabName
}

/** 认证设置保存处理 */
const handleSaveAuth = async (data: any) => {
  try {
    loading.auth = true
    // 这里调用保存认证设置的API
    console.log("保存认证设置:", data)
    ElMessage.success("认证设置保存成功")
    hasUnsavedChanges.value = false
  } catch (error: any) {
    console.error("保存认证设置失败:", error)
    ElMessage.error("保存认证设置失败")
  } finally {
    loading.auth = false
  }
}

/** 安全配置保存处理 */
const handleSaveSecurity = async (data: any) => {
  try {
    loading.security = true
    // 这里调用保存安全配置的API
    console.log("保存安全配置:", data)
    ElMessage.success("安全配置保存成功")
    hasUnsavedChanges.value = false
  } catch (error: any) {
    console.error("保存安全配置失败:", error)
    ElMessage.error("保存安全配置失败")
  } finally {
    loading.security = false
  }
}

/** 系统首选项保存处理 */
const handleSavePreferences = async (data: any) => {
  try {
    loading.preferences = true
    // 这里调用保存系统首选项的API
    console.log("保存系统首选项:", data)
    ElMessage.success("系统首选项保存成功")
    hasUnsavedChanges.value = false
  } catch (error: any) {
    console.error("保存系统首选项失败:", error)
    ElMessage.error("保存系统首选项失败")
  } finally {
    loading.preferences = false
  }
}

/** 邮件配置保存处理 */
const handleSaveEmail = async (data: any) => {
  try {
    loading.email = true
    // 这里调用保存邮件配置的API
    console.log("保存邮件配置:", data)
    ElMessage.success("邮件配置保存成功")
    hasUnsavedChanges.value = false
  } catch (error: any) {
    console.error("保存邮件配置失败:", error)
    ElMessage.error("保存邮件配置失败")
  } finally {
    loading.email = false
  }
}

/** 放弃更改 */
const discardChanges = async () => {
  try {
    await ElMessageBox.confirm(
      "确定要放弃所有未保存的更改吗？",
      "确认放弃",
      {
        confirmButtonText: "确定放弃",
        cancelButtonText: "取消",
        type: "warning"
      }
    )
    hasUnsavedChanges.value = false
    // 重新加载当前标签页的数据
    window.location.reload()
  } catch {
    // 用户取消
  }
}

/** 保存所有更改 */
const saveAllChanges = async () => {
  // 根据当前标签页保存对应的设置
  switch (activeTab.value) {
    case 'auth':
      // 触发认证设置保存
      break
    case 'security':
      // 触发安全配置保存
      break
    case 'preferences':
      // 触发系统首选项保存
      break
    case 'email':
      // 触发邮件配置保存
      break
  }
}

/** 页面离开前检查未保存的更改 */
const beforeUnloadHandler = (event: BeforeUnloadEvent) => {
  if (hasUnsavedChanges.value) {
    event.preventDefault()
    event.returnValue = "您有未保存的更改，确定要离开吗？"
    return "您有未保存的更改，确定要离开吗？"
  }
}

/** 组件挂载时添加页面离开监听 */
onMounted(() => {
  window.addEventListener("beforeunload", beforeUnloadHandler)
})

/** 组件卸载时移除监听 */
onBeforeUnmount(() => {
  window.removeEventListener("beforeunload", beforeUnloadHandler)
})
</script>

<style lang="scss" scoped>
.system-settings {
  padding: 24px;
  background-color: var(--el-bg-color-page);
  min-height: calc(100vh - 60px);
}

.page-header {
  margin-bottom: 24px;
  
  h1 {
    margin: 0 0 8px 0;
    color: var(--el-text-color-primary);
    font-size: 28px;
    font-weight: 600;
  }
  
  .subtitle {
    margin: 0;
    color: var(--el-text-color-regular);
    font-size: 16px;
  }
}

.settings-nav {
  margin-bottom: 24px;
  
  .tab-label {
    display: flex;
    align-items: center;
    gap: 8px;
    
    .el-icon {
      font-size: 16px;
    }
  }
  
  :deep(.el-tabs__header) {
    margin-bottom: 0;
  }
  
  :deep(.el-tabs__content) {
    display: none;
  }
}

.settings-content {
  .settings-panel {
    background: var(--el-bg-color);
    border-radius: 8px;
    padding: 24px;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  }
}

.save-reminder {
  position: fixed;
  bottom: 24px;
  right: 24px;
  max-width: 400px;
  z-index: 1000;
  
  .reminder-actions {
    margin-top: 12px;
    display: flex;
    gap: 8px;
  }
}

@media (max-width: 768px) {
  .system-settings {
    padding: 16px;
  }
  
  .settings-content {
    .settings-panel {
      padding: 16px;
    }
  }
  
  .save-reminder {
    bottom: 16px;
    right: 16px;
    left: 16px;
    max-width: none;
  }
}
</style>