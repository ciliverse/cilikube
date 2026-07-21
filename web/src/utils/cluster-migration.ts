/**
 * 集群存储迁移工具
 * 帮助用户从name-based系统迁移到ID-based系统
 */

import { useClusterStore } from "@/store/modules/clusterStore"

const OLD_STORAGE_KEY = "selectedClusterName"
const NEW_STORAGE_KEY = "selectedClusterId"

/**
 * 执行集群存储迁移
 * 将旧的name-based存储迁移到新的ID-based存储
 */
export async function migrateClusterStorage(): Promise<boolean> {
  try {
    // 检查是否需要迁移
    const oldClusterName = localStorage.getItem(OLD_STORAGE_KEY)
    const newClusterId = localStorage.getItem(NEW_STORAGE_KEY)
    
    // 如果已经有新的ID存储，或者没有旧的name存储，则不需要迁移
    if (newClusterId || !oldClusterName) {
      return false
    }
    
    console.log(`开始迁移集群存储: ${oldClusterName} -> ID`)
    
    // 获取集群存储实例
    const clusterStore = useClusterStore()
    
    // 确保集群列表已加载
    if (clusterStore.availableClusters.length === 0) {
      await clusterStore.fetchAvailableClusters()
    }
    
    // 根据名称查找对应的集群ID
    const targetCluster = clusterStore.availableClusters.find(
      cluster => cluster.name === oldClusterName
    )
    
    if (targetCluster) {
      // 设置新的ID-based存储
      await clusterStore.setSelectedClusterId(targetCluster.id)
      
      // 清除旧的name-based存储
      localStorage.removeItem(OLD_STORAGE_KEY)
      
      console.log(`集群存储迁移成功: ${oldClusterName} -> ${targetCluster.id}`)
      return true
    } else {
      console.warn(`迁移失败: 未找到名为 "${oldClusterName}" 的集群`)
      // 清除无效的旧存储
      localStorage.removeItem(OLD_STORAGE_KEY)
      return false
    }
  } catch (error) {
    console.error("集群存储迁移失败:", error)
    return false
  }
}

/**
 * 检查是否需要迁移
 */
export function needsMigration(): boolean {
  const oldClusterName = localStorage.getItem(OLD_STORAGE_KEY)
  const newClusterId = localStorage.getItem(NEW_STORAGE_KEY)
  
  return !!(oldClusterName && !newClusterId)
}

/**
 * 清理旧的存储数据
 */
export function cleanupOldStorage(): void {
  localStorage.removeItem(OLD_STORAGE_KEY)
  console.log("已清理旧的集群存储数据")
}

/**
 * 获取迁移状态信息
 */
export function getMigrationStatus(): {
  needsMigration: boolean
  oldClusterName: string | null
  newClusterId: string | null
} {
  return {
    needsMigration: needsMigration(),
    oldClusterName: localStorage.getItem(OLD_STORAGE_KEY),
    newClusterId: localStorage.getItem(NEW_STORAGE_KEY)
  }
}