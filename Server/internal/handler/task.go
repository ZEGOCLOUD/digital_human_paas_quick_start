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

// CreateDigitalHumanStreamTaskRequest 创建数字人视频流任务请求
type CreateDigitalHumanStreamTaskRequest struct {
	OutputMode int    `json:"OutputMode"`                // 1-大图模式(web), 2-小图模式(mobile)
	UserId     string `json:"UserId" binding:"required"` // 用户ID，必选
}

// createDigitalHumanStreamTaskData ZEGO 创建流任务返回的数据结构
type createDigitalHumanStreamTaskData struct {
	TaskId       string `json:"TaskId,omitempty"`
	RoomId       string `json:"RoomId,omitempty"`
	StreamId     string `json:"StreamId,omitempty"`
	Base64Config string `json:"Base64Config,omitempty"`
	AppId        string `json:"AppId,omitempty"`
	Token        string `json:"Token,omitempty"`
}

// createDigitalHumanStreamTaskAPIResp ZEGO 创建流任务完整响应
type createDigitalHumanStreamTaskAPIResp struct {
	Code      int                               `json:"Code"`
	Message   string                            `json:"Message"`
	RequestId string                            `json:"RequestId,omitempty"`
	Data      *createDigitalHumanStreamTaskData `json:"Data,omitempty"`
	Raw       string                            `json:"-"`
}

// CreateDigitalHumanStreamTask 创建数字人视频流任务
func CreateDigitalHumanStreamTask(c *gin.Context) {
	logger.LogInfo("[CreateDigitalHumanStreamTask] 收到请求------------------")

	//第一步: 参数校验
	bodyParams, code, msg := parseCreateDigitalHumanStreamTaskRequest(c)
	if code != 0 {
		logger.LogErrorf("Error,请求信息校验失败,code:%d,msg:%s", code, msg)
		c.JSON(http.StatusInternalServerError, response.CommonResponse{Code: code, Message: msg, Data: map[string]interface{}{}})
		return
	}
	logger.LogInfof("请求信息,bodyParams:%v", bodyParams)

	//第二步: 配置校验
	zegoConfig, code, msg := ensureZegoConfig()
	if code != 0 {
		logger.LogErrorf("Error,配置信息获取失败,code:%d,msg:%s", code, msg)
		c.JSON(http.StatusInternalServerError, response.CommonResponse{Code: code, Message: msg, Data: map[string]interface{}{}})
		return
	}
	logger.LogInfof("配置信息,zegoConfig:%v", zegoConfig)

	//第三步: 生成房间与流 ID
	roomId, streamId := generateStreamIdentifiers()
	requestBody := buildCreateDigitalHumanStreamTaskBody(bodyParams.OutputMode, roomId, streamId)
	logger.LogInfof("请求Zego Paas信息:requestBody:%v", requestBody)

	//第四步: 调用 ZEGO 创建任务接口
	apiResp, code, msg := callCreateDigitalHumanStreamTaskAPI(zegoConfig, requestBody)
	if code != 0 || apiResp.Code != 0 {
		logger.LogErrorf("Error,调用ZEGO PAAS创建任务接口失败,code:%d,msg:%s", code, msg)
		c.JSON(http.StatusInternalServerError, response.CommonResponse{Code: code, Message: msg, Data: map[string]interface{}{}})
		return
	}
	logger.LogInfof("ZEGO PAAS响应信息,apiResp:%v", apiResp)

	if apiResp.Data == nil {
		apiResp.Data = &createDigitalHumanStreamTaskData{}
	}

	//第五步: 生成 base64ConfigString
	base64ConfigString, base64Code, base64Msg := generateBase64Config(zegoConfig, roomId, streamId, bodyParams.OutputMode)
	if base64Code != 0 {
		logger.LogErrorf("Error,生成base64Config失败,code:%d,msg:%s", base64Code, base64Msg)
		c.JSON(http.StatusInternalServerError, response.CommonResponse{Code: base64Code, Message: base64Msg, Data: map[string]interface{}{}})	
		return
	}
	logger.LogInfof("生成的base64Config:%v", base64ConfigString)

	//第六步: 生成 token
	effectiveTimeInSeconds := int64(3600)
	token, tokenErr := zego.GenerateToken04(
		zegoConfig.AppID,
		bodyParams.UserId,
		zegoConfig.ServerSecret,
		effectiveTimeInSeconds,
		"",
	)
	if tokenErr != nil || token == "" {
		logger.LogError("Error,生成 Token 失败:", tokenErr)
		c.JSON(http.StatusInternalServerError, response.CommonResponse{Code: 500, Message: "生成 Token 失败: " + tokenErr.Error(), Data: map[string]interface{}{}})
		return
	}

	logger.LogInfof("生成的token:%v", token)

	apiResp.Data.Base64Config = base64ConfigString
	apiResp.Data.AppId = strconv.FormatInt(zegoConfig.AppID, 10)
	apiResp.Data.RoomId = roomId
	apiResp.Data.StreamId = streamId
	apiResp.Data.Token = token

	//最后一步 返回结果
	logger.LogInfof("返回结果,apiResp:%v", apiResp)
	c.JSON(http.StatusOK, response.CommonResponse{
		Code:      apiResp.Code,
		Message:   apiResp.Message,
		Data:      apiResp.Data,
		RequestId: apiResp.RequestId,
	})
}

