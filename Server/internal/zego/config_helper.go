package zego

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"zego-digital-human-server/internal/config"
)

// GetZegoConfig 获取 ZEGO 配置
// 返回 appId (string) 和 serverSecret
func GetZegoConfig() (string, string) {
	zegoConfig := config.GetZegoConfig()
	return strconv.FormatInt(zegoConfig.AppID, 10), zegoConfig.ServerSecret
}

// ValidateZegoConfig 验证 ZEGO 配置是否完整
func ValidateZegoConfig(appId string, serverSecret string) bool {
	return appId != "" && serverSecret != ""
}

// GetDefaultDigitalHumanId 获取默认数字人ID
func GetDefaultDigitalHumanId() string {
	return config.GetDefaultDigitalHumanId()
}

// GenerateId 生成唯一ID 
// 规则: prefix + Date.now().toString(36) + Math.random().toString(36).substr(2)
func GenerateId(prefix string) string {
	// 获取当前时间戳（毫秒）
	timestamp := time.Now().UnixMilli()
	
	// 转换为 36 进制字符串
	timestamp36 := toBase36(timestamp)
	
	// 生成随机字符串
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Float64()
	randomStr := toBase36Float(randomNum)
	if len(randomStr) > 2 {
		randomStr = randomStr[2:]
	}
	
	return prefix + timestamp36 + randomStr
}

// toBase36 将数字转换为 36 进制字符串
func toBase36(n int64) string {
	if n == 0 {
		return "0"
	}
	
	chars := "0123456789abcdefghijklmnopqrstuvwxyz"
	result := ""
	
	for n > 0 {
		result = string(chars[n%36]) + result
		n /= 36
	}
	
	return result
}

// toBase36Float 将浮点数转换为 36 进制字符串（模拟 JavaScript 的 toString(36)）
func toBase36Float(f float64) string {
	// JavaScript 的 Math.random().toString(36) 会生成类似 "0.abc123" 的字符串
	// 我们模拟这个过程
	chars := "0123456789abcdefghijklmnopqrstuvwxyz"
	
	// 将浮点数转换为字符串表示
	str := fmt.Sprintf("%.15f", f)
	
	// 提取小数部分并转换为 36 进制
	result := ""
	for i := 2; i < len(str) && len(result) < 10; i++ {
		if str[i] >= '0' && str[i] <= '9' {
			idx := int(str[i] - '0')
			if idx < 36 {
				result += string(chars[idx])
			}
		}
	}
	
	// 如果结果为空，生成一个随机字符串
	if result == "" {
		rand.Seed(time.Now().UnixNano())
		for i := 0; i < 10; i++ {
			result += string(chars[rand.Intn(36)])
		}
	}
	
	return result
}

