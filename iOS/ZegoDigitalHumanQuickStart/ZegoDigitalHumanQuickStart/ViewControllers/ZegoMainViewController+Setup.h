//
//  ZegoMainViewController+Setup.h
//  ZegoDigitalHumanModelQuickStart
//
//  Created by Zego.
//

#import "ZegoMainViewController.h"
#import "ZegoTaskControlView.h"
#import "ZegoDriveControlView.h"
#import <ZegoDigitalMobile/ZegoDigitalMobile.h>

NS_ASSUME_NONNULL_BEGIN

@interface ZegoMainViewController (Setup) <ZegoDigitalMobileDelegate>

// Setup methods
- (void)loadConfig;
- (void)loadInitialData;
- (void)setupUI;
- (void)setupGestureRecognizers;
- (void)handleTapGesture:(UITapGestureRecognizer *)gesture;
- (void)setupControlPanel;
- (void)viewDidLayoutSubviews;
- (void)toggleControlPanel;
- (void)updateStatus:(NSString *)status;
- (void)updatePlaceholderView;

// 数字人相关
- (void)startDigitalHumanWithConfig:(NSString *)base64Config;
- (void)stopDigitalHuman;

@end

NS_ASSUME_NONNULL_END

