package tts

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"

	"zego-digital-human-server/internal/config"
	"zego-digital-human-server/internal/logger"
	"zego-digital-human-server/internal/tts/protocols"

	"github.com/google/uuid"
)

// https://www.volcengine.com/docs/6561/1719100
// 整个过程比较简单：
// 1. Connect( 建立 websocket 连接
// 2. Request( 发送请求，文本，语音参数等都包含在这个请求里面
// 3. 解析回包
// 4. Close（ 关闭连接

const (
	// 这个ResourceId表示大模型语音合成，固定是这个字符串，见文档
	HuoShan_V3_ResourceId = "volc.service_type.10029"
	TTS_URL_HUOSHAN_V3    = "wss://openspeech.bytedance.com/api/v3/tts/unidirectional/stream"
	HuoShan_V3_uid        = "zego"
)

// HuoShanTTSWebSocketV3 火山TTS WebSocket V3
type HuoShanTTSWebSocketV3 struct {
	TTSWebSocketConnBase
	retryConnectCount int // 重试连接次数
	appID             string
	accessToken       string
	logID             string
	sampleRate        int
}

// NewHuoShanTTSWSV3 创建新的火山TTS V3客户端
func NewHuoShanTTSWSV3() *HuoShanTTSWebSocketV3 {
	return &HuoShanTTSWebSocketV3{
		sampleRate: 24000, // 默认值
	}
}

// UpdateConfig 更新配置
func (c *HuoShanTTSWebSocketV3) UpdateConfig(cfg *config.TTSConfig) {
	c.config = cfg
	c.appID = cfg.AppID
	c.accessToken = cfg.Token
	c.sampleRate = cfg.SampleRate
	if c.sampleRate == 0 {
		c.sampleRate = 24000
	}
}

// IsConnected 返回当前TTS连接是否已建立
func (c *HuoShanTTSWebSocketV3) IsConnected() bool {
	return c.ttsConnDone.Load() == 2 && c.websocket != nil
}

// Connect 建立WebSocket连接
func (c *HuoShanTTSWebSocketV3) Connect(ctx context.Context, callSource string) error {
	start := time.Now()
	// 如果连接已经建立或者正在创建中
	if !c.ttsConnDone.CompareAndSwap(0, 1) {
		logger.LogInfo(fmt.Sprintf("[HuoShanTTSV3] connect already present, ttsConnDone=%d", c.ttsConnDone.Load()))
		return nil
	}

	// 连接完成后，清理状态
	defer func() {
		logger.LogInfo(fmt.Sprintf("[HuoShanTTSV3] 创建连接耗时 cost: %v", time.Since(start)))
	}()

	c.retryConnectCount += 1 // 重试次数加1
	// 建立WebSocket连接
	header := c.getHuoShanTTSV3Header()

	resp, err := c.TTSWebSocketConnBase.BuildWebSocket(ctx, TTS_URL_HUOSHAN_V3, header)

	defer func() {
		if resp != nil {
			closeErr := resp.Body.Close()
			if closeErr != nil {
				logger.LogError(fmt.Sprintf("[HuoShanTTSV3] resp.Body.Close err:%v", closeErr))
			}
		}
	}()

	if err != nil {
		c.ttsConnDone.Store(0) // 重置连接状态
		if IsContextCanceled(err) {
			return fmt.Errorf("context canceled: %w", err)
		} else {
			// 尝试从HTTP响应中提取更详细的错误信息
			errorMsg := err.Error()
			if resp != nil {
				statusCode := resp.StatusCode
				errorMsg = fmt.Sprintf("HTTP %d: %s", statusCode, errorMsg)
				
				// 尝试读取响应体获取错误详情
				if resp.Body != nil {
					bodyBytes, readErr := io.ReadAll(resp.Body)
					if readErr == nil && len(bodyBytes) > 0 {
						bodyStr := string(bodyBytes)
						// 限制响应体长度，避免日志过长
						if len(bodyStr) > 500 {
							bodyStr = bodyStr[:500] + "...(truncated)"
						}
						errorMsg = fmt.Sprintf("%s, response: %s", errorMsg, bodyStr)
					}
				}
				
				logger.LogError(fmt.Sprintf("[HuoShanTTSV3] connectInternal error: %v, status: %d", err, statusCode))
			} else {
				logger.LogError(fmt.Sprintf("[HuoShanTTSV3] connectInternal error: %v", err))
			}
			return fmt.Errorf("failed to connect to %s: %s", callSource, errorMsg)
		}
	}

	// 获取 logID
	logID := resp.Header.Get("X-Tt-Logid")
	c.logID = logID
	c.ttsConnDone.Store(2) // 设置连接状态为已连接

	logger.LogInfo(fmt.Sprintf("[HuoShanTTSV3] buildTTSConnection logId=%s", logID))

	return nil
}

