//
//  ZegoNetworkManager.h
//  ZegoDigitalHumanQuickStart
//
//  Created by Zego.
//

#import <Foundation/Foundation.h>

NS_ASSUME_NONNULL_BEGIN

typedef void(^ZegoNetworkSuccessBlock)(NSDictionary * _Nullable response);
typedef void(^ZegoNetworkFailureBlock)(NSError *error);

/// 网络请求管理器
@interface ZegoNetworkManager : NSObject

+ (instancetype)sharedManager;

/// 发送GET请求
/// @param urlString URL字符串
/// @param parameters 请求参数
/// @param headers 请求头
/// @param success 成功回调
/// @param failure 失败回调
- (void)GET:(NSString *)urlString
 parameters:(nullable NSDictionary *)parameters
    headers:(nullable NSDictionary *)headers
    success:(ZegoNetworkSuccessBlock)success
    failure:(ZegoNetworkFailureBlock)failure;

/// 发送POST请求
/// @param urlString URL字符串
/// @param parameters 请求参数
/// @param headers 请求头
/// @param success 成功回调
/// @param failure 失败回调
- (void)POST:(NSString *)urlString
  parameters:(nullable NSDictionary *)parameters
     headers:(nullable NSDictionary *)headers
     success:(ZegoNetworkSuccessBlock)success
     failure:(ZegoNetworkFailureBlock)failure;

@end

NS_ASSUME_NONNULL_END

