<template>
  <div v-if="hasPermission">
    <slot />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { checkPermission, checkResourcePermission, isAdmin, canEdit, isViewer } from '@/utils/permission'

interface Props {
  /** 角色权限检查 */
  roles?: string[]
  /** 资源权限检查 */
  resource?: string
  /** 操作类型 */
  action?: 'read' | 'write' | 'delete'
  /** 快速权限检查 */
  admin?: boolean
  editor?: boolean
  viewer?: boolean
  /** 自定义权限检查函数 */
  check?: () => boolean
}

const props = withDefaults(defineProps<Props>(), {
  roles: () => [],
  admin: false,
  editor: false,
  viewer: false
})

const hasPermission = computed(() => {
  // 如果提供了自定义检查函数，优先使用
  if (props.check) {
    return props.check()
  }
  
  // 快速角色检查
  if (props.admin) {
    return isAdmin()
  }
  
  if (props.editor) {
    return canEdit()
  }
  
  if (props.viewer) {
    return isViewer()
  }
  
  // 资源权限检查
  if (props.resource && props.action) {
    return checkResourcePermission(props.resource, props.action)
  }
  
  // 角色权限检查
  if (props.roles.length > 0) {
    return checkPermission(props.roles)
  }
  
  // 默认允许访问
  return true
})
</script>