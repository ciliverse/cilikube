<template>
    <div class="navigation-container">
      <!-- 页面标题 -->
      <div class="header">
        <h1>运维导航</h1>
        <el-input
          v-model="searchQuery"
          placeholder="搜索运维网站..."
          clearable
          prefix-icon="Search"
          class="search-input"
          @input="filterSites"
        />
      </div>
  
      <!-- 导航卡片列表 -->
      <div class="card-list">
        <el-row :gutter="20">
          <el-col
            :xs="24"
            :sm="12"
            :md="8"
            :lg="6"
            v-for="(site, index) in filteredSites"
            :key="index"
          >
            <div class="nav-card" @click="navigateTo(site.url)">
              <div class="card-icon" :style="{ backgroundColor: site.color }">
                <el-icon :size="30" color="#fff">
                  <component :is="site.icon" />
                </el-icon>
              </div>
              <div class="card-content">
                <h3>{{ site.name }}</h3>
                <p>{{ site.description }}</p>
              </div>
            </div>
          </el-col>
        </el-row>
      </div>
  
      <!-- 空状态提示 -->
      <el-empty v-if="filteredSites.length === 0" description="未找到匹配的网站" />
    </div>
  </template>
  
  <script setup lang="ts">
  import { ref, computed, onMounted } from 'vue'
  import {
    Search,
    Monitor,
    Tools,
    DataAnalysis,
    Setting,
    Link,
    Compass,
    Collection,
    ChatLineSquare,
    Tickets,
    Histogram
  } from '@element-plus/icons-vue'
  
  // 导航网站数据
  const sites = ref([
    {
      name: 'Prometheus',
      description: '开源监控和告警系统',
      url: 'https://prometheus.io/',
      icon: Monitor,
      color: '#E6522C'
    },
    {
      name: 'Grafana',
      description: '数据可视化与监控平台',
      url: 'https://grafana.com/',
      icon: DataAnalysis,
      color: '#F46800'
    },
    {
      name: 'Jenkins',
      description: '自动化构建与部署工具',
      url: 'https://www.jenkins.io/',
      icon: Tools,
      color: '#D24939'
    },
    {
      name: 'Kubernetes Dashboard',
      description: 'Kubernetes 集群管理界面',
      url: 'https://kubernetes.io/docs/tasks/access-application-cluster/web-ui-dashboard/',
      icon: Setting,
      color: '#326CE5'
    },
    {
      name: 'ELK Stack',
      description: '日志收集与分析平台',
      url: 'https://www.elastic.co/elastic-stack/',
      icon: Collection,
      color: '#00BFB3'
    },
    {
      name: 'Zabbix',
      description: '企业级监控解决方案',
      url: 'https://www.zabbix.com/',
      icon: Histogram,
      color: '#D81E1E'
    },
    {
      name: 'Nagios',
      description: '经典网络监控工具',
      url: 'https://www.nagios.org/',
      icon: Compass,
      color: '#F9A825'
    },
    {
      name: 'GitLab',
      description: '代码托管与CI/CD平台',
      url: 'https://gitlab.com/',
      icon: Link,
      color: '#FC6D26'
    },
    {
      name: 'Slack',
      description: '团队协作与通知工具',
      url: 'https://slack.com/',
      icon: ChatLineSquare,
      color: '#4A154B'
    },
    {
      name: 'Jira',
      description: '问题跟踪与项目管理',
      url: 'https://www.atlassian.com/software/jira',
      icon: Tickets,
      color: '#0052CC'
    }
  ])
  
  // 搜索功能
  const searchQuery = ref('')
  const filteredSites = computed(() => {
    if (!searchQuery.value) return sites.value
    return sites.value.filter(site =>
      site.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      site.description.toLowerCase().includes(searchQuery.value.toLowerCase())
    )
  })
  
  const filterSites = () => {
    // 触发计算属性更新
  }
  
  // 跳转方法
  const navigateTo = (url: string) => {
    window.open(url, '_blank')
  }
  
  // 初始化逻辑
  onMounted(() => {
    console.log('导航页面已加载')
  })
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
      border-radius: 8px;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  
      h1 {
        margin: 0;
        font-size: 24px;
        color: #212529;
        font-weight: 600;
      }
  
      .search-input {
        width: 300px;
        :deep(.el-input__inner) {
          border-radius: 20px;
        }
      }
    }
  
    .card-list {
      .nav-card {
        display: flex;
        align-items: center;
        padding: 16px;
        background: #fff;
        border-radius: 8px;
        box-shadow: 0 2px 6px rgba(0, 0, 0, 0.05);
        margin-bottom: 20px;
        cursor: pointer;
        transition: all 0.3s ease;
  
        &:hover {
          transform: translateY(-4px);
          box-shadow: 0 6px 12px rgba(0, 0, 0, 0.1);
        }
  
        .card-icon {
          width: 50px;
          height: 50px;
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          margin-right: 16px;
          flex-shrink: 0;
        }
  
        .card-content {
          flex: 1;
  
          h3 {
            margin: 0 0 8px 0;
            font-size: 18px;
            color: #343a40;
            font-weight: 500;
          }
  
          p {
            margin: 0;
            font-size: 14px;
            color: #6c757d;
            line-height: 1.4;
          }
        }
      }
    }
  
    .el-empty {
      padding: 40px;
      background: #fff;
      border-radius: 8px;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
    }
  
    @media (max-width: 768px) {
      .header {
        flex-direction: column;
        align-items: stretch;
        gap: 16px;
  
        .search-input {
          width: 100%;
        }
      }
  
      .card-list {
        .nav-card {
          margin-bottom: 16px;
        }
      }
    }
  }
  </style>