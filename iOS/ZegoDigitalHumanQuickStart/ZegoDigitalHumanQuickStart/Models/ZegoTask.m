//
//  ZegoTask.m
//  ZegoDigitalHumanQuickStart
//
//  Created by Zego.
//

#import "ZegoTask.h"

@implementation ZegoTask

- (instancetype)initWithDictionary:(NSDictionary *)dict {
    self = [super init];
    if (self) {
        // 边界检查
        if (![dict isKindOfClass:[NSDictionary class]]) {
            return nil;
        }
        
        _taskId = [dict[@"TaskId"] copy] ?: @"";
        _roomId = [dict[@"RoomId"] copy] ?: @"";
        _streamId = [dict[@"StreamId"] copy] ?: @"";
        _userId = [dict[@"UserID"] copy] ?: [dict[@"UserId"] copy] ?: @"";
        _userName = [dict[@"UserName"] copy] ?: @"";
        _token = [dict[@"Token"] copy];
        _appId = [dict[@"appId"] integerValue];
        _server = [dict[@"server"] copy];
        _status = ZegoTaskStatusIdle;
    }
    return self;
}

@end

