package com.example.zegodigitalhumanquickstart.model;

import com.google.gson.JsonObject;
import java.io.Serializable;

/**
 * 视频配置模型类
 */
public class ZegoQuickStartVideoConfig implements Serializable {
    
    private String codec;      // 编码格式（如"H264"）
    private int width;         // 宽度
    private int height;        // 高度
    private int bitrate;       // 码率
    
    public ZegoQuickStartVideoConfig() {
        // 默认小图模式配置
        this.codec = "H264";
        this.width = 320;
        this.height = 400;
        this.bitrate = 800;
    }
    
    public ZegoQuickStartVideoConfig(String codec, int width, int height, int bitrate) {
        this.codec = codec;
        this.width = width;
        this.height = height;
        this.bitrate = bitrate;
    }
    
    // Getters and Setters
    public String getCodec() {
        return codec != null ? codec : "H264";
    }
    
    public void setCodec(String codec) {
        this.codec = codec;
    }
    
    public int getWidth() {
        return width;
    }
    
    public void setWidth(int width) {
        this.width = width;
    }
    
    public int getHeight() {
        return height;
    }
    
    public void setHeight(int height) {
        this.height = height;
    }
    
    public int getBitrate() {
        return bitrate;
    }
    
    public void setBitrate(int bitrate) {
        this.bitrate = bitrate;
    }
    
    /**
     * 转换为JSON对象（用于API请求）
     */
    public JsonObject toJsonObject() {
        JsonObject json = new JsonObject();
        json.addProperty("Codec", getCodec());
        json.addProperty("Width", width);
        json.addProperty("Height", height);
        json.addProperty("Bitrate", bitrate);
        return json;
    }
    
    @Override
    public String toString() {
        return "ZegoQuickStartVideoConfig{" +
                "codec='" + codec + '\'' +
                ", width=" + width +
                ", height=" + height +
                ", bitrate=" + bitrate +
                '}';
    }
}

