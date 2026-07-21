<template>
    <div class="pod-management">
        <!-- 页面头部 -->
        <div class="page-header">
            <h1 class="page-title">{{ $t('podManagement.title') }}</h1>
            <div class="header-actions">
                <el-button :icon="Refresh" @click="fetchPodData" :loading="loading.pods" circle />
                <el-button type="primary" :icon="Plus" @click="handleAddPod" :disabled="!selectedNamespace">
                    {{ $t('podManagement.actions.create') }}
                </el-button>
            </div>
        </div>

        <!-- 命名空间选择提示 -->
        <el-alert v-if="!selectedNamespace && !loading.namespaces" :title="$t('podManagement.messages.selectNamespace')"
            type="info" show-icon :closable="false" style="margin-bottom: 24px;" />

        <!-- 工具栏 -->
        <div class="toolbar" v-if="selectedNamespace">
            <div class="search-filters">
                <el-select v-model="selectedNamespace" :placeholder="$t('podManagement.filterNamespace')" filterable
                    :loading="loading.namespaces" class="namespace-select">
                    <el-option :label="$t('podManagement.allNamespaces')" value="__ALL__" />
                    <el-option v-for="ns in namespaces" :key="ns" :label="ns" :value="ns" />
                </el-select>
                <el-input v-model="searchQuery" :placeholder="$t('podManagement.searchPlaceholder')"
                    :prefix-icon="Search" clearable class="search-input" @input="handleSearchDebounced" />
                <el-select v-model="filterStatus" :placeholder="$t('podManagement.filterStatus')" clearable
                    class="filter-select">
                    <el-option :label="$t('podManagement.allStatuses')" value="" />
                    <el-option :label="$t('podManagement.statuses.running')" value="running" />
                    <el-option :label="$t('podManagement.statuses.pending')" value="pending" />
                    <el-option :label="$t('podManagement.statuses.failed')" value="failed" />
                    <el-option :label="$t('podManagement.statuses.succeeded')" value="succeeded" />
                </el-select>
            </div>
            <div class="toolbar-right">
                <el-button-group class="view-toggle">
                    <el-button :type="viewMode === 'card' ? 'primary' : ''" :icon="Grid" @click="viewMode = 'card'" />
                    <el-button :type="viewMode === 'list' ? 'primary' : ''" :icon="List" @click="viewMode = 'list'" />
                </el-button-group>
            </div>
        </div>

        <!-- Pod 列表 -->
        <div class="pod-list" v-loading="loading.pods">
            <!-- 卡片视图 -->
            <div v-if="filteredPods.length > 0 && viewMode === 'card'" class="pod-grid">
                <div v-for="pod in paginatedData" :key="pod.uid" class="pod-card" @click="handlePodClick(pod)">
                    <div class="card-header">
                        <div class="pod-info">
                            <div class="pod-name-row">
                                <div class="status-dot" :class="getStatusClass(pod.status)"></div>
                                <div class="pod-name">{{ pod.name }}</div>
                            </div>
                            <div class="pod-namespace">{{ pod.namespace }}</div>
                        </div>
                    </div>

                    <div class="card-body">
                        <div class="pod-meta">
                            <div class="meta-item">
                                <span class="meta-label">{{ $t('podManagement.columns.status') }}</span>
                                <span class="meta-value">{{ getStatusText(pod.status) }}</span>
                            </div>
                            <div class="meta-item" v-if="pod.ip">
                                <span class="meta-label">{{ $t('podManagement.columns.ip') }}</span>
                                <span class="meta-value">{{ pod.ip }}</span>
                            </div>
                            <div class="meta-item" v-if="pod.node">
                                <span class="meta-label">{{ $t('podManagement.columns.node') }}</span>
                                <span class="meta-value">{{ pod.node }}</span>
                            </div>
                            <div class="meta-item">
                                <span class="meta-label">{{ $t('podManagement.columns.created') }}</span>
                                <span class="meta-value">{{ formatTimestamp(pod.createdAt) }}</span>
                            </div>
                        </div>
                    </div>

                    <div class="card-footer">
                        <div class="pod-actions">
                            <el-button type="primary" size="small" @click.stop="handleViewLogs(pod)">
                                {{ $t('podManagement.actions.viewLogs') }}
                            </el-button>
                            <el-button type="success" size="small" @click.stop="handleExec(pod)">
                                {{ $t('podManagement.actions.exec') }}
                            </el-button>
                            <el-button type="danger" size="small" @click.stop="handleDeletePod(pod)">
                                {{ $t('podManagement.actions.delete') }}
                            </el-button>
                        </div>
                    </div>
                </div>
            </div>

            <!-- 列表视图 -->
            <div v-else-if="filteredPods.length > 0 && viewMode === 'list'" class="pod-table">
                <el-table :data="paginatedData" @row-click="handlePodClick">
                    <el-table-column width="60">
                        <template #default="{ row }">
                            <div class="status-dot" :class="getStatusClass(row.status)"></div>
                        </template>
                    </el-table-column>
                    <el-table-column prop="name" :label="$t('podManagement.columns.name')" min-width="200">
                        <template #default="{ row }">
                            <div class="table-pod-name">{{ row.name }}</div>
                        </template>
                    </el-table-column>
                    <el-table-column prop="namespace" :label="$t('podManagement.columns.namespace')" width="150" />
                    <el-table-column :label="$t('podManagement.columns.status')" width="120">
                        <template #default="{ row }">
                            <span>{{ getStatusText(row.status) }}</span>
                        </template>
                    </el-table-column>
                    <el-table-column prop="ip" :label="$t('podManagement.columns.ip')" width="140">
                        <template #default="{ row }">
                            <code class="ip-address" v-if="row.ip">{{ row.ip }}</code>
                            <span v-else>-</span>
                        </template>
                    </el-table-column>
                    <el-table-column prop="node" :label="$t('podManagement.columns.node')" min-width="180"
                        show-overflow-tooltip />
                    <el-table-column :label="$t('podManagement.columns.created')" width="180">
                        <template #default="{ row }">
                            <span>{{ formatTimestamp(row.createdAt) }}</span>
                        </template>
                    </el-table-column>
                    <el-table-column :label="$t('common.actions')" width="240" fixed="right">
                        <template #default="{ row }">
                            <div class="table-actions">
                                <el-button type="primary" size="small" @click.stop="handleViewLogs(row)">
                                    {{ $t('podManagement.actions.viewLogs') }}
                                </el-button>
                                <el-button type="success" size="small" @click.stop="handleExec(row)">
                                    {{ $t('podManagement.actions.exec') }}
                                </el-button>
                                <el-button type="danger" size="small" @click.stop="handleDeletePod(row)">
                                    {{ $t('podManagement.actions.delete') }}
                                </el-button>
                            </div>
                        </template>
                    </el-table-column>
                </el-table>
            </div>

            <div v-else-if="!loading.pods && selectedNamespace" class="empty-state">
                <el-empty :description="searchQuery || filterStatus ?
                    $t('podManagement.messages.noMatchingPods') :
                    $t('podManagement.messages.noPods')" />
            </div>
        </div>

        <!-- 分页 -->
        <div class="pagination-container" v-if="filteredPods.length > pageSize">
            <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize"
                :page-sizes="[12, 24, 36, 48]" :total="filteredPods.length"
                layout="total, sizes, prev, pager, next, jumper" background @size-change="handleSizeChange"
                @current-change="handlePageChange" />
        </div>

        <!-- 创建 Pod 对话框 -->
        <el-dialog v-model="yamlDialogConfig.visible" :title="$t('podManagement.form.createFromYaml')" width="70%"
            :close-on-click-modal="false" @closed="handleYamlDialogClose" class="yaml-dialog">
            <div class="yaml-editor-container">
                <el-input v-model="yamlDialogConfig.content" type="textarea" :rows="20"
                    :placeholder="$t('podManagement.form.yamlContent')" />
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

        <!-- 查看日志对话框 -->
        <el-dialog v-model="logDialogConfig.visible" :title="$t('podManagement.actions.viewLogs')" width="80%"
            :close-on-click-modal="false" @closed="handleLogDialogClose" class="log-dialog">
            <div style="height: 70vh;">
                <PodLogs v-if="logDialogConfig.targetPod" :namespace="logDialogConfig.targetPod.namespace"
                    :pod-name="logDialogConfig.targetPod.name" :containers="logDialogConfig.containers"
                    ref="podLogsRef" />
            </div>
            <template #footer>
                <div class="dialog-footer">
                    <el-button @click="logDialogConfig.visible = false" size="large">
                        {{ $t('common.cancel') }}
                    </el-button>
                </div>
            </template>
        </el-dialog>

        <!-- 进入容器对话框 -->
        <el-dialog v-model="execDialogConfig.visible" :title="$t('podManagement.actions.exec')" width="80%"
            :close-on-click-modal="false" @closed="handleExecDialogClose" class="exec-dialog">
            <div style="height: 70vh;">
                <PodTerminal v-if="execDialogConfig.targetPod" :namespace="execDialogConfig.targetPod.namespace"
                    :pod-name="execDialogConfig.targetPod.name" :containers="execDialogConfig.containers"
                    ref="podTerminalRef" />
            </div>
            <template #footer>
                <div class="dialog-footer">
                    <el-button @click="execDialogConfig.visible = false" size="large">
                        {{ $t('common.cancel') }}
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
import { kubernetesRequest, fetchNamespaces } from "@/utils/api-config"
import dayjs from "dayjs"
import { debounce } from 'lodash-es'
import {
    Plus,
    Search,
    Refresh,
    Grid,
    List
} from '@element-plus/icons-vue'
import PodLogs from '@/components/PodLogs.vue'
import PodTerminal from '@/components/PodTerminal.vue'

