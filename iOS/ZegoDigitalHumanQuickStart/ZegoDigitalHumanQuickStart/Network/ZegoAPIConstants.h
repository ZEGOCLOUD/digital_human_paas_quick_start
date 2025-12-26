//
//  ZegoAPIConstants.h
//  ZegoDigitalHumanQuickStart
//
//  Created by Zego.
//

#ifndef ZegoAPIConstants_h
#define ZegoAPIConstants_h

#import <Foundation/Foundation.h>


// 默认应用配置
// 注意：请替换为您的实际服务器地址
// 示例:"http://192.168.88.213:3000/api
static NSString * const kZegoDefaultServerURL = 你的服务端地址; // 请替换为您的服务器地址

// API Action 名称
static NSString * const kZegoActionGetDigitalHumanInfo = @"GetDigitalHumanInfo";
static NSString * const kZegoActionCreateDigitalHumanStreamTask = @"CreateDigitalHumanStreamTask";
static NSString * const kZegoActionStopDigitalHumanStreamTask = @"StopDigitalHumanStreamTask";
static NSString * const kZegoActionQueryDigitalHumanStreamTasks = @"QueryDigitalHumanStreamTasks";
static NSString * const kZegoActionDriveByText = @"DriveByText";
static NSString * const kZegoActionDriveByAudio = @"DriveByAudio";
static NSString * const kZegoActionDriveByWsStreamWithTTS = @"DriveByWsStreamWithTTS";
static NSString * const kZegoActionInterruptDriveTask = @"InterruptDriveTask";

// HTTP Method
static NSString * const kZegoHTTPMethodGET = @"GET";
static NSString * const kZegoHTTPMethodPOST = @"POST";

// HTTP Header Keys
static NSString * const kZegoHeaderAppId = @"X-App-Id";
static NSString * const kZegoHeaderContentType = @"Content-Type";

// Content Type
static NSString * const kZegoContentTypeJSON = @"application/json";

#endif /* ZegoAPIConstants_h */

