package com.example.zegodigitalhumanquickstart.model;

import com.example.zegodigitalhumanquickstart.ZegoQuickStartConstants;
import java.io.Serializable;

/**
 * 应用配置模型类
 */
public class ZegoQuickStartConfig implements Serializable {
    
    private String serverURL;           // 服务器URL
    private ZegoQuickStartOutputMode outputMode;      // 输出模式
    private ZegoQuickStartVideoConfig videoConfig;    // 视频配置
    
    public ZegoQuickStartConfig() {
        // 使用默认值初始化
        this.serverURL = ZegoQuickStartConstants.DEFAULT_SERVER_URL;
        this.outputMode = ZegoQuickStartOutputMode.SMALL;
        this.videoConfig = new ZegoQuickStartVideoConfig();
    }
    
    // Getters and Setters
    public String getServerURL() {
        return serverURL != null ? serverURL : ZegoQuickStartConstants.DEFAULT_SERVER_URL;
    }
    
    public void setServerURL(String serverURL) {
        this.serverURL = serverURL;
    }
    
    public ZegoQuickStartOutputMode getOutputMode() {
        return outputMode != null ? outputMode : ZegoQuickStartOutputMode.SMALL;
    }
    
    public void setOutputMode(ZegoQuickStartOutputMode outputMode) {
        this.outputMode = outputMode;
    }
    
    public ZegoQuickStartVideoConfig getVideoConfig() {
        if (videoConfig == null) {
            videoConfig = new ZegoQuickStartVideoConfig();
        }
        return videoConfig;
    }
    
    public void setVideoConfig(ZegoQuickStartVideoConfig videoConfig) {
        this.videoConfig = videoConfig;
    }
    
    /**
     * 根据输出模式更新视频配置
     */
    public void updateVideoConfigByOutputMode() {
        if (videoConfig == null) {
            videoConfig = new ZegoQuickStartVideoConfig();
        }
        
        if (outputMode == ZegoQuickStartOutputMode.LARGE) {
            videoConfig.setWidth(ZegoQuickStartConstants.DEFAULT_VIDEO_WIDTH_LARGE);
            videoConfig.setHeight(ZegoQuickStartConstants.DEFAULT_VIDEO_HEIGHT_LARGE);
        } else {
            videoConfig.setWidth(ZegoQuickStartConstants.DEFAULT_VIDEO_WIDTH_SMALL);
            videoConfig.setHeight(ZegoQuickStartConstants.DEFAULT_VIDEO_HEIGHT_SMALL);
        }
    }
    
    @Override
    public String toString() {
        return "ZegoQuickStartConfig{" +
                "serverURL='" + serverURL + '\'' +
                ", outputMode=" + outputMode +
                ", videoConfig=" + videoConfig +
                '}';
    }
}
