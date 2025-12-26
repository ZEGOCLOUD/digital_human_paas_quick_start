package com.example.zegodigitalhumanquickstart.view;

import android.content.Context;
import android.graphics.drawable.GradientDrawable;
import android.text.TextUtils;
import android.util.AttributeSet;
import android.view.Gravity;
import android.view.View;
import android.view.animation.AlphaAnimation;
import android.view.animation.Animation;
import android.widget.FrameLayout;
import android.widget.ImageView;
import android.widget.TextView;

import androidx.annotation.NonNull;
import androidx.annotation.Nullable;

import com.bumptech.glide.Glide;
import com.bumptech.glide.load.engine.DiskCacheStrategy;
import com.bumptech.glide.request.RequestOptions;
import com.example.zegodigitalhumanquickstart.R;

/**
 * 数字人占位视图
 * 在数字人未运行时显示封面和名称
 */
public class ZegoQuickStartDigitalHumanPlaceholderView extends FrameLayout {
    
    private static final String TAG = "PlaceholderView";
    private static final int ANIMATION_DURATION = 300;
    
    private ImageView backgroundImageView;  // 渐变背景
    private ImageView coverImageView;       // 头像（全屏平铺）
    private TextView nameTextView;          // 名称文本
    
    public ZegoQuickStartDigitalHumanPlaceholderView(@NonNull Context context) {
        super(context);
        init(context);
    }
    
    public ZegoQuickStartDigitalHumanPlaceholderView(@NonNull Context context, @Nullable AttributeSet attrs) {
        super(context, attrs);
        init(context);
    }
    
    public ZegoQuickStartDigitalHumanPlaceholderView(@NonNull Context context, @Nullable AttributeSet attrs, int defStyleAttr) {
        super(context, attrs, defStyleAttr);
        init(context);
    }
    
    private void init(Context context) {
        // 创建渐变背景ImageView
        backgroundImageView = new ImageView(context);
        backgroundImageView.setScaleType(ImageView.ScaleType.CENTER_CROP);
        
        // 创建渐变背景
        GradientDrawable gradientDrawable = new GradientDrawable(
                GradientDrawable.Orientation.TOP_BOTTOM,
                new int[]{
                        0xFF6670E6,  // 浅紫色 (0.4, 0.5, 0.9)
                        0xFF4D33B3   // 深紫色 (0.3, 0.2, 0.7)
                }
        );
        backgroundImageView.setImageDrawable(gradientDrawable);
        
        FrameLayout.LayoutParams backgroundParams = new FrameLayout.LayoutParams(
                LayoutParams.MATCH_PARENT,
                LayoutParams.MATCH_PARENT
        );
        addView(backgroundImageView, backgroundParams);
        
        // 创建头像ImageView（全屏平铺，aspect_fill方式）
        coverImageView = new ImageView(context);
        coverImageView.setScaleType(ImageView.ScaleType.CENTER_CROP); // aspect_fill方式
        coverImageView.setBackgroundColor(0x00000000); // 透明背景
        coverImageView.setClickable(false); // 禁用手势，避免遮挡父view的手势
        
        FrameLayout.LayoutParams coverParams = new FrameLayout.LayoutParams(
                LayoutParams.MATCH_PARENT,
                LayoutParams.MATCH_PARENT
        );
        addView(coverImageView, coverParams);
        
        // 设置默认占位图
        coverImageView.setImageResource(R.drawable.image_placeholder);
        
        // 创建名称文本（显示在底部）
        nameTextView = new TextView(context);
        nameTextView.setText("点击创建任务\n开始体验数字人");
        nameTextView.setTextColor(0xFF333333); // 深灰色
        nameTextView.setTextSize(24);
        nameTextView.setGravity(Gravity.CENTER);
        nameTextView.setTypeface(null, android.graphics.Typeface.BOLD);
        
        FrameLayout.LayoutParams textParams = new FrameLayout.LayoutParams(
                LayoutParams.MATCH_PARENT,
                LayoutParams.WRAP_CONTENT
        );
        textParams.gravity = Gravity.BOTTOM;
        textParams.bottomMargin = (int) (60 * context.getResources().getDisplayMetrics().density);
        textParams.leftMargin = (int) (20 * context.getResources().getDisplayMetrics().density);
        textParams.rightMargin = (int) (20 * context.getResources().getDisplayMetrics().density);
        addView(nameTextView, textParams);
    }
    
