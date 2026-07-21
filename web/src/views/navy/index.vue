<template>
  <div class="navigation-container">
    <!-- 页面标题 -->
    <div class="header">
      <h1>运维导航</h1>
      <div class="header-actions">
        <el-input v-model="searchQuery" placeholder="搜索运维网站..." clearable prefix-icon="Search" class="search-input"
          @input="filterSites" />
        <el-button type="primary" :icon="Plus" @click="showAddDialog">
          添加网站
        </el-button>
      </div>
    </div>

    <!-- 分类标签 -->
    <div class="category-tabs">
      <el-tag v-for="category in categories" :key="category" :type="selectedCategory === category ? 'primary' : 'info'"
        :effect="selectedCategory === category ? 'dark' : 'plain'" class="category-tag"
        @click="filterByCategory(category)">
        {{ category }}
      </el-tag>
    </div>

    <!-- 导航卡片列表 -->
    <div class="card-list">
      <el-row :gutter="20">
        <el-col :xs="24" :sm="12" :md="8" :lg="6" v-for="(site, index) in filteredSites" :key="site.id">
          <div class="nav-card">
            <div class="card-actions">
              <el-dropdown trigger="click" @command="handleCommand">
                <el-icon class="action-icon">
                  <MoreFilled />
                </el-icon>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item :command="{ action: 'edit', site }">
                      <el-icon>
                        <Edit />
                      </el-icon>编辑
                    </el-dropdown-item>
                    <el-dropdown-item :command="{ action: 'delete', site }" divided>
                      <el-icon>
                        <Delete />
                      </el-icon>删除
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
            <div class="card-main" @click="navigateTo(site.url)">
              <div class="card-icon" :style="{ backgroundColor: site.color }">
                <img v-if="site.logo" :src="site.logo" :alt="site.name" class="site-logo" />
                <el-icon v-else :size="30" color="#fff">
                  <component :is="site.icon" />
                </el-icon>
              </div>
              <div class="card-content">
                <h3>{{ site.name }}</h3>
                <p>{{ site.description }}</p>
                <div class="card-meta">
                  <el-tag size="small" type="info">{{ site.category }}</el-tag>
                  <span class="visit-count">访问: {{ site.visitCount || 0 }}</span>
                </div>
              </div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 空状态提示 -->
    <el-empty v-if="filteredSites.length === 0" description="未找到匹配的网站" />

    <!-- 添加/编辑网站对话框 -->
    <el-dialog v-model="dialogVisible" :title="isEditing ? '编辑网站' : '添加网站'" width="500px" @close="resetForm">
      <el-form :model="siteForm" :rules="formRules" ref="formRef" label-width="80px">
        <el-form-item label="网站名称" prop="name">
          <el-input v-model="siteForm.name" placeholder="请输入网站名称" />
        </el-form-item>
        <el-form-item label="网站描述" prop="description">
          <el-input v-model="siteForm.description" placeholder="请输入网站描述" />
        </el-form-item>
        <el-form-item label="网站地址" prop="url">
          <el-input v-model="siteForm.url" placeholder="请输入网站地址" />
        </el-form-item>
        <el-form-item label="Logo地址" prop="logo">
          <el-input v-model="siteForm.logo" placeholder="请输入Logo地址（可选）" />
        </el-form-item>
        <el-form-item label="分类" prop="category">
          <el-select v-model="siteForm.category" placeholder="请选择分类">
            <el-option v-for="category in categories.slice(1)" :key="category" :label="category" :value="category" />
          </el-select>
        </el-form-item>
        <el-form-item label="主题色" prop="color">
          <el-color-picker v-model="siteForm.color" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">
            {{ isEditing ? '更新' : '添加' }}
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Monitor,
  Tools,
  DataAnalysis,
  Setting,
  Link,
  Compass,
  Collection,
  ChatLineSquare,
  Tickets,
  Histogram,
  Plus,
  Edit,
  Delete,
  MoreFilled,
  Box,
  Upload,
  Folder,
  Lock
} from '@element-plus/icons-vue'

// 网站数据接口
interface Site {
  id: number
  name: string
  description: string
  url: string
  logo?: string
  icon: any
  color: string
  category: string
  visitCount?: number
}

