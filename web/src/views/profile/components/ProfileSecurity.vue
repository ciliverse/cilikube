<template>
  <div class="profile-security">
    <!-- 密码修改 -->
    <el-card class="security-card">
      <template #header>
        <div class="card-header">
          <span>修改密码</span>
        </div>
      </template>

      <el-form
        ref="passwordFormRef"
        :model="passwordForm"
        :rules="passwordRules"
        label-width="120px"
      >
        <el-form-item label="当前密码" prop="currentPassword">
          <el-input
            v-model="passwordForm.currentPassword"
            type="password"
            placeholder="请输入当前密码"
            show-password
            autocomplete="current-password"
          />
        </el-form-item>

        <el-form-item label="新密码" prop="newPassword">
          <el-input
            v-model="passwordForm.newPassword"
            type="password"
            placeholder="请输入新密码"
            show-password
            autocomplete="new-password"
          />
          <div class="password-strength">
            <div class="strength-bar">
              <div 
                class="strength-fill" 
                :class="passwordStrengthClass"
                :style="{ width: passwordStrengthWidth }"
              ></div>
            </div>
            <span class="strength-text">{{ passwordStrengthText }}</span>
          </div>
        </el-form-item>

        <el-form-item label="确认新密码" prop="confirmPassword">
          <el-input
            v-model="passwordForm.confirmPassword"
            type="password"
            placeholder="请再次输入新密码"
            show-password
            autocomplete="new-password"
          />
        </el-form-item>

        <el-form-item>
          <el-button 
            type="primary" 
            @click="changePassword"
            :loading="passwordLoading"
          >
            修改密码
          </el-button>
          <el-button @click="resetPasswordForm">
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 登录会话管理 -->
    <el-card class="security-card">
      <template #header>
        <div class="card-header">
          <span>登录会话</span>
          <el-button 
            type="danger" 
            size="small" 
            @click="logoutAllSessions"
            :loading="logoutAllLoading"
          >
            注销所有会话
          </el-button>
        </div>
      </template>

      <div class="sessions-list">
        <div v-if="sessionsLoading" class="loading-container">
          <el-icon class="loading-icon"><Loading /></el-icon>
          <span>加载中...</span>
        </div>
        
        <div v-else-if="sessions.length === 0" class="empty-sessions">
          <el-empty description="暂无活跃会话" />
        </div>
        
        <div v-else>
          <div 
            v-for="session in sessions" 
            :key="session.id" 
            class="session-item"
          >
            <div class="session-info">
              <div class="session-header">
                <el-icon class="device-icon">
                  <Monitor v-if="session.device_type === 'desktop'" />
                  <Cellphone v-else-if="session.device_type === 'mobile'" />
                  <Monitor v-else />
                </el-icon>
                <span class="device-name">{{ session.device_name }}</span>
                <el-tag 
                  v-if="session.is_current" 
                  type="success" 
                  size="small"
                >
                  当前会话
                </el-tag>
              </div>
              
              <div class="session-details">
                <p><strong>IP地址:</strong> {{ session.ip_address }}</p>
                <p><strong>位置:</strong> {{ session.location || '未知' }}</p>
                <p><strong>浏览器:</strong> {{ session.user_agent }}</p>
                <p><strong>最后活跃:</strong> {{ formatDate(session.last_activity) }}</p>
              </div>
            </div>
            
            <div class="session-actions">
              <el-button 
                v-if="!session.is_current"
                type="danger" 
                size="small" 
                @click="logoutSession(session.id)"
                :loading="logoutSessionLoading[session.id]"
              >
                注销
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 安全日志 -->
    <el-card class="security-card">
      <template #header>
        <div class="card-header">
          <span>安全日志</span>
          <el-button 
            size="small" 
            @click="refreshSecurityLogs"
            :loading="securityLogsLoading"
          >
            刷新
          </el-button>
        </div>
      </template>

      <div class="security-logs">
        <div v-if="securityLogsLoading" class="loading-container">
          <el-icon class="loading-icon"><Loading /></el-icon>
          <span>加载中...</span>
        </div>
        
        <div v-else-if="securityLogs.length === 0" class="empty-logs">
          <el-empty description="暂无安全日志" />
        </div>
        
        <el-timeline v-else>
          <el-timeline-item
            v-for="log in securityLogs"
            :key="log.id"
            :timestamp="formatDate(log.created_at)"
            :type="getLogType(log.action)"
          >
            <div class="log-content">
              <h4>{{ getLogTitle(log.action) }}</h4>
              <p>{{ log.description }}</p>
              <div class="log-meta">
                <span>IP: {{ log.ip_address }}</span>
                <span v-if="log.location">位置: {{ log.location }}</span>
              </div>
            </div>
          </el-timeline-item>
        </el-timeline>
      </div>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, computed, onMounted } from "vue"
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from "element-plus"
import { 
  Loading, 
  Monitor, 
  Cellphone 
} from "@element-plus/icons-vue"
import { 
  changePasswordApi, 
  getUserSessionsApi, 
  logoutSessionApi, 
  logoutAllSessionsApi,
  getSecurityLogsApi 
} from "@/api/profile"

