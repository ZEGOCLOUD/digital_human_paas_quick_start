//
//  ZegoAPIService.m
//  ZegoDigitalHumanModelQuickStart
//
//  Created by Zego.
//

#import "ZegoAPIService.h"
#import "ZegoNetworkManager.h"
#import "ZegoAPIConstants.h"
#import "../Utils/ZegoCommonDefines.h"
#import "../Models/ZegoDigitalHumanInfo.h"

@interface ZegoAPIService ()

@property (nonatomic, copy) NSString *serverURL;


@end

@implementation ZegoAPIService

+ (instancetype)sharedService {
    static ZegoAPIService *instance = nil;
    static dispatch_once_t onceToken;
    dispatch_once(&onceToken, ^{
        instance = [[ZegoAPIService alloc] init];
    });
    return instance;
}

- (instancetype)init {
    self = [super init];
    if (self) {
        _serverURL = kZegoDefaultServerURL;
    }
    return self;
}

- (void)setServerURL:(NSString *)serverURL {
    // 边界检查
    if (serverURL && serverURL.length > 0) {
        _serverURL = [serverURL copy];
    }
}



- (NSDictionary *)buildHeaders {
    return @{
        kZegoHeaderContentType: kZegoContentTypeJSON
    };
}

- (NSString *)buildURL:(NSString *)action {
    // 边界检查
    if (!action || action.length == 0) {
        return self.serverURL;
    }
    
    // 移除serverURL末尾的斜杠
    NSString *baseURL = self.serverURL;
    if ([baseURL hasSuffix:@"/"]) {
        baseURL = [baseURL substringToIndex:baseURL.length - 1];
    }
    
    // 移除action开头的斜杠
    NSString *actionPath = action;
    if ([actionPath hasPrefix:@"/"]) {
        actionPath = [actionPath substringFromIndex:1];
    }
    
    return [NSString stringWithFormat:@"%@/%@", baseURL, actionPath];
}

- (void)handleResponse:(NSDictionary *)response
               success:(ZegoAPISuccessBlock)success
               failure:(ZegoAPIFailureBlock)failure {
    // 边界检查
    if (![response isKindOfClass:[NSDictionary class]]) {
        if (failure) {
            failure(nil, ZegoAppErrorCodeUnknown, @"响应格式错误");
        }
        return;
    }
    
    // 统一解析格式: {Code: 0, Message: "...", Data: {...}}
    NSInteger code = [response[@"Code"] integerValue];
    NSString *message = response[@"Message"];
    id data = response[@"Data"];
    
    // 检查 Code 是否为 0
    if (code != 0) {
        NSString *errorMessage = message ?: @"操作失败";
        if (failure) {
            NSError *error = [NSError errorWithDomain:ZegoErrorDomain
                                               code:code
                                           userInfo:@{NSLocalizedDescriptionKey: errorMessage}];
            failure(error, code, errorMessage);
        }
        return;
    }
    
    // Code 为 0，提取 Data 并传递给 success
    // Data 必须是字典类型
    NSDictionary *resultData = nil;
    if (data) {
        if ([data isKindOfClass:[NSDictionary class]]) {
            resultData = (NSDictionary *)data;
        } else {
            // Data 不是字典类型，格式错误
            if (failure) {
                failure(nil, ZegoAppErrorCodeUnknown, @"响应格式错误：Data 字段必须是字典类型");
            }
            return;
        }
    } else {
        // Data 为 nil，传递空字典
        resultData = @{};
    }
    
    if (success) {
        success(resultData);
    }
}

#pragma mark - Digital Human API

- (void)getDigitalHumanInfo:(NSString *)userId
                     success:(void (^)(ZegoDigitalHumanInfoModel * _Nonnull))success
                     failure:(ZegoAPIFailureBlock)failure {
    // 边界检查
    if (!userId || userId.length == 0) {
        if (failure) {
            failure(nil, ZegoAppErrorCodeInvalidParameter, @"用户ID不能为空");
        }
        return;
    }
    
    NSString *url = [self buildURL:kZegoActionGetDigitalHumanInfo];
    NSDictionary *params = @{@"UserId": userId};
    
    [[ZegoNetworkManager sharedManager] POST:url parameters:params headers:[self buildHeaders] success:^(NSDictionary * _Nullable response) {
        [self handleResponse:response success:^(NSDictionary * _Nullable data) {
            // handleResponse 已统一解析，data 就是 Data 字段的内容
            ZegoDigitalHumanInfoModel *digitalHuman = [[ZegoDigitalHumanInfoModel alloc] initWithDictionary:data];
            if (digitalHuman && success) {
                success(digitalHuman);
            } else if (failure) {
                failure(nil, ZegoAppErrorCodeUnknown, @"解析数字人信息失败");
            }
        } failure:failure];
    } failure:^(NSError * _Nonnull error) {
        if (failure) {
            failure(error, error.code, error.localizedDescription);
        }
    }];
}