func (c *HuoShanTTSWebSocketV3) getHuoShanTTSV3Header() http.Header {
	resourceId := HuoShan_V3_ResourceId

	header := http.Header{}
	header.Set("X-Api-App-Key", c.appID)
	header.Set("X-Api-Access-Key", c.accessToken)
	header.Set("X-Api-Resource-Id", resourceId)
	// 随机字符串，非必须
	header.Set("X-Api-Connect-Id", uuid.New().String())
	return header
}

// getHuoShanTTSV3Body 构建请求体
func (c *HuoShanTTSWebSocketV3) getHuoShanTTSV3Body(ttsRequest *TTSRequest) map[string]any {
	bodyMap := make(map[string]any)

	// 初始化 req_params
	reqParams := make(map[string]any)
	bodyMap["req_params"] = reqParams

	// 初始化 audio_params
	audioParams := make(map[string]any)
	reqParams["audio_params"] = audioParams

	// 设置 audio_params
	audioParams["format"] = "pcm"
	audioParams["sample_rate"] = c.sampleRate

	// 设置 speaker，使用配置中的VoiceType
	if c.config != nil && c.config.VoiceType != "" {
		reqParams["speaker"] = c.config.VoiceType
	} else {
		reqParams["speaker"] = "BV001_streaming"
	}

	// 设置 text
	reqParams["text"] = ttsRequest.Text

	// 设置 user
	bodyMap["user"] = map[string]any{
		"uid": HuoShan_V3_uid,
	}

	return bodyMap
}

func (c *HuoShanTTSWebSocketV3) isSupportedSpeed(speed int) bool {
	return speed >= -50 && speed <= 100
}

func (c *HuoShanTTSWebSocketV3) isSupportedEmotionScale(emotion float64) bool {
	return emotion >= 1 && emotion <= 5
}

// SpeedRateToSpeed 将倍速小数转换为语速参数整数
// multiplier: [0.5, 2.0]
// 返回: rate [-50, 100], valid
func (c *HuoShanTTSWebSocketV3) SpeedRateToSpeed(speed float64) (rate int, valid bool) {
	if speed < 0.5 || speed > 2.0 {
		return 0, false
	}

	// 使用 math.Round 进行标准四舍五入
	rate = int(math.Round((speed - 1.0) * 100.0))

	// 额外保险：限制范围
	if rate < -50 {
		rate = -50
	} else if rate > 100 {
		rate = 100
	}

	return rate, true
}

func (c *HuoShanTTSWebSocketV3) sendTTSRequest(ctx context.Context, request *TTSRequest) error {
	requestBody := c.getHuoShanTTSV3Body(request)

	payload, err := json.Marshal(&requestBody)
	if err != nil {
		return fmt.Errorf("marshal request body error: %w", err)
	}

	logger.LogInfo(fmt.Sprintf("[HuoShanTTSV3] sendTTSRequest : %v", requestBody))

	// 发送文本请求
	if err := protocols.FullClientRequest(c.websocket, payload); err != nil {
		logger.LogError(fmt.Sprintf("[HuoShanTTSV3] [round:%d] tts sending error: %s", request.RequestID, err.Error()))
		return fmt.Errorf("failed to send request to %s: %w", c.logID, err)
	}

	logger.LogInfo(fmt.Sprintf("[HuoShanTTSV3] WebSocket response logID: %s, text: %s", c.logID, request.Text))

	return nil
}

// mapErrorCode 映射错误码
func (c *HuoShanTTSWebSocketV3) mapErrorCode(code uint32, payload []byte) error {
	switch code {
	case 55000000:
		return fmt.Errorf("server error: %s", string(payload))
	case 45000000:
		payloadStr := string(payload)
		if strings.Contains(payloadStr, "concurrency") {
			return fmt.Errorf("concurrent limit exceeded")
		}
		if strings.Contains(payloadStr, "speaker") {
			return fmt.Errorf("speaker permission denied")
		}
		return fmt.Errorf("client error: %s", payloadStr)
	default:
		return fmt.Errorf("unknown error code: %d, message: %s", code, string(payload))
	}
}

