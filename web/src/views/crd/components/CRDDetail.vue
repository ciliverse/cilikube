<template>
  <div class="crd-detail" v-loading="loading">
    <template v-if="crdDetail">
      <!-- 基本信息 -->
      <div class="detail-section">
        <h3 class="section-title">{{ $t('crd.basicInfo') }}</h3>
        <div class="info-grid">
          <div class="info-item">
            <span class="label">{{ $t('crd.name') }}:</span>
            <span class="value">{{ crdDetail.name }}</span>
          </div>
          <div class="info-item">
            <span class="label">{{ $t('crd.group') }}:</span>
            <span class="value">{{ crdDetail.group }}</span>
          </div>
          <div class="info-item">
            <span class="label">{{ $t('crd.kind') }}:</span>
            <span class="value">{{ crdDetail.kind }}</span>
          </div>
          <div class="info-item">
            <span class="label">{{ $t('crd.scope') }}:</span>
            <el-tag :type="crdDetail.scope === 'Namespaced' ? 'primary' : 'info'" size="small">
              {{ crdDetail.scope }}
            </el-tag>
          </div>
          <div class="info-item">
            <span class="label">{{ $t('crd.plural') }}:</span>
            <span class="value">{{ crdDetail.plural }}</span>
          </div>
          <div class="info-item">
            <span class="label">{{ $t('crd.singular') }}:</span>
            <span class="value">{{ crdDetail.singular }}</span>
          </div>
        </div>
      </div>

      <!-- 版本信息 -->
      <div class="detail-section" v-if="crdDetail.versions?.length">
        <h3 class="section-title">{{ $t('crd.versions') }}</h3>
        <el-table :data="crdDetail.versions" size="small">
          <el-table-column prop="name" :label="$t('crd.version')" />
          <el-table-column prop="served" :label="$t('crd.served')">
            <template #default="{ row }">
              <el-tag :type="row.served ? 'success' : 'danger'" size="small">
                {{ row.served ? $t('common.yes') : $t('common.no') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="storage" :label="$t('crd.storage')">
            <template #default="{ row }">
              <el-tag :type="row.storage ? 'success' : 'info'" size="small">
                {{ row.storage ? $t('common.yes') : $t('common.no') }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 状态条件 -->
      <div class="detail-section" v-if="crdDetail.conditions?.length">
        <h3 class="section-title">{{ $t('crd.conditions') }}</h3>
        <el-table :data="crdDetail.conditions" size="small">
          <el-table-column prop="type" :label="$t('crd.conditionType')" />
          <el-table-column prop="status" :label="$t('crd.conditionStatus')">
            <template #default="{ row }">
              <el-tag :type="getConditionStatusType(row.status)" size="small">
                {{ row.status }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="reason" :label="$t('crd.conditionReason')" />
          <el-table-column prop="message" :label="$t('crd.conditionMessage')" show-overflow-tooltip />
          <el-table-column prop="lastTransitionTime" :label="$t('crd.lastTransitionTime')">
            <template #default="{ row }">
              {{ formatTime(row.lastTransitionTime) }}
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 标签和注解 -->
      <div class="detail-section" v-if="crdDetail.labels || crdDetail.annotations">
        <h3 class="section-title">{{ $t('crd.metadata') }}</h3>
        
        <div v-if="crdDetail.labels && Object.keys(crdDetail.labels).length > 0" class="metadata-section">
          <h4 class="subsection-title">{{ $t('crd.labels') }}</h4>
          <div class="metadata-tags">
            <el-tag
              v-for="(value, key) in crdDetail.labels"
              :key="key"
              size="small"
              class="metadata-tag"
            >
              {{ key }}: {{ value }}
            </el-tag>
          </div>
        </div>

        <div v-if="crdDetail.annotations && Object.keys(crdDetail.annotations).length > 0" class="metadata-section">
          <h4 class="subsection-title">{{ $t('crd.annotations') }}</h4>
          <div class="metadata-tags">
            <el-tag
              v-for="(value, key) in crdDetail.annotations"
              :key="key"
              size="small"
              type="info"
              class="metadata-tag"
            >
              {{ key }}: {{ value }}
            </el-tag>
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="actions">
        <el-button type="primary" @click="handleViewResources">
          {{ $t('crd.viewResources') }}
        </el-button>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
import { useCRD } from '@/composables/useCRD'
import type { CRDItem, CRDDetailResponse } from '@/api/crd'

interface Props {
  crd: CRDItem
}

const props = defineProps<Props>()
const emit = defineEmits<{
  viewResources: [crd: CRDItem]
}>()

const { t, locale } = useI18n()
const { fetchCRDDetail } = useCRD()

const loading = ref(false)
const crdDetail = ref<CRDDetailResponse | null>(null)

// 获取条件状态类型
const getConditionStatusType = (status: string) => {
  switch (status.toLowerCase()) {
    case 'true':
      return 'success'
    case 'false':
      return 'danger'
    default:
      return 'warning'
  }
}

// 格式化时间
const formatTime = (timeStr: string) => {
  const currentLocale = locale.value === 'zh' ? 'zh-cn' : 'en'
  dayjs.locale(currentLocale)
  return dayjs(timeStr).fromNow()
}

// 处理查看资源
const handleViewResources = () => {
  emit('viewResources', props.crd)
}

// 获取CRD详情
const loadCRDDetail = async () => {
  loading.value = true
  try {
    crdDetail.value = await fetchCRDDetail(props.crd.name)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadCRDDetail()
})
</script>

<style lang="scss" scoped>
.crd-detail {
  .detail-section {
    margin-bottom: 24px;

    .section-title {
      margin: 0 0 16px 0;
      font-size: 16px;
      font-weight: 600;
      color: #303133;
      border-bottom: 1px solid #e4e7ed;
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
          color: #909399;
          margin-right: 8px;
          min-width: 80px;
          font-size: 14px;
        }

        .value {
          color: #606266;
          font-size: 14px;
        }
      }
    }

    .metadata-section {
      margin-bottom: 16px;

      .subsection-title {
        margin: 0 0 8px 0;
        font-size: 14px;
        font-weight: 500;
        color: #606266;
      }

      .metadata-tags {
        display: flex;
        flex-wrap: wrap;
        gap: 8px;

        .metadata-tag {
          max-width: 300px;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
      }
    }
  }

  .actions {
    padding-top: 16px;
    border-top: 1px solid #e4e7ed;
    text-align: right;
  }
}
</style>