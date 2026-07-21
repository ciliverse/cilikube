<template>
  <div class="profile-container">
    <div class="profile-header">
      <h1>个人资料</h1>
      <p class="subtitle">管理您的账户信息和偏好设置</p>
    </div>

    <div class="profile-content">
      <el-row :gutter="24">
        <!-- 左侧：头像和基本信息 -->
        <el-col :xs="24" :sm="24" :md="8" :lg="6">
          <el-card class="profile-card">
            <div class="avatar-section">
              <div class="avatar-container">
                <el-avatar 
                  :size="120" 
                  :src="fullAvatarUrl" 
                  class="user-avatar"
                >
                  <el-icon><User /></el-icon>
                </el-avatar>
                <div class="avatar-overlay" @click="showAvatarUpload = true">
                  <el-icon><Camera /></el-icon>
                  <span>更换头像</span>
                </div>
              </div>
              
              <div class="user-basic-info">
                <h3>{{ userInfo.display_name || userInfo.username }}</h3>
                <p class="username">@{{ userInfo.username }}</p>
                <p class="email">{{ userInfo.email }}</p>
                
                <div class="user-status">
                  <el-tag 
                    :type="userInfo.is_active ? 'success' : 'danger'" 
                    size="small"
                  >
                    {{ userInfo.is_active ? '活跃' : '已停用' }}
                  </el-tag>
                  <el-tag 
                    :type="userInfo.email_verified ? 'success' : 'warning'" 
                    size="small"
                  >
                    {{ userInfo.email_verified ? '邮箱已验证' : '邮箱未验证' }}
                  </el-tag>
                </div>
                
                <div class="user-meta">
                  <p><strong>角色:</strong> {{ userInfo.role }}</p>
                  <p><strong>注册时间:</strong> {{ formatDate(userInfo.created_at) }}</p>
                  <p v-if="userInfo.last_login">
                    <strong>最后登录:</strong> {{ formatDate(userInfo.last_login) }}
                  </p>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>

        <!-- 右侧：详细信息和设置 -->
        <el-col :xs="24" :sm="24" :md="16" :lg="18">
          <el-tabs v-model="activeTab" class="profile-tabs">
            <!-- 基本信息标签页 -->
            <el-tab-pane label="基本信息" name="basic">
              <ProfileBasicInfo 
                :user-info="userInfo"
                @update="handleUserInfoUpdate"
                @loading="handleLoading"
              />
            </el-tab-pane>

            <!-- 安全设置标签页 -->
            <el-tab-pane label="安全设置" name="security">
              <ProfileSecurity 
                @loading="handleLoading"
              />
            </el-tab-pane>

            <!-- OAuth账户标签页 -->
            <el-tab-pane label="关联账户" name="oauth">
              <ProfileOAuth 
                :user-info="userInfo"
                @update="handleUserInfoUpdate"
                @loading="handleLoading"
              />
            </el-tab-pane>

            <!-- 偏好设置标签页 -->
            <el-tab-pane label="偏好设置" name="preferences">
              <ProfilePreferences 
                @loading="handleLoading"
              />
            </el-tab-pane>
          </el-tabs>
        </el-col>
      </el-row>
    </div>

    <!-- 头像上传对话框 -->
    <AvatarUploadDialog 
      v-model="showAvatarUpload"
      :current-avatar="fullAvatarUrl"
      :current-user="{ username: userInfo.username, display_name: userInfo.display_name, email: userInfo.email }"
      @uploaded="handleAvatarUploaded"
    />

    <!-- 全局加载遮罩 -->
    <div v-if="globalLoading" class="loading-overlay">
      <el-icon class="loading-icon" :size="40">
        <Loading />
      </el-icon>
      <p>处理中...</p>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, computed } from "vue"
import { ElMessage } from "element-plus"
import { User, Camera, Loading } from "@element-plus/icons-vue"
import { useUserStore } from "@/store/modules/user"
import { getUserProfileApi } from "@/api/profile"
import { buildAvatarUrl } from "@/utils/avatar"
import ProfileBasicInfo from "./components/ProfileBasicInfo.vue"
import ProfileSecurity from "./components/ProfileSecurity.vue"
import ProfileOAuth from "./components/ProfileOAuth.vue"
import ProfilePreferences from "./components/ProfilePreferences.vue"
import AvatarUploadDialog from "./components/AvatarUploadDialog.vue"

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
}

