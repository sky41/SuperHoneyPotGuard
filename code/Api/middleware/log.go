package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"superhoneypotguard/database"
	"superhoneypotguard/models"

	"github.com/gin-gonic/gin"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r *responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

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

		database.DB.Create(&log)

		c.Next()
	}
}
