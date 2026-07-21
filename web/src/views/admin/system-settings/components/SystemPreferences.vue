<template>
  <div class="system-preferences">
    <div class="settings-header">
      <h2>系统首选项</h2>
      <p>配置系统界面、行为和默认设置</p>
    </div>

    <el-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-width="140px"
      class="settings-form"
    >
      <!-- 界面设置 -->
      <div class="settings-section">
        <h3>界面设置</h3>
        
        <el-form-item label="系统标题" prop="system_title">
          <el-input
            v-model="formData.system_title"
            placeholder="请输入系统标题"
            maxlength="100"
            @change="handleFormChange"
          />
          <div class="form-help">
            显示在浏览器标题栏和登录页面的系统名称
          </div>
        </el-form-item>

        <el-form-item label="系统描述" prop="system_description">
          <el-input
            v-model="formData.system_description"
            type="textarea"
            :rows="3"
            placeholder="请输入系统描述"
            maxlength="500"
            @change="handleFormChange"
          />
          <div class="form-help">
            显示在登录页面的系统描述信息
          </div>
        </el-form-item>

        <el-form-item label="默认语言" prop="default_language">
          <el-select v-model="formData.default_language" @change="handleFormChange">
            <el-option label="中文简体" value="zh-CN" />
            <el-option label="English" value="en-US" />
          </el-select>
        </el-form-item>

        <el-form-item label="默认主题" prop="default_theme">
          <el-select v-model="formData.default_theme" @change="handleFormChange">
            <el-option label="浅色主题" value="light" />
            <el-option label="深色主题" value="dark" />
            <el-option label="跟随系统" value="auto" />
          </el-select>
        </el-form-item>

        <el-form-item label="侧边栏默认状态" prop="sidebar_collapsed">
          <el-radio-group v-model="formData.sidebar_collapsed" @change="handleFormChange">
            <el-radio :label="false">展开</el-radio>
            <el-radio :label="true">收起</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="显示欢迎消息" prop="show_welcome_message">
          <el-switch
            v-model="formData.show_welcome_message"
            active-text="显示"
            inactive-text="隐藏"
            @change="handleFormChange"
          />
        </el-form-item>
      </div>

      <!-- 分页设置 -->
      <div class="settings-section">
        <h3>分页设置</h3>
        
        <el-form-item label="默认页面大小" prop="default_page_size">
          <el-select v-model="formData.default_page_size" @change="handleFormChange">
            <el-option label="10条/页" :value="10" />
            <el-option label="20条/页" :value="20" />
            <el-option label="50条/页" :value="50" />
            <el-option label="100条/页" :value="100" />
          </el-select>
        </el-form-item>

        <el-form-item label="可选页面大小">
          <el-checkbox-group v-model="formData.available_page_sizes" @change="handleFormChange">
            <el-checkbox :label="10">10</el-checkbox>
            <el-checkbox :label="20">20</el-checkbox>
            <el-checkbox :label="50">50</el-checkbox>
            <el-checkbox :label="100">100</el-checkbox>
            <el-checkbox :label="200">200</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
      </div>

      <!-- 时间和日期 -->
      <div class="settings-section">
        <h3>时间和日期</h3>
        
        <el-form-item label="时区" prop="timezone">
          <el-select
            v-model="formData.timezone"
            filterable
            placeholder="选择时区"
            @change="handleFormChange"
          >
            <el-option
              v-for="tz in timezones"
              :key="tz.value"
              :label="tz.label"
              :value="tz.value"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="日期格式" prop="date_format">
          <el-select v-model="formData.date_format" @change="handleFormChange">
            <el-option label="YYYY-MM-DD" value="YYYY-MM-DD" />
            <el-option label="MM/DD/YYYY" value="MM/DD/YYYY" />
            <el-option label="DD/MM/YYYY" value="DD/MM/YYYY" />
            <el-option label="YYYY年MM月DD日" value="YYYY年MM月DD日" />
          </el-select>
          <div class="form-help">
            当前预览: {{ formatDatePreview() }}
          </div>
        </el-form-item>

        <el-form-item label="时间格式" prop="time_format">
          <el-radio-group v-model="formData.time_format" @change="handleFormChange">
            <el-radio label="24">24小时制</el-radio>
            <el-radio label="12">12小时制</el-radio>
          </el-radio-group>
          <div class="form-help">
            当前预览: {{ formatTimePreview() }}
          </div>
        </el-form-item>
      </div>

      <!-- 数据刷新 -->
      <div class="settings-section">
        <h3>数据刷新</h3>
        
        <el-form-item label="启用自动刷新" prop="auto_refresh_enabled">
          <el-switch
            v-model="formData.auto_refresh_enabled"
            active-text="启用"
            inactive-text="禁用"
            @change="handleFormChange"
          />
        </el-form-item>

        <template v-if="formData.auto_refresh_enabled">
          <el-form-item label="默认刷新间隔" prop="default_refresh_interval">
            <el-select v-model="formData.default_refresh_interval" @change="handleFormChange">
              <el-option label="5秒" :value="5" />
              <el-option label="10秒" :value="10" />
              <el-option label="30秒" :value="30" />
              <el-option label="1分钟" :value="60" />
              <el-option label="5分钟" :value="300" />
            </el-select>
          </el-form-item>

          <el-form-item label="可选刷新间隔">
            <el-checkbox-group v-model="formData.available_refresh_intervals" @change="handleFormChange">
              <el-checkbox :label="5">5秒</el-checkbox>
              <el-checkbox :label="10">10秒</el-checkbox>
              <el-checkbox :label="30">30秒</el-checkbox>
              <el-checkbox :label="60">1分钟</el-checkbox>
              <el-checkbox :label="300">5分钟</el-checkbox>
            </el-checkbox-group>
          </el-form-item>
        </template>
      </div>

      <!-- 文件上传 -->
      <div class="settings-section">
        <h3>文件上传</h3>
        
        <el-form-item label="最大文件大小" prop="max_file_size">
          <el-input-number
            v-model="formData.max_file_size"
            :min="1"
            :max="100"
            @change="handleFormChange"
          />
          <span class="input-suffix">MB</span>
        </el-form-item>

        <el-form-item label="允许的文件类型">
          <el-checkbox-group v-model="formData.allowed_file_types" @change="handleFormChange">
            <el-checkbox label="image">图片文件 (jpg, png, gif)</el-checkbox>
            <el-checkbox label="document">文档文件 (pdf, doc, txt)</el-checkbox>
            <el-checkbox label="archive">压缩文件 (zip, tar, gz)</el-checkbox>
            <el-checkbox label="yaml">配置文件 (yaml, json, xml)</el-checkbox>
          </el-checkbox-group>
        </el-form-item>

        <el-form-item label="头像文件大小限制" prop="avatar_max_size">
          <el-input-number
            v-model="formData.avatar_max_size"
            :min="0.1"
            :max="10"
            :step="0.1"
            :precision="1"
            @change="handleFormChange"
          />
          <span class="input-suffix">MB</span>
        </el-form-item>
      </div>

      <!-- 通知设置 -->
      <div class="settings-section">
        <h3>通知设置</h3>
        
        <el-form-item label="启用系统通知" prop="notifications_enabled">
          <el-switch
            v-model="formData.notifications_enabled"
            active-text="启用"
            inactive-text="禁用"
            @change="handleFormChange"
          />
        </el-form-item>

        <template v-if="formData.notifications_enabled">
          <el-form-item label="通知类型">
            <el-checkbox-group v-model="formData.notification_types" @change="handleFormChange">
              <el-checkbox label="system">系统消息</el-checkbox>
              <el-checkbox label="security">安全警告</el-checkbox>
              <el-checkbox label="update">更新提醒</el-checkbox>
              <el-checkbox label="maintenance">维护通知</el-checkbox>
            </el-checkbox-group>
          </el-form-item>

          <el-form-item label="通知保留时间" prop="notification_retention_days">
            <el-input-number
              v-model="formData.notification_retention_days"
              :min="1"
              :max="90"
              @change="handleFormChange"
            />
            <span class="input-suffix">天</span>
          </el-form-item>
        </template>
      </div>

      <!-- 性能设置 -->
      <div class="settings-section">
        <h3>性能设置</h3>
        
        <el-form-item label="启用缓存" prop="cache_enabled">
          <el-switch
            v-model="formData.cache_enabled"
            active-text="启用"
            inactive-text="禁用"
            @change="handleFormChange"
          />
        </el-form-item>

        <el-form-item label="缓存过期时间" prop="cache_ttl">
          <el-input-number
            v-model="formData.cache_ttl"
            :min="60"
            :max="3600"
            @change="handleFormChange"
          />
          <span class="input-suffix">秒</span>
        </el-form-item>

        <el-form-item label="启用数据压缩" prop="compression_enabled">
          <el-switch
            v-model="formData.compression_enabled"
            active-text="启用"
            inactive-text="禁用"
            @change="handleFormChange"
          />
        </el-form-item>

        <el-form-item label="API请求超时时间" prop="api_timeout">
          <el-input-number
            v-model="formData.api_timeout"
            :min="5"
            :max="120"
            @change="handleFormChange"
          />
          <span class="input-suffix">秒</span>
        </el-form-item>
      </div>
    </el-form>

    <!-- 操作按钮 -->
    <div class="settings-actions">
      <el-button @click="resetForm">重置</el-button>
      <el-button
        type="primary"
        @click="saveSettings"
        :loading="loading"
      >
        保存设置
      </el-button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted } from "vue"
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from "element-plus"

