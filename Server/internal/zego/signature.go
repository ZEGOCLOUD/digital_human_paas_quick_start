package zego

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// GenerateSignature 生成 ZEGO API 签名
// Signature=md5(AppId + SignatureNonce + ServerSecret + Timestamp)
func GenerateSignature(appId string, signatureNonce string, serverSecret string, timestamp int64) string {
	str := fmt.Sprintf("%s%s%s%d", appId, signatureNonce, serverSecret, timestamp)
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// GenerateSignatureNonce 生成16位十六进制随机字符串
func GenerateSignatureNonce() string {
	nonceByte := make([]byte, 8)
	_, err := rand.Read(nonceByte)
	if err != nil {
		// 如果随机数生成失败，使用时间戳作为后备
		return fmt.Sprintf("%016x", time.Now().UnixNano())
	}
	return hex.EncodeToString(nonceByte)
}

// GenerateQueryParams 生成查询参数
func GenerateQueryParams(action string, appId string, serverSecret string) map[string]string {
	signatureNonce := GenerateSignatureNonce()
	timestamp := time.Now().Unix()
	signature := GenerateSignature(appId, signatureNonce, serverSecret, timestamp)

	return map[string]string{
		"Action":          action,
		"AppId":           appId,
		"Signature":       signature,
		"SignatureNonce":  signatureNonce,
		"SignatureVersion": "2.0",
		"Timestamp":       strconv.FormatInt(timestamp, 10),
	}
}

// GenerateQueryParamsString 生成查询参数字符串
func GenerateQueryParamsString(action string, appId string, serverSecret string) string {
	params := GenerateQueryParams(action, appId, serverSecret)
	
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}
	
	return values.Encode()
}

