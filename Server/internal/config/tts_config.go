package config

import (
	"strconv"
)

// TTSConfig 火山TTS配置
type TTSConfig struct {
	AppID      string
	Token      string
	Cluster    string
	VoiceType  string
	SampleRate int
}

var ttsConfig *TTSConfig

// GetTTSConfig 获取 TTS 配置
func GetTTSConfig() *TTSConfig {
	if ttsConfig == nil {
		ttsConfig = &TTSConfig{
			AppID:      getEnv("TTS_BYTEDANCE_APP_ID", ""),
			Token:      getEnv("TTS_BYTEDANCE_TOKEN", ""),
			Cluster:    getEnv("TTS_BYTEDANCE_CLUSTER", ""),
			VoiceType:  getEnv("TTS_BYTEDANCE_VOICE_TYPE", ""),
			SampleRate: getSampleRate(),
		}
	}
	return ttsConfig
}

// ValidateTTSConfig 验证 TTS 配置是否完整
func ValidateTTSConfig() bool {
	config := GetTTSConfig()
	return config.AppID != "" && config.Token != "" && config.Cluster != "" && config.VoiceType != ""
}

// getSampleRate 从环境变量获取采样率，默认24000
func getSampleRate() int {
	sampleRateStr := getEnv("TTS_BYTEDANCE_SAMPLE_RATE", "24000")
	sampleRate, err := strconv.Atoi(sampleRateStr)
	if err != nil {
		return 24000
	}
	return sampleRate
}