const { t } = useI18n()

// 接口定义
interface PodDisplayItem {
    uid: string
    name: string
    namespace: string
    status: string
    reason?: string
    message?: string
    ip?: string
    node?: string
    createdAt: string
    labels?: { [key: string]: string }
}

interface K8sContainer {
    name: string
    image: string
}

interface PodDetail {
    name: string
    namespace: string
    uid: string
    spec: {
        containers: K8sContainer[]
        initContainers?: K8sContainer[]
    }
    status?: {
        phase?: string
        containerStatuses?: Array<{ name: string; ready: boolean; state: any }>
    }
}

// 响应式状态
const allPods = ref<PodDisplayItem[]>([])
const namespaces = ref<string[]>([])
const selectedNamespace = ref<string>("")
const currentPage = ref(1)
const pageSize = ref(12)
const searchQuery = ref("")
const filterStatus = ref("")
const viewMode = ref<'card' | 'list'>('card')

const loading = reactive({
    namespaces: false,
    pods: false
})

// 对话框状态
const yamlDialogConfig = reactive({
    visible: false,
    content: '',
    saving: false
})

const logDialogConfig = reactive({
    visible: false,
    targetPod: null as PodDisplayItem | null,
    containers: [] as K8sContainer[]
})

const execDialogConfig = reactive({
    visible: false,
    targetPod: null as PodDisplayItem | null,
    containers: [] as K8sContainer[]
})

