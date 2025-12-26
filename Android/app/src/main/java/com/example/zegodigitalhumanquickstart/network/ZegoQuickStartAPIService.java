package com.example.zegodigitalhumanquickstart.network;

import android.text.TextUtils;
import android.util.Log;

import com.example.zegodigitalhumanquickstart.model.ZegoQuickStartDigitalHuman;
import com.example.zegodigitalhumanquickstart.model.ZegoQuickStartTask;
import com.example.zegodigitalhumanquickstart.ZegoQuickStartConstants;
import com.google.gson.Gson;
import com.google.gson.JsonArray;
import com.google.gson.JsonObject;
import com.google.gson.reflect.TypeToken;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * API服务类
 * 封装所有数字人相关的API接口
 */
public class ZegoQuickStartAPIService {
    
    private static final String TAG = "ZegoQuickStartAPIService";
    
    private static ZegoQuickStartAPIService instance;
    private final ZegoQuickStartNetworkManager networkManager;
    private final Gson gson;
    
    private String serverURL;
    private long appId;  // 动态设置，从服务端返回的AppId
    
    private ZegoQuickStartAPIService() {
        networkManager = ZegoQuickStartNetworkManager.getInstance();
        gson = new Gson();
    }
    
    public static synchronized ZegoQuickStartAPIService getInstance() {
        if (instance == null) {
            instance = new ZegoQuickStartAPIService();
        }
        return instance;
    }
    
    /**
     * 设置服务器URL
     */
    public void setServerURL(String serverURL) {
        this.serverURL = serverURL;
    }
    
    /**
     * 设置AppID（动态设置，从服务端返回的AppId）
     * 注意：创建任务后，从服务端返回的AppId会调用此方法设置
     */
    public void setAppId(long appId) {
        this.appId = appId;
    }
    
    /**
     * 获取AppID
     */
    public long getAppId() {
        return this.appId;
    }
    
    /**
     * 构建完整的URL
     */
    private String buildURL(String action) {
        if (TextUtils.isEmpty(action)) {
            return serverURL;
        }
        
        String baseURL = serverURL;
        if (baseURL.endsWith("/")) {
            baseURL = baseURL.substring(0, baseURL.length() - 1);
        }
        
        String actionPath = action;
        if (actionPath.startsWith("/")) {
            actionPath = actionPath.substring(1);
        }
        
        return baseURL + "/" + actionPath;
    }
    
    /**
     * 构建请求头
     */
    private Map<String, String> buildHeaders() {
        Map<String, String> headers = new HashMap<>();
        headers.put(ZegoQuickStartAPIConstants.HEADER_CONTENT_TYPE, ZegoQuickStartAPIConstants.CONTENT_TYPE_JSON);
        return headers;
    }
    
    // ==================== API接口回调 ====================
    
    public interface DigitalHumanInfoCallback {
        void onSuccess(ZegoQuickStartDigitalHuman digitalHuman);
        void onFailure(int code, String message);
    }
    
    public interface TaskCallback {
        void onSuccess(JsonObject data);
        void onFailure(int code, String message);
    }
    
    public interface TaskListCallback {
        void onSuccess(List<ZegoQuickStartTask> tasks);
        void onFailure(int code, String message);
    }
    
    public interface CommonCallback {
        void onSuccess(JsonObject data);
        void onFailure(int code, String message);
    }
    
    /**
     * 统一响应解析结果
     */
    private static class ParsedResponse {
        int code;
        String message;
        JsonObject data;
    }
    
