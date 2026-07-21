import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { crdApi, type CRDItem, type CRDDetailResponse, type CustomResourceItem } from '@/api/crd'
import { useClusterStore } from '@/store/modules/clusterStore'

export function useCRD() {
  const { t } = useI18n()
  const clusterStore = useClusterStore()
  
  const crds = ref<CRDItem[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Get CRD list
  const fetchCRDs = async () => {
    loading.value = true
    error.value = null
    
    try {
      const clusterId = clusterStore.selectedClusterId || undefined
      const response = await crdApi.getCRDs(clusterId)
      
      if (response.code === 200 || response.code === 0) {
        crds.value = response.data.items
      } else {
        error.value = response.message || t('crd.fetchCRDsFailed')
        ElMessage.error(error.value)
      }
    } catch (err) {
      error.value = t('crd.fetchCRDsFailed')
      ElMessage.error(error.value)
      console.error('Failed to fetch CRDs:', err)
    } finally {
      loading.value = false
    }
  }

  // 获取CRD详情
  const fetchCRDDetail = async (name: string): Promise<CRDDetailResponse | null> => {
    try {
      const clusterId = clusterStore.selectedClusterId || undefined
      const response = await crdApi.getCRD(name, clusterId)
      
      if (response.code === 200 || response.code === 0) {
        return response.data
      } else {
        ElMessage.error(response.message || t('crd.fetchCRDDetailFailed'))
        return null
      }
    } catch (err) {
      ElMessage.error(t('crd.fetchCRDDetailFailed'))
      console.error('Failed to fetch CRD detail:', err)
      return null
    }
  }

  // 按组分类的CRD
  const crdsByGroup = computed(() => {
    const groups: Record<string, CRDItem[]> = {}
    
    crds.value.forEach(crd => {
      const group = crd.group || 'core'
      if (!groups[group]) {
        groups[group] = []
      }
      groups[group].push(crd)
    })
    
    return groups
  })

  // 搜索CRD
  const searchCRDs = (keyword: string) => {
    if (!keyword) return crds.value
    
    const lowerKeyword = keyword.toLowerCase()
    return crds.value.filter(crd => 
      crd.name.toLowerCase().includes(lowerKeyword) ||
      crd.kind.toLowerCase().includes(lowerKeyword) ||
      crd.group.toLowerCase().includes(lowerKeyword)
    )
  }

  return {
    crds,
    loading,
    error,
    crdsByGroup,
    fetchCRDs,
    fetchCRDDetail,
    searchCRDs
  }
}

export function useCustomResource() {
  const { t } = useI18n()
  const clusterStore = useClusterStore()
  
  const resources = ref<CustomResourceItem[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Get custom resources list
  const fetchCustomResources = async (
    group: string,
    version: string,
    plural: string,
    namespace?: string
  ) => {
    loading.value = true
    error.value = null
    
    try {
      const clusterId = clusterStore.selectedClusterId || undefined
      const response = await crdApi.getCustomResources(group, version, plural, {
        namespace,
        clusterId
      })
      
      if (response.code === 200 || response.code === 0) {
        resources.value = response.data.items
      } else {
        error.value = response.message || t('crd.fetchResourcesFailed')
        ElMessage.error(error.value)
      }
    } catch (err) {
      error.value = t('crd.fetchResourcesFailed')
      ElMessage.error(error.value)
      console.error('Failed to fetch custom resources:', err)
    } finally {
      loading.value = false
    }
  }

  // 获取自定义资源详情
  const fetchCustomResource = async (
    group: string,
    version: string,
    plural: string,
    name: string,
    namespace?: string
  ): Promise<CustomResourceItem | null> => {
    try {
      const clusterId = clusterStore.selectedClusterId || undefined
      const response = await crdApi.getCustomResource(group, version, plural, name, {
        namespace,
        clusterId
      })
      
      if (response.code === 200 || response.code === 0) {
        return response.data
      } else {
        ElMessage.error(response.message || t('crd.fetchResourceDetailFailed'))
        return null
      }
    } catch (err) {
      ElMessage.error(t('crd.fetchResourceDetailFailed'))
      console.error('Failed to fetch custom resource:', err)
      return null
    }
  }

  // 创建自定义资源
  const createCustomResource = async (
    group: string,
    version: string,
    plural: string,
    resource: any,
    namespace?: string
  ): Promise<boolean> => {
    try {
      const clusterId = clusterStore.selectedClusterId || undefined
      const response = await crdApi.createCustomResource(group, version, plural, resource, {
        namespace,
        clusterId
      })
      
      if (response.code === 200 || response.code === 0) {
        ElMessage.success(t('crd.createResourceSuccess'))
        return true
      } else {
        ElMessage.error(response.message || t('crd.createResourceFailed'))
        return false
      }
    } catch (err) {
      ElMessage.error(t('crd.createResourceFailed'))
      console.error('Failed to create custom resource:', err)
      return false
    }
  }

  // 更新自定义资源
  const updateCustomResource = async (
    group: string,
    version: string,
    plural: string,
    name: string,
    resource: any,
    namespace?: string
  ): Promise<boolean> => {
    try {
      const clusterId = clusterStore.selectedClusterId || undefined
      const response = await crdApi.updateCustomResource(group, version, plural, name, resource, {
        namespace,
        clusterId
      })
      
      if (response.code === 200 || response.code === 0) {
        ElMessage.success(t('crd.updateResourceSuccess'))
        return true
      } else {
        ElMessage.error(response.message || t('crd.updateResourceFailed'))
        return false
      }
    } catch (err) {
      ElMessage.error(t('crd.updateResourceFailed'))
      console.error('Failed to update custom resource:', err)
      return false
    }
  }

  // 删除自定义资源
  const deleteCustomResource = async (
    group: string,
    version: string,
    plural: string,
    name: string,
    namespace?: string
  ): Promise<boolean> => {
    try {
      const clusterId = clusterStore.selectedClusterId || undefined
      const response = await crdApi.deleteCustomResource(group, version, plural, name, {
        namespace,
        clusterId
      })
      
      if (response.code === 200 || response.code === 0) {
        ElMessage.success(t('crd.deleteResourceSuccess'))
        return true
      } else {
        ElMessage.error(response.message || t('crd.deleteResourceFailed'))
        return false
      }
    } catch (err) {
      ElMessage.error(t('crd.deleteResourceFailed'))
      console.error('Failed to delete custom resource:', err)
      return false
    }
  }

  return {
    resources,
    loading,
    error,
    fetchCustomResources,
    fetchCustomResource,
    createCustomResource,
    updateCustomResource,
    deleteCustomResource
  }
}