package tts

import (
	"encoding/base64"
	"fmt"
	"time"

	"zego-digital-human-server/internal/config"
	"zego-digital-human-server/internal/logger"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

const (
	// TTS API端点
	huoshanTTSURL = "https://openspeech.bytedance.com/api/v1/tts"
)

// TTSClient TTS客户端
type TTSClient struct {
	config *config.TTSConfig
	client *resty.Client
}

// NewTTSClient 创建新的TTS客户端
func NewTTSClient(cfg *config.TTSConfig) *TTSClient {
	return &TTSClient{
		config: cfg,
		client: resty.New().SetTimeout(30 * time.Second),
	}
}

// HuoshanRequest 火山TTS请求结构
type HuoshanRequest struct {
	App     AppConfig     `json:"app"`
	User    UserConfig    `json:"user"`
	Audio   AudioConfig   `json:"audio"`
	Request RequestConfig `json:"request"`
}

// AppConfig 应用配置
type AppConfig struct {
	AppID   string `json:"appid"`
	Token   string `json:"token"`
	Cluster string `json:"cluster"`
}

// UserConfig 用户配置
type UserConfig struct {
	UID string `json:"uid"`
}

// AudioConfig 音频配置
type AudioConfig struct {
	VoiceType        string  `json:"voice_type"`
	Encoding         string  `json:"encoding,omitempty"`
	SpeedRatio       float32 `json:"speed_ratio,omitempty"`
	Rate             int     `json:"rate,omitempty"`
	Bitrate          int     `json:"bitrate,omitempty"`
	ExplicitLanguage string  `json:"explicit_language,omitempty"`
	ContextLanguage  string  `json:"context_language,omitempty"`
	LoudnessRatio    float32 `json:"loudness_ratio,omitempty"`
}

// RequestConfig 请求配置
type RequestConfig struct {
	ReqID    string `json:"reqid"`
	Text     string `json:"text"`
	TextType string `json:"text_type,omitempty"`
	Operation string `json:"operation"`
}

// HuoshanResponse 火山TTS响应结构
type HuoshanResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"` // Base64编码的音频数据
}

// ConvertTextToPCM 将文本转换为PCM音频
func (c *TTSClient) ConvertTextToPCM(text string, sampleRate int) ([]byte, error) {
	if text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}

	if c.config == nil {
		return nil, fmt.Errorf("TTS config is nil")
	}

	// 构建请求体
	reqID := uuid.New().String()
	request := HuoshanRequest{
		App: AppConfig{
			AppID:   c.config.AppID,
			Token:   c.config.Token,
			Cluster: c.config.Cluster,
		},
		User: UserConfig{
			UID: "server_user",
		},
		Audio: AudioConfig{
			VoiceType:        c.config.VoiceType,
			Encoding:         "pcm",
			SpeedRatio:       1.0,
			Rate:             sampleRate,
			Bitrate:          160,
			ExplicitLanguage: "zh",
			ContextLanguage:  "zh",
			LoudnessRatio:    1.0,
		},
		Request: RequestConfig{
			ReqID:     reqID,
			Text:      text,
			TextType:  "plain",
			Operation: "query",
		},
	}

	// 发送HTTP POST请求

	var resp HuoshanResponse
	httpResp, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("Bearer;%s", c.config.Token)).
		SetBody(request).
		SetResult(&resp).
		Post(huoshanTTSURL)

	if err != nil {
		logger.LogError(fmt.Sprintf("[TTSClient] HTTP request error: %v", err))
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	if httpResp.StatusCode() != 200 {
		logger.LogError(fmt.Sprintf("[TTSClient] HTTP status error: %d, body: %s", httpResp.StatusCode(), httpResp.String()))
		return nil, fmt.Errorf("HTTP request failed with status %d", httpResp.StatusCode())
	}

	// 打印完整响应用于调试
	logger.LogInfo(fmt.Sprintf("[TTSClient] TTS API response: code=%d, message=%s, data length=%d", resp.Code, resp.Message, len(resp.Data)))

	// 检查响应码
	// 火山引擎TTS API: code=3000表示成功，code=0也可能表示成功
	// 只有当code不是3000且不是0时才认为是错误
	if resp.Code != 0 && resp.Code != 3000 {
		logger.LogError(fmt.Sprintf("[TTSClient] TTS API error: code=%d, message=%s", resp.Code, resp.Message))
		return nil, fmt.Errorf("TTS API error: code=%d, message=%s", resp.Code, resp.Message)
	}
	
	// code=3000且message=Success表示成功
	if resp.Code == 3000 && resp.Message == "Success" {
		logger.LogInfo(fmt.Sprintf("[TTSClient] TTS API success: code=%d, message=%s", resp.Code, resp.Message))
	}

	// 解码Base64音频数据
	if resp.Data == "" {
		return nil, fmt.Errorf("empty audio data in response")
	}

	audioData, err := base64.StdEncoding.DecodeString(resp.Data)
	if err != nil {
		logger.LogError(fmt.Sprintf("[TTSClient] Base64 decode error: %v", err))
		return nil, fmt.Errorf("failed to decode base64 audio data: %w", err)
	}

	if len(audioData) == 0 {
		return nil, fmt.Errorf("decoded audio data is empty")
	}

	logger.LogInfo(fmt.Sprintf("[TTSClient] Successfully converted text to PCM, text length: %d, audio length: %d bytes", len(text), len(audioData)))
	return audioData, nil
}

