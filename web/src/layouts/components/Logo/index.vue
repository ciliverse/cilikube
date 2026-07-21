<script lang="ts" setup>
// 导入布局模式钩子和 logo 图片资源
import { useLayoutMode } from "@/hooks/useLayoutMode"
import logo from "@/assets/layouts/logo.png?url"
// 组件属性接口定义
interface Props {
  collapse?: boolean // 是否折叠状态，true=折叠（只显示图标），false=展开（显示图标+文字）
}

// 设置默认属性值
const props = withDefaults(defineProps<Props>(), {
  collapse: true
})

// 获取布局模式信息（左侧布局、顶部布局等）
const { isTop, isLeftTop } = useLayoutMode()

// 获取版本号（使用 Vite 定义的全局变量）
const version = __APP_VERSION__
</script>

<template>
  <!-- Logo 容器，根据折叠状态和布局模式应用不同的 CSS 类 -->
  <div class="layout-logo-container" :class="{ collapse: props.collapse, 'layout-mode-top': isTop, 'layout-mode-lefttop': isLeftTop }">
    <!-- 淡入淡出过渡动画 -->
    <transition name="layout-logo-fade">
      <!-- 折叠状态：只显示 logo 图标 -->
      <router-link v-if="props.collapse" key="collapse" to="/" class="logo-link">
        <img :src="logo" class="layout-logo" />
      </router-link>
      <!-- 展开状态：显示 logo 图标 + CiliKube 文字 + 版本号 -->
      <router-link v-else key="expand" to="/" class="logo-link logo-with-text">
        <img :src="logo" class="layout-logo-icon" />
        <div class="logo-text-container">
          <span class="logo-text">CiliKube</span>
          <span class="logo-version">v{{ version }}</span>
        </div>
      </router-link>
    </transition>
  </div>
</template>

<style lang="scss" scoped>
/* Logo 组件样式 */
.layout-logo-container {
  position: relative;
  width: 100%;
  height: var(--ck-header-height);
  /* 使用全局头部高度变量 */
  line-height: var(--ck-header-height);
  text-align: center;
  overflow: hidden;

  /* Logo 链接基础样式 */
  .logo-link {
    display: inline-flex;
    /* 使用 flex 布局让图标和文字水平排列 */
    align-items: center;
    /* 垂直居中对齐 */
    text-decoration: none;
    /* 移除链接下划线 */
    color: inherit;
    /* 继承父元素颜色 */
  }

  /* 折叠状态下的 logo（默认隐藏，在 .collapse 类中显示） */
  .layout-logo {
    display: none;
  }

  /* 展开状态：logo + 文字组合 */
  .logo-with-text {
    gap: 12px;
    /* 图标和文字容器之间的间距，可调整 */

    /* 展开状态下的 logo 图标样式 */
    .layout-logo-icon {
      width: 80px;
      /* logo 宽度，可调整 */
      height: 80px;
      /* logo 高度，可调整 */
    }

    /* 文字容器样式 */
    .logo-text-container {
      display: flex;
      flex-direction: column;
      align-items: flex-start;
      justify-content: center;
      line-height: 1;
    }

    /* CiliKube 文字样式 */
    .logo-text {
      font-size: 25px;
      /* 文字大小，可调整 */
      font-weight: 700;
      /* 文字粗细，可调整 */
      color: var(--el-text-color-primary);
      /* 文字颜色，使用 Element Plus 主题变量 */
      letter-spacing: 0.5px;
      /* 字母间距，可调整 */
      margin-bottom: 2px;
      /* 主标题和版本号之间的间距 */
    }

    /* 版本号样式 */
    .logo-version {
      font-size: 11px;
      /* 版本号文字大小 */
      font-weight: 500;
      /* 版本号文字粗细 */
      color: var(--el-text-color-regular);
      /* 版本号颜色，稍微淡一些 */
      letter-spacing: 0.3px;
      /* 版本号字母间距 */
      opacity: 0.8;
      /* 版本号透明度 */
    }
  }
}

/* 左侧模式下的特殊样式（纯左侧模式，不包括混合模式） */
.layout-logo-container:not(.layout-mode-top):not(.layout-mode-lefttop) {
  .logo-with-text {
    .logo-text {
      color: #ffffff;
      /* 左侧模式下文字为白色 */
    }
    
    .logo-version {
      color: rgba(255, 255, 255, 0.7);
      /* 左侧模式下版本号为半透明白色 */
    }
  }
}

/* 顶部布局模式样式 */
.layout-mode-top {
  height: var(--ck-navigationbar-height);
  /* 顶部导航栏高度 */
  line-height: var(--ck-navigationbar-height);
  
  /* 顶部模式下的logo和文字组合样式优化 */
  .logo-with-text {
    justify-content: flex-start;
    /* 左对齐 */
    padding-left: 16px;
    /* 左侧内边距 */
    gap: 8px;
    /* 顶部模式下图标和文字间距稍小 */
    
    .layout-logo-icon {
      width: 40px;
      /* 顶部模式下logo稍小一些 */
      height: 40px;
    }
    
    .logo-text-container {
      align-items: flex-start;
      /* 顶部模式下文字左对齐 */
    }
    
    .logo-text {
      font-size: 18px;
      /* 顶部模式下主文字稍小一些 */
      margin-bottom: 1px;
      /* 顶部模式下间距更紧凑 */
    }
    
    .logo-version {
      font-size: 10px;
      /* 顶部模式下版本号更小 */
    }
  }
}

/* 混合布局模式样式（左侧+顶部） */
.layout-mode-lefttop {
  height: var(--ck-header-height);
  /* 混合模式使用完整的header高度 */
  line-height: var(--ck-header-height);
  
  /* 混合模式下的logo和文字组合样式，适合完整的header高度 */
  .logo-with-text {
    justify-content: center;
    /* 居中对齐，确保有足够空间显示 */
    padding-left: 0;
    /* 移除左侧内边距，给文字更多空间 */
    
    .layout-logo-icon {
      width: 80px;
      /* 混合模式下logo与左侧模式一样大 */
      height: 80px;
    }
    
    .logo-text-container {
      align-items: flex-start;
      /* 混合模式下文字左对齐 */
    }
    
    .logo-text {
      font-size: 25px;
      /* 混合模式下文字与左侧模式一样大 */
      color: var(--el-text-color-primary);
      /* 混合模式下使用 Element Plus 主题变量 */
    }
    
    .logo-version {
      font-size: 11px;
      /* 混合模式下版本号大小 */
    }
  }
}

/* 折叠状态样式 */
.collapse {

  /* 折叠时显示小 logo */
  .layout-logo {
    width: 32px;
    /* 折叠状态 logo 宽度，可调整 */
    height: 32px;
    /* 折叠状态 logo 高度，可调整 */
    vertical-align: middle;
    display: inline-block;
  }

  /* 折叠时隐藏文字组合 */
  .logo-with-text {
    display: none;
  }
}

/* Logo 切换动画效果 */
.layout-logo-fade-enter-active,
.layout-logo-fade-leave-active {
  transition: opacity 0.3s ease;
  /* 0.3秒淡入淡出动画 */
}

.layout-logo-fade-enter-from,
.layout-logo-fade-leave-to {
  opacity: 0;
  /* 动画起始和结束透明度 */
}
</style>
