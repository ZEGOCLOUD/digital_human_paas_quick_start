package com.example.zegodigitalhumanquickstart.network;

import android.os.Handler;
import android.os.Looper;
import android.util.Log;

import com.example.zegodigitalhumanquickstart.util.ZegoQuickStartConstants;
import com.google.gson.Gson;
import com.google.gson.JsonObject;

import java.io.IOException;
import java.util.Map;
import java.util.concurrent.TimeUnit;

import okhttp3.Call;
import okhttp3.Callback;
import okhttp3.MediaType;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.RequestBody;
import okhttp3.Response;

/**
 * 网络请求管理器
 * 封装OkHttp的GET和POST请求，提供统一的网络请求接口
 */
public class ZegoQuickStartNetworkManager {
    
    private static final String TAG = "ZegoQuickStartNetworkManager";
    private static final MediaType JSON = MediaType.get("application/json; charset=utf-8");
    
    private static ZegoQuickStartNetworkManager instance;
    private final OkHttpClient okHttpClient;
    private final Gson gson;
    private final Handler mainHandler;
    
    private ZegoQuickStartNetworkManager() {
        // 配置OkHttpClient
        okHttpClient = new OkHttpClient.Builder()
                .connectTimeout(ZegoQuickStartConstants.NETWORK_TIMEOUT, TimeUnit.SECONDS)
                .readTimeout(ZegoQuickStartConstants.NETWORK_TIMEOUT, TimeUnit.SECONDS)
                .writeTimeout(ZegoQuickStartConstants.NETWORK_TIMEOUT, TimeUnit.SECONDS)
                .retryOnConnectionFailure(true)
                .build();
        
        gson = new Gson();
        mainHandler = new Handler(Looper.getMainLooper());
    }
    
    public static synchronized ZegoQuickStartNetworkManager getInstance() {
        if (instance == null) {
            instance = new ZegoQuickStartNetworkManager();
        }
        return instance;
    }
    
    /**
     * 网络请求回调接口
     */
    public interface NetworkCallback {
        void onSuccess(JsonObject response);
        void onFailure(int code, String message);
    }
    
    /**
     * 发送GET请求
     *
     * @param url      请求URL
     * @param headers  请求头
     * @param callback 回调接口
     */
    public void get(String url, Map<String, String> headers, NetworkCallback callback) {
        // 边界检查
        if (url == null || url.isEmpty()) {
            if (callback != null) {
                mainHandler.post(() -> callback.onFailure(ZegoQuickStartConstants.ERROR_CODE_INVALID_PARAMETER, "URL不能为空"));
            }
            return;
        }
        
        try {
            Request.Builder requestBuilder = new Request.Builder().url(url);
            
            // 添加请求头
            if (headers != null && !headers.isEmpty()) {
                for (Map.Entry<String, String> entry : headers.entrySet()) {
                    if (entry.getKey() != null && entry.getValue() != null) {
                        requestBuilder.addHeader(entry.getKey(), entry.getValue());
                    }
                }
            }
            
            Request request = requestBuilder.build();
            
            Log.d(TAG, "GET请求: " + url);
            
            okHttpClient.newCall(request).enqueue(new Callback() {
                @Override
                public void onFailure(Call call, IOException e) {
                    Log.e(TAG, "GET请求失败: " + url, e);
                    if (callback != null) {
                        mainHandler.post(() -> callback.onFailure(
                                ZegoQuickStartConstants.ERROR_CODE_NETWORK_ERROR,
                                "网络请求失败: " + e.getMessage()
                        ));
                    }
                }
                
                @Override
                public void onResponse(Call call, Response response) {
                    handleResponse(response, callback);
                }
            });
            
        } catch (Exception e) {
            Log.e(TAG, "GET请求异常: " + url, e);
            if (callback != null) {
                mainHandler.post(() -> callback.onFailure(
                        ZegoQuickStartConstants.ERROR_CODE_UNKNOWN,
                        "请求异常: " + e.getMessage()
                ));
            }
        }
    }
    
    /**
     * 发送POST请求
     *
     * @param url        请求URL
     * @param parameters 请求参数（JSON对象）
     * @param headers    请求头
     * @param callback   回调接口
     */
    public void post(String url, JsonObject parameters, Map<String, String> headers, NetworkCallback callback) {
        // 边界检查
        if (url == null || url.isEmpty()) {
            if (callback != null) {
                mainHandler.post(() -> callback.onFailure(ZegoQuickStartConstants.ERROR_CODE_INVALID_PARAMETER, "URL不能为空"));
            }
            return;
        }
        
        try {
            // 构建请求体
            String jsonBody = parameters != null ? parameters.toString() : "{}";
            RequestBody body = RequestBody.create(jsonBody, JSON);
            
            Request.Builder requestBuilder = new Request.Builder()
                    .url(url)
                    .post(body);
            
            // 添加请求头
            if (headers != null && !headers.isEmpty()) {
                for (Map.Entry<String, String> entry : headers.entrySet()) {
                    if (entry.getKey() != null && entry.getValue() != null) {
                        requestBuilder.addHeader(entry.getKey(), entry.getValue());
                    }
                }
            }
            
            Request request = requestBuilder.build();
            
            Log.d(TAG, "POST请求: " + url);
            Log.d(TAG, "请求体: " + jsonBody);
            
            okHttpClient.newCall(request).enqueue(new Callback() {
                @Override
                public void onFailure(Call call, IOException e) {
                    Log.e(TAG, "POST请求失败: " + url, e);
                    if (callback != null) {
                        mainHandler.post(() -> callback.onFailure(
                                ZegoQuickStartConstants.ERROR_CODE_NETWORK_ERROR,
                                "网络请求失败: " + e.getMessage()
                        ));
                    }
                }
                
                @Override
                public void onResponse(Call call, Response response) {
                    handleResponse(response, callback);
                }
            });
            
        } catch (Exception e) {
            Log.e(TAG, "POST请求异常: " + url, e);
            if (callback != null) {
                mainHandler.post(() -> callback.onFailure(
                        ZegoQuickStartConstants.ERROR_CODE_UNKNOWN,
                        "请求异常: " + e.getMessage()
                ));
            }
        }
    }
    
    /**
     * 处理响应
     */
    private void handleResponse(Response response, NetworkCallback callback) {
        try {
            if (!response.isSuccessful()) {
                Log.e(TAG, "响应失败: code=" + response.code());
                if (callback != null) {
                    int code = response.code();
                    mainHandler.post(() -> callback.onFailure(code, "HTTP错误: " + code));
                }
                return;
            }
            
            String responseBody = response.body() != null ? response.body().string() : null;
            
            // 边界检查
            if (responseBody == null || responseBody.isEmpty()) {
                if (callback != null) {
                    mainHandler.post(() -> callback.onFailure(
                            ZegoQuickStartConstants.ERROR_CODE_PARSE_ERROR,
                            "响应体为空"
                    ));
                }
                return;
            }
            
            Log.d(TAG, "响应成功: " + responseBody);
            
            // 解析JSON
            JsonObject jsonObject = gson.fromJson(responseBody, JsonObject.class);
            
            if (callback != null) {
                mainHandler.post(() -> callback.onSuccess(jsonObject));
            }
            
        } catch (Exception e) {
            Log.e(TAG, "响应解析失败", e);
            if (callback != null) {
                mainHandler.post(() -> callback.onFailure(
                        ZegoQuickStartConstants.ERROR_CODE_PARSE_ERROR,
                        "响应解析失败: " + e.getMessage()
                ));
            }
        } finally {
            if (response.body() != null) {
                response.body().close();
            }
        }
    }
}

