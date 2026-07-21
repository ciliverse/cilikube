import { httpClient } from '@/utils/service'

// CRD related interface type definitions
export interface CRDItem {
  name: string
  group: string
  version: string
  kind: string
  plural: string
  singular: string
  scope: string
  categories?: string[]
  shortNames?: string[]
  createdAt: string
  labels?: Record<string, string>
  annotations?: Record<string, string>
}

export interface CRDListResponse {
  items: CRDItem[]
  total: number
}

export interface CRDVersion {
  name: string
  served: boolean
  storage: boolean
}

export interface CRDCondition {
  type: string
  status: string
  lastTransitionTime: string
  reason?: string
  message?: string
}

export interface CRDDetailResponse extends CRDItem {
  schema?: Record<string, any>
  versions?: CRDVersion[]
  conditions?: CRDCondition[]
  description?: string
}

export interface CustomResourceItem {
  name: string
  namespace?: string
  kind: string
  apiVersion: string
  createdAt: string
  labels?: Record<string, string>
  annotations?: Record<string, string>
  spec?: Record<string, any>
  status?: Record<string, any>
}

export interface CustomResourceListResponse {
  items: CustomResourceItem[]
  total: number
}

export interface CustomResourceRequest {
  apiVersion: string
  kind: string
  metadata: Record<string, any>
  spec?: Record<string, any>
}

// CRD API
export const crdApi = {
  // Get CRD list
  getCRDs: (clusterId?: string): Promise<ApiResponseData<CRDListResponse>> => {
    const params = clusterId ? { clusterId } : {}
    return httpClient.standard.get('/api/v1/crds', { params })
  },

  // Get CRD details
  getCRD: (name: string, clusterId?: string): Promise<ApiResponseData<CRDDetailResponse>> => {
    const params = clusterId ? { clusterId } : {}
    return httpClient.standard.get(`/api/v1/crds/definition/${name}`, { params })
  },

  // Get custom resources list
  getCustomResources: (
    group: string,
    version: string,
    plural: string,
    options?: {
      namespace?: string
      limit?: number
      continue?: string
      clusterId?: string
    }
  ): Promise<ApiResponseData<CustomResourceListResponse>> => {
    const params = { ...options }
    return httpClient.standard.get(`/api/v1/crds/resources/${group}/${version}/${plural}`, { params })
  },

  // Get custom resource details
  getCustomResource: (
    group: string,
    version: string,
    plural: string,
    name: string,
    options?: {
      namespace?: string
      clusterId?: string
    }
  ): Promise<ApiResponseData<CustomResourceItem>> => {
    const params = { ...options }
    return httpClient.standard.get(`/api/v1/crds/resources/${group}/${version}/${plural}/${name}`, { params })
  },

  // Create custom resource
  createCustomResource: (
    group: string,
    version: string,
    plural: string,
    resource: CustomResourceRequest,
    options?: {
      namespace?: string
      clusterId?: string
    }
  ): Promise<ApiResponseData<CustomResourceItem>> => {
    const params = { ...options }
    return httpClient.standard.post(`/api/v1/crds/resources/${group}/${version}/${plural}`, resource, { params })
  },

  // Update custom resource
  updateCustomResource: (
    group: string,
    version: string,
    plural: string,
    name: string,
    resource: CustomResourceRequest,
    options?: {
      namespace?: string
      clusterId?: string
    }
  ): Promise<ApiResponseData<CustomResourceItem>> => {
    const params = { ...options }
    return httpClient.standard.put(`/api/v1/crds/resources/${group}/${version}/${plural}/${name}`, resource, { params })
  },

  // Delete custom resource
  deleteCustomResource: (
    group: string,
    version: string,
    plural: string,
    name: string,
    options?: {
      namespace?: string
      clusterId?: string
    }
  ): Promise<ApiResponseData<null>> => {
    const params = { ...options }
    return httpClient.standard.delete(`/api/v1/crds/resources/${group}/${version}/${plural}/${name}`, { params })
  }
}