<template>
  <div class="tech-stack-container">
    <!-- Breadcrumbs -->
    <el-breadcrumb separator="/" class="page-breadcrumb">
      <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
      <el-breadcrumb-item>关于</el-breadcrumb-item>
      <el-breadcrumb-item>技术栈</el-breadcrumb-item>
    </el-breadcrumb>

    <!-- Page Header -->
    <div class="page-header">
      <h1 class="page-title">希里安项目技术栈 (CILIKUBE v{{ appVersion }})</h1>
       <el-tooltip content="刷新后端依赖" placement="top">
         <el-button :icon="Refresh" circle text @click="fetchBackendStack" :loading="loading.backend"/>
       </el-tooltip>
    </div>

    <el-row :gutter="24">
      <!-- Frontend Section -->
      <el-col :xs="24" :lg="8" class="stack-section">
        <el-card class="box-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon><Monitor /></el-icon>
              <span>前端 (UI)</span>
            </div>
          </template>
          <!-- Use computed property for frontend stack -->
          <el-table :data="formattedFrontendStack" stripe style="width: 100%" class="tech-table" height="400px">
            <el-table-column prop="name" label="库名" min-width="160">
                <template #default="{ row }">
                    <!-- Basic link generation (improve if needed) -->
                    <el-link v-if="row.link" :href="row.link" target="_blank" type="primary">{{ row.name }}</el-link>
                     <span v-else>{{ row.name }}</span>
                </template>
            </el-table-column>
            <el-table-column prop="version" label="版本 (package.json)" min-width="140" align="center">
                 <template #default="{ row }">
                    <el-tag size="small" effect="plain">{{ row.version }}</el-tag>
                 </template>
            </el-table-column>
             <!-- Optional: Add a description column manually mapping known libraries -->
             <el-table-column prop="description" label="描述" min-width="180" show-overflow-tooltip>
                   <template #default="{ row }">
                       {{ getKnownDescription(row.name) }}
                   </template>
             </el-table-column>
          </el-table>
           <div class="card-footer-note">前端依赖版本来自 <code>package.json</code></div>
        </el-card>
      </el-col>

      <!-- Backend Section -->
      <el-col :xs="24" :lg="8" class="stack-section">
        <el-card class="box-card" shadow="hover" v-loading="loading.backend">
          <template #header>
            <div class="card-header">
               <el-icon><Platform /></el-icon>
              <span>后端 (API)</span>
            </div>
          </template>
           <div v-if="backendError" class="error-message">
              <el-alert type="error" :closable="false" show-icon>
                加载后端依赖失败: {{ backendError }}
              </el-alert>
           </div>
           <!-- Use reactive ref for backend stack -->
          <el-table v-else :data="backendStack" stripe style="width: 100%" class="tech-table" height="400px">
             <el-table-column prop="path" label="模块路径" min-width="250" show-overflow-tooltip>
                 <template #default="{ row }">
                    <!-- You might want to format path or link to pkg.go.dev -->
                    <el-link :href="`https://pkg.go.dev/${row.path}`" target="_blank" type="primary">{{ row.path }}</el-link>
                 </template>
             </el-table-column>
            <el-table-column prop="version" label="版本 (go.mod)" min-width="150" align="center">
                 <template #default="{ row }">
                     <el-tag v-if="row.version" size="small" effect="plain">{{ row.version }}</el-tag>
                     <span v-else>-</span>
                 </template>
            </el-table-column>
          </el-table>
          <div v-if="!backendError" class="card-footer-note">后端依赖版本来自 <code>go.mod</code></div>
        </el-card>
      </el-col>

        <!-- DevOps/Tooling Section (Manual) -->
      <el-col :xs="24" :lg="8" class="stack-section">
        <el-card class="box-card" shadow="hover">
          <template #header>
            <div class="card-header">
               <el-icon><Setting /></el-icon>
              <span>开发与部署工具 (手动维护)</span>
            </div>
          </template>
          <el-table :data="toolingStack" stripe style="width: 100%" class="tech-table" height="400px">
            <el-table-column prop="name" label="工具" min-width="150">
                 <template #default="{ row }">
                    <el-link v-if="row.link" :href="row.link" target="_blank" type="primary">{{ row.name }}</el-link>
                    <span v-else>{{ row.name }}</span>
                </template>
            </el-table-column>
            <el-table-column prop="version" label="版本 (参考)" min-width="120" align="center">
                 <template #default="{ row }">
                     <el-tag v-if="row.version" size="small" effect="plain">{{ row.version }}</el-tag>
                      <span v-else>-</span>
                 </template>
            </el-table-column>
            <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip/>
          </el-table>
           <div class="card-footer-note">工具版本需手动更新</div>
        </el-card>
      </el-col>

    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue';
