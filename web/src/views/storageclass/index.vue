<template>
  <div class="resource-page">
    <div class="page-header">
      <h1 class="page-title">{{ $t('storageClassManagement.title') }}</h1>
    </div>

    <div class="filter-bar">
      <el-input
        v-model="searchQuery"
        :placeholder="$t('storageClassManagement.searchPlaceholder')"
        :prefix-icon="SearchIcon"
        clearable
        class="search-input"
        :disabled="loading"
      />
      <el-tooltip :content="$t('actions.refresh')" placement="top">
        <el-button :icon="RefreshIcon" circle @click="fetchData" :loading="loading" />
      </el-tooltip>
    </div>

    <el-card shadow="hover" class="table-card">
      <el-table
        :data="paginatedData"
        v-loading="loading"
        border
        stripe
        row-key="uid"
        :default-sort="{ prop: 'createdAt', order: 'descending' }"
        @sort-change="handleSortChange"
      >
        <el-table-column prop="name" :label="$t('storageClassManagement.columns.name')" min-width="200" sortable="custom" fixed show-overflow-tooltip>
          <template #default="{ row }">
            <el-icon class="resource-icon"><Box /></el-icon>
            <span>{{ row.name }}</span>
            <el-tag v-if="row.isDefault" size="small" type="success" style="margin-left: 8px">default</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="provisioner" :label="$t('storageClassManagement.columns.provisioner')" min-width="220" show-overflow-tooltip />
        <el-table-column prop="reclaimPolicy" :label="$t('storageClassManagement.columns.reclaimPolicy')" min-width="120" align="center">
          <template #default="{ row }">
            <el-tag size="small" :type="row.reclaimPolicy === 'Retain' ? 'success' : row.reclaimPolicy === 'Delete' ? 'danger' : 'info'">
              {{ row.reclaimPolicy || '-' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="volumeBindingMode" :label="$t('storageClassManagement.columns.volumeBindingMode')" min-width="160" show-overflow-tooltip />
        <el-table-column prop="allowVolumeExpansion" :label="$t('storageClassManagement.columns.allowVolumeExpansion')" min-width="120" align="center">
          <template #default="{ row }">
            <el-tag size="small" :type="row.allowVolumeExpansion ? 'success' : 'info'">
              {{ row.allowVolumeExpansion ? $t('common.yes') : $t('common.no') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" :label="$t('storageClassManagement.columns.createdAt')" min-width="180" sortable="custom" />
        <el-table-column :label="$t('storageClassManagement.columns.actions')" width="100" align="center" fixed="right">
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
          <el-empty :description="$t('storageClassManagement.messages.noData')" :image-size="100" />
        </template>
      </el-table>

      <div class="pagination-container" v-if="!loading && filteredData.length > 0">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
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
  Box
} from '@element-plus/icons-vue'

interface DisplayItem {
  uid: string
  name: string
  provisioner: string
  reclaimPolicy: string
  volumeBindingMode: string
  allowVolumeExpansion: boolean
  isDefault: boolean
  createdAt: string
  rawData: any
}

const { t } = useI18n()
const loading = ref(false)
const allItems = ref<DisplayItem[]>([])
const searchQuery = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const sortKey = ref('createdAt')
const sortOrder = ref<'ascending' | 'descending' | null>('descending')
const detailVisible = ref(false)
const detailItem = ref<DisplayItem | null>(null)

const filteredData = computed(() => {
  let list = [...allItems.value]
  const q = searchQuery.value.trim().toLowerCase()
  if (q) {
    list = list.filter(
      (item) =>
        item.name.toLowerCase().includes(q) ||
        item.provisioner.toLowerCase().includes(q)
    )
  }
  if (sortKey.value && sortOrder.value) {
    const dir = sortOrder.value === 'ascending' ? 1 : -1
    list.sort((a: any, b: any) => {
      const av = a[sortKey.value] ?? ''
      const bv = b[sortKey.value] ?? ''
      if (av < bv) return -1 * dir
      if (av > bv) return 1 * dir
      return 0
    })
  }
  return list
})

const paginatedData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  return filteredData.value.slice(start, start + pageSize.value)
})

const fetchData = async () => {
  loading.value = true
  try {
    const response = await kubernetesRequest<{ code: number; data: any; message: string }>({
      url: '/api/v1/storageclasses',
      method: 'GET'
    })
    if (response.code === 200 && response.data?.items) {
      allItems.value = response.data.items.map((item: any) => ({
        uid: item.metadata?.uid || item.metadata?.name,
        name: item.metadata?.name || '',
        provisioner: item.provisioner || '',
        reclaimPolicy: item.reclaimPolicy || 'Delete',
        volumeBindingMode: item.volumeBindingMode || 'Immediate',
        allowVolumeExpansion: !!item.allowVolumeExpansion,
        isDefault:
          item.metadata?.annotations?.['storageclass.kubernetes.io/is-default-class'] === 'true',
        createdAt: item.metadata?.creationTimestamp
          ? dayjs(item.metadata.creationTimestamp).format('YYYY-MM-DD HH:mm:ss')
          : '-',
        rawData: item
      }))
    } else {
      allItems.value = []
      ElMessage.error(response.message || t('storageClassManagement.messages.fetchFailed'))
    }
  } catch (error: any) {
    allItems.value = []
    ElMessage.error(error.message || t('storageClassManagement.messages.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSortChange = ({
  prop,
  order
}: {
  prop: string | null
  order: 'ascending' | 'descending' | null
}) => {
  sortKey.value = prop || 'createdAt'
  sortOrder.value = order
  currentPage.value = 1
}

const openDetail = (row: DisplayItem) => {
  detailItem.value = row
  detailVisible.value = true
}

const handleDelete = (row: DisplayItem) => {
  ElMessageBox.confirm(
    t('storageClassManagement.messages.deleteConfirm', { name: row.name }),
    t('common.confirm'),
    { type: 'warning' }
  )
    .then(async () => {
      loading.value = true
      try {
        const response = await kubernetesRequest<{ code: number; message: string }>({
          url: `/api/v1/storageclasses/${row.name}`,
          method: 'DELETE'
        })
        if (response.code === 200 || response.code === 204) {
          ElMessage.success(t('storageClassManagement.messages.deleteSuccess'))
          await fetchData()
        } else {
          ElMessage.error(response.message || t('storageClassManagement.messages.deleteFailed'))
          loading.value = false
        }
      } catch (error: any) {
        ElMessage.error(error.message || t('storageClassManagement.messages.deleteFailed'))
        loading.value = false
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
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}
.page-title {
  margin: 0;
  font-size: 28px;
  font-weight: 600;
  color: var(--el-text-color-primary);
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
  color: #0ea5e9;
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
  line-height: 1.5;
  overflow: auto;
  white-space: pre-wrap;
  word-break: break-word;
}
</style>
