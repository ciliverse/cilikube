import { requestGo } from "@/utils/service"
import type { Node } from "./types/node"

/** 获取节点列表 (使用当前活动集群) */
export function getNodeList() {
  // 后端直接返回 Node[]，使用当前活动集群
  return requestGo<Node[]>({
    url: `/api/v1/nodes`,
    method: "get"
  })
}