<script lang="ts" setup>
import { ref } from "vue"
import { Setting } from "@element-plus/icons-vue"

interface Props {
  buttonTop?: number
}

const props = withDefaults(defineProps<Props>(), {
  buttonTop: 350
})

const buttonTopCss = props.buttonTop + "px"
const show = ref(false)
</script>

<template>
  <div class="handle-button" @click="show = true">
    <el-icon :size="24">
      <Setting />
    </el-icon>
  </div>
  <el-drawer v-model="show" size="300px" :with-header="false">
    <slot />
  </el-drawer>
</template>

<style lang="scss" scoped>
.handle-button {
  width: 48px;
  height: 48px;
  background-color: var(--ck-rightpanel-button-bg-color);
  position: fixed;
  top: v-bind(buttonTopCss);
  right: 0;
  border-radius: 6px 0 0 6px;
  z-index: 1000; /* 提高 z-index 确保按钮在最上层 */
  cursor: pointer;
  pointer-events: auto;
  color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15); /* 添加阴影增加可见性 */
  transition: all 0.3s ease; /* 添加过渡动画 */

  /* 悬停效果 */
  &:hover {
    background-color: var(--el-color-primary-light-3);
    transform: translateX(-2px); /* 悬停时稍微向左移动 */
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
  }

  /* 确保图标居中 */
  .el-icon {
    display: flex;
    align-items: center;
    justify-content: center;
  }
}
</style>
