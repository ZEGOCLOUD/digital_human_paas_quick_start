package com.example.zegodigitalhumanquickstart.callback;

/**
 * 驱动控制视图回调接口
 */
public interface ZegoQuickStartDriveControlViewCallback {
    
    /**
     * 点击文本驱动按钮
     */
    void onTextDriveClicked();
    
    /**
     * 点击音频驱动按钮
     */
    void onAudioDriveClicked();
    
    /**
     * 点击WebSocket TTS驱动按钮
     */
    void onWsTTSDriveClicked();
}
