<template>
  <el-dialog v-model="dialogVisible" title="用户详情" width="700px" :close-on-click-modal="false" @close="handleClose">
    <div v-if="user" class="user-detail">
      <!-- 用户基本信息 -->
      <div class="user-header">
        <el-avatar :size="80" :src="buildAvatarUrlForTemplate(user.avatar_url)">
          {{ user.username.charAt(0).toUpperCase() }}
        </el-avatar>
        <div class="user-basic-info">
          <h3>{{ user.display_name || user.username }}</h3>
          <p class="username">@{{ user.username }}</p>
          <div class="user-tags">
            <el-tag :type="getRoleTagType(user.role)" size="small">
              {{ getRoleDisplayName(user.role) }}
            </el-tag>
            <el-tag :type="user.is_active ? 'success' : 'danger'" size="small">
              {{ user.is_active ? '活跃' : '停用' }}
            </el-tag>
            <el-tag :type="user.email_verified ? 'success' : 'warning'" size="small">
              {{ user.email_verified ? '邮箱已验证' : '邮箱未验证' }}
            </el-tag>
          </div>
        </div>
      </div>

      <!-- 详细信息标签页 -->
      <el-tabs v-model="activeTab" class="user-detail-tabs">
        <!-- 基本信息 -->
        <el-tab-pane label="基本信息" name="basic">
          <div class="info-section">
            <div class="info-grid">
              <div class="info-item">
                <label>用户ID</label>
                <span>{{ user.id }}</span>
              </div>
              <div class="info-item">
                <label>用户名</label>
                <span>{{ user.username }}</span>
              </div>
              <div class="info-item">
                <label>显示名称</label>
                <span>{{ user.display_name || '未设置' }}</span>
              </div>
              <div class="info-item">
                <label>邮箱地址</label>
                <span>{{ user.email }}</span>
              </div>
              <div class="info-item">
                <label>用户角色</label>
                <span>{{ getRoleDisplayName(user.role) }}</span>
              </div>
              <div class="info-item">
                <label>账户状态</label>
                <span>{{ user.is_active ? '活跃' : '停用' }}</span>
              </div>
              <div class="info-item">
                <label>邮箱验证</label>
                <span>{{ user.email_verified ? '已验证' : '未验证' }}</span>
              </div>
              <div class="info-item">
                <label>创建时间</label>
                <span>{{ formatDate(user.created_at) }}</span>
              </div>
              <div class="info-item">
                <label>最后登录</label>
                <span>{{ formatDate(user.last_login) }}</span>
              </div>
            </div>
          </div>
        </el-tab-pane>

        <!-- 登录记录 -->
        <el-tab-pane label="登录记录" name="login-history">
          <div class="login-history">
            <div v-if="loadingLoginHistory" class="loading-container">
              <el-icon class="loading-icon">
                <Loading />
              </el-icon>
              <span>加载登录记录中...</span>
            </div>

            <div v-else-if="loginHistory.length === 0" class="empty-container">
              <el-empty description="暂无登录记录" />
            </div>

            <el-timeline v-else>
              <el-timeline-item v-for="record in loginHistory" :key="record.id"
                :timestamp="formatDate(record.login_time)" :type="getLoginStatusType(record.status)">
                <div class="login-record">
                  <div class="record-header">
                    <span class="login-status">{{ getLoginStatusText(record.status) }}</span>
                    <span class="ip-address">{{ record.ip_address }}</span>
                  </div>
                  <div class="record-details">
                    <p><strong>位置:</strong> {{ record.location || '未知' }}</p>
                    <p><strong>设备:</strong> {{ record.user_agent }}</p>
                    <p v-if="record.failure_reason"><strong>失败原因:</strong> {{ record.failure_reason }}</p>
                  </div>
                </div>
              </el-timeline-item>
            </el-timeline>
          </div>
        </el-tab-pane>

        <!-- 权限信息 -->
        <el-tab-pane label="权限信息" name="permissions">
          <div class="permissions-section">
            <div v-if="loadingPermissions" class="loading-container">
              <el-icon class="loading-icon">
                <Loading />
              </el-icon>
              <span>加载权限信息中...</span>
            </div>

            <div v-else>
              <div class="role-permissions">
                <h4>角色权限</h4>
                <div class="permission-grid">
                  <div v-for="permission in rolePermissions" :key="permission.resource" class="permission-item">
                    <div class="permission-resource">{{ permission.resource }}</div>
                    <div class="permission-actions">
                      <el-tag v-for="action in permission.actions" :key="action" size="small" type="success">
                        {{ action }}
                      </el-tag>
                    </div>
                  </div>
                </div>
              </div>

              <div v-if="customPermissions.length > 0" class="custom-permissions">
                <h4>自定义权限</h4>
                <div class="permission-grid">
                  <div v-for="permission in customPermissions" :key="permission.id" class="permission-item">
                    <div class="permission-resource">{{ permission.resource }}</div>
                    <div class="permission-actions">
                      <el-tag v-for="action in permission.actions" :key="action" size="small" type="primary">
                        {{ action }}
                      </el-tag>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </el-tab-pane>

        <!-- 活动日志 -->
        <el-tab-pane label="活动日志" name="activity-log">
          <div class="activity-log">
            <div v-if="loadingActivityLog" class="loading-container">
              <el-icon class="loading-icon">
                <Loading />
              </el-icon>
              <span>加载活动日志中...</span>
            </div>

            <div v-else-if="activityLog.length === 0" class="empty-container">
              <el-empty description="暂无活动记录" />
            </div>

            <el-timeline v-else>
              <el-timeline-item v-for="activity in activityLog" :key="activity.id"
                :timestamp="formatDate(activity.created_at)" :type="getActivityType(activity.action)">
                <div class="activity-record">
                  <div class="activity-header">
                    <span class="activity-action">{{ getActivityTitle(activity.action) }}</span>
                  </div>
                  <div class="activity-details">
                    <p>{{ activity.description }}</p>
                    <div class="activity-meta">
                      <span>IP: {{ activity.ip_address }}</span>
                      <span v-if="activity.location">位置: {{ activity.location }}</span>
                    </div>
                  </div>
                </div>
              </el-timeline-item>
            </el-timeline>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">关闭</el-button>
        <el-button type="primary" @click="editUser">
          编辑用户
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import { ref, watch, onMounted } from "vue"
import { Loading } from "@element-plus/icons-vue"
import { useI18n } from "vue-i18n"
import { buildAvatarUrl } from "@/utils/avatar"
import {
  getUserLoginHistoryApi,
  getUserPermissionsApi,
  getUserActivityLogApi
} from "@/api/admin/user-management"

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
  edit: [user: User]
}>()

