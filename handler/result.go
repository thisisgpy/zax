package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ZaxResult struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, ZaxResult{Success: true, Data: data, Message: ""})
}

func Error(c *gin.Context, message string) {
	c.JSON(http.StatusOK, ZaxResult{Success: false, Data: nil, Message: message})
}
