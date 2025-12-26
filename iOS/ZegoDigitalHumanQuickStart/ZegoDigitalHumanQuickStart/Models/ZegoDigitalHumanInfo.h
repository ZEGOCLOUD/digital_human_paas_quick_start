//
//  ZegoDigitalHuman.h
//  ZegoDigitalHumanQuickStart
//
//  Created by Zego.
//

#import <Foundation/Foundation.h>

NS_ASSUME_NONNULL_BEGIN

/// 数字人形象模型
@interface ZegoDigitalHumanInfoModel : NSObject

@property (nonatomic, copy) NSString *digitalHumanId;      // 数字人ID
@property (nonatomic, copy) NSString *name;                // 数字人名称
@property (nonatomic, copy, nullable) NSString *coverUrl;  // 封面URL
@property (nonatomic, copy, nullable) NSString *previewUrl; // 预览URL
@property (nonatomic, assign) BOOL isPublic;               // 是否为公共数字人
@property (nonatomic, assign) NSInteger appId;             // AppID
@property (nonatomic, copy, nullable) NSString *token;      // Token
@property (nonatomic, assign) long long expireTime;        // Token过期时间

- (instancetype)initWithDictionary:(NSDictionary *)dict;

@end

NS_ASSUME_NONNULL_END