interface Props {
  loading: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
  "update:loading": [value: boolean]
  save: [data: any]
}>()

const formRef = ref<FormInstance>()

const formData = reactive({
  // 界面设置
  system_title: "CiliKube",
  system_description: "Kubernetes集群管理平台",
  default_language: "zh-CN",
  default_theme: "light",
  sidebar_collapsed: false,
  show_welcome_message: true,
  
  // 分页设置
  default_page_size: 20,
  available_page_sizes: [10, 20, 50, 100],
  
  // 时间和日期
  timezone: "Asia/Shanghai",
  date_format: "YYYY-MM-DD",
  time_format: "24",
  
  // 数据刷新
  auto_refresh_enabled: true,
  default_refresh_interval: 30,
  available_refresh_intervals: [10, 30, 60, 300],
  
  // 文件上传
  max_file_size: 10,
  allowed_file_types: ["image", "document", "yaml"],
  avatar_max_size: 2.0,
  
  // 通知设置
  notifications_enabled: true,
  notification_types: ["system", "security", "update"],
  notification_retention_days: 30,
  
  // 性能设置
  cache_enabled: true,
  cache_ttl: 300,
  compression_enabled: true,
  api_timeout: 30
})

const timezones = [
  { label: "Asia/Shanghai (UTC+8)", value: "Asia/Shanghai" },
  { label: "Asia/Tokyo (UTC+9)", value: "Asia/Tokyo" },
  { label: "Asia/Seoul (UTC+9)", value: "Asia/Seoul" },
  { label: "Asia/Singapore (UTC+8)", value: "Asia/Singapore" },
  { label: "Asia/Hong_Kong (UTC+8)", value: "Asia/Hong_Kong" },
  { label: "Europe/London (UTC+0)", value: "Europe/London" },
  { label: "Europe/Paris (UTC+1)", value: "Europe/Paris" },
  { label: "America/New_York (UTC-5)", value: "America/New_York" },
  { label: "America/Los_Angeles (UTC-8)", value: "America/Los_Angeles" },
  { label: "UTC", value: "UTC" }
]

