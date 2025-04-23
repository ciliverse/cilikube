<template>
  <div class="tech-navigation-dashboard-v2">
    <!-- 背景动画层 -->
    <div class="background-animation"></div>

    <!-- 顶部信息/广告栏 -->
    <div v-if="showBanner" class="top-info-banner">
      <div class="banner-content">
        <el-icon><Bell /></el-icon>
        <span>特别通知：CILIKUBE v0.1.0 版本现已发布，带来性能优化与新功能！ <a href="https://cilikube.cillian.website" target="_blank">查看详情</a></span>
      </div>
      <el-icon class="close-banner-icon" @click="hideBanner"><Close /></el-icon>
    </div>

    <!-- 右上角操作按钮 -->
    <div class="top-actions">
      <el-button v-if="!isLoggedIn" type="primary" round :icon="User" @click="handleLogin">
        登录
      </el-button>
      <el-button v-if="isLoggedIn" type="danger" round :icon="SwitchButton" @click="handleLogout">
        退出登录
      </el-button>
    </div>

    <!-- 主要内容区域 -->
    <div class="main-content">
      <!-- Logo 和 标题区域 -->
      <div class="header-section">
         <el-image
           style="width: 100px; height: 100px; margin-right: 18px; box-shadow: 0 0 15px rgba(0, 191, 255, 0.5);"
           src="favicon.ico"
           fit="cover"
           lazy>
           <template #error>
             <div class="image-slot">
               <el-icon><Picture /></el-icon>
             </div>
           </template>
         </el-image>
         <div class="title-container">
           <h1 class="dashboard-title">CILIKUBE</h1>
           <p class="dashboard-subtitle">Kubernetes初学和开发入门的优质开源项目</p>
         </div>
       </div>

      <!-- 卡片网格 -->
      <el-row :gutter="30" class="card-row">
        <!-- 卡片: 集群安装 -->
        <el-col :xs="24" :sm="12" :md="8" :lg="6" class="card-col">
          <div class="card-wrapper" @click="navigateTo('minikube')">
            <el-card shadow="never" class="feature-card">
              <div class="card-glow"></div> <!-- 发光层 -->
              <div class="card-content">
                <div class="card-icon-wrapper bg-cyan">
                  <el-icon :size="30"><Download /></el-icon>
                </div>
                <div class="card-text">
                  <h3>集群安装</h3>
                  <p>本地快速启动 Minikube 测试环境。</p>
                </div>
              </div>
            </el-card>
          </div>
        </el-col>

        <!-- 卡片: 集群管理 -->
        <el-col :xs="24" :sm="12" :md="8" :lg="6" class="card-col">
           <div class="card-wrapper" @click="navigateTo('Dashboard')">
             <el-card shadow="never" class="feature-card">
              <div class="card-glow"></div>
              <div class="card-content">
                <div class="card-icon-wrapper bg-blue">
                  <el-icon :size="30"><DataAnalysis /></el-icon>
                </div>
                <div class="card-text">
                  <h3>集群管理</h3>
                  <p>概览集群状态、资源监控与管理。</p>
                </div>
              </div>
            </el-card>
          </div>
        </el-col>

        <!-- 卡片: 工作负载 -->
         <el-col :xs="24" :sm="12" :md="8" :lg="6" class="card-col">
          <div class="card-wrapper" @click="navigateTo('pods')">
            <el-card shadow="never" class="feature-card">
              <div class="card-glow"></div>
              <div class="card-content">
                <div class="card-icon-wrapper bg-purple">
                   <el-icon :size="30"><Box /></el-icon>
                </div>
                <div class="card-text">
                  <h3>工作负载</h3>
                  <p>管理 Pods, Deployments 等核心应用。</p>
                </div>
              </div>
            </el-card>
          </div>
        </el-col>

        <!-- 卡片: 网络配置 -->
         <el-col :xs="24" :sm="12" :md="8" :lg="6" class="card-col">
          <div class="card-wrapper" @click="navigateTo('service')">
            <el-card shadow="never" class="feature-card">
               <div class="card-glow"></div>
              <div class="card-content">
                <div class="card-icon-wrapper bg-green">
                   <el-icon :size="30"><Connection /></el-icon>
                </div>
                <div class="card-text">
                  <h3>网络配置</h3>
                  <p>配置 Service 发现与 Ingress 规则。</p>
                </div>
              </div>
            </el-card>
          </div>
        </el-col>

        <!-- 卡片: 存储管理 -->
         <el-col :xs="24" :sm="12" :md="8" :lg="6" class="card-col">
          <div class="card-wrapper" @click="navigateTo('Notice')">
            <el-card shadow="never" class="feature-card">
              <div class="card-glow"></div>
              <div class="card-content">
                <div class="card-icon-wrapper bg-orange">
                   <el-icon :size="30"><Coin /></el-icon>
                </div>
                <div class="card-text">
                  <h3>存储管理</h3>
                  <p>管理 PV 与 PVC 持久化数据卷。</p>
                </div>
              </div>
            </el-card>
          </div>
        </el-col>

         <!-- 卡片: 配置管理 -->
         <el-col :xs="24" :sm="12" :md="8" :lg="6" class="card-col">
          <div class="card-wrapper" @click="navigateTo('ElementPlus')">
            <el-card shadow="never" class="feature-card">
              <div class="card-glow"></div>
              <div class="card-content">
                <div class="card-icon-wrapper bg-red">
                   <el-icon :size="30"><Tickets /></el-icon>
                </div>
                <div class="card-text">
                  <h3>配置管理</h3>
                  <p>管理 ConfigMaps 与 Secrets 配置数据。</p>
                </div>
              </div>
            </el-card>
          </div>
        </el-col>

         <!-- 卡片: 功能待开发 -->
        <el-col :xs="24" :sm="12" :md="8" :lg="6" class="card-col">
           <div class="card-wrapper coming-soon" @click="notifyComingSoon">
            <el-card shadow="never" class="feature-card">
               <div class="card-glow"></div>
              <div class="card-content">
                <div class="card-icon-wrapper bg-grey">
                  <el-icon :size="30"><Monitor /></el-icon>
                </div>
                <div class="card-text">
                  <h3>日志与监控</h3>
                  <p>集中查看应用日志与高级监控。</p>
                </div>
              </div>
              <div class="coming-soon-overlay">
                <span>即将推出</span>
              </div>
            </el-card>
          </div>
        </el-col>
         <!-- 可以继续添加更多卡片 -->
      </el-row>
    </div> <!-- main-content 结束 -->

    <!-- 底部版权信息 -->
    <footer class="dashboard-footer">
       <span>Copyright © {{ new Date().getFullYear() }} CILIVERSE . All Rights Reserved.</span>
       <!-- 添加触发弹窗的链接 -->
       <el-link type="primary" class="footer-link" @click="dialogVisible = true">加入社群 & 支持我们</el-link>
    </footer>

     <!-- 社群与捐赠弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      title="加入社群 & 支持我们"
      width="80%"
      :close-on-click-modal="true"
      draggable
      :append-to-body="true"
       center
       custom-class="community-dialog"
    >
      <div class="community-dialog-content">
        <el-row :gutter="30">
          <!-- 微信交流群 -->
          <el-col :xs="24" :sm="8" class="qr-code-col">
             <el-image :src="wechatGroupQR" fit="contain" class="qr-code-image" lazy>
                <template #placeholder><div class="image-slot">加载中...</div></template>
                <template #error><div class="image-slot">图片加载失败</div></template>
             </el-image>
             <p class="qr-code-description">扫码加入微信交流群</p>
          </el-col>

           <!-- 个人微信号 -->
          <el-col :xs="24" :sm="8" class="qr-code-col">
             <el-image :src="personalWechatQR" fit="contain" class="qr-code-image" lazy>
                <template #placeholder><div class="image-slot">加载中...</div></template>
                <template #error><div class="image-slot">图片加载失败</div></template>
             </el-image>
             <p class="qr-code-description">添加开发者微信</p>
          </el-col>

           <!-- 捐赠/支持 -->
          <el-col :xs="24" :sm="8" class="qr-code-col">
             <el-image :src="donationQR" fit="contain" class="qr-code-image" lazy>
               <template #placeholder><div class="image-slot">加载中...</div></template>
                <template #error><div class="image-slot">图片加载失败</div></template>
             </el-image>
             <p class="qr-code-description">您的支持是我们前进的动力！</p>
          </el-col>
        </el-row>
        <p class="dialog-tip">感谢您对 CILIKUBE 项目的关注与支持！</p>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>

  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';
