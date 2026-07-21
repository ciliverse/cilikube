<template>
  <div class="security-settings">
    <div class="settings-header">
      <h2>安全配置</h2>
      <p>配置系统安全策略和访问控制</p>
    </div>

    <el-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-width="160px"
      class="settings-form"
    >
      <!-- 密码策略 -->
      <div class="settings-section">
        <h3>密码策略</h3>
        
        <el-form-item label="启用密码复杂度检查" prop="password_complexity_enabled">
          <el-switch
            v-model="formData.password_complexity_enabled"
            active-text="启用"
            inactive-text="禁用"
            @change="handleFormChange"
          />
        </el-form-item>

        <template v-if="formData.password_complexity_enabled">
          <el-form-item label="最小密码长度" prop="password_min_length">
            <el-input-number
              v-model="formData.password_min_length"
              :min="6"
              :max="32"
              @change="handleFormChange"
            />
            <span class="input-suffix">字符</span>
          </el-form-item>

          <el-form-item label="密码复杂度要求">
            <el-checkbox-group v-model="formData.password_requirements" @change="handleFormChange">
              <el-checkbox label="uppercase">包含大写字母</el-checkbox>
              <el-checkbox label="lowercase">包含小写字母</el-checkbox>
              <el-checkbox label="numbers">包含数字</el-checkbox>
              <el-checkbox label="symbols">包含特殊字符</el-checkbox>
            </el-checkbox-group>
          </el-form-item>

          <el-form-item label="密码过期时间" prop="password_expiry_days">
            <el-input-number
              v-model="formData.password_expiry_days"
              :min="0"
              :max="365"
              @change="handleFormChange"
            />
            <span class="input-suffix">天（0表示永不过期）</span>
          </el-form-item>

          <el-form-item label="密码历史记录" prop="password_history_count">
            <el-input-number
              v-model="formData.password_history_count"
              :min="0"
              :max="10"
              @change="handleFormChange"
            />
            <span class="input-suffix">个（防止重复使用最近的密码）</span>
          </el-form-item>
        </template>
      </div>

      <!-- 账户锁定策略 -->
      <div class="settings-section">
        <h3>账户锁定策略</h3>
        
        <el-form-item label="启用账户锁定" prop="account_lockout_enabled">
          <el-switch
            v-model="formData.account_lockout_enabled"
            active-text="启用"
            inactive-text="禁用"
            @change="handleFormChange"
          />
        </el-form-item>

        <template v-if="formData.account_lockout_enabled">
          <el-form-item label="最大失败尝试次数" prop="max_login_attempts">
            <el-input-number
              v-model="formData.max_login_attempts"
              :min="3"
              :max="10"
              @change="handleFormChange"
            />
            <span class="input-suffix">次</span>
          </el-form-item>

          <el-form-item label="锁定时间" prop="lockout_duration">
            <el-input-number
              v-model="formData.lockout_duration"
              :min="5"
              :max="1440"
              @change="handleFormChange"
            />
            <span class="input-suffix">分钟</span>
          </el-form-item>

          <el-form-item label="重置计数器时间" prop="lockout_reset_time">
            <el-input-number
              v-model="formData.lockout_reset_time"
              :min="5"
              :max="60"
              @change="handleFormChange"
            />
            <span class="input-suffix">分钟</span>
            <div class="form-help">
              在此时间内没有失败尝试，则重置失败计数器
            </div>
          </el-form-item>
        </template>
      </div>

      <!-- 会话管理 -->
      <div class="settings-section">
        <h3>会话管理</h3>
        
        <el-form-item label="会话超时时间" prop="session_timeout">
          <el-input-number
            v-model="formData.session_timeout"
            :min="5"
            :max="480"
            @change="handleFormChange"
          />
          <span class="input-suffix">分钟</span>
        </el-form-item>

        <el-form-item label="记住登录状态" prop="remember_me_enabled">
          <el-switch
            v-model="formData.remember_me_enabled"
            active-text="允许"
            inactive-text="禁止"
            @change="handleFormChange"
          />
        </el-form-item>

        <el-form-item v-if="formData.remember_me_enabled" label="记住登录时长" prop="remember_me_duration">
          <el-input-number
            v-model="formData.remember_me_duration"
            :min="1"
            :max="90"
            @change="handleFormChange"
          />
          <span class="input-suffix">天</span>
        </el-form-item>

        <el-form-item label="并发会话限制" prop="max_concurrent_sessions">
          <el-input-number
            v-model="formData.max_concurrent_sessions"
            :min="1"
            :max="10"
            @change="handleFormChange"
          />
          <span class="input-suffix">个</span>
          <div class="form-help">
            每个用户允许的最大并发会话数
          </div>
        </el-form-item>
      </div>

      <!-- IP访问控制 -->
      <div class="settings-section">
        <h3>IP访问控制</h3>
        
        <el-form-item label="启用IP白名单" prop="ip_whitelist_enabled">
          <el-switch
            v-model="formData.ip_whitelist_enabled"
            active-text="启用"
            inactive-text="禁用"
            @change="handleFormChange"
          />
        </el-form-item>

        <el-form-item v-if="formData.ip_whitelist_enabled" label="IP白名单">
          <div class="ip-list-container">
            <div
              v-for="(ip, index) in formData.ip_whitelist"
              :key="index"
              class="ip-item"
            >
              <el-input
                v-model="formData.ip_whitelist[index]"
                placeholder="请输入IP地址或CIDR"
                @change="handleFormChange"
              />
              <el-button
                type="danger"
                size="small"
                @click="removeIP(index)"
              >
                删除
              </el-button>
            </div>
            <el-button
              type="primary"
              size="small"
              @click="addIP"
            >
              添加IP
            </el-button>
          </div>
          <div class="form-help">
            支持单个IP（如：192.168.1.1）或CIDR格式（如：192.168.1.0/24）
          </div>
        </el-form-item>

        <el-form-item label="启用IP黑名单" prop="ip_blacklist_enabled">
          <el-switch
            v-model="formData.ip_blacklist_enabled"
            active-text="启用"
            inactive-text="禁用"
            @change="handleFormChange"
          />
        </el-form-item>

        <el-form-item v-if="formData.ip_blacklist_enabled" label="IP黑名单">
          <div class="ip-list-container">
            <div
              v-for="(ip, index) in formData.ip_blacklist"
              :key="index"
              class="ip-item"
            >
              <el-input
                v-model="formData.ip_blacklist[index]"
                placeholder="请输入IP地址或CIDR"
                @change="handleFormChange"
              />
              <el-button
                type="danger"
                size="small"
                @click="removeBlacklistIP(index)"
              >
                删除
              </el-button>
            </div>
            <el-button
              type="primary"
              size="small"
              @click="addBlacklistIP"
            >
              添加IP
            </el-button>
          </div>
        </el-form-item>
      </div>

      <!-- 审计日志 -->
      <div class="settings-section">
        <h3>审计日志</h3>
        
        <el-form-item label="启用审计日志" prop="audit_log_enabled">
          <el-switch
            v-model="formData.audit_log_enabled"
            active-text="启用"
            inactive-text="禁用"
            @change="handleFormChange"
          />
        </el-form-item>

        <template v-if="formData.audit_log_enabled">
          <el-form-item label="日志级别" prop="audit_log_level">
            <el-select v-model="formData.audit_log_level" @change="handleFormChange">
              <el-option label="基础（登录/登出）" value="basic" />
              <el-option label="标准（包含权限变更）" value="standard" />
              <el-option label="详细（包含所有操作）" value="detailed" />
            </el-select>
          </el-form-item>

          <el-form-item label="日志保留时间" prop="audit_log_retention_days">
            <el-input-number
              v-model="formData.audit_log_retention_days"
              :min="30"
              :max="365"
              @change="handleFormChange"
            />
            <span class="input-suffix">天</span>
          </el-form-item>

          <el-form-item label="记录的事件类型">
            <el-checkbox-group v-model="formData.audit_events" @change="handleFormChange">
              <el-checkbox label="login">用户登录</el-checkbox>
              <el-checkbox label="logout">用户登出</el-checkbox>
              <el-checkbox label="password_change">密码修改</el-checkbox>
              <el-checkbox label="role_change">角色变更</el-checkbox>
              <el-checkbox label="permission_change">权限变更</el-checkbox>
              <el-checkbox label="system_config">系统配置变更</el-checkbox>
              <el-checkbox label="resource_access">资源访问</el-checkbox>
            </el-checkbox-group>
          </el-form-item>
        </template>
      </div>

      <!-- 安全头部 -->
      <div class="settings-section">
        <h3>HTTP安全头部</h3>
        
        <el-form-item label="启用安全头部" prop="security_headers_enabled">
          <el-switch
            v-model="formData.security_headers_enabled"
            active-text="启用"
            inactive-text="禁用"
            @change="handleFormChange"
          />
        </el-form-item>

        <template v-if="formData.security_headers_enabled">
          <el-form-item label="启用的安全头部">
            <el-checkbox-group v-model="formData.security_headers" @change="handleFormChange">
              <el-checkbox label="x-frame-options">X-Frame-Options</el-checkbox>
              <el-checkbox label="x-content-type-options">X-Content-Type-Options</el-checkbox>
              <el-checkbox label="x-xss-protection">X-XSS-Protection</el-checkbox>
              <el-checkbox label="strict-transport-security">Strict-Transport-Security</el-checkbox>
              <el-checkbox label="content-security-policy">Content-Security-Policy</el-checkbox>
            </el-checkbox-group>
          </el-form-item>

          <el-form-item label="CSP策略" prop="csp_policy">
            <el-input
              v-model="formData.csp_policy"
              type="textarea"
              :rows="3"
              placeholder="default-src 'self'; script-src 'self' 'unsafe-inline'"
              @change="handleFormChange"
            />
            <div class="form-help">
              Content Security Policy策略配置
            </div>
          </el-form-item>
        </template>
      </div>
    </el-form>

    <!-- 操作按钮 -->
    <div class="settings-actions">
      <el-button @click="resetForm">重置</el-button>
      <el-button
        type="primary"
        @click="saveSettings"
        :loading="loading"
      >
        保存设置
      </el-button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted } from "vue"
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from "element-plus"