// 组件引用
const podLogsRef = ref()
const podTerminalRef = ref()

// 过滤后的 Pod 列表
const filteredPods = computed(() => {
    return allPods.value.filter(pod => {
        // 搜索过滤
        const matchesSearch = !searchQuery.value ||
            pod.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
            (pod.ip && pod.ip.toLowerCase().includes(searchQuery.value.toLowerCase())) ||
            (pod.node && pod.node.toLowerCase().includes(searchQuery.value.toLowerCase())) ||
            pod.status.toLowerCase().includes(searchQuery.value.toLowerCase())

        // 状态过滤
        const matchesStatus = !filterStatus.value ||
            pod.status.toLowerCase() === filterStatus.value.toLowerCase()

        return matchesSearch && matchesStatus
    })
})

// 当前页数据
const paginatedData = computed(() => {
    const start = (currentPage.value - 1) * pageSize.value
    const end = start + pageSize.value
    return filteredPods.value.slice(start, end)
})

// 获取命名空间列表
const fetchNamespacesList = async () => {
    loading.namespaces = true
    try {
        const namespaceList = await fetchNamespaces()
        namespaces.value = namespaceList
        if (namespaceList.length > 0 && !selectedNamespace.value) {
            // 默认选择 "All" 选项来查看所有命名空间的 Pod
            selectedNamespace.value = '__ALL__'
        }
    } catch (error: any) {
        ElMessage.error(error.message)
        namespaces.value = []
        selectedNamespace.value = ""
        allPods.value = []
    } finally {
        loading.namespaces = false
    }
}

