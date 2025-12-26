//
//  ZegoDigitalHumanPlaceholderView.h
//  ZegoDigitalHumanQuickStart
//
//  Created by Zego.
//

#import <UIKit/UIKit.h>

NS_ASSUME_NONNULL_BEGIN

/// 头像点击回调
typedef void(^ZegoDigitalHumanPlaceholderViewTapBlock)(void);

/// 数字人占位视图
/// 在数字人未启动或停止时显示，包含图标和名称
@interface ZegoDigitalHumanPlaceholderView : UIView

/// 更新数字人信息
/// @param name 数字人名称
/// @param coverUrl 封面URL（可选，如果为nil则使用默认图标）
- (void)updateWithName:(NSString *)name coverUrl:(nullable NSString *)coverUrl;

/// 显示占位视图
- (void)show;

/// 隐藏占位视图
- (void)hide;

/// 头像点击回调（点击头像时触发）
@property (nonatomic, copy, nullable) ZegoDigitalHumanPlaceholderViewTapBlock avatarTapBlock;

@end

NS_ASSUME_NONNULL_END

