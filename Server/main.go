package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"zego-digital-human-server/internal/handler"
	"zego-digital-human-server/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 加载 .env 文件到系统环境变量
	// 如果 .env 文件不存在，不会报错（适合生产环境使用系统环境变量）
	if err := godotenv.Load(); err != nil {
		log.Println("未找到 .env 文件，将使用系统环境变量")
	}
	
	// 设置 Gin 模式
	gin.SetMode(gin.ReleaseMode)
	
	// 创建 Gin 路由
	r := gin.New()
	
	// 添加中间件
	r.Use(ginLogger())
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())
	
	// 注册 API 路由
	api := r.Group("/api")
	{
		// 数字人管理 API
		api.POST("/GetDigitalHumanInfo", handler.GetDigitalHumanInfo)

		// 数字人驱动 API
		api.POST("/DriveByText", handler.DriveByText)
		api.POST("/DriveByAudio", handler.DriveByAudio)
		api.POST("/DriveByWsStreamWithTTS", handler.DriveByWsStreamWithTTS)
		
		// 流任务管理 API
		api.POST("/CreateDigitalHumanStreamTask", handler.CreateDigitalHumanStreamTask)
		api.POST("/QueryDigitalHumanStreamTasks", handler.QueryDigitalHumanStreamTasks)
		api.POST("/StopDigitalHumanStreamTask", handler.StopDigitalHumanStreamTask)
		api.POST("/InterruptDriveTask", handler.InterruptDriveTask)
		
		// 其他 API
		api.GET("/ZegoToken", handler.ZegoToken)
	}
	
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
	
	// 获取端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	
	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
	
	// 在 goroutine 中启动服务器
	go func() {
		logger.LogInfo("服务器启动在端口:", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.LogError("服务器启动失败:", err)
			os.Exit(1)
		}
	}()
	
	// 等待打断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	logger.LogInfo("正在关闭服务器...")
	
	// 优雅关闭，等待最多 5 秒
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		logger.LogError("服务器强制关闭:", err)
		os.Exit(1)
	}
	
	logger.LogInfo("服务器已关闭")
}

// ginLogger 自定义日志中间件
func ginLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		
		// 处理请求
		c.Next()
		
		// 计算延迟
		latency := time.Since(start)
		
		// 构建日志信息
		if raw != "" {
			path = path + "?" + raw
		}
		
		logger.LogInfo(
			"[GIN]",
			"status:", c.Writer.Status(),
			"method:", c.Request.Method,
			"path:", path,
			"latency:", latency,
			"ip:", c.ClientIP(),
		)
	}
}

// corsMiddleware CORS 中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	}
}

