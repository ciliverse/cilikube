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
    title: "notifications.newYear.title",
    datetime: "2025-01-01",
    description: "notifications.newYear.description"
  },
  {
    avatar: "",
    title: "notifications.cilikubeLaunch.title",
    datetime: "2025-05-01",
    description: "notifications.cilikubeLaunch.description"
  },
  {
    avatar: "",
    title: "notifications.newVersion.title",
    datetime: "2025-06-01",
    description: "notifications.newVersion.description"
  }
]

export const messageData: ListItem[] = [
  {
    avatar: "",
    title: "messages.morning.title",
    description: "messages.morning.description",
    datetime: "2025-01-01"
  },
  {
    avatar: "",
    title: "messages.noon.title",
    description: "messages.noon.description",
    datetime: "2025-06-01"
  },
  {
    avatar: "",
    title: "messages.evening.title",
    description: "messages.evening.description",
    datetime: "2025-12-01"
  }
]

export const todoData: ListItem[] = [
  {
    title: "todos.task1.title",
    description: "todos.task1.description",
    extra: "todos.status.notStarted",
    status: "info"
  },
  {
    title: "todos.task2.title",
    description: "todos.task2.description",
    extra: "todos.status.inProgress",
    status: "info"
  },
  {
    title: "todos.task3.title",
    description: "todos.task3.description",
    extra: "todos.status.overdue",
    status: "danger"
  }
]
