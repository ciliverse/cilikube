<template>
    <div class="deployment-management">
        <!-- 页面头部 -->
        <div class="page-header">
            <h1 class="page-title">{{ $t('deploymentManagement.title') }}</h1>
            <div class="header-actions">
                <el-button @click="fetchDeploymentData" :loading="loading.deployments" circle>
                    {{ $t('common.refresh') }}
                </el-button>
                <el-button type="primary" @click="handleAddDeployment"
                    :disabled="!selectedNamespace || selectedNamespace === '__ALL__'">
                    {{ $t('deploymentManagement.actions.create') }}
                </el-button>
            </div>
        </div>

        <!-- 命名空间选择提示 -->
        <el-alert v-if="!selectedNamespace && !loading.namespaces"
            :title="$t('deploymentManagement.messages.selectNamespace')" type="info" show-icon :closable="false"
            style="margin-bottom: 24px;" />

        <!-- 工具栏 -->
        <div class="toolbar" v-if="selectedNamespace">
            <div class="toolbar-left">
                <div class="search-filters">
                    <el-select v-model="selectedNamespace" :placeholder="$t('deploymentManagement.filterNamespace')"
                        filterable :loading="loading.namespaces" class="namespace-select" @change="handleNamespaceChange">
                        <el-option :label="$t('deploymentManagement.allNamespaces')" value="__ALL__" />
                        <el-option v-for="ns in namespaces" :key="ns" :label="ns" :value="ns" />
                    </el-select>
                    <el-input v-model="searchQuery" :placeholder="$t('deploymentManagement.searchPlaceholder')" clearable
                        class="search-input" @input="handleSearchDebounced" />
                    <el-select v-model="filterStatus" :placeholder="$t('deploymentManagement.filterStatus')" clearable
                        class="filter-select">
                        <el-option :label="$t('deploymentManagement.allStatuses')" value="" />
                        <el-option :label="$t('deploymentManagement.statuses.available')" value="available" />
                        <el-option :label="$t('deploymentManagement.statuses.progressing')" value="progressing" />
                        <el-option :label="$t('deploymentManagement.statuses.degraded')" value="degraded" />
                    </el-select>
                </div>
            </div>
            <div class="toolbar-right">
                <el-button-group class="view-toggle">
                    <el-button :type="viewMode === 'card' ? 'primary' : ''" @click="viewMode = 'card'">
                        {{ $t('deploymentManagement.viewToggle.card') }}
                    </el-button>
                    <el-button :type="viewMode === 'list' ? 'primary' : ''" @click="viewMode = 'list'">
                        {{ $t('deploymentManagement.viewToggle.list') }}
                    </el-button>
                </el-button-group>
            </div>
        </div>

        <!-- Deployment 列表 -->
        <div class="deployment-list" v-loading="loading.deployments">
            <!-- 卡片视图 -->
            <div v-if="filteredDeployments.length > 0 && viewMode === 'card'" class="deployment-grid">
                <div v-for="deployment in paginatedData" :key="deployment.uid" class="deployment-card"
                    @click="handleDeploymentClick(deployment)">
                    <div class="card-header">
                        <div class="deployment-info">
                            <div class="deployment-name-row">
                                <div class="status-dot" :class="getStatusClass(deployment.status)"></div>
                                <div class="deployment-name">{{ deployment.name }}</div>
                            </div>
                            <div class="deployment-namespace">{{ deployment.namespace }}</div>
                        </div>
                    </div>

                    <div class="card-body">
                        <div class="deployment-meta">
                            <div class="meta-item">
                                <span class="meta-label">{{ $t('deploymentManagement.columns.status') }}</span>
                                <span class="meta-value">{{ getStatusText(deployment.status) }}</span>
                            </div>
                            <div class="meta-item">
                                <span class="meta-label">{{ $t('deploymentManagement.columns.replicas') }}</span>
                                <span class="meta-value" :class="getReplicasClass(deployment)">
                                    {{ deployment.readyReplicas }} / {{ deployment.desiredReplicas }}
                                </span>
                            </div>
                            <div class="meta-item" v-if="deployment.images && deployment.images.length > 0">
                                <span class="meta-label">{{ $t('deploymentManagement.columns.images') }}</span>
                                <span class="meta-value">{{ deployment.images[0] }}</span>
                            </div>
                            <div class="meta-item">
                                <span class="meta-label">{{ $t('common.createdAt') }}</span>
                                <span class="meta-value">{{ formatTimestamp(deployment.createdAt) }}</span>
                            </div>
                        </div>
                    </div>

                    <div class="card-footer">
                        <div class="deployment-actions">
                            <el-button type="primary" size="small" plain @click.stop="handleScale(deployment)">
                                {{ $t('deploymentManagement.actions.scale') }}
                            </el-button>
                            <el-button type="success" size="small" plain @click.stop="handleViewPods(deployment)">
                                {{ $t('deploymentManagement.actions.viewPods') }}
                            </el-button>
                            <el-button type="danger" size="small" plain
                                @click.stop="handleDeleteDeployment(deployment)">
                                {{ $t('common.delete') }}
                            </el-button>
                        </div>
                    </div>
                </div>
            </div>

            <!-- 列表视图 -->
            <div v-else-if="filteredDeployments.length > 0 && viewMode === 'list'" class="deployment-table">
                <el-table :data="paginatedData" @row-click="handleDeploymentClick">
                    <el-table-column width="50" align="center">
                        <template #default="{ row }">
                            <div class="status-dot" :class="getStatusClass(row.status)"></div>
                        </template>
                    </el-table-column>
                    <el-table-column prop="name" :label="$t('deploymentManagement.columns.name')" min-width="150">
                        <template #default="{ row }">
                            <div class="table-deployment-name">{{ row.name }}</div>
                        </template>
                    </el-table-column>
                    <el-table-column prop="namespace" :label="$t('deploymentManagement.columns.namespace')"
                        min-width="100" />
                    <el-table-column :label="$t('deploymentManagement.columns.status')" min-width="100" align="center">
                        <template #default="{ row }">
                            <span>{{ getStatusText(row.status) }}</span>
                        </template>
                    </el-table-column>
                    <el-table-column :label="$t('deploymentManagement.columns.replicas')" min-width="100" align="center">
                        <template #default="{ row }">
                            <span :class="getReplicasClass(row)">{{ row.readyReplicas }} / {{ row.desiredReplicas
                                }}</span>
                        </template>
                    </el-table-column>
                    <el-table-column :label="$t('deploymentManagement.columns.images')" min-width="200"
                        show-overflow-tooltip>
                        <template #default="{ row }">
                            <div v-if="row.images && row.images.length > 0" class="images-container">
                                <el-tag v-for="(image, index) in row.images" :key="index" size="small" type="info"
                                    effect="plain" style="margin-right: 4px; margin-bottom: 2px;">
                                    {{ image }}
                                </el-tag>
                            </div>
                            <span v-else>-</span>
                        </template>
                    </el-table-column>
                    <el-table-column :label="$t('common.createdAt')" min-width="140">
                        <template #default="{ row }">
                            <span>{{ formatTimestamp(row.createdAt) }}</span>
                        </template>
                    </el-table-column>
                    <el-table-column :label="$t('common.actions')" min-width="180" align="center" fixed="right">
                        <template #default="{ row }">
                            <div class="table-actions">
                                <el-button type="primary" size="small" plain @click.stop="handleScale(row)">
                                    {{ $t('deploymentManagement.actions.scale') }}
                                </el-button>
                                <el-button type="success" size="small" plain @click.stop="handleViewPods(row)">
                                    {{ $t('deploymentManagement.actions.viewPods') }}
                                </el-button>
                                <el-button type="danger" size="small" plain @click.stop="handleDeleteDeployment(row)">
                                    {{ $t('common.delete') }}
                                </el-button>
                            </div>
                        </template>
                    </el-table-column>
                </el-table>
            </div>

            <div v-else-if="!loading.deployments && selectedNamespace" class="empty-state">
                <el-empty :description="searchQuery || filterStatus ?
                    $t('deploymentManagement.messages.noMatchingDeployments') :
                    $t('deploymentManagement.messages.noDeployments')" />
            </div>
        </div>

        <!-- 分页 -->
        <div class="pagination-container" v-if="filteredDeployments.length > pageSize">
            <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize"
                :page-sizes="[12, 24, 36, 48]" :total="filteredDeployments.length"
                layout="total, sizes, prev, pager, next, jumper" background @size-change="handleSizeChange"
                @current-change="handlePageChange" />
        </div>

        <!-- 创建 Deployment 对话框 -->
        <el-dialog v-model="yamlDialogConfig.visible" :title="$t('deploymentManagement.form.createFromYaml')"
            width="70%" :close-on-click-modal="false" @closed="handleYamlDialogClose" class="yaml-dialog">
            <div class="yaml-editor-container">
                <el-input v-model="yamlDialogConfig.content" type="textarea" :rows="20"
                    :placeholder="$t('deploymentManagement.form.yamlContent')" />
            </div>
            <template #footer>
                <div class="dialog-footer">
                    <el-button @click="yamlDialogConfig.visible = false" size="large">
                        {{ $t('common.cancel') }}
                    </el-button>
                    <el-button type="primary" @click="handleSaveYaml" :loading="yamlDialogConfig.saving" size="large">
                        {{ $t('common.confirm') }}
                    </el-button>
                </div>
            </template>
        </el-dialog>

        <!-- 扩缩容对话框 -->
        <el-dialog v-model="scaleDialogConfig.visible" :title="$t('deploymentManagement.form.scaleReplicas')"
            width="400px" :close-on-click-modal="false" @closed="handleScaleDialogClose" class="scale-dialog">
            <el-form :model="scaleForm" label-width="120px" ref="scaleFormRef">
                <el-form-item :label="$t('deploymentManagement.columns.name')">
                    <el-input :model-value="scaleForm.name" disabled />
                </el-form-item>
                <el-form-item :label="$t('deploymentManagement.columns.namespace')">
                    <el-input :model-value="scaleForm.namespace" disabled />
                </el-form-item>
                <el-form-item :label="$t('deploymentManagement.form.currentReplicas')">
                    <el-input :model-value="scaleForm.currentReplicas" disabled />
                </el-form-item>
                <el-form-item :label="$t('deploymentManagement.form.targetReplicas')" prop="replicas"
                    :rules="[{ required: true, message: $t('deploymentManagement.form.replicasRequired') }]">
                    <el-input-number v-model="scaleForm.replicas" :min="0" style="width: 100%" />
                </el-form-item>
            </el-form>
            <template #footer>
                <div class="dialog-footer">
                    <el-button @click="scaleDialogConfig.visible = false" size="large">
                        {{ $t('common.cancel') }}
                    </el-button>
                    <el-button type="primary" @click="handleConfirmScale" :loading="scaleDialogConfig.saving"
                        size="large">
                        {{ $t('common.confirm') }}
                    </el-button>
                </div>
            </template>
        </el-dialog>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, reactive, watch, nextTick } from "vue"
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from "element-plus"
import { kubernetesRequest, fetchNamespaces, fetchDeployments } from "@/utils/api-config"
import { useClusterStore } from "@/store/modules/clusterStore"
import dayjs from "dayjs"
import { debounce } from 'lodash-es'
import type { FormInstance } from 'element-plus'

