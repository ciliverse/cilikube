<template>
  <div class="permission-example">
    <el-card header="权限控制示例">
      <!-- 使用权限组件 -->
      <div class="example-section">
        <h4>权限组件示例</h4>
        
        <!-- 只有管理员可以看到 -->
        <Permission admin>
          <el-alert title="管理员专用功能" type="success" show-icon />
        </Permission>
        
        <!-- 编辑者和管理员可以看到 -->
        <Permission :roles="['admin', 'editor']">
          <el-alert title="编辑者和管理员可见" type="info" show-icon />
        </Permission>
        
        <!-- 资源权限检查 -->
        <Permission resource="pod" action="write">
          <el-alert title="可以编辑 Pod" type="warning" show-icon />
        </Permission>
        
        <!-- 自定义权限检查 -->
        <Permission :check="() => canAccessSpecialFeature">
          <el-alert title="特殊功能访问权限" type="error" show-icon />
        </Permission>
      </div>

      <!-- 使用权限指令 -->
      <div class="example-section">
        <h4>权限指令示例</h4>
        
        <!-- 角色权限指令 -->
        <el-button v-role-permission="'admin'" type="primary">
          管理员按钮
        </el-button>
        
        <el-button v-role-permission="['admin', 'editor']" type="success">
          编辑者按钮
        </el-button>
        
        <!-- 资源权限指令 -->
        <el-button v-resource-permission:write="'deployment'" type="warning">
          编辑 Deployment
        </el-button>
        
        <el-button v-resource-permission="{ resource: 'secret', action: 'read' }" type="danger">
          查看 Secret
        </el-button>
        
        <!-- 传统权限指令 -->
        <el-button v-permission="['admin']" type="info">
          传统权限指令
        </el-button>
      </div>

      <!-- 使用权限 Composable -->
      <div class="example-section">
        <h4>权限 Composable 示例</h4>
        
        <el-descriptions :column="2" border>
          <el-descriptions-item label="当前角色">
            {{ currentRoleDisplayName }}
          </el-descriptions-item>
          <el-descriptions-item label="是否管理员">
            <el-tag :type="permissions.isAdmin ? 'success' : 'info'">
              {{ permissions.isAdmin ? '是' : '否' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="可以编辑">
            <el-tag :type="permissions.canEdit ? 'success' : 'info'">
              {{ permissions.canEdit ? '是' : '否' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="只能查看">
            <el-tag :type="permissions.isViewer ? 'warning' : 'info'">
              {{ permissions.isViewer ? '是' : '否' }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>
        
        <div style="margin-top: 16px;">
          <h5>资源操作权限</h5>
          <el-space wrap>
            <el-tag 
              v-for="resource in testResources" 
              :key="resource"
              :type="hasResourcePermission(resource, 'read') ? 'success' : 'danger'"
            >
              {{ resource }} 读取: {{ hasResourcePermission(resource, 'read') ? '✓' : '✗' }}
            </el-tag>
            <el-tag 
              v-for="resource in testResources" 
              :key="`${resource}-write`"
              :type="hasResourcePermission(resource, 'write') ? 'success' : 'danger'"
            >
              {{ resource }} 写入: {{ hasResourcePermission(resource, 'write') ? '✓' : '✗' }}
            </el-tag>
          </el-space>
        </div>
        
        <div style="margin-top: 16px;">
          <h5>动态操作按钮</h5>
          <el-space>
            <el-button 
              v-for="button in getActionButtons('pod')" 
              :key="button.action"
              :type="button.type"
              size="small"
            >
              {{ button.label }}
            </el-button>
          </el-space>
        </div>
      </div>

      <!-- 权限检查函数示例 -->
      <div class="example-section">
        <h4>权限检查函数示例</h4>
        
        <el-space direction="vertical" style="width: 100%;">
          <div v-if="hasRole('admin')">
            <el-alert title="管理员功能区域" type="success" />
          </div>
          
          <div v-if="hasRole(['admin', 'editor'])">
            <el-alert title="编辑功能区域" type="info" />
          </div>
          
          <div v-if="canAccessAdmin">
            <el-button type="primary" @click="goToAdmin">
              进入系统管理
            </el-button>
          </div>
          
          <div v-if="canEditResource">
            <el-button type="warning">
              编辑资源
            </el-button>
          </div>
          
          <div v-if="isReadOnly">
            <el-alert title="您当前只有查看权限" type="warning" />
          </div>
        </el-space>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { usePermission } from '@/composables/usePermission'
import Permission from '@/components/Permission/index.vue'

const router = useRouter()

// 使用权限管理 Composable
const {
  currentRoleDisplayName,
  permissions,
  canAccessAdmin,
  canEditResource,
  isReadOnly,
  hasRole,
  hasResourcePermission,
  getActionButtons
} = usePermission()

// 测试资源列表
const testResources = ['pod', 'deployment', 'service', 'secret', 'configmap']

// 自定义权限检查示例
const canAccessSpecialFeature = computed(() => {
  // 这里可以实现复杂的权限逻辑
  return permissions.value.isAdmin || (permissions.value.canEdit && hasResourcePermission('pod', 'write'))
})

const goToAdmin = () => {
  router.push('/admin')
}
</script>

<style scoped>
.permission-example {
  padding: 20px;
}

.example-section {
  margin-bottom: 24px;
  padding: 16px;
  border: 1px solid var(--el-border-color-light);
  border-radius: 4px;
}

.example-section h4 {
  margin-top: 0;
  margin-bottom: 16px;
  color: var(--el-text-color-primary);
}

.example-section h5 {
  margin-bottom: 8px;
  color: var(--el-text-color-regular);
}
</style>