import { ElMessage } from 'element-plus';
import { request } from '@/utils/service'; // Adjust path
import { Platform, Monitor, Setting, Refresh } from '@element-plus/icons-vue';

const VITE_API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://192.168.100:8080";
// --- Data Interfaces ---
interface TechItem {
  name: string;
  version: string;
  description?: string; // Description might be manual
  link?: string;
}

// Backend dependency structure from Go service
interface BackendDependency {
  path: string;
  version: string;
}


// --- Frontend Data (Injected by Vite) ---
// Access the globally defined variables from vite.config.js
// Use declare to inform TypeScript about these global variables
declare const __APP_DEPENDENCIES__: Record<string, string>;
declare const __APP_DEV_DEPENDENCIES__: Record<string, string>;
declare const __APP_VERSION__: string;

const appVersion = ref(__APP_VERSION__ || 'N/A');
const dependencies = ref(__APP_DEPENDENCIES__ || {});
const devDependencies = ref(__APP_DEV_DEPENDENCIES__ || {});

// --- Backend Data ---
const backendStack = ref<BackendDependency[]>([]); // Store fetched backend deps
const loading = reactive({
    page: false, // For initial load indicators if needed
    backend: false,
});
const backendError = ref<string | null>(null);

// --- Tooling Data (Manual) ---
const toolingStack = ref<TechItem[]>([
  { name: 'Docker', version: '2x.x+', description: '容器化平台', link: 'https://www.docker.com/' },
  { name: 'Kubernetes', version: 'v1.xx', description: '容器编排平台', link: 'https://kubernetes.io/' },
  { name: 'Nginx Ingress', version: 'v1.x', description: 'Ingress 控制器', link: 'https://kubernetes.github.io/ingress-nginx/' },
  { name: 'Node.js', version: 'v20+', description: 'JS 运行时 (for build)', link: 'https://nodejs.org/' },
  { name: 'pnpm', version: 'v8+', description: '包管理器', link: 'https://pnpm.io/' },
  { name: 'ESLint', version: 'latest', description: '代码检查', link: 'https://eslint.org/' },
  { name: 'Prettier', version: 'latest', description: '代码格式化', link: 'https://prettier.io/' },
  // Add others (Prometheus, Grafana, CI/CD...) and update versions manually
]);

// --- Computed Properties for Formatted Stacks ---
const formattedFrontendStack = computed(() => {
    const deps = Object.entries(dependencies.value).map(([name, version]) => ({
        name,
        version,
        description: getKnownDescription(name), // Add manual description
        link: getKnownLink(name), // Add manual link
        type: 'prod'
    }));
     const devDeps = Object.entries(devDependencies.value).map(([name, version]) => ({
        name,
        version,
        description: getKnownDescription(name),
        link: getKnownLink(name),
        type: 'dev'
    }));

    // Combine and sort, maybe prioritize prod deps or sort alphabetically
    return [...deps, ...devDeps].sort((a, b) => a.name.localeCompare(b.name));
});