const formRules: FormRules = {
  system_title: [
    { required: true, message: "请输入系统标题", trigger: "blur" },
    { min: 2, max: 100, message: "标题长度在 2 到 100 个字符", trigger: "blur" }
  ],
  default_language: [
    { required: true, message: "请选择默认语言", trigger: "change" }
  ],
  default_theme: [
    { required: true, message: "请选择默认主题", trigger: "change" }
  ],
  timezone: [
    { required: true, message: "请选择时区", trigger: "change" }
  ],
  max_file_size: [
    { required: true, message: "请设置最大文件大小", trigger: "blur" }
  ],
  api_timeout: [
    { required: true, message: "请设置API超时时间", trigger: "blur" }
  ]
}

/** 表单变更处理 */
const handleFormChange = () => {
  emit("update:loading", false)
}

/** 格式化日期预览 */
const formatDatePreview = () => {
  const now = new Date()
  const formats: Record<string, string> = {
    "YYYY-MM-DD": now.toISOString().split('T')[0],
    "MM/DD/YYYY": `${(now.getMonth() + 1).toString().padStart(2, '0')}/${now.getDate().toString().padStart(2, '0')}/${now.getFullYear()}`,
    "DD/MM/YYYY": `${now.getDate().toString().padStart(2, '0')}/${(now.getMonth() + 1).toString().padStart(2, '0')}/${now.getFullYear()}`,
    "YYYY年MM月DD日": `${now.getFullYear()}年${(now.getMonth() + 1).toString().padStart(2, '0')}月${now.getDate().toString().padStart(2, '0')}日`
  }
  return formats[formData.date_format] || formData.date_format
}

