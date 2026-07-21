<template>
  <div class="email-settings">
    <div class="settings-header">
      <h2>邮件配置</h2>
      <p>配置系统邮件发送和通知设置</p>
    </div>

    <el-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-width="140px"
      class="settings-form"
    >
      <!-- 邮件服务器设置 -->
      <div class="settings-section">
        <h3>SMTP服务器设置</h3>
        
        <el-form-item label="启用邮件服务" prop="email_enabled">
          <el-switch
            v-model="formData.email_enabled"
            active-text="启用"
            inactive-text="禁用"
            @change="handleFormChange"
          />
        </el-form-item>

        <template v-if="formData.email_enabled">
          <el-form-item label="SMTP服务器" prop="smtp_host">
            <el-input
              v-model="formData.smtp_host"
              placeholder="请输入SMTP服务器地址"
              @change="handleFormChange"
            />
          </el-form-item>

          <el-form-item label="SMTP端口" prop="smtp_port">
            <el-input-number
              v-model="formData.smtp_port"
              :min="1"
              :max="65535"
              @change="handleFormChange"
            />
            <div class="form-help">
              常用端口：25 (不加密), 465 (SSL), 587 (TLS)
            </div>
          </el-form-item>

          <el-form-item label="加密方式" prop="smtp_encryption">
            <el-radio-group v-model="formData.smtp_encryption" @change="handleFormChange">
              <el-radio label="none">无加密</el-radio>
              <el-radio label="tls">TLS</el-radio>
              <el-radio label="ssl">SSL</el-radio>
            </el-radio-group>
          </el-form-item>

          <el-form-item label="用户名" prop="smtp_username">
            <el-input
              v-model="formData.smtp_username"
              placeholder="请输入SMTP用户名"
              @change="handleFormChange"
            />
          </el-form-item>

          <el-form-item label="密码" prop="smtp_password">
            <el-input
              v-model="formData.smtp_password"
              type="password"
              placeholder="请输入SMTP密码"
              show-password
              @change="handleFormChange"
            />
          </el-form-item>

          <el-form-item label="发件人邮箱" prop="from_email">
            <el-input
              v-model="formData.from_email"
              placeholder="请输入发件人邮箱"
              @change="handleFormChange"
            />
          </el-form-item>

          <el-form-item label="发件人名称" prop="from_name">
            <el-input
              v-model="formData.from_name"
              placeholder="请输入发件人名称"
              @change="handleFormChange"
            />
          </el-form-item>
        </template>
      </div>

      <!-- 邮件模板设置 -->
      <div v-if="formData.email_enabled" class="settings-section">
        <h3>邮件模板设置</h3>
        
        <el-form-item label="邮件模板类型">
          <el-tabs v-model="activeTemplate" type="card">
            <el-tab-pane label="欢迎邮件" name="welcome" />
            <el-tab-pane label="密码重置" name="password_reset" />
            <el-tab-pane label="账户激活" name="account_activation" />
            <el-tab-pane label="安全警告" name="security_alert" />
          </el-tabs>
        </el-form-item>

        <div class="template-editor">
          <el-form-item label="邮件主题" :prop="`templates.${activeTemplate}.subject`">
            <el-input
              v-model="formData.templates[activeTemplate].subject"
              placeholder="请输入邮件主题"
              @change="handleFormChange"
            />
          </el-form-item>

          <el-form-item label="邮件内容" :prop="`templates.${activeTemplate}.content`">
            <el-input
              v-model="formData.templates[activeTemplate].content"
              type="textarea"
              :rows="10"
              placeholder="请输入邮件内容，支持HTML格式"
              @change="handleFormChange"
            />
            <div class="form-help">
              可用变量：{{username}}, {{email}}, {{system_name}}, {{reset_link}}, {{activation_link}}
            </div>
          </el-form-item>

          <el-form-item>
            <el-button @click="previewTemplate">预览邮件</el-button>
            <el-button @click="resetTemplate">重置为默认</el-button>
          </el-form-item>
        </div>
      </div>

      <!-- 邮件通知设置 -->
      <div v-if="formData.email_enabled" class="settings-section">
        <h3>邮件通知设置</h3>
        
        <el-form-item label="启用邮件通知" prop="notifications_enabled">
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
              <el-checkbox label="user_registration">用户注册</el-checkbox>
              <el-checkbox label="password_change">密码修改</el-checkbox>
              <el-checkbox label="login_alert">异常登录</el-checkbox>
              <el-checkbox label="system_maintenance">系统维护</el-checkbox>
              <el-checkbox label="security_incident">安全事件</el-checkbox>
            </el-checkbox-group>
          </el-form-item>

          <el-form-item label="管理员邮箱" prop="admin_emails">
            <div class="email-list-container">
              <div
                v-for="(email, index) in formData.admin_emails"
                :key="index"
                class="email-item"
              >
                <el-input
                  v-model="formData.admin_emails[index]"
                  placeholder="请输入管理员邮箱"
                  @change="handleFormChange"
                />
                <el-button
                  type="danger"
                  size="small"
                  @click="removeAdminEmail(index)"
                >
                  删除
                </el-button>
              </div>
              <el-button
                type="primary"
                size="small"
                @click="addAdminEmail"
              >
                添加邮箱
              </el-button>
            </div>
            <div class="form-help">
              接收系统通知的管理员邮箱地址
            </div>
          </el-form-item>

          <el-form-item label="邮件发送频率限制" prop="rate_limit">
            <el-input-number
              v-model="formData.rate_limit"
              :min="1"
              :max="100"
              @change="handleFormChange"
            />
            <span class="input-suffix">封/小时</span>
            <div class="form-help">
              防止邮件轰炸的频率限制
            </div>
          </el-form-item>
        </template>
      </div>

      <!-- 邮件队列设置 -->
      <div v-if="formData.email_enabled" class="settings-section">
        <h3>邮件队列设置</h3>
        
        <el-form-item label="启用邮件队列" prop="queue_enabled">
          <el-switch
            v-model="formData.queue_enabled"
            active-text="启用"
            inactive-text="禁用"
            @change="handleFormChange"
          />
          <div class="form-help">
            启用后邮件将异步发送，提高系统响应速度
          </div>
        </el-form-item>

        <template v-if="formData.queue_enabled">
          <el-form-item label="队列处理间隔" prop="queue_interval">
            <el-input-number
              v-model="formData.queue_interval"
              :min="10"
              :max="300"
              @change="handleFormChange"
            />
            <span class="input-suffix">秒</span>
          </el-form-item>

          <el-form-item label="批量发送数量" prop="batch_size">
            <el-input-number
              v-model="formData.batch_size"
              :min="1"
              :max="50"
              @change="handleFormChange"
            />
            <span class="input-suffix">封/批次</span>
          </el-form-item>

          <el-form-item label="重试次数" prop="max_retries">
            <el-input-number
              v-model="formData.max_retries"
              :min="0"
              :max="5"
              @change="handleFormChange"
            />
            <span class="input-suffix">次</span>
          </el-form-item>
        </template>
      </div>
    </el-form>

    <!-- 测试邮件 -->
    <div v-if="formData.email_enabled" class="test-email-section">
      <h3>测试邮件发送</h3>
      <div class="test-email-form">
        <el-input
          v-model="testEmail"
          placeholder="请输入测试邮箱地址"
          style="width: 300px; margin-right: 12px;"
        />
        <el-button
          type="primary"
          @click="sendTestEmail"
          :loading="testingEmail"
        >
          发送测试邮件
        </el-button>
      </div>
      <div class="form-help">
        发送测试邮件以验证SMTP配置是否正确
      </div>
    </div>

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

    <!-- 邮件预览对话框 -->
    <el-dialog
      v-model="showPreview"
      title="邮件预览"
      width="600px"
    >
      <div class="email-preview">
        <div class="preview-header">
          <strong>主题：</strong>{{ previewData.subject }}
        </div>
        <div class="preview-content" v-html="previewData.content"></div>
      </div>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted } from "vue"
