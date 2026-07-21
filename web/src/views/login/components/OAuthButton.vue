<template>
  <el-button
    :loading="loading"
    :size="size"
    :class="buttonClass"
    @click="handleOAuthLogin"
  >
    <component :is="iconComponent" class="oauth-icon" v-if="!loading" />
    {{ buttonText }}
  </el-button>
</template>

<script lang="ts" setup>
import { ref, computed } from "vue"
import { useRouter, useRoute } from "vue-router"
import { useUserStore } from "@/store/modules/user"
import { ElMessage } from "element-plus"

interface Props {
  provider: string
  size?: 'large' | 'default' | 'small'
  block?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  size: 'large',
  block: false
})

const emit = defineEmits<{
  error: [message: string]
}>()

const router = useRouter()
const route = useRoute()
const loading = ref(false)

// OAuth provider configurations
const providerConfigs = {
  github: {
    name: 'GitHub',
    icon: 'github-icon',
    class: 'github-oauth-button'
  },
  google: {
    name: 'Google',
    icon: 'google-icon', 
    class: 'google-oauth-button'
  },
  microsoft: {
    name: 'Microsoft',
    icon: 'microsoft-icon',
    class: 'microsoft-oauth-button'
  }
}

const config = computed(() => providerConfigs[props.provider as keyof typeof providerConfigs])

const buttonText = computed(() => `使用 ${config.value?.name || props.provider} 登录`)

const buttonClass = computed(() => [
  config.value?.class || 'oauth-button',
  { 'w-100': props.block }
])

const iconComponent = computed(() => {
  switch (props.provider) {
    case 'github':
      return 'GitHubIcon'
    case 'google':
      return 'GoogleIcon'
    case 'microsoft':
      return 'MicrosoftIcon'
    default:
      return 'DefaultIcon'
  }
})

const handleOAuthLogin = async () => {
  try {
    loading.value = true
    
    // 保存当前页面的重定向信息（如果有的话）
    const redirect = route.query.redirect as string
    if (redirect) {
      localStorage.setItem('oauth_redirect', redirect)
    }
    
    const userStore = useUserStore()
    const { auth_url } = await userStore.getOAuthAuthUrl(props.provider)
    
    // 重定向到OAuth授权页面
    window.location.href = auth_url
  } catch (error: any) {
    console.error(`${props.provider} OAuth登录失败:`, error)
    const errorMessage = error.response?.data?.message || `${config.value?.name || props.provider}登录失败，请稍后重试`
    emit('error', errorMessage)
    loading.value = false
  }
}
</script>

<script lang="ts">
// Icon components
import { defineComponent, h } from 'vue'

const GitHubIcon = defineComponent({
  name: 'GitHubIcon',
  render() {
    return h('svg', {
      class: 'oauth-icon',
      viewBox: '0 0 24 24',
      width: '20',
      height: '20'
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
      class: 'oauth-icon',
      viewBox: '0 0 24 24',
      width: '20',
      height: '20'
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
      class: 'oauth-icon',
      viewBox: '0 0 24 24',
      width: '20',
      height: '20'
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
      class: 'oauth-icon default-icon'
    }, '🔐')
  }
})

export { GitHubIcon, GoogleIcon, MicrosoftIcon, DefaultIcon }
</script>

<style lang="scss" scoped>
.oauth-icon {
  flex-shrink: 0;
  margin-right: 8px;
}

.github-oauth-button {
  background-color: #24292e;
  border-color: #24292e;
  color: white;
  
  &:hover {
    background-color: #1b1f23;
    border-color: #1b1f23;
  }
  
  &:focus {
    background-color: #1b1f23;
    border-color: #1b1f23;
  }
}

.google-oauth-button {
  background-color: white;
  border-color: #dadce0;
  color: #3c4043;
  
  &:hover {
    background-color: #f8f9fa;
    border-color: #dadce0;
    color: #3c4043;
  }
  
  &:focus {
    background-color: #f8f9fa;
    border-color: #4285f4;
    color: #3c4043;
  }
}

.microsoft-oauth-button {
  background-color: #0078d4;
  border-color: #0078d4;
  color: white;
  
  &:hover {
    background-color: #106ebe;
    border-color: #106ebe;
  }
  
  &:focus {
    background-color: #106ebe;
    border-color: #106ebe;
  }
}

.oauth-button {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.default-icon {
  font-size: 16px;
}
</style>