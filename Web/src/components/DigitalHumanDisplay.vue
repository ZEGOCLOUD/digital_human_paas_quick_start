<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from "vue";

// 接收 props
const props = defineProps({
  appConfig: Object,
  taskState: Object,
  deviceType: String,
});

// 响应式状态
const isPlaying = ref(false);
const volume = ref(50);
const showControls = ref(false);
const windowWidth = ref(window.innerWidth);
const windowHeight = ref(window.innerHeight);

// 计算属性
const currentDigitalHuman = computed(() => {
  if (props.appConfig.selectedDigitalHuman) {
    return props.appConfig.selectedDigitalHuman;
  } else {
    return {
      name: "请选择数字人形象",
      avatarUrl: "",
    };
  }
});

const displayClass = computed(() => {
  return `display-${props.deviceType}`;
});

// 根据设备类型计算最佳宽高比
const aspectRatio = computed(() => {
  switch (props.deviceType) {
    case 'mobile':
      return '9/16'; // 移动端保持9:16
    case 'tablet':
      return '9/16'; // 平板端也使用9:16，保持一致性
    case 'desktop':
    case 'large-desktop':
      return '9/16'; // 桌面端保持9:16
    default:
      return '9/16';
  }
});

// 计算容器最大宽度 - 优化为尽可能撑满
const maxWidth = computed(() => {
  switch (props.deviceType) {
    case 'mobile':
      return '100vw'; // 移动端全屏
    case 'tablet':
      return '100%'; // 平板端撑满可用空间
    case 'desktop':
      return '100%'; // 桌面端撑满可用空间
    case 'large-desktop':
      return '100%'; // 大桌面端撑满可用空间
    default:
      return '100%';
  }
});

// 计算容器最大高度
const maxHeight = computed(() => {
  switch (props.deviceType) {
    case 'mobile':
      return '100vh'; // 移动端全屏
    case 'tablet':
      return '100%'; // 平板端撑满可用空间
    case 'desktop':
      return '100%'; // 桌面端撑满可用空间
    case 'large-desktop':
      return '100%'; // 大桌面端撑满可用空间
    default:
      return '100%';
  }
});

// 监听窗口尺寸变化
const handleResize = () => {
  windowWidth.value = window.innerWidth;
  windowHeight.value = window.innerHeight;
};

// 组件挂载和卸载
onMounted(() => {
  window.addEventListener('resize', handleResize);
  handleResize(); // 初始化尺寸

  // 延迟检查，确保DOM完全渲染
  setTimeout(() => {
    // 调试信息：确认视频容器存在
    const remoteVideoContainer = document.getElementById('remote-video');
    if (remoteVideoContainer) {
      console.log('✅ 视频容器已找到:', remoteVideoContainer);
      console.log('✅ 容器尺寸:', {
        width: remoteVideoContainer.offsetWidth,
        height: remoteVideoContainer.offsetHeight,
        style: remoteVideoContainer.style.cssText
      });

      // 检查父容器尺寸
      const videoContainer = remoteVideoContainer.closest('.video-container');
      const videoWrapper = remoteVideoContainer.closest('.video-wrapper');
      if (videoContainer) {
        console.log('📦 父容器尺寸:', {
          width: videoContainer.offsetWidth,
          height: videoContainer.offsetHeight,
          style: videoContainer.style.cssText
        });
      }
      if (videoWrapper) {
        console.log('📦 包装器尺寸:', {
          width: videoWrapper.offsetWidth,
          height: videoWrapper.offsetHeight,
          style: videoWrapper.style.cssText
        });
      }

      // 如果容器尺寸为0，尝试强制设置
      if (remoteVideoContainer.offsetWidth === 0 || remoteVideoContainer.offsetHeight === 0) {
        console.warn('⚠️ 容器尺寸为0，尝试强制设置');
        forceContainerSize(remoteVideoContainer);
      }
    } else {
      console.warn('⚠️ 视频容器未找到');
    }
  }, 100);
});