// 获取 Pod 数据
const fetchPodData = async () => {
    if (!selectedNamespace.value) {
        allPods.value = []
        return
    }

    loading.pods = true
    try {
        let url = ''
        if (selectedNamespace.value === '__ALL__') {
            // 获取所有命名空间的 Pod
            url = '/api/v1/pods'
        } else {
            // 获取指定命名空间的 Pod
            url = `/api/v1/namespaces/${selectedNamespace.value}/pods`
        }

        const response = await kubernetesRequest({
            url: url,
            method: "get"
        })

        if (response.code === 200 && response.data?.items) {
            allPods.value = response.data.items.map((pod: any) => ({
                uid: pod.metadata.uid,
                name: pod.metadata.name,
                namespace: pod.metadata.namespace,
                labels: pod.metadata.labels,
                status: pod.status?.phase || 'Unknown',
                reason: pod.status?.reason,
                message: pod.status?.message,
                ip: pod.status?.podIP,
                node: pod.spec?.nodeName,
                createdAt: pod.metadata.creationTimestamp
            }))

            // 调整页码
            nextTick(() => {
                const totalPages = Math.ceil(filteredPods.value.length / pageSize.value)
                if (currentPage.value > totalPages && totalPages > 0) {
                    currentPage.value = totalPages
                }
            })
        } else {
            ElMessage.error(`获取 Pod 数据失败: ${response.message || '未知错误'}`)
            allPods.value = []
        }
    } catch (error: any) {
        console.error("获取 Pod 数据失败:", error)
        ElMessage.error(`获取 Pod 数据出错: ${error.message || '网络请求失败'}`)
        allPods.value = []
    } finally {
        loading.pods = false
    }
}

// 获取 Pod 详情
const fetchPodDetails = async (namespace: string, name: string): Promise<PodDetail | null> => {
    try {
        const response = await kubernetesRequest({
            url: `/api/v1/namespaces/${namespace}/pods/${name}`,
            method: "get"
        })
        if (response.code === 200 && response.data) {
            return response.data
        } else {
            ElMessage.error(`获取 Pod 详情失败: ${response.message || '未知错误'}`)
            return null
        }
    } catch (error: any) {
        console.error(`获取 Pod 详情失败:`, error)
        ElMessage.error(`获取 Pod 详情出错: ${error.message || '网络错误'}`)
        return null
    }
}

// 辅助函数
const getStatusClass = (status: string) => {
    const lowerStatus = status.toLowerCase()
    if (lowerStatus === 'running' || lowerStatus === 'succeeded') return 'status-healthy'
    if (lowerStatus.includes('failed') || lowerStatus.includes('error') || lowerStatus.includes('crashloop')) return 'status-unhealthy'
    return 'status-unknown'
}

const getStatusText = (status: string) => {
    const statusMap: Record<string, string> = {
        running: t('podManagement.statuses.running'),
        pending: t('podManagement.statuses.pending'),
        succeeded: t('podManagement.statuses.succeeded'),
        failed: t('podManagement.statuses.failed'),
        unknown: t('podManagement.statuses.unknown'),
        terminating: t('podManagement.statuses.terminating')
    }
    return statusMap[status.toLowerCase()] || status
}

const formatTimestamp = (timestamp: string): string => {
    if (!timestamp) return 'N/A'
    const d = dayjs(timestamp)
    return d.isValid() ? d.format("YYYY-MM-DD HH:mm:ss") : 'Invalid Date'
}

