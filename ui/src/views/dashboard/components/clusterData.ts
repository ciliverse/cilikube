<script lang="ts">
// 这是应该在单独文件中定义的模拟数据
export const mockClusterData = () => {
  const namespaces = ['default', 'kube-system', 'monitoring', 'logging', 'dev', 'prod']
  
  const overviewData = [
    {
      title: '集群节点',
      value: '5',
      percent: 95,
      total: '',
      icon: 'DataBoard',
      color: '#409EFF'
    },
    {
      title: '命名空间',
      value: '6',
      percent: 75,
      total: '8',
      icon: 'Collection',
      color: '#67C23A'
    },
    {
      title: '运行Pods',
      value: '148',
      percent: 82,
      total: '180',
      icon 'Box',
      color: '#E6A23C'
    },
    {
      title: '服务数量',
      value: '23',
      percent: 65,
      total: '35',
      icon: 'Connection',
      color: '#F56C6C'
    }
  ]
  
  const cpuUsageOption = {
    tooltip: {
      trigger: 'item'
    },
    legend: {
      top: '5%',
      left: 'center'
    },
    series: [
      {
        name: 'CPU使用率',
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: '18',
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: [
          { value: 70, name: '已使用 (70%)' },
          { value: 30, name: '剩余 (30%)' }
        ],
        color: ['#F566C', '#E4E7ED']
      }
    ]
  }
  
  const memoryUsageOption = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'value',
      boundaryGap: [0, 0.01],
      axisLine: {
        show: false
      },
      axisTick: {
        show: false
      }
    },
    yAxis: {
      type: 'category',
     : ['总内存', '已使用', '缓存', '可用'],
      axisLine: {
        show: false
      },
      axisTick: {
        show: false
      }
    },
    series: [
      {
        name: '内存分配',
        type: 'bar',
        data: [
          { value: 64, itemStyle: { color: '#E4E7ED' } },
          { value: 45, itemStyle: { color: '#F56C6C' } },
          { value: 12, itemStyle: { color: '#E6A23C' } },
          { value: 19, itemStyle: { color: '#67C23A' } }
        ],
        label: {
          show: true,
          position: 'right',
          formatter: '{@value} GB'
        }
      }
    ]
  }
  
  const storageUsageOption = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type 'cross',
        label: {
          backgroundColor: '#6a7985'
        }
      }
    },
    legend: {
      data: ['已使用', '总量']
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
 xAxis: [
      {
        type: 'category',
        boundaryGap: false,
        data: ['Master-1', 'Node-1', 'Node-2', 'Node-3']
      }
    ],
    yAxis: [
      {
        type: '',
        name: '存储 (GB)'
      }
    ],
    series: [
      {
        name: '已使用',
        type: 'line',
        stack: '总量',
        areaStyle: {
          color: '#F56C6C'
        },
        emphasis: {
          focus: 'series'
        },
        data: [120, 80, 90, 100]
      },
      {
        name: '总量',
        type: 'line',
        stack: '总量',
        areaStyle: {
          color: '#E4E7ED'
        },
        emphasis: {
          focus: 'series'
        },
        data: [200, 200, 200, 200]
      }
    ]
  }
  
  const networkUsageOption = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: ['08:00', '10:00', '12:00', '14:00', '16:00', '18:00', '20:00'],
      axisLine: {
        show: false
      },
      axisTick: {
        show: false
      }
    },
    yAxis: {
      type: 'value',
      name: '网络流量 (Mbps)',
      axisLine: {
        show: false
      },
      axisTick: {
        show: false
      }
    },
    series: [
      {
        name: '入站',
        type: 'bar',
        stack: '流量',
        data: [120, 200, 150, 180, 210, 190, 230],
        itemStyle: {
          color: '#409EFF'
        }
      },
      {
        name: '出站',
        type: 'bar',
        stack: '流量',
        data: [100, 170, 130, 150, 180, 160, 190],
        itemStyle: {
          color: '#67C23A'
        }
      }
    ]
  }
  
  const nodes = [
    {
      name: 'master-1',
      role: 'master',
      status: 'Ready',
      cpuUsage: 65,
      memoryUsage: 75,
      totalPods: 64,
      runningPods: 60
    },
    {
      name: 'node-1',
      role: 'worker',
      status: 'Ready',
      cpuUsage: 45,
      memoryUsage: 60,
      totalPods: 128,
      runningPods: 110
    },
    {
      name: 'node-2',
      role: 'worker',
      status: 'Ready',
      cpuUsage: 30,
      memoryUsage: 55,
      totalPods: 128,
      runningPods: 98
    },
    {
      name: 'node-3',
      role: 'worker',
      status: 'NotReady',
      cpuUsage 85,
      memoryUsage: 90,
      totalPods: 128,
      runningPods: 128
    },
    {
      name: 'node-4',
      role: 'worker',
      status: 'Ready',
      cpuUsage: 60,
      memoryUsage: 65,
      totalPods: 128,
      runningPods: 105
    }
  ]
  
  const recentEvents = [
    {
      timestamp: '2023-10-15T09:23:17Z',
      type: 'Warning',
      object: 'pod/nginx-deployment75675f5897-58xj7',
      namespace: 'default',
      reason: 'FailedScheduling',
      message: '0/5 nodes available: 3 Insufficient cpu, 2 Insufficient memory.'
    },
    {
      timestamp: '2023-10-15T09:20:45Z',
      type: 'Normal',
      object: 'deployment/nginx-deployment',
      namespace 'default',
      reason: 'ScalingReplicaSet',
      message: 'Scaled up replica set nginx-deployment-75675f5897 to 1'
    },
    {
      timestamp: '2023-10-15T09:18:32Z',
      type: 'Warning',
      object: 'node/node-3',
      namespace: '',
      reason: 'NodeNotReady',
      message: 'Node node-3 status is now: NodeNotReady'
    },
    {
      timestamp: '2023-10-15T09:15:12Z',
      type: 'Normal',
      object: 'service/redis-master',
      namespace: 'prod',
     : 'UpdatedLoadBalancer',
      message: 'Updated load balancer with new hosts'
    },
    {
      timestamp: '2023-10-15T09:12:58Z',
      type: 'Normal',
      object: 'pod/mysql-56f4d7f65d-f5z2n',
      namespace: 'prod',
      reason: 'Pulled',
      message: 'Successfully pulled image "mysql:5.7"'
    },
    {
      timestamp: '2023-10-15T09:10:31Z',
      type: 'Warning',
      object: 'pod/cron-job-1625868000-abcde',
      namespace: 'dev',
      reason: 'Failed',
      message: 'Job has reached the specified backoff limit'
    }
  ]
  
  const eventStatistics = [
    { type: 'Normal', count: 42, percentage: 70, trend: 'up' },
    { type: 'Warning', count: 18, percentage: 30, trend: 'down' }
  ]
  
  const eventStatisticsOption = {
    tooltip: {
      trigger: 'item'
       legend: {
      top: '5%',
      left: 'center'
    },
    series: [
      {
        name: '事件统计',
        type: 'pie',
        radius: ['50%', '70%'],
        center: ['50%', '60%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: '18',
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: [
          { value: 42, name: '正常事件 (70%)', itemStyle: { color: '#67C23A' } },
          { value: 18, name: '警告事件 (30%)', itemStyle: { color: '#E6A23C' } }
        ]
      }
    ]
  }
  
  return {
    namespaces,
    overviewData,
    cpuUsageOption,
    memoryUsageOption,
    storageUsageOption,
    networkUsageOption,
    nodes,
    recentEvents,
    eventStatistics,
    eventStatisticsOption
  }
}
</script>