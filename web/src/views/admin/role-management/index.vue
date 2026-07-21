<template>
  <div class="role-management">
    <div class="page-header">
      <h1>{{ t('menu.roleManagement') }}</h1>
      <p class="subtitle">{{ t('admin.roleManagement.description') }}</p>
    </div>

    <!-- 操作栏 -->
    <div class="action-bar">
      <div class="search-section">
        <el-input
          v-model="searchQuery"
          :placeholder="t('admin.roleManagement.searchPlaceholder')"
          :prefix-icon="Search"
          clearable
          style="width: 300px"
          @input="handleSearch"
        />
        
        <el-select
          v-model="typeFilter"
          :placeholder="t('admin.roleManagement.typeFilter')"
          clearable
          style="width: 120px; margin-left: 12px"
          @change="handleFilter"
        >
          <el-option :label="t('common.all')" value="" />
          <el-option :label="t('admin.roleManagement.types.system')" value="system" />
          <el-option :label="t('admin.roleManagement.types.custom')" value="custom" />
        </el-select>

        <el-button
          :icon="Refresh"
          @click="refreshRoles"
          :loading="loading"
          style="margin-left: 12px"
        >
          {{ t('common.refresh') }}
        </el-button>
      </div>

      <div class="action-buttons">
        <el-button
          type="primary"
          :icon="Plus"
          @click="showCreateDialog = true"
        >
          {{ t('admin.roleManagement.actions.createRole') }}
        </el-button>
        
        <el-button
          :icon="Setting"
          @click="showPermissionConfig = true"
        >
          {{ t('admin.roleManagement.actions.permissionConfig') }}
        </el-button>
      </div>
    </div>

    <!-- 角色卡片列表 -->
    <div v-if="loading" class="loading-container">
      <el-icon class="loading-icon" :size="40">
        <Loading />
      </el-icon>
      <p>{{ t('admin.roleManagement.loading') }}</p>
    </div>

    <div v-else class="roles-grid">
      <div
        v-for="role in filteredRoles"
        :key="role.id"
        class="role-card"
        :class="{ 'system-role': role.type === 'system' }"
      >
        <el-card>
          <div class="role-header">
            <div class="role-info">
              <div class="role-name">
                {{ role.name }}
                <el-tag
                  v-if="role.type === 'system'"
                  type="info"
                  size="small"
                >
                  {{ t('admin.roleManagement.types.system') }}
                </el-tag>
              </div>
              <div class="role-description">{{ role.description }}</div>
            </div>
            
            <div class="role-actions">
              <el-dropdown
                @command="(command) => handleRoleAction(command, role)"
                trigger="click"
              >
                <el-button size="small" type="info">
                  {{ t('common.actions') }}<el-icon class="el-icon--right"><ArrowDown /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="view">
                      {{ t('admin.roleManagement.actions.viewDetails') }}
                    </el-dropdown-item>
                    <el-dropdown-item
                      command="edit"
                      :disabled="role.type === 'system'"
                    >
                      {{ t('admin.roleManagement.actions.editRole') }}
                    </el-dropdown-item>
                    <el-dropdown-item command="permissions">
                      {{ t('admin.roleManagement.actions.managePermissions') }}
                    </el-dropdown-item>
                    <el-dropdown-item command="users">
                      {{ t('admin.roleManagement.actions.viewUsers') }}
                    </el-dropdown-item>
                    <el-dropdown-item
                      command="delete"
                      :disabled="role.type === 'system' || role.user_count > 0"
                      divided
                      style="color: var(--el-color-danger)"
                    >
                      {{ t('admin.roleManagement.actions.deleteRole') }}
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </div>

          <div class="role-stats">
            <div class="stat-item">
              <el-icon><User /></el-icon>
              <span>{{ t('admin.roleManagement.stats.users', { count: role.user_count }) }}</span>
            </div>
            <div class="stat-item">
              <el-icon><Key /></el-icon>
              <span>{{ t('admin.roleManagement.stats.permissions', { count: role.permission_count }) }}</span>
            </div>
          </div>

          <div class="role-permissions">
            <div class="permissions-header">
              <span>{{ t('admin.roleManagement.mainPermissions') }}</span>
              <el-button
                size="small"
                type="primary"
                link
                @click="viewRolePermissions(role)"
              >
                {{ t('admin.roleManagement.viewAll') }}
              </el-button>
            </div>
            <div class="permissions-list">
              <el-tag
                v-for="permission in (role.main_permissions || []).slice(0, 3)"
                :key="permission"
                size="small"
                type="success"
              >
                {{ getPermissionDisplayName(permission) }}
              </el-tag>
              <el-tag
                v-if="(role.main_permissions || []).length > 3"
                size="small"
                type="info"
              >
                {{ t('admin.roleManagement.morePermissions', { count: (role.main_permissions || []).length - 3 }) }}
              </el-tag>
            </div>
          </div>

          <div class="role-footer">
            <div class="role-meta">
              <span>{{ t('admin.roleManagement.createdAt') }}: {{ formatDate(role.created_at) }}</span>
              <span v-if="role.updated_at">{{ t('admin.roleManagement.updatedAt') }}: {{ formatDate(role.updated_at) }}</span>
            </div>
          </div>
        </el-card>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-if="!loading && filteredRoles.length === 0" class="empty-container">
      <el-empty :description="t('admin.roleManagement.noRoles')">
        <el-button type="primary" @click="showCreateDialog = true">
          {{ t('admin.roleManagement.createFirstRole') }}
        </el-button>
      </el-empty>
    </div>

    <!-- 创建角色对话框 -->
    <CreateRoleDialog
      v-model="showCreateDialog"
      @created="handleRoleCreated"
    />

    <!-- 编辑角色对话框 -->
    <EditRoleDialog
      v-model="showEditDialog"
      :role="selectedRole"
      @updated="handleRoleUpdated"
    />

    <!-- 角色详情对话框 -->
    <RoleDetailDialog
      v-model="showDetailDialog"
      :role="selectedRole"
    />

    <!-- 权限管理对话框 -->
    <RolePermissionDialog
      v-model="showPermissionDialog"
      :role="selectedRole"
      @updated="handleRoleUpdated"
    />

    <!-- 角色用户列表对话框 -->
    <RoleUsersDialog
      v-model="showUsersDialog"
      :role="selectedRole"
    />

    <!-- 权限配置对话框 -->
    <PermissionConfigDialog
      v-model="showPermissionConfig"
      @updated="refreshRoles"
    />
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
  Setting,
  Loading,
  ArrowDown,
  User,
  Key
} from "@element-plus/icons-vue"
import {
  getRolesApi,
  deleteRoleApi
} from "@/api/admin/role-management"
import CreateRoleDialog from "./components/CreateRoleDialog.vue"
import EditRoleDialog from "./components/EditRoleDialog.vue"
import RoleDetailDialog from "./components/RoleDetailDialog.vue"
import RolePermissionDialog from "./components/RolePermissionDialog.vue"
import RoleUsersDialog from "./components/RoleUsersDialog.vue"
import PermissionConfigDialog from "./components/PermissionConfigDialog.vue"

