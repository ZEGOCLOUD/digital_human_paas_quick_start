package com.example.zegodigitalhumanquickstart.network;

/**
 * API接口常量定义类
 * 定义所有API路径
 */
public class ZegoQuickStartAPIConstants {
    
    // ==================== API路径 ====================
    
    // 数字人相关API
    public static final String ACTION_GET_DIGITAL_HUMAN_INFO = "/GetDigitalHumanInfo";
    
    // 任务管理API
    public static final String ACTION_CREATE_DIGITAL_HUMAN_STREAM_TASK = "/CreateDigitalHumanStreamTask";
    public static final String ACTION_STOP_DIGITAL_HUMAN_STREAM_TASK = "/StopDigitalHumanStreamTask";
    public static final String ACTION_QUERY_DIGITAL_HUMAN_STREAM_TASKS = "/QueryDigitalHumanStreamTasks";
    
    // 驱动相关API
    public static final String ACTION_DRIVE_BY_TEXT = "/DriveByText";
    public static final String ACTION_DRIVE_BY_AUDIO = "/DriveByAudio";
    public static final String ACTION_DRIVE_BY_WS_STREAM_WITH_TTS = "/DriveByWsStreamWithTTS";
    public static final String ACTION_INTERRUPT_DRIVE_TASK = "/InterruptDriveTask";
    
    // ==================== HTTP请求头 ====================
    public static final String HEADER_CONTENT_TYPE = "Content-Type";
    public static final String CONTENT_TYPE_JSON = "application/json";
    
    // 私有构造函数，防止实例化
    private ZegoQuickStartAPIConstants() {
        throw new AssertionError("Cannot instantiate ZegoQuickStartAPIConstants");
    }
}

