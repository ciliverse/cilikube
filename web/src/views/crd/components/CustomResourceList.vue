<template>
  <div class="custom-resource-list">
    <!-- 筛选器 -->
    <div class="filters">
      <el-select
        v-if="crd.scope === 'Namespaced'"
        v-model="selectedNamespace"
        :placeholder="$t('crd.selectNamespace')"
        clearable
        style="width: 200px"
        @change="handleNamespaceChange"
      >
        <el-option
          v-for="ns in namespaces"
          :key="ns"
          :label="ns"
          :value="ns"
        />
      </el-select>
      
      <el-button
        type="primary"
        :icon="Refresh"
        @click="handleRefresh"
        :loading="loading"
        style="margin-left: 12px"
      >
        {{ $t('common.refresh') }}
      </el-button>

      <el-button
        type="success"
        :icon="Plus"
        @click="handleCreate"
        style="margin-left: 12px"
      >
        {{ $t('crd.createResource') }}
      </el-button>
    </div>

    <!-- 资源列表 -->
    <div class="resource-content">
      <el-table
        :data="resources"
        v-loading="loading"
        style="width: 100%"
        @row-click="handleRowClick"
      >
        <el-table-column prop="name" :label="$t('common.name')" min-width="150" />
        <el-table-column 
          v-if="crd.scope === 'Namespaced'"
          prop="namespace" 
          :label="$t('common.namespace')" 
          width="120" 
        />
        <el-table-column prop="createdAt" :label="$t('common.createdAt')" width="180">
          <template #default="{ row }">
            {{ formatTime(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('crd.labels')" min-width="200">
          <template #default="{ row }">
            <div v-if="row.labels && Object.keys(row.labels).length > 0" class="labels">
              <el-tag
                v-for="(value, key) in row.labels"
                :key="key"
                size="small"
                class="label-tag"
              >
                {{ key }}={{ value }}
              </el-tag>
            </div>
            <span v-else class="no-data">-</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('common.actions')" width="150" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              size="small"
              @click.stop="handleView(row)"
            >
              {{ $t('common.view') }}
            </el-button>
            <el-button
              type="danger"
              size="small"
              @click.stop="handleDelete(row)"
            >
              {{ $t('common.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 空状态 -->
      <el-empty
        v-if="!loading && resources.length === 0"
        :description="$t('crd.noResources')"
      />
    </div>

    <!-- 资源详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      :title="`${selectedResource?.name} - ${crd.kind}`"
      width="800px"
      destroy-on-close
    >
      <CustomResourceDetail
        v-if="selectedResource"
        :resource="selectedResource"
        :crd="crd"
        @updated="handleResourceUpdated"
      />
    </el-dialog>

    <!-- 创建/编辑资源对话框 -->
    <el-dialog
      v-model="editDialogVisible"
      :title="editMode === 'create' ? $t('crd.createResource') : $t('crd.editResource')"
      width="800px"
      destroy-on-close
    >
      <CustomResourceEditor
        v-if="editDialogVisible"
        :crd="crd"
        :resource="editingResource"
        :mode="editMode"
        :namespace="selectedNamespace"
        @saved="handleResourceSaved"
        @cancel="editDialogVisible = false"
      />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh, Plus } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
import { useCustomResource } from '@/composables/useCRD'
import { useNamespaces } from '@/composables/useNamespaces'
import type { CRDItem, CustomResourceItem } from '@/api/crd'
import CustomResourceDetail from './CustomResourceDetail.vue'
import CustomResourceEditor from './CustomResourceEditor.vue'

interface Props {
  crd: CRDItem
}

const props = defineProps<Props>()

const { t, locale } = useI18n()
const { 
  resources, 
  loading, 
  fetchCustomResources, 
  deleteCustomResource 
} = useCustomResource()
const { namespaces, fetchNamespaces } = useNamespaces()

const selectedNamespace = ref('')
const detailDialogVisible = ref(false)
const editDialogVisible = ref(false)
const selectedResource = ref<CustomResourceItem | null>(null)
const editingResource = ref<CustomResourceItem | null>(null)
const editMode = ref<'create' | 'edit'>('create')

// 格式化时间
const formatTime = (timeStr: string) => {
  const currentLocale = locale.value === 'zh' ? 'zh-cn' : 'en'
  dayjs.locale(currentLocale)
  return dayjs(timeStr).fromNow()
}

// 处理命名空间变化
const handleNamespaceChange = () => {
  loadResources()
}

// 处理刷新
const handleRefresh = () => {
  loadResources()
}

// 处理创建
const handleCreate = () => {
  editingResource.value = null
  editMode.value = 'create'
  editDialogVisible.value = true
}

// 处理行点击
const handleRowClick = (row: CustomResourceItem) => {
  handleView(row)
}

// 处理查看
const handleView = (resource: CustomResourceItem) => {
  selectedResource.value = resource
  detailDialogVisible.value = true
}

// 处理删除
const handleDelete = async (resource: CustomResourceItem) => {
  try {
    await ElMessageBox.confirm(
      t('crd.deleteConfirm', { name: resource.name }),
      t('common.warning'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      }
    )

    const success = await deleteCustomResource(
      props.crd.group,
      props.crd.version,
      props.crd.plural,
      resource.name,
      resource.namespace
    )

    if (success) {
      loadResources()
    }
  } catch (error) {
    // 用户取消删除
  }
}

// 处理资源更新
const handleResourceUpdated = () => {
  loadResources()
  detailDialogVisible.value = false
}

// 处理资源保存
const handleResourceSaved = () => {
  loadResources()
  editDialogVisible.value = false
}

// 加载资源列表
const loadResources = async () => {
  await fetchCustomResources(
    props.crd.group,
    props.crd.version,
    props.crd.plural,
    selectedNamespace.value
  )
}

// 初始化
onMounted(async () => {
  if (props.crd.scope === 'Namespaced') {
    await fetchNamespaces()
  }
  await loadResources()
})

// 监听CRD变化
watch(() => props.crd, () => {
  loadResources()
}, { deep: true })
</script>

<style lang="scss" scoped>
.custom-resource-list {
  height: 100%;
  display: flex;
  flex-direction: column;
  
  .filters {
    margin-bottom: 16px;
    display: flex;
    align-items: center;
    flex-shrink: 0;
  }

  .resource-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    
    .el-table {
      flex: 1;
    }
    
    .labels {
      display: flex;
      flex-wrap: wrap;
      gap: 4px;

      .label-tag {
        max-width: 150px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }

    .no-data {
      color: var(--el-text-color-placeholder);
    }
    
    .el-empty {
      flex: 1;
      display: flex;
      align-items: center;
      justify-content: center;
    }
  }
}
</style>