const { t } = useI18n()

// 集群 store
const clusterStore = useClusterStore()

// 接口定义
interface DeploymentDisplayItem {
    uid: string
    name: string
    namespace: string
    status: string
    desiredReplicas: number
    availableReplicas: number
    readyReplicas: number
    updatedReplicas: number
    createdAt: string
    images: string[]
    rawData: any
}

interface ScaleForm {
    name: string
    namespace: string
    currentReplicas: number
    replicas: number
}

// 响应式状态
const allDeployments = ref<DeploymentDisplayItem[]>([])
const namespaces = ref<string[]>([])
const selectedNamespace = ref<string>("")
const currentPage = ref(1)
const pageSize = ref(12)
const searchQuery = ref("")
const filterStatus = ref("")
const viewMode = ref<'card' | 'list'>('card')

const loading = reactive({
    namespaces: false,
    deployments: false
})

// 对话框状态
const yamlDialogConfig = reactive({
    visible: false,
    content: '',
    saving: false
})

const scaleDialogConfig = reactive({
    visible: false,
    saving: false
})

const scaleForm = reactive<ScaleForm>({
    name: '',
    namespace: '',
    currentReplicas: 0,
    replicas: 0
})

// 组件引用
const scaleFormRef = ref<FormInstance>()

// 过滤后的 Deployment 列表
const filteredDeployments = computed(() => {
    return allDeployments.value.filter(deployment => {
        // 搜索过滤
        const matchesSearch = !searchQuery.value ||
            deployment.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
            deployment.namespace.toLowerCase().includes(searchQuery.value.toLowerCase())

        // 状态过滤
        const matchesStatus = !filterStatus.value ||
            deployment.status.toLowerCase() === filterStatus.value.toLowerCase()

        return matchesSearch && matchesStatus
    })
})

