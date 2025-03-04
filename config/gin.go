package config

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GinInit(logger *zap.SugaredLogger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(GinRecovery(logger))
	r.Use(GinLogger(logger))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true, "data": "pong", "message": "pong"})
	})

	return r
}

func GinRecovery(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("Gin初始化发生错误: %v", err)
				debug.PrintStack()
				var message string
				switch t := err.(type) {
				case error:
					message = t.Error()
				default:
					message = err.(string)
				}
				c.JSON(http.StatusOK, gin.H{"success": false, "data": nil, "message": message})
			}
		}()
		c.Next()
	}
}

func GinLogger(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start).Milliseconds()
		logger.Infof("%d [%s] %s Elapse:%dms", c.Writer.Status(), c.Request.Method, c.Request.URL.Path, latency)
	}
}
