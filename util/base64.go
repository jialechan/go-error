package util

import (
	"encoding/base64"

	"github.com/pkg/errors"
)

func DecodeBase64(base64Str string) (resultStr string, err error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		// 如果是调用其他库（标准库、企业公共库、开源第三方库等）获取到错误时，请使用 errors.Wrap 添加堆栈信息
		return "", errors.Wrap(err, "不是一个合法的base64编码字符串")
	}
	return string(decodeBytes), nil
}