const userStore = useUserStore()
const activeTab = ref("basic")
const showAvatarUpload = ref(false)
const globalLoading = ref(false)

const userInfo = ref<UserInfo>({
  id: 0,
  username: "",
  email: "",
  display_name: "",
  avatar_url: "",
  role: "",
  is_active: true,
  email_verified: false,
  last_login: null,
  created_at: ""
})

/** 计算完整的头像URL */
const fullAvatarUrl = computed(() => buildAvatarUrl(userInfo.value.avatar_url))

/** 格式化日期 */
const formatDate = (dateString: string) => {
  if (!dateString) return "未知"
  return new Date(dateString).toLocaleString("zh-CN")
}

/** 加载用户资料 */
const loadUserProfile = async () => {
  try {
    globalLoading.value = true
    const { data } = await getUserProfileApi()
    userInfo.value = data
  } catch (error: any) {
    console.error("加载用户资料失败:", error)
    ElMessage.error("加载用户资料失败")
  } finally {
    globalLoading.value = false
  }
}

/** 处理用户信息更新 */
const handleUserInfoUpdate = (updatedInfo: Partial<UserInfo>) => {
  userInfo.value = { ...userInfo.value, ...updatedInfo }
}

/** 处理头像上传完成 */
const handleAvatarUploaded = (avatarUrl: string) => {
  userInfo.value.avatar_url = avatarUrl
  // 同时更新用户store中的头像信息
  if (userStore.userInfo) {
    userStore.userInfo.avatar_url = avatarUrl
  }
  ElMessage.success("头像更新成功")
}

/** 处理加载状态 */
const handleLoading = (loading: boolean) => {
  globalLoading.value = loading
}

/** 组件挂载时加载数据 */
onMounted(() => {
  loadUserProfile()
})
</script>

<style lang="scss" scoped>
.profile-container {
  padding: 24px;
  background-color: var(--el-bg-color-page);
  min-height: calc(100vh - 60px);
  position: relative;
}

.profile-header {
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

.profile-content {
  .profile-card {
    margin-bottom: 24px;
    
    .avatar-section {
      text-align: center;
      
      .avatar-container {
        position: relative;
        display: inline-block;
        margin-bottom: 20px;
        
        .user-avatar {
          border: 4px solid var(--el-color-primary-light-8);
          transition: all 0.3s ease;
        }
        
        .avatar-overlay {
          position: absolute;
          top: 0;
          left: 0;
          right: 0;
          bottom: 0;
          background: rgba(0, 0, 0, 0.6);
          border-radius: 50%;
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          opacity: 0;
          transition: opacity 0.3s ease;
          cursor: pointer;
          color: white;
          font-size: 12px;
          
          &:hover {
            opacity: 1;
          }
          
          .el-icon {
            font-size: 24px;
            margin-bottom: 4px;
          }
        }
      }
      
      .user-basic-info {
        h3 {
          margin: 0 0 8px 0;
          color: var(--el-text-color-primary);
          font-size: 20px;
          font-weight: 600;
        }
        
        .username {
          margin: 0 0 4px 0;
          color: var(--el-text-color-regular);
          font-size: 14px;
        }
        
        .email {
          margin: 0 0 16px 0;
          color: var(--el-text-color-secondary);
          font-size: 14px;
        }
        
        .user-status {
          margin-bottom: 16px;
          
          .el-tag {
            margin: 0 4px 4px 0;
          }
        }
        
        .user-meta {
          text-align: left;
          
          p {
            margin: 8px 0;
            font-size: 14px;
            color: var(--el-text-color-regular);
            
            strong {
              color: var(--el-text-color-primary);
            }
          }
        }
      }
    }
  }
  
  .profile-tabs {
    :deep(.el-tabs__header) {
      margin-bottom: 20px;
    }
    
    :deep(.el-tabs__content) {
      padding: 0;
    }
  }
}

.loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.8);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  
  .loading-icon {
    color: var(--el-color-primary);
    animation: rotate 2s linear infinite;
    margin-bottom: 16px;
  }
  
  p {
    color: var(--el-text-color-primary);
    font-size: 16px;
    margin: 0;
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

@media (max-width: 768px) {
  .profile-container {
    padding: 16px;
  }
  
  .profile-header {
    h1 {
      font-size: 24px;
    }
    
    .subtitle {
      font-size: 14px;
    }
  }
}
</style>