#pragma mark - Stream Task API

- (void)createDigitalHumanStreamTask:(NSDictionary *)config
                             success:(void (^)(NSDictionary * _Nonnull))success
                             failure:(ZegoAPIFailureBlock)failure {
    // 边界检查
    if (!config || ![config isKindOfClass:[NSDictionary class]]) {
        if (failure) {
            failure(nil, ZegoAppErrorCodeInvalidParameter, @"配置参数不能为空");
        }
        return;
    }
    
    NSString *url = [self buildURL:kZegoActionCreateDigitalHumanStreamTask];
    
    [[ZegoNetworkManager sharedManager] POST:url parameters:config headers:[self buildHeaders] success:^(NSDictionary * _Nullable response) {
        [self handleResponse:response success:^(NSDictionary * _Nullable data) {
            NSString *taskId = data[@"TaskId"];
            if (taskId && taskId.length > 0) {
                if (success) {
                    success(data);
                }
            } else {
                if (failure) {
                    failure(nil, ZegoAppErrorCodeUnknown, @"创建任务失败：未返回TaskId");
                }
            }
        } failure:failure];
    } failure:^(NSError * _Nonnull error) {
        if (failure) {
            failure(error, error.code, error.localizedDescription);
        }
    }];
}

- (void)stopDigitalHumanStreamTask:(NSString *)taskId
                           success:(ZegoAPISuccessBlock)success
                           failure:(ZegoAPIFailureBlock)failure {
    // 边界检查
    if (!taskId || taskId.length == 0) {
        if (failure) {
            failure(nil, ZegoAppErrorCodeInvalidParameter, @"任务ID不能为空");
        }
        return;
    }
    
    NSString *url = [self buildURL:kZegoActionStopDigitalHumanStreamTask];
    NSDictionary *params = @{@"TaskId": taskId};
    
    [[ZegoNetworkManager sharedManager] POST:url parameters:params headers:[self buildHeaders] success:^(NSDictionary * _Nullable response) {
        [self handleResponse:response success:^(NSDictionary * _Nullable data) {
            // handleResponse 已统一解析，data 就是 Data 字段的内容
            if (success) {
                success(data);
            }
        } failure:failure];
    } failure:^(NSError * _Nonnull error) {
        if (failure) {
            failure(error, error.code, error.localizedDescription);
        }
    }];
}

- (void)queryDigitalHumanStreamTasks:(void (^)(NSArray<ZegoTask *> * _Nonnull))success
                             failure:(ZegoAPIFailureBlock)failure {
    NSString *url = [self buildURL:kZegoActionQueryDigitalHumanStreamTasks];
    
    [[ZegoNetworkManager sharedManager] POST:url parameters:@{} headers:[self buildHeaders] success:^(NSDictionary * _Nullable response) {
        [self handleResponse:response success:^(NSDictionary * _Nullable data) {
            // handleResponse 已统一解析，data 就是 Data 字段的内容
            // Data 格式: {Tasks: [...]} 或 {TaskList: [...]}
            NSArray *tasksArray = data[@"Tasks"] ?: data[@"TaskList"];
            
            // 边界检查
            if (![tasksArray isKindOfClass:[NSArray class]]) {
                if (success) {
                    success(@[]);
                }
                return;
            }
            
            NSMutableArray<ZegoTask *> *tasks = [NSMutableArray array];
            for (id item in tasksArray) {
                if ([item isKindOfClass:[NSDictionary class]]) {
                    ZegoTask *task = [[ZegoTask alloc] initWithDictionary:item];
                    if (task) {
                        [tasks addObject:task];
                    }
                }
            }
            
            if (success) {
                success(tasks);
            }
        } failure:failure];
    } failure:^(NSError * _Nonnull error) {
        if (failure) {
            failure(error, error.code, error.localizedDescription);
        }
    }];
}

