//
//  ZegoConfigManager.h
//  ZegoDigitalHumanQuickStart
//
//  Created by Zego.
//

#import <Foundation/Foundation.h>
#import "../Models/ZegoConfig.h"

NS_ASSUME_NONNULL_BEGIN

/// 配置管理器
@interface ZegoConfigManager : NSObject

+ (instancetype)sharedManager;

/// 当前配置
@property (nonatomic, strong, readonly) ZegoConfig *currentConfig;

/// 加载配置
- (void)loadConfig;

/// 保存配置
- (void)saveConfig:(ZegoConfig *)config;

/// 重置为默认配置
- (void)resetToDefaultConfig;

@end

NS_ASSUME_NONNULL_END

