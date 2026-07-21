<template>
  <div class="profile-basic-info">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>基本信息</span>
          <el-button 
            v-if="!isEditing" 
            type="primary" 
            size="small" 
            @click="startEdit"
          >
            编辑
          </el-button>
        </div>
      </template>

      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
        :disabled="!isEditing"
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

        <el-form-item label="显示名称" prop="display_name">
          <el-input 
            v-model="formData.display_name" 
            placeholder="请输入显示名称"
            maxlength="50"
            show-word-limit
          />
        </el-form-item>

        <el-form-item label="邮箱地址" prop="email">
          <el-input 
            v-model="formData.email" 
            placeholder="请输入邮箱地址"
            type="email"
          >
            <template #suffix>
              <el-tag 
                v-if="userInfo.email_verified" 
                type="success" 
                size="small"
              >
                已验证
              </el-tag>
              <el-tag 
                v-else 
                type="warning" 
                size="small"
              >
                未验证
              </el-tag>
            </template>
          </el-input>
          <div v-if="!userInfo.email_verified" class="email-verify-hint">
            <el-button 
              type="text" 
              size="small" 
              @click="sendVerificationEmail"
              :loading="verifyEmailLoading"
            >
              发送验证邮件
            </el-button>
          </div>
        </el-form-item>

        <el-form-item label="个人简介" prop="bio">
          <el-input 
            v-model="formData.bio" 
            type="textarea" 
            :rows="4"
            placeholder="介绍一下自己..."
            maxlength="200"
            show-word-limit
          />
        </el-form-item>

        <el-form-item label="所在地区" prop="location">
          <el-input 
            v-model="formData.location" 
            placeholder="请输入所在地区"
            maxlength="100"
          />
        </el-form-item>

        <el-form-item label="个人网站" prop="website">
          <el-input 
            v-model="formData.website" 
            placeholder="https://example.com"
            type="url"
          />
        </el-form-item>

        <el-form-item v-if="isEditing">
          <el-button 
            type="primary" 
            @click="saveChanges"
            :loading="saveLoading"
          >
            保存更改
          </el-button>
          <el-button @click="cancelEdit">
            取消
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, watch } from "vue"
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from "element-plus"
import { Lock } from "@element-plus/icons-vue"
import { updateUserProfileApi, sendEmailVerificationApi } from "@/api/profile"

interface UserInfo {
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
  bio?: string
  location?: string
  website?: string
}

interface Props {
  userInfo: UserInfo
}

const props = defineProps<Props>()
const emit = defineEmits<{
  update: [userInfo: Partial<UserInfo>]
  loading: [loading: boolean]
}>()

const formRef = ref<FormInstance>()
const isEditing = ref(false)
const saveLoading = ref(false)
const verifyEmailLoading = ref(false)

const formData = reactive({
  username: "",
  display_name: "",
  email: "",
  bio: "",
  location: "",
  website: ""
})

const originalData = reactive({
  username: "",
  display_name: "",
  email: "",
  bio: "",
  location: "",
  website: ""
})

const formRules: FormRules = {
  display_name: [
    { max: 50, message: "显示名称不能超过50个字符", trigger: "blur" }
  ],
  email: [
    { required: true, message: "请输入邮箱地址", trigger: "blur" },
    { type: "email", message: "请输入正确的邮箱格式", trigger: "blur" }
  ],
  bio: [
    { max: 200, message: "个人简介不能超过200个字符", trigger: "blur" }
  ],
  location: [
    { max: 100, message: "所在地区不能超过100个字符", trigger: "blur" }
  ],
  website: [
    { type: "url", message: "请输入正确的网址格式", trigger: "blur" }
  ]
}

/** 更新表单数据 */
const updateFormData = (userInfo: UserInfo) => {
  formData.username = userInfo.username || ""
  formData.display_name = userInfo.display_name || ""
  formData.email = userInfo.email || ""
  formData.bio = userInfo.bio || ""
  formData.location = userInfo.location || ""
  formData.website = userInfo.website || ""
  
  // 保存原始数据
  Object.assign(originalData, formData)
}

/** 监听用户信息变化，更新表单数据 */
watch(() => props.userInfo, (newUserInfo) => {
  updateFormData(newUserInfo)
}, { immediate: true, deep: true })

/** 开始编辑 */
const startEdit = () => {
  isEditing.value = true
  // 重新保存原始数据
  Object.assign(originalData, formData)
}

/** 取消编辑 */
const cancelEdit = async () => {
  // 检查是否有未保存的更改
  const hasChanges = Object.keys(formData).some(
    key => formData[key as keyof typeof formData] !== originalData[key as keyof typeof originalData]
  )
  
  if (hasChanges) {
    try {
      await ElMessageBox.confirm(
        "您有未保存的更改，确定要取消吗？",
        "确认取消",
        {
          confirmButtonText: "确定",
          cancelButtonText: "继续编辑",
          type: "warning"
        }
      )
    } catch {
      return // 用户选择继续编辑
    }
  }
  
  // 恢复原始数据
  Object.assign(formData, originalData)
  isEditing.value = false
}

/** 保存更改 */
const saveChanges = async () => {
  if (!formRef.value) return
  
  try {
    // 表单验证
    await formRef.value.validate()
    
    saveLoading.value = true
    emit("loading", true)
    
    // 准备更新数据
    const updateData = {
      display_name: formData.display_name,
      email: formData.email,
      bio: formData.bio,
      location: formData.location,
      website: formData.website
    }
    
    // 调用API更新用户信息
    const { data } = await updateUserProfileApi(updateData)
    
    // 更新本地数据
    emit("update", data)
    Object.assign(originalData, formData)
    
    isEditing.value = false
    ElMessage.success("个人信息更新成功")
    
  } catch (error: any) {
    console.error("更新个人信息失败:", error)
    ElMessage.error(error.response?.data?.message || "更新失败，请重试")
  } finally {
    saveLoading.value = false
    emit("loading", false)
  }
}

/** 发送邮箱验证邮件 */
const sendVerificationEmail = async () => {
  try {
    verifyEmailLoading.value = true
    
    await sendEmailVerificationApi()
    
    ElMessage.success("验证邮件已发送，请检查您的邮箱")
  } catch (error: any) {
    console.error("发送验证邮件失败:", error)
    ElMessage.error(error.response?.data?.message || "发送失败，请重试")
  } finally {
    verifyEmailLoading.value = false
  }
}
</script>

<style lang="scss" scoped>
.profile-basic-info {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    
    span {
      font-weight: 600;
      font-size: 16px;
    }
  }
  
  .email-verify-hint {
    margin-top: 4px;
    font-size: 12px;
    color: var(--el-text-color-secondary);
  }
  
  :deep(.el-form-item__label) {
    font-weight: 500;
  }
  
  :deep(.el-input.is-disabled .el-input__wrapper) {
    background-color: var(--el-fill-color-light);
  }
}
</style>