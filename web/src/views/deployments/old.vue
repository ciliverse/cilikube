<template>
  <div class="deployment-container">

        <!-- 集群知识卡片 -->
        <el-row :gutter="20" class="card-row explanation-card">
      <el-col :span="24">
        <el-card class="custom-card">
          <template #header>
            <div class="card-header">
              <span>集群知识</span>
            </div>
          </template>
          <div class="card-content">
            <p>Deployment 用于管理运行一个应用负载的一组 Pod，通常适用于不保持状态的负载。</p>
          </div>
        </el-card>
      </el-col>
    </el-row>
    <!-- 命名空间筛选框 -->
    <el-select v-model="selectedNamespace" placeholder="选择命名空间" @change="handleNamespaceChange" class="namespace-select">
      <el-option v-for="namespace in namespaces" :key="namespace" :label="namespace" :value="namespace" />
    </el-select>

    <!-- 搜索框 -->
    <el-input
      v-model="searchQuery"
      placeholder="搜索Deployment"
      clearable
      @clear="handleSearch"
      @input="handleSearch"
      class="search-input"
    />

    <!-- 新增Deployment按钮 -->
    <el-button type="success" @click="handleAdd" class="add-deployment-button">新增Deployment</el-button>

    <el-table
      :data="currentPageData"
      border
      stripe
      style="width: 100%; margin-top: 20px;"
      @sort-change="handleSortChange"
    >
      <el-table-column prop="metadata.name" label="Deployment名称" min-width="180" sortable />
      <el-table-column prop="metadata.namespace" label="命名空间" min-width="180" sortable />
      <el-table-column prop="spec.replicas" label="副本数" min-width="120" sortable />
      <el-table-column prop="status.availableReplicas" label="可用副本数" min-width="120" sortable />
      <el-table-column prop="metadata.creationTimestamp" label="创建时间" min-width="200" sortable />
      <el-table-column label="操作" min-width="200">
        <template #default="{ row }">
          <el-button type="primary" size="small" @click="handleEdit(row)">编辑</el-button>
          <el-button type="danger" size="small" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页组件 -->
    <el-pagination
      v-model:currentPage="currentPage"
      :page-size="pageSize"
      :total="totalDeployments"
      layout="total, prev, pager, next, jumper"
      @current-change="handlePageChange"
    />

    <!-- 新增/编辑Deployment对话框 -->
    <el-dialog :title="dialogTitle" v-model="dialogVisible">
      <el-form :model="currentDeployment" label-width="120px">
        <el-form-item label="Deployment名称">
          <el-input v-model="currentDeployment.metadata.name" />
        </el-form-item>
        <el-form-item label="命名空间">
          <el-input v-model="currentDeployment.metadata.namespace" />
        </el-form-item>
        <el-form-item label="副本数">
          <el-input v-model="currentDeployment.spec.replicas" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted, ref, computed } from "vue"
import { ElMessage, ElMessageBox, ElLoading } from "element-plus"
import { request } from "@/utils/service"
import dayjs from "dayjs"

interface Deployment {
  metadata: {
    name: string
    namespace: string
    creationTimestamp: string
  }
  spec: {
    replicas: number
  }
  status: {
    availableReplicas: number
  }
}

interface DeploymentResponse {
  code: number
  data: {
    items: Deployment[] | null
    total: number
  }
  message: string
}

