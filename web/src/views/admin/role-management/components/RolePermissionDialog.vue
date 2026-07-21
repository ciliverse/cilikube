<template>
  <el-dialog
    v-model="dialogVisible"
    :title="`管理权限 - ${role?.display_name}`"
    width="800px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <div class="permissions-container">
      <!-- 权限统计 -->
      <div class="permissions-stats">
        <div class="stat-card">
          <div class="stat-icon">
            <el-icon><Key /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ selectedPermissions.length }}</div>
            <div class="stat-label">已分配权限</div>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">
            <el-icon><User /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ role?.user_count || 0 }}</div>
            <div class="stat-label">影响用户</div>
          </div>
        </div>
      </div>

      <!-- 权限搜索 -->
      <div class="permissions-search">
        <el-input
          v-model="searchQuery"
          placeholder="搜索权限..."
          :prefix-icon="Search"
          clearable
        />
        <el-button-group style="margin-left: 12px;">
          <el-button
            :type="showMode === 'all' ? 'primary' : 'default'"
            @click="showMode = 'all'"
          >
            全部权限
          </el-button>
          <el-button
            :type="showMode === 'assigned' ? 'primary' : 'default'"
            @click="showMode = 'assigned'"
          >
            已分配
          </el-button>
          <el-button
            :type="showMode === 'unassigned' ? 'primary' : 'default'"
            @click="showMode = 'unassigned'"
          >
            未分配
          </el-button>
        </el-button-group>
      </div>

      <!-- 权限分类 -->
      <div class="permissions-content">
        <div v-if="loading" class="loading-container">
          <el-icon class="loading-icon" :size="40">
            <Loading />
          </el-icon>
          <p>加载权限数据中...</p>
        </div>

        <div v-else class="permissions-categories">
          <div
            v-for="category in filteredCategories"
            :key="category.name"
            class="permission-category"
          >
            <div class="category-header">
              <div class="category-info">
                <el-icon class="category-icon">
                  <component :is="getCategoryIcon(category.name)" />
                </el-icon>
                <div class="category-details">
                  <h3>{{ category.display_name }}</h3>
                  <p>{{ category.description }}</p>
                </div>
              </div>
              <div class="category-actions">
                <el-button
                  size="small"
                  @click="toggleCategoryPermissions(category)"
                >
                  {{ isCategoryFullySelected(category) ? '取消全选' : '全选' }}
                </el-button>
              </div>
            </div>

            <div class="permissions-grid">
              <div
                v-for="permission in category.permissions"
                :key="permission.name"
                class="permission-item"
                :class="{ 'selected': isPermissionSelected(permission.name) }"
              >
                <el-checkbox
                  :model-value="isPermissionSelected(permission.name)"
                  @change="togglePermission(permission.name, $event)"
                  :disabled="role?.is_system && permission.system_required"
                >
                  <div class="permission-content">
                    <div class="permission-info">
                      <span class="permission-name">{{ permission.display_name }}</span>
                      <span class="permission-code">{{ permission.name }}</span>
                    </div>
                    <div class="permission-description">
                      {{ permission.description }}
                    </div>
                    <div v-if="permission.system_required" class="permission-tags">
                      <el-tag size="small" type="warning">系统必需</el-tag>
                    </div>
                  </div>
                </el-checkbox>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <div class="footer-info">
          <span>已选择 {{ selectedPermissions.length }} 个权限</span>
        </div>
        <div class="footer-actions">
          <el-button @click="handleClose">取消</el-button>
          <el-button
            type="primary"
            @click="handleSubmit"
            :loading="saving"
          >
            保存权限
          </el-button>
        </div>
      </div>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import { ref, computed, watch, onMounted } from "vue"
import { ElMessage } from "element-plus"
import {
  Search,
  Loading,
  Key,
  User,
  Lock,
  View,
  Edit,
  Delete,
  Setting,
  Monitor
} from "@element-plus/icons-vue"
import {
  getRolePermissionsApi,
  updateRolePermissionsApi,
  getAvailablePermissionsApi
} from "@/api/admin/role-management"

interface Role {
  id: number
  name: string
  display_name: string
  is_system: boolean
  user_count: number
}

interface Permission {
  name: string
  display_name: string
  description: string
  system_required?: boolean
}

interface PermissionCategory {
  name: string
  display_name: string
  description: string
  permissions: Permission[]
}

interface Props {
  modelValue: boolean
  role: Role | null
}

const props = defineProps<Props>()
const emit = defineEmits<{
  "update:modelValue": [value: boolean]
  updated: []
}>()

const loading = ref(false)
const saving = ref(false)
const dialogVisible = ref(false)
const searchQuery = ref("")
const showMode = ref<"all" | "assigned" | "unassigned">("all")
const selectedPermissions = ref<string[]>([])
const availableCategories = ref<PermissionCategory[]>([])

