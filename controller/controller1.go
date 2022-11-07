package controller

import (
	"go-error/bizerr"
	"go-error/service"
	"go-error/web"

	"github.com/gin-gonic/gin"
)

func init() {
	web.Engine.GET("/controller1/url1", web.RespWrapper(controller1Url1Handler))
}

func controller1Url1Handler(c *gin.Context) (interface{}, *bizerr.APIException) {
	dateBase64Str := c.Query("dateBase64Str")
	return service.Fun1(dateBase64Str)
}
