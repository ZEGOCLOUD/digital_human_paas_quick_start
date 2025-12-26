import { CONFIG } from "../config";

const apiRequest = async (path, data = {}, method = "POST", params = {}) => {
  // 清理路径，移除开头的 / 和 api/
  const cleanPath = path.replace(/^\//, "").replace(/^api\//, "");
  let baseUrl = CONFIG.API_BASE_URL.replace(/\/$/, "");
  
  // 检查 baseUrl 是否已经以 /api 结尾，如果是就不再加 api 前缀
  if (!baseUrl.endsWith("/api")) {
    baseUrl = `${baseUrl}/api`;
  }
  
  const url = `${baseUrl}/${cleanPath}`;
  const fullUrl = url + (method === "GET" ? "?" + new URLSearchParams(params).toString() : "");
  
  // 构建请求配置
  const requestHeaders = {
    "Content-Type": "application/json",
  };
  const requestBody = method === "POST" ? JSON.stringify(data) : null;
  
  // 打印完整的请求信息
  console.group(`🔵 [API请求] ${method} ${path}`);
  console.log("📋 完整URL:", fullUrl);
  console.log("🔧 请求方法:", method);
  console.log("📦 请求头:", requestHeaders);
  if (method === "GET" && Object.keys(params).length > 0) {
    console.log("🔍 查询参数:", params);
  }
  if (method === "POST" && data && Object.keys(data).length > 0) {
    console.log("📤 请求体:", data);
    console.log("📤 请求体(JSON字符串):", requestBody);
  }
  console.groupEnd();
  
  try {
    const response = await fetch(fullUrl, {
      method,
      headers: requestHeaders,
      body: requestBody,
    });
    
    // 打印响应信息
    console.group(`🟢 [API响应] ${method} ${path}`);
    console.log("📋 完整URL:", fullUrl);
    console.log("📊 状态码:", response.status);
    console.log("📝 状态文本:", response.statusText);
    console.log("📦 响应头:", Object.fromEntries(response.headers.entries()));
    
    const result = await response.json();
    console.log("📥 响应体:", result);
    console.groupEnd();

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return result;
  } catch (error) {
    console.group(`🔴 [API错误] ${method} ${path}`);
    console.error("📋 完整URL:", fullUrl);
    console.error("❌ 错误信息:", error);
    console.error("📝 错误详情:", {
      message: error.message,
      stack: error.stack,
    });
    console.groupEnd();
    throw error;
  }
}

export const post = (path, data = {}) => apiRequest(path, data, "POST");
export const get = (path, params = {}) => apiRequest(path, {}, "GET", params);

// 数字人资产相关 API
export const digitalHumanAPI = {
  async getDigitalHumanInfo(userId) {
    // 后端已要求传递用户ID，参考 iOS 端实现
    return post(
      "GetDigitalHumanInfo",
      {
        UserId: userId
      }
    );
  }
};

// 实时流相关API
export const streamAPI = {
   // 查询所有运行中的数字人视频流任务
   async queryStreamTasks() {
    return post(
      "QueryDigitalHumanStreamTasks",
      {}
    );
  },
  // 创建数字人视频流任务
  async createStreamTask(config) {
    return post("CreateDigitalHumanStreamTask", config);
  },

  // 停止数字人视频流任务
  async stopStreamTask(taskId) {
    return post("StopDigitalHumanStreamTask", { TaskId: taskId });
  },

  // 获取视频流任务状态
  async getStreamTaskStatus(taskId) {
    return post("GetDigitalHumanStreamTaskStatus", { TaskId: taskId });
  },
};

// 驱动相关API
export const driveAPI = {
  // 文本驱动数字人 - 只传递TaskId，所有参数在服务端设置
  async driveByText(taskId) {
    return post(
      "DriveByText",
      {
        TaskId: taskId
      }
    );
  },

  // 音频驱动数字人 - 只传递TaskId，所有参数在服务端设置
  async driveByAudio(taskId) {
    return post(
      "DriveByAudio",
      {
        TaskId: taskId
      }
    );
  },

  // WebSocket TTS驱动数字人 - 只传递TaskId，所有参数在服务端设置
  async driveByWsStreamWithTTS(taskId) {
    return post(
      "DriveByWsStreamWithTTS",
      {
        TaskId: taskId
      }
    );
  },

  // 打断驱动任务
  async interruptDriveTask(taskId) {
    return post(
      "InterruptDriveTask",
      {
        TaskId: taskId
      }
    );
  },
};
