package com.example.zegodigitalhumanquickstart.model;

import com.google.gson.annotations.SerializedName;
import java.io.Serializable;

/**
 * 数字人模型类
 */
public class ZegoQuickStartDigitalHuman implements Serializable {
    
    @SerializedName("DigitalHumanId")
    private String digitalHumanId;
    
    @SerializedName("Name")
    private String name;
    
    @SerializedName("AvatarUrl")
    private String avatarUrl;
    
    @SerializedName("PreviewUrl")
    private String previewUrl;
    
    @SerializedName("IsPublic")
    private boolean isPublic;
    
    @SerializedName("AppId")
    private long appId;
    
    @SerializedName("Token")
    private String token;

    @SerializedName("ExpireTime")
    private long expireTime;
    
    public ZegoQuickStartDigitalHuman() {
    }
    
    // Getters and Setters
    public String getDigitalHumanId() {
        return digitalHumanId != null ? digitalHumanId : "";
    }
    
    public void setDigitalHumanId(String digitalHumanId) {
        this.digitalHumanId = digitalHumanId;
    }
    
    public String getName() {
        return name != null ? name : "";
    }
    
    public void setName(String name) {
        this.name = name;
    }
    
    public String getAvatarUrl() {
        return avatarUrl;
    }

    public void setAvatarUrl(String avatarUrl) {
        this.avatarUrl = avatarUrl;
    }
    
    public String getPreviewUrl() {
        return previewUrl;
    }
    
    public void setPreviewUrl(String previewUrl) {
        this.previewUrl = previewUrl;
    }
    
    public boolean isPublic() {
        return isPublic;
    }
    
    public void setPublic(boolean isPublic) {
        this.isPublic = isPublic;
    }
    
    public long getAppId() {
        return appId;
    }
    
    public void setAppId(long appId) {
        this.appId = appId;
    }
    
    public String getToken() {
        return token != null ? token : "";
    }
    
    public void setToken(String token) {
        this.token = token;
    }
    
    public long getExpireTime() {
        return expireTime;
    }
    
    public void setExpireTime(long expireTime) {
        this.expireTime = expireTime;
    }
    
    @Override
    public String toString() {
        return "ZegoQuickStartDigitalHuman{" +
                "digitalHumanId='" + digitalHumanId + '\'' +
                ", name='" + name + '\'' +
                ", isPublic=" + isPublic +
                '}';
    }
}

