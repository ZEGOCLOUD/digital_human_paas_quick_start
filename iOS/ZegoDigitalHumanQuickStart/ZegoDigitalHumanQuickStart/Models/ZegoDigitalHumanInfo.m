//
//  ZegoDigitalHuman.m
//  ZegoDigitalHumanQuickStart
//
//  Created by Zego.
//

#import "ZegoDigitalHumanInfo.h"

@implementation ZegoDigitalHumanInfoModel

- (instancetype)initWithDictionary:(NSDictionary *)dict {
    self = [super init];
    if (self) {
        // 边界检查
        if (![dict isKindOfClass:[NSDictionary class]]) {
            return nil;
        }
        
        _digitalHumanId = [dict[@"DigitalHumanId"] copy] ?: @"";
        _name = [dict[@"Name"] copy] ?: @"";
        _coverUrl = [dict[@"AvatarUrl"] copy];
        _previewUrl = [dict[@"PreviewUrl"] copy];
        _isPublic = [dict[@"IsPublic"] boolValue];
        _appId = [dict[@"AppId"] integerValue];
        _token = [dict[@"Token"] copy];
        _expireTime = [dict[@"ExpireTime"] longLongValue];
    }
    return self;
}

@end