// 导航网站数据 - 使用可靠的图标源
const sites = ref<Site[]>([
  {
    id: 1,
    name: 'Prometheus',
    description: '开源监控和告警系统',
    url: 'https://prometheus.io/',
    logo: 'https://raw.githubusercontent.com/prometheus/prometheus/main/web/ui/static/img/prometheus_logo.svg',
    icon: Monitor,
    color: '#E6522C',
    category: '监控告警',
    visitCount: 0
  },
  {
    id: 2,
    name: 'Grafana',
    description: '数据可视化与监控平台',
    url: 'https://grafana.com/',
    logo: 'https://cdn.jsdelivr.net/gh/devicons/devicon/icons/grafana/grafana-original.svg',
    icon: DataAnalysis,
    color: '#F46800',
    category: '监控告警',
    visitCount: 0
  },
  {
    id: 3,
    name: 'Jenkins',
    description: '自动化构建与部署工具',
    url: 'https://www.jenkins.io/',
    logo: 'https://cdn.jsdelivr.net/gh/devicons/devicon/icons/jenkins/jenkins-original.svg',
    icon: Tools,
    color: '#D24939',
    category: 'CI/CD',
    visitCount: 0
  },
  {
    id: 4,
    name: 'Kubernetes',
    description: 'Kubernetes 官方文档',
    url: 'https://kubernetes.io/',
    logo: 'https://raw.githubusercontent.com/kubernetes/kubernetes/master/logo/logo.svg',
    icon: Setting,
    color: '#326CE5',
    category: '容器编排',
    visitCount: 0
  },
  {
    id: 5,
    name: 'Docker',
    description: '容器化平台官方网站',
    url: 'https://www.docker.com/',
    logo: 'https://cdn.jsdelivr.net/gh/devicons/devicon/icons/docker/docker-original.svg',
    icon: Box,
    color: '#2496ED',
    category: '容器编排',
    visitCount: 0
  },
  {
    id: 6,
    name: 'Elasticsearch',
    description: '分布式搜索和分析引擎',
    url: 'https://www.elastic.co/',
    logo: 'https://cdn.jsdelivr.net/gh/simple-icons/simple-icons/icons/elasticsearch.svg',
    icon: Collection,
    color: '#005571',
    category: '日志分析',
    visitCount: 0
  },
  {
    id: 7,
    name: 'Zabbix',
    description: '企业级监控解决方案',
    url: 'https://www.zabbix.com/',
    icon: Histogram,
    color: '#D81E1E',
    category: '监控告警',
    visitCount: 0
  },
  {
    id: 8,
    name: 'Nagios',
    description: '经典网络监控工具',
    url: 'https://www.nagios.org/',
    icon: Compass,
    color: '#F9A825',
    category: '监控告警',
    visitCount: 0
  },
  {
    id: 9,
    name: 'GitLab',
    description: '代码托管与CI/CD平台',
    url: 'https://gitlab.com/',
    logo: 'https://cdn.jsdelivr.net/gh/devicons/devicon/icons/gitlab/gitlab-original.svg',
    icon: Link,
    color: '#FC6D26',
    category: 'CI/CD',
    visitCount: 0
  },
  {
    id: 10,
    name: 'GitHub',
    description: 'GitHub 代码托管平台',
    url: 'https://github.com/',
    logo: 'https://cdn.jsdelivr.net/gh/devicons/devicon/icons/github/github-original.svg',
    icon: Upload,
    color: '#24292e',
    category: 'CI/CD',
    visitCount: 0
  },
  {
    id: 11,
    name: 'Ansible',
    description: '自动化配置管理工具',
    url: 'https://www.ansible.com/',
    logo: 'https://cdn.jsdelivr.net/gh/devicons/devicon/icons/ansible/ansible-original.svg',
    icon: Setting,
    color: '#EE0000',
    category: '配置管理',
    visitCount: 0
  },
  {
    id: 12,
    name: 'Terraform',
    description: '基础设施即代码工具',
    url: 'https://www.terraform.io/',
    logo: 'https://cdn.jsdelivr.net/gh/devicons/devicon/icons/terraform/terraform-original.svg',
    icon: Folder,
    color: '#623CE4',
    category: '配置管理',
    visitCount: 0
  },
  {
    id: 13,
    name: 'Redis',
    description: '内存数据结构存储',
    url: 'https://redis.io/',
    logo: 'https://cdn.jsdelivr.net/gh/devicons/devicon/icons/redis/redis-original.svg',
    icon: Link,
    color: '#DC382D',
    category: '数据库',
    visitCount: 0
  },
  {
    id: 14,
    name: 'MongoDB',
    description: 'NoSQL 文档数据库',
    url: 'https://www.mongodb.com/',
    logo: 'https://cdn.jsdelivr.net/gh/devicons/devicon/icons/mongodb/mongodb-original.svg',
    icon: Lock,
    color: '#47A248',
    category: '数据库',
    visitCount: 0
  },
  {
    id: 15,
    name: 'Nginx',
    description: '高性能Web服务器',
    url: 'https://nginx.org/',
    logo: 'https://cdn.jsdelivr.net/gh/devicons/devicon/icons/nginx/nginx-original.svg',
    icon: ChatLineSquare,
    color: '#009639',
    category: '服务器',
    visitCount: 0
  },
  {
    id: 16,
    name: 'Apache',
    description: 'Apache HTTP 服务器',
    url: 'https://httpd.apache.org/',
    logo: 'https://cdn.jsdelivr.net/gh/devicons/devicon/icons/apache/apache-original.svg',
    icon: Tickets,
    color: '#D22128',
    category: '服务器',
    visitCount: 0
  }
])

