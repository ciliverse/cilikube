<template>
  <el-dialog
    v-model="dialogVisible"
    title="创建角色"
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
          maxlength="50"
          show-word-limit
        />
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

      <el-form-item label="预设模板">
        <div class="role-templates">
          <el-button
            v-for="template in roleTemplates"
            :key="template.key"
            size="small"
            @click="applyTemplate(template)"
          >
            {{ template.name }}
          </el-button>
        </div>
        <div class="form-help">
          点击预设模板快速配置权限
        </div>
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
          创建角色
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import { ref, reactive, computed, watch } from "vue"
import { ElMessage, type FormInstance, type FormRules } from "element-plus"
import { createRoleApi } from "@/api/admin/role-management"

interface Props {
  modelValue: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
  "update:modelValue": [value: boolean]
  created: []
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

// 权限分类
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

// 角色模板
const roleTemplates = [
  {
    key: "admin",
    name: "管理员",
    permissions: permissionCategories.flatMap(cat => cat.permissions.map(p => p.key))
  },
  {
    key: "developer",
    name: "开发者",
    permissions: [
      "clusters:read", "pods:read", "pods:write", "pods:exec", "pods:logs",
      "deployments:read", "deployments:write", "services:read", "services:write",
      "configmaps:read", "configmaps:write", "secrets:read"
    ]
  },
  {
    key: "viewer",
    name: "查看者",
    permissions: [
      "clusters:read", "pods:read", "pods:logs", "deployments:read",
      "services:read", "ingress:read", "configmaps:read", "secrets:read"
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

/** 应用角色模板 */
const applyTemplate = (template: any) => {
  formData.permissions = [...template.permissions]
  ElMessage.success(`已应用 ${template.name} 模板`)
}

/** 监听对话框显示状态 */
watch(() => props.modelValue, (val) => {
  dialogVisible.value = val
})

watch(dialogVisible, (val) => {
  emit("update:modelValue", val)
})

/** 提交表单 */
const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    
    if (formData.permissions.length === 0) {
      ElMessage.warning("请至少选择一个权限")
      return
    }
    
    loading.value = true
    
    await createRoleApi({
      name: formData.name,
      display_name: formData.display_name,
      description: formData.description,
      permissions: formData.permissions
    })
    
    ElMessage.success("角色创建成功")
    emit("created")
    handleClose()
    
  } catch (error: any) {
    console.error("创建角色失败:", error)
    ElMessage.error(error.response?.data?.message || "创建角色失败")
  } finally {
    loading.value = false
  }
}

/** 关闭对话框 */
const handleClose = () => {
  // 重置表单
  formRef.value?.resetFields()
  Object.assign(formData, {
    name: "",
    display_name: "",
    description: "",
    permissions: []
  })
  
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

.role-templates {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 8px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
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