/**
 * 头像URL处理工具函数
 */

/**
 * 构建完整的头像URL
 * @param avatarUrl 头像URL（可能是相对路径或完整URL）
 * @returns 完整的头像URL
 */
export function buildAvatarUrl(avatarUrl: string | null | undefined): string {
  if (!avatarUrl) {
    return ''
  }

  // 如果已经是完整的URL（包含协议），直接返回
  if (avatarUrl.startsWith('http://') || avatarUrl.startsWith('https://') || avatarUrl.startsWith('data:')) {
    return avatarUrl
  }

  // 如果是相对路径，构建完整URL
  if (avatarUrl.startsWith('/')) {
    const baseUrl = `${window.location.protocol}//${window.location.hostname}:8080`
    return `${baseUrl}${avatarUrl}`
  }

  // 其他情况直接返回
  return avatarUrl
}

/**
 * 获取用户名首字母作为头像占位符
 * @param username 用户名
 * @returns 首字母
 */
export function getAvatarInitials(username: string): string {
  if (!username) return 'U'
  
  // 如果是中文名，取前两个字符
  if (/[\u4e00-\u9fa5]/.test(username)) {
    return username.slice(0, 2).toUpperCase()
  }
  
  // 如果是英文名，取首字母
  const words = username.split(/\s+/)
  if (words.length >= 2) {
    return (words[0][0] + words[1][0]).toUpperCase()
  }
  
  return username.slice(0, 2).toUpperCase()
}