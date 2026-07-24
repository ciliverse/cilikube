import { apiGet } from '@/lib/api'

export type ShowcaseAccount = {
  username: string
  password: string
  role: string
  note?: string
}

export type ShowcaseInfo = {
  showcase: boolean
  message?: string
  cluster?: string
  accounts?: ShowcaseAccount[]
}

/** Public endpoint — credentials only when server runs with CILIKUBE_SHOWCASE=1. */
export async function fetchShowcaseInfo(): Promise<ShowcaseInfo> {
  try {
    const data = await apiGet<ShowcaseInfo>('/api/v1/showcase/info')
    return data || { showcase: false }
  } catch {
    return { showcase: false }
  }
}
