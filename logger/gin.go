package logger

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// GinLogger is a custom logger that can be used with gin Framework. This function will
// generate a custom log format as follows:
//
// [Date] [GIN] [INFO] [StatusCode] [Method] [Path] [ClientIP] [Latency] [UserAgent] [ErrorMessage]
//
// The log format is customizable, but this function will always return a string
// that ends with a newline character. The function will also always return a
// string, regardless of whether the error is nil or not.
func GinLogger(param gin.LogFormatterParams) string {
	var now = time.Now().Format("2006/01/02 15:04:05")
	return fmt.Sprintf("[%s] [GIN] [INFO] [%d] [%s] [%s] [%s] [%dms] [%s] %s \n",
		now,
		param.StatusCode,
		param.Method,
		param.Path,
		param.ClientIP,
		param.Latency.Milliseconds(),
		param.Request.UserAgent(),
		param.ErrorMessage,
	)
}

// RecoveryLogger is a middleware that will recover from any panic that occurs during
// the execution of the request and return a 500 status code with a JSON response.
// The error message will be logged with the "RECOVER" log level.
func RecoveryLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log := New("RECOVER")
				log.ErrorWithoutTrace(err)

				c.AbortWithStatusJSON(500, gin.H{
					"error": "An unexpected error occurred. Please try again later / contact admin",
				})
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
// The log message includes the current timestamp, formatted as specified,
// followed by the provided values. It is intended for use with the Gin
// framework to output debug information during the request processing.
//
// Parameters:
//   - format: The format string for the log message.
//   - values: A variadic parameter representing the values to be logged.
func GinDebugPrint(format string, values ...interface{}) {
	var now = time.Now().Format("2006/01/02 15:04:05")
	fmt.Printf("[%s] [GIN] [INFO] %v \n", now, values)
}