import { ElMessage, ElMessageBox } from "element-plus"
import type { FormInstance, FormRules } from "element-plus"

interface Props {
  loading: boolean
}

const props = defineProps<Props>()
const emit = defineEmits<{
  "update:loading": [value: boolean]
  save: [data: any]
}>()

const formRef = ref<FormInstance>()
const activeTemplate = ref("welcome")
const showPreview = ref(false)
const testEmail = ref("")
const testingEmail = ref(false)

const previewData = reactive({
  subject: "",
  content: ""
})

const formData = reactive({
  // SMTP设置
  email_enabled: false,
  smtp_host: "",
  smtp_port: 587,
  smtp_encryption: "tls",
  smtp_username: "",
  smtp_password: "",
  from_email: "",
  from_name: "CiliKube System",
  
  // 邮件模板
  templates: {
    welcome: {
      subject: "欢迎使用 {{system_name}}",
      content: `<h2>欢迎，{{username}}！</h2>
<p>感谢您注册 {{system_name}}。您的账户已成功创建。</p>
<p>您可以使用以下信息登录：</p>
<ul>
  <li>邮箱：{{email}}</li>
  <li>用户名：{{username}}</li>
</ul>
<p>如有任何问题，请联系管理员。</p>`
    },
    password_reset: {
      subject: "{{system_name}} - 密码重置请求",
      content: `<h2>密码重置请求</h2>
<p>您好，{{username}}！</p>
<p>我们收到了您的密码重置请求。请点击下面的链接重置您的密码：</p>
<p><a href="{{reset_link}}" style="background: #409EFF; color: white; padding: 10px 20px; text-decoration: none; border-radius: 4px;">重置密码</a></p>
<p>如果您没有请求重置密码，请忽略此邮件。</p>
<p>此链接将在24小时后失效。</p>`
    },
    account_activation: {
      subject: "{{system_name}} - 账户激活",
      content: `<h2>账户激活</h2>
<p>您好，{{username}}！</p>
<p>请点击下面的链接激活您的账户：</p>
<p><a href="{{activation_link}}" style="background: #67C23A; color: white; padding: 10px 20px; text-decoration: none; border-radius: 4px;">激活账户</a></p>
<p>如果您没有注册此账户，请忽略此邮件。</p>`
    },
    security_alert: {
      subject: "{{system_name}} - 安全警告",
      content: `<h2>安全警告</h2>
<p>您好，{{username}}！</p>
<p>我们检测到您的账户存在异常活动。如果这不是您的操作，请立即：</p>
<ul>
  <li>修改您的密码</li>
  <li>检查账户设置</li>
  <li>联系管理员</li>
</ul>
<p>如果这是您的正常操作，请忽略此邮件。</p>`
    }
  },
  
  // 通知设置
  notifications_enabled: true,
  notification_types: ["user_registration", "password_change", "login_alert"],
  admin_emails: [] as string[],
  rate_limit: 10,
  
  // 队列设置
  queue_enabled: true,
  queue_interval: 30,
  batch_size: 10,
  max_retries: 3
})

