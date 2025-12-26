//
//  ZegoDriveControlView.h
//  ZegoDigitalHumanQuickStart
//
//  Created by Zego.
//

#import <UIKit/UIKit.h>
#import "ZegoCommonDefines.h"

NS_ASSUME_NONNULL_BEGIN

@class ZegoDriveControlView;

/// 驱动控制视图代理
@protocol ZegoDriveControlViewDelegate <NSObject>

@optional
/// 点击文本驱动按钮
/// @param view 视图
- (void)driveControlViewDidTapTextDrive:(ZegoDriveControlView *)view;

/// 点击音频驱动按钮
/// @param view 视图
- (void)driveControlViewDidTapAudioDrive:(ZegoDriveControlView *)view;

/// 点击WebSocket TTS驱动按钮
/// @param view 视图
- (void)driveControlViewDidTapWsTTSDrive:(ZegoDriveControlView *)view;

@end

/// 驱动控制视图
@interface ZegoDriveControlView : UIView

/// 代理
@property (nonatomic, weak, nullable) id<ZegoDriveControlViewDelegate> delegate;

/// 当前驱动类型
@property (nonatomic, assign, readonly) ZegoDriveType currentDriveType;


/// 设置按钮Loading状态
/// @param loading 是否Loading
/// @param driveType 驱动类型
- (void)setLoading:(BOOL)loading forDriveType:(ZegoDriveType)driveType;

/// 设置按钮可用状态
/// @param enabled 是否可用
- (void)setDriveButtonsEnabled:(BOOL)enabled;

@end

NS_ASSUME_NONNULL_END
