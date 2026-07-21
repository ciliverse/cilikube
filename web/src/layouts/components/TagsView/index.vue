<script lang="ts" setup>
import { getCurrentInstance, onMounted, ref, watch } from "vue"
import { type RouteLocationNormalizedLoaded, type RouteRecordRaw, RouterLink, useRoute, useRouter } from "vue-router"
import { useI18n } from "vue-i18n"
import { type TagView, useTagsViewStore } from "@/store/modules/tags-view"
import { usePermissionStore } from "@/store/modules/permission"
import { useRouteListener } from "@/hooks/useRouteListener"
import path from "path-browserify"
import ScrollPane from "./ScrollPane.vue"
import { Close } from "@element-plus/icons-vue"

const instance = getCurrentInstance()
const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const tagsViewStore = useTagsViewStore()
const permissionStore = usePermissionStore()
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
  '权限': 'menu.permission',
  '页面级权限': 'menu.pageLevel',
  '页面级': 'menu.pageLevel',
  '按钮级权限': 'menu.buttonLevel',
  '按钮级': 'menu.buttonLevel',
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

/** 标签页组件元素的引用数组 */
const tagRefs = ref<InstanceType<typeof RouterLink>[]>([])

/** 右键菜单的状态 */
const visible = ref(false)
/** 右键菜单的 top 位置 */
const top = ref(0)
/** 右键菜单的 left 位置 */
const left = ref(0)
/** 当前正在右键操作的标签页 */
const selectedTag = ref<TagView>({})
/** 固定的标签页 */
let affixTags: TagView[] = []

/** 判断标签页是否激活 */
const isActive = (tag: TagView) => {
  return tag.path === route.path
}

/** 判断标签页是否固定 */
const isAffix = (tag: TagView) => {
  return tag.meta?.affix
}

/** 筛选出固定标签页 */
const filterAffixTags = (routes: RouteRecordRaw[], basePath = "/") => {
  const tags: TagView[] = []
  routes.forEach((route) => {
    if (isAffix(route)) {
      const tagPath = path.resolve(basePath, route.path)
      tags.push({
        fullPath: tagPath,
        path: tagPath,
        name: route.name,
        meta: { ...route.meta }
      })
    }
    if (route.children) {
      const childTags = filterAffixTags(route.children, route.path)
      tags.push(...childTags)
    }
  })
  return tags
}

/** 初始化标签页 */
const initTags = () => {
  affixTags = filterAffixTags(permissionStore.routes)
  for (const tag of affixTags) {
    // 必须含有 name 属性
    tag.name && tagsViewStore.addVisitedView(tag)
  }
}

/** 添加标签页 */
const addTags = (route: RouteLocationNormalizedLoaded) => {
  if (route.name) {
    tagsViewStore.addVisitedView(route)
    tagsViewStore.addCachedView(route)
  }
}

/** 刷新当前正在右键操作的标签页 */
const refreshSelectedTag = (view: TagView) => {
  tagsViewStore.delCachedView(view)
  router.replace({ path: "/redirect" + view.path, query: view.query })
}

/** 关闭当前正在右键操作的标签页 */
const closeSelectedTag = (view: TagView) => {
  tagsViewStore.delVisitedView(view)
  tagsViewStore.delCachedView(view)
  isActive(view) && toLastView(tagsViewStore.visitedViews, view)
}

/** 关闭其他标签页 */
const closeOthersTags = () => {
  const fullPath = selectedTag.value.fullPath
  if (fullPath !== route.path && fullPath !== undefined) {
    router.push(fullPath)
  }
  tagsViewStore.delOthersVisitedViews(selectedTag.value)
  tagsViewStore.delOthersCachedViews(selectedTag.value)
}

/** 关闭所有标签页 */
const closeAllTags = (view: TagView) => {
  tagsViewStore.delAllVisitedViews()
  tagsViewStore.delAllCachedViews()
  if (affixTags.some((tag) => tag.path === route.path)) return
  toLastView(tagsViewStore.visitedViews, view)
}

/** 跳转到最后一个标签页 */
const toLastView = (visitedViews: TagView[], view: TagView) => {
  const latestView = visitedViews.slice(-1)[0]
  const fullPath = latestView?.fullPath
  if (fullPath !== undefined) {
    router.push(fullPath)
  } else {
    // 如果 TagsView 全部被关闭了，则默认重定向到主页
    if (view.name === "Dashboard") {
      // 重新加载主页
      router.push({ path: "/redirect" + view.path, query: view.query })
    } else {
      router.push("/")
    }
  }
}

