package logger

import (
	"github.com/gin-gonic/gin"
)

func Examplemain() {
	r := gin.New()

	// Add RequestID middleware first — it sets the request ID in context
	r.Use(RequestIDMiddleware())
	r.Use(RecoveryLogger(false, nil))
	r.Use(gin.LoggerWithFormatter(GinLogger))

	gin.DebugPrintRouteFunc = GinDebugRoute
	gin.DebugPrintFunc = GinDebugPrint

	r.GET("/ping", func(c *gin.Context) {
		l := New("PING")
		ctx := c.Request.Context()

		l.Log(ctx, "Test log ping")
		l.Error(ctx, "something went wrong")

		// Test recovery
		var a any
		a = 1
		a = a.(string)

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(":8082")
}
