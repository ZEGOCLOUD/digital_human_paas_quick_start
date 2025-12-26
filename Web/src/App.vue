<script setup>
import { ref, reactive, onMounted, onUnmounted, computed, provide, watch } from "vue";
import { ZegoExpressEngine } from 'zego-express-engine-webrtc';
import { device } from "./utils";
import { CONFIG } from "./config";
import { digitalHumanAPI } from "./utils/api";
import MessageToast from "./components/MessageToast.vue";
import HeaderConfig from "./components/HeaderConfig.vue";
import DigitalHumanDisplay from "./components/DigitalHumanDisplay.vue";
import ControlPanel from "./components/ControlPanel.vue";

// ZegoExpressEngine 延迟初始化
const zg = ref(null);
const currentAppId = ref(null);

// 初始化ZegoExpressEngine的函数
const initZegoEngine = (appId, server) => {
  // 如果已经初始化且appId相同，直接返回
  if (zg.value && currentAppId.value === appId) {
    return zg.value;
  }
  
  // 如果已经初始化但appId不同，需要先销毁再重新初始化
  if (zg.value && currentAppId.value !== appId) {
    console.log('[ZegoEngine] AppId已变更，重新初始化引擎');
    // 注意：ZegoExpressEngine可能没有destroy方法，这里只是记录日志
    // 实际使用时可能需要根据SDK文档处理
  }
  
  // 初始化新的引擎实例
  console.log('[ZegoEngine] 初始化引擎，AppId:', appId, 'Server:', server);
  zg.value = new ZegoExpressEngine(Number(appId), server);
  currentAppId.value = appId;
  
  return zg.value;
};

// provide zg实例和初始化函数
provide('zg', computed(() => zg.value));
provide('initZegoEngine', initZegoEngine);

// 响应式状态
const deviceType = ref(device.getDeviceType());
const userId = ref("");
const isMobile = computed(() => deviceType.value === "mobile");
const showMobileControls = ref(false); // 移动端默认隐藏控制面板

// 获取布局配置
const getLayoutConfig = (type) => {
  const configMap = {
    mobile: CONFIG.LAYOUT_CONFIG.MOBILE,
    tablet: CONFIG.LAYOUT_CONFIG.TABLET,
    desktop: CONFIG.LAYOUT_CONFIG.DESKTOP,
    "large-desktop": CONFIG.LAYOUT_CONFIG.LARGE_DESKTOP,
  };
  return configMap[type] || CONFIG.LAYOUT_CONFIG.DESKTOP;
};

// 表单数据
const settings = reactive({
  appId: "",
});

// 应用配置状态
const appConfig = reactive({
  appId: "",
  selectedDigitalHuman: null,
  videoConfig: { ...CONFIG.VIDEO_CONFIG },
  layoutConfig: getLayoutConfig(deviceType.value),
});

// 任务状态
const taskState = reactive({
  currentTask: null,
  taskStatus: null,
  driveList: [],
  isStreaming: false,
});

// 消息提示状态
const messageState = reactive({
  show: false,
  type: "info", // info, success, warning, error
  message: "",
});

// 计算属性
const layoutClass = computed(() => {
  return `layout-${deviceType.value}`;
});

// 监听窗口大小变化
const handleResize = () => {
  const newDeviceType = device.getDeviceType();
  if (newDeviceType !== deviceType.value) {
    deviceType.value = newDeviceType;
    appConfig.layoutConfig = getLayoutConfig(newDeviceType);
  }
};

// 显示消息提示
const showMessage = (type, message, duration = 3000) => {
  messageState.type = type;
  messageState.message = message;
  messageState.show = true;

  setTimeout(() => {
    messageState.show = false;
  }, duration);
};

// 获取/生成用户ID
const getCurrentUserId = () => {
  if (!userId.value) {
    const storedId = typeof localStorage !== 'undefined'
      ? localStorage.getItem('digital_human_user_id')
      : "";
    if (storedId) {
      userId.value = storedId;
    } else {
      userId.value = `user_${Math.random().toString(36).substr(2, 8)}`;
      if (typeof localStorage !== 'undefined') {
        localStorage.setItem('digital_human_user_id', userId.value);
      }
      console.log('[UserId] 生成新的 userId:', userId.value);
    }
  }
  return userId.value;
};

// 更新应用配置
const updateConfig = (newConfig) => {
  Object.assign(appConfig, newConfig);
};

// 更新表单数据
const updateSettings = (newSettings) => {
  Object.assign(settings, newSettings);
};

// 更新任务状态
const updateTaskState = (newState) => {
  Object.assign(taskState, newState);
};