// 分类数据
const categories = ref(['全部', '监控告警', 'CI/CD', '容器编排', '日志分析', '配置管理', '数据库', '服务器'])
const selectedCategory = ref('全部')

// 对话框相关
const dialogVisible = ref(false)
const isEditing = ref(false)
const formRef = ref()

// 表单数据
const siteForm = ref<Partial<Site>>({
  name: '',
  description: '',
  url: '',
  logo: '',
  category: '',
  color: '#409EFF'
})

// 表单验证规则
const formRules = {
  name: [
    { required: true, message: '请输入网站名称', trigger: 'blur' }
  ],
  description: [
    { required: true, message: '请输入网站描述', trigger: 'blur' }
  ],
  url: [
    { required: true, message: '请输入网站地址', trigger: 'blur' }
  ],
  category: [
    { required: true, message: '请选择分类', trigger: 'change' }
  ]
}

// 搜索功能
const searchQuery = ref('')

// 过滤后的网站列表
const filteredSites = computed(() => {
  let result = sites.value

  // 按分类过滤
  if (selectedCategory.value !== '全部') {
    result = result.filter(site => site.category === selectedCategory.value)
  }

  // 按搜索关键词过滤
  if (searchQuery.value) {
    result = result.filter(site =>
      site.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      site.description.toLowerCase().includes(searchQuery.value.toLowerCase())
    )
  }

  return result
})

const filterSites = () => {
  // 触发计算属性更新
}

// 按分类过滤
const filterByCategory = (category: string) => {
  selectedCategory.value = category
}

// 跳转方法
const navigateTo = (url: string) => {
  // 增加访问计数
  const site = sites.value.find(s => s.url === url)
  if (site) {
    site.visitCount = (site.visitCount || 0) + 1
    // 保存到本地存储
    localStorage.setItem('navy-sites-data', JSON.stringify(sites.value))
  }
  window.open(url, '_blank')
}

// 显示添加对话框
const showAddDialog = () => {
  isEditing.value = false
  resetForm()
  dialogVisible.value = true
}

