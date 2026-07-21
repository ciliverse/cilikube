<template>
  <el-dialog
    v-model="dialogVisible"
    :title="isEditing ? '编辑OAuth提供商' : '添加OAuth提供商'"
    width="600px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-width="120px"
    >
      <el-form-item label="提供商类型" prop="type">
        <el-select
          v-model="formData.type"
          placeholder="选择OAuth提供商类型"
          :disabled="isEditing"
          @change="handleTypeChange"
        >
          <el-option
            v-for="type in providerTypes"
            :key="type.value"
            :label="type.label"
            :value="type.value"
          >
            <div class="provider-option">
              <el-icon class="provider-icon">
                <component :is="type.icon" />
              </el-icon>
              <span>{{ type.label }}</span>
            </div>
          </el-option>
        </el-select>
      </el-form-item>

      <el-form-item label="提供商名称" prop="name">
        <el-input
          v-model="formData.name"
          placeholder="请输入提供商名称"
          maxlength="100"
        />
      </el-form-item>

      <el-form-item label="Client ID" prop="client_id">
        <el-input
          v-model="formData.client_id"
          placeholder="请输入Client ID"
          maxlength="200"
        />
      </el-form-item>

      <el-form-item label="Client Secret" prop="client_secret">
        <el-input
          v-model="formData.client_secret"
          type="password"
          placeholder="请输入Client Secret"
          maxlength="200"
          show-password
        />
      </el-form-item>

      <el-form-item label="回调URL" prop="redirect_uri">
        <el-input
          v-model="formData.redirect_uri"
          placeholder="请输入回调URL"
          maxlength="500"
        />
        <div class="form-help">
          OAuth认证成功后的回调地址
        </div>
      </el-form-item>

      <el-form-item label="权限范围" prop="scopes">
        <div class="scopes-container">
          <el-tag
            v-for="scope in formData.scopes"
            :key="scope"
            closable
            @close="removeScope(scope)"
          >
            {{ scope }}
          </el-tag>
          <el-input
            v-if="inputVisible"
            ref="inputRef"
            v-model="inputValue"
            size="small"
            style="width: 120px"
            @keyup.enter="handleInputConfirm"
            @blur="handleInputConfirm"
          />
          <el-button
            v-else
            size="small"
            @click="showInput"
          >
            + 添加权限范围
          </el-button>
        </div>
        <div class="form-help">
          请求的OAuth权限范围，如：user:email, read:user
        </div>
      </el-form-item>

      <!-- 提供商特定配置 -->
      <div v-if="formData.type" class="provider-specific-config">
        <h4>{{ getProviderName(formData.type) }} 特定配置</h4>
        
        <!-- GitHub 配置 -->
        <template v-if="formData.type === 'github'">
          <el-form-item label="GitHub URL">
            <el-input
              v-model="formData.config.github_url"
              placeholder="https://github.com"
            />
            <div class="form-help">
              GitHub服务器地址，企业版可自定义
            </div>
          </el-form-item>
        </template>

        <!-- Google 配置 -->
        <template v-if="formData.type === 'google'">
          <el-form-item label="Hosted Domain">
            <el-input
              v-model="formData.config.hosted_domain"
              placeholder="example.com"
            />
            <div class="form-help">
              限制登录的Google Workspace域名（可选）
            </div>
          </el-form-item>
        </template>

        <!-- Microsoft 配置 -->
        <template v-if="formData.type === 'microsoft'">
          <el-form-item label="Tenant ID">
            <el-input
              v-model="formData.config.tenant_id"
              placeholder="请输入Tenant ID"
            />
            <div class="form-help">
              Azure AD租户ID
            </div>
          </el-form-item>
        </template>

        <!-- GitLab 配置 -->
        <template v-if="formData.type === 'gitlab'">
          <el-form-item label="GitLab URL">
            <el-input
              v-model="formData.config.gitlab_url"
              placeholder="https://gitlab.com"
            />
            <div class="form-help">
              GitLab服务器地址
            </div>
          </el-form-item>
        </template>
      </div>

      <el-form-item label="启用状态">
        <el-switch
          v-model="formData.enabled"
          active-text="启用"
          inactive-text="禁用"
        />
      </el-form-item>
    </el-form>

    <!-- 配置预览 -->
    <div v-if="formData.type" class="config-preview">
      <h4>配置预览</h4>
      <el-descriptions :column="1" size="small" border>
        <el-descriptions-item label="授权URL">
          {{ generateAuthUrl() }}
        </el-descriptions-item>
        <el-descriptions-item label="令牌URL">
          {{ generateTokenUrl() }}
        </el-descriptions-item>
        <el-descriptions-item label="用户信息URL">
          {{ generateUserInfoUrl() }}
        </el-descriptions-item>
      </el-descriptions>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button
          type="primary"
          @click="handleSubmit"
          :loading="loading"
        >
          {{ isEditing ? '更新' : '添加' }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import { ref, reactive, computed, watch, nextTick } from "vue"
import { ElMessage, type FormInstance, type FormRules } from "element-plus"

interface OAuthProvider {
  id?: string
  type: string
  name: string
  client_id: string
  client_secret: string
  redirect_uri: string
  scopes: string[]
  enabled: boolean
  config: Record<string, any>
}

interface Props {
  modelValue: boolean
  provider?: OAuthProvider | null
}

const props = defineProps<Props>()
const emit = defineEmits<{
  "update:modelValue": [value: boolean]
  save: [provider: OAuthProvider]
}>()

const formRef = ref<FormInstance>()
const inputRef = ref()
const loading = ref(false)
const dialogVisible = ref(false)
const inputVisible = ref(false)
const inputValue = ref("")

const formData = reactive<OAuthProvider>({
  type: "",
  name: "",
  client_id: "",
  client_secret: "",
  redirect_uri: "",
  scopes: [],
  enabled: true,
  config: {}
})

const providerTypes = [
  { label: "GitHub", value: "github", icon: "Github" },
  { label: "Google", value: "google", icon: "Google" },
  { label: "Microsoft", value: "microsoft", icon: "Microsoft" },
  { label: "GitLab", value: "gitlab", icon: "Gitlab" }
]

const formRules: FormRules = {
  type: [
    { required: true, message: "请选择提供商类型", trigger: "change" }
  ],
  name: [
    { required: true, message: "请输入提供商名称", trigger: "blur" },
    { min: 2, max: 100, message: "名称长度在 2 到 100 个字符", trigger: "blur" }
  ],
  client_id: [
    { required: true, message: "请输入Client ID", trigger: "blur" }
  ],
  client_secret: [
    { required: true, message: "请输入Client Secret", trigger: "blur" }
  ],
  redirect_uri: [
    { required: true, message: "请输入回调URL", trigger: "blur" },
    { type: "url", message: "请输入有效的URL", trigger: "blur" }
  ]
}

const isEditing = computed(() => !!props.provider?.id)

/** 获取提供商名称 */
const getProviderName = (type: string) => {
  const provider = providerTypes.find(p => p.value === type)
  return provider?.label || type
}

/** 监听对话框显示状态 */
watch(() => props.modelValue, (val) => {
  dialogVisible.value = val
  if (val && props.provider) {
    // 编辑模式，填充数据
    Object.assign(formData, {
      ...props.provider,
      config: { ...props.provider.config }
    })
  }
})

watch(dialogVisible, (val) => {
  emit("update:modelValue", val)
})

/** 提供商类型变更处理 */
const handleTypeChange = (type: string) => {
  // 重置配置
  formData.config = {}
  
  // 设置默认值
  const defaults = getProviderDefaults(type)
  Object.assign(formData, defaults)
}

/** 获取提供商默认配置 */
const getProviderDefaults = (type: string) => {
  const baseUrl = window.location.origin
  
  const defaults: Record<string, any> = {
    github: {
      name: "GitHub OAuth",
      redirect_uri: `${baseUrl}/oauth/callback/github`,
      scopes: ["user:email", "read:user"],
      config: { github_url: "https://github.com" }
    },
    google: {
      name: "Google OAuth",
      redirect_uri: `${baseUrl}/oauth/callback/google`,
      scopes: ["openid", "email", "profile"],
      config: {}
    },
    microsoft: {
      name: "Microsoft OAuth",
      redirect_uri: `${baseUrl}/oauth/callback/microsoft`,
      scopes: ["openid", "email", "profile"],
      config: {}
    },
    gitlab: {
      name: "GitLab OAuth",
      redirect_uri: `${baseUrl}/oauth/callback/gitlab`,
      scopes: ["read_user", "read_repository"],
      config: { gitlab_url: "https://gitlab.com" }
    }
  }
  
  return defaults[type] || {}
}

/** 生成授权URL */
const generateAuthUrl = () => {
  const urls: Record<string, string> = {
    github: `${formData.config.github_url || 'https://github.com'}/login/oauth/authorize`,
    google: "https://accounts.google.com/oauth2/auth",
    microsoft: `https://login.microsoftonline.com/${formData.config.tenant_id || 'common'}/oauth2/v2.0/authorize`,
    gitlab: `${formData.config.gitlab_url || 'https://gitlab.com'}/oauth/authorize`
  }
  return urls[formData.type] || ""
}

/** 生成令牌URL */
const generateTokenUrl = () => {
  const urls: Record<string, string> = {
    github: `${formData.config.github_url || 'https://github.com'}/login/oauth/access_token`,
    google: "https://oauth2.googleapis.com/token",
    microsoft: `https://login.microsoftonline.com/${formData.config.tenant_id || 'common'}/oauth2/v2.0/token`,
    gitlab: `${formData.config.gitlab_url || 'https://gitlab.com'}/oauth/token`
  }
  return urls[formData.type] || ""
}

/** 生成用户信息URL */
const generateUserInfoUrl = () => {
  const urls: Record<string, string> = {
    github: "https://api.github.com/user",
    google: "https://www.googleapis.com/oauth2/v2/userinfo",
    microsoft: "https://graph.microsoft.com/v1.0/me",
    gitlab: `${formData.config.gitlab_url || 'https://gitlab.com'}/api/v4/user`
  }
  return urls[formData.type] || ""
}

/** 移除权限范围 */
const removeScope = (scope: string) => {
  const index = formData.scopes.indexOf(scope)
  if (index > -1) {
    formData.scopes.splice(index, 1)
  }
}

/** 显示输入框 */
const showInput = () => {
  inputVisible.value = true
  nextTick(() => {
    inputRef.value?.focus()
  })
}

/** 确认输入 */
const handleInputConfirm = () => {
  if (inputValue.value && !formData.scopes.includes(inputValue.value)) {
    formData.scopes.push(inputValue.value)
  }
  inputVisible.value = false
  inputValue.value = ""
}

/** 提交表单 */
const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    
    loading.value = true
    
    // 构建提供商数据
    const providerData: OAuthProvider = {
      ...formData,
      scopes: [...formData.scopes],
      config: { ...formData.config }
    }
    
    emit("save", providerData)
    handleClose()
    
  } catch (error) {
    ElMessage.error("请检查表单输入")
  } finally {
    loading.value = false
  }
}

/** 关闭对话框 */
const handleClose = () => {
  // 重置表单
  formRef.value?.resetFields()
  Object.assign(formData, {
    type: "",
    name: "",
    client_id: "",
    client_secret: "",
    redirect_uri: "",
    scopes: [],
    enabled: true,
    config: {}
  })
  
  inputVisible.value = false
  inputValue.value = ""
  
  dialogVisible.value = false
}
</script>

<style lang="scss" scoped>
.provider-option {
  display: flex;
  align-items: center;
  gap: 8px;
  
  .provider-icon {
    font-size: 16px;
  }
}

.form-help {
  margin-top: 4px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
  line-height: 1.4;
}

.scopes-container {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}

.provider-specific-config {
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid var(--el-border-color-lighter);
  
  h4 {
    margin: 0 0 16px 0;
    color: var(--el-text-color-primary);
    font-size: 14px;
    font-weight: 600;
  }
}

.config-preview {
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid var(--el-border-color-lighter);
  
  h4 {
    margin: 0 0 16px 0;
    color: var(--el-text-color-primary);
    font-size: 14px;
    font-weight: 600;
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>