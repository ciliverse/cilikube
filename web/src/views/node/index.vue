<template>
  <div class="node-management">
    <!-- 页面头部 -->
    <div class="page-header">
      <h1 class="page-title">{{ $t('nodeManagement.title') }}</h1>
      <div class="header-actions">
        <el-button :icon="Refresh" @click="fetchNodeData" :loading="loading" circle />
      </div>
    </div>

    <!-- 集群选择提示 -->
    <el-alert
      v-if="!clusterStore.selectedClusterName && !clusterStore.loadingClusters"
      :title="$t('nodeManagement.messages.selectCluster')"
      type="info" show-icon :closable="false" style="margin-bottom: 24px;"
    />

    <!-- 工具栏 -->
    <div class="toolbar" v-if="clusterStore.selectedClusterName">
      <div class="search-filters">
        <el-input v-model="searchQuery" :placeholder="$t('nodeManagement.searchPlaceholder')" :prefix-icon="Search"
          clearable class="search-input" />
        <el-select v-model="filterRole" :placeholder="$t('nodeManagement.filterRole')" clearable
          class="filter-select">
          <el-option :label="$t('nodeManagement.allRoles')" value="" />
          <el-option :label="$t('nodeManagement.roles.controlPlane')" value="control-plane" />
          <el-option :label="$t('nodeManagement.roles.worker')" value="worker" />
        </el-select>
        <el-select v-model="filterStatus" :placeholder="$t('nodeManagement.filterStatus')" clearable
          class="filter-select">
          <el-option :label="$t('nodeManagement.allStatuses')" value="" />
          <el-option :label="$t('nodeManagement.statuses.ready')" value="ready" />
          <el-option :label="$t('nodeManagement.statuses.notReady')" value="notready" />
        </el-select>
      </div>
      <div class="toolbar-right">
        <el-button-group class="view-toggle">
          <el-button :type="viewMode === 'card' ? 'primary' : ''" :icon="Grid" @click="viewMode = 'card'" />
          <el-button :type="viewMode === 'list' ? 'primary' : ''" :icon="List" @click="viewMode = 'list'" />
        </el-button-group>
      </div>
    </div>

    <!-- 节点列表 -->
    <div class="node-list" v-loading="loading">
      <!-- 卡片视图 -->
      <div v-if="filteredNodes.length > 0 && viewMode === 'card'" class="node-grid">
        <div v-for="node in currentPageData" :key="node.metadata.uid" class="node-card"
          @click="handleNodeClick(node)">
          <div class="card-header">
            <div class="node-info">
              <div class="node-name-row">
                <div class="status-dot" :class="getStatusClass(node.status?.conditions)"></div>
                <div class="node-name">{{ node.metadata.name }}</div>
                <el-tag v-if="isControlPlaneNode(node)" type="primary" size="small" class="role-badge">
                  {{ $t('nodeManagement.roles.controlPlane') }}
                </el-tag>
              </div>
              <div class="node-address">{{ getNodeAddress(node.status?.addresses) }}</div>
            </div>
          </div>

          <div class="card-body">
            <div class="node-meta">
              <div class="meta-item">
                <span class="meta-label">{{ $t('nodeManagement.info.kubeletVersion') }}</span>
                <span class="meta-value">{{ node.status?.nodeInfo?.kubeletVersion || 'N/A' }}</span>
              </div>
              <div class="meta-item">
                <span class="meta-label">{{ $t('nodeManagement.info.osImage') }}</span>
                <span class="meta-value">{{ node.status?.nodeInfo?.osImage || 'N/A' }}</span>
              </div>
              <div class="meta-item">
                <span class="meta-label">{{ $t('nodeManagement.info.containerRuntime') }}</span>
                <span class="meta-value">{{ formatRuntime(node.status?.nodeInfo?.containerRuntimeVersion) }}</span>
              </div>
              <div class="meta-item">
                <span class="meta-label">{{ $t('nodeManagement.columns.cpu') }}</span>
                <span class="meta-value">{{ node.status?.capacity?.cpu || 'N/A' }} {{ $t('nodeMetrics.cores') }}</span>
              </div>
              <div class="meta-item">
                <span class="meta-label">{{ $t('nodeManagement.columns.memory') }}</span>
                <span class="meta-value">{{ formatMemory(node.status?.capacity?.memory) }}</span>
              </div>
            </div>
          </div>

          <div class="card-footer">
            <div class="node-actions">
              <el-button type="primary" size="small" @click.stop="handleViewDetails(node)">
                {{ $t('nodeManagement.actions.viewDetails') }}
              </el-button>
            </div>
          </div>
        </div>
      </div>

      <!-- 列表视图 -->
      <div v-else-if="filteredNodes.length > 0 && viewMode === 'list'" class="node-table">
        <el-table :data="currentPageData" @row-click="handleNodeClick">
          <el-table-column width="60">
            <template #default="{ row }">
              <div class="status-dot" :class="getStatusClass(row.status?.conditions)"></div>
            </template>
          </el-table-column>
          <el-table-column prop="metadata.name" :label="$t('nodeManagement.columns.name')" min-width="150">
            <template #default="{ row }">
              <div class="table-node-name">
                {{ row.metadata.name }}
                <el-tag v-if="isControlPlaneNode(row)" type="primary" size="small" class="role-badge">
                  {{ $t('nodeManagement.roles.controlPlane') }}
                </el-tag>
              </div>
            </template>
          </el-table-column>
          <el-table-column :label="$t('nodeManagement.columns.status')" width="100">
            <template #default="{ row }">
              <span>{{ getNodeStatusText(row.status?.conditions) }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="$t('nodeManagement.columns.internalIP')" width="140">
            <template #default="{ row }">
              <code class="ip-address">{{ getNodeAddress(row.status?.addresses) }}</code>
            </template>
          </el-table-column>
          <el-table-column :label="$t('nodeManagement.columns.version')" width="120">
            <template #default="{ row }">
              <code class="version-text">{{ row.status?.nodeInfo?.kubeletVersion || 'N/A' }}</code>
            </template>
          </el-table-column>
          <el-table-column :label="$t('nodeManagement.columns.os')" min-width="200" show-overflow-tooltip>
            <template #default="{ row }">
              <span>{{ row.status?.nodeInfo?.osImage || 'N/A' }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="$t('nodeManagement.columns.cpu')" width="100">
            <template #default="{ row }">
              <span>{{ row.status?.capacity?.cpu || 'N/A' }} {{ $t('nodeMetrics.cores') }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="$t('nodeManagement.columns.memory')" width="120">
            <template #default="{ row }">
              <span>{{ formatMemory(row.status?.capacity?.memory) }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="$t('common.actions')" width="120" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" size="small" @click.stop="handleViewDetails(row)">
                {{ $t('nodeManagement.actions.viewDetails') }}
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <div v-else-if="!loading && clusterStore.selectedClusterName" class="empty-state">
        <el-empty :description="searchQuery || filterRole || filterStatus ? 
          $t('nodeManagement.messages.noMatchingNodes') : 
          $t('nodeManagement.messages.noNodes')" />
      </div>
    </div>

    <!-- 分页 -->
    <div class="pagination-container" v-if="filteredNodes.length > pageSize">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[12, 24, 36, 48]"
        :total="filteredNodes.length"
        layout="total, sizes, prev, pager, next, jumper"
        background
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue"
import { useI18n } from 'vue-i18n'
import { ElMessage } from "element-plus"
import { useClusterStore } from "@/store/modules/clusterStore"
import { getNodeList } from "@/api/node" 
import type { Node, NodeCondition, NodeAddress } from "@/api/types/node"
import {
  Refresh, Search, Grid, List
} from '@element-plus/icons-vue'

const { t } = useI18n()

// Store
const clusterStore = useClusterStore()

// 响应式状态
const nodeData = ref<Node[]>([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(12)

// 筛选状态
const searchQuery = ref('')
const filterRole = ref('')
const filterStatus = ref('')
const viewMode = ref<'card' | 'list'>('card')

// 过滤后的节点列表
const filteredNodes = computed(() => {
  return nodeData.value.filter(node => {
    // 搜索过滤
    const matchesSearch = !searchQuery.value ||
      node.metadata.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      getNodeAddress(node.status?.addresses).toLowerCase().includes(searchQuery.value.toLowerCase())

    // 角色过滤
    const matchesRole = !filterRole.value ||
      (filterRole.value === 'control-plane' && isControlPlaneNode(node)) ||
      (filterRole.value === 'worker' && !isControlPlaneNode(node))

    // 状态过滤
    const nodeStatus = getNodeStatus(node.status?.conditions).toLowerCase()
    const matchesStatus = !filterStatus.value ||
      (filterStatus.value === 'ready' && nodeStatus === 'ready') ||
      (filterStatus.value === 'notready' && nodeStatus === 'notready')

    return matchesSearch && matchesRole && matchesStatus
  })
})

// 当前页数据
const currentPageData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredNodes.value.slice(start, end)
})

// 获取节点数据
const fetchNodeData = async () => {
  const clusterName = clusterStore.selectedClusterName
  if (!clusterName) {
    nodeData.value = []
    return
  }

  loading.value = true
  try {
    const res = await getNodeList()
    if (res && res.data && Array.isArray(res.data.items)) {
      nodeData.value = res.data.items
    } else {
      throw new Error('节点数据格式错误')
    }
  } catch (error: any) {
    console.error(`获取节点数据时发生错误:`, error)
    ElMessage.error(`获取节点数据失败: ${error.message || '未知错误'}`)
    nodeData.value = []
  } finally {
    loading.value = false
  }
}

// 监听集群变化
watch(() => clusterStore.selectedClusterName, (newClusterName) => {
  if (newClusterName) {
    ElMessage.info(`${t('nodeManagement.messages.loadingNodes')} '${clusterStore.currentClusterDisplayName}'...`)
    fetchNodeData()
  } else {
    nodeData.value = []
  }
}, { immediate: true })

// 组件挂载
onMounted(() => {
  if (clusterStore.availableClusters.length === 0) {
    clusterStore.fetchAvailableClusters()
  }
})

// 辅助函数
const getNodeStatus = (conditions: NodeCondition[] | undefined): 'Ready' | 'NotReady' | 'Unknown' => {
  if (!conditions) return 'Unknown'
  const ready = conditions.find(c => c.type === 'Ready')
  return ready?.status === 'True' ? 'Ready' : 'NotReady'
}

const getNodeStatusText = (conditions: NodeCondition[] | undefined) => {
  const status = getNodeStatus(conditions)
  if (status === 'Ready') return t('nodeManagement.statuses.ready')
  if (status === 'NotReady') return t('nodeManagement.statuses.notReady')
  return t('nodeManagement.statuses.unknown')
}

const getStatusClass = (conditions: NodeCondition[] | undefined) => {
  const status = getNodeStatus(conditions)
  if (status === 'Ready') return 'status-healthy'
  if (status === 'NotReady') return 'status-unhealthy'
  return 'status-unknown'
}

const isControlPlaneNode = (node: Node) => {
  return node.metadata.labels ? 'node-role.kubernetes.io/control-plane' in node.metadata.labels : false
}

const getNodeAddress = (addresses: NodeAddress[] | undefined) => {
  if (!addresses) return "N/A"
  const internal = addresses.find(a => a.type === 'InternalIP')
  return internal ? internal.address : 'N/A'
}

const formatMemory = (memoryStr: string | undefined) => {
  if (!memoryStr) return 'N/A'
  const memoryKi = parseInt(memoryStr.replace('Ki', ''), 10)
  const memoryGi = memoryKi / 1024 / 1024
  return `${memoryGi.toFixed(1)} GB`
}

const formatRuntime = (runtime: string | undefined) => {
  if (!runtime) return 'N/A'
  // 简化容器运行时显示，例如 "containerd://1.6.6" -> "containerd 1.6.6"
  return runtime.replace('://', ' ')
}

// 事件处理
const handlePageChange = (page: number) => { 
  currentPage.value = page 
}

const handleSizeChange = (size: number) => { 
  pageSize.value = size
  currentPage.value = 1 
}

const handleNodeClick = (node: Node) => {
  console.log('Node clicked:', node)
}

const handleViewDetails = (node: Node) => {
  ElMessage.info(`查看节点 ${node.metadata.name} 的详细信息`)
}
</script>

<style scoped>
.node-management {
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

.view-toggle {
  border-radius: 8px;
}

/* 节点列表 */
.node-list {
  min-height: 400px;
}

.node-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(380px, 1fr));
  gap: 20px;
}

.node-card {
  background: white;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  transition: all 0.2s ease;
  cursor: pointer;
  overflow: hidden;
}

.node-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border-color: #3b82f6;
}

.card-header {
  padding: 20px 20px 16px;
  border-bottom: 1px solid #f1f5f9;
}

.node-info {
  margin-bottom: 12px;
}

.node-name {
  font-size: 18px;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 4px;
}

.node-address {
  font-size: 13px;
  color: #475569;
  font-family: 'Maple Mono', 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-weight: 600;
  margin-top: 4px;
}

.card-body {
  padding: 16px 20px;
}

.node-meta {
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

.node-actions {
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

/* 节点名称行 */
.node-name-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.role-badge {
  margin-left: auto;
  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%) !important;
  border: none !important;
  color: white !important;
  font-weight: 600;
  box-shadow: 0 2px 4px rgba(59, 130, 246, 0.3);
}

/* 列表视图 */
.node-table {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid #e2e8f0;
}

.table-node-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.ip-address {
  font-family: 'Maple Mono', 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 12px;
  color: #475569;
  font-weight: 600;
}

.version-text {
  font-family: 'Maple Mono', 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 12px;
  color: #475569;
  font-weight: 600;
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

/* 分页 */
.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 24px;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .node-grid {
    grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  }
}

@media (max-width: 768px) {
  .node-management {
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

  .node-grid {
    grid-template-columns: 1fr;
  }

  .meta-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }

  .node-actions {
    flex-direction: column;
  }
}
</style>
