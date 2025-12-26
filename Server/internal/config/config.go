package config

import (
	"os"
	"strconv"
)

// ZegoConfig ZEGO 配置
type ZegoConfig struct {
	AppID          int64
	ServerSecret   string
	APIHost        string
	DefaultDigitalHumanID string
}

var zegoConfig *ZegoConfig

// GetZegoConfig 获取 ZEGO 配置
func GetZegoConfig() *ZegoConfig {
	if zegoConfig == nil {
		zegoConfig = &ZegoConfig{
			AppID:          getAppID(),
			ServerSecret:   getEnv("ZEGO_SERVER_SECRET", ""),
			APIHost:        getEnv("ZEGO_API_HOST", "aigc-digital-human-api.zegotech.cn"),
			DefaultDigitalHumanID: getEnv("DEFAULT_DIGITAL_HUMAN_ID", "your_default_digital_human_id"),
		}
	}
	return zegoConfig
}

// ValidateZegoConfig 验证 ZEGO 配置是否完整
func ValidateZegoConfig() bool {
	config := GetZegoConfig()
	return config.AppID > 0 && config.ServerSecret != ""
}

// GetDefaultDigitalHumanId 获取默认数字人ID
func GetDefaultDigitalHumanId() string {
	return GetZegoConfig().DefaultDigitalHumanID
}

// getAppID 从环境变量获取 AppID
func getAppID() int64 {
	appIDStr := getEnv("NEXT_PUBLIC_ZEGO_APP_ID", "")
	if appIDStr == "" {
		return 0
	}
	appID, err := strconv.ParseInt(appIDStr, 10, 64)
	if err != nil {
		return 0
	}
	return appID
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

