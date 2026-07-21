import { type App } from "vue"
import { permission, resourcePermission, rolePermission } from "./permission"

/** 挂载自定义指令 */
export function loadDirectives(app: App) {
  app.directive("permission", permission)
  app.directive("resource-permission", resourcePermission)
  app.directive("role-permission", rolePermission)
}
