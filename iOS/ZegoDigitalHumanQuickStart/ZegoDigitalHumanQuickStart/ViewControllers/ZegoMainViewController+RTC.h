//
//  ZegoMainViewController+RTC.h
//  ZegoDigitalHumanModelQuickStart
//
//  Created by Zego.
//

#import "ZegoMainViewController.h"
#import <ZegoExpressEngine/ZegoExpressEngine.h>

NS_ASSUME_NONNULL_BEGIN

@interface ZegoMainViewController (RTC) <ZegoEventHandler, ZegoCustomVideoRenderHandler>

// RTC methods
- (void)initExpressEngineWithAppId:(NSInteger)appId;
- (void)destroyExpressEngine;
- (void)loginRoomWithCompletion:(void(^)(BOOL success))completion;
- (void)startPlayingStream:(NSString *)streamID;

@end

NS_ASSUME_NONNULL_END

