//
//  ZegoDigitalHumanPlaceholderView.m
//  ZegoDigitalHumanQuickStart
//
//  Created by Zego.
//

#import "ZegoDigitalHumanPlaceholderView.h"
#import <Masonry/Masonry.h>
#import <SDWebImage/SDWebImage.h>
@interface ZegoDigitalHumanPlaceholderView ()

@property (nonatomic, strong) UIImageView *backgroundImageView;
@property (nonatomic, strong) UIImageView *iconImageView;
@property (nonatomic, strong) UILabel *titleLabel;

@end

@implementation ZegoDigitalHumanPlaceholderView

#pragma mark - Lifecycle

- (instancetype)initWithFrame:(CGRect)frame {
    self = [super initWithFrame:frame];
    if (self) {
        [self setupViews];
    }
    return self;
}

- (instancetype)initWithCoder:(NSCoder *)coder {
    self = [super initWithCoder:coder];
    if (self) {
        [self setupViews];
    }
    return self;
}

#pragma mark - Setup

- (void)setupViews {
    // 创建渐变背景
    self.backgroundImageView = [[UIImageView alloc] init];
    self.backgroundImageView.contentMode = UIViewContentModeScaleAspectFill;
    self.backgroundImageView.clipsToBounds = YES;
    [self addSubview:self.backgroundImageView];
    
    // 创建渐变图层
    CAGradientLayer *gradientLayer = [CAGradientLayer layer];
    gradientLayer.colors = @[
        (id)[UIColor colorWithRed:0.4 green:0.5 blue:0.9 alpha:1.0].CGColor,  // 浅紫色
        (id)[UIColor colorWithRed:0.3 green:0.2 blue:0.7 alpha:1.0].CGColor  // 深紫色
    ];
    gradientLayer.startPoint = CGPointMake(0.5, 0.0);
    gradientLayer.endPoint = CGPointMake(0.5, 1.0);
    gradientLayer.frame = self.bounds;
    [self.backgroundImageView.layer insertSublayer:gradientLayer atIndex:0];
    
    // 图标（全屏显示，aspect_fill方式）
    self.iconImageView = [[UIImageView alloc] init];
    self.iconImageView.contentMode = UIViewContentModeScaleAspectFill; // aspect_fill方式
    self.iconImageView.clipsToBounds = YES;
    self.iconImageView.backgroundColor = [UIColor clearColor];
    self.iconImageView.userInteractionEnabled = NO; // 禁用手势，避免遮挡父view的手势
    [self addSubview:self.iconImageView];
    
    // 标题
    self.titleLabel = [[UILabel alloc] init];
    self.titleLabel.text = @"点击创建任务\n开始体验数字人";
    self.titleLabel.textColor = [UIColor colorWithWhite:0.2 alpha:1.0];
    self.titleLabel.font = [UIFont boldSystemFontOfSize:24];
    self.titleLabel.textAlignment = NSTextAlignmentCenter;
    self.titleLabel.numberOfLines = 0;
    [self addSubview:self.titleLabel];
    
    // 使用 Masonry 布局
    [self.backgroundImageView mas_makeConstraints:^(MASConstraintMaker *make) {
        make.edges.equalTo(self);
    }];
    
    // 图标占据全屏
    [self.iconImageView mas_makeConstraints:^(MASConstraintMaker *make) {
        make.edges.equalTo(self);
    }];
    
    [self.titleLabel mas_makeConstraints:^(MASConstraintMaker *make) {
        make.bottom.equalTo(self).offset(-60);
        make.leading.equalTo(self).offset(20);
        make.trailing.equalTo(self).offset(-20);
    }];
    
    // 设置默认图标
    [self setDefaultIcon];
    
    // 不再为头像添加点击手势，避免遮挡父view的手势
}

- (void)layoutSubviews {
    [super layoutSubviews];
    
    // 更新渐变图层大小
    for (CALayer *layer in self.backgroundImageView.layer.sublayers) {
        if ([layer isKindOfClass:[CAGradientLayer class]]) {
            layer.frame = self.backgroundImageView.bounds;
        }
    }
}

#pragma mark - Public Methods

- (void)updateWithName:(NSString *)name coverUrl:(nullable NSString *)coverUrl {
    // 边界检查
    if (!name || name.length == 0) {
        name = @"请选择数字人";
    }
    
    self.titleLabel.text = name;
    
    // 加载图标
    UIImage *placeholderImage = [self defaultPlaceholderImage];
    if (coverUrl && coverUrl.length > 0) {
        // 使用 SDWebImage 加载网络图片
        NSURL *url = [NSURL URLWithString:coverUrl];
        if (url) {
            [self.iconImageView sd_setImageWithURL:url
                                  placeholderImage:placeholderImage
                                           options:SDWebImageRetryFailed | SDWebImageHighPriority];
        } else {
            [self setDefaultIcon];
        }
    } else {
        [self setDefaultIcon];
    }
}

- (void)show {
    self.hidden = NO;
    self.alpha = 0;
    [UIView animateWithDuration:0.3 animations:^{
        self.alpha = 1.0;
    }];
}

- (void)hide {
    [UIView animateWithDuration:0.3 animations:^{
        self.alpha = 0;
    } completion:^(BOOL finished) {
        if (finished) {
            self.hidden = YES;
        }
    }];
}

#pragma mark - Private Methods

- (void)setDefaultIcon {
    UIImage *placeholderImage = [self defaultPlaceholderImage];
    if (placeholderImage) {
        self.iconImageView.image = placeholderImage;
    }
}

- (UIImage *)defaultPlaceholderImage {
    // 从 Resources 目录加载占位图片
    // 首先尝试使用 imageNamed（如果图片已添加到 bundle）
    UIImage *image = [UIImage imageNamed:@"image_placeholder"];
    if (image) {
        return image;
    }
    
    // 如果 imageNamed 失败，尝试从指定路径加载
    NSString *imagePath = [[NSBundle mainBundle] pathForResource:@"image_placeholder" ofType:@"png"];
    if (imagePath && [[NSFileManager defaultManager] fileExistsAtPath:imagePath]) {
        image = [UIImage imageWithContentsOfFile:imagePath];
        if (image) {
            return image;
        }
    }
    
    // 如果加载失败，返回 nil（保持当前状态）
    return nil;
}

@end

