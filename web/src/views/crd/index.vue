<template>
  <div class="crd-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <h2 class="page-title">{{ $t('menu.crd') }}</h2>
        <p class="page-description">{{ $t('crd.description') }}</p>
      </div>
      <div class="header-actions">
        <el-button 
          type="primary" 
          :icon="Refresh" 
          @click="handleRefresh"
          :loading="loading"
        >
          {{ $t('common.refresh') }}
        </el-button>
      </div>
    </div>

    <!-- 搜索和筛选 -->
    <div class="filters">
      <el-input
        v-model="searchKeyword"
        :placeholder="$t('crd.searchPlaceholder')"
        :prefix-icon="Search"
        clearable
        style="width: 300px"
        @input="handleSearch"
      />
      <el-select
        v-model="selectedGroup"
        :placeholder="$t('crd.selectGroup')"
        clearable
        style="width: 200px; margin-left: 12px"
        @change="handleGroupFilter"
      >
        <el-option
          v-for="group in availableGroups"
          :key="group"
          :label="group"
          :value="group"
        />
      </el-select>
    </div>

    <!-- CRD列表 -->
    <div class="crd-content">
      <el-card v-loading="loading" shadow="never">
        <template v-if="filteredCRDs.length > 0">
          <div class="crd-grid">
            <div
              v-for="crd in filteredCRDs"
              :key="crd.name"
              class="crd-card"
              @click="handleCRDClick(crd)"
            >
              <div class="crd-card-header">
                <div class="crd-info">
                  <h3 class="crd-name">{{ crd.kind }}</h3>
                  <p class="crd-group">{{ crd.group }}/{{ crd.version }}</p>
                </div>
                <div class="crd-scope">
                  <el-tag :type="crd.scope === 'Namespaced' ? 'primary' : 'info'" size="small">
                    {{ crd.scope }}
                  </el-tag>
                </div>
              </div>
              
              <div class="crd-card-body">
                <div class="crd-details">
                  <div class="detail-item">
                    <span class="label">{{ $t('crd.plural') }}:</span>
                    <span class="value">{{ crd.plural }}</span>
                  </div>
                  <div class="detail-item" v-if="crd.shortNames?.length">
                    <span class="label">{{ $t('crd.shortNames') }}:</span>
                    <span class="value">{{ crd.shortNames.join(', ') }}</span>
                  </div>
                  <div class="detail-item" v-if="crd.categories?.length">
                    <span class="label">{{ $t('crd.categories') }}:</span>
                    <span class="value">{{ crd.categories.join(', ') }}</span>
                  </div>
                </div>
              </div>

              <div class="crd-card-footer">
                <div class="created-time">
                  {{ $t('common.createdAt') }}: {{ formatTime(crd.createdAt) }}
                </div>
                <div class="actions">
                  <el-button
                    type="primary"
                    size="small"
                    @click.stop="handleViewResources(crd)"
                  >
                    {{ $t('crd.viewResources') }}
                  </el-button>
                </div>
              </div>
            </div>
          </div>
        </template>
        
        <template v-else-if="!loading">
          <el-empty :description="$t('crd.noCRDs')" />
        </template>
      </el-card>
    </div>

    <!-- CRD详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      :title="selectedCRD?.kind"
      width="800px"
      destroy-on-close
    >
      <div v-if="selectedCRD" class="crd-detail-simple">
        <div class="detail-section">
          <h3 class="section-title">{{ $t('crd.basicInfo') }}</h3>
          <div class="info-grid">
            <div class="info-item">
              <span class="label">{{ $t('crd.name') }}:</span>
              <span class="value">{{ selectedCRD.name }}</span>
            </div>
            <div class="info-item">
              <span class="label">{{ $t('crd.group') }}:</span>
              <span class="value">{{ selectedCRD.group }}</span>
            </div>
            <div class="info-item">
              <span class="label">{{ $t('crd.kind') }}:</span>
              <span class="value">{{ selectedCRD.kind }}</span>
            </div>
            <div class="info-item">
              <span class="label">{{ $t('crd.scope') }}:</span>
              <el-tag :type="selectedCRD.scope === 'Namespaced' ? 'primary' : 'info'" size="small">
                {{ selectedCRD.scope }}
              </el-tag>
            </div>
          </div>
        </div>
      </div>
    </el-dialog>

    <!-- 自定义资源列表对话框 -->
    <el-dialog
      v-model="resourcesDialogVisible"
      :title="`${currentCRD?.kind} ${$t('crd.resources')}`"
      width="1200px"
      top="5vh"
      destroy-on-close
      class="resource-dialog"
    >
      <div class="resource-dialog-content">
        <CustomResourceList
          v-if="currentCRD"
          :crd="currentCRD"
        />
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Refresh, Search } from '@element-plus/icons-vue'
import { useCRD } from '@/composables/useCRD'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
import type { CRDItem } from '@/api/crd'
import CustomResourceList from './components/CustomResourceList.vue'

const { t, locale } = useI18n()
const { crds, loading, fetchCRDs } = useCRD()

// 搜索和筛选
const searchKeyword = ref('')
const selectedGroup = ref('')
const detailDialogVisible = ref(false)
const resourcesDialogVisible = ref(false)
const selectedCRD = ref<CRDItem | null>(null)
const currentCRD = ref<CRDItem | null>(null)