// 权限分类定义
const defaultCategories: PermissionCategory[] = [
  {
    name: "cluster",
    display_name: "集群管理",
    description: "Kubernetes集群相关操作权限",
    permissions: [
      { name: "read:clusters", display_name: "查看集群", description: "查看集群信息和状态" },
      { name: "write:clusters", display_name: "管理集群", description: "创建、更新、删除集群" },
      { name: "read:nodes", display_name: "查看节点", description: "查看集群节点信息" },
      { name: "write:nodes", display_name: "管理节点", description: "管理集群节点" }
    ]
  },
  {
    name: "workload",
    display_name: "工作负载",
    description: "Pod、Deployment等工作负载权限",
    permissions: [
      { name: "read:pods", display_name: "查看Pod", description: "查看Pod信息和日志" },
      { name: "write:pods", display_name: "管理Pod", description: "创建、更新、删除Pod" },
      { name: "read:deployments", display_name: "查看Deployment", description: "查看部署信息" },
      { name: "write:deployments", display_name: "管理Deployment", description: "管理部署" },
      { name: "exec:pods", display_name: "Pod终端", description: "进入Pod执行命令" }
    ]
  },
  {
    name: "network",
    display_name: "网络资源",
    description: "Service、Ingress等网络资源权限",
    permissions: [
      { name: "read:services", display_name: "查看Service", description: "查看服务信息" },
      { name: "write:services", display_name: "管理Service", description: "管理服务" },
      { name: "read:ingress", display_name: "查看Ingress", description: "查看入口规则" },
      { name: "write:ingress", display_name: "管理Ingress", description: "管理入口规则" }
    ]
  },
  {
    name: "storage",
    display_name: "存储资源",
    description: "PV、PVC、ConfigMap等存储权限",
    permissions: [
      { name: "read:configmaps", display_name: "查看ConfigMap", description: "查看配置映射" },
      { name: "write:configmaps", display_name: "管理ConfigMap", description: "管理配置映射" },
      { name: "read:secrets", display_name: "查看Secret", description: "查看密钥" },
      { name: "write:secrets", display_name: "管理Secret", description: "管理密钥" },
      { name: "read:pv", display_name: "查看PV", description: "查看持久卷" },
      { name: "write:pv", display_name: "管理PV", description: "管理持久卷" }
    ]
  },
  {
    name: "admin",
    display_name: "系统管理",
    description: "用户、角色等系统管理权限",
    permissions: [
      { name: "admin:users", display_name: "用户管理", description: "管理系统用户", system_required: true },
      { name: "admin:roles", display_name: "角色管理", description: "管理系统角色", system_required: true },
      { name: "admin:system", display_name: "系统设置", description: "管理系统配置", system_required: true },
      { name: "admin:audit", display_name: "审计日志", description: "查看审计日志" }
    ]
  }
]

/** 计算属性 */
const filteredCategories = computed(() => {
  let categories = availableCategories.value

  // 根据显示模式过滤
  if (showMode.value === "assigned") {
    categories = categories.map(category => ({
      ...category,
      permissions: category.permissions.filter(p => 
        selectedPermissions.value.includes(p.name)
      )
    })).filter(category => category.permissions.length > 0)
  } else if (showMode.value === "unassigned") {
    categories = categories.map(category => ({
      ...category,
      permissions: category.permissions.filter(p => 
        !selectedPermissions.value.includes(p.name)
      )
    })).filter(category => category.permissions.length > 0)
  }

  // 根据搜索查询过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    categories = categories.map(category => ({
      ...category,
      permissions: category.permissions.filter(p =>
        p.name.toLowerCase().includes(query) ||
        p.display_name.toLowerCase().includes(query) ||
        p.description.toLowerCase().includes(query)
      )
    })).filter(category => category.permissions.length > 0)
  }

  return categories
})

/** 获取分类图标 */
const getCategoryIcon = (categoryName: string) => {
  const icons: Record<string, string> = {
    cluster: "Monitor",
    workload: "Setting",
    network: "Link",
    storage: "FolderOpened",
    admin: "Lock"
  }
  return icons[categoryName] || "Setting"
}

/** 检查权限是否已选择 */
const isPermissionSelected = (permissionName: string) => {
  return selectedPermissions.value.includes(permissionName)
}

/** 检查分类是否全选 */
const isCategoryFullySelected = (category: PermissionCategory) => {
  return category.permissions.every(p => selectedPermissions.value.includes(p.name))
}

/** 切换权限选择 */
const togglePermission = (permissionName: string, selected: boolean) => {
  if (selected) {
    if (!selectedPermissions.value.includes(permissionName)) {
      selectedPermissions.value.push(permissionName)
    }
  } else {
    const index = selectedPermissions.value.indexOf(permissionName)
    if (index > -1) {
      selectedPermissions.value.splice(index, 1)
    }
  }
}

/** 切换分类权限 */
const toggleCategoryPermissions = (category: PermissionCategory) => {
  const isFullySelected = isCategoryFullySelected(category)
  
  category.permissions.forEach(permission => {
    if (isFullySelected) {
      // 取消全选
      const index = selectedPermissions.value.indexOf(permission.name)
      if (index > -1) {
        selectedPermissions.value.splice(index, 1)
      }
    } else {
      // 全选
      if (!selectedPermissions.value.includes(permission.name)) {
        selectedPermissions.value.push(permission.name)
      }
    }
  })
}

