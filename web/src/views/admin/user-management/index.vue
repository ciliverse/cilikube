<template>
  <div class="user-management">
    <div class="page-header">
      <h1>{{ t('menu.userManagement') }}</h1>
      <p class="subtitle">{{ t('admin.userManagement.description') }}</p>
    </div>

    <!-- 操作栏 -->
    <div class="action-bar">
      <div class="search-section">
        <el-input v-model="searchQuery" :placeholder="t('admin.userManagement.searchPlaceholder')" :prefix-icon="Search" clearable
          style="width: 300px" @input="handleSearch" />

        <el-select v-model="statusFilter" :placeholder="t('admin.userManagement.statusFilter')" clearable style="width: 120px; margin-left: 12px"
          @change="handleFilter">
          <el-option :label="t('common.all')" value="" />
          <el-option :label="t('admin.userManagement.status.active')" value="active" />
          <el-option :label="t('admin.userManagement.status.inactive')" value="inactive" />
        </el-select>

        <el-select v-model="roleFilter" :placeholder="t('admin.userManagement.roleFilter')" clearable style="width: 120px; margin-left: 12px"
          @change="handleFilter">
          <el-option :label="t('common.all')" value="" />
          <el-option :label="t('admin.userManagement.roles.admin')" value="admin" />
          <el-option :label="t('admin.userManagement.roles.editor')" value="editor" />
          <el-option :label="t('admin.userManagement.roles.viewer')" value="viewer" />
        </el-select>

        <el-button :icon="Refresh" @click="refreshUsers" :loading="loading" style="margin-left: 12px">
          {{ t('common.refresh') }}
        </el-button>
      </div>

      <div class="action-buttons">
        <el-button type="primary" :icon="Plus" @click="showCreateDialog = true">
          {{ t('admin.userManagement.actions.addUser') }}
        </el-button>

        <el-button :icon="Download" @click="exportUsers" :loading="exportLoading">
          {{ t('admin.userManagement.actions.exportUsers') }}
        </el-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-cards">
      <el-row :gutter="16">
        <el-col :xs="12" :sm="6">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon total">
                <el-icon>
                  <User />
                </el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-number">{{ userStats.total }}</div>
                <div class="stat-label">{{ t('admin.userManagement.stats.totalUsers') }}</div>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :xs="12" :sm="6">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon active">
                <el-icon>
                  <UserFilled />
                </el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-number">{{ userStats.active }}</div>
                <div class="stat-label">{{ t('admin.userManagement.stats.activeUsers') }}</div>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :xs="12" :sm="6">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon admin">
                <el-icon>
                  <Avatar />
                </el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-number">{{ userStats.admins }}</div>
                <div class="stat-label">{{ t('admin.userManagement.stats.admins') }}</div>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :xs="12" :sm="6">
          <el-card class="stat-card">
            <div class="stat-content">
              <div class="stat-icon new">
                <el-icon>
                  <Calendar />
                </el-icon>
              </div>
              <div class="stat-info">
                <div class="stat-number">{{ userStats.newThisMonth }}</div>
                <div class="stat-label">{{ t('admin.userManagement.stats.newThisMonth') }}</div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 用户列表 -->
    <el-card class="user-table-card">
      <div v-if="loading" class="loading-container">
        <el-icon class="loading-icon" :size="40">
          <Loading />
        </el-icon>
        <p>{{ t('admin.userManagement.loading') }}</p>
      </div>

      <el-table v-else :data="filteredUsers" style="width: 100%"
        :default-sort="{ prop: 'created_at', order: 'descending' }" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />

        <el-table-column prop="avatar_url" :label="t('admin.userManagement.columns.avatar')" width="80">
          <template #default="{ row }">
            <el-avatar :size="40" :src="row.avatar_url">
              {{ row.username.charAt(0).toUpperCase() }}
            </el-avatar>
          </template>
        </el-table-column>

        <el-table-column prop="username" :label="t('admin.userManagement.columns.username')" min-width="120" sortable>
          <template #default="{ row }">
            <div class="user-info">
              <div class="username">{{ row.username }}</div>
              <div class="display-name">{{ row.display_name || t('admin.userManagement.notSet') }}</div>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="email" :label="t('admin.userManagement.columns.email')" min-width="200" sortable>
          <template #default="{ row }">
            <div class="email-info">
              <div>{{ row.email }}</div>
              <el-tag v-if="row.email_verified" type="success" size="small">
                {{ t('admin.userManagement.emailStatus.verified') }}
              </el-tag>
              <el-tag v-else type="warning" size="small">
                {{ t('admin.userManagement.emailStatus.unverified') }}
              </el-tag>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="role" :label="t('admin.userManagement.columns.role')" width="100" sortable>
          <template #default="{ row }">
            <el-tag :type="getRoleTagType(getUserRole(row))" size="small">
              {{ getRoleDisplayName(getUserRole(row)) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="is_active" :label="t('admin.userManagement.columns.status')" width="80" sortable>
          <template #default="{ row }">
            <el-tag :type="row.is_active ? 'success' : 'danger'" size="small">
              {{ row.is_active ? t('admin.userManagement.status.active') : t('admin.userManagement.status.inactive') }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="last_login" :label="t('admin.userManagement.columns.lastLogin')" width="160" sortable>
          <template #default="{ row }">
            {{ formatDate(row.last_login) }}
          </template>
        </el-table-column>

        <el-table-column prop="created_at" :label="t('admin.userManagement.columns.createdAt')" width="160" sortable>
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>

        <el-table-column :label="t('common.actions')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="viewUser(row)">
              {{ t('common.view') }}
            </el-button>
            <el-button size="small" type="primary" @click="editUser(row)">
              {{ t('common.edit') }}
            </el-button>
            <el-dropdown @command="(command) => handleUserAction(command, row)" trigger="click">
              <el-button size="small" type="info">
                {{ t('admin.userManagement.actions.more') }}<el-icon class="el-icon--right">
                  <ArrowDown />
                </el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item :command="row.is_active ? 'deactivate' : 'activate'">
                    {{ row.is_active ? t('admin.userManagement.actions.deactivate') : t('admin.userManagement.actions.activate') }}
                  </el-dropdown-item>
                  <el-dropdown-item command="reset-password">
                    {{ t('admin.userManagement.actions.resetPassword') }}
                  </el-dropdown-item>
                  <el-dropdown-item command="send-verification">
                    {{ t('admin.userManagement.actions.sendVerification') }}
                  </el-dropdown-item>
                  <el-dropdown-item command="delete" divided style="color: var(--el-color-danger)">
                    {{ t('admin.userManagement.actions.deleteUser') }}
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize" :page-sizes="[10, 20, 50, 100]"
          :total="totalUsers" layout="total, sizes, prev, pager, next, jumper" @size-change="handleSizeChange"
          @current-change="handleCurrentChange" />
      </div>
    </el-card>

    <!-- 批量操作栏 -->
    <div v-if="selectedUsers.length > 0" class="batch-actions">
      <el-card>
        <div class="batch-content">
          <span>{{ t('admin.userManagement.batch.selected', { count: selectedUsers.length }) }}</span>
          <div class="batch-buttons">
            <el-button size="small" @click="batchActivate">
              {{ t('admin.userManagement.batch.activate') }}
            </el-button>
            <el-button size="small" @click="batchDeactivate">
              {{ t('admin.userManagement.batch.deactivate') }}
            </el-button>
            <el-button size="small" type="danger" @click="batchDelete">
              {{ t('admin.userManagement.batch.delete') }}
            </el-button>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 创建用户对话框 -->
    <CreateUserDialog v-model="showCreateDialog" @created="handleUserCreated" />

    <!-- 编辑用户对话框 -->
    <EditUserDialog v-model="showEditDialog" :user="selectedUser" @updated="handleUserUpdated" />

    <!-- 用户详情对话框 -->
    <UserDetailDialog v-model="showDetailDialog" :user="selectedUser" />
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, computed, onMounted } from "vue"
import { useI18n } from "vue-i18n"
import { ElMessage, ElMessageBox } from "element-plus"
import {
  Search,
  Refresh,
  Plus,
  Download,
  User,
  UserFilled,
  Avatar,
  Calendar,
  Loading,
  ArrowDown
} from "@element-plus/icons-vue"
import {
  getUsersApi,
  createUserApi,
  updateUserApi,
  deleteUserApi,
  batchUpdateUsersApi,
  exportUsersApi,
  getUserStatsApi
} from "@/api/admin/user-management"
import CreateUserDialog from "./components/CreateUserDialog.vue"
import EditUserDialog from "./components/EditUserDialog.vue"
import UserDetailDialog from "./components/UserDetailDialog.vue"

const { t } = useI18n()

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

interface UserStats {
  total: number
  active: number
  admins: number
  newThisMonth: number
}

// 响应式数据
const loading = ref(false)
const exportLoading = ref(false)
const users = ref<User[]>([])
const selectedUsers = ref<User[]>([])
const selectedUser = ref<User | null>(null)

// 对话框状态
const showCreateDialog = ref(false)
const showEditDialog = ref(false)
const showDetailDialog = ref(false)

// 搜索和筛选
const searchQuery = ref("")
const statusFilter = ref("")
const roleFilter = ref("")

// 分页
const currentPage = ref(1)
const pageSize = ref(20)
const totalUsers = ref(0)

// 统计数据
const userStats = reactive<UserStats>({
  total: 0,
  active: 0,
  admins: 0,
  newThisMonth: 0
})

// 计算属性
const filteredUsers = computed(() => {
  let filtered = users.value || []

  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(user =>
      user.username.toLowerCase().includes(query) ||
      user.email.toLowerCase().includes(query) ||
      (user.display_name && user.display_name.toLowerCase().includes(query))
    )
  }

  // 状态过滤
  if (statusFilter.value) {
    filtered = filtered.filter(user => {
      if (statusFilter.value === 'active') return user.is_active
      if (statusFilter.value === 'inactive') return !user.is_active
      return true
    })
  }

  // 角色过滤
  if (roleFilter.value) {
    filtered = filtered.filter(user => getUserRole(user) === roleFilter.value)
  }

  totalUsers.value = filtered.length
  return filtered
})

/** 格式化日期 */
const formatDate = (dateString: string | null) => {
  if (!dateString) return t('admin.userManagement.never')
  return new Date(dateString).toLocaleString()
}

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
  const names: Record<string, string> = {
    admin: t('admin.userManagement.roles.admin'),
    editor: t('admin.userManagement.roles.editor'),
    viewer: t('admin.userManagement.roles.viewer')
  }
  return names[role] || role
}

