package com.example.zegodigitalhumanquickstart.view;

import android.content.Context;
import android.graphics.Color;
import android.util.AttributeSet;
import android.view.Gravity;
import android.view.View;
import android.widget.Button;
import android.widget.FrameLayout;
import android.widget.LinearLayout;
import android.widget.ProgressBar;
import android.widget.TextView;

import androidx.annotation.NonNull;
import androidx.annotation.Nullable;

import com.example.zegodigitalhumanquickstart.callback.ZegoQuickStartDriveControlViewCallback;
import com.example.zegodigitalhumanquickstart.model.ZegoQuickStartDriveType;

/**
 * 驱动控制视图
 * 包含文本、音频、WebSocket TTS三种驱动方式
 */
public class ZegoQuickStartDriveControlView extends LinearLayout {
    
    private static final String TAG = "ZegoQuickStartDriveControlView";
    
    // 标题
    private TextView titleTextView;
    
    // 文本驱动按钮容器
    private FrameLayout textDriveButtonContainer;
    private Button textDriveButton;
    private ProgressBar textLoadingView;
    
    // 音频驱动按钮容器
    private FrameLayout audioDriveButtonContainer;
    private Button audioDriveButton;
    private ProgressBar audioLoadingView;
    
    // WebSocket TTS驱动按钮容器
    private FrameLayout wsTTSDriveButtonContainer;
    private Button wsTTSDriveButton;
    private ProgressBar wsTTSLoadingView;
    
    private ZegoQuickStartDriveControlViewCallback callback;
    
    public ZegoQuickStartDriveControlView(@NonNull Context context) {
        super(context);
        init(context);
    }
    
    public ZegoQuickStartDriveControlView(@NonNull Context context, @Nullable AttributeSet attrs) {
        super(context, attrs);
        init(context);
    }
    
    public ZegoQuickStartDriveControlView(@NonNull Context context, @Nullable AttributeSet attrs, int defStyleAttr) {
        super(context, attrs, defStyleAttr);
        init(context);
    }
    
    private void init(Context context) {
        setOrientation(LinearLayout.VERTICAL);
        setBackgroundColor(Color.TRANSPARENT);
        
        float density = context.getResources().getDisplayMetrics().density;
        int padding = (int) (15 * density);
        setPadding(padding, padding, padding, padding);
        
        // 标题
        titleTextView = new TextView(context);
        titleTextView.setText("驱动控制");
        titleTextView.setTextColor(Color.parseColor("#333333"));
        titleTextView.setTextSize(16);
        titleTextView.setTypeface(null, android.graphics.Typeface.BOLD);
        LinearLayout.LayoutParams titleParams = new LinearLayout.LayoutParams(
                LayoutParams.MATCH_PARENT,
                LayoutParams.WRAP_CONTENT
        );
        titleParams.bottomMargin = (int) (15 * density);
        addView(titleTextView, titleParams);
        
        // 文本驱动按钮容器
        textDriveButtonContainer = createButtonContainer(context, density);
        textDriveButton = createDriveButton(context, "文本驱动", density);
        textDriveButton.setOnClickListener(v -> {
            if (callback != null) {
                callback.onTextDriveClicked();
            }
        });
        textDriveButtonContainer.addView(textDriveButton);
        
        textLoadingView = createLoadingIndicator(context);
        textDriveButtonContainer.addView(textLoadingView);
        addView(textDriveButtonContainer, createButtonLayoutParams(density));
        
        // 音频驱动按钮容器
        audioDriveButtonContainer = createButtonContainer(context, density);
        audioDriveButton = createDriveButton(context, "音频驱动", density);
        audioDriveButton.setOnClickListener(v -> {
            if (callback != null) {
                callback.onAudioDriveClicked();
            }
        });
        audioDriveButtonContainer.addView(audioDriveButton);
        
        audioLoadingView = createLoadingIndicator(context);
        audioDriveButtonContainer.addView(audioLoadingView);
        addView(audioDriveButtonContainer, createButtonLayoutParams(density));
        
        // WebSocket TTS驱动按钮容器
        wsTTSDriveButtonContainer = createButtonContainer(context, density);
        wsTTSDriveButton = createDriveButton(context, "WebSocket TTS驱动", density);
        wsTTSDriveButton.setOnClickListener(v -> {
            if (callback != null) {
                callback.onWsTTSDriveClicked();
            }
        });
        wsTTSDriveButtonContainer.addView(wsTTSDriveButton);
        
        wsTTSLoadingView = createLoadingIndicator(context);
        wsTTSDriveButtonContainer.addView(wsTTSLoadingView);
        LinearLayout.LayoutParams wsTTSParams = createButtonLayoutParams(density);
        wsTTSParams.bottomMargin = (int) (15 * density);
        addView(wsTTSDriveButtonContainer, wsTTSParams);
    }
    
