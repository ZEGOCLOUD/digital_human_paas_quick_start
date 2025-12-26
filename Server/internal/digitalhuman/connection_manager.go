package digitalhuman

import (
	"context"
	"fmt"
	"sync"
	"time"

	"zego-digital-human-server/internal/config"
	"zego-digital-human-server/internal/logger"
	"zego-digital-human-server/internal/tts"
)

const (
	// 忙碌占用超时（防止异常未释放 busy）
	staleBusyTimeout = 5 * time.Minute
)

// TaskConnection 任务连接结构
type TaskConnection struct {
	TaskID     string
	TTSClient  *tts.HuoShanTTSWebSocketV3
	WSDriver   *WSDriver
	WSInfo     *DriveByWsStreamingRes
	LastUsed   time.Time
	mu         sync.RWMutex
	responseCh chan *tts.TTSResponse
	busy       bool       // 是否有驱动在执行
	busyMu     sync.Mutex // 保护 busy
}

// ConnectionManager 连接管理器
type ConnectionManager struct {
	connections map[string]*TaskConnection
	mu          sync.RWMutex
}

var globalConnectionManager *ConnectionManager
var once sync.Once

// GetConnectionManager 获取全局连接管理器单例
func GetConnectionManager() *ConnectionManager {
	once.Do(func() {
		globalConnectionManager = &ConnectionManager{
			connections: make(map[string]*TaskConnection),
		}
	})
	return globalConnectionManager
}

// GetOrCreateConnection 获取或创建任务连接
func (cm *ConnectionManager) GetOrCreateConnection(taskId string) (*TaskConnection, error) {
	cm.mu.RLock()
	conn, exists := cm.connections[taskId]
	cm.mu.RUnlock()

	if exists {
		//复用连接，更新最后使用时间
		conn.mu.Lock()
		conn.LastUsed = time.Now()
		conn.mu.Unlock()
		// 检查 Zego WebSocket, 保持长连接
		if !cm.isZegoConnectionValid(conn) {
			logger.LogWarn(fmt.Sprintf("[ConnectionManager] Zego WS invalid for taskId: %s, recreating full connection", taskId))
			cm.removeConnection(taskId)
		} else {
			// 确保 TTS 连接可用；若不可用则重建 TTS
			if err := cm.ensureTTSConnection(taskId, conn); err == nil {
				logger.LogInfo(fmt.Sprintf("[ConnectionManager] ReUsing connection for taskId: %s", taskId))
				// 返回复用连接
				return conn, nil
			}
			logger.LogWarn(fmt.Sprintf("[ConnectionManager] Recreate TTS failed for taskId: %s, recreating full connection", taskId))
			cm.removeConnection(taskId)
		}
	}

	// 创建新连接
	logger.LogInfo(fmt.Sprintf("[ConnectionManager] Creating new connection for taskId: %s", taskId))
	return cm.createConnection(taskId)
}

// createConnection 创建新连接
func (cm *ConnectionManager) createConnection(taskId string) (*TaskConnection, error) {
	ctx := context.Background()

	// 创建TTS客户端
	ttsConfig := config.GetTTSConfig()
	if !config.ValidateTTSConfig() {
		return nil, fmt.Errorf("TTS config is invalid")
	}

	ttsWSClient := tts.NewHuoShanTTSWSV3()
	ttsWSClient.UpdateConfig(ttsConfig)

	// 创建WSDriver
	wsDriver := NewWSDriver()

	// 获取WebSocket连接信息
	wsInfo, err := wsDriver.GetWSInfo(taskId)
	if err != nil {
		return nil, fmt.Errorf("failed to get websocket info: %w", err)
	}

	// 建立数字人WebSocket连接
	if err := wsDriver.Connect(wsInfo.WssAddress); err != nil {
		return nil, fmt.Errorf("failed to connect websocket: %w", err)
	}

	// 建立TTS WebSocket连接
	if err := ttsWSClient.Connect(ctx, "ConnectionManager"); err != nil {
		wsDriver.Close() // 清理已建立的连接
		return nil, fmt.Errorf("failed to connect TTS websocket: %w", err)
	}

	// 创建响应通道
	responseCh := make(chan *tts.TTSResponse, 100)
	ttsWSClient.Setup(responseCh)

	conn := &TaskConnection{
		TaskID:     taskId,
		TTSClient:  ttsWSClient,
		WSDriver:   wsDriver,
		WSInfo:     wsInfo,
		LastUsed:   time.Now(),
		responseCh: responseCh,
	}

	cm.mu.Lock()
	cm.connections[taskId] = conn
	cm.mu.Unlock()

	logger.LogInfo(fmt.Sprintf("[ConnectionManager] Created connection for taskId: %s", taskId))
	return conn, nil
}

// isConnectionValid 检查连接是否有效
func (cm *ConnectionManager) isConnectionValid(conn *TaskConnection) bool {
	return cm.isZegoConnectionValid(conn) && conn.TTSClient != nil && conn.TTSClient.IsConnected()
}

// isZegoConnectionValid 检查数字人WebSocket是否有效
func (cm *ConnectionManager) isZegoConnectionValid(conn *TaskConnection) bool {
	conn.mu.RLock()
	defer conn.mu.RUnlock()
	return conn.WSDriver != nil && conn.WSDriver.conn != nil
}

