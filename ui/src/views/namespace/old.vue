<template>
  <div class="namespace-container">
    <!-- 卡片 -->
    <el-row :gutter="20" class="card-row explanation-card">
      <el-col :span="24">
        <el-card class="custom-card">
        <template #header>
          <div class="card-header">
            <span>集群知识</span>
          </div>
        </template>
        <div class="card-content">
            <p>在 Kubernetes 中，名字空间（Namespace） 提供一种机制，将同一集群中的资源划分为相互隔离的组。 同一名字空间内的资源名称要唯一，但跨名字空间时没有这个要求。 名字空间作用域仅针对带有名字空间的对象， （例如 Deployment、Service 等），这种作用域对集群范围的对象 （例如 StorageClass、Node、PersistentVolume 等）不适用。</p>
        </div>
    </el-card>
      </el-col>
    </el-row>
    <!-- 搜索框和新增按钮 -->
    <el-row :gutter="20" style="margin-bottom: 20px;">
      <el-col :span="18">
        <el-input v-model="searchQuery" placeholder="搜索命名空间" @input="handleSearch" />
      </el-col>
      <el-col :span="6">
        <el-button type="primary" @click="showCreateDialog" style="width: 100%;">新增命名空间</el-button>
      </el-col>
    </el-row>

    <!-- 命名空间列表 -->
    <el-table :data="paginatedData" style="width: 100%;">
      <el-table-column prop="name" label="命名空间名称" />
      <el-table-column prop="status" label="状态">
        <template #default="{ row }">
          <el-tag :type="row.status === 'Active' ? 'success' : 'danger'">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="creationTimestamp" label="创建时间" />
    </el-table>

    <!-- 分页组件 -->
    <el-pagination
      v-model:currentPage="currentPage"
      :page-size="pageSize"
      :total="filteredData.length"
      layout="total, prev, pager, next, jumper"
      @current-change="handlePageChange"
    />

    <!-- 新增命名空间对话框 -->
    <el-dialog v-model="isDialogVisible" title="新增命名空间">
      <el-form :model="form">
        <el-form-item label="命名空间名称">
          <el-input v-model="form.name" />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="isDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="createNamespace">确定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted, ref, computed } from "vue"
import { ElMessage, ElLoading } from "element-plus"
import { request } from "@/utils/service"
import dayjs from "dayjs"

// 定义 Namespace 接口
interface Namespace {
  name: string
  status: string
  creationTimestamp: string
}

export default defineComponent({
  name: "Namespace",
  setup() {
    // 定义响应式变量
    const namespaceData = ref<Namespace[]>([])
    const currentPage = ref(1)
    const pageSize = ref(10)
    const searchQuery = ref("")
    const isDialogVisible = ref(false)
    const form = ref<Namespace>({ name: "", status: "Active", creationTimestamp: "" })

    // 获取命名空间数据
    const fetchNamespaceData = async () => {
      const loading = ElLoading.service({
        lock: true,
        text: "加载中",
        background: "rgba(0, 0, 0, 0.7)",
      })
      try {
        const response = await request<{ code: number; data: { items: any[] }; message: string }>({
          url: "/api/v1/namespace",
          method: "get",
          baseURL: "VITE_API_BASE_URL",
        })
        if (response.code === 200) {
          namespaceData.value = response.data.items.map(item => ({
            name: item.metadata.name,
            status: item.status.phase,
            creationTimestamp: dayjs(item.metadata.creationTimestamp).format("YYYY-MM-DD HH:mm:ss")
          }))
        } else {
          ElMessage.error("获取命名空间数据失败: " + response.message)
        }
      } catch (error) {
        console.error("获取命名空间数据失败:", error)
        ElMessage.error("获取命名空间数据失败")
      } finally {
        loading.close()
      }
    }

    // 创建命名空间
    const createNamespace = async () => {
      const loading = ElLoading.service({
        lock: true,
        text: "创建中",
        background: "rgba(0, 0, 0, 0.7)",
      })
      try {
        const response = await request<{ code: number; message: string }>({
          url: "/api/v1/namespace",
          method: "post",
          baseURL: "VITE_API_BASE_URL",
          data: form.value
        })
        if (response.code === 200) {
          ElMessage.success("命名空间创建成功")
          fetchNamespaceData()
          isDialogVisible.value = false
        } else {
          ElMessage.error("命名空间创建失败: " + response.message)
        }
      } catch (error) {
        console.error("命名空间创建失败:", error)
        ElMessage.error("命名空间创建失败")
      } finally {
        loading.close()
      }
    }

    // 显示新增命名空间对话框
    const showCreateDialog = () => {
      form.value = { name: "", status: "Active", creationTimestamp: "" }
      isDialogVisible.value = true
    }

    // 处理分页变化
    const handlePageChange = (page: number) => {
      currentPage.value = page
    }

    // 处理搜索
    const handleSearch = () => {
      currentPage.value = 1
    }

    // 过滤数据
    const filteredData = computed(() => {
      return namespaceData.value.filter((ns) => ns.name.includes(searchQuery.value))
    })

    // 分页数据
    const paginatedData = computed(() => {
      const start = (currentPage.value - 1) * pageSize.value
      const end = start + pageSize.value
      return filteredData.value.slice(start, end)
    })

    // 组件挂载时获取命名空间数据
    onMounted(() => {
      fetchNamespaceData()
    })

    return {
      namespaceData,
      currentPage,
      pageSize,
      searchQuery,
      isDialogVisible,
      form,
      fetchNamespaceData,
      handlePageChange,
      handleSearch,
      filteredData,
      paginatedData,
      showCreateDialog,
      createNamespace
    }
  }
})
</script>

<style scoped>
.namespace-container {
  padding: 20px;
}

.info-card {
  margin-bottom: 20px;
}

.card-header {
  font-size: 18px;
  font-weight: bold;
}

.text {
  margin: 10px 0;
}

.card-row {
  display: flex;
  flex-wrap: wrap;
}

.explanation-card {
  margin-bottom: 20px;
}
</style>