const { t } = useI18n()
const dialogVisible = ref(false)
const activeTab = ref("basic")

// 加载状态
const loadingLoginHistory = ref(false)
const loadingPermissions = ref(false)
const loadingActivityLog = ref(false)

// 数据
const loginHistory = ref<any[]>([])
const rolePermissions = ref<any[]>([])
const customPermissions = ref<any[]>([])
const activityLog = ref<any[]>([])

/** 获取角色标签类型 */
const getRoleTagType = (role: string) => {
  const types: Record<string, string> = {
    admin: "danger",
    editor: "warning",
    viewer: "info"
  }
  return types[role] || "info"
}

/** 获取角色显示名称 */
const getRoleDisplayName = (role: string) => {
  const roleKeys: Record<string, string> = {
    admin: 'admin.userManagement.roles.admin',
    editor: 'admin.userManagement.roles.editor',
    viewer: 'admin.userManagement.roles.viewer'
  }
  return roleKeys[role] ? t(roleKeys[role]) : role
}

/** 暴露 buildAvatarUrl 函数给模板使用 */
const buildAvatarUrlForTemplate = buildAvatarUrl

/** 格式化日期 */
const formatDate = (dateString: string | null) => {
  if (!dateString) return "从未"
  return new Date(dateString).toLocaleString("zh-CN")
}

/** 获取登录状态类型 */
const getLoginStatusType = (status: string) => {
  return status === "success" ? "success" : "danger"
}

/** 获取登录状态文本 */
const getLoginStatusText = (status: string) => {
  return status === "success" ? "登录成功" : "登录失败"
}

/** 获取活动类型 */
const getActivityType = (action: string) => {
  const types: Record<string, string> = {
    login: "success",
    logout: "info",
    update_profile: "primary",
    change_password: "warning",
    delete_resource: "danger"
  }
  return types[action] || "primary"
}

/** 获取活动标题 */
const getActivityTitle = (action: string) => {
  const titles: Record<string, string> = {
    login: "用户登录",
    logout: "用户登出",
    update_profile: "更新资料",
    change_password: "修改密码",
    create_resource: "创建资源",
    update_resource: "更新资源",
    delete_resource: "删除资源"
  }
  return titles[action] || action
}

/** 加载登录记录 */
const loadLoginHistory = async () => {
  if (!props.user) return

  try {
    loadingLoginHistory.value = true
    const { data } = await getUserLoginHistoryApi(props.user.id)
    loginHistory.value = data
  } catch (error: any) {
    console.error("加载登录记录失败:", error)
    loginHistory.value = []
  } finally {
    loadingLoginHistory.value = false
  }
}

/** 加载权限信息 */
const loadPermissions = async () => {
  if (!props.user) return

  try {
    loadingPermissions.value = true
    const { data } = await getUserPermissionsApi(props.user.id)
    rolePermissions.value = data.role_permissions || []
    customPermissions.value = data.custom_permissions || []
  } catch (error: any) {
    console.error("加载权限信息失败:", error)
    rolePermissions.value = []
    customPermissions.value = []
  } finally {
    loadingPermissions.value = false
  }
}