interface Props {
  loading: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
  "update:loading": [value: boolean]
  save: [data: any]
}>()

const formRef = ref<FormInstance>()

const formData = reactive({
  // 密码策略
  password_complexity_enabled: true,
  password_min_length: 8,
  password_requirements: ["lowercase", "numbers"],
  password_expiry_days: 90,
  password_history_count: 3,
  
  // 账户锁定
  account_lockout_enabled: true,
  max_login_attempts: 5,
  lockout_duration: 30,
  lockout_reset_time: 15,
  
  // 会话管理
  session_timeout: 60,
  remember_me_enabled: true,
  remember_me_duration: 30,
  max_concurrent_sessions: 3,
  
  // IP访问控制
  ip_whitelist_enabled: false,
  ip_whitelist: [] as string[],
  ip_blacklist_enabled: false,
  ip_blacklist: [] as string[],
  
  // 审计日志
  audit_log_enabled: true,
  audit_log_level: "standard",
  audit_log_retention_days: 90,
  audit_events: ["login", "logout", "password_change", "role_change"],
  
  // 安全头部
  security_headers_enabled: true,
  security_headers: ["x-frame-options", "x-content-type-options", "x-xss-protection"],
  csp_policy: "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'"
})

const formRules: FormRules = {
  password_min_length: [
    { required: true, message: "请设置最小密码长度", trigger: "blur" }
  ],
  password_expiry_days: [
    { required: true, message: "请设置密码过期时间", trigger: "blur" }
  ],
  max_login_attempts: [
    { required: true, message: "请设置最大失败尝试次数", trigger: "blur" }
  ],
  lockout_duration: [
    { required: true, message: "请设置锁定时间", trigger: "blur" }
  ],
  session_timeout: [
    { required: true, message: "请设置会话超时时间", trigger: "blur" }
  ],
  audit_log_retention_days: [
    { required: true, message: "请设置日志保留时间", trigger: "blur" }
  ]
}

