package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"sync"
	"time"

	"superhoneypotguard/database"
	"superhoneypotguard/models"

	"github.com/gin-gonic/gin"
)

// 提交目的：实现异步日志写入，提升API响应性能
// 提交内容：使用channel和goroutine实现异步日志记录，避免阻塞主流程
// 提交时间：2026-01-19

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r *responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// 日志缓冲区
var logBuffer = make(chan models.OperationLog, 1000)
var logBatchSize = 100
var logFlushInterval = 5 * time.Second
var wg sync.WaitGroup

// 启动日志处理协程
func init() {
	wg.Add(1)
	go processLogs()
}

// 异步处理日志
func processLogs() {
	defer wg.Done()

	batch := make([]models.OperationLog, 0, logBatchSize)
	ticker := time.NewTicker(logFlushInterval)
	defer ticker.Stop()

	for {
		select {
		case log := <-logBuffer:
			batch = append(batch, log)
			if len(batch) >= logBatchSize {
				flushLogs(batch)
				batch = batch[:0]
			}

		case <-ticker.C:
			if len(batch) > 0 {
				flushLogs(batch)
				batch = batch[:0]
			}
		}
	}
}

// 批量写入日志
func flushLogs(logs []models.OperationLog) {
	if len(logs) == 0 {
		return
	}

	if err := database.DB.CreateInBatches(logs, logBatchSize).Error; err != nil {
		// 如果批量写入失败，尝试逐条写入
		for _, log := range logs {
			database.DB.Create(&log)
		}
	}
}

// LogMiddleware 日志中间件（异步版本）
func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		w := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = w

		c.Next()

		duration := time.Since(startTime).Milliseconds()
		responseBody := w.body.String()

		var responseData map[string]interface{}
		status := 1
		if err := json.Unmarshal([]byte(responseBody), &responseData); err == nil {
			if success, ok := responseData["success"].(bool); ok && !success {
				status = 0
			}
		}

		userId, _ := c.Get("userId")
		username, _ := c.Get("username")

		var userIdPtr *int
		if uid, ok := userId.(int); ok {
			userIdPtr = &uid
		}

		var usernamePtr *string
		if uname, ok := username.(string); ok {
			usernamePtr = &uname
		}

		method := c.Request.Method
		path := c.FullPath()
		ip := c.ClientIP()

		var paramsStr string
		if len(requestBody) > 0 {
			if len(requestBody) > 500 {
				paramsStr = string(requestBody[:500])
			} else {
				paramsStr = string(requestBody)
			}
		}

		var resultStr string
		if len(responseBody) > 0 {
			if len(responseBody) > 500 {
				resultStr = responseBody[:500]
			} else {
				resultStr = responseBody
			}
		}

		log := models.OperationLog{
			UserID:      userIdPtr,
			Username:    usernamePtr,
			Operation:   path,
			Method:      &method,
			URL:         &path,
			IP:          &ip,
			Params:      &paramsStr,
			Result:      &resultStr,
			Status:      status,
			ExecuteTime: int(duration),
		}

		// 异步写入日志
		select {
		case logBuffer <- log:
		default:
			// 如果缓冲区满，同步写入
			database.DB.Create(&log)
		}

		c.Next()
	}
}

// Shutdown 优雅关闭日志处理
func Shutdown() {
	close(logBuffer)
	wg.Wait()
}