/** 获取用户角色 */
const getUserRole = (user: any) => {
  // 处理后端返回的roles数组格式
  if (user.roles && Array.isArray(user.roles) && user.roles.length > 0) {
    return user.roles[0] // 取第一个角色
  }
  // 兼容直接的role字段
  return user.role || 'viewer'
}

/** 加载用户列表 */
const loadUsers = async () => {
  try {
    loading.value = true
    console.log("开始加载用户列表...")
    
    const response = await getUsersApi({
      page: currentPage.value,
      page_size: pageSize.value,
      search: searchQuery.value,
      status: statusFilter.value,
      role: roleFilter.value
    })

    console.log("API响应:", response)

    // 处理后端返回的数据格式
    if (response.data && response.data.data) {
      users.value = response.data.data || []
      totalUsers.value = response.data.pagination?.total || 0
      console.log("用户数据:", users.value)
      console.log("总用户数:", totalUsers.value)
    } else {
      users.value = []
      totalUsers.value = 0
      console.log("没有用户数据")
    }
    
    // 加载完用户数据后更新统计信息
    loadUserStats()
  } catch (error: any) {
    console.error("加载用户列表失败:", error)
    ElMessage.error(t('admin.userManagement.messages.loadFailed', { error: error.message || error }))
    users.value = []
    totalUsers.value = 0
  } finally {
    loading.value = false
  }
}

