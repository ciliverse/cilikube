<template>
  <el-dialog
    v-model="dialogVisible"
    title="更换头像"
    width="500px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <div class="avatar-upload-container">
      <!-- 当前头像预览 -->
      <div class="current-avatar">
        <h4>当前头像</h4>
        <el-avatar :size="80" :src="currentAvatar">
          <el-icon><User /></el-icon>
        </el-avatar>
      </div>

      <!-- 上传区域 -->
      <div class="upload-section">
        <h4>选择新头像</h4>
        
        <!-- 文件上传 -->
        <el-upload
          ref="uploadRef"
          class="avatar-uploader"
          :action="uploadAction"
          :headers="uploadHeaders"
          :show-file-list="false"
          :before-upload="beforeUpload"
          :on-success="handleUploadSuccess"
          :on-error="handleUploadError"
          :on-progress="handleUploadProgress"
          name="avatar"
          accept="image/*"
          drag
        >
          <div v-if="!previewUrl" class="upload-placeholder">
            <el-icon class="upload-icon"><Plus /></el-icon>
            <div class="upload-text">
              <p>点击或拖拽图片到此处上传</p>
              <p class="upload-hint">支持 JPG、PNG、GIF 格式，文件大小不超过 2MB</p>
            </div>
          </div>
          
          <div v-else class="preview-container">
            <img :src="previewUrl" class="preview-image" alt="预览" />
            <div class="preview-overlay">
              <el-icon><Edit /></el-icon>
              <span>点击更换</span>
            </div>
          </div>
        </el-upload>

        <!-- 上传进度 -->
        <div v-if="uploading" class="upload-progress">
          <el-progress 
            :percentage="uploadProgress" 
            :status="uploadStatus"
            :stroke-width="6"
          />
          <p class="progress-text">{{ uploadProgressText }}</p>
        </div>
      </div>

      <!-- 图片裁剪 -->
      <div v-if="previewUrl && !uploading" class="crop-section">
        <h4>调整头像</h4>
        <div class="crop-container">
          <div class="crop-preview">
            <div class="crop-preview-item">
              <span>大头像 (120x120)</span>
              <div class="preview-avatar large">
                <img :src="previewUrl" alt="大头像预览" />
              </div>
            </div>
            <div class="crop-preview-item">
              <span>中头像 (80x80)</span>
              <div class="preview-avatar medium">
                <img :src="previewUrl" alt="中头像预览" />
              </div>
            </div>
            <div class="crop-preview-item">
              <span>小头像 (40x40)</span>
              <div class="preview-avatar small">
                <img :src="previewUrl" alt="小头像预览" />
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 颜色头像选择 -->
      <div class="color-avatars">
        <h4>或选择颜色头像</h4>
        <div class="color-grid">
          <div 
            v-for="(color, index) in colorAvatars" 
            :key="index"
            class="color-item"
            :class="{ active: selectedColor === color }"
            @click="selectColorAvatar(color)"
          >
            <div 
              class="color-avatar"
              :style="{ backgroundColor: color.bg, color: color.text }"
            >
              {{ getInitials(props.currentUser?.username || 'U') }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button 
          type="primary" 
          @click="confirmUpload"
          :loading="confirming"
          :disabled="!previewUrl && !selectedPreset && !selectedColor"
        >
          确认更换
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import { ref, computed, watch } from "vue"
import { ElMessage, type UploadInstance } from "element-plus"
import { User, Plus, Edit } from "@element-plus/icons-vue"
import { getToken } from "@/utils/cache/cookies"
import { updateAvatarApi } from "@/api/profile"

interface Props {
  modelValue: boolean
  currentAvatar?: string
  currentUser?: {
    username: string
    display_name?: string
    email?: string
  }
}

const props = withDefaults(defineProps<Props>(), {
  currentAvatar: "",
  currentUser: () => ({ username: "User" })
})

const emit = defineEmits<{
  "update:modelValue": [value: boolean]
  uploaded: [avatarUrl: string]
}>()

const uploadRef = ref<UploadInstance>()
const dialogVisible = ref(false)
const previewUrl = ref("")
const selectedPreset = ref("")
const uploading = ref(false)
const uploadProgress = ref(0)
const uploadStatus = ref<"" | "success" | "exception" | "warning">("")
const uploadProgressText = ref("")
const confirming = ref(false)

// 上传配置
const uploadAction = "/api/v1/profile/avatar"
const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${getToken()}`
}))

// 颜色头像配置
const colorAvatars = [
  { bg: "#1890ff", text: "#ffffff" }, // 蓝色
  { bg: "#52c41a", text: "#ffffff" }, // 绿色
  { bg: "#fa541c", text: "#ffffff" }, // 橙色
  { bg: "#722ed1", text: "#ffffff" }, // 紫色
  { bg: "#eb2f96", text: "#ffffff" }, // 粉色
  { bg: "#13c2c2", text: "#ffffff" }, // 青色
  { bg: "#faad14", text: "#ffffff" }, // 黄色
  { bg: "#f5222d", text: "#ffffff" }, // 红色
  { bg: "#2f54eb", text: "#ffffff" }, // 靛蓝
  { bg: "#a0d911", text: "#ffffff" }, // 青绿
  { bg: "#fa8c16", text: "#ffffff" }, // 金色
  { bg: "#531dab", text: "#ffffff" }  // 深紫
]

const selectedColor = ref<{ bg: string; text: string } | null>(null)

/** 监听对话框显示状态 */
watch(() => props.modelValue, (val) => {
  dialogVisible.value = val
})

watch(dialogVisible, (val) => {
  emit("update:modelValue", val)
})

/** 上传前验证 */
const beforeUpload = (file: File) => {
  // 检查文件类型
  const isImage = file.type.startsWith("image/")
  if (!isImage) {
    ElMessage.error("只能上传图片文件！")
    return false
  }

  // 检查文件大小
  const isLt2M = file.size / 1024 / 1024 < 2
  if (!isLt2M) {
    ElMessage.error("图片大小不能超过 2MB！")
    return false
  }

  // 生成预览URL
  previewUrl.value = URL.createObjectURL(file)
  selectedPreset.value = ""
  
  uploading.value = true
  uploadProgress.value = 0
  uploadStatus.value = ""
  uploadProgressText.value = "准备上传..."

  return true
}

/** 上传进度 */
const handleUploadProgress = (event: any) => {
  uploadProgress.value = Math.round((event.loaded / event.total) * 100)
  uploadProgressText.value = `上传中... ${uploadProgress.value}%`
}

/** 上传成功 */
const handleUploadSuccess = (response: any) => {
  uploading.value = false
  uploadStatus.value = "success"
  uploadProgressText.value = "上传成功！"
  
  if (response.code === 200) {
    // 构建完整的头像URL用于预览
    const baseUrl = `${window.location.protocol}//${window.location.hostname}:8080`
    const fullAvatarUrl = `${baseUrl}${response.data.avatar_url}`
    previewUrl.value = fullAvatarUrl
    
    // 通知父组件更新头像 - 传递后端返回的相对路径
    emit("uploaded", response.data.avatar_url)
    handleClose()
    ElMessage.success("头像上传成功")
  } else {
    ElMessage.error(response.message || "上传失败")
  }
}

