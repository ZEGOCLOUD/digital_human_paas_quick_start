package zego

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

const (
	DigitalHumanAPIEndpoint     = "https://aigc-digitalhuman-api.zegotech.cn"
	DigitalHumanConfigID_Mobile = "mobile"
	DigitalHumanConfigID_Web    = "web"

	DigitalHuman_Paas_ErrCode_NotFoundDigitalHuman int = 400000003 // 数字人 paas 错误：数字人不存在
)

// DigitalHumanEncodeStream 流配置
type DigitalHumanEncodeStream struct {
	RoomId                  string `json:"RoomId"`
	StreamId                string `json:"StreamId"`
	EncodeCode              string `json:"EncodeCode"`
	PackageUrl              string `json:"PackageUrl,omitempty"`
	ConfigId                string `json:"ConfigId"`
	IsSupportSmallImageMode bool   `json:"-"` // 是否支持小图模式，使用 json:"-" 标签来忽略该字段
}

// DigitalHumanEncodeConfig 编码配置
type DigitalHumanEncodeConfig struct {
	DigitalHumanId string                     `json:"DigitalHumanId"`
	Streams        []DigitalHumanEncodeStream `json:"Streams"`
}

// GetRenderInfoReq 获取渲染信息请求
type GetRenderInfoReq struct {
	DigitalHumanId string `json:"DigitalHumanId"`
}

// GetRenderInfoRsp 获取渲染信息响应
type GetRenderInfoRsp struct {
	ClientInferencePackageUrl string `json:"ClientInferencePackageUrl"`
	IsSupportSmallImageMode   bool   `json:"IsSupportSmallImageMode"`
}

// CommonResp 通用响应
type CommonResp struct {
	Code      int         `json:"Code"`
	Message   string      `json:"Message"`
	Data      interface{} `json:"Data"`
	RequestId string      `json:"RequestId"`
}

// GetDigitalHumanRenderInfo 获取数字人渲染信息
func GetDigitalHumanRenderInfo(appId int64, serverSecret string, digitalHumanId string, apiHost string) (*GetRenderInfoRsp, error) {
	// 确保 apiHost 包含协议
	if apiHost == "" {
		apiHost = DigitalHumanAPIEndpoint
	} else if len(apiHost) > 0 && apiHost[0:4] != "http" {
		// 如果不包含协议，添加 https://
		apiHost = "https://" + apiHost
	}

	// 生成查询参数
	queryString := GenerateQueryParamsString("GetDigitalHumanRenderInfo", fmt.Sprintf("%d", appId), serverSecret)
	fullUrl := fmt.Sprintf("%s/?%s", apiHost, queryString)

	// 构建请求体
	req := &GetRenderInfoReq{
		DigitalHumanId: digitalHumanId,
	}
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	// 发送请求
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		Post(fullUrl)

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("HTTP error: %s", resp.Status())
	}

	// 解析响应
	var commonResp CommonResp
	if err := json.Unmarshal(resp.Body(), &commonResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}

	if commonResp.Code != 0 {
		if commonResp.Code == DigitalHuman_Paas_ErrCode_NotFoundDigitalHuman {
			return nil, fmt.Errorf("digital human ID not found (code: %d, msg: %s)", commonResp.Code, commonResp.Message)
		}
		return nil, fmt.Errorf("get render info response error (code: %d, msg: %s)", commonResp.Code, commonResp.Message)
	}

	// 将 Data 转换为 GetRenderInfoRsp
	if commonResp.Data == nil {
		return nil, fmt.Errorf("response data is nil")
	}

	var renderInfoRsp GetRenderInfoRsp
	// 尝试类型断言
	if dataMap, ok := commonResp.Data.(map[string]interface{}); ok {
		// 如果是 map，转换为 JSON 再解析
		dataBytes, err := json.Marshal(dataMap)
		if err != nil {
			return nil, fmt.Errorf("marshal response data failed: %w", err)
		}
		if err := json.Unmarshal(dataBytes, &renderInfoRsp); err != nil {
			return nil, fmt.Errorf("unmarshal response data failed: %w", err)
		}
	} else {
		// 否则使用 JSON 序列化/反序列化
		dataBytes, err := json.Marshal(commonResp.Data)
		if err != nil {
			return nil, fmt.Errorf("marshal response data failed: %w", err)
		}
		if err := json.Unmarshal(dataBytes, &renderInfoRsp); err != nil {
			return nil, fmt.Errorf("unmarshal response data failed: %w", err)
		}
	}

	return &renderInfoRsp, nil
}

