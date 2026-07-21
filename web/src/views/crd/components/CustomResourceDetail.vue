<template>
  <div class="custom-resource-detail">
    <!-- 基本信息 -->
    <div class="detail-section">
      <h3 class="section-title">{{ $t('crd.basicInfo') }}</h3>
      <div class="info-grid">
        <div class="info-item">
          <span class="label">{{ $t('common.name') }}:</span>
          <span class="value">{{ resource.name }}</span>
        </div>
        <div class="info-item" v-if="resource.namespace">
          <span class="label">{{ $t('common.namespace') }}:</span>
          <span class="value">{{ resource.namespace }}</span>
        </div>
        <div class="info-item">
          <span class="label">{{ $t('crd.kind') }}:</span>
          <span class="value">{{ resource.kind }}</span>
        </div>
        <div class="info-item">
          <span class="label">{{ $t('crd.apiVersion') }}:</span>
          <span class="value">{{ resource.apiVersion }}</span>
        </div>
        <div class="info-item">
          <span class="label">{{ $t('common.createdAt') }}:</span>
          <span class="value">{{ formatTime(resource.createdAt) }}</span>
        </div>
      </div>
    </div>

    <!-- 标签和注解 -->
    <div class="detail-section" v-if="resource.labels || resource.annotations">
      <h3 class="section-title">{{ $t('crd.metadata') }}</h3>
      
      <div v-if="resource.labels && Object.keys(resource.labels).length > 0" class="metadata-section">
        <h4 class="subsection-title">{{ $t('crd.labels') }}</h4>
        <div class="metadata-tags">
          <el-tag
            v-for="(value, key) in resource.labels"
            :key="key"
            size="small"
            class="metadata-tag"
          >
            {{ key }}: {{ value }}
          </el-tag>
        </div>
      </div>

      <div v-if="resource.annotations && Object.keys(resource.annotations).length > 0" class="metadata-section">
        <h4 class="subsection-title">{{ $t('crd.annotations') }}</h4>
        <div class="metadata-tags">
          <el-tag
            v-for="(value, key) in resource.annotations"
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

    <!-- Spec -->
    <div class="detail-section" v-if="resource.spec">
      <h3 class="section-title">{{ $t('crd.spec') }}</h3>
      <div class="json-viewer">
        <el-input
          type="textarea"
          :model-value="formatJSON(resource.spec)"
          readonly
          :rows="10"
          resize="vertical"
        />
      </div>
    </div>

    <!-- Status -->
    <div class="detail-section" v-if="resource.status">
      <h3 class="section-title">{{ $t('crd.status') }}</h3>
      <div class="json-viewer">
        <el-input
          type="textarea"
          :model-value="formatJSON(resource.status)"
          readonly
          :rows="8"
          resize="vertical"
        />
      </div>
    </div>

    <!-- 操作按钮 -->
    <div class="actions">
      <el-button type="primary" @click="handleEdit">
        {{ $t('common.edit') }}
      </el-button>
      <el-button type="danger" @click="handleDelete">
        {{ $t('common.delete') }}
      </el-button>
    </div>

    <!-- 编辑对话框 -->
    <el-dialog
      v-model="editDialogVisible"
      :title="$t('crd.editResource')"
      width="800px"
      destroy-on-close
    >
      <CustomResourceEditor
        v-if="editDialogVisible"
        :crd="crd"
        :resource="resource"
        mode="edit"
        @saved="handleSaved"
        @cancel="editDialogVisible = false"
      />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
import { useCustomResource } from '@/composables/useCRD'
import type { CRDItem, CustomResourceItem } from '@/api/crd'
import CustomResourceEditor from './CustomResourceEditor.vue'

interface Props {
  resource: CustomResourceItem
  crd: CRDItem
}

const props = defineProps<Props>()
const emit = defineEmits<{
  updated: []
}>()

const { t, locale } = useI18n()
const { deleteCustomResource } = useCustomResource()

const editDialogVisible = ref(false)

// 格式化时间
const formatTime = (timeStr: string) => {
  const currentLocale = locale.value === 'zh' ? 'zh-cn' : 'en'
  dayjs.locale(currentLocale)
  return dayjs(timeStr).fromNow()
}

// 格式化JSON
const formatJSON = (obj: any) => {
  try {
    return JSON.stringify(obj, null, 2)
  } catch (error) {
    return String(obj)
  }
}

// 处理编辑
const handleEdit = () => {
  editDialogVisible.value = true
}

// 处理删除
const handleDelete = async () => {
  try {
    await ElMessageBox.confirm(
      t('crd.deleteConfirm', { name: props.resource.name }),
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
      props.resource.name,
      props.resource.namespace
    )

    if (success) {
      emit('updated')
    }
  } catch (error) {
    // 用户取消删除
  }
}

// 处理保存
const handleSaved = () => {
  editDialogVisible.value = false
  emit('updated')
}
</script>

<style lang="scss" scoped>
.custom-resource-detail {
  padding: 20px;
  background-color: #fff;
  border-radius: 8px;
  margin-bottom: 20px;
  
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
          min-width: 100px;
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

    .json-viewer {
      :deep(.el-textarea__inner) {
        font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
        font-size: 12px;
        line-height: 1.4;
      }
    }
  }

  .actions {
    padding-top: 16px;
    border-top: 1px solid #e4e7ed;
    text-align: right;

    .el-button + .el-button {
      margin-left: 12px;
    }
  }
}
</style>