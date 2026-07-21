export default {
  // 通用
  common: {
    confirm: 'Confirm',
    cancel: 'Cancel',
    save: 'Save',
    delete: 'Delete',
    edit: 'Edit',
    add: 'Add',
    search: 'Search',
    reset: 'Reset',
    refresh: 'Refresh',
    loading: 'Loading...',
    success: 'Success',
    error: 'Error',
    warning: 'Warning',
    info: 'Info',
    retry: 'Retry',
    yes: 'Yes',
    no: 'No',
    view: 'View',
    name: 'Name',
    namespace: 'Namespace',
    createdAt: 'Created At',
    actions: 'Actions',
    all: 'All'
  },

  // 设置
  settings: {
    title: 'Settings',
    language: 'Language',
    font: 'Font',
    theme: 'Theme',
    reset: 'Reset',
    
    // Layout Settings
    layout: {
      title: 'Layout Configuration',
      showTagsView: 'Show Tags View',
      showLogo: 'Show Logo',
      fixedHeader: 'Fixed Header',
      showFooter: 'Show Footer',
      showNotify: 'Show Notifications',
      showThemeSwitch: 'Show Theme Switch',
      showScreenfull: 'Show Fullscreen',
      showSearchMenu: 'Show Search Menu',
      showFontSelector: 'Show Font Selector',
      cacheTagsView: 'Cache Tags View',
      showWatermark: 'Show Watermark',
      showGreyMode: 'Show Grey Mode',
      showColorWeakness: 'Show Color Weakness',
      modes: {
        left: 'Left Mode',
        top: 'Top Mode',
        mixed: 'Mixed Mode'
      }
    },
    
    // Feature Settings
    features: {
      title: 'Feature Configuration'
    }
  },

  // 语言选项
  language: {
    english: 'English',
    chinese: 'Chinese'
  },

  // 字体选项
  font: {
    'maple-mono': 'Maple Mono',
    'victor-mono': 'Victor Mono'
  },

  // 主题选项
  theme: {
    default: 'Default',
    dark: 'Dark',
    darkBlue: 'Dark Blue'
  },

  // 导航
  nav: {
    home: 'Home',
    dashboard: 'Dashboard',
    settings: 'Settings'
  },

  // 菜单
  menu: {
    // 看板
    board: 'Dashboard',
    clusterInstall: 'Cluster Install',
    clusterOverview: 'Cluster Overview',
    clusterMonitor: 'Cluster Monitor',
    opsNavigation: 'Ops Navigation',

    // 集群
    cluster: 'Cluster',
    clusterManagement: 'Cluster Management',
    clusterEvents: 'Cluster Events',
    node: 'Node',
    namespace: 'Namespace',

    // 工作负载
    workloads: 'Workloads',
    pod: 'Pod',
    deployment: 'Deployment',

    // 存储
    storage: 'Storage',
    pv: 'PV',
    pvc: 'PVC',

    // 网络
    network: 'Network',
    service: 'Service',
    ingress: 'Ingress',

    // 配置管理
    config: 'Configuration',
    configmap: 'ConfigMap',
    secret: 'Secret',

    // 项目管理
    project: 'Project Management',
    chineseDocs: 'Chinese Docs',
    myBlog: 'My Blog',
    techStack: 'Tech Stack',

    // 权限
    permission: 'Permission Demo',
    pageLevel: 'Page Level Permission',
    buttonLevel: 'Button Level Permission',

    // 系统管理
    admin: 'System Management',
    userManagement: 'User Management',
    roleManagement: 'Role Management',
    systemSettings: 'System Settings',

    // CRD
    crd: 'Custom Resources'
  },

  // 操作
  actions: {
    profile: 'Profile',
    logout: 'Logout',
    refresh: 'Refresh',
    refreshing: 'Refreshing cluster list...',
    switchedToCluster: 'Switched to cluster: {name}',
    unknownCluster: 'Unknown cluster',
    noClusterSelected: 'No cluster selected',
    noClustersAvailable: 'No clusters available',
    loading: 'Loading...'
  },

  // 标签页操作
  tagsView: {
    refresh: 'Refresh',
    close: 'Close',
    closeOthers: 'Close Others',
    closeAll: 'Close All'
  },

  // 应用消息
  app: {
    welcomeTitle: 'Hello! Friend!',
    welcomeMessage: 'This is a passion project, please give it a free star! If possible, follow the WeChat account - Cillian, let\'s learn and progress together!',
    migrationSuccess: 'Your cluster selection settings have been automatically migrated to the new version'
  },

  // UI 组件
  ui: {
    themeSwitch: 'Theme Switch',
    notifications: 'Notifications',
    messages: 'Messages',
    todos: 'Todos',
    viewHistory: 'View {type} History',
    searchMenu: 'Search Menu',
    fullscreen: {
      enter: 'Fullscreen',
      exit: 'Exit Fullscreen',
      browserNotSupported: 'Your browser does not support fullscreen',
      contentEnlarge: 'Enlarge Content',
      contentRestore: 'Restore Content',
      contentFullscreen: 'Content Fullscreen'
    }
  },

  // 通知消息
  notifications: {
    newYear: {
      title: 'New Year 2025',
      description: 'Happy New Year! Let\'s work hard in the new year!'
    },
    cilikubeLaunch: {
      title: 'CiliKube is Online',
      description: 'A free and open-source k8s management system solution with front-end and back-end separation, using the latest technology stack'
    },
    newVersion: {
      title: 'New Version Released',
      description: 'The long-awaited v1.0 version is finally released, thank you for your support!'
    }
  },

  // 消息
  messages: {
    morning: {
      title: 'Good Morning Worker',
      description: 'If you don\'t work hard today, your position won\'t be stable tomorrow. Good morning, worker!'
    },
    noon: {
      title: 'Good Noon Worker',
      description: 'If you don\'t work hard today, your position won\'t be stable tomorrow. Good noon, worker!'
    },
    evening: {
      title: 'Good Evening Worker',
      description: 'If you don\'t work hard today, your position won\'t be stable tomorrow. Good evening, worker!'
    }
  },

  // 待办事项
  todos: {
    task1: {
      title: 'Task One',
      description: 'Follow Cillian WeChat Official Account'
    },
    task2: {
      title: 'Task Two',
      description: 'Add project team WeChat ciliverse'
    },
    task3: {
      title: 'Task Three',
      description: 'Join Cillian tech discussion group'
    },
    status: {
      notStarted: 'Not Started',
      inProgress: 'In Progress',
      overdue: 'Overdue'
    }
  }, 
 // 集群概览
  clusterOverview: {
    description: 'View overall cluster status and resource usage',
    clusterStatus: 'Cluster Status',
    healthy: 'Healthy',
    running: 'Running',
    totalNodes: 'Total Nodes',
    runningPods: 'Running Pods',
    namespaces: 'Namespaces',
    resourceOverview: 'Resource Overview',
    resourceDescription: 'Cluster resource usage and performance metrics',
    totalResources: 'Total Resources',
    activeTypes: 'Active Types',
    noResourceData: 'No resource data available',
    nodeMetrics: 'Node Metrics'
  },

  // 节点监控
  nodeMetrics: {
    avgCpuUsage: 'Avg CPU Usage',
    avgMemoryUsage: 'Avg Memory Usage',
    healthyNodes: 'Healthy Nodes',
    totalCpuCores: 'Total CPU Cores',
    totalMemory: 'Total Memory',
    diskUsage: 'Disk Usage',
    networkIO: 'Network I/O',
    cores: 'Cores',
    cpuUsage: 'CPU Usage',
    memoryUsage: 'Memory Usage',
    noData: 'No monitoring data available',
    resourceRequests: 'Resource Requests',
    cpuRequests: 'CPU Requests',
    memoryRequests: 'Memory Requests'
  },

  // 集群事件
  clusterEvents: {
    title: 'Cluster Events',
    realtime: 'Real-time',
    refresh: 'Refresh',
    viewAll: 'View All Events',
    noEvents: 'No events',
    fetchError: 'Failed to fetch events',
    eventList: 'Event List',
    totalEvents: 'Total {count} events',
    count: '{count} times',
    firstTime: 'First Time',
    filters: {
      namespace: 'Namespace',
      type: 'Type',
      limit: 'Limit',
      allNamespaces: 'All Namespaces',
      allTypes: 'All Types'
    },
    eventTypes: {
      normal: 'Normal',
      warning: 'Warning',
      error: 'Error',
      info: 'Info'
    },
    timeAgo: {
      minutesAgo: 'minutes ago',
      hoursAgo: 'hours ago',
      daysAgo: 'days ago',
      justNow: 'just now'
    }
  },

  // 集群管理
  clusterManagement: {
    title: 'Cluster Management',
    description: 'Manage your Kubernetes clusters, support multi-cluster switching and monitoring',
    addCluster: 'Add Cluster',
    searchPlaceholder: 'Search cluster name, environment or server address...',
    filterEnvironment: 'Filter Environment',
    filterStatus: 'Filter Status',
    allEnvironments: 'All Environments',
    allStatuses: 'All Statuses',
    cardView: 'Card View',
    tableView: 'Table View',

    // 统计
    stats: {
      totalClusters: 'Total Clusters',
      activeClusters: 'Active Clusters',
      healthyClusters: 'Healthy Clusters',
      unhealthyClusters: 'Unhealthy Clusters'
    },

    // 环境
    environments: {
      production: 'Production',
      staging: 'Staging',
      development: 'Development',
      testing: 'Testing'
    },

    // 状态
    statuses: {
      available: 'Available',
      unavailable: 'Unavailable'
    },

    // 表格列
    columns: {
      name: 'Cluster Name',
      clusterName: 'Cluster Name',
      status: 'Status',
      server: 'Server Address',
      apiServer: 'API Server',
      version: 'K8s Version',
      environment: 'Environment',
      source: 'Source',
      actions: 'Actions'
    },

    // 操作
    actions: {
      setActive: 'Set Active',
      edit: 'Edit',
      delete: 'Delete',
      active: 'Active'
    },

    // 来源
    sources: {
      database: 'Database',
      file: 'File'
    },

    // 表单
    form: {
      name: 'Cluster Name',
      namePlaceholder: 'Enter unique cluster name',
      environment: 'Environment',
      environmentPlaceholder: 'Select environment type',
      provider: 'Provider',
      providerPlaceholder: 'e.g: minikube, aws, gcp, aliyun',
      region: 'Region',
      regionPlaceholder: 'e.g: us-west-2, cn-beijing-1',
      description: 'Description',
      descriptionPlaceholder: 'Enter cluster description',
      kubeconfig: 'Kubeconfig',
      kubeconfigTitle: 'Kubeconfig Configuration',
      kubeconfigPlaceholder: 'Paste your Kubeconfig file content here',
      kubeconfigTips: 'Please ensure your Kubeconfig file is valid and has proper access permissions',
      clusterNameRequired: 'Cluster name is required',
      kubeconfigRequired: 'Kubeconfig content is required',
      environmentRequired: 'Please select environment type'
    },

    // 基本信息
    basicInfo: 'Basic Information',
    clustersFound: 'clusters found',

    // 消息
    messages: {
      addSuccess: 'Cluster added successfully',
      addFailed: 'Failed to add cluster',
      deleteConfirm: 'Are you sure you want to delete this cluster?',
      deleteSuccess: 'Cluster deleted successfully',
      deleteFailed: 'Failed to delete cluster',
      setActiveSuccess: 'Cluster "{name}" has been set as active cluster',
      setActiveFailed: 'Failed to switch active cluster',
      fetchFailed: 'Failed to fetch cluster data, please check network or contact administrator',
      noMatchingClusters: 'No matching clusters found',
      canOnlyDeleteDatabaseClusters: 'Only clusters added through database can be deleted'
    }
  },

  // CRD
  crd: {
    description: 'Manage Custom Resource Definitions and custom resources',
    searchPlaceholder: 'Search CRDs...',
    selectGroup: 'Select Group',
    selectNamespace: 'Select Namespace',
    noCRDs: 'No Custom Resource Definitions found',
    noResources: 'No custom resources found',

    // CRD fields
    name: 'Name',
    group: 'Group',
    version: 'Version',
    kind: 'Kind',
    plural: 'Plural',
    singular: 'Singular',
    scope: 'Scope',
    categories: 'Categories',
    shortNames: 'Short Names',
    apiVersion: 'API Version',

    // CRD details
    basicInfo: 'Basic Information',
    versions: 'Versions',
    conditions: 'Conditions',
    metadata: 'Metadata',
    labels: 'Labels',
    annotations: 'Annotations',
    spec: 'Spec',
    status: 'Status',

    // Version fields
    served: 'Served',
    storage: 'Storage',

    // Condition fields
    conditionType: 'Type',
    conditionStatus: 'Status',
    conditionReason: 'Reason',
    conditionMessage: 'Message',
    lastTransitionTime: 'Last Transition Time',

    // Actions
    viewResources: 'View Resources',
    createResource: 'Create Resource',
    editResource: 'Edit Resource',
    deleteResource: 'Delete Resource',
    resources: 'Resources',

    // Forms
    enterName: 'Enter resource name',
    nameRequired: 'Name is required',
    nameFormat: 'Name must be lowercase alphanumeric with hyphens',
    namespaceRequired: 'Namespace is required',
    labelKey: 'Label Key',
    labelValue: 'Label Value',
    addLabel: 'Add Label',
    specPlaceholder: 'Enter resource specification in JSON format',
    invalidJSON: 'Invalid JSON format',

    // Messages
    deleteConfirm: 'Are you sure you want to delete resource "{name}"?',
    
    // Error messages
    fetchCRDsFailed: 'Failed to fetch CRD list',
    fetchCRDDetailFailed: 'Failed to fetch CRD details',
    fetchResourcesFailed: 'Failed to fetch custom resources list',
    fetchResourceDetailFailed: 'Failed to fetch custom resource details',
    createResourceFailed: 'Failed to create custom resource',
    updateResourceFailed: 'Failed to update custom resource',
    deleteResourceFailed: 'Failed to delete custom resource',
    
    // Success messages
    createResourceSuccess: 'Custom resource created successfully',
    updateResourceSuccess: 'Custom resource updated successfully',
    deleteResourceSuccess: 'Custom resource deleted successfully'
  },

  // 节点管理
  nodeManagement: {
    title: 'Nodes',
    searchPlaceholder: 'Search nodes...',
    filterRole: 'Filter by Role',
    filterStatus: 'Filter by Status',
    allRoles: 'All Roles',
    allStatuses: 'All Statuses',
    viewToggle: {
      card: 'Card View',
      list: 'List View'
    },
    columns: {
      name: 'Name',
      status: 'Status',
      role: 'Role',
      version: 'Version',
      internalIP: 'Internal IP',
      externalIP: 'External IP',
      os: 'OS',
      kernel: 'Kernel',
      runtime: 'Runtime',
      cpu: 'CPU',
      memory: 'Memory',
      pods: 'Pods',
      created: 'Created'
    },
    statuses: {
      ready: 'Ready',
      notReady: 'Not Ready',
      unknown: 'Unknown'
    },
    roles: {
      controlPlane: 'Control Plane',
      worker: 'Worker'
    },
    info: {
      addresses: 'Addresses',
      capacity: 'Capacity',
      allocatable: 'Allocatable',
      systemInfo: 'System Info',
      conditions: 'Conditions',
      taints: 'Taints',
      kubeletVersion: 'Kubelet Version',
      kubeProxyVersion: 'Kube-proxy Version',
      containerRuntime: 'Container Runtime',
      kernelVersion: 'Kernel Version',
      osImage: 'OS Image',
      podCIDR: 'Pod CIDR',
      allAddresses: 'All Addresses'
    },
    actions: {
      cordon: 'Cordon',
      uncordon: 'Uncordon',
      drain: 'Drain',
      viewDetails: 'View Details'
    },
    messages: {
      noNodes: 'No nodes found in current cluster',
      noMatchingNodes: 'No nodes match your search criteria',
      selectCluster: 'Please select a cluster from the navigation bar to view nodes.',
      loadingNodes: 'Loading nodes for cluster'
    }
  },

  // 命名空间管理
  namespaceManagement: {
    title: 'Namespaces',
    searchPlaceholder: 'Search namespaces...',
    filterStatus: 'Filter by Status',
    allStatuses: 'All Statuses',
    viewToggle: {
      card: 'Card View',
      list: 'List View'
    },
    columns: {
      name: 'Name',
      status: 'Status',
      created: 'Created',
      labels: 'Labels',
      annotations: 'Annotations'
    },
    statuses: {
      active: 'Active',
      terminating: 'Terminating',
      unknown: 'Unknown'
    },
    info: {
      systemNamespace: 'System Namespace',
      userNamespace: 'User Namespace',
      resourceCount: 'Resource Count',
      description: 'Description'
    },
    actions: {
      create: 'Create Namespace',
      viewDetails: 'View Details',
      editMetadata: 'Edit Metadata',
      delete: 'Delete'
    },
    form: {
      name: 'Name',
      namePlaceholder: 'Enter namespace name',
      nameRequired: 'Namespace name is required',
      nameFormat: 'Name must contain only lowercase letters, numbers, or hyphens, and start/end with alphanumeric characters',
      nameLength: 'Name cannot exceed 63 characters',
      labels: 'Labels',
      labelsPlaceholder: 'Enter labels in format: key=value',
      annotations: 'Annotations',
      annotationsPlaceholder: 'Enter annotations in format: key=value',
      description: 'Description',
      descriptionPlaceholder: 'Enter namespace description'
    },
    messages: {
      noNamespaces: 'No namespaces found',
      noMatchingNamespaces: 'No namespaces match your search criteria',
      selectCluster: 'Please select a cluster to view namespaces.',
      createSuccess: 'Namespace created successfully',
      createFailed: 'Failed to create namespace',
      deleteConfirm: 'Are you sure you want to delete namespace "{name}"? This will delete all resources in this namespace and cannot be undone!',
      deleteSuccess: 'Namespace deleted successfully',
      deleteFailed: 'Failed to delete namespace',
      cannotDeleteSystem: 'Cannot delete system namespace',
      loadingNamespaces: 'Loading namespaces'
    }
  },

  // Deployment 管理
  deploymentManagement: {
    title: 'Deployments',
    searchPlaceholder: 'Search Deployment name...',
    filterNamespace: 'Filter Namespace',
    filterStatus: 'Filter by Status',
    allNamespaces: 'All Namespaces',
    allStatuses: 'All Statuses',
    viewToggle: {
      card: 'Card View',
      list: 'List View'
    },
    columns: {
      name: 'Name',
      namespace: 'Namespace',
      status: 'Status',
      replicas: 'Replicas',
      ready: 'Ready Replicas',
      upToDate: 'Up-to-Date',
      available: 'Available',
      age: 'Age',
      images: 'Images'
    },
    statuses: {
      available: 'Available',
      progressing: 'Progressing',
      replicaFailure: 'Replica Failure',
      degraded: 'Degraded',
      unknown: 'Unknown'
    },
    info: {
      strategy: 'Update Strategy',
      selector: 'Selector',
      template: 'Pod Template',
      conditions: 'Conditions'
    },
    actions: {
      create: 'Create Deployment',
      scale: 'Scale',
      rollout: 'Rollout',
      restart: 'Restart',
      edit: 'Edit YAML',
      delete: 'Delete',
      viewDetails: 'View Details',
      viewPods: 'View Pods'
    },
    form: {
      createFromYaml: 'Create from YAML',
      editYaml: 'Edit YAML',
      yamlContent: 'YAML Content',
      name: 'Name',
      namespace: 'Namespace',
      replicas: 'Replicas',
      image: 'Image',
      nameRequired: 'Name is required',
      namespaceRequired: 'Namespace is required',
      replicasRequired: 'Replicas count is required',
      imageRequired: 'Image is required',
      scaleReplicas: 'Scale Replicas',
      currentReplicas: 'Current Replicas',
      targetReplicas: 'Target Replicas'
    },
    messages: {
      noDeployments: 'No Deployments found in current namespace',
      noMatchingDeployments: 'No Deployments match your search criteria',
      selectNamespace: 'Please select a namespace first to view Deployments.',
      createSuccess: 'Deployment created successfully',
      createFailed: 'Failed to create Deployment',
      updateSuccess: 'Deployment updated successfully',
      updateFailed: 'Failed to update Deployment',
      scaleSuccess: 'Deployment scaled successfully',
      scaleFailed: 'Failed to scale Deployment',
      deleteConfirm: 'Are you sure you want to delete Deployment "{name}"?',
      deleteSuccess: 'Deployment deleted successfully',
      deleteFailed: 'Failed to delete Deployment',
      loadingDeployments: 'Loading Deployments'
    }
  },

  // Pod 管理
  podManagement: {
    title: 'Pods',
    searchPlaceholder: 'Search Pod name, IP, node or status...',
    filterNamespace: 'Filter Namespace',
    filterStatus: 'Filter by Status',
    allNamespaces: 'All Namespaces',
    allStatuses: 'All Statuses',
    viewToggle: {
      card: 'Card View',
      list: 'List View'
    },
    columns: {
      name: 'Name',
      namespace: 'Namespace',
      status: 'Status',
      ip: 'Pod IP',
      node: 'Node',
      created: 'Created',
      restarts: 'Restarts',
      ready: 'Ready'
    },
    statuses: {
      running: 'Running',
      pending: 'Pending',
      succeeded: 'Succeeded',
      failed: 'Failed',
      unknown: 'Unknown',
      terminating: 'Terminating',
      crashLoopBackOff: 'CrashLoopBackOff',
      imagePullBackOff: 'ImagePullBackOff',
      containerCreating: 'ContainerCreating'
    },
    info: {
      containers: 'Containers',
      restartPolicy: 'Restart Policy',
      qosClass: 'QoS Class',
      phase: 'Phase',
      conditions: 'Conditions'
    },
    actions: {
      create: 'Create Pod',
      viewLogs: 'View Logs',
      exec: 'Exec',
      edit: 'Edit YAML',
      delete: 'Delete',
      viewDetails: 'View Details'
    },
    form: {
      createFromYaml: 'Create from YAML',
      editYaml: 'Edit YAML',
      yamlContent: 'YAML Content',
      container: 'Container',
      command: 'Command',
      tailLines: 'Tail Lines',
      follow: 'Follow Logs'
    },
    messages: {
      noPods: 'No Pods found in current namespace',
      noMatchingPods: 'No Pods match your search criteria',
      selectNamespace: 'Please select a namespace to view Pods.',
      createSuccess: 'Pod created successfully',
      createFailed: 'Failed to create Pod',
      updateSuccess: 'Pod updated successfully',
      updateFailed: 'Failed to update Pod',
      deleteConfirm: 'Are you sure you want to delete Pod "{name}"?',
      deleteSuccess: 'Pod deleted successfully',
      deleteFailed: 'Failed to delete Pod',
      loadingPods: 'Loading Pods',
      loadingContainers: 'Loading container list',
      connectingTerminal: 'Connecting to terminal',
      terminalConnected: 'Terminal connected',
      terminalDisconnected: 'Terminal disconnected'
    }
  }, 
 // PersistentVolume 管理
  persistentVolumeManagement: {
    title: 'PersistentVolumes',
    searchPlaceholder: 'Search PV name / StorageClass',
    columns: {
      name: 'Name',
      status: 'Status',
      capacity: 'Capacity',
      accessModes: 'Access Modes',
      reclaimPolicy: 'Reclaim Policy',
      volumeMode: 'Volume Mode',
      storageClass: 'StorageClass',
      boundClaim: 'Bound to (PVC)',
      createdAt: 'Created At',
      actions: 'Actions'
    },
    actions: {
      create: 'Create PV',
      edit: 'Edit',
      delete: 'Delete',
      editYaml: 'Edit YAML',
      applyYaml: 'Apply YAML'
    },
    messages: {
      noPVs: 'No PersistentVolumes found',
      deleteConfirm: 'Are you sure you want to delete PersistentVolume "{name}"?',
      deleteSuccess: 'PersistentVolume deleted successfully',
      deleteFailed: 'Failed to delete PersistentVolume',
      yamlEditTip: 'Paste or edit the YAML configuration of the PersistentVolume here. PV is a cluster-scoped resource and does not belong to any namespace.',
      yamlEditorPlaceholder: 'YAML editor placeholder (e.g., using Monaco Editor or Codemirror)',
      createPVTitle: 'Create PV (YAML)',
      fetchingYaml: 'Simulating: Fetching YAML for PV "{name}"'
    }
  },

  // PersistentVolumeClaim 管理
  persistentVolumeClaimManagement: {
    title: 'PersistentVolumeClaims',
    filterNamespace: 'Select Namespace',
    searchPlaceholder: 'Search PVC name / Bound PV',
    columns: {
      name: 'Name',
      namespace: 'Namespace',
      status: 'Status',
      boundVolume: 'Bound Volume (PV)',
      requestCapacity: 'Request Capacity',
      actualCapacity: 'Actual Capacity',
      accessModes: 'Access Modes',
      volumeMode: 'Volume Mode',
      storageClass: 'StorageClass',
      createdAt: 'Created At',
      actions: 'Actions'
    },
    actions: {
      create: 'Create PVC',
      edit: 'Edit',
      delete: 'Delete',
      editYaml: 'Edit YAML',
      applyYaml: 'Apply YAML'
    },
    messages: {
      noPVCs: 'No PersistentVolumeClaims found',
      noNamespaces: 'No available namespaces',
      selectNamespace: 'Please select a namespace to load PVCs',
      selectNamespaceWarning: 'Please select a namespace first',
      deleteConfirm: 'Are you sure you want to delete PersistentVolumeClaim "{name}"?',
      deleteConfirmDetail: 'Are you sure you want to delete PVC "{name}" (namespace: {namespace})? Please note the PV\'s reclaim policy.',
      deleteSuccess: 'PersistentVolumeClaim deleted successfully',
      deleteFailed: 'Failed to delete PersistentVolumeClaim',
      yamlEditTip: 'Paste or edit the YAML configuration of the PersistentVolumeClaim here. Ensure the namespace in the YAML matches the currently selected namespace or is omitted.',
      yamlEditorPlaceholder: 'YAML editor placeholder (e.g., using Monaco Editor or Codemirror)',
      createPVCTitle: 'Create PVC (YAML)',
      fetchingYaml: 'Fetching YAML for PVC "{name}"',
      editError: 'Cannot edit, missing raw data',
      fetchNamespacesError: 'Failed to fetch namespaces',
      fetchPVCsError: 'Failed to fetch PVC data',
      operationSuccess: 'Successfully {action} PVC',
      operationCancelled: 'Delete operation cancelled',
      noNamespace: 'No available namespace',
      deleteCanceled: 'Delete operation cancelled',
      confirmDelete: 'Confirm Delete',
      editSecretTitle: 'Edit Secret',
      fetchSecretDetailsFailed: 'Failed to fetch Secret details',
      fetchNamespacesFailed: 'Failed to fetch namespaces',
      create: 'Create',
      update: 'Update',
      cannotDetermineNamespace: 'Cannot determine target namespace',
      dataFormatError: 'Data format error',
      networkError: 'Network request failed',
      invalidResponse: 'Invalid response data',
      unknownError: 'Unknown error',
      requestFailed: 'Request failed',
      dataKeys: 'Data Entries (Keys)',
      noSecretsInNamespace: 'No Secrets found in namespace',
      aboutPVCTitle: 'About PersistentVolumeClaims',
      aboutPVCDescription: 'A PersistentVolumeClaim (PVC) is a request for storage by a user. It is similar to a Pod. Pods consume node resources and PVCs consume PV resources. A Pod can request a specific amount of resources (CPU and memory); similarly, a PVC can request a specific size and access mode (for example, it can be mounted as ReadWriteOnce, ReadOnlyMany, ReadWriteMany or ReadWriteOncePod).'
    }
  },

  // Service 管理
  serviceManagement: {
    title: 'Services',
    filterNamespace: 'Select Namespace',
    searchPlaceholder: 'Search Service name / type / ClusterIP',
    columns: {
      name: 'Name',
      namespace: 'Namespace',
      type: 'Type',
      clusterIP: 'ClusterIP',
      externalIP: 'External IP',
      ports: 'Ports',
      selector: 'Selector',
      createdAt: 'Created At',
      actions: 'Actions'
    },
    actions: {
      create: 'Create Service',
      edit: 'Edit',
      delete: 'Delete',
      editYaml: 'Edit YAML',
      applyYaml: 'Apply YAML'
    },
    messages: {
      noServices: 'No Services found',
      noNamespaces: 'No available namespaces',
      selectNamespace: 'Please select a namespace to load Services',
      selectNamespaceWarning: 'Please select a namespace first',
      deleteConfirm: 'Are you sure you want to delete Service "{name}"?',
      deleteConfirmDetail: 'Are you sure you want to delete Service "{name}" (namespace: {namespace})?',
      deleteSuccess: 'Service deleted successfully',
      deleteFailed: 'Failed to delete Service',
      yamlEditTip: 'Paste or edit the YAML configuration of the Service here. Ensure the namespace in the YAML matches the currently selected namespace or is omitted.',
      yamlEditorPlaceholder: 'YAML editor placeholder (e.g., using Monaco Editor or Codemirror)',
      createServiceTitle: 'Create Service (YAML)',
      fetchingYaml: 'Fetching YAML for Service "{name}"',
      fetchServicesError: 'Failed to fetch Service data',
      operationSuccess: 'Successfully {action} Service',
      operationCancelled: 'Delete operation cancelled',
      noNamespace: 'No available namespace',
      deleteCanceled: 'Delete operation cancelled',
      confirmDelete: 'Confirm Delete',
      editSecretTitle: 'Edit Secret',
      fetchSecretDetailsFailed: 'Failed to fetch Secret details',
      fetchNamespacesFailed: 'Failed to fetch namespaces',
      create: 'Create',
      update: 'Update',
      cannotDetermineNamespace: 'Cannot determine target namespace',
      dataFormatError: 'Data format error',
      networkError: 'Network request failed',
      invalidResponse: 'Invalid response data',
      unknownError: 'Unknown error',
      requestFailed: 'Request failed',
      dataKeys: 'Data Entries (Keys)',
      noSecretsInNamespace: 'No Secrets found in namespace'
    }
  },

  // Ingress 管理
  ingressManagement: {
    title: 'Ingresses',
    filterNamespace: 'Select Namespace',
    searchPlaceholder: 'Search Ingress name / host / class',
    columns: {
      name: 'Name',
      namespace: 'Namespace',
      class: 'Class',
      hosts: 'Hosts',
      address: 'Address (LoadBalancer IP)',
      ports: 'Ports',
      rules: 'Rules',
      createdAt: 'Created At',
      actions: 'Actions'
    },
    actions: {
      create: 'Create Ingress',
      edit: 'Edit',
      delete: 'Delete',
      editYaml: 'Edit YAML',
      applyYaml: 'Apply YAML'
    },
    messages: {
      noIngresses: 'No Ingresses found',
      noNamespaces: 'No available namespaces',
      selectNamespace: 'Please select a namespace to load Ingresses',
      selectNamespaceWarning: 'Please select a namespace first',
      deleteConfirm: 'Are you sure you want to delete Ingress "{name}"?',
      deleteConfirmDetail: 'Are you sure you want to delete Ingress "{name}" (namespace: {namespace})?',
      deleteSuccess: 'Ingress deleted successfully',
      deleteFailed: 'Failed to delete Ingress',
      yamlEditTip: 'Paste or edit the YAML configuration of the Ingress here. Ensure the namespace in the YAML matches the currently selected namespace or is omitted.',
      yamlEditorPlaceholder: 'YAML editor placeholder (e.g., using Monaco Editor or Codemirror)',
      createIngressTitle: 'Create Ingress (YAML)',
      fetchingYaml: 'Fetching YAML for Ingress "{name}"',
      fetchIngressesError: 'Failed to fetch Ingress data',
      operationSuccess: 'Successfully {action} Ingress',
      operationCancelled: 'Delete operation cancelled',
      noNamespace: 'No available namespace',
      deleteCanceled: 'Delete operation cancelled',
      confirmDelete: 'Confirm Delete',
      editSecretTitle: 'Edit Secret',
      fetchSecretDetailsFailed: 'Failed to fetch Secret details',
      fetchNamespacesFailed: 'Failed to fetch namespaces',
      create: 'Create',
      update: 'Update',
      cannotDetermineNamespace: 'Cannot determine target namespace',
      dataFormatError: 'Data format error',
      networkError: 'Network request failed',
      invalidResponse: 'Invalid response data',
      unknownError: 'Unknown error',
      requestFailed: 'Request failed',
      dataKeys: 'Data Entries (Keys)',
      noSecretsInNamespace: 'No Secrets found in namespace'
    }
  },

  // ConfigMap 管理
  configmapManagement: {
    title: 'ConfigMaps',
    filterNamespace: 'Select Namespace',
    searchPlaceholder: 'Search ConfigMap name',
    columns: {
      name: 'Name',
      namespace: 'Namespace',
      dataCount: 'Data Entries (Keys)',
      createdAt: 'Created At',
      actions: 'Actions'
    },
    actions: {
      create: 'Create ConfigMap',
      edit: 'Edit',
      delete: 'Delete',
      editYaml: 'Edit YAML',
      applyYaml: 'Apply YAML'
    },
    messages: {
      noConfigMaps: 'No ConfigMaps found',
      noNamespaces: 'No available namespaces',
      selectNamespace: 'Please select a namespace to load ConfigMaps',
      selectNamespaceWarning: 'Please select a namespace first',
      deleteConfirm: 'Are you sure you want to delete ConfigMap "{name}"?',
      deleteConfirmDetail: 'Are you sure you want to delete ConfigMap "{name}" (namespace: {namespace})?',
      deleteSuccess: 'ConfigMap deleted successfully',
      deleteFailed: 'Failed to delete ConfigMap',
      yamlEditTip: 'Paste or edit the YAML configuration of the ConfigMap here. Ensure the namespace in the YAML matches the currently selected namespace or is omitted.',
      yamlEditorPlaceholder: 'YAML editor placeholder (e.g., using Monaco Editor or Codemirror)',
      createConfigMapTitle: 'Create ConfigMap (YAML)',
      fetchingYaml: 'Fetching YAML for ConfigMap "{name}"',
      fetchConfigMapsError: 'Failed to fetch ConfigMap data',
      operationSuccess: 'Successfully {action} ConfigMap',
      operationCancelled: 'Delete operation cancelled',
      noNamespace: 'No available namespace',
      deleteCanceled: 'Delete operation cancelled',
      confirmDelete: 'Confirm Delete',
      editSecretTitle: 'Edit Secret',
      fetchSecretDetailsFailed: 'Failed to fetch Secret details',
      fetchNamespacesFailed: 'Failed to fetch namespaces',
      create: 'Create',
      update: 'Update',
      cannotDetermineNamespace: 'Cannot determine target namespace',
      dataFormatError: 'Data format error',
      networkError: 'Network request failed',
      invalidResponse: 'Invalid response data',
      unknownError: 'Unknown error',
      requestFailed: 'Request failed',
      dataKeys: 'Data Entries (Keys)',
      noSecretsInNamespace: 'No Secrets found in namespace'
    }
  },

  // Secret 管理
  secretManagement: {
    title: 'Secrets',
    filterNamespace: 'Select Namespace',
    searchPlaceholder: 'Search Secret name / type',
    columns: {
      name: 'Name',
      namespace: 'Namespace',
      type: 'Type',
      dataCount: 'Data Entries (Keys)',
      createdAt: 'Created At',
      actions: 'Actions'
    },
    actions: {
      create: 'Create Secret',
      edit: 'Edit',
      delete: 'Delete',
      editYaml: 'Edit YAML',
      applyYaml: 'Apply YAML'
    },
    messages: {
      noSecrets: 'No Secrets found',
      noNamespaces: 'No available namespaces',
      selectNamespace: 'Please select a namespace to load Secrets',
      selectNamespaceWarning: 'Please select a namespace first',
      deleteConfirm: 'Are you sure you want to delete Secret "{name}"?',
      deleteConfirmDetail: 'Are you sure you want to delete Secret "{name}" (namespace: {namespace})?',
      deleteSuccess: 'Secret deleted successfully',
      deleteFailed: 'Failed to delete Secret',
      yamlEditTip: 'Paste or edit the YAML configuration of the Secret here. Ensure the namespace in the YAML matches the currently selected namespace or is omitted.',
      yamlEditorPlaceholder: 'YAML editor placeholder (e.g., using Monaco Editor or Codemirror)',
      securityWarning: 'Security Warning: When editing Secrets, values in the `data` field should be Base64 encoded. The `stringData` field allows plain text strings (Kubernetes will automatically encode them). Please handle sensitive information carefully.',
      createSecretTitle: 'Create Secret (YAML)',
      fetchingYaml: 'Fetching YAML for Secret "{name}"',
      fetchSecretsError: 'Failed to fetch Secret data',
      operationSuccess: 'Successfully {action} Secret',
      operationCancelled: 'Delete operation cancelled',
      noNamespace: 'No available namespace',
      deleteCanceled: 'Delete operation cancelled',
      confirmDelete: 'Confirm Delete',
      editSecretTitle: 'Edit Secret',
      fetchSecretDetailsFailed: 'Failed to fetch Secret details',
      fetchNamespacesFailed: 'Failed to fetch namespaces',
      create: 'Create',
      update: 'Update',
      cannotDetermineNamespace: 'Cannot determine target namespace',
      dataFormatError: 'Data format error',
      networkError: 'Network request failed',
      invalidResponse: 'Invalid response data',
      unknownError: 'Unknown error',
      requestFailed: 'Request failed',
      dataKeys: 'Data Entries (Keys)',
      noSecretsInNamespace: 'No Secrets found in namespace'
    }
  },

  // 系统管理
  admin: {
    // 用户管理
    userManagement: {
      description: 'Manage system user accounts, roles and permissions',
      searchPlaceholder: 'Search username, email or display name...',
      statusFilter: 'Status Filter',
      roleFilter: 'Role Filter',
      loading: 'Loading user data...',
      notSet: 'Not Set',
      never: 'Never',
      
      // 统计
      stats: {
        totalUsers: 'Total Users',
        activeUsers: 'Active Users',
        admins: 'Administrators',
        newThisMonth: 'New This Month'
      },
      
      // 列
      columns: {
        avatar: 'Avatar',
        username: 'Username',
        email: 'Email',
        role: 'Role',
        status: 'Status',
        lastLogin: 'Last Login',
        createdAt: 'Created At'
      },
      
      // 状态
      status: {
        active: 'Active',
        inactive: 'Inactive'
      },
      
      // 角色
      roles: {
        admin: 'Administrator',
        editor: 'Editor',
        viewer: 'Viewer',
        unknown: 'Unknown Role'
      },

      // 角色描述
      roleDescriptions: {
        admin: 'Has all system permissions',
        editor: 'Can edit and manage resources',
        viewer: 'Can only view resources, cannot modify'
      },
      
      // 邮箱状态
      emailStatus: {
        verified: 'Verified',
        unverified: 'Unverified'
      },
      
      // 操作
      actions: {
        addUser: 'Add User',
        exportUsers: 'Export Users',
        more: 'More',
        activate: 'Activate Account',
        deactivate: 'Deactivate Account',
        resetPassword: 'Reset Password',
        sendVerification: 'Send Verification Email',
        deleteUser: 'Delete User'
      },
      
      // 批量操作
      batch: {
        selected: 'Selected {count} users',
        activate: 'Batch Activate',
        deactivate: 'Batch Deactivate',
        delete: 'Batch Delete'
      },
      
      // 消息
      messages: {
        loadFailed: 'Failed to load user list: {error}',
        invalidUser: 'Invalid user information',
        confirmToggleStatus: 'Are you sure you want to {action} user "{username}"?',
        confirmAction: 'Confirm {action}',
        toggleStatusSuccess: 'User {action} successful',
        toggleStatusFailed: 'Failed to update user status',
        confirmResetPassword: 'Are you sure you want to reset password for user "{username}"? A new password will be sent to the user via email.',
        confirmResetPasswordTitle: 'Confirm Reset Password',
        resetPasswordSuccess: 'Password reset successful, new password has been sent to user email',
        resetPasswordFailed: 'Password reset failed',
        verificationEmailSent: 'Verification email has been sent to {email}',
        verificationEmailFailed: 'Failed to send verification email',
        confirmDeleteUser: 'Are you sure you want to delete user "{username}"? This operation cannot be undone!',
        confirmDeleteTitle: 'Confirm Delete',
        confirmDeleteButton: 'Confirm Delete',
        deleteUserSuccess: 'User deleted successfully',
        deleteUserFailed: 'Failed to delete user',
        batchActivateSuccess: 'Successfully activated {count} users',
        batchActivateFailed: 'Batch activation failed',
        batchDeactivateSuccess: 'Successfully deactivated {count} users',
        batchDeactivateFailed: 'Batch deactivation failed',
        confirmBatchDelete: 'Are you sure you want to delete the selected {count} users? This operation cannot be undone!',
        confirmBatchDeleteTitle: 'Confirm Batch Delete',
        batchDeleteSuccess: 'Successfully deleted {count} users',
        batchDeleteFailed: 'Batch delete failed',
        exportSuccess: 'User data exported successfully',
        exportFailed: 'Failed to export user data'
      }
    },
    
    // 角色管理
    roleManagement: {
      description: 'Manage system roles and permission configurations',
      searchPlaceholder: 'Search role name or description...',
      typeFilter: 'Type Filter',
      loading: 'Loading role data...',
      unknown: 'Unknown',
      mainPermissions: 'Main Permissions',
      viewAll: 'View All',
      morePermissions: '+{count} more',
      createdAt: 'Created At',
      updatedAt: 'Updated At',
      noRoles: 'No role data available',
      createFirstRole: 'Create First Role',
      
      // 类型
      types: {
        system: 'System Role',
        custom: 'Custom Role'
      },
      
      // 统计
      stats: {
        users: '{count} users',
        permissions: '{count} permissions'
      },
      
      // 操作
      actions: {
        createRole: 'Create Role',
        permissionConfig: 'Permission Config',
        viewDetails: 'View Details',
        editRole: 'Edit Role',
        managePermissions: 'Manage Permissions',
        viewUsers: 'View Users',
        deleteRole: 'Delete Role'
      },
      
      // 消息
      messages: {
        loadFailed: 'Failed to load role list',
        confirmDelete: 'Are you sure you want to delete role "{name}"? This operation cannot be undone!',
        confirmDeleteTitle: 'Confirm Delete',
        confirmDeleteButton: 'Confirm Delete',
        deleteSuccess: 'Role deleted successfully',
        deleteFailed: 'Failed to delete role'
      }
    },
    
    // 权限
    permissions: {
      users: {
        read: 'View Users',
        write: 'Manage Users'
      },
      roles: {
        read: 'View Roles',
        write: 'Manage Roles'
      },
      clusters: {
        read: 'View Clusters',
        write: 'Manage Clusters'
      },
      pods: {
        read: 'View Pods',
        write: 'Manage Pods'
      },
      deployments: {
        read: 'View Deployments',
        write: 'Manage Deployments'
      },
      services: {
        read: 'View Services',
        write: 'Manage Services'
      },
      configmaps: {
        read: 'View ConfigMaps',
        write: 'Manage ConfigMaps'
      },
      secrets: {
        read: 'View Secrets',
        write: 'Manage Secrets'
      }
    }
  }
}