//
//  ZegoDriveControlView.m
//  ZegoDigitalHumanModelQuickStart
//
//  Created by Zego.
//

#import "ZegoDriveControlView.h"

@interface ZegoDriveControlView ()

// 标题
@property (nonatomic, strong) UILabel *titleLabel;

// 文本驱动按钮
@property (nonatomic, strong) UIButton *textDriveButton;
@property (nonatomic, strong) UIActivityIndicatorView *textLoadingView;

// 音频驱动按钮
@property (nonatomic, strong) UIButton *audioDriveButton;
@property (nonatomic, strong) UIActivityIndicatorView *audioLoadingView;

// WebSocket TTS驱动按钮
@property (nonatomic, strong) UIButton *wsTTSDriveButton;
@property (nonatomic, strong) UIActivityIndicatorView *wsTTSLoadingView;

@property (nonatomic, assign, readwrite) ZegoDriveType currentDriveType;

@end

@implementation ZegoDriveControlView

#pragma mark - Lifecycle

- (instancetype)initWithFrame:(CGRect)frame {
    self = [super initWithFrame:frame];
    if (self) {
        [self setupViews];
        _currentDriveType = ZegoDriveTypeText;
    }
    return self;
}

- (instancetype)initWithCoder:(NSCoder *)coder {
    self = [super initWithCoder:coder];
    if (self) {
        [self setupViews];
        _currentDriveType = ZegoDriveTypeText;
    }
    return self;
}

#pragma mark - Setup Views

- (void)setupViews {
    self.backgroundColor = [UIColor clearColor];
    self.layer.cornerRadius = 12;
    
    // 标题
    self.titleLabel = [[UILabel alloc] init];
    self.titleLabel.text = @"驱动控制";
    self.titleLabel.font = [UIFont boldSystemFontOfSize:16];
    self.titleLabel.textColor = [UIColor colorWithRed:0.2 green:0.2 blue:0.2 alpha:1.0];
    self.titleLabel.translatesAutoresizingMaskIntoConstraints = NO;
    [self addSubview:self.titleLabel];
    
    // 文本驱动按钮
    self.textDriveButton = [self createDriveButton:@"文本驱动" action:@selector(textDriveTapped)];
    self.textDriveButton.translatesAutoresizingMaskIntoConstraints = NO;
    [self addSubview:self.textDriveButton];
    
    self.textLoadingView = [self createLoadingIndicator];
    self.textLoadingView.translatesAutoresizingMaskIntoConstraints = NO;
    [self.textDriveButton addSubview:self.textLoadingView];
    
    // 音频驱动按钮
    self.audioDriveButton = [self createDriveButton:@"音频驱动" action:@selector(audioDriveTapped)];
    self.audioDriveButton.translatesAutoresizingMaskIntoConstraints = NO;
    [self addSubview:self.audioDriveButton];
    
    self.audioLoadingView = [self createLoadingIndicator];
    self.audioLoadingView.translatesAutoresizingMaskIntoConstraints = NO;
    [self.audioDriveButton addSubview:self.audioLoadingView];
    
    // WebSocket TTS驱动按钮
    self.wsTTSDriveButton = [self createDriveButton:@"WebSocket TTS驱动" action:@selector(wsTTSDriveTapped)];
    self.wsTTSDriveButton.translatesAutoresizingMaskIntoConstraints = NO;
    [self addSubview:self.wsTTSDriveButton];
    
    self.wsTTSLoadingView = [self createLoadingIndicator];
    self.wsTTSLoadingView.translatesAutoresizingMaskIntoConstraints = NO;
    [self.wsTTSDriveButton addSubview:self.wsTTSLoadingView];
    
    // 布局
    [NSLayoutConstraint activateConstraints:@[
        [self.titleLabel.topAnchor constraintEqualToAnchor:self.topAnchor constant:15],
        [self.titleLabel.leadingAnchor constraintEqualToAnchor:self.leadingAnchor constant:15],
        [self.titleLabel.trailingAnchor constraintEqualToAnchor:self.trailingAnchor constant:-15],
        
        [self.textDriveButton.topAnchor constraintEqualToAnchor:self.titleLabel.bottomAnchor constant:15],
        [self.textDriveButton.leadingAnchor constraintEqualToAnchor:self.leadingAnchor constant:15],
        [self.textDriveButton.trailingAnchor constraintEqualToAnchor:self.trailingAnchor constant:-15],
        [self.textDriveButton.heightAnchor constraintEqualToConstant:44],
        
        [self.audioDriveButton.topAnchor constraintEqualToAnchor:self.textDriveButton.bottomAnchor constant:10],
        [self.audioDriveButton.leadingAnchor constraintEqualToAnchor:self.leadingAnchor constant:15],
        [self.audioDriveButton.trailingAnchor constraintEqualToAnchor:self.trailingAnchor constant:-15],
        [self.audioDriveButton.heightAnchor constraintEqualToConstant:44],
        
        [self.wsTTSDriveButton.topAnchor constraintEqualToAnchor:self.audioDriveButton.bottomAnchor constant:10],
        [self.wsTTSDriveButton.leadingAnchor constraintEqualToAnchor:self.leadingAnchor constant:15],
        [self.wsTTSDriveButton.trailingAnchor constraintEqualToAnchor:self.trailingAnchor constant:-15],
        [self.wsTTSDriveButton.heightAnchor constraintEqualToConstant:44],
        [self.wsTTSDriveButton.bottomAnchor constraintEqualToAnchor:self.bottomAnchor constant:-15],
        
        [self.textLoadingView.centerXAnchor constraintEqualToAnchor:self.textDriveButton.centerXAnchor],
        [self.textLoadingView.centerYAnchor constraintEqualToAnchor:self.textDriveButton.centerYAnchor],
        
        [self.audioLoadingView.centerXAnchor constraintEqualToAnchor:self.audioDriveButton.centerXAnchor],
        [self.audioLoadingView.centerYAnchor constraintEqualToAnchor:self.audioDriveButton.centerYAnchor],
        
        [self.wsTTSLoadingView.centerXAnchor constraintEqualToAnchor:self.wsTTSDriveButton.centerXAnchor],
        [self.wsTTSLoadingView.centerYAnchor constraintEqualToAnchor:self.wsTTSDriveButton.centerYAnchor]
    ]];
    
    _currentDriveType = ZegoDriveTypeText;
}


