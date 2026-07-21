<template>
  <div class="authentication-settings">
    <div class="settings-header">
      <h2>认证设置</h2>
      <p>配置系统支持的认证方式和OAuth提供商</p>
    </div>

    <el-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-width="140px"
      class="settings-form"
    >
      <!-- 基础认证设置 -->
      <div class="settings-section">
        <h3>基础认证</h3>
        
        <el-form-item label="启用本地认证" prop="local_auth_enabled">
          <el-switch
            v-model="formData.local_auth_enabled"
            active-text="启用"
            inactive-text="禁用"
            @change="handleFormChange"
          />
          <div class="form-help">
            允许用户使用用户名和密码登录
          </div>
        </el-form-item>

        <el-form-item label="允许用户注册" prop="registration_enabled">
          <el-switch
            v-model="formData.registration_enabled"
            active-text="允许"
            inactive-text="禁止"
            @change="handleFormChange"
          />
          <div class="form-help">
            允许新用户自行注册账户
          </div>
        </el-form-item>

        <el-form-item label="默认用户角色" prop="default_user_role">
          <el-select
            v-model="formData.default_user_role"
            placeholder="选择默认角色"
            @change="handleFormChange"
          >
            <el-option
              v-for="role in availableRoles"
              :key="role.value"
              :label="role.label"
              :value="role.value"
            />
          </el-select>
          <div class="form-help">
            新注册用户的默认角色
          </div>
        </el-form-item>
      </div>

      <!-- OAuth 提供商设置 -->
      <div class="settings-section">
        <div class="section-header">
          <h3>OAuth 提供商</h3>
          <el-button
            type="primary"
            size="small"
            :icon="Plus"
            @click="showAddProviderDialog = true"
          >
            添加提供商
          </el-button>
        </div>

        <div class="oauth-providers">
          <div
            v-for="provider in formData.oauth_providers"
            :key="provider.id"
            class="provider-card"
          >
            <div class="provider-header">
              <div class="provider-info">
                <div class="provider-icon">
                  <el-icon :size="24">
                    <component :is="getProviderIcon(provider.type)" />
                  </el-icon>
                </div>
                <div class="provider-details">
                  <h4>{{ provider.name }}</h4>
                  <span class="provider-type">{{ getProviderTypeName(provider.type) }}</span>
                </div>
              </div>
              
              <div class="provider-actions">
                <el-switch
                  v-model="provider.enabled"
                  size="small"
                  @change="handleFormChange"
                />
                <el-button
                  size="small"
                  type="text"
                  @click="editProvider(provider)"
                >
                  编辑
                </el-button>
                <el-button
                  size="small"
                  type="text"
                  style="color: var(--el-color-danger)"
                  @click="deleteProvider(provider)"
                >
                  删除
                </el-button>
              </div>
            </div>

            <div class="provider-config">
              <div class="config-item">
                <span class="config-label">Client ID:</span>
                <span class="config-value">{{ maskClientId(provider.client_id) }}</span>
              </div>
              <div class="config-item">
                <span class="config-label">回调URL:</span>
                <span class="config-value">{{ provider.redirect_uri }}</span>
              </div>
              <div v-if="provider.scopes" class="config-item">
                <span class="config-label">权限范围:</span>
                <div class="config-scopes">
                  <el-tag
                    v-for="scope in provider.scopes"
                    :key="scope"
                    size="small"
                    type="info"
                  >
                    {{ scope }}
                  </el-tag>
                </div>
              </div>
            </div>
          </div>

          <div v-if="formData.oauth_providers.length === 0" class="empty-providers">
            <el-empty
              :image-size="80"
              description="暂未配置OAuth提供商"
            >
              <el-button
                type="primary"
                @click="showAddProviderDialog = true"
              >
                添加第一个提供商
              </el-button>
            </el-empty>
          </div>
        </div>
      </div>

      <!-- JWT 设置 -->
      <div class="settings-section">
        <h3>JWT 令牌设置</h3>
        
        <el-form-item label="访问令牌过期时间" prop="jwt_access_token_ttl">
          <el-input-number
            v-model="formData.jwt_access_token_ttl"
            :min="5"
            :max="1440"
            :step="5"
            @change="handleFormChange"
          />
          <span class="input-suffix">分钟</span>
          <div class="form-help">
            访问令牌的有效期，建议15-60分钟
          </div>
        </el-form-item>

        <el-form-item label="刷新令牌过期时间" prop="jwt_refresh_token_ttl">
          <el-input-number
            v-model="formData.jwt_refresh_token_ttl"
            :min="1"
            :max="30"
            :step="1"
            @change="handleFormChange"
          />
          <span class="input-suffix">天</span>
          <div class="form-help">
            刷新令牌的有效期，建议7-30天
          </div>
        </el-form-item>

        <el-form-item label="令牌签名算法" prop="jwt_algorithm">
          <el-select
            v-model="formData.jwt_algorithm"
            @change="handleFormChange"
          >
            <el-option label="HS256" value="HS256" />
            <el-option label="RS256" value="RS256" />
            <el-option label="ES256" value="ES256" />
          </el-select>
          <div class="form-help">
            JWT令牌的签名算法
          </div>
        </el-form-item>
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

    <!-- 添加OAuth提供商对话框 -->
    <OAuthProviderDialog
      v-model="showAddProviderDialog"
      :provider="editingProvider"
      @save="handleProviderSave"
    />
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted } from "vue"
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from "element-plus"
import { Plus } from "@element-plus/icons-vue"
import OAuthProviderDialog from "./OAuthProviderDialog.vue"