// 监听任务状态变化
const watchTaskState = () => {
  if (props.taskState && props.taskState.isStreaming) {
    console.log('🎥 视频流状态: 正在直播');
    // 确保视频容器可见
    const remoteVideoContainer = document.getElementById('remote-video');
    if (remoteVideoContainer) {
      remoteVideoContainer.style.display = 'flex';
      console.log('✅ 视频容器已设置为可见');

      // 强制设置容器尺寸
      forceContainerSize(remoteVideoContainer);

      // 检查视频元素
      setTimeout(() => {
        const videoElements = remoteVideoContainer.querySelectorAll('canvas, video');
        console.log('📹 找到视频元素数量:', videoElements.length);
        videoElements.forEach((el, index) => {
          console.log(`📹 视频元素 ${index + 1}:`, {
            tagName: el.tagName,
            width: el.offsetWidth,
            height: el.offsetHeight,
            naturalWidth: el.naturalWidth || el.width,
            naturalHeight: el.naturalHeight || el.height,
            objectFit: getComputedStyle(el).objectFit,
            aspectRatio: (el.offsetWidth / el.offsetHeight).toFixed(2)
          });
        });
      }, 500);

      // 再次检查尺寸
      setTimeout(() => {
        console.log('📏 直播后容器尺寸:', {
          width: remoteVideoContainer.offsetWidth,
          height: remoteVideoContainer.offsetHeight
        });
      }, 100);
    }
  }
};

// 强制设置容器尺寸
const forceContainerSize = (container) => {
  if (!container) return;

  // 获取父容器尺寸
  const videoContainer = container.closest('.video-container');
  const videoWrapper = container.closest('.video-wrapper');

  if (videoContainer && videoWrapper) {
    const containerWidth = videoContainer.offsetWidth;
    const containerHeight = videoContainer.offsetHeight;

    console.log('🔧 强制设置容器尺寸:', { containerWidth, containerHeight });

    // 确保容器有尺寸
    if (containerWidth > 0 && containerHeight > 0) {
      container.style.width = `${containerWidth}px`;
      container.style.height = `${containerHeight}px`;
      console.log('✅ 容器尺寸已强制设置');
    } else {
      console.warn('⚠️ 父容器尺寸为0，无法设置');
    }
  }
};

// 监听props变化
watch(() => props.taskState?.isStreaming, (newVal, oldVal) => {
  if (newVal !== oldVal) {
    watchTaskState();
  }
});

onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
  if (zg) {
    try {
      if (isRoomLoggedIn && props.taskState && props.taskState.currentTask) {
        zg.logoutRoom(props.taskState.currentTask.RoomId);
      }
      zg.destroyEngine();
    } catch {}
    zg = null;
    isRoomLoggedIn = false;
  }
});

// 处理音量变化
const handleVolumeChange = (event) => {
  volume.value = event.target.value;
  // 音量控制需要通过Zego引擎来实现
  // 这里可以添加对Zego引擎的音量控制
  console.log('音量设置为:', volume.value);
};

// 切换播放状态
const togglePlay = () => {
  isPlaying.value = !isPlaying.value;
  // 播放/暂停控制需要通过Zego引擎来实现
  // 这里可以添加对Zego引擎的播放控制
  console.log('播放状态切换为:', isPlaying.value ? '播放' : '暂停');
};

// 显示/隐藏控制栏
const toggleControls = () => {
  showControls.value = !showControls.value;
};

let zg = null;
let isRoomLoggedIn = false;
</script>