    /**
     * 统一解析服务端响应，规则：{Code, Message, Data:{}}
     */
    private ParsedResponse parseResponse(JsonObject response) {
        ParsedResponse result = new ParsedResponse();
        if (response == null) {
            result.code = ZegoQuickStartConstants.ERROR_CODE_PARSE_ERROR;
            result.message = "响应为空";
            result.data = new JsonObject();
            return result;
        }
        
        try {
            result.code = response.has("Code") ? response.get("Code").getAsInt() : 0;
            result.message = response.has("Message") ? response.get("Message").getAsString() : "";
            if (response.has("Data") && response.get("Data").isJsonObject()) {
                result.data = response.getAsJsonObject("Data");
            }
            
            if (result.data == null) {
                result.data = new JsonObject();
            }
        } catch (Exception e) {
            Log.e(TAG, "统一解析响应失败", e);
            result.code = ZegoQuickStartConstants.ERROR_CODE_PARSE_ERROR;
            result.message = "解析响应失败";
            result.data = new JsonObject();
        }
        
        return result;
    }
    
    // ==================== 2. 获取数字人详情 ====================
    
    /**
     * 获取数字人信息
     * @param userId 用户ID
     * @param callback 回调
     */
    public void getDigitalHumanInfo(String userId, DigitalHumanInfoCallback callback) {
        // 边界检查
        if (TextUtils.isEmpty(userId)) {
            if (callback != null) {
                callback.onFailure(ZegoQuickStartConstants.ERROR_CODE_INVALID_PARAMETER, "用户ID不能为空");
            }
            return;
        }
        
        String url = buildURL(ZegoQuickStartAPIConstants.ACTION_GET_DIGITAL_HUMAN_INFO);
        JsonObject params = new JsonObject();
        params.addProperty("UserId", userId);
        
        networkManager.post(url, params, buildHeaders(), new ZegoQuickStartNetworkManager.NetworkCallback() {
            @Override
            public void onSuccess(JsonObject response) {
                try {
                    ParsedResponse parsed = parseResponse(response);
                    if (parsed.code != 0) {
                        if (callback != null) {
                            callback.onFailure(parsed.code, TextUtils.isEmpty(parsed.message) ? "获取数字人信息失败" : parsed.message);
                        }
                        return;
                    }
                    JsonObject data = parsed.data != null ? parsed.data : new JsonObject();
                    
                    ZegoQuickStartDigitalHuman digitalHuman = gson.fromJson(data, ZegoQuickStartDigitalHuman.class);
                    if (callback != null) {
                        callback.onSuccess(digitalHuman);
                    }
                } catch (Exception e) {
                    Log.e(TAG, "解析数字人信息失败", e);
                    if (callback != null) {
                        callback.onFailure(ZegoQuickStartConstants.ERROR_CODE_PARSE_ERROR, "解析响应失败");
                    }
                }
            }
            
            @Override
            public void onFailure(int code, String message) {
                if (callback != null) {
                    callback.onFailure(code, message);
                }
            }
        });
    }
    
    // ==================== 3. 创建数字人流任务 ====================
    
    public void createDigitalHumanStreamTask(JsonObject config, TaskCallback callback) {
        // 边界检查
        if (config == null) {
            if (callback != null) {
                callback.onFailure(ZegoQuickStartConstants.ERROR_CODE_INVALID_PARAMETER, "配置参数不能为空");
            }
            return;
        }
        
        String url = buildURL(ZegoQuickStartAPIConstants.ACTION_CREATE_DIGITAL_HUMAN_STREAM_TASK);
        
        networkManager.post(url, config, buildHeaders(), new ZegoQuickStartNetworkManager.NetworkCallback() {
            @Override
            public void onSuccess(JsonObject response) {
                try {
                    ParsedResponse parsed = parseResponse(response);
                    if (parsed.code != 0) {
                        if (callback != null) {
                            callback.onFailure(parsed.code, TextUtils.isEmpty(parsed.message) ? "创建任务失败" : parsed.message);
                        }
                        return;
                    }
                    JsonObject data = parsed.data != null ? parsed.data : new JsonObject();
                    
                    if (callback != null) {
                        callback.onSuccess(data);
                    }
                } catch (Exception e) {
                    Log.e(TAG, "解析创建任务响应失败", e);
                    if (callback != null) {
                        callback.onFailure(ZegoQuickStartConstants.ERROR_CODE_PARSE_ERROR, "解析响应失败");
                    }
                }
            }
            
            @Override
            public void onFailure(int code, String message) {
                if (callback != null) {
                    callback.onFailure(code, message);
                }
            }
        });
    }
    
