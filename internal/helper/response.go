package helper

import "github.com/gin-gonic/gin"

// APIResponse adalah format standar response JSON
func APIResponse(ctx *gin.Context, message string, statusCode int, data interface{}) {
	jsonResponse := gin.H{
		"status_code": statusCode,
		"message":     message,
		"data":        data,
	}
	ctx.JSON(statusCode, jsonResponse)
}

// ErrorResponse adalah format standar untuk response error
func ErrorResponse(ctx *gin.Context, message string, statusCode int, err error) {
	jsonResponse := gin.H{
		"status_code": statusCode,
		"message":     message,
		"error":       err.Error(),
	}
	ctx.JSON(statusCode, jsonResponse)
}