// 当前页数据
const paginatedData = computed(() => {
    const start = (currentPage.value - 1) * pageSize.value
    const end = start + pageSize.value
    return filteredDeployments.value.slice(start, end)
})

// 获取命名空间列表
const fetchNamespacesList = async () => {
    loading.namespaces = true
    try {
        const namespaceList = await fetchNamespaces()
        namespaces.value = namespaceList
        if (namespaceList.length > 0 && !selectedNamespace.value) {
            selectedNamespace.value = '__ALL__'
        }
    } catch (error: any) {
        ElMessage.error(error.message)
        namespaces.value = []
        selectedNamespace.value = ""
        allDeployments.value = []
    } finally {
        loading.namespaces = false
    }
}

// 获取 Deployment 数据
const fetchDeploymentData = async () => {
    if (!selectedNamespace.value) {
        allDeployments.value = []
        return
    }

    loading.deployments = true
    try {
        let response
        if (selectedNamespace.value === '__ALL__') {
            // 获取所有命名空间的 Deployment
            response = await fetchDeployments()
        } else {
            // 获取指定命名空间的 Deployment
            response = await fetchDeployments(selectedNamespace.value)
        }

        // 更强大的数据验证和处理
        if (response && Array.isArray(response.items)) {
            const processedDeployments = response.items
                .filter((item: any) => {
                    // 过滤掉无效的items
                    return item && item.metadata && item.metadata.name
                })
                .map((deployment: any, index: number) => {
                    try {
                        return {
                            uid: deployment.metadata?.uid || `${deployment.metadata?.namespace}-${deployment.metadata?.name}-${index}`,
                            name: deployment.metadata?.name || 'Unknown',
                            namespace: deployment.metadata?.namespace || 'default',
                            status: getDeploymentStatus(deployment),
                            desiredReplicas: deployment.spec?.replicas ?? 0,
                            availableReplicas: deployment.status?.availableReplicas ?? 0,
                            readyReplicas: deployment.status?.readyReplicas ?? 0,
                            updatedReplicas: deployment.status?.updatedReplicas ?? 0,
                            createdAt: deployment.metadata?.creationTimestamp || new Date().toISOString(),
                            images: extractImages(deployment.spec),
                            rawData: deployment
                        }
                    } catch (itemError) {
                        console.warn(`Failed to process deployment item:`, itemError, deployment)
                        return null
                    }
                })
                .filter((item: any) => item !== null) as DeploymentDisplayItem[]
            
            allDeployments.value = processedDeployments

            console.log(`Successfully loaded ${allDeployments.value.length} deployments`)

            // 调整页码
            nextTick(() => {
                const totalPages = Math.ceil(filteredDeployments.value.length / pageSize.value)
                if (currentPage.value > totalPages && totalPages > 0) {
                    currentPage.value = totalPages
                }
            })
        } else {
            console.warn('Invalid response structure:', response)
            ElMessage.error(`${t('deploymentManagement.messages.loadingDeployments')}失败: 响应数据格式不正确`)
            allDeployments.value = []
        }
    } catch (error: any) {
        console.error("获取 Deployment 数据失败:", error)
        const errorMsg = error?.response?.data?.message || error.message || '网络请求失败'
        ElMessage.error(`${t('deploymentManagement.messages.loadingDeployments')}出错: ${errorMsg}`)
        allDeployments.value = []
    } finally {
        loading.deployments = false
    }
}

