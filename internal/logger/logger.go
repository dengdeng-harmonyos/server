package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	// InfoLogger 信息日志
	InfoLogger *log.Logger
	// ErrorLogger 错误日志
	ErrorLogger *log.Logger
	// DebugLogger 调试日志
	DebugLogger *log.Logger
	// AccessLogger 访问日志
	AccessLogger *log.Logger
)

// Init 初始化日志系统
func Init() {
	// 创建多个writer：标准输出
	infoWriter := io.MultiWriter(os.Stdout)
	errorWriter := io.MultiWriter(os.Stderr)

	// 初始化不同级别的logger
	InfoLogger = log.New(infoWriter, "[INFO] ", log.LstdFlags|log.Lshortfile)
	ErrorLogger = log.New(errorWriter, "[ERROR] ", log.LstdFlags|log.Lshortfile)
	DebugLogger = log.New(infoWriter, "[DEBUG] ", log.LstdFlags|log.Lshortfile)
	AccessLogger = log.New(infoWriter, "[ACCESS] ", log.LstdFlags)
}

// Info 记录信息日志
func Info(format string, v ...interface{}) {
	if InfoLogger == nil {
		Init()
	}
	InfoLogger.Output(2, fmt.Sprintf(format, v...))
}

// Error 记录错误日志
func Error(format string, v ...interface{}) {
	if ErrorLogger == nil {
		Init()
	}
	ErrorLogger.Output(2, fmt.Sprintf(format, v...))
}

// ErrorWithStack 记录错误日志并打印堆栈信息
func ErrorWithStack(err error, format string, v ...interface{}) {
	if ErrorLogger == nil {
		Init()
	}
	msg := fmt.Sprintf(format, v...)
	if err != nil {
		msg = fmt.Sprintf("%s: %v", msg, err)
	}
	ErrorLogger.Output(2, msg)
	// 打印堆栈信息
	if err != nil {
		ErrorLogger.Output(2, fmt.Sprintf("Stack trace:\n%s", debug.Stack()))
	}
}

// Debug 记录调试日志
func Debug(format string, v ...interface{}) {
	if DebugLogger == nil {
		Init()
	}
	DebugLogger.Output(2, fmt.Sprintf(format, v...))
}

// Access 记录访问日志
func Access(format string, v ...interface{}) {
	if AccessLogger == nil {
		Init()
	}
	AccessLogger.Output(2, fmt.Sprintf(format, v...))
}

// GinLogger Gin框架的日志中间件
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 请求日志
		Access("→ %s %s from %s", c.Request.Method, c.Request.URL.Path, c.ClientIP())
		if len(c.Request.URL.RawQuery) > 0 {
			Debug("  Query: %s", c.Request.URL.RawQuery)
		}

		// 处理请求
		c.Next()

		// 响应日志
		duration := time.Since(startTime)
		statusCode := c.Writer.Status()

		statusEmoji := "✓"
		logFunc := Access
		if statusCode >= 400 {
			statusEmoji = "✗"
			logFunc = Error
		}

		logFunc("← %s %s %s %d %s [%v]",
			statusEmoji,
			c.Request.Method,
			c.Request.URL.Path,
			statusCode,
			c.ClientIP(),
			duration,
		)

		// 如果有错误，记录详细信息
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				Error("Request error: %v", err.Err)
				Error("Error type: %v, Meta: %v", err.Type, err.Meta)
			}
		}
	}
}

// GinRecovery Gin框架的恢复中间件
func GinRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				Error("========================================")
				Error("PANIC recovered: %v", err)
				Error("Request: %s %s", c.Request.Method, c.Request.URL.Path)
				Error("Client IP: %s", c.ClientIP())
				Error("Stack trace:\n%s", debug.Stack())
				Error("========================================")
				c.JSON(500, gin.H{
					"success": false,
					"error":   "Internal server error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
