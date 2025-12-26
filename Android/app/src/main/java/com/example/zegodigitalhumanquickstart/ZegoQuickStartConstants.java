package com.example.zegodigitalhumanquickstart;

/**
 * 应用常量定义类
 * 包含错误码、默认值、配置项等常量
 */
public class ZegoQuickStartConstants {
    
    // ==================== 错误码 ====================
    public static final int ERROR_CODE_SUCCESS = 0;
    public static final int ERROR_CODE_UNKNOWN = -1;
    public static final int ERROR_CODE_INVALID_PARAMETER = -2;
    public static final int ERROR_CODE_NETWORK_ERROR = -3;
    public static final int ERROR_CODE_PARSE_ERROR = -4;
    
    // ==================== 默认视频配置 ====================
    public static final String DEFAULT_VIDEO_CODEC = "H264";
    public static final int DEFAULT_VIDEO_WIDTH_SMALL = 320;
    public static final int DEFAULT_VIDEO_HEIGHT_SMALL = 400;
    public static final int DEFAULT_VIDEO_WIDTH_LARGE = 720;
    public static final int DEFAULT_VIDEO_HEIGHT_LARGE = 900;
    public static final int DEFAULT_VIDEO_BITRATE = 800;
    
    // ==================== 文本驱动配置 ====================
    public static final int MAX_TEXT_LENGTH = 1800;
    public static final int MIN_SPEECH_RATE = -500;
    public static final int MAX_SPEECH_RATE = 500;
    public static final int DEFAULT_SPEECH_RATE = 0;
    public static final int MIN_PITCH_RATE = -500;
    public static final int MAX_PITCH_RATE = 500;
    public static final int DEFAULT_PITCH_RATE = 0;
    public static final int MIN_VOLUME = 1;
    public static final int MAX_VOLUME = 100;
    public static final int DEFAULT_VOLUME = 50;
    
    // ==================== 任务房间ID前缀 ====================
    public static final String TASK_ROOM_PREFIX = "test_room_";
    public static final String TASK_STREAM_PREFIX = "stream_";
    public static final String TASK_PUBLISH_STREAM_PREFIX = "local_stream_";
    public static final String TASK_USER_PREFIX = "user_";
    
    // ==================== 请求超时时间 ====================
    public static final int NETWORK_TIMEOUT = 30; // 秒
    
    // ==================== 服务器配置 ====================
    // 示例:"http://192.168.88.213:3000/api"
    public static final String DEFAULT_SERVER_URL = 你的服务端地址;
    
    // 私有构造函数，防止实例化
    private ZegoQuickStartConstants() {
        throw new AssertionError("Cannot instantiate ZegoQuickStartConstants");
    }
}

