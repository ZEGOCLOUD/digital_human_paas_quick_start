//
//  ZegoTask.h
//  ZegoDigitalHumanQuickStart
//
//  Created by Zego.
//

#import <Foundation/Foundation.h>
#import "../Utils/ZegoCommonDefines.h"

NS_ASSUME_NONNULL_BEGIN

/// 任务模型
@interface ZegoTask : NSObject

@property (nonatomic, copy) NSString *taskId;              // 任务ID
@property (nonatomic, copy) NSString *roomId;              // 房间ID
@property (nonatomic, copy) NSString *streamId;            // 流ID
@property (nonatomic, copy) NSString *userId;              // 用户ID
@property (nonatomic, copy) NSString *userName;            // 用户名
@property (nonatomic, copy, nullable) NSString *token;     // Token
@property (nonatomic, assign) NSInteger appId;             // AppID
@property (nonatomic, copy, nullable) NSString *server;    // 服务器地址
@property (nonatomic, assign) ZegoTaskStatus status;       // 任务状态

- (instancetype)initWithDictionary:(NSDictionary *)dict;

@end

NS_ASSUME_NONNULL_END

