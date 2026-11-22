/** Waiting / terminated reasons worth explaining next to a pod phase. */

const NOISE = new Set(['ContainerCreating', 'PodInitializing', 'Completed'])

type ContainerStatus = {
  state?: {
    waiting?: { reason?: string; message?: string }
    terminated?: { reason?: string; message?: string }
  }
}

export function podProblem(pod: {
  status?: { initContainerStatuses?: ContainerStatus[]; containerStatuses?: ContainerStatus[] }
} | null | undefined): { reason: string; message: string } {
  const lists = [pod?.status?.initContainerStatuses, pod?.status?.containerStatuses]
  for (const list of lists) {
    for (const container of list || []) {
      const waiting = container?.state?.waiting
      if (waiting?.reason && !NOISE.has(waiting.reason)) {
        return { reason: waiting.reason, message: waiting.message || '' }
      }
      const terminated = container?.state?.terminated
      if (terminated?.reason && !NOISE.has(terminated.reason)) {
        return { reason: terminated.reason, message: terminated.message || '' }
      }
    }
  }
  return { reason: '', message: '' }
}
