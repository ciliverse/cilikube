<template>
  <el-dialog
    v-model="dialogVisible"
    title="编辑角色"
    width="600px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-width="100px"
    >
      <el-form-item label="角色名称" prop="name">
        <el-input
          v-model="formData.name"
          placeholder="请输入角色名称"
          :disabled="role?.type === 'system'"
          maxlength="50"
          show-word-limit
        >
          <template v-if="role?.type === 'system'" #suffix>
            <el-tooltip content="系统角色名称不可修改" placement="top">
              <el-icon><Lock /></el-icon>
            </el-tooltip>
          </template>
        </el-input>
        <div class="form-help">
          角色名称用于标识角色，建议使用英文和下划线
        </div>
      </el-form-item>

      <el-form-item label="显示名称" prop="display_name">
        <el-input
          v-model="formData.display_name"
          placeholder="请输入显示名称"
          maxlength="100"
          show-word-limit
        />
        <div class="form-help">
          显示名称用于界面展示，支持中文
        </div>
      </el-form-item>

      <el-form-item label="角色描述" prop="description">
        <el-input
          v-model="formData.description"
          type="textarea"
          :rows="3"
          placeholder="请输入角色描述"
          maxlength="200"
          show-word-limit
        />
      </el-form-item>

      <el-form-item label="角色类型">
        <el-tag :type="role?.type === 'system' ? 'info' : 'success'">
          {{ role?.type === 'system' ? '系统角色' : '自定义角色' }}
        </el-tag>
        <div class="form-help">
          系统角色的某些属性不可修改
        </div>
      </el-form-item>

      <el-form-item label="用户数量">
        <span class="info-text">{{ role?.user_count || 0 }} 个用户</span>
        <el-button
          v-if="role && role.user_count > 0"
          size="small"
          type="primary"
          link
          @click="viewRoleUsers"
        >
          查看用户列表
        </el-button>
      </el-form-item>

      <el-form-item label="权限配置">
        <div class="permissions-section">
          <div class="permission-tabs">
            <el-tabs v-model="activePermissionTab">
              <el-tab-pane
                v-for="category in permissionCategories"
                :key="category.key"
                :label="category.label"
                :name="category.key"
              >
                <div class="permission-category">
                  <div class="category-header">
                    <el-checkbox
                      :model-value="isCategoryAllSelected(category.key)"
                      :indeterminate="isCategoryIndeterminate(category.key)"
                      @change="handleCategorySelectAll(category.key, $event)"
                    >
                      全选
                    </el-checkbox>
                  </div>
                  
                  <div class="permission-grid">
                    <div
                      v-for="permission in category.permissions"
                      :key="permission.key"
                      class="permission-item"
                    >
                      <el-checkbox
                        v-model="formData.permissions"
                        :label="permission.key"
                      >
                        <div class="permission-info">
                          <span class="permission-name">{{ permission.name }}</span>
                          <span class="permission-desc">{{ permission.description }}</span>
                        </div>
                      </el-checkbox>
                    </div>
                  </div>
                </div>
              </el-tab-pane>
            </el-tabs>
          </div>
        </div>
      </el-form-item>

      <el-form-item label="创建时间">
        <span class="info-text">{{ formatDate(role?.created_at) }}</span>
      </el-form-item>

      <el-form-item label="更新时间">
        <span class="info-text">{{ formatDate(role?.updated_at) }}</span>
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button
          type="primary"
          @click="handleSubmit"
          :loading="loading"
        >
          保存更改
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import { ref, reactive, watch } from "vue"
import { ElMessage, type FormInstance, type FormRules } from "element-plus"
import { Lock } from "@element-plus/icons-vue"
import { updateRoleApi, getRolePermissionsApi } from "@/api/admin/role-management"