const { t } = useI18n()

interface Role {
  id: number
  name: string
  description: string
  type: "system" | "custom"
  user_count: number
  permission_count: number
  main_permissions: string[]
  created_at: string
  updated_at?: string
}

// 响应式数据
const loading = ref(false)
const roles = ref<Role[]>([])
const selectedRole = ref<Role | null>(null)

// 对话框状态
const showCreateDialog = ref(false)
const showEditDialog = ref(false)
const showDetailDialog = ref(false)
const showPermissionDialog = ref(false)
const showUsersDialog = ref(false)
const showPermissionConfig = ref(false)

// 搜索和筛选
const searchQuery = ref("")
const typeFilter = ref("")

// 计算属性
const filteredRoles = computed(() => {
  let filtered = roles.value

  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(role =>
      role.name.toLowerCase().includes(query) ||
      role.description.toLowerCase().includes(query)
    )
  }

  // 类型过滤
  if (typeFilter.value) {
    filtered = filtered.filter(role => role.type === typeFilter.value)
  }

  return filtered
})

/** 格式化日期 */
const formatDate = (dateString: string) => {
  if (!dateString) return t('admin.roleManagement.unknown')
  return new Date(dateString).toLocaleDateString()
}

/** 获取权限显示名称 */
const getPermissionDisplayName = (permission: string) => {
  const permissionNames: Record<string, string> = {
    "users:read": t('admin.permissions.users.read'),
    "users:write": t('admin.permissions.users.write'),
    "roles:read": t('admin.permissions.roles.read'),
    "roles:write": t('admin.permissions.roles.write'),
    "clusters:read": t('admin.permissions.clusters.read'),
    "clusters:write": t('admin.permissions.clusters.write'),
    "pods:read": t('admin.permissions.pods.read'),
    "pods:write": t('admin.permissions.pods.write'),
    "deployments:read": t('admin.permissions.deployments.read'),
    "deployments:write": t('admin.permissions.deployments.write'),
    "services:read": t('admin.permissions.services.read'),
    "services:write": t('admin.permissions.services.write'),
    "configmaps:read": t('admin.permissions.configmaps.read'),
    "configmaps:write": t('admin.permissions.configmaps.write'),
    "secrets:read": t('admin.permissions.secrets.read'),
    "secrets:write": t('admin.permissions.secrets.write')
  }
  return permissionNames[permission] || permission
}

