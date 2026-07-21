<template>
  <el-dialog
    v-model="dialogVisible"
    title="添加用户"
    width="600px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-width="100px"
    >
      <el-form-item label="用户名" prop="username">
        <el-input
          v-model="formData.username"
          placeholder="请输入用户名"
          maxlength="50"
          show-word-limit
        />
        <div class="form-help">
          用户名用于登录，只能包含字母、数字和下划线
        </div>
      </el-form-item>

      <el-form-item label="邮箱地址" prop="email">
        <el-input
          v-model="formData.email"
          type="email"
          placeholder="请输入邮箱地址"
        />
        <div class="form-help">
          邮箱地址用于接收系统通知和密码重置
        </div>
      </el-form-item>

      <el-form-item label="显示名称" prop="display_name">
        <el-input
          v-model="formData.display_name"
          placeholder="请输入显示名称（可选）"
          maxlength="100"
          show-word-limit
        />
      </el-form-item>

      <el-form-item label="初始密码" prop="password">
        <el-input
          v-model="formData.password"
          type="password"
          placeholder="请输入初始密码"
          show-password
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

      <el-form-item label="确认密码" prop="confirmPassword">
        <el-input
          v-model="formData.confirmPassword"
          type="password"
          placeholder="请再次输入密码"
          show-password
        />
      </el-form-item>

      <el-form-item label="用户角色" prop="role">
        <el-select
          v-model="formData.role"
          placeholder="请选择用户角色"
          style="width: 100%"
        >
          <el-option
            v-for="role in roleOptions"
            :key="role.value"
            :label="role.label"
            :value="role.value"
          >
            <div class="role-option">
              <span class="role-name">{{ role.label }}</span>
              <span class="role-desc">{{ role.description }}</span>
            </div>
          </el-option>
        </el-select>
      </el-form-item>

      <el-form-item label="账户状态">
        <el-switch
          v-model="formData.is_active"
          active-text="激活"
          inactive-text="停用"
        />
        <div class="form-help">
          停用的账户无法登录系统
        </div>
      </el-form-item>

      <el-form-item label="发送通知">
        <el-checkbox v-model="formData.send_welcome_email">
          发送欢迎邮件给新用户
        </el-checkbox>
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button
          type="primary"
          @click="handleSubmit"
          :loading="loading"
        >
          创建用户
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import { ref, reactive, computed, watch } from "vue"
import { ElMessage, type FormInstance, type FormRules } from "element-plus"
import { useI18n } from "vue-i18n"
import { createUserApi } from "@/api/admin/user-management"

interface Props {
  modelValue: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
  "update:modelValue": [value: boolean]
  created: []
}>()

const { t } = useI18n()
const formRef = ref<FormInstance>()
const loading = ref(false)
const dialogVisible = ref(false)

const formData = reactive({
  username: "",
  email: "",
  display_name: "",
  password: "",
  confirmPassword: "",
  role: "viewer",
  is_active: true,
  send_welcome_email: true
})

const roleOptions = computed(() => [
  {
    value: "admin",
    label: t('admin.userManagement.roles.admin'),
    description: t('admin.userManagement.roleDescriptions.admin')
  },
  {
    value: "editor",
    label: t('admin.userManagement.roles.editor'),
    description: t('admin.userManagement.roleDescriptions.editor')
  },
  {
    value: "viewer",
    label: t('admin.userManagement.roles.viewer'),
    description: t('admin.userManagement.roleDescriptions.viewer')
  }
])

const formRules: FormRules = {
  username: [
    { required: true, message: "请输入用户名", trigger: "blur" },
    { min: 3, max: 50, message: "用户名长度在 3 到 50 个字符", trigger: "blur" },
    { 
      pattern: /^[a-zA-Z0-9_]+$/, 
      message: "用户名只能包含字母、数字和下划线", 
      trigger: "blur" 
    }
  ],
  email: [
    { required: true, message: "请输入邮箱地址", trigger: "blur" },
    { type: "email", message: "请输入正确的邮箱格式", trigger: "blur" }
  ],
  display_name: [
    { max: 100, message: "显示名称不能超过100个字符", trigger: "blur" }
  ],
  password: [
    { required: true, message: "请输入初始密码", trigger: "blur" },
    { min: 6, message: "密码长度至少6个字符", trigger: "blur" }
  ],
  confirmPassword: [
    { required: true, message: "请确认密码", trigger: "blur" },
    { 
      validator: (rule, value, callback) => {
        if (value !== formData.password) {
          callback(new Error("两次输入的密码不一致"))
        } else {
          callback()
        }
      }, 
      trigger: "blur" 
    }
  ],
  role: [
    { required: true, message: "请选择用户角色", trigger: "change" }
  ]
}

// 密码强度计算
const passwordStrength = computed(() => {
  const password = formData.password
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

/** 监听对话框显示状态 */
watch(() => props.modelValue, (val) => {
  dialogVisible.value = val
})

watch(dialogVisible, (val) => {
  emit("update:modelValue", val)
})

/** 提交表单 */
const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    
    loading.value = true
    
    await createUserApi({
      username: formData.username,
      email: formData.email,
      display_name: formData.display_name || undefined,
      password: formData.password,
      role: formData.role,
      is_active: formData.is_active,
      send_welcome_email: formData.send_welcome_email
    })
    
    ElMessage.success("用户创建成功")
    emit("created")
    handleClose()
    
  } catch (error: any) {
    console.error("创建用户失败:", error)
    ElMessage.error(error.response?.data?.message || "创建用户失败")
  } finally {
    loading.value = false
  }
}

/** 关闭对话框 */
const handleClose = () => {
  // 重置表单
  formRef.value?.resetFields()
  Object.assign(formData, {
    username: "",
    email: "",
    display_name: "",
    password: "",
    confirmPassword: "",
    role: "viewer",
    is_active: true,
    send_welcome_email: true
  })
  
  dialogVisible.value = false
}
</script>

<style lang="scss" scoped>
.form-help {
  margin-top: 4px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
  line-height: 1.4;
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

.role-option {
  display: flex;
  flex-direction: column;
  
  .role-name {
    font-weight: 500;
  }
  
  .role-desc {
    font-size: 12px;
    color: var(--el-text-color-secondary);
    margin-top: 2px;
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>