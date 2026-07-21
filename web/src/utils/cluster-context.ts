// 全局集群上下文管理
// 这个文件提供全局的集群上下文，避免在service中直接访问pinia store

let currentClusterId: string | null = null

/** 设置当前集群ID */
export function setCurrentClusterId(id: string | null) {
  currentClusterId = id
}

/** 获取当前集群ID */
export function getCurrentClusterId(): string | null {
  return currentClusterId
}

/** 从localStorage获取保存的集群ID */
export function getSavedClusterId(): string | null {
  return localStorage.getItem('selectedClusterId')
}

/** 从localStorage获取保存的集群名称（向后兼容） */
export function getSavedClusterName(): string | null {
  return localStorage.getItem('selectedClusterName')
}