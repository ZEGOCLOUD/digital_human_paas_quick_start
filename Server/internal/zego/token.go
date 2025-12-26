package zego

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"time"
)

// ErrorCode 错误码
type ErrorCode int

const (
	ErrorCodeSuccess                  ErrorCode = 0
	ErrorCodeAppIDInvalid             ErrorCode = 1
	ErrorCodeUserIDInvalid            ErrorCode = 3
	ErrorCodeSecretInvalid            ErrorCode = 5
	ErrorCodeEffectiveTimeInSecondsInvalid ErrorCode = 6
)

// TokenError Token 生成错误
type TokenError struct {
	ErrorCode    ErrorCode
	ErrorMessage string
}

func (e *TokenError) Error() string {
	return e.ErrorMessage
}

// TokenInfo Token 信息
type TokenInfo struct {
	AppID   int64  `json:"app_id"`
	UserID  string `json:"user_id"`
	Nonce   int32  `json:"nonce"`
	Ctime   int64  `json:"ctime"`
	Expire  int64  `json:"expire"`
	Payload string `json:"payload"`
}

// makeNonce 生成 int32 范围的随机数
func makeNonce() int32 {
	// 生成 -2147483648 到 2147483647 之间的随机数
	max := big.NewInt(math.MaxInt32)
	min := big.NewInt(math.MinInt32)
	rangeVal := new(big.Int).Sub(max, min)
	
	n, err := rand.Int(rand.Reader, rangeVal)
	if err != nil {
		// 如果生成失败，返回一个默认值
		return 0
	}
	
	result := new(big.Int).Add(n, min)
	return int32(result.Int64())
}

// makeRandomIv 生成随机 IV (16 字节字符串)
func makeRandomIv() string {
	const chars = "0123456789abcdefghijklmnopqrstuvwxyz"
	result := make([]byte, 16)
	
	for i := 0; i < 16; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			// 如果生成失败，使用时间戳作为后备
			result[i] = chars[time.Now().UnixNano()%int64(len(chars))]
		} else {
			result[i] = chars[n.Int64()]
		}
	}
	
	return string(result)
}

// getAlgorithm 根据密钥长度确定算法
func getAlgorithm(keyLength int) string {
	switch keyLength {
	case 16:
		return "aes-128-cbc"
	case 24:
		return "aes-192-cbc"
	case 32:
		return "aes-256-cbc"
	default:
		return ""
	}
}

// aesEncrypt AES加密，使用模式: CBC/PKCS5Padding
func aesEncrypt(plainText string, key string, iv string) ([]byte, error) {
	keyBytes := []byte(key)
	ivBytes := []byte(iv)
	
	// 验证密钥长度
	keyLen := len(keyBytes)
	if keyLen != 16 && keyLen != 24 && keyLen != 32 {
		return nil, fmt.Errorf("invalid key length: %d", keyLen)
	}
	
	// 验证 IV 长度
	if len(ivBytes) != 16 {
		return nil, fmt.Errorf("invalid IV length: %d", len(ivBytes))
	}
	
	// 创建 AES cipher
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}
	
	// PKCS5Padding
	plainBytes := []byte(plainText)
	padLen := aes.BlockSize - len(plainBytes)%aes.BlockSize
	padText := bytes.Repeat([]byte{byte(padLen)}, padLen)
	plainBytes = append(plainBytes, padText...)
	
	// 创建 CBC mode
	mode := cipher.NewCBCEncrypter(block, ivBytes)
	
	// 加密
	ciphertext := make([]byte, len(plainBytes))
	mode.CryptBlocks(ciphertext, plainBytes)
	
	return ciphertext, nil
}

// GenerateToken04 生成 ZEGO Token04
func GenerateToken04(
	appId int64,
	userId string,
	secret string,
	effectiveTimeInSeconds int64,
	payload string,
) (string, error) {
	// 参数验证
	if appId == 0 {
		return "", &TokenError{
			ErrorCode:    ErrorCodeAppIDInvalid,
			ErrorMessage: "appID invalid",
		}
	}
	
	if userId == "" {
		return "", &TokenError{
			ErrorCode:    ErrorCodeUserIDInvalid,
			ErrorMessage: "userId invalid",
		}
	}
	
	if secret == "" || len(secret) != 32 {
		return "", &TokenError{
			ErrorCode:    ErrorCodeSecretInvalid,
			ErrorMessage: "secret must be a 32 byte string",
		}
	}
	
	if effectiveTimeInSeconds <= 0 {
		return "", &TokenError{
			ErrorCode:    ErrorCodeEffectiveTimeInSecondsInvalid,
			ErrorMessage: "effectiveTimeInSeconds invalid",
		}
	}
	
	// 创建时间戳
	createTime := time.Now().Unix()
	
	// 构建 Token 信息
	tokenInfo := TokenInfo{
		AppID:   appId,
		UserID:  userId,
		Nonce:   makeNonce(),
		Ctime:   createTime,
		Expire:  createTime + effectiveTimeInSeconds,
		Payload: payload,
	}
	
	// 将 Token 信息转为 JSON
	plainTextBytes, err := json.Marshal(tokenInfo)
	if err != nil {
		return "", fmt.Errorf("marshal token info failed: %w", err)
	}
	plainText := string(plainTextBytes)
	
	// 随机生成 16 字节 IV
	iv := makeRandomIv()
	
	// 进行 AES 加密
	encryptBuf, err := aesEncrypt(plainText, secret, iv)
	if err != nil {
		return "", fmt.Errorf("aes encrypt failed: %w", err)
	}
	
	// Token 二进制拼接: 过期时间(8字节) + IV长度(2字节) + IV(16字节) + 加密信息长度(2字节) + 加密信息
	buf := new(bytes.Buffer)
	
	// 过期时间 (8 字节, big-endian)
	binary.Write(buf, binary.BigEndian, tokenInfo.Expire)
	
	// IV 长度 (2 字节, big-endian)
	binary.Write(buf, binary.BigEndian, uint16(len(iv)))
	
	// IV (16 字节)
	buf.WriteString(iv)
	
	// 加密信息长度 (2 字节, big-endian)
	binary.Write(buf, binary.BigEndian, uint16(len(encryptBuf)))
	
	// 加密信息
	buf.Write(encryptBuf)
	
	// Base64 编码
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	
	// 返回 "04" + Base64 编码
	return "04" + encoded, nil
}

