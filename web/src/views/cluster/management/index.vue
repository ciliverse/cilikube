<template>
  <div class="cluster-management">
    <!-- 页面头部 -->
    <div class="page-header">
      <h1 class="page-title">{{ $t('clusterManagement.title') }}</h1>
      <div class="header-actions">
        <el-button :icon="Refresh" @click="fetchAllData" :loading="isLoading" circle />
        <el-button type="primary" :icon="Plus" @click="handleOpenDialog">
          {{ $t('clusterManagement.addCluster') }}
        </el-button>
      </div>
    </div>

    <!-- 工具栏 -->
    <div class="toolbar">
      <div class="search-filters">
        <el-input v-model="searchQuery" :placeholder="$t('clusterManagement.searchPlaceholder')" :prefix-icon="Search"
          clearable class="search-input" />
        <el-select v-model="filterEnvironment" :placeholder="$t('clusterManagement.filterEnvironment')" clearable
          class="filter-select">
          <el-option :label="$t('clusterManagement.allEnvironments')" value="" />
          <el-option :label="$t('clusterManagement.environments.production')" value="production" />
          <el-option :label="$t('clusterManagement.environments.staging')" value="staging" />
          <el-option :label="$t('clusterManagement.environments.development')" value="development" />
        </el-select>
        <el-select v-model="filterStatus" :placeholder="$t('clusterManagement.filterStatus')" clearable
          class="filter-select">
          <el-option :label="$t('clusterManagement.allStatuses')" value="" />
          <el-option :label="$t('clusterManagement.statuses.available')" value="available" />
          <el-option :label="$t('clusterManagement.statuses.unavailable')" value="unavailable" />
        </el-select>
      </div>
      <div class="toolbar-right">
        <el-button-group class="view-toggle">
          <el-button :type="viewMode === 'card' ? 'primary' : ''" :icon="Grid" @click="viewMode = 'card'" />
          <el-button :type="viewMode === 'list' ? 'primary' : ''" :icon="List" @click="viewMode = 'list'" />
        </el-button-group>
      </div>
    </div>

    <!-- 集群列表 -->
    <div class="cluster-list" v-loading="isLoading">
      <!-- 卡片视图 -->
      <div v-if="filteredClusters.length > 0 && viewMode === 'card'" class="cluster-grid">
        <div v-for="cluster in filteredClusters" :key="cluster.id" class="cluster-card"
          :class="{ 'active-cluster': isCurrentCluster(cluster) }" @click="handleClusterClick(cluster)">
          <div class="card-header">
            <div class="cluster-info">
              <div class="cluster-name-row">
                <div class="status-dot" :class="getStatusClass(cluster.status)"></div>
                <div class="cluster-name">{{ cluster.name }}</div>
                <el-tag v-if="isCurrentCluster(cluster)" type="success" size="small" class="active-badge">
                  {{ $t('clusterManagement.actions.active') }}
                </el-tag>
              </div>
              <div class="cluster-server">{{ cluster.server }}</div>
            </div>
          </div>

          <div class="card-body">
            <div class="cluster-meta">
              <div class="meta-item" v-if="cluster.version">
                <span class="meta-label">{{ $t('clusterManagement.columns.version') }}</span>
                <span class="meta-value">{{ cluster.version }}</span>
              </div>
              <div class="meta-item" v-if="cluster.environment">
                <span class="meta-label">{{ $t('clusterManagement.columns.environment') }}</span>
                <span class="meta-value">{{ getEnvironmentLabel(cluster.environment) }}</span>
              </div>
              <div class="meta-item">
                <span class="meta-label">{{ $t('clusterManagement.columns.source') }}</span>
                <span class="meta-value">{{ cluster.source === 'database' ? $t('clusterManagement.sources.database') :
                  $t('clusterManagement.sources.file') }}</span>
              </div>
            </div>
          </div>

          <div class="card-footer">
            <div class="cluster-actions">
              <el-button type="primary" size="small"
                :disabled="isCurrentCluster(cluster) || !isAvailable(cluster.status)"
                @click.stop="handleSetActive(cluster)">
                {{ $t('clusterManagement.actions.setActive') }}
              </el-button>
              <el-popconfirm :title="$t('clusterManagement.messages.deleteConfirm')" @confirm="handleDelete(cluster)"
                :disabled="cluster.source === 'file'">
                <template #reference>
                  <el-button type="danger" size="small" :disabled="cluster.source === 'file'" @click.stop>
                    {{ $t('common.delete') }}
                  </el-button>
                </template>
              </el-popconfirm>
            </div>
          </div>
        </div>
      </div>

      <!-- 列表视图 -->
      <div v-else-if="filteredClusters.length > 0 && viewMode === 'list'" class="cluster-table">
        <el-table :data="filteredClusters" @row-click="handleClusterClick">
          <el-table-column width="60">
            <template #default="{ row }">
              <div class="status-dot" :class="getStatusClass(row.status)"></div>
            </template>
          </el-table-column>
          <el-table-column prop="name" :label="$t('clusterManagement.columns.name')" min-width="150">
            <template #default="{ row }">
              <div class="table-cluster-name">
                {{ row.name }}
                <el-tag v-if="isCurrentCluster(row)" type="success" size="small" class="active-badge">
                  {{ $t('clusterManagement.actions.active') }}
                </el-tag>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="server" :label="$t('clusterManagement.columns.server')" min-width="200"
            show-overflow-tooltip>
            <template #default="{ row }">
              <code class="server-url">{{ row.server }}</code>
            </template>
          </el-table-column>
          <el-table-column prop="version" :label="$t('clusterManagement.columns.version')" width="120" />
          <el-table-column prop="environment" :label="$t('clusterManagement.columns.environment')" width="120">
            <template #default="{ row }">
              <span v-if="row.environment">{{ getEnvironmentLabel(row.environment) }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="source" :label="$t('clusterManagement.columns.source')" width="100">
            <template #default="{ row }">
              <span>{{ row.source === 'database' ? $t('clusterManagement.sources.database') :
                $t('clusterManagement.sources.file') }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="$t('common.actions')" width="200" fixed="right">
            <template #default="{ row }">
              <div class="table-actions">
                <el-button type="primary" size="small" :disabled="isCurrentCluster(row) || !isAvailable(row.status)"
                  @click.stop="handleSetActive(row)">
                  {{ $t('clusterManagement.actions.setActive') }}
                </el-button>
                <el-popconfirm :title="$t('clusterManagement.messages.deleteConfirm')" @confirm="handleDelete(row)"
                  :disabled="row.source === 'file'">
                  <template #reference>
                    <el-button type="danger" size="small" :disabled="row.source === 'file'" @click.stop>
                      {{ $t('common.delete') }}
                    </el-button>
                  </template>
                </el-popconfirm>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <div v-else-if="!isLoading" class="empty-state">
        <el-empty :description="$t('clusterManagement.messages.noMatchingClusters')" />
      </div>
    </div>

    <!-- 添加集群对话框 -->
    <el-dialog v-model="dialogVisible" :title="$t('clusterManagement.addCluster')" width="700px"
      :close-on-click-modal="false" class="cluster-dialog">
      <el-form :model="form" ref="formRef" :rules="rules" label-width="120px" class="cluster-form">
        <div class="form-section">
          <h4 class="section-title">{{ $t('clusterManagement.basicInfo') }}</h4>
          <el-form-item :label="$t('clusterManagement.form.name')" prop="name">
            <el-input v-model="form.name" :placeholder="$t('clusterManagement.form.namePlaceholder')" clearable />
          </el-form-item>

          <el-row :gutter="16">
            <el-col :span="12">
              <el-form-item :label="$t('clusterManagement.form.environment')" prop="environment">
                <el-select v-model="form.environment" :placeholder="$t('clusterManagement.form.environmentPlaceholder')"
                  style="width: 100%">
                  <el-option :label="$t('clusterManagement.environments.production')" value="production" />
                  <el-option :label="$t('clusterManagement.environments.staging')" value="staging" />
                  <el-option :label="$t('clusterManagement.environments.development')" value="development" />
                  <el-option :label="$t('clusterManagement.environments.testing')" value="testing" />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item :label="$t('clusterManagement.form.provider')" prop="provider">
                <el-input v-model="form.provider" :placeholder="$t('clusterManagement.form.providerPlaceholder')"
                  clearable />
              </el-form-item>
            </el-col>
          </el-row>

          <el-form-item :label="$t('clusterManagement.form.region')" prop="region">
            <el-input v-model="form.region" :placeholder="$t('clusterManagement.form.regionPlaceholder')" clearable />
          </el-form-item>

          <el-form-item :label="$t('clusterManagement.form.description')" prop="description">
            <el-input v-model="form.description" type="textarea" :rows="3"
              :placeholder="$t('clusterManagement.form.descriptionPlaceholder')" />
          </el-form-item>
        </div>

        <div class="form-section">
          <h4 class="section-title">{{ $t('clusterManagement.form.kubeconfigTitle') }}</h4>
          <el-form-item :label="$t('clusterManagement.form.kubeconfig')" prop="kubeconfigData">
            <div class="kubeconfig-input">
              <el-input v-model="form.kubeconfigData" type="textarea" :rows="12"
                :placeholder="$t('clusterManagement.form.kubeconfigPlaceholder')" class="kubeconfig-textarea" />
              <div class="kubeconfig-tips">
                <el-alert :title="$t('clusterManagement.form.kubeconfigTips')" type="info" :closable="false"
                  show-icon />
              </div>
            </div>
          </el-form-item>
        </div>
      </el-form>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogVisible = false" size="large">
            {{ $t('common.cancel') }}
          </el-button>
          <el-button type="primary" @click="submitForm" :loading="isSubmitting" size="large">
            {{ $t('common.confirm') }}
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from "vue"
import { useI18n } from 'vue-i18n'
import { type FormInstance, type FormRules, ElMessage } from "element-plus"
import {
  Plus,
  Search,
  Refresh,
  Grid,
  List
} from "@element-plus/icons-vue"
// 移除未使用的 dayjs 导入
import {
  getClusterList,
  createCluster,
  deleteCluster,
  setActiveCluster,
  getActiveCluster,
  type ClusterInfo,
  type CreateClusterRequest
} from "@/api/cluster"
import { getCurrentClusterId, getSavedClusterId } from "@/utils/cluster-context"

const { t } = useI18n()

// 响应式状态
const clusterList = ref<ClusterInfo[]>([])
const activeClusterName = ref<string>("")
const isLoading = ref(false)
const isSubmitting = ref(false)
const dialogVisible = ref(false)
const formRef = ref<FormInstance | null>(null)

// 筛选状态
const searchQuery = ref('')
const filterEnvironment = ref('')
const filterStatus = ref('')
const viewMode = ref<'card' | 'list'>('card')

const form = ref<CreateClusterRequest>({
  name: "",
  kubeconfigData: "",
  provider: "",
  description: "",
  environment: "",
  region: ""
})

// 表单验证规则
const rules: FormRules = {
  name: [{ required: true, message: t('clusterManagement.validation.nameRequired'), trigger: "blur" }],
  kubeconfigData: [{ required: true, message: t('clusterManagement.validation.kubeconfigRequired'), trigger: "blur" }],
  environment: [{ required: true, message: t('clusterManagement.validation.environmentRequired'), trigger: "change" }]
}

// 移除未使用的统计数据

// 过滤后的集群列表
const filteredClusters = computed(() => {
  return clusterList.value.filter(cluster => {
    // 搜索过滤
    const matchesSearch = !searchQuery.value ||
      cluster.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      cluster.server.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      (cluster.environment && cluster.environment.toLowerCase().includes(searchQuery.value.toLowerCase()))

    // 环境过滤
    const matchesEnvironment = !filterEnvironment.value ||
      cluster.environment?.toLowerCase() === filterEnvironment.value.toLowerCase()

    // 状态过滤
    const matchesStatus = !filterStatus.value ||
      (filterStatus.value === 'available' && cluster.status.startsWith('可用')) ||
      (filterStatus.value === 'unavailable' && cluster.status.startsWith('不可用'))

    return matchesSearch && matchesEnvironment && matchesStatus
  })
})

// 获取集群列表和活动集群
const fetchAllData = async () => {
  isLoading.value = true
  try {
    const [listRes, activeRes] = await Promise.all([getClusterList(), getActiveCluster()])
    clusterList.value = listRes.data
    activeClusterName.value = activeRes.data
  } catch (error) {
    console.error("获取集群数据失败:", error)
    ElMessage.error("获取集群数据失败，请检查网络或联系管理员")
  } finally {
    isLoading.value = false
  }
}

// 移除未使用的 getStatusTagType 函数

// 辅助函数：根据环境获取标签类型
const getEnvironmentType = (environment: string | undefined) => {
  if (!environment) return 'info'
  const env = environment.toLowerCase()
  if (env === 'production' || env === 'prod') return 'danger'
  if (env === 'staging' || env === 'stage') return 'warning'
  if (env === 'development' || env === 'dev') return 'success'
  return 'info'
}

// 辅助函数：判断集群是否可用
const isAvailable = (status: string) => status.startsWith("可用")

// 获取状态指示器样式类
const getStatusClass = (status: string) => {
  // 只有明确标识为不可用的才显示红色，其他情况都显示绿色（正常）
  if (status.startsWith("不可用") || status.toLowerCase().includes("error") || status.toLowerCase().includes("failed")) {
    return "status-unhealthy"
  }
  // 发现的集群默认都是正常的（绿色）
  return "status-healthy"
}

// 判断是否为当前选中的集群（与导航栏联动）
const isCurrentCluster = (cluster: ClusterInfo) => {
  // 优先使用集群上下文中的ID
  const currentId = getCurrentClusterId() || getSavedClusterId()
  if (currentId) {
    return cluster.id === currentId
  }
  // 向后兼容：使用名称比较
  return cluster.name === activeClusterName.value
}

// 获取环境标签文本
const getEnvironmentLabel = (environment: string) => {
  const envMap: Record<string, string> = {
    production: t('clusterManagement.environments.production'),
    staging: t('clusterManagement.environments.staging'),
    development: t('clusterManagement.environments.development'),
    testing: t('clusterManagement.environments.testing')
  }
  return envMap[environment] || environment
}

// 处理集群点击
const handleClusterClick = (cluster: ClusterInfo) => {
  // 可以在这里添加集群详情查看逻辑
  console.log('Cluster clicked:', cluster)
}

// 处理操作：设为活动
const handleSetActive = async (row: ClusterInfo) => {
  try {
    await setActiveCluster(row.id)
    ElMessage.success(`集群 '${row.name}' 已设为活动集群`)
    await fetchAllData()
  } catch (error) {
    console.error("切换活动集群失败:", error)
    ElMessage.error("切换活动集群失败")
  }
}

// 处理操作：删除
const handleDelete = async (row: ClusterInfo) => {
  try {
    await deleteCluster(row.id)
    ElMessage.success(`集群 '${row.name}' 已删除`)
    await fetchAllData()
  } catch (error) {
    console.error("删除集群失败:", error)
    ElMessage.error("删除集群失败")
  }
}

// 处理操作：打开对话框
const handleOpenDialog = () => {
  if (formRef.value) {
    formRef.value.resetFields()
  }
  form.value = {
    name: "",
    kubeconfigData: "",
    provider: "",
    description: "",
    environment: "",
    region: ""
  }
  dialogVisible.value = true
}

// 验证 kubeconfig 格式
const validateKubeconfig = (kubeconfigData: string): boolean => {
  try {
    // 检查是否包含必要的 kubeconfig 字段
    const requiredFields = ['apiVersion', 'clusters', 'contexts', 'users']
    const hasAllFields = requiredFields.every(field => 
      kubeconfigData.includes(field + ':') || kubeconfigData.includes('"' + field + '"')
    )
    
    if (!hasAllFields) {
      ElMessage.error("kubeconfig 格式不正确，缺少必要字段")
      return false
    }
    
    // 检查是否为空或只有空白字符
    if (!kubeconfigData.trim()) {
      ElMessage.error("kubeconfig 不能为空")
      return false
    }
    
    return true
  } catch (error) {
    ElMessage.error("kubeconfig 格式验证失败")
    return false
  }
}

// 处理操作：提交表单
const submitForm = () => {
  formRef.value?.validate(async (valid) => {
    if (valid) {
      // 验证 kubeconfig 格式
      if (!validateKubeconfig(form.value.kubeconfigData)) {
        return
      }
      
      isSubmitting.value = true
      try {
        console.log('准备发送的数据:', {
          name: form.value.name,
          environment: form.value.environment,
          kubeconfigLength: form.value.kubeconfigData.length
        })
        
        await createCluster(form.value)
        ElMessage.success("集群添加成功")
        dialogVisible.value = false
        await fetchAllData()
      } catch (error: any) {
        console.error("添加集群失败:", error)
        console.error("错误详情:", error.response?.data)
        
        // 显示更详细的错误信息
        let errorMsg = "添加集群失败"
        if (error.response?.data?.message) {
          errorMsg = error.response.data.message
        } else if (error.message) {
          errorMsg = error.message
        }
        
        // 如果是 kubeconfig 相关错误，提供更友好的提示
        if (errorMsg.includes('kubeconfig') || errorMsg.includes('parse')) {
          errorMsg += "\n\n请检查：\n1. kubeconfig 格式是否正确\n2. 是否包含完整的集群配置信息\n3. 证书和密钥是否有效"
        }
        
        ElMessage.error(errorMsg)
      } finally {
        isSubmitting.value = false
      }
    }
  })
}

// 组件挂载后加载数据
onMounted(() => {
  fetchAllData()
})
</script>

<style scoped>
.cluster-management {
  padding: 24px;
  background: #f8fafc;
  min-height: 100vh;
}

/* 页面头部 */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-title {
  font-size: 32px;
  font-weight: 600;
  color: #1e293b;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

/* 工具栏 */
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding: 16px 20px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border: 1px solid #e2e8f0;
}

.search-filters {
  display: flex;
  gap: 12px;
  align-items: center;
  flex: 1;
}

.search-input {
  width: 280px;
}

.filter-select {
  width: 160px;
}

.toolbar-right {
  display: flex;
  gap: 16px;
  align-items: center;
}

.cluster-count {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 14px;
  color: #64748b;
  font-weight: 500;
  white-space: nowrap;
}

.count-number {
  font-weight: 600;
  color: #1e293b;
}

.count-text {
  font-weight: 500;
}

.view-toggle {
  border-radius: 8px;
}

/* 集群列表 */
.cluster-list {
  min-height: 400px;
}

.cluster-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(380px, 1fr));
  gap: 20px;
}