/** 上传失败 */
const handleUploadError = (error: any) => {
  uploading.value = false
  uploadStatus.value = "exception"
  uploadProgressText.value = "上传失败"
  
  console.error("头像上传失败:", error)
  ElMessage.error("头像上传失败，请重试")
}

/** 获取用户名首字母 */
const getInitials = (username: string) => {
  if (!username) return 'U'
  
  // 如果是中文名，取前两个字符
  if (/[\u4e00-\u9fa5]/.test(username)) {
    return username.slice(0, 2).toUpperCase()
  }
  
  // 如果是英文名，取首字母
  const words = username.split(/\s+/)
  if (words.length >= 2) {
    return (words[0][0] + words[1][0]).toUpperCase()
  }
  
  return username.slice(0, 2).toUpperCase()
}

/** 选择颜色头像 */
const selectColorAvatar = (color: { bg: string; text: string }) => {
  selectedColor.value = color
  previewUrl.value = ""
  selectedPreset.value = ""
  
  // 清理上传状态
  uploading.value = false
  uploadProgress.value = 0
}

/** 选择预设头像 */
const selectPresetAvatar = (avatar: string) => {
  selectedPreset.value = avatar
  previewUrl.value = ""
  selectedColor.value = null
  
  // 清理上传状态
  uploading.value = false
  uploadProgress.value = 0
}

