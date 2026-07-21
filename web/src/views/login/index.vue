<script lang="ts" setup>
import { reactive, ref, onMounted } from "vue"
import { useRouter, useRoute } from "vue-router"
import { useUserStore } from "@/store/modules/user"
import { type FormInstance, type FormRules } from "element-plus"
import { User, Lock } from "@element-plus/icons-vue"
import { ElMessage } from "element-plus"
import { type LoginRequestData } from "@/api/login/types/login"
import ThemeSwitch from "@/components/ThemeSwitch/index.vue"
import Owl from "./components/Owl.vue"
import OAuthButton from "./components/OAuthButton.vue"
import { useFocus } from "./hooks/useFocus"

const router = useRouter()
const route = useRoute()
const { isFocus, handleBlur, handleFocus } = useFocus()

/** 登录表单元素的引用 */
const loginFormRef = ref<FormInstance | null>(null)

/** 登录按钮 Loading */
const loading = ref(false)

/** 错误消息 */
const errorMessage = ref("")

/** 登录表单数据 */
const loginFormData: LoginRequestData = reactive({
  username: "",
  password: ""
})

/** 登录表单校验规则 */
const loginFormRules: FormRules = {
  username: [{ required: true, message: "请输入用户名", trigger: "blur" }],
  password: [
    { required: true, message: "请输入密码", trigger: "blur" },
    { min: 6, message: "密码长度至少6个字符", trigger: "blur" }
  ]
}

/** 用户名密码登录逻辑 */
const handleLogin = () => {
  loginFormRef.value?.validate((valid: boolean, fields) => {
    if (valid) {
      loading.value = true
      errorMessage.value = ""
      
      useUserStore()
        .login(loginFormData)
        .then(() => {
          ElMessage.success("登录成功")
          router.push({ path: "/board/dashboard" })
        })
        .catch((error) => {
          console.error("登录失败:", error)
          errorMessage.value = error.response?.data?.message || "登录失败，请检查用户名和密码"
          loginFormData.password = ""
        })
        .finally(() => {
          loading.value = false
        })
    } else {
      console.error("表单校验不通过", fields)
    }
  })
}

/** OAuth登录错误处理 */
const handleOAuthError = (message: string) => {
  errorMessage.value = message
}


</script>

<template>
  <div class="login-container">
    <ThemeSwitch class="theme-switch" />
    <Owl :close-eyes="isFocus" />
    <div class="login-card">
      <div class="title">
        <img src="@/assets/layouts/logo-text-1.png" />
      </div>
      <div class="content">
        <!-- 错误消息显示 -->
        <el-alert
          v-if="errorMessage"
          :title="errorMessage"
          type="error"
          :closable="false"
          show-icon
          class="error-alert"
        />
        
        <!-- 用户名密码登录表单 -->
        <el-form ref="loginFormRef" :model="loginFormData" :rules="loginFormRules" @keyup.enter="handleLogin">
          <el-form-item prop="username">
            <el-input
              v-model.trim="loginFormData.username"
              placeholder="用户名"
              type="text"
              tabindex="1"
              :prefix-icon="User"
              size="large"
            />
          </el-form-item>
          <el-form-item prop="password">
            <el-input
              v-model.trim="loginFormData.password"
              placeholder="密码"
              type="password"
              tabindex="2"
              :prefix-icon="Lock"
              size="large"
              show-password
              @blur="handleBlur"
              @focus="handleFocus"
            />
          </el-form-item>
          <el-button 
            :loading="loading" 
            type="primary" 
            size="large" 
            @click.prevent="handleLogin"
            class="login-button"
          >
            登 录
          </el-button>
        </el-form>

        <!-- 分割线 -->
        <div class="divider">
          <span>或</span>
        </div>

        <!-- OAuth登录按钮 -->
        <div class="oauth-login">
          <OAuthButton 
            provider="github" 
            :block="true"
            @error="handleOAuthError"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.login-container {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  width: 100%;
  min-height: 100%;
  
  .theme-switch {
    position: fixed;
    top: 5%;
    right: 5%;
    cursor: pointer;
  }
  
  .login-card {
    width: 480px;
    max-width: 90%;
    border-radius: 20px;
    box-shadow: 0 0 50px #054cf1;
    background-color: var(--el-bg-color);
    overflow: hidden;
    
    .title {
      display: flex;
      justify-content: center;
      vertical-align: middle;
      align-items: center;
      height: 120px;
      
      img {
        height: 100%;
      }
    }
    
    .content {
      padding: 10px 50px 50px 50px;
      
      .error-alert {
        margin-bottom: 20px;
      }
      
      .login-button {
        width: 100%;
        margin-top: 10px;
      }
      
      .divider {
        display: flex;
        align-items: center;
        margin: 30px 0 20px 0;
        
        &::before,
        &::after {
          content: '';
          flex: 1;
          height: 1px;
          background: var(--el-border-color-light);
        }
        
        span {
          padding: 0 15px;
          color: var(--el-text-color-regular);
          font-size: 14px;
        }
      }
      
      .oauth-login {
        // OAuth button styles are handled by the OAuthButton component
      }
    }
  }
}
</style>
