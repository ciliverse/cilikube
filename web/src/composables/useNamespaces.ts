import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { httpClient } from '@/utils/service'
import { useClusterStore } from '@/store/modules/clusterStore'

export function useNamespaces() {
  const clusterStore = useClusterStore()
  
  const namespaces = ref<string[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  // 获取命名空间列表
  const fetchNamespaces = async () => {
    loading.value = true
    error.value = null
    
    try {
      const clusterId = clusterStore.selectedClusterId || undefined
      const params = clusterId ? { clusterId } : {}
      
      const response = await httpClient.standard.get('/api/v1/namespaces', { params })
      
      if ((response.code === 200 || response.code === 0) && response.data?.items) {
        namespaces.value = response.data.items.map((item: any) => item.metadata?.name).filter(Boolean)
      } else {
        error.value = response.message || '获取命名空间列表失败'
        ElMessage.error(error.value)
      }
    } catch (err) {
      error.value = '获取命名空间列表失败'
      ElMessage.error(error.value)
      console.error('Failed to fetch namespaces:', err)
    } finally {
      loading.value = false
    }
  }

  return {
    namespaces,
    loading,
    error,
    fetchNamespaces
  }
}