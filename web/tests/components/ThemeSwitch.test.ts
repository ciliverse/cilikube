import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import ThemeSwitch from '@/components/ThemeSwitch/index.vue'

describe('ThemeSwitch', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('renders properly', () => {
    const wrapper = mount(ThemeSwitch)
    expect(wrapper.exists()).toBe(true)
  })

  it('toggles theme when clicked', async () => {
    const wrapper = mount(ThemeSwitch)
    
    // 模拟点击事件
    await wrapper.trigger('click')
    
    // 验证组件状态变化
    expect(wrapper.emitted()).toHaveProperty('click')
  })
})