// 事件处理
const handleNamespaceChange = () => {
    currentPage.value = 1
    searchQuery.value = ''
    filterStatus.value = ''
    fetchPodData()
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

const handlePodClick = (pod: PodDisplayItem) => {
    console.log('Pod clicked:', pod)
}

// Pod 操作
const handleAddPod = () => {
    if (!selectedNamespace.value || selectedNamespace.value === '__ALL__') {
        ElMessage.warning("请先选择一个具体的命名空间")
        return
    }
    yamlDialogConfig.content = getPlaceholderYaml()
    yamlDialogConfig.visible = true
}

const getPlaceholderYaml = () => {
    return `apiVersion: v1
kind: Pod
metadata:
  generateName: new-pod-
  namespace: ${selectedNamespace.value || 'default'}
  labels:
    app: myapp
    created-by: cilikube-ui
spec:
  containers:
  - name: main-container
    image: nginx:alpine
    ports:
    - containerPort: 80
      protocol: TCP
    resources:
      requests:
        memory: "64Mi"
        cpu: "100m"
      limits:
        memory: "128Mi"
        cpu: "200m"
  restartPolicy: Never`
}

const handleSaveYaml = async () => {
    if (!selectedNamespace.value) {
        ElMessage.error("未选择命名空间")
        return
    }

    yamlDialogConfig.saving = true
    try {
        const response = await kubernetesRequest({
            url: `/api/v1/namespaces/${selectedNamespace.value}/pods`,
            method: 'post',
            headers: { 'Content-Type': 'application/yaml' },
            data: yamlDialogConfig.content
        })

        if (response.code === 200 || response.code === 201) {
            ElMessage.success(t('podManagement.messages.createSuccess'))
            yamlDialogConfig.visible = false
            fetchPodData()
        } else {
            ElMessage.error(`操作失败: ${response.message || '未知错误'}`)
        }
    } catch (error: any) {
        console.error("应用 YAML 失败:", error)
        const errMsg = error.response?.data?.message || error.message || '请求失败'
        ElMessage.error(`应用 YAML 出错: ${errMsg}`)
    } finally {
        yamlDialogConfig.saving = false
    }
}

const handleYamlDialogClose = () => {
    yamlDialogConfig.content = ''
}

const handleDeletePod = (pod: PodDisplayItem) => {
    ElMessageBox.confirm(
        t('podManagement.messages.deleteConfirm', { name: pod.name }),
        t('common.confirm'),
        {
            confirmButtonText: t('common.confirm'),
            cancelButtonText: t('common.cancel'),
            type: 'warning'
        }
    ).then(async () => {
        try {
            await kubernetesRequest({
                url: `/api/v1/namespaces/${pod.namespace}/pods/${pod.name}`,
                method: "delete"
            })
            ElMessage.success(t('podManagement.messages.deleteSuccess'))
            await fetchPodData()
        } catch (error: any) {
            console.error("删除 Pod 失败:", error)
            const errMsg = error.response?.data?.message || error.message || '请求失败'
            ElMessage.error(`${t('podManagement.messages.deleteFailed')}: ${errMsg}`)
        }
    }).catch(() => {
        ElMessage.info('删除操作已取消')
    })
}

// 查看日志
const handleViewLogs = async (pod: PodDisplayItem) => {
    logDialogConfig.targetPod = pod
    logDialogConfig.visible = true
    logDialogConfig.containers = []

    const details = await fetchPodDetails(pod.namespace, pod.name)
    if (details?.spec) {
        const containers = details.spec.containers || []
        logDialogConfig.containers = containers
    }
}

const handleLogDialogClose = () => {
    logDialogConfig.targetPod = null
    logDialogConfig.containers = []
}

// 进入容器
const handleExec = async (pod: PodDisplayItem) => {
    execDialogConfig.targetPod = pod
    execDialogConfig.visible = true
    execDialogConfig.containers = []

    const details = await fetchPodDetails(pod.namespace, pod.name)
    if (details?.spec) {
        const containers = details.spec.containers || []
        execDialogConfig.containers = containers
    }
}

const handleExecDialogClose = () => {
    execDialogConfig.targetPod = null
    execDialogConfig.containers = []
}

// 监听命名空间变化
watch(() => selectedNamespace.value, (newNamespace) => {
    if (newNamespace) {
        fetchPodData()
    } else {
        allPods.value = []
    }
}, { immediate: true })

// 组件挂载
onMounted(async () => {
    await fetchNamespacesList()
    if (selectedNamespace.value) {
        await fetchPodData()
    }
})
</script>

<style scoped>
.pod-management {
    padding: 24px;
    background: #f8fafc;
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
    color: #1e293b;
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
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    border: 1px solid #e2e8f0;
}

.search-filters {
    display: flex;
    gap: 12px;
    align-items: center;
    flex: 1;
}

.namespace-select {
    width: 200px;
}

.search-input {
    width: 280px;
}

.filter-select {
    width: 160px;
}

.toolbar-right {
    display: flex;
    gap: 16px;
    align-items: center;
}

.view-toggle {
    border-radius: 8px;
}

/* 按钮样式 - 透明背景 + 彩色边框 */
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

:deep(.el-button--primary:active),
:deep(.el-button--primary.is-active) {
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

:deep(.el-button--success:active),
:deep(.el-button--success.is-active) {
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

:deep(.el-button--danger:active),
:deep(.el-button--danger.is-active) {
    background-color: #ef4444;
    border-color: #ef4444;
    color: #ffffff;
}

/* Pod 列表 */
.pod-list {
    min-height: 400px;
}

.pod-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(380px, 1fr));
    gap: 20px;
}

.pod-card {
    background: white;
    border-radius: 12px;
    border: 1px solid #e2e8f0;
    transition: all 0.2s ease;
    cursor: pointer;
    overflow: hidden;
}

.pod-card:hover {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    border-color: #3b82f6;
}

.card-header {
    padding: 20px 20px 16px;
    border-bottom: 1px solid #f1f5f9;
}

.pod-info {
    margin-bottom: 12px;
}

.pod-name {
    font-size: 18px;
    font-weight: 600;
    color: #1e293b;
    margin-bottom: 4px;
}

.pod-namespace {
    font-size: 13px;
    color: #475569;
    font-family: 'Maple Mono', 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
    font-weight: 600;
    margin-top: 4px;
}

.card-body {
    padding: 16px 20px;
}

.pod-meta {
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
    color: #1e293b;
    font-family: 'Maple Mono', 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
    font-weight: 600;
}

.card-footer {
    padding: 16px 20px;
    border-top: 1px solid #f1f5f9;
    background: white;
}

.pod-actions {
    display: flex;
    gap: 8px;
    justify-content: flex-end;
}

/* 状态指示器 - 简洁小圆点样式 */
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

/* Pod 名称行 */
.pod-name-row {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 8px;
}

/* 列表视图 */
.pod-table {
    background: white;
    border-radius: 12px;
    overflow: hidden;
    border: 1px solid #e2e8f0;
}

.table-pod-name {
    font-weight: 500;
    color: #1e293b;
}

.ip-address {
    font-family: 'Maple Mono', 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
    font-size: 12px;
    color: #475569;
    font-weight: 600;
}

.table-actions {
    display: flex;
    gap: 8px;
}

.empty-state {
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 80px 20px;
    background: white;
    border-radius: 12px;
    border: 1px solid #e2e8f0;
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
    border: 1px solid #e2e8f0;
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

.log-dialog :deep(.el-dialog__body) {
    padding: 20px 24px;
}

.log-options {
    margin-bottom: 16px;
    padding-bottom: 16px;
    border-bottom: 1px solid #e2e8f0;
}

.log-content-container {
    height: 60vh;
    border: 1px solid #e2e8f0;
    border-radius: 8px;
    overflow: hidden;
    background: #1e293b;
}

.log-content {
    height: 100%;
    overflow: auto;
    padding: 16px;
    margin: 0;
    font-family: 'Maple Mono', 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
    font-size: 12px;
    line-height: 1.5;
    color: #e2e8f0;
    white-space: pre-wrap;
    word-break: break-all;
    background: #1e293b;
}

.exec-dialog :deep(.el-dialog__body) {
    padding: 20px 24px;
}

.exec-options {
    margin-bottom: 16px;
    padding-bottom: 16px;
    border-bottom: 1px solid #e2e8f0;
}

.terminal-container {
    width: 100%;
    height: 60vh;
    background: #1e293b;
    border-radius: 8px;
    border: 1px solid #e2e8f0;
    display: flex;
    align-items: center;
    justify-content: center;
}

.terminal-placeholder {
    color: #64748b;
    font-family: 'Maple Mono', 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
    font-size: 14px;
}

.dialog-footer {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    padding: 16px 0 0;
    border-top: 1px solid #e2e8f0;
}

/* 响应式设计 */
@media (max-width: 1200px) {
    .pod-grid {
        grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
    }
}

@media (max-width: 768px) {
    .pod-management {
        padding: 16px;
    }

    .page-header {
        flex-direction: column;
        align-items: flex-start;
        gap: 16px;
    }

    .page-title {
        font-size: 28px;
    }

    .toolbar {
        flex-direction: column;
        align-items: stretch;
        gap: 16px;
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
        justify-content: space-between;
        width: 100%;
    }

    .pod-grid {
        grid-template-columns: 1fr;
    }

    .meta-item {
        flex-direction: column;
        align-items: flex-start;
        gap: 4px;
    }

    .pod-actions {
        flex-direction: column;
    }

    .table-actions {
        flex-direction: column;
        gap: 4px;
    }
}
</style>