/** 打开右键菜单面板 */
const openMenu = (tag: TagView, e: MouseEvent) => {
  const menuMinWidth = 105
  // 当前组件距离浏览器左端的距离
  const offsetLeft = instance!.proxy!.$el.getBoundingClientRect().left
  // 当前组件宽度
  const offsetWidth = instance!.proxy!.$el.offsetWidth
  // 面板的最大左边距
  const maxLeft = offsetWidth - menuMinWidth
  // 面板距离鼠标指针的距离
  const left15 = e.clientX - offsetLeft + 15
  left.value = left15 > maxLeft ? maxLeft : left15
  top.value = e.clientY
  // 显示面板
  visible.value = true
  // 更新当前正在右键操作的标签页
  selectedTag.value = tag
}

/** 关闭右键菜单面板 */
const closeMenu = () => {
  visible.value = false
}

watch(visible, (value) => {
  value ? document.body.addEventListener("click", closeMenu) : document.body.removeEventListener("click", closeMenu)
})

onMounted(() => {
  initTags()
  /** 监听路由变化 */
  listenerRouteChange(async (route) => {
    addTags(route)
  }, true)
})
</script>

<template>
  <div class="tags-view-container">
    <ScrollPane class="tags-view-wrapper" :tag-refs="tagRefs">
      <router-link ref="tagRefs" v-for="tag in tagsViewStore.visitedViews" :key="tag.path"
        :class="{ active: isActive(tag) }" class="tags-view-item" :to="{ path: tag.path, query: tag.query }"
        @click.middle="!isAffix(tag) && closeSelectedTag(tag)" @contextmenu.prevent="openMenu(tag, $event)">
        {{ getI18nTitle(tag.meta?.title || '') }}
        <el-icon v-if="!isAffix(tag)" :size="12" @click.prevent.stop="closeSelectedTag(tag)">
          <Close />
        </el-icon>
      </router-link>
    </ScrollPane>
    <ul v-show="visible" class="contextmenu" :style="{ left: left + 'px', top: top + 'px' }">
      <li @click="refreshSelectedTag(selectedTag)">{{ t('tagsView.refresh') }}</li>
      <li v-if="!isAffix(selectedTag)" @click="closeSelectedTag(selectedTag)">{{ t('tagsView.close') }}</li>
      <li @click="closeOthersTags">{{ t('tagsView.closeOthers') }}</li>
      <li @click="closeAllTags(selectedTag)">{{ t('tagsView.closeAll') }}</li>
    </ul>
  </div>
</template>

<style lang="scss" scoped>
.tags-view-container {
  height: var(--ck-tagsview-height);
  width: 100%;
  color: var(--ck-tagsview-text-color);
  overflow: hidden;

  .tags-view-wrapper {
    .tags-view-item {
      display: inline-block;
      position: relative;
      cursor: pointer;
      height: 26px;
      line-height: 26px;
      border: 1px solid var(--ck-tagsview-tag-border-color);
      border-radius: var(--ck-tagsview-tag-border-radius);
      background-color: var(--ck-tagsview-tag-bg-color);
      padding: 0 8px;
      font-size: 12px;
      margin-left: 5px;
      margin-top: 4px;

      &:first-of-type {
        margin-left: 5px;
      }

      &:last-of-type {
        margin-right: 5px;
      }

      &.active {
        background-color: var(--ck-tagsview-tag-active-bg-color);
        color: var(--ck-tagsview-tag-active-text-color);
        border-color: var(--ck-tagsview-tag-active-border-color);
      }

      .el-icon {
        margin: 0 2px;
        vertical-align: middle;
        border-radius: 50%;

        &:hover {
          background-color: var(--ck-tagsview-tag-icon-hover-bg-color);
          color: var(--ck-tagsview-tag-icon-hover-color);
        }
      }
    }
  }

  .contextmenu {
    margin: 0;
    z-index: 3000;
    position: absolute;
    list-style-type: none;
    padding: 5px 0;
    border-radius: 4px;
    font-size: 12px;
    color: var(--ck-tagsview-contextmenu-text-color);
    background-color: var(--ck-tagsview-contextmenu-bg-color);
    box-shadow: var(--ck-tagsview-contextmenu-box-shadow);

    li {
      margin: 0;
      padding: 7px 16px;
      cursor: pointer;

      &:hover {
        color: var(--ck-tagsview-contextmenu-hover-text-color);
        background-color: var(--ck-tagsview-contextmenu-hover-bg-color);
      }
    }
  }
}
</style>
