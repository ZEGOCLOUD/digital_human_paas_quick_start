<script setup>
import { computed } from "vue";

// 接收 props
const props = defineProps({
  show: Boolean,
  type: {
    type: String,
    default: "info",
    validator: (value) =>
      ["info", "success", "warning", "error"].includes(value),
  },
  message: String,
});

// 计算属性
const toastClass = computed(() => {
  return ["message-toast", `toast-${props.type}`, { show: props.show }];
});

const iconName = computed(() => {
  const icons = {
    info: "M12,2C6.48,2 2,6.48 2,12C2,17.52 6.48,22 12,22C17.52,22 22,17.52 22,12C22,6.48 17.52,2 12,2ZM13,17H11V15H13V17ZM13,13H11V7H13V13Z",
    success:
      "M12,2C6.48,2 2,6.48 2,12C2,17.52 6.48,22 12,22C17.52,22 22,17.52 22,12C22,6.48 17.52,2 12,2ZM10,17L5,12L6.41,10.59L10,14.17L17.59,6.58L19,8L10,17Z",
    warning: "M1,21H23L12,2L1,21ZM13,18H11V16H13V18ZM13,14H11V10H13V14Z",
    error:
      "M12,2C6.47,2 2,6.47 2,12C2,17.53 6.47,22 12,22C17.53,22 22,17.53 22,12C22,6.47 17.53,2 12,2ZM17,15.59L15.59,17L12,13.41L8.41,17L7,15.59L10.59,12L7,8.41L8.41,7L12,10.59L15.59,7L17,8.41L13.41,12L17,15.59Z",
  };
  return icons[props.type] || icons.info;
});

const typeConfig = computed(() => {
  const configs = {
    info: {
      color: "#1890ff",
      backgroundColor: "rgba(24, 144, 255, 0.1)",
      borderColor: "rgba(24, 144, 255, 0.3)",
    },
    success: {
      color: "#52c41a",
      backgroundColor: "rgba(82, 196, 26, 0.1)",
      borderColor: "rgba(82, 196, 26, 0.3)",
    },
    warning: {
      color: "#faad14",
      backgroundColor: "rgba(250, 173, 20, 0.1)",
      borderColor: "rgba(250, 173, 20, 0.3)",
    },
    error: {
      color: "#ff4d4f",
      backgroundColor: "rgba(255, 77, 79, 0.1)",
      borderColor: "rgba(255, 77, 79, 0.3)",
    },
  };
  return configs[props.type] || configs.info;
});
</script>

<template>
  <Transition name="toast">
    <div
      v-if="show"
      :class="toastClass"
      :style="{
        '--toast-color': typeConfig.color,
        '--toast-bg': typeConfig.backgroundColor,
        '--toast-border': typeConfig.borderColor,
      }"
    >
      <div class="toast-icon">
        <svg viewBox="0 0 24 24" fill="currentColor">
          <path :d="iconName" />
        </svg>
      </div>

      <div class="toast-content">
        <span class="toast-message">{{ message }}</span>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.message-toast {
  position: fixed;
  top: 80px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 2000;

  display: flex;
  align-items: center;
  gap: 12px;

  min-width: 300px;
  max-width: 500px;
  padding: 16px 20px;

  background: var(--toast-bg);
  backdrop-filter: blur(10px);
  border: 1px solid var(--toast-border);
  border-radius: 12px;

  color: var(--toast-color);
  font-size: 14px;
  font-weight: 500;

  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);

  /* 确保在移动端也居中 */
  margin: 0 20px;
}

.toast-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  flex-shrink: 0;
}

.toast-icon svg {
  width: 20px;
  height: 20px;
  color: var(--toast-color);
}

.toast-content {
  flex: 1;
  min-width: 0; /* 允许文本换行 */
}

.toast-message {
  line-height: 1.4;
  word-break: break-word;
}

/* 不同类型的样式 */
.toast-info {
  /* 使用 CSS 变量，在模板中动态设置 */
}

.toast-success {
  /* 使用 CSS 变量，在模板中动态设置 */
}

.toast-warning {
  /* 使用 CSS 变量，在模板中动态设置 */
}

.toast-error {
  /* 使用 CSS 变量，在模板中动态设置 */
}

/* 动画效果 */
.toast-enter-active {
  transition: all 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}

.toast-leave-active {
  transition: all 0.3s ease-in;
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(-50%) translateY(-30px) scale(0.9);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(-50%) translateY(-30px) scale(0.9);
}

/* 移动端适配 */
@media (max-width: 768px) {
  .message-toast {
    top: 70px;
    min-width: 280px;
    max-width: calc(100vw - 40px);
    padding: 14px 18px;
    font-size: 13px;
  }

  .toast-icon {
    width: 20px;
    height: 20px;
  }

  .toast-icon svg {
    width: 18px;
    height: 18px;
  }
}

/* 高对比度模式 */
@media (prefers-contrast: high) {
  .message-toast {
    border-width: 2px;
    font-weight: 600;
  }
}

/* 减少动画效果（为有运动敏感的用户） */
@media (prefers-reduced-motion: reduce) {
  .toast-enter-active,
  .toast-leave-active {
    transition: opacity 0.2s ease;
  }

  .toast-enter-from,
  .toast-leave-to {
    transform: translateX(-50%);
  }
}
</style>