// 辅助函数
const getDeploymentStatus = (deployment: any): string => {
    const desired = deployment.spec?.replicas ?? 0
    const available = deployment.status?.availableReplicas ?? 0
    const ready = deployment.status?.readyReplicas ?? 0
    const updated = deployment.status?.updatedReplicas ?? 0

    // 检查Conditions以获得更准确的状态
    const conditions = deployment.status?.conditions || []
    const availableCondition = conditions.find((c: any) => c.type === 'Available')
    const progressingCondition = conditions.find((c: any) => c.type === 'Progressing')

    // 如果所有副本都就绪且可用
    if (desired > 0 && ready === desired && available === desired && updated === desired) {
        return 'available'
    }

    // 如果Progressing条件为False，表示部署失败
    if (progressingCondition?.status === 'False') {
        return 'degraded'
    }

    // 如果副本数不匹配
    if (desired > 0 && (available < desired || ready < desired)) {
        // 检查是否在更新中
        if (progressingCondition?.status === 'True') {
            return 'progressing'
        }
        return 'degraded'
    }

    // 如果没有desired replicas
    if (desired === 0) {
        return 'available'
    }

    return 'unknown'
}

const getStatusClass = (status: string) => {
    const lowerStatus = status.toLowerCase()
    if (lowerStatus === 'available') return 'status-healthy'
    if (lowerStatus === 'degraded' || lowerStatus === 'replicafailure') return 'status-unhealthy'
    if (lowerStatus === 'progressing') return 'status-progressing'
    return 'status-unknown'
}