<template>
  <div :class="displayClass" :style="{ '--aspect-ratio': aspectRatio, '--max-width': maxWidth, '--max-height': maxHeight }">
    <!-- 视频显示区域 -->
    <div class="video-container">
      <div class="video-wrapper">
        <!-- 远端流播放容器 - Zego引擎会将视频流渲染到这里 -->
        <div
          id="remote-video"
          class="remote-video-container"
          :class="{ 'streaming': taskState.isStreaming }"
        ></div>

        <!-- 预览覆盖层 -->
        <div class="preview-overlay" v-if="!taskState.isStreaming">
          <!-- 全屏头像 -->
          <div class="digital-human-avatar-fullscreen">
            <div v-if="!currentDigitalHuman.AvatarUrl" class="avatar-placeholder-fullscreen">
              <svg viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z"/>
              </svg>
            </div>
            <div v-else class="avatar-img-fullscreen">
              <img :src="currentDigitalHuman.AvatarUrl" alt="AvatarUrl" />
            </div>
          </div>
          <!-- 文字信息 -->
          <div class="preview-text">
            <h3>{{ currentDigitalHuman.Name || currentDigitalHuman.name }}</h3>
            <p v-if="!taskState.currentTask">点击"创建任务"开始体验数字人</p>
            <p v-else>准备中，请稍候...</p>
          </div>
        </div>

        <!-- 流状态指示器 -->
        <div class="stream-indicator" v-if="taskState.isStreaming">
          <div class="indicator-dot"></div>
          <span>直播中</span>
        </div>

        <!-- 视频控制栏 -->
        <div
          class="video-controls"
          v-show="showControls && taskState.isStreaming"
          @click.stop
        >
          <button class="play-btn" @click="togglePlay">
            <svg v-if="!isPlaying" viewBox="0 0 24 24" fill="currentColor">
              <path d="M8 5v14l11-7z" />
            </svg>
            <svg v-else viewBox="0 0 24 24" fill="currentColor">
              <path d="M6 19h4V5H6v14zm8-14v14h4V5h-4z" />
            </svg>
          </button>

          <div class="volume-control">
            <svg viewBox="0 0 24 24" fill="currentColor">
              <path d="M3 9v6h4l5 5V4L7 9H3zm13.5 3c0-1.77-1.02-3.29-2.5-4.03v8.05c1.48-.73 2.5-2.25 2.5-4.02zM14 3.23v2.06c2.89.86 5 3.54 5 6.71s-2.11 5.85-5 6.71v2.06c4.01-.91 7-4.49 7-8.77s-2.99-7.86-7-8.77z"/>
            </svg>
            <input
              type="range"
              min="0"
              max="100"
              :value="volume"
              @input="handleVolumeChange"
              class="volume-slider"
            />
          </div>
        </div>

        <!-- 点击显示控制栏的透明层 -->
        <div
          class="controls-trigger"
          v-if="taskState.isStreaming && !showControls"
          @click="toggleControls"
        ></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* CSS变量定义 */
:root {
  --primary-color: #667eea;
  --secondary-color: #764ba2;
  --background-dark: #020202;
  --background-light: #f8f9fa;
  --text-light: #ffffff;
  --text-dark: #333333;
  --border-radius: 16px;
  --shadow-light: 0 4px 20px rgba(0, 0, 0, 0.1);
  --shadow-heavy: 0 8px 32px rgba(0, 0, 0, 0.3);
  --transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/* 通用显示容器样式 */
.display-mobile,
.display-tablet,
.display-desktop,
.display-large-desktop {
  display: flex;
  flex-direction: column;
  height: 100%;
  justify-content: center;
  align-items: center;
  overflow: hidden;
  background: linear-gradient(135deg, var(--primary-color) 0%, var(--secondary-color) 100%);
  border-radius: var(--border-radius);
  position: relative;
  box-shadow: var(--shadow-heavy);
  /* 撑满可用空间 */
  width: 100%;
  max-width: var(--max-width);
  max-height: var(--max-height);
}

/* 移动端特殊样式 */
.display-mobile {
  flex: 1;
  border-radius: 0;
  background: var(--background-dark);
  width: 100vw;
  height: 100%; /* 改为100%，不固定为100vh */
  max-width: 100vw;
  max-height: 100%; /* 改为100% */
  overflow: hidden; /* 确保不产生滚动 */
  position: relative; /* 改为相对定位 */
}

.display-mobile .video-container {
  width: 100%;
  height: 100%;
  max-width: 100%;
  max-height: 100%;
  min-height: 0; /* 移除最小高度限制 */
  overflow: hidden; /* 防止内容溢出 */
}

.display-mobile .video-wrapper {
  border-radius: 0;
  width: 100%;
  height: 100%;
  aspect-ratio: auto; /* 移动端不使用固定比例，全屏显示 */
  overflow: hidden; /* 防止内容溢出 */
}

/* 平板端样式 */
.display-tablet .video-container {
  width: 100%;
  height: 100%;
  max-width: 100%;
  max-height: 100%;
  min-height: 500px;
  padding: 0;
}

/* 桌面端样式 */
.display-desktop .video-container,
.display-large-desktop .video-container {
  width: 100%;
  height: 100%;
  max-width: 100%;
  max-height: 100%;
  min-height: 600px;
  padding: 0;
}

/* 视频容器 */
.video-container {
  position: relative;
  width: 100%;
  max-width: var(--max-width);
  max-height: var(--max-height);
  flex-shrink: 0;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  /* 设置明确的高度 */
  height: 100%;
  min-height: 400px;
}

/* 视频包装器 */
.video-wrapper {
  position: relative;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--border-radius);
  overflow: hidden;
  background: var(--background-dark);
  box-shadow: var(--shadow-heavy);
  /* 确保有最小尺寸 */
  min-width: 200px;
  min-height: 300px;
  /* 使用aspect-ratio保持比例 */
  aspect-ratio: var(--aspect-ratio);
}