/** 加载角色列表 */
const loadRoles = async () => {
  try {
    loading.value = true
    const { data } = await getRolesApi()
    roles.value = data
  } catch (error: any) {
    console.error("加载角色列表失败:", error)
    ElMessage.error(t('admin.roleManagement.messages.loadFailed'))
  } finally {
    loading.value = false
  }
}

/** 搜索处理 */
const handleSearch = () => {
  // 实时搜索，无需额外处理
}

/** 筛选处理 */
const handleFilter = () => {
  // 实时筛选，无需额外处理
}

/** 刷新角色列表 */
const refreshRoles = () => {
  loadRoles()
}

/** 角色操作处理 */
const handleRoleAction = async (command: string, role: Role) => {
  selectedRole.value = role
  
  switch (command) {
    case 'view':
      showDetailDialog.value = true
      break
    case 'edit':
      showEditDialog.value = true
      break
    case 'permissions':
      showPermissionDialog.value = true
      break
    case 'users':
      showUsersDialog.value = true
      break
    case 'delete':
      await deleteRole(role)
      break
  }
}

/** 查看角色权限 */
const viewRolePermissions = (role: Role) => {
  selectedRole.value = role
  showPermissionDialog.value = true
}

/** 删除角色 */
const deleteRole = async (role: Role) => {
  try {
    await ElMessageBox.confirm(
      t('admin.roleManagement.messages.confirmDelete', { name: role.name }),
      t('admin.roleManagement.messages.confirmDeleteTitle'),
      {
        confirmButtonText: t('admin.roleManagement.messages.confirmDeleteButton'),
        cancelButtonText: t('common.cancel'),
        type: "error"
      }
    )

    await deleteRoleApi(role.id)
    ElMessage.success(t('admin.roleManagement.messages.deleteSuccess'))
    refreshRoles()
  } catch (error: any) {
    if (error !== "cancel") {
      console.error("删除角色失败:", error)
      ElMessage.error(error.response?.data?.message || t('admin.roleManagement.messages.deleteFailed'))
    }
  }
}

/** 角色创建成功处理 */
const handleRoleCreated = () => {
  refreshRoles()
}

/** 角色更新成功处理 */
const handleRoleUpdated = () => {
  refreshRoles()
}

/** 组件挂载时加载数据 */
onMounted(() => {
  loadRoles()
})
</script>

<style lang="scss" scoped>
.role-management {
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

.roles-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: 24px;
  
  .role-card {
    transition: all 0.3s ease;
    
    &:hover {
      transform: translateY(-2px);
      
      :deep(.el-card) {
        box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
      }
    }
    
    &.system-role {
      :deep(.el-card) {
        border-left: 4px solid var(--el-color-primary);
      }
    }
    
    .role-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 16px;
      
      .role-info {
        flex: 1;
        
        .role-name {
          display: flex;
          align-items: center;
          gap: 8px;
          font-size: 18px;
          font-weight: 600;
          color: var(--el-text-color-primary);
          margin-bottom: 4px;
        }
        
        .role-description {
          font-size: 14px;
          color: var(--el-text-color-regular);
          line-height: 1.5;
        }
      }
      
      .role-actions {
        margin-left: 12px;
      }
    }
    
    .role-stats {
      display: flex;
      gap: 16px;
      margin-bottom: 16px;
      
      .stat-item {
        display: flex;
        align-items: center;
        gap: 4px;
        font-size: 14px;
        color: var(--el-text-color-regular);
        
        .el-icon {
          color: var(--el-color-primary);
        }
      }
    }
    
    .role-permissions {
      margin-bottom: 16px;
      
      .permissions-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 8px;
        
        span {
          font-size: 14px;
          font-weight: 500;
          color: var(--el-text-color-primary);
        }
      }
      
      .permissions-list {
        display: flex;
        flex-wrap: wrap;
        gap: 4px;
      }
    }
    
    .role-footer {
      .role-meta {
        display: flex;
        flex-direction: column;
        gap: 2px;
        font-size: 12px;
        color: var(--el-text-color-secondary);
      }
    }
  }
}

.empty-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 300px;
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
  .role-management {
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
  
  .roles-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }
}
</style>