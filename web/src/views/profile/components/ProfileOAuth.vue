<template>
  <div class="profile-oauth">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>关联账户</span>
          <el-button 
            size="small" 
            @click="refreshOAuthAccounts"
            :loading="loading"
          >
            刷新
          </el-button>
        </div>
      </template>

      <div class="oauth-accounts">
        <div v-if="loading" class="loading-container">
          <el-icon class="loading-icon"><Loading /></el-icon>
          <span>加载中...</span>
        </div>
        
        <div v-else>
          <div 
            v-for="provider in oauthProviders" 
            :key="provider.name"
            class="oauth-provider"
          >
            <div class="provider-info">
              <div class="provider-header">
                <div class="provider-icon">
                  <component :is="getProviderIcon(provider.name)" />
                </div>
                <div class="provider-details">
                  <h4>{{ provider.display_name }}</h4>
                  <p>{{ provider.description }}</p>
                </div>
              </div>
              
              <div v-if="provider.connected" class="connected-info">
                <div class="connected-account">
                  <el-avatar 
                    :size="32" 
                    :src="provider.account_info?.avatar_url"
                  >
                    {{ provider.account_info?.username?.charAt(0) }}
                  </el-avatar>
                  <div class="account-details">
                    <span class="account-name">
                      {{ provider.account_info?.display_name || provider.account_info?.username }}
                    </span>
                    <span class="account-email">{{ provider.account_info?.email }}</span>
                  </div>
                </div>
                
                <div class="connection-meta">
                  <p><strong>关联时间:</strong> {{ formatDate(provider.connected_at) }}</p>
                  <p><strong>最后同步:</strong> {{ formatDate(provider.last_sync) }}</p>
                </div>
              </div>
            </div>
            
            <div class="provider-actions">
              <el-button
                v-if="!provider.connected"
                type="primary"
                @click="connectOAuth(provider.name)"
                :loading="connectLoading[provider.name]"
              >
                关联账户
              </el-button>
              
              <div v-else class="connected-actions">
                <el-button
                  size="small"
                  @click="syncOAuthAccount(provider.name)"
                  :loading="syncLoading[provider.name]"
                >
                  同步信息
                </el-button>
                <el-button
                  type="danger"
                  size="small"
                  @click="disconnectOAuth(provider.name)"
                  :loading="disconnectLoading[provider.name]"
                >
                  解除关联
                </el-button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </el-card>

    <!-- OAuth权限说明 -->
    <el-card class="permissions-card">
      <template #header>
        <span>权限说明</span>
      </template>
      
      <div class="permissions-info">
        <el-alert
          title="关联账户权限"
          type="info"
          :closable="false"
          show-icon
        >
          <template #default>
            <p>关联第三方账户时，我们会请求以下权限：</p>
            <ul>
              <li><strong>基本信息:</strong> 获取您的用户名、头像和邮箱地址</li>
              <li><strong>身份验证:</strong> 用于快速登录和身份验证</li>
              <li><strong>信息同步:</strong> 可选择同步头像和显示名称</li>
            </ul>
            <p class="note">
              <el-icon><InfoFilled /></el-icon>
              我们不会获取您的密码或其他敏感信息，您可以随时解除关联。
            </p>
          </template>
        </el-alert>
      </div>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted } from "vue"
import { ElMessage, ElMessageBox } from "element-plus"
import { Loading, InfoFilled } from "@element-plus/icons-vue"
import { 
  getOAuthAccountsApi, 
  connectOAuthApi, 
  disconnectOAuthApi,
  syncOAuthAccountApi 
} from "@/api/profile"
import { useUserStore } from "@/store/modules/user"

interface UserInfo {
  id: number
  username: string
  email: string
  display_name: string
  avatar_url: string
  role: string
  is_active: boolean
  email_verified: boolean
  last_login: string | null
  created_at: string
}

interface Props {
  userInfo: UserInfo
}

const props = defineProps<Props>()
const emit = defineEmits<{
  update: [userInfo: Partial<UserInfo>]
  loading: [loading: boolean]
}>()

const userStore = useUserStore()
const loading = ref(false)
const connectLoading = reactive<Record<string, boolean>>({})
const disconnectLoading = reactive<Record<string, boolean>>({})
const syncLoading = reactive<Record<string, boolean>>({})

const oauthProviders = ref<any[]>([])

