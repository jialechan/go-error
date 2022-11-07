package bizerr

type bizErrConfig struct {
	code int
	msg  string
}

var (
	unknownError  = bizErrConfig{code: 1002, msg: "未知错误"}
	paramError    = bizErrConfig{code: 1003, msg: "参数错误"}
	userNameError = bizErrConfig{code: 1004, msg: "用户名错"}
)

// api错误的结构体
type APIException struct {
	ErrorCode int
	Msg       string
	Err       error
}

func newAPIException(bec bizErrConfig, err error) *APIException {
	return &APIException{
		ErrorCode: bec.code,
		Msg:       bec.msg,
		Err:       err,
	}
}

// 未知错误
func UnknownError(err error) *APIException {
	return newAPIException(unknownError, err)
}

// 参数错误
func ParameterError(err error) *APIException {
	return newAPIException(paramError, err)
}

// 用户错误
func UserNameError(err error) *APIException {
	return newAPIException(userNameError, err)
}
