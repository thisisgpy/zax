package config

import (
	"net/http"
	"regexp"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GinInit(logger *zap.SugaredLogger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// 替换输出的时间格式 yyyy-MM-dd HH:mm:ss
	r.Use(GinJson())
	// 404
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "data": nil, "message": "404 Not Found"})
	})
	// 405
	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"success": false, "data": nil, "message": "405 Method Not Allowed"})
	})
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
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "data": nil, "message": message})
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

func GinContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-Token")
		if token != "" {
			c.Set("username", "admin")
		} else {
			c.JSON(http.StatusOK, gin.H{"success": false, "data": nil, "message": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func GinJson() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer = &GinJsonWriter{
			ResponseWriter: c.Writer,
		}
		c.Next()
	}
}

type GinJsonWriter struct {
	gin.ResponseWriter
}

// 匹配时间格式  2025-03-04T10:00:00+08:00
var time_regexp_1 = regexp.MustCompile(`(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2}:\d{2})([+-]\d{2}:\d{2})`)

// 匹配时间格式  2025-03-04T10:00:00.000+08:00
var time_regexp_2 = regexp.MustCompile(`(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2}:\d{2}\.\d+)([+-]\d{2}:\d{2})`)

func (w GinJsonWriter) Write(data []byte) (int, error) {
	s := string(data)
	res := time_regexp_1.ReplaceAllString(s, "$1 $2")
	res = time_regexp_2.ReplaceAllString(res, "$1 $2")
	res = strings.ReplaceAll(res, "0001-01-01T00:00:00Z", "")
	return w.ResponseWriter.WriteString(res)
}