    // ==================== 4. 停止数字人流任务 ====================
    
    public void stopDigitalHumanStreamTask(String taskId, CommonCallback callback) {
        // 边界检查
        if (TextUtils.isEmpty(taskId)) {
            if (callback != null) {
                callback.onFailure(ZegoQuickStartConstants.ERROR_CODE_INVALID_PARAMETER, "任务ID不能为空");
            }
            return;
        }
        
        String url = buildURL(ZegoQuickStartAPIConstants.ACTION_STOP_DIGITAL_HUMAN_STREAM_TASK);
        JsonObject params = new JsonObject();
        params.addProperty("TaskId", taskId);
        
        networkManager.post(url, params, buildHeaders(), new ZegoQuickStartNetworkManager.NetworkCallback() {
            @Override
            public void onSuccess(JsonObject response) {
                handleCommonResponse(response, callback);
            }
            
            @Override
            public void onFailure(int code, String message) {
                if (callback != null) {
                    callback.onFailure(code, message);
                }
            }
        });
    }
    
    // ==================== 5. 查询数字人流任务列表 ====================
    
    public void queryDigitalHumanStreamTasks(TaskListCallback callback) {
        String url = buildURL(ZegoQuickStartAPIConstants.ACTION_QUERY_DIGITAL_HUMAN_STREAM_TASKS);
        
        networkManager.post(url, new JsonObject(), buildHeaders(), new ZegoQuickStartNetworkManager.NetworkCallback() {
            @Override
            public void onSuccess(JsonObject response) {
                try {
                    ParsedResponse parsed = parseResponse(response);
                    if (parsed.code != 0) {
                        if (callback != null) {
                            callback.onFailure(parsed.code, TextUtils.isEmpty(parsed.message) ? "查询任务列表失败" : parsed.message);
                        }
                        return;
                    }
                    JsonObject data = parsed.data != null ? parsed.data : new JsonObject();
                    
                    JsonArray tasksArray = null;
                    if (data.has("Tasks") && data.get("Tasks").isJsonArray()) {
                        tasksArray = data.getAsJsonArray("Tasks");
                    } else if (data.has("TaskList") && data.get("TaskList").isJsonArray()) {
                        tasksArray = data.getAsJsonArray("TaskList");
                    } else if (data.isJsonArray()) {
                        tasksArray = data.getAsJsonArray();
                    }
                    
                    if (tasksArray == null) {
                        tasksArray = new JsonArray();
                    }
                    
                    List<ZegoQuickStartTask> tasks = gson.fromJson(tasksArray, new TypeToken<List<ZegoQuickStartTask>>(){}.getType());
                    if (tasks == null) {
                        tasks = new ArrayList<>();
                    }
                    
                    if (callback != null) {
                        callback.onSuccess(tasks);
                    }
                } catch (Exception e) {
                    Log.e(TAG, "解析任务列表失败", e);
                    if (callback != null) {
                        callback.onFailure(ZegoQuickStartConstants.ERROR_CODE_PARSE_ERROR, "解析响应失败");
                    }
                }
            }
            
            @Override
            public void onFailure(int code, String message) {
                if (callback != null) {
                    callback.onFailure(code, message);
                }
            }
        });
    }
    
    // ==================== 6. 文本驱动 ====================
    
    public void driveByText(String taskId, CommonCallback callback) {
        // 边界检查
        if (TextUtils.isEmpty(taskId)) {
            if (callback != null) {
                callback.onFailure(ZegoQuickStartConstants.ERROR_CODE_INVALID_PARAMETER, "任务ID不能为空");
            }
            return;
        }
        
        String url = buildURL(ZegoQuickStartAPIConstants.ACTION_DRIVE_BY_TEXT);
        JsonObject params = new JsonObject();
        params.addProperty("TaskId", taskId);
        
        networkManager.post(url, params, buildHeaders(), new ZegoQuickStartNetworkManager.NetworkCallback() {
            @Override
            public void onSuccess(JsonObject response) {
                handleCommonResponse(response, callback);
            }
            
            @Override
            public void onFailure(int code, String message) {
                if (callback != null) {
                    callback.onFailure(code, message);
                }
            }
        });
    }
    