// 重置表单
const resetForm = () => {
  siteForm.value = {
    name: '',
    description: '',
    url: '',
    logo: '',
    category: '',
    color: '#409EFF'
  }
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

// 处理下拉菜单命令
const handleCommand = (command: { action: string; site: Site }) => {
  const { action, site } = command
  if (action === 'edit') {
    editSite(site)
  } else if (action === 'delete') {
    deleteSite(site)
  }
}

// 编辑网站
const editSite = (site: Site) => {
  isEditing.value = true
  siteForm.value = { ...site }
  dialogVisible.value = true
}

// 删除网站
const deleteSite = async (site: Site) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除 "${site.name}" 吗？`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )

    const index = sites.value.findIndex(s => s.id === site.id)
    if (index > -1) {
      sites.value.splice(index, 1)
      ElMessage.success('删除成功')
    }
  } catch {
    // 用户取消删除
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()

    if (isEditing.value) {
      // 更新网站
      const index = sites.value.findIndex(s => s.id === siteForm.value.id)
      if (index > -1) {
        sites.value[index] = {
          ...sites.value[index],
          ...siteForm.value,
          icon: getDefaultIcon(siteForm.value.category || '')
        }
        // 保存到本地存储
        localStorage.setItem('navy-sites-data', JSON.stringify(sites.value))
        ElMessage.success('更新成功')
      }
    } else {
      // 添加新网站
      const newSite: Site = {
        id: Date.now(),
        name: siteForm.value.name!,
        description: siteForm.value.description!,
        url: siteForm.value.url!,
        logo: siteForm.value.logo,
        icon: getDefaultIcon(siteForm.value.category!),
        color: siteForm.value.color!,
        category: siteForm.value.category!,
        visitCount: 0
      }
      sites.value.push(newSite)
      // 保存到本地存储
      localStorage.setItem('navy-sites-data', JSON.stringify(sites.value))
      ElMessage.success('添加成功')
    }

    dialogVisible.value = false
    resetForm()
  } catch {
    // 表单验证失败
  }
}

// 根据分类获取默认图标
const getDefaultIcon = (category: string) => {
  const iconMap: Record<string, any> = {
    '监控告警': Monitor,
    'CI/CD': Tools,
    '容器编排': Box,
    '日志分析': Collection,
    '配置管理': Setting,
    '服务发现': Link,
    '安全管理': Lock,
    '协作工具': ChatLineSquare,
    '项目管理': Tickets
  }
  return iconMap[category] || Link
}

// 初始化逻辑
onMounted(() => {
  console.log('运维导航页面已加载')
  // 从本地存储加载数据
  loadFromStorage()
})

// 从本地存储加载数据
const loadFromStorage = () => {
  try {
    const savedData = localStorage.getItem('navy-sites-data')
    if (savedData) {
      const parsedData = JSON.parse(savedData)
      // 合并保存的数据到默认数据，包括访问计数、主题色、Logo等所有用户修改
      sites.value.forEach(site => {
        const savedSite = parsedData.find((s: Site) => s.id === site.id)
        if (savedSite) {
          // 保存访问计数
          site.visitCount = savedSite.visitCount || 0
          
          // 保存用户修改的主题色
          if (savedSite.color && savedSite.color !== site.color) {
            site.color = savedSite.color
          }
          
          // 保存用户修改的Logo地址（重要：这里修复了Logo不保存的问题）
          if (savedSite.logo !== undefined && savedSite.logo !== site.logo) {
            site.logo = savedSite.logo
          }
          
          // 保存用户修改的其他字段
          if (savedSite.name && savedSite.name !== site.name) {
            site.name = savedSite.name
          }
          if (savedSite.description && savedSite.description !== site.description) {
            site.description = savedSite.description
          }
          if (savedSite.url && savedSite.url !== site.url) {
            site.url = savedSite.url
          }
          if (savedSite.category && savedSite.category !== site.category) {
            site.category = savedSite.category
          }
        }
      })

      // 加载用户添加的自定义网站
      const customSites = parsedData.filter((s: Site) => s.id > 1000) // 自定义网站ID大于1000
      customSites.forEach((customSite: Site) => {
        if (!sites.value.find(s => s.id === customSite.id)) {
          sites.value.push(customSite)
        }
      })
    }
  } catch (error) {
    console.error('加载本地数据失败:', error)
  }
}
</script>

<style lang="scss" scoped>
.navigation-container {
  padding: 24px;
  background: #f5f7fa;
  min-height: calc(100vh - 64px);
  font-family: 'Arial', sans-serif;

  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24px;
    padding: 20px;
    background: #fff;
    border-radius: 12px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);

    h1 {
      margin: 0;
      font-size: 28px;
      color: #1a202c;
      font-weight: 700;
      background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
    }

    .header-actions {
      display: flex;
      align-items: center;
      gap: 16px;
    }

    .search-input {
      width: 300px;

      :deep(.el-input__wrapper) {
        border-radius: 25px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
      }
    }
  }

  .category-tabs {
    margin-bottom: 24px;
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
    padding: 16px 20px;
    background: #fff;
    border-radius: 12px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);

    .category-tag {
      cursor: pointer;
      transition: all 0.3s ease;
      border-radius: 20px;
      padding: 8px 16px;
      font-weight: 500;

      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
      }
    }
  }

  .card-list {
    .nav-card {
      position: relative;
      background: #fff;
      border-radius: 16px;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
      margin-bottom: 20px;
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
      overflow: hidden;
      border: 1px solid rgba(0, 0, 0, 0.05);

      &:hover {
        transform: translateY(-8px);
        box-shadow: 0 12px 24px rgba(0, 0, 0, 0.15);
        border-color: rgba(102, 126, 234, 0.3);
      }

      .card-actions {
        position: absolute;
        top: 12px;
        right: 12px;
        z-index: 10;
        opacity: 0;
        transition: opacity 0.3s ease;

        .action-icon {
          width: 24px;
          height: 24px;
          cursor: pointer;
          color: #6b7280;
          transition: color 0.2s ease;

          &:hover {
            color: #374151;
          }
        }
      }

      &:hover .card-actions {
        opacity: 1;
      }

      .card-main {
        display: flex;
        align-items: flex-start;
        padding: 20px;
        cursor: pointer;
        height: 100%;
      }

      .card-icon {
        width: 60px;
        height: 60px;
        border-radius: 16px;
        display: flex;
        align-items: center;
        justify-content: center;
        margin-right: 20px;
        flex-shrink: 0;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
        transition: transform 0.3s ease;
        position: relative;
        overflow: hidden;

        .site-logo {
          width: 40px;
          height: 40px;
          object-fit: contain;
          background: rgba(255, 255, 255, 0.9);
          border-radius: 8px;
          padding: 6px;
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
          // 保持原始颜色，提高对比度
        }

        .el-icon {
          color: #fff !important;
        }
      }

      &:hover .card-icon {
        transform: scale(1.05);
      }

      .card-content {
        flex: 1;
        display: flex;
        flex-direction: column;
        justify-content: space-between;
        min-height: 80px;

        h3 {
          margin: 0 0 8px 0;
          font-size: 20px;
          color: #1a202c;
          font-weight: 600;
          line-height: 1.3;
        }

        p {
          margin: 0 0 12px 0;
          font-size: 14px;
          color: #6b7280;
          line-height: 1.5;
          flex-grow: 1;
        }

        .card-meta {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-top: auto;

          .visit-count {
            font-size: 12px;
            color: #9ca3af;
            font-weight: 500;
          }
        }
      }
    }
  }

  .el-empty {
    padding: 60px;
    background: #fff;
    border-radius: 16px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);

    :deep(.el-empty__description) {
      color: #6b7280;
      font-size: 16px;
    }
  }

  // 对话框样式优化
  :deep(.el-dialog) {
    border-radius: 16px;
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.15);

    .el-dialog__header {
      padding: 24px 24px 16px;
      border-bottom: 1px solid #f3f4f6;

      .el-dialog__title {
        font-size: 20px;
        font-weight: 600;
        color: #1a202c;
      }
    }

    .el-dialog__body {
      padding: 24px;
    }

    .el-dialog__footer {
      padding: 16px 24px 24px;
      border-top: 1px solid #f3f4f6;
    }
  }

  :deep(.el-form-item__label) {
    font-weight: 500;
    color: #374151;
  }

  :deep(.el-input__wrapper) {
    border-radius: 8px;
    transition: all 0.3s ease;

    &:hover {
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    }
  }

  :deep(.el-select .el-input__wrapper) {
    border-radius: 8px;
  }

  @media (max-width: 768px) {
    padding: 16px;

    .header {
      flex-direction: column;
      align-items: stretch;
      gap: 16px;
      padding: 16px;

      h1 {
        font-size: 24px;
        text-align: center;
      }

      .header-actions {
        flex-direction: column;
        gap: 12px;
      }

      .search-input {
        width: 100%;
      }
    }

    .category-tabs {
      padding: 12px 16px;

      .category-tag {
        padding: 6px 12px;
        font-size: 12px;
      }
    }

    .card-list {
      .nav-card {
        margin-bottom: 16px;

        .card-main {
          padding: 16px;
        }

        .card-icon {
          width: 50px;
          height: 50px;
          margin-right: 16px;

          .site-logo {
            width: 30px;
            height: 30px;
          }
        }

        .card-content {
          min-height: 70px;

          h3 {
            font-size: 18px;
          }

          p {
            font-size: 13px;
          }

          .card-meta {
            .visit-count {
              font-size: 11px;
            }
          }
        }
      }
    }

    :deep(.el-dialog) {
      width: 90% !important;
      margin: 5vh auto;
    }
  }

  @media (max-width: 480px) {
    .category-tabs {
      .category-tag {
        padding: 4px 8px;
        font-size: 11px;
      }
    }

    .card-list {
      .nav-card {
        .card-main {
          flex-direction: column;
          text-align: center;
        }

        .card-icon {
          margin: 0 auto 12px;
        }

        .card-content {
          .card-meta {
            justify-content: center;
            gap: 12px;
          }
        }
      }
    }
  }
}
</style>