<template>
  <div class="node-container">
    <!-- 解释文字卡片 -->
    <el-row :gutter="20" class="card-row explanation-card">
      <el-col :span="24">
        <el-card class="custom-card">
          <template #header>
            <div class="card-header">
              <span>集群知识</span>
            </div>
          </template>
          <div class="card-content">
            <p>Kubernetes 通过将容器放入在节点（Node）上运行的 Pod 中来执行你的工作负载。 节点可以是一个虚拟机或者物理机器，取决于所在的集群配置。 每个节点包含运行 Pod 所需的服务； 这些节点由控制面负责管理。通常集群中会有若干个节点；而在一个学习所用或者资源受限的环境中，你的集群中也可能只有一个节点。</p>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 卡片布局 -->
    <el-row :gutter="20" class="card-row">
      <!-- 集群信息卡片 -->
      <el-col :span="8" v-for="node in currentPageData" :key="node.metadata.uid">
        <el-card class="custom-card">
          <template #header>
            <div class="card-header">
              <span>集群信息</span>
            </div>
          </template>
          <div class="card-content">
            <p>集群名称：<el-tag type="primary" class="custom-tag">{{ node.metadata.name }}</el-tag></p>
            <p>状态: <el-tag :type="getNodeStatusTagType(node.status.conditions)" class="custom-tag">{{ getNodeStatus(node.status.conditions) }}</el-tag></p>
            <p>节点版本: <el-tag type="primary" class="custom-tag">{{ node.status.nodeInfo.kubeletVersion }}</el-tag></p>
            <p>CPU: <el-tag type="warning" class="custom-tag">{{ node.status.capacity.cpu }}核</el-tag></p>
            <p>内存: <el-tag type="danger" class="custom-tag">{{ formatMemory(node.status.capacity.memory) }}</el-tag></p>
          </div>
        </el-card>
      </el-col>

      <!-- 其他信息卡片 -->
      <el-col :span="8" v-for="node in currentPageData" :key="node.metadata.uid + '-extra'">
        <el-card class="custom-card">
          <template #header>
            <div class="card-header">
              <span>版本信息</span>
            </div>
          </template>
          <div class="card-content">
            <p>操作系统: <el-tag type="info" class="custom-tag">{{ node.status.nodeInfo.osImage }}</el-tag></p>
            <p>内核版本: <el-tag type="info" class="custom-tag">{{ node.status.nodeInfo.kernelVersion }}</el-tag></p>
            <p>容器运行时版本: <el-tag type="info" class="custom-tag">{{ node.status.nodeInfo.containerRuntimeVersion }}</el-tag></p>
            <p>kube-proxy 版本: <el-tag type="info" class="custom-tag">{{ node.status.nodeInfo.kubeProxyVersion }}</el-tag></p>
            
          </div>
        </el-card>
      </el-col>

      <!-- 资源信息卡片 -->
      <el-col :span="8" v-for="node in currentPageData" :key="node.metadata.uid + '-resources'">
        <el-card class="custom-card">
          <template #header>
            <div class="card-header">
              <span>节点信息</span>
            </div>
          </template>
          <div class="card-content">
            <p>UID: <el-tag type="info" class="custom-tag">{{ node.metadata.uid }}</el-tag></p>
            <p>创建时间: <el-tag type="success" class="custom-tag">{{ formatTimestamp(node.metadata.creationTimestamp) }}</el-tag></p>
            <p>Pod CIDR: <el-tag type="info" class="custom-tag">{{ node.spec.podCIDR }}</el-tag></p>
            <p>IP地址: <el-tag type="primary" class="custom-tag">{{ getNodeAddress(node.status.addresses) }}</el-tag></p>
            <!-- <p>分配的 CPU: <el-tag type="warning" class="custom-tag">{{ node.status.allocatable.cpu }}</el-tag></p>
            <p>分配的内存: <el-tag type="danger" class="custom-tag">{{ formatMemory(node.status.allocatable.memory) }}</el-tag></p> -->
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 按钮组 -->
    <div class="button-group">
      <el-button type="primary" @click="handleCreateNode">新建节点</el-button>
      <el-button type="success" @click="handleMonitor">监控</el-button>
      <el-button type="info" @click="handleAction1">升级</el-button>
      <el-button type="warning" @click="handleAction2">更多操作</el-button>
    </div>

    <!-- 表格展示详细信息 -->
    <el-table :data="currentPageData" style="width: 100%; margin-top: 20px;">
      <el-table-column prop="metadata.uid" label="节点UID" min-width="200" />
      <el-table-column prop="metadata.creationTimestamp" label="创建时间" min-width="180" :formatter="formatTimestamp" />
      <el-table-column prop="status.nodeInfo.kernelVersion" label="内核版本" min-width="180" />
      <el-table-column prop="status.nodeInfo.containerRuntimeVersion" label="容器运行时版本" min-width="200" />
      <el-table-column prop="status.addresses" label="IP地址" min-width="200">
        <template #default="{ row }">
          <span>{{ getNodeAddress(row.status.addresses) }}</span>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页组件 -->
    <el-pagination
      v-model:currentPage="currentPage"
      :page-size="pageSize"
      :total="totalNodes"
      layout="total, prev, pager, next, jumper"
      @current-change="handlePageChange"
    />
  </div>
