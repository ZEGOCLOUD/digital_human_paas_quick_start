package tts

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"zego-digital-human-server/internal/config"
	"zego-digital-human-server/internal/logger"

	"github.com/gorilla/websocket"
)

var defaultTTSDialer = &websocket.Dialer{
	Proxy:            http.ProxyFromEnvironment,
	HandshakeTimeout: 5 * time.Second,
}

// TTSWebSocketConnBase TTS WebSocket连接基础结构
type TTSWebSocketConnBase struct {
	ttsConnDone atomic.Int32 // 标识连接状态: 0=未连接, 1=连接中, 2=已连接
	responseCh  chan<- *TTSResponse
	websocket   *websocket.Conn
	config      *config.TTSConfig
}

// Setup 设置响应通道
func (c *TTSWebSocketConnBase) Setup(responseCh chan<- *TTSResponse) {
	c.responseCh = responseCh
}

// BuildWebSocket 建立WebSocket连接
func (c *TTSWebSocketConnBase) BuildWebSocket(ctx context.Context, url string, header http.Header) (*http.Response, error) {
	startTime := time.Now()
	conn, rsp, err := defaultTTSDialer.Dial(url, header)

	if err != nil {
		if rsp != nil {
			logger.LogError(fmt.Sprintf("[TTSWebSocketBase] BuildWebSocket error: %v, status: %d", err, rsp.StatusCode))
		} else {
			logger.LogError(fmt.Sprintf("[TTSWebSocketBase] BuildWebSocket error: %v", err))
		}
		if c.responseCh != nil {
			select {
			case c.responseCh <- &TTSResponse{
				Type:  TTSResponseTypeError,
				Error: fmt.Errorf("build WebSocket error: %v", err),
			}:
			default:
			}
		}
		return rsp, err
	}

	c.websocket = conn
	logger.LogInfo(fmt.Sprintf("[TTSWebSocketBase] BuildWebSocket success, url: %s, cost: %v", url, time.Since(startTime)))
	return rsp, nil
}

// writeResponse 写入响应到通道
func (c *TTSWebSocketConnBase) writeResponse(resp *TTSResponse) {
	if c.responseCh == nil {
		return
	}
	select {
	case c.responseCh <- resp:
	default:
		logger.LogWarn("[TTSWebSocketBase] writeResponse: channel is full, dropping response")
	}
}

// Close 关闭WebSocket连接
func (c *TTSWebSocketConnBase) Close() error {
	if c.websocket == nil {
		return nil
	}
	startTime := time.Now()
	err := c.websocket.Close()
	c.websocket = nil
	c.ttsConnDone.Store(0)
	logger.LogInfo(fmt.Sprintf("[TTSWebSocketBase] Close WebSocket, cost: %v", time.Since(startTime)))
	return err
}

// IsContextCanceled 检查错误是否为上下文取消
func IsContextCanceled(err error) bool {
	if err == nil {
		return false
	}
	return err == context.Canceled || err == context.DeadlineExceeded
}

