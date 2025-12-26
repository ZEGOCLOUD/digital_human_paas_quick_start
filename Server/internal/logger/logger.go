package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// LogLevel 日志等级
type LogLevel string

const (
	LogLevelDEBUG LogLevel = "DEBUG"
	LogLevelINFO  LogLevel = "INFO"
	LogLevelWARN  LogLevel = "WARN"
	LogLevelERROR LogLevel = "ERROR"
)

var (
	beijingLocation *time.Location
)

func init() {
	// 初始化北京时区 (UTC+8)
	var err error
	beijingLocation, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		// 如果加载失败，使用固定偏移
		beijingLocation = time.FixedZone("CST", 8*60*60)
	}
}

// logWithLevel 带日期打印的日志方法
func logWithLevel(level LogLevel, args ...interface{}) {
	// 获取当前时间并转换为北京时间
	now := time.Now().In(beijingLocation)
	timestamp := now.Format(time.RFC3339)

	// 构建日志消息，使用 fmt.Sprintln 确保参数之间有空格
	// fmt.Sprintln 会在参数之间自动添加空格
	var message string
	if len(args) > 0 {
		// 使用 fmt.Sprintln 格式化，会自动在参数之间添加空格
		message = fmt.Sprintln(args...)
		// 移除末尾的换行符（因为 log.Println 会自动添加）
		if len(message) > 0 && message[len(message)-1] == '\n' {
			message = message[:len(message)-1]
		}
	}
	logMessage := fmt.Sprintf("[%s] [%s] %s", timestamp, level, message)

	// 根据日志等级选择对应的输出
	switch level {
	case LogLevelDEBUG:
		log.Println(logMessage)
	case LogLevelINFO:
		log.Println(logMessage)
	case LogLevelWARN:
		log.Println(logMessage)
	case LogLevelERROR:
		log.New(os.Stderr, "", 0).Println(logMessage)
	default:
		log.Println(logMessage)
	}
}

// Log 通用日志方法，自动判断第一个参数是否为日志等级
func Log(levelOrFirstArg interface{}, args ...interface{}) {
	// 判断第一个参数是否为日志等级
	if level, ok := levelOrFirstArg.(LogLevel); ok {
		logWithLevel(level, args...)
	} else {
		// 如果不是日志等级，则将其加入到打印参数中，使用 INFO 级别
		allArgs := append([]interface{}{levelOrFirstArg}, args...)
		logWithLevel(LogLevelINFO, allArgs...)
	}
}

// LogInfo 信息日志
func LogInfo(args ...interface{}) {
	logWithLevel(LogLevelINFO, args...)
}

// LogDebug 调试日志
func LogDebug(args ...interface{}) {
	logWithLevel(LogLevelDEBUG, args...)
}

// LogWarn 警告日志
func LogWarn(args ...interface{}) {
	logWithLevel(LogLevelWARN, args...)
}

// LogError 错误日志
func LogError(args ...interface{}) {
	logWithLevel(LogLevelERROR, args...)
}

// LogInfof 格式化信息日志
func LogInfof(format string, args ...interface{}) {
	logWithLevel(LogLevelINFO, fmt.Sprintf(format, args...))
}

// LogDebugf 格式化调试日志
func LogDebugf(format string, args ...interface{}) {
	logWithLevel(LogLevelDEBUG, fmt.Sprintf(format, args...))
}

// LogWarnf 格式化警告日志
func LogWarnf(format string, args ...interface{}) {
	logWithLevel(LogLevelWARN, fmt.Sprintf(format, args...))
}

// LogErrorf 格式化错误日志
func LogErrorf(format string, args ...interface{}) {
	logWithLevel(LogLevelERROR, fmt.Sprintf(format, args...))
}

