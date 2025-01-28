package logger

import (
	"github.com/gin-gonic/gin"
)

func Examplemain() {
	r := gin.New()
	r.Use(RecoveryLogger())
	r.Use(gin.LoggerWithFormatter(GinLogger))

	gin.LoggerWithFormatter(GinLogger)

	gin.DebugPrintRouteFunc = GinDebugRoute
	gin.DebugPrintFunc = GinDebugPrint

	r.GET("/ping", func(c *gin.Context) {
		log := New("PING")

		log.Log("Test log ping")

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
