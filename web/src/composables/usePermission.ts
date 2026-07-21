import { computed } from 'vue'
import { useUserStoreHook } from '@/store/modules/user'
import { 
  checkPermission, 
  checkResourcePermission, 
  isAdmin, 
  canEdit, 
  isViewer,
  getCurrentRole 
} from '@/utils/permission'

/**
 * 权限管理 Composable
 * 提供各种权限检查功能
 */
export function usePermission() {
  const userStore = useUserStoreHook()
  
  // 当前用户角色
  const currentRoles = computed(() => userStore.roles)
  
  // 当前用户角色
  const currentRole = computed(() => getCurrentRole())
  
  // 快速角色检查
  const permissions = computed(() => ({
    isAdmin: isAdmin(),
    canEdit: canEdit(),
    isViewer: isViewer()
  }))
  
  /**
   * 检查角色权限
   * @param roles 需要的角色列表
   */
  const hasRole = (roles: string | string[]) => {
    const roleArray = Array.isArray(roles) ? roles : [roles]
    return checkPermission(roleArray)
  }
  
  /**
   * 检查资源操作权限
   * @param resource 资源名称
   * @param action 操作类型
   */
  const hasResourcePermission = (resource: string, action: 'read' | 'write' | 'delete') => {
    return checkResourcePermission(resource, action)
  }
  
  /**
   * 检查是否可以访问管理功能
   */
  const canAccessAdmin = computed(() => isAdmin())
  
  /**
   * 检查是否可以编辑资源
   */
  const canEditResource = computed(() => canEdit())
  
  /**
   * 检查是否只能查看资源
   */
  const isReadOnly = computed(() => isViewer())
  
  /**
   * 获取用户可访问的菜单项
   * @param menuItems 菜单项列表
   */
  const getAccessibleMenuItems = (menuItems: any[]) => {
    return menuItems.filter(item => {
      if (!item.meta?.roles) return true
      return checkPermission(item.meta.roles)
    })
  }
  
  /**
   * 检查是否可以执行危险操作（删除等）
   */
  const canPerformDangerousAction = computed(() => {
    // 只有管理员可以执行危险操作
    return isAdmin()
  })
  
  /**
   * 根据权限获取操作按钮配置
   */
  const getActionButtons = (resource: string) => {
    const buttons = []
    
    if (hasResourcePermission(resource, 'read')) {
      buttons.push({ action: 'view', label: '查看', type: 'primary' })
    }
    
    if (hasResourcePermission(resource, 'write')) {
      buttons.push({ action: 'edit', label: '编辑', type: 'warning' })
    }
    
    if (hasResourcePermission(resource, 'delete')) {
      buttons.push({ action: 'delete', label: '删除', type: 'danger' })
    }
    
    return buttons
  }
  
  return {
    // 状态
    currentRoles,
    currentRole,
    permissions,
    canAccessAdmin,
    canEditResource,
    isReadOnly,
    canPerformDangerousAction,
    
    // 方法
    hasRole,
    hasResourcePermission,
    getAccessibleMenuItems,
    getActionButtons
  }
}

/**
 * 权限检查 Hook，用于在组件中快速检查权限
 */
export function usePermissionCheck() {
  return {
    checkPermission,
    checkResourcePermission,
    isAdmin,
    canEdit,
    isViewer
  }
}