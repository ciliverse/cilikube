<template>
  <div class="cluster-card" :class="{ 'is-active': isActive, 'is-unavailable': !isAvailable }">
    <div class="card-header">
      <div class="cluster-info">
        <div class="cluster-icon" :style="{ background: providerBg }">
          <component :is="providerIcon" class="icon" />
        </div>
        <div class="cluster-meta">
          <div class="cluster-name">
            {{ cluster.name }}
            <el-tag v-if="isActive" type="success" size="small" effect="dark">Active</el-tag>
          </div>
          <div class="cluster-server">{{ cluster.server }}</div>
        </div>
      </div>
      <el-dropdown trigger="click" @command="handleCommand">
        <el-button circle :icon="MoreFilled" size="small" />
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="setActive" :disabled="isActive || !isAvailable">
              <el-icon><Check /></el-icon>
              设为活动
            </el-dropdown-item>
            <el-dropdown-item command="edit" disabled>
              <el-icon><Edit /></el-icon>
              编辑
            </el-dropdown-item>
            <el-dropdown-item command="delete" :disabled="cluster.source === 'file'" divided>
              <el-icon><Delete /></el-icon>
              删除
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>

    <div class="card-body">
      <div class="cluster-stats">
        <div class="stat-item">
          <el-icon class="stat-icon"><Platform /></el-icon>
          <div class="stat-content">
            <div class="stat-label">版本</div>
            <div class="stat-value">{{ cluster.version || 'N/A' }}</div>
          </div>
        </div>
        <div class="stat-item">
          <el-icon class="stat-icon"><Location /></el-icon>
          <div class="stat-content">
            <div class="stat-label">环境</div>
            <div class="stat-value">
              <el-tag :type="environmentType" size="small">{{ cluster.environment || 'default' }}</el-tag>
            </div>
          </div>
        </div>
      </div>

      <div class="cluster-status">
        <div class="status-indicator" :class="statusClass">
          <span class="status-dot"></span>
          <span class="status-text">{{ cluster.status }}</span>
        </div>
        <el-tag :type="cluster.source === 'database' ? 'primary' : 'info'" size="small" effect="plain">
          {{ cluster.source === 'database' ? '数据库' : '文件' }}
        </el-tag>
      </div>
    </div>

    <div v-if="!isActive && isAvailable" class="card-footer">
      <el-button type="primary" size="small" :icon="Check" @click="handleSetActive" block>
        切换到此集群
      </el-button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { ElMessageBox } from 'element-plus'
import {
  MoreFilled,
  Check,
  Edit,
  Delete,
  Platform,
  Location
} from '@element-plus/icons-vue'
import type { ClusterInfo } from '@/api/cluster'

interface Props {
  cluster: ClusterInfo
  activeClusterName: string
}

interface Emits {
  (e: 'setActive', cluster: ClusterInfo): void
  (e: 'delete', cluster: ClusterInfo): void
  (e: 'edit', cluster: ClusterInfo): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const isActive = computed(() => props.cluster.name === props.activeClusterName)
const isAvailable = computed(() => props.cluster.status.startsWith('可用'))

const statusClass = computed(() => ({
  'status-available': isAvailable.value,
  'status-unavailable': !isAvailable.value
}))

const environmentType = computed(() => {
  const env = props.cluster.environment?.toLowerCase()
  if (env === 'production' || env === 'prod') return 'danger'
  if (env === 'staging' || env === 'stage') return 'warning'
  if (env === 'development' || env === 'dev') return 'success'
  return 'info'
})

const providerIcon = computed(() => {
  // 根据不同的提供商返回不同的图标
  return Platform
})

const providerBg = computed(() => {
  const provider = props.cluster.environment?.toLowerCase() || 'default'
  const gradients: Record<string, string> = {
    production: 'linear-gradient(135deg, #409EFF 0%, #3A8FE8 100%)',
    staging: 'linear-gradient(135deg, #66B1FF 0%, #5DADE2 100%)',
    development: 'linear-gradient(135deg, #79BBFF 0%, #A0CFFF 100%)',
    testing: 'linear-gradient(135deg, #ECF5FF 0%, #C6E2FF 100%)',
    default: 'linear-gradient(135deg, #409EFF 0%, #66B1FF 100%)'
  }
  return gradients[provider] || gradients.default
})

const handleCommand = (command: string) => {
  switch (command) {
    case 'setActive':
      handleSetActive()
      break
    case 'edit':
      emit('edit', props.cluster)
      break
    case 'delete':
      handleDelete()
      break
  }
}

const handleSetActive = () => {
  emit('setActive', props.cluster)
}

const handleDelete = () => {
  ElMessageBox.confirm(
    `确定要删除集群 "${props.cluster.name}" 吗？此操作不可恢复。`,
    '删除集群',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(() => {
    emit('delete', props.cluster)
  }).catch(() => {
    // 用户取消
  })
}
</script>

<style scoped>
.cluster-card {
  background: #fff;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(64, 158, 255, 0.08);
  transition: all 0.3s ease;
  border: 2px solid #E8F4FF;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.cluster-card:hover {
  box-shadow: 0 8px 24px rgba(64, 158, 255, 0.15);
  transform: translateY(-4px);
  border-color: #409EFF;
}

.cluster-card.is-active {
  border-color: #409EFF;
  background: linear-gradient(135deg, #FFFFFF 0%, #F0F9FF 100%);
  box-shadow: 0 4px 20px rgba(64, 158, 255, 0.2);
}

.cluster-card.is-unavailable {
  opacity: 0.7;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}

.cluster-info {
  display: flex;
  gap: 12px;
  flex: 1;
  min-width: 0;
}

.cluster-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.cluster-icon .icon {
  width: 24px;
  height: 24px;
  color: #fff;
}

.cluster-meta {
  flex: 1;
  min-width: 0;
}

.cluster-name {
  font-size: 18px;
  font-weight: 600;
  color: #1D2129;
  margin-bottom: 4px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.cluster-server {
  font-size: 13px;
  color: #4E5969;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.card-body {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.cluster-stats {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.stat-item {
  display: flex;
  gap: 8px;
  padding: 12px;
  background: #F7FBFF;
  border-radius: 8px;
  border: 1px solid #E8F4FF;
}

.stat-icon {
  font-size: 20px;
  color: #409EFF;
  flex-shrink: 0;
}

.stat-content {
  flex: 1;
  min-width: 0;
}

.stat-label {
  font-size: 12px;
  color: #86909C;
  margin-bottom: 2px;
}

.stat-value {
  font-size: 14px;
  font-weight: 600;
  color: #1D2129;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.cluster-status {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid #E8F4FF;
}

.status-indicator {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  font-weight: 500;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

.status-available .status-dot {
  background: #67c23a;
}

.status-available .status-text {
  color: #67c23a;
}

.status-unavailable .status-dot {
  background: #f56c6c;
}

.status-unavailable .status-text {
  color: #f56c6c;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

.card-footer {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #E8F4FF;
}

@media (max-width: 768px) {
  .cluster-stats {
    grid-template-columns: 1fr;
  }
}
</style>