func (c *HuoShanTTSWebSocketV3) requestInternal(ctx context.Context, request *TTSRequest) error {
	start := time.Now()
	firstAudioPacket := true
	requestID := request.RequestID
	text := request.Text

	err := c.sendTTSRequest(ctx, request)
	if err != nil {
		return err
	}

	// 处理响应
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		msg, err := protocols.ReceiveMessage(c.websocket)
		if err != nil {
			c.writeResponse(&TTSResponse{
				Type:      TTSResponseTypeError,
				RequestID: requestID,
				Error:     fmt.Errorf("failed to read response from %s: %w", c.logID, err),
				ID:        c.logID,
			})
			return fmt.Errorf("failed to read response: %w", err)
		}

		if msg.MsgType == protocols.MsgTypeError {
			err := c.mapErrorCode(msg.ErrorCode, msg.Payload)
			c.writeResponse(&TTSResponse{
				Type:      TTSResponseTypeError,
				RequestID: requestID,
				Error:     err,
				ID:        c.logID,
			})
			return err
		}

		// 1. 先来 MsgType_FullServerResponse && EventType_TTSSentenceStart
		// 2. 然后是 MsgType_AudioOnlyServer
		// 3. 然后是 MsgType_FullServerResponse && EventType_TTSSentenceEnd
		// 4. 然后是 MsgType_FullServerResponse && EventType_SessionFinished，对，没有 Session_start。
		// 5. 如果继续有文本，重复 start-audio-end-finish过程

		if msg.MsgType == protocols.MsgTypeAudioOnlyServer {
			if firstAudioPacket {
				logger.LogInfo(fmt.Sprintf("[HuoShanTTSV3] [requestID:%d] 当前请求首个音频包耗时 cost: %v, text:%s, len(msg.Payload):%d", requestID, time.Since(start), text, len(msg.Payload)))
				firstAudioPacket = false
			}
			ttsResponse := &TTSResponse{
				Type:       TTSResponseTypeAudio,
				RequestID:  requestID,
				AudioData:  msg.Payload,
				SampleRate: c.sampleRate,
				Channel:    1,
				ID:         c.logID,
			}
			c.writeResponse(ttsResponse)
		}

		if msg.MsgType == protocols.MsgTypeFullServerResponse && msg.EventType == protocols.EventType_SessionFinished {
			// 主动发出结束信号，便于驱动方退出循环并释放占用
			c.writeResponse(&TTSResponse{
				Type:      TTSResponseTypeEnd,
				RequestID: requestID,
				ID:        c.logID,
			})
			break
		}
	}
	return nil
}

// Request 发送TTS请求
func (c *HuoShanTTSWebSocketV3) Request(ctx context.Context, request *TTSRequest) error {
	if c.responseCh == nil || request.Text == "" {
		logger.LogInfo("[HuoShanTTSV3] responseCh is nil or request.Text is empty")
		return nil
	}

	c.writeResponse(&TTSResponse{
		RequestID: request.RequestID,
		Type:      TTSResponseTypeStart,
	})

	err := c.requestInternal(ctx, request)
	if err != nil {
		c.writeResponse(&TTSResponse{
			Type:      TTSResponseTypeError,
			RequestID: request.RequestID,
			Error:     err,
		})
	}
	return err
}

// Close 关闭连接
func (c *HuoShanTTSWebSocketV3) Close(ctx context.Context) error {
	if c.ttsConnDone.Load() == 0 {
		logger.LogInfo("[HuoShanTTSV3] WebSocket connection is not established, skipping close")
		return nil
	}
	err := c.TTSWebSocketConnBase.Close()
	if err != nil {
		logger.LogInfo(fmt.Sprintf("[HuoShanTTSV3] failed to close connection: %v", err))
	}
	c.ttsConnDone.Store(0) // 重置连接状态
	return err
}

// Flush 刷新，发送结束响应
func (c *HuoShanTTSWebSocketV3) Flush(ctx context.Context) error {
	c.writeResponse(&TTSResponse{
		Type: TTSResponseTypeEnd,
	})
	return nil
}
