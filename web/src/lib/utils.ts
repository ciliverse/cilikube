import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatPercent(value?: string | number) {
  if (value === undefined || value === null || value === '') return '-'
  const n = typeof value === 'number' ? value : parseFloat(String(value).replace('%', ''))
  if (Number.isNaN(n)) return String(value)
  return `${n.toFixed(1)}%`
}