// EncodeConfig 将配置编码为 base64 字符串
func EncodeConfig(config *DigitalHumanEncodeConfig) (string, error) {
	jsonBytes, err := json.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("marshal digital human config failed: %w", err)
	}
	// 将 JSON 字节转换为 base64 编码的字符串
	encodedConfig := base64.StdEncoding.EncodeToString(jsonBytes)
	return encodedConfig, nil
}

// GetDigitalHumanEncodedConfig 生成数字人配置并编码为 base64 字符串
// outputMode: 1-(web), 2-(mobile)
func GetDigitalHumanEncodedConfig(
	appId int64,
	serverSecret string,
	digitalHumanId string,
	roomId string,
	streamId string,
	encodeCode string,
	outputMode int,
	apiHost string,
) (string, error) {
	// 参数验证
	// appid 验证：必须大于0
	if appId == 0 {
		return "", fmt.Errorf("appid cannot be empty or zero")
	}

	// serverSecret 验证：不能为空或只包含空白字符
	serverSecret = strings.TrimSpace(serverSecret)
	if serverSecret == "" {
		return "", fmt.Errorf("serverSecret cannot be empty")
	}

	// digitalHumanId 验证：不能为空或只包含空白字符
	digitalHumanId = strings.TrimSpace(digitalHumanId)
	if digitalHumanId == "" {
		return "", fmt.Errorf("digitalHumanId cannot be empty")
	}

	// roomId 验证：不能为空或只包含空白字符
	roomId = strings.TrimSpace(roomId)
	if roomId == "" {
		return "", fmt.Errorf("roomId cannot be empty")
	}

	// streamId 验证：不能为空或只包含空白字符
	streamId = strings.TrimSpace(streamId)
	if streamId == "" {
		return "", fmt.Errorf("streamId cannot be empty")
	}

	// encodeCode 验证：如果为空或只包含空白字符，使用默认值
	encodeCode = strings.TrimSpace(encodeCode)
	if encodeCode == "" {
		encodeCode = "H264"
	}

	// outputMode 验证：必须是1或2
	if outputMode != 1 && outputMode != 2 {
		return "", fmt.Errorf("outputMode must be 1 (web) or 2 (mobile), got: %d", outputMode)
	}

	// 获取素材包 RenderInfo
	getRenderInfoRsp, err := GetDigitalHumanRenderInfo(appId, serverSecret, digitalHumanId, apiHost)
	if err != nil {
		return "", fmt.Errorf("get render info error: %w", err)
	}

	// 根据 outputMode 转换为 configId
	// 1: (web), 2: (mobile)
	var configId string
	if outputMode == 1 {
		configId = DigitalHumanConfigID_Web
	} else {
		configId = DigitalHumanConfigID_Mobile
	}

	// 编码配置
	config := &DigitalHumanEncodeConfig{
		DigitalHumanId: digitalHumanId,
		Streams: []DigitalHumanEncodeStream{
			{
				RoomId:                  roomId,
				StreamId:                streamId,
				EncodeCode:              encodeCode,
				ConfigId:                configId,
				PackageUrl:              getRenderInfoRsp.ClientInferencePackageUrl,
				IsSupportSmallImageMode: getRenderInfoRsp.IsSupportSmallImageMode,
			},
		},
	}

	encodedConfig, encodeErr := EncodeConfig(config)
	if encodeErr != nil {
		return "", fmt.Errorf("encode config error: %w", encodeErr)
	}

	return encodedConfig, nil
}


