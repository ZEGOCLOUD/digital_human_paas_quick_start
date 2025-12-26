//
//  ZegoMainViewController.h
//  ZegoDigitalHumanQuickStart
//
//  Created by Zego.
//

#import <UIKit/UIKit.h>

@class ZegoDigitalView;
@class ZegoTaskControlView;
@class ZegoDriveControlView;
@class ZegoDigitalHumanPlaceholderView;
@class ZegoTask;
@class ZegoConfig;

NS_ASSUME_NONNULL_BEGIN

@protocol IZegoDigitalMobile;

@interface ZegoMainViewController : UIViewController

// UI组件
@property (nonatomic, strong) UIButton *toggleControlButton;
@property (nonatomic, strong) UILabel *statusLabel;
@property (nonatomic, strong) UIView *controlPanelContainerView;

// 数字人视图和逻辑
@property (nonatomic, strong) ZegoDigitalView *digitalHumanView;  // 纯视图
@property (nonatomic, strong, nullable) id<IZegoDigitalMobile> digitalMobile;  // 逻辑管理

@property (nonatomic, strong) ZegoTaskControlView *taskControlView;
@property (nonatomic, strong) ZegoDriveControlView *driveControlView;
@property (nonatomic, strong) ZegoDigitalHumanPlaceholderView *placeholderView;

// RTC引擎（不持有，使用sharedEngine）
@property (nonatomic, assign) BOOL rtcEngineCreated;

// 任务状态
@property (nonatomic, strong, nullable) ZegoTask *currentTask;
@property (nonatomic, copy, nullable) NSString *currentStreamId;
@property (nonatomic, copy, nullable) NSString *currentRoomId;
@property (nonatomic, copy, nullable) NSString *currentUserId;
@property (nonatomic, copy, nullable) NSString *currentToken;
@property (nonatomic, assign) NSInteger currentAppId;
@property (nonatomic, assign) BOOL isControlPanelVisible;
@property (nonatomic, assign) BOOL isRoomLogined;

// 配置和数据
@property (nonatomic, strong) ZegoConfig *config;

// 手势识别
@property (nonatomic, strong) UITapGestureRecognizer *tapGesture;

// 方法
/// 获取当前用户ID，如果不存在则生成新的.实际业务中请结合具体业务需求.
- (NSString *)getCurrentUserId;

@end

NS_ASSUME_NONNULL_END

