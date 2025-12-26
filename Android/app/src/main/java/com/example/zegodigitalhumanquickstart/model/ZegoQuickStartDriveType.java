package com.example.zegodigitalhumanquickstart.model;

/**
 * 驱动类型枚举
 */
public enum ZegoQuickStartDriveType {
    TEXT(0, "文本驱动"),
    AUDIO(1, "音频驱动"),
    WS_TTS(2, "WebSocket TTS驱动");
    
    private final int value;
    private final String description;
    
    ZegoQuickStartDriveType(int value, String description) {
        this.value = value;
        this.description = description;
    }
    
    public int getValue() {
        return value;
    }
    
    public String getDescription() {
        return description;
    }
    
    public static ZegoQuickStartDriveType fromValue(int value) {
        for (ZegoQuickStartDriveType type : ZegoQuickStartDriveType.values()) {
            if (type.value == value) {
                return type;
            }
        }
        return TEXT;
    }
}

