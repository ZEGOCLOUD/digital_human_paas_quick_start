package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"zego-digital-human-server/internal/config"
	"zego-digital-human-server/internal/logger"
	"zego-digital-human-server/internal/zego"
	"zego-digital-human-server/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

// GetDigitalHumanInfoRequest 获取数字人信息请求
type GetDigitalHumanInfoRequest struct {
	UserId string `json:"UserId" binding:"required"` // 用户ID，必选
}

// GetDigitalHumanInfoData ZEGO API 返回的数字人信息数据结构
type GetDigitalHumanInfoData struct {
	DigitalHumanId string `json:"DigitalHumanId,omitempty"`
	Name           string `json:"Name,omitempty"`
	AvatarUrl      string `json:"AvatarUrl,omitempty"`
	PreviewUrl     string `json:"PreviewUrl,omitempty"`
	IsPublic       bool   `json:"IsPublic,omitempty"`
}

// GetDigitalHumanInfoAPIResp ZEGO API 返回的完整响应结构
type GetDigitalHumanInfoAPIResp struct {
	Code      int                        `json:"Code"`
	Message   string                     `json:"Message"`
	RequestId string                     `json:"RequestId,omitempty"`
	Data      *GetDigitalHumanInfoResponse `json:"Data,omitempty"`
}

// GetDigitalHumanInfoResponse 返回给客户端的数据结构（包含数字人信息和Token）
type GetDigitalHumanInfoResponse struct {
	GetDigitalHumanInfoData
	AppId      int64  `json:"AppId"`      // AppID，用于预下载
	Token      string `json:"Token"`      // Token，用于预下载
	UserId     string `json:"UserId"`     // UserID，用于预下载
	ExpireTime int64  `json:"ExpireTime"` // Token过期时间
}