#pragma mark - Drive API

- (void)driveByText:(NSString *)taskId
            success:(ZegoAPISuccessBlock)success
            failure:(ZegoAPIFailureBlock)failure {
    // 边界检查
    if (!taskId || taskId.length == 0) {
        if (failure) {
            failure(nil, ZegoAppErrorCodeInvalidParameter, @"任务ID不能为空");
        }
        return;
    }
    
    NSString *url = [self buildURL:kZegoActionDriveByText];
    NSDictionary *params = @{
        @"TaskId": taskId
    };
    
    [[ZegoNetworkManager sharedManager] POST:url parameters:params headers:[self buildHeaders] success:^(NSDictionary * _Nullable response) {
        [self handleResponse:response success:^(NSDictionary * _Nullable data) {
            // handleResponse 已统一解析，data 就是 Data 字段的内容
            if (success) {
                success(data);
            }
        } failure:failure];
    } failure:^(NSError * _Nonnull error) {
        if (failure) {
            failure(error, error.code, error.localizedDescription);
        }
    }];
}

- (void)driveByAudio:(NSString *)taskId
             success:(ZegoAPISuccessBlock)success
             failure:(ZegoAPIFailureBlock)failure {
    // 边界检查
    if (!taskId || taskId.length == 0) {
        if (failure) {
            failure(nil, ZegoAppErrorCodeInvalidParameter, @"任务ID不能为空");
        }
        return;
    }
    
    NSString *url = [self buildURL:kZegoActionDriveByAudio];
    NSDictionary *params = @{
        @"TaskId": taskId
    };
    
    [[ZegoNetworkManager sharedManager] POST:url parameters:params headers:[self buildHeaders] success:^(NSDictionary * _Nullable response) {
        [self handleResponse:response success:^(NSDictionary * _Nullable data) {
            // handleResponse 已统一解析，data 就是 Data 字段的内容
            if (success) {
                success(data);
            }
        } failure:failure];
    } failure:^(NSError * _Nonnull error) {
        if (failure) {
            failure(error, error.code, error.localizedDescription);
        }
    }];
}

- (void)driveByWsStreamWithTTS:(NSString *)taskId
                       success:(ZegoAPISuccessBlock)success
                       failure:(ZegoAPIFailureBlock)failure {
    // 边界检查
    if (!taskId || taskId.length == 0) {
        if (failure) {
            failure(nil, ZegoAppErrorCodeInvalidParameter, @"任务ID不能为空");
        }
        return;
    }
    
    NSString *url = [self buildURL:kZegoActionDriveByWsStreamWithTTS];
    NSDictionary *params = @{
        @"TaskId": taskId
    };
    
    [[ZegoNetworkManager sharedManager] POST:url parameters:params headers:[self buildHeaders] success:^(NSDictionary * _Nullable response) {
        [self handleResponse:response success:^(NSDictionary * _Nullable data) {
            // handleResponse 已统一解析，data 就是 Data 字段的内容
            if (success) {
                success(data);
            }
        } failure:failure];
    } failure:^(NSError * _Nonnull error) {
        if (failure) {
            failure(error, error.code, error.localizedDescription);
        }
    }];
}

- (void)interruptDriveTask:(NSString *)taskId
                   success:(ZegoAPISuccessBlock)success
                   failure:(ZegoAPIFailureBlock)failure {
    // 边界检查
    if (!taskId || taskId.length == 0) {
        if (failure) {
            failure(nil, ZegoAppErrorCodeInvalidParameter, @"任务ID不能为空");
        }
        return;
    }
    
    NSString *url = [self buildURL:kZegoActionInterruptDriveTask];
    NSDictionary *params = @{@"TaskId": taskId};
    
    [[ZegoNetworkManager sharedManager] POST:url parameters:params headers:[self buildHeaders] success:^(NSDictionary * _Nullable response) {
        [self handleResponse:response success:^(NSDictionary * _Nullable data) {
            // handleResponse 已统一解析，data 就是 Data 字段的内容
            if (success) {
                success(data);
            }
        } failure:failure];
    } failure:^(NSError * _Nonnull error) {
        if (failure) {
            failure(error, error.code, error.localizedDescription);
        }
    }];
}

- (void)dealloc {
    // 内存清理
    _serverURL = nil;
}

@end

