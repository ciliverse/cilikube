// ./mockData/kubernetesMock.ts (Create this file)
import {
    HealthStatusItem, AlertSummaryItem, OverviewDataItem, NodeData,
    EventItem, EventStatisticItem, ResourceAlertItem
  } from '@/types/kubernetes'; // Adjust path if needed
  import {
    DataBoard as DataBoardIcon,
    Collection as CollectionIcon,
    Box as BoxIcon,
    Connection as ConnectionIcon,
  } from '@element-plus/icons-vue' // Import icons here or pass them as strings
  
  export const getMockNamespaces = (): string[] => ['default', 'kube-system', 'monitoring', 'logging', 'dev', 'prod'];
  
  export const getMockHealthStatus = (): HealthStatusItem[] => ([
    { label: '集群状态', value: '健康 (Mock)', color: '#67C23A' },
    { label: '节点在线', value: '4/5 (Mock)', color: '#E6A23C' },
    { label: 'Pod运行', value: '148/180 (Mock)', color: '#409EFF' },
    { label: '警报', value: '3 (Mock)', color: '#F56C6C' }
  ]);
  
  export const getMockAlertSummary = (): AlertSummaryItem[] => ([
    { label: '严重警报 (Mock)', count: 1, color: '#F56C6C', trend: 12 },
    { label: '警告 (Mock)', count: 3, color: '#E6A23C', trend: -5 },
    { label: '通知 (Mock)', count: 8, color: '#909399', trend: 3 }
  ]);
  
  // If passing icons directly:
  // export const getMockOverviewData = (): OverviewDataItem[] => ([
  //   { title: '集群节点 (Mock)', value: '5', percent: 95, total: '', icon: DataBoardIcon, color: '#409EFF' },
  //   { title: '命名空间 (Mock)', value: '6', percent: 75, total: '8', icon: CollectionIcon, color: '#67C23A' },
  //   { title: '运行Pods (Mock)', value: '148', percent: 82, total: '180', icon: BoxIcon, color: '#E6A23C' },
  //   { title: '服务数量 (Mock)', value: '23', percent: 65, total: '35', icon: ConnectionIcon, color: '#F56C6C' }
  // ]);
  
  // If passing icon names as strings:
  export const getMockOverviewData = (): OverviewDataItem[] => ([
      // Adjust the 'icon' property to store a string identifier if you prefer
      // You'll need a way to map this string back to the component in the template
      { title: '集群节点 (Mock)', value: '5', percent: 95, total: '', icon: 'DataBoardIcon', color: '#409EFF' },
      { title: '命名空间 (Mock)', value: '6', percent: 75, total: '8', icon: 'CollectionIcon', color: '#67C23A' },
      { title: '运行Pods (Mock)', value: '148', percent: 82, total: '180', icon: 'BoxIcon', color: '#E6A23C' },
      { title: '服务数量 (Mock)', value: '23', percent: 65, total: '35', icon: 'ConnectionIcon', color: '#F56C6C' }
  ]);
  
  
  // Mock chart data functions (returning just the data part)
  export const getMockCpuUsageChartData = () => ({
    labels: ['Master-1 M', 'Node-1 M', 'Node-2 M', 'Node-3 M', 'Node-4 M'],
    values: [65, 45, 30, 85, 60]
  });
  
  export const getMockMemoryUsageChartData = () => ({
    labels: ['Master-1 M', 'Node-1 M', 'Node-2 M', 'Node-3 M', 'Node-4 M'],
    // Structure for stacked bar - an array of arrays
    values: [
      [6, 10, 8, 15, 12], // Used
      [3, 4, 2, 5, 3],   // Cached
      [15, 10, 14, 4, 9] // Available
    ]
  });
  
  export const getMockStorageUsageChartData = () => ({
      labels: ['Master-1 M', 'Node-1 M', 'Node-2 M', 'Node-3 M', 'Node-4 M'],
      // Structure for stacked area - an array of arrays
      values: [
        [120, 80, 90, 150, 100], // Used
        [200, 200, 200, 200, 200] // Total (can be calculated or fetched separately if API provides only used)
      ]
  });
  
  export const getMockNetworkUsageChartData = () => ({
    labels: ['08:00 M', '10:00 M', '12:00 M', '14:00 M', '16:00 M', '18:00 M', '20:00 M'],
    values: [
        [120, 200, 150, 180, 210, 190, 230], // Inbound
        [100, 170, 130, 150, 180, 160, 190]  // Outbound
    ]
  });
  
  
  export const getMockNodes = (): NodeData[] => ([
    { name: 'master-1-mock', role: 'master', status: 'Ready', cpuUsage: 65, memoryUsage: 75, totalPods: 64, runningPods: 60, isNew: false },
    { name: 'node-1-mock', role: 'worker', status: 'Ready', cpuUsage: 45, memoryUsage: 60, totalPods: 128, runningPods: 110, isNew: true },
    { name: 'node-2-mock', role: 'worker', status: 'Ready', cpuUsage: 30, memoryUsage: 55, totalPods: 128, runningPods: 98, isNew: false },
    { name: 'node-3-mock', role: 'worker', status: 'NotReady', cpuUsage: 85, memoryUsage: 90, totalPods: 128, runningPods: 128, isNew: false },
    { name: 'node-4-mock', role: 'worker', status: 'Ready', cpuUsage: 60, memoryUsage: 65, totalPods: 128, runningPods: 105, isNew: false }
  ]);
  
  export const getMockRecentEvents = (): EventItem[] => ([
    { timestamp: '2023-10-15T09:23:17Z', type: 'Warning', object: 'pod/nginx-mock-58xj7', namespace: 'default', reason: 'FailedScheduling', message: 'Mock: 0/4 nodes available: 3 Insufficient cpu, 1 Insufficient memory.' },
    { timestamp: '2023-10-15T09:20:45Z', type: 'Normal', object: 'deployment/nginx-mock', namespace: 'default', reason: 'ScalingReplicaSet', message: 'Mock: Scaled up replica set nginx-deployment-75675f5897 to 1' },
    // ... add more mock events if needed
  ]);
  
  export const getMockEventStatistics = (): EventStatisticItem[] => ([
    { type: 'Normal (Mock)', count: 42, percentage: 70, trend: 'up' },
    { type: 'Warning (Mock)', count: 18, percentage: 30, trend: 'down' },
    // ... add more mock stats
  ]);
  
  // For the Pie chart based on statistics
  export const getMockEventStatisticsChartData = (stats: EventStatisticItem[]) => {
    // Example: Process the stats to fit the pie chart data structure
    const pieData = stats
      .filter(s => ['Normal (Mock)', 'Warning (Mock)', 'Critical (Mock)', 'Info (Mock)'].includes(s.type)) // Filter relevant types
      .map(s => {
          let color = '#909399'; // Default
          if (s.type.includes('Normal')) color = '#67C23A';
          if (s.type.includes('Warning')) color = '#E6A23C';
          if (s.type.includes('Critical')) color = '#F56C6C';
          return {
              value: s.count,
              name: `${s.type} (${s.percentage}%)`,
              itemStyle: { color: color }
          };
      });
      // Add dummy data if needed for the pie chart structure shown in the original code
      if (!pieData.find(p => p.name.includes('Critical'))) pieData.push({ value: 5, name: '严重事件 (10%) (Mock)', itemStyle: { color: '#F56C6C' }});
      if (!pieData.find(p => p.name.includes('Info'))) pieData.push({ value: 10, name: '通知事件 (20%) (Mock)', itemStyle: { color: '#909399' }});
  
      return pieData;
  };
  
  
  export const getMockResourceAlerts = (): ResourceAlertItem[] => ([
    { name: 'node-3 CPU过载 (Mock)', message: 'Mock: 节点node-3 CPU使用率达85%，超过阈值80%', level: 'critical', levelText: '严重', time: '5分钟前' },
    { name: 'node-3 内存不足 (Mock)', message: 'Mock: 节点node-3 内存使用率达90%，超过阈值85%', level: 'warning', levelText: '警告', time: '10分钟前' },
  ]);