<template>
  <div class="stat-card" :class="{ 'stat-card-active': active }">
    <div class="stat-icon" :style="{ background: iconBg }">
      <component :is="icon" class="icon" />
    </div>
    <div class="stat-content">
      <div class="stat-value">{{ value }}</div>
      <div class="stat-label">{{ label }}</div>
      <div v-if="trend" class="stat-trend" :class="trendClass">
        <el-icon><CaretTop v-if="trend > 0" /><CaretBottom v-else /></el-icon>
        <span>{{ Math.abs(trend) }}%</span>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { CaretTop, CaretBottom } from '@element-plus/icons-vue'

interface Props {
  icon: any
  label: string
  value: number | string
  iconBg?: string
  trend?: number
  active?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  iconBg: 'linear-gradient(135deg, #409EFF 0%, #3A8FE8 100%)',
  trend: 0,
  active: false
})

const trendClass = computed(() => ({
  'trend-up': props.trend > 0,
  'trend-down': props.trend < 0
}))
</script>

<style scoped>
.stat-card {
  position: relative;
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 24px;
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 2px 12px rgba(64, 158, 255, 0.08);
  transition: all 0.3s ease;
  overflow: hidden;
  border: 1px solid #E8F4FF;
}

.stat-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, #409EFF 0%, #66B1FF 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(64, 158, 255, 0.15);
  border-color: #409EFF;
}

.stat-card:hover::before {
  opacity: 1;
}

.stat-card-active {
  border: 2px solid #409EFF;
  box-shadow: 0 4px 20px rgba(64, 158, 255, 0.25);
  background: linear-gradient(135deg, #FFFFFF 0%, #F0F9FF 100%);
}

.stat-icon {
  width: 64px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 16px;
  flex-shrink: 0;
}

.icon {
  width: 32px;
  height: 32px;
  color: #fff;
}

.stat-content {
  flex: 1;
  min-width: 0;
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  color: #1D2129;
  line-height: 1.2;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #4E5969;
  margin-bottom: 8px;
}

.stat-trend {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 12px;
}

.trend-up {
  color: #203316;
  background: rgba(103, 194, 58, 0.1);
}

.trend-down {
  color: #f56c6c;
  background: rgba(245, 108, 108, 0.1);
}

@media (max-width: 768px) {
  .stat-card {
    padding: 16px;
    gap: 12px;
  }

  .stat-icon {
    width: 48px;
    height: 48px;
  }

  .icon {
    width: 24px;
    height: 24px;
  }

  .stat-value {
    font-size: 24px;
  }
}
</style>
