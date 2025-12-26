<script setup>
import { computed, ref } from "vue";

// 接收 props
const props = defineProps({
  appConfig: Object,
  taskState: Object,
  isMobile: Boolean,
  deviceType: String,
  showMobileControls: Boolean,
  showMessage: Function,
});

// 定义 emits
const emit = defineEmits(["toggle-mobile-controls"]);

const showStatusInfo = ref(false);

// 计算属性
const currentStatus = computed(() => {
  if (props.taskState.isStreaming) {
    return { text: "运行中", color: "#52c41a" };
  }
  return { text: "待机中,点击控制面板\"创建任务\"", color: "#d9d9d9" };
});

const headerClass = computed(() => {
  return `header-${props.deviceType}`;
});
</script>

<template>
  <header :class="headerClass">
    <div class="header-left">
      <div class="logo">
        <h1>🤖 数字人演示</h1>
      </div>
    </div>

    <div class="header-right">
      <!-- 状态信息按钮 -->
      <button
        class="status-btn"
        @click="showStatusInfo = true"
        :title="isMobile ? '状态信息' : '查看状态信息'"
      >
        <svg viewBox="0 0 24 24" fill="currentColor">
          <path d="M13,9H11V7H13M13,17H11V11H13M12,2A10,10 0 0,0 2,12A10,10 0 0,0 12,22A10,10 0 0,0 22,12A10,10 0 0,0 12,2Z"/>
        </svg>
        {{ isMobile ? "" : "状态" }}
      </button>

      <!-- 移动端控制面板按钮 -->
      <button
        v-if="isMobile"
        class="control-btn"
        @click="emit('toggle-mobile-controls')"
        :title="showMobileControls ? '收起控制面板' : '打开控制面板'"
        :class="{ active: showMobileControls }"
      >
        <svg v-if="!showMobileControls" viewBox="0 0 24 24" fill="currentColor">
          <!-- 菜单图标 -->
          <path d="M3 18h18v-2H3v2zm0-5h18v-2H3v2zm0-7v2h18V6H3z"/>
        </svg>
        <svg v-else viewBox="0 0 24 24" fill="currentColor">
          <!-- 关闭图标 -->
          <path d="M19,6.41L17.59,5L12,10.59L6.41,5L5,6.41L10.59,12L5,17.59L6.41,19L12,13.41L17.59,19L19,17.59L13.41,12L19,6.41Z"/>
        </svg>
      </button>
    </div>
  </header>

  <template v-if="showStatusInfo">
    <div class="status-overlay" @click="showStatusInfo = false"></div>
    <div class="status-info">
      <button class="modal-close-btn" @click="showStatusInfo = false" title="关闭" aria-label="关闭">
        <svg viewBox="0 0 24 24" fill="currentColor">
          <path d="M19,6.41L17.59,5L12,10.59L6.41,5L5,6.41L10.59,12L5,17.59L6.41,19L12,13.41L17.59,19L19,17.59L13.41,12L19,6.41Z"/>
        </svg>
      </button>
      <div class="status-item">
        <span class="status" :style="{ color: currentStatus.color }">
          {{ currentStatus.text }}
        </span>
      </div>

      <div class="status-item" v-if="taskState.currentTaskId">
        <span class="label">任务ID:</span>
        <span class="task-id">
          {{ taskState.currentTaskId }}
        </span>
      </div>

      <div class="status-item" v-if="taskState.currentTask && taskState.currentTask.RoomId">
        <span class="label">房间ID:</span>
        <span class="task-id">{{ taskState.currentTask.RoomId }}</span>
      </div>

      <div class="status-item" v-if="taskState.currentTask && taskState.currentTask.UserID">
        <span class="label">用户ID:</span>
        <span class="task-id">{{ taskState.currentTask.UserID }}</span>
      </div>

      <div class="status-item">
        <span class="label">AppID:</span>
        <span class="app-id">{{ appConfig.appId }}</span>
      </div>
    </div>
  </template>
</template>

<style scoped>
/* 通用头部样式 */
.header-mobile,
.header-tablet,
.header-desktop,
.header-large-desktop {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.2);
  transition: all 0.3s ease;
}

/* 移动端头部样式 */
.header-mobile {
  padding: 10px 16px;
  min-height: 60px;
}

/* 平板头部样式 */
.header-tablet {
  padding: 12px 20px;
  min-height: 65px;
}

/* 桌面端头部样式 */
.header-desktop {
  padding: 15px;
  min-height: 70px;
}