import { useRouter } from 'vue-router';
// 导入 Element Plus 组件
import {
    ElCard, ElRow, ElCol, ElIcon, ElNotification, ElImage, ElButton,
    ElMessageBox, ElMessage, ElDialog, ElLink
} from 'element-plus';
import { useUserStore } from '@/store/modules/user';
// 导入图标
import {
  Platform, Download, Setting, Cpu, Files, Grid, Picture, User, SwitchButton, Bell, Close,
  DataAnalysis, Box, Connection, Coin, Tickets, Monitor
} from '@element-plus/icons-vue';

// --- 导入图片资源 (请确保路径正确！) ---
import wechatGroupQR from '/logo.png'; // 替换为你的实际图片路径
import personalWechatQR from '/logo.png'; // 替换为你的实际图片路径
import donationQR from '/logo.png'; // 替换为你的实际图片路径

const router = useRouter();
const userStore = useUserStore();

// 控制顶部信息栏显示
const showBanner = ref(true);
const hideBanner = () => {
  showBanner.value = false;
};

// --- 控制社群/捐赠弹窗显示 ---
const dialogVisible = ref(false);

// --- Computed Property for Login Status ---
const isLoggedIn = computed(() => !!userStore.token);

// --- Navigation Function ---
const navigateTo = (routeName: string) => {
  if (!routeName) {
    console.warn('导航尝试使用了空路由名称。');
    return;
  }
  try {
    if (router.hasRoute(routeName)) {
         router.push({ name: routeName });
    } else {
        console.error(`路由跳转失败：名为 "${routeName}" 的路由未在路由配置中定义。`);
        ElNotification({
          title: '导航错误',
          message: `目标路由 "${routeName}" 未找到，请检查 src/router/index.ts 或相关路由配置文件。`,
          type: 'error',
          zIndex: 3000,
        });
    }
  } catch (error) {
      console.error(`导航到路由 "${routeName}" 时发生意外错误:`, error);
       ElNotification({
          title: '导航错误',
          message: '页面跳转时发生未知错误。',
          type: 'error',
          zIndex: 3000,
        });
  }
};