// 可用的组列表
const availableGroups = computed(() => {
  const groups = new Set(crds.value.map(crd => crd.group || 'core'))
  return Array.from(groups).sort()
})

// 过滤后的CRD列表
const filteredCRDs = computed(() => {
  let result = crds.value

  // 按关键词搜索
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    result = result.filter(crd =>
      crd.name.toLowerCase().includes(keyword) ||
      crd.kind.toLowerCase().includes(keyword) ||
      crd.group.toLowerCase().includes(keyword) ||
      crd.plural.toLowerCase().includes(keyword)
    )
  }

  // 按组筛选
  if (selectedGroup.value) {
    result = result.filter(crd => crd.group === selectedGroup.value)
  }

  return result
})

// 格式化时间
const formatTime = (timeStr: string) => {
  const currentLocale = locale.value === 'zh' ? 'zh-cn' : 'en'
  dayjs.locale(currentLocale)
  return dayjs(timeStr).fromNow()
}

// 处理刷新
const handleRefresh = async () => {
  await fetchCRDs()
}

// 处理搜索
const handleSearch = () => {
  // 搜索逻辑在computed中处理
}

// 处理组筛选
const handleGroupFilter = () => {
  // 筛选逻辑在computed中处理
}

// 处理CRD点击
const handleCRDClick = async (crd: CRDItem) => {
  selectedCRD.value = crd
  detailDialogVisible.value = true
}

// 处理查看资源
const handleViewResources = (crd: CRDItem) => {
  currentCRD.value = crd
  resourcesDialogVisible.value = true
}

// 初始化
onMounted(() => {
  fetchCRDs()
})
</script>

<style lang="scss" scoped>
.crd-container {
  padding: 20px;
  background-color: var(--el-bg-color-page);
  flex: 1;
  display: flex;
  flex-direction: column;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
  background: var(--el-bg-color);
  padding: 20px;
  border-radius: 8px;
  box-shadow: var(--el-box-shadow-light);

  .header-content {
    .page-title {
      margin: 0 0 8px 0;
      font-size: 24px;
      font-weight: 600;
      color: var(--el-text-color-primary);
    }

    .page-description {
      margin: 0;
      color: var(--el-text-color-regular);
      font-size: 14px;
    }
  }
}

.filters {
  margin-bottom: 20px;
  display: flex;
  align-items: center;
}

.crd-content {
  .crd-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
    gap: 16px;
  }

  .crd-card {
    background: var(--el-bg-color);
    border-radius: 8px;
    padding: 16px;
    cursor: pointer;
    transition: all 0.3s ease;
    border: 1px solid var(--el-border-color-light);

    &:hover {
      box-shadow: var(--el-box-shadow);
      transform: translateY(-2px);
    }

    .crd-card-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 12px;

      .crd-info {
        .crd-name {
          margin: 0 0 4px 0;
          font-size: 18px;
          font-weight: 600;
          color: var(--el-text-color-primary);
        }

        .crd-group {
          margin: 0;
          font-size: 12px;
          color: var(--el-text-color-secondary);
        }
      }
    }

    .crd-card-body {
      margin-bottom: 12px;

      .crd-details {
        .detail-item {
          display: flex;
          margin-bottom: 4px;
          font-size: 12px;

          .label {
            color: var(--el-text-color-secondary);
            margin-right: 8px;
            min-width: 60px;
          }

          .value {
            color: var(--el-text-color-regular);
            flex: 1;
          }
        }
      }
    }

    .crd-card-footer {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding-top: 12px;
      border-top: 1px solid var(--el-border-color-lighter);

      .created-time {
        font-size: 12px;
        color: var(--el-text-color-placeholder);
      }
    }
  }
}

.crd-detail-simple {
  .detail-section {
    margin-bottom: 24px;

    .section-title {
      margin: 0 0 16px 0;
      font-size: 16px;
      font-weight: 600;
      color: var(--el-text-color-primary);
      border-bottom: 1px solid var(--el-border-color-light);
      padding-bottom: 8px;
    }

    .info-grid {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 12px;

      .info-item {
        display: flex;
        align-items: center;

        .label {
          color: var(--el-text-color-secondary);
          margin-right: 8px;
          min-width: 80px;
          font-size: 14px;
        }

        .value {
          color: var(--el-text-color-regular);
          font-size: 14px;
        }
      }
    }
  }
}

.resources-simple {
  padding: 20px;
  text-align: center;
}

// 资源对话框样式
:deep(.resource-dialog) {
  .el-dialog {
    max-height: 90vh;
    margin-top: 5vh !important;
    display: flex;
    flex-direction: column;
  }
  
  .el-dialog__header {
    flex-shrink: 0;
  }
  
  .el-dialog__body {
    flex: 1;
    padding: 0;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }
  
  .el-dialog__footer {
    flex-shrink: 0;
  }
}

.resource-dialog-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  padding: 20px;
  
  .custom-resource-list {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    
    .filters {
      flex-shrink: 0;
      margin-bottom: 16px;
    }
    
    .resource-content {
      flex: 1;
      overflow: auto;
      
      .el-table {
        height: 100%;
      }
    }
  }
}
</style>