const getStatusText = (status: string) => {
    const statusMap: Record<string, string> = {
        available: t('deploymentManagement.statuses.available'),
        progressing: t('deploymentManagement.statuses.progressing'),
        degraded: t('deploymentManagement.statuses.degraded'),
        replicafailure: t('deploymentManagement.statuses.replicaFailure'),
        unknown: t('deploymentManagement.statuses.unknown')
    }
    return statusMap[status.toLowerCase()] || status
}

const getReplicasClass = (deployment: DeploymentDisplayItem): string => {
    if (deployment.readyReplicas === deployment.desiredReplicas && deployment.desiredReplicas > 0) return 'replicas-ok'
    if (deployment.desiredReplicas > 0 && deployment.readyReplicas < deployment.desiredReplicas) return 'replicas-warning'
    if (deployment.desiredReplicas === 0) return 'replicas-scaled-down'
    return ''
}

const extractImages = (spec: any): string[] => {
    if (!spec) return []
    
    const images: Set<string> = new Set()
    
    // 从主容器中获取镜像
    if (spec.template?.spec?.containers) {
        spec.template.spec.containers.forEach((c: any) => {
            if (c.image) {
                images.add(c.image)
            }
        })
    }
    
    // 从initContainers中获取镜像
    if (spec.template?.spec?.initContainers) {
        spec.template.spec.initContainers.forEach((c: any) => {
            if (c.image) {
                images.add(c.image)
            }
        })
    }
    
    return Array.from(images)
}

const formatTimestamp = (timestamp: string | Date | undefined): string => {
    if (!timestamp) return 'N/A'
    
    try {
        let date
        if (typeof timestamp === 'string') {
            // 尝试多种时间戳格式
            date = dayjs(timestamp)
        } else if (timestamp instanceof Date) {
            date = dayjs(timestamp)
        } else {
            return 'Invalid Date'
        }
        
        if (!date.isValid()) {
            return 'Invalid Date'
        }
        
        return date.format("YYYY-MM-DD HH:mm:ss")
    } catch (error) {
        console.warn('Failed to format timestamp:', timestamp, error)
        return 'N/A'
    }
}

// 事件处理
const handleNamespaceChange = () => {
    currentPage.value = 1
    searchQuery.value = ''
    filterStatus.value = ''
    fetchDeploymentData()
}

const handlePageChange = (page: number) => {
    currentPage.value = page
}

const handleSizeChange = (size: number) => {
    pageSize.value = size
    currentPage.value = 1
}

const handleSearchDebounced = debounce(() => {
    currentPage.value = 1
}, 300)

const handleDeploymentClick = (deployment: DeploymentDisplayItem) => {
    console.log('Deployment clicked:', deployment)
}

// Deployment 操作
const handleAddDeployment = () => {
    if (!selectedNamespace.value || selectedNamespace.value === '__ALL__') {
        ElMessage.warning(t('deploymentManagement.messages.selectNamespace'))
        return
    }
    yamlDialogConfig.content = getPlaceholderYaml()
    yamlDialogConfig.visible = true
}

const getPlaceholderYaml = () => {
    return `apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
  namespace: ${selectedNamespace.value || 'default'}
  labels:
    app: myapp
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
      - name: main-container
        image: nginx:alpine
        ports:
        - containerPort: 80
        resources:
          requests:
            memory: "64Mi"
            cpu: "100m"
          limits:
            memory: "128Mi"
            cpu: "200m"`
}