const getDigitalHumanInfo = async () => {
  try {
    const currentUserId = getCurrentUserId();
    const res = await digitalHumanAPI.getDigitalHumanInfo(currentUserId);
    console.log("getDigitalHumanInfo 返回：", res);
    
    // 更新数字人信息到占位视图
    if (res?.Data) {
      const digitalHuman = res.Data;
      appConfig.selectedDigitalHuman = {
        DigitalHumanId: digitalHuman.DigitalHumanId || digitalHuman.digitalHumanId,
        Name: digitalHuman.Name || digitalHuman.name,
        AvatarUrl: digitalHuman.AvatarUrl || digitalHuman.avatarUrl || digitalHuman.CoverUrl || digitalHuman.coverUrl,
        IsPublic: digitalHuman.IsPublic || digitalHuman.isPublic || false
      };
      console.log("已更新数字人信息到占位视图:", appConfig.selectedDigitalHuman);
    }

    if (res?.Data?.AppId) {
      appConfig.appId = res.Data.AppId;
    }
  } catch (error) {
    console.error("获取数字人信息失败：", error);
  }
}


const init = async () => {
  getCurrentUserId();
  window.addEventListener("resize", handleResize);
  // 在程序启动时加载数字人信息
  await getDigitalHumanInfo();
}


// 组件挂载
onMounted(() => {
  init();
});

// 监听页面刷新/可见性变化，重新加载数字人信息
const handleVisibilityChange = () => {
  if (!document.hidden) {
    // 页面变为可见时，重新加载数字人信息
    getDigitalHumanInfo();
  }
};

// 监听页面可见性变化
if (typeof document !== 'undefined') {
  document.addEventListener('visibilitychange', handleVisibilityChange);
}

// 组件卸载时清理事件监听器
onUnmounted(() => {
  if (typeof document !== 'undefined') {
    document.removeEventListener('visibilitychange', handleVisibilityChange);
  }
});

watch(() => appConfig.appId, () => {
  // AppId 变化时的处理逻辑
})

watch(() => appConfig.selectedDigitalHuman, () => {
  // 数字人选择变化时的处理逻辑
})
</script>

<template>
  <div id="app" :class="layoutClass">
    <!-- 头部配置区 -->
    <HeaderConfig
      :appConfig="appConfig"
      :taskState="taskState"
      :isMobile="isMobile"
      :deviceType="deviceType"
      :showMessage="showMessage"
      :updateTaskState="updateTaskState"
      @toggle-mobile-controls="showMobileControls = !showMobileControls"
    />

    <!-- 主内容区 -->
    <main class="main-content">
      <!-- 数字人展示区 -->
      <section class="display-section">
        <DigitalHumanDisplay
          :appConfig="appConfig"
          :taskState="taskState"
          :isMobile="isMobile"
          :showMessage="showMessage"
          :deviceType="deviceType"
          :updateTaskState="updateTaskState"
        />
      </section>

      <!-- 控制面板 -->
      <section
        class="control-section"
        v-if="!isMobile || showMobileControls"
        :class="{ show: showMobileControls }"
      >
        <ControlPanel
          :appConfig="appConfig"
          :userId="userId"
          :taskState="taskState"
          :isMobile="isMobile"
          :showMessage="showMessage"
          :updateTaskState="updateTaskState"
        />
      </section>
    </main>

    <!-- 消息提示 -->
    <MessageToast
      :show="messageState.show"
      :type="messageState.type"
      :message="messageState.message"
    />
  </div>
</template>

<style scoped>
/* 全局样式重置 - 防止移动端滚动 */
:global(html), :global(body) {
  margin: 0;
  padding: 0;
  height: 100%;
  overflow: hidden; /* 防止页面滚动 */
}

:global(#app) {
  height: 100vh;
  overflow: hidden;
}

#app {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  font-family: "Segoe UI", Tahoma, Geneva, Verdana, sans-serif;
  position: relative;
  overflow-x: hidden;
}

/* 主内容区布局 */
.main-content {
  display: flex;
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  height: calc(100vh - 80px);
  justify-content: space-between;
  gap: 20px;
  padding: 0 20px;
}

.display-section {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.control-section {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 0;
}

/* 移动端布局 */
.layout-mobile {
  display: flex;
  flex-direction: column;
  height: 100vh; /* 改为固定高度 */
  overflow: hidden; /* 确保不产生滚动 */
  position: fixed; /* 固定定位，防止滚动 */
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}

.layout-mobile .main-content {
  display: flex;
  flex-direction: column;
  flex: 1;
  padding: 0;
  gap: 0;
  overflow: hidden;
  max-width: 100vw;
  height: calc(100vh - 80px); /* 减去头部高度 */
  min-height: 0; /* 允许收缩 */
}

.layout-mobile .display-section {
  flex: 1;
  padding: 0;
  min-height: 0; /* 允许收缩 */
  overflow: hidden; /* 防止内容溢出 */
}

.layout-mobile .control-section {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  max-height: 85vh;
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(20px);
  border-radius: 20px 20px 0 0;
  z-index: 100;
  transform: translateY(100%);
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow-y: auto;
  box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.15);
  padding: 20px;
}

