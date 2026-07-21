<template>
  <el-dialog
    v-model="visible"
    title="权限配置"
    width="800px"
    :before-close="handleClose"
  >
    <div class="permission-config">
      <div class="config-header">
        <p class="description">
          在这里可以配置系统的权限项目，这些权限可以分配给不同的角色。
        </p>
      </div>

      <!-- 权限分类 -->
      <div class="permission-categories">
        <div
          v-for="category in permissionCategories"
          :key="category.key"
          class="category-section"
        >
          <div class="category-header">
            <h3>{{ category.name }}</h3>
            <span class="category-count">{{ category.permissions.length }} 个权限</span>
          </div>
          
          <div class="permissions-grid">
            <div
              v-for="permission in category.permissions"
              :key="permission.key"
              class="permission-item"
            >
              <div class="permission-info">
                <div class="permission-name">{{ permission.name }}</div>
                <div class="permission-description">{{ permission.description }}</div>
                <el-tag size="small" :type="getPermissionTypeColor(permission.type)">
                  {{ permission.type === 'read' ? '只读' : '读写' }}
                </el-tag>
              </div>
              
              <div class="permission-actions">
                <el-switch
                  v-model="permission.enabled"
                  @change="handlePermissionToggle(permission)"
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" @click="handleSave" :loading="saving">
          保存配置
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from "vue"
import { ElMessage } from "element-plus"

interface Permission {
  key: string
  name: string
  description: string
  type: "read" | "write"
  enabled: boolean
}

interface PermissionCategory {
  key: string
  name: string
  permissions: Permission[]
}

interface Props {
  modelValue: boolean
}

interface Emits {
  (e: "update:modelValue", value: boolean): void
  (e: "updated"): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit("update:modelValue", value)
})

const saving = ref(false)

// 权限分类数据
const permissionCategories = ref<PermissionCategory[]>([
  {
    key: "users",
    name: "用户管理",
    permissions: [
      {
        key: "users:read",
        name: "查看用户",
        description: "查看用户列表和用户详情",
        type: "read",
        enabled: true
      },
      {
        key: "users:write",
        name: "管理用户",
        description: "创建、编辑、删除用户",
        type: "write",
        enabled: true
      }
    ]
  },
  {
    key: "roles",
    name: "角色管理",
    permissions: [
      {
        key: "roles:read",
        name: "查看角色",
        description: "查看角色列表和角色详情",
        type: "read",
        enabled: true
      },
      {
        key: "roles:write",
        name: "管理角色",
        description: "创建、编辑、删除角色",
        type: "write",
        enabled: true
      }
    ]
  },
  {
    key: "clusters",
    name: "集群管理",
    permissions: [
      {
        key: "clusters:read",
        name: "查看集群",
        description: "查看集群信息和状态",
        type: "read",
        enabled: true
      },
      {
        key: "clusters:write",
        name: "管理集群",
        description: "创建、配置、删除集群",
        type: "write",
        enabled: true
      }
    ]
  },
  {
    key: "workloads",
    name: "工作负载",
    permissions: [
      {
        key: "pods:read",
        name: "查看Pod",
        description: "查看Pod列表和详情",
        type: "read",
        enabled: true
      },
      {
        key: "pods:write",
        name: "管理Pod",
        description: "创建、删除、重启Pod",
        type: "write",
        enabled: true
      },
      {
        key: "deployments:read",
        name: "查看Deployment",
        description: "查看Deployment列表和详情",
        type: "read",
        enabled: true
      },
      {
        key: "deployments:write",
        name: "管理Deployment",
        description: "创建、更新、删除Deployment",
        type: "write",
        enabled: true
      }
    ]
  },
  {
    key: "services",
    name: "服务网络",
    permissions: [
      {
        key: "services:read",
        name: "查看Service",
        description: "查看Service列表和详情",
        type: "read",
        enabled: true
      },
      {
        key: "services:write",
        name: "管理Service",
        description: "创建、更新、删除Service",
        type: "write",
        enabled: true
      }
    ]
  },
  {
    key: "configs",
    name: "配置管理",
    permissions: [
      {
        key: "configmaps:read",
        name: "查看ConfigMap",
        description: "查看ConfigMap列表和内容",
        type: "read",
        enabled: true
      },
      {
        key: "configmaps:write",
        name: "管理ConfigMap",
        description: "创建、更新、删除ConfigMap",
        type: "write",
        enabled: true
      },
      {
        key: "secrets:read",
        name: "查看Secret",
        description: "查看Secret列表（不包含敏感内容）",
        type: "read",
        enabled: true
      },
      {
        key: "secrets:write",
        name: "管理Secret",
        description: "创建、更新、删除Secret",
        type: "write",
        enabled: true
      }
    ]
  }
])

/** 获取权限类型颜色 */
const getPermissionTypeColor = (type: string) => {
  return type === "read" ? "info" : "warning"
}

/** 权限开关切换 */
const handlePermissionToggle = (permission: Permission) => {
  console.log(`权限 ${permission.key} 状态变更为: ${permission.enabled}`)
}

/** 保存配置 */
const handleSave = async () => {
  try {
    saving.value = true
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    ElMessage.success("权限配置保存成功")
    emit("updated")
    handleClose()
  } catch (error: any) {
    console.error("保存权限配置失败:", error)
    ElMessage.error("保存权限配置失败")
  } finally {
    saving.value = false
  }
}

/** 关闭对话框 */
const handleClose = () => {
  visible.value = false
}

/** 加载权限配置 */
const loadPermissionConfig = async () => {
  try {
    // 这里可以从API加载实际的权限配置
    console.log("加载权限配置")
  } catch (error: any) {
    console.error("加载权限配置失败:", error)
    ElMessage.error("加载权限配置失败")
  }
}

onMounted(() => {
  loadPermissionConfig()
})
</script>

<style lang="scss" scoped>
.permission-config {
  .config-header {
    margin-bottom: 24px;
    
    .description {
      margin: 0;
      color: var(--el-text-color-regular);
      line-height: 1.6;
    }
  }
  
  .permission-categories {
    .category-section {
      margin-bottom: 32px;
      
      &:last-child {
        margin-bottom: 0;
      }
      
      .category-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 16px;
        padding-bottom: 8px;
        border-bottom: 1px solid var(--el-border-color-light);
        
        h3 {
          margin: 0;
          font-size: 16px;
          font-weight: 600;
          color: var(--el-text-color-primary);
        }
        
        .category-count {
          font-size: 14px;
          color: var(--el-text-color-secondary);
        }
      }
      
      .permissions-grid {
        display: grid;
        gap: 12px;
        
        .permission-item {
          display: flex;
          justify-content: space-between;
          align-items: center;
          padding: 16px;
          background: var(--el-bg-color-page);
          border: 1px solid var(--el-border-color-lighter);
          border-radius: 6px;
          transition: all 0.3s ease;
          
          &:hover {
            border-color: var(--el-color-primary-light-7);
            background: var(--el-color-primary-light-9);
          }
          
          .permission-info {
            flex: 1;
            
            .permission-name {
              font-size: 14px;
              font-weight: 500;
              color: var(--el-text-color-primary);
              margin-bottom: 4px;
            }
            
            .permission-description {
              font-size: 13px;
              color: var(--el-text-color-regular);
              margin-bottom: 8px;
              line-height: 1.4;
            }
          }
          
          .permission-actions {
            margin-left: 16px;
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
</style>