/** 生成颜色头像数据URL */
const generateColorAvatarDataUrl = (color: { bg: string; text: string }, initials: string) => {
  const canvas = document.createElement('canvas')
  const ctx = canvas.getContext('2d')
  const size = 120
  
  canvas.width = size
  canvas.height = size
  
  if (ctx) {
    // 绘制背景
    ctx.fillStyle = color.bg
    ctx.fillRect(0, 0, size, size)
    
    // 绘制文字
    ctx.fillStyle = color.text
    ctx.font = `${size * 0.4}px Arial, sans-serif`
    ctx.textAlign = 'center'
    ctx.textBaseline = 'middle'
    ctx.fillText(initials, size / 2, size / 2)
  }
  
  return canvas.toDataURL('image/png')
}

/** 确认更换头像 */
const confirmUpload = async () => {
  try {
    confirming.value = true
    
    let avatarUrl = ""
    
    if (previewUrl.value) {
      // 使用上传的图片
      avatarUrl = previewUrl.value
    } else if (selectedColor.value) {
      // 使用颜色头像
      const initials = getInitials(props.currentUser?.username || 'U')
      const colorAvatarDataUrl = generateColorAvatarDataUrl(selectedColor.value, initials)
      
      // 将颜色头像信息发送到后端
      const { data } = await updateAvatarApi({ 
        email: props.currentUser?.email || '',
        display_name: props.currentUser?.display_name || '',
        avatar_url: colorAvatarDataUrl
      })
      avatarUrl = data.avatar_url
    } else if (selectedPreset.value) {
      // 使用预设头像
      const { data } = await updateAvatarApi({ 
        email: props.currentUser?.email || '',
        display_name: props.currentUser?.display_name || '',
        avatar_url: selectedPreset.value 
      })
      avatarUrl = data.avatar_url
    }
    
    if (avatarUrl) {
      emit("uploaded", avatarUrl)
      handleClose()
      ElMessage.success("头像更换成功")
    }
    
  } catch (error: any) {
    console.error("更换头像失败:", error)
    ElMessage.error(error.response?.data?.message || "更换头像失败")
  } finally {
    confirming.value = false
  }
}

/** 关闭对话框 */
const handleClose = () => {
  // 清理状态
  previewUrl.value = ""
  selectedPreset.value = ""
  uploading.value = false
  uploadProgress.value = 0
  uploadStatus.value = ""
  uploadProgressText.value = ""
  confirming.value = false
  
  // 清理预览URL
  if (previewUrl.value && previewUrl.value.startsWith("blob:")) {
    URL.revokeObjectURL(previewUrl.value)
  }
  
  dialogVisible.value = false
}
</script>

