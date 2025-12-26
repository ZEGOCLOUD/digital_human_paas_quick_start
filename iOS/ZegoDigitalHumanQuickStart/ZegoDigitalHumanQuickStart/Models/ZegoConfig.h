//
//  ZegoConfig.h
//  ZegoDigitalHumanQuickStart
//
//  Created by Zego.
//

#import <Foundation/Foundation.h>

NS_ASSUME_NONNULL_BEGIN

/// 应用配置
@interface ZegoConfig : NSObject

@property (nonatomic, copy) NSString *serverURL;           // 服务器URL

+ (instancetype)defaultConfig;
- (void)save;
- (void)load;

@end

NS_ASSUME_NONNULL_END

