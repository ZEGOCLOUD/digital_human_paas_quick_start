//
//  ZegoAPIService.h
//  ZegoDigitalHumanModelQuickStart
//
//  Created by Zego.
//

#import <Foundation/Foundation.h>
#import "../Models/ZegoTask.h"
#import "../Models/ZegoConfig.h"

NS_ASSUME_NONNULL_BEGIN

@class ZegoDigitalHumanInfoModel;

typedef void(^ZegoAPISuccessBlock)(NSDictionary * _Nullable data);
typedef void(^ZegoAPIFailureBlock)(NSError *error, NSInteger code, NSString * _Nullable message);

/// API服务类
@interface ZegoAPIService : NSObject

+ (instancetype)sharedService;

/// 设置服务器URL
- (void)setServerURL:(NSString *)serverURL;


#pragma mark - Digital Human API

/// 获取数字人信息
/// @param userId 用户ID
/// @param success 成功回调
/// @param failure 失败回调
- (void)getDigitalHumanInfo:(NSString *)userId
                     success:(void(^)(ZegoDigitalHumanInfoModel *digitalHuman))success
                     failure:(ZegoAPIFailureBlock)failure;

#pragma mark - Stream Task API

/// 创建数字人视频流任务
/// @param config 配置参数
/// @param success 成功回调，返回字典包含：TaskId（任务ID）和 Base64Config（客户端渲染配置）
/// @param failure 失败回调
- (void)createDigitalHumanStreamTask:(NSDictionary *)config
                             success:(void(^)(NSDictionary *taskData))success
                             failure:(ZegoAPIFailureBlock)failure;

/// 停止数字人视频流任务
/// @param taskId 任务ID
/// @param success 成功回调
/// @param failure 失败回调
- (void)stopDigitalHumanStreamTask:(NSString *)taskId
                           success:(ZegoAPISuccessBlock)success
                           failure:(ZegoAPIFailureBlock)failure;

/// 查询数字人视频流任务列表
/// @param success 成功回调
/// @param failure 失败回调
- (void)queryDigitalHumanStreamTasks:(void(^)(NSArray<ZegoTask *> *tasks))success
                             failure:(ZegoAPIFailureBlock)failure;

#pragma mark - Drive API

/// 文本驱动数字人
/// @param taskId 任务ID
/// @param success 成功回调
/// @param failure 失败回调
- (void)driveByText:(NSString *)taskId
            success:(ZegoAPISuccessBlock)success
            failure:(ZegoAPIFailureBlock)failure;

/// 音频驱动数字人
/// @param taskId 任务ID
/// @param success 成功回调
/// @param failure 失败回调
- (void)driveByAudio:(NSString *)taskId
             success:(ZegoAPISuccessBlock)success
             failure:(ZegoAPIFailureBlock)failure;

/// WebSocket TTS驱动数字人
/// @param taskId 任务ID
/// @param success 成功回调
/// @param failure 失败回调
- (void)driveByWsStreamWithTTS:(NSString *)taskId
                       success:(ZegoAPISuccessBlock)success
                       failure:(ZegoAPIFailureBlock)failure;

/// 打断驱动任务
/// @param taskId 任务ID
/// @param success 成功回调
/// @param failure 失败回调
- (void)interruptDriveTask:(NSString *)taskId
                   success:(ZegoAPISuccessBlock)success
                   failure:(ZegoAPIFailureBlock)failure;

@end

NS_ASSUME_NONNULL_END

