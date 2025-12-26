package digitalhuman

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"zego-digital-human-server/internal/config"
	"zego-digital-human-server/internal/logger"
	"zego-digital-human-server/internal/zego"
	"zego-digital-human-server/pkg/response"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// WSDriver WebSocket驱动管理器
type WSDriver struct {
	conn     *websocket.Conn
	taskId   string
	driveId  string
	sampleRate int
	writeMu  sync.Mutex // 串行化写操作，避免并发写 panic
}

// DriveByWsStreamingRes WebSocket流驱动响应
type DriveByWsStreamingRes struct {
	DriveId    string `json:"DriveId"`
	WssAddress string `json:"WssAddress"`
}

// BaseMessage WebSocket基础消息
type BaseMessage struct {
	SequenceID string                 `json:"SequenceID"`
	Action     string                 `json:"Action"`
	Payload    interface{}            `json:"Payload"`
	DataPoint  map[string]interface{} `json:"DataPoint,omitempty"`
}

// DriveStartPayload 驱动开始负载
type DriveStartPayload struct {
	DriveId    string `json:"DriveId"`
	SampleRate int    `json:"SampleRate"`
}

// DriveEndPayload 驱动结束负载
type DriveEndPayload struct {
	DriveId string `json:"DriveId"`
}

// NewWSDriver 创建新的WebSocket驱动管理器
func NewWSDriver() *WSDriver {
	return &WSDriver{
		sampleRate: 24000, // 默认采样率
	}
}

// GetWSInfo 获取WebSocket连接信息
func (d *WSDriver) GetWSInfo(taskId string) (*DriveByWsStreamingRes, error) {
	zegoConfig := config.GetZegoConfig()
	if !config.ValidateZegoConfig() {
		return nil, fmt.Errorf("ZEGO config is invalid")
	}

	queryString := zego.GenerateQueryParamsString(
		"DriveByWsStream",
		strconv.FormatInt(zegoConfig.AppID, 10),
		zegoConfig.ServerSecret,
	)
	fullUrl := "https://" + zegoConfig.APIHost + "/?" + queryString

	requestBody := map[string]interface{}{
		"TaskId": taskId,
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		Post(fullUrl)

	if err != nil {
		logger.LogError(fmt.Sprintf("[WSDriver] GetWSInfo HTTP error: %v", err))
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		logger.LogError(fmt.Sprintf("[WSDriver] GetWSInfo HTTP status error: %d, body: %s", resp.StatusCode(), resp.String()))
		return nil, fmt.Errorf("HTTP request failed with status %d", resp.StatusCode())
	}

	var apiResp response.CommonResponse
	if err := json.Unmarshal(resp.Body(), &apiResp); err != nil {
		logger.LogError(fmt.Sprintf("[WSDriver] GetWSInfo unmarshal error: %v", err))
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if apiResp.Code != 0 {
		logger.LogError(fmt.Sprintf("[WSDriver] GetWSInfo API error: code=%d, message=%s", apiResp.Code, apiResp.Message))
		return nil, fmt.Errorf("API error: code=%d, message=%s", apiResp.Code, apiResp.Message)
	}

	// 解析Data字段
	dataBytes, err := json.Marshal(apiResp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var wsInfo DriveByWsStreamingRes
	if err := json.Unmarshal(dataBytes, &wsInfo); err != nil {
		logger.LogError(fmt.Sprintf("[WSDriver] GetWSInfo unmarshal data error: %v", err))
		return nil, fmt.Errorf("failed to unmarshal ws info: %w", err)
	}

	d.taskId = taskId
	d.driveId = wsInfo.DriveId
	logger.LogInfo(fmt.Sprintf("[WSDriver] GetWSInfo success, taskId: %s, driveId: %s, wssAddress: %s", taskId, wsInfo.DriveId, wsInfo.WssAddress))
	return &wsInfo, nil
}

// Connect 建立WebSocket连接
func (d *WSDriver) Connect(wssAddress string) error {
	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
		// 跳过TLS证书验证
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	conn, resp, err := dialer.Dial(wssAddress, nil)
	if err != nil {
		if resp != nil {
			logger.LogError(fmt.Sprintf("[WSDriver] Connect error: %v, status: %d", err, resp.StatusCode))
		} else {
			logger.LogError(fmt.Sprintf("[WSDriver] Connect error: %v", err))
		}
		return fmt.Errorf("failed to connect websocket: %w", err)
	}

	d.conn = conn
	logger.LogInfo(fmt.Sprintf("[WSDriver] Connect success, wssAddress: %s", wssAddress))
	return nil
}

// SendStart 发送Start指令
func (d *WSDriver) SendStart(driveId string, sampleRate int) error {
	d.writeMu.Lock()
	defer d.writeMu.Unlock()

	if d.conn == nil {
		return fmt.Errorf("websocket connection is nil")
	}

	d.driveId = driveId
	d.sampleRate = sampleRate

	message := BaseMessage{
		SequenceID: uuid.New().String(),
		Action:     "Start",
		Payload: DriveStartPayload{
			DriveId:    driveId,
			SampleRate: sampleRate,
		},
	}

	if err := d.conn.WriteJSON(message); err != nil {
		logger.LogError(fmt.Sprintf("[WSDriver] SendStart error: %v", err))
		return fmt.Errorf("failed to send start message: %w", err)
	}

	logger.LogInfo(fmt.Sprintf("[WSDriver] SendStart success, driveId: %s, sampleRate: %d", driveId, sampleRate))
	return nil
}

// SendAudio 发送PCM音频数据
func (d *WSDriver) SendAudio(pcmData []byte) error {
	d.writeMu.Lock()
	defer d.writeMu.Unlock()

	if d.conn == nil {
		return fmt.Errorf("websocket connection is nil")
	}

	if len(pcmData) == 0 {
		return nil // 空数据不发送，但不报错
	}

	if err := d.conn.WriteMessage(websocket.BinaryMessage, pcmData); err != nil {
		logger.LogError(fmt.Sprintf("[WSDriver] SendAudio error: %v", err))
		return fmt.Errorf("failed to send audio data: %w", err)
	}

	return nil
}

// SendStop 发送Stop指令
func (d *WSDriver) SendStop(driveId string) error {
	d.writeMu.Lock()
	defer d.writeMu.Unlock()

	if d.conn == nil {
		return fmt.Errorf("websocket connection is nil")
	}

	message := BaseMessage{
		SequenceID: uuid.New().String(),
		Action:     "Stop",
		Payload: DriveEndPayload{
			DriveId: driveId,
		},
	}

	if err := d.conn.WriteJSON(message); err != nil {
		logger.LogError(fmt.Sprintf("[WSDriver] SendStop error: %v", err))
		return fmt.Errorf("failed to send stop message: %w", err)
	}

	logger.LogInfo(fmt.Sprintf("[WSDriver] SendStop success, driveId: %s", driveId))
	return nil
}

// Close 关闭WebSocket连接
func (d *WSDriver) Close() error {
	if d.conn == nil {
		return nil
	}

	err := d.conn.Close()
	d.conn = nil
	if err != nil {
		logger.LogError(fmt.Sprintf("[WSDriver] Close error: %v", err))
		return fmt.Errorf("failed to close websocket: %w", err)
	}

	logger.LogInfo("[WSDriver] Close success")
	return nil
}

// GetDriveId 获取驱动ID
func (d *WSDriver) GetDriveId() string {
	return d.driveId
}

