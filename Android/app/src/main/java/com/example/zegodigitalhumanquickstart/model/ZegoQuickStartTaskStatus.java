package com.example.zegodigitalhumanquickstart.model;

/**
 * 任务状态枚举
 */
public enum ZegoQuickStartTaskStatus {
    IDLE(0, "空闲"),
    RUNNING(1, "运行中"),
    STOPPED(2, "已停止"),
    ERROR(3, "错误");
    
    private final int value;
    private final String description;
    
    ZegoQuickStartTaskStatus(int value, String description) {
        this.value = value;
        this.description = description;
    }
    
    public int getValue() {
        return value;
    }
    
    public String getDescription() {
        return description;
    }
    
    public static ZegoQuickStartTaskStatus fromValue(int value) {
        for (ZegoQuickStartTaskStatus status : ZegoQuickStartTaskStatus.values()) {
            if (status.value == value) {
                return status;
            }
        }
        return IDLE;
    }
}

