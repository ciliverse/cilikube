<template>
  <el-dialog v-model="dialogVisible" :title="`角色用户 - ${role?.display_name}`" width="700px" :close-on-click-modal="false"
    @close="handleClose">
    <div class="role-users-container">
      <!-- 用户统计 -->
      <div class="users-stats">
        <div class="stat-item">
          <el-icon class="stat-icon">
            <User />
          </el-icon>
          <div class="stat-info">
            <div class="stat-value">{{ users.length }}</div>
            <div class="stat-label">总用户数</div>
          </div>
        </div>
        <div class="stat-item">
          <el-icon class="stat-icon">
            <UserFilled />
          </el-icon>
          <div class="stat-info">
            <div class="stat-value">{{ activeUsers }}</div>
            <div class="stat-label">活跃用户</div>
          </div>
        </div>
      </div>

      <!-- 搜索和操作 -->
      <div class="users-actions">
        <el-input v-model="searchQuery" placeholder="搜索用户名或邮箱..." :prefix-icon="Search" clearable
          style="width: 300px" />
        <el-button :icon="Refresh" @click="loadRoleUsers" :loading="loading">
          刷新
        </el-button>
      </div>

      <!-- 用户列表 -->
      <div class="users-content">
        <div v-if="loading" class="loading-container">
          <el-icon class="loading-icon" :size="40">
            <Loading />
          </el-icon>
          <p>加载用户数据中...</p>
        </div>

        <div v-else-if="filteredUsers.length === 0" class="empty-container">
          <el-empty :image-size="80" description="暂无用户数据" />
        </div>

        <div v-else class="users-list">
          <div v-for="user in filteredUsers" :key="user.id" class="user-item">
            <div class="user-avatar">
              <el-avatar :size="40" :src="buildAvatarUrlForTemplate(user.avatar_url)"
                :alt="user.display_name || user.username">
                {{ (user.display_name || user.username).charAt(0).toUpperCase() }}
              </el-avatar>
            </div>

            <div class="user-info">
              <div class="user-basic">
                <h4>{{ user.display_name || user.username }}</h4>
                <span class="user-username">@{{ user.username }}</span>
              </div>
              <div class="user-details">
                <span v-if="user.email" class="user-email">{{ user.email }}</span>
                <span class="user-login">最后登录: {{ formatDate(user.last_login_at) }}</span>
              </div>
            </div>

            <div class="user-status">
              <el-tag :type="user.is_active ? 'success' : 'danger'" size="small">
                {{ user.is_active ? '活跃' : '禁用' }}
              </el-tag>
              <el-tag v-if="user.is_admin" type="warning" size="small">
                管理员
              </el-tag>
            </div>

            <div class="user-actions">
              <el-dropdown @command="(command) => handleUserAction(command, user)" trigger="click">
                <el-button size="small" type="text">
                  <el-icon>
                    <MoreFilled />
                  </el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="view">
                      查看详情
                    </el-dropdown-item>
                    <el-dropdown-item command="edit">
                      编辑用户
                    </el-dropdown-item>
                    <el-dropdown-item v-if="!role?.is_system" command="remove" divided
                      style="color: var(--el-color-danger)">
                      移除角色
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </div>
        </div>
      </div>

      <!-- 分页 -->
      <div v-if="!loading && filteredUsers.length > 0" class="users-pagination">
        <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize" :page-sizes="[10, 20, 50, 100]"
          :total="filteredUsers.length" layout="total, sizes, prev, pager, next, jumper" @size-change="handleSizeChange"
          @current-change="handleCurrentChange" />
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">关闭</el-button>
        <el-button type="primary" @click="exportUsers" :loading="exporting">
          导出用户列表
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import { ref, computed, watch } from "vue"
import { ElMessage, ElMessageBox } from "element-plus"
import {
  Search,
  Refresh,
  Loading,
  User,
  UserFilled,
  MoreFilled
} from "@element-plus/icons-vue"
import { buildAvatarUrl } from "@/utils/avatar"
import {
  getRoleUsersApi,
  removeUserFromRoleApi
} from "@/api/admin/role-management"

interface Role {
  id: number
  name: string
  display_name: string
  is_system: boolean
}

interface RoleUser {
  id: number
  username: string
  display_name: string
  email: string
  avatar_url: string
  is_active: boolean
  is_admin: boolean
  last_login_at: string
  created_at: string
}

interface Props {
  modelValue: boolean
  role: Role | null
}

const props = defineProps<Props>()
const emit = defineEmits<{
  "update:modelValue": [value: boolean]
  viewUser: [user: RoleUser]
  editUser: [user: RoleUser]
}>()

const loading = ref(false)
const exporting = ref(false)
const dialogVisible = ref(false)
const searchQuery = ref("")
const users = ref<RoleUser[]>([])
const currentPage = ref(1)
const pageSize = ref(20)

/** 计算属性 */
const activeUsers = computed(() => {
  return users.value.filter(user => user.is_active).length
})

const filteredUsers = computed(() => {
  if (!searchQuery.value) return users.value

  const query = searchQuery.value.toLowerCase()
  return users.value.filter(user =>
    user.username.toLowerCase().includes(query) ||
    (user.display_name && user.display_name.toLowerCase().includes(query)) ||
    (user.email && user.email.toLowerCase().includes(query))
  )
})

/** 格式化日期 */
const formatDate = (dateString: string) => {
  if (!dateString) return "从未登录"
  return new Date(dateString).toLocaleString("zh-CN")
}

/** 暴露 buildAvatarUrl 函数给模板使用 */
const buildAvatarUrlForTemplate = buildAvatarUrl

