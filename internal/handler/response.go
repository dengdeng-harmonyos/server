package handler

import (
	"github.com/dengdeng-harmenyos/server/internal/models"
	"github.com/gin-gonic/gin"
)

// RespondSuccess 返回成功响应
func RespondSuccess(c *gin.Context, httpStatus int, data interface{}) {
	response := models.UnifiedApiResponse{
		Code: 0,
		Msg:  "success",
		Data: data,
	}
	c.JSON(httpStatus, response)
}

// RespondError 返回错误响应
func RespondError(c *gin.Context, httpStatus int, errorCode int, message string) {
	response := models.UnifiedApiResponse{
		Code: errorCode,
		Msg:  message,
		Data: nil,
	}
	c.JSON(httpStatus, response)
}
