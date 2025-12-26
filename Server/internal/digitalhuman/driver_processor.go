package digitalhuman

import (
	"context"
	"fmt"
	"time"

	"zego-digital-human-server/internal/config"
	"zego-digital-human-server/internal/logger"
	"zego-digital-human-server/internal/tts"
)

const (
	// 音频数据分块大小（字节）
	audioChunkSize = 2048
	// 请求超时时间（秒）
	requestTimeout = 180 // 3分钟
)

// DriverProcessor 驱动处理器
type DriverProcessor struct {
	connManager *ConnectionManager
}

// NewDriverProcessor 创建新的驱动处理器
func NewDriverProcessor() (*DriverProcessor, error) {
	connManager := GetConnectionManager()
	return &DriverProcessor{
		connManager: connManager,
	}, nil
}

// PrepareConnection 仅做连接准备，用于提前验证WS/TTS是否就绪
func (p *DriverProcessor) PrepareConnection(taskId string) (*TaskConnection, error) {
	logger.LogInfo(fmt.Sprintf("[DriverProcessor] PrepareConnection for taskId: %s", taskId))
	conn, err := p.connManager.GetOrCreateConnection(taskId)
	if err != nil {
		logger.LogError(fmt.Sprintf("[DriverProcessor] PrepareConnection failed: %v", err))
		return nil, fmt.Errorf("failed to get or create connection: %w", err)
	}
	return conn, nil
}

// Process 使用已存在的连接执行驱动流程（驱动前会加 busy 锁）
func (p *DriverProcessor) Process(conn *TaskConnection, text string) error {
	if conn == nil {
		return fmt.Errorf("connection is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout*time.Second)
	defer cancel()

	// 驱动前再加 busy 锁，避免并发执行
	if err := p.connManager.LockTask(ctx, conn); err != nil {
		return fmt.Errorf("failed to lock task: %w", err)
	}
	defer p.connManager.UnlockTask(conn)
	return p.runProcess(ctx, conn, conn.TaskID, text)
}

// runProcess 执行核心驱动流程，要求传入的连接已被占用
func (p *DriverProcessor) runProcess(ctx context.Context, conn *TaskConnection, taskId string, text string) error {
	if conn == nil {
		return fmt.Errorf("connection is nil")
	}

	// 更新最后使用时间
	conn.mu.Lock()
	conn.LastUsed = time.Now()
	conn.mu.Unlock()

	// 步骤2: 发送Start指令到数字人WebSocket
	sampleRate := config.GetTTSConfig().SampleRate
	if sampleRate == 0 {
		sampleRate = 24000
	}
	logger.LogInfo(fmt.Sprintf("[DriverProcessor] Step 2: SendStart, driveId: %s, sampleRate: %d", conn.WSInfo.DriveId, sampleRate))
	if err := conn.WSDriver.SendStart(conn.WSInfo.DriveId, sampleRate); err != nil {
		logger.LogError(fmt.Sprintf("[DriverProcessor] SendStart failed: %v", err))
		// 连接可能已断开，移除并重试
		p.connManager.RemoveConnection(taskId)
		return fmt.Errorf("failed to send start: %w", err)
	}

	// 步骤3: 发送TTS请求
	logger.LogInfo(fmt.Sprintf("[DriverProcessor] Step 3: Send TTS request, text length: %d", len(text)))
	ttsRequest := &tts.TTSRequest{
		RequestID: uint64(time.Now().UnixNano()), // 使用时间戳作为请求ID
		Text:      text,
	}

	// 在goroutine中处理TTS请求
	ttsErrCh := make(chan error, 1)
	go func() {
		ttsErrCh <- conn.TTSClient.Request(ctx, ttsRequest)
	}()

	// 步骤4: 实时接收音频并转发到数字人WebSocket
	logger.LogInfo("[DriverProcessor] Step 4: Start streaming audio")
	audioReceived := false
	ttsRequestCompleted := false
	responseEndReceived := false

	for {
		select {
		case <-ctx.Done():
			logger.LogError(fmt.Sprintf("[DriverProcessor] Context timeout or canceled: %v", ctx.Err()))
			p.connManager.RemoveConnection(taskId)
			return fmt.Errorf("request timeout: %w", ctx.Err())

		case err := <-ttsErrCh:
			if err != nil {
				logger.LogError(fmt.Sprintf("[DriverProcessor] TTS request failed: %v", err))
				p.connManager.RemoveConnection(taskId)
				return fmt.Errorf("TTS request failed: %w", err)
			}
			// TTS请求完成
			ttsRequestCompleted = true
			logger.LogInfo("[DriverProcessor] TTS request completed")
			// 如果已经收到End响应，可以退出
			if responseEndReceived {
				goto finish
			}

		case response, ok := <-conn.responseCh:
			if !ok {
				// 通道关闭
				logger.LogInfo("[DriverProcessor] Response channel closed")
				if !audioReceived {
					return fmt.Errorf("no audio data received")
				}
				goto finish
			}

			if response == nil {
				continue
			}

			switch response.Type {
			case tts.TTSResponseTypeStart:
				logger.LogInfo("[DriverProcessor] TTS response: Start")

			case tts.TTSResponseTypeAudio:
				audioReceived = true
				// 实时转发音频数据到数字人WebSocket
				if len(response.AudioData) > 0 {
					if err := conn.WSDriver.SendAudio(response.AudioData); err != nil {
						logger.LogError(fmt.Sprintf("[DriverProcessor] SendAudio failed: %v", err))
						// 连接可能已断开，移除连接
						p.connManager.RemoveConnection(taskId)
						return fmt.Errorf("failed to send audio: %w", err)
					}
				}

			case tts.TTSResponseTypeError:
				logger.LogError(fmt.Sprintf("[DriverProcessor] TTS response error: %v", response.Error))
				p.connManager.RemoveConnection(taskId)
				return fmt.Errorf("TTS error: %w", response.Error)

			case tts.TTSResponseTypeEnd:
				logger.LogInfo("[DriverProcessor] TTS response: End")
				responseEndReceived = true
				// 如果TTS请求也已完成，可以退出
				if ttsRequestCompleted {
					goto finish
				}
			}
		}
	}

finish:
	// 步骤5: 发送Stop指令
	logger.LogInfo(fmt.Sprintf("[DriverProcessor] Step 5: SendStop, driveId: %s", conn.WSInfo.DriveId))
	if err := conn.WSDriver.SendStop(conn.WSInfo.DriveId); err != nil {
		logger.LogError(fmt.Sprintf("[DriverProcessor] SendStop failed: %v", err))
		// 连接可能已断开，移除连接
		p.connManager.RemoveConnection(taskId)
		return fmt.Errorf("failed to send stop: %w", err)
	}

	logger.LogInfo("[DriverProcessor] Process completed successfully")
	return nil
}