/* 大屏桌面头部样式 */
.header-large-desktop {
  padding: 18px 30px;
  min-height: 75px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
  max-width: calc(100vw - 90px);
}

.logo {
  min-width: 120px;
  max-width: 200px;
}

.logo h1 {
  margin: 0;
  color: white;
  font-size: 18px;
  font-weight: 600;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.status-info {
  position: fixed;
  top: 30%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 2000;
  display: flex;
  flex-direction: column;
  align-items: stretch;
  gap: 12px;
  min-width: 320px;
  width: min(420px, 90vw);
  max-height: 80vh;
  padding: 24px 24px 40px;
  background: rgba(0, 0, 0, 0.75);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  box-shadow: 0 12px 30px rgba(0, 0, 0, 0.45);
  overflow-y: auto;
}

.status-item {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 4px;
  font-size: 14px;
  width: 100%;
  min-width: 0;
  white-space: normal;
}

.status-item:first-child {
  width: 100%;
  min-width: 0;
}

.label {
  color: rgba(255, 255, 255, 0.8);
  font-weight: 500;
}

.status {
  font-weight: 600;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.task-id,
.app-id {
  color: rgba(255, 255, 255, 0.9);
  font-family: "Courier New", monospace;
  background: rgba(255, 255, 255, 0.1);
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 12px;
  width: 100%;
  display: block;
  word-break: break-all;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.status-btn,
.control-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.12);
  border: 1px solid rgba(255, 255, 255, 0.18);
  border-radius: 8px;
  color: white;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  backdrop-filter: blur(4px);
  min-height: 36px;
}

.status-btn:hover,
.control-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  border-color: rgba(255, 255, 255, 0.25);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.status-btn:active,
.control-btn:active {
  transform: translateY(0);
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);
}

.control-btn.active {
  background: rgba(255, 255, 255, 0.25);
  border-color: rgba(255, 255, 255, 0.3);
  box-shadow: 0 0 0 2px rgba(255, 255, 255, 0.1);
}

.status-btn svg,
.control-btn svg {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
}

/* 移动端适配 */
.header-mobile .logo h1 {
  font-size: 18px;
}

.header-mobile .header-left {
  gap: 10px;
}

.header-mobile .status-info {
  top: 72px;
  left: 50%;
  transform: translateX(-50%);
  width: calc(100vw - 24px);
  min-width: 0;
  max-width: 520px;
  max-height: 65vh;
  padding: 12px 14px;
  gap: 10px;
  border-radius: 10px;
  box-shadow: 0 10px 24px rgba(0, 0, 0, 0.35);
  -webkit-overflow-scrolling: touch;
}

.header-mobile .status-item {
  font-size: 12px;
  gap: 6px;
}

.header-mobile .task-id,
.header-mobile .app-id {
  font-size: 11px;
  padding: 1px 4px;
}

.header-mobile .header-right {
  gap: 6px;
}

.header-mobile .status-btn,
.header-mobile .control-btn {
  padding: 10px;
  min-width: 44px;
  min-height: 44px;
  font-size: 12px;
}

.header-mobile .status-btn svg,
.header-mobile .control-btn svg {
  width: 18px;
  height: 18px;
}

/* 平板适配 */
.header-tablet .logo h1 {
  font-size: 22px;
}

.header-tablet .status-info {
  gap: 15px;
}

/* 大屏桌面适配 */
.header-large-desktop .logo h1 {
  font-size: 26px;
}

.header-large-desktop .status-info {
  gap: 25px;
}

/* 动画效果 */
.status {
  position: relative;
}

.status::before {
  content: "";
  position: absolute;
  left: -12px;
  top: 50%;
  transform: translateY(-50%);
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: currentColor;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% {
    box-shadow: 0 0 0 0 currentColor;
  }
  70% {
    box-shadow: 0 0 0 4px transparent;
  }
  100% {
    box-shadow: 0 0 0 0 transparent;
  }
}

.status-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.45);
  backdrop-filter: blur(2px);
  z-index: 1500;
}

.modal-close-btn {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.2);
  color: #fff;
  border-radius: 6px;
  cursor: pointer;
}

.modal-close-btn:hover {
  background: rgba(255, 255, 255, 0.25);
}

.modal-close-btn svg {
  width: 18px;
  height: 18px;
}

.header-mobile .modal-close-btn {
  width: 32px;
  height: 32px;
  top: 6px;
  right: 6px;
}
</style>
