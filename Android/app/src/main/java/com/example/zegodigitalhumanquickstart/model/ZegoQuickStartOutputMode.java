package com.example.zegodigitalhumanquickstart.model;

/**
 * 输出模式枚举
 */
public enum ZegoQuickStartOutputMode {
    SMALL(0, "小图模式"),
    LARGE(1, "大图模式");
    
    private final int value;
    private final String description;
    
    ZegoQuickStartOutputMode(int value, String description) {
        this.value = value;
        this.description = description;
    }
    
    public int getValue() {
        return value;
    }
    
    public String getDescription() {
        return description;
    }
    
    public static ZegoQuickStartOutputMode fromValue(int value) {
        for (ZegoQuickStartOutputMode mode : ZegoQuickStartOutputMode.values()) {
            if (mode.value == value) {
                return mode;
            }
        }
        return SMALL;
    }
}

