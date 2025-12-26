# ZEGO 数字人 PaaS 客户端快速开始

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-Android%20%7C%20iOS%20%7C%20Web-lightgrey.svg)]()

一个完整的 ZEGO 数字人 PaaS 客户端快速启动示例项目，支持 Android、iOS、Web,Server(Go) 四个平台，提供数字人渲染、多种驱动方式、任务管理等完整功能。

## 项目简介

本项目是 ZEGO 数字人 PaaS 服务的客户端快速启动示例，帮助开发者快速集成数字人功能到自己的应用中。项目包含客户端应用（Android、iOS、Web）和服务端示例（Go），展示了完整的数字人集成方案。

### 核心功能

- ✅ **数字人渲染** - 基于 ZEGO DigitalMobile SDK 实现高质量数字人渲染
- ✅ **多种驱动方式** - 支持文本驱动、音频驱动、WebSocket驱动
- ✅ **任务管理** - 完整的流任务生命周期管理（创建、查询、停止、中断）
- ✅ **配置管理** - 灵活的数字人、音色、输出模式配置
- ✅ **跨平台支持** - Android、iOS、Web 三端统一 API 设计

## 功能特性

### 数字人管理
- **获取数字人列表** - 支持分页查询数字人形象(服务端)
- **获取数字人信息** - 获取指定数字人的详细信息
- **获取音色列表** - 获取可用的音色选项(服务端)

### 数字人驱动
- **文本驱动** - 通过文本内容驱动数字人说话，支持语速、语调、音量调节
- **音频驱动** - 通过音频文件驱动数字人
- **WebSocket TTS 驱动** - 通过 WebSocket 实时 TTS 驱动数字人

### 流任务管理
- **创建流任务** - 创建数字人直播流任务
- **查询流任务** - 查询当前运行的流任务状态
- **停止流任务** - 停止指定的流任务
- **打断驱动任务** - 打断正在执行的驱动任务

### 其他功能
- **Token 生成** - 服务端 Token 生成接口
- **配置管理** - 支持 AppID、ServerSecret、数字人、音色等配置

## 技术栈

### Android
- **语言**: Java 11
- **SDK**: 
  - `im.zego:digitalmobile:1.4.0.88` - 数字人渲染 SDK(最新版本请联系技术支持)
  - `im.zego:express-private:3.22.0.46522` - RTC 引擎(最新版本请联系技术支持)


### iOS
- **语言**: Objective-C
- **SDK**: 
  - `ZegoDigitalMobile` - 数字人渲染 SDK
  - `ZegoExpressEngine` - RTC 引擎
- **依赖管理**: CocoaPods


### Web
- **框架**: Vue 3.4.0
- **构建工具**: Vite 5.0.0
- **SDK**: `zego-express-engine-webrtc:3.9.123`(最新版本请联系技术支持)

### Server
- **语言**: Go 1.21+
- **框架**: Gin 1.10.0
- **HTTP 客户端**: resty/v2
- **环境配置**: godotenv

## 📁 项目结构

```
digital_human_paas_client_quick_start/
├── Android/                    # Android 客户端
│   ├── app/
│   │   ├── src/main/
│   │   │   ├── java/com/example/zegodigitalhumanquickstart/
│   │   │   │   ├── model/          # 数据模型
│   │   │   │   ├── callback/        # 回调接口
│   │   │   │   ├── util/           # 工具类
│   │   │   │   └── ...             # 其他业务代码
│   │   │   └── res/                # 资源文件
│   │   └── build.gradle
│   └── build.gradle
│
├── iOS/                        # iOS 客户端
│   ├── ZegoDigitalHumanQuickStart/
│   │   ├── Models/              # 数据模型
│   │   ├── Network/             # 网络层
│   │   ├── Views/               # 视图组件
│   │   ├── ViewControllers/     # 视图控制器
│   │   └── Utils/               # 工具类
│   └── docs/                    # iOS 文档
│
├── Web/                        # Web 客户端
│   ├── src/
│   │   ├── components/          # Vue 组件
│   │   ├── config/              # 配置文件
│   │   └── utils/               # 工具函数
│   ├── package.json
│   └── vite.config.js
│
├── Server/                     # 服务端示例（Go）
│   ├── internal/
│   │   ├── config/              # 配置管理
│   │   ├── handler/             # API 处理器
│   │   ├── digitalhuman/        # 数字人业务逻辑
│   │   ├── tts/                 # TTS 服务
│   │   ├── zego/                # ZEGO 工具函数
│   │   └── logger/              # 日志模块
│   ├── pkg/response/            # 响应格式
│   ├── main.go                  # 入口文件
│   └── go.mod
│
└── README.md                   # 项目说明文档
```

## 快速开始

### 前置要求

1. **ZEGO 账号** - 需要有效的 ZEGO AppID 和 ServerSecret
2. **服务端运行** - 需要先启动服务端（Server 目录）

### 1. 启动服务端

```bash
cd Server

# 安装依赖
go mod download

# 配置环境变量（创建 .env 文件或设置系统环境变量）
# ZEGO_API_HOST=aigc-digital-human-api.zegotech.cn
# NEXT_PUBLIC_ZEGO_APP_ID=your_app_id
# ZEGO_SERVER_SECRET=your_server_secret
# PORT=3000

# 运行服务
go run main.go
```

