<template>
  <div class="deployment-container">
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
      :data="filteredAndSortedData"
      border
      stripe
      style="width: 100%; margin-top: 20px;"
      @sort-change="handleSortChange"
    >
      <el-table-column prop="name" label="Deployment名称" min-width="180" sortable />
      <el-table-column prop="namespace" label="命名空间" min-width="180" sortable />
      <el-table-column prop="replicas" label="副本数" min-width="100" sortable />
      <el-table-column prop="deploymentSpec.template.spec.containers[0].image" label="镜像" min-width="300" sortable />
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
          <el-input v-model="currentDeployment.name" />
        </el-form-item>
        <el-form-item label="命名空间">
          <el-input v-model="currentDeployment.namespace" />
        </el-form-item>
        <el-form-item label="副本数">
          <el-input-number v-model="currentDeployment.replicas" :min="1" />
        </el-form-item>
        <el-form-item label="镜像">
          <el-input v-model="currentDeployment.deploymentSpec.template.spec.containers[0].image" />
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
import { ElMessage, ElMessageBox } from "element-plus"
import { request } from "@/utils/service"

interface Deployment {
  name: string
  namespace: string
  replicas: number
  deploymentSpec: {
    replicas: number
    selector: {
      matchLabels: { [key: string]: string }
    }
    template: {
      metadata: {
        creationTimestamp: string | null
        labels: { [key: string]: string }
      }
      spec: {
        containers: {
          name: string
          image: string
          args?: string[]
          ports?: {
            containerPort: number
            protocol: string
          }[]
          resources?: {
            limits?: { [key: string]: string }
            requests?: { [key: string]: string }
          }
          livenessProbe?: {
            httpGet: {
              path: string
              port: number
              scheme: string
            }
            initialDelaySeconds: number
            timeoutSeconds: number
            periodSeconds: number
            successThreshold: number
            failureThreshold: number
          }
          terminationMessagePath: string
          terminationMessagePolicy: string
          imagePullPolicy: string
          securityContext?: {
            capabilities?: {
              add?: string[]
              drop?: string[]
            }
            readOnlyRootFilesystem?: boolean
            allowPrivilegeEscalation?: boolean
          }
        }[]
        restartPolicy: string
        terminationGracePeriodSeconds: number
        dnsPolicy: string
        nodeSelector?: { [key: string]: string }
        serviceAccountName: string
        serviceAccount: string
        securityContext: { [key: string]: any }
        schedulerName: string
        tolerations?: {
          key: string
          operator: string
          effect: string
        }[]
      }
    }
    strategy: {
      type: string
      rollingUpdate: {
        maxUnavailable: string
        maxSurge: string
      }
    }
    revisionHistoryLimit: number
    progressDeadlineSeconds: number
  }
  deploymentStatus: {
    observedGeneration: number
    replicas: number
    updatedReplicas: number
    readyReplicas: number
    availableReplicas: number
    conditions: {
      type: string
      status: string
      lastUpdateTime: string
      lastTransitionTime: string
      reason: string
      message: string
    }[]
  }
}

