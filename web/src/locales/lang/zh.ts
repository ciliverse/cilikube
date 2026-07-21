export default {
  // 通用
  common: {
    confirm: '确认',
    cancel: '取消',
    save: '保存',
    delete: '删除',
    edit: '编辑',
    add: '添加',
    search: '搜索',
    reset: '重置',
    refresh: '刷新',
    loading: '加载中...',
    success: '成功',
    error: '错误',
    warning: '警告',
    info: '信息',
    retry: '重试',
    yes: '是',
    no: '否',
    view: '查看',
    name: '名称',
    namespace: '命名空间',
    createdAt: '创建时间',
    actions: '操作',
    all: '全部'
  },
  
  // 设置
  settings: {
    title: '设置',
    language: '语言',
    font: '字体',
    theme: '主题',
    reset: '重置',
    
    // 布局设置
    layout: {
      title: '布局配置',
      showTagsView: '显示标签栏',
      showLogo: '显示 Logo',
      fixedHeader: '固定 Header',
      showFooter: '显示页脚 Footer',
      showNotify: '显示消息通知',
      showThemeSwitch: '显示切换主题按钮',
      showScreenfull: '显示全屏按钮',
      showSearchMenu: '显示搜索按钮',
      showFontSelector: '显示字体选择器',
      cacheTagsView: '是否缓存标签栏',
      showWatermark: '开启系统水印',
      showGreyMode: '显示灰色模式',
      showColorWeakness: '显示色弱模式',
      modes: {
        left: '左侧模式',
        top: '顶部模式',
        mixed: '混合模式'
      }
    },
    
    // 功能设置
    features: {
      title: '功能配置'
    }
  },
  
  // 语言选项
  language: {
    english: 'English',
    chinese: '中文'
  },
  
  // 字体选项
  font: {
    'maple-mono': 'Maple Mono',
    'kai': '楷体'
  },

  // 主题选项
  theme: {
    default: '默认',
    dark: '黑暗',
    darkBlue: '深蓝'
  },
  
  // 导航
  nav: {
    home: '首页',
    dashboard: '仪表板',
    settings: '设置'
  },

  // 菜单
  menu: {
    // 看板
    board: '看板',
    clusterInstall: '集群安装',
    clusterOverview: '集群概览',
    clusterMonitor: '集群监控',
    opsNavigation: '运维导航',
    
    // 集群
    cluster: '集群',
    clusterManagement: '集群管理',
    clusterEvents: '集群事件',
    node: '节点',
    namespace: '命名空间',
    
    // 工作负载
    workloads: '工作负载',
    pod: 'Pod',
    deployment: 'Deployment',
    
    // 存储
    storage: '存储',
    pv: 'PV',
    pvc: 'PVC',
    
    // 网络
    network: '网络',
    service: 'Service',
    ingress: 'Ingress',
    
    // 配置管理
    config: '配置管理',
    configmap: 'ConfigMap',
    secret: 'Secret',
    
    // 项目管理
    project: '项目管理',
    chineseDocs: '中文文档',
    myBlog: '我的博客',
    techStack: '技术栈',
    
    // 权限
    permission: '权限演示',
    pageLevel: '页面级权限',
    buttonLevel: '按钮级权限',
    
    // 系统管理
    admin: '系统管理',
    userManagement: '用户管理',
    roleManagement: '角色管理',
    systemSettings: '系统设置',
    
    // CRD
    crd: '自定义资源'
  },

  // 操作
  actions: {
    profile: '个人资料',
    logout: '退出登录',
    refresh: '刷新列表',
    refreshing: '正在刷新集群列表...',
    switchedToCluster: '已切换到集群: {name}',
    unknownCluster: '未知集群',
    noClusterSelected: '未选择集群',
    noClustersAvailable: '无可用集群',
    loading: '加载中...'
  },

  // 标签页操作
  tagsView: {
    refresh: '刷新',
    close: '关闭',
    closeOthers: '关闭其它',
    closeAll: '关闭所有'
  },

  // 应用消息
  app: {
    welcomeTitle: '你好呀！朋友！',
    welcomeMessage: '这是纯爱发电项目,点个免费star支持一下呀！如果可以，关注下公众号-希里安，共同学习进步！',
    migrationSuccess: '已自动迁移您的集群选择设置到新版本'
  },

  // UI 组件
  ui: {
    themeSwitch: '主题切换',
    notifications: '消息通知',
    messages: '消息',
    todos: '待办',
    viewHistory: '查看{type}历史',
    searchMenu: '搜索菜单',
    fullscreen: {
      enter: '全屏',
      exit: '退出全屏',
      browserNotSupported: '您的浏览器无法工作',
      contentEnlarge: '内容区放大',
      contentRestore: '内容区复原',
      contentFullscreen: '内容区全屏'
    }
  },

  // 通知消息
  notifications: {
    newYear: {
      title: '2025年新的启程',
      description: '新年快乐，新的一年好好努力啊！'
    },
    cilikubeLaunch: {
      title: 'CiliKube 上线啦',
      description: '一个免费开源的k8s管理系统基础解决方案，前后端分离、均采用最新的技术栈'
    },
    newVersion: {
      title: '新版本发布',
      description: '期待已久的v1.0版本终于发布了，感谢大家的支持！'
    }
  },

  // 消息
  messages: {
    morning: {
      title: '打工人早安',
      description: '今天搬砖不狠，明天地位不稳，早安，打工人'
    },
    noon: {
      title: '打工人午安',
      description: '今天搬砖不狠，明天地位不稳，午安，打工人'
    },
    evening: {
      title: '打工人晚安',
      description: '今天搬砖不狠，明天地位不稳，晚安，打工人'
    }
  },

  // 待办事项
  todos: {
    task1: {
      title: '任务一',
      description: '关注希里安公众号'
    },
    task2: {
      title: '任务二',
      description: '添加项目团队微信ciliverse'
    },
    task3: {
      title: '任务三',
      description: '加希里安技术交流群吹牛'
    },
    status: {
      notStarted: '未开始',
      inProgress: '进行中',
      overdue: '已超时'
    }
  },

  // 集群概览
  clusterOverview: {
    description: '查看集群的整体状态和资源使用情况',
    clusterStatus: '集群状态',
    healthy: '健康',
    running: '运行中',
    totalNodes: '节点总数',
    runningPods: '运行中的Pod',
    namespaces: '命名空间',
    resourceOverview: '资源概览',
    resourceDescription: '集群资源使用情况和性能指标',
    totalResources: '总资源',
    activeTypes: '活跃类型',
    noResourceData: '暂无资源数据',
    nodeMetrics: '节点监控'
  },

  // 节点监控
  nodeMetrics: {
    avgCpuUsage: '平均CPU使用率',
    avgMemoryUsage: '平均内存使用率',
    healthyNodes: '健康节点',
    totalCpuCores: '总CPU核心',
    totalMemory: '总内存',
    diskUsage: '磁盘使用率',
    networkIO: '网络I/O',
    cores: '核',
    cpuUsage: 'CPU使用率',
    memoryUsage: '内存使用率',
    noData: '暂无监控数据',
    resourceRequests: '资源请求',
    cpuRequests: 'CPU请求',
    memoryRequests: '内存请求'
  },

  // 集群事件
  clusterEvents: {
    title: '集群事件',
    realtime: '实时',
    refresh: '刷新',
    viewAll: '查看全部事件',
    noEvents: '暂无事件',
    fetchError: '获取事件失败',
    eventList: '事件列表',
    totalEvents: '共 {count} 个事件',
    count: '{count} 次',
    firstTime: '首次时间',
    filters: {
      namespace: '命名空间',
      type: '类型',
      limit: '数量限制',
      allNamespaces: '所有命名空间',
      allTypes: '所有类型'
    },
    eventTypes: {
      normal: '正常',
      warning: '警告',
      error: '错误',
      info: '信息'
    },
    timeAgo: {
      minutesAgo: '分钟前',
      hoursAgo: '小时前',
      daysAgo: '天前',
      justNow: '刚刚'
    }
  },

  // 集群管理
  clusterManagement: {
    title: '集群管理',
    description: '管理您的 Kubernetes 集群，支持多集群切换和监控',
    addCluster: '添加集群',
    searchPlaceholder: '搜索集群名称、环境或服务器地址...',
    filterEnvironment: '筛选环境',
    filterStatus: '筛选状态',
    allEnvironments: '全部环境',
    allStatuses: '全部状态',
    cardView: '卡片视图',
    tableView: '列表视图',
    
    // 统计
    stats: {
      totalClusters: '集群总数',
      activeClusters: '活跃集群',
      healthyClusters: '健康集群',
      unhealthyClusters: '异常集群'
    },
    
    // 环境
    environments: {
      production: '生产环境',
      staging: '预发布',
      development: '开发环境',
      testing: '测试环境'
    },
    
    // 状态
    statuses: {
      available: '可用',
      unavailable: '不可用'
    },
    
    // 表格列
    columns: {
      name: '集群名称',
      clusterName: '集群名称',
      status: '状态',
      server: '服务器地址',
      apiServer: 'API 服务器地址',
      version: 'K8s 版本',
      environment: '环境',
      source: '来源',
      actions: '操作'
    },
    
    // 操作
    actions: {
      setActive: '设为活动',
      edit: '编辑',
      delete: '删除',
      active: 'Active'
    },
    
    // 来源
    sources: {
      database: '数据库',
      file: '文件'
    },
    
    // 表单
    form: {
      name: '集群名称',
      namePlaceholder: '请输入唯一的集群名称',
      environment: '环境',
      environmentPlaceholder: '选择环境类型',
      provider: '提供商',
      providerPlaceholder: '例如: minikube, aws, gcp, aliyun',
      region: '区域',
      regionPlaceholder: '例如: us-west-2, cn-beijing-1',
      description: '描述',
      descriptionPlaceholder: '请输入集群的描述信息',
      kubeconfig: 'Kubeconfig',
      kubeconfigTitle: 'Kubeconfig 配置',
      kubeconfigPlaceholder: '请粘贴您的 Kubeconfig 文件内容',
      kubeconfigTips: '请确保您的 Kubeconfig 文件有效且具有适当的访问权限',
      clusterNameRequired: '集群名称不能为空',
      kubeconfigRequired: 'Kubeconfig 内容不能为空',
      environmentRequired: '请选择环境类型'
    },
    
    // 基本信息
    basicInfo: '基本信息',
    clustersFound: '个集群',
    
    // 验证
    validation: {
      nameRequired: '集群名称不能为空',
      kubeconfigRequired: 'Kubeconfig 内容不能为空',
      environmentRequired: '请选择环境类型'
    },
    
    // 消息
    messages: {
      addSuccess: '集群添加成功',
      addFailed: '添加集群失败',
      deleteConfirm: '确定要删除这个集群吗？',
      deleteSuccess: '集群删除成功',
      deleteFailed: '删除集群失败',
      setActiveSuccess: '集群 "{name}" 已设为活动集群',
      setActiveFailed: '切换活动集群失败',
      fetchFailed: '获取集群数据失败，请检查网络或联系管理员',
      noMatchingClusters: '没有找到匹配的集群',
      canOnlyDeleteDatabaseClusters: '只能删除通过数据库添加的集群'
    }
  },

  // CRD
  crd: {
    description: '管理自定义资源定义和自定义资源',
    searchPlaceholder: '搜索CRD...',
    selectGroup: '选择组',
    selectNamespace: '选择命名空间',
    noCRDs: '未找到自定义资源定义',
    noResources: '未找到自定义资源',
    
    // CRD字段
    name: '名称',
    group: '组',
    version: '版本',
    kind: '类型',
    plural: '复数形式',
    singular: '单数形式',
    scope: '作用域',
    categories: '分类',
    shortNames: '简称',
    apiVersion: 'API版本',
    
    // CRD详情
    basicInfo: '基本信息',
    versions: '版本信息',
    conditions: '状态条件',
    metadata: '元数据',
    labels: '标签',
    annotations: '注解',
    spec: '规格',
    status: '状态',
    
    // 版本字段
    served: '已启用',
    storage: '存储版本',
    
    // 条件字段
    conditionType: '类型',
    conditionStatus: '状态',
    conditionReason: '原因',
    conditionMessage: '消息',
    lastTransitionTime: '最后转换时间',
    
    // 操作
    viewResources: '查看资源',
    createResource: '创建资源',
    editResource: '编辑资源',
    deleteResource: '删除资源',
    resources: '资源',
    
    // 表单
    enterName: '输入资源名称',
    nameRequired: '名称是必填项',
    nameFormat: '名称必须是小写字母数字和连字符',
    namespaceRequired: '命名空间是必填项',
    labelKey: '标签键',
    labelValue: '标签值',
    addLabel: '添加标签',
    specPlaceholder: '以JSON格式输入资源规格',
    invalidJSON: '无效的JSON格式',
    
    // 消息
    deleteConfirm: '确定要删除资源 "{name}" 吗？',
    
    // 错误消息
    fetchCRDsFailed: '获取CRD列表失败',
    fetchCRDDetailFailed: '获取CRD详情失败',
    fetchResourcesFailed: '获取自定义资源列表失败',
    fetchResourceDetailFailed: '获取自定义资源详情失败',
    createResourceFailed: '创建自定义资源失败',
    updateResourceFailed: '更新自定义资源失败',
    deleteResourceFailed: '删除自定义资源失败',
    
    // 成功消息
    createResourceSuccess: '自定义资源创建成功',
    updateResourceSuccess: '自定义资源更新成功',
    deleteResourceSuccess: '自定义资源删除成功'
  },

  // 节点管理
  nodeManagement: {
    title: '节点',
    searchPlaceholder: '搜索节点...',
    filterRole: '按角色筛选',
    filterStatus: '按状态筛选',
    allRoles: '所有角色',
    allStatuses: '所有状态',
    viewToggle: {
      card: '卡片视图',
      list: '列表视图'
    },
    columns: {
      name: '名称',
      status: '状态',
      role: '角色',
      version: '版本',
      internalIP: '内部IP',
      externalIP: '外部IP',
      os: '操作系统',
      kernel: '内核',
      runtime: '运行时',
      cpu: 'CPU',
      memory: '内存',
      pods: 'Pod数量',
      created: '创建时间'
    },
    statuses: {
      ready: '就绪',
      notReady: '未就绪',
      unknown: '未知'
    },
    roles: {
      controlPlane: '控制平面',
      worker: '工作节点'
    },
    info: {
      addresses: '地址信息',
      capacity: '容量',
      allocatable: '可分配',
      systemInfo: '系统信息',
      conditions: '状态条件',
      taints: '污点',
      kubeletVersion: 'Kubelet 版本',
      kubeProxyVersion: 'Kube-proxy 版本',
      containerRuntime: '容器运行时',
      kernelVersion: '内核版本',
      osImage: '操作系统镜像',
      podCIDR: 'Pod CIDR',
      allAddresses: '所有地址'
    },
    actions: {
      cordon: '封锁',
      uncordon: '解封',
      drain: '驱逐',
      viewDetails: '查看详情'
    },
    messages: {
      noNodes: '当前集群中未找到任何节点',
      noMatchingNodes: '没有符合搜索条件的节点',
      selectCluster: '请先在导航栏选择一个集群以查看节点信息。',
      loadingNodes: '正在为集群加载节点'
    }
  },

  // 命名空间管理
  namespaceManagement: {
    title: '命名空间',
    searchPlaceholder: '搜索命名空间...',
    filterStatus: '按状态筛选',
    allStatuses: '所有状态',
    viewToggle: {
      card: '卡片视图',
      list: '列表视图'
    },
    columns: {
      name: '名称',
      status: '状态',
      created: '创建时间',
      labels: '标签',
      annotations: '注解'
    },
    statuses: {
      active: '活跃',
      terminating: '终止中',
      unknown: '未知'
    },
    info: {
      systemNamespace: '系统命名空间',
      userNamespace: '用户命名空间',
      resourceCount: '资源数量',
      description: '描述'
    },
    actions: {
      create: '创建命名空间',
      viewDetails: '查看详情',
      editMetadata: '编辑元数据',
      delete: '删除'
    },
    form: {
      name: '名称',
      namePlaceholder: '输入命名空间名称',
      nameRequired: '命名空间名称不能为空',
      nameFormat: '名称必须由小写字母、数字或"-"组成，并以字母或数字开头和结尾',
      nameLength: '名称不能超过63个字符',
      labels: '标签',
      labelsPlaceholder: '输入标签，格式：key=value',
      annotations: '注解',
      annotationsPlaceholder: '输入注解，格式：key=value',
      description: '描述',
      descriptionPlaceholder: '输入命名空间描述'
    },
    messages: {
      noNamespaces: '未找到任何命名空间',
      noMatchingNamespaces: '没有符合搜索条件的命名空间',
      selectCluster: '请先选择一个集群以查看命名空间。',
      createSuccess: '命名空间创建成功',
      createFailed: '命名空间创建失败',
      deleteConfirm: '确定要删除命名空间 "{name}" 吗？此操作将删除该空间下的所有资源且不可恢复！',
      deleteSuccess: '命名空间删除成功',
      deleteFailed: '命名空间删除失败',
      cannotDeleteSystem: '不能删除系统命名空间',
      loadingNamespaces: '正在加载命名空间'
    }
  },

  // Deployment 管理
  deploymentManagement: {
    title: 'Deployments',
    searchPlaceholder: '搜索 Deployment 名称...',
    filterNamespace: '筛选命名空间',
    filterStatus: '按状态筛选',
    allNamespaces: '所有命名空间',
    allStatuses: '所有状态',
    viewToggle: {
      card: '卡片视图',
      list: '列表视图'
    },
    columns: {
      name: '名称',
      namespace: '命名空间',
      status: '状态',
      replicas: '副本数',
      ready: '就绪副本',
      upToDate: '最新副本',
      available: '可用副本',
      age: '运行时间',
      images: '镜像'
    },
    statuses: {
      available: '可用',
      progressing: '进行中',
      replicaFailure: '副本失败',
      degraded: '降级',
      unknown: '未知'
    },
    info: {
      strategy: '更新策略',
      selector: '选择器',
      template: 'Pod 模板',
      conditions: '状态条件'
    },
    actions: {
      create: '创建 Deployment',
      scale: '扩缩容',
      rollout: '滚动更新',
      restart: '重启',
      edit: '编辑 YAML',
      delete: '删除',
      viewDetails: '查看详情',
      viewPods: '查看 Pods'
    },
    form: {
      createFromYaml: '从 YAML 创建',
      editYaml: '编辑 YAML',
      yamlContent: 'YAML 内容',
      name: '名称',
      namespace: '命名空间',
      replicas: '副本数',
      image: '镜像',
      nameRequired: '名称不能为空',
      namespaceRequired: '命名空间不能为空',
      replicasRequired: '副本数不能为空',
      imageRequired: '镜像不能为空',
      scaleReplicas: '调整副本数',
      currentReplicas: '当前副本数',
      targetReplicas: '目标副本数'
    },
    messages: {
      noDeployments: '当前命名空间中未找到任何 Deployment',
      noMatchingDeployments: '没有符合搜索条件的 Deployment',
      selectNamespace: '请先选择一个命名空间以查看 Deployment。',
      createSuccess: 'Deployment 创建成功',
      createFailed: 'Deployment 创建失败',
      updateSuccess: 'Deployment 更新成功',
      updateFailed: 'Deployment 更新失败',
      scaleSuccess: 'Deployment 扩缩容成功',
      scaleFailed: 'Deployment 扩缩容失败',
      deleteConfirm: '确定要删除 Deployment "{name}" 吗？',
      deleteSuccess: 'Deployment 删除成功',
      deleteFailed: 'Deployment 删除失败',
      loadingDeployments: '正在加载 Deployment'
    }
  },

  // Pod 管理
  podManagement: {
    title: 'Pods',
    searchPlaceholder: '搜索 Pod 名称、IP、节点或状态...',
    filterNamespace: '筛选命名空间',
    filterStatus: '按状态筛选',
    allNamespaces: '所有命名空间',
    allStatuses: '所有状态',
    viewToggle: {
      card: '卡片视图',
      list: '列表视图'
    },
    columns: {
      name: '名称',
      namespace: '命名空间',
      status: '状态',
      ip: 'Pod IP',
      node: '所在节点',
      created: '创建时间',
      restarts: '重启次数',
      ready: '就绪状态'
    },
    statuses: {
      running: '运行中',
      pending: '等待中',
      succeeded: '成功',
      failed: '失败',
      unknown: '未知',
      terminating: '终止中',
      crashLoopBackOff: '崩溃循环',
      imagePullBackOff: '镜像拉取失败',
      containerCreating: '容器创建中'
    },
    info: {
      containers: '容器',
      restartPolicy: '重启策略',
      qosClass: 'QoS 类别',
      phase: '阶段',
      conditions: '状态条件'
    },
    actions: {
      create: '创建 Pod',
      viewLogs: '查看日志',
      exec: '进入容器',
      edit: '编辑 YAML',
      delete: '删除',
      viewDetails: '查看详情'
    },
    form: {
      createFromYaml: '从 YAML 创建',
      editYaml: '编辑 YAML',
      yamlContent: 'YAML 内容',
      container: '容器',
      command: '命令',
      tailLines: '显示行数',
      follow: '跟踪日志'
    },
    messages: {
      noPods: '当前命名空间中未找到任何 Pod',
      noMatchingPods: '没有符合搜索条件的 Pod',
      selectNamespace: '请先选择一个命名空间以查看 Pod。',
      createSuccess: 'Pod 创建成功',
      createFailed: 'Pod 创建失败',
      updateSuccess: 'Pod 更新成功',
      updateFailed: 'Pod 更新失败',
      deleteConfirm: '确定要删除 Pod "{name}" 吗？',
      deleteSuccess: 'Pod 删除成功',
      deleteFailed: 'Pod 删除失败',
      loadingPods: '正在加载 Pod',
      loadingContainers: '正在加载容器列表',
      connectingTerminal: '正在连接终端',
      terminalConnected: '终端已连接',
      terminalDisconnected: '终端连接已断开'
    }
  },

  // PersistentVolume 管理
  persistentVolumeManagement: {
    title: 'PersistentVolumes',
    searchPlaceholder: '搜索 PV 名称 / StorageClass',
    columns: {
      name: '名称',
      status: '状态',
      capacity: '容量',
      accessModes: '访问模式',
      reclaimPolicy: '回收策略',
      volumeMode: '卷模式',
      storageClass: 'StorageClass',
      boundClaim: '绑定到 (PVC)',
      createdAt: '创建时间',
      actions: '操作'
    },
    actions: {
      create: '创建 PV',
      edit: '编辑',
      delete: '删除',
      editYaml: '编辑 YAML',
      applyYaml: '应用 YAML'
    },
    messages: {
      noPVs: '未找到 PersistentVolumes',
      deleteConfirm: '确定要删除 PersistentVolume "{name}" 吗？',
      deleteSuccess: 'PersistentVolume 删除成功',
      deleteFailed: 'PersistentVolume 删除失败',
      yamlEditTip: '请在此处粘贴或编辑 PersistentVolume 的 YAML 配置。PV 是集群范围资源，不属于任何命名空间。',
      yamlEditorPlaceholder: 'YAML 编辑器占位符 (例如: 使用 Monaco Editor 或 Codemirror)',
      createPVTitle: '创建 PV (YAML)',
      fetchingYaml: '模拟: 获取 PV "{name}" 的 YAML'
    }
  },

  // PersistentVolumeClaim 管理
  persistentVolumeClaimManagement: {
    title: 'PersistentVolumeClaims',
    filterNamespace: '选择命名空间',
    searchPlaceholder: '搜索 PVC 名称 / 绑定 PV',
    columns: {
      name: '名称',
      namespace: '命名空间',
      status: '状态',
      boundVolume: '绑定卷 (PV)',
      requestCapacity: '请求容量',
      actualCapacity: '实际容量',
      accessModes: '访问模式',
      volumeMode: '卷模式',
      storageClass: 'StorageClass',
      createdAt: '创建时间',
      actions: '操作'
    },
    actions: {
      create: '创建 PVC',
      edit: '编辑',
      delete: '删除',
      editYaml: '编辑 YAML',
      applyYaml: '应用 YAML'
    },
    messages: {
      noPVCs: '未找到 PersistentVolumeClaims',
      noNamespaces: '无可用命名空间',
      selectNamespace: '请先选择一个命名空间以加载 PVCs',
      selectNamespaceWarning: '请先选择一个命名空间',
      deleteConfirm: '确定要删除 PersistentVolumeClaim "{name}" 吗？',
      deleteConfirmDetail: '确定要删除 PVC "{name}" (命名空间: {namespace}) 吗？请注意 PV 的回收策略。',
      deleteSuccess: 'PersistentVolumeClaim 删除成功',
      deleteFailed: 'PersistentVolumeClaim 删除失败',
      yamlEditTip: '请在此处粘贴或编辑 PersistentVolumeClaim 的 YAML 配置。确保 YAML 中的 namespace 与当前选定的命名空间匹配或省略。',
      yamlEditorPlaceholder: 'YAML 编辑器占位符 (例如: 使用 Monaco Editor 或 Codemirror)',
      createPVCTitle: '创建 PVC (YAML)',
      fetchingYaml: '获取 PVC "{name}" 的 YAML',
      editError: '无法编辑，缺少原始数据',
      fetchNamespacesError: '获取命名空间失败',
      fetchPVCsError: '获取 PVC 数据失败',
      operationSuccess: '模拟 PVC {action} 成功',
      operationCancelled: '删除操作已取消',
      noNamespace: '无可用命名空间',
      deleteCanceled: '删除操作已取消',
      confirmDelete: '确认删除',
      editSecretTitle: '编辑 Secret',
      fetchSecretDetailsFailed: '获取 Secret 详情失败',
      fetchNamespacesFailed: '获取命名空间失败',
      create: '创建',
      update: '更新',
      cannotDetermineNamespace: '无法确定目标命名空间',
      dataFormatError: '数据格式错误',
      networkError: '网络请求失败',
      invalidResponse: '无效的响应数据',
      unknownError: '未知错误',
      requestFailed: '请求失败',
      dataKeys: '数据条目 (Keys)',
      noSecretsInNamespace: '在命名空间',
      aboutPVCTitle: '关于 PersistentVolumeClaims',
      aboutPVCDescription: 'PersistentVolumeClaim (PVC) 是用户对存储的请求。它与 Pod 相似，Pod 会耗用节点资源，而 PVC 耗用 PV 资源。Pod 可以请求特定数量的资源（CPU 和内存）；同样 PVC 可以请求特定的大小和访问模式（例如，可以要求 PV 挂载为 ReadWriteOnce、ReadOnlyMany、ReadWriteMany 或 ReadWriteOncePod）。'
    }
  },

  // Service 管理
  serviceManagement: {
    title: 'Services',
    filterNamespace: '选择命名空间',
    searchPlaceholder: '搜索 Service 名称 / 类型 / ClusterIP',
    columns: {
      name: '名称',
      namespace: '命名空间',
      type: '类型',
      clusterIP: 'ClusterIP',
      externalIP: 'External IP',
      ports: '端口',
      selector: '选择器',
      createdAt: '创建时间',
      actions: '操作'
    },
    actions: {
      create: '创建 Service',
      edit: '编辑',
      delete: '删除',
      editYaml: '编辑 YAML',
      applyYaml: '应用 YAML'
    },
    messages: {
      noServices: '未找到 Services',
      noNamespaces: '无可用命名空间',
      selectNamespace: '请先选择一个命名空间以加载 Services',
      selectNamespaceWarning: '请先选择一个命名空间',
      deleteConfirm: '确定要删除 Service "{name}" 吗？',
      deleteConfirmDetail: '确定要删除 Service "{name}" (命名空间: {namespace}) 吗？',
      deleteSuccess: 'Service 删除成功',
      deleteFailed: 'Service 删除失败',
      yamlEditTip: '请在此处粘贴或编辑 Service 的 YAML 配置。确保 YAML 中的 namespace 与当前选定的命名空间匹配或省略。',
      yamlEditorPlaceholder: 'YAML 编辑器占位符 (例如: 使用 Monaco Editor 或 Codemirror)',
      createServiceTitle: '创建 Service (YAML)',
      fetchingYaml: '获取 Service "{name}" 的 YAML',
      fetchServicesError: '获取 Service 数据失败',
      operationSuccess: '模拟 Service {action} 成功',
      operationCancelled: '删除操作已取消',
      noNamespace: '无可用命名空间',
      deleteCanceled: '删除操作已取消',
      confirmDelete: '确认删除',
      editSecretTitle: '编辑 Secret',
      fetchSecretDetailsFailed: '获取 Secret 详情失败',
      fetchNamespacesFailed: '获取命名空间失败',
      create: '创建',
      update: '更新',
      cannotDetermineNamespace: '无法确定目标命名空间',
      dataFormatError: '数据格式错误',
      networkError: '网络请求失败',
      invalidResponse: '无效的响应数据',
      unknownError: '未知错误',
      requestFailed: '请求失败',
      dataKeys: '数据条目 (Keys)',
      noSecretsInNamespace: '在命名空间'
    }
  },

  // Ingress 管理
  ingressManagement: {
    title: 'Ingresses',
    filterNamespace: '选择命名空间',
    searchPlaceholder: '搜索 Ingress 名称 / Host / Class',
    columns: {
      name: '名称',
      namespace: '命名空间',
      class: 'Class',
      hosts: 'Hosts',
      address: 'Address (LoadBalancer IP)',
      ports: '端口',
      rules: '规则',
      createdAt: '创建时间',
      actions: '操作'
    },
    actions: {
      create: '创建 Ingress',
      edit: '编辑',
      delete: '删除',
      editYaml: '编辑 YAML',
      applyYaml: '应用 YAML'
    },
    messages: {
      noIngresses: '未找到 Ingresses',
      noNamespaces: '无可用命名空间',
      selectNamespace: '请先选择一个命名空间以加载 Ingresses',
      selectNamespaceWarning: '请先选择一个命名空间',
      deleteConfirm: '确定要删除 Ingress "{name}" 吗？',
      deleteConfirmDetail: '确定要删除 Ingress "{name}" (命名空间: {namespace}) 吗？',
      deleteSuccess: 'Ingress 删除成功',
      deleteFailed: 'Ingress 删除失败',
      yamlEditTip: '请在此处粘贴或编辑 Ingress 的 YAML 配置。确保 YAML 中的 namespace 与当前选定的命名空间匹配或省略。',
      yamlEditorPlaceholder: 'YAML 编辑器占位符 (例如: 使用 Monaco Editor 或 Codemirror)',
      createIngressTitle: '创建 Ingress (YAML)',
      fetchingYaml: '获取 Ingress "{name}" 的 YAML',
      fetchIngressesError: '获取 Ingress 数据失败',
      operationSuccess: '模拟 Ingress {action} 成功',
      operationCancelled: '删除操作已取消',
      noNamespace: '无可用命名空间',
      deleteCanceled: '删除操作已取消',
      confirmDelete: '确认删除',
      editSecretTitle: '编辑 Secret',
      fetchSecretDetailsFailed: '获取 Secret 详情失败',
      fetchNamespacesFailed: '获取命名空间失败',
      create: '创建',
      update: '更新',
      cannotDetermineNamespace: '无法确定目标命名空间',
      dataFormatError: '数据格式错误',
      networkError: '网络请求失败',
      invalidResponse: '无效的响应数据',
      unknownError: '未知错误',
      requestFailed: '请求失败',
      dataKeys: '数据条目 (Keys)',
      noSecretsInNamespace: '在命名空间'
    }
  },

  // ConfigMap 管理
  configmapManagement: {
    title: 'ConfigMaps',
    filterNamespace: '选择命名空间',
    searchPlaceholder: '搜索 ConfigMap 名称',
    columns: {
      name: '名称',
      namespace: '命名空间',
      dataCount: '数据条目 (Keys)',
      createdAt: '创建时间',
      actions: '操作'
    },
    actions: {
      create: '创建 ConfigMap',
      edit: '编辑',
      delete: '删除',
      editYaml: '编辑 YAML',
      applyYaml: '应用 YAML'
    },
    messages: {
      noConfigMaps: '未找到 ConfigMaps',
      noNamespaces: '无可用命名空间',
      selectNamespace: '请先选择一个命名空间以加载 ConfigMaps',
      selectNamespaceWarning: '请先选择一个命名空间',
      deleteConfirm: '确定要删除 ConfigMap "{name}" 吗？',
      deleteConfirmDetail: '确定要删除 ConfigMap "{name}" (命名空间: {namespace}) 吗？',
      deleteSuccess: 'ConfigMap 删除成功',
      deleteFailed: 'ConfigMap 删除失败',
      yamlEditTip: '请在此处粘贴或编辑 ConfigMap 的 YAML 配置。确保 YAML 中的 namespace 与当前选定的命名空间匹配或省略。',
      yamlEditorPlaceholder: 'YAML 编辑器占位符 (例如: 使用 Monaco Editor 或 Codemirror)',
      createConfigMapTitle: '创建 ConfigMap (YAML)',
      fetchingYaml: '获取 ConfigMap "{name}" 的 YAML',
      fetchConfigMapsError: '获取 ConfigMap 数据失败',
      operationSuccess: '模拟 ConfigMap {action} 成功',
      operationCancelled: '删除操作已取消',
      noNamespace: '无可用命名空间',
      deleteCanceled: '删除操作已取消',
      confirmDelete: '确认删除',
      editSecretTitle: '编辑 Secret',
      fetchSecretDetailsFailed: '获取 Secret 详情失败',
      fetchNamespacesFailed: '获取命名空间失败',
      create: '创建',
      update: '更新',
      cannotDetermineNamespace: '无法确定目标命名空间',
      dataFormatError: '数据格式错误',
      networkError: '网络请求失败',
      invalidResponse: '无效的响应数据',
      unknownError: '未知错误',
      requestFailed: '请求失败',
      dataKeys: '数据条目 (Keys)',
      noSecretsInNamespace: '在命名空间'
    }
  },

  // Secret 管理
  secretManagement: {
    title: 'Secrets',
    filterNamespace: '选择命名空间',
    searchPlaceholder: '搜索 Secret 名称 / 类型',
    columns: {
      name: '名称',
      namespace: '命名空间',
      type: '类型',
      dataCount: '数据条目 (Keys)',
      createdAt: '创建时间',
      actions: '操作'
    },
    actions: {
      create: '创建 Secret',
      edit: '编辑',
      delete: '删除',
      editYaml: '编辑 YAML',
      applyYaml: '应用 YAML'
    },
    messages: {
      noSecrets: '未找到 Secrets',
      noNamespaces: '无可用命名空间',
      selectNamespace: '请先选择一个命名空间以加载 Secrets',
      selectNamespaceWarning: '请先选择一个命名空间',
      deleteConfirm: '确定要删除 Secret "{name}" 吗？',
      deleteConfirmDetail: '确定要删除 Secret "{name}" (命名空间: {namespace}) 吗？',
      deleteSuccess: 'Secret 删除成功',
      deleteFailed: 'Secret 删除失败',
      yamlEditTip: '请在此处粘贴或编辑 Secret 的 YAML 配置。确保 YAML 中的 namespace 与当前选定的命名空间匹配或省略。',
      yamlEditorPlaceholder: 'YAML 编辑器占位符 (例如: 使用 Monaco Editor 或 Codemirror)',
      securityWarning: '安全警告: 编辑 Secret 时，`data` 字段中的值应为 Base64 编码格式。`stringData` 字段允许使用明文字符串（Kubernetes 会自动编码）。请谨慎处理敏感信息。',
      createSecretTitle: '创建 Secret (YAML)',
      fetchingYaml: '获取 Secret "{name}" 的 YAML',
      fetchSecretsError: '获取 Secret 数据失败',
      operationSuccess: '模拟 Secret {action} 成功',
      operationCancelled: '删除操作已取消',
      noNamespace: '无可用命名空间',
      deleteCanceled: '删除操作已取消',
      confirmDelete: '确认删除',
      editSecretTitle: '编辑 Secret',
      fetchSecretDetailsFailed: '获取 Secret 详情失败',
      fetchNamespacesFailed: '获取命名空间失败',
      create: '创建',
      update: '更新',
      cannotDetermineNamespace: '无法确定目标命名空间',
      dataFormatError: '数据格式错误',
      networkError: '网络请求失败',
      invalidResponse: '无效的响应数据',
      unknownError: '未知错误',
      requestFailed: '请求失败',
      dataKeys: '数据条目 (Keys)',
      noSecretsInNamespace: '在命名空间'
    }
  },

  // 系统管理
  admin: {
    // 用户管理
    userManagement: {
      description: '管理系统用户账户、角色和权限',
      searchPlaceholder: '搜索用户名、邮箱或显示名称...',
      statusFilter: '状态筛选',
      roleFilter: '角色筛选',
      loading: '加载用户数据中...',
      notSet: '未设置',
      never: '从未',
      
      // 统计
      stats: {
        totalUsers: '总用户数',
        activeUsers: '活跃用户',
        admins: '管理员',
        newThisMonth: '本月新增'
      },
      
      // 列
      columns: {
        avatar: '头像',
        username: '用户名',
        email: '邮箱',
        role: '角色',
        status: '状态',
        lastLogin: '最后登录',
        createdAt: '创建时间'
      },
      
      // 状态
      status: {
        active: '活跃',
        inactive: '已停用'
      },
      
      // 角色
      roles: {
        admin: '管理员',
        editor: '编辑者',
        viewer: '查看者',
        unknown: '未知角色'
      },

      // 角色描述
      roleDescriptions: {
        admin: '拥有系统所有权限',
        editor: '可以编辑和管理资源',
        viewer: '只能查看资源，无法修改'
      },
      
      // 邮箱状态
      emailStatus: {
        verified: '已验证',
        unverified: '未验证'
      },
      
      // 操作
      actions: {
        addUser: '添加用户',
        exportUsers: '导出用户',
        more: '更多',
        activate: '激活账户',
        deactivate: '停用账户',
        resetPassword: '重置密码',
        sendVerification: '发送验证邮件',
        deleteUser: '删除用户'
      },
      
      // 批量操作
      batch: {
        selected: '已选择 {count} 个用户',
        activate: '批量激活',
        deactivate: '批量停用',
        delete: '批量删除'
      },
      
      // 消息
      messages: {
        loadFailed: '加载用户列表失败: {error}',
        invalidUser: '用户信息无效',
        confirmToggleStatus: '确定要{action}用户 "{username}" 吗？',
        confirmAction: '确认{action}',
        toggleStatusSuccess: '用户{action}成功',
        toggleStatusFailed: '用户状态更新失败',
        confirmResetPassword: '确定要重置用户 "{username}" 的密码吗？新密码将通过邮件发送给用户。',
        confirmResetPasswordTitle: '确认重置密码',
        resetPasswordSuccess: '密码重置成功，新密码已发送到用户邮箱',
        resetPasswordFailed: '密码重置失败',
        verificationEmailSent: '验证邮件已发送到 {email}',
        verificationEmailFailed: '发送验证邮件失败',
        confirmDeleteUser: '确定要删除用户 "{username}" 吗？此操作不可恢复！',
        confirmDeleteTitle: '确认删除',
        confirmDeleteButton: '确定删除',
        deleteUserSuccess: '用户删除成功',
        deleteUserFailed: '用户删除失败',
        batchActivateSuccess: '成功激活 {count} 个用户',
        batchActivateFailed: '批量激活失败',
        batchDeactivateSuccess: '成功停用 {count} 个用户',
        batchDeactivateFailed: '批量停用失败',
        confirmBatchDelete: '确定要删除选中的 {count} 个用户吗？此操作不可恢复！',
        confirmBatchDeleteTitle: '确认批量删除',
        batchDeleteSuccess: '成功删除 {count} 个用户',
        batchDeleteFailed: '批量删除失败',
        exportSuccess: '用户数据导出成功',
        exportFailed: '用户数据导出失败'
      }
    },
    
    // 角色管理
    roleManagement: {
      description: '管理系统角色和权限配置',
      searchPlaceholder: '搜索角色名称或描述...',
      typeFilter: '类型筛选',
      loading: '加载角色数据中...',
      unknown: '未知',
      mainPermissions: '主要权限',
      viewAll: '查看全部',
      morePermissions: '+{count} 更多',
      createdAt: '创建时间',
      updatedAt: '更新时间',
      noRoles: '暂无角色数据',
      createFirstRole: '创建第一个角色',
      
      // 类型
      types: {
        system: '系统角色',
        custom: '自定义角色'
      },
      
      // 统计
      stats: {
        users: '{count} 个用户',
        permissions: '{count} 个权限'
      },
      
      // 操作
      actions: {
        createRole: '创建角色',
        permissionConfig: '权限配置',
        viewDetails: '查看详情',
        editRole: '编辑角色',
        managePermissions: '管理权限',
        viewUsers: '查看用户',
        deleteRole: '删除角色'
      },
      
      // 消息
      messages: {
        loadFailed: '加载角色列表失败',
        confirmDelete: '确定要删除角色 "{name}" 吗？此操作不可恢复！',
        confirmDeleteTitle: '确认删除',
        confirmDeleteButton: '确定删除',
        deleteSuccess: '角色删除成功',
        deleteFailed: '删除角色失败'
      }
    },
    
    // 权限
    permissions: {
      users: {
        read: '查看用户',
        write: '管理用户'
      },
      roles: {
        read: '查看角色',
        write: '管理角色'
      },
      clusters: {
        read: '查看集群',
        write: '管理集群'
      },
      pods: {
        read: '查看Pod',
        write: '管理Pod'
      },
      deployments: {
        read: '查看Deployment',
        write: '管理Deployment'
      },
      services: {
        read: '查看Service',
        write: '管理Service'
      },
      configmaps: {
        read: '查看ConfigMap',
        write: '管理ConfigMap'
      },
      secrets: {
        read: '查看Secret',
        write: '管理Secret'
      }
    }
  }
}