/** 格式化时间预览 */
const formatTimePreview = () => {
  const now = new Date()
  if (formData.time_format === "12") {
    return now.toLocaleTimeString('en-US', { hour12: true })
  } else {
    return now.toLocaleTimeString('en-US', { hour12: false })
  }
}

/** 重置表单 */
const resetForm = async () => {
  try {
    await ElMessageBox.confirm(
      "确定要重置所有系统首选项吗？未保存的更改将丢失。",
      "确认重置",
      {
        confirmButtonText: "确定重置",
        cancelButtonText: "取消",
        type: "warning"
      }
    )
    
    loadSettings()
  } catch {
    // 用户取消
  }
}

/** 保存设置 */
const saveSettings = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    emit("save", { ...formData })
  } catch (error) {
    ElMessage.error("请检查表单输入")
  }
}

/** 加载设置 */
const loadSettings = async () => {
  try {
    emit("update:loading", true)
    
    // 这里调用API加载系统首选项
    // const { data } = await getSystemPreferencesApi()
    // Object.assign(formData, data)
    
    // 模拟数据已在reactive中设置
    
  } catch (error: any) {
    console.error("加载系统首选项失败:", error)
    ElMessage.error("加载系统首选项失败")
  } finally {
    emit("update:loading", false)
  }
}

/** 组件挂载时加载设置 */
onMounted(() => {
  loadSettings()
})
</script>

<style lang="scss" scoped>
.system-preferences {
  .settings-header {
    margin-bottom: 32px;
    
    h2 {
      margin: 0 0 8px 0;
      color: var(--el-text-color-primary);
      font-size: 20px;
      font-weight: 600;
    }
    
    p {
      margin: 0;
      color: var(--el-text-color-regular);
    }
  }
  
  .settings-form {
    .settings-section {
      margin-bottom: 40px;
      padding-bottom: 32px;
      border-bottom: 1px solid var(--el-border-color-lighter);
      
      &:last-child {
        border-bottom: none;
        margin-bottom: 0;
      }
      
      h3 {
        margin: 0 0 24px 0;
        color: var(--el-text-color-primary);
        font-size: 16px;
        font-weight: 600;
      }
    }
    
    .form-help {
      margin-top: 4px;
      font-size: 12px;
      color: var(--el-text-color-secondary);
      line-height: 1.4;
    }
    
    .input-suffix {
      margin-left: 8px;
      color: var(--el-text-color-regular);
      font-size: 14px;
    }
  }
  
  .settings-actions {
    margin-top: 32px;
    padding-top: 24px;
    border-top: 1px solid var(--el-border-color-lighter);
    display: flex;
    justify-content: flex-end;
    gap: 12px;
  }
}

:deep(.el-checkbox-group) {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>