    // ==================== 辅助方法 ====================
    
    private FrameLayout createButtonContainer(Context context, float density) {
        FrameLayout container = new FrameLayout(context);
        container.setLayoutParams(new LinearLayout.LayoutParams(
                LayoutParams.MATCH_PARENT,
                (int) (44 * density)
        ));
        return container;
    }
    
    private Button createDriveButton(Context context, String text, float density) {
        Button button = new Button(context);
        button.setText(text);
        button.setTextColor(Color.WHITE);
        button.setTextSize(16);
        button.setAllCaps(false);
        
        // 设置背景颜色和圆角
        float radius = 8 * density;
        android.graphics.drawable.GradientDrawable drawable = new android.graphics.drawable.GradientDrawable();
        drawable.setColor(Color.parseColor("#1791FF"));
        drawable.setCornerRadius(radius);
        button.setBackground(drawable);
        
        // 设置按钮布局参数
        FrameLayout.LayoutParams params = new FrameLayout.LayoutParams(
                LayoutParams.MATCH_PARENT,
                LayoutParams.MATCH_PARENT
        );
        button.setLayoutParams(params);
        
        return button;
    }
    
    private ProgressBar createLoadingIndicator(Context context) {
        ProgressBar progressBar = new ProgressBar(context);
        progressBar.setVisibility(GONE);
        progressBar.setIndeterminate(true);
        progressBar.setIndeterminateTintList(android.content.res.ColorStateList.valueOf(Color.WHITE));
        
        // 设置居中
        FrameLayout.LayoutParams params = new FrameLayout.LayoutParams(
                LayoutParams.WRAP_CONTENT,
                LayoutParams.WRAP_CONTENT
        );
        params.gravity = Gravity.CENTER;
        progressBar.setLayoutParams(params);
        
        return progressBar;
    }
    
    private LinearLayout.LayoutParams createButtonLayoutParams(float density) {
        LinearLayout.LayoutParams params = new LinearLayout.LayoutParams(
                LayoutParams.MATCH_PARENT,
                (int) (44 * density)
        );
        params.topMargin = (int) (10 * density);
        return params;
    }
    
    // ==================== 公共方法 ====================
    
    public void setCallback(ZegoQuickStartDriveControlViewCallback callback) {
        this.callback = callback;
    }
    
    public void setLoading(ZegoQuickStartDriveType driveType, boolean loading) {
        ProgressBar loadingView = null;
        Button button = null;
        String buttonText = "";
        
        switch (driveType) {
            case TEXT:
                loadingView = textLoadingView;
                button = textDriveButton;
                buttonText = "文本驱动";
                break;
            case AUDIO:
                loadingView = audioLoadingView;
                button = audioDriveButton;
                buttonText = "音频驱动";
                break;
            case WS_TTS:
                loadingView = wsTTSLoadingView;
                button = wsTTSDriveButton;
                buttonText = "WebSocket TTS驱动";
                break;
            default:
                return;
        }
        
        if (loadingView != null && button != null) {
            if (loading) {
                loadingView.setVisibility(VISIBLE);
                button.setEnabled(false);
                button.setText("");
            } else {
                loadingView.setVisibility(GONE);
                button.setEnabled(true);
                button.setText(buttonText);
            }
        }
    }
    
    public void setDriveButtonsEnabled(boolean enabled) {
        textDriveButton.setEnabled(enabled);
        audioDriveButton.setEnabled(enabled);
        wsTTSDriveButton.setEnabled(enabled);
    }
}