详细配置请参考 [Server/README.md](Server/README.md)

### 2. Android 客户端

#### 环境要求
- Android Studio Hedgehog | 2023.1.1 或更高版本
- JDK 11 或更高版本
- Android SDK API 26+ (Android 8.0+)
- Gradle 8.0+

#### 运行步骤

```bash
cd Android

# 使用 Android Studio 打开项目
# 或使用命令行构建
./gradlew assembleDebug

# 安装到设备
./gradlew installDebug
```


### 3. iOS 客户端

#### 环境要求
- Xcode 12.0+
- iOS 11.0+
- CocoaPods 1.10.0+

#### 运行步骤

```bash
cd iOS/ZegoDigitalHumanQuickStart

# 安装依赖
pod install

# 打开工作空间
open ZegoDigitalHumanQuickStart.xcworkspace

# 在 Xcode 中配置签名并运行
```


### 4. Web 客户端

#### 环境要求
- Node.js 16.0+ 或更高版本
- npm 或 yarn 或 pnpm

#### 运行步骤

```bash
cd Web

# 安装依赖
npm install
# 或
yarn install
# 或
pnpm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build
```

## API 文档

### 服务端 API

服务端提供以下主要 API 接口：

#### 数字人管理
- `POST /api/GetDigitalHumanList` - 获取数字人列表
- `POST /api/GetDigitalHumanInfo` - 获取数字人信息
- `POST /api/GetTimbreList` - 获取音色列表

#### 数字人驱动
- `POST /api/DriveByText` - 文本驱动
- `POST /api/DriveByAudio` - 音频驱动
- `POST /api/DriveByRTCStream` - RTC 流驱动
- `POST /api/DriveByWsStreamWithTTS` - WebSocket TTS 驱动
- `POST /api/DoAction` - 执行动作

#### 流任务管理
- `POST /api/CreateDigitalHumanStreamTask` - 创建流任务
- `POST /api/QueryDigitalHumanStreamTasks` - 查询流任务
- `POST /api/StopDigitalHumanStreamTask` - 停止流任务
- `POST /api/InterruptDriveTask` - 打断驱动任务

#### 其他
- `GET /api/ZegoToken?userId=xxx` - 生成 Token


详细的 API 文档请参考 [Server/README.md](Server/README.md)

## 配置说明

### 环境变量

服务端需要配置以下环境变量：

```env
# ZEGO API 配置
ZEGO_API_HOST=aigc-digital-human-api.zegotech.cn
NEXT_PUBLIC_ZEGO_APP_ID=your_app_id
ZEGO_SERVER_SECRET=your_server_secret

# 默认数字人ID（可选）
DEFAULT_DIGITAL_HUMAN_ID=your_digital_human_id

# 服务器配置
PORT=3000
```

### 客户端配置

各平台客户端需要在应用内配置：
- **ServerURL**: 服务端地址（如：`http://你的服务端地址ip:3000`）
- **OutputMode**: 输出模式

## 🔧 开发指南

### 代码规范

- **边界检查**: 所有参数都需要进行边界检查和空值检查
- **内存安全**: 注意内存泄漏，及时释放资源
- **错误处理**: 完整的错误处理机制
- **代码注释**: 关键逻辑添加必要注释

### 平台特定说明

#### Android
- 最低支持 Android 8.0 (API 26)
- 使用 Java 11 语法
- 注意 Activity 生命周期管理

#### iOS
- 最低支持 iOS 12.0
- 使用 Objective-C 开发
- 注意内存管理和线程安全
- RTC 驱动需要麦克风权限

#### Web
- 使用 Vue 3 Composition API
- 支持现代浏览器（Chrome、Firefox、Safari、Edge）
- 注意 RTC 权限处理

## 常见问题

### 1. Token 获取失败
- 检查 ServerURL 是否正确
- 检查 AppID 和 ServerSecret 是否正确
- 检查服务端是否正常运行

### 2. 数字人不显示
- 检查任务是否创建成功
- 检查流是否正常拉取
- 检查 SDK 是否正确集成
- 检查网络连接是否正常

### 3. 服务端启动失败
- 检查 Go 版本是否符合要求（1.21+）
- 检查环境变量是否配置正确
- 检查端口是否被占用

## 📝 更新日志

### v1.0.0
- ✅ 支持 Android、iOS、Web 三端
- ✅ 完整的数字人驱动功能
- ✅ 服务端示例（Go）
- ✅ 完整的任务管理功能

## 贡献

欢迎提交 Issue 和 Pull Request！

### 贡献指南

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'feat: Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request



##  许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

## 支持与反馈

- **ZEGO 官方文档**: [https://doc.zego.im/](https://doc.zego.im/)
- **问题反馈**: 请在 GitHub Issues 中提交
- **技术支持**: 联系 ZEGO 技术支持团队


**注意**: 本项目仅供学习和参考使用。在生产环境中使用前，请确保：
1. 妥善保管 AppID 和 ServerSecret，不要提交到版本控制系统
2. 注意 API 调用频率限制
3. 配置适当的日志记录和监控
4. 进行充分的安全测试
