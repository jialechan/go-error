package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Engine 是gin的Engine指针实例变量
var Engine *gin.Engine

func init() {
	Engine = gin.New()
}

func StartWeb() {
	listen()
}

func listen() {
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", 80),
		Handler: Engine,
	}

	_ = s.ListenAndServe()
}
