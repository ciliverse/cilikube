/** Skip opacity enter animations on mobile / reduced-motion (avoids stuck invisible UI). */
export function shouldSkipEnterAnim(): boolean {
  if (typeof window === 'undefined') return true
  try {
    return (
      window.matchMedia('(prefers-reduced-motion: reduce)').matches ||
      window.matchMedia('(max-width: 768px), (pointer: coarse)').matches
    )
  } catch {
    return false
  }
}
