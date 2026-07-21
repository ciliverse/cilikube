import { useUserStoreHook } from "@/store/modules/user"

/** 全局权限判断函数，和权限指令 v-permission 功能类似 */
export const checkPermission = (permissionRoles: string[]): boolean => {
  if (Array.isArray(permissionRoles) && permissionRoles.length > 0) {
    const { roles } = useUserStoreHook()
    return roles.some((role) => permissionRoles.includes(role))
  } else {
    console.error("need roles! Like checkPermission(['admin','editor'])")
    return false
  }
}

/** 检查用户是否有管理员权限 */
export const isAdmin = (): boolean => {
  const { roles } = useUserStoreHook()
  return roles.includes("admin")
}

/** 检查用户是否有编辑权限（管理员或编辑者） */
export const canEdit = (): boolean => {
  const { roles } = useUserStoreHook()
  return roles.some(role => ["admin", "editor"].includes(role))
}

/** 检查用户是否只有查看权限 */
export const isViewer = (): boolean => {
  const { roles } = useUserStoreHook()
  return roles.includes("viewer") && !canEdit()
}

/** 检查用户对特定资源的操作权限 */
export const checkResourcePermission = (resource: string, action: 'read' | 'write' | 'delete'): boolean => {
  const { roles } = useUserStoreHook()
  
  // 管理员拥有所有权限
  if (roles.includes("admin")) {
    return true
  }
  
  // 编辑者权限规则
  if (roles.includes("editor")) {
    // 编辑者不能删除集群或进行集群管理操作
    if (resource === "cluster" && action === "delete") {
      return false
    }
    if (resource === "cluster-management") {
      return false
    }
    // 编辑者不能查看或操作 Secret
    if (resource === "secret") {
      return false
    }
    // 其他资源编辑者可以读写
    return action !== "delete" || !["namespace", "node"].includes(resource)
  }
  
  // 查看者只有读权限
  if (roles.includes("viewer")) {
    if (action === "read") {
      // 查看者不能查看 Secret
      return resource !== "secret"
    }
    return false
  }
  
  return false
}

/** 获取当前用户角色 */
export const getCurrentRole = (): string => {
  const { roles } = useUserStoreHook()
  
  if (roles.includes("admin")) {
    return "admin"
  } else if (roles.includes("editor")) {
    return "editor"
  } else if (roles.includes("viewer")) {
    return "viewer"
  }
  
  return "unknown"
}