const handleSaveYaml = async () => {
    if (!selectedNamespace.value || selectedNamespace.value === '__ALL__') {
        ElMessage.error(t('deploymentManagement.form.namespaceRequired'))
        return
    }

    yamlDialogConfig.saving = true
    try {
        const response = await kubernetesRequest({
            url: `/api/v1/namespaces/${selectedNamespace.value}/deployments`,
            method: 'POST',
            headers: { 'Content-Type': 'application/yaml' },
            data: yamlDialogConfig.content
        })

        if (response.code === 200 || response.code === 201) {
            ElMessage.success(t('deploymentManagement.messages.createSuccess'))
            yamlDialogConfig.visible = false
            fetchDeploymentData()
        } else {
            ElMessage.error(`${t('deploymentManagement.messages.createFailed')}: ${response.message || '未知错误'}`)
        }
    } catch (error: any) {
        console.error("应用 YAML 失败:", error)
        const errMsg = error.response?.data?.message || error.message || '请求失败'
        ElMessage.error(`${t('deploymentManagement.messages.createFailed')}: ${errMsg}`)
    } finally {
        yamlDialogConfig.saving = false
    }
}

const handleYamlDialogClose = () => {
    yamlDialogConfig.content = ''
}

const handleScale = (deployment: DeploymentDisplayItem) => {
    Object.assign(scaleForm, {
        name: deployment.name,
        namespace: deployment.namespace,
        currentReplicas: deployment.desiredReplicas,
        replicas: deployment.desiredReplicas
    })
    scaleDialogConfig.visible = true
}

const handleConfirmScale = async () => {
    if (!scaleFormRef.value) return

    const valid = await scaleFormRef.value.validate().catch(() => false)
    if (!valid) return

    scaleDialogConfig.saving = true
    try {
        // 使用 PATCH 方法进行扩缩容
        const patchData = {
            spec: {
                replicas: scaleForm.replicas
            }
        }

        const response = await kubernetesRequest({
            url: `/api/v1/namespaces/${scaleForm.namespace}/deployments/${scaleForm.name}`,
            method: 'PATCH',
            headers: { 'Content-Type': 'application/merge-patch+json' },
            data: patchData
        })

        if (response.code === 200) {
            ElMessage.success(t('deploymentManagement.messages.scaleSuccess'))
            scaleDialogConfig.visible = false
            fetchDeploymentData()
        } else {
            ElMessage.error(`${t('deploymentManagement.messages.scaleFailed')}: ${response.message || '未知错误'}`)
        }
    } catch (error: any) {
        ElMessage.error(`${t('deploymentManagement.messages.scaleFailed')}: ${error.message}`)
    } finally {
        scaleDialogConfig.saving = false
    }
}

const handleScaleDialogClose = () => {
    Object.assign(scaleForm, {
        name: '',
        namespace: '',
        currentReplicas: 0,
        replicas: 0
    })
}

const handleViewPods = (deployment: DeploymentDisplayItem) => {
    ElMessage.info(`查看 ${deployment.name} 的 Pods`)
}

const handleDeleteDeployment = (deployment: DeploymentDisplayItem) => {
    ElMessageBox.confirm(
        t('deploymentManagement.messages.deleteConfirm', { name: deployment.name }),
        t('common.confirm'),
        {
            confirmButtonText: t('common.confirm'),
            cancelButtonText: t('common.cancel'),
            type: 'warning'
        }
    ).then(async () => {
        try {
            const response = await kubernetesRequest({
                url: `/api/v1/namespaces/${deployment.namespace}/deployments/${deployment.name}`,
                method: "DELETE"
            })

            if (response.code === 200 || response.code === 202) {
                ElMessage.success(t('deploymentManagement.messages.deleteSuccess'))
                await fetchDeploymentData()
            } else {
                ElMessage.error(`${t('deploymentManagement.messages.deleteFailed')}: ${response.message || '未知错误'}`)
            }
        } catch (error: any) {
            console.error("删除 Deployment 失败:", error)
            const errMsg = error.response?.data?.message || error.message || '请求失败'
            ElMessage.error(`${t('deploymentManagement.messages.deleteFailed')}: ${errMsg}`)
        }
    }).catch(() => {
        ElMessage.info('删除操作已取消')
    })
}

