//
//  ZegoNetworkManager.m
//  ZegoDigitalHumanQuickStart
//
//  Created by Zego.
//

#import "ZegoNetworkManager.h"
#import "../Utils/ZegoCommonDefines.h"

@interface ZegoNetworkManager () <NSURLSessionDelegate>

@property (nonatomic, strong) NSURLSession *session;

@end

@implementation ZegoNetworkManager

+ (instancetype)sharedManager {
    static ZegoNetworkManager *instance = nil;
    static dispatch_once_t onceToken;
    dispatch_once(&onceToken, ^{
        instance = [[ZegoNetworkManager alloc] init];
    });
    return instance;
}

- (instancetype)init {
    self = [super init];
    if (self) {
        NSURLSessionConfiguration *config = [NSURLSessionConfiguration defaultSessionConfiguration];
        config.timeoutIntervalForRequest = 30.0;
        config.timeoutIntervalForResource = 60.0;
        _session = [NSURLSession sessionWithConfiguration:config delegate:self delegateQueue:[NSOperationQueue mainQueue]];
    }
    return self;
}

- (void)GET:(NSString *)urlString
 parameters:(nullable NSDictionary *)parameters
    headers:(nullable NSDictionary *)headers
    success:(ZegoNetworkSuccessBlock)success
    failure:(ZegoNetworkFailureBlock)failure {
    
    // 边界检查
    if (!urlString || urlString.length == 0) {
        if (failure) {
            NSError *error = [NSError errorWithDomain:ZegoErrorDomain
                                               code:ZegoAppErrorCodeInvalidParameter
                                           userInfo:@{NSLocalizedDescriptionKey: @"URL不能为空"}];
            failure(error);
        }
        return;
    }
    
    // 构建URL
    NSURLComponents *components = [NSURLComponents componentsWithString:urlString];
    if (parameters && parameters.count > 0) {
        NSMutableArray *queryItems = [NSMutableArray array];
        for (NSString *key in parameters) {
            id value = parameters[key];
            NSString *valueString = [NSString stringWithFormat:@"%@", value];
            [queryItems addObject:[NSURLQueryItem queryItemWithName:key value:valueString]];
        }
        components.queryItems = queryItems;
    }
    
    NSURL *url = components.URL;
    if (!url) {
        if (failure) {
            NSError *error = [NSError errorWithDomain:ZegoErrorDomain
                                               code:ZegoAppErrorCodeInvalidParameter
                                           userInfo:@{NSLocalizedDescriptionKey: @"URL格式错误"}];
            failure(error);
        }
        return;
    }
    
    // 创建请求
    NSMutableURLRequest *request = [NSMutableURLRequest requestWithURL:url];
    request.HTTPMethod = @"GET";
    
    // 设置headers
    if (headers) {
        for (NSString *key in headers) {
            [request setValue:headers[key] forHTTPHeaderField:key];
        }
    }
    
    // 发送请求
    [self sendRequest:request success:success failure:failure];
}

- (void)POST:(NSString *)urlString
  parameters:(nullable NSDictionary *)parameters
     headers:(nullable NSDictionary *)headers
     success:(ZegoNetworkSuccessBlock)success
     failure:(ZegoNetworkFailureBlock)failure {
    
    // 边界检查
    if (!urlString || urlString.length == 0) {
        if (failure) {
            NSError *error = [NSError errorWithDomain:ZegoErrorDomain
                                               code:ZegoAppErrorCodeInvalidParameter
                                           userInfo:@{NSLocalizedDescriptionKey: @"URL不能为空"}];
            failure(error);
        }
        return;
    }
    
    NSURL *url = [NSURL URLWithString:urlString];
    if (!url) {
        if (failure) {
            NSError *error = [NSError errorWithDomain:ZegoErrorDomain
                                               code:ZegoAppErrorCodeInvalidParameter
                                           userInfo:@{NSLocalizedDescriptionKey: @"URL格式错误"}];
            failure(error);
        }
        return;
    }
    
    // 创建请求
    NSMutableURLRequest *request = [NSMutableURLRequest requestWithURL:url];
    request.HTTPMethod = @"POST";
    
    // 设置headers
    [request setValue:@"application/json" forHTTPHeaderField:@"Content-Type"];
    if (headers) {
        for (NSString *key in headers) {
            [request setValue:headers[key] forHTTPHeaderField:key];
        }
    }
    
    // 设置body
    if (parameters) {
        NSError *jsonError = nil;
        NSData *jsonData = [NSJSONSerialization dataWithJSONObject:parameters options:0 error:&jsonError];
        if (jsonError) {
            if (failure) {
                failure(jsonError);
            }
            return;
        }
        request.HTTPBody = jsonData;
    }
    
    // 发送请求
    [self sendRequest:request success:success failure:failure];
}

- (void)sendRequest:(NSURLRequest *)request
            success:(ZegoNetworkSuccessBlock)success
            failure:(ZegoNetworkFailureBlock)failure {
    
    __weak typeof(self) weakSelf = self;
    NSURLSessionDataTask *task = [self.session dataTaskWithRequest:request completionHandler:^(NSData * _Nullable data, NSURLResponse * _Nullable response, NSError * _Nullable error) {
        __strong typeof(weakSelf) strongSelf = weakSelf;
        if (!strongSelf) {
            return;
        }
        
        // 检查网络错误
        if (error) {
            if (failure) {
                failure(error);
            }
            return;
        }
        
        // 检查HTTP响应
        NSHTTPURLResponse *httpResponse = (NSHTTPURLResponse *)response;
        if (httpResponse.statusCode < 200 || httpResponse.statusCode >= 300) {
            if (failure) {
                NSError *httpError = [NSError errorWithDomain:ZegoErrorDomain
                                                        code:ZegoAppErrorCodeNetworkError
                                                    userInfo:@{NSLocalizedDescriptionKey: [NSString stringWithFormat:@"HTTP错误: %ld", (long)httpResponse.statusCode]}];
                failure(httpError);
            }
            return;
        }
        
        // 解析JSON
        if (data && data.length > 0) {
            NSError *jsonError = nil;
            id jsonObject = [NSJSONSerialization JSONObjectWithData:data options:0 error:&jsonError];
            if (jsonError) {
                if (failure) {
                    failure(jsonError);
                }
                return;
            }
            
            if ([jsonObject isKindOfClass:[NSDictionary class]]) {
                if (success) {
                    success((NSDictionary *)jsonObject);
                }
            } else {
                if (failure) {
                    NSError *typeError = [NSError errorWithDomain:ZegoErrorDomain
                                                           code:ZegoAppErrorCodeUnknown
                                                       userInfo:@{NSLocalizedDescriptionKey: @"响应数据格式错误"}];
                    failure(typeError);
                }
            }
        } else {
            if (success) {
                success(@{});
            }
        }
    }];
    
    [task resume];
}

- (void)dealloc {
    [_session invalidateAndCancel];
    _session = nil;
}

@end