export default defineComponent({
  name: "Deployment",
  setup() {
    const deploymentData = ref<Deployment[]>([])
    const currentPage = ref(1)
    const pageSize = ref(10)
    const totalDeployments = ref(0)
    const dialogVisible = ref(false)
    const dialogTitle = ref("新增Deployment")
    const currentDeployment = ref<Deployment>({
      name: "",
      namespace: "default",
      replicas: 1,
      deploymentSpec: {
        replicas: 1,
        selector: {
          matchLabels: {
            app: "default"
          }
        },
        template: {
          metadata: {
            creationTimestamp: null,
            labels: {
              app: "default"
            }
          },
          spec: {
            containers: [
              {
                name: "container-name",
                image: "nginx:latest",
                terminationMessagePath: "/dev/termination-log",
                terminationMessagePolicy: "File",
                imagePullPolicy: "IfNotPresent"
              }
            ],
            restartPolicy: "Always",
            terminationGracePeriodSeconds: 30,
            dnsPolicy: "ClusterFirst",
            serviceAccountName: "default",
            serviceAccount: "default",
            securityContext: {},
            schedulerName: "default-scheduler"
          }
        },
        strategy: {
          type: "RollingUpdate",
          rollingUpdate: {
            maxUnavailable: "25%",
            maxSurge: "25%"
          }
        },
        revisionHistoryLimit: 10,
        progressDeadlineSeconds: 600
      },
      deploymentStatus: {
        observedGeneration: 1,
        replicas: 1,
        updatedReplicas: 1,
        readyReplicas: 1,
        availableReplicas: 1,
        conditions: []
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
            deployment.name.toLowerCase().includes(query) ||
            deployment.namespace.toLowerCase().includes(query) ||
            deployment.deploymentSpec.template.spec.containers[0].image.toLowerCase().includes(query)
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
      try {
        const response = await request({
          url: "/api/v1/namespaces/kube-system/deployments",
          method: "get",
          baseURL: "http://192.168.1.100:8080" // 可根据需要调整 baseURL
        })
        console.log("API response:", response)
        deploymentData.value = response
        totalDeployments.value = response.length
      } catch (error) {
        console.error("获取Deployment数据失败:", error)
        ElMessage.error("获取Deployment数据失败")
      }
    }

    const handlePageChange = (page: number) => {
      currentPage.value = page
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
        name: "",
        namespace: "default",
        replicas: 1,
        deploymentSpec: {
          replicas: 1,
          selector: {
            matchLabels: {
              app: "default"
            }
          },
          template: {
            metadata: {
              creationTimestamp: null,
              labels: {
                app: "default"
              }
            },
            spec: {
              containers: [
                {
                  name: "container-name",
                  image: "nginx:latest",
                  terminationMessagePath: "/dev/termination-log",
                  terminationMessagePolicy: "File",
                  imagePullPolicy: "IfNotPresent"
                }
              ],
              restartPolicy: "Always",
              terminationGracePeriodSeconds: 30,
              dnsPolicy: "ClusterFirst",
              serviceAccountName: "default",
              serviceAccount: "default",
              securityContext: {},
              schedulerName: "default-scheduler"
            }
          },
          strategy: {
            type: "RollingUpdate",
            rollingUpdate: {
              maxUnavailable: "25%",
              maxSurge: "25%"
            }
          },
          revisionHistoryLimit: 10,
          progressDeadlineSeconds: 600
        },
        deploymentStatus: {
          observedGeneration: 1,
          replicas: 1,
          updatedReplicas: 1,
          readyReplicas: 1,
          availableReplicas: 1,
          conditions: []
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
            url: "/apis/v1/k8s/deployments",
            method: "post",
            data: currentDeployment.value,
            baseURL: "http://192.168.1.200:8080" // 可根据需要调整 baseURL
          })
          ElMessage.success("Deployment新增成功")
        } else {
          await request({
            url: `/apis/v1/k8s/deployments/${currentDeployment.value.name}`,
            method: "put",
            data: currentDeployment.value,
            baseURL: "http://192.168.1.200:8080" // 可根据需要调整 baseURL
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
      ElMessageBox.confirm(`确定删除Deployment ${deployment.name} 吗？`, "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning"
      }).then(async () => {
        try {
          await request({
            url: `/apis/v1/k8s/deployments/${deployment.name}`,
            method: "delete",
            baseURL: "http://192.168.1.200:8080" // 可根据需要调整 baseURL
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
      fetchDeploymentData()
    })

    return {
      deploymentData,
      currentPage,
      pageSize,
      totalDeployments,
      currentPageData,
      filteredAndSortedData,
      searchQuery,
      sortKey,
      sortOrder,
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
</style>