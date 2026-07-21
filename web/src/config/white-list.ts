import { type RouteLocationNormalized } from "vue-router"

/** 免登录白名单（匹配路由 path） */
const whiteListByPath: string[] = ["/login", "/", "/navi", "/board", "/board/dashboard", "/board/summary", "/board/minikube", "/board/navy", "/cluster", "/cluster/management", "/cluster/node", "/cluster/namespace", "/workloads", "/workloads/pods", "/workloads/deployments", "/storage", "/storage/persistentvolume", "/storage/persistentvolumeclaim", "/network", "/network/service", "/network/ingress", "/config", "/config/configmap", "/config/secret", "/link", "/link/techstack", "/permission", "/permission/page", "/permission/directive"]

/** 免登录白名单（匹配路由 name） */
const whiteListByName: string[] = []

/** 判断是否在白名单 */
const isWhiteList = (to: RouteLocationNormalized) => {
  // path 和 name 任意一个匹配上即可
  return whiteListByPath.indexOf(to.path) !== -1 || whiteListByName.indexOf(to.name as any) !== -1
}

export default isWhiteList