// ensureTTSConnection 确保TTS连接可用（断开则重建）
func (cm *ConnectionManager) ensureTTSConnection(taskId string, conn *TaskConnection) error {
	conn.mu.RLock()
	ttsClient := conn.TTSClient
	conn.mu.RUnlock()

	if ttsClient != nil && ttsClient.IsConnected() {
		return nil
	}
	return cm.recreateTTS(taskId, conn)
}

// recreateTTS 重建TTS连接，复用或新建响应通道
func (cm *ConnectionManager) recreateTTS(taskId string, conn *TaskConnection) error {
	ctx := context.Background()

	ttsConfig := config.GetTTSConfig()
	if !config.ValidateTTSConfig() {
		return fmt.Errorf("TTS config is invalid")
	}

	newTTS := tts.NewHuoShanTTSWSV3()
	newTTS.UpdateConfig(ttsConfig)

	conn.mu.Lock()
	if conn.responseCh == nil {
		conn.responseCh = make(chan *tts.TTSResponse, 100)
	}
	respCh := conn.responseCh
	conn.mu.Unlock()

	newTTS.Setup(respCh)
	if err := newTTS.Connect(ctx, "ConnectionManagerRecreate"); err != nil {
		return fmt.Errorf("failed to connect TTS websocket: %w", err)
	}

	// 替换旧的TTS client
	conn.mu.Lock()
	old := conn.TTSClient
	conn.TTSClient = newTTS
	conn.mu.Unlock()

	if old != nil {
		_ = old.Close(ctx)
	}
	return nil
}

// removeConnection 移除连接
func (cm *ConnectionManager) removeConnection(taskId string) {
	cm.mu.Lock()
	conn, exists := cm.connections[taskId]
	if exists {
		delete(cm.connections, taskId)
	}
	cm.mu.Unlock()

	if exists {
		cm.closeConnection(conn)
		logger.LogInfo(fmt.Sprintf("[ConnectionManager] Removed connection for taskId: %s", taskId))
	}
}

// closeConnection 关闭连接
func (cm *ConnectionManager) closeConnection(conn *TaskConnection) {
	ctx := context.Background()

	if conn.TTSClient != nil {
		if err := conn.TTSClient.Close(ctx); err != nil {
			logger.LogError(fmt.Sprintf("[ConnectionManager] Close TTS websocket error: %v", err))
		}
	}

	if conn.WSDriver != nil {
		if err := conn.WSDriver.Close(); err != nil {
			logger.LogError(fmt.Sprintf("[ConnectionManager] Close websocket error: %v", err))
		}
	}

	if conn.responseCh != nil {
		close(conn.responseCh)
	}
}

// RemoveConnection 公开方法：移除连接
func (cm *ConnectionManager) RemoveConnection(taskId string) {
	cm.removeConnection(taskId)
}

// CloseAll 关闭所有连接
func (cm *ConnectionManager) CloseAll() {
	cm.mu.Lock()
	connections := make(map[string]*TaskConnection)
	for k, v := range cm.connections {
		connections[k] = v
	}
	cm.connections = make(map[string]*TaskConnection)
	cm.mu.Unlock()

	for _, conn := range connections {
		cm.closeConnection(conn)
	}

	logger.LogInfo("[ConnectionManager] Closed all connections")
}

// GetConnection 获取连接（不创建）
func (cm *ConnectionManager) GetConnection(taskId string) (*TaskConnection, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	conn, exists := cm.connections[taskId]
	return conn, exists
}

// TryLockTask 获取或创建连接并占用任务（等待可用，避免同任务并发驱动）
func (cm *ConnectionManager) TryLockTask(ctx context.Context, taskId string) (*TaskConnection, error) {
	for {
		conn, err := cm.GetOrCreateConnection(taskId)
		if err != nil {
			return nil, err
		}

		if err := cm.LockTask(ctx, conn); err != nil {
			return nil, err
		}
		return conn, nil
	}
}

// UnlockTask 释放任务占用
func (cm *ConnectionManager) UnlockTask(conn *TaskConnection) {
	if conn == nil {
		return
	}

	conn.busyMu.Lock()
	conn.busy = false
	conn.mu.Lock()
	conn.LastUsed = time.Now()
	conn.mu.Unlock()
	conn.busyMu.Unlock()
}

// LockTask 仅对现有连接加忙碌锁，等待可用
func (cm *ConnectionManager) LockTask(ctx context.Context, conn *TaskConnection) error {
	for {
		if cm.tryLockConn(conn) {
			return nil
		}

		logger.LogWarn(fmt.Sprintf("[ConnectionManager] task %s busy, waiting for release", conn.TaskID))
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(50 * time.Millisecond):
		}
	}
}

// tryLockConn 尝试占用任务连接
func (cm *ConnectionManager) tryLockConn(conn *TaskConnection) bool {
	conn.busyMu.Lock()
	defer conn.busyMu.Unlock()

	if conn.busy {
		conn.mu.RLock()
		lastUsed := conn.LastUsed
		conn.mu.RUnlock()

		// 避免异常占用导致永久忙碌
		if time.Since(lastUsed) > staleBusyTimeout {
			logger.LogWarn(fmt.Sprintf("[ConnectionManager] stale busy lock detected for taskId: %s, force releasing", conn.TaskID))
			conn.busy = false
		} else {
			return false
		}
	}

	conn.busy = true
	conn.mu.Lock()
	conn.LastUsed = time.Now()
	conn.mu.Unlock()
	return true
}
