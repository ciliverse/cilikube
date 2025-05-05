<template>
    <el-card class="summary-card" shadow="never" v-loading="loading">
      <template #header>
        <div class="card-header">
          <el-icon><DataAnalysis /></el-icon>
          <span>集群资源概览</span>
           <el-tooltip content="刷新" placement="top">
              <el-button :icon="Refresh" circle text size="small" @click="fetchSummary" :loading="loading" class="refresh-btn"/>
           </el-tooltip>
        </div>
      </template>
      <div v-if="error" class="error-message">
        <el-alert type="error" :closable="false" show-icon>
          加载资源概览失败: {{ error }}
        </el-alert>
      </div>
      <div v-else-if="!summaryData" class="loading-placeholder">
        正在加载资源数量...
      </div>
      <el-row v-else :gutter="15" class="summary-grid">
        <el-col v-for="item in displayItems" :key="item.key" :xs="12" :sm="8" :md="6" :lg="4">
          <div class="summary-item">
            <div class="item-icon" :style="{ backgroundColor: item.color + '1A' }">
              <el-icon :size="24" :color="item.color">
                <component :is="item.icon" />
              </el-icon>
            </div>
            <div class="item-content">
              <div class="item-label">{{ item.label }}</div>
              <div class="item-value">{{ formatCount(summaryData[item.key]) }}</div>
            </div>
          </div>
        </el-col>
         <el-col v-if="displayItems.length === 0 && !loading" :span="24">
              <el-empty description="暂无资源数据" :image-size="60" />
         </el-col>
      </el-row>
    </el-card>
  </template>
  
  <script setup lang="ts">
  import { ref, onMounted, computed } from 'vue';
  import { ElMessage } from 'element-plus';
  import { request } from '@/utils/service'; // Adjust path
  import {
      DataAnalysis, Refresh, Platform, CollectionTag, Files, TakeawayBox, Service,
      Coin, MessageBox, SetUp, Notebook, Key, Lock, Connection, Warning // Added more icons
  } from '@element-plus/icons-vue';
  
  interface ResourceSummaryData {
    nodes?: number | null;
    namespaces?: number | null;
    pods?: number | null;
    deployments?: number | null;
    services?: number | null;
    persistentVolumes?: number | null;
    pvcs?: number | null;
    statefulSets?: number | null;
    daemonSets?: number | null;
    configMaps?: number | null;
    secrets?: number | null;
    ingresses?: number | null;
    // Add keys corresponding to the Go model
  }
  
  interface DisplayConfig {
      key: keyof ResourceSummaryData;
      label: string;
      icon: any; // Vue component type
      color: string;
  }
  
  const loading = ref(false);
  const summaryData = ref<ResourceSummaryData | null>(null);
  const error = ref<string | null>(null);
  
  // Define the display order, labels, icons, and colors
  const displayConfig: DisplayConfig[] = [
      { key: 'nodes', label: '节点', icon: Platform, color: '#409EFF' },
      { key: 'namespaces', label: '命名空间', icon: CollectionTag, color: '#67C23A' },
      { key: 'pods', label: 'Pods', icon: Files, color: '#E6A23C' },
      { key: 'deployments', label: 'Deployments', icon: TakeawayBox, color: '#F56C6C' },
      { key: 'statefulSets', label: 'StatefulSets', icon: SetUp, color: '#a774d1' }, // Purple
      { key: 'daemonSets', label: 'DaemonSets', icon: Notebook, color: '#7f8c8d' }, // Grey
      { key: 'services', label: 'Services', icon: Service, color: '#3498DB' }, // Blue
      { key: 'ingresses', label: 'Ingresses', icon: Connection, color: '#1ABC9C' }, // Turquoise
      { key: 'persistentVolumes', label: 'PVs', icon: Coin, color: '#f39c12' }, // Orange
      { key: 'pvcs', label: 'PVCs', icon: MessageBox, color: '#f1c40f' }, // Yellow
      { key: 'configMaps', label: 'ConfigMaps', icon: Key, color: '#2ecc71' }, // Emerald
      { key: 'secrets', label: 'Secrets', icon: Lock, color: '#e74c3c' }, // Red
      // Add more resources here
  ];
  
  // Filter config based on available data keys from backend
  const displayItems = computed(() => {
      if (!summaryData.value) return [];
      return displayConfig.filter(item => summaryData.value?.[item.key] !== undefined && summaryData.value?.[item.key] !== null);
  });
  
  const VITE_API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://192.168.1.100:8080";
  const fetchSummary = async () => {
    loading.value = true;
    error.value = null;
    summaryData.value = null; // Clear previous data
    try {
      const response = await request<{ code: number; data: ResourceSummaryData; message: string }>({
        url: '/api/v1/summary/resources', // Match your Go route
        method: 'get',
        baseURL: VITE_API_BASE_URL // If needed
        // VITE_API_BASE_URL/api/v1/summary/resources
      });
      if (response.code === 200 && response.data) {
        summaryData.value = response.data;
      } else {
        throw new Error(response.message || '获取数据格式错误');
      }
    } catch (err: any) {
      console.error("Failed to fetch resource summary:", err);
      error.value = err.message || '网络请求失败';
      // Optionally show ElMessage
      // ElMessage.error(`加载资源概览失败: ${error.value}`);
    } finally {
      loading.value = false;
    }
  };
  
  const formatCount = (count: number | null | undefined): string | number => {
    // Backend uses pointers, so check for null explicitly
    if (count === null || count === undefined) {
      return '?'; // Indicate data wasn't fetched or error occurred for this type
    }
    return count;
  };
  
  onMounted(() => {
    fetchSummary();
  });
  
  // Expose fetchSummary if you want to call it from parent
  // defineExpose({ fetchSummary });
  
  </script>
  
  <style lang="scss" scoped>
  .summary-card {
    border: 1px solid var(--el-border-color-lighter);
    background-color: var(--el-bg-color);
  
    .card-header {
      display: flex;
      align-items: center;
      font-size: 16px;
      font-weight: 600;
      color: var(--el-text-color-primary);
      gap: 8px; // Space between icon and text
  
       .refresh-btn {
          margin-left: auto; // Push refresh button to the right
       }
    }
  
    .error-message {
      padding: 10px 0;
    }
    .loading-placeholder {
       padding: 20px;
       text-align: center;
       color: var(--el-text-color-secondary);
    }
  
    .summary-grid {
      padding-top: 5px; // Small space below header
    }
  
    .summary-item {
      display: flex;
      align-items: center;
      padding: 12px 8px;
      margin-bottom: 10px;
      background-color: var(--el-bg-color); // Match card background or slightly lighter/darker
      border-radius: 4px;
      // border: 1px solid var(--el-border-color-extra-light);
       transition: background-color 0.2s ease;
       &:hover {
           background-color: var(--el-fill-color-lighter);
       }
    }
  
    .item-icon {
      flex-shrink: 0;
      width: 48px;
      height: 48px;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      margin-right: 12px;
      // Background color is set inline
    }
  
    .item-content {
      overflow: hidden; // Prevent text overflow
    }
  
    .item-label {
      font-size: 13px;
      color: var(--el-text-color-secondary);
      margin-bottom: 4px;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
  
    .item-value {
      font-size: 20px;
      font-weight: 600;
      color: var(--el-text-color-primary);
      line-height: 1.2;
    }
  }
  </style>