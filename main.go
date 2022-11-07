package main

import (
	_ "go-error/controller"
	"go-error/web"
)

func main() {
	// 测试url:
	// curl "http://localhost/controller1/url1?dateBase64Str="
	// curl "http://localhost/controller1/url1?dateBase64Str=abc"
	// curl "http://localhost/controller1/url1?dateBase64Str=amlhbGUxMTIy"
	// curl "http://localhost/controller1/url1?dateBase64Str=SGVsbG8%3D"
	web.StartWeb()
}
