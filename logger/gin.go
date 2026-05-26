package logger

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// GinLogger is a custom logger that can be used with gin Framework. This function will
// generate a custom log format as follows:
//
// [Date] [GIN] [INFO] [RequestID] [StatusCode] [Method] [Path] [ClientIP] [Latency] [UserAgent] [ErrorMessage]
func GinLogger(param gin.LogFormatterParams) string {
	var now = time.Now().Format("2006/01/02 15:04:05")
	requestID := GetRequestID(param.Request.Context())
	return fmt.Sprintf("[%s] [GIN] [INFO] [%d] [%s] [%s] [%s] [%dms] [%s] %s [%s]\n",
		now,
		param.StatusCode,
		param.Method,
		param.Path,
		param.ClientIP,
		param.Latency.Milliseconds(),
		param.Request.UserAgent(),
		param.ErrorMessage,
		requestID,
	)
}

// RecoveryLogger is a middleware that will recover from any panic that occurs during
// the execution of the request and return a 500 status code with a JSON response.
// The error message will be logged with the "RECOVER" log level.
func RecoveryLogger(withTrace bool, response map[string]any) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				l := New("RECOVER")
				if withTrace {
					l.Error(c.Request.Context(), err)
				} else {
					l.ErrorWithoutTrace(c.Request.Context(), err)
				}

				jsonObj := gin.H{
					"error": "An unexpected error occurred. Please try again later / contact admin",
				}
				if response != nil {
					jsonObj = response
				}
				c.AbortWithStatusJSON(500, jsonObj)
			}
		}()
		c.Next()
	}
}

// GinDebugRoute logs information about a Gin route during the debugging process.
// It prints the HTTP method, absolute path, handler name, and the number of handlers
// associated with the route.
func GinDebugRoute(httpMethod, absolutePath, handlerName string, nuHandlers int) {
	var now = time.Now().Format("2006/01/02 15:04:05")
	fmt.Printf("[%s] [GIN] [INFO] %v %v %v %v \n", now, httpMethod, absolutePath, handlerName, nuHandlers)
}

// GinDebugPrint logs debug information with a custom format.
func GinDebugPrint(format string, values ...interface{}) {
	var now = time.Now().Format("2006/01/02 15:04:05")
	fmt.Printf("[%s] [GIN] [INFO] %v \n", now, values)
}

// NewContextFromGin creates a standard context.Context from a gin.Context.
// Useful when passing context to service/repository layers that don't depend on gin.
func NewContextFromGin(c *gin.Context) context.Context {
	return c.Request.Context()
}
