//
//  ZegoConfig.m
//  ZegoDigitalHumanQuickStart
//
//  Created by Zego.
//

#import "ZegoConfig.h"
#import "../Utils/ZegoCommonDefines.h"
#import "ZegoAPIConstants.h"
static NSString * const kZegoConfigKey = @"ZegoConfigKey";

@implementation ZegoConfig

+ (instancetype)defaultConfig {
    ZegoConfig *config = [[ZegoConfig alloc] init];
    config.serverURL = kZegoDefaultServerURL;
    return config;
}

- (void)save {
    NSMutableDictionary *dict = [NSMutableDictionary dictionary];
    dict[@"serverURL"] = self.serverURL ?: @"";
    
    [[NSUserDefaults standardUserDefaults] setObject:dict forKey:kZegoConfigKey];
    [[NSUserDefaults standardUserDefaults] synchronize];
}

- (void)load {
    NSDictionary *dict = [[NSUserDefaults standardUserDefaults] objectForKey:kZegoConfigKey];
    
    // 边界检查
    if (![dict isKindOfClass:[NSDictionary class]]) {
        return;
    }
    
    // 优先使用保存的配置，如果没有则使用默认值
    NSString *savedServerURL = dict[@"serverURL"];
    if (savedServerURL && [savedServerURL isKindOfClass:[NSString class]] && savedServerURL.length > 0) {
        self.serverURL = savedServerURL;
    } else {
        self.serverURL = kZegoDefaultServerURL;
    }
}

- (void)dealloc {
    // 内存清理
}

@end

