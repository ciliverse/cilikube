<template>
  <div class="resource-page">
    <div class="page-header">
      <h1 class="page-title">{{ $t('roleResourceManagement.title') }}</h1>
    </div>

    <div class="filter-bar">
      <el-select
        v-model="selectedNamespace"
        :placeholder="$t('roleResourceManagement.filterNamespace')"
        filterable
        :loading="loadingNamespaces"
        class="namespace-select"
        @change="handleNamespaceChange"
      >
        <el-option v-for="ns in namespaces" :key="ns" :label="ns" :value="ns" />
      </el-select>
      <el-input
        v-model="searchQuery"
        :placeholder="$t('roleResourceManagement.searchPlaceholder')"
        :prefix-icon="SearchIcon"
        clearable
        class="search-input"
        :disabled="!selectedNamespace"
      />
      <el-tooltip :content="$t('actions.refresh')" placement="top">
        <el-button
          :icon="RefreshIcon"
          circle
          @click="fetchData"
          :loading="loading"
          :disabled="!selectedNamespace"
        />
      </el-tooltip>
    </div>

    <el-card shadow="hover" class="table-card">
      <el-table :data="paginatedData" v-loading="loading" border stripe row-key="uid">
        <el-table-column prop="name" :label="$t('roleResourceManagement.columns.name')" min-width="200" show-overflow-tooltip fixed>
          <template #default="{ row }">
            <el-icon class="resource-icon"><Key /></el-icon>
            <span>{{ row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="namespace" :label="$t('roleResourceManagement.columns.namespace')" min-width="140" />
        <el-table-column prop="rulesCount" :label="$t('roleResourceManagement.columns.rulesCount')" width="120" align="center" />
        <el-table-column prop="rulesSummary" :label="$t('roleResourceManagement.columns.rulesSummary')" min-width="300" show-overflow-tooltip />
        <el-table-column prop="createdAt" :label="$t('roleResourceManagement.columns.createdAt')" min-width="180" />
        <el-table-column :label="$t('roleResourceManagement.columns.actions')" width="100" align="center" fixed="right">
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
          <el-empty
            :description="
              selectedNamespace
                ? $t('roleResourceManagement.messages.noData')
                : $t('roleResourceManagement.messages.selectNamespace')
            "
          />
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
import { kubernetesRequest, fetchNamespaces } from '@/utils/api-config'
import {
  Search as SearchIcon,
  Refresh as RefreshIcon,
  Delete as DeleteIcon,
  View as ViewIcon,
  Key
} from '@element-plus/icons-vue'

interface DisplayItem {
  uid: string
  name: string
  namespace: string
  rulesCount: number
  rulesSummary: string
  createdAt: string
  rawData: any
}

const { t } = useI18n()
const loading = ref(false)
const loadingNamespaces = ref(false)
const namespaces = ref<string[]>([])
const selectedNamespace = ref('')
const allItems = ref<DisplayItem[]>([])
const searchQuery = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const detailVisible = ref(false)
const detailItem = ref<DisplayItem | null>(null)

function summarizeRules(rules: any[] = []): string {
  return rules
    .slice(0, 3)
    .map((rule) => `${(rule.verbs || []).join(',')}→${(rule.resources || ['*']).join(',')}`)
    .join('; ')
}

const filteredData = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return allItems.value
  return allItems.value.filter((item) => item.name.toLowerCase().includes(q))
})

const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  return filteredData.value.slice(start, start + pageSize.value)
})

const loadNamespaces = async () => {
  loadingNamespaces.value = true
  try {
    namespaces.value = await fetchNamespaces()
    if (!selectedNamespace.value && namespaces.value.length > 0) {
      selectedNamespace.value = namespaces.value.includes('default')
        ? 'default'
        : namespaces.value[0]
      await fetchData()
    }
  } catch (error: any) {
    ElMessage.error(error.message || t('roleResourceManagement.messages.fetchNamespacesFailed'))
  } finally {
    loadingNamespaces.value = false
  }
}

const fetchData = async () => {
  if (!selectedNamespace.value) return
  loading.value = true
  try {
    const response = await kubernetesRequest<{ code: number; data: any; message: string }>({
      url: `/api/v1/namespaces/${selectedNamespace.value}/roles`,
      method: 'GET'
    })
    if (response.code === 200 && response.data?.items) {
      allItems.value = response.data.items.map((item: any) => ({
        uid: item.metadata?.uid || item.metadata?.name,
        name: item.metadata?.name || '',
        namespace: item.metadata?.namespace || selectedNamespace.value,
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
    ElMessage.error(error.message || t('roleResourceManagement.messages.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleNamespaceChange = () => {
  currentPage.value = 1
  fetchData()
}

const openDetail = (row: DisplayItem) => {
  detailItem.value = row
  detailVisible.value = true
}

const handleDelete = (row: DisplayItem) => {
  ElMessageBox.confirm(
    t('roleResourceManagement.messages.deleteConfirm', { name: row.name, namespace: row.namespace }),
    t('common.confirm'),
    { type: 'warning' }
  )
    .then(async () => {
      try {
        const response = await kubernetesRequest<{ code: number; message: string }>({
          url: `/api/v1/namespaces/${row.namespace}/roles/${row.name}`,
          method: 'DELETE'
        })
        if (response.code === 200 || response.code === 204) {
          ElMessage.success(t('roleResourceManagement.messages.deleteSuccess'))
          await fetchData()
        }
      } catch (error: any) {
        ElMessage.error(error.message || t('roleResourceManagement.messages.deleteFailed'))
      }
    })
    .catch(() => undefined)
}

onMounted(loadNamespaces)
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
.namespace-select {
  width: 220px;
}
.search-input {
  width: 280px;
}
.table-card :deep(.el-card__body) {
  padding: 0;
}
.resource-icon {
  margin-right: 8px;
  color: #8b5cf6;
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