/* 远端视频容器 */
.remote-video-container {
  width: 100% !important;
  height: 100% !important;
  position: absolute;
  top: 0;
  left: 0;
  z-index: 2;
  background: var(--background-dark);
  border-radius: var(--border-radius);
  transition: var(--transition);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  min-width: 0;
  min-height: 0;
}

.remote-video-container.streaming {
  background: transparent;
}

/* 确保Zego视频流正确显示 */
.remote-video-container > * {
  width: 100% !important;
  height: 100% !important;
  object-fit: contain; /* 保持视频比例 */
  border-radius: var(--border-radius);
  display: block;
  max-width: 100%;
  max-height: 100%;
}

/* 确保Zego创建的canvas或video元素正确显示 */
.remote-video-container canvas,
.remote-video-container video {
  width: 100% !important;
  height: 100% !important;
  object-fit: contain; /* 保持视频比例 */
  border-radius: var(--border-radius);
  display: block;
  max-width: 100%;
  max-height: 100%;
}

/* 移动端特殊处理 - 使用cover填充 */
.display-mobile .remote-video-container > *,
.display-mobile .remote-video-container canvas,
.display-mobile .remote-video-container video {
  object-fit: cover; /* 移动端使用cover填充 */
}

/* 预览覆盖层 */
.preview-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  align-items: center;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.9) 0%, rgba(118, 75, 162, 0.9) 100%);
  backdrop-filter: blur(10px);
  color: var(--text-light);
  text-align: center;
  z-index: 3;
  transition: var(--transition);
}

/* 全屏数字人头像 */
.digital-human-avatar-fullscreen {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
}

