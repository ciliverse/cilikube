<script lang="ts" setup>
import { ref } from "vue"
import { type RouteLocationMatched, useRoute, useRouter } from "vue-router"
import { useI18n } from "vue-i18n"
import { useRouteListener } from "@/hooks/useRouteListener"
import { compile } from "path-to-regexp"

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const { listenerRouteChange } = useRouteListener()

// 标题映射表，将中文标题映射到国际化key
const titleMap: Record<string, string> = {
  '看板': 'menu.board',
  '集群安装': 'menu.clusterInstall',
  '集群概览': 'menu.clusterOverview',
  '集群监控': 'menu.clusterMonitor',
  '运维导航': 'menu.opsNavigation',
  '集群': 'menu.cluster',
  '集群管理': 'menu.clusterManagement',
  '集群事件': 'menu.clusterEvents',
  '自定义资源': 'menu.crd',
  '节点': 'menu.node',
  '命名空间': 'menu.namespace',
  '工作负载': 'menu.workloads',
  'Pod': 'menu.pod',
  'pod': 'menu.pod',
  'Deployment': 'menu.deployment',
  'deployment': 'menu.deployment',
  '存储': 'menu.storage',
  'PV': 'menu.pv',
  'pv': 'menu.pv',
  'PVC': 'menu.pvc',
  'pvc': 'menu.pvc',
  '网络': 'menu.network',
  'Service': 'menu.service',
  'service': 'menu.service',
  'Ingress': 'menu.ingress',
  'ingress': 'menu.ingress',
  '配置管理': 'menu.config',
  'ConfigMap': 'menu.configmap',
  'configmap': 'menu.configmap',
  'Secret': 'menu.secret',
  'secret': 'menu.secret',
  '项目管理': 'menu.project',
  '中文文档': 'menu.chineseDocs',
  '我的博客': 'menu.myBlog',
  '技术栈': 'menu.techStack',
  '权限演示': 'menu.permission',
  '页面级权限': 'menu.pageLevel',
  '按钮级权限': 'menu.buttonLevel',
  '系统管理': 'menu.admin',
  '用户管理': 'menu.userManagement',
  '角色管理': 'menu.roleManagement',
  '系统设置': 'menu.systemSettings'
}

// 获取国际化标题
const getI18nTitle = (title: string) => {
  const i18nKey = titleMap[title]
  return i18nKey ? t(i18nKey) : title
}

/** 定义响应式数据 breadcrumbs，用于存储面包屑导航信息 */
const breadcrumbs = ref<RouteLocationMatched[]>([])

/** 获取面包屑导航信息 */
const getBreadcrumb = () => {
  breadcrumbs.value = route.matched.filter((item) => item.meta?.title && item.meta?.breadcrumb !== false)
}

/** 编译路由路径 */
const pathCompile = (path: string) => {
  const toPath = compile(path)
  return toPath(route.params)
}

/** 处理面包屑导航点击事件 */
const handleLink = (item: RouteLocationMatched) => {
  const { redirect, path } = item
  if (redirect) {
    router.push(redirect as string)
    return
  }
  router.push(pathCompile(path))
}

/** 监听路由变化，更新面包屑导航信息 */
listenerRouteChange((route) => {
  if (route.path.startsWith("/redirect/")) return
  getBreadcrumb()
}, true)
</script>

<template>
  <el-breadcrumb>
    <el-breadcrumb-item v-for="(item, index) in breadcrumbs" :key="item.path">
      <span v-if="item.redirect === 'noRedirect' || index === breadcrumbs.length - 1" class="no-redirect">
        {{ getI18nTitle(item.meta.title || '') }}
      </span>
      <a v-else @click.prevent="handleLink(item)">
        {{ getI18nTitle(item.meta.title || '') }}
      </a>
    </el-breadcrumb-item>
  </el-breadcrumb>
</template>

<style lang="scss" scoped>
.el-breadcrumb {
  line-height: var(--ck-navigationbar-height);
  .no-redirect {
    color: var(--el-text-color-placeholder);
  }
  a {
    font-weight: normal;
  }
}
</style>
