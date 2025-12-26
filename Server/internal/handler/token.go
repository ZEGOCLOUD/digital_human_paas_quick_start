package handler

import (
	"net/http"
	"time"

	"zego-digital-human-server/internal/config"
	"zego-digital-human-server/internal/logger"
	"zego-digital-human-server/internal/zego"
	"zego-digital-human-server/pkg/response"

	"github.com/gin-gonic/gin"
)

// ZegoToken 生成 ZEGO Token
func ZegoToken(c *gin.Context) {
	logger.LogInfo("[ZegoToken] 收到请求")
	
	// 获取URL参数
	userId := c.Query("userId")
	
	logger.LogInfo("Request parameters:", map[string]interface{}{
		"url":    c.Request.URL.String(),
		"userId": userId,
	})
	
	// 验证必要参数
	if userId == "" {
		logger.LogWarn("Error: userId is missing")
		c.JSON(http.StatusBadRequest, response.CommonResponse{
			Code:    400,
			Message: "userId is required",
			Data:    map[string]interface{}{},
		})
		return
	}
	
	zegoConfig := config.GetZegoConfig()
	appID := zegoConfig.AppID
	
	logger.LogInfo("AppId:", appID, "ServerSecret:", zegoConfig.ServerSecret)
	
	if !config.ValidateZegoConfig() {
		logger.LogWarn("Error: Server configuration missing:", map[string]interface{}{
			"hasAppID":       appID > 0,
			"hasServerSecret": zegoConfig.ServerSecret != "",
		})
		c.JSON(http.StatusInternalServerError, response.CommonResponse{
			Code:    500,
			Message: "Server configuration error",
			Data:    map[string]interface{}{},
		})
		return
	}
	
	// 设置token有效期（1小时）
	effectiveTimeInSeconds := int64(3600)
	
	logger.LogInfo("Generating token with parameters:", map[string]interface{}{
		"appID":                 appID,
		"userId":                userId,
		"effectiveTimeInSeconds": effectiveTimeInSeconds,
	})
	
	// 生成token
	token, err := zego.GenerateToken04(
		appID,
		userId,
		zegoConfig.ServerSecret,
		effectiveTimeInSeconds,
		"", // payload为空字符串
	)
	
	if err != nil {
		logger.LogError("Token generation error:", err)
		c.JSON(http.StatusInternalServerError, response.CommonResponse{
			Code:    500,
			Message: "Failed to generate token",
			Data:    map[string]interface{}{},
		})
		return
	}
	
	logger.LogInfo("Token generated successfully")
	
	// 返回token，使用统一格式: {Code: 0, Message: "...", Data: {...}}
	// expire_time 应该是当前时间加上有效期的毫秒时间戳
	expireTime := time.Now().UnixMilli() + effectiveTimeInSeconds*1000
	data := map[string]interface{}{
		"token":       token,
		"user_id":     userId,
		"expire_time": expireTime,
	}
	
	responseData := response.CommonResponse{
		Code:    0,
		Message: "Generate token success",
		Data:    data,
	}
	
	logger.LogInfo("Sending response:", map[string]interface{}{
		"hasToken":   token != "",
		"userId":     userId,
		"expireTime": expireTime,
	})
	
	c.JSON(http.StatusOK, responseData)
}

