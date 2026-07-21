<template>
  <div class="oauth-callback-container">
    <div class="callback-card">
      <div class="loading-content" v-if="loading">
        <el-icon class="loading-icon" :size="40">
          <Loading />
        </el-icon>
        <h3>正在处理登录...</h3>
        <p>请稍候，我们正在验证您的身份</p>
      </div>
      
      <div class="error-content" v-else-if="error">
        <el-icon class="error-icon" :size="40">
          <CircleClose />
        </el-icon>
        <h3>登录失败</h3>
        <p>{{ error }}</p>
        <el-button type="primary" @click="backToLogin">返回登录</el-button>
      </div>
      
      <div class="success-content" v-else-if="success">
        <el-icon class="success-icon" :size="40">
          <CircleCheck />
        </el-icon>
        <h3>登录成功</h3>
        <p>正在跳转到主页...</p>
        <div class="progress-bar">
          <div class="progress-fill"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from "vue"
import { useRouter, useRoute } from "vue-router"
import { useUserStore } from "@/store/modules/user"
import { Loading, CircleClose, CircleCheck } from "@element-plus/icons-vue"
import { ElMessage } from "element-plus"

const router = useRouter()
const route = useRoute()

const loading = ref(true)
const error = ref("")
const success = ref(false)

/** 处理OAuth回调 */
const handleOAuthCallback = async () => {
  try {
    const code = route.query.code as string
    const state = route.query.state as string
    const provider = (route.params.provider as string) || "github"
    const errorParam = route.query.error as string
    const errorDescription = route.query.error_description as string
    
    // 检查是否有错误参数
    if (errorParam) {
      const errorMsg = errorDescription || errorParam
      throw new Error(`OAuth授权失败: ${errorMsg}`)
    }
    
    // 检查必需参数
    if (!code) {
      throw new Error("缺少授权码参数")
    }
    
    if (!state) {
      throw new Error("缺少状态参数")
    }
    
    // 验证state参数
    const storedStateStr = localStorage.getItem('oauth_state')
    if (!storedStateStr) {
      throw new Error("未找到存储的状态参数，请重新登录")
    }
    
    let storedStateData
    try {
      storedStateData = JSON.parse(storedStateStr)
    } catch {
      throw new Error("状态参数格式错误，请重新登录")
    }
    
    // 检查state是否匹配
    if (state !== storedStateData.state) {
      throw new Error("无效的状态参数，可能存在安全风险")
    }
    
    // 检查provider是否匹配
    if (provider !== storedStateData.provider) {
      throw new Error("OAuth提供商不匹配")
    }
    
    // 检查是否过期（10分钟）
    const now = Date.now()
    const stateAge = now - storedStateData.timestamp
    if (stateAge > 10 * 60 * 1000) {
      throw new Error("授权已过期，请重新登录")
    }
    
    // 执行OAuth登录
    const userStore = useUserStore()
    await userStore.oauthLogin({ provider, code, state })
    
    success.value = true
    ElMessage.success("登录成功")
    
    // 获取重定向URL（如果有的话）
    const redirectUrl = localStorage.getItem('oauth_redirect') || "/board/dashboard"
    localStorage.removeItem('oauth_redirect')
    
    // 延迟跳转，让用户看到成功消息
    setTimeout(() => {
      router.replace({ path: redirectUrl })
    }, 1500)
    
  } catch (err: any) {
    console.error("OAuth登录失败:", err)
    error.value = err.message || "OAuth登录失败，请重试"
  } finally {
    loading.value = false
    // 清理存储的state
    localStorage.removeItem('oauth_state')
  }
}

/** 返回登录页面 */
const backToLogin = () => {
  router.replace({ path: "/login" })
}

/** 组件挂载时处理回调 */
onMounted(() => {
  // 添加超时机制，防止页面卡住
  const timeout = setTimeout(() => {
    if (loading.value) {
      error.value = "OAuth登录超时，请重试"
      loading.value = false
    }
  }, 30000) // 30秒超时

  handleOAuthCallback().finally(() => {
    clearTimeout(timeout)
  })
})
</script>

<style lang="scss" scoped>
.oauth-callback-container {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  
  .callback-card {
    width: 400px;
    max-width: 90%;
    padding: 40px;
    border-radius: 20px;
    background-color: var(--el-bg-color);
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
    text-align: center;
    
    .loading-content,
    .error-content,
    .success-content {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 16px;
      
      h3 {
        margin: 0;
        color: var(--el-text-color-primary);
        font-size: 24px;
        font-weight: 600;
      }
      
      p {
        margin: 0;
        color: var(--el-text-color-regular);
        font-size: 16px;
        line-height: 1.5;
      }
    }
    
    .loading-icon {
      color: var(--el-color-primary);
      animation: rotate 2s linear infinite;
    }
    
    .error-icon {
      color: var(--el-color-danger);
    }
    
    .success-icon {
      color: var(--el-color-success);
    }
    
    .progress-bar {
      width: 100%;
      height: 4px;
      background-color: var(--el-border-color-light);
      border-radius: 2px;
      overflow: hidden;
      margin-top: 8px;
      
      .progress-fill {
        height: 100%;
        background: linear-gradient(90deg, var(--el-color-primary), var(--el-color-success));
        border-radius: 2px;
        animation: progress 1.5s ease-in-out;
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

@keyframes progress {
  from {
    width: 0%;
  }
  to {
    width: 100%;
  }
}
</style>