<style lang="scss" scoped>
.avatar-upload-container {
  .current-avatar {
    text-align: center;
    margin-bottom: 24px;
    
    h4 {
      margin: 0 0 12px 0;
      font-size: 14px;
      color: var(--el-text-color-regular);
    }
  }
  
  .upload-section {
    margin-bottom: 24px;
    
    h4 {
      margin: 0 0 12px 0;
      font-size: 14px;
      color: var(--el-text-color-regular);
    }
    
    .avatar-uploader {
      :deep(.el-upload) {
        border: 2px dashed var(--el-border-color);
        border-radius: 8px;
        cursor: pointer;
        position: relative;
        overflow: hidden;
        transition: all 0.3s ease;
        
        &:hover {
          border-color: var(--el-color-primary);
        }
      }
      
      :deep(.el-upload-dragger) {
        background: transparent;
        border: none;
        width: 100%;
        height: 200px;
        display: flex;
        align-items: center;
        justify-content: center;
      }
      
      .upload-placeholder {
        text-align: center;
        
        .upload-icon {
          font-size: 48px;
          color: var(--el-text-color-placeholder);
          margin-bottom: 16px;
        }
        
        .upload-text {
          p {
            margin: 8px 0;
            color: var(--el-text-color-regular);
          }
          
          .upload-hint {
            font-size: 12px;
            color: var(--el-text-color-secondary);
          }
        }
      }
      
      .preview-container {
        position: relative;
        width: 100%;
        height: 200px;
        
        .preview-image {
          width: 100%;
          height: 100%;
          object-fit: cover;
        }
        
        .preview-overlay {
          position: absolute;
          top: 0;
          left: 0;
          right: 0;
          bottom: 0;
          background: rgba(0, 0, 0, 0.6);
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          opacity: 0;
          transition: opacity 0.3s ease;
          color: white;
          
          .el-icon {
            font-size: 24px;
            margin-bottom: 8px;
          }
        }
        
        &:hover .preview-overlay {
          opacity: 1;
        }
      }
    }
    
    .upload-progress {
      margin-top: 16px;
      
      .progress-text {
        margin: 8px 0 0 0;
        font-size: 14px;
        color: var(--el-text-color-regular);
        text-align: center;
      }
    }
  }
  
  .crop-section {
    margin-bottom: 24px;
    
    h4 {
      margin: 0 0 12px 0;
      font-size: 14px;
      color: var(--el-text-color-regular);
    }
    
    .crop-preview {
      display: flex;
      justify-content: space-around;
      align-items: center;
      padding: 20px;
      background: var(--el-fill-color-extra-light);
      border-radius: 8px;
      
      .crop-preview-item {
        text-align: center;
        
        span {
          display: block;
          margin-bottom: 8px;
          font-size: 12px;
          color: var(--el-text-color-secondary);
        }
        
        .preview-avatar {
          border-radius: 50%;
          overflow: hidden;
          border: 2px solid var(--el-border-color-light);
          
          &.large {
            width: 60px;
            height: 60px;
          }
          
          &.medium {
            width: 40px;
            height: 40px;
          }
          
          &.small {
            width: 24px;
            height: 24px;
          }
          
          img {
            width: 100%;
            height: 100%;
            object-fit: cover;
          }
        }
      }
    }
  }
  
  .color-avatars {
    h4 {
      margin: 0 0 12px 0;
      font-size: 14px;
      color: var(--el-text-color-regular);
    }
    
    .color-grid {
      display: grid;
      grid-template-columns: repeat(6, 1fr);
      gap: 12px;
      
      .color-item {
        cursor: pointer;
        transition: all 0.3s ease;
        
        &:hover {
          transform: scale(1.05);
        }
        
        &.active {
          .color-avatar {
            box-shadow: 0 0 0 3px var(--el-color-primary-light-8);
            border: 2px solid var(--el-color-primary);
          }
        }
        
        .color-avatar {
          width: 50px;
          height: 50px;
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          font-size: 16px;
          font-weight: 600;
          border: 2px solid transparent;
          transition: all 0.3s ease;
        }
      }
    }
  }
  
  .preset-avatars {
    h4 {
      margin: 0 0 12px 0;
      font-size: 14px;
      color: var(--el-text-color-regular);
    }
    
    .preset-grid {
      display: grid;
      grid-template-columns: repeat(4, 1fr);
      gap: 12px;
      
      .preset-item {
        width: 60px;
        height: 60px;
        border-radius: 50%;
        overflow: hidden;
        border: 2px solid var(--el-border-color-light);
        cursor: pointer;
        transition: all 0.3s ease;
        
        &:hover {
          border-color: var(--el-color-primary-light-5);
          transform: scale(1.05);
        }
        
        &.active {
          border-color: var(--el-color-primary);
          box-shadow: 0 0 0 2px var(--el-color-primary-light-8);
        }
        
        img {
          width: 100%;
          height: 100%;
          object-fit: cover;
        }
      }
    }
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

@media (max-width: 768px) {
  .crop-preview {
    flex-direction: column;
    gap: 16px;
    
    .crop-preview-item {
      display: flex;
      align-items: center;
      gap: 12px;
      
      span {
        margin-bottom: 0;
      }
    }
  }
  
  .preset-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}
</style>