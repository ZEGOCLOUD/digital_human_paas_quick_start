package tts

// TTSResponseType TTS响应类型
type TTSResponseType int

const (
	TTSResponseTypeStart TTSResponseType = iota
	TTSResponseTypeAudio
	TTSResponseTypeMessage
	TTSResponseTypeError
	TTSResponseTypeEnd
)

// TTSRequest TTS请求结构
type TTSRequest struct {
	RequestID uint64
	Text      string // 合成文本
}

// TTSResponse TTS合成响应
type TTSResponse struct {
	Error     error
	RequestID uint64

	ID         string          // 任务ID, 用于给厂商排查问题
	Type       TTSResponseType // 响应类型
	Message    string          // 提示消息
	AudioData  []byte          // 音频数据
	SampleRate int             // 采样率
	Channel    int             // 声道数
	VoiceType  string          // 音色
}

// TTSWebsocketConn TTS WebSocket连接接口
type TTSWebsocketConn interface {
	Connect(ctx interface{}, callSource string) error
	Request(ctx interface{}, request *TTSRequest) error
	Close(ctx interface{}) error
	UpdateConfig(ctx interface{}, cfg interface{})
}