/** 加载用户统计 */
const loadUserStats = async () => {
  try {
    // 暂时使用用户列表数据来计算统计信息
    // 后续可以实现专门的统计API
    if (users.value && users.value.length > 0) {
      const total = users.value.length
      const active = users.value.filter(user => user.is_active).length
      const admins = users.value.filter(user => getUserRole(user) === 'admin').length

      // 计算本月新增用户（简单实现）
      const thisMonth = new Date()
      thisMonth.setDate(1)
      thisMonth.setHours(0, 0, 0, 0)
      const newThisMonth = users.value.filter(user => {
        const createdAt = new Date(user.created_at)
        return createdAt >= thisMonth
      }).length

      Object.assign(userStats, {
        total,
        active,
        admins,
        newThisMonth
      })
    } else {
      // 如果没有用户数据，设置默认值
      Object.assign(userStats, {
        total: 0,
        active: 0,
        admins: 0,
        newThisMonth: 0
      })
    }
  } catch (error: any) {
    console.error("加载用户统计失败:", error)
    // 设置默认统计值
    Object.assign(userStats, {
      total: 0,
      active: 0,
      admins: 0,
      newThisMonth: 0
    })
  }
}

/** 搜索处理 */
const handleSearch = () => {
  currentPage.value = 1
  loadUsers()
}

/** 筛选处理 */
const handleFilter = () => {
  currentPage.value = 1
  loadUsers()
}