</template>

<script lang="ts">
import { defineComponent, onMounted, ref, computed } from "vue"
import { ElMessage } from "element-plus"
import { request } from "@/utils/service"
import dayjs from "dayjs"

// 定义 Node 接口
interface Node {
  metadata: {
    name: string
    uid: string
    creationTimestamp: string
  }
  status: {
    capacity: {
      cpu: string
      memory: string
    }
    allocatable: {
      cpu: string
      memory: string
    }
    conditions: Array<{
      type: string
      status: string
    }>
    addresses: Array<{
      type: string
      address: string
    }>
    nodeInfo: {
      osImage: string
      kernelVersion: string
      containerRuntimeVersion: string
      kubeletVersion: string
      kubeProxyVersion: string
    }
  }
  spec: {
    podCIDR: string
    podCIDRs: string[]
  }
}

// 定义 NodeResponse 接口
interface NodeResponse {
  code: number
  data: {
    items: Node[]
  }
  message: string
}

export default defineComponent({
  name: "Node",
  setup() {
    // 定义响应式变量
    const nodeData = ref<Node[]>([])
    const currentPage = ref(1)
    const pageSize = ref(10)
    const totalNodes = ref(0)

    // 计算当前页的数据
    const currentPageData = computed(() => {
      const start = (currentPage.value - 1) * pageSize.value
      const end = start + pageSize.value
      return nodeData.value.slice(start, end)
    })

    // 获取节点数据
    const fetchNodeData = async () => {
      try {
        const response = await request<NodeResponse>({
          url: "/api/v1/nodes",
          method: "get",
          baseURL: "VITE_API_BASE_URL" // 可根据需要调整 baseURL
        })
        // console.log("API response:", response)
        if (response.code === 200) {
          nodeData.value = response.data.items
          totalNodes.value = response.data.items.length
        } else {
          ElMessage.error("获取节点数据失败: " + response.message)
        }
      } catch (error) {
        console.error("获取节点数据失败:", error)
        ElMessage.error("获取节点数据失败")
      }
    }

    // 处理分页变化
    const handlePageChange = (page: number) => {
      currentPage.value = page
    }

    // 获取节点状态
    const getNodeStatus = (conditions: Array<{ type: string; status: string }>) => {
      const readyCondition = conditions.find(condition => condition.type === "Ready")
      return readyCondition && readyCondition.status === "True" ? "健康" : "不健康"
    }

    // 获取节点状态标签类型
    const getNodeStatusTagType = (conditions: Array<{ type: string; status: string }>) => {
      const readyCondition = conditions.find(condition => condition.type === "Ready")
      return readyCondition && readyCondition.status === "True" ? "success" : "danger"
    }

    // 获取节点地址
    const getNodeAddress = (addresses: Array<{ type: string; address: string }>) => {
      const internalIP = addresses.find(address => address.type === "InternalIP")
      return internalIP ? internalIP.address : "Unknown"
    }

    // 格式化时间戳
    const formatTimestamp = (row: any, column: any, cellValue: string) => {
      return dayjs(cellValue).format("YYYY-MM-DD HH:mm:ss")
    }

    // 格式化内存
    const formatMemory = (memory: string) => {
      const memoryInKi = parseInt(memory.replace("Ki", ""), 10)
      const memoryInGi = (memoryInKi / 1024 / 1024).toFixed(1)
      return `${memoryInGi} GB`
    }

    // 处理新建节点
    const handleCreateNode = () => {
      ElMessage.info("新建节点功能待实现")
    }

    // 处理监控
    const handleMonitor = () => {
      ElMessage.info("监控功能待实现")
    }

    // 处理更多操作
    const handleAction1 = () => {
      ElMessage.info("更多功能待实现")
    }

    const handleAction2 = () => {
      ElMessage.info("更多功能待实现")
    }

    // 组件挂载时获取节点数据
    onMounted(() => {
      fetchNodeData()
    })

    return {
      nodeData,
      currentPage,
      pageSize,
      totalNodes,
      currentPageData,
      getNodeStatus,
      getNodeStatusTagType,
      getNodeAddress,
      formatTimestamp,
      formatMemory,
      handlePageChange,
      handleCreateNode,
      handleMonitor,
      handleAction1,
      handleAction2
    }
  }
})
</script>

<style scoped>
.node-container {
  padding: 20px;
}

.card-row {
  display: flex;
  flex-wrap: wrap;
}

.explanation-card {
  margin-bottom: 20px;
}

.el-col {
  display: flex;
}

.custom-card {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.el-card {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.card-header {
  font-size: 16px;
  font-weight: bold;
}

.card-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  font-size: 14px;
  /* gap: 10px; */
}

.custom-tag {
  font-size: 14px;
  text-align: left;
}

.button-group {
  margin: 20px 0;
  display: flex;
  gap: 10px;
}
</style>