const formRules: FormRules = {
  smtp_host: [
    { required: true, message: "请输入SMTP服务器地址", trigger: "blur" }
  ],
  smtp_port: [
    { required: true, message: "请输入SMTP端口", trigger: "blur" }
  ],
  smtp_username: [
    { required: true, message: "请输入SMTP用户名", trigger: "blur" }
  ],
  smtp_password: [
    { required: true, message: "请输入SMTP密码", trigger: "blur" }
  ],
  from_email: [
    { required: true, message: "请输入发件人邮箱", trigger: "blur" },
    { type: "email", message: "请输入有效的邮箱地址", trigger: "blur" }
  ],
  from_name: [
    { required: true, message: "请输入发件人名称", trigger: "blur" }
  ]
}

/** 表单变更处理 */
const handleFormChange = () => {
  emit("update:loading", false)
}

/** 添加管理员邮箱 */
const addAdminEmail = () => {
  formData.admin_emails.push("")
  handleFormChange()
}

/** 移除管理员邮箱 */
const removeAdminEmail = (index: number) => {
  formData.admin_emails.splice(index, 1)
  handleFormChange()
}

/** 预览邮件模板 */
const previewTemplate = () => {
  const template = formData.templates[activeTemplate.value]
  previewData.subject = template.subject.replace(/\{\{(\w+)\}\}/g, (match, key) => {
    const replacements: Record<string, string> = {
      username: "示例用户",
      email: "user@example.com",
      system_name: "CiliKube",
      reset_link: "https://example.com/reset-password?token=example",
      activation_link: "https://example.com/activate?token=example"
    }
    return replacements[key] || match
  })
  
  previewData.content = template.content.replace(/\{\{(\w+)\}\}/g, (match, key) => {
    const replacements: Record<string, string> = {
      username: "示例用户",
      email: "user@example.com",
      system_name: "CiliKube",
      reset_link: "https://example.com/reset-password?token=example",
      activation_link: "https://example.com/activate?token=example"
    }
    return replacements[key] || match
  })
  
  showPreview.value = true
}

