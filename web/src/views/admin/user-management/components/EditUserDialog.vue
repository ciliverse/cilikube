<template>
  <el-dialog
    v-model="dialogVisible"
    title="编辑用户"
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
          :disabled="true"
        >
          <template #suffix>
            <el-tooltip content="用户名不可修改" placement="top">
              <el-icon><Lock /></el-icon>
            </el-tooltip>
          </template>
        </el-input>
      </el-form-item>

      <el-form-item label="邮箱地址" prop="email">
        <el-input
          v-model="formData.email"
          type="email"
          placeholder="请输入邮箱地址"
        />
        <div class="form-help">
          修改邮箱后需要重新验证
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

      <el-form-item label="邮箱验证">
        <div class="email-verification">
          <el-tag
            :type="formData.email_verified ? 'success' : 'warning'"
            size="small"
          >
            {{ formData.email_verified ? '已验证' : '未验证' }}
          </el-tag>
          <el-button
            v-if="!formData.email_verified"
            size="small"
            type="primary"
            link
            @click="sendVerificationEmail"
            :loading="sendingVerification"
          >
            发送验证邮件
          </el-button>
        </div>
      </el-form-item>

      <el-form-item label="密码重置">
        <el-button
          size="small"
          @click="showResetPassword = true"
        >
          重置密码
        </el-button>
        <div class="form-help">
          重置后的新密码将通过邮件发送给用户
        </div>
      </el-form-item>

      <el-form-item label="最后登录">
        <span class="info-text">{{ formatDate(user?.last_login) }}</span>
      </el-form-item>

      <el-form-item label="创建时间">
        <span class="info-text">{{ formatDate(user?.created_at) }}</span>
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
          保存更改
        </el-button>
      </div>
    </template>

    <!-- 重置密码确认对话框 -->
    <el-dialog
      v-model="showResetPassword"
      title="重置密码"
      width="400px"
      :close-on-click-modal="false"
    >
      <div class="reset-password-content">
        <el-alert
          title="重置密码确认"
          type="warning"
          :closable="false"
          show-icon
        >
          <template #default>
            <p>确定要重置用户 <strong>{{ user?.username }}</strong> 的密码吗？</p>
            <p>新密码将通过邮件发送到：<strong>{{ user?.email }}</strong></p>
          </template>
        </el-alert>
      </div>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="showResetPassword = false">取消</el-button>
          <el-button
            type="primary"
            @click="handleResetPassword"
            :loading="resettingPassword"
          >
            确认重置
          </el-button>
        </div>
      </template>
    </el-dialog>
  </el-dialog>
</template>

<script lang="ts" setup>
import { ref, reactive, watch, computed } from "vue"
import { ElMessage, type FormInstance, type FormRules } from "element-plus"
import { Lock } from "@element-plus/icons-vue"
import { useI18n } from "vue-i18n"
import { updateUserApi, resetUserPasswordApi, sendUserVerificationApi } from "@/api/admin/user-management"

interface User {
  id: number
  username: string
  email: string
  display_name: string
  avatar_url: string
  role: string
  is_active: boolean
  email_verified: boolean
  last_login: string | null
  created_at: string
}

interface Props {
  modelValue: boolean
  user: User | null
}

const props = defineProps<Props>()
const emit = defineEmits<{
  "update:modelValue": [value: boolean]
  updated: []
}>()

const { t } = useI18n()
const formRef = ref<FormInstance>()
const loading = ref(false)
const sendingVerification = ref(false)
const resettingPassword = ref(false)
const dialogVisible = ref(false)
const showResetPassword = ref(false)

const formData = reactive({
  username: "",
  email: "",
  display_name: "",
  role: "",
  is_active: true,
  email_verified: false
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
  email: [
    { required: true, message: "请输入邮箱地址", trigger: "blur" },
    { type: "email", message: "请输入正确的邮箱格式", trigger: "blur" }
  ],
  display_name: [
    { max: 100, message: "显示名称不能超过100个字符", trigger: "blur" }
  ],
  role: [
    { required: true, message: "请选择用户角色", trigger: "change" }
  ]
}

/** 格式化日期 */
const formatDate = (dateString: string | null | undefined) => {
  if (!dateString) return "从未"
  return new Date(dateString).toLocaleString("zh-CN")
}

/** 监听对话框显示状态 */
watch(() => props.modelValue, (val) => {
  dialogVisible.value = val
})

watch(dialogVisible, (val) => {
  emit("update:modelValue", val)
})

/** 监听用户数据变化 */
watch(() => props.user, (user) => {
  if (user) {
    formData.username = user.username
    formData.email = user.email
    formData.display_name = user.display_name || ""
    formData.role = user.role
    formData.is_active = user.is_active
    formData.email_verified = user.email_verified
  }
}, { immediate: true })

/** 发送验证邮件 */
const sendVerificationEmail = async () => {
  if (!props.user) return
  
  try {
    sendingVerification.value = true
    await sendUserVerificationApi(props.user.id)
    ElMessage.success("验证邮件发送成功")
  } catch (error: any) {
    console.error("发送验证邮件失败:", error)
    ElMessage.error("发送验证邮件失败")
  } finally {
    sendingVerification.value = false
  }
}

/** 重置密码 */
const handleResetPassword = async () => {
  if (!props.user) return
  
  try {
    resettingPassword.value = true
    await resetUserPasswordApi(props.user.id)
    ElMessage.success("密码重置成功，新密码已发送到用户邮箱")
    showResetPassword.value = false
  } catch (error: any) {
    console.error("重置密码失败:", error)
    ElMessage.error("重置密码失败")
  } finally {
    resettingPassword.value = false
  }
}

/** 提交表单 */
const handleSubmit = async () => {
  if (!formRef.value || !props.user) return
  
  try {
    await formRef.value.validate()
    
    loading.value = true
    
    await updateUserApi(props.user.id, {
      email: formData.email,
      display_name: formData.display_name || undefined,
      role: formData.role,
      is_active: formData.is_active
    })
    
    ElMessage.success("用户信息更新成功")
    emit("updated")
    handleClose()
    
  } catch (error: any) {
    console.error("更新用户失败:", error)
    ElMessage.error(error.response?.data?.message || "更新用户失败")
  } finally {
    loading.value = false
  }
}

/** 关闭对话框 */
const handleClose = () => {
  formRef.value?.clearValidate()
  showResetPassword.value = false
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

.email-verification {
  display: flex;
  align-items: center;
  gap: 12px;
}

.info-text {
  color: var(--el-text-color-regular);
  font-size: 14px;
}

.reset-password-content {
  margin-bottom: 20px;
  
  :deep(.el-alert__content) {
    p {
      margin: 8px 0;
      line-height: 1.5;
    }
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

:deep(.el-input.is-disabled .el-input__wrapper) {
  background-color: var(--el-fill-color-light);
}
</style>