<template>
  <div class="resource-page">
    <div class="page-header">
      <h1 class="page-title">{{ $t('clusterRoleManagement.title') }}</h1>
    </div>

    <div class="filter-bar">
      <el-input
        v-model="searchQuery"
        :placeholder="$t('clusterRoleManagement.searchPlaceholder')"
        :prefix-icon="SearchIcon"
        clearable
        class="search-input"
      />
      <el-tooltip :content="$t('actions.refresh')" placement="top">
        <el-button :icon="RefreshIcon" circle @click="fetchData" :loading="loading" />
      </el-tooltip>
    </div>

    <el-card shadow="hover" class="table-card">
      <el-table :data="paginatedData" v-loading="loading" border stripe row-key="uid">
        <el-table-column prop="name" :label="$t('clusterRoleManagement.columns.name')" min-width="220" show-overflow-tooltip fixed>
          <template #default="{ row }">
            <el-icon class="resource-icon"><Lock /></el-icon>
            <span>{{ row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="rulesCount" :label="$t('clusterRoleManagement.columns.rulesCount')" width="120" align="center" />
        <el-table-column prop="rulesSummary" :label="$t('clusterRoleManagement.columns.rulesSummary')" min-width="320" show-overflow-tooltip />
        <el-table-column prop="createdAt" :label="$t('clusterRoleManagement.columns.createdAt')" min-width="180" />
        <el-table-column :label="$t('clusterRoleManagement.columns.actions')" width="100" align="center" fixed="right">
          <template #default="{ row }">
            <el-tooltip :content="$t('common.view')" placement="top">
              <el-button link type="primary" :icon="ViewIcon" @click="openDetail(row)" />
            </el-tooltip>
            <el-tooltip :content="$t('common.delete')" placement="top">
              <el-button link type="danger" :icon="DeleteIcon" @click="handleDelete(row)" />
            </el-tooltip>
          </template>
        </el-table-column>
        <template #empty>
          <el-empty :description="$t('clusterRoleManagement.messages.noData')" />
        </template>
      </el-table>
      <div class="pagination-container" v-if="filteredData.length > 0">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50]"
          :total="filteredData.length"
          layout="total, sizes, prev, pager, next"
          background
        />
      </div>
    </el-card>

    <el-drawer v-model="detailVisible" :title="detailItem?.name" size="50%">
      <pre v-if="detailItem" class="json-view">{{ JSON.stringify(detailItem.rawData, null, 2) }}</pre>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'
import { kubernetesRequest } from '@/utils/api-config'
import {
  Search as SearchIcon,
  Refresh as RefreshIcon,
  Delete as DeleteIcon,
  View as ViewIcon,
  Lock
} from '@element-plus/icons-vue'

interface DisplayItem {
  uid: string
  name: string
  rulesCount: number
  rulesSummary: string
  createdAt: string
  rawData: any
}

const { t } = useI18n()
const loading = ref(false)
const allItems = ref<DisplayItem[]>([])
const searchQuery = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const detailVisible = ref(false)
const detailItem = ref<DisplayItem | null>(null)

function summarizeRules(rules: any[] = []): string {
  return rules
    .slice(0, 3)
    .map((rule) => {
      const verbs = (rule.verbs || []).join(',')
      const resources = (rule.resources || ['*']).join(',')
      return `${verbs}→${resources}`
    })
    .join('; ')
}

const filteredData = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return allItems.value
  return allItems.value.filter(
    (item) => item.name.toLowerCase().includes(q) || item.rulesSummary.toLowerCase().includes(q)
  )
})

const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  return filteredData.value.slice(start, start + pageSize.value)
})

const fetchData = async () => {
  loading.value = true
  try {
    const response = await kubernetesRequest<{ code: number; data: any; message: string }>({
      url: '/api/v1/clusterroles',
      method: 'GET'
    })
    if (response.code === 200 && response.data?.items) {
      allItems.value = response.data.items.map((item: any) => ({
        uid: item.metadata?.uid || item.metadata?.name,
        name: item.metadata?.name || '',
        rulesCount: item.rules?.length || 0,
        rulesSummary: summarizeRules(item.rules) || '-',
        createdAt: item.metadata?.creationTimestamp
          ? dayjs(item.metadata.creationTimestamp).format('YYYY-MM-DD HH:mm:ss')
          : '-',
        rawData: item
      }))
    } else {
      allItems.value = []
    }
  } catch (error: any) {
    allItems.value = []
    ElMessage.error(error.message || t('clusterRoleManagement.messages.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const openDetail = (row: DisplayItem) => {
  detailItem.value = row
  detailVisible.value = true
}

const handleDelete = (row: DisplayItem) => {
  ElMessageBox.confirm(
    t('clusterRoleManagement.messages.deleteConfirm', { name: row.name }),
    t('common.confirm'),
    { type: 'warning' }
  )
    .then(async () => {
      try {
        const response = await kubernetesRequest<{ code: number; message: string }>({
          url: `/api/v1/clusterroles/${row.name}`,
          method: 'DELETE'
        })
        if (response.code === 200 || response.code === 204) {
          ElMessage.success(t('clusterRoleManagement.messages.deleteSuccess'))
          await fetchData()
        }
      } catch (error: any) {
        ElMessage.error(error.message || t('clusterRoleManagement.messages.deleteFailed'))
      }
    })
    .catch(() => undefined)
}

onMounted(fetchData)
</script>

<style lang="scss" scoped>
.resource-page {
  padding: 24px;
  background-color: var(--el-bg-color-page);
}
.page-header {
  margin-bottom: 24px;
}
.page-title {
  margin: 0;
  font-size: 28px;
  font-weight: 600;
}
.filter-bar {
  display: flex;
  gap: 16px;
  margin-bottom: 24px;
  padding: 16px;
  background: var(--el-bg-color);
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 4px;
}
.search-input {
  width: 320px;
}
.table-card :deep(.el-card__body) {
  padding: 0;
}
.resource-icon {
  margin-right: 8px;
  color: #6366f1;
  vertical-align: middle;
}
.pagination-container {
  display: flex;
  justify-content: flex-end;
  padding: 16px;
}
.json-view {
  margin: 0;
  padding: 16px;
  background: var(--el-fill-color-lighter);
  border-radius: 4px;
  font-size: 13px;
  overflow: auto;
  white-space: pre-wrap;
  word-break: break-word;
}
</style>