export default defineComponent({
  name: "Deployment",
  setup() {
    const deploymentData = ref<Deployment[]>([])
    const namespaces = ref<string[]>([])
    const selectedNamespace = ref<string>("")
    const currentPage = ref(1)
    const pageSize = ref(10)
    const totalDeployments = ref(0)
    const dialogVisible = ref(false)
    const dialogTitle = ref("新增Deployment")
    const currentDeployment = ref<Deployment>({
      metadata: {
        name: "",
        namespace: "default",
        creationTimestamp: ""
      },
      spec: {
        replicas: 1
      },
      status: {
        availableReplicas: 0
      }
    })
    const searchQuery = ref("")
    const sortKey = ref("")
    const sortOrder = ref("")

    const currentPageData = computed(() => {
      const start = (currentPage.value - 1) * pageSize.value
      const end = start + pageSize.value
      return filteredAndSortedData.value.slice(start, end)
    })

    const filteredAndSortedData = computed(() => {
      let data = deploymentData.value

      // 搜索过滤
      if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase()
        data = data.filter(
          (deployment) =>
            deployment.metadata.name.toLowerCase().includes(query) ||
            deployment.metadata.namespace.toLowerCase().includes(query)
        )
      }

      // 排序
      if (sortKey.value) {
        data = data.slice().sort((a, b) => {
          const aValue = a[sortKey.value]
          const bValue = b[sortKey.value]
          if (sortOrder.value === "ascending") {
            return aValue > bValue ? 1 : -1
          } else if (sortOrder.value === "descending") {
            return aValue < bValue ? 1 : -1
          }
          return 0
        })
      }

      return data
    })

    const fetchDeploymentData = async () => {
      if (!selectedNamespace.value) {
        ElMessage.warning("请选择命名空间")
        return
      }

      const loading = ElLoading.service({
        lock: true,
        text: "加载中",
        background: "rgba(0, 0, 0, 0.7)",
      })
      try {
        const response = await request({
          url: `/api/v1/namespaces/${selectedNamespace.value}/deployments`,
          method: "get",
          baseURL: "VITE_API_BASE_URL"
        }) as DeploymentResponse
        console.log("Deployment response:", response) // 添加日志
        if (response.code === 200) {
          deploymentData.value = response.data.items ? response.data.items.map(item => ({
            ...item,
            metadata: {
              ...item.metadata,
              creationTimestamp: dayjs(item.metadata.creationTimestamp).format("YYYY-MM-DD HH:mm:ss")
            }
          })) : []
          totalDeployments.value = response.data.items ? response.data.items.length : 0
        } else {
          ElMessage.error("获取Deployment数据失败: " + response.message)
        }
      } catch (error) {
        console.error("获取Deployment数据失败:", error)
        ElMessage.error("获取Deployment数据失败: " + (error instanceof Error ? error.message : String(error)))
      } finally {
        loading.close()
      }
    }

    const fetchNamespaces = async () => {
      const loading = ElLoading.service({
        lock: true,
        text: "加载中",
        background: "rgba(0, 0, 0, 0.7)",
      })
      try {
        const response = await request<{ code: number; data: string[]; message: string }>({
          url: "/api/v1/namespaces",
          method: "get",
          baseURL: "VITE_API_BASE_URL"
        })
        console.log("Namespace response:", response) // 添加日志
        if (response.code === 200 && response.data) {
          namespaces.value = response.data
          if (namespaces.value.length > 0) {
            selectedNamespace.value = namespaces.value[0]
            await fetchDeploymentData() // 等待 Deployment 数据加载完成
          } else {
            ElMessage.warning("没有可用的命名空间")
          }
        } else {
          ElMessage.error(`获取命名空间失败: ${response.message || '未知错误'}`)
        }
      } catch (error) {
        console.error("获取命名空间失败:", error)
        ElMessage.error("获取命名空间失败: " + (error instanceof Error ? error.message : String(error)))
      } finally {
        loading.close()
      }
    }

    const handleNamespaceChange = async () => {
      currentPage.value = 1
      await fetchDeploymentData()
    }

    const handlePageChange = (page: number) => {
      currentPage.value = page
      fetchDeploymentData()
    }

    const handleSearch = () => {
      currentPage.value = 1
    }

    const handleSortChange = ({ prop, order }: { prop: string; order: string }) => {
      sortKey.value = prop
      sortOrder.value = order
    }

    const handleAdd = () => {
      dialogTitle.value = "新增Deployment"
      currentDeployment.value = {
        metadata: {
          name: "",
          namespace: selectedNamespace.value,
          creationTimestamp: ""
        },
        spec: {
          replicas: 1
        },
        status: {
          availableReplicas: 0
        }
      }
      dialogVisible.value = true
    }

    const handleEdit = (deployment: Deployment) => {
      dialogTitle.value = "编辑Deployment"
      currentDeployment.value = JSON.parse(JSON.stringify(deployment))
      dialogVisible.value = true
    }

    const handleSave = async () => {
      try {
        if (dialogTitle.value === "新增Deployment") {
          await request({
            url: `/api/v1/namespaces/${currentDeployment.value.metadata.namespace}/deployments`,
            method: "post",
            data: currentDeployment.value,
            baseURL: "VITE_API_BASE_URL"
          })
          ElMessage.success("Deployment新增成功")
        } else {
          await request({
            url: `/api/v1/namespaces/${currentDeployment.value.metadata.namespace}/deployments/${currentDeployment.value.metadata.name}`,
            method: "put",
            data: currentDeployment.value,
            baseURL: "VITE_API_BASE_URL"
          })
          ElMessage.success("Deployment编辑成功")
        }
        dialogVisible.value = false
        fetchDeploymentData()
      } catch (error) {
        console.error("保存Deployment失败:", error)
        ElMessage.error("保存Deployment失败")
      }
    }

    const handleDelete = (deployment: Deployment) => {
      ElMessageBox.confirm(`确定删除Deployment ${deployment.metadata.name} 吗？`, "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning"
      }).then(async () => {
        try {
          await request({
            url: `/api/v1/namespaces/${deployment.metadata.namespace}/deployments/${deployment.metadata.name}`,
            method: "delete",
            baseURL: "VITE_API_BASE_URL"
          })
          ElMessage.success("Deployment删除成功")
          fetchDeploymentData()
        } catch (error) {
          console.error("删除Deployment失败:", error)
          ElMessage.error("删除Deployment失败")
        }
      })
    }

    onMounted(() => {
      fetchNamespaces()
    })

    return {
      deploymentData,
      namespaces,
      selectedNamespace,
      currentPage,
      pageSize,
      totalDeployments,
      currentPageData,
      filteredAndSortedData,
      searchQuery,
      sortKey,
      sortOrder,
      handleNamespaceChange,
      handlePageChange,
      handleSearch,
      handleSortChange,
      dialogVisible,
      dialogTitle,
      currentDeployment,
      handleAdd,
      handleEdit,
      handleSave,
      handleDelete
    }
  }
})
</script>

<style scoped>
.deployment-container {
  padding: 20px;
}

.el-tag {
  font-size: 12px;
}

.el-table-column .cell {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.add-deployment-button {
  margin-bottom: 20px;
}

.search-input {
  margin-bottom: 20px;
  width: 300px;
}

.namespace-select {
  margin-bottom: 20px;
  width: 300px;
}

.card-row {
  margin-bottom: 20px;
}

.custom-card {
  border: 1px solid #ebeef5;
  border-radius: 4px;
}

.card-header {
  font-size: 16px;
  font-weight: bold;
}

.card-content {
  font-size: 14px;
  color: #606266;
}
</style>