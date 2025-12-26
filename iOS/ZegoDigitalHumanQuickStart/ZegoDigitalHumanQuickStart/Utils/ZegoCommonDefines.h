//
//  ZegoCommonDefines.h
//  ZegoDigitalHumanQuickStart
//
//  Created by Zego.
//

#ifndef ZegoCommonDefines_h
#define ZegoCommonDefines_h

#import <Foundation/Foundation.h>

// 输出模式
typedef NS_ENUM(NSInteger, ZegoOutputMode) {
    ZegoOutputModeLarge = 1,  // 大图模式（web）
    ZegoOutputModeSmall = 2   // 小图模式（mobile）
};

// 任务状态
typedef NS_ENUM(NSInteger, ZegoTaskStatus) {
    ZegoTaskStatusIdle = 0,      // 空闲
    ZegoTaskStatusRunning = 1,   // 运行中
    ZegoTaskStatusStopped = 2    // 已停止
};

// 驱动类型
typedef NS_ENUM(NSInteger, ZegoDriveType) {
    ZegoDriveTypeText = 0,    // 文本驱动
    ZegoDriveTypeAudio = 1,   // 音频驱动
    ZegoDriveTypeWsTTS = 2    // WebSocket TTS驱动
};

// 通知名称
extern NSString * const ZegoConfigDidChangeNotification;
extern NSString * const ZegoTaskStateDidChangeNotification;
extern NSString * const ZegoDigitalHumanDidStartNotification;
extern NSString * const ZegoDigitalHumanDidStopNotification;

// 默认视频配置
static const NSInteger kZegoDefaultVideoWidth = 720;
static const NSInteger kZegoDefaultVideoHeight = 1280;
static const NSInteger kZegoDefaultVideoBitrate = 3000000;

// 错误域
extern NSString * const ZegoErrorDomain;

// 应用错误码（避免与 ZegoExpressEngine SDK 的 ZegoErrorCode 冲突）
typedef NS_ENUM(NSInteger, ZegoAppErrorCode) {
    ZegoAppErrorCodeSuccess = 0,
    ZegoAppErrorCodeNetworkError = -1,
    ZegoAppErrorCodeInvalidParameter = -2,
    ZegoAppErrorCodeTaskNotFound = -3,
    ZegoAppErrorCodeSDKError = -4,
    ZegoAppErrorCodeUnknown = -999
};

#endif /* ZegoCommonDefines_h */

