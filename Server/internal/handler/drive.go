package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"zego-digital-human-server/internal/config"
	"zego-digital-human-server/internal/logger"
	"zego-digital-human-server/internal/zego"
	"zego-digital-human-server/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

// DriveByTextRequest 文本驱动请求
type DriveByTextRequest struct {
	TaskId string `json:"TaskId" binding:"required"`
}

// DriveByText 文本驱动数字人
func DriveByText(c *gin.Context) {
	logger.LogInfo("[DriveByText] 收到请求------------------")

	var req DriveByTextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.CommonResponse{
			Code:    400,
			Message: "TaskId不能为空",
			Data:    map[string]interface{}{},
		})
		return
	}

	zegoConfig := config.GetZegoConfig()
	if !config.ValidateZegoConfig() {
		c.JSON(http.StatusInternalServerError, response.CommonResponse{
			Code:    500,
			Message: "服务端配置缺失",
			Data:    map[string]interface{}{},
		})
		return
	}

	// 使用服务端默认参数，完全忽略客户端传递的Text和TTSConfig
	defaultText := "你好，我是数字人小助手，很高兴为您服务！"
	defaultTTSConfig := map[string]interface{}{
		"TimbreId":   "0f50d026-aeed-4af9-9b31-ada92084a41a",
		"SpeechRate": 0,
		"PitchRate":  0,
		"Volume":     50,
	}

	// 构建请求参数，只使用服务端默认值
	requestParams := map[string]interface{}{
		"TaskId":    req.TaskId,
		"Text":      defaultText,
		"TTSConfig": defaultTTSConfig,
	}

	queryString := zego.GenerateQueryParamsString(
		"DriveByText",
		strconv.FormatInt(zegoConfig.AppID, 10),
		zegoConfig.ServerSecret,
	)
	fullUrl := "https://" + zegoConfig.APIHost + "/?" + queryString

	bodyBytes, _ := json.Marshal(requestParams)

	logger.LogInfo("完整请求URL:", fullUrl)
	logger.LogInfo("POST body (使用服务端默认参数):", string(bodyBytes))

	// 发送请求
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(bodyBytes).
		Post(fullUrl)

	if err != nil {
		logger.LogError("代理异常:", err)
		c.JSON(http.StatusInternalServerError, response.CommonResponse{
			Code:    500,
			Message: "服务端代理请求失败",
			Data:    map[string]interface{}{},
		})
		return
	}

	logger.LogInfo("ZEGO原始响应:", resp.String(), "状态码:", resp.StatusCode())

	var apiData map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &apiData); err != nil {
		apiData = map[string]interface{}{
			"raw": resp.String(),
		}
	}

	c.JSON(resp.StatusCode(), buildCommonResponseFromAPIData(apiData))
}

// DriveByAudioRequest 音频驱动请求
type DriveByAudioRequest struct {
	TaskId string `json:"TaskId" binding:"required"`
}

// DriveByAudio 音频驱动数字人
func DriveByAudio(c *gin.Context) {
	logger.LogInfo("[DriveByAudio] 收到请求------------------")

	var req DriveByAudioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.CommonResponse{
			Code:    400,
			Message: "TaskId不能为空",
			Data:    map[string]interface{}{},
		})
		return
	}

	zegoConfig := config.GetZegoConfig()
	if !config.ValidateZegoConfig() {
		c.JSON(http.StatusInternalServerError, response.CommonResponse{
			Code:    500,
			Message: "服务端配置缺失",
			Data:    map[string]interface{}{},
		})
		return
	}

	// 使用服务端默认参数，完全忽略客户端传递的AudioUrl
	defaultAudioUrl := "https://zego-aigc-test.oss-cn-shanghai.aliyuncs.com/resource_audio/weather.wav"

	// 构建请求参数，只使用服务端默认值
	requestParams := map[string]interface{}{
		"TaskId":   req.TaskId,
		"AudioUrl": defaultAudioUrl,
	}

	queryString := zego.GenerateQueryParamsString(
		"DriveByAudio",
		strconv.FormatInt(zegoConfig.AppID, 10),
		zegoConfig.ServerSecret,
	)
	fullUrl := "https://" + zegoConfig.APIHost + "/?" + queryString

	bodyBytes, _ := json.Marshal(requestParams)

	logger.LogInfo("完整请求URL:", fullUrl)
	logger.LogInfo("POST body (使用服务端默认参数):", string(bodyBytes))

	// 发送请求
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(bodyBytes).
		Post(fullUrl)

	if err != nil {
		logger.LogError("代理异常:", err)
		c.JSON(http.StatusInternalServerError, response.CommonResponse{
			Code:    500,
			Message: "服务端代理请求失败",
			Data:    map[string]interface{}{},
		})
		return
	}

	logger.LogInfo("ZEGO原始响应:", resp.String(), "状态码:", resp.StatusCode())

	var apiData map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &apiData); err != nil {
		apiData = map[string]interface{}{
			"raw": resp.String(),
		}
	}

	c.JSON(resp.StatusCode(), buildCommonResponseFromAPIData(apiData))
}