// GetDigitalHumanInfo 获取数字人信息
func GetDigitalHumanInfo(c *gin.Context) {
	logger.LogInfo("[GetDigitalHumanInfo] 收到请求------------------")

	// 解析请求参数
	var req GetDigitalHumanInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.LogWarn("解析请求体失败:", err)
		c.JSON(http.StatusBadRequest, response.CommonResponse{
			Code:    400,
			Message: "UserId is required",
			Data:    map[string]interface{}{},
		})
		return
	}
	
	userId := req.UserId
	if userId == "" {
		c.JSON(http.StatusBadRequest, response.CommonResponse{
			Code:    400,
			Message: "UserId is required",
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

	// 从配置中获取默认数字人ID
	digitalHumanId := config.GetDefaultDigitalHumanId()
	if digitalHumanId == "" {
		c.JSON(http.StatusInternalServerError, response.CommonResponse{
			Code:    500,
			Message: "服务端未配置默认数字人ID",
			Data:    map[string]interface{}{},
		})
		return
	}

	// 步骤1: 调用ZEGO API获取数字人信息
	// 向ZEGO服务发送HTTP请求，获取指定数字人的详细信息（包括名称、头像、预览图等）
	apiResp, err := fetchDigitalHumanInfoFromZego(zegoConfig, digitalHumanId, c)
	if err != nil {
		// 错误已在 fetchDigitalHumanInfoFromZego 中处理并返回响应
		return
	}

	// 步骤2: 生成用户Token
	// 为客户端传入的userId生成ZEGO Token04，用于后续的预下载和数字人交互操作
	// Token有效期为1小时，过期后需要重新获取
	token, expireTime, err := generateUserToken(zegoConfig, userId)
	if err != nil {
		logger.LogError("生成Token失败:", err)
		c.JSON(http.StatusInternalServerError, response.CommonResponse{
			Code:    500,
			Message: "生成Token失败",
			Data:    map[string]interface{}{},
		})
		return
	}

	// 步骤3: 构建响应数据
	// 合并ZEGO API返回的数字人信息和生成的Token信息，构建完整的响应数据返回给客户端
	responseData := buildGetDigitalHumanInfoResponse(zegoConfig.AppID, apiResp, userId, token, expireTime)

	logger.LogInfo("添加 AppId 和 Token 到返回数据, userId:", userId)

	// 返回统一格式的响应
	c.JSON(http.StatusOK, response.CommonResponse{
		Code:    0,
		Message: "Success",
		Data:    responseData,
	})
}

// fetchDigitalHumanInfoFromZego 调用ZEGO API获取数字人信息
// 构建查询字符串和URL，发送HTTP POST请求，解析响应并返回数字人信息
func fetchDigitalHumanInfoFromZego(zegoConfig *config.ZegoConfig, digitalHumanId string, c *gin.Context) (*GetDigitalHumanInfoAPIResp, error) {
	// 构建查询参数字符串，包含API名称、AppID和签名信息
	queryString := zego.GenerateQueryParamsString(
		"GetDigitalHumanInfo",
		strconv.FormatInt(zegoConfig.AppID, 10),
		zegoConfig.ServerSecret,
	)
	fullUrl := "https://" + zegoConfig.APIHost + "/?" + queryString

	// 构建请求体，包含要查询的数字人ID
	bodyParams := map[string]interface{}{
		"DigitalHumanId": digitalHumanId,
	}
	bodyBytes, _ := json.Marshal(bodyParams)

	logger.LogInfo("完整请求URL:", fullUrl)
	logger.LogInfo("POST body:", string(bodyBytes))

	// 发送HTTP POST请求到ZEGO API
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(bodyBytes).
		Post(fullUrl)

	if err != nil {
		logger.LogError("请求ZEGO获取数字人信息异常:", err)
		c.JSON(http.StatusInternalServerError, response.CommonResponse{
			Code:    500,
			Message: "服务端请求失败",
			Data:    map[string]interface{}{},
		})
		return nil, err
	}

	logger.LogInfo("ZEGO原始响应:", resp.String(), "状态码:", resp.StatusCode())

	// 解析ZEGO API返回的JSON响应
	var apiResp GetDigitalHumanInfoAPIResp
	if err := json.Unmarshal(resp.Body(), &apiResp); err != nil {
		logger.LogError("解析ZEGO API响应失败:", err)
		c.JSON(http.StatusInternalServerError, response.CommonResponse{
			Code:    500,
			Message: "解析响应失败",
			Data:    map[string]interface{}{},
		})
		return nil, err
	}

	// 检查ZEGO API返回的业务状态码，Code不为0表示请求失败
	if apiResp.Code != 0 {
		c.JSON(resp.StatusCode(), response.CommonResponse{
			Code:    apiResp.Code,
			Message: apiResp.Message,
			Data:    map[string]interface{}{},
		})
		return nil, fmt.Errorf("ZEGO API返回错误: Code=%d, Message=%s", apiResp.Code, apiResp.Message)
	}

	return &apiResp, nil
}

// generateUserToken 生成用户Token
// 使用ZEGO Token04算法生成用户认证Token，用于后续的预下载和数字人交互操作
func generateUserToken(zegoConfig *config.ZegoConfig, userId string) (string, int64, error) {
	// 设置token有效期（1小时）
	effectiveTimeInSeconds := int64(3600)

	// 生成Token04，包含AppID、用户ID、服务端密钥和有效期信息
	token, err := zego.GenerateToken04(
		zegoConfig.AppID,
		userId,
		zegoConfig.ServerSecret,
		effectiveTimeInSeconds,
		"", // payload为空字符串
	)

	if err != nil {
		return "", 0, err
	}

	// 计算Token过期时间（当前时间 + 有效期，单位：毫秒）
	expireTime := time.Now().UnixMilli() + effectiveTimeInSeconds*1000

	return token, expireTime, nil
}

// buildGetDigitalHumanInfoResponse 构建响应数据
// 合并ZEGO API返回的数字人信息和生成的Token信息，构建完整的响应数据结构
func buildGetDigitalHumanInfoResponse(appId int64, apiResp *GetDigitalHumanInfoAPIResp, userId, token string, expireTime int64) GetDigitalHumanInfoResponse {
	// 初始化响应数据，设置AppID
	responseData := GetDigitalHumanInfoResponse{
		AppId: appId,
	}

	// 复制ZEGO API返回的数字人信息（名称、头像、预览图等）
	if apiResp.Data != nil {
		responseData.GetDigitalHumanInfoData = apiResp.Data.GetDigitalHumanInfoData
	}

	// 添加Token相关信息（Token、用户ID、过期时间）
	responseData.Token = token
	responseData.UserId = userId
	responseData.ExpireTime = expireTime

	return responseData
}
