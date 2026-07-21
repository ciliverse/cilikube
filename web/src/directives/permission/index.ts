import { type Directive } from "vue"
import { useUserStoreHook } from "@/store/modules/user"
import { checkResourcePermission } from "@/utils/permission"

/** 权限指令，和权限判断函数 checkPermission 功能类似 */
export const permission: Directive = {
  mounted(el, binding) {
    const { value: permissionRoles } = binding
    const { roles } = useUserStoreHook()
    if (Array.isArray(permissionRoles) && permissionRoles.length > 0) {
      const hasPermission = roles.some((role) => permissionRoles.includes(role))
      // hasPermission || (el.style.display = "none") // 隐藏
      hasPermission || el.parentNode?.removeChild(el) // 销毁
    } else {
      throw new Error(`need roles! Like v-permission="['admin','editor']"`)
    }
  }
}

/** 资源权限指令，用于检查对特定资源的操作权限 */
export const resourcePermission: Directive = {
  mounted(el, binding) {
    const { value, arg } = binding
    
    if (typeof value === 'string' && arg) {
      // 使用方式: v-resource-permission:read="'pod'" 或 v-resource-permission:write="'deployment'"
      const resource = value
      const action = arg as 'read' | 'write' | 'delete'
      const hasPermission = checkResourcePermission(resource, action)
      
      if (!hasPermission) {
        el.parentNode?.removeChild(el)
      }
    } else if (typeof value === 'object' && value.resource && value.action) {
      // 使用方式: v-resource-permission="{ resource: 'pod', action: 'write' }"
      const hasPermission = checkResourcePermission(value.resource, value.action)
      
      if (!hasPermission) {
        el.parentNode?.removeChild(el)
      }
    } else {
      throw new Error(`need resource and action! Like v-resource-permission:read="'pod'" or v-resource-permission="{ resource: 'pod', action: 'write' }"`)
    }
  }
}

/** 角色权限指令，用于快速检查特定角色 */
export const rolePermission: Directive = {
  mounted(el, binding) {
    const { value } = binding
    const { roles } = useUserStoreHook()
    
    let hasPermission = false
    
    if (typeof value === 'string') {
      // 使用方式: v-role-permission="'admin'"
      hasPermission = roles.includes(value)
    } else if (Array.isArray(value)) {
      // 使用方式: v-role-permission="['admin', 'editor']"
      hasPermission = roles.some(role => value.includes(role))
    } else {
      throw new Error(`need role string or array! Like v-role-permission="'admin'" or v-role-permission="['admin','editor']"`)
    }
    
    if (!hasPermission) {
      el.parentNode?.removeChild(el)
    }
  }
}