// 监听命名空间变化
watch(() => selectedNamespace.value, (newNamespace) => {
    if (newNamespace) {
        currentPage.value = 1
        fetchDeploymentData()
    } else {
        allDeployments.value = []
    }
}, { immediate: false })

// 组件挂载
onMounted(async () => {
    // 初始化集群选择
    if (clusterStore.availableClusters.length === 0) {
        try {
            await clusterStore.fetchAvailableClusters()
        } catch (error) {
            console.error('Failed to fetch clusters:', error)
        }
    }
    
    // 获取命名空间列表
    await fetchNamespacesList()
    
    // 如果成功获取命名空间，则获取deployment数据
    if (selectedNamespace.value) {
        await fetchDeploymentData()
    }
})
</script>

<style scoped>
.deployment-management {
    padding: 24px;
    background: linear-gradient(135deg, #ffffff 0%, #f8faff 100%);
    min-height: 100vh;
}

/* 页面头部 */
.page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
}

.page-title {
    font-size: 32px;
    font-weight: 600;
    color: #1f2937;
    margin: 0;
}

.header-actions {
    display: flex;
    gap: 12px;
    align-items: center;
}

/* 工具栏 */
.toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
    padding: 16px 20px;
    background: white;
    border-radius: 12px;
    box-shadow: 0 2px 8px rgba(30, 58, 138, 0.08);
    border: 1px solid #e1e8f0;
    gap: 16px;
    flex-wrap: wrap;
}

.toolbar-left {
    flex: 1;
    min-width: 300px;
}

.search-filters {
    display: flex;
    gap: 12px;
    align-items: center;
    flex-wrap: wrap;
}

.namespace-select {
    min-width: 180px;
    flex-shrink: 0;
}

.search-input {
    min-width: 200px;
    flex: 1;
}

.filter-select {
    min-width: 140px;
    flex-shrink: 0;
}

.toolbar-right {
    display: flex;
    gap: 16px;
    align-items: center;
    flex-shrink: 0;
}

.view-toggle {
    border-radius: 8px;
}

/* Deployment 列表 */
.deployment-list {
    min-height: 400px;
}

.deployment-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(380px, 1fr));
    gap: 20px;
}

.deployment-card {
    background: white;
    border-radius: 12px;
    border: 1px solid #e1e8f0;
    transition: all 0.2s ease;
    cursor: pointer;
    overflow: hidden;
}

.deployment-card:hover {
    box-shadow: 0 4px 12px rgba(30, 58, 138, 0.15);
    border-color: #3b82f6;
}

.card-header {
    padding: 20px 20px 16px;
    border-bottom: 1px solid #f1f5f9;
}

.deployment-info {
    margin-bottom: 12px;
}

.deployment-name {
    font-size: 18px;
    font-weight: 600;
    color: #1f2937;
    margin-bottom: 4px;
    word-break: break-word;
    overflow-wrap: break-word;
    white-space: normal;
}

.deployment-namespace {
    font-size: 13px;
    color: #6b7280;
    font-family: 'Maple Mono', 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
    font-weight: 600;
    margin-top: 4px;
}

.card-body {
    padding: 16px 20px;
}

.deployment-meta {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.meta-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.meta-label {
    font-size: 13px;
    color: #64748b;
    font-weight: 500;
}

.meta-value {
    font-size: 13px;
    color: #1f2937;
    font-family: 'Maple Mono', 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
    font-weight: 600;
}

.card-footer {
    padding: 16px 20px;
    border-top: 1px solid #f1f5f9;
    background: white;
}

.deployment-actions {
    display: flex;
    gap: 8px;
    justify-content: flex-end;
}

/* 状态指示器 */
.status-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    flex-shrink: 0;
    transition: all 0.2s ease;
}

.status-healthy {
    background: #10b981;
    box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.2);
}

.status-unhealthy {
    background: #ef4444;
    box-shadow: 0 0 0 2px rgba(239, 68, 68, 0.2);
}

.status-unknown {
    background: #64748b;
    box-shadow: 0 0 0 2px rgba(100, 116, 139, 0.2);
}

/* Deployment 名称行 */
.deployment-name-row {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 8px;
}

/* 副本数状态 */
.replicas-ok {
    color: #10b981;
    font-weight: 600;
}

.replicas-warning {
    color: #f59e0b;
    font-weight: 600;
}

.replicas-scaled-down {
    color: #64748b;
}