/** 表单变更处理 */
const handleFormChange = () => {
  emit("update:loading", false)
}

/** 添加IP到白名单 */
const addIP = () => {
  formData.ip_whitelist.push("")
  handleFormChange()
}

/** 从白名单移除IP */
const removeIP = (index: number) => {
  formData.ip_whitelist.splice(index, 1)
  handleFormChange()
}

/** 添加IP到黑名单 */
const addBlacklistIP = () => {
  formData.ip_blacklist.push("")
  handleFormChange()
}

/** 从黑名单移除IP */
const removeBlacklistIP = (index: number) => {
  formData.ip_blacklist.splice(index, 1)
  handleFormChange()
}

/** 重置表单 */
const resetForm = async () => {
  try {
    await ElMessageBox.confirm(
      "确定要重置所有安全设置吗？未保存的更改将丢失。",
      "确认重置",
      {
        confirmButtonText: "确定重置",
        cancelButtonText: "取消",
        type: "warning"
      }
    )
    
    loadSettings()
  } catch {
    // 用户取消
  }
}

/** 保存设置 */
const saveSettings = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    emit("save", { ...formData })
  } catch (error) {
    ElMessage.error("请检查表单输入")
  }
}

/** 加载设置 */
const loadSettings = async () => {
  try {
    emit("update:loading", true)
    
    // 这里调用API加载安全设置
    // const { data } = await getSecuritySettingsApi()
    // Object.assign(formData, data)
    
    // 模拟数据已在reactive中设置
    
  } catch (error: any) {
    console.error("加载安全设置失败:", error)
    ElMessage.error("加载安全设置失败")
  } finally {
    emit("update:loading", false)
  }
}