interface OAuthProvider {
  id?: string
  type: string
  name: string
  client_id: string
  client_secret: string
  redirect_uri: string
  scopes: string[]
  enabled: boolean
  config?: Record<string, any>
}

interface Props {
  loading: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
  "update:loading": [value: boolean]
  save: [data: any]
}>()

const formRef = ref<FormInstance>()
const showAddProviderDialog = ref(false)
const editingProvider = ref<OAuthProvider | null>(null)

const formData = reactive({
  local_auth_enabled: true,
  registration_enabled: false,
  default_user_role: "viewer",
  oauth_providers: [] as OAuthProvider[],
  jwt_access_token_ttl: 15,
  jwt_refresh_token_ttl: 7,
  jwt_algorithm: "HS256"
})

const availableRoles = [
  { label: "查看者", value: "viewer" },
  { label: "编辑者", value: "editor" },
  { label: "管理员", value: "admin" }
]

const formRules: FormRules = {
  default_user_role: [
    { required: true, message: "请选择默认用户角色", trigger: "change" }
  ],
  jwt_access_token_ttl: [
    { required: true, message: "请设置访问令牌过期时间", trigger: "blur" }
  ],
  jwt_refresh_token_ttl: [
    { required: true, message: "请设置刷新令牌过期时间", trigger: "blur" }
  ],
  jwt_algorithm: [
    { required: true, message: "请选择令牌签名算法", trigger: "change" }
  ]
}

/** 获取提供商图标 */
const getProviderIcon = (type: string) => {
  const icons: Record<string, string> = {
    github: "Github",
    google: "Google",
    microsoft: "Microsoft",
    gitlab: "Gitlab"
  }
  return icons[type] || "Key"
}

/** 获取提供商类型名称 */
const getProviderTypeName = (type: string) => {
  const names: Record<string, string> = {
    github: "GitHub",
    google: "Google",
    microsoft: "Microsoft",
    gitlab: "GitLab"
  }
  return names[type] || type
}

/** 遮蔽Client ID */
const maskClientId = (clientId: string) => {
  if (!clientId || clientId.length < 8) return clientId
  return clientId.substring(0, 4) + "****" + clientId.substring(clientId.length - 4)
}

/** 表单变更处理 */
const handleFormChange = () => {
  // 标记有未保存的更改
  emit("update:loading", false)
}

/** 编辑提供商 */
const editProvider = (provider: OAuthProvider) => {
  editingProvider.value = { ...provider }
  showAddProviderDialog.value = true
}