interface Role {
  id: number
  name: string
  display_name?: string
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

const props = defineProps<Props>()
const emit = defineEmits<{
  "update:modelValue": [value: boolean]
  updated: []
  viewUsers: [role: Role]
}>()

const formRef = ref<FormInstance>()
const loading = ref(false)
const dialogVisible = ref(false)
const activePermissionTab = ref("user")

const formData = reactive({
  name: "",
  display_name: "",
  description: "",
  permissions: [] as string[]
})

// 权限分类（与创建对话框相同）
const permissionCategories = [
  {
    key: "user",
    label: "用户管理",
    permissions: [
      { key: "users:read", name: "查看用户", description: "查看用户列表和详情" },
      { key: "users:write", name: "管理用户", description: "创建、编辑、删除用户" },
      { key: "users:password", name: "重置密码", description: "重置用户密码" },
      { key: "users:status", name: "用户状态", description: "激活或停用用户账户" }
    ]
  },
  {
    key: "role",
    label: "角色权限",
    permissions: [
      { key: "roles:read", name: "查看角色", description: "查看角色列表和权限" },
      { key: "roles:write", name: "管理角色", description: "创建、编辑、删除角色" },
      { key: "permissions:assign", name: "分配权限", description: "为用户分配角色和权限" }
    ]
  },
  {
    key: "cluster",
    label: "集群管理",
    permissions: [
      { key: "clusters:read", name: "查看集群", description: "查看集群信息和状态" },
      { key: "clusters:write", name: "管理集群", description: "添加、编辑、删除集群" },
      { key: "clusters:connect", name: "连接集群", description: "连接和切换集群" }
    ]
  },
  {
    key: "workload",
    label: "工作负载",
    permissions: [
      { key: "pods:read", name: "查看Pod", description: "查看Pod列表和详情" },
      { key: "pods:write", name: "管理Pod", description: "创建、编辑、删除Pod" },
      { key: "pods:exec", name: "Pod终端", description: "进入Pod容器执行命令" },
      { key: "pods:logs", name: "Pod日志", description: "查看Pod日志" },
      { key: "deployments:read", name: "查看Deployment", description: "查看Deployment列表和详情" },
      { key: "deployments:write", name: "管理Deployment", description: "创建、编辑、删除Deployment" }
    ]
  },
  {
    key: "network",
    label: "网络资源",
    permissions: [
      { key: "services:read", name: "查看Service", description: "查看Service列表和详情" },
      { key: "services:write", name: "管理Service", description: "创建、编辑、删除Service" },
      { key: "ingress:read", name: "查看Ingress", description: "查看Ingress列表和详情" },
      { key: "ingress:write", name: "管理Ingress", description: "创建、编辑、删除Ingress" }
    ]
  },
  {
    key: "config",
    label: "配置管理",
    permissions: [
      { key: "configmaps:read", name: "查看ConfigMap", description: "查看ConfigMap列表和详情" },
      { key: "configmaps:write", name: "管理ConfigMap", description: "创建、编辑、删除ConfigMap" },
      { key: "secrets:read", name: "查看Secret", description: "查看Secret列表和详情" },
      { key: "secrets:write", name: "管理Secret", description: "创建、编辑、删除Secret" }
    ]
  }
]

const formRules: FormRules = {
  name: [
    { required: true, message: "请输入角色名称", trigger: "blur" },
    { min: 2, max: 50, message: "角色名称长度在 2 到 50 个字符", trigger: "blur" },
    { 
      pattern: /^[a-zA-Z0-9_-]+$/, 
      message: "角色名称只能包含字母、数字、下划线和连字符", 
      trigger: "blur" 
    }
  ],
  display_name: [
    { required: true, message: "请输入显示名称", trigger: "blur" },
    { max: 100, message: "显示名称不能超过100个字符", trigger: "blur" }
  ],
  description: [
    { max: 200, message: "描述不能超过200个字符", trigger: "blur" }
  ]
}

/** 格式化日期 */
const formatDate = (dateString: string | undefined) => {
  if (!dateString) return "未知"
  return new Date(dateString).toLocaleString("zh-CN")
}

/** 检查分类是否全选 */
const isCategoryAllSelected = (categoryKey: string) => {
  const category = permissionCategories.find(cat => cat.key === categoryKey)
  if (!category) return false
  
  return category.permissions.every(permission => 
    formData.permissions.includes(permission.key)
  )
}

/** 检查分类是否部分选中 */
const isCategoryIndeterminate = (categoryKey: string) => {
  const category = permissionCategories.find(cat => cat.key === categoryKey)
  if (!category) return false
  
  const selectedCount = category.permissions.filter(permission => 
    formData.permissions.includes(permission.key)
  ).length
  
  return selectedCount > 0 && selectedCount < category.permissions.length
}

/** 处理分类全选 */
const handleCategorySelectAll = (categoryKey: string, checked: boolean) => {
  const category = permissionCategories.find(cat => cat.key === categoryKey)
  if (!category) return
  
  if (checked) {
    // 添加该分类下的所有权限
    category.permissions.forEach(permission => {
      if (!formData.permissions.includes(permission.key)) {
        formData.permissions.push(permission.key)
      }
    })
  } else {
    // 移除该分类下的所有权限
    category.permissions.forEach(permission => {
      const index = formData.permissions.indexOf(permission.key)
      if (index > -1) {
        formData.permissions.splice(index, 1)
      }
    })
  }
}

/** 查看角色用户 */
const viewRoleUsers = () => {
  if (props.role) {
    emit("viewUsers", props.role)
  }
}

/** 加载角色权限 */
const loadRolePermissions = async () => {
  if (!props.role) return
  
  try {
    const { data } = await getRolePermissionsApi(props.role.id)
    formData.permissions = data.permissions || []
  } catch (error: any) {
    console.error("加载角色权限失败:", error)
    // 使用主要权限作为备选
    formData.permissions = props.role.main_permissions || []
  }
}

/** 监听对话框显示状态 */
watch(() => props.modelValue, (val) => {
  dialogVisible.value = val
})

watch(dialogVisible, (val) => {
  emit("update:modelValue", val)
})

/** 监听角色数据变化 */
watch(() => props.role, (role) => {
  if (role) {
    formData.name = role.name
    formData.display_name = role.display_name || ""
    formData.description = role.description
    loadRolePermissions()
  }
}, { immediate: true })

/** 提交表单 */
const handleSubmit = async () => {
  if (!formRef.value || !props.role) return
  
  try {
    await formRef.value.validate()
    
    if (formData.permissions.length === 0) {
      ElMessage.warning("请至少选择一个权限")
      return
    }
    
    loading.value = true
    
    await updateRoleApi(props.role.id, {
      name: formData.name,
      display_name: formData.display_name,
      description: formData.description,
      permissions: formData.permissions
    })
    
    ElMessage.success("角色更新成功")
    emit("updated")
    handleClose()
    
  } catch (error: any) {
    console.error("更新角色失败:", error)
    ElMessage.error(error.response?.data?.message || "更新角色失败")
  } finally {
    loading.value = false
  }
}

/** 关闭对话框 */
const handleClose = () => {
  formRef.value?.clearValidate()
  activePermissionTab.value = "user"
  dialogVisible.value = false
}
</script>

<style lang="scss" scoped>
.form-help {
  margin-top: 4px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
  line-height: 1.4;
}

.info-text {
  color: var(--el-text-color-regular);
  font-size: 14px;
  margin-right: 12px;
}

.permissions-section {
  border: 1px solid var(--el-border-color-light);
  border-radius: 6px;
  padding: 16px;
  
  .permission-tabs {
    :deep(.el-tabs__content) {
      padding: 16px 0 0 0;
    }
  }
  
  .permission-category {
    .category-header {
      margin-bottom: 16px;
      padding-bottom: 8px;
      border-bottom: 1px solid var(--el-border-color-lighter);
    }
    
    .permission-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
      gap: 12px;
      
      .permission-item {
        padding: 8px;
        border: 1px solid var(--el-border-color-lighter);
        border-radius: 4px;
        transition: all 0.2s ease;
        
        &:hover {
          border-color: var(--el-color-primary-light-7);
          background-color: var(--el-fill-color-extra-light);
        }
        
        .permission-info {
          display: flex;
          flex-direction: column;
          margin-left: 8px;
          
          .permission-name {
            font-weight: 500;
            color: var(--el-text-color-primary);
          }
          
          .permission-desc {
            font-size: 12px;
            color: var(--el-text-color-secondary);
            margin-top: 2px;
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

:deep(.el-input.is-disabled .el-input__wrapper) {
  background-color: var(--el-fill-color-light);
}

@media (max-width: 768px) {
  .permissions-section {
    .permission-category {
      .permission-grid {
        grid-template-columns: 1fr;
      }
    }
  }
}
</style>