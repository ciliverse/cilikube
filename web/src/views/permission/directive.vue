<script lang="ts" setup>
import { checkPermission, checkResourcePermission, isAdmin, canEdit, isViewer } from "@/utils/permission"
import { usePermission } from "@/composables/usePermission"
import Permission from "@/components/Permission/index.vue"
import SwitchRoles from "./components/SwitchRoles.vue"

// 使用权限管理 Composable
const { currentRoleDisplayName, permissions } = usePermission()

// 资源权限测试数据
const resourcePermissionData = [
  { resource: 'pod' },
  { resource: 'deployment' },
  { resource: 'service' },
  { resource: 'secret' },
  { resource: 'configmap' },
  { resource: 'cluster-management' }
]
</script>

<template>
  <div class="app-container">
    <SwitchRoles />
    
    <!-- 当前用户信息 -->
    <div class="margin-top-30">
      <el-card header="当前用户权限信息">
        <el-descriptions :column="3" border>
          <el-descriptions-item label="当前角色">
            <el-tag type="primary">{{ currentRoleDisplayName }}</el-tag>
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
        </el-descriptions>
      </el-card>
    </div>

    <!-- 传统 v-permission 指令示例 -->
    <div class="margin-top-30">
      <el-card header="传统 v-permission 指令示例">
        <div>
          <el-tag v-permission="['admin']" type="success" size="large" effect="plain">
            这里采用了 v-permission="['admin']" 所以只有 admin 可以看见这句话
          </el-tag>
        </div>
        <div class="margin-top-15">
          <el-tag v-permission="['editor']" type="success" size="large" effect="plain">
            这里采用了 v-permission="['editor']" 所以只有 editor 可以看见这句话
          </el-tag>
        </div>
        <div class="margin-top-15">
          <el-tag v-permission="['admin', 'editor']" type="success" size="large" effect="plain">
            这里采用了 v-permission="['admin', 'editor']" 所以 admin 和 editor 都可以看见这句话
          </el-tag>
        </div>
        <div class="margin-top-15">
          <el-tag v-permission="['viewer']" type="info" size="large" effect="plain">
            这里采用了 v-permission="['viewer']" 所以只有 viewer 可以看见这句话
          </el-tag>
        </div>
      </el-card>
    </div>

    <!-- 新增角色权限指令示例 -->
    <div class="margin-top-30">
      <el-card header="新增角色权限指令示例">
        <el-space wrap>
          <el-button v-role-permission="'admin'" type="primary">
            管理员专用按钮
          </el-button>
          <el-button v-role-permission="'editor'" type="success">
            编辑者专用按钮
          </el-button>
          <el-button v-role-permission="'viewer'" type="info">
            查看者专用按钮
          </el-button>
          <el-button v-role-permission="['admin', 'editor']" type="warning">
            管理员和编辑者可用
          </el-button>
        </el-space>
      </el-card>
    </div>

    <!-- 资源权限指令示例 -->
    <div class="margin-top-30">
      <el-card header="资源权限指令示例">
        <el-space wrap>
          <el-button v-resource-permission:read="'pod'" type="primary">
            查看 Pod
          </el-button>
          <el-button v-resource-permission:write="'pod'" type="success">
            编辑 Pod
          </el-button>
          <el-button v-resource-permission:delete="'pod'" type="danger">
            删除 Pod
          </el-button>
          <el-button v-resource-permission="{ resource: 'secret', action: 'read' }" type="warning">
            查看 Secret
          </el-button>
          <el-button v-resource-permission="{ resource: 'cluster-management', action: 'write' }" type="info">
            集群管理
          </el-button>
        </el-space>
      </el-card>
    </div>

    <!-- Permission 组件示例 -->
    <div class="margin-top-30">
      <el-card header="Permission 组件示例">
        <Permission admin>
          <el-alert title="管理员专用功能区域" type="success" show-icon />
        </Permission>
        
        <Permission :roles="['admin', 'editor']" class="margin-top-15">
          <el-alert title="编辑者和管理员可见区域" type="info" show-icon />
        </Permission>
        
        <Permission resource="pod" action="write" class="margin-top-15">
          <el-alert title="可以编辑 Pod 的用户可见" type="warning" show-icon />
        </Permission>
        
        <Permission viewer class="margin-top-15">
          <el-alert title="查看者专用提示" type="error" show-icon />
        </Permission>
      </el-card>
    </div>

    <!-- checkPermission 函数示例 -->
    <div class="margin-top-30">
      <el-card header="权限检查函数示例">
        <el-tag type="warning" size="large">
          例如 Element Plus 的 el-tab-pane 或 el-table-column 以及其它动态渲染 Dom 的场景不适合使用
          v-permission，这种情况下你可以通过 v-if 和权限检查函数来实现：
        </el-tag>
        <el-tabs type="border-card" class="margin-top-15">
          <el-tab-pane v-if="checkPermission(['admin'])" label="管理员">
            这里采用了 <el-tag>v-if="checkPermission(['admin'])"</el-tag> 所以只有 admin 可以看见这句话
          </el-tab-pane>
          <el-tab-pane v-if="checkPermission(['editor'])" label="编辑者">
            这里采用了 <el-tag>v-if="checkPermission(['editor'])"</el-tag> 所以只有 editor 可以看见这句话
          </el-tab-pane>
          <el-tab-pane v-if="checkPermission(['viewer'])" label="查看者">
            这里采用了 <el-tag>v-if="checkPermission(['viewer'])"</el-tag> 所以只有 viewer 可以看见这句话
          </el-tab-pane>
          <el-tab-pane v-if="checkPermission(['admin', 'editor'])" label="管理员和编辑者">
            这里采用了 <el-tag>v-if="checkPermission(['admin', 'editor'])"</el-tag> 所以 admin 和 editor 都可以看见这句话
          </el-tab-pane>
        </el-tabs>
      </el-card>
    </div>

    <!-- 资源权限检查示例 -->
    <div class="margin-top-30">
      <el-card header="资源权限检查示例">
        <el-table :data="resourcePermissionData" border>
          <el-table-column prop="resource" label="资源" width="120" />
          <el-table-column label="读取权限" width="100">
            <template #default="{ row }">
              <el-tag :type="checkResourcePermission(row.resource, 'read') ? 'success' : 'danger'">
                {{ checkResourcePermission(row.resource, 'read') ? '✓' : '✗' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="写入权限" width="100">
            <template #default="{ row }">
              <el-tag :type="checkResourcePermission(row.resource, 'write') ? 'success' : 'danger'">
                {{ checkResourcePermission(row.resource, 'write') ? '✓' : '✗' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="删除权限" width="100">
            <template #default="{ row }">
              <el-tag :type="checkResourcePermission(row.resource, 'delete') ? 'success' : 'danger'">
                {{ checkResourcePermission(row.resource, 'delete') ? '✓' : '✗' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作">
            <template #default="{ row }">
              <el-space>
                <el-button v-if="checkResourcePermission(row.resource, 'read')" size="small" type="primary">
                  查看
                </el-button>
                <el-button v-if="checkResourcePermission(row.resource, 'write')" size="small" type="success">
                  编辑
                </el-button>
                <el-button v-if="checkResourcePermission(row.resource, 'delete')" size="small" type="danger">
                  删除
                </el-button>
              </el-space>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </div>

    <!-- 快速权限检查示例 -->
    <div class="margin-top-30">
      <el-card header="快速权限检查示例">
        <el-space direction="vertical" style="width: 100%;">
          <div v-if="isAdmin()">
            <el-alert title="管理员功能：您拥有系统的完全控制权限" type="success" show-icon />
          </div>
          
          <div v-if="canEdit()">
            <el-alert title="编辑权限：您可以修改大部分资源" type="info" show-icon />
          </div>
          
          <div v-if="isViewer()">
            <el-alert title="查看权限：您只能查看资源，无法进行修改" type="warning" show-icon />
          </div>
        </el-space>
      </el-card>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.margin-top-15 {
  margin-top: 15px;
}

.margin-top-30 {
  margin-top: 30px;
}
</style>