// OAuth提供商配置
const providerConfigs = {
  github: {
    display_name: "GitHub",
    description: "使用GitHub账户快速登录",
    icon: "GitHubIcon"
  },
  google: {
    display_name: "Google",
    description: "使用Google账户快速登录",
    icon: "GoogleIcon"
  },
  microsoft: {
    display_name: "Microsoft",
    description: "使用Microsoft账户快速登录",
    icon: "MicrosoftIcon"
  }
}

/** 获取提供商图标组件 */
const getProviderIcon = (provider: string) => {
  const config = providerConfigs[provider as keyof typeof providerConfigs]
  return config?.icon || "DefaultIcon"
}

/** 加载OAuth账户信息 */
const loadOAuthAccounts = async () => {
  try {
    loading.value = true
    const { data } = await getOAuthAccountsApi()
    
    // 合并配置信息和账户状态
    oauthProviders.value = Object.keys(providerConfigs).map(providerName => {
      const config = providerConfigs[providerName as keyof typeof providerConfigs]
      const accountData = data.find((account: any) => account.provider === providerName)
      
      return {
        name: providerName,
        ...config,
        connected: !!accountData,
        account_info: accountData?.account_info,
        connected_at: accountData?.connected_at,
        last_sync: accountData?.last_sync
      }
    })
    
  } catch (error: any) {
    console.error("加载OAuth账户失败:", error)
    ElMessage.error("加载OAuth账户失败")
  } finally {
    loading.value = false
  }
}

/** 刷新OAuth账户 */
const refreshOAuthAccounts = () => {
  loadOAuthAccounts()
}

/** 关联OAuth账户 */
const connectOAuth = async (provider: string) => {
  try {
    connectLoading[provider] = true
    
    // 保存当前页面信息，OAuth完成后返回
    localStorage.setItem('oauth_redirect', '/profile')
    
    // 获取OAuth授权URL并跳转
    const { auth_url } = await userStore.getOAuthAuthUrl(provider)
    window.location.href = auth_url
    
  } catch (error: any) {
    console.error(`关联${provider}失败:`, error)
    ElMessage.error(`关联${provider}失败`)
    connectLoading[provider] = false
  }
}

/** 解除OAuth关联 */
const disconnectOAuth = async (provider: string) => {
  try {
    const config = providerConfigs[provider as keyof typeof providerConfigs]
    
    await ElMessageBox.confirm(
      `确定要解除与${config.display_name}的关联吗？解除后将无法使用该账户登录。`,
      "确认解除关联",
      {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning"
      }
    )
    
    disconnectLoading[provider] = true
    emit("loading", true)
    
    await disconnectOAuthApi(provider)
    
    ElMessage.success(`已解除与${config.display_name}的关联`)
    await loadOAuthAccounts()
    
  } catch (error: any) {
    if (error !== "cancel") {
      console.error(`解除${provider}关联失败:`, error)
      ElMessage.error(`解除${provider}关联失败`)
    }
  } finally {
    disconnectLoading[provider] = false
    emit("loading", false)
  }
}

/** 同步OAuth账户信息 */
const syncOAuthAccount = async (provider: string) => {
  try {
    const config = providerConfigs[provider as keyof typeof providerConfigs]
    
    syncLoading[provider] = true
    emit("loading", true)
    
    const { data } = await syncOAuthAccountApi(provider)
    
    ElMessage.success(`${config.display_name}信息同步成功`)
    
    // 如果同步了用户信息，更新本地数据
    if (data.user_updated) {
      emit("update", data.user_info)
    }
    
    await loadOAuthAccounts()
    
  } catch (error: any) {
    console.error(`同步${provider}信息失败:`, error)
    ElMessage.error(`同步${provider}信息失败`)
  } finally {
    syncLoading[provider] = false
    emit("loading", false)
  }
}

/** 格式化日期 */
const formatDate = (dateString: string) => {
  if (!dateString) return "未知"
  return new Date(dateString).toLocaleString("zh-CN")
}

/** 组件挂载时加载数据 */
onMounted(() => {
  // TODO: 暂时禁用OAuth账户加载，等后端实现后再启用
  // loadOAuthAccounts()
})
</script>

<script lang="ts">
// Icon components (reuse from OAuthButton component)
import { defineComponent, h } from 'vue'

const GitHubIcon = defineComponent({
  name: 'GitHubIcon',
  render() {
    return h('svg', {
      viewBox: '0 0 24 24',
      width: '24',
      height: '24'
    }, [
      h('path', {
        fill: 'currentColor',
        d: 'M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z'
      })
    ])
  }
})