// parseCreateDigitalHumanStreamTaskRequest 校验/解析请求参数
func parseCreateDigitalHumanStreamTaskRequest(c *gin.Context) (CreateDigitalHumanStreamTaskRequest, int, string) {
	var bodyParams CreateDigitalHumanStreamTaskRequest
	if err := c.ShouldBindJSON(&bodyParams); err != nil {
		logger.LogWarn("请求体解析失败:", err)
		return bodyParams, 400, "请求参数错误: " + err.Error()
	}

	if bodyParams.UserId == "" {
		logger.LogWarn("UserId 参数缺失")
		return bodyParams, 400, "UserId 参数必选"
	}

	if bodyParams.OutputMode != 1 && bodyParams.OutputMode != 2 {
		return bodyParams, 400, "OutputMode 参数必须为 1(web) 或 2(mobile)"
	}

	return bodyParams, 0, ""
}

// ensureZegoConfig 校验 ZEGO 配置
func ensureZegoConfig() (*config.ZegoConfig, int, string) {
	zegoConfig := config.GetZegoConfig()
	if !config.ValidateZegoConfig() {
		return zegoConfig, 500, "服务端配置缺失,请检查.env文件配置"
	}
	return zegoConfig, 0, ""
}

// generateStreamIdentifiers 生成房间与流 ID
func generateStreamIdentifiers() (string, string) {
	return zego.GenerateId("test_room_"), zego.GenerateId("stream_")
}

// buildCreateDigitalHumanStreamTaskBody 构造创建任务请求体
func buildCreateDigitalHumanStreamTaskBody(outputMode int, roomId, streamId string) map[string]interface{} {
	requestBody := make(map[string]interface{})
	requestBody["RTCConfig"] = map[string]interface{}{
		"RoomId":   roomId,
		"StreamId": streamId,
	}

	backgroundColor := "#00000000"
	if outputMode == 1 {
		backgroundColor = "#000000"
	}

	requestBody["DigitalHumanConfig"] = map[string]interface{}{
		"DigitalHumanId":  config.GetDefaultDigitalHumanId(),
		"BackgroundColor": backgroundColor,
	}
	requestBody["ExtraConfig"] = map[string]interface{}{
		"OutputMode": outputMode,
	}
	return requestBody
}

// callCreateDigitalHumanStreamTaskAPI 调用 ZEGO 创建任务接口
func callCreateDigitalHumanStreamTaskAPI(zegoConfig *config.ZegoConfig, requestBody map[string]interface{}) (createDigitalHumanStreamTaskAPIResp, int, string) {
	queryString := zego.GenerateQueryParamsString(
		"CreateDigitalHumanStreamTask",
		strconv.FormatInt(zegoConfig.AppID, 10),
		zegoConfig.ServerSecret,
	)
	fullUrl := "https://" + zegoConfig.APIHost + "/?" + queryString

	bodyBytes, _ := json.Marshal(requestBody)

	logger.LogInfo("完整请求URL:", fullUrl)
	logger.LogInfo("POST body :", string(bodyBytes))

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(bodyBytes).
		Post(fullUrl)
	if err != nil {
		return createDigitalHumanStreamTaskAPIResp{}, 500, "调用 ZEGO 创建任务失败: " + err.Error()
	}

	logger.LogInfo("ZEGO原始响应:", resp.String(), "状态码:", resp.StatusCode())

	var apiResp createDigitalHumanStreamTaskAPIResp
	if err := json.Unmarshal(resp.Body(), &apiResp); err != nil {
		return createDigitalHumanStreamTaskAPIResp{}, 500, "解析 ZEGO 响应失败: " + err.Error()
	}

	return apiResp, 0, ""
}