const emit = defineEmits<{
  loading: [loading: boolean]
}>()

// 密码修改相关
const passwordFormRef = ref<FormInstance>()
const passwordLoading = ref(false)

const passwordForm = reactive({
  currentPassword: "",
  newPassword: "",
  confirmPassword: ""
})

const passwordRules: FormRules = {
  currentPassword: [
    { required: true, message: "请输入当前密码", trigger: "blur" }
  ],
  newPassword: [
    { required: true, message: "请输入新密码", trigger: "blur" },
    { min: 6, message: "密码长度至少6个字符", trigger: "blur" },
    { 
      validator: (rule, value, callback) => {
        if (value === passwordForm.currentPassword) {
          callback(new Error("新密码不能与当前密码相同"))
        } else {
          callback()
        }
      }, 
      trigger: "blur" 
    }
  ],
  confirmPassword: [
    { required: true, message: "请确认新密码", trigger: "blur" },
    { 
      validator: (rule, value, callback) => {
        if (value !== passwordForm.newPassword) {
          callback(new Error("两次输入的密码不一致"))
        } else {
          callback()
        }
      }, 
      trigger: "blur" 
    }
  ]
}

// 密码强度计算
const passwordStrength = computed(() => {
  const password = passwordForm.newPassword
  if (!password) return 0
  
  let score = 0
  
  // 长度检查
  if (password.length >= 6) score += 20
  if (password.length >= 8) score += 10
  if (password.length >= 12) score += 10
  
  // 字符类型检查
  if (/[a-z]/.test(password)) score += 15
  if (/[A-Z]/.test(password)) score += 15
  if (/[0-9]/.test(password)) score += 15
  if (/[^A-Za-z0-9]/.test(password)) score += 15
  
  return Math.min(score, 100)
})

const passwordStrengthClass = computed(() => {
  const strength = passwordStrength.value
  if (strength < 30) return "weak"
  if (strength < 60) return "medium"
  if (strength < 80) return "good"
  return "strong"
})

const passwordStrengthWidth = computed(() => `${passwordStrength.value}%`)

const passwordStrengthText = computed(() => {
  const strength = passwordStrength.value
  if (strength === 0) return ""
  if (strength < 30) return "弱"
  if (strength < 60) return "中等"
  if (strength < 80) return "良好"
  return "强"
})

// 会话管理相关
const sessions = ref<any[]>([])
const sessionsLoading = ref(false)
const logoutAllLoading = ref(false)
const logoutSessionLoading = ref<Record<string, boolean>>({})

// 安全日志相关
const securityLogs = ref<any[]>([])
const securityLogsLoading = ref(false)

/** 修改密码 */
const changePassword = async () => {
  if (!passwordFormRef.value) return
  
  try {
    await passwordFormRef.value.validate()
    
    passwordLoading.value = true
    emit("loading", true)
    
    await changePasswordApi({
      current_password: passwordForm.currentPassword,
      new_password: passwordForm.newPassword
    })
    
    ElMessage.success("密码修改成功")
    resetPasswordForm()
    
  } catch (error: any) {
    console.error("修改密码失败:", error)
    ElMessage.error(error.response?.data?.message || "修改密码失败")
  } finally {
    passwordLoading.value = false
    emit("loading", false)
  }
}

/** 重置密码表单 */
const resetPasswordForm = () => {
  passwordForm.currentPassword = ""
  passwordForm.newPassword = ""
  passwordForm.confirmPassword = ""
  passwordFormRef.value?.clearValidate()
}

/** 加载用户会话 */
const loadUserSessions = async () => {
  try {
    sessionsLoading.value = true
    const { data } = await getUserSessionsApi()
    sessions.value = data
  } catch (error: any) {
    console.error("加载会话失败:", error)
    ElMessage.error("加载会话失败")
  } finally {
    sessionsLoading.value = false
  }
}

/** 注销单个会话 */
const logoutSession = async (sessionId: string) => {
  try {
    await ElMessageBox.confirm(
      "确定要注销这个会话吗？",
      "确认注销",
      {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning"
      }
    )
    
    logoutSessionLoading.value[sessionId] = true
    
    await logoutSessionApi(sessionId)
    
    ElMessage.success("会话已注销")
    await loadUserSessions()
    
  } catch (error: any) {
    if (error !== "cancel") {
      console.error("注销会话失败:", error)
      ElMessage.error("注销会话失败")
    }
  } finally {
    logoutSessionLoading.value[sessionId] = false
  }
}

