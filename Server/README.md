# ZEGO 数字人快速启动服务器 (Go 版本)

这是一个基于 `Go` 框架构建的数字人快速启动服务器，集成了 ZEGO 数字人 API 服务，提供完整的数字人驱动和管理功能。

## 功能特性

### 数字人管理
- **获取数字人信息** - 获取指定数字人的详细信息

### 数字人驱动
- **文本驱动**       - 通过文本内容驱动数字人说话
- **音频驱动**       - 通过音频文件驱动数字人
- **WebSocket驱动** - 通过websocket数据驱动数字人

### 流任务管理
- **创建流任务** - 创建数字人直播流任务
- **查询流任务** - 查询当前运行的流任务状态
- **停止流任务** - 停止指定的流任务
- **打断驱动任务** - 打断正在执行的驱动任务

### 其他功能
- **Token生成** - 生成ZEGO服务端Token

## 技术栈

- **框架**: Gin 1.10.0
- **语言**: Go 1.21+
- **HTTP 客户端**: resty/v2
- **UUID**: google/uuid

## 安装和运行

### 环境要求
- Go 1.21 或更高版本

### 安装依赖
```bash
go mod download
```

### 环境配置

项目使用 `godotenv` 库自动加载 `.env` 文件中的环境变量到系统环境变量中。

#### 方式一：使用 .env 文件（推荐用于开发环境）

1. 复制 `.env.example` 为 `.env`：
```bash
cp .env.example .env
```

2. 编辑 `.env` 文件，填入实际配置值：
```env
# ZEGO API 配置
ZEGO_API_HOST=aigc-digital-human-api.zegotech.cn
NEXT_PUBLIC_ZEGO_APP_ID=your_app_id
ZEGO_SERVER_SECRET=your_server_secret

# 默认数字人ID
DEFAULT_DIGITAL_HUMAN_ID=7a08e6e1-d6a9-43c9-a99f-8be184468d1b

# 服务器配置
PORT=3000
```

程序启动时会自动加载 `.env` 文件中的变量到系统环境变量中，然后通过 `os.Getenv()` 读取。

#### 方式二：使用系统环境变量（推荐用于生产环境）

在生产环境中，可以直接设置系统环境变量，无需 `.env` 文件：

```bash
export NEXT_PUBLIC_ZEGO_APP_ID=your_app_id
export ZEGO_SERVER_SECRET=your_server_secret
export ZEGO_API_HOST=aigc-digital-human-api.zegotech.cn
export PORT=3000
# ... 其他环境变量
```

或者使用 Docker 的 `-e` 参数或 Kubernetes 的 ConfigMap/Secret。

### 开发环境运行
```bash
# 启动开发服务器
go run main.go

# 或构建后运行
go build -o server
./server
```

### 生产环境部署
```bash
# 构建项目
go build -o server

# 运行
./server
```

## API 接口

### 数字人管理接口

#### 获取数字人信息
```
POST /api/GetDigitalHumanInfo
```
参数：
- `DigitalHumanId` (必填): 数字人ID

### 数字人驱动接口

#### 文本驱动
```
POST /api/DriveByText
```
参数：
- `TaskId` (必填): 任务ID

#### 音频驱动
```
POST /api/DriveByAudio
```
参数：
- `TaskId` (必填): 任务ID

### 流任务管理接口

#### 创建流任务
```
POST /api/CreateDigitalHumanStreamTask
```

#### 查询流任务
```
POST /api/QueryDigitalHumanStreamTasks
```

#### 停止流任务
```
POST /api/StopDigitalHumanStreamTask
```
参数：
- `TaskId` (必填): 任务ID

#### 打断驱动任务
```
POST /api/InterruptDriveTask
```

### 其他接口

#### 生成Token
```
GET /api/ZegoToken?userId=xxx
```
参数：
- `userId` (必填): 用户ID

## 开发指南

### 项目结构
```
ServerGo/
├── main.go                    # 入口文件
├── go.mod                     # 依赖管理
├── internal/
│   ├── config/                # 配置管理
│   ├── handler/               # API 处理器
│   ├── zego/                  # ZEGO 相关工具函数
│   ├── store/                 # 存储层
│   └── logger/                # 日志模块
└── pkg/
    └── response/              # 响应格式
```

## 注意事项

1. **配置安全**: 请妥善保管 ZEGO_APP_ID 和 ZEGO_SERVER_SECRET，不要提交到版本控制系统
2. **API限制**: 注意 ZEGO API 的调用频率限制
3. **错误处理**: 所有API都包含完整的错误处理机制
4. **日志记录**: 生产环境建议配置日志记录

## 许可证

本项目采用 MIT 许可证。

## 支持

如有问题，请联系 ZEGO 技术支持或查看 [ZEGO 官方文档](https://doc.zego.im/)。