.cluster-card {
  background: white;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  transition: all 0.2s ease;
  cursor: pointer;
  overflow: hidden;
}

.cluster-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border-color: #3b82f6;
}

.cluster-card.active-cluster {
  border-color: #10b981;
  border-width: 2px;
  box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.2);
}

.card-header {
  padding: 20px 20px 16px;
  border-bottom: 1px solid #f1f5f9;
}

.cluster-info {
  margin-bottom: 12px;
}

.cluster-name {
  font-size: 18px;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 4px;
}

.cluster-server {
  font-size: 13px;
  color: #475569;
  font-family: 'Maple Mono', 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-weight: 600;
  margin-top: 4px;
}

.status-badges {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.card-body {
  padding: 16px 20px;
}

.cluster-meta {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.meta-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.meta-label {
  font-size: 13px;
  color: #64748b;
  font-weight: 500;
}

.meta-value {
  font-size: 13px;
  color: #1e293b;
  font-family: 'Maple Mono', 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-weight: 600;
}

.card-footer {
  padding: 16px 20px;
  border-top: 1px solid #f1f5f9;
  background: white;
}

.cluster-actions {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
}

/* 状态指示器 - 简洁小圆点样式 */
.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
  transition: all 0.2s ease;
}

.status-healthy {
  background: #10b981;
  box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.2);
}