/** 刷新用户列表 */
const refreshUsers = () => {
  loadUsers()
  loadUserStats()
}

/** 分页大小变化 */
const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  loadUsers()
}

/** 当前页变化 */
const handleCurrentChange = (page: number) => {
  currentPage.value = page
  loadUsers()
}

/** 选择变化 */
const handleSelectionChange = (selection: User[]) => {
  selectedUsers.value = selection
}

/** 查看用户 */
const viewUser = (user: User) => {
  selectedUser.value = user
  showDetailDialog.value = true
}

/** 编辑用户 */
const editUser = (user: User) => {
  selectedUser.value = user
  showEditDialog.value = true
}

/** 用户操作处理 */
const handleUserAction = async (command: string, user: User) => {
  switch (command) {
    case 'activate':
    case 'deactivate':
      await toggleUserStatus(user)
      break
    case 'reset-password':
      await resetUserPassword(user)
      break
    case 'send-verification':
      await sendVerificationEmail(user)
      break
    case 'delete':
      await deleteUser(user)
      break
  }
}

/** 切换用户状态 */
const toggleUserStatus = async (user: User) => {
  try {
    if (!user || !user.id) {
      ElMessage.error(t('admin.userManagement.messages.invalidUser'))
      return
    }

    const action = user.is_active ? t('admin.userManagement.actions.deactivate') : t('admin.userManagement.actions.activate')
    await ElMessageBox.confirm(
      t('admin.userManagement.messages.confirmToggleStatus', { action, username: user.username }),
      t('admin.userManagement.messages.confirmAction', { action }),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: "warning"
      }
    )

    await updateUserApi(user.id, { is_active: !user.is_active })
    ElMessage.success(t('admin.userManagement.messages.toggleStatusSuccess', { action }))
    refreshUsers()
  } catch (error: any) {
    if (error !== "cancel") {
      ElMessage.error(t('admin.userManagement.messages.toggleStatusFailed'))
    }
  }
}

/** 重置用户密码 */
const resetUserPassword = async (user: User) => {
  try {
    if (!user || !user.id) {
      ElMessage.error(t('admin.userManagement.messages.invalidUser'))
      return
    }

    await ElMessageBox.confirm(
      t('admin.userManagement.messages.confirmResetPassword', { username: user.username }),
      t('admin.userManagement.messages.confirmResetPasswordTitle'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: "warning"
      }
    )

    // 这里应该调用重置密码的API
    ElMessage.success(t('admin.userManagement.messages.resetPasswordSuccess'))
  } catch (error: any) {
    if (error !== "cancel") {
      ElMessage.error(t('admin.userManagement.messages.resetPasswordFailed'))
    }
  }
}

/** 发送验证邮件 */
const sendVerificationEmail = async (user: User) => {
  try {
    if (!user || !user.id || !user.email) {
      ElMessage.error(t('admin.userManagement.messages.invalidUser'))
      return
    }

    // 这里应该调用发送验证邮件的API
    ElMessage.success(t('admin.userManagement.messages.verificationEmailSent', { email: user.email }))
  } catch (error: any) {
    ElMessage.error(t('admin.userManagement.messages.verificationEmailFailed'))
  }
}

/** 删除用户 */
const deleteUser = async (user: User) => {
  try {
    await ElMessageBox.confirm(
      t('admin.userManagement.messages.confirmDeleteUser', { username: user.username }),
      t('admin.userManagement.messages.confirmDeleteTitle'),
      {
        confirmButtonText: t('admin.userManagement.messages.confirmDeleteButton'),
        cancelButtonText: t('common.cancel'),
        type: "error"
      }
    )

    await deleteUserApi(user.id)
    ElMessage.success(t('admin.userManagement.messages.deleteUserSuccess'))
    refreshUsers()
  } catch (error: any) {
    if (error !== "cancel") {
      ElMessage.error(t('admin.userManagement.messages.deleteUserFailed'))
    }
  }
}

/** 批量激活 */
const batchActivate = async () => {
  try {
    const userIds = selectedUsers.value.map(user => user.id)
    await batchUpdateUsersApi(userIds, { is_active: true })
    ElMessage.success(t('admin.userManagement.messages.batchActivateSuccess', { count: userIds.length }))
    selectedUsers.value = []
    refreshUsers()
  } catch (error: any) {
    ElMessage.error(t('admin.userManagement.messages.batchActivateFailed'))
  }
}

