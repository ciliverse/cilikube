// ./types/kubernetes.ts (Create this file if it doesn't exist)
export interface HealthStatusItem {
    label: string;
    value: string;
    color: string;
  }
  
  export interface AlertSummaryItem {
    label: string;
    count: number;
    color: string;
    trend: number;
  }
  
  export interface OverviewDataItem {
    title: string;
    value: string;
    percent: number;
    total?: string | number; // Make total optional or always string/number
    icon: string; // Store icon name as string
    color: string;
  }
  
  // Basic chart data structure
  export interface ChartData {
    labels: string[];
    values: number[] | number[][]; // Can be single series or multi-series
  }
  
  export interface NodeData {
    name: string;
    role: 'master' | 'worker';
    status: 'Ready' | 'NotReady' | string; // Allow other statuses
    cpuUsage: number;
    memoryUsage: number;
    totalPods: number;
    runningPods: number;
    isNew?: boolean; // Optional flag
  }
  
  export interface EventItem {
    timestamp: string; // ISO date string
    type: 'Normal' | 'Warning' | string;
    object: string;
    namespace: string;
    reason: string;
    message: string;
  }
  
  export interface EventStatisticItem {
    type: string;
    count: number;
    percentage: number;
    trend: 'up' | 'down';
  }
  
  export interface ResourceAlertItem {
    name: string;
    message: string;
    level: 'critical' | 'warning' | 'notice' | string;
    levelText: string;
    time: string; // Relative time like "5分钟前"
  }
  
  // Example for a consolidated API response (optional but recommended)
  export interface DashboardDataResponse {
    namespaces: string[];
    healthStatus: HealthStatusItem[];
    alertSummary: AlertSummaryItem[];
    overviewData: OverviewDataItem[];
    cpuUsageData: ChartData;
    memoryUsageData: ChartData; // Assuming structure for stacked bar
    storageUsageData: ChartData; // Assuming structure for area chart
    networkUsageData: ChartData;
    nodes: NodeData[];
    recentEvents: EventItem[];
    eventStatistics: EventStatisticItem[];
    resourceAlerts: ResourceAlertItem[];
  }