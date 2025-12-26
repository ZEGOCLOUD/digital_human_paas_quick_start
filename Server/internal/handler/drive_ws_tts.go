package handler

import (
	"fmt"
	"net/http"

	"zego-digital-human-server/internal/config"
	"zego-digital-human-server/internal/digitalhuman"
	"zego-digital-human-server/internal/logger"
	"zego-digital-human-server/pkg/response"

	"github.com/gin-gonic/gin"
)

// DriveByWsStreamWithTTSRequest WebSocket TTS驱动请求
type DriveByWsStreamWithTTSRequest struct {
	TaskId string `json:"TaskId" binding:"required"`
}

// DriveByWsStreamWithTTS WebSocket TTS驱动数字人
func DriveByWsStreamWithTTS(c *gin.Context) {
	logger.LogInfo("[DriveByWsStreamWithTTS] 收到请求------------------")

	// 步骤1: 解析并验证请求参数
	// 从HTTP请求体中解析JSON数据，验证TaskId字段是否提供
	var req DriveByWsStreamWithTTSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.CommonResponse{
			Code:    400,
			Message: "TaskId不能为空",
			Data:    map[string]interface{}{},
		})
		return
	}

	// 步骤2: 验证ZEGO服务配置
	// 检查ZEGO相关的配置项是否完整（如AppID、ServerSecret等）
	if !config.ValidateZegoConfig() {
		c.JSON(http.StatusInternalServerError, response.CommonResponse{
			Code:    500,
			Message: "服务端配置缺失",
			Data:    map[string]interface{}{},
		})
		return
	}

	// 步骤3: 验证TTS服务配置
	// 检查TTS（文本转语音）相关的配置项是否完整（如火山引擎TTS配置）
	if !config.ValidateTTSConfig() {
		c.JSON(http.StatusInternalServerError, response.CommonResponse{
			Code:    500,
			Message: "TTS配置缺失，请检查.env文件中的TTS_BYTEDANCE_*配置",
			Data:    map[string]interface{}{},
		})
		return
	}

	// 步骤4: 准备待转换的文本内容
	// 使用服务端默认文本，这里根据业务需求可修改
	defaultText := "你好，我是数字人小助手，很高兴为您服务！"
	logger.LogInfo(fmt.Sprintf("[DriveByWsStreamWithTTS] 使用默认文本，长度: %d", len(defaultText)))

	// 步骤5: 创建驱动处理器实例
	// 初始化驱动处理器，用于管理WebSocket连接和驱动流程（使用全局连接管理器）
	processor, err := digitalhuman.NewDriverProcessor()
	if err != nil {
		logger.LogError(fmt.Sprintf("[DriveByWsStreamWithTTS] 创建驱动处理器失败: %v", err))
		c.JSON(http.StatusInternalServerError, response.CommonResponse{
			Code:    500,
			Message: fmt.Sprintf("创建驱动处理器失败: %v", err),
			Data:    map[string]interface{}{},
		})
		return
	}


	logger.LogInfo(fmt.Sprintf("[DriveByWsStreamWithTTS] 开始执行驱动流程，taskId: %s", req.TaskId))
	// 步骤6: 准备WebSocket连接
	// 根据TaskId获取或创建WebSocket连接，验证连接状态和配置是否就绪
	conn, err := processor.PrepareConnection(req.TaskId)
	if err != nil {
		logger.LogError(fmt.Sprintf("[DriveByWsStreamWithTTS] 连接准备失败: %v", err))
		c.JSON(http.StatusInternalServerError, response.CommonResponse{
			Code:    500,
			Message: fmt.Sprintf("驱动准备失败: %v", err),
			Data: map[string]interface{}{
				"TaskId": req.TaskId,
			},
		})
		return
	}

	// 步骤7: 立即返回成功响应给客户端
	// 在连接准备成功后立即返回，避免客户端长时间等待
	c.JSON(http.StatusOK, response.CommonResponse{
		Code:    0,
		Message: "驱动请求已接收，连接已就绪，开始驱动",
		Data: map[string]interface{}{
			"TaskId": req.TaskId,
		},
	})

	// 步骤8: 异步执行驱动流程
	// 在独立的goroutine中执行实际的驱动处理：
	// - 将文本通过TTS转换为音频流
	// - 通过WebSocket将音频流发送给数字人驱动服务
	// - 驱动数字人进行相应的动作和表情
	go func() {
		if err := processor.Process(conn, defaultText); err != nil {
			logger.LogError(fmt.Sprintf("[DriveByWsStreamWithTTS] 驱动流程失败: %v", err))
		} else {
			logger.LogInfo(fmt.Sprintf("[DriveByWsStreamWithTTS] 驱动流程成功完成，taskId: %s", req.TaskId))
		}
	}()
}

