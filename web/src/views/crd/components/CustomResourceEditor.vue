<template>
  <div class="custom-resource-editor">
    <el-form
      ref="formRef"
      :model="formData"
      :rules="rules"
      label-width="120px"
      @submit.prevent
    >
      <!-- 基本信息 -->
      <div class="form-section">
        <h3 class="section-title">{{ $t('crd.basicInfo') }}</h3>
        
        <el-form-item :label="$t('common.name')" prop="name">
          <el-input
            v-model="formData.name"
            :disabled="mode === 'edit'"
            :placeholder="$t('crd.enterName')"
          />
        </el-form-item>

        <el-form-item
          v-if="crd.scope === 'Namespaced'"
          :label="$t('common.namespace')"
          prop="namespace"
        >
          <el-select
            v-model="formData.namespace"
            :disabled="mode === 'edit'"
            :placeholder="$t('crd.selectNamespace')"
            style="width: 100%"
          >
            <el-option
              v-for="ns in namespaces"
              :key="ns"
              :label="ns"
              :value="ns"
            />
          </el-select>
        </el-form-item>
      </div>

      <!-- 标签 -->
      <div class="form-section">
        <h3 class="section-title">{{ $t('crd.labels') }}</h3>
        <div class="key-value-editor">
          <div
            v-for="(label, index) in formData.labels"
            :key="index"
            class="key-value-row"
          >
            <el-input
              v-model="label.key"
              :placeholder="$t('crd.labelKey')"
              style="width: 200px"
            />
            <el-input
              v-model="label.value"
              :placeholder="$t('crd.labelValue')"
              style="width: 200px; margin-left: 8px"
            />
            <el-button
              type="danger"
              size="small"
              :icon="Delete"
              @click="removeLabel(index)"
              style="margin-left: 8px"
            />
          </div>
          <el-button
            type="primary"
            size="small"
            :icon="Plus"
            @click="addLabel"
          >
            {{ $t('crd.addLabel') }}
          </el-button>
        </div>
      </div>

      <!-- Spec编辑器 -->
      <div class="form-section">
        <h3 class="section-title">{{ $t('crd.spec') }}</h3>
        <div class="editor-container">
          <el-input
            v-model="specText"
            type="textarea"
            :rows="15"
            :placeholder="$t('crd.specPlaceholder')"
            resize="vertical"
            @blur="validateSpec"
          />
          <div v-if="specError" class="error-message">
            {{ specError }}
          </div>
        </div>
      </div>
    </el-form>

    <!-- 操作按钮 -->
    <div class="actions">
      <el-button @click="handleCancel">
        {{ $t('common.cancel') }}
      </el-button>
      <el-button
        type="primary"
        @click="handleSave"
        :loading="saving"
      >
        {{ $t('common.save') }}
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { Plus, Delete } from '@element-plus/icons-vue'
import { useCustomResource } from '@/composables/useCRD'
import { useNamespaces } from '@/composables/useNamespaces'
import type { CRDItem, CustomResourceItem } from '@/api/crd'

interface Props {
  crd: CRDItem
  resource?: CustomResourceItem | null
  mode: 'create' | 'edit'
  namespace?: string
}

const props = defineProps<Props>()
const emit = defineEmits<{
  saved: []
  cancel: []
}>()

const { t } = useI18n()
const { createCustomResource, updateCustomResource } = useCustomResource()
const { namespaces, fetchNamespaces } = useNamespaces()

const formRef = ref<FormInstance>()
const saving = ref(false)
const specError = ref('')

// 表单数据
const formData = reactive({
  name: '',
  namespace: '',
  labels: [] as Array<{ key: string; value: string }>
})

// Spec文本
const specText = ref('')

// 表单验证规则
const rules: FormRules = {
  name: [
    { required: true, message: t('crd.nameRequired'), trigger: 'blur' },
    { 
      pattern: /^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/, 
      message: t('crd.nameFormat'), 
      trigger: 'blur' 
    }
  ],
  namespace: props.crd.scope === 'Namespaced' ? [
    { required: true, message: t('crd.namespaceRequired'), trigger: 'change' }
  ] : []
}

// 添加标签
const addLabel = () => {
  formData.labels.push({ key: '', value: '' })
}

// 删除标签
const removeLabel = (index: number) => {
  formData.labels.splice(index, 1)
}

// 验证Spec
const validateSpec = () => {
  specError.value = ''
  if (specText.value.trim()) {
    try {
      JSON.parse(specText.value)
    } catch (error) {
      specError.value = t('crd.invalidJSON')
    }
  }
}

// 处理取消
const handleCancel = () => {
  emit('cancel')
}

// 处理保存
const handleSave = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    
    // 验证Spec
    validateSpec()
    if (specError.value) {
      ElMessage.error(specError.value)
      return
    }

    saving.value = true

    // 构建资源对象
    const metadata: any = {
      name: formData.name
    }

    if (props.crd.scope === 'Namespaced') {
      metadata.namespace = formData.namespace
    }

    // 添加标签
    const labels: Record<string, string> = {}
    formData.labels.forEach(label => {
      if (label.key && label.value) {
        labels[label.key] = label.value
      }
    })
    if (Object.keys(labels).length > 0) {
      metadata.labels = labels
    }

    const resource = {
      apiVersion: `${props.crd.group}/${props.crd.version}`,
      kind: props.crd.kind,
      metadata,
      spec: specText.value.trim() ? JSON.parse(specText.value) : {}
    }

    let success = false
    if (props.mode === 'create') {
      success = await createCustomResource(
        props.crd.group,
        props.crd.version,
        props.crd.plural,
        resource,
        formData.namespace
      )
    } else {
      success = await updateCustomResource(
        props.crd.group,
        props.crd.version,
        props.crd.plural,
        formData.name,
        resource,
        formData.namespace
      )
    }

    if (success) {
      emit('saved')
    }
  } catch (error) {
    console.error('Form validation failed:', error)
  } finally {
    saving.value = false
  }
}

// 初始化表单数据
const initFormData = () => {
  if (props.resource) {
    formData.name = props.resource.name
    formData.namespace = props.resource.namespace || ''
    
    // 初始化标签
    formData.labels = []
    if (props.resource.labels) {
      Object.entries(props.resource.labels).forEach(([key, value]) => {
        formData.labels.push({ key, value })
      })
    }

    // 初始化Spec
    if (props.resource.spec) {
      specText.value = JSON.stringify(props.resource.spec, null, 2)
    }
  } else {
    // 创建模式
    formData.name = ''
    formData.namespace = props.namespace || ''
    formData.labels = []
    specText.value = JSON.stringify({}, null, 2)
  }
}

// 初始化
onMounted(async () => {
  if (props.crd.scope === 'Namespaced') {
    await fetchNamespaces()
  }
  initFormData()
})

// 监听资源变化
watch(() => props.resource, () => {
  initFormData()
}, { deep: true })
</script>

<style lang="scss" scoped>
.custom-resource-editor {
  .form-section {
    margin-bottom: 24px;

    .section-title {
      margin: 0 0 16px 0;
      font-size: 16px;
      font-weight: 600;
      color: #303133;
      border-bottom: 1px solid #e4e7ed;
      padding-bottom: 8px;
    }

    .key-value-editor {
      .key-value-row {
        display: flex;
        align-items: center;
        margin-bottom: 8px;
      }
    }

    .editor-container {
      .error-message {
        color: #f56c6c;
        font-size: 12px;
        margin-top: 4px;
      }

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