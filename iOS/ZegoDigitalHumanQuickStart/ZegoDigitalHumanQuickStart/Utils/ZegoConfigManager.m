//
//  ZegoConfigManager.m
//  ZegoDigitalHumanQuickStart
//
//  Created by Zego.
//

#import "ZegoConfigManager.h"
#import "ZegoCommonDefines.h"

@interface ZegoConfigManager ()

@property (nonatomic, strong, readwrite) ZegoConfig *currentConfig;

@end

@implementation ZegoConfigManager

+ (instancetype)sharedManager {
    static ZegoConfigManager *instance = nil;
    static dispatch_once_t onceToken;
    dispatch_once(&onceToken, ^{
        instance = [[ZegoConfigManager alloc] init];
        [instance loadConfig];
    });
    return instance;
}

- (void)loadConfig {
    ZegoConfig *config = [ZegoConfig defaultConfig];
    [config load];
    self.currentConfig = config;
    
    // 发送配置变更通知
    [[NSNotificationCenter defaultCenter] postNotificationName:ZegoConfigDidChangeNotification object:config];
}

- (void)saveConfig:(ZegoConfig *)config {
    // 边界检查
    if (!config) {
        return;
    }
    
    [config save];
    self.currentConfig = config;
    
    // 发送配置变更通知
    [[NSNotificationCenter defaultCenter] postNotificationName:ZegoConfigDidChangeNotification object:config];
}

- (void)resetToDefaultConfig {
    ZegoConfig *config = [ZegoConfig defaultConfig];
    [self saveConfig:config];
}

- (void)dealloc {
    // 内存清理
    _currentConfig = nil;
}

@end

