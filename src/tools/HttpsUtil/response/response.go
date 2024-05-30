package response

import (
	"github.com/gin-gonic/gin"
	"selfWeb/src/configuration/structs"
)

func ReturnContextNoBody(context *gin.Context, code int, message string) *gin.Context {
	response := structs.Response{
		Code:    code,
		Message: message,
		Data:    nil,
	}
	context.JSON(code, response)
	return context
}

func ReturnContext(context *gin.Context, code int, message string, body any) *gin.Context {
	response := structs.Response{
		Code:    code,
		Message: message,
		Data:    body,
	}
	context.JSON(code, response)
	return context
}
