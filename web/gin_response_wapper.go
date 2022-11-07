package web

import (
	"fmt"
	"net/http"

	"go-error/bizerr"

	"github.com/gin-gonic/gin"
)

type ginxHandleFunction func(*gin.Context) (interface{}, *bizerr.APIException)

func RespWrapper(handler ginxHandleFunction) func(c *gin.Context) {
	return func(c *gin.Context) {
		result, err := handler(c)
		if err != nil {
			// 打印错误
			fmt.Printf("err: %+v\n", err.Err)
			c.JSON(http.StatusOK, map[string]interface{}{"code": err.ErrorCode, "msg": err.Msg})
			return
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{"code": 0, "data": result})
		}
	}
}