/** 加载活动日志 */
const loadActivityLog = async () => {
  if (!props.user) return

  try {
    loadingActivityLog.value = true
    const { data } = await getUserActivityLogApi(props.user.id)
    activityLog.value = data
  } catch (error: any) {
    console.error("加载活动日志失败:", error)
    activityLog.value = []
  } finally {
    loadingActivityLog.value = false
  }
}

/** 监听对话框显示状态 */
watch(() => props.modelValue, (val) => {
  dialogVisible.value = val
})

watch(dialogVisible, (val) => {
  emit("update:modelValue", val)
})

/** 监听标签页切换 */
watch(activeTab, (tab) => {
  if (!props.user) return

  switch (tab) {
    case "login-history":
      if (loginHistory.value.length === 0) {
        loadLoginHistory()
      }
      break
    case "permissions":
      if (rolePermissions.value.length === 0) {
        loadPermissions()
      }
      break
    case "activity-log":
      if (activityLog.value.length === 0) {
        loadActivityLog()
      }
      break
  }
})

/** 编辑用户 */
const editUser = () => {
  if (props.user) {
    emit("edit", props.user)
    handleClose()
  }
}

/** 关闭对话框 */
const handleClose = () => {
  activeTab.value = "basic"
  dialogVisible.value = false
}
</script>

<style lang="scss" scoped>
.user-detail {
  .user-header {
    display: flex;
    align-items: center;
    margin-bottom: 24px;
    padding: 20px;
    background: var(--el-fill-color-extra-light);
    border-radius: 8px;

    .user-basic-info {
      margin-left: 20px;

      h3 {
        margin: 0 0 8px 0;
        font-size: 20px;
        font-weight: 600;
        color: var(--el-text-color-primary);
      }

      .username {
        margin: 0 0 12px 0;
        color: var(--el-text-color-regular);
        font-size: 14px;
      }

      .user-tags {
        display: flex;
        gap: 8px;
        flex-wrap: wrap;
      }
    }
  }

  .user-detail-tabs {
    :deep(.el-tabs__content) {
      padding: 0;
    }
  }

  .info-section {
    .info-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
      gap: 16px;

      .info-item {
        display: flex;
        flex-direction: column;
        gap: 4px;

        label {
          font-weight: 500;
          color: var(--el-text-color-secondary);
          font-size: 12px;
          text-transform: uppercase;
        }

        span {
          color: var(--el-text-color-primary);
          font-size: 14px;
        }
      }
    }
  }

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

  .empty-container {
    padding: 20px;
  }

  .login-history,
  .activity-log {

    .login-record,
    .activity-record {

      .record-header,
      .activity-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 8px;

        .login-status,
        .activity-action {
          font-weight: 500;
          color: var(--el-text-color-primary);
        }

        .ip-address {
          font-size: 12px;
          color: var(--el-text-color-secondary);
          font-family: monospace;
        }
      }

      .record-details,
      .activity-details {
        p {
          margin: 4px 0;
          font-size: 14px;
          color: var(--el-text-color-regular);

          strong {
            color: var(--el-text-color-primary);
          }
        }

        .activity-meta {
          margin-top: 8px;
          font-size: 12px;
          color: var(--el-text-color-secondary);

          span {
            margin-right: 16px;
          }
        }
      }
    }
  }

  .permissions-section {

    .role-permissions,
    .custom-permissions {
      margin-bottom: 24px;

      h4 {
        margin: 0 0 16px 0;
        font-size: 16px;
        font-weight: 600;
        color: var(--el-text-color-primary);
      }

      .permission-grid {
        display: flex;
        flex-direction: column;
        gap: 12px;

        .permission-item {
          display: flex;
          justify-content: space-between;
          align-items: center;
          padding: 12px;
          border: 1px solid var(--el-border-color-light);
          border-radius: 6px;

          .permission-resource {
            font-weight: 500;
            color: var(--el-text-color-primary);
          }

          .permission-actions {
            display: flex;
            gap: 4px;
            flex-wrap: wrap;
          }
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

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }

  to {
    transform: rotate(360deg);
  }
}

@media (max-width: 768px) {
  .user-detail {
    .user-header {
      flex-direction: column;
      text-align: center;

      .user-basic-info {
        margin-left: 0;
        margin-top: 16px;
      }
    }

    .info-section {
      .info-grid {
        grid-template-columns: 1fr;
      }
    }

    .permissions-section {
      .permission-grid {
        .permission-item {
          flex-direction: column;
          align-items: flex-start;
          gap: 8px;
        }
      }
    }
  }
}
</style>