.avatar-placeholder-fullscreen {
  width: 200px;
  height: 200px;
  color: rgba(255, 255, 255, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar-placeholder-fullscreen svg {
  width: 100%;
  height: 100%;
}

.avatar-img-fullscreen {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar-img-fullscreen img {
  width: 100%;
  height: 100%;
  object-fit: contain; /* aspect_fit方式 */
}

/* 预览文字 */
.preview-text {
  position: relative;
  z-index: 4;
  padding: 20px;
  background: linear-gradient(to top, rgba(0, 0, 0, 0.7) 0%, transparent 100%);
  width: 100%;
  margin-bottom: 0;
  color: #ffffff !important;
}

.preview-text h3 {
  margin: 0 0 12px 0;
  font-size: clamp(20px, 4vw, 28px);
  font-weight: 600;
  line-height: 1.2;
  color: #ffffff !important;
}

.preview-text p {
  margin: 0;
  font-size: clamp(14px, 3vw, 16px);
  opacity: 0.9;
  line-height: 1.5;
  color: #ffffff !important;
}

/* 流状态指示器 */
.stream-indicator {
  position: absolute;
  top: 16px;
  right: 16px;
  display: flex;
  align-items: center;
  gap: 8px;
  background: rgba(0, 0, 0, 0.7);
  color: var(--text-light);
  padding: 8px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
  z-index: 4;
  backdrop-filter: blur(10px);
}

.indicator-dot {
  width: 8px;
  height: 8px;
  background: #ff4757;
  border-radius: 50%;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

/* 视频控制栏 */
.video-controls {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: linear-gradient(transparent, rgba(0, 0, 0, 0.8));
  color: var(--text-light);
  z-index: 4;
  transition: var(--transition);
}

.play-btn {
  background: rgba(255, 255, 255, 0.2);
  border: none;
  border-radius: 50%;
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: var(--transition);
  backdrop-filter: blur(10px);
}

.play-btn:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: scale(1.05);
}

.play-btn svg {
  width: 24px;
  height: 24px;
  fill: var(--text-light);
}

.volume-control {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
}

.volume-control svg {
  width: 20px;
  height: 20px;
  fill: var(--text-light);
  opacity: 0.8;
}

.volume-slider {
  flex: 1;
  height: 4px;
  background: rgba(255, 255, 255, 0.3);
  border-radius: 2px;
  outline: none;
  appearance: none;
  -webkit-appearance: none;
  cursor: pointer;
  transition: var(--transition);
}

.volume-slider:hover {
  background: rgba(255, 255, 255, 0.4);
}

.volume-slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: var(--text-light);
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
  transition: var(--transition);
}

.volume-slider::-webkit-slider-thumb:hover {
  transform: scale(1.2);
}

/* 控制栏触发器 */
.controls-trigger {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 60px;
  z-index: 3;
  cursor: pointer;
}

/* 响应式适配 */
@media (max-width: 768px) {
  .display-mobile .video-container {
    width: 100%;
    height: 100%;
    max-width: 100%;
    max-height: 100%;
    padding: 0;
    overflow: hidden;
  }

  .display-mobile .video-wrapper {
    border-radius: 0;
    width: 100%;
    height: 100%;
    aspect-ratio: auto; /* 移动端不使用固定比例 */
    overflow: hidden;
  }

  .preview-text {
    padding: 16px;
    color: #ffffff !important;
  }

  .avatar-placeholder-fullscreen {
    width: 150px;
    height: 150px;
  }

  .preview-text h3 {
    font-size: 20px;
    color: #ffffff !important;
  }

  .preview-text p {
    font-size: 14px;
    color: #ffffff !important;
  }

  .video-controls {
    padding: 16px;
  }

  .play-btn {
    width: 40px;
    height: 40px;
  }

  .play-btn svg {
    width: 20px;
    height: 20px;
  }

  .volume-slider {
    max-width: 80px;
  }

  .stream-indicator {
    top: 12px;
    right: 12px;
    padding: 6px 10px;
    font-size: 11px;
  }
}

@media (min-width: 769px) and (max-width: 1024px) {
  .display-tablet .video-container {
    padding: 0;
    width: 100%;
    height: 100%;
  }

  .display-tablet .video-wrapper {
    aspect-ratio: var(--aspect-ratio);
    max-width: 100%;
    max-height: 100%;
  }

  .preview-text {
    padding: 20px;
    color: #ffffff !important;
  }

  .preview-text h3 {
    color: #ffffff !important;
  }

  .preview-text p {
    color: #ffffff !important;
  }

  .avatar-placeholder-fullscreen {
    width: 180px;
    height: 180px;
  }
}

@media (min-width: 1025px) {
  .display-desktop .video-container,
  .display-large-desktop .video-container {
    padding: 0;
    width: 100%;
    height: 100%;
  }

  .display-desktop .video-wrapper,
  .display-large-desktop .video-wrapper {
    aspect-ratio: var(--aspect-ratio);
    max-width: 100%;
    max-height: 100%;
  }

  .preview-content {
    padding: 32px;
    gap: 24px;
  }

  .digital-human-avatar {
    width: 120px;
    height: 120px;
  }

  .avatar-placeholder {
    width: 60px;
    height: 60px;
  }
}

/* 动画效果 */
.preview-overlay {
  animation: fadeIn 0.5s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.video-wrapper {
  transition: var(--transition);
}

.remote-video-container {
  transition: var(--transition);
}

/* 高分辨率屏幕优化 */
@media (-webkit-min-device-pixel-ratio: 2), (min-resolution: 192dpi) {
  .video-wrapper {
    border-radius: calc(var(--border-radius) * 0.5);
  }

  .play-btn {
    border-radius: 50%;
  }
}

/* 减少动画偏好 */
@media (prefers-reduced-motion: reduce) {
  * {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
  }
}
</style>