/** 删除提供商 */
const deleteProvider = async (provider: OAuthProvider) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除OAuth提供商 "${provider.name}" 吗？`,
      "确认删除",
      {
        confirmButtonText: "确定删除",
        cancelButtonText: "取消",
        type: "warning"
      }
    )
    
    const index = formData.oauth_providers.findIndex(p => p.id === provider.id)
    if (index > -1) {
      formData.oauth_providers.splice(index, 1)
      handleFormChange()
      ElMessage.success("OAuth提供商删除成功")
    }
  } catch {
    // 用户取消
  }
}

/** 提供商保存处理 */
const handleProviderSave = (provider: OAuthProvider) => {
  if (editingProvider.value?.id) {
    // 编辑模式
    const index = formData.oauth_providers.findIndex(p => p.id === editingProvider.value?.id)
    if (index > -1) {
      formData.oauth_providers[index] = provider
    }
  } else {
    // 新增模式
    provider.id = Date.now().toString()
    formData.oauth_providers.push(provider)
  }
  
  editingProvider.value = null
  handleFormChange()
  ElMessage.success("OAuth提供商保存成功")
}

/** 重置表单 */
const resetForm = async () => {
  try {
    await ElMessageBox.confirm(
      "确定要重置所有设置吗？未保存的更改将丢失。",
      "确认重置",
      {
        confirmButtonText: "确定重置",
        cancelButtonText: "取消",
        type: "warning"
      }
    )
    
    // 重新加载设置
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
    
    // 这里调用API加载认证设置
    // const { data } = await getAuthSettingsApi()
    // Object.assign(formData, data)
    
    // 模拟数据
    Object.assign(formData, {
      local_auth_enabled: true,
      registration_enabled: false,
      default_user_role: "viewer",
      oauth_providers: [
        {
          id: "1",
          type: "github",
          name: "GitHub OAuth",
          client_id: "github_client_id_example",
          client_secret: "github_client_secret",
          redirect_uri: "http://localhost:3000/oauth/callback/github",
          scopes: ["user:email", "read:user"],
          enabled: true
        }
      ],
      jwt_access_token_ttl: 15,
      jwt_refresh_token_ttl: 7,
      jwt_algorithm: "HS256"
    })
    
  } catch (error: any) {
    console.error("加载认证设置失败:", error)
    ElMessage.error("加载认证设置失败")
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
.authentication-settings {
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
      
      .section-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 24px;
        
        h3 {
          margin: 0;
        }
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
  
  .oauth-providers {
    .provider-card {
      border: 1px solid var(--el-border-color-light);
      border-radius: 8px;
      padding: 16px;
      margin-bottom: 16px;
      
      .provider-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 12px;
        
        .provider-info {
          display: flex;
          align-items: center;
          
          .provider-icon {
            width: 40px;
            height: 40px;
            border-radius: 8px;
            background: var(--el-color-primary-light-8);
            color: var(--el-color-primary);
            display: flex;
            align-items: center;
            justify-content: center;
            margin-right: 12px;
          }
          
          .provider-details {
            h4 {
              margin: 0 0 4px 0;
              font-size: 14px;
              font-weight: 600;
              color: var(--el-text-color-primary);
            }
            
            .provider-type {
              font-size: 12px;
              color: var(--el-text-color-secondary);
            }
          }
        }
        
        .provider-actions {
          display: flex;
          align-items: center;
          gap: 8px;
        }
      }
      
      .provider-config {
        .config-item {
          display: flex;
          align-items: center;
          margin-bottom: 8px;
          
          &:last-child {
            margin-bottom: 0;
          }
          
          .config-label {
            font-size: 12px;
            color: var(--el-text-color-secondary);
            min-width: 80px;
          }
          
          .config-value {
            font-size: 12px;
            color: var(--el-text-color-regular);
            font-family: monospace;
          }
          
          .config-scopes {
            display: flex;
            flex-wrap: wrap;
            gap: 4px;
          }
        }
      }
    }
    
    .empty-providers {
      text-align: center;
      padding: 40px;
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
</style>