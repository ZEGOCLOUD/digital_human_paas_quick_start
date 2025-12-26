package com.example.zegodigitalhumanquickstart.model;

import com.google.gson.annotations.SerializedName;
import java.io.Serializable;

/**
 * 任务模型类
 */
public class ZegoQuickStartTask implements Serializable {
    
    @SerializedName("TaskId")
    private String taskId;
    
    @SerializedName("RoomId")
    private String roomId;
    
    @SerializedName("StreamId")
    private String streamId;
    
    @SerializedName("UserId")
    private String userId;
    
    @SerializedName("AppId")
    private long appId;
    
    private ZegoQuickStartTaskStatus status;  // 任务状态
    
    public ZegoQuickStartTask() {
        this.status = ZegoQuickStartTaskStatus.IDLE;
    }
    
    // Getters and Setters
    public String getTaskId() {
        return taskId != null ? taskId : "";
    }
    
    public void setTaskId(String taskId) {
        this.taskId = taskId;
    }
    
    public String getRoomId() {
        return roomId != null ? roomId : "";
    }
    
    public void setRoomId(String roomId) {
        this.roomId = roomId;
    }
    
    public String getStreamId() {
        return streamId != null ? streamId : "";
    }
    
    public void setStreamId(String streamId) {
        this.streamId = streamId;
    }
    
    public String getUserId() {
        return userId != null ? userId : "";
    }
    
    public void setUserId(String userId) {
        this.userId = userId;
    }
    
    public long getAppId() {
        return appId;
    }
    
    public void setAppId(long appId) {
        this.appId = appId;
    }
    
    public ZegoQuickStartTaskStatus getStatus() {
        return status != null ? status : ZegoQuickStartTaskStatus.IDLE;
    }
    
    public void setStatus(ZegoQuickStartTaskStatus status) {
        this.status = status;
    }
    
    @Override
    public String toString() {
        return "ZegoQuickStartTask{" +
                "taskId='" + taskId + '\'' +
                ", roomId='" + roomId + '\'' +
                ", streamId='" + streamId + '\'' +
                ", userId='" + userId + '\'' +
                ", status=" + status +
                '}';
    }
}