/** 重置模板为默认 */
const resetTemplate = async () => {
  try {
    await ElMessageBox.confirm(
      "确定要重置当前模板为默认内容吗？",
      "确认重置",
      {
        confirmButtonText: "确定重置",
        cancelButtonText: "取消",
        type: "warning"
      }
    )
    
    // 重置为默认模板（这里可以从API获取默认模板）
    ElMessage.success("模板已重置为默认内容")
    handleFormChange()
  } catch {
    // 用户取消
  }
}

/** 发送测试邮件 */
const sendTestEmail = async () => {
  if (!testEmail.value) {
    ElMessage.error("请输入测试邮箱地址")
    return
  }
  
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(testEmail.value)) {
    ElMessage.error("请输入有效的邮箱地址")
    return
  }
  
  try {
    testingEmail.value = true
    
    // 这里调用发送测试邮件的API
    // await sendTestEmailApi({ email: testEmail.value })
    
    // 模拟发送
    await new Promise(resolve => setTimeout(resolve, 2000))
    
    ElMessage.success("测试邮件发送成功，请检查邮箱")
  } catch (error: any) {
    console.error("发送测试邮件失败:", error)
    ElMessage.error("发送测试邮件失败")
  } finally {
    testingEmail.value = false
  }
}

/** 重置表单 */
const resetForm = async () => {
  try {
    await ElMessageBox.confirm(
      "确定要重置所有邮件设置吗？未保存的更改将丢失。",
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
    if (formData.email_enabled) {
      await formRef.value.validate()
    }
    emit("save", { ...formData })
  } catch (error) {
    ElMessage.error("请检查表单输入")
  }
}

/** 加载设置 */
const loadSettings = async () => {
  try {
    emit("update:loading", true)
    
    // 这里调用API加载邮件设置
    // const { data } = await getEmailSettingsApi()
    // Object.assign(formData, data)
    
    // 模拟数据
    Object.assign(formData, {
      admin_emails: ["admin@example.com"]
    })
    
  } catch (error: any) {
    console.error("加载邮件设置失败:", error)
    ElMessage.error("加载邮件设置失败")
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
.email-settings {
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
  
  .template-editor {
    margin-top: 16px;
  }
  
  .email-list-container {
    .email-item {
      display: flex;
      align-items: center;
      gap: 8px;
      margin-bottom: 8px;
      
      .el-input {
        flex: 1;
      }
    }
  }
  
  .test-email-section {
    margin-top: 32px;
    padding-top: 24px;
    border-top: 1px solid var(--el-border-color-lighter);
    
    h3 {
      margin: 0 0 16px 0;
      color: var(--el-text-color-primary);
      font-size: 16px;
      font-weight: 600;
    }
    
    .test-email-form {
      margin-bottom: 8px;
    }
    
    .form-help {
      font-size: 12px;
      color: var(--el-text-color-secondary);
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

.email-preview {
  .preview-header {
    padding: 12px;
    background: var(--el-fill-color-extra-light);
    border-radius: 4px;
    margin-bottom: 16px;
    font-size: 14px;
  }
  
  .preview-content {
    padding: 16px;
    border: 1px solid var(--el-border-color-light);
    border-radius: 4px;
    background: white;
    min-height: 200px;
    
    :deep(h2) {
      margin-top: 0;
      color: var(--el-text-color-primary);
    }
    
    :deep(a) {
      display: inline-block;
      margin: 8px 0;
    }
  }
}

:deep(.el-checkbox-group) {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>