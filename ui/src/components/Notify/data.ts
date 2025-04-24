export interface ListItem {
  avatar?: string
  title: string
  datetime?: string
  description?: string
  status?: "success" | "warning" | "info" | "danger" | "primary"
  extra?: string
}

export const notifyData: ListItem[] = [
  {
    avatar: "",
    title: "2025年新的启程",
    datetime: "2025年1月1日",
    description: "新年快乐，新的一年好好努力啊！"
  },
  {
    avatar: "",
    title: "CiliKube 上线啦",
    datetime: "2025年5月1日",
    description: "一个免费开源的k8s管理系统基础解决方案，前后端分离、均采用最新的技术栈"
  },
  {
    avatar: "",
    title: "新版本发布",
    datetime: "2025年6月1日",
    description: "期待已久的v1.0版本终于发布了，感谢大家的支持！"
  }
]

export const messageData: ListItem[] = [
  {
    avatar: "https://gw.alipayobjects.com/zos/rmsportal/ThXAXghbEsBCCSDihZxY.png",
    title: "打工人早安",
    description: "今天搬砖不狠，明天地位不稳，早安，打工人",
    datetime: "2025-1-1"
  },
  {
    avatar: "https://gw.alipayobjects.com/zos/rmsportal/ThXAXghbEsBCCSDihZxY.png",
    title: "打工人午安",
    description: "今天搬砖不狠，明天地位不稳，午安，打工人",
    datetime: "2025-6-1"
  },
  {
    avatar: "https://gw.alipayobjects.com/zos/rmsportal/ThXAXghbEsBCCSDihZxY.png",
    title: "打工人晚安",
    description: "今天搬砖不狠，明天地位不稳，晚安，打工人",
    datetime: "2025-12-1"
  }
]

export const todoData: ListItem[] = [
  {
    title: "任务一",
    description: "关注希里安公众号",
    extra: "未开始",
    status: "info"
  },
  {
    title: "任务二",
    description: "添加项目团队微信ciliverse",
    extra: "进行中",
    status: "info"
  },
  {
    title: "任务三",
    description: "加希里安技术交流群吹牛",
    extra: "已超时",
    status: "danger"
  }
]