/** 加载角色权限 */
const loadRolePermissions = async () => {
  if (!props.role) return
  
  try {
    loading.value = true
    const { data } = await getRolePermissionsApi(props.role.id)
    selectedPermissions.value = data.permissions || []
  } catch (error: any) {
    console.error("加载角色权限失败:", error)
    ElMessage.error("加载角色权限失败")
  } finally {
    loading.value = false
  }
}

/** 加载可用权限 */
const loadAvailablePermissions = async () => {
  try {
    const { data } = await getAvailablePermissionsApi()
    availableCategories.value = data.categories || defaultCategories
  } catch (error: any) {
    console.error("加载可用权限失败:", error)
    // 使用默认权限分类
    availableCategories.value = defaultCategories
  }
}

/** 监听对话框显示状态 */
watch(() => props.modelValue, (val) => {
  dialogVisible.value = val
  if (val && props.role) {
    loadRolePermissions()
  }
})

watch(dialogVisible, (val) => {
  emit("update:modelValue", val)
})

/** 提交权限更改 */
const handleSubmit = async () => {
  if (!props.role) return
  
  try {
    saving.value = true
    
    await updateRolePermissionsApi(props.role.id, {
      permissions: selectedPermissions.value
    })
    
    ElMessage.success("权限更新成功")
    emit("updated")
    handleClose()
    
  } catch (error: any) {
    console.error("更新权限失败:", error)
    ElMessage.error(error.response?.data?.message || "更新权限失败")
  } finally {
    saving.value = false
  }
}

/** 关闭对话框 */
const handleClose = () => {
  dialogVisible.value = false
  searchQuery.value = ""
  showMode.value = "all"
}

/** 组件挂载时加载可用权限 */
onMounted(() => {
  loadAvailablePermissions()
})
</script>

<style lang="scss" scoped>
.permissions-container {
  .permissions-stats {
    display: flex;
    gap: 16px;
    margin-bottom: 24px;
    
    .stat-card {
      display: flex;
      align-items: center;
      padding: 16px;
      background: var(--el-fill-color-extra-light);
      border-radius: 8px;
      flex: 1;
      
      .stat-icon {
        width: 40px;
        height: 40px;
        border-radius: 8px;
        background: var(--el-color-primary-light-8);
        color: var(--el-color-primary);
        display: flex;
        align-items: center;
        justify-content: center;
        margin-right: 12px;
      }
      
      .stat-info {
        .stat-value {
          font-size: 20px;
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
  
  .permissions-search {
    display: flex;
    align-items: center;
    margin-bottom: 24px;
  }
  
  .permissions-content {
    max-height: 500px;
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
  }
}

.permissions-categories {
  .permission-category {
    margin-bottom: 24px;
    border: 1px solid var(--el-border-color-light);
    border-radius: 8px;
    overflow: hidden;
    
    .category-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 16px;
      background: var(--el-fill-color-extra-light);
      border-bottom: 1px solid var(--el-border-color-light);
      
      .category-info {
        display: flex;
        align-items: center;
        
        .category-icon {
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
        
        .category-details {
          h3 {
            margin: 0 0 4px 0;
            font-size: 16px;
            font-weight: 600;
            color: var(--el-text-color-primary);
          }
          
          p {
            margin: 0;
            font-size: 12px;
            color: var(--el-text-color-secondary);
          }
        }
      }
    }
    
    .permissions-grid {
      padding: 16px;
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
      gap: 12px;
      
      .permission-item {
        padding: 12px;
        border: 1px solid var(--el-border-color-lighter);
        border-radius: 6px;
        transition: all 0.3s ease;
        
        &:hover {
          border-color: var(--el-color-primary-light-7);
          background: var(--el-color-primary-light-9);
        }
        
        &.selected {
          border-color: var(--el-color-primary);
          background: var(--el-color-primary-light-9);
        }
        
        :deep(.el-checkbox) {
          width: 100%;
          
          .el-checkbox__label {
            width: 100%;
          }
        }
        
        .permission-content {
          .permission-info {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 8px;
            
            .permission-name {
              font-weight: 500;
              color: var(--el-text-color-primary);
            }
            
            .permission-code {
              font-size: 11px;
              color: var(--el-text-color-secondary);
              font-family: monospace;
              background: var(--el-fill-color-light);
              padding: 2px 6px;
              border-radius: 3px;
            }
          }
          
          .permission-description {
            font-size: 12px;
            color: var(--el-text-color-regular);
            line-height: 1.4;
            margin-bottom: 8px;
          }
          
          .permission-tags {
            display: flex;
            gap: 4px;
          }
        }
      }
    }
  }
}

.dialog-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  
  .footer-info {
    font-size: 14px;
    color: var(--el-text-color-regular);
  }
  
  .footer-actions {
    display: flex;
    gap: 12px;
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
</style>