/** 注销所有会话 */
const logoutAllSessions = async () => {
  try {
    await ElMessageBox.confirm(
      "确定要注销所有其他会话吗？这将强制其他设备重新登录。",
      "确认注销",
      {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning"
      }
    )
    
    logoutAllLoading.value = true
    
    await logoutAllSessionsApi()
    
    ElMessage.success("所有其他会话已注销")
    await loadUserSessions()
    
  } catch (error: any) {
    if (error !== "cancel") {
      console.error("注销所有会话失败:", error)
      ElMessage.error("注销所有会话失败")
    }
  } finally {
    logoutAllLoading.value = false
  }
}

/** 加载安全日志 */
const loadSecurityLogs = async () => {
  try {
    securityLogsLoading.value = true
    const { data } = await getSecurityLogsApi()
    securityLogs.value = data
  } catch (error: any) {
    console.error("加载安全日志失败:", error)
    ElMessage.error("加载安全日志失败")
  } finally {
    securityLogsLoading.value = false
  }
}

/** 刷新安全日志 */
const refreshSecurityLogs = () => {
  loadSecurityLogs()
}

/** 格式化日期 */
const formatDate = (dateString: string) => {
  if (!dateString) return "未知"
  return new Date(dateString).toLocaleString("zh-CN")
}

/** 获取日志类型 */
const getLogType = (action: string) => {
  switch (action) {
    case "login":
    case "oauth_login":
      return "success"
    case "logout":
      return "info"
    case "password_change":
      return "warning"
    case "failed_login":
    case "suspicious_activity":
      return "danger"
    default:
      return "primary"
  }
}

/** 获取日志标题 */
const getLogTitle = (action: string) => {
  const titles: Record<string, string> = {
    login: "用户登录",
    oauth_login: "OAuth登录",
    logout: "用户登出",
    password_change: "密码修改",
    failed_login: "登录失败",
    suspicious_activity: "可疑活动"
  }
  return titles[action] || action
}

/** 组件挂载时加载数据 */
onMounted(() => {
  // TODO: 暂时禁用这些API调用，等后端实现后再启用
  // loadUserSessions()
  // loadSecurityLogs()
})
</script>

<style lang="scss" scoped>
.profile-security {
  .security-card {
    margin-bottom: 24px;
    
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      
      span {
        font-weight: 600;
        font-size: 16px;
      }
    }
  }
  
  .password-strength {
    margin-top: 8px;
    
    .strength-bar {
      width: 100%;
      height: 4px;
      background-color: var(--el-border-color-light);
      border-radius: 2px;
      overflow: hidden;
      margin-bottom: 4px;
      
      .strength-fill {
        height: 100%;
        transition: all 0.3s ease;
        
        &.weak {
          background-color: var(--el-color-danger);
        }
        
        &.medium {
          background-color: var(--el-color-warning);
        }
        
        &.good {
          background-color: var(--el-color-primary);
        }
        
        &.strong {
          background-color: var(--el-color-success);
        }
      }
    }
    
    .strength-text {
      font-size: 12px;
      color: var(--el-text-color-secondary);
    }
  }
  
  .sessions-list,
  .security-logs {
    .loading-container {
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 40px;
      color: var(--el-text-color-secondary);
      
      .loading-icon {
        margin-right: 8px;
        animation: rotate 2s linear infinite;
      }
    }
    
    .empty-sessions,
    .empty-logs {
      padding: 20px;
    }
  }
  
  .session-item {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    padding: 16px;
    border: 1px solid var(--el-border-color-light);
    border-radius: 8px;
    margin-bottom: 12px;
    
    .session-info {
      flex: 1;
      
      .session-header {
        display: flex;
        align-items: center;
        margin-bottom: 8px;
        
        .device-icon {
          margin-right: 8px;
          color: var(--el-color-primary);
        }
        
        .device-name {
          font-weight: 500;
          margin-right: 8px;
        }
      }
      
      .session-details {
        p {
          margin: 4px 0;
          font-size: 14px;
          color: var(--el-text-color-regular);
          
          strong {
            color: var(--el-text-color-primary);
          }
        }
      }
    }
    
    .session-actions {
      margin-left: 16px;
    }
  }
  
  .log-content {
    h4 {
      margin: 0 0 8px 0;
      font-size: 16px;
      color: var(--el-text-color-primary);
    }
    
    p {
      margin: 0 0 8px 0;
      color: var(--el-text-color-regular);
    }
    
    .log-meta {
      font-size: 12px;
      color: var(--el-text-color-secondary);
      
      span {
        margin-right: 16px;
      }
    }
  }
}

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>