const GoogleIcon = defineComponent({
  name: 'GoogleIcon',
  render() {
    return h('svg', {
      viewBox: '0 0 24 24',
      width: '24',
      height: '24'
    }, [
      h('path', {
        fill: '#4285F4',
        d: 'M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z'
      }),
      h('path', {
        fill: '#34A853',
        d: 'M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z'
      }),
      h('path', {
        fill: '#FBBC05',
        d: 'M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z'
      }),
      h('path', {
        fill: '#EA4335',
        d: 'M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z'
      })
    ])
  }
})

const MicrosoftIcon = defineComponent({
  name: 'MicrosoftIcon',
  render() {
    return h('svg', {
      viewBox: '0 0 24 24',
      width: '24',
      height: '24'
    }, [
      h('path', {
        fill: '#f25022',
        d: 'M1 1h10v10H1z'
      }),
      h('path', {
        fill: '#00a4ef',
        d: 'M13 1h10v10H13z'
      }),
      h('path', {
        fill: '#7fba00',
        d: 'M1 13h10v10H1z'
      }),
      h('path', {
        fill: '#ffb900',
        d: 'M13 13h10v10H13z'
      })
    ])
  }
})

const DefaultIcon = defineComponent({
  name: 'DefaultIcon',
  render() {
    return h('div', {
      style: {
        width: '24px',
        height: '24px',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        fontSize: '16px'
      }
    }, '🔐')
  }
})

export { GitHubIcon, GoogleIcon, MicrosoftIcon, DefaultIcon }
</script>

<style lang="scss" scoped>
.profile-oauth {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    
    span {
      font-weight: 600;
      font-size: 16px;
    }
  }
  
  .oauth-accounts {
    .loading-container {
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 40px;
      color: var(--el-text-color-secondary);
      
      .loading-icon {
        margin-right: 8px;
        animation: rotate 2s linear infinite;
      }
    }
  }
  
  .oauth-provider {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    padding: 20px;
    border: 1px solid var(--el-border-color-light);
    border-radius: 8px;
    margin-bottom: 16px;
    transition: all 0.3s ease;
    
    &:hover {
      border-color: var(--el-color-primary-light-7);
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    }
    
    .provider-info {
      flex: 1;
      
      .provider-header {
        display: flex;
        align-items: center;
        margin-bottom: 12px;
        
        .provider-icon {
          width: 40px;
          height: 40px;
          display: flex;
          align-items: center;
          justify-content: center;
          background: var(--el-fill-color-light);
          border-radius: 8px;
          margin-right: 12px;
          color: var(--el-text-color-primary);
        }
        
        .provider-details {
          h4 {
            margin: 0 0 4px 0;
            font-size: 16px;
            font-weight: 600;
            color: var(--el-text-color-primary);
          }
          
          p {
            margin: 0;
            font-size: 14px;
            color: var(--el-text-color-regular);
          }
        }
      }
      
      .connected-info {
        .connected-account {
          display: flex;
          align-items: center;
          margin-bottom: 12px;
          
          .account-details {
            margin-left: 12px;
            
            .account-name {
              display: block;
              font-weight: 500;
              color: var(--el-text-color-primary);
            }
            
            .account-email {
              display: block;
              font-size: 12px;
              color: var(--el-text-color-secondary);
            }
          }
        }
        
        .connection-meta {
          p {
            margin: 4px 0;
            font-size: 12px;
            color: var(--el-text-color-regular);
            
            strong {
              color: var(--el-text-color-primary);
            }
          }
        }
      }
    }
    
    .provider-actions {
      margin-left: 20px;
      
      .connected-actions {
        display: flex;
        flex-direction: column;
        gap: 8px;
      }
    }
  }
  
  .permissions-card {
    margin-top: 24px;
    
    .permissions-info {
      :deep(.el-alert__content) {
        ul {
          margin: 12px 0;
          padding-left: 20px;
          
          li {
            margin: 8px 0;
            line-height: 1.5;
          }
        }
        
        .note {
          display: flex;
          align-items: center;
          margin-top: 16px;
          padding: 12px;
          background: var(--el-fill-color-extra-light);
          border-radius: 6px;
          font-size: 14px;
          
          .el-icon {
            margin-right: 8px;
            color: var(--el-color-primary);
          }
        }
      }
    }
  }
}

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

@media (max-width: 768px) {
  .oauth-provider {
    flex-direction: column;
    
    .provider-actions {
      margin-left: 0;
      margin-top: 16px;
      width: 100%;
      
      .connected-actions {
        flex-direction: row;
        justify-content: flex-start;
      }
    }
  }
}
</style>