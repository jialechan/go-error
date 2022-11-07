package service

import (
	"go-error/bizerr"
	"go-error/util"

	"github.com/pkg/errors"
)

func Fun1(dateBase64Str string) (interface{}, *bizerr.APIException) {

	if len(dateBase64Str) == 0 {
		// 应用程序中出现错误时，使用 errors.New 或者 errors.Errorf  返回错误
		return nil, bizerr.ParameterError(errors.New("参数不能为空"))
	}

	dataStr, err := util.DecodeBase64(dateBase64Str)
	if err != nil {
		// 如果是调用应用程序的其他函数出现错误，请直接返回，如果需要携带信息，请使用 errors.WithMessage
		return nil, bizerr.ParameterError(errors.WithMessage(err, "不是一个合法的base64字符串: "+dateBase64Str))
	}

	if dataStr != "Hello" {
		// 应用程序中出现错误时，使用 errors.New 或者 errors.Errorf  返回错误
		return nil, bizerr.UserNameError(errors.New("期望值不对"))
	}

	return nil, nil
}