/* 列表视图 */
.deployment-table {
    background: white;
    border-radius: 12px;
    overflow: auto;
    border: 1px solid #e1e8f0;
    width: 100%;
}

.table-deployment-name {
    font-weight: 500;
    color: #1f2937;
    word-break: break-word;
    overflow-wrap: break-word;
    white-space: normal;
}

.table-actions {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
    justify-content: center;
}

.images-container {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
}

.empty-state {
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 80px 20px;
    background: white;
    border-radius: 12px;
    border: 1px solid #e1e8f0;
}

/* 分页 */
.pagination-container {
    display: flex;
    justify-content: flex-end;
    margin-top: 24px;
}

/* 对话框样式 */
.yaml-dialog :deep(.el-dialog__body) {
    padding: 20px 24px;
}

.yaml-editor-container {
    border: 1px solid #e1e8f0;
    border-radius: 8px;
    overflow: hidden;
}

.yaml-editor-container :deep(.el-textarea__inner) {
    font-family: 'Maple Mono', 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
    font-size: 12px;
    line-height: 1.5;
    border: none;
    border-radius: 0;
}

.dialog-footer {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    padding: 16px 0 0;
    border-top: 1px solid #e1e8f0;
}

/* 蓝白主题 - 按钮边框样式 */
:deep(.el-button--primary) {
    background-color: transparent;
    border-color: #3b82f6;
    color: #3b82f6;
}

:deep(.el-button--primary:hover) {
    background-color: #3b82f6;
    border-color: #3b82f6;
    color: #ffffff;
}

:deep(.el-button--success) {
    background-color: transparent;
    border-color: #10b981;
    color: #10b981;
}

:deep(.el-button--success:hover) {
    background-color: #10b981;
    border-color: #10b981;
    color: #ffffff;
}

:deep(.el-button--danger) {
    background-color: transparent;
    border-color: #ef4444;
    color: #ef4444;
}

:deep(.el-button--danger:hover) {
    background-color: #ef4444;
    border-color: #ef4444;
    color: #ffffff;
}

/* 其他组件样式 */
:deep(.el-table) {
    border: 1px solid #e1e8f0;
    border-radius: 8px;
    overflow: hidden;
}

:deep(.el-table th) {
    background-color: #f8faff;
    color: #1f2937;
    font-weight: 600;
}

:deep(.el-table--striped .el-table__body tr.el-table__row--striped td) {
    background-color: #f8faff;
}

:deep(.el-select .el-input__wrapper) {
    border-color: #d1d5db;
}

:deep(.el-select .el-input__wrapper:hover) {
    border-color: #3b82f6;
}

:deep(.el-input__wrapper) {
    border-color: #d1d5db;
}

:deep(.el-input__wrapper:hover) {
    border-color: #3b82f6;
}

:deep(.el-pagination .el-pager li.is-active) {
    background-color: #3b82f6;
    color: #ffffff;
}

:deep(.el-dialog) {
    border-radius: 8px;
}

:deep(.el-dialog__header) {
    background-color: #f8faff;
    border-bottom: 1px solid #e1e8f0;
}

:deep(.el-dialog__title) {
    color: #1f2937;
    font-weight: 600;
}

/* 响应式设计 */
@media (max-width: 1200px) {
    .deployment-grid {
        grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
    }
}

@media (max-width: 768px) {
    .deployment-management {
        padding: 16px;
    }

    .page-header {
        flex-direction: column;
        align-items: flex-start;
        gap: 16px;
    }

    .page-title {
        font-size: 24px;
    }

    .toolbar {
        flex-direction: column;
        align-items: stretch;
        gap: 12px;
    }

    .toolbar-left,
    .toolbar-right {
        width: 100%;
    }

    .search-filters {
        flex-direction: column;
        align-items: stretch;
    }

    .namespace-select,
    .search-input,
    .filter-select {
        width: 100%;
    }

    .toolbar-right {
        justify-content: stretch;
    }

    .view-toggle {
        width: 100%;
    }

    .deployment-grid {
        grid-template-columns: 1fr;
    }

    .meta-item {
        flex-direction: column;
        align-items: flex-start;
        gap: 4px;
    }

    :deep(.el-table) {
        font-size: 12px;
    }

    :deep(.el-table th),
    :deep(.el-table td) {
        padding: 8px 0;
    }
}
</style>