/** 加载角色用户 */
const loadRoleUsers = async () => {
  if (!props.role) return

  try {
    loading.value = true
    const { data } = await getRoleUsersApi(props.role.id)
    users.value = data.users || []
  } catch (error: any) {
    console.error("加载角色用户失败:", error)
    ElMessage.error("加载角色用户失败")
  } finally {
    loading.value = false
  }
}

/** 用户操作处理 */
const handleUserAction = async (command: string, user: RoleUser) => {
  switch (command) {
    case 'view':
      emit("viewUser", user)
      break
    case 'edit':
      emit("editUser", user)
      break
    case 'remove':
      await removeUserFromRole(user)
      break
  }
}

/** 从角色中移除用户 */
const removeUserFromRole = async (user: RoleUser) => {
  if (!props.role) return

  if (props.role.is_system) {
    ElMessage.warning("系统角色不能移除用户")
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要将用户 "${user.display_name || user.username}" 从角色 "${props.role.display_name}" 中移除吗？`,
      "确认移除",
      {
        confirmButtonText: "确定移除",
        cancelButtonText: "取消",
        type: "warning"
      }
    )

    await removeUserFromRoleApi(props.role.id, user.id)
    ElMessage.success("用户移除成功")
    loadRoleUsers()

  } catch (error: any) {
    if (error !== "cancel") {
      console.error("移除用户失败:", error)
      ElMessage.error("移除用户失败")
    }
  }
}

/** 导出用户列表 */
const exportUsers = async () => {
  try {
    exporting.value = true

    // 准备导出数据
    const exportData = filteredUsers.value.map(user => ({
      用户名: user.username,
      显示名称: user.display_name || "",
      邮箱: user.email || "",
      状态: user.is_active ? "活跃" : "禁用",
      管理员: user.is_admin ? "是" : "否",
      最后登录: formatDate(user.last_login_at),
      创建时间: formatDate(user.created_at)
    }))

    // 转换为CSV格式
    const headers = Object.keys(exportData[0])
    const csvContent = [
      headers.join(","),
      ...exportData.map(row =>
        headers.map(header => `"${row[header as keyof typeof row]}"`).join(",")
      )
    ].join("\\n")

    // 创建下载链接
    const blob = new Blob([csvContent], { type: "text/csv;charset=utf-8;" })
    const link = document.createElement("a")
    const url = URL.createObjectURL(blob)
    link.setAttribute("href", url)
    link.setAttribute("download", `${props.role?.name}_users_${new Date().toISOString().split('T')[0]}.csv`)
    link.style.visibility = "hidden"
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)

    ElMessage.success("用户列表导出成功")

  } catch (error: any) {
    console.error("导出用户列表失败:", error)
    ElMessage.error("导出用户列表失败")
  } finally {
    exporting.value = false
  }
}

/** 分页处理 */
const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
}

/** 监听对话框显示状态 */
watch(() => props.modelValue, (val) => {
  dialogVisible.value = val
  if (val && props.role) {
    loadRoleUsers()
  }
})

watch(dialogVisible, (val) => {
  emit("update:modelValue", val)
})

/** 关闭对话框 */
const handleClose = () => {
  dialogVisible.value = false
  searchQuery.value = ""
  currentPage.value = 1
}
</script>

<style lang="scss" scoped>
.role-users-container {
  .users-stats {
    display: flex;
    gap: 16px;
    margin-bottom: 24px;

    .stat-item {
      display: flex;
      align-items: center;
      padding: 16px;
      background: var(--el-fill-color-extra-light);
      border-radius: 8px;
      flex: 1;

      .stat-icon {
        width: 32px;
        height: 32px;
        border-radius: 6px;
        background: var(--el-color-primary-light-8);
        color: var(--el-color-primary);
        display: flex;
        align-items: center;
        justify-content: center;
        margin-right: 12px;
      }

      .stat-info {
        .stat-value {
          font-size: 18px;
          font-weight: 600;
          color: var(--el-text-color-primary);
          line-height: 1;
        }

        .stat-label {
          font-size: 12px;
          color: var(--el-text-color-secondary);
          margin-top: 4px;
        }
      }
    }
  }

  .users-actions {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
  }

  .users-content {
    min-height: 300px;
    max-height: 400px;
    overflow-y: auto;

    .loading-container {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 60px;
      color: var(--el-text-color-secondary);

      .loading-icon {
        margin-bottom: 16px;
        animation: rotate 2s linear infinite;
      }
    }

    .empty-container {
      display: flex;
      justify-content: center;
      align-items: center;
      padding: 40px;
    }
  }
}

.users-list {
  .user-item {
    display: flex;
    align-items: center;
    padding: 12px;
    border: 1px solid var(--el-border-color-lighter);
    border-radius: 8px;
    margin-bottom: 8px;
    transition: all 0.3s ease;

    &:hover {
      border-color: var(--el-color-primary-light-7);
      background: var(--el-color-primary-light-9);
    }

    .user-avatar {
      margin-right: 12px;
    }

    .user-info {
      flex: 1;

      .user-basic {
        display: flex;
        align-items: center;
        gap: 8px;
        margin-bottom: 4px;

        h4 {
          margin: 0;
          font-size: 14px;
          font-weight: 500;
          color: var(--el-text-color-primary);
        }

        .user-username {
          font-size: 12px;
          color: var(--el-text-color-secondary);
          font-family: monospace;
        }
      }

      .user-details {
        display: flex;
        gap: 16px;
        font-size: 12px;
        color: var(--el-text-color-regular);

        .user-email {
          color: var(--el-color-primary);
        }
      }
    }

    .user-status {
      display: flex;
      flex-direction: column;
      gap: 4px;
      margin-right: 12px;
    }

    .user-actions {
      display: flex;
      align-items: center;
    }
  }
}

.users-pagination {
  margin-top: 16px;
  display: flex;
  justify-content: center;
}

.dialog-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
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