#pragma mark - Helper Methods


- (UIButton *)createDriveButton:(NSString *)title action:(SEL)action {
    UIButton *button = [UIButton buttonWithType:UIButtonTypeSystem];
    [button setTitle:title forState:UIControlStateNormal];
    button.backgroundColor = [UIColor colorWithRed:0.09 green:0.57 blue:1.0 alpha:1.0];
    button.tintColor = [UIColor whiteColor];
    [button setTitleColor:[UIColor whiteColor] forState:UIControlStateNormal];
    button.titleLabel.font = [UIFont systemFontOfSize:16 weight:UIFontWeightMedium];
    button.layer.cornerRadius = 8;
    [button addTarget:self action:action forControlEvents:UIControlEventTouchUpInside];
    return button;
}

- (UIActivityIndicatorView *)createLoadingIndicator {
    UIActivityIndicatorView *indicator;
    if (@available(iOS 13.0, *)) {
        indicator = [[UIActivityIndicatorView alloc] initWithActivityIndicatorStyle:UIActivityIndicatorViewStyleMedium];
    } else {
        indicator = [[UIActivityIndicatorView alloc] initWithActivityIndicatorStyle:UIActivityIndicatorViewStyleWhite];
    }
    indicator.hidesWhenStopped = YES;
    return indicator;
}


#pragma mark - Public Methods

- (void)setLoading:(BOOL)loading forDriveType:(ZegoDriveType)driveType {
    UIActivityIndicatorView *loadingView = nil;
    UIButton *button = nil;
    NSString *title = nil;
    
    switch (driveType) {
        case ZegoDriveTypeText:
            loadingView = self.textLoadingView;
            button = self.textDriveButton;
            title = @"文本驱动";
            break;
        case ZegoDriveTypeAudio:
            loadingView = self.audioLoadingView;
            button = self.audioDriveButton;
            title = @"音频驱动";
            break;
        case ZegoDriveTypeWsTTS:
            loadingView = self.wsTTSLoadingView;
            button = self.wsTTSDriveButton;
            title = @"WebSocket TTS驱动";
            break;
        default:
            return;
    }
    
    if (loading) {
        [loadingView startAnimating];
        button.enabled = NO;
        [button setTitle:@"" forState:UIControlStateNormal];
    } else {
        [loadingView stopAnimating];
        button.enabled = YES;
        [button setTitle:title forState:UIControlStateNormal];
    }
}

- (void)setDriveButtonsEnabled:(BOOL)enabled {
    self.textDriveButton.enabled = enabled;
    self.audioDriveButton.enabled = enabled;
    self.wsTTSDriveButton.enabled = enabled;
}

#pragma mark - Actions

- (void)textDriveTapped {
    if ([self.delegate respondsToSelector:@selector(driveControlViewDidTapTextDrive:)]) {
        [self.delegate driveControlViewDidTapTextDrive:self];
    }
}

- (void)audioDriveTapped {
    if ([self.delegate respondsToSelector:@selector(driveControlViewDidTapAudioDrive:)]) {
        [self.delegate driveControlViewDidTapAudioDrive:self];
    }
}

- (void)wsTTSDriveTapped {
    if ([self.delegate respondsToSelector:@selector(driveControlViewDidTapWsTTSDrive:)]) {
        [self.delegate driveControlViewDidTapWsTTSDrive:self];
    }
}

#pragma mark - Memory Management

- (void)dealloc {
    // 清理资源
}

@end