/** 组件挂载时加载设置 */
onMounted(() => {
  loadSettings()
})
</script>

<style lang="scss" scoped>
.security-settings {
  .settings-header {
    margin-bottom: 32px;
    
    h2 {
      margin: 0 0 8px 0;
      color: var(--el-text-color-primary);
      font-size: 20px;
      font-weight: 600;
    }
    
    p {
      margin: 0;
      color: var(--el-text-color-regular);
    }
  }
  
  .settings-form {
    .settings-section {
      margin-bottom: 40px;
      padding-bottom: 32px;
      border-bottom: 1px solid var(--el-border-color-lighter);
      
      &:last-child {
        border-bottom: none;
        margin-bottom: 0;
      }
      
      h3 {
        margin: 0 0 24px 0;
        color: var(--el-text-color-primary);
        font-size: 16px;
        font-weight: 600;
      }
    }
    
    .form-help {
      margin-top: 4px;
      font-size: 12px;
      color: var(--el-text-color-secondary);
      line-height: 1.4;
    }
    
    .input-suffix {
      margin-left: 8px;
      color: var(--el-text-color-regular);
      font-size: 14px;
    }
  }
  
  .ip-list-container {
    .ip-item {
      display: flex;
      align-items: center;
      gap: 8px;
      margin-bottom: 8px;
      
      .el-input {
        flex: 1;
      }
    }
  }
  
  .settings-actions {
    margin-top: 32px;
    padding-top: 24px;
    border-top: 1px solid var(--el-border-color-lighter);
    display: flex;
    justify-content: flex-end;
    gap: 12px;
  }
}

:deep(.el-checkbox-group) {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>