// generateBase64Config 生成 base64
func generateBase64Config(
	zegoConfig *config.ZegoConfig,
	roomId string,
	streamId string,
	outputMode int,
) (string, int, string) {
	if config.GetDefaultDigitalHumanId() == "" || roomId == "" || streamId == "" {
		logger.LogError("生成 Base64Config 前置参数缺失")
		return "", 500, "生成 Base64Config 的参数缺失"
	}

	apiHost := "https://" + zegoConfig.APIHost
	base64Config, err := zego.GetDigitalHumanEncodedConfig(
		zegoConfig.AppID,
		zegoConfig.ServerSecret,
		config.GetDefaultDigitalHumanId(),
		roomId,
		streamId,
		"H264",
		outputMode,
		apiHost,
	)
	if err != nil {
		logger.LogError("生成 Base64Config 失败:", err)
		return "", 500, "生成 Base64Config 失败: " + err.Error()
	}
	return base64Config, 0, ""
}

// base64,
// appid
// token
// roomid,
// taskid,
// streamid,

//getinfo,返回dhid

//tod refactor: 结构化

// QueryDigitalHumanStreamTasks 查询所有运行中的数字人视频流任务
func QueryDigitalHumanStreamTasks(c *gin.Context) {
	logger.LogInfo("[QueryDigitalHumanStreamTasks] 收到请求------------------")

	zegoConfig := config.GetZegoConfig()
	if !config.ValidateZegoConfig() {
		c.JSON(http.StatusInternalServerError, response.CommonResponse{
			Code:    500,
			Message: "服务端配置缺失",
			Data:    map[string]interface{}{},
		})
		return
	}

	queryString := zego.GenerateQueryParamsString(
		"QueryDigitalHumanStreamTasks",
		strconv.FormatInt(zegoConfig.AppID, 10),
		zegoConfig.ServerSecret,
	)
	fullUrl := "https://" + zegoConfig.APIHost + "/?" + queryString

	logger.LogInfo("完整请求URL:", fullUrl)

	// 发送请求
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody([]byte("{}")). // 无需业务参数
		Post(fullUrl)

	if err != nil {
		logger.LogError("请求异常:", err)
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

// StopDigitalHumanStreamTaskRequest 停止数字人视频流任务请求
type StopDigitalHumanStreamTaskRequest struct {
	TaskId string `json:"TaskId" binding:"required"`
}

// StopDigitalHumanStreamTask 停止数字人视频流任务
func StopDigitalHumanStreamTask(c *gin.Context) {
	logger.LogInfo("[StopDigitalHumanStreamTask] 收到请求------------------")

	var req StopDigitalHumanStreamTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.CommonResponse{
			Code:    400,
			Message: "缺少 TaskId 参数",
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

	queryString := zego.GenerateQueryParamsString(
		"StopDigitalHumanStreamTask",
		strconv.FormatInt(zegoConfig.AppID, 10),
		zegoConfig.ServerSecret,
	)
	fullUrl := "https://" + zegoConfig.APIHost + "/?" + queryString

	bodyParams := map[string]interface{}{
		"TaskId": req.TaskId,
	}
	bodyBytes, _ := json.Marshal(bodyParams)

	logger.LogInfo("完整请求URL:", fullUrl)
	logger.LogInfo("POST body:", string(bodyBytes))

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

// InterruptDriveTaskRequest 打断驱动任务请求
type InterruptDriveTaskRequest struct {
	TaskId string `json:"TaskId"`
}

// InterruptDriveTask 打断数字人驱动行为
func InterruptDriveTask(c *gin.Context) {
	logger.LogInfo("[InterruptDriveTask] 收到请求------------------")

	var bodyParams map[string]interface{}
	if err := c.ShouldBindJSON(&bodyParams); err != nil {
		c.JSON(http.StatusBadRequest, response.CommonResponse{
			Code:    400,
			Message: "请求参数错误",
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

	queryString := zego.GenerateQueryParamsString(
		"InterruptDriveTask",
		strconv.FormatInt(zegoConfig.AppID, 10),
		zegoConfig.ServerSecret,
	)
	fullUrl := "https://" + zegoConfig.APIHost + "/?" + queryString

	bodyBytes, _ := json.Marshal(bodyParams)

	logger.LogInfo("完整请求URL:", fullUrl)
	logger.LogInfo("POST body:", string(bodyBytes))

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
