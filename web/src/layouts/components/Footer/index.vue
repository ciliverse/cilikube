<script lang="ts" setup>
const VITE_APP_TITLE = import.meta.env.VITE_APP_TITLE
</script>

<template>
  <footer class="layout-footer">
    <a href="https://github.com/ciliverse/cilikube" target="_blank" rel="noopener noreferrer" class="footer-link">
      cilliantech©2025 {{ VITE_APP_TITLE }}
    </a>
  </footer>
</template>

<style lang="scss" scoped>
.layout-footer {
  width: 100%;
  min-height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--el-text-color-placeholder);
  margin-top: auto;
  flex-shrink: 0;
  background-color: var(--el-bg-color);
  border-top: 1px solid var(--el-border-color-lighter);
}

.footer-link {
  color: var(--el-color-primary);
  text-decoration: none;
  transition: all 0.3s ease;

  &:hover {
    color: var(--el-color-primary-light-3);
    text-decoration: underline;
  }

  &:active {
    color: var(--el-color-primary-dark-2);
  }
}
</style>
