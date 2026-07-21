<template>
  <el-dialog
    v-model="visible"
    title="角色详情"
    width="600px"
    :before-close="handleClose"
  >
    <div v-if="role" class="role-detail">
      <!-- 基本信息 -->
      <div class="detail-section">
        <h3>基本信息</h3>
        <div class="info-grid">
          <div class="info-item">
            <label>角色名称:</label>
            <span>{{ role.name }}</span>
          </div>
          <div class="info-item">
            <label>角色类型:</label>
            <el-tag :type="role.type === 'system' ? 'info' : 'success'">
              {{ role.type === 'system' ? '系统角色' : '自定义角色' }}
            </el-tag>
          </div>
          <div class="info-item">
            <label>角色描述:</label>
            <span>{{ role.description || '暂无描述' }}</span>
          </div>
          <div class="info-item">
            <label>用户数量:</label>
            <span>{{ role.user_count }} 个用户</span>
          </div>
          <div class="info-item">
            <label>权限数量:</label>
            <span>{{ role.permission_count }} 个权限</span>
          </div>
          <div class="info-item">
            <label>创建时间:</label>
            <span>{{ formatDate(role.created_at) }}</span>
          </div>
          <div v-if="role.updated_at" class="info-item">
            <label>更新时间:</label>
            <span>{{ formatDate(role.updated_at) }}</span>
          </div>
        </div>
      </div>

      <!-- 权限列表 -->
      <div class="detail-section">
        <h3>权限列表</h3>
        <div v-if="role.main_permissions && role.main_permissions.length > 0" class="permissions-list">
          <el-tag
            v-for="permission in role.main_permissions"
            :key="permission"
            type="success"
            class="permission-tag"
          >
            {{ getPermissionDisplayName(permission) }}
          </el-tag>
        </div>
        <div v-else class="empty-permissions">
          <el-empty description="暂无权限" :image-size="80" />
        </div>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">关闭</el-button>
        <el-button
          v-if="role && role.type !== 'system'"
          type="primary"
          @click="handleEdit"
        >
          编辑角色
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import { computed } from "vue"

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

interface Props {
  modelValue: boolean
  role: Role | null
}

interface Emits {
  (e: "update:modelValue", value: boolean): void
  (e: "edit", role: Role): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const visible = computed({
  get: () => props.modelValue,
  set: (value) => emit("update:modelValue", value)
})

/** 格式化日期 */
const formatDate = (dateString: string) => {
  if (!dateString) return "未知"
  return new Date(dateString).toLocaleString("zh-CN")
}

/** 获取权限显示名称 */
const getPermissionDisplayName = (permission: string) => {
  const permissionNames: Record<string, string> = {
    "users:read": "查看用户",
    "users:write": "管理用户",
    "roles:read": "查看角色",
    "roles:write": "管理角色",
    "clusters:read": "查看集群",
    "clusters:write": "管理集群",
    "pods:read": "查看Pod",
    "pods:write": "管理Pod",
    "deployments:read": "查看Deployment",
    "deployments:write": "管理Deployment",
    "services:read": "查看Service",
    "services:write": "管理Service",
    "configmaps:read": "查看ConfigMap",
    "configmaps:write": "管理ConfigMap",
    "secrets:read": "查看Secret",
    "secrets:write": "管理Secret"
  }
  return permissionNames[permission] || permission
}

/** 关闭对话框 */
const handleClose = () => {
  visible.value = false
}

/** 编辑角色 */
const handleEdit = () => {
  if (props.role) {
    emit("edit", props.role)
    handleClose()
  }
}
</script>

<style lang="scss" scoped>
.role-detail {
  .detail-section {
    margin-bottom: 24px;
    
    &:last-child {
      margin-bottom: 0;
    }
    
    h3 {
      margin: 0 0 16px 0;
      font-size: 16px;
      font-weight: 600;
      color: var(--el-text-color-primary);
      border-bottom: 1px solid var(--el-border-color-light);
      padding-bottom: 8px;
    }
  }
  
  .info-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 16px;
    
    .info-item {
      display: flex;
      align-items: center;
      
      label {
        font-weight: 500;
        color: var(--el-text-color-regular);
        margin-right: 8px;
        min-width: 80px;
      }
      
      span {
        color: var(--el-text-color-primary);
      }
    }
  }
  
  .permissions-list {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    
    .permission-tag {
      margin: 0;
    }
  }
  
  .empty-permissions {
    text-align: center;
    padding: 20px;
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

@media (max-width: 768px) {
  .info-grid {
    grid-template-columns: 1fr;
  }
}
</style>