.status-unhealthy {
  background: #ef4444;
  box-shadow: 0 0 0 2px rgba(239, 68, 68, 0.2);
}

.status-unknown {
  background: #64748b;
  box-shadow: 0 0 0 2px rgba(100, 116, 139, 0.2);
}

/* 集群名称行 */
.cluster-name-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.active-badge {
  margin-left: auto;
  background: linear-gradient(135deg, #10b981 0%, #059669 100%) !important;
  border: none !important;
  color: white !important;
  font-weight: 600;
  box-shadow: 0 2px 4px rgba(16, 185, 129, 0.3);
  animation: glow-green 2s ease-in-out infinite alternate;
}

/* 活动集群特殊样式 - 简化颜色，只保留绿色边框 */
.cluster-card.active-cluster .cluster-name {
  color: #1e293b;
  font-weight: 700;
}

/* 活动标签动画 */
@keyframes glow-green {
  0% {
    box-shadow: 0 2px 4px rgba(16, 185, 129, 0.3);
  }

  100% {
    box-shadow: 0 2px 8px rgba(16, 185, 129, 0.6), 0 0 12px rgba(16, 185, 129, 0.3);
  }
}

/* 列表视图 */
.cluster-table {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid #e2e8f0;
}

.table-cluster-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.server-url {
  font-family: 'Maple Mono', 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 12px;
  color: #475569;
  font-weight: 600;
}

.table-actions {
  display: flex;
  gap: 8px;
}

.empty-state {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 80px 20px;
  background: white;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
}

/* 对话框样式 */
.cluster-dialog :deep(.el-dialog__body) {
  padding: 20px 24px;
}

.cluster-form {
  max-height: 70vh;
  overflow-y: auto;
}

.form-section {
  margin-bottom: 32px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #1e293b;
  margin: 0 0 16px 0;
  padding-bottom: 8px;
  border-bottom: 1px solid #e2e8f0;
}

.kubeconfig-input {
  width: 100%;
}

.kubeconfig-textarea :deep(.el-textarea__inner) {
  font-family: 'Maple Mono', 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 12px;
  line-height: 1.5;
}

.kubeconfig-tips {
  margin-top: 12px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 0 0;
  border-top: 1px solid #e2e8f0;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .cluster-grid {
    grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  }
}

@media (max-width: 768px) {
  .cluster-management {
    padding: 16px;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .page-title {
    font-size: 28px;
  }

  .toolbar {
    flex-direction: column;
    align-items: stretch;
    gap: 16px;
  }

  .search-filters {
    flex-direction: column;
    align-items: stretch;
  }

  .search-input,
  .filter-select {
    width: 100%;
  }

  .toolbar-right {
    justify-content: space-between;
    width: 100%;
  }

  .cluster-count {
    text-align: left;
  }

  .count-text {
    display: none;
  }

  .cluster-grid {
    grid-template-columns: 1fr;
  }

  .meta-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }

  .cluster-actions {
    flex-direction: column;
  }
}
</style>