    // ==================== 7. 音频驱动 ====================
    
    public void driveByAudio(String taskId, CommonCallback callback) {
        // 边界检查
        if (TextUtils.isEmpty(taskId)) {
            if (callback != null) {
                callback.onFailure(ZegoQuickStartConstants.ERROR_CODE_INVALID_PARAMETER, "任务ID不能为空");
            }
            return;
        }
        
        String url = buildURL(ZegoQuickStartAPIConstants.ACTION_DRIVE_BY_AUDIO);
        JsonObject params = new JsonObject();
        params.addProperty("TaskId", taskId);
        
        networkManager.post(url, params, buildHeaders(), new ZegoQuickStartNetworkManager.NetworkCallback() {
            @Override
            public void onSuccess(JsonObject response) {
                handleCommonResponse(response, callback);
            }
            
            @Override
            public void onFailure(int code, String message) {
                if (callback != null) {
                    callback.onFailure(code, message);
                }
            }
        });
    }
    
    // ==================== 7.1 WebSocket TTS驱动 ====================
    
    public void driveByWsStreamWithTTS(String taskId, CommonCallback callback) {
        // 边界检查
        if (TextUtils.isEmpty(taskId)) {
            if (callback != null) {
                callback.onFailure(ZegoQuickStartConstants.ERROR_CODE_INVALID_PARAMETER, "任务ID不能为空");
            }
            return;
        }
        
        String url = buildURL(ZegoQuickStartAPIConstants.ACTION_DRIVE_BY_WS_STREAM_WITH_TTS);
        JsonObject params = new JsonObject();
        params.addProperty("TaskId", taskId);
        
        networkManager.post(url, params, buildHeaders(), new ZegoQuickStartNetworkManager.NetworkCallback() {
            @Override
            public void onSuccess(JsonObject response) {
                handleCommonResponse(response, callback);
            }
            
            @Override
            public void onFailure(int code, String message) {
                if (callback != null) {
                    callback.onFailure(code, message);
                }
            }
        });
    }
    
    // ==================== 8. 打断驱动任务 ====================
    
    public void interruptDriveTask(String taskId, CommonCallback callback) {
        // 边界检查
        if (TextUtils.isEmpty(taskId)) {
            if (callback != null) {
                callback.onFailure(ZegoQuickStartConstants.ERROR_CODE_INVALID_PARAMETER, "任务ID不能为空");
            }
            return;
        }
        
        String url = buildURL(ZegoQuickStartAPIConstants.ACTION_INTERRUPT_DRIVE_TASK);
        JsonObject params = new JsonObject();
        params.addProperty("TaskId", taskId);
        
        networkManager.post(url, params, buildHeaders(), new ZegoQuickStartNetworkManager.NetworkCallback() {
            @Override
            public void onSuccess(JsonObject response) {
                handleCommonResponse(response, callback);
            }
            
            @Override
            public void onFailure(int code, String message) {
                if (callback != null) {
                    callback.onFailure(code, message);
                }
            }
        });
    }
    
    // ==================== 通用响应处理 ====================
    
    private void handleCommonResponse(JsonObject response, CommonCallback callback) {
        try {
            ParsedResponse parsed = parseResponse(response);
            if (parsed.code != 0) {
                if (callback != null) {
                    callback.onFailure(parsed.code, TextUtils.isEmpty(parsed.message) ? "操作失败" : parsed.message);
                }
                return;
            }
            if (callback != null) {
                callback.onSuccess(parsed.data != null ? parsed.data : new JsonObject());
            }
        } catch (Exception e) {
            Log.e(TAG, "解析响应失败", e);
            if (callback != null) {
                callback.onFailure(ZegoQuickStartConstants.ERROR_CODE_PARSE_ERROR, "解析响应失败");
            }
        }
    }
}