// --- Manual Mappings for Descriptions/Links (Optional) ---
// Maintain a map for known libraries to add descriptions or links
const knownDescriptions: Record<string, string> = {
    'vue': '渐进式 JavaScript 框架',
    'vite': '下一代前端开发与构建工具',
    'typescript': 'JavaScript 的超集',
    '@element-plus/icons-vue': 'Element Plus 图标库',
    'element-plus': '基于 Vue 3 的 UI 组件库',
    'vue-router': 'Vue.js 官方路由',
    'pinia': 'Vue.js 状态管理库',
    'axios': 'HTTP 客户端',
    'echarts': '图表库',
    'vue-echarts': 'ECharts Vue 封装',
    'dayjs': '日期时间处理库',
    'lodash-es': '实用工具库',
    'js-yaml': 'YAML 解析/序列化',
    'js-base64': 'Base64 编解码',
    'js-quantities': '单位转换与解析',
    'sass': 'CSS 预处理器',
    'unocss': '原子化 CSS 引擎',
    // Add more as needed
};
const knownLinks: Record<string, string> = {
    'vue': 'https://vuejs.org/',
    'vite': 'https://vitejs.dev/',
    'typescript': 'https://www.typescriptlang.org/',
    'element-plus': 'https://element-plus.org/',
    'vue-router': 'https://router.vuejs.org/',
    'pinia': 'https://pinia.vuejs.org/',
    'axios': 'https://axios-http.com/',
    'echarts': 'https://echarts.apache.org/',
    'dayjs': 'https://day.js.org/',
    'lodash-es': 'https://lodash.com/',
    'js-yaml': 'https://github.com/nodeca/js-yaml',
    'js-base64': 'https://github.com/dankogai/js-base64',
    'js-quantities': 'https://github.com/gentooboontoo/js-quantities',
    'sass': 'https://sass-lang.com/',
    'unocss': 'https://unocss.dev/',
}

const getKnownDescription = (name: string): string => knownDescriptions[name] || '';
const getKnownLink = (name: string): string | undefined => knownLinks[name];


// --- API Fetching ---
const fetchBackendStack = async () => {
    loading.backend = true;
    backendError.value = null;
    try {
        const response = await request<{ code: number; data: BackendDependency[]; message: string }>({
            url: '/api/v1/summary/backend-dependencies', // Match Go route
            method: 'get',
            baseURL: 'VITE_API_BASE_URL', // If needed
        });
        if (response.code === 200 && Array.isArray(response.data)) {
            // Sort fetched data alphabetically by path
             backendStack.value = response.data.sort((a, b) => a.path.localeCompare(b.path));
        } else {
            throw new Error(response.message || '获取后端依赖数据格式错误');
        }
    } catch (err: any) {
        console.error("Failed to fetch backend dependencies:", err);
        backendError.value = err.message || '网络请求失败';
        backendStack.value = []; // Clear data on error
    } finally {
        loading.backend = false;
    }
};

// --- Lifecycle ---
onMounted(() => {
    fetchBackendStack(); // Fetch backend data when component mounts
});

</script>

<style lang="scss" scoped>
// Using fallback variables directly
$page-padding: 20px; $spacing-md: 15px; $spacing-lg: 20px; $border-radius-base: 4px;

.tech-stack-container { padding: $page-padding; background-color: var(--el-bg-color-page); }
.page-breadcrumb { margin-bottom: $spacing-lg; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: $spacing-lg + 5px; .el-button { margin-left: $spacing-md; } }
.page-title { font-size: 24px; font-weight: 600; color: var(--el-text-color-primary); margin: 0; }
.stack-section { margin-bottom: $spacing-lg; }

.box-card {
    border: 1px solid var(--el-border-color-lighter); height: 100%; display: flex; flex-direction: column;
    :deep(.el-card__header) { padding: 14px 20px; background-color: var(--el-fill-color-light); border-bottom: 1px solid var(--el-border-color-lighter); }
    :deep(.el-card__body) { padding: 0; flex-grow: 1; overflow: hidden; /* Prevent table breaking card layout */ }
}
.card-header { display: flex; align-items: center; font-size: 16px; font-weight: 600; color: var(--el-text-color-primary); .el-icon { margin-right: 8px; font-size: 18px; } }
.tech-table {
     :deep(th.el-table__cell) { background-color: var(--el-fill-color-lighter); color: var(--el-text-color-secondary); font-weight: 500; font-size: 13px; }
    :deep(td.el-table__cell) { padding: 10px 0; font-size: 14px; vertical-align: middle; }
    .el-link { font-size: 14px; }
}
.error-message { padding: 20px; }
.card-footer-note {
    padding: 8px 15px;
    font-size: 12px;
    color: var(--el-text-color-secondary);
    background-color: var(--el-fill-color-lighter);
    border-top: 1px solid var(--el-border-color-extra-light);
    code {
        background-color: var(--el-color-info-light-9);
        padding: 2px 4px;
        border-radius: 3px;
    }
}
</style>