.layout-mobile .control-section.show {
  transform: translateY(0);
}

/* 平板端布局 */
.layout-tablet {
  display: grid;
  grid-template-rows: auto 1fr;
  height: 100vh;
}

.layout-tablet .main-content {
  display: flex;
  flex-direction: column;
  padding: 16px;
  gap: 16px;
  overflow: hidden;
  max-width: 100vw;
  height: calc(100vh - 80px);
}

.layout-tablet .display-section {
  flex: 1;
  min-height: 0;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.layout-tablet .control-section {
  flex: 0 0 auto;
  max-height: 40vh;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  overflow-y: auto;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

/* 桌面端布局 */
.layout-desktop .main-content {
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  height: calc(100vh - 80px);
  justify-content: space-between;
  gap: 24px;
  padding: 20px 24px;
}

.layout-desktop .display-section {
  flex: 0 0 400px;
  max-width: 400px;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.layout-desktop .control-section {
  flex: 1;
  max-width: 600px;
}

/* 大桌面端布局 */
.layout-large-desktop .main-content {
  width: 100%;
  max-width: 1400px;
  margin: 0 auto;
  height: calc(100vh - 80px);
  justify-content: space-between;
  gap: 32px;
  padding: 20px 32px;
}

.layout-large-desktop .display-section {
  flex: 0 0 450px;
  max-width: 450px;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.layout-large-desktop .control-section {
  flex: 1;
  max-width: 700px;
}

/* 响应式断点 */
@media (max-width: 768px) {
  .main-content {
    flex-direction: column;
    width: 100vw;
    min-width: 0;
    height: auto;
    padding: 0;
    justify-content: flex-start;
    max-width: 100vw;
  }

  .display-section {
    width: 100vw;
    max-width: 100vw;
    padding: 0;
    margin: 0;
    flex: 1;
    min-height: 0;
  }

  .control-section {
    width: 100vw;
    max-width: 100vw;
    position: fixed;
    left: 0;
    right: 0;
    bottom: 0;
    border-radius: 20px 20px 0 0;
    background: rgba(255, 255, 255, 0.98);
    backdrop-filter: blur(20px);
    z-index: 100;
    box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.15);
    transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    transform: translateY(100%);
    max-height: 85vh;
    overflow-y: auto;
  }

  .control-section.show {
    transform: translateY(0);
  }

  /* 头部适配 */
  .header-mobile {
    position: sticky;
    top: 0;
    z-index: 101;
    background: inherit;
  }
}

@media (min-width: 769px) and (max-width: 1024px) {
  .main-content {
    width: 100%;
    max-width: 100vw;
    padding: 16px;
    gap: 16px;
  }

  .display-section {
    flex: 1;
    min-width: 0;
  }

  .control-section {
    flex: 1;
    min-width: 0;
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(10px);
    border-radius: 16px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  }
}

@media (min-width: 1025px) and (max-width: 1440px) {
  .main-content {
    width: 100%;
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 24px;
    gap: 24px;
  }

  .display-section {
    flex: 0 0 400px;
    max-width: 400px;
  }

  .control-section {
    flex: 1;
    max-width: 600px;
  }
}

@media (min-width: 1441px) {
  .main-content {
    width: 100%;
    max-width: 1400px;
    margin: 0 auto;
    padding: 0 32px;
    gap: 32px;
  }

  .display-section {
    flex: 0 0 450px;
    max-width: 450px;
  }

  .control-section {
    flex: 1;
    max-width: 700px;
  }
}

/* 滚动条样式 */
.control-section::-webkit-scrollbar {
  width: 6px;
}

.control-section::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.1);
  border-radius: 3px;
}

.control-section::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.3);
  border-radius: 3px;
}

.control-section::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 0, 0, 0.5);
}

/* 动画效果 */
.main-content {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.display-section,
.control-section {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/* 高分辨率屏幕优化 */
@media (-webkit-min-device-pixel-ratio: 2), (min-resolution: 192dpi) {
  .control-section {
    border-radius: 12px;
  }

  .layout-mobile .control-section {
    border-radius: 16px 16px 0 0;
  }
}

/* 减少动画偏好 */
@media (prefers-reduced-motion: reduce) {
  .main-content,
  .display-section,
  .control-section {
    transition: none;
  }

  .layout-mobile .control-section {
    transition: none;
  }
}
</style>