    /**
     * 更新显示内容
     *
     * @param name     数字人名称
     * @param coverUrl 封面URL
     */
    public void updateContent(String name, String coverUrl) {
        // 边界检查
        if (!TextUtils.isEmpty(name)) {
            nameTextView.setText(name);
        } else {
            nameTextView.setText("数字人");
        }
        
        // 加载封面图片（全屏平铺，aspect_fill方式）
        try {
            RequestOptions requestOptions = new RequestOptions()
                    .centerCrop() // aspect_fill方式，使用centerCrop而不是circleCrop
                    .placeholder(R.drawable.image_placeholder) // 默认占位图
                    .error(R.drawable.image_placeholder) // 错误时显示占位图
                    .diskCacheStrategy(DiskCacheStrategy.ALL);
            
            if (!TextUtils.isEmpty(coverUrl)) {
                Glide.with(getContext())
                        .load(coverUrl)
                        .apply(requestOptions)
                        .into(coverImageView);
            } else {
                // 没有URL时显示默认占位图
                Glide.with(getContext())
                        .load(R.drawable.image_placeholder)
                        .apply(requestOptions)
                        .into(coverImageView);
            }
        } catch (Exception e) {
            // 防止Glide在某些情况下崩溃
            android.util.Log.e(TAG, "加载图片失败", e);
            coverImageView.setImageResource(R.drawable.image_placeholder);
        }
    }
    
    /**
     * 显示占位视图（带动画）
     */
    public void show() {
        android.util.Log.d(TAG, "show() 被调用，当前visibility=" + getVisibility() + ", alpha=" + getAlpha());
        
        if (getVisibility() == VISIBLE && getAlpha() == 1.0f) {
            // 已经完全显示，无需重复
            android.util.Log.d(TAG, "已经完全显示，跳过");
            return;
        }
        
        setVisibility(VISIBLE);
        setAlpha(0.0f); // 从完全透明开始
        
        // 渐入动画
        AlphaAnimation fadeIn = new AlphaAnimation(0.0f, 1.0f);
        fadeIn.setDuration(ANIMATION_DURATION);
        fadeIn.setFillAfter(true); // 保持动画结束后的状态
        fadeIn.setAnimationListener(new Animation.AnimationListener() {
            @Override
            public void onAnimationStart(Animation animation) {
                android.util.Log.d(TAG, "动画开始");
            }
            
            @Override
            public void onAnimationEnd(Animation animation) {
                setAlpha(1.0f); // 确保alpha为1
                android.util.Log.d(TAG, "动画结束，alpha=" + getAlpha());
            }
            
            @Override
            public void onAnimationRepeat(Animation animation) {
            }
        });
        startAnimation(fadeIn);
        
        android.util.Log.d(TAG, "占位视图显示完成");
    }
    
    /**
     * 隐藏占位视图（带动画）
     */
    public void hide() {
        if (getVisibility() == GONE || getVisibility() == INVISIBLE) {
            return;
        }
        
        // 渐出动画
        AlphaAnimation fadeOut = new AlphaAnimation(1.0f, 0.0f);
        fadeOut.setDuration(ANIMATION_DURATION);
        fadeOut.setAnimationListener(new Animation.AnimationListener() {
            @Override
            public void onAnimationStart(Animation animation) {
            }
            
            @Override
            public void onAnimationEnd(Animation animation) {
                setVisibility(GONE);
            }
            
            @Override
            public void onAnimationRepeat(Animation animation) {
            }
        });
        startAnimation(fadeOut);
    }
}
