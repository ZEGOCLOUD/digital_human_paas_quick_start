const { VITE_API_BASE_URL, VITE_APP_SERVER } = import.meta.env;

// 应用配置
export const CONFIG = {
  API_BASE_URL: VITE_API_BASE_URL,
  DEFAULT_SERVER: VITE_APP_SERVER,
  DEFAULT_SIGNATURE_VERSION: "2.0", // 签名版本，按实际后端要求填写

  // 视频配置
  VIDEO_CONFIG: {
    codec: "H264",
    width: 0,
    height: 0,
    bitrate: 3000000,
  },

  // 布局配置
  LAYOUT_CONFIG: {
    MOBILE: {
      digitalHuman: {
        top: 20,
        left: 20,
        width: 280,
        height: 360,
        layer: 1,
      },
    },
    TABLET: {
      digitalHuman: {
        top: 40,
        left: 60,
        width: 0,
        height: 0,
        layer: 1,
      },
    },
    DESKTOP: {
      digitalHuman: {
        top: 50,
        left: 100,
        width: 0,
        height: 0,
        layer: 1,
      },
    },
    LARGE_DESKTOP: {
      digitalHuman: {
        top: 60,
        left: 120,
        width: 0,
        height: 0,
        layer: 1,
      },
    },
  },

  // TTS配置
  TTS_CONFIG: {
    speechRate: 0,
    pitchRate: 0,
    volume: 50,
  },
};