/** 批量停用 */
const batchDeactivate = async () => {
  try {
    const userIds = selectedUsers.value.map(user => user.id)
    await batchUpdateUsersApi(userIds, { is_active: false })
    ElMessage.success(t('admin.userManagement.messages.batchDeactivateSuccess', { count: userIds.length }))
    selectedUsers.value = []
    refreshUsers()
  } catch (error: any) {
    ElMessage.error(t('admin.userManagement.messages.batchDeactivateFailed'))
  }
}

/** 批量删除 */
const batchDelete = async () => {
  try {
    await ElMessageBox.confirm(
      t('admin.userManagement.messages.confirmBatchDelete', { count: selectedUsers.value.length }),
      t('admin.userManagement.messages.confirmBatchDeleteTitle'),
      {
        confirmButtonText: t('admin.userManagement.messages.confirmDeleteButton'),
        cancelButtonText: t('common.cancel'),
        type: "error"
      }
    )

    const userIds = selectedUsers.value.map(user => user.id)
    // 这里应该调用批量删除的API
    ElMessage.success(t('admin.userManagement.messages.batchDeleteSuccess', { count: userIds.length }))
    selectedUsers.value = []
    refreshUsers()
  } catch (error: any) {
    if (error !== "cancel") {
      ElMessage.error(t('admin.userManagement.messages.batchDeleteFailed'))
    }
  }
}

/** 导出用户 */
const exportUsers = async () => {
  try {
    exportLoading.value = true
    const response = await exportUsersApi({
      search: searchQuery.value,
      status: statusFilter.value,
      role: roleFilter.value
    })

    // 创建下载链接
    const blob = new Blob([response.data], {
      type: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
    })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement("a")
    link.href = url
    link.download = `users-${new Date().toISOString().split('T')[0]}.xlsx`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)

    ElMessage.success(t('admin.userManagement.messages.exportSuccess'))
  } catch (error: any) {
    ElMessage.error(t('admin.userManagement.messages.exportFailed'))
  } finally {
    exportLoading.value = false
  }
}

/** 用户创建成功处理 */
const handleUserCreated = () => {
  refreshUsers()
}

/** 用户更新成功处理 */
const handleUserUpdated = () => {
  refreshUsers()
}

/** 组件挂载时加载数据 */
onMounted(() => {
  loadUsers()
  loadUserStats()
})
</script>

<style lang="scss" scoped>
.user-management {
  padding: 24px;
  background-color: var(--el-bg-color-page);
  min-height: calc(100vh - 60px);
}

.page-header {
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

.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;

  .search-section {
    display: flex;
    align-items: center;
  }

  .action-buttons {
    display: flex;
    gap: 12px;
  }
}

.stats-cards {
  margin-bottom: 24px;

  .stat-card {
    .stat-content {
      display: flex;
      align-items: center;

      .stat-icon {
        width: 48px;
        height: 48px;
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        margin-right: 16px;
        font-size: 24px;

        &.total {
          background: var(--el-color-primary-light-8);
          color: var(--el-color-primary);
        }

        &.active {
          background: var(--el-color-success-light-8);
          color: var(--el-color-success);
        }

        &.admin {
          background: var(--el-color-warning-light-8);
          color: var(--el-color-warning);
        }

        &.new {
          background: var(--el-color-info-light-8);
          color: var(--el-color-info);
        }
      }

      .stat-info {
        .stat-number {
          font-size: 24px;
          font-weight: 600;
          color: var(--el-text-color-primary);
          line-height: 1;
          margin-bottom: 4px;
        }

        .stat-label {
          font-size: 14px;
          color: var(--el-text-color-regular);
        }
      }
    }
  }
}

.user-table-card {
  margin-bottom: 24px;

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

  .user-info {
    .username {
      font-weight: 500;
      color: var(--el-text-color-primary);
    }

    .display-name {
      font-size: 12px;
      color: var(--el-text-color-secondary);
      margin-top: 2px;
    }
  }

  .email-info {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .pagination-container {
    display: flex;
    justify-content: center;
    margin-top: 24px;
  }
}

.batch-actions {
  position: fixed;
  bottom: 24px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 1000;

  .batch-content {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;

    .batch-buttons {
      display: flex;
      gap: 8px;
    }
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
  .user-management {
    padding: 16px;
  }

  .action-bar {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;

    .search-section {
      flex-wrap: wrap;
      gap: 8px;
    }

    .action-buttons {
      justify-content: center;
    }
  }

  .stats-cards {
    :deep(.el-col) {
      margin-bottom: 16px;
    }
  }
}
</style>