// --- Handler for Coming Soon Cards ---
const notifyComingSoon = () => {
  ElNotification({
    title: '敬请期待',
    message: '该功能正在紧张开发中...',
    type: 'info',
    duration: 2500,
    zIndex: 3000,
  });
};

// --- Login Handler ---
const handleLogin = () => {
  const loginRouteName = 'Login';
  if (router.hasRoute(loginRouteName)) {
    router.push({ name: loginRouteName });
  } else {
     console.error(`登录失败：名为 "${loginRouteName}" 的路由未定义！请检查路由配置。`);
     ElMessage.error(`无法导航到登录页面，请联系管理员或检查路由配置 (需要名为 '${loginRouteName}' 的路由)。`);
  }
};

// --- Logout Handler ---
const handleLogout = () => {
  ElMessageBox.confirm(
    '您确定要退出登录吗？',
    '退出确认',
    {
      confirmButtonText: '确定退出',
      cancelButtonText: '取消',
      type: 'warning',
      draggable: true,
    }
  )
  .then(() => {
    userStore.logout();
    ElMessage({
      type: 'success',
      message: '已成功退出登录',
    });
    // Optional redirect logic can be placed inside userStore.logout() or here
  })
  .catch(() => {
    // User canceled the logout
  });
};

</script>

<style scoped>
/* --- 动画 --- */
@keyframes gradientBG {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

@keyframes fadeInSlideUp {
  from { opacity: 0; transform: translateY(30px); }
  to { opacity: 1; transform: translateY(0); }
}

/* --- 主容器 Flex 布局 --- */
.tech-navigation-dashboard-v2 {
  position: relative;
  padding: 0 60px; /* 移除顶部和底部 padding，由 flex 控制 */
  min-height: 100vh;
  background: linear-gradient(-45deg, #0a101f, #11182f, #1a233f, #2a3b5f);
  background-size: 400% 400%;
  animation: gradientBG 25s ease infinite;
  color: #e0e6f0;
  box-sizing: border-box;
  overflow-x: hidden;
  /* 添加 Flex 布局 */
  display: flex;
  flex-direction: column;
}

/* --- 主要内容区域 --- */
.main-content {
  flex-grow: 1; /* 占据所有可用垂直空间，将 footer 推到底部 */
  /* 为顶部的绝对定位元素留出空间 */
  padding-top: 80px; /* 大约等于 banner 高度 + actions 区域高度 + 一些间距 */
  padding-bottom: 40px; /* 内容区域底部的间距 */
  position: relative; /* 使得内部的 z-index 生效 */
  z-index: 1; /* 确保在背景之上 */
}

/* --- 背景网格 --- */
.background-animation::before {
    content: "";
    position: absolute; top: 0; left: 0; right: 0; bottom: 0;
    background-image: linear-gradient(rgba(255,255,255,0.02) 1px, transparent 1px),
                      linear-gradient(90deg, rgba(255,255,255,0.02) 1px, transparent 1px);
    background-size: 50px 50px;
    opacity: 0.3;
    z-index: 0;
}

/* --- 顶部信息栏 --- */
.top-info-banner {
  position: fixed; /* 固定在顶部 */
  top: 0; left: 0; right: 0;
  background-color: rgba(17, 24, 47, 0.85); /* 稍微不透明一点 */
  backdrop-filter: blur(8px);
  color: #a0aec0;
  padding: 8px 60px; /* 左右内边距与主容器一致 */
  font-size: 0.85rem;
  z-index: 15; /* 比 actions 高 */
  display: flex;
  align-items: center;
  justify-content: space-between; /* 让关闭按钮自然到右边 */
  box-shadow: 0 2px 5px rgba(0,0,0,0.2);
}
.banner-content {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%; /* 占满宽度 */
  margin-bottom: 5px; /* 底部间距 */
  font-size: 0.9rem; /* 调整字体大小 */
  color: #e0e6f0; /* 内容颜色 */
  font-weight: 400; /* 内容字体粗细 */
  line-height: 1.7; /* 行高 */
}
.top-info-banner a { color: #60a5fa; text-decoration: none; font-weight: 500; }
.top-info-banner a:hover { text-decoration: underline; }
.close-banner-icon {
  cursor: pointer;
  font-size: 1.1rem;
  color: #9ca3af; /* 关闭图标颜色 */
  transition: color 0.2s ease;
  margin-left: 15px; /* 与内容保持距离 */
}
.close-banner-icon:hover {
  color: #e0e6f0; /* 悬停时变亮 */
}


/* --- 右上角按钮 --- */
.top-actions {
  position: absolute;
  /* 顶部位置需要考虑 banner 的高度 */
  top: 50px; /* 调整 top 值，使其在 banner 下方 */
  right: 40px;
  z-index: 10;
}
.top-actions .el-button { box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2); }
.top-actions .el-button--primary { background-color: #409eff; border-color: #409eff; }
.top-actions .el-button--danger { background-color: #f56c6c; border-color: #f56c6c; }

/* --- Header --- */
.header-section {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 60px;
  text-align: left;
  position: relative;
  z-index: 1;
}
.title-container { display: flex; flex-direction: column; }
.dashboard-title { font-size: 2.8rem; font-weight: 600; color: #ffffff; margin: 0 0 8px 0; letter-spacing: 1.5px; text-shadow: 0 0 10px rgba(255, 255, 255, 0.3); }
.dashboard-subtitle { font-size: 1.2rem; color: #a0aec0; margin: 0; font-weight: 300; }
.image-slot { display: flex; justify-content: center; align-items: center; width: 100%; height: 100%; background: #1f2937; color: #6b7280; font-size: 26px; border-radius: 50%; }

/* --- 卡片网格和动画 --- */
.card-row { justify-content: center; }
.card-col { opacity: 0; animation: fadeInSlideUp 0.6s ease-out forwards; }
.card-col:nth-child(1) { animation-delay: 0.1s; }
.card-col:nth-child(2) { animation-delay: 0.2s; }
.card-col:nth-child(3) { animation-delay: 0.3s; }
.card-col:nth-child(4) { animation-delay: 0.4s; }
.card-col:nth-child(5) { animation-delay: 0.5s; }
.card-col:nth-child(6) { animation-delay: 0.6s; }
.card-col:nth-child(7) { animation-delay: 0.7s; }

/* --- 卡片样式 --- */
.card-wrapper { margin-bottom: 30px; transition: transform 0.35s cubic-bezier(0.25, 0.8, 0.25, 1), box-shadow 0.35s cubic-bezier(0.25, 0.8, 0.25, 1); border-radius: 14px; position: relative; cursor: pointer; }
.feature-card { border-radius: 14px; border: 1px solid rgba(129, 140, 160, 0.2); background: rgba(27, 38, 63, 0.6); backdrop-filter: blur(12px) saturate(150%); overflow: hidden; height: 100%; display: flex; flex-direction: column; position: relative; transition: border-color 0.35s ease, background-color 0.35s ease; }
.card-wrapper:hover { transform: translateY(-10px); }
.card-glow { content: ''; position: absolute; top: -2px; left: -2px; right: -2px; bottom: -2px; border-radius: 16px; background: radial-gradient(ellipse at center, rgba(56, 189, 248, 0.2), transparent 70%); opacity: 0; transition: opacity 0.4s ease-out; z-index: -1; pointer-events: none; }
.card-wrapper:hover .card-glow { opacity: 1; }
.card-wrapper:hover .feature-card { border-color: rgba(96, 165, 250, 0.5); background: rgba(36, 48, 75, 0.7); }
.card-content { display: flex; align-items: center; padding: 25px 30px; gap: 22px; flex-grow: 1; position: relative; z-index: 1; }
.card-icon-wrapper { flex-shrink: 0; width: 58px; height: 58px; border-radius: 50%; display: flex; align-items: center; justify-content: center; box-shadow: 0 5px 15px rgba(0, 0, 0, 0.2); transition: transform 0.3s ease; }
.card-wrapper:hover .card-icon-wrapper { transform: scale(1.1); }
.bg-cyan { background: linear-gradient(135deg, #2dd4bf, #14b8a6); }
.bg-blue { background: linear-gradient(135deg, #60a5fa, #2563eb); }
.bg-purple { background: linear-gradient(135deg, #c084fc, #9333ea); }
.bg-green { background: linear-gradient(135deg, #4ade80, #16a34a); }
.bg-orange { background: linear-gradient(135deg, #fbbf24, #d97706); }
.bg-red { background: linear-gradient(135deg, #f87171, #dc2626); }
.bg-grey { background: linear-gradient(135deg, #9ca3af, #6b7280); }
.card-text h3 { margin: 0 0 10px 0; font-size: 1.25rem; color: #f0f8ff; font-weight: 500; }
.card-text p { margin: 0; font-size: 0.9rem; color: #b0c4de; line-height: 1.65; font-weight: 300; }

/* --- Coming Soon --- */
.card-wrapper.coming-soon { cursor: not-allowed; }
.card-wrapper.coming-soon:hover { transform: none; }
.card-wrapper.coming-soon .card-glow { opacity: 0; }
.card-wrapper.coming-soon .feature-card { border-color: rgba(129, 140, 160, 0.2); background: rgba(27, 38, 63, 0.6); }
.card-wrapper.coming-soon:hover .card-icon-wrapper { transform: none; }
.coming-soon-overlay { position: absolute; top: 0; left: 0; right: 0; bottom: 0; background-color: rgba(11, 15, 28, 0.1); display: flex; align-items: center; justify-content: center; font-weight: 500; color: #ffffff; font-size: 1.1rem; backdrop-filter: blur(4px); opacity: 0; transition: opacity 0.3s ease; pointer-events: none; border-radius: 14px; z-index: 2; }
.card-wrapper.coming-soon .coming-soon-overlay { opacity: 1; pointer-events: auto; }

/* --- 底部版权信息与链接 --- */
.dashboard-footer {
  text-align: center;
  padding: 20px 0;
  color: #6b7280;
  font-size: 0.85rem;
  position: relative;
  z-index: 1;
  border-top: 1px solid rgba(129, 140, 160, 0.1);
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 15px;
  flex-wrap: wrap;
}
.footer-link {
  font-size: 0.85rem;
}
.footer-link .el-link--inner { /* 确保链接颜色应用 */
    color: var(--el-color-primary);
}

/* --- 社群/捐赠弹窗样式 --- */
:deep(.el-dialog.community-dialog) {
  background-color: #1f2940;
  border-radius: 10px;
  max-width: 900px; /* 限制最大宽度 */
}
:deep(.el-dialog.community-dialog .el-dialog__title) {
  color: #e0e6f0;
}
:deep(.el-dialog.community-dialog .el-dialog__header) {
    margin-right: 0; /* 修正 Element Plus 标题默认样式 */
    border-bottom: 1px solid rgba(129, 140, 160, 0.2); /* 添加头部底部分割线 */
    padding-bottom: 10px;
}
:deep(.el-dialog.community-dialog .el-dialog__body) {
    padding: 20px 25px; /* 调整 body 内边距 */
}
:deep(.el-dialog.community-dialog .el-dialog__footer) {
    border-top: 1px solid rgba(129, 140, 160, 0.2); /* 添加脚部顶部分割线 */
    padding-top: 15px;
}
:deep(.el-dialog.community-dialog .el-dialog__headerbtn .el-dialog__close) {
    color: #a0aec0;
}
:deep(.el-dialog.community-dialog .el-dialog__headerbtn .el-dialog__close:hover) {
    color: #ffffff;
}

.community-dialog-content {
  /* padding: 10px 20px; 由 dialog body 控制 */
}

.qr-code-col {
  text-align: center;
  margin-bottom: 20px;
}

.qr-code-image {
  width: 100%;
  max-width: 200px;
  height: auto;
  max-height: 200px; /* 限制高度 */
  margin-bottom: 10px;
  border: 1px solid rgba(129, 140, 160, 0.2);
  border-radius: 5px;
  display: inline-block; /* 防止下方多余间隙 */
}
.qr-code-image .image-slot {
    display: flex;
    justify-content: center;
    align-items: center;
    width: 100%;
    min-height: 150px;
    height: 100%; /* 占满父容器高度 */
    max-height: 200px;
    background: #161e2e;
    color: #6b7280;
    font-size: 14px;
    border-radius: 5px; /* 保持圆角一致 */
}

.qr-code-description {
  font-size: 0.9rem;
  color: #b0c4de;
  margin-top: 8px;
}

.dialog-tip {
  text-align: center;
  margin-top: 25px;
  font-size: 0.95rem;
  color: #a0aec0;
}


/* --- 响应式调整 --- */
@media (max-width: 992px) {
    .tech-navigation-dashboard-v2 { padding: 0 40px; }
    .main-content { padding-top: 75px; padding-bottom: 30px; }
    .top-actions { top: 50px; right: 20px; } /* 顶部信息栏高了点，按钮下移 */
    .top-info-banner { padding: 8px 40px; }
}

@media (max-width: 768px) {
  .tech-navigation-dashboard-v2 { padding: 0 20px; }
  .main-content { padding-top: 80px; padding-bottom: 20px; } /* 信息栏换行后可能更高，padding 调整 */

  .header-section { flex-direction: column; text-align: center; margin-bottom: 40px; }
  .el-image { margin-bottom: 20px; margin-right: 0 !important; }
  .dashboard-title { font-size: 2.2rem; }
  .dashboard-subtitle { font-size: 1rem; }
  .card-content { padding: 20px; gap: 15px; }
  .card-text h3 { font-size: 1.15rem; }
  .card-text p { font-size: 0.88rem; }
  .top-actions { top: 55px; right: 15px; } /* 顶部信息栏换行后可能更高，按钮下移 */
  .top-actions .el-button--small { padding: 8px 12px; }
  .top-info-banner {
       padding: 8px 20px;
       font-size: 0.8rem;
       flex-direction: column;
       align-items: flex-start;
  }
   .banner-content { width: 100%; margin-bottom: 5px; }
   .close-banner-icon { align-self: flex-end; }

   .dashboard-footer {
       flex-direction: column;
       gap: 8px;
       font-size: 0.8rem;
       padding: 15px 0;
   }
   :deep(.el-dialog.community-dialog) {
       width: 90% !important; /* 移动端弹窗宽度 */
   }
    :deep(.el-dialog.community-dialog .el-dialog__body) {
        padding: 15px; /* 移动端 body 内边距减小 */
    }
   .qr-code-image {
        max-width: 150px; /* 移动端图片稍小 */
        max-height: 150px;
   }
    .qr-code-image .image-slot {
        min-height: 120px;
        max-height: 150px;
    }
   .qr-code-description { font-size: 0.85rem; }
   .dialog-tip { margin-top: